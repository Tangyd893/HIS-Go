# HIS-Go

[中文说明] | [English](README-en.md)

HIS-Go 是一个基于 Go 语言和 Vue 3 的全链路医院信息系统（Hospital Information System，HIS）。采用前后端分离微服务架构，覆盖院内诊疗流程与院外患者服务，所有服务统一通过 Docker 容器化部署。

> 本项目是对原 [Hospital-Information-System](https://github.com/Tangyd893/Hospital-Information-System)（Spring Cloud Alibaba + Vue.js）的 **Go 语言重构版本**，将 Java 生态全面迁移至 Go 生态（Gin + gRPC + GORM）。

## 项目能力

- 18 个微服务模块，覆盖挂号→就诊→处方→收费→发药→住院全流程
- 6 个院外服务（在线问诊、慢病管理、健康档案、随访管理）
- Gateway 统一入口，路由转发，JWT 鉴权中间件已接入
- PostgreSQL 17 分库设计（Database per Service），17 个独立数据库（Gateway 无状态）
- GORM 持久化 + 乐观锁 + 逻辑删除 + 自动填充
- RabbitMQ 消息可靠性设计（Publisher Confirm、手动 ACK、死信队列、本地消息表）
- Redis 缓存策略（号源缓存、分布式锁、排队 Sorted Set、滑动窗口限流）
- CDSS 临床决策支持（药物过敏/相互作用/剂量校验）
- SOAP 结构化病历（模板引擎、质控流程、FHIR R4 对接）
- Docker Compose 一键部署（PostgreSQL、Redis、RabbitMQ、Nacos、MinIO、Nginx）

## 当前完成度

> 本项目处于**交付就绪**阶段。18个服务可编译、158+个测试保护核心链路、gRPC/HTTP双协议完整、前端管理端+患者端已搭建（66个Vue/TS文件）、Docker Compose+K8s部署方案齐全。

| 维度          | 状态     | 说明                                                      |
| ----------- | ------ | ------------------------------------------------------- |
| 后端骨架        | 已完成    | 18 个服务入口、17 个领域模块，handler/service/repository/model 分层完整 |
| 构建基线        | 已通过    | `go build ./cmd/...` 通过，`gofmt -l .` 零输出，`go.sum` 已生成 |
| gRPC        | 已完成    | 18 个 `.proto` 已编写，gRPC Server 注册和 `.pb.go` 代码已生成           |
| 数据库         | 已补齐    | 15 个迁移脚本覆盖全部 17 个数据库，表结构完整匹配 Go 模型                  |
| Docker      | 配置就绪   | Dockerfile 和 docker-compose.yml/docker-compose.prod.yml 已就绪，Compose 配置解析已通过              |
| Gateway JWT | 已接入    | JWT 鉴权中间件已接入 Gateway，白名单路由放行，用户上下文透传至下游服务               |
| 健康检查        | 已增强    | `/health` 存活检查 + `/ready` 就绪检查（含 DB、Redis 连通性）          |
| 数据库迁移       | 已补齐    | 15 个迁移脚本覆盖全部 17 个数据库，`db_init.sh` 支持按序执行             |
| 测试          | 已建立    | 158+ 个单元测试覆盖 JWT/Redis锁/Auth/挂号/收费/药房/用户/排班/门诊/处方/随访/慢病/EMR/检查/错误码 |
| 质量检查       | 已固化    | `gofmt`零输出 + `go vet`通过 + `go test`全绿 + `go build`通过        |
| 前端          | 已完成    | `his-web-admin`(管理端19模块) + `his-web-patient`(患者端8模块)，Vue3+TS+Vite6 |
| K8s 部署      | 已就绪    | `k8s/base/` 含 11 个 YAML 清单，覆盖全部基础设施和 18 微服务 |

## 技术栈

- **后端：** Go 1.25+、Gin 1.10+、gRPC 1.70+、GORM 2.x、PostgreSQL 17、Redis 7、RabbitMQ 4、Nacos Go SDK、MinIO
- **前端：** Vue 3.5、TypeScript、Ant Design Vue 4、Element Plus 2、ECharts 5、Vite 6、Pinia 2
- **依赖注入：** Wire 0.6+
- **定时任务：** robfig/cron 3.x
- **日志：** Zap
- **配置：** Viper + Nacos 配置中心
- **可观测性：** 待接入 OpenTelemetry + Prometheus + Grafana

## 默认端口

| 组件                | 地址                            | 说明                                   |
| ----------------- | ----------------------------- | ------------------------------------ |
| API Gateway       | `http://localhost:8080`       | 统一入口，路由 `/api/*`                     |
| his-auth          | `http://localhost:8081`       | 认证授权 (gRPC:9081)                     |
| his-user          | `http://localhost:8082`       | 用户/患者/科室管理 (gRPC:9082)               |
| his-registration  | `http://localhost:8083`       | 挂号预约、排队叫号 (gRPC:9083)                |
| his-clinic        | `http://localhost:8084`       | 门诊诊疗 (gRPC:9084)                     |
| his-prescription  | `http://localhost:8085`       | 处方管理 (gRPC:9085)                     |
| his-billing       | `http://localhost:8086`       | 收费结算 (gRPC:9086)                     |
| his-pharmacy      | `http://localhost:8087`       | 药房管理 (gRPC:9087)                     |
| his-examination   | `http://localhost:8088`       | 检查检验 (gRPC:9088)                     |
| his-inpatient     | `http://localhost:8089`       | 住院管理 (gRPC:9089)                     |
| his-schedule      | `http://localhost:8090`       | 排班管理 (gRPC:9090)                     |
| his-outpatient    | `http://localhost:8091`       | 院外患者服务 (gRPC:9091)                   |
| his-followup      | `http://localhost:8092`       | 随访管理 (gRPC:9092)                     |
| his-health-record | `http://localhost:8093`       | 健康档案 (gRPC:9093)                     |
| his-notification  | `http://localhost:8094`       | 消息通知 (gRPC:9094)                     |
| his-statistics    | `http://localhost:8095`       | 数据统计 (gRPC:9095)                     |
| his-system        | `http://localhost:8096`       | 系统管理 (gRPC:9096)                     |
| his-emr           | `http://localhost:8097`       | 电子病历 (gRPC:9097)                     |
| PostgreSQL        | `localhost:5432`              | `his_admin / change_me_123`          |
| Redis             | `localhost:6379`              | 密码 `change_me_456`                   |
| RabbitMQ          | `localhost:5672`              | 管理端口 `15672`，`admin / change_me_789` |
| Nacos             | `http://localhost:8848/nacos` | `nacos / nacos`                      |
| MinIO Console     | `http://localhost:9001`       | `minioadmin / change_me_012`         |

## 微服务边界

| 服务                | 数据库                 | 主要职责                           |
| ----------------- | ------------------- | ------------------------------ |
| Gateway           | 无状态                 | 统一入口、路由转发、JWT 认证、CORS、限流       |
| his-auth          | `his_auth`          | 登录认证、Token 签发/刷新（RS256）、角色权限管理 |
| his-user          | `his_user`          | 患者档案、员工管理、科室树                  |
| his-registration  | `his_registration`  | 号源管理、挂号预约、排队叫号、Redis 分布式锁      |
| his-clinic        | `his_clinic`        | 接诊登记、诊断录入（ICD-10）、检查申请、转诊      |
| his-emr           | `his_emr`           | SOAP 结构化病历、模板引擎、三级质控、CDSS      |
| his-prescription  | `his_prescription`  | 处方开具、审核、退回、处方状态流转              |
| his-billing       | `his_billing`       | 多类型费用合并结算、支付、退费审批、日报表          |
| his-pharmacy      | `his_pharmacy`      | 药品库存、入库、发药、效期预警（cron 定时）       |
| his-examination   | `his_examination`   | 检查执行、报告录入、审核流程                 |
| his-inpatient     | `his_inpatient`     | 入院登记、床位分配、医嘱下达、护理记录、出院结算       |
| his-schedule      | `his_schedule`      | 医生排班、诊室安排、号源生成（乐观锁）            |
| his-outpatient    | `his_outpatient`    | 在线问诊、消息记录、慢病签约、健康自测            |
| his-followup      | `his_followup`      | 随访计划自动生成、执行记录、满意度调查            |
| his-health-record | `his_health_record` | 全生命周期健康档案总览、时间轴                |
| his-notification  | `his_notification`  | 通知模板管理、SMS/邮件/站内信发送            |
| his-statistics    | `his_statistics`    | 运营报表、挂号/收入趋势、科室工作量、医疗质量        |
| his-system        | `his_system`        | 字典类型/字典项管理、参数配置、操作日志审计         |

## 快速开始

### 环境要求

| 工具      | 最低版本     | 用途                                       |
| ------- | -------- | ---------------------------------------- |
| Docker  | 20.10+   | PostgreSQL、Redis、RabbitMQ、Nacos、MinIO 容器 |
| Go      | 1.25+    | Go 编译与运行                                 |
| Node.js | 24 (LTS) | 前端构建与开发服务器                               |

一条命令检查所有必需工具：

```bash
go version && node --version && docker --version
```

### 启动完整后端技术栈

```bash
cd docker
cp .env.example .env
# 编辑 .env 设置数据库密码和 JWT 密钥

# 启动基础设施
docker compose up -d postgresql redis rabbitmq nacos minio

# 初始化数据库（等 PostgreSQL 就绪后，推荐方式）
cd ../backend
# 如修改了 .env 中的 DB_PASSWORD，请在此同步导出：export DB_PASSWORD=你的密码
bash scripts/db_init.sh
cd ../docker

# 构建 Go 服务并启动全部服务
docker compose up -d --build
```

备选：也可以在项目根目录手动执行 SQL：

```bash
docker exec -i his-postgres psql -U his_admin < backend/sql/init_all.sql
docker exec -i his-postgres psql -U his_admin < backend/sql/seed_data.sql
```

> 当前已验证 Compose 配置可解析；完整容器启动和集成测试仍需在有 Docker 运行环境时执行。

### 逐个启动后端服务（开发调试）

先启动基础设施，再分别在独立终端运行各服务：

```bash
docker compose -f docker/docker-compose.yml up -d postgresql redis rabbitmq nacos minio

# 下载依赖
cd backend
go mod tidy

# 启动各服务
go run ./cmd/gateway
go run ./cmd/auth
go run ./cmd/user
go run ./cmd/registration
go run ./cmd/clinic
go run ./cmd/prescription
go run ./cmd/billing
go run ./cmd/pharmacy
go run ./cmd/examination
go run ./cmd/inpatient
go run ./cmd/schedule
go run ./cmd/outpatient
go run ./cmd/followup
go run ./cmd/health_record
go run ./cmd/notification
go run ./cmd/statistics
go run ./cmd/system
go run ./cmd/emr
```

### 启动前端

```bash
# 管理端
cd frontend/his-web-admin
npm install
npm run dev      # 开发服务器 → http://localhost:5173

# 患者端
cd frontend/his-web-patient
npm install
npm run dev      # 开发服务器 → http://localhost:5174
```

> 前端通过 Vite 代理将 `/api` 请求转发至 `http://localhost:8080`（Gateway），确保后端已启动。

## 默认验收账号

| 角色  | 用户名           | 密码        |
| --- | ------------- | --------- |
| 医生  | `demo-doctor` | `demo123` |
| 护士  | `demo-nurse`  | `demo123` |
| 管理员 | `demo-admin`  | `demo123` |

> 更多角色和权限数据见 `backend/sql/seed_data.sql`

## 健康检查自检

```bash
# Gateway 健康检查
curl http://localhost:8080/health

# 各服务独立健康检查
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8097/health

# 测试 API 连通性
curl http://localhost:8080/api/ping
```

## 项目目录结构

本项目按职责划分为四个顶层子项目：后端、前端、Docker 部署和测试。

```
HIS-Go/
├── backend/                        # Go 后端子项目
│   ├── api/                        # Proto 接口定义
│   │   └── proto/                  # 18 个服务的 gRPC 接口定义
│   ├── cmd/                        # 各服务入口 main.go
│   │   ├── gateway/                # 网关服务
│   │   ├── auth/                   # 认证服务
│   │   └── ...                     # 其他 15 个微服务
│   ├── internal/                   # 各服务内部实现（不对外暴露）
│   │   ├── auth/                   # handler / service / repository / model
│   │   └── ...
│   ├── pkg/                        # 公共模块（可被外部导入）
│   │   ├── common/                 # 雪花算法、校验、加密
│   │   ├── security/               # JWT 解析、鉴权中间件
│   │   ├── database/               # PostgreSQL 连接管理
│   │   ├── redis/                  # Redis 缓存、分布式锁
│   │   ├── mq/                     # RabbitMQ 封装
│   │   ├── grpc/                   # gRPC 客户端/拦截器
│   │   ├── emr/                    # EMR 公共组件
│   │   ├── logger/                 # Zap 日志
│   │   ├── config/                 # Viper 配置
│   │   ├── middleware/              # CORS、限流、链路追踪、恢复
│   │   ├── errors/                 # 统一错误码
│   │   ├── response/               # 统一响应
│   │   └── health/                 # 健康检查
│   ├── configs/                    # 配置文件
│   │   └── config.yaml             # 默认配置
│   ├── sql/                        # 数据库初始化脚本
│   │   ├── init_all.sql            # 全量建库脚本
│   │   └── seed_data.sql           # 基础数据/字典数据
│   ├── scripts/                    # 辅助脚本
│   │   ├── proto_gen.sh            # Proto 代码生成
│   │   ├── db_init.sh              # 数据库初始化（迁移+种子）
│   │   └── migrate.sh              # 数据库迁移执行器
│   ├── migrations/                 # 版本化数据库迁移
│   │   ├── 001_init_all_db.sql     # 创建所有数据库
│   │   └── ...                     # 15 个迁移脚本
│   ├── go.mod                      # Go Module 定义
│   └── go.sum                      # 依赖校验
├── docker/                         # Docker 部署配置
│   ├── docker-compose.yml          # 开发环境编排
│   ├── Dockerfile                  # Go 多阶段构建
│   ├── nginx/                      # Nginx 配置
│   │   └── nginx.conf
│   └── .env.example                # 环境变量模板
├── frontend/                       # 前端子项目
│   ├── his-web-admin/              # 管理端 (Vue3 + Ant Design Vue4)
│   │   ├── src/views/              # 19 个功能模块页面
│   │   └── ...
│   └── his-web-patient/            # 患者端 (Vue3 + Ant Design Vue4, H5)
│       ├── src/views/              # 8 个功能模块页面
│       └── ...
├── k8s/                            # Kubernetes 部署清单
│   └── base/
│       ├── namespace.yaml
│       ├── configmap.yaml
│       ├── secrets.yaml
│       ├── postgresql.yaml
│       ├── redis.yaml
│       ├── rabbitmq.yaml
│       ├── nacos.yaml
│       ├── minio.yaml
│       ├── services.yaml           # 18 微服务 Deployment+Service
│       ├── nginx.yaml              # Nginx + Ingress
│       └── kustomization.yaml
├── scripts/                        # 顶层脚本
│   ├── check.sh                    # Linux/macOS 质量检查
│   └── check.ps1                   # Windows 质量检查
├── Makefile                        # 构建/检查快捷命令
├── testing/                        # API 集成测试
│   ├── go.mod                      # 独立 Go module
│   ├── api/                        # 集成测试用例
│   │   ├── client.go               # HTTP 客户端封装
│   │   └── auth_flow_test.go       # 认证流程 / 鉴权验收
│   └── run.sh                      # 集成测试运行脚本
├── docs/                           # 项目文档
│   └── 项目架构设计文档.md
├── .gitignore
└── README.md
```

## 数据库说明

每个微服务拥有独立的 PostgreSQL 数据库（Database per Service），共 17 个 database：

`his_auth` `his_user` `his_registration` `his_clinic` `his_emr` `his_prescription` `his_billing` `his_pharmacy` `his_examination` `his_inpatient` `his_schedule` `his_outpatient` `his_followup` `his_health_record` `his_notification` `his_statistics` `his_system`

建表脚本位于 `backend/sql/init_all.sql`，种子数据位于 `backend/sql/seed_data.sql`。

版本化迁移脚本位于 `backend/migrations/`，按编号顺序执行即可初始化全部表结构。

## RabbitMQ 消息可靠性设计要点

本项目在架构设计中重点阐述了消息可靠性的三大保障策略：

1. **生产者可靠发送：** 本地消息表（Transactional Outbox）+ Publisher Confirm + 定时补偿重试
2. **防止消息丢失：** 持久化交换机/队列/消息 + Quorum Queue + 消费者手动 ACK + 死信队列
3. **防止重复消费：** 消息唯一 ID（雪花算法）+ 数据库唯一约束 + Redis 防重标记 + 业务幂等设计

详见 [项目架构设计文档](docs/项目架构设计文档.md) 第八章。

## 测试命令

### 后端单元测试

```bash
cd backend

# 运行所有测试
go test ./...

# 运行特定服务测试
go test ./pkg/...
go test ./internal/auth/...

# 编译检查所有服务
go build ./cmd/...

# 代码质量一键检查
bash scripts/check.sh
```

### 集成测试（需 Docker 环境启动后）

```bash
# 启动全部服务后执行 API 验收测试
cd testing
bash run.sh

# 或手动指定目标地址
HIS_BASE_URL=http://localhost:8080 HIS_INTEGRATION_TEST=true go test -v ./...
```

### 编译

```bash
cd backend

# 编译所有服务
go build -o bin/ ./cmd/...

# 单个服务编译
go build -o bin/gateway ./cmd/gateway
```

## 与原项目技术栈映射

| 原项目（Java）            | HIS-Go（Go）          | 说明          |
| -------------------- | ------------------- | ----------- |
| Spring Boot 3.3      | Gin 1.10+           | HTTP 框架     |
| Spring Cloud Gateway | Gin + 自定义路由         | API 网关      |
| Nacos                | Nacos Go SDK        | 服务注册与配置中心   |
| OpenFeign            | gRPC                | 微服务间 RPC 通信 |
| MyBatis-Plus         | GORM 2.x            | ORM 框架      |
| Seata                | — (后续引入 DTM)        | 分布式事务       |
| Sentinel             | 自定义限流中间件            | 熔断降级        |
| XXL-JOB              | robfig/cron         | 定时任务        |
| Maven                | Go Modules          | 依赖管理        |
| Java 21              | Go 1.25+            | 运行语言        |
| Spring Security      | golang-jwt + Gin中间件 | 安全认证        |

## 生产部署注意事项

- 生产环境必须修改 `docker/.env` 中所有默认密码（数据库、Redis、RabbitMQ、MinIO）
- JWT 密钥对（RS256）需通过环境变量注入，生成方式：
  
  ```bash
  openssl genrsa -out private.pem 2048
  openssl rsa -in private.pem -pubout -out public.pem
  ```
- 不要提交 `docker/.env` 到版本仓库
- 生产环境 RabbitMQ 管理端口（15672）不应对外暴露
- 敏感字段（手机号、身份证号）写入日志前需脱敏处理
- 操作日志表建议定期归档
- 生产环境建议启用 HTTPS（Nginx 终结 TLS）

## 相关文档

| 文档                           | 说明                             |
| ---------------------------- | ------------------------------ |
| [项目架构设计文档](docs/项目架构设计文档.md) | 系统总体架构、技术选型、微服务划分、RabbitMQ 可靠性 |

## 后续规划

1. **第一阶段**：搭建基础框架，完成网关、认证、用户、挂号、诊疗、处方核心服务（当前阶段）
2. **第二阶段**：完善院外患者服务（在线问诊、慢病管理、随访）
3. **第三阶段**：对接医保接口、第三方支付
4. **第四阶段**：引入分布式事务方案（DTM 或 SAGA 模式）
5. **第五阶段**：Kubernetes 容器编排迁移，CI/CD 流水线搭建（GitHub Actions）
6. **第六阶段**：性能压测与 Go 协程调优
