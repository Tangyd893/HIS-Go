生成一个16页的中文PPT，主题：HIS-Go 医院信息系统 小程序患者端答辩演示。

风格要求：面向患者的医疗科技产品感，主色 #1890ff 蓝色 + 绿色 #00A870 点缀，浅灰背景，16:9 宽屏。移动端截图占位框，每页不超过5个要点。比管理端更亲和、更有 C 端产品感。

内容来源：GitHub 仓库 Tangyd893/HIS-Go。突出「小程序壳 + H5」混合架构与「就诊助手 RAG」亮点。

---

第1页 封面
标题：HIS-Go 患者端 — 微信小程序 + 智能就诊助手
副标题：院外患者服务 · 自助就医 · AI 分诊辅助
底部：GitHub: Tangyd893/HIS-Go
备注：本专题与管理端分开，聚焦患者视角与 C 端体验。

第2页 患者端定位
标题：为什么需要独立患者端？
要点：
- 院内：管理端服务医护人员（挂号、诊疗、发药）
- 院外：患者端服务就医人群（预约、查报告、慢病、随访）
- 渠道：H5 浏览器 + 微信小程序（覆盖主流入口）
- 目标：减少大厅排队、提升就医透明度

第3页 混合架构
标题：一套代码，双端运行
内容：画一个架构示意图：
左：微信小程序框 → web-view 箭头 → 右：Vue3 H5 页面
要点：
- 小程序壳：frontend/mp-webview — 仅负责 <web-view> 容器
- 业务 H5：frontend/patient — Vue3 + Ant Design Vue + Vite
- 配置：H5_MODE 切换 docker/vite/cloud
- Docker 模式：http://127.0.0.1/patient/（Nginx 托管）

第4页 技术栈
标题：技术选型
要点：
- 前端：Vue 3 · TypeScript · Ant Design Vue · Vue Router · Pinia
- 路由 base：/patient/
- 布局：顶部 Header + 底部 5 Tab 导航
- API：统一走 Gateway http://localhost:8080/api/...

第5页 功能地图
标题：患者端功能模块
用表格展示：
| 模块 | 路由 | 入口 |
| 首页 | /dashboard | 底部 Tab |
| 就诊助手 | /triage | 底部 Tab「助手」 |
| 预约挂号 | /appointment | 底部 Tab「挂号」 |
| 我的处方 | /prescription | 底部 Tab「处方」 |
| 健康档案 | /health-record | 底部 Tab「档案」 |
| 检查报告 | /report | 首页网格 |
| 慢病管理 | /chronic | 首页网格 |
| 我的随访 | /followup | 首页网格 |
备注：底部仅 5 Tab，报告/随访需从首页网格进入。

第6页 后端服务
标题：demo-patient.yml 服务清单
要点：
- 基础设施：PostgreSQL · Redis（无 RabbitMQ）
- 核心服务：gateway · auth · user · registration · schedule
- 患者特色：outpatient（就诊助手）· examination · followup · health_record
- 静态托管：Nginx 提供 /patient/ 与 /api/ 反代

第7页 演示环境
标题：如何启动患者端演示
要点：
- 推荐入口：微信开发者工具模拟器
- 备用入口：浏览器 http://localhost/patient
- 账号：demo-patient / demo123（绑定王小明）
- ⚠️ 微信开发者工具需勾选「不校验合法域名」
启动命令：docker compose -f deploy/compose/demo-patient.yml --env-file deploy/config/demo.env up -d --build

第8页 核心亮点：就诊助手
标题：RAG + DeepSeek 智能分诊
内容：用流程图展示：
用户输入症状 → 本地知识库语义检索（16类病症）→ 结合本院科室表 → DeepSeek 生成建议 → 输出推荐科室+就医建议
要点：
- 知识库：16 类病症（knowledge.json）
- 嵌入：SiliconFlow BGE-M3 语义检索
- 召回：关键词 + 语义双路 Top-K 合并
- 生成：DeepSeek Chat + 本院科室交集
- 降级：无 API Key 时关键词模式仍可用

第9页 现场演示1
标题：【现场演示】小程序加载与登录
表格：
| 步骤 | 环境 | 操作 |
| 1 | 微信开发者工具 | 展示 mp-webview 项目结构 |
| 2 | 模拟器 | 编译运行，显示登录页 |
| 3 | 登录 | demo-patient / demo123 |
| 4 | 首页 | 确认进入患者服务中心 |
备注：备用方案：浏览器打开 http://localhost/patient

第10页 现场演示2（重点页）
标题：【现场演示】AI 就诊助手
表格：
| 步骤 | 页面 | 操作 |
| 1 | 助手 Tab | 输入「头痛发热三天，伴有咳嗽」 |
| 2 | 聊天页 | 等待 AI 回复（3-8秒） |
| 3 | 回复内容 | 推荐科室、就医建议、免责声明 |
备注：mode=llm 完整 RAG；mode=keyword 为降级。强调：辅助分诊，不替代医生诊断。

第11页 现场演示3
标题：【现场演示】预约挂号与处方
表格：
| 步骤 | Tab | 操作 |
| 1 | 挂号 | 浏览今日排班（内科/张医生） |
| 2 | 处方 | 查看处方列表与状态 |
备注：与管理端共用 schedule/prescription 服务，数据一致。

第12页 现场演示4
标题：【现场演示】报告与慢病
表格：
| 步骤 | 入口 | 操作 |
| 1 | 首页网格 | 检查报告 → 展示 seed 报告 |
| 2 | 首页网格 | 慢病管理 → 郑十/刘十二签约 |
| 3 | 档案 Tab | 健康档案时间轴 |
备注：覆盖诊后查结果、慢病长期管理场景。

第13页 数据联动
标题：内外双环 — 与管理端数据联动
要点：
- 管理端开方 → 患者端「我的处方」可见
- 管理端排班 → 患者端「预约挂号」可见
- 统一 Gateway :8080 · 统一 JWT 体系
- 联演示建议：管理端 8 分钟 + 患者端 8 分钟
配图建议：管理端操作 → 共享 DB/服务 → 患者端展示

第14页 真机预览
标题：从模拟器到真机
要点：
- 模拟器：开发者工具本地编译（推荐答辩）
- 真机：预览扫码，H5 地址改局域网 IP
- 上线需配置：微信业务域名 HTTPS、合法 request 域名
- 当前范围：本地/局域网演示，不上云

第15页 待完善项
标题：诚实口径 — 患者端待完善项
要点：
- 底部 Tab 未覆盖报告/随访（需首页网格进入）
- 部分路由网关未注册（followup/health-record list）
- DEMO_PATIENT_MAP 仅精确映射 demo-patient 账号
- 就诊助手依赖外部 API Key，需降级预案

第16页 总结
标题：患者端价值总结
要点：
- ✅ 小程序 + H5 混合架构，一套业务双端复用
- ✅ 就诊助手 RAG 分诊，差异化亮点
- ✅ 预约、处方、报告、慢病、随访覆盖院外场景
- ⚠️ 演示版，上线需域名合规与安全加固
- 🔗 github.com/Tangyd893/HIS-Go
备注：邀请体验模拟器；可与管理端 PPT 组合为完整 HIS 汇报。
