#!/bin/bash
# 数据库初始化脚本
# 用法: bash scripts/db_init.sh
# 依赖: psql, Docker 运行中的 his-postgres 容器

set -e

echo "=== 开始初始化数据库 ==="

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-his_admin}"
DB_PASSWORD="${DB_PASSWORD:-change_me_123}"
SQL_DIR="sql"

# 检查 psql
command -v psql >/dev/null 2>&1 || { echo "错误: psql 未安装"; exit 1; }

export PGPASSWORD="${DB_PASSWORD}"

# 执行建表脚本
if [ -f "${SQL_DIR}/init_all.sql" ]; then
    echo "执行 init_all.sql ..."
    psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d postgres -f "${SQL_DIR}/init_all.sql"
    echo "建表完成"
else
    echo "警告: ${SQL_DIR}/init_all.sql 不存在，跳过建表"
fi

# 执行种子数据脚本
if [ -f "${SQL_DIR}/seed_data.sql" ]; then
    echo "执行 seed_data.sql ..."
    psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d postgres -f "${SQL_DIR}/seed_data.sql"
    echo "种子数据导入完成"
else
    echo "信息: ${SQL_DIR}/seed_data.sql 不存在，跳过种子数据"
fi

unset PGPASSWORD
echo "=== 数据库初始化完成 ==="
