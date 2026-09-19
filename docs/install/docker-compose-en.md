---
title: Docker Compose Deploy
description: Compose One API Pro with a database and Redis for production.
category: install
order: 3
---

# Docker Compose Deploy

> For production with a single node + MySQL/Redis.

## One API Pro only (SQLite)

The minimum compose file:

`docker-compose.yml`:

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

Start:

```bash
docker compose up -d
```

## One API Pro + MySQL (recommended for production)

`docker-compose.yml`:

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

Start:

```bash
docker compose up -d
```

DB schema auto-migrates.

## One API Pro + MySQL + Redis (production)

```yaml
services:
  one-api-pro:
    # ... same as above ...
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

## Key env vars

| Var | Notes |
|---|---|
| `TZ` | Time zone |
| `SQL_DSN` | MySQL DSN; blank = SQLite |
| `REDIS_CONN_STRING` | Redis URL; blank = no Redis |
| `CLUSTER_NODE_ID` / `CLUSTER_NODE_SECRET` / `CLUSTER_NODE_PORT` | Cluster mode (multi-node) |
| More | See [Configuration](./config) |

## Data persistence

Three host mounts: `./data`, `./mysql-data`, `./redis-data`. **Container restarts / image upgrades do NOT lose data.**

## Upgrade

```bash
docker compose pull one-api-pro
docker compose up -d
```

Only the one-api-pro container restarts.

## Backup

```bash
docker compose stop one-api-pro
cp -r ./mysql-data ./backup-$(date +%Y%m%d)
docker compose start one-api-pro
```

Full backup strategy: [Backup & Restore](./backup-restore).

## Related

- [Docker Single-instance](./docker-deploy)
- [Configuration](./config)
- [Backup & Restore](./backup-restore)
- [Reverse Proxy](./reverse-proxy)