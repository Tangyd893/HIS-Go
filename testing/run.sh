#!/bin/bash
# HIS-Go API 集成测试运行脚本
# 用法: 确保 Docker 服务已启动后执行 bash testing/run.sh
# 环境变量: HIS_INTEGRATION_TEST=true HIS_BASE_URL=http://localhost:8080

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TESTING_DIR="${SCRIPT_DIR}"

echo "=== HIS-Go API 集成测试 ==="
echo ""

export HIS_INTEGRATION_TEST=true
export HIS_BASE_URL="${HIS_BASE_URL:-http://localhost:8080}"

echo "目标地址: ${HIS_BASE_URL}"
echo ""

cd "${TESTING_DIR}"

go test -v -count=1 ./...

echo ""
echo "=== 集成测试完成 ==="
