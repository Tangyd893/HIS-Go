# HIS-Go 代码质量检查脚本 (Windows PowerShell)
# 用法: .\scripts\check.ps1

$ErrorActionPreference = "Stop"
$BackendDir = Join-Path $PSScriptRoot "..\backend"

Write-Host "=== HIS-Go 代码质量检查 ===" -ForegroundColor Cyan
Write-Host ""

Push-Location $BackendDir

try {
    Write-Host "[1/4] 代码格式检查 (gofmt)..." -ForegroundColor Yellow
    $fmtResult = gofmt -l .
    if ($fmtResult) {
        Write-Host "以下文件格式不符合规范:" -ForegroundColor Red
        Write-Host $fmtResult
        throw "gofmt 检查未通过"
    }
    Write-Host "  ✓ 通过" -ForegroundColor Green

    Write-Host "[2/4] 静态代码检查 (go vet)..." -ForegroundColor Yellow
    go vet ./...
    Write-Host "  ✓ 通过" -ForegroundColor Green

    Write-Host "[3/4] 运行测试 (go test)..." -ForegroundColor Yellow
    go test -count=1 ./...
    Write-Host "  ✓ 通过" -ForegroundColor Green

    Write-Host "[4/4] 编译检查 (go build)..." -ForegroundColor Yellow
    go build ./cmd/...
    Write-Host "  ✓ 通过" -ForegroundColor Green

    Write-Host ""
    Write-Host "=== 全部检查通过 ===" -ForegroundColor Green
} catch {
    Write-Host ""
    Write-Host "=== 检查失败: $_ ===" -ForegroundColor Red
    exit 1
} finally {
    Pop-Location
}
