# HIS-Go 演示资料索引

本目录存放 **Web 管理端** 与 **小程序患者端** 两套独立演示材料，供现场操作与后续 PPT 生成使用。

## 目录结构

```
演示/
├── README.md                          ← 本文件
├── 管理端/
│   ├── 演示流程-15分钟.md              完整版（含业务闭环讲解词）
│   ├── 演示流程-8分钟精简.md           时间紧张时使用
│   └── PPT生成提示词.md                交给 AI/PPT 工具生成答辩幻灯片
└── 患者端/
    ├── 演示流程.md                     微信小程序 + H5 患者端
    └── PPT生成提示词.md                患者端专题幻灯片提示词
```

## 使用方式

| 场景 | 使用文件 |
|------|----------|
| 答辩/汇报前彩排（管理端） | `管理端/演示流程-15分钟.md` 或 `8分钟精简.md` |
| 答辩/汇报前彩排（患者端） | `患者端/演示流程.md` |
| 用工具根据仓库生成 PPT | 将对应 `PPT生成提示词.md` 全文粘贴给生成工具，并附上 GitHub 仓库地址 |

## 仓库与入口

- **GitHub**：https://github.com/Tangyd893/HIS-Go
- **管理端**：http://localhost/admin（`demo-admin` / `demo123`）
- **患者 H5**：http://localhost/patient（`demo-patient` / `demo123`）
- **小程序壳**：`frontend/mp-webview`（微信开发者工具导入）

## 启动命令（演示前）

```bash
# 管理端全套
docker compose -f deploy/compose/demo-admin.yml --env-file deploy/config/demo.env up -d --build
cd frontend/admin && npm ci && npm run build

# 患者端（可与管理端共用基础设施，或单独 profile）
docker compose -f deploy/compose/demo-patient.yml --env-file deploy/config/demo.env up -d --build
cd frontend/patient && npm ci && npm run build
```

## 重要说明

- 演示密码 `demo123` **仅限本地演示**，不可用于生产。
- 管理端侧边栏 18 项菜单中，部分模块在 `demo-admin.yml` 未启动后端，演示时勿点 EMR/住院等。
- 患者端底部 Tab 仅 5 项；检查报告、随访等需从首页网格进入。
