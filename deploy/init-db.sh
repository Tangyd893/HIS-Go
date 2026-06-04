#!/bin/sh
# 数据库初始化脚本 (Docker 容器内运行)

echo "=== 初始化所有数据库 ==="
psql -h postgresql -U "${DB_USER:-his_admin}" -f /init_all.sql
echo "[OK] 数据库创建完成"

echo ""
echo "=== 导入种子数据 ==="
psql -h postgresql -U "${DB_USER:-his_admin}" -f /seed_data.sql
echo "[OK] 种子数据导入完成"

echo ""
echo "=== 数据库初始化完成 ==="
