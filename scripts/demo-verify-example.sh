#!/bin/bash
# HIS-Go 演示验证工具使用示例
# 用法: ./scripts/demo-verify-example.sh

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

# 示例 1: 验证管理端演示
example_verify_admin() {
    print_info "示例 1: 验证管理端演示"
    echo ""
    
    print_info "步骤 1: 启动管理端演示服务"
    echo "make demo-admin"
    echo ""
    
    print_info "步骤 2: 等待服务启动完成"
    echo "sleep 30"
    echo ""
    
    print_info "步骤 3: 验证管理端演示环境"
    echo "make verify-admin"
    echo ""
    
    print_info "步骤 4: 查看验证结果"
    echo "# 如果所有检查都通过，会显示:"
    echo "# [SUCCESS] 管理端演示环境验证通过！"
    echo ""
}

# 示例 2: 验证患者端演示
example_verify_patient() {
    print_info "示例 2: 验证患者端演示"
    echo ""
    
    print_info "步骤 1: 启动患者端演示服务"
    echo "make demo-patient"
    echo ""
    
    print_info "步骤 2: 等待服务启动完成"
    echo "sleep 30"
    echo ""
    
    print_info "步骤 3: 验证患者端演示环境"
    echo "make verify-patient"
    echo ""
    
    print_info "步骤 4: 查看验证结果"
    echo "# 如果所有检查都通过，会显示:"
    echo "# [SUCCESS] 患者端演示环境验证通过！"
    echo ""
}

# 示例 3: 测试 API 接口
example_test_api() {
    print_info "示例 3: 测试 API 接口"
    echo ""
    
    print_info "步骤 1: 确保演示服务已启动"
    echo "make demo-admin"
    echo ""
    
    print_info "步骤 2: 测试 API 接口"
    echo "make demo-verify"
    echo ""
    
    print_info "步骤 3: 查看测试结果"
    echo "# 如果所有测试都通过，会显示:"
    echo "# [SUCCESS] API 接口测试通过！"
    echo ""
}

# 示例 4: 验证所有演示环境
example_verify_all() {
    print_info "示例 4: 验证所有演示环境"
    echo ""
    
    print_info "步骤 1: 启动管理端演示服务"
    echo "make demo-admin"
    echo ""
    
    print_info "步骤 2: 等待服务启动完成"
    echo "sleep 30"
    echo ""
    
    print_info "步骤 3: 验证所有演示环境"
    echo "make demo-verify"
    echo ""
    
    print_info "步骤 4: 查看验证结果"
    echo "# 会同时验证管理端和患者端演示环境"
    echo ""
}

# 示例 5: 自定义验证
example_custom_verify() {
    print_info "示例 5: 自定义验证"
    echo ""
    
    print_info "步骤 1: 编辑验证脚本"
    echo "vim scripts/demo-verify.sh"
    echo ""
    
    print_info "步骤 2: 添加自定义验证函数"
    echo "# 在脚本中添加:"
    echo "# verify_custom_feature() {"
    echo "#     print_info \"验证自定义功能...\""
    echo "#     # 添加验证逻辑"
    echo "# }"
    echo ""
    
    print_info "步骤 3: 在主函数中添加自定义命令"
    echo "# 在 main() 函数中添加:"
    echo "# custom)"
    echo "#     verify_custom_feature"
    echo "#     ;;"
    echo ""
    
    print_info "步骤 4: 运行自定义验证"
    echo "./scripts/demo-verify.sh custom"
    echo ""
}

# 显示帮助
show_help() {
    echo "HIS-Go 演示验证工具使用示例"
    echo ""
    echo "用法: $0 [命令]"
    echo ""
    echo "命令:"
    echo "  admin      验证管理端演示示例"
    echo "  patient    验证患者端演示示例"
    echo "  api        测试 API 接口示例"
    echo "  all        验证所有演示环境示例"
    echo "  custom     自定义验证示例"
    echo "  help       显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 admin    # 显示验证管理端演示示例"
    echo "  $0 patient  # 显示验证患者端演示示例"
    echo "  $0 api      # 显示测试 API 接口示例"
    echo "  $0 all      # 显示验证所有演示环境示例"
    echo "  $0 custom   # 显示自定义验证示例"
}

# 主函数
main() {
    case "${1:-help}" in
        admin)
            example_verify_admin
            ;;
        patient)
            example_verify_patient
            ;;
        api)
            example_test_api
            ;;
        all)
            example_verify_all
            ;;
        custom)
            example_custom_verify
            ;;
        help|*)
            show_help
            ;;
    esac
}

main "$@"
