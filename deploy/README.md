# HIS-Go 部署目录

演示与生产相关的 Docker 配置集中在此，**不再使用根目录 `docker/`**（避免与工具链混淆，且 `.env` 不再放在易被误加载的路径）。

## 目录说明

| 路径 | 用途 |
|------|------|
| `compose/demo-admin.yml` | 管理端演示：PG + Redis + RabbitMQ + 9 微服务 + Gateway + Nginx |
| `compose/demo-patient.yml` | 患者端演示：PG + Redis + 8 微服务 + Gateway + Nginx（无 MQ） |
| `compose/stack.yml` | 全量 18 服务（非演示默认） |
| `compose/stack.prod.yml` | 生产向全量栈 |
| `config/demo.env.example` | 演示环境变量模板 → 复制为 `demo.env` |
| `config/stack.env.example` | 全量栈环境变量模板 |
| `nginx/nginx.conf` | 静态 `/admin`、`/patient` + `/api/` 反代 |
| `Dockerfile` | Go 微服务多阶段构建 |
| `init-db.sh` | 数据库初始化辅助脚本 |

## 快速命令

```bash
cp config/demo.env.example config/demo.env   # 在 deploy/ 目录下或从项目根指定路径

# 从项目根
make demo-admin
make demo-patient
```

详细步骤见 **`docs/云端部署指南.md`**。

## 环境变量说明

- 使用 **`deploy/config/demo.env`**，通过 `--env-file` 显式传入，**不要**在 `deploy/config/` 下创建名为 `.env` 的文件（Docker Compose 可能自动加载导致意外覆盖）。
- 切勿在仓库根目录创建名为 `%USERPROFILE%` 等含 `%` 的文件夹（Windows 下会与用户环境变量冲突）。

## PostgreSQL 端口说明

> **容器内端口 5432，宿主机映射为 5433**（避免与 Windows 本机 PostgreSQL 冲突）。

- 容器间通信：各微服务通过 Docker 网络 `his-net` 连接 `postgresql:5432`
- 宿主机直连：使用 `localhost:5433`（如 `go run` 本地开发时）
- `deploy/compose/*.yml` 中统一配置为 `5433:5432`
- `backend/configs/config.yaml` 默认端口已改为 5433
