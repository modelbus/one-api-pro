---
title: 系统要求
description: "硬件、操作系统、端口、数据库与依赖说明。"
category: install
order: 1
---

# 系统要求

> 硬件、操作系统、端口、数据库与依赖说明。
> Hardware, OS, ports, supported databases and dependencies.

## 硬件与操作系统 / Hardware & OS

One API Pro 是单一 Go 二进制，自带内嵌 SQLite 与 Redis 客户端，无需额外系统组件即可运行。

One API Pro is a single Go binary that bundles an embedded SQLite driver and Redis client; no extra system services are required to start.

| 资源 / Resource | 最低 / Minimum | 推荐 / Recommended |
| --- | --- | --- |
| CPU | 1 vCPU | 2 vCPU 及以上 / 2+ vCPU |
| 内存 / RAM | 512 MB | 1 GB 及以上 / 1 GB+ |
| 磁盘 / Disk | 1 GB | 10 GB+（日志与缓存按需扩容 / leave headroom for logs & cache） |
| OS | Linux、macOS、Windows | Linux（生产建议 / recommended for production） |

Docker 镜像基于 `alpine:latest`；二进制可在 Linux（amd64 / arm64）、macOS（amd64 / arm64）、Windows（amd64）上原生运行。

The Docker image is based on `alpine:latest`; native binaries are published for Linux (amd64 / arm64), macOS (amd64 / arm64) and Windows (amd64).

## 端口 / Ports

默认监听 **`3000`**（环境变量 `PORT`，CLI 参数 `--port` 可覆盖）。
容器内已 `EXPOSE 3000`；通过反向代理对外暴露 443 时，记得开启 WebSocket 升级以支持流式接口。

By default the service listens on **`3000`** (override via `PORT` or `--port`).
The container `EXPOSE`s `3000`; when exposing 443 through a reverse proxy, remember to enable WebSocket upgrade for streaming endpoints.

## 数据库 / Databases

主库 `SQL_DSN` 不设置时使用 **SQLite**（默认 `one-api-pro.db`，容器内 `/app/data/one-api-pro.db`）；设置后自动切换为：

When `SQL_DSN` is unset, the service uses **SQLite** (default `one-api-pro.db`; in the container `/app/data/one-api-pro.db`). Once set, it switches to:

- **MySQL** — `SQL_DSN=root:123456@tcp(localhost:3306)/oneapi`，需要预先创建空库 `oneapi`，表由 `AutoMigrate` 自动建好。
- **PostgreSQL** — `SQL_DSN=postgres://postgres:123456@localhost:5432/oneapi`。

`LOG_SQL_DSN` 可独立指定 `logs` 表所用的数据库（推荐 MySQL / PostgreSQL，因为 `logs` 表写入频繁，SQLite 在并发下易触发 `SQLITE_BUSY`）。

`LOG_SQL_DSN` (optional) splits the `logs` table onto a dedicated database; MySQL / PostgreSQL are recommended here because the `logs` table is write-heavy and SQLite can hit `SQLITE_BUSY` under concurrency.

集群模式下要求每个节点使用各自独立的 MySQL 实例（受 `auto_increment_offset` 实例级特性约束，详见 `decentralization/deployment`）。

In cluster mode every node must run its **own** MySQL instance (see `decentralization/deployment` for the rationale around `auto_increment_offset`).

## Redis（可选 / Optional）

不启用 Redis 服务即可运行（走进程内缓存 + 直读 DB）。启用 `REDIS_CONN_STRING` 后会作为分布式缓存层生效，可显著降低 DB 读压力；当 DB 延迟本身很低时，启用 Redis 可能带来短期数据陈旧，请按场景权衡。

The service runs without Redis (in-process cache + direct DB read). Enabling `REDIS_CONN_STRING` activates a distributed cache layer that can offload DB reads; if the DB is already low-latency, enabling Redis may introduce short-term staleness — choose based on your trade-off.

支持单节点、哨兵（`REDIS_MASTER_NAME`）与 Cluster 模式（节点列表逗号分隔 + `REDIS_PASSWORD`）。

Supports single-node, Sentinel (set `REDIS_MASTER_NAME`) and Cluster mode (comma-separated nodes + `REDIS_PASSWORD`).

## 网络 / Network

自 v0.0.21 起，tiktoken 通用模型的 BPE 编码已 **内嵌** 进二进制，启动时默认不再访问网络；只有在用户通过 `TIKTOKEN_CACHE_DIR` 显式提供自定义编码、且 URL 不在内嵌清单中时，才会触发 HTTP 下载兜底。

Since v0.0.21 the tiktoken encodings for common models are **embedded** into the binary; startup no longer requires network access by default. An HTTP download fallback is only triggered when the user supplies a custom encoding via `TIKTOKEN_CACHE_DIR` whose URL is not in the embedded allow-list.

支付回调、对外代理等出站流量仍需正常出网；离线部署时评估白名单。

Outbound traffic for payment callbacks and upstream relays still requires network; whitelist accordingly for air-gapped deployments.

## 开发依赖 / Dev Dependencies

构建 Go 后端与前端主题所需工具，与 `AGENTS.md §2` 保持一致 / Toolchain required to build the Go backend and web theme, matching `AGENTS.md §2`:

| 工具 / Tool | 版本 / Version |
| --- | --- |
| Go | 1.25.0（`go.mod` 中声明，toolchain 自动拉取） |
| Node.js | 22+ |
| pnpm | 9（仓库提交 `pnpm-lock.yaml`） |

下一步 / Next: [Docker 单实例部署](/zh/install/docker-deploy) · [Docker Compose](/zh/install/docker-compose) · [源码编译](/zh/install/source-build)。
