<p align="center">
  <img src="../docs/logo.png" width="150" height="150" alt="one-api-pro logo">
</p>

<p align="center">
  One Api Pro · An enterprise-grade AI API Gateway, built with Go
</p>
<p align="center">
  Deep refactor and extension of <a href="https://github.com/songquanpeng/one-api">one-api</a> (by <a href="https://github.com/songquanpeng">JustSong</a>) — many thanks to the original author.
</p>

<p align="center">
  <a href="../LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="license"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/language-Go-00ADD8.svg?logo=go&logoColor=white" alt="language"></a>
  <a href="https://gin-gonic.com/"><img src="https://img.shields.io/badge/framework-Gin-008080.svg?logo=go&logoColor=white" alt="framework"></a>
  <a href="https://vuejs.org/"><img src="https://img.shields.io/badge/frontend-Vue%203-42B883.svg?logo=vue.js&logoColor=white" alt="frontend"></a>
  <a href="https://arco.design/vue"><img src="https://img.shields.io/badge/ui-Arco%20Design-165DFF.svg" alt="ui"></a>
  <a href="https://vitejs.dev/"><img src="https://img.shields.io/badge/build-Vite-646CFF.svg?logo=vite&logoColor=white" alt="build"></a>
  <a href="https://gorm.io/"><img src="https://img.shields.io/badge/database-MySQL%20%7C%20PostgreSQL%20%7C%20SQLite-4479A1.svg?logo=mysql&logoColor=white" alt="database"></a>
  <a href="https://github.com/modelbus/one-api-pro"><img src="https://img.shields.io/badge/cluster-decentralized-FF6B6B.svg" alt="cluster"></a>
</p>

<p align="center">
  <a href="../README.md">简体中文</a>
  &nbsp;·&nbsp;
  <a href="README.en.md">English</a>
  &nbsp;·&nbsp;
  <a href="README.zh-TW.md">繁體中文</a>
  &nbsp;·&nbsp;
  <a href="README.ja.md">日本語</a>
  &nbsp;·&nbsp;
  <a href="README.ru.md">Русский</a>
  &nbsp;·&nbsp;
  <a href="README.ko.md">한국어</a>
  &nbsp;·&nbsp;
  <a href="README.ar.md">العربية</a>
  &nbsp;·&nbsp;
  <a href="README.de.md">Deutsch</a>
</p>

---

## 📑 Contents

- [🚀 Quick Start](#-quick-start)
- [🔧 Tech Stack](#-tech-stack)
  - [Go Backend](#go-backend)
  - [Vue 3 Frontend](#vue-3-frontend)
- [✨ Feature Highlights](#-feature-highlights)
- [🔥 One Api Pro vs. one-api](#-one-api-pro-vs-one-api)
- [📸 Screenshots](#-screenshots)
- [⚙️ Configuration](#%EF%B8%8F-configuration)
  - [🔧 Environment Variables](#-environment-variables)
  - [⌨️ Command-line Arguments](#%EF%B8%8F-command-line-arguments)
- [📖 API Documentation](#-api-documentation)
- [📦 Deployment](#-deployment)
  - [🔨 Manual Deployment](#-manual-deployment)
  - [🏢 Multi-host Deployment](#-multi-host-deployment)
  - [🌐 Decentralized Cluster Deployment](#-decentralized-cluster-deployment)
- [🗺️ Roadmap](#%EF%B8%8F-roadmap)
- [📄 License](#-license)

---

## 🚀 Quick Start

### 1. Get the binary

Either download the prebuilt binary from the [GitHub Releases](https://github.com/modelbus/one-api-pro/releases/latest) page, or build from source:

```bash
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro
```

### 2. (Source build) Build the Vue 3 admin console

```bash
cd web
sh build.sh        # builds every theme listed in web/THEMES (default: default-pro)
cd ..
```

### 3. (Source build) Build the Go binary

> The backend must be built **after** the frontend, so the latest UI assets are embedded.

```bash
go build -ldflags "-s -w" -o one-api-pro
```

### 4. (Optional) One-click multi-platform packaging

Use the `release.sh` script in the repo root to download dependencies, build the frontend and cross-compile all platforms in one step:

```bash
./release.sh                          # version from the VERSION file
./release.sh v0.1.0                   # specify a version
./release.sh v0.1.0 --skip-frontend   # skip the frontend build (reuse web/build)
```

> Prerequisites: `go`, `node`, `npm`. The version is read from the root `VERSION` file (with or without the `v` prefix).

The outputs are **statically linked bare executables** (no extraction needed, run directly) written to `dist/`:

```
dist/one-api-pro-linux-amd64
dist/one-api-pro-linux-arm64
dist/one-api-pro-windows-amd64.exe
dist/one-api-pro-darwin-amd64
dist/one-api-pro-darwin-arm64
```

> The `linux-*` binaries are statically linked and work on both CentOS and Ubuntu. GitHub Releases are built and published automatically by `.github/workflows/release.yml` when a `v*` tag is pushed, mirroring the local `release.sh` output.

### 5. Run

```bash
./one-api-pro --port 3000 --log-dir ./logs
```

Visit `http://localhost:3000` and log in with the default account `root / 123456`.

> For production deployments, see [📦 Deployment](#-deployment). For API references, see [📖 API Documentation](#-api-documentation).

---

## 🔧 Tech Stack

One Api Pro is built with the help of these outstanding open-source projects — many thanks to their authors.

### Go Backend

| Library | Purpose |
| --- | --- |
| [Gin](https://github.com/gin-gonic/gin) | HTTP web framework |
| [GORM](https://gorm.io) | ORM, supports SQLite / MySQL / PostgreSQL |
| [go-redis/redis](https://github.com/go-redis/redis) | Redis client |
| [golang-jwt/jwt](https://github.com/golang-jwt/jwt) | JWT authentication |
| [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2) | AWS Bedrock integration |
| [Google API Go Client](https://github.com/googleapis/google-api-go-client) | Google Gemini / PaLM2 integration |
| [pkoukk/tiktoken-go](https://github.com/pkoukk/tiktoken-go) | Token counting |
| [gorilla/websocket](https://github.com/gorilla/websocket) | WebSocket support (e.g. iFlytek) |
| [joho/godotenv](https://github.com/joho/godotenv) | `.env` configuration file parsing |

### Vue 3 Frontend

| Library | Purpose |
| --- | --- |
| [Vue 3](https://vuejs.org) | Frontend framework (Composition API) |
| [Vite](https://vitejs.dev) | Build tool |
| [Arco Design Vue](https://arco.design/vue) | UI component library |
| [Pinia](https://pinia.vuejs.org) | State management |
| [Vue Router 4](https://router.vuejs.org) | Routing |
| [Axios](https://axios-http.com) | HTTP client |
| [ECharts](https://echarts.apache.org) | Data visualization |
| [vue-i18n](https://vue-i18n.intlify.dev) | Internationalization |

---

## ✨ Feature Highlights

One Api Pro is an **enterprise-grade AI API Gateway** written in Go + Vue 3, deep refactored on top of [one-api](https://github.com/songquanpeng/one-api) with architecture-level changes and enterprise-grade enhancements, while keeping every feature of the original.

### 🖥️ Visual Dashboard

A brand-new Vue 3 + Arco Design admin console delivers a visual dashboard with core metrics, usage trends and per-model distribution at a glance.

| Key-metric cards | Usage trend chart |
|:---:|:---:|
| ![Dashboard](../docs/Demo-Index.png) | ![Dashboard](../docs/Demo-Index.png) |

### 🔑 Fine-grained Token Management

Multi-dimensional token controls: per-model allowlists, sub-net IP restrictions, quota caps, expiration, unlimited quotas. Permissions can be scoped all the way down to a single model.

| Token management |
|:---:|
| ![Token management](../docs/Demo-Token.png) |

### 📦 Subscription & Plan System

A complete plan & subscription system: token- or request-based billing, period-based rate limits (hourly / weekly / monthly), per-model controls, plus recommended plans and pricing.

| Plan management | Subscription management |
|:---:|:---:|
| ![Plan management](../docs/Demo-Plan.png) | ![Subscription management](../docs/Demo-Subscribe.png) |

### 💳 Orders & Real Payments

Every subscription checkout leaves a full **order audit trail** (order number, user, plan snapshot JSON, amount, payment channel, status, paid-at timestamp, provider trade number) supporting two order types — subscription and top-up. Native integrations for **WeChat Pay Native** (PC QR) and **Alipay Face-to-Face** (TradePrecreate) ship out of the box, with `bank` / `offline` / `free` reserved for admin-side channels. Upgrade pricing is calculated automatically from remaining days, and an "operate in stack mode" toggle is hot-switchable from **Settings → Operations → Plan** to keep the old subscription running alongside the new one.

| Order center | Payment config |
|:---:|:---:|
| ![Order center](../docs/Demo-Order.png) | ![Payment config](../docs/Demo-Payment.png) |

### 🌐 Decentralized Active-Active Cluster

Each node runs independent MySQL + Redis. Application-layer events keep data in sync across nodes — no shared database required, naturally supporting low-latency access from any region.

| Cluster node management |
|:---:|
| ![Cluster node management](../docs/Demo-cluster.png) |

### 🧩 Other Core Capabilities

- **30+ model providers** — OpenAI / Anthropic / Gemini / DeepSeek / Qwen / Wenxin / iFlytek / Zhipu, all unified behind an OpenAI-compatible interface.
- **Precise cost accounting** — token- or request-based billing, independent Prompt / Completion / Cached pricing, stacked group discounts, period usage tracking.
- **Channel load balancing** — weighted random distribution, automatic failover, cool-down / disable policies, channel concurrency and RPM rate limits.
- **Multi-level permissions** — Guest / User / Admin / Root, with upstream API authorization loophole fixed and granular admin action controls.
- **Enterprise-grade security** — full-link HTTPS, token authentication, sub-net IP restrictions, real-time auditing logs.

---

## 🔥 One Api Pro vs. one-api

| Dimension | one-api | one-api-pro |
| --- | --- | --- |
| Project name | one-api | one-api-pro |
| Adaptor architecture | Centralized constants: `channeltype/define.go` (56-line `iota`) + parallel arrays in `url.go` + 2-tier `switch` in `helper.go`; adding a provider requires edits across 4 framework files | Self-registration (`registry` + `register.go`); adding a provider only needs a new package + one `register` call, zero framework changes |
| Permission granularity | Admin / user boundaries are blurry — anyone can change system settings through the API | Tiered permission system with the API authorization loophole fixed and granular admin controls |
| Subscription mode | No plan / subscription system | Full plan & subscription + period rate limits + per-model controls |
| Decentralized cluster | No independent cluster support; multi-host deployments share MySQL | Decentralized active-active cluster — each node runs its own MySQL + Redis, kept in sync through application-layer events with no shared database |
| Directory layout | `relay/adaptor/` flattens 40 directories; base protocols and providers are mixed together; `relay/model/` collides with root `model/` | `adaptor/openai/`, `adaptor/anthropic/` as base protocols; `adaptor/provider/` unifies 37 providers; `relay/schema/` removes the naming collision |
| Admin console | 3 React themes (`default` / `berry` / `air`), basic management features | Brand-new Vue 3 + Arco Design console with a visual dashboard |
| Maintenance status | Upstream project officially stopped in 2024 | Actively maintained, optimized for enterprise-grade scenarios |

---

## 📸 Screenshots

### 🖥️ Dashboard
![Dashboard](../docs/Demo-Index.png)

### 🔑 Token Management
![Token management](../docs/Demo-Token.png)

### 📦 Plan Management
![Plan management](../docs/Demo-Plan.png)

### 🔄 Subscription Management
![Subscription management](../docs/Demo-Subscribe.png)

### 🌐 Cluster Node Management
![Cluster nodes](../docs/Demo-cluster.png)

---

## ⚙️ Configuration

One Api Pro is ready to use out of the box.

You can further customize it through environment variables or command-line arguments, then log in as `root` and continue configuration inside the admin console.

> **Tip:** if you are unsure about the meaning of a setting, delete its value temporarily to see the in-app hint.

### 🔧 Environment Variables

> One Api Pro reads environment variables from a `.env` file. See `.env.example` for the full list — copy it to `.env` before use. You can also point to any path with `--env` (see below).

1. `REDIS_CONN_STRING` — enables Redis as a caching layer.
   - Example: `REDIS_CONN_STRING=redis://default:redispw@localhost:49153`
   - If your database latency is already low, Redis is not strictly required — enabling it may even introduce short-term staleness.
   - For Redis Sentinel / cluster mode, set this to a comma-separated node list, e.g. `localhost:49153,localhost:49154,localhost:49155`, and additionally set:
     - `REDIS_PASSWORD` — password for Sentinel / cluster mode.
     - `REDIS_MASTER_NAME` — master name for Sentinel mode.
2. `SESSION_SECRET` — pins the session key, keeping existing cookies valid across restarts.
   - Example: `SESSION_SECRET=random_string`
3. `SQL_DSN` — uses an external database instead of SQLite. MySQL or PostgreSQL.
   - Examples:
     - MySQL: `SQL_DSN=root:123456@tcp(localhost:3306)/oneapi`
     - PostgreSQL: `SQL_DSN=postgres://postgres:123456@localhost:5432/oneapi` (in-progress, feedback welcome)
   - The `oneapi` database must exist beforehand; tables are auto-created on first run.
   - Cloud databases that require identity verification may need `?tls=skip-verify` in the DSN.
   - Tunable pool parameters (defaults shown):
     - `SQL_MAX_IDLE_CONNS` — max idle connections, default `100`.
     - `SQL_MAX_OPEN_CONNS` — max open connections, default `1000`. Reduce if you see `Error 1040: Too many connections`.
     - `SQL_CONN_MAX_LIFETIME` — connection lifetime in minutes, default `60`.
4. `LOG_SQL_DSN` — uses a dedicated database for the `logs` table.
5. `FRONTEND_BASE_URL` — when set on a slave node, redirects page requests to this address.
   - Example: `FRONTEND_BASE_URL=https://openai.justsong.cn`
6. `MEMORY_CACHE_ENABLED` — enables in-memory caching. Causes minor staleness on quota updates. Default `false`.
   - Example: `MEMORY_CACHE_ENABLED=true`
7. `SYNC_FREQUENCY` — when caching is enabled, how often (in seconds) configurations are synced from the database. Default `600`.
   - Example: `SYNC_FREQUENCY=60`
8. `NODE_TYPE` — `master` or `slave`. Default `master`.
   - Example: `NODE_TYPE=slave`
9. `CHANNEL_UPDATE_FREQUENCY` — interval in minutes for periodic channel-balance refresh. No refresh if unset.
   - Example: `CHANNEL_UPDATE_FREQUENCY=1440`
10. `CHANNEL_TEST_FREQUENCY` — interval in minutes for periodic channel health checks. No check if unset.
    - Example: `CHANNEL_TEST_FREQUENCY=1440`
11. `POLLING_INTERVAL` — interval (seconds) between batch channel-balance / health requests. Default `0` (no spacing).
    - Example: `POLLING_INTERVAL=5`
12. `BATCH_UPDATE_ENABLED` — enables batched DB updates for user-quota changes. Causes minor staleness. Default `false`.
    - Example: `BATCH_UPDATE_ENABLED=true`
    - Helpful when you hit too many DB connections.
13. `BATCH_UPDATE_INTERVAL` — batch update aggregation window in seconds, default `5`.
    - Example: `BATCH_UPDATE_INTERVAL=5`
14. Rate limits:
    - `GLOBAL_API_RATE_LIMIT` — global API rate limit (excluding relay requests), per IP, per 3 minutes. Default `180`.
    - `GLOBAL_WEB_RATE_LIMIT` — global Web rate limit, per IP, per 3 minutes. Default `60`.
15. Tokenizer caches:
    - `TIKTOKEN_CACHE_DIR` — caches tiktoken encodings (e.g. `gpt-3.5-turbo`, `gpt-4`, `gpt-4o`) locally. On startup the encodings are downloaded on demand; if the network is restricted or offline, the download times out (about 30 seconds) and the system automatically falls back to approximate token counting (about `0.38 × character count`), so the service still starts normally. For precise billing, pre-download the encodings into this directory on a networked machine, then move them into the offline environment.
    - `DATA_GYM_CACHE_DIR` — same purpose, lower priority than `TIKTOKEN_CACHE_DIR`.
16. `RELAY_TIMEOUT` — relay request timeout in seconds. No timeout by default.
17. `RELAY_PROXY` — proxy URL for upstream API requests.
18. `USER_CONTENT_REQUEST_TIMEOUT` — user content download timeout in seconds.
19. `USER_CONTENT_REQUEST_PROXY` — proxy for user-uploaded content (e.g. images).
20. `SQLITE_BUSY_TIMEOUT` — SQLite lock-wait timeout in milliseconds. Default `3000`.
21. `GEMINI_SAFETY_SETTING` — Gemini safety setting. Default `BLOCK_NONE`.
22. `GEMINI_VERSION` — Gemini API version used by One Api Pro. Default `v1`.
23. `THEME` — admin theme. Default `default-pro` (Vue 3). Other options: `default` / `berry` / `air` (legacy React themes). See [web/README.md](../web/README.md).
24. `ENABLE_METRIC` — auto-disable channels with low success rate. Default `false`.
25. `METRIC_QUEUE_SIZE` — sliding window size for success-rate statistics. Default `10`.
26. `METRIC_SUCCESS_RATE_THRESHOLD` — success-rate threshold. Default `0.8`.
27. `INITIAL_ROOT_TOKEN` — if set, a root token with this value is auto-created on first startup.
28. `INITIAL_ROOT_ACCESS_TOKEN` — if set, a root access token with this value is auto-created on first startup.
29. `ENFORCE_INCLUDE_USAGE` — force `usage` in stream responses. Default `false`.
30. `TEST_PROMPT` — prompt used by channel/model self-tests. Default `Print your model name exactly and do not output without any other text.`.

#### 🌐 Decentralized Cluster Variables

> When none of the following are set, One Api Pro runs in single-node mode with no side effects.

1. `CLUSTER_ENABLED` — enable cluster mode. Default off.
   - Example: `CLUSTER_ENABLED=true`
2. `CLUSTER_NODE_ID` — node ID (1–49). Must match MySQL's `auto_increment_offset`, unique across nodes.
   - Example: `CLUSTER_NODE_ID=1`
3. `CLUSTER_NODE_NAME` — display name. Default `node-{NODE_ID}`.
   - Example: `CLUSTER_NODE_NAME=node-cn`
4. `CLUSTER_NODE_ADDRESS` — public address other nodes use to reach this one (must include protocol).
   - Example: `CLUSTER_NODE_ADDRESS=https://cn.example.com`
5. `CLUSTER_SECRET` — initial per-node secret. Used as seed on first start, then stored in the DB; admins can rotate it later.
   - Example: `CLUSTER_SECRET=MyClusterSecret123abc`
6. `CLUSTER_SEEDS` — comma-separated seed nodes used for first-time discovery; one reachable node is enough. The very first node can leave this empty (or use its own address).
   - Example: `CLUSTER_SEEDS=https://cn.example.com`
   - Multiple: `CLUSTER_SEEDS=https://cn.example.com,https://us.example.com`
7. `CLUSTER_PUSH_INTERVAL` — sync event push interval in seconds. Default `3`.
8. `CLUSTER_DISCOVERY_INTERVAL` — node discovery interval in seconds. Default `30`.
9. `CLUSTER_DEAD_PING_INTERVAL` — ping interval for unreachable nodes in seconds. Default `120`.
10. `CLUSTER_MAX_PING_FAILURES` — consecutive failures before a node is marked failed. Default `3`.
11. `CLUSTER_SYNC_LOGS` — sync the `logs` table across nodes. Default `true`.
    - Example: `CLUSTER_SYNC_LOGS=false`
12. `CLUSTER_BATCH_SIZE` — max events per push. Default `50`.

### ⌨️ Command-line Arguments

1. `--port <port_number>` — listen port. Default `3000`.
   - Example: `--port 3000`
2. `--log-dir <log_dir>` — log directory. Default `./logs` in the working directory.
   - Example: `--log-dir ./logs`
3. `--env <env_file_path>` — configuration file path; supports relative and absolute paths. Auto-loads `./.env` when unset.
   - Example: `--env ./config.env`
   - Example: `--env /etc/one-api-pro/production.env`
   - Multiple-instance example:
     ```bash
     ./one-api-pro --env ./instances/instance1.env --port 3001 &
     ./one-api-pro --env ./instances/instance2.env --port 3002 &
     ```
   - Precedence: CLI args > system env vars > `--env` file > defaults
4. `--version` — print version and exit.
   - Example: `./one-api-pro --version`
   - Version source (highest priority first):
     1. The `VERSION` file in the current working directory or next to the executable (auto-detects the `v` prefix, e.g. `0.0.2` or `v0.0.2`);
     2. The version injected at build time via `-ldflags "-X .../common.Version=..."` (both `release.sh` and CI do this automatically);
     3. The hard-coded default in `common/constants.go`.
   - Therefore you only need to maintain the root `VERSION` file to keep `--version`, the startup log, the `/api/status` endpoint and the dashboard badge consistent.
5. `--help` — print command-line help.
   - Example: `./one-api-pro --help`

---

## 📖 API Documentation

The full API reference lives at [docs/API.md](../docs/API.md), covering:

- **Authentication** — Cookie Session / Access Token / API Key (Bearer Token)
- **Admin endpoints** — full CRUD for model pricing, group discounts, channels, tokens, users, logs, redemption codes, plans, subscriptions, etc.
- **OpenAI-compatible endpoints** — `/v1/models`, `/v1/chat/completions`, `/v1/embeddings`, image, audio, moderation, etc.
- **Cluster management API** — node discovery, heartbeats, data sync endpoints for decentralized clusters

👉 [View the full API documentation →](../docs/API.md)

---

## 📦 Deployment

### 🔨 Manual Deployment

#### 1. Get the executable

Choose one of the following:

**Option 1: Download a prebuilt binary (recommended)**

Download the bare executable for your platform (Linux / macOS / Windows) from [GitHub Releases](https://github.com/modelbus/one-api-pro/releases/latest) — no extraction needed, run it directly.

**Option 2: One-click packaging with release.sh**

```shell
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro
./release.sh            # multi-platform packaging, outputs to dist/
```

**Option 3: Build from source**

```shell
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro

# Build the Vue 3 admin console (per web/THEMES)
cd web
sh build.sh

# Build the backend (must run AFTER the frontend build, so the latest assets are embedded)
cd ..
go build -ldflags "-s -w" -o one-api-pro
```

#### 2. Run

```shell
chmod u+x one-api-pro
./one-api-pro --port 3000 --log-dir ./logs
```

#### 3. Access

Open [http://localhost:3000/](http://localhost:3000/) and log in. The default account is `root` / `123456`.

### 🏢 Multi-host Deployment

1. Set the same `SESSION_SECRET` on every host.
2. Set `SQL_DSN` so all hosts share the same MySQL database (SQLite won't work).
3. Set `NODE_TYPE=slave` on every secondary host (default is `master`).
4. Set `SYNC_FREQUENCY` so every host periodically syncs configuration from the database. With a remote database this is recommended regardless of master/slave, and Redis should be enabled.
5. Secondary hosts may set `FRONTEND_BASE_URL` to redirect page requests to the master.
6. Install **separate** Redis instances on every secondary host and set `REDIS_CONN_STRING` accordingly. This can hit the DB zero times while the cache is warm, reducing latency (see env vars for Sentinel / cluster support).
7. If the master also has high DB latency, enable Redis + `SYNC_FREQUENCY` on the master as well.

See [Environment Variables](#-environment-variables) for the full set of options.

### 🌐 Decentralized Cluster Deployment

The cluster mode lets multiple nodes each run their own One Api Pro + MySQL, kept in sync through application-layer events. No shared database is required.

> **Use cases:** global multi-region deployment, low-latency local access, HA / DR, multi-node load balancing.

#### 🗺️ Architecture

```
                    ┌─────────────┐
                    │  Nginx/LB   │  (single entry, ip_hash load balancing)
                    └──────┬──────┘
                           │
            ┌──────────────┼──────────────┐
            │              │              │
     ┌──────┴──────┐ ┌────┴───────┐ ┌───┴────────┐
     │  Node A     │ │  Node B     │ │  Node C     │
     │ (one-api-pro)   │ │ (one-api-pro)   │ │ (one-api-pro)   │
     │ + MySQL     │ │ + MySQL     │ │ + MySQL     │
     │ + Redis     │ │ + Redis     │ │ + Redis     │
     └──────┬──────┘ └─────┬──────┘ └────┬────────┘
            │              │              │
            └────── HTTP push of sync events ──────┘
```

#### ⭐ Core Characteristics

- **Decentralized** — every node is equal; any data change is actively pushed to all alive nodes.
- **Zero-invasion** — GORM callbacks capture data changes; no business-code modifications.
- **Async push** — sync runs in a background goroutine, never blocking the main flow.
- **Conflict resolution** — comparison is based on `updated_at`; only newer data is written.
- **Rate-limit sync** — channel concurrency & RPM counters sync across nodes via DB tables.
- **Single-node compatible** — without the cluster env vars, the system runs in plain single-node mode.

#### 📊 Sync Scope

| Table | Synced? | Notes |
| --- | --- | --- |
| users | ✅ | user accounts |
| tokens | ✅ | API tokens |
| channels | ✅ | provider channels |
| abilities | ✅ | channel abilities |
| options | ✅ | system settings |
| redemptions | ✅ | redemption codes |
| plans | ✅ | subscription plans |
| user_plans | ✅ | user subscriptions |
| plan_usages | ✅ | plan usage |
| channel_counters | ✅ | channel rate-limit counters |
| cluster_nodes | 🔄 Discovery | maintained by the discovery mechanism, not data sync |
| logs | ⚠️ Optional | controlled by `CLUSTER_SYNC_LOGS` |

#### 🚀 Deployment Steps

**1. MySQL — every node must use an independent MySQL instance.**

Each node needs its **own MySQL instance** — you cannot deploy multiple nodes by simply creating multiple databases in a single MySQL instance, because `auto_increment_offset` is an instance-level variable.

```ini
# Node 1 my.cnf
[mysqld]
server-id = 1
auto_increment_increment = 50
auto_increment_offset = 1
log_bin = mysql-bin
binlog_format = ROW

# Node 2 my.cnf
[mysqld]
server-id = 2
auto_increment_increment = 50
auto_increment_offset = 2
log_bin = mysql-bin
binlog_format = ROW

# Node 3 my.cnf
[mysqld]
server-id = 3
auto_increment_increment = 50
auto_increment_offset = 3
log_bin = mysql-bin
binlog_format = ROW
```

> `auto_increment_increment` is set to 50 — supports up to 50 nodes. Each node's `offset` must match its `CLUSTER_NODE_ID` and be unique.

> **Important:** `auto_increment_increment` and `auto_increment_offset` are MySQL **instance-level** variables. They apply to every database inside an instance; they cannot be set per database or per table (MySQL's table option only supports an `AUTO_INCREMENT` starting value, not step size). So every node **must run its own MySQL instance** — you cannot deploy multiple nodes by creating separate databases in the same MySQL. To run multiple instances on one machine, start multiple `mysqld` processes on different ports or use multiple Docker containers.

> **About `server-id` and binlog:** `server-id` must be unique across all MySQL instances in the cluster. `log_bin` and `binlog_format=ROW` are strongly recommended — they enable future master-slave replication and point-in-time recovery. The cluster data sync itself doesn't depend on binlog (it goes through GORM callbacks), but binlog adds an extra layer of reliability.

**2. Redis — every node must use an independent Redis instance.**

Each node needs **its own Redis instance** (different port or different machine). Redis is used only for local cache and rate limiting here; cluster traffic does not flow through Redis.

**3. Bootstrap a new node**

A newly added node must first obtain a data snapshot from an existing node:

```bash
# Option 1: dump + import
mysqldump -h existing-node -u root -p oneapi > backup.sql
mysql -u root -p oneapi < backup.sql

# Option 2: snapshot API (after the existing node is up)
curl -H "X-Cluster-Secret: your-secret" \
  "https://existing-node/api/cluster/snapshot?tables=users,tokens,channels,abilities,options,redemptions,plans,user_plans,plan_usages" \
  -o snapshot.json
```

**4. Environment configuration (full example)**

Below is a complete `.env` for a 3-node cluster. Every node uses its own MySQL & Redis instances, with distinct ports and paths.

**Node 1 — China (`/opt/one-api-pro/node1/.env`):**
```bash
# ========================
# Basics
# ========================
PORT=3000
SYSTEM_NAME=One Api Pro Cluster

# ========================
# Database (independent MySQL)
# ========================
SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node1?charset=utf8mb4&parseTime=True&loc=Local

# ========================
# Redis (independent)
# ========================
REDIS_CONN_STRING=redis://127.0.0.1:6379/0

# ========================
# Cluster
# ========================
CLUSTER_ENABLED=true
CLUSTER_NODE_ID=1
CLUSTER_NODE_NAME=node-cn
CLUSTER_NODE_ADDRESS=https://cn.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me

# Seed nodes (only needed for first-time discovery)
# First node: leave empty or use its own address
# Later nodes: any alive node is enough
CLUSTER_SEEDS=https://cn.example.com,https://us.example.com,https://eu.example.com

# ========================
# Cluster tuning (optional)
# ========================
CLUSTER_DISCOVERY_INTERVAL=30
CLUSTER_DEAD_PING_INTERVAL=120
CLUSTER_MAX_PING_FAILURES=3
CLUSTER_PUSH_INTERVAL=3
CLUSTER_SYNC_LOGS=true
CLUSTER_BATCH_SIZE=50
```

**Node 2 — US (`/opt/one-api-pro/node2/.env`):**
```bash
PORT=3001
SYSTEM_NAME=One Api Pro Cluster

SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node2?charset=utf8mb4&parseTime=True&loc=Local

REDIS_CONN_STRING=redis://127.0.0.1:6380/0

CLUSTER_ENABLED=true
CLUSTER_NODE_ID=2
CLUSTER_NODE_NAME=node-us
CLUSTER_NODE_ADDRESS=https://us.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me   # must match node 1

# One reachable node is enough
CLUSTER_SEEDS=https://cn.example.com
```

**Node 3 — Europe (`/opt/one-api-pro/node3/.env`):**
```bash
PORT=3002
SYSTEM_NAME=One Api Pro Cluster

SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node3?charset=utf8mb4&parseTime=True&loc=Local

REDIS_CONN_STRING=redis://127.0.0.1:6381/0

CLUSTER_ENABLED=true
CLUSTER_NODE_ID=3
CLUSTER_NODE_NAME=node-eu
CLUSTER_NODE_ADDRESS=https://eu.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me   # must match all nodes

CLUSTER_SEEDS=https://cn.example.com
```

**Configuration cheat-sheet:**

| Variable | Node 1 | Node 2 | Node 3 | Notes |
| --- | --- | --- | --- | --- |
| `PORT` | 3000 | 3001 | 3002 | distinct on same host |
| `SQL_DSN` | ...oneapi_node1 | ...oneapi_node2 | ...oneapi_node3 | independent MySQL |
| `REDIS_CONN_STRING` | :6379/0 | :6380/0 | :6381/0 | independent Redis |
| `CLUSTER_NODE_ID` | 1 | 2 | 3 | matches `auto_increment_offset` |
| `CLUSTER_NODE_NAME` | node-cn | node-us | node-eu | display name |
| `CLUSTER_NODE_ADDRESS` | https://cn.example.com | https://us.example.com | https://eu.example.com | public URL |
| `CLUSTER_SECRET` | same | same | same | identical across all nodes |
| `CLUSTER_SEEDS` | own address or empty | any alive node | any alive node | bootstraps discovery |

**5. Start commands**

Each node loads its own config via `--env`:

```bash
# Node 1
./one-api-pro --env /opt/one-api-pro/node1/.env --port 3000

# Node 2
./one-api-pro --env /opt/one-api-pro/node2/.env --port 3001

# Node 3
./one-api-pro --env /opt/one-api-pro/node3/.env --port 3002
```

**6. Startup order**

1. Start the first node (Node A); leave `CLUSTER_SEEDS` empty or use its own address.
2. Wait for Node A to fully start (~5–10 s; look for `cluster module initialized` in the logs).
3. Start the remaining nodes, pointing `CLUSTER_SEEDS` at any alive node.
4. New nodes ping the seed, then transitively learn all the rest.
5. Once every node is up, view the cluster under **Settings → Node Management** in the admin console.

**7. Optional Nginx load-balancing example**

```nginx
upstream one_api_cluster {
    ip_hash;  # same client always hits the same node — preserves session & cache
    server cn.example.com:3000;
    server us.example.com:3000;
    server eu.example.com:3000;
}

server {
    listen 443 ssl;
    server_name api.example.com;

    location / {
        proxy_pass http://one_api_cluster;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 300s;
    }
}
```

> **`ip_hash` is critical** — it pins a client to a single node so plan rate-limits and Redis cache stay consistent.

**8. Verify cluster status**

```bash
# List nodes (call any node)
curl -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  https://cn.example.com/api/cluster_node/
# Returns every node with status, last_heartbeat, ping_failures, …
```

Or open the admin console → **Settings → Node Management** to inspect nodes, status and last heartbeats.

> 💡 Cluster API details: [docs/API.md Appendix E — Cluster Management API](../docs/API.md#appendix-e-cluster-management-api)

#### ⚠️ Operational Notes

- Each node must have its own MySQL & Redis instances.
- `CLUSTER_SECRET` must be identical across all nodes. Use a strong value and keep it safe.
- `CLUSTER_NODE_ID` must be unique across nodes and match MySQL's `auto_increment_offset`.
- `CLUSTER_NODE_ADDRESS` must be reachable by other nodes (include protocol prefix, e.g. `https://`).
- New nodes must be bootstrapped manually (pull a snapshot from a live node).
- Logs are large by nature — disable with `CLUSTER_SYNC_LOGS=false` if needed.
- MySQL's `auto_increment_increment` & `auto_increment_offset` must match `CLUSTER_NODE_ID`.
- Discovery uses a bi-directional ping; failed nodes are not deleted, only marked `status=2`. They auto-resurrect when the network is back.
- `CLUSTER_SEEDS` is only the bootstrap. Once peers are discovered, SEEDS is no longer consulted.
- Changes that happen while a node is offline are **not back-filled** automatically. The node must pull a snapshot after coming back.

#### 📝 Why the "self-registration" of nodes?

Each node writes a self-record (`node_id == CLUSTER_NODE_ID`) into its own `cluster_nodes` table on startup. This is **intentional**:

1. **Admin visibility** — Settings → Node Management should display local info (address, status, heartbeat) for troubleshooting.
2. **Transitive discovery** — when Node B receives a ping from Node A, A returns its full node list (including itself). B merges it. So C can learn about A through B.
3. **Liveness signal** — local `last_heartbeat` is refreshed every 30 s by `discoverOnce`, reflecting this node's own health.

The design prevents loops with five layers of guards:

| Guard | Effect |
| --- | --- |
| ① `GetAllRemoteNodes` SQL filter | `WHERE node_id != ?` excludes self during discovery |
| ② `GetAliveNodesForSync` SQL filter | `WHERE node_id != ?` excludes self during push |
| ③ `handlePing` rejects self-ping | `req.NodeId == NodeID` is refused explicitly |
| ④ `mergeDiscoveredNodes` skips self | local merge skips self |
| ⑤ `ApplyEvents` skips self events | receiver ignores events produced by itself |

Data flow is one-way: local → remote, remote → local — **no loops**.

The admin console displays a blue "self" badge next to the local node and disables **Delete** / **Manual Ping** on it (they have no meaning for the local node).

#### 🔐 Per-node Secrets

Each node carries **its own secret**, replacing the previous global-shared-secret design. Reasons:

1. **Security** — one node's secret leaking does not affect others.
2. **Flexibility** — every node can rotate its own secret independently.
3. **Auto-discovery** — peers share their secret on every ping.

**Secret lifecycle:**
- First start: uses `CLUSTER_SECRET` env var as the initial value, written to `cluster_nodes.secret_key`.
- Subsequent starts: read from `cluster_nodes.secret_key`.
- Admins can update any other node's secret from the Node Management page.
- `X-Cluster-Secret` header on ping = **target node's** secret (looked up locally).

**Adding a new node:**
1. On Node A, add Node B's record using B's `CLUSTER_SECRET`.
2. On Node B, add Node A's record using A's `CLUSTER_SECRET`.
3. A → B ping uses B's secret; B verifies with its own secret ✓.
4. B's response carries both A's and B's secrets — A updates its local copy.

#### 🗑️ "Soft-delete" for Nodes

Removing a node sets `disabled = true` instead of physically deleting the record:

- Prevents the deleted node from "regrowing" through ping-driven re-registration.
- Disabled nodes still respond to pings (so peers know they're online) but won't fetch this node's info.
- Hard delete requires manual SQL: `DELETE FROM cluster_nodes WHERE node_id = ?`.

#### 🔄 Data Sync Mechanism (Important)

Cluster data sync relies entirely on **GORM callbacks + HTTP active push**:
- Every INSERT/UPDATE/DELETE on a business table → captured by the GORM callback → written to `sync_events` → Pusher goroutine pushes to every alive node.
- Receivers write through `WithSkipHook` (no loop-back).
- Receivers also skip events whose `event.NodeId == local NodeID` (belt-and-braces).

**Design trade-off:** this design **does not implement cross-node active pull**, for these reasons:
1. **Business-intrusive** — pull would require knowing each table's business-unique field, polluting business code.
2. **Primary-key conflicts** — auto-increment IDs differ across nodes (different `auto_increment_offset`); using the source's ID would break the offset design.
3. **Complexity** — maintenance cost is high, reliability gain is limited.
4. **Push is enough** — push covers ~95% of normal scenarios (alive nodes, normal traffic).

**Known limits & operational needs:**
- Changes made while a node is offline → **lost permanently** (push is realtime only).
- After a node comes back, missed data is not auto-recovered.
- New nodes only see changes from the moment they join — no history.
- **Operator's recourse:** `mysqldump` an existing node and restore.

**Scenario cheat-sheet:**

| Scenario | Pull needed? | How to handle |
| --- | --- | --- |
| Node permanently online | ❌ | push is enough |
| Node reboots occasionally (minutes) | ⚠️ | minor loss during downtime, usually acceptable |
| Node under frequent maintenance | ❌ | push keeps going; resume is instant |
| New node joins cluster | ❌ | DBA runs `mysqldump` to seed |
| Long offline node recovers | ❌ | DBA runs `mysqldump` to repair |

If a freshly deployed cluster shows a blank page, see [#97](https://github.com/modelbus/one-api-pro/issues/97).

---

## 🗺️ Roadmap

### ✅ Done

- [x] **Architectural refactor** — adaptor self-registration; adding a provider needs zero framework changes.
- [x] **Vue 3 admin console** — Arco Design + visual dashboard + 30+ provider icons.
- [x] **Subscription & plan system** — token- or request-based billing, period rate limits, per-model controls.
- [x] **Decentralized active-active cluster** — GORM-event-driven + HTTP active push; no shared DB.
- [x] **Precise cost accounting** — independent Prompt / Completion / Cached pricing, stacked group discounts.
- [x] **Multi-level permissions** — Guest / User / Admin / Root, upstream API authorization loophole fixed.
- [x] **OpenAI-compatible API** — `models` / `chat` / `completions` / `embeddings` / `images` / `audio` / `moderations`.
- [x] **Subscription checkout & upgrade flow** — native `POST /api/order/plan` creates the subscription order, supports `stack` (additive) and `price_diff` (upgrade) modes, auto-computes the upgrade price from remaining days, and rejects same-tier / downgrade attempts.
- [x] **Order audit & order center** — new `orders` table (type/source/order_no/plan_info/amount/status/pay_status/pay_method/pay_time/pay_trade_no) persists every checkout / admin grant. Frontend `/plans` and `/orders` pages render the full lifecycle.
- [x] **Real payment integration (gopay)** — native WeChat Pay Native (PC QR) and Alipay Face-to-Face (TradePrecreate); async callbacks at `/api/payment/{wechat,alipay}/notify` complete the verify → mark paid → activate loop.
- [x] **Payment / plan-operations settings** — two new sub-tabs under **Operations Settings**: `Plan` (price_diff vs stack upgrade toggle) and `Payment` (independent enable/disable + cert upload + notify URL per channel: WeChat / Alipay / Bank).

### 🔄 In Progress

- [ ] **Richer channel diagnostics & intelligent routing** — `CooldownFilter`, `FallbackFilter` and the `monitor`'s auto-disable on low success rate are already in place. Remaining: standalone channel-health dashboard, per-node ping endpoint, and a manual-review workflow.
- [ ] Richer usage analytics reports and exports.
- [ ] Improved i18n coverage.

### 🔭 Planned

- [ ] **More payment channels** — Apple Pay, UnionPay, Stripe, etc.; async refund API + automated refund ledger.
- [ ] **Online quota top-up** — users can top up their account balance from the personal area; this stays independent of subscription plans.
- [ ] **Finance-system integration** — sync top-ups, consumption and refunds with mainstream finance / reconciliation platforms.
- [ ] **Token-low alerting** — multi-channel notifications when account / token Token balance drops below threshold.
- [ ] **Audit logs & reports** — full operation audit trail and visual reports for compliance.
- [ ] **AI-powered analytics** — LLM-driven analysis of usage, costs and channel health, with recommendations.
- [ ] Plugin extension mechanism.
- [ ] Enterprise SSO / LDAP integration.
- [ ] Alerting & notification channels (DingTalk / Feishu / WeCom).
- [ ] More model providers.

> 💡 PRs and issues are welcome — see [issues](https://github.com/modelbus/one-api-pro/issues).

---

## 📄 License

This project is released under the [MIT License](../LICENSE).
