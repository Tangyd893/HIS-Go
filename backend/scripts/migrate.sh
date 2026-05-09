#!/bin/bash
# 数据库迁移脚本 - 按编号顺序执行所有迁移脚本
# 用法: bash migrate.sh [up|down|status]
# 依赖: psql

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
MIGRATIONS_DIR="${SCRIPT_DIR}/../migrations"

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-his_admin}"
DB_PASSWORD="${DB_PASSWORD:-change_me_123}"

command -v psql >/dev/null 2>&1 || { echo "错误: psql 未安装"; exit 1; }

export PGPASSWORD="${DB_PASSWORD}"

ACTION="${1:-up}"

case "$ACTION" in
    up)
        echo "=== 执行全部迁移 (up) ==="
        if [ ! -d "$MIGRATIONS_DIR" ]; then
            echo "错误: 迁移目录不存在: $MIGRATIONS_DIR"
            exit 1
        fi

        for migration in "${MIGRATIONS_DIR}"/*.sql; do
            if [ -f "$migration" ]; then
                name=$(basename "$migration")
                echo "  → 执行 ${name} ..."
                psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d postgres -f "$migration" -v ON_ERROR_STOP=1
            fi
        done
        echo "=== 迁移执行完成 ==="
        ;;

    status)
        echo "=== 迁移文件列表 ==="
        if [ -d "$MIGRATIONS_DIR" ]; then
            ls -1 "${MIGRATIONS_DIR}"/*.sql 2>/dev/null | while read f; do
                echo "  $(basename "$f")"
            done
        else
            echo "  (迁移目录不存在)"
        fi
        ;;

    *)
        echo "用法: bash migrate.sh [up|status]"
        echo "  up      - 按顺序执行全部迁移脚本"
        echo "  status  - 列出迁移文件"
        exit 1
        ;;
esac

unset PGPASSWORD
