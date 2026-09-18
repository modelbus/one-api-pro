---
title: Docker Deploy
description: "Run a single instance via the official image."
category: install
order: 2
---

# Docker Deploy

> Run a single instance via the official image.

## Pull the image

```bash
docker pull ghcr.io/modelbus/one-api-pro:latest
```

> The multi-arch image covers both `linux/amd64` and `linux/arm64`.

## Start a container

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

Volume notes:

| Path inside container | Purpose |
| --- | --- |
| `/app/config` | `.env` directory; auto-loaded by `docker-entrypoint.sh` via `--env` |
| `/app/data` | SQLite DB, logs, upload cache; `SQLITE_PATH` points here by default |

Entry-point behaviour:

- Loads `$CONFIG_DIR/.env` (default `/app/config/.env`) via `--env` when present.
- Pass-through CLI arguments are forwarded to the binary.
- Healthcheck: `wget http://localhost:3000/api/status` (built-in `wget`; `HEALTHCHECK` declared in the image).

## Custom port

```bash
docker run -d \
  --name one-api-pro \
  -p 8080:8080 \
  -v /opt/one-api-pro/config:/app/config \
  -v /opt/one-api-pro/data:/app/data \
  -e PORT=8080 \
  ghcr.io/modelbus/one-api-pro:latest
```

The `PORT` variable is consumed by both the binary and the `HEALTHCHECK`.

## External database

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

Pre-create the empty `oneapi` and `oneapi_logs` databases; tables are created by `AutoMigrate`.

## Upgrade & rollback

```bash
# Pull the new version
docker pull ghcr.io/modelbus/one-api-pro:latest
# Recreate the container (data volume is preserved)
docker stop one-api-pro && docker rm one-api-pro
docker run -d --name one-api-pro \
  -p 3000:3000 \
  -v /opt/one-api-pro/config:/app/config \
  -v /opt/one-api-pro/data:/app/data \
  ghcr.io/modelbus/one-api-pro:latest
```

Mounting the same volumes preserves SQLite and `.env`; see [Upgrade](/en/install/upgrade) for the full procedure.

Next: [Docker Compose](/en/install/docker-compose) · [Reverse Proxy](/en/install/reverse-proxy).
