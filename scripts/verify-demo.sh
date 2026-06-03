#!/bin/bash
# HIS-Go 演示环境验证脚本
# 用法: ./scripts/verify-demo.sh [admin|patient]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

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

# 检查服务健康状态
check_health() {
    local service=$1
    local port=$2
    local url="http://localhost:${port}/health"
    
    if curl -s -f "$url" > /dev/null 2>&1; then
        print_success "$service 健康检查通过"
        return 0
    else
        print_error "$service 健康检查失败"
        return 1
    fi
}

# 验证管理端演示
verify_admin() {
    print_info "验证管理端演示环境..."
    
    local services=(
        "Gateway:8080"
        "Auth:8081"
        "User:8082"
        "Registration:8083"
        "Clinic:8084"
        "Prescription:8085"
        "Billing:8086"
        "Pharmacy:8087"
        "Schedule:8090"
        "System:8096"
    )
    
    local failed=0
    
    for service_port in "${services[@]}"; do
        IFS=':' read -r service port <<< "$service_port"
        if ! check_health "$service" "$port"; then
            ((failed++))
        fi
    done
    
    echo ""
    print_info "检查 Nginx 访问..."
    if curl -s -f "http://localhost/admin" > /dev/null 2>&1; then
        print_success "管理端前端可访问"
    else
        print_error "管理端前端不可访问"
        ((failed++))
    fi
    
    echo ""
    if [ $failed -eq 0 ]; then
        print_success "管理端演示环境验证通过！"
        print_info "访问地址: http://localhost/admin"
        print_info "演示账号: demo-admin / demo123"
    else
        print_error "管理端演示环境验证失败，共 $failed 个服务异常"
        return 1
    fi
}

# 验证患者端演示
verify_patient() {
    print_info "验证患者端演示环境..."
    
    local services=(
        "Gateway:8080"
        "Auth:8081"
        "User:8082"
        "Registration:8083"
        "Schedule:8090"
        "Prescription:8085"
        "Examination:8088"
        "Followup:8092"
        "Health Record:8093"
    )
    
    local failed=0
    
    for service_port in "${services[@]}"; do
        IFS=':' read -r service port <<< "$service_port"
        if ! check_health "$service" "$port"; then
            ((failed++))
        fi
    done
    
    echo ""
    print_info "检查 Nginx 访问..."
    if curl -s -f "http://localhost/patient" > /dev/null 2>&1; then
        print_success "患者端前端可访问"
    else
        print_error "患者端前端不可访问"
        ((failed++))
    fi
    
    echo ""
    if [ $failed -eq 0 ]; then
        print_success "患者端演示环境验证通过！"
        print_info "访问地址: http://localhost/patient"
        print_info "小程序配置: frontend/his-mp-webview/pages/index/index.js"
    else
        print_error "患者端演示环境验证失败，共 $failed 个服务异常"
        return 1
    fi
}

# 显示帮助
show_help() {
    echo "HIS-Go 演示环境验证脚本"
    echo ""
    echo "用法: $0 [命令]"
    echo ""
    echo "命令:"
    echo "  admin    验证管理端演示环境"
    echo "  patient  验证患者端演示环境"
    echo "  help     显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 admin    # 验证管理端演示"
    echo "  $0 patient  # 验证患者端演示"
}

# 主函数
main() {
    case "${1:-help}" in
        admin)
            verify_admin
            ;;
        patient)
            verify_patient
            ;;
        help|*)
            show_help
            ;;
    esac
}

main "$@"
