# HIS-Go Web 管理端 — PPT 生成提示词

> **用途**：将本文全文提供给 PPT/幻灯片生成工具（如 Gamma、MindShow、Copilot 等），并附上 GitHub 仓库 https://github.com/Tangyd893/HIS-Go 作为内容来源。  
> **配套演示脚本**：`演示流程-15分钟.md`、`演示流程-8分钟精简.md`  
> **建议总页数**：18–22 页 · **风格**：医疗科技、蓝白主色、简洁商务 · **比例**：16:9

---

## 全局生成指令（放在提示词最前面）

```
请根据以下分页提示，为「HIS-Go 医院信息系统 — Web 管理端」生成答辩/汇报 PPT。
内容来源：GitHub 仓库 Tangyd893/HIS-Go（Go 微服务 HIS 重构项目）。
风格：医疗信息化、专业可信；主色 #1890ff + 白底；每页不超过 5 个要点；关键术语中英文并列一次即可。
每页需包含：标题、3-5 条正文要点、1 条「演讲备注」、若标注了「现场演示」则加「演示动作」框。
不要虚构仓库中不存在的功能；未实现模块标注为「规划中/演示占位」。
```

---

## 第 1 页｜封面

**标题**：HIS-Go 医院信息系统 — Web 管理端演示

**要点**
- 全链路医院信息系统 · Go 微服务重构版
- 汇报人 / 日期 / 单位（留空占位）
- GitHub：Tangyd893/HIS-Go

**演讲备注**：开场 30 秒，说明本项目是对原 Java 版 HIS 的 Go 语言重构。

**仓库引用**：`README.md` 项目标题与徽章区

---

## 第 2 页｜项目背景与重构动机

**标题**：为什么用 Go 重构 HIS？

**要点**
- 原系统：Spring Cloud Alibaba + Vue（Hospital-Information-System）
- 目标：17+ 微服务迁移至 Go，保持业务逻辑与分库设计
- 收益：更低资源占用、更简部署、gRPC 强类型通信
- 范围：院内诊疗 + 院外患者服务

**演讲备注**：强调「重构而非从零」，降低评审对业务完整性的质疑。

**仓库引用**：`docs/项目架构设计文档.md` §1.1 重构目标

---

## 第 3 页｜系统总体架构

**标题**：分层架构一览

**要点**
- 客户端：Vue3 管理端（`frontend/admin`）
- 接入：Nginx :80 → Gateway :8080（JWT / 限流 / 路由）
- 服务层：18 个 Go 微服务，gRPC 内部通信
- 数据层：PostgreSQL 分库 · Redis · RabbitMQ

**图示建议**：引用 README 架构 ASCII 图重绘为分层框图

**仓库引用**：`README.md` 系统架构图、`backend/cmd/` 目录列表

---

## 第 4 页｜微服务拆分（管理端相关）

**标题**：Database per Service — 管理端核心服务

**要点**
- `auth` 认证 · `user` 患者/员工/科室
- `schedule` 排班 · `registration` 挂号
- `clinic` 门诊诊疗 · `prescription` 处方
- `billing` 收费 · `pharmacy` 药房 · `system` 系统

**演讲备注**：每个服务独立数据库，独立 `cmd/<name>/main.go` 入口。

**仓库引用**：`deploy/compose/demo-admin.yml` 服务列表

---

## 第 5 页｜技术栈

**标题**：技术选型

**要点**
- 后端：Go 1.25 · Gin · gRPC · GORM · Wire
- 前端：Vue 3 · TypeScript · Ant Design Vue · Vite
- 基础设施：PostgreSQL 17 · Redis 7 · RabbitMQ 4
- 部署：Docker Compose 多阶段构建（`deploy/Dockerfile`）

**仓库引用**：`backend/go.mod`、`frontend/admin/package.json`

---

## 第 6 页｜部署形态（演示环境）

**标题**：本地演示一键启动

