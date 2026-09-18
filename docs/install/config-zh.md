---
title: 配置项与环境变量
description: "所有可配置项及其默认值说明。"
category: install
order: 5
---

# 配置项与环境变量

> 所有可配置项及其默认值说明。

优先级（高 → 低）/ Precedence (high → low): **CLI 参数 / CLI flags** &gt; **进程环境变量 / process env** &gt; **`--env` 文件 / `--env` file** &gt; **默认值 / defaults**。

## 基础

| 变量| 默认| 说明|
| --- | --- | --- |
| `PORT` | `3000` | HTTP 监听端口；同时被镜像 `HEALTHCHECK` 读取 |
| `DEBUG` | `false` | 设为 `true` 开启 Gin debug 日志 |
| `DEBUG_SQL` | `false` | 打印 GORM SQL |
| `TZ` | `UTC` | 时区；推荐 `Asia/Shanghai` 等本地值 |
| `THEME` | `default-pro` | 主题；可选（仅维护 `default-pro`，其他为历史保留） |

## 数据库

| 变量| 默认| 说明|
| --- | --- | --- |
| `SQL_DSN` | 空| MySQL: `root:pw@tcp(host:3306)/oneapi`；PostgreSQL: `postgres://user:pw@host:5432/oneapi` |
| `LOG_SQL_DSN` | 空| 为 `logs` 表指定独立数据库（推荐 MySQL/PostgreSQL） |
| `SQLITE_PATH` | `one-api-pro.db` | SQLite 文件路径；容器默认 `/app/data/one-api-pro.db` |
| `SQLITE_BUSY_TIMEOUT` | `3000` | SQLite 锁等待超时（毫秒） |

连接池（仅当 `SQL_DSN` 非空时生效，参考 README）：`SQL_MAX_IDLE_CONNS=100`、`SQL_MAX_OPEN_CONNS=1000`、`SQL_CONN_MAX_LIFETIME=60`（分钟）。

## Redis

| 变量| 默认| 说明|
| --- | --- | --- |
| `REDIS_CONN_STRING` | 空| `redis://default:pw@host:6379`；Sentinel/Cluster 模式用逗号分隔多个节点 |
| `REDIS_PASSWORD` | 空 | Sentinel|
| `REDIS_MASTER_NAME` | 空 | Sentinel 模式专用 |
| `MEMORY_CACHE_ENABLED` | `false` | 进程内缓存；启用会带来短期数据陈旧 |
| `SYNC_FREQUENCY` | `600` | 配置从 DB 同步进内存的间隔（秒） |

## 会话与限流

| 变量| 默认| 说明|
| --- | --- | --- |
| `SESSION_SECRET` | 随机 UUID（首次启动） | 所有集群/多实例节点必须一致，否则会话不互通 |
| `GLOBAL_API_RATE_LIMIT` | `480` | 全局 API 限流（每 IP，3 分钟窗口，README 旧版本曾为 180） |
| `GLOBAL_WEB_RATE_LIMIT` | `240` | 全局 Web 限流（每 IP，3 分钟窗口，README 旧版本曾为 60） |
| `FRONTEND_BASE_URL` | 空 | 从节点设置后，页面请求会被重定向到该地址 |
| `NODE_TYPE` | `master` | `master` 或 `slave`；用于多实例共享 DB 模式 |

## 渠道与轮询

| 变量| 默认| 说明|
| --- | --- | --- |
| `CHANNEL_TEST_FREQUENCY` | 空（不测） | 渠道可用性定期测试（分钟） |
| `CHANNEL_UPDATE_FREQUENCY` | 空（不更） | 渠道余额定期更新（分钟） |
| `CHANNEL_DEFAULT_COOLDOWN_SECONDS` | `60` | 渠道默认冷却时间（秒） |
| `CHANNEL_MAX_COOLDOWN_SECONDS` | `600` | 渠道最大冷却时间（秒） |
| `CHANNEL_CONCURRENCY_ENABLED` | `false` | 是否启用渠道并发限制 |
| `CHANNEL_STICKY_SESSION_ENABLED` | `false` | 是否启用渠道粘性会话 |
| `POLLING_INTERVAL` | `0`（无间隔） | 批量更新渠道余额/可用性请求之间的间隔（秒） |
| `BATCH_UPDATE_ENABLED` | `false` | 数据库批量更新聚合 |
| `BATCH_UPDATE_INTERVAL` | `5` | 聚合窗口（秒） |

## 转发与代理

| 变量| 默认| 说明|
| --- | --- | --- |
| `RELAY_TIMEOUT` | `0`（无超时） | 上游 LLM 请求超时（秒） |
| `RELAY_PROXY` | 空 | 上游 LLM 请求的代理 URL |
| `USER_CONTENT_REQUEST_TIMEOUT` | `30` | 用户内容（如图片）下载超时（秒） |
| `USER_CONTENT_REQUEST_PROXY` | 空 | 用户内容下载的代理 URL |
| `ENFORCE_INCLUDE_USAGE` | `false` | 强制 stream 模式返回 `usage` |
| `TEST_PROMPT` | `Output only your specific model name with no additional text.` | 渠道/模型自测 prompt |

