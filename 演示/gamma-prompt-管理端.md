生成一个18页的中文PPT，主题：HIS-Go 医院信息系统 Web 管理端答辩演示。

风格要求：医疗科技风，主色 #1890ff 蓝色，白底，简洁商务，16:9 宽屏。每页不超过5个要点，文字精炼。关键页面配扁平插画或图标。

内容来源：GitHub 仓库 Tangyd893/HIS-Go，这是一个用 Go 语言重构的医院信息系统微服务项目。

---

第1页 封面
标题：HIS-Go 医院信息系统 — Web 管理端演示
副标题：全链路 Go 微服务重构版 · 门诊业务闭环
底部：GitHub: Tangyd893/HIS-Go
备注：开场30秒，说明本项目是对原 Java 版 HIS 的 Go 语言重构。

第2页 项目背景
标题：为什么用 Go 重构 HIS？
要点：
- 原系统：Spring Cloud Alibaba + Vue（Hospital-Information-System）
- 目标：17+ 微服务迁移至 Go，保持业务逻辑与分库设计
- 收益：更低资源占用、更简部署、gRPC 强类型通信
- 范围：院内诊疗 + 院外患者服务
备注：强调「重构而非从零」，降低评审对业务完整性的质疑。

第3页 系统架构
标题：分层架构一览
要点：
- 客户端：Vue3 管理端（frontend/admin）
- 接入层：Nginx :80 → Gateway :8080（JWT / 限流 / 路由）
- 服务层：18 个 Go 微服务，gRPC 内部通信
- 数据层：PostgreSQL 分库 · Redis · RabbitMQ
配图建议：画一个分层架构图，从上到下：客户端 → 网关 → 微服务 → 数据层

第4页 微服务拆分
标题：Database per Service — 核心服务
要点：
- auth 认证 · user 患者/员工/科室
- schedule 排班 · registration 挂号
- clinic 门诊诊疗 · prescription 处方
- billing 收费 · pharmacy 药房 · system 系统
备注：每个服务独立数据库，独立 cmd/<name>/main.go 入口。

第5页 技术栈
标题：技术选型
要点：
- 后端：Go 1.25 · Gin · gRPC · GORM · Wire
- 前端：Vue 3 · TypeScript · Ant Design Vue · Vite
- 基础设施：PostgreSQL 17 · Redis 7 · RabbitMQ 4
- 部署：Docker Compose 多阶段构建

第6页 部署
标题：本地演示一键启动
要点：
- Profile：deploy/compose/demo-admin.yml
- 含 PG + Redis + RabbitMQ + 10 微服务 + Gateway + Nginx
- 环境变量：deploy/config/demo.env
- 一键命令：docker compose up -d --build

第7页 演示环境
标题：演示环境说明
要点：
- 入口：http://localhost/admin
- 管理员：demo-admin / demo123
- 医生：demo-doctor / demo123
- ⚠️ 仅限本地演示，密码不可用于生产
- 演示数据来自 seed_data.sql（10 患者、5 科室、8 药品）

第8页 业务闭环总览
标题：门诊主路径 — 六步闭环
内容：用流程图展示6个步骤：
① 排班管理（发布号源）→ ② 挂号管理（预约+签到）→ ③ 门诊诊疗（主诉/诊断）→ ④ 处方管理（开方+审核）→ ⑤ 收费结算（账单支付）→ ⑥ 药房管理（发药出库）
每步标注对应微服务名。

第9页 现场演示1
标题：【现场演示】挂号与签到
内容：用表格展示操作步骤：
| 步骤 | 页面 | 操作 |
| 1 | 登录 | demo-admin 登录 |
| 2 | 挂号管理 | 今天号源 → 挂王小明 |
| 3 | 挂号记录 | 签到 → 开始接诊 |
备注：患者搜索姓名模糊匹配；号源「剩余」随挂号扣减。seed 患者：王小明（patient_001）

第10页 现场演示2
标题：【现场演示】诊疗与处方审核
表格：
| 步骤 | 页面 | 操作 |
| 4 | 门诊诊疗 | 主诉/诊断 → 保存 |
| 5 | 处方管理 | 开方（阿莫西林）→ 审核 |
备注：处方状态流转：待审核 → 已审核 → 已发药。药品来自 seed：阿莫西林胶囊。

第11页 现场演示3
标题：【现场演示】收费发药闭环
表格：
| 步骤 | 页面 | 操作 |
| 6 | 收费结算 | 待支付 → 收费 |
| 7 | 药房管理 | 选处方 → 发药 |
| 8 | 处方管理 | 确认「已发药」 |
备注：支付渠道为演示 Stub；发药扣减库存，闭环完成。

第12页 网关与安全
标题：API 网关与鉴权
要点：
- Gateway 路由：/api/auth → auth:8081，/api/user → user:8082 …
- JWT 签发：auth 服务；网关校验 Token
- 健康检查：GET /health
- 演示阶段 RBAC 简化，生产需完善

第13页 数据层
标题：PostgreSQL 分库 + Seed 数据
要点：
- 每服务独立库：his_auth、his_user、his_clinic …
- 初始化：backend/migrations/ 自动建库建表
- 演示数据：10 患者、5 科室、8 药品
- 宿主机 PG 端口：5433

第14页 前端架构
标题：Vue 3 管理端前端结构
要点：
- 路由：18 个菜单项（router/index.ts）
- 状态管理：Pinia store/auth.ts
- API 封装：api/client.ts 统一 Bearer Token
- 未启动菜单显示占位页
备注：演示时勿点 EMR/住院/检查等未启动项。

第15页 完成度
标题：功能完成度（诚实口径）
用三种颜色标签分类：
✅ 可演示：挂号、排班、诊疗、处方、收费、药房、患者/科室
⚠️ 占位 UI：电子病历、住院、检查检验、随访、通知等
🔧 演示级：支付 Stub、计费金额简化、部分详情按钮占位
备注：定位为「可脚本化本地演示」，非生产就绪。

第16页 质量保障
标题：测试与 CI
要点：
- 单元测试：go test ./... 通过
- 集成测试：testing/api/demo_admin_flow_test.go
- CI：GitHub Actions — vet / test / 前端 build / SonarCloud
- 前端构建：admin + patient 双端 npm run build 校验

第17页 后续优化
标题：已知限制与改进路线
要点：
- 首次部署：Docker 构建耗时、seed 自动化
- UX：API 失败提示、账单金额准确性
- 安全：JWT 密钥生产化、Sonar 质量门
- 集成测试纳入 CI Gate

第18页 总结
标题：总结
要点：
- ✅ Go 微服务重构 HIS，管理端门诊闭环可演示
- ✅ Docker Compose 本地可复现，前后端分离清晰
- ⚠️ 演示版，部分菜单为占位，支付/安全需生产加固
- 🔗 开源地址：github.com/Tangyd893/HIS-Go
备注：邀请提问；可备用 8 分钟精简脚本救场。
