---
title: docker-compose 部署
description: "通过 docker-compose 编排单实例或带依赖的部署。"
category: install
order: 3
---

# docker-compose 部署

> 通过 docker-compose 编排单实例或带依赖的部署。
> Orchestrate a single instance or a stack with dependencies via docker-compose.

## 仅 One API Pro / One API Pro only

最简编排，只使用容器内嵌的 SQLite。

The minimal stack — relies on the embedded SQLite inside the container.

```yaml
# compose.yaml
services:
  one-api-pro:
    image: ghcr.io/modelbus/one-api-pro:latest
    container_name: one-api-pro
    restart: unless-stopped
    ports:
      - "3000:3000"
    environment:
      TZ: Asia/Shanghai
      SESSION_SECRET: please-change-me
    volumes:
      - ./config:/app/config
      - ./data:/app/data
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:3000/api/status"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 20s
```

启动 / Start:

```bash
docker compose up -d
```

健康检查复用镜像内 `HEALTHCHECK` 指令；`wget` 已在镜像中预装。

The `wget`-based healthcheck matches the image's built-in `HEALTHCHECK`.

## MySQL + Redis 全栈 / Full stack with MySQL & Redis

生产推荐组合：MySQL 持久化业务数据，Redis 做缓存与限流。

Production-ready stack: MySQL for persistence, Redis for caching and rate-limit state.

```yaml
# compose.yaml
services:
  one-api-pro:
    image: ghcr.io/modelbus/one-api-pro:latest
    container_name: one-api-pro
    restart: unless-stopped
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
    ports:
      - "3000:3000"
    environment:
      TZ: Asia/Shanghai
      SQL_DSN: "root:oneapi_pw@tcp(mysql:3306)/oneapi?charset=utf8mb4&parseTime=True&loc=Local"
      LOG_SQL_DSN: "root:oneapi_pw@tcp(mysql:3306)/oneapi_logs?charset=utf8mb4&parseTime=True&loc=Local"
      REDIS_CONN_STRING: "redis://default:redispw@redis:6379/0"
      SESSION_SECRET: "please-change-me"
      SYNC_FREQUENCY: "60"
    volumes:
      - ./config:/app/config
      - ./data:/app/data
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:3000/api/status"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 30s

  mysql:
    image: mysql:8.0
    container_name: one-api-mysql
    restart: unless-stopped
    command:
      - --default-authentication-plugin=mysql_native_password
      - --character-set-server=utf8mb4
      - --collation-server=utf8mb4_unicode_ci
    environment:
      MYSQL_ROOT_PASSWORD: oneapi_pw
      MYSQL_DATABASE: oneapi
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "3306:3306"
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "127.0.0.1", "-u", "root", "-poneapi_pw"]
      interval: 10s
      timeout: 5s
      retries: 10
      start_period: 30s

  redis:
    image: redis:7-alpine
    container_name: one-api-redis
    restart: unless-stopped
    command:
      - redis-server
      - --requirepass
      - redispw
      - --appendonly
      - "yes"
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "-a", "redispw", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 10s

volumes:
  mysql_data:
  redis_data:
```

启动 / Start:

```bash
# 准备一个空目录，把上面保存为 compose.yaml
docker compose up -d
docker compose logs -f one-api-pro
```

看到 `cluster module initialized`（仅在 `CLUSTER_ENABLED=true` 时打印）或 `using MySQL as database` 字样即表示初始化成功。

Look for `using MySQL as database` (and `cluster module initialized` when `CLUSTER_ENABLED=true`) in the logs to confirm a successful boot.

## 多实例共用 MySQL / Redis（共享库版本）

> 与多节点 **集群（去中心化）模式不同**：此处仅是多实例共享 DB 与 Redis 的传统主备/水平扩展方案，仍通过 `SESSION_SECRET` 与 `SYNC_FREQUENCY` 维持一致性。
> This is **not** the decentralized cluster mode. It is the classic horizontal scale-out where every instance shares the same MySQL / Redis; consistency is maintained via `SESSION_SECRET` and `SYNC_FREQUENCY`.

```yaml
services:
  one-api-pro-1:
    image: ghcr.io/modelbus/one-api-pro:latest
    restart: unless-stopped
    ports:
      - "3001:3000"
    environment:
      SQL_DSN: "root:oneapi_pw@tcp(mysql:3306)/oneapi?charset=utf8mb4&parseTime=True&loc=Local"
      REDIS_CONN_STRING: "redis://default:redispw@redis:6379/0"
      SESSION_SECRET: "all-instances-must-share-this"
      SYNC_FREQUENCY: "60"
      NODE_TYPE: master
    depends_on:
      mysql: { condition: service_healthy }
      redis: { condition: service_healthy }

  one-api-pro-2:
    image: ghcr.io/modelbus/one-api-pro:latest
    restart: unless-stopped
    ports:
      - "3002:3000"
    environment:
      SQL_DSN: "root:oneapi_pw@tcp(mysql:3306)/oneapi?charset=utf8mb4&parseTime=True&loc=Local"
      REDIS_CONN_STRING: "redis://default:redispw@redis:6379/0"
      SESSION_SECRET: "all-instances-must-share-this"
      SYNC_FREQUENCY: "60"
      NODE_TYPE: slave
      FRONTEND_BASE_URL: "http://host-master:3001"
    depends_on:
      mysql: { condition: service_healthy }
      redis: { condition: service_healthy }
```

所有实例必须共用相同的 `SESSION_SECRET`；从节点可选 `FRONTEND_BASE_URL` 将页面请求重定向到主节点。

All instances must share the same `SESSION_SECRET`; slave instances may set `FRONTEND_BASE_URL` to redirect page requests to the master.

## 卸载 / Uninstall

```bash
docker compose down            # 停止并删除容器
docker compose down -v         # 同时删除数据卷（谨慎）
```

数据卷 `mysql_data` / `redis_data` 与宿主机 `./data` 目录需手动清理。

Data volumes (`mysql_data`, `redis_data`) and the host `./data` directory must be removed manually.

下一步 / Next: [源码编译](/zh/install/source-build) · [备份与恢复](/zh/install/backup-restore)。
