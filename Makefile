.PHONY: help fmt vet test build proto clean check all deps demo-admin demo-patient demo-admin-logs demo-patient-logs verify-admin verify-patient demo-verify

help: ## 显示帮助信息
	@echo "HIS-Go 项目 Makefile"
	@echo ""
	@echo "可用目标:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

deps: ## 安装 Go 测试依赖
	cd backend && go get github.com/alicebob/miniredis/v2 gorm.io/driver/sqlite && go mod tidy

fmt: ## 检查代码格式
	cd backend && gofmt -l .

vet: ## 静态代码检查
	cd backend && go vet ./...

test: ## 运行所有测试
	cd backend && go test -v -count=1 ./...

build: ## 编译所有微服务
	cd backend && go build ./cmd/...

proto: ## 生成 Proto 代码
	cd backend && bash scripts/proto_gen.sh

mod: ## 整理 Go 依赖
	cd backend && go mod tidy

check: fmt vet test build ## 执行全部质量检查（fmt + vet + test + build）

docker-config: ## 验证 Docker Compose 配置
	docker compose -f deploy/compose/stack.yml config

build-frontend: build-admin build-patient ## 构建 Vue 前端（管理端 + 患者端）

build-admin: ## 构建 Vue 管理端
	cd frontend/admin && npm ci && npm run build

build-patient: ## 构建 Vue 患者端
	cd frontend/patient && npm ci && npm run build

# [deprecated] React 版本已归档至 frontend/archive/
build-admin-react: ## [deprecated] 构建 React 管理端
	cd frontend/archive/admin-react && npm ci --legacy-peer-deps && npm run build

build-patient-react: ## [deprecated] 构建 React 患者端
	cd frontend/archive/patient-react && npm ci --legacy-peer-deps && npm run build

clean: ## 清理构建产物
	rm -rf backend/bin/
	rm -rf frontend/admin/dist
	rm -rf frontend/patient/dist
	rm -rf frontend/archive/admin-react/dist
	rm -rf frontend/archive/patient-react/dist

demo-admin: ## 启动管理端演示服务
	./scripts/demo-admin.sh start

demo-admin-stop: ## 停止管理端演示服务
	./scripts/demo-admin.sh stop

demo-admin-logs: ## 查看管理端演示服务日志
	./scripts/demo-admin.sh logs

demo-patient: ## 启动患者端演示服务
	./scripts/demo-patient.sh start

demo-patient-stop: ## 停止患者端演示服务
	./scripts/demo-patient.sh stop

demo-patient-logs: ## 查看患者端演示服务日志
	./scripts/demo-patient.sh logs

demo-status: ## 查看演示服务状态
	@echo "管理端演示服务状态:"
	@-./scripts/demo-admin.sh status 2>/dev/null || echo "未运行"
	@echo ""
	@echo "患者端演示服务状态:"
	@-./scripts/demo-patient.sh status 2>/dev/null || echo "未运行"

verify-admin: ## 验证管理端演示环境
	./scripts/verify-demo.sh admin

verify-patient: ## 验证患者端演示环境
	./scripts/verify-demo.sh patient

demo-verify: ## 演示验证工具
	./scripts/demo-verify.sh all

all: check ## 默认：执行全部质量检查
