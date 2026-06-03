# HIS 患者端 — 微信开发者工具演示壳

本目录**不是**重写业务的小程序，仅用于在**微信开发者工具**里演示已有的 `his-web-patient-react` H5。

> **定位**：本子项目只服务于患者端（`his-web-patient-react`），不包含管理端（`his-web-admin-react`）。  
> **职责**：通过 `<web-view>` 加载已运行的患者 H5；不实现挂号、登录等业务逻辑。

## 使用步骤

### 1. 前置条件（必须先完成）

```bash
# 1.1 启动基础设施
cd docker && docker compose up -d

# 1.2 初始化数据库
cd ../backend && bash scripts/db_init.sh

# 1.3 启动后端服务（至少 Gateway）
go run ./cmd/gateway

# 1.4 启动患者 H5
cd frontend/his-web-patient-react
npm install --legacy-peer-deps
npm run dev -- --host    # → http://127.0.0.1:5174
```

### 2. 导入微信开发者工具

1. 下载并安装 [微信开发者工具](https://developers.weixin.qq.com/miniprogram/dev/devtools/download.html)
2. 打开工具 → **导入项目** → 选择本目录 `frontend/his-mp-webview`
3. 进入 **详情 → 本地设置**，勾选：
   - ✅ 不校验合法域名、web-view（业务域名）、TLS 版本以及 HTTPS 证书
4. 编译运行，模拟器中应显示患者端页面

### 3. 真机预览（可选）

1. 修改 `pages/index/index.js` 中 `DEFAULT_PATIENT_H5` 为 PC 局域网 IP：
   ```javascript
   const DEFAULT_PATIENT_H5 = 'http://192.168.x.x:5174'
   ```
2. 患者 H5 以 host 模式启动：`npm run dev -- --host`
3. 手机与 PC 连接同一 WiFi
4. 微信开发者工具 → **预览** → 手机扫码

## 文件说明

| 文件 | 说明 |
|------|------|
| `app.json` | 小程序配置，仅一个页面 |
| `pages/index/index.js` | 页面逻辑，配置 H5 地址 |
| `pages/index/index.wxml` | 页面模板，包含 web-view 和加载/错误提示 |
| `pages/index/index.wxss` | 页面样式 |
| `project.config.json` | 项目配置，`urlCheck: false` |

## 常见问题

| 现象 | 原因 | 解决 |
|------|------|------|
| 白屏/无法加载 | H5 未启动 | 确认 `npm run dev` 已运行，`http://127.0.0.1:5174` 可访问 |
| 提示无法打开 web-view | 未勾选不校验 | 详情 → 本地设置 → 勾选相关选项 |
| 真机白屏 | 使用了 127.0.0.1 | 改为 PC 局域网 IP |
| 接口报错 | 后端未启动 | 确认 Gateway:8080 及相关服务已启动 |

## 说明

- 演示阶段无需小程序上架、无需配置业务域名（依赖本地「不校验」选项）
- 业务逻辑仍在 `his-web-patient-react`；本壳仅一个 `web-view` 页面
- `touristappid` 用于本机演示；真机受限时需改为测试号/正式 AppID
