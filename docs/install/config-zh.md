---
title: 配置项与环境变量
description: 所有可配置项及其含义、默认值与影响。
category: install
order: 5
---

# 配置项与环境变量

> 改某个配置前要看本页：「这是什么、从哪里来、设置后会改变什么」。

## 优先级

```
CLI 参数 > 进程环境变量 > --env 文件 > 默认值
```

## 配置位置

- **环境变量**：在 `docker-compose.yml` 的 `environment` 或 `.env` 文件里设置
- **系统设置**：后台 → 系统设置，部分常用配置存数据库（如注册开关、邮件签名）
- **配置文件**：通过 `--config /path/to/config.yaml` 加载（高级用法）

## 数据库与缓存

### `SQL_DSN`

- **是什么**：MySQL 连接串。留空则使用内嵌 SQLite
- **格式**：`user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local`
- **设置后影响**：启用 MySQL；不设置就用 SQLite

### `REDIS_CONN_STRING`

- **是什么**：Redis 连接串
- **格式**：`redis://host:port/db`
- **设置后影响**：启用 Redis。多节点 Cluster 模式**必须**设置

## 服务监听

### `PORT`

- **是什么**：HTTP 主端口
- **默认值**：`3000`
- **设置后影响**：影响 `docker run -p` 映射、反向代理 upstream

### `LISTEN`

- **是什么**：监听地址
- **默认值**：`0.0.0.0`（全部接口）
- **设置后影响**：默认就 OK；反向代理场景无需修改

## 时区与时间

### `TZ`

- **是什么**：时区
- **推荐**：`Asia/Shanghai`（或其他你的业务时区）
- **设置后影响**：日志时间戳、CRON 任务、日切

## 安全与限流

### `SESSION_SECRET`

- **是什么**：Cookie / Session 签名密钥
- **默认值**：自动生成（重启后变化，会导致所有用户被登出）
- **建议**：**生产必填**，设一个长随机字符串
- **设置后影响**：所有现有 Session 失效

### `JWT_SECRET`

- **是什么**：管理后台某些接口用的 JWT 密钥
- **建议**：**生产必填**

### `RATE_LIMIT_*`

- **是什么**：全局限流参数
- **默认值**：激进，足够日常
- **设置后影响**：超出阈值直接 429

## 集群（Cluster）

仅多节点部署需要：

- `CLUSTER_NODE_ID`：全集群唯一（每节点不同）
- `CLUSTER_NODE_SECRET`：所有节点相同
- `CLUSTER_NODE_PORT`：节点间通信端口，默认 3000
- `CLUSTER_DISCOVERY_INTERVAL`：心跳 / 互 ping 周期（秒），默认 30
- `CLUSTER_DEAD_PING_INTERVAL`：失败节点重试间隔，默认 120
- `CLUSTER_MAX_PING_FAILURES`：连续失败多少次后自动禁用节点，默认 3

详见 [去中心化多节点](../decentralization/deployment)。

## 邮件（SMTP）

- `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS`, `SMTP_FROM`
- 不设置 → 关闭邮件发送（密码重置 / 兑换码 / 订单通知全部失效）

## 日志

- `LOG_CONSOLE_OUTPUT`：是否输出到 stdout（Docker 必备）
- `LOG_CONSOLE_LEVEL`：`debug` / `info` / `warn` / `error`
- `LOG_DB_RETENTION_DAYS`：调用日志保留天数（默认 90）

## Tiktoken（token 计数）

- `TIKTOKEN_CACHE_DIR`：覆盖内嵌的 BPE 文件路径。默认不用管；自定义编码时设置
- 内嵌 4 个 BPE 文件（cl100k_base / o200k_base / p50k_base / r50k_base），无需网络

## 完整列表

CLI 参数通过 `./one-api-pro --help` 查看完整列表。

## 相关文档

- [Docker 单实例部署](./docker-deploy)
- [docker-compose 部署](./docker-compose)
- [去中心化多节点](../decentralization/deployment)