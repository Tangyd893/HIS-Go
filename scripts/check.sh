#!/bin/bash
# HIS-Go 代码质量检查脚本 (Linux/macOS)
# 用法: bash scripts/check.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="${SCRIPT_DIR}/../backend"

echo "=== HIS-Go 代码质量检查 ==="
echo ""

cd "${BACKEND_DIR}"

echo "[1/4] 代码格式检查 (gofmt)..."
fmt_result=$(gofmt -l .)
if [[[ -n "${fmt_result}" ]]; then
    echo "以下文件格式不符合规范:"
    echo "${fmt_result}"
    echo "ERROR: gofmt 检查未通过"
    exit 1
fi
echo "  ✓ 通过"

echo "[2/4] 静态代码检查 (go vet)..."
go vet ./...
echo "  ✓ 通过"

echo "[3/4] 运行测试 (go test)..."
go test -count=1 ./...
echo "  ✓ 通过"

echo "[4/4] 编译检查 (go build)..."
go build ./cmd/...
echo "  ✓ 通过"

echo ""
echo "=== 全部检查通过 ==="
