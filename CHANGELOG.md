# HIS-Go 更新日志

## [0.2.0] - 2026-06-03

### 新增

#### 后端

- 添加 `GET /api/auth/current` 接口，用于获取当前登录用户信息
- 在 `AuthService` 中添加 `GetCurrentUserInfo` 方法

#### DevOps

- 创建 `docker/docker-compose.demo-admin.yml` - 管理端演示 Profile
- 创建 `docker/docker-compose.demo-patient.yml` - 患者端演示 Profile
- 创建 `docker/.env.demo.example` - 演示环境变量模板
- 创建 `scripts/demo-admin.sh` - 管理端演示一键启停脚本
- 创建 `scripts/demo-patient.sh` - 患者端演示一键启停脚本
- 创建 `scripts/verify-demo.sh` - 演示环境验证脚本
- 创建 `scripts/demo-verify.sh` - 演示验证工具脚本
- 创建 `scripts/demo-verify-example.sh` - 演示验证工具使用示例脚本

#### 文档

- 创建 `docs/演示部署-管理端.md` - 管理端演示部署指南
- 创建 `docs/演示部署-患者端.md` - 患者端演示部署指南（含小程序配置）
- 创建 `docs/演示环境验证指南.md` - 演示环境验证指南
- 创建 `docs/演示验证报告模板.md` - 演示验证报告模板
- 创建 `docs/演示验证检查清单.md` - 演示验证检查清单
- 创建 `docs/演示验证总结.md` - 演示验证总结
- 创建 `docs/演示验证流程.md` - 演示验证流程
- 创建 `docs/演示验证工具.md` - 演示验证工具

#### 前端

- 配置管理端 `vite.config.ts` 的 `base: '/admin/'`
- 配置患者端 `vite.config.ts` 的 `base: '/patient/'`
- 配置管理端 `BrowserRouter` 的 `basename="/admin"`
- 配置患者端 `BrowserRouter` 的 `basename="/patient"`

### 更新

- 更新 `Makefile`，添加演示相关命令：
  - `make demo-admin` - 启动管理端演示
  - `make demo-admin-stop` - 停止管理端演示
  - `make demo-admin-logs` - 查看管理端演示日志
  - `make demo-patient` - 启动患者端演示
  - `make demo-patient-stop` - 停止患者端演示
  - `make demo-patient-logs` - 查看患者端演示日志
  - `make demo-status` - 查看演示服务状态

### 改进

- 项目总体完成度提升：
  - DevOps 全量编排：75% → 85%
  - 云演示部署：10% → 50%

---

## [0.1.0] - 2026-06-01

### 新增

- 初始项目结构
- 18 个 Go 微服务骨架
- PostgreSQL 数据库迁移脚本
- Redis 缓存集成
- RabbitMQ 消息队列集成
- Gateway API 网关
- JWT 认证授权
- gRPC 服务间通信
- Nginx 反向代理配置
- Docker Compose 全量编排
- React 管理端前端
- React 患者端前端
- 微信小程序壳

---

## 版本说明

- **主版本号**：重大功能变更或架构调整
- **次版本号**：新功能添加或重要改进
- **修订号**：Bug 修复或小改动
