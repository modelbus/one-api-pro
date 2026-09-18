---
title: Docker 单实例部署
description: "通过官方镜像运行单实例。"
category: install
order: 2
---

# Docker 单实例部署

> 通过官方镜像运行单实例。

## 拉取镜像

```bash
docker pull ghcr.io/modelbus/one-api-pro:latest
```

> 多架构镜像自动覆盖 `linux/amd64` 与 `linux/arm64`。

## 启动容器

```bash
docker run -d \
  --name one-api-pro \
  --restart unless-stopped \
  -p 3000:3000 \
  -v /opt/one-api-pro/config:/app/config \
  -v /opt/one-api-pro/data:/app/data \
  -e TZ=Asia/Shanghai \
  ghcr.io/modelbus/one-api-pro:latest
```

挂载说明 / Volume notes:

| 容器路径| 用途|
| --- | --- |
| `/app/config` | `.env` 配置文件目录，由 `docker-entrypoint.sh` 自动 `--env` 加载 |
| `/app/data` | SQLite 数据库、日志、上传缓存；`SQLITE_PATH` 默认指向此目录 |

入口脚本默认行为 / Entry-point behaviour:

- 若 `$CONFIG_DIR/.env`（默认 `/app/config/.env`）存在，则以 `--env` 形式加载；
- 任何用户传入的 CLI 参数都会被透传给二进制；
- 健康检查：`wget http://localhost:3000/api/status`（容器内置 `wget`，镜像已声明 `HEALTHCHECK`）。

## 自定义端口

```bash
docker run -d \
  --name one-api-pro \
  -p 8080:8080 \
  -v /opt/one-api-pro/config:/app/config \
  -v /opt/one-api-pro/data:/app/data \
  -e PORT=8080 \
  ghcr.io/modelbus/one-api-pro:latest
```

`PORT` 环境变量会被二进制与 `HEALTHCHECK` 同时读取。

## 接入外部数据库

```bash
docker run -d \
  --name one-api-pro \
  -p 3000:3000 \
  -v /opt/one-api-pro/config:/app/config \
  -v /opt/one-api-pro/data:/app/data \
  -e SQL_DSN='root:secret@tcp(mysql:3306)/oneapi?charset=utf8mb4&parseTime=True&loc=Local' \
  -e LOG_SQL_DSN='root:secret@tcp(mysql:3306)/oneapi_logs?charset=utf8mb4&parseTime=True&loc=Local' \
  -e REDIS_CONN_STRING='redis://default:redispw@redis:6379/0' \
  -e SESSION_SECRET='please-change-me' \
  ghcr.io/modelbus/one-api-pro:latest
```

数据库需预先创建空库 `oneapi` 与 `oneapi_logs`，表由 `AutoMigrate` 自动建好。

## 升级与回滚

```bash
# 拉取新版本
docker pull ghcr.io/modelbus/one-api-pro:latest
# 重建容器（数据卷保持不变）
docker stop one-api-pro && docker rm one-api-pro
docker run -d --name one-api-pro \
  -p 3000:3000 \
  -v /opt/one-api-pro/config:/app/config \
  -v /opt/one-api-pro/data:/app/data \
  ghcr.io/modelbus/one-api-pro:latest
```

数据卷与配置卷复用即可保留 SQLite / `.env` 内容；详细流程见 [版本升级](/zh/install/upgrade)。

下一步 / Next: [docker-compose 部署](/zh/install/docker-compose) · [反向代理](/zh/install/reverse-proxy)。

