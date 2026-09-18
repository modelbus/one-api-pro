---
title: Requirements
description: "Hardware, OS, ports, supported databases and dependencies."
category: install
order: 1
---

# Requirements

> Hardware, OS, ports, supported databases and dependencies.

## Hardware & OS

One API Pro ships as a single Go binary with an embedded SQLite driver and Redis client — no extra system services are required to start.

| Resource | Minimum | Recommended |
| --- | --- | --- |
| CPU | 1 vCPU | 2+ vCPU |
| RAM | 512 MB | 1 GB+ |
| Disk | 1 GB | 10 GB+ (leave headroom for logs & cache) |
| OS | Linux, macOS, Windows | Linux (recommended for production) |

The Docker image is based on `alpine:latest`; native binaries are published for Linux (amd64 / arm64), macOS (amd64 / arm64) and Windows (amd64).

## Ports

The service listens on **`3000`** by default (override with `PORT` or `--port`).
The container `EXPOSE`s `3000`; when exposing 443 through a reverse proxy, enable WebSocket upgrade for streaming endpoints.

## Databases

When `SQL_DSN` is unset, **SQLite** is used (default file `one-api-pro.db`; in the container `/app/data/one-api-pro.db`). Once set, the service auto-switches to:

- **MySQL** — `SQL_DSN=root:123456@tcp(localhost:3306)/oneapi`. Create the empty `oneapi` database up front; tables are created automatically by `AutoMigrate`.
- **PostgreSQL** — `SQL_DSN=postgres://postgres:123456@localhost:5432/oneapi`.

`LOG_SQL_DSN` (optional) splits the `logs` table onto a dedicated database. MySQL or PostgreSQL is recommended because the `logs` table is write-heavy and SQLite can hit `SQLITE_BUSY` under concurrency.

In cluster mode every node must run its **own** MySQL instance (see `decentralization/deployment` for the rationale around the instance-level `auto_increment_offset`).

## Redis (Optional)

The service runs without Redis (in-process cache + direct DB read). Setting `REDIS_CONN_STRING` activates a distributed cache layer that can offload DB reads; if the DB is already low-latency, enabling Redis may introduce short-term staleness — choose based on your trade-off.

Supports single-node, Sentinel (set `REDIS_MASTER_NAME`) and Cluster mode (comma-separated nodes + `REDIS_PASSWORD`).

## Network

Since v0.0.21 the tiktoken encodings for common models are **embedded** into the binary; startup no longer requires network access by default. An HTTP download fallback is only triggered when the user supplies a custom encoding via `TIKTOKEN_CACHE_DIR` whose URL is not in the embedded allow-list.

Outbound traffic for payment callbacks and upstream relays still requires network; whitelist accordingly for air-gapped deployments.

## Dev Dependencies

Toolchain required to build the Go backend and the web theme, matching `AGENTS.md §2`:

| Tool | Version |
| --- | --- |
| Go | 1.25.0 (declared in `go.mod`; toolchain auto-fetched) |
| Node.js | 22+ |
| pnpm | 9 (`pnpm-lock.yaml` is committed) |

Next: [Docker Deploy](/en/install/docker-deploy) · [Docker Compose](/en/install/docker-compose) · [Source Build](/en/install/source-build).
