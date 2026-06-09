# HIS-Go 小程序患者端 — PPT 生成提示词

> **用途**：将本文全文提供给 PPT 生成工具，并附上 GitHub 仓库 https://github.com/Tangyd893/HIS-Go  
> **配套演示脚本**：`演示流程.md`  
> **建议总页数**：16–20 页 · **风格**：移动端医疗、蓝绿主色、偏 C 端亲和 · **比例**：16:9

---

## 全局生成指令（放在提示词最前面）

```
请根据以下分页提示，为「HIS-Go 医院信息系统 — 小程序患者端」生成答辩/汇报 PPT。
内容来源：GitHub 仓库 Tangyd893/HIS-Go。
风格：面向患者的医疗科技产品感；主色 #1890ff + 浅灰背景；移动端截图占位框；每页不超过 5 要点。
每页需包含：标题、正文要点、演讲备注；标注「现场演示」的页面必须含「演示动作」表。
突出「小程序壳 + H5」混合架构与「就诊助手 RAG」亮点。
不要虚构未实现功能；网关未注册路由标注为「演示预留」。
```

---

## 第 1 页｜封面

**标题**：HIS-Go 患者端 — 微信小程序 + 智能就诊助手

**要点**
- 院外患者服务 · 自助就医 · AI 分诊辅助
- 小程序 WebView + Vue3 H5 混合架构
- GitHub：Tangyd893/HIS-Go

**演讲备注**：本专题与管理端分开，聚焦患者视角与 C 端体验。

**仓库引用**：`README.md`、`frontend/patient/`、`frontend/mp-webview/`

---

## 第 2 页｜患者端定位

**标题**：为什么需要独立患者端？

**要点**
- 院内：管理端服务医护人员（挂号、诊疗、发药）
- 院外：患者端服务就医人群（预约、查报告、慢病、随访）
- 渠道：H5 浏览器 + 微信小程序（覆盖主流入口）
- 目标：减少大厅排队、提升就医透明度

**仓库引用**：`docs/项目架构设计文档.md` 院外患者服务章节

---

## 第 3 页｜混合架构：小程序壳 + H5

**标题**：一套代码，双端运行

**要点**
- 小程序壳：`frontend/mp-webview` — 仅负责 `<web-view>` 容器
- 业务 H5：`frontend/patient` — Vue3 + Ant Design Vue + Vite
- 配置：`pages/index/index.js` 中 `H5_MODE` 切换 docker/vite/cloud
- Docker 模式：`http://127.0.0.1/patient/`（Nginx 托管 dist）

**图示建议**：左小程序框 → web-view 箭头 → 右 H5 页面

**仓库引用**：`frontend/mp-webview/pages/index/index.js`

---

## 第 4 页｜患者端技术栈

**标题**：技术选型

**要点**
- 前端：Vue 3 · TypeScript · Ant Design Vue · Vue Router · Pinia
- 路由 base：`/patient/`（`router/index.ts`）
- 布局：顶部 Header + 底部 5 Tab 导航（`DefaultLayout.vue`）
- API：统一走 Gateway `http://localhost:8080/api/...`

**仓库引用**：`frontend/patient/package.json`、`src/layouts/DefaultLayout.vue`

---

## 第 5 页｜功能地图

**标题**：患者端功能模块

**要点**
| 模块 | 路由 | 入口 |
|------|------|------|
| 首页 | /dashboard | 底部 Tab |
| 就诊助手 | /triage | 底部 Tab「助手」 |
| 预约挂号 | /appointment | 底部 Tab「挂号」 |
| 我的处方 | /prescription | 底部 Tab「处方」 |
| 健康档案 | /health-record | 底部 Tab「档案」 |
| 检查报告 | /report | 首页网格 |
| 慢病管理 | /chronic | 首页网格 |
| 我的随访 | /followup | 首页网格 |

**演讲备注**：底部仅 5 Tab，报告/随访需从首页进入——演示时提前熟悉。

**仓库引用**：`frontend/patient/src/router/index.ts`

---

## 第 6 页｜后端服务依赖（患者端 Profile）

**标题**：demo-patient.yml 服务清单

**要点**
- 基础设施：PostgreSQL · Redis（无 RabbitMQ）
- 核心服务：gateway · auth · user · registration · schedule
- 患者特色：outpatient（就诊助手）· examination · followup · health_record
- 静态托管：Nginx 同时提供 `/patient/` 与 `/api/` 反代

**仓库引用**：`deploy/compose/demo-patient.yml`

---

## 第 7 页｜演示环境与账号

**标题**：如何启动患者端演示

**要点**
```bash
docker compose -f deploy/compose/demo-patient.yml --env-file deploy/config/demo.env up -d --build
cd frontend/patient && npm ci && npm run build
# 微信开发者工具导入 frontend/mp-webview
```
- 账号：`demo-patient` / `demo123`
- 绑定患者：王小明（patient_001）
- ⚠️ 勾选「不校验合法域名」

**仓库引用**：`docs/快速启动-患者端.md`

---

## 第 8 页｜核心亮点：就诊助手架构

**标题**：RAG + DeepSeek 智能分诊

**要点**
- 知识库：16 类病症（`backend/data/triage/knowledge.json`）
- 嵌入：SiliconFlow BGE-M3 语义检索（`embeddings.json`）
- 召回：关键词 + 语义双路 Top-K 合并
- 生成：DeepSeek Chat + 本院科室交集（user 服务 departments API）
- 降级：无 API Key 时关键词模式仍可用

**图示建议**：用户输入 → 检索 → 科室匹配 → LLM → 建议输出

**仓库引用**：`backend/internal/outpatient/assistant/`、`docs/待办清单.md` TA 章节

---

## 第 9 页｜【现场演示】小程序加载与登录

**标题**：现场演示 ① — 进入患者端

**演示动作**
| 步骤 | 环境 | 操作 |
|------|------|------|
| 1 | 微信开发者工具 | 展示 mp-webview 项目结构 |
| 2 | 模拟器 | 编译运行，显示登录页 |
| 3 | 登录 | demo-patient / demo123 |
| 4 | 首页 | 确认进入患者服务中心 |

**讲解要点**：web-view 加载 Nginx 托管的 H5；同一套后端 Gateway。

**备用**：浏览器打开 http://localhost/patient

---

## 第 10 页｜【现场演示】就诊助手

**标题**：现场演示 ② — AI 就诊助手（重点）

**演示动作**
| 步骤 | 页面 | 操作 |
|------|------|------|
| 1 | 助手 Tab | 输入「头痛发热三天，伴有咳嗽」 |
| 2 | 聊天页 | 等待 AI 回复 |
| 3 | 回复内容 | 指出推荐科室、就医建议、免责声明 |

**讲解要点**
- mode=`llm` 完整 RAG；mode=`keyword` 为降级
- 免责声明：辅助分诊，不替代医生

**仓库引用**：`frontend/patient/src/views/triage/TriageChatView.vue`

---

## 第 11 页｜【现场演示】预约与处方

**标题**：现场演示 ③ — 自助就医

**演示动作**
| 步骤 | Tab | 操作 |
|------|-----|------|
| 1 | 挂号 | 浏览今日排班（内科/张医生） |
| 2 | 处方 | 查看处方列表与状态 |
| 3 | （可选） | 若管理端刚演示过开方，指新处方 |

**讲解要点**：与管理端共用 schedule/prescription 服务，数据一致。

---

## 第 12 页｜【现场演示】报告与慢病

**标题**：现场演示 ④ — 院外健康管理

**演示动作**
| 步骤 | 入口 | 操作 |
|------|------|------|
| 1 | 首页网格 | 检查报告 → 展示 seed 报告 |
| 2 | 首页网格 | 慢病管理 → 郑十/刘十二签约 |
| 3 | 档案 Tab | 健康档案时间轴 |

**讲解要点**：覆盖诊后查结果、慢病长期管理场景。

**仓库引用**：`backend/sql/seed_data.sql` 慢病与报告数据

---

## 第 13 页｜与管理端数据联动

**标题**：内外双环 — 同一后端

**要点**
- 管理端开方 → 患者端「我的处方」可见
- 管理端排班 → 患者端「预约挂号」可见
- 统一 Gateway :8080 · 统一 JWT 体系
- 联演示建议：管理端 8 分钟 + 患者端 8 分钟

**图示建议**：管理端操作箭头 → 共享 DB/服务 → 患者端展示

**仓库引用**：`演示/管理端/演示流程-8分钟精简.md`

---

## 第 14 页｜真机预览与部署说明

**标题**：从模拟器到真机

**要点**
- 模拟器：开发者工具本地编译（推荐答辩）
- 真机：预览扫码，H5 地址改局域网 IP（`192.168.x.x/patient/`）
- 上线需配置：微信业务域名 HTTPS、合法 request 域名
- 当前范围：本地/局域网演示，不上云

**仓库引用**：`docs/快速启动-患者端.md` 真机预览章节

---

## 第 15 页｜已知限制

**标题**：诚实口径 — 患者端待完善项

**要点**
- 底部 Tab 未覆盖报告/随访（需首页网格进入）
- 部分路由网关未注册（followup/health-record list）
- `DEMO_PATIENT_MAP` 仅精确映射 demo-patient 账号
- 就诊助手依赖外部 API Key，需降级预案

**仓库引用**：`docs/待办清单.md` UX-6、UX-7

---

## 第 16 页｜总结

**标题**：患者端价值总结

**要点**
- ✅ 小程序 + H5 混合架构，一套业务双端复用
- ✅ 就诊助手 RAG 分诊，差异化亮点
- ✅ 预约、处方、报告、慢病、随访覆盖院外场景
- ⚠️ 演示版，上线需域名合规与安全加固
- 🔗 github.com/Tangyd893/HIS-Go

**演讲备注**：邀请体验模拟器；可与管理端 PPT 组合为完整 HIS 汇报。

---

## 附录：8 分钟精简版页序

合并第 9–12 页为 **1 页「患者端演示速查」**：

| 时间 | 动作 |
|------|------|
| 0:30 | 小程序编译 + 登录 |
| 1:00 | 助手：输入症状 |
| 3:00 | 挂号 Tab |
| 4:00 | 处方 Tab |
| 5:00 | 首页 → 检查报告 |
| 6:00 | 档案 Tab |
| 7:00 | 总结 |

配套：`演示流程.md` 文末 8 分钟精简路径

---

## 附录：与管理端 PPT 组合建议

若生成 **完整 HIS 汇报**（30 页），建议结构：

1. 封面 + 背景（2 页，共用）
2. 总体架构（2 页，共用）
3. **管理端专题**（8 页）← 使用 `../管理端/PPT生成提示词.md`
4. **患者端专题**（8 页）← 本文件
5. 质量保障 + 总结（2 页，共用）

---

## 附录：生成工具参数建议

```
语言：简体中文
页数：16
风格：移动端医疗 C 端
配图：手机 mockup 框 + 微信绿色点缀（小程序段落）
截图占位：登录页、助手聊天、挂号列表、处方列表
避免：过于冰冷的 B 端表格风
```