**要点**
- Profile：`deploy/compose/demo-admin.yml`
- 含 PG + Redis + RabbitMQ + 10 微服务 + Gateway + Nginx
- 环境变量：`deploy/config/demo.env`
- 管理端静态资源：`frontend/admin/dist` 由 Nginx 托管

**代码块（简化）**：
```bash
docker compose -f deploy/compose/demo-admin.yml --env-file deploy/config/demo.env up -d --build
```

**仓库引用**：`deploy/README.md`、`docs/快速启动-管理端.md`

---

## 第 7 页｜演示账号与环境

**标题**：演示环境说明

**要点**
- 入口：http://localhost/admin
- 账号：`demo-admin` / `demo123`（管理员）
- 辅助账号：`demo-doctor` / `demo123`（医生）
- ⚠️ 仅限本地演示，密码不可用于生产

**演讲备注**：提前说明演示数据来自 `backend/sql/seed_data*.sql`。

**仓库引用**：`docs/快速启动-管理端.md` 演示账号表

---

## 第 8 页｜业务闭环总览

**标题**：门诊主路径 — 六步闭环

**要点**
1. 排班管理 — 发布号源
2. 挂号管理 — 预约 + 签到
3. 门诊诊疗 — 主诉/诊断
4. 处方管理 — 开方 + 审核
5. 收费结算 — 账单支付
6. 药房管理 — 发药出库

**图示建议**：横向流程箭头图，每步标注对应微服务名

**仓库引用**：`演示/管理端/演示流程-15分钟.md`

---

## 第 9 页｜现场演示 ① — 登录与挂号

**标题**：【现场演示】挂号与签到

**演示动作**（与 PPT 同屏或下一屏）
| 步骤 | 页面 | 操作 |
|------|------|------|
| 1 | 登录 | demo-admin 登录 |
| 2 | 挂号管理 | 今天号源 → 挂王小明 |
| 3 | 挂号记录 | 签到 → 开始接诊 |

**讲解要点**
- 患者搜索：姓名模糊匹配
- 号源「剩余」随挂号扣减
- seed 患者：王小明（patient_001）

**仓库引用**：`frontend/admin/src/views/registration/RegistrationView.vue`

---

## 第 10 页｜现场演示 ② — 接诊与处方

**标题**：【现场演示】诊疗与处方审核

**演示动作**
| 步骤 | 页面 | 操作 |
|------|------|------|
| 4 | 门诊诊疗 | 主诉/诊断 → 保存 |
| 5 | 处方管理 | 开方（阿莫西林）→ 审核 |

**讲解要点**
- 处方状态：待审核 → 已审核 → 已发药
- 药品数据来自 seed：`drug_001` 阿莫西林胶囊

**仓库引用**：`frontend/admin/src/views/clinic/ClinicView.vue`、`prescription/PrescriptionView.vue`

---

## 第 11 页｜现场演示 ③ — 收费与发药

**标题**：【现场演示】收费发药闭环

**演示动作**
| 步骤 | 页面 | 操作 |
|------|------|------|
| 6 | 收费结算 | 待支付 → 收费 |
| 7 | 药房管理 | 选处方 → 发药 |
| 8 | 处方管理 | 确认「已发药」 |

**讲解要点**
- 支付渠道为演示 Stub（`pkg/payment`）
- 发药扣减库存，闭环完成

**仓库引用**：`frontend/admin/src/views/billing/BillingView.vue`、`pharmacy/PharmacyView.vue`

---

## 第 12 页｜API 网关与鉴权

**标题**：统一入口与安全

**要点**
- Gateway 路由：`/api/auth` → auth:8081，`/api/user` → user:8082 …
- JWT 签发：auth 服务；网关校验 Token
- 健康检查：`GET /health`
- 演示阶段 RBAC 简化，生产需完善 `RequireRole`

**仓库引用**：`backend/cmd/gateway/main.go`

---

## 第 13 页｜数据层设计

