---
title: docker-compose 部署
description: 通过 docker-compose 编排 One API Pro + 数据库，适合生产环境。
category: install
order: 3
---

# docker-compose 部署

> 适合生产、单节点有 MySQL/Redis 的标准部署。

## 仅 One API Pro

最简编排，只用容器内嵌的 SQLite。

`docker-compose.yml`：

```yaml
services:
  one-api-pro:
    image: ghcr.io/modelbus/one-api-pro:latest
    container_name: one-api-pro
    restart: always
    ports:
      - "3000:3000"
    volumes:
      - ./data:/app/data
    environment:
      TZ: Asia/Shanghai
```

启动：

```bash
docker compose up -d
```

## One API Pro + MySQL

推荐生产配置。

`docker-compose.yml`：

```yaml
services:
  one-api-pro:
    image: ghcr.io/modelbus/one-api-pro:latest
    container_name: one-api-pro
    restart: always
    ports:
      - "3000:3000"
    volumes:
      - ./data:/app/data
    environment:
      TZ: Asia/Shanghai
      SQL_DSN: "oneapi:oneapi-pass@tcp(mysql:3306)/oneapi?charset=utf8mb4&parseTime=True&loc=Local"
    depends_on:
      mysql:
        condition: service_healthy

  mysql:
    image: mysql:8.0
    container_name: one-api-pro-mysql
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: root-pass
      MYSQL_DATABASE: oneapi
      MYSQL_USER: oneapi
      MYSQL_PASSWORD: oneapi-pass
    command:
      - --character-set-server=utf8mb4
      - --collation-server=utf8mb4_unicode_ci
    volumes:
      - ./mysql-data:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 5s
      timeout: 3s
      retries: 20
```

启动：

```bash
docker compose up -d
```

数据库自动初始化。

## One API Pro + MySQL + Redis（生产推荐）

加一个 Redis 服务，并打开 Redis 依赖：

```yaml
services:
  one-api-pro:
    # ... 同上 ...
    environment:
      SQL_DSN: "oneapi:oneapi-pass@tcp(mysql:3306)/oneapi?charset=utf8mb4&parseTime=True&loc=Local"
      REDIS_CONN_STRING: "redis://redis:6379/0"
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy

  redis:
    image: redis:7-alpine
    container_name: one-api-pro-redis
    restart: always
    volumes:
      - ./redis-data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 10
```

## 关键环境变量

| 变量 | 说明 |
|---|---|
| `TZ` | 时区 |
| `SQL_DSN` | MySQL 连接串；留空用 SQLite |
| `REDIS_CONN_STRING` | Redis 连接串；留空不用 Redis |
| `CLUSTER_NODE_ID` / `CLUSTER_NODE_SECRET` / `CLUSTER_NODE_PORT` | Cluster 模式（多节点时需要） |
| 其他高级项 | 见 [配置项](./config) |

## 数据持久化

通过 `./data`、`./mysql-data`、`./redis-data` 三个挂载目录持久化。**升级或重启容器数据不丢**。

## 升级版本

```bash
docker compose pull one-api-pro
docker compose up -d
```

只重启 one-api-pro 容器，其他服务不受影响。

## 备份

```bash
# 停服后备份（确保一致性）
docker compose stop one-api-pro
cp -r ./mysql-data ./backup-$(date +%Y%m%d)
docker compose start one-api-pro
```

更详细的备份恢复策略见 [备份与恢复](./backup-restore)。

## 相关文档

- [Docker 单实例部署](./docker-deploy)
- [配置项与环境变量](./config)
- [备份与恢复](./backup-restore)
- [反向代理](./reverse-proxy)