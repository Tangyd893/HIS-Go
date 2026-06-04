# HIS-Go 患者端演示启动脚本 (Windows PowerShell)
# 用法: .\scripts\demo-patient.ps1 [start|stop|restart|status|logs|build]

param(
    [ValidateSet("start","stop","restart","status","logs","build","help")]
    [string]$Command = "help",
    [string]$Service = ""
)

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir
$ComposeFile = "$ProjectRoot\deploy\compose\demo-patient.yml"
$EnvFile = "$ProjectRoot\deploy\config\demo.env"
$PatientDist = "$ProjectRoot\frontend\patient\dist"

function Write-Info  { Write-Host "[INFO] $args" -ForegroundColor Blue }
function Write-OK    { Write-Host "[OK] $args" -ForegroundColor Green }
function Write-Warn  { Write-Host "[WARN] $args" -ForegroundColor Yellow }
function Write-Err   { Write-Host "[ERROR] $args" -ForegroundColor Red }

function Check-Docker {
    try { docker info | Out-Null }
    catch { Write-Err "Docker 未运行"; exit 1 }
}

function Check-ComposeFile {
    if (-not (Test-Path $ComposeFile)) {
        Write-Err "找不到 compose 文件: $ComposeFile"; exit 1
    }
}

function Build-Frontend {
    Write-Info "构建 Vue 患者端..."
    Push-Location "$ProjectRoot\frontend\patient"
    npm ci
    npm run build
    Pop-Location
    Write-OK "构建完成: $PatientDist"
}

function Start-Services {
    if (-not (Test-Path "$PatientDist\index.html")) {
        Write-Warn "前端构建产物不存在，请先运行 build"
        return
    }
    Write-Info "启动患者端演示服务..."
    $cmd = "docker compose -f `"$ComposeFile`""
    if (Test-Path $EnvFile) { $cmd += " --env-file `"$EnvFile`"" }
    $cmd += " up -d"
    Invoke-Expression $cmd
    Write-OK "患者端服务已启动（无 RabbitMQ）"
    Write-Info "访问: http://localhost/patient  账号: demo-patient / demo123"
}

function Stop-Services {
    Write-Info "停止患者端服务..."
    $cmd = "docker compose -f `"$ComposeFile`""
    if (Test-Path $EnvFile) { $cmd += " --env-file `"$EnvFile`"" }
    $cmd += " down"
    Invoke-Expression $cmd
    Write-OK "已停止"
}

function Show-Status {
    $cmd = "docker compose -f `"$ComposeFile`""
    if (Test-Path $EnvFile) { $cmd += " --env-file `"$EnvFile`"" }
    $cmd += " ps"
    Invoke-Expression $cmd
}

function Show-Logs {
    $cmd = "docker compose -f `"$ComposeFile`""
    if (Test-Path $EnvFile) { $cmd += " --env-file `"$EnvFile`"" }
    if ($Service) { $cmd += " logs -f $Service" }
    else { $cmd += " logs -f" }
    Invoke-Expression $cmd
}

# === Main ===
Check-Docker
Check-ComposeFile

switch ($Command) {
    "build"   { Build-Frontend }
    "start"   { Start-Services }
    "stop"    { Stop-Services }
    "restart" { Stop-Services; Start-Services }
    "status"  { Show-Status }
    "logs"    { Show-Logs }
    default   {
        Write-Host "HIS-Go 患者端演示脚本 (PowerShell)"
        Write-Host "用法: .\demo-patient.ps1 [build|start|stop|restart|status|logs]"
    }
}
