# HIS-Go 管理端演示启动脚本 (Windows PowerShell)
# 用法: .\scripts\demo-admin.ps1 [start|stop|restart|status|logs|build]

param(
    [ValidateSet("start","stop","restart","status","logs","build","help")]
    [string]$Command = "help",
    [string]$Service = ""
)

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir
$ComposeFile = "$ProjectRoot\deploy\compose\demo-admin.yml"
$EnvFile = "$ProjectRoot\deploy\config\demo.env"
$AdminDist = "$ProjectRoot\frontend\admin\dist"

function Write-Info  { Write-Information "[INFO] $args" -InformationAction Continue }
function Write-OK    { Write-Information "[OK] $args" -InformationAction Continue }
function Write-Warn  { Write-Information "[WARN] $args" -InformationAction Continue }
function Write-Err   { Write-Error "$args" }

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
    Write-Info "构建 Vue 管理端..."
    Push-Location "$ProjectRoot\frontend\admin"
    npm ci
    npm run build
    Pop-Location
    Write-OK "构建完成: $AdminDist"
}

function Start-Services {
    if (-not (Test-Path "$AdminDist\index.html")) {
        Write-Warn "前端构建产物不存在，请先运行 build"
        return
    }
    Write-Info "启动管理端演示服务..."
    $args = @("compose", "-f", $ComposeFile)
    if (Test-Path $EnvFile) { $args += @("--env-file", $EnvFile) }
    $args += @("up", "-d")
    & docker @args
    Write-OK "管理端服务已启动"
    Write-Info "访问: http://localhost/admin  账号: demo-admin / demo123"
}

function Stop-Services {
    Write-Info "停止管理端服务..."
    $args = @("compose", "-f", $ComposeFile)
    if (Test-Path $EnvFile) { $args += @("--env-file", $EnvFile) }
    $args += @("down")
    & docker @args
    Write-OK "已停止"
}

function Show-Status {
    $args = @("compose", "-f", $ComposeFile)
    if (Test-Path $EnvFile) { $args += @("--env-file", $EnvFile) }
    $args += @("ps")
    & docker @args
}

function Show-Logs {
    $args = @("compose", "-f", $ComposeFile)
    if (Test-Path $EnvFile) { $args += @("--env-file", $EnvFile) }
    $args += @("logs", "-f")
    if ($Service) { $args += @($Service) }
    & docker @args
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
        Write-Output "HIS-Go 管理端演示脚本 (PowerShell)"
        Write-Output "用法: .\demo-admin.ps1 [build|start|stop|restart|status|logs]"
    }
}
