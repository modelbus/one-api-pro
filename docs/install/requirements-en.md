---
title: System Requirements
description: Hardware, OS, ports, and database needed to run One API Pro.
category: install
order: 1
---

# System Requirements

> Pre-flight checklist before you deploy.

## Hardware

| Scale | Recommended |
|---|---|
| Trial / personal | 1 vCPU / 1 GB RAM / 5 GB disk |
| Small team (< 100 users) | 2 vCPU / 4 GB RAM / 20 GB SSD |
| Production (hundreds–thousands of users) | 4+ vCPU / 8+ GB RAM / 50 GB+ SSD (depends on log retention) |

One API Pro is a single Go binary. CPU scales with request forwarding; memory with log cache and concurrent in-flight calls.

## OS

Linux (Ubuntu 22.04 / Debian 12 / CentOS Stream 9), macOS 12+, Windows 10+. Docker ignores this.

## Ports

| Port | Purpose | Public? |
|---|---|---|
| `3000` | HTTP main port (admin + `/api/*` + `/v1/*`) | **Yes** |
| inter-node | `CLUSTER_NODE_PORT` (default `3000`) | Only for multi-node |

## Database

Pick one:

### Embedded SQLite (default)

Zero config. Good for single instance, < 1000 users, < 50 QPS.

- DB lives at `/app/data/one-api.db` inside the container.
- Backup = copy that file.

### External MySQL (production)

For multi-instance / high QPS / large data.

- MySQL ≥ 5.7 / 8.0
- Create the empty DB + user in advance.
- Connection string goes in env `SQL_DSN`.

> Decide before going live. Migrating SQLite → MySQL later requires a downtime window.

## Redis (optional but recommended)

Multi-node cluster **requires** Redis:

- Redis ≥ 6.0
- Used for inter-node state and rate limiting.
- Single instance can skip Redis.

## Reverse proxy (recommended in production)

Don't expose the container port directly. Put something in front:

- Nginx / Caddy / Cloudflare.
- Handles TLS termination + HSTS + real client IP.
- See [Reverse Proxy](./reverse-proxy).

## Other

- **Domain**: a real domain + TLS certificate.
- **SMTP** (optional): for password reset / redemption / order mails.
- **Payment channels**: to enable top-up / plan checkout, configure [Payment Settings](../pricing/payment-settings) first.

## Pre-deploy checklist

- [ ] Server ≥ 2 vCPU / 4 GB RAM
- [ ] OS Linux / macOS / Windows
- [ ] Port 80/443 (reverse proxy) or 3000 (direct) reachable
- [ ] Database: SQLite default works; or MySQL ≥ 5.7
- [ ] Redis ≥ 6.0 (multi-node only)
- [ ] Reverse proxy (Nginx / Caddy / Cloudflare)
- [ ] Domain + TLS
- [ ] (Optional) SMTP
- [ ] (Optional) Payment channel credentials

## Next

- Fastest start → [Docker Deploy](./docker-deploy)
- Production → [Docker Compose](./docker-compose)
- Build from source → [Source Build](./source-build)