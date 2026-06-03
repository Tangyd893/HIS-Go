# HIS-Go 小程序壳集成实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 从上游仓库引入微信小程序壳 `his-mp-webview`，配置患者端 H5 端口为 5174，并完成本地联调验证。

**Architecture:** 复用现有小程序壳（仅 web-view 容器），通过修改患者端 Vite 配置解决端口冲突，实现小程序 → H5 → Gateway 的完整链路。

**Tech Stack:** 微信小程序、React、Vite、Gateway

---

## 文件结构

### 新增文件
- `frontend/his-mp-webview/app.js` - 小程序入口
- `frontend/his-mp-webview/app.json` - 小程序配置
- `frontend/his-mp-webview/app.wxss` - 全局样式
- `frontend/his-mp-webview/project.config.json` - 项目配置
- `frontend/his-mp-webview/pages/index/index.js` - 页面逻辑
- `frontend/his-mp-webview/pages/index/index.wxml` - 页面模板
- `frontend/his-mp-webview/pages/index/index.wxss` - 页面样式
- `frontend/his-mp-webview/pages/index/index.json` - 页面配置
- `frontend/his-mp-webview/README.md` - 使用说明

### 修改文件
- `frontend/his-web-patient-react/vite.config.ts` - 修改端口为 5174，启用 host
- `README.md` - 更新患者端端口说明

---

## Task 1: 创建小程序壳目录结构

**Files:**
- Create: `frontend/his-mp-webview/`

- [ ] **Step 1: 创建目录结构**

```bash
mkdir -p frontend/his-mp-webview/pages/index
```

- [ ] **Step 2: 验证目录创建**

```bash
ls -la frontend/his-mp-webview/
ls -la frontend/his-mp-webview/pages/
```

Expected: 目录结构已创建

---

## Task 2: 创建小程序配置文件

**Files:**
- Create: `frontend/his-mp-webview/app.json`
- Create: `frontend/his-mp-webview/app.js`
- Create: `frontend/his-mp-webview/app.wxss`
- Create: `frontend/his-mp-webview/project.config.json`

- [ ] **Step 1: 创建 app.json**

```json
{
  "pages": ["pages/index/index"],
  "window": {
    "navigationBarTitleText": "HIS 患者端",
    "navigationBarBackgroundColor": "#1890ff",
    "navigationBarTextStyle": "white"
  }
}
```

- [ ] **Step 2: 创建 app.js**

```javascript
App({
  onLaunch() {},
})
```

- [ ] **Step 3: 创建 app.wxss**

```css
page {
  height: 100%;
}
```

- [ ] **Step 4: 创建 project.config.json**

```json
{
  "description": "HIS 患者端 H5 的微信开发者工具演示壳（web-view）",
  "packOptions": {
    "ignore": [],
    "include": []
  },
  "setting": {
    "urlCheck": false,
    "es6": true,
    "postcss": true,
    "minified": false
  },
  "compileType": "miniprogram",
  "appid": "touristappid",
  "projectname": "his-mp-webview",
  "miniprogramRoot": "./"
}
```

- [ ] **Step 5: 验证配置文件**

```bash
cat frontend/his-mp-webview/app.json
cat frontend/his-mp-webview/project.config.json
```

Expected: 文件内容正确

---

## Task 3: 创建小程序页面文件

**Files:**
- Create: `frontend/his-mp-webview/pages/index/index.js`
- Create: `frontend/his-mp-webview/pages/index/index.wxml`
- Create: `frontend/his-mp-webview/pages/index/index.wxss`
- Create: `frontend/his-mp-webview/pages/index/index.json`

- [ ] **Step 1: 创建 index.js（已修改 DEFAULT_PATIENT_H5 为 5174）**

