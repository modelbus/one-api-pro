<!--
  README.md OneApi Pro 英文首页
  OneApi Pro English landing README (Quick Start / Deployment / vs oneapi)

  版本: v0.0.21
  日期: 2026-09-18
  作者: modelbus
-->

<!-- SEO keywords: oneapi, one-api, newapi, sub2api, AI API Gateway, LLM gateway, OpenAI compatible, oneapi pro -->
<!-- keywords: oneapi,one-api-pro,one-api,newapi,sub2api AI API Gateway LLM gateway OpenAI compatible Claude Gemini DeepSeek relay -->

<p align="center">
  <a href="https://github.com/modelbus/one-api-pro"><img src="docs/logo.png" width="150" height="150" alt="oneapi"></a>
</p>

<h1 align="center">OneAPI Pro · Enterprise AI API Gateway</h1>

<p align="center">
  The self-hosted, OpenAI-compatible AI API Gateway for <strong>OpenAI</strong>, <strong>Claude</strong>, <strong>Gemini</strong>, <strong>DeepSeek</strong>, <strong>Qwen</strong>, and 30+ providers.<br>
  A next-generation successor to <a href="https://github.com/songquanpeng/one-api" alt="oneapi">oneapi</a>, built with Go + Vue 3 + Arco Design.
</p>

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="license"></a>
  <a href="https://github.com/modelbus/one-api-pro/releases/latest"><img src="https://img.shields.io/github/v/release/modelbus/one-api-pro?color=00ADD8&label=release" alt="oneapi pro release"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/language-Go-00ADD8.svg?logo=go&logoColor=white" alt="language"></a>
  <a href="https://gin-gonic.com/"><img src="https://img.shields.io/badge/framework-Gin-008080.svg?logo=go&logoColor=white" alt="framework"></a>
  <a href="https://vuejs.org/"><img src="https://img.shields.io/badge/frontend-Vue%203-42B883.svg?logo=vue.js&logoColor=white" alt="frontend"></a>
  <a href="https://arco.design/vue"><img src="https://img.shields.io/badge/ui-Arco%20Design-165DFF.svg" alt="ui"></a>
  <a href="https://vitejs.dev/"><img src="https://img.shields.io/badge/build-Vite-646CFF.svg?logo=vite&logoColor=white" alt="build"></a>
  <a href="https://gorm.io/"><img src="https://img.shields.io/badge/database-MySQL%20%7C%20PostgreSQL%20%7C%20SQLite-4479A1.svg?logo=mysql&logoColor=white" alt="database"></a>
</p>

<p align="center">
  👉 <strong>Live Demo</strong>: <a href="http://demo.one-api.pro" alt="oneapi pro demo">http://demo.one-api.pro</a>
  &nbsp;·&nbsp;
  👉 <strong>Demo Account</strong>: <code>root</code> / <code>123456</code>
  &nbsp;·&nbsp;
  👉 <strong>QQ Group</strong>: 1102851586
</p>

<p align="center">
  <strong>English</strong>
  &nbsp;·&nbsp;
  <a href="readme/README.zh.md">简体中文</a>
  &nbsp;·&nbsp;
  <a href="readme/README.zh-TW.md">繁體中文</a>
  &nbsp;·&nbsp;
  <a href="readme/README.ja.md">日本語</a>
  &nbsp;·&nbsp;
  <a href="readme/README.ru.md">Русский</a>
  &nbsp;·&nbsp;
  <a href="readme/README.ko.md">한국어</a>
  &nbsp;·&nbsp;
  <a href="readme/README.ar.md">العربية</a>
  &nbsp;·&nbsp;
  <a href="readme/README.de.md">Deutsch</a>
</p>

---

## What is OneAPI Pro?

**OneAPI Pro** (`one-api-pro`) is an enterprise-grade, self-hosted **AI API Gateway** and **LLM routing platform** written in **Go** with a **Vue 3** admin console. It speaks the **OpenAI-compatible HTTP API** at the front, fans out to **30+ upstream LLM providers** at the back, and adds the missing pieces for production use:

- Per-token quota, expiry, IP allowlist, model allowlist
- Plan / subscription / top-up billing with **WeChat Pay** and **Alipay**
- Decentralized multi-active cluster (no shared database)
- Live channel health checks, cooldown, fallback routing
- A modern Vue 3 dashboard with charts, plans, orders, top-up, cluster nodes

