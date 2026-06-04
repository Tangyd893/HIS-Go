#!/bin/bash
# HIS-Go 患者端演示启动脚本
# 用法: ./scripts/demo-patient.sh [start|stop|restart|status|logs]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
COMPOSE_FILE="$PROJECT_ROOT/deploy/compose/demo-patient.yml"
ENV_FILE="$PROJECT_ROOT/deploy/config/demo.env"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查 Docker 是否运行
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker 未运行，请先启动 Docker"
        exit 1
    fi
}

# 检查 docker-compose 文件是否存在
check_compose_file() {
    if [ ! -f "$COMPOSE_FILE" ]; then
        print_error "找不到 docker-compose 文件: $COMPOSE_FILE"
        exit 1
    fi
}

# 检查前端构建产物
PATIENT_DIST="$PROJECT_ROOT/frontend/patient/dist"

check_frontend_build() {
    if [ ! -f "$PATIENT_DIST/index.html" ]; then
        print_warning "前端构建产物不存在: $PATIENT_DIST/index.html"
        print_info "请先运行: $0 build  或  make build-patient"
        return 1
    fi
    return 0
}

# 构建前端
build_frontend() {
    print_info "构建 Vue 患者端..."
    cd "$PROJECT_ROOT/frontend/patient"
    npm ci
    npm run build
    cd "$PROJECT_ROOT"
    print_success "患者端前端构建完成: $PATIENT_DIST"
}

# 检查环境变量文件
check_env_file() {
    if [ ! -f "$ENV_FILE" ]; then
        print_warning "找不到环境变量文件: $ENV_FILE"
        print_info "将使用默认环境变量"
    fi
}

# 启动服务
start_services() {
    print_info "启动患者端演示服务..."
    
    local compose_cmd="docker compose -f $COMPOSE_FILE"
    if [ -f "$ENV_FILE" ]; then
        compose_cmd="$compose_cmd --env-file $ENV_FILE"
    fi
    
    $compose_cmd up -d
    
    print_success "患者端演示服务已启动"
    print_info "服务列表："
    echo "  - PostgreSQL: localhost:5432"
    echo "  - Redis: localhost:6379"
    echo "  - Gateway: http://localhost:8080"
    echo "  - Auth: http://localhost:8081"
    echo "  - User: http://localhost:8082"
    echo "  - Registration: http://localhost:8083"
    echo "  - Schedule: http://localhost:8090"
    echo "  - Prescription: http://localhost:8085"
    echo "  - Examination: http://localhost:8088"
    echo "  - Followup: http://localhost:8092"
    echo "  - Health Record: http://localhost:8093"
    echo "  - Nginx: http://localhost:80"
    echo ""
    print_info "患者端访问地址: http://localhost/patient"
    print_info "演示账号: demo-patient / demo123"
}

# 停止服务
stop_services() {
    print_info "停止患者端演示服务..."
    
    local compose_cmd="docker compose -f $COMPOSE_FILE"
    if [ -f "$ENV_FILE" ]; then
        compose_cmd="$compose_cmd --env-file $ENV_FILE"
    fi
    
    $compose_cmd down
    
    print_success "患者端演示服务已停止"
}

# 重启服务
restart_services() {
    print_info "重启患者端演示服务..."
    stop_services
    start_services
}

# 查看服务状态
show_status() {
    print_info "患者端演示服务状态："
    
    local compose_cmd="docker compose -f $COMPOSE_FILE"
    if [ -f "$ENV_FILE" ]; then
        compose_cmd="$compose_cmd --env-file $ENV_FILE"
    fi
    
    $compose_cmd ps
}

# 查看日志
show_logs() {
    local service=$1
    local compose_cmd="docker compose -f $COMPOSE_FILE"
    if [ -f "$ENV_FILE" ]; then
        compose_cmd="$compose_cmd --env-file $ENV_FILE"
    fi
    
    if [ -n "$service" ]; then
        $compose_cmd logs -f "$service"
    else
        $compose_cmd logs -f
    fi
}

# 显示帮助
show_help() {
    echo "HIS-Go 患者端演示启动脚本"
    echo ""
    echo "用法: $0 [命令]"
    echo ""
    echo "命令:"
    echo "  build    构建 Vue 患者端前端"
    echo "  start    启动所有患者端演示服务 (自动检查前端是否已构建)"
    echo "  stop     停止所有患者端演示服务"
    echo "  restart  重启所有患者端演示服务"
    echo "  status   查看服务状态"
    echo "  logs     查看服务日志 (可指定服务名，如: $0 logs his-auth)"
    echo "  help     显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 start          # 启动所有服务"
    echo "  $0 logs his-auth  # 查看 auth 服务日志"
    echo "  $0 stop           # 停止所有服务"
}

# 主函数
main() {
    check_docker
    check_compose_file
    check_env_file
    
    case "${1:-help}" in
        build)
            build_frontend
            ;;
        start)
            if check_frontend_build; then
                start_services
            else
                print_info "是否立即构建？(y/N)"
                read -r answer
                if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
                    build_frontend
                    start_services
                fi
            fi
            ;;
        stop)
            stop_services
            ;;
        restart)
            if check_frontend_build; then
                restart_services
            else
                build_frontend
                restart_services
            fi
            ;;
        status)
            show_status
            ;;
        logs)
            show_logs "$2"
            ;;
        help|*)
            show_help
            ;;
    esac
}

main "$@"