```javascript
// 本地演示：先启动 frontend/his-web-patient-react（npm run dev，默认 5174）
// 微信开发者工具 → 详情 → 本地设置 → 勾选「不校验合法域名、web-view、TLS…」

// 环境区分：模拟器用 127.0.0.1，真机需改为 PC 局域网 IP
const DEFAULT_PATIENT_H5 = 'http://127.0.0.1:5174'

Page({
  data: {
    patientUrl: DEFAULT_PATIENT_H5,
    loading: true,
    loadError: false,
  },

  onLoad() {
    // 模拟加载延迟，避免白屏闪烁
    setTimeout(() => {
      this.setData({ loading: false })
    }, 500)
  },

  onWebviewError(e) {
    console.error('web-view 加载失败:', e)
    this.setData({
      loading: false,
      loadError: true,
    })
  },

  retryLoad() {
    this.setData({
      loading: true,
      loadError: false,
      patientUrl: '',
    })
    // 重新设置 URL 触发重新加载
    setTimeout(() => {
      this.setData({ patientUrl: DEFAULT_PATIENT_H5 })
    }, 100)
  },
})
```

- [ ] **Step 2: 创建 index.wxml**

```html
<view class="container">
  <!-- 加载中提示 -->
  <view class="loading-container" wx:if="{{loading}}">
    <view class="loading-icon">⏳</view>
    <view class="loading-text">正在连接患者服务...</view>
  </view>

  <!-- 加载失败提示 -->
  <view class="error-container" wx:elif="{{loadError}}">
    <view class="error-icon">❌</view>
    <view class="error-title">无法加载患者端页面</view>
    <view class="error-desc">请检查以下事项：</view>
    <view class="error-list">
      <view>1. 患者 H5 是否已启动（npm run dev）</view>
      <view>2. 后端 Gateway 是否运行（:8080）</view>
      <view>3. 微信开发者工具是否勾选「不校验合法域名」</view>
    </view>
    <button class="retry-btn" bindtap="retryLoad">重新加载</button>
  </view>

  <!-- web-view -->
  <web-view wx:else src="{{patientUrl}}" binderror="onWebviewError"></web-view>
</view>
```

- [ ] **Step 3: 创建 index.wxss**

```css
.container {
  min-height: 100vh;
  background: #f5f5f5;
}

.loading-container,
.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 40rpx;
}

.loading-icon,
.error-icon {
  font-size: 80rpx;
  margin-bottom: 30rpx;
}

.loading-text {
  font-size: 32rpx;
  color: #666;
}

.error-title {
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 20rpx;
}

.error-desc {
  font-size: 28rpx;
  color: #666;
  margin-bottom: 20rpx;
}

.error-list {
  font-size: 26rpx;
  color: #999;
  line-height: 1.8;
  text-align: left;
  width: 100%;
  padding: 0 40rpx;
  margin-bottom: 40rpx;
}

.retry-btn {
  background: #1890ff;
  color: #fff;
  border-radius: 8rpx;
  font-size: 30rpx;
  padding: 20rpx 60rpx;
}
```

- [ ] **Step 4: 创建 index.json**

```json
{
  "navigationBarTitleText": "HIS 患者端"
}
```

- [ ] **Step 5: 验证页面文件**

```bash
ls -la frontend/his-mp-webview/pages/index/
```

Expected: 4 个文件已创建

---

## Task 4: 创建小程序壳 README

**Files:**
- Create: `frontend/his-mp-webview/README.md`

- [ ] **Step 1: 创建 README.md**

```markdown
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
```

- [ ] **Step 2: 验证 README**

```bash
cat frontend/his-mp-webview/README.md
```

Expected: README 内容正确，包含 HIS-Go 启动步骤

---

## Task 5: 修改患者端 Vite 配置

**Files:**
- Modify: `frontend/his-web-patient-react/vite.config.ts`

- [ ] **Step 1: 读取当前配置**

```bash
cat frontend/his-web-patient-react/vite.config.ts
```

Expected: 当前端口为 5173

- [ ] **Step 2: 修改 vite.config.ts**

将 `server.port` 从 5173 改为 5174，添加 `host: true`：

