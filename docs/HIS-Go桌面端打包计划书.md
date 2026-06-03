# HIS-Go Windows 桌面端打包计划书

> 版本：v2.0  
> 日期：2026-05-15  
> 状态：规划中  
> 原则：**不影响已有 Web 版本任何代码与部署流程**

---

## 目录

1. [需求定位](#1-需求定位)
2. [可行方案对比](#2-可行方案对比)
3. [推荐方案：Tauri v2 纯壳](#3-推荐方案tauri-v2-纯壳)
4. [架构设计](#4-架构设计)
5. [实施计划](#5-实施计划)
6. [风险与对策](#6-风险与对策)
7. [交付物清单](#7-交付物清单)

---

## 1. 需求定位

### 1.1 目标

将 HIS-Go **管理端**（`his-web-admin-react`）打包为 Windows 原生桌面应用，**仅作为瘦客户端**运行：

- **只服务管理端**，不含患者端
- **所有业务逻辑由医院内网服务器提供**，桌面端不内置任何后端服务
- **零本地数据库**，不内置 SQLite、不内置认证
- **通过 HTTP/HTTPS 连接内网 Gateway**（如 `http://192.168.1.100:8080`）
- **无需用户安装 Node.js / Docker / 数据库**
- 系统托盘常驻 + 开机自启（可选）

### 1.2 约束

| 约束项 | 说明 |
|--------|------|
| 不侵入现有代码 | `backend/`、`frontend/`、`docker/` 零改动 |
| 复用现有 React 前端 | 直接引用 `his-web-admin-react/` 构建产物 |
| 不内置业务后端 | 桌面端无 Go 逻辑、无数据库、无 MQ |
| Windows 10+ 支持 | 依赖系统自带 WebView2 |

### 1.3 运行模型

```
┌──────────────────────┐         HTTP/HTTPS         ┌──────────────────────┐
│   桌面端 (瘦客户端)     │ ──────────────────────────→ │  医院内网服务器         │
│                       │                            │                      │
│  WebView2             │    POST /api/auth/login    │  Gateway :8080       │
│  (React 管理端)        │ ←────────────────────────── │  ↓                   │
│                       │       JWT Token            │  18 微服务            │
│  无本地数据库           │                            │  PostgreSQL / Redis  │
│  无业务逻辑             │    GET /api/user/patients  │  RabbitMQ            │
│                       │ ←────────────────────────── │                      │
└──────────────────────┘                            └──────────────────────┘
```

---

## 2. 可行方案对比

需求本质：**一个带窗口外壳的浏览器**，只访问内网地址。

| 方案 | 体积 | 开发量 | 系统托盘 | 自动更新 | 推荐 |
|------|------|--------|---------|---------|------|
| **Tauri v2** | ~3 MB | 低 | ✅ 内置 | ✅ 内置 | ★★★★★ |
| Electron | ~120 MB | 极低 | ✅ | ✅ | ★★★☆ |
| WebView2 直装 | ~2 MB | 中 | 需自建 | 需自建 | ★★★★☆ |
| PWA (Edge/Chrome) | 0 MB | 零 | ❌ | ❌ | ★★☆☆ |

### 逐项分析

#### Tauri v2（推荐）

- 纯 Rust 壳，后端代码 < 50 行（仅窗口配置 + 加载 URL）
- 体积 ~3 MB（不含前端静态资源）
- 内置系统托盘、自动更新、窗口定制
- `tauri.conf.json` 中配置 `"url": "http://192.168.1.100:8080/admin"` 即可
- 前端静态资源可内置（离线打开更快），也可完全从服务器加载

#### Electron

- 体积 120 MB+，对瘦客户端来说过重
- 优势是零学习成本，`main.js` 十几行即可
- 适合快速验证，不适合长期分发

#### WebView2 直接封装

- 用 C# / Rust 调用 Windows WebView2 API，手动创建窗口
- 体积最小，但需自行实现托盘/更新/安装包
- 适合有 Windows 原生开发经验的团队

#### PWA

- 浏览器地址栏输入即可，无需打包
- 但没有系统托盘、没有 EXE 图标、用户感知弱

### 结论

**选择 Tauri v2**。理由：
1. 体积 3 MB，分发轻量
2. 内置系统托盘、自动更新，无需从零造轮子
3. Rust 后端仅 30–50 行代码（窗口 + 加载 URL），无业务逻辑
4. 可内嵌前端静态资源作为启动页，再代理到内网服务器

---

## 3. 推荐方案：Tauri v2 纯壳

### 3.1 技术栈

| 层级 | 技术 | 说明 |
|------|------|------|
| 桌面外壳 | Tauri v2.x | Rust 框架，仅管理窗口生命周期 |
| 渲染引擎 | WebView2 | Windows 10/11 自带 |
| 前端 UI | React 19 + Semantic UI | 复用 `his-web-admin-react` **构建产物** |
| 业务后端 | **无** | 全部由内网 HIS 服务器提供 |
| 数据库 | **无** | 全部由内网 PostgreSQL 提供 |
| 打包 | Tauri CLI + NSIS | 生成 `.exe` / `.msi` |

### 3.2 Tauri 后端代码（全部）

```rust
// desktop/src-tauri/src/main.rs  —— 完整代码，仅此而已
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

fn main() {
    tauri::Builder::default()
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
```

```json
// desktop/src-tauri/tauri.conf.json  —— 关键配置
{
  "app": {
    "windows": [{
      "title": "HIS-Go 医院管理系统",
      "width": 1280,
      "height": 800,
      "url": "http://192.168.1.100:8080/admin"
    }]
  }
}
```

> 无需 Go 代码、无需 SQLite、无需 JWT 签发。Rust 代码仅 9 行。

### 3.3 新增目录结构

```
HIS-Go/
├── desktop/                            # ★ 新增
│   ├── src-tauri/                      # Tauri Rust 后端
│   │   ├── Cargo.toml                  # Rust 依赖（仅 tauri）
│   │   ├── tauri.conf.json             # 窗口/URL/图标配置
│   │   ├── icons/                      # 应用图标（自动生成）
│   │   └── src/
│   │       └── main.rs                 # 入口（< 10 行）
│   ├── src/                            # 前端壳（可指向 React dist）
│   │   └── index.html                  # 内嵌首页或重定向
│   ├── build/                          # 打包脚本
│   │   └── installer.nsi               # NSIS 安装包
│   ├── package.json                    # 仅含 tauri-cli 依赖
│   └── README.md
├── frontend/
│   └── his-web-admin-react/            # 现有管理端（不动）
├── backend/                            # 现有后端（不动）
└── ...
```

### 3.4 两种加载模式

| 模式 | `tauri.conf.json` 配置 | 适用场景 |
|------|----------------------|---------|
| **纯远程** | `"url": "http://192.168.1.100:8080/admin"` | 客户端零部署，服务器更新即时生效 |
| **本地壳 + 远程 API** | 内嵌 `dist/`，前端代码本地加载，API 走远程代理 | 启动更快，离线显示错误页 |

**推荐纯远程模式**：桌面端 EXE 只有 3 MB，所有前端更新在服务器端一次完成。

---

## 4. 实施计划

### 第一阶段：环境准备（0.5 天）

| 任务 | 产出 |
|------|------|
| 1.1 安装 Rust + Tauri CLI | `cargo install tauri-cli` |
| 1.2 安装 WebView2 Runtime（Win10 已自带） | 验证 `webview2` 可用 |
| 1.3 创建 Tauri 骨架项目 | `desktop/` 目录结构 |

### 第二阶段：壳开发（1 天）

| 任务 | 产出 |
|------|------|
| 2.1 配置窗口标题/尺寸/图标 | `tauri.conf.json` |
| 2.2 配置加载 URL（内网服务器地址） | 窗口显示管理端登录页 |
| 2.3 系统托盘 + 右键菜单 | 最小化到托盘、退出 |
| 2.4 窗口状态记忆（位置/尺寸） | 下次启动恢复 |
| 2.5 内网地址配置文件（`config.json`） | 用户可修改服务器地址 |

### 第三阶段：分发增强（1 天）

| 任务 | 产出 |
|------|------|
| 3.1 NSIS 安装包脚本 | `HIS-Go-Setup-1.0.0.exe` |
| 3.2 应用图标设计（256×256 ICO） | 桌面/任务栏图标 |
| 3.3 自动更新（Tauri updater） | 版本检测 + 增量更新 |
| 3.4 开机自启（Tauri 内置） | 注册表写入 |

### 第四阶段：测试与发布（0.5 天）

| 任务 | 产出 |
|------|------|
| 4.1 多版本 Windows 兼容测试 | Win10 / Win11 均可运行 |
| 4.2 CI/CD（GitHub Actions） | Release 自动构建 `.exe` |
| 4.3 数字签名（可选） | 避免 SmartScreen 拦截 |

> **总工期：3 天**（极简，因无业务后端开发）

---

## 5. 风险与对策

| 风险 | 概率 | 影响 | 对策 |
|------|------|------|------|
| WebView2 未安装（Win10 1803 之前） | 极低 | 高 | 安装包内置 Evergreen Bootstrapper（~2 MB），自动安装 |
| 内网地址变更 | 中 | 中 | 提供 `config.json` 让用户修改；或启动时弹出地址配置框 |
| 服务器宕机时白屏 | 中 | 低 | 本地内嵌错误页："无法连接服务器，请检查网络" |
| Tauri 安全策略阻止跨域请求 | 低 | 中 | 配置 CSP 允许内网地址；API 走 HTTP 不受 WebView 跨域限制 |
| 杀毒软件误报 | 低 | 低 | 数字签名后可消除 |

---

## 6. 交付物清单

| 文件 | 说明 |
|------|------|
| `desktop/src-tauri/src/main.rs` | Tauri 入口（< 10 行 Rust） |
| `desktop/src-tauri/tauri.conf.json` | 窗口/URL/图标配置 |
| `desktop/src-tauri/Cargo.toml` | Rust 依赖 |
| `desktop/package.json` | Tauri CLI + 构建脚本 |
| `desktop/build/installer.nsi` | NSIS 安装包脚本 |
| `docs/HIS-Go桌面端打包计划书.md` | 本文件 |
| GitHub Release: `HIS-Go-Setup-1.0.0.exe` | Windows 安装包（~5 MB） |

---

## 附录 A：Tauri 快速上手命令

```bash
# 安装 Rust（如未安装）
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# 安装 Tauri CLI
cargo install tauri-cli

# 创建项目
cd HIS-Go/desktop
npm create tauri-app@latest . -- --template vanilla
# 选择: TypeScript / 否(不需要前端框架，直接用 WebView 加载远程 URL)

# 开发
cargo tauri dev

# 构建 Windows 安装包
cargo tauri build
```

## 附录 B：与 Web 版关系

| 项目 | Web 版本（现有） | 桌面版本（计划） |
|------|-----------------|------------------|
| 目标用户 | 浏览器访问 | 双击 EXE 启动 |
| 管理端代码 | `his-web-admin-react/` | **完全复用**，由服务器提供 |
| 患者端 | `his-web-patient-react/` | 不包含 |
| 后端 | 18 微服务（Docker） | 无后端，纯客户端 |
| 数据库 | PostgreSQL 17 | 无数据库 |
| 认证 | JWT + Gateway | 同上，走 HTTP |
| 部署 | Docker Compose | `.exe`（3 MB） + 内网服务器 |
| 维护 | 服务器端更新 | 前端更新仅在服务器端 |

---

> **本方案遵循"零侵入现有代码"原则。**  
> 桌面端作为独立子项目 `desktop/` 存在，Tauri 壳仅 30 行 Rust 代码，  
> 所有业务能力由医院内网服务器提供。前端更新一次，所有桌面端即时生效。