**标题**：PostgreSQL 分库 + Seed 数据

**要点**
- 每服务独立库：`his_auth`、`his_user`、`his_clinic` …
- 初始化：`backend/migrations/` 自动建库建表
- 演示数据：`backend/sql/seed_data.sql`（10 患者、5 科室、8 药品）
- 宿主机 PG 端口：5433（避免与本机冲突）

**仓库引用**：`backend/sql/seed_data.sql`、`deploy/compose/demo-admin.yml`

---

## 第 14 页｜前端架构（管理端）

**标题**：Vue 3 管理端结构

**要点**
- 路由：`frontend/admin/src/router/index.ts`（18 菜单）
- 状态：Pinia `store/auth.ts`
- API 封装：`api/client.ts` 统一 Bearer Token
- 演示 Profile 仅 10 个后端 — 未启动菜单显示占位页

**演讲备注**：演示时勿点 EMR/住院/检查等未启动项。

**仓库引用**：`frontend/admin/src/router/index.ts`、`views/ServicePlaceholder.vue`

---

## 第 15 页｜已实现 vs 演示占位

**标题**：功能完成度（诚实口径）

**要点**
| 状态 | 模块 |
|------|------|
| ✅ 可演示 | 挂号、排班、诊疗、处方、收费、药房、患者/科室 |
| ⚠️ 占位 UI | 电子病历、住院、检查检验、随访、通知等 |
| 🔧 演示级 | 支付 Stub、计费金额简化、部分详情按钮占位 |

**演讲备注**：定位为「可脚本化本地演示」，非生产就绪。

**仓库引用**：`docs/待办清单.md` PM 验收挑刺章节

---

## 第 16 页｜质量保障

**标题**：测试与 CI

**要点**
- 单元测试：`go test ./...` 通过
- 集成测试：`testing/api/demo_admin_flow_test.go`（需 Docker）
- CI：`.github/workflows/ci.yml` — vet / test / 前端 build / SonarCloud
- 前端构建：admin + patient 双端 `npm run build` 校验

**仓库引用**：`.github/workflows/ci.yml`、`testing/api/`

---

## 第 17 页｜已知限制与改进路线

**标题**：后续优化方向

**要点**
- 首次部署：Docker 构建耗时、seed 自动化（db-init）
- UX：API 失败提示、账单金额、Dashboard 服务状态准确性
- 安全：演示密码声明、JWT 密钥生产化、Sonar 质量门
- 集成测试纳入 CI Gate

**仓库引用**：`docs/待办清单.md`

---

## 第 18 页｜总结

**标题**：总结

**要点**
- ✅ Go 微服务重构 HIS，管理端门诊闭环可演示
- ✅ Docker Compose 本地可复现，前后端分离清晰
- ⚠️ 演示版，部分菜单为占位，支付/安全需生产加固
- 🔗 开源地址：github.com/Tangyd893/HIS-Go

**演讲备注**：邀请提问；可备用 8 分钟精简脚本救场。

---

## 附录：8 分钟答辩精简版页序（可选替换 8–11 页）

若时间只允许 8 分钟，将第 8–11 页合并为 **1 页「现场演示速查表」**：

| 时间 | 菜单 | 关键动作 |
|------|------|----------|
| 0:30 | 挂号管理 | 挂王小明 → 签到 → 接诊 |
| 2:30 | 门诊诊疗 | 保存主诉诊断 |
| 3:30 | 处方管理 | 开方审核 |
| 5:00 | 收费结算 | 收费 |
| 6:00 | 药房管理 | 发药 |

配套文件：`演示流程-8分钟精简.md`

---

## 附录：生成工具参数建议

```
语言：简体中文
页数：18
模板：科技/医疗
字体：标题 微软雅黑 Bold / 正文 微软雅黑 Regular
每页动画：淡入（可选，答辩现场建议无动画）
配图：架构图用扁平图标；流程图用箭头；避免_stock 手术照片
```