```typescript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: { '@': '/src' },
  },
  css: {
    transformer: 'postcss',
  },
  build: {
    cssMinify: 'esbuild',
  },
  server: {
    port: 5174,
    host: true,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
```

- [ ] **Step 3: 验证修改**

```bash
cat frontend/his-web-patient-react/vite.config.ts
```

Expected: 端口已改为 5174，host 已启用

---

## Task 6: 更新 README.md 端口说明

**Files:**
- Modify: `README.md`

- [ ] **Step 1: 检查当前 README 中患者端端口说明**

```bash
grep -n "5174" README.md
```

Expected: 找到患者端端口说明

- [ ] **Step 2: 确认端口说明一致**

README.md 中患者端端口应为 5174，与 vite.config.ts 一致。

- [ ] **Step 3: 验证 README**

```bash
grep -A2 -B2 "患者端" README.md | grep -E "517[0-9]"
```

Expected: 患者端端口为 5174

---

## Task 7: 本地验证

**Files:**
- None（验证步骤）

- [ ] **Step 1: 启动后端 Gateway**

```bash
cd backend && go run ./cmd/gateway
```

Expected: Gateway 启动成功，监听 8080 端口

- [ ] **Step 2: 启动患者端 H5**

```bash
cd frontend/his-web-patient-react
npm run dev -- --host
```

Expected: 开发服务器启动，监听 5174 端口

- [ ] **Step 3: 验证 H5 可访问**

```bash
curl -I http://127.0.0.1:5174
```

Expected: 返回 200 OK

- [ ] **Step 4: 微信开发者工具验证**

1. 打开微信开发者工具
2. 导入项目 `frontend/his-mp-webview`
3. 勾选「不校验合法域名」
4. 编译运行

Expected: 模拟器中显示患者端登录/首页

- [ ] **Step 5: 记录验证结果**

在 `docs/待办清单.md` 中记录验证结果：
- [x] P0-1 引入 `frontend/his-mp-webview` ✅
- [x] P0-2 修改 `pages/index/index.js` 中 `DEFAULT_PATIENT_H5` 指向 `5174` ✅
- [x] P0-3 患者端 `vite.config.ts`：`port: 5174`、`host: true` ✅
- [x] P0-4 本地验证：Gateway 健康 → 患者 H5 可访问 → 微信开发者工具 web-view 正常显示 ✅
- [x] P0-5 文档：在本文档记录验证结果与 AppID 说明 ✅

---

## 执行顺序

1. Task 1 → Task 2 → Task 3 → Task 4（创建小程序壳）
2. Task 5（修改患者端配置）
3. Task 6（更新文档）
4. Task 7（验证）

---

## 验证清单

- [ ] `frontend/his-mp-webview` 目录结构完整（9 个文件）
- [ ] `pages/index/index.js` 中 `DEFAULT_PATIENT_H5` 为 `http://127.0.0.1:5174`
- [ ] `frontend/his-web-patient-react/vite.config.ts` 端口为 5174，host 为 true
- [ ] 微信开发者工具可导入项目并无报错
- [ ] 模拟器中 web-view 加载患者端页面
- [ ] README.md 中患者端端口说明与配置一致

---

## 风险与对策

| 风险 | 影响 | 对策 |
|------|------|------|
| 患者 H5 仍为 mock | 小程序能打开但业务不可用 | 后续 P1 患者端 API 联调 |
| 5173 端口冲突 | 管理端/患者端无法同时 dev | 已通过 P0-3 解决 |
| web-view 白屏 | 域名/证书/未勾选不校验 | 按 README 检查；真机改用局域网 IP |

---

## 下一步（P1 任务）

完成 P0 后，继续：
- P1-1 配置 `base: '/patient/'` + `BrowserRouter basename`
- P1-2 患者端 API 模块对齐接口规范
- P1-3 去除 mock 数据，接入真实接口
- P1-4 后端补充 `GET /api/auth/current`
- P1-5 小程序生产环境配置文档