It is the spiritual successor to [`songquanpeng/one-api`](https://github.com/songquanpeng/one-api) — the most widely deployed **one-api** gateway — rebuilt from the ground up with a self-registering Adaptor architecture, stricter permission boundaries, a real subscription engine, and a cluster mode designed for global multi-region deployment.

If you are searching for **oneapi**, **one-api**, **newapi**, or **sub2api**, you have come to the right place.

---

## 📑 Table of Contents

- [🚀 Quick Start](#-quick-start)
- [📦 Deployment](#-deployment)
  - [Manual / Binary](#-manual--binary-deployment)
  - [Docker](#-docker-deployment)
  - [Docker Compose](#-docker-compose)
  - [Multi-node (shared DB)](#-multi-node-shared-database)
  - [Decentralized Cluster](#-decentralized-cluster-multi-active)
- [🔥 vs oneapi](#-vs-oneapi)
- [📸 Screenshots](#-screenshots)
- [📖 Documentation](#-documentation)
- [🗺️ Roadmap](#-roadmap)
- [License](#license)

---

## 🚀 Quick Start

The fastest way to try OneAPI Pro is the pre-built binary. It works on Linux / macOS / Windows with **no external dependencies** (SQLite is embedded).

### 1. Download a release

Grab the binary for your platform from
[**GitHub Releases →**](https://github.com/modelbus/one-api-pro/releases/latest)

Pre-built artifacts are statically linked, **no extraction needed**, just `chmod +x` and run.

### 2. Run

```bash
chmod u+x one-api-pro
./one-api-pro --port 3000 --log-dir ./logs
```

### 3. Open the admin console

Visit **http://localhost:3000** and log in with the default admin account:

| Field    | Value    |
|----------|----------|
| Username | `root`   |
| Password | `123456` |

> ⚠️ **Change the default password immediately** in `Settings → Users` after first login.

### 4. Try the public Demo

If you don't want to install anything, point your browser at the hosted demo:

- URL: <http://demo.one-api.pro>
- Username: **`root`**
- Password: **`123456`**

The demo runs the latest `main` build, resets nightly, and is safe to click around in.

### Build from source (optional)

```bash
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro

# Build the Vue 3 frontend (embeds into the Go binary)
cd web && sh build.sh && cd ..

# Build the Go backend (must run AFTER the frontend build)
go build -ldflags "-s -w" -o one-api-pro
```

For multi-platform packaging, use the bundled `release.sh`:

```bash
./release.sh                # uses the VERSION file
./release.sh v0.1.0         # explicit version
```

Outputs statically linked binaries to `dist/` for `linux-amd64`, `linux-arm64`, `darwin-amd64`, `darwin-arm64`, `windows-amd64.exe`.

---

## 📦 Deployment

### 🔨 Manual / Binary deployment

Best for: single server, VMs without Docker, air-gapped environments.

1. Download from [GitHub Releases](https://github.com/modelbus/one-api-pro/releases/latest) or run `./release.sh`.
2. Start the binary:

   ```bash
   ./one-api-pro --port 3000 --log-dir ./logs
   ```

3. Open `http://<server-ip>:3000` and log in with `root` / `123456`.

SQLite is the default database — zero setup. For MySQL / PostgreSQL, set `SQL_DSN`; see the full env-var reference in [`readme/README.zh.md` → 配置 → 环境变量](readme/README.zh.md).

### 🐳 Docker deployment

Images are published to GitHub Container Registry, multi-arch (`linux/amd64` + `linux/arm64`):

| Tag pattern            | Image                                                |
|------------------------|------------------------------------------------------|
| Latest stable          | `ghcr.io/modelbus/one-api-pro:latest`                |
| Specific version       | `ghcr.io/modelbus/one-api-pro:v0.1.0`                |
| Major version          | `ghcr.io/modelbus/one-api-pro:0`                     |

```bash
mkdir -p ./one-api-pro/config ./one-api-pro/data

docker run -d --name one-api-pro --restart unless-stopped \
  -p 3000:3000 \
  -v $(pwd)/one-api-pro/config:/app/config \
  -v $(pwd)/one-api-pro/data:/app/data \
  ghcr.io/modelbus/one-api-pro:latest
```

Open <http://localhost:3000> and log in with `root` / `123456`. Data persists in `./one-api-pro/data`, config (optional) in `./one-api-pro/config`.

**Switch to MySQL:**

```bash
docker run -d --name one-api-pro -p 3000:3000 \
  -e SQL_DSN='root:123456@tcp(mysql:3306)/oneapi?charset=utf8mb4&parseTime=True&loc=Local' \
  -v $(pwd)/one-api-pro/config:/app/config \
  -v $(pwd)/one-api-pro/data:/app/data \
  ghcr.io/modelbus/one-api-pro:latest
```

### 🐳 Docker Compose

```yaml
# docker-compose.yml
services:
  one-api-pro:
    image: ghcr.io/modelbus/one-api-pro:latest
    container_name: one-api-pro
    restart: unless-stopped
    ports:
      - "3000:3000"
    volumes:
      - ./config:/app/config
      - ./data:/app/data
    environment:
      - TZ=Asia/Shanghai
```

```bash
docker compose up -d
```

### 🏢 Multi-node (shared database)

Same MySQL/Redis in front of N OneAPI Pro instances. Set on every node:

| Variable             | Value                                                      |
|----------------------|------------------------------------------------------------|
| `SESSION_SECRET`     | **Same value** on every node                               |
| `SQL_DSN`            | **Same MySQL** DSN on every node                           |
| `NODE_TYPE`          | `master` on one node, `slave` on the rest                  |
| `SYNC_FREQUENCY`     | e.g. `60` — pull config from DB                            |
| `REDIS_CONN_STRING`  | Per-node Redis (cluster/sentinel supported)                |
| `FRONTEND_BASE_URL`  | On slaves, redirect browser requests to master             |

### 🌐 Decentralized Cluster (multi-active)

Run **N independent OneAPI Pro + MySQL + Redis** nodes and let them sync state at the application layer via GORM event hooks — **no shared database**, ideal for global multi-region, geo-routing, and disaster recovery.

```
                ┌─────────────┐
                │  Nginx/LB   │  (ip_hash → sticky session per user)
                └──────┬──────┘
                       │
       ┌───────────────┼───────────────┐
       │               │               │
┌──────┴──────┐  ┌─────┴──────┐  ┌─────┴──────┐
│  Node A     │  │  Node B    │  │  Node C    │
│  one-api-pro│  │ one-api-pro│  │ one-api-pro│
│  + MySQL    │  │ + MySQL    │  │ + MySQL    │
│  + Redis    │  │ + Redis    │  │ + Redis    │
└──────┬──────┘  └─────┬──────┘  └─────┬──────┘
       └──────── HTTP push events ─────┘
```

Highlights:

- **No shared DB** — every node owns its MySQL/Redis
- **Decentralized** — peers sync via HTTP; no master/slave
- **Conflict-free** — `updated_at` timestamp wins
- **Backwards compatible** — no env vars set = single-node mode

Each node must set `auto_increment_increment = 50` and a unique `auto_increment_offset = CLUSTER_NODE_ID` in MySQL. Full multi-region env templates and Nginx config are in [`readme/README.zh.md` → 集群部署](readme/README.zh.md).

---

## 🔥 vs oneapi

OneAPI Pro is a from-scratch re-architecture of the original [`songquanpeng/one-api`](https://github.com/songquanpeng/one-api). The table below summarizes what changed and why teams migrate.

| Dimension          | one-api                                       | one-api-pro                                                                                  |
|--------------------|-----------------------------------------------|----------------------------------------------------------------------------------------------|
| Project name       | one-api                                       | one-api-pro                                                                                  |
| Adaptor architecture | Centralized constants — `channeltype/define.go` 56-line `iota`, parallel arrays in `url.go`, dual-layer switches in `helper.go`. **Adding a provider requires editing 4 framework files.** | **Self-registering registry** — drop a package, call `Register`, zero framework edits.       |
| Permission model   | Blurred admin/user boundary; any user could mutate settings via API | Tiered `Guest / User / Admin / Root`, **fixed API permission bugs**, fine-grained admin actions |
| Subscription system | None                                          | Full plan engine — per-token / per-call billing, period rate limits, per-model control      |
| Decentralized cluster | No built-in cluster; multi-node needs a shared MySQL | **Decentralized multi-active cluster** — each node has its own MySQL + Redis, synced via app-layer events |
| Directory layout   | `relay/adaptor/` flat 40 dirs, base protocols mixed with providers; `relay/model/` collides with root `model/` | `adaptor/openai/`, `adaptor/anthropic/` as base protocols; `adaptor/provider/` collects 37 vendors; `relay/schema/` removes the collision |
| Admin console      | 3 React themes (default / berry / air), basic CRUD | **Vue 3 + Arco Design**, visualization dashboard, 30+ provider icons                         |
| Maintenance        | Original project stopped updating in 2024     | **Actively maintained**, enterprise-focused roadmap                                          |
| Payment & orders   | None built-in                                 | WeChat Pay Native + Alipay TradePrecreate, plan & top-up orders with full audit trail       |
| Online top-up      | None                                          | Per-user balance top-up via the same WeChat / Alipay channels                                |

> If you are still running `one-api` and hitting the limits above, the migration path is straightforward — same OpenAI-compatible API surface, drop-in replacement at the client side. See [`readme/README.zh.md`](readme/README.zh.md) for the full feature comparison.

---

## 📸 Screenshots

| Dashboard | Token Management |
|:-:|:-:|
| ![Dashboard](docs/Demo-Index.png) | ![Tokens](docs/Demo-Token.png) |

| Plans | Subscriptions |
|:-:|:-:|
| ![Plans](docs/Demo-Plan.png) | ![Subscriptions](docs/Demo-Subscribe.png) |

| Cluster Nodes |
|:-:|
| ![Cluster](docs/Demo-cluster.png) |

---

## 📖 Documentation

| Doc | Description |
| --- | --- |
| [`readme/README.zh.md`](readme/README.zh.md)            | Full reference — install, configure, deploy, env vars, CLI flags, cluster guide, roadmap (Simplified Chinese) |
| [`readme/README.en.md`](readme/README.en.md)            | English full reference                                                  |
| [`readme/README.zh-TW.md`](readme/README.zh-TW.md)      | 繁體中文                                                               |
| [`readme/README.ja.md`](readme/README.ja.md)            | 日本語                                                                  |
| [`readme/README.ru.md`](readme/README.ru.md)            | Русский                                                                 |
| [`readme/README.ko.md`](readme/README.ko.md)            | 한국어                                                                  |
| [`readme/README.ar.md`](readme/README.ar.md)            | العربية                                                                 |
| [`readme/README.de.md`](readme/README.de.md)            | Deutsch                                                                 |
| [`docs/API.md`](docs/API.md)                            | HTTP API reference (auth, admin endpoints, OpenAI-compatible, cluster APIs) |

---

## 🗺️ Roadmap

**Shipped**

- Architecture-level refactor — self-registering Adaptor, zero framework edits to add a provider
- Vue 3 admin console with Arco Design, visualization dashboard, 30+ provider icons
- Plan & subscription engine — per-token / per-call billing, period rate limits, per-model control
- Decentralized multi-active cluster — GORM-event-driven + HTTP push, no shared DB
- Fine-grained cost accounting — Prompt / Completion / Cached priced independently, group discounts stack
- Tiered permissions — `Guest / User / Admin / Root`, fixed API permission bugs in the original
- OpenAI-compatible API — models / chat / completions / embeddings / images / audio / moderations
- Real payments (WeChat Native + Alipay TradePrecreate) — plan orders + top-up orders with idempotent activation
- Top-up settings center — global switch, preset amounts, custom amounts, exchange rate
- Order center (admin view) — list / search / detail / edit / delete, multi-dimensional filters
- Channel diagnostics & smart routing — auto cooldown, fallback, auto-disable on low success rate
- i18n — Chinese / English across admin, landing page, and legal pages

**In progress**

- Enhanced channel diagnostics — independent diagnostics panel, node-level ping, manual review
- Richer usage analytics & export

**Planned**

- More languages — 繁體中文 / 日本語 / 한국어 / Русский / Deutsch / العربية in the admin console
- More payment channels — Apple Pay, UnionPay, Stripe
- Refund engine — async refund API + subscription / quota rollback + visual refund ledger
- Finance integrations — sync top-up / spend / refund to mainstream accounting platforms
- Token balance alerts — low-balance notifications via multi-channel
- Audit logs & reports — full operation audit + visual reports for compliance
- AI-driven analytics — LLM-based cost / channel-health recommendations
- Plugin extension mechanism
- Enterprise SSO / LDAP

> PRs and Issues are welcome — see [Issues](https://github.com/modelbus/one-api-pro/issues).

---

## License

[MIT](LICENSE)
