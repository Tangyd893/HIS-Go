# HIS-Go 代码质量检查脚本 (Windows PowerShell)
# 用法: .\scripts\check.ps1

$ErrorActionPreference = "Stop"
$BackendDir = Join-Path $PSScriptRoot "..\backend"

Write-Output "=== HIS-Go 代码质量检查 ==="
Write-Output ""

Push-Location $BackendDir

try {
    Write-Output "[1/4] 代码格式检查 (gofmt)..."
    $fmtResult = gofmt -l .
    if ($fmtResult) {
        Write-Output "以下文件格式不符合规范:"
        Write-Output $fmtResult
        throw "gofmt 检查未通过"
    }
    Write-Output "  ✓ 通过"

    Write-Output "[2/4] 静态代码检查 (go vet)..."
    go vet ./...
    Write-Output "  ✓ 通过"

    Write-Output "[3/4] 运行测试 (go test)..."
    go test -count=1 ./...
    Write-Output "  ✓ 通过"

    Write-Output "[4/4] 编译检查 (go build)..."
    go build ./cmd/...
    Write-Output "  ✓ 通过"

    Write-Output ""
    Write-Output "=== 全部检查通过 ==="
} catch {
    Write-Output ""
    Write-Output "=== 检查失败: $_ ==="
    exit 1
} finally {
    Pop-Location
}
