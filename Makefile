.PHONY: help fmt vet test build proto clean check all deps

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
	docker compose -f docker/docker-compose.yml config

build-frontend: ## 构建 React 前端（管理端 + 患者端）
	cd frontend/his-web-admin-react && npm ci --legacy-peer-deps && npm run build
	cd frontend/his-web-patient-react && npm ci --legacy-peer-deps && npm run build

build-admin-react: ## 构建 React 管理端
	cd frontend/his-web-admin-react && npm ci --legacy-peer-deps && npm run build

build-patient-react: ## 构建 React 患者端
	cd frontend/his-web-patient-react && npm ci --legacy-peer-deps && npm run build

clean: ## 清理构建产物
	rm -rf backend/bin/
	rm -rf frontend/his-web-admin-react/dist
	rm -rf frontend/his-web-patient-react/dist

all: check ## 默认：执行全部质量检查