## Gemini

| 变量| 默认| 说明|
| --- | --- | --- |
| `GEMINI_SAFETY_SETTING` | `BLOCK_NONE` | Gemini 安全等级 |
| `GEMINI_VERSION` | `v1` | Gemini API 版本 |

## Tokenizer 缓存

自 v0.0.21 起，tiktoken BPE 编码已默认内嵌进二进制，启动零网络依赖；下述目录仅在高级用户希望替换编码时使用。

| 变量| 默认| 说明|
| --- | --- | --- |
| `TIKTOKEN_CACHE_DIR` | 空 | 显式提供 BPE 文件目录时优先于此目录 |
| `DATA_GYM_CACHE_DIR` | 空 | 与 `TIKTOKEN_CACHE_DIR` 同义，优先级更低 |

## 指标监控

| 变量| 默认| 说明|
| --- | --- | --- |
| `ENABLE_METRIC` | `false` | 启用自动下线低成功率渠道 |
| `METRIC_QUEUE_SIZE` | `10` | 滑窗大小 |
| `METRIC_SUCCESS_RATE_THRESHOLD` | `0.8` | 成功率阈值 |
| `METRIC_SUCCESS_CHAN_SIZE` | `1024` | 成功事件通道容量 |
| `METRIC_FAIL_CHAN_SIZE` | `128` | 失败事件通道容量 |
| `ONLY_ONE_LOG_FILE` | `false` | 是否只保留一个日志文件 |

## 初始化令牌

首次启动时若设置了以下两个变量，会自动创建对应 root 用户的令牌/管理令牌；留空则保持默认。

| 变量| 默认| 说明|
| --- | --- | --- |
| `INITIAL_ROOT_TOKEN` | 空 | 首次启动自动创建的 root 用户 API token |
| `INITIAL_ROOT_ACCESS_TOKEN` | 空 | 首次启动自动创建的 root 系统管理 token |

## 集群模式

| 变量| 默认| 说明|
| --- | --- | --- |
| `CLUSTER_ENABLED` | `false` | 启用去中心化集群模式 |
| `CLUSTER_NODE_ID` | 必填| 节点编号 `1-49`，必须与 MySQL `auto_increment_offset` 一致 |
| `CLUSTER_NODE_NAME` | `node-<ID>` | 节点显示名 |
| `CLUSTER_NODE_ADDRESS` | 必填| 本节点公网访问地址（含协议前缀，如 `https://cn.example.com`） |
| `CLUSTER_SECRET` | 必填| 节点间通信认证密钥（首次启动种子，可后续在 UI 轮换） |
| `CLUSTER_SEEDS` | 空 | 逗号分隔种子节点；至少一个可达即可引导发现 |
| `CLUSTER_PUSH_INTERVAL` | `3` | 同步事件推送间隔（秒） |
| `CLUSTER_DISCOVERY_INTERVAL` | `30` | 节点发现间隔（秒） |
| `CLUSTER_DEAD_PING_INTERVAL` | `120` | 失败节点 ping 间隔（秒） |
| `CLUSTER_MAX_PING_FAILURES` | `3` | 连续失败次数阈值 |
| `CLUSTER_SYNC_LOGS` | `true` | 是否同步 `logs` 表；数据量大时可关闭 |
| `CLUSTER_BATCH_SIZE` | `50` | 每次推送最大事件数 |

集群模式细节见 [Cluster 概览](/zh/decentralization/overview)。
See [Cluster Overview](/en/decentralization/overview) for cluster-mode details.

## CLI 参数

```text
--port <port_number>         # 监听端口，默认 3000
--log-dir <log_dir>          # 日志目录，默认 ./logs
--env <env_file_path>        # 配置 .env 路径；不传则自动加载 ./.env
--version                    # 打印版本号并退出
--help                       # 打印帮助
```

多实例同机示例：

```

## .env 加载时机

- 源码启动：`--env <path>` 显式指定；不指定时自动加载当前目录下的 `./.env`。
  Binary: load with `--env <path>`; auto-loads `./.env` in CWD.
- Docker 启动：`docker-entrypoint.sh` 自动加载 `$CONFIG_DIR/.env`（默认 `/app/config/.env`）。
  Docker: `docker-entrypoint.sh` auto-loads `$CONFIG_DIR/.env` (default `/app/config/.env`).

> 部分变量（如 `GLOBAL_API_RATE_LIMIT`、`GLOBAL_WEB_RATE_LIMIT`）README 历史版本记录了更小的默认值（`180` / `60`），源码当前默认已上调（`480` / `240`），请以仓库 `common/config/config.go` 与本表为准。

下一步 / Next: [版本升级](/zh/install/upgrade) · [备份与恢复](/zh/install/backup-restore)。

