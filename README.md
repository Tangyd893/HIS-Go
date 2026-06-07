<p align="center">
  <img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/gRPC-1.70+-2447B6?style=for-the-badge&logo=grpc&logoColor=white" alt="gRPC">
  <img src="https://img.shields.io/badge/PostgreSQL-17-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/Docker-Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker">
  <img src="https://img.shields.io/badge/SonarCloud-enabled-blue?style=for-the-badge&logo=sonarcloud&logoColor=white" alt="SonarCloud">
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License">
</p>

<h1 align="center">🏥 HIS-Go</h1>
<h3 align="center">全链路医院信息系统 · Go 微服务重构版</h3>
<p align="center">从挂号到出院，18 个微服务覆盖诊疗全流程 + 院外患者服务</p>
<p align="center">
  <b>微服务架构 → gRPC 双协议 → PostgreSQL 分库 → Docker/K8s 一键部署</b>
</p>

---

## 💡 核心理念

> **用 Go 生态重写医院信息系统，用微服务拆解医疗业务边界**

本项目是对原 [Hospital-Information-System](https://github.com/Tangyd893/Hospital-Information-System)（Spring Cloud Alibaba + Vue.js）的 **Go + Vue 3 全面重构**。将 Java 生态迁移至 Go 生态（Gin + gRPC + GORM），保留业务逻辑的同时获得更优的性能和更简洁的部署形态。

18 个微服务模块严格遵循 **Database per Service** 原则，每个服务独立数据库、独立部署、通过 gRPC 通信，Gateway 统一入口并承担 JWT 鉴权与限流。

## 🏗️ 系统架构

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  Vue3 管理端  │     │  Vue3 患者端  │     │  小程序 WebView │
│  :5173       │     │  :5174       │     │              │
└──────┬───────┘     └──────┬───────┘     └──────┬───────┘
       │                    │                    │
       └────────────┬───────┴────────────────────┘
                    ▼
            ┌──────────────┐
            │   Nginx      │
            │   :80 / :443 │
            └──────┬───────┘
                   ▼
            ┌──────────────┐
            │   Gateway    │ ← JWT 鉴权 + 限流 + CORS
            │   :8080      │
            └──────┬───────┘
                   │ gRPC
       ┌───────────┼───────────┬───────────┐
       ▼           ▼           ▼           ▼
  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐
  │  Auth   │ │  User   │ │ Clinic  │ │  ...    │
  │  :8081  │ │  :8082  │ │  :8084  │ │ 15 more │
  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘
       │           │           │           │
       ▼           ▼           ▼           ▼
  ┌─────────────────────────────────────────────┐
  │  PostgreSQL (17 DB) · Redis · RabbitMQ      │
  │  Nacos · MinIO                               │
  └─────────────────────────────────────────────┘
```

## 📁 项目结构

```
HIS-Go/
├── 📖 README.md
├── 🚫 .gitignore
├── ⚙️  Makefile
├── 📦 backend/                     # Go 后端（18 微服务）
│   ├── api/proto/                  # 18 个 gRPC Proto 定义
│   ├── cmd/                        # 18 个服务入口 main.go
│   │   ├── gateway/                #   API 网关
│   │   ├── auth/                   #   认证授权
│   │   ├── user/                   #   用户/患者/科室
│   │   ├── registration/           #   挂号预约
│   │   ├── clinic/                 #   门诊诊疗
│   │   ├── emr/                    #   电子病历 (SOAP)
│   │   ├── prescription/           #   处方管理
│   │   ├── billing/                #   收费结算
│   │   ├── pharmacy/               #   药房管理
│   │   ├── examination/            #   检查检验
│   │   ├── inpatient/              #   住院管理
│   │   ├── schedule/               #   排班管理
│   │   ├── outpatient/             #   院外患者服务
│   │   ├── followup/               #   随访管理
│   │   ├── health_record/          #   健康档案
│   │   ├── notification/           #   消息通知
│   │   ├── statistics/             #   数据统计
│   │   └── system/                 #   系统管理
│   ├── internal/                   # 各服务内部实现
│   ├── pkg/                        # 公共模块（16 个）
│   │   ├── common/                 #   雪花算法、校验、加密
│   │   ├── security/               #   JWT 解析、鉴权中间件
│   │   ├── database/               #   PostgreSQL 连接管理
│   │   ├── redis/                  #   缓存、分布式锁
│   │   ├── mq/                     #   RabbitMQ 封装
│   │   ├── grpc/                   #   gRPC 客户端/拦截器
│   │   ├── emr/                    #   EMR 公共组件
│   │   ├── middleware/             #   CORS、限流、链路追踪
│   │   ├── errors/                 #   统一错误码
│   │   ├── response/               #   统一响应
│   │   └── health/                 #   健康检查
│   ├── configs/                    # 配置文件
│   ├── sql/                        # 建库 + 种子数据
│   ├── migrations/                 # 15 个版本化迁移脚本
│   ├── scripts/                    # proto_gen / db_init / migrate
│   └── go.mod
├── 🖥️  frontend/                    # 前端
│   ├── admin/                      # ★ Vue3 管理端
│   ├── patient/                    # ★ Vue3 患者 H5
│   ├── mp-webview/                 # 微信小程序壳
│   └── archive/                    # React 版（归档）
├── 🐳 deploy/                      # Docker 部署
│   ├── compose/                    # stack.yml / demo-*.yml
│   ├── config/                     # 环境变量模板
│   ├── nginx/                      # Nginx 配置
│   └── Dockerfile
├── ☸️  k8s/                         # Kubernetes 部署清单
│   └── base/                       # 11 个 YAML（基础设施 + 微服务）
├── 🧪 testing/                     # API 集成测试
├── 📝 docs/                        # 项目文档
└── 🔧 scripts/                     # check.sh / check.ps1
```

## 🚀 快速开始

### ✅ 前置条件

| 依赖 | 必需 | 获取方式 |
|:-----|:----:|---------|
| **Go 1.25+** | ✅ | [go.dev/dl](https://go.dev/dl/) |
| **Docker 20.10+** | ✅ | [docker.com](https://docker.com) |
| **Node.js 24 LTS** | ✅ | [nodejs.org](https://nodejs.org/) |
| **Git** | ✅ | [git-scm.com](https://git-scm.com/) |

> 💡 **一条命令检查：** `go version && node --version && docker --version`

### 📦 Step 1 — 克隆 & 准备

```bash
git clone https://github.com/Tangyd893/HIS-Go.git
cd HIS-Go
```

### ⚙️ Step 2 — 配置

**Linux / macOS：**

```bash
cd deploy
cp config/stack.env.example config/stack.env
# 编辑 stack.env 设置数据库/Redis/RabbitMQ 密码
```

**Windows PowerShell：**

```powershell
cd deploy
Copy-Item config/stack.env.example config/stack.env
# 编辑 stack.env 设置密码
```

> 🤖 **让 AI 帮你：** 让 AI 根据你的环境生成 `stack.env` 配置

| 字段 | 必填 | 说明 |
|:-----|:----:|------|
| `POSTGRES_PASSWORD` | ✅ | PostgreSQL 超级用户密码 |
| `REDIS_PASSWORD` | ✅ | Redis 访问密码 |
| `RABBITMQ_PASSWORD` | ✅ | RabbitMQ 管理员密码 |
| `JWT_PRIVATE_KEY` | ⭐ | RS256 私钥（不填使用默认） |
| `MINIO_ACCESS_KEY` | ❌ | MinIO 访问密钥（有默认值） |

### ▶️ Step 3 — 启动

**一键启动全部基础设施 + 18 个微服务（最常用）：**

```bash
# 启动基础设施
docker compose -f compose/stack.yml up -d postgresql redis rabbitmq

# 初始化数据库
cd ../backend && bash scripts/db_init.sh && cd ../deploy

# 构建并启动全部服务
docker compose -f compose/stack.yml up -d --build
```

**开发模式 — 逐个启动后端服务：**

```bash
docker compose -f compose/stack.yml up -d postgresql redis rabbitmq nacos minio
cd backend && go mod tidy

# 按需启动单个服务
go run ./cmd/gateway
go run ./cmd/auth
go run ./cmd/registration
# ... 其他服务同理
```

**启动前端（开发调试）：**

```bash
# 管理端 → http://localhost:5173
cd frontend/admin && npm install && npm run dev

# 患者端 → http://localhost:5174
cd frontend/patient && npm install && npm run dev -- --host
```

## 📂 默认端口

| 组件 | 地址 | 说明 |
|:-----|------|:----|
| 🌐 API Gateway | `http://localhost:8080` | 统一入口，路由 `/api/*` |
| 🔐 his-auth | `:8081` / gRPC `:9081` | 认证授权 |
| 👤 his-user | `:8082` / gRPC `:9082` | 用户/患者/科室 |
| 📋 his-registration | `:8083` / gRPC `:9083` | 挂号预约、排队叫号 |
| 🩺 his-clinic | `:8084` / gRPC `:9084` | 门诊诊疗 |
| 📝 his-emr | `:8097` / gRPC `:9097` | 电子病历 (SOAP) |
| 💊 his-prescription | `:8085` / gRPC `:9085` | 处方管理 |
| 💰 his-billing | `:8086` / gRPC `:9086` | 收费结算 |
| 🏥 his-pharmacy | `:8087` / gRPC `:9087` | 药房管理 |
| 🔬 his-examination | `:8088` / gRPC `:9088` | 检查检验 |
| 🛏️  his-inpatient | `:8089` / gRPC `:9089` | 住院管理 |
| 📅 his-schedule | `:8090` / gRPC `:9090` | 排班管理 |
| 📱 his-outpatient | `:8091` / gRPC `:9091` | 院外患者服务 |
| 🔄 his-followup | `:8092` / gRPC `:9092` | 随访管理 |
| 📊 his-health-record | `:8093` / gRPC `:9093` | 健康档案 |
| 📮 his-notification | `:8094` / gRPC `:9094` | 消息通知 |
| 📈 his-statistics | `:8095` / gRPC `:9095` | 数据统计 |
| ⚙️  his-system | `:8096` / gRPC `:9096` | 系统管理 |

### 🏗️ 基础设施端口

| 组件 | 端口 | 说明 |
|:-----|:----:|------|
| 🐘 PostgreSQL | `5433→5432` | 宿主机 5433 映射容器 5432（避免与本机 PG 冲突） |
| 🔴 Redis | `6379` | 缓存服务 |
| 🐰 RabbitMQ | `5672` | 消息队列（管理界面 `:15672`） |
| 🪣 MinIO | `9000` | 对象存储（控制台 `:9001`） |

> **注意：** Demo Profile 将 PostgreSQL 映射到宿主机 **5433** 端口，避免与本地已安装的 PostgreSQL（默认 5432）冲突。若本机无 PG 可直接使用 5432。

## 🗄️ 数据库架构

每个微服务拥有独立的 PostgreSQL 数据库（Database per Service），共 **17 个 database**：

| 数据库 | 服务 | 数据库 | 服务 |
|:-------|:-----|:-------|:-----|
| `his_auth` | 认证授权 | `his_billing` | 收费结算 |
| `his_user` | 用户管理 | `his_pharmacy` | 药房管理 |
| `his_registration` | 挂号预约 | `his_examination` | 检查检验 |
| `his_clinic` | 门诊诊疗 | `his_inpatient` | 住院管理 |
| `his_emr` | 电子病历 | `his_schedule` | 排班管理 |
| `his_prescription` | 处方管理 | `his_outpatient` | 院外服务 |
| `his_followup` | 随访管理 | `his_health_record` | 健康档案 |
| `his_notification` | 消息通知 | `his_statistics` | 数据统计 |
| `his_system` | 系统管理 | | |

> 15 个版本化迁移脚本位于 `backend/migrations/`，`db_init.sh` 按序执行全部初始化。

## ✨ 特性

### 🔐 Gateway 统一鉴权

```
客户端 → Gateway (:8080) → JWT 验证 → 路由转发 → 下游服务
                           ├─ 白名单路由放行
                           ├─ 用户上下文透传
                           └─ 限流 + CORS
```

RS256 非对称加密，私钥签发、公钥验证，各服务零信任独立校验。

### 📨 RabbitMQ 消息可靠性

```
生产者                      Broker                      消费者
  │                          │                           │
  ├─ 本地消息表 ─────────────►│                           │
  ├─ Publisher Confirm ──────►│─ 持久化 + Quorum Queue ──►│
  │                          │                           ├─ 手动 ACK
  └─ 定时补偿重试 ◄───────────│◄── 死信队列 ◄─────────────┘
```

- **防丢失：** 持久化交换机/队列 + Quorum Queue + 手动 ACK
- **防重复：** 消息唯一 ID（雪花算法）+ DB 唯一约束 + Redis 防重

### 🏥 CDSS 临床决策支持

- 药物过敏交叉检查
- 药物相互作用检测
- 剂量合理性校验
- SOAP 结构化病历 + 模板引擎 + 三级质控

### 📊 健康检查

```bash
# Gateway 存活检查
curl http://localhost:8080/health

# 各服务就绪检查（含 DB、Redis 连通性）
curl http://localhost:8081/ready
curl http://localhost:8082/ready
```

## 🧪 测试 & 质量

```bash
cd backend

# 单元测试（158+ 用例）
go test ./...

# 编译检查
go build ./cmd/...

# 代码质量一键检查
make check
# 或
bash scripts/check.sh
```

> 质量基线：`gofmt` 零输出 ✅ | `go vet` 通过 ✅ | `go test` 全绿 ✅ | `go build` 通过 ✅

## 🔧 技术栈映射

| 原项目（Java） | HIS-Go（Go） | 说明 |
|:---------------|:-------------|:-----|
| Spring Boot 3.3 | **Gin 1.10+** | HTTP 框架 |
| Spring Cloud Gateway | **Gin + 自定义路由** | API 网关 |
| Nacos | **Nacos Go SDK** | 服务注册与配置中心 |
| OpenFeign | **gRPC 1.70+** | 微服务间 RPC |
| MyBatis-Plus | **GORM 2.x** | ORM 框架 |
| Spring Security | **golang-jwt + Gin 中间件** | 安全认证 |
| XXL-JOB | **robfig/cron 3.x** | 定时任务 |
| Maven | **Go Modules** | 依赖管理 |
| Java 21 | **Go 1.25+** | 运行语言 |

## 🖥️ 前端说明

| 维度 | Vue 3 版（当前） | React 版（归档） |
|:-----|:----------------|:----------------|
| 目录 | `frontend/admin` / `frontend/patient` | `frontend/archive/` |
| 框架 | Vue 3.5 | React 19 |
| UI 库 | Ant Design Vue 4 | Semantic UI React |
| 状态管理 | Pinia | Zustand |
| 路由 | Vue Router 4 | React Router 7 |

> 两套前端共享同一套后端 API（`/api/*`），可互相替换。

## ☸️ K8s 部署

```bash
# 部署到 Kubernetes
kubectl apply -k k8s/base/
```

`k8s/base/` 包含 11 个 YAML 清单：Namespace、ConfigMap、Secrets、PostgreSQL、Redis、RabbitMQ、Nacos、MinIO、18 微服务 Deployment+Service、Nginx Ingress。

## 🔒 生产部署注意事项

- ⚠️ 必须修改 `deploy/config/stack.env` 中所有默认密码
- 🔐 JWT 密钥对（RS256）需通过环境变量注入：
  ```bash
  openssl genrsa -out private.pem 2048
  openssl rsa -in private.pem -pubout -out public.pem
  ```
- 🚫 不要提交 `.env` 到版本仓库
- 🔒 RabbitMQ 管理端口（15672）不应对外暴露
- 🛡️ 敏感字段（手机号、身份证号）日志脱敏
- 📜 操作日志表建议定期归档
- 🌐 生产环境启用 HTTPS（Nginx 终结 TLS）

## 📋 默认验收账号

| 角色 | 用户名 | 密码 |
|:-----|:-------|:-----|
| 医生 | `demo-doctor` | `demo123` |
| 护士 | `demo-nurse` | `demo123` |
| 管理员 | `demo-admin` | `demo123` |

> 更多角色和权限数据见 `backend/sql/seed_data.sql`

## 📖 相关文档

| 文档 | 说明 |
|:-----|:-----|
| [项目架构设计文档](docs/项目架构设计文档.md) | 系统总体架构、技术选型、微服务划分 |
| [SonarCloud 质量报告](docs/sonarcloud-report.md) | 代码质量分析（待运行首次分析） |
| [云端部署指南](docs/云端部署指南.md) | 云服务器演示部署 |

## 🗺️ 后续规划

| 阶段 | 内容 | 状态 |
|:-----|:-----|:----:|
| 第一阶段 | 基础框架 + 核心服务 | ✅ 完成 |
| 第二阶段 | 院外患者服务（在线问诊、慢病管理、随访） | ✅ 完成 |
| 第三阶段 | 对接医保接口、第三方支付 | 📋 计划中 |
| 第四阶段 | 分布式事务（DTM / SAGA） | 📋 计划中 |
| 第五阶段 | K8s 编排 + CI/CD（GitHub Actions） | 📋 计划中 |
| 第六阶段 | 性能压测 + Go 协程调优 | 📋 计划中 |

---

<p align="center">
  <b>⭐ 觉得有用？点个 Star 支持一下！</b>
</p>
