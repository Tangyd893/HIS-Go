#!/bin/bash
# HIS-Go 云服务器初始化脚本
# 适用系统: Ubuntu 22.04 / 24.04 LTS
# 用途: 从空 VPS → Docker + Compose + 防火墙就绪
#
# 用法:
#   curl -fsSL https://... | sudo bash
#   或
#   sudo bash scripts/cloud-init.sh
#
# 非 root 用户请使用 sudo 执行。

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_info()  { echo -e "${BLUE}[INFO]${NC} $1"; }
print_ok()    { echo -e "${GREEN}[OK]${NC} $1"; }
print_warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
print_err()   { echo -e "${RED}[ERROR]${NC} $1"; }

# ============================================================
# 0. 权限检查
# ============================================================
if [[ "$(id -u)" -ne 0 ]]; then
    print_err "请使用 root 或 sudo 执行此脚本"
    exit 1
fi

print_info "HIS-Go 云服务器初始化开始"
print_info "目标系统: $(lsb_release -ds 2>/dev/null || cat /etc/os-release | grep PRETTY | cut -d= -f2)"

# ============================================================
# 1. 系统更新 + 基础工具
# ============================================================
print_info "更新系统包..."
apt-get update -qq
apt-get upgrade -y -qq
apt-get install -y -qq \
    curl \
    wget \
    gnupg \
    ca-certificates \
    lsb-release \
    ufw \
    htop \
    git

print_ok "系统更新完成"

# ============================================================
# 2. 安装 Docker CE
# ============================================================
print_info "安装 Docker..."

# 卸载旧版本（如果存在）
for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do
    apt-get remove -y $pkg 2>/dev/null || true
done

# 添加 Docker 官方 GPG key
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
chmod a+r /etc/apt/keyrings/docker.gpg

# 添加 apt 仓库
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

apt-get update -qq
apt-get install -y -qq docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# 验证
docker --version
docker compose version

print_ok "Docker 安装完成"

# ============================================================
# 3. Docker 加速器（中国大陆云服务器）
# ============================================================
read -p "是否配置 Docker 镜像加速器？(中国大陆服务器建议配置) [y/N]: " setup_mirror
if [[ "$setup_mirror" = "y" ] || [[ "$setup_mirror" = "Y" ]]; then
    mkdir -p /etc/docker
    cat > /etc/docker/daemon.json <<'EOF'
{
  "registry-mirrors": [
    "https://docker.1ms.run",
    "https://docker.xuanyuan.me"
  ],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
EOF
    systemctl restart docker
    print_ok "Docker 镜像加速器已配置"
fi

# ============================================================
# 4. 配置防火墙 (UFW)
# ============================================================
print_info "配置防火墙..."

# 默认策略：拒绝入站，允许出站
ufw default deny incoming
ufw default allow outgoing

# 开放 SSH（防止把自己锁在门外）
ufw allow 22/tcp comment 'SSH'

# 开放 HTTP / HTTPS
ufw allow 80/tcp comment 'HTTP (Nginx)'
ufw allow 443/tcp comment 'HTTPS (Nginx)'

# 不开放数据库端口（仅内部使用）
# 5432 (PostgreSQL)、6379 (Redis)、5672 (RabbitMQ) 等均不允许外部访问

ufw --force enable

print_ok "防火墙已启用"
ufw status verbose

# ============================================================
# 5. 创建项目目录结构
# ============================================================
print_info "创建项目目录..."

mkdir -p /opt/his-go
mkdir -p /opt/his-go/nginx/html/admin
mkdir -p /opt/his-go/nginx/html/patient
mkdir -p /opt/his-go/data/postgresql
mkdir -p /opt/his-go/data/redis
mkdir -p /opt/his-go/data/rabbitmq

print_ok "项目目录已创建: /opt/his-go/"

# ============================================================
# 6. 配置 Docker 日志轮转
# ============================================================
if [[ ! -f /etc/docker/daemon.json ]]; then
    cat > /etc/docker/daemon.json <<'EOF'
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
EOF
    systemctl restart docker
    print_ok "Docker 日志轮转已配置 (10MB x 3)"
fi

# ============================================================
# 7. 验证
# ============================================================
print_info "验证安装..."

echo ""
echo "=== 安装摘要 ==="
echo "Docker:         $(docker --version)"
echo "Docker Compose: $(docker compose version)"
echo "OS:             $(lsb_release -ds 2>/dev/null || echo 'Unknown')"
echo "内核:           $(uname -r)"
echo "内存:           $(free -h | awk '/^Mem:/ {print $2}')"
echo "磁盘:           $(df -h / | awk 'NR==2 {print $2}')"
echo ""

print_ok "HIS-Go 云服务器初始化完成！"
echo ""
echo "下一步："
echo "  1. 上传项目文件到 /opt/his-go/"
echo "  2. 复制 deploy/config/demo.env.example 为 deploy/config/demo.env 并修改密码"
echo "  3. 管理端: docker compose -f deploy/compose/demo-admin.yml up -d"
echo "  4. 患者端: docker compose -f deploy/compose/demo-patient.yml up -d"
echo "  5. 参见: docs/演示部署-管理端.md  和  docs/演示部署-患者端.md"
