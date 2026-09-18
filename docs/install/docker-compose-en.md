---
title: Docker Compose
description: "Orchestrate a single instance or a stack with dependencies via docker-compose."
category: install
order: 3
---

# Docker Compose

> Orchestrate a single instance or a stack with dependencies via docker-compose.

## One API Pro only

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

Start:

```bash
docker compose up -d
```

The `wget`-based healthcheck matches the image's built-in `HEALTHCHECK`.

## Full stack with MySQL & Redis

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

Start:

```bash
# Drop the file above into an empty directory as compose.yaml
docker compose up -d
docker compose logs -f one-api-pro
```

Look for `using MySQL as database` (and `cluster module initialized` when `CLUSTER_ENABLED=true`) in the logs to confirm a successful boot.

## Multi-instance sharing MySQL / Redis

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

All instances must share the same `SESSION_SECRET`; slave instances may set `FRONTEND_BASE_URL` to redirect page requests to the master.

## Uninstall

```bash
docker compose down            # stop and remove containers
docker compose down -v         # also remove volumes (irreversible)
```

Data volumes (`mysql_data`, `redis_data`) and the host `./data` directory must be removed manually.

Next: [Source Build](/en/install/source-build) · [Backup & Restore](/en/install/backup-restore).
