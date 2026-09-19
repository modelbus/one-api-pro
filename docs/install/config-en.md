---
title: Configuration
description: All env vars and their meaning, defaults, and effects.
category: install
order: 5
---

# Configuration

> Before changing a setting, read here: "what is it, where does it come from, what does it change".

## Precedence

```
CLI flag > process env > --env file > default
```

## Where to configure

- **Env vars**: in `docker-compose.yml` `environment:` or `.env`.
- **System settings**: Admin → System Settings. Common runtime options live in the DB (signup switch, mail signature, …).
- **Config file**: `--config /path/to/config.yaml` (advanced).

## Database & cache

### `SQL_DSN`

- **What**: MySQL DSN. Empty = embedded SQLite.
- **Format**: `user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local`
- **Effect**: enables MySQL. Leave empty for SQLite.

### `REDIS_CONN_STRING`

- **What**: Redis URL.
- **Format**: `redis://host:port/db`
- **Effect**: enables Redis. **Required** for multi-node cluster.

## Listener

### `PORT`

- **Default**: `3000`
- **Effect**: changes what to map / proxy upstream.

### `LISTEN`

- **Default**: `0.0.0.0` (all interfaces). Default is fine.

## Time zone

### `TZ`

- **Recommended**: `Asia/Shanghai` (or your business TZ)
- **Effect**: log timestamps, cron, day boundaries.

## Security & rate limit

### `SESSION_SECRET`

- **What**: cookie/session signing key.
- **Default**: auto-generated (changes on restart → all users logged out).
- **Recommended**: long random string in production.
- **Effect**: invalidates all existing sessions when changed.

### `JWT_SECRET`

- JWT signing key for some admin endpoints.
- **Recommended**: set in production.

### `RATE_LIMIT_*`

- Global rate limit knobs.
- **Default**: aggressive, fine for normal use.
- **Effect**: exceeding → 429.

## Cluster

Multi-node only:

- `CLUSTER_NODE_ID`: globally unique per node.
- `CLUSTER_NODE_SECRET`: same across all nodes.
- `CLUSTER_NODE_PORT`: inter-node port (default 3000).
- `CLUSTER_DISCOVERY_INTERVAL`: heartbeat / ping interval (s, default 30).
- `CLUSTER_DEAD_PING_INTERVAL`: retry interval for failed nodes (default 120).
- `CLUSTER_MAX_PING_FAILURES`: failures before auto-disable (default 3).

See [Multi-node Deployment](../decentralization/deployment).

## SMTP

- `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS`, `SMTP_FROM`.
- Not set → mail is disabled (password reset / redemption / order mails all fail).

## Logging

- `LOG_CONSOLE_OUTPUT`: tee logs to stdout (required in Docker).
- `LOG_CONSOLE_LEVEL`: `debug` / `info` / `warn` / `error`.
- `LOG_DB_RETENTION_DAYS`: call log retention (default 90).

## Tiktoken

- `TIKTOKEN_CACHE_DIR`: override the embedded BPE path. Leave empty.
- 4 BPE files embedded (cl100k_base / o200k_base / p50k_base / r50k_base). No network needed.

## Full list

Run `./one-api-pro --help` for every flag.

## Related

- [Docker Single-instance](./docker-deploy)
- [Docker Compose](./docker-compose)
- [Multi-node Deployment](../decentralization/deployment)