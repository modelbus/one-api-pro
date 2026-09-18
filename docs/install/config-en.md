---
title: Configuration
description: "Every configuration option and its default."
category: install
order: 5
---

# Configuration

> Every configuration option and its default.

Precedence (high → low): **CLI flags** &gt; **process env** &gt; **`--env` file** &gt; **defaults**.

## Basics

| Variable | Default | Description |
| --- | --- | --- |
| `PORT` | `3000` | HTTP listen port; also consumed by the image `HEALTHCHECK` |
| `DEBUG` | `false` | Set to `true` to enable Gin debug logging |
| `DEBUG_SQL` | `false` | Print GORM SQL statements |
| `TZ` | `UTC` | Timezone; set e.g. `Asia/Shanghai` |
| `THEME` | `default-pro` | Web theme; only `default-pro` is actively maintained |

## Database

| Variable | Default | Description |
| --- | --- | --- |
| `SQL_DSN` | unset → SQLite | MySQL: `root:pw@tcp(host:3306)/oneapi`; PostgreSQL: `postgres://user:pw@host:5432/oneapi` |
| `LOG_SQL_DSN` | unset → reuse primary | Dedicated DB for the `logs` table (MySQL/PostgreSQL recommended) |
| `SQLITE_PATH` | `one-api-pro.db` | SQLite file path; in the container `/app/data/one-api-pro.db` |
| `SQLITE_BUSY_TIMEOUT` | `3000` | SQLite lock-wait timeout in milliseconds |

Connection pool (effective only when `SQL_DSN` is set, per README): `SQL_MAX_IDLE_CONNS=100`, `SQL_MAX_OPEN_CONNS=1000`, `SQL_CONN_MAX_LIFETIME=60` (minutes).

## Redis / Cache

| Variable | Default | Description |
| --- | --- | --- |
| `REDIS_CONN_STRING` | unset → disabled | `redis://default:pw@host:6379`; for Sentinel/Cluster use a comma-separated node list |
| `REDIS_PASSWORD` | unset | Required for Sentinel / Cluster mode |
| `REDIS_MASTER_NAME` | unset | Required for Sentinel mode |
| `MEMORY_CACHE_ENABLED` | `false` | In-process cache; introduces short-term staleness |
| `SYNC_FREQUENCY` | `600` | Interval (seconds) for syncing config from DB into memory |

## Session & Rate Limits

| Variable | Default | Description |
| --- | --- | --- |
| `SESSION_SECRET` | random UUID on first boot | Must be identical across every node/instance for shared sessions |
| `GLOBAL_API_RATE_LIMIT` | `480` | Global API rate limit (per IP, 3-minute window; older README revisions documented `180`) |
| `GLOBAL_WEB_RATE_LIMIT` | `240` | Global Web rate limit (per IP, 3-minute window; older README revisions documented `60`) |
| `FRONTEND_BASE_URL` | unset | When set on a slave node, page requests are redirected here |
| `NODE_TYPE` | `master` | `master` or `slave`; used by the multi-instance-shared-DB pattern |

## Channels & Polling

| Variable | Default | Description |
| --- | --- | --- |
| `CHANNEL_TEST_FREQUENCY` | unset (disabled) | Channel health-check interval in minutes |
| `CHANNEL_UPDATE_FREQUENCY` | unset (disabled) | Channel-balance refresh interval in minutes |
| `CHANNEL_DEFAULT_COOLDOWN_SECONDS` | `60` | Default channel cooldown in seconds |
| `CHANNEL_MAX_COOLDOWN_SECONDS` | `600` | Max channel cooldown in seconds |
| `CHANNEL_CONCURRENCY_ENABLED` | `false` | Enable per-channel concurrency limit |
| `CHANNEL_STICKY_SESSION_ENABLED` | `false` | Enable sticky session on channels |
| `POLLING_INTERVAL` | `0` (no spacing) | Spacing (seconds) between batched channel balance/health requests |
| `BATCH_UPDATE_ENABLED` | `false` | Aggregate user-quota writes in batch |
| `BATCH_UPDATE_INTERVAL` | `5` | Batch aggregation window in seconds |

## Relay & Proxy

| Variable | Default | Description |
| --- | --- | --- |
| `RELAY_TIMEOUT` | `0` (no timeout) | Upstream LLM request timeout in seconds |
| `RELAY_PROXY` | unset | Proxy URL for upstream LLM requests |
| `USER_CONTENT_REQUEST_TIMEOUT` | `30` | User-content (e.g. image) download timeout in seconds |
| `USER_CONTENT_REQUEST_PROXY` | unset | Proxy for user-content downloads |
| `ENFORCE_INCLUDE_USAGE` | `false` | Force `usage` to be returned in stream responses |
| `TEST_PROMPT` | `Output only your specific model name with no additional text.` | Channel/model self-test prompt |

## Gemini / Model Adaptors

| Variable | Default | Description |
| --- | --- | --- |
| `GEMINI_SAFETY_SETTING` | `BLOCK_NONE` | Gemini safety level |
| `GEMINI_VERSION` | `v1` | Gemini API version |

## Tokenizer Cache

Since v0.0.21 the tiktoken BPE encodings are embedded by default; the directories below are only for advanced overrides.

| Variable | Default | Description |
| --- | --- | --- |
| `TIKTOKEN_CACHE_DIR` | unset | Explicit BPE directory overrides embedded + HTTP fallback |
| `DATA_GYM_CACHE_DIR` | unset | Same as `TIKTOKEN_CACHE_DIR`, lower priority |

## Metric Monitor

| Variable | Default | Description |
| --- | --- | --- |
| `ENABLE_METRIC` | `false` | Auto-disable channels below the success-rate threshold |
| `METRIC_QUEUE_SIZE` | `10` | Sliding-window size |
| `METRIC_SUCCESS_RATE_THRESHOLD` | `0.8` | Success-rate threshold |
| `METRIC_SUCCESS_CHAN_SIZE` | `1024` | Success event channel capacity |
| `METRIC_FAIL_CHAN_SIZE` | `128` | Failure event channel capacity |
| `ONLY_ONE_LOG_FILE` | `false` | Keep only a single log file |

## Initial Root Tokens

If set on first startup, the corresponding root tokens are auto-created.

| Variable | Default | Description |
| --- | --- | --- |
| `INITIAL_ROOT_TOKEN` | unset | Auto-created root API token on first boot |
| `INITIAL_ROOT_ACCESS_TOKEN` | unset | Auto-created root system-management token on first boot |

## Cluster Mode

| Variable | Default | Description |
| --- | --- | --- |
| `CLUSTER_ENABLED` | `false` | Enable decentralized cluster mode |
| `CLUSTER_NODE_ID` | required | Node ID `1-49`; must match MySQL `auto_increment_offset` |
| `CLUSTER_NODE_NAME` | `node-<ID>` | Display name |
| `CLUSTER_NODE_ADDRESS` | required | Public URL other nodes use to reach this node (include scheme, e.g. `https://cn.example.com`) |
| `CLUSTER_SECRET` | required | Node-to-node auth secret (seed on first boot; rotatable from the admin UI) |
| `CLUSTER_SEEDS` | unset | Comma-separated seed nodes; one reachable node is enough |
| `CLUSTER_PUSH_INTERVAL` | `3` | Sync event push interval (seconds) |
| `CLUSTER_DISCOVERY_INTERVAL` | `30` | Node discovery interval (seconds) |
| `CLUSTER_DEAD_PING_INTERVAL` | `120` | Ping interval for unreachable nodes (seconds) |
| `CLUSTER_MAX_PING_FAILURES` | `3` | Consecutive failure threshold |
| `CLUSTER_SYNC_LOGS` | `true` | Sync the `logs` table; disable if it grows too large |
| `CLUSTER_BATCH_SIZE` | `50` | Max events per push |

See [Cluster Overview](/en/decentralization/overview) for cluster-mode details.

## CLI Flags

```text
--port <port_number>         # listen port, default 3000
--log-dir <log_dir>          # log directory, default ./logs
--env <env_file_path>        # .env file path; auto-loads ./.env when unset
--version                    # print version and exit
--help                       # print help
```

Multi-instance on one host:

```bash
./one-api-pro --env ./instances/instance1.env --port 3001 &
./one-api-pro --env ./instances/instance2.env --port 3002 &
```

## .env Loading

- Binary: load with `--env <path>`; auto-loads `./.env` in CWD.
- Docker: `docker-entrypoint.sh` auto-loads `$CONFIG_DIR/.env` (default `/app/config/.env`).

> A few variables (e.g. `GLOBAL_API_RATE_LIMIT`, `GLOBAL_WEB_RATE_LIMIT`) used to have lower defaults in older README revisions (`180` / `60`); the current source defaults are higher (`480` / `240`). Use `common/config/config.go` and this table as the source of truth.

Next: [Upgrade](/en/install/upgrade) · [Backup & Restore](/en/install/backup-restore).
