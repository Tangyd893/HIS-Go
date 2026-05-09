#!/bin/bash
# 数据库初始化脚本
# 用法: bash scripts/db_init.sh
# 依赖: psql, Docker 运行中的 his-postgres 容器
# 按顺序执行：migrations/*.sql（建库+建表） → sql/seed_data.sql（种子数据）

set -e

echo "=== 开始初始化数据库 ==="

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-his_admin}"
DB_PASSWORD="${DB_PASSWORD:-change_me_123}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
MIGRATIONS_DIR="${SCRIPT_DIR}/../migrations"
SQL_DIR="${SCRIPT_DIR}/../sql"

# 检查 psql
command -v psql >/dev/null 2>&1 || { echo "错误: psql 未安装"; exit 1; }

export PGPASSWORD="${DB_PASSWORD}"

# 第一步：按编号顺序执行迁移脚本（建库 + 建表）
if [ -d "${MIGRATIONS_DIR}" ]; then
    echo "执行版本化迁移脚本 (${MIGRATIONS_DIR}) ..."
    for migration in "${MIGRATIONS_DIR}"/*.sql; do
        if [ -f "${migration}" ]; then
            migration_name=$(basename "${migration}")
            echo "  → 执行 ${migration_name} ..."
            psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d postgres -f "${migration}"
        fi
    done
    echo "迁移脚本执行完成"
else
    echo "警告: migrations 目录不存在，跳过迁移"
fi

# 第二步：同时执行 init_all.sql（兼容旧流程，含幂等检查）
if [ -f "${SQL_DIR}/init_all.sql" ]; then
    echo "执行 init_all.sql（兼容旧流程） ..."
    psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d postgres -f "${SQL_DIR}/init_all.sql"
    echo "建表完成"
else
    echo "信息: ${SQL_DIR}/init_all.sql 不存在，跳过"
fi

# 第三步：导入种子数据
if [ -f "${SQL_DIR}/seed_data.sql" ]; then
    echo "执行 seed_data.sql ..."
    psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d postgres -f "${SQL_DIR}/seed_data.sql"
    echo "种子数据导入完成"
else
    echo "信息: ${SQL_DIR}/seed_data.sql 不存在，跳过种子数据"
fi

unset PGPASSWORD
echo "=== 数据库初始化完成 ==="
