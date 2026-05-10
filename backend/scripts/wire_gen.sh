#!/bin/bash
# Wire 依赖注入代码生成脚本
# 为所有 17 个微服务模块生成 wire_gen.go

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$SCRIPT_DIR")"

cd "$BACKEND_DIR"

echo "=== 安装 Wire 工具 ==="
go install github.com/google/wire/cmd/wire@latest

echo ""
echo "=== 生成依赖注入代码 ==="

MODULES=(
    "internal/auth"
    "internal/user"
    "internal/registration"
    "internal/clinic"
    "internal/emr"
    "internal/prescription"
    "internal/billing"
    "internal/pharmacy"
    "internal/examination"
    "internal/inpatient"
    "internal/schedule"
    "internal/outpatient"
    "internal/followup"
    "internal/health_record"
    "internal/notification"
    "internal/statistics"
    "internal/system"
)

FAILED=0
for module in "${MODULES[@]}"; do
    echo -n "  $module ... "
    if wire ./"$module" 2>&1; then
        echo "OK"
    else
        echo "FAILED"
        FAILED=$((FAILED + 1))
    fi
done

echo ""
if [ $FAILED -eq 0 ]; then
    echo "=== 全部 17 个模块 Wire 代码生成成功 ==="
else
    echo "=== $FAILED 个模块生成失败，请检查 wire.go 配置 ==="
    exit 1
fi
