# HIS 患者端 — 微信开发者工具演示壳

本目录**不是**重写业务的小程序，仅用于在**微信开发者工具**里演示 **`frontend/patient`（Vue 3）** H5。

> **定位**：仅服务患者端，不包含管理端。  
> **职责**：通过 `<web-view>` 加载患者 H5；业务逻辑在 H5 + 后端 API。

## 使用步骤

### 1. 前置条件

```bash
# 启动患者端 Profile
make demo-patient
# 或本地 dev：
cd frontend/patient && npm run dev -- --host   # → http://127.0.0.1:5174
```

### 2. 导入微信开发者工具

1. 导入目录 **`frontend/mp-webview`**
2. **详情 → 本地设置** → 勾选「不校验合法域名、web-view…」
3. 编译运行

### 3. 云端 H5

编辑 `pages/index/index.js`：`USE_CLOUD = true`，`CLOUD_PATIENT_H5 = 'https://域名/patient/'`  
详见 **`docs/云端部署指南.md`**。

## 演示账号

`demo-patient` / `demo123`
