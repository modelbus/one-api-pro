---
title: 多节点部署
description: "生产环境多节点部署最佳实践。"
category: decentralization
order: 5
---

# 多节点部署

> 生产环境多节点部署最佳实践。

## 总体拓扑

```text
                ┌─────────────┐
                │  Nginx/LB   │   (单一入口, ip_hash)
                └──────┬──────┘
                       │
       ┌───────────────┼───────────────┐
       │               │               │
 ┌─────┴─────┐   ┌─────┴─────┐   ┌─────┴─────┐
 │  Node 1   │   │  Node 2   │   │  Node 3   │
 │ one-api   │   │ one-api   │   │ one-api   │
 │ + MySQL   │   │ + MySQL   │   │ + MySQL   │
 │ + Redis   │   │ + Redis   │   │ + Redis   │
 └─────┬─────┘   └─────┬─────┘   └─────┬─────┘
       │               │               │
       └────── HTTP push of sync events ──────┘
```

每个节点都是 **完全独立** 的：
Each node is **fully independent**:

- 各自运行一个 MySQL 实例；
- 各自运行一个 Redis 实例；
- 自己的 `CLUSTER_NODE_ID`、`CLUSTER_NODE_ADDRESS`、独立 secret。

## MySQL：每节点独立实例

> **Why**: `auto_increment_offset` 是 MySQL **实例级**变量，多个数据库共用同一个实例时无法分别为它们设置不同 offset。

`my.cnf` 模板 / Sample `my.cnf` (three-node example):

```

要点 / Highlights:

- `auto_increment_increment = 50` — 支持最多 50 个节点；
  `auto_increment_increment = 50` supports up to 50 nodes.
- 每个节点的 `offset` = `CLUSTER_NODE_ID`，全集群唯一；
  Each node's `offset` equals `CLUSTER_NODE_ID`, unique across the cluster.
- `server-id` 必须唯一；
  `server-id` must be unique across all MySQL instances.
- `log_bin` + `binlog_format=ROW` 推荐开启 — 便于未来主从复制与 PITR；
  `log_bin` + `binlog_format=ROW` recommended — enables future master/slave replication and point-in-time recovery.
- 集群数据同步不依赖 binlog（走 GORM 回调），binlog 只作额外保险。
  Cluster data sync does not depend on binlog (it uses GORM callbacks); binlog is just an extra safety net.

预创建空库 / Pre-create the empty database:

```

## Redis：每节点独立实例

Redis 在此仅作 **本地缓存 + 限流**，不承载集群流量；每节点独立即可。

```bash
# 每节点各自启动
redis-server --port 6379 --appendonly yes
```

## 引导新节点

新加入节点必须从已有节点同步一次快照，否则从 0 启动后只能看到加入之后的变更。

### 方法 1：mysqldump 导入（推荐）/ Option 1: mysqldump import (recommended)

```bash
# 从一个存活节点导出
mysqldump -h <alive-node> -u root -p \
  --single-transaction --routines --triggers \
  --databases oneapi \
  > /opt/one-api-pro/seed/oneapi.sql

# 导入到新节点（数据库必须已 CREATE DATABASE）
mysql -u root -p oneapi < /opt/one-api-pro/seed/oneapi.sql
```

### 方法 2：快照 API

```bash
curl -H "X-Cluster-Secret: <target-node-secret>" \
  "https://<alive-node>/api/cluster/snapshot?tables=users,tokens,channels,abilities,options,redemptions,plans,user_plans,plan_usages" \
  -o /opt/one-api-pro/seed/snapshot.json
```

> 启动顺序：先启第一个节点 → 等到日志打印 `cluster module initialized` → 再启其他节点（其 `CLUSTER_SEEDS` 指向第一个节点）。

## .env 全模板

> 所有节点共用 `CLUSTER_SECRET` 初始值；secret 在首次启动时写入各自的 `cluster_nodes.secret_key`，之后由管理后台或 API 独立轮换。

### 节点 1 — China

```bash
# /opt/one-api-pro/node1/.env
PORT=3000
SYSTEM_NAME=One Api Pro Cluster

# Database (independent MySQL)
SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node1?charset=utf8mb4&parseTime=True&loc=Local

# Redis (independent)
REDIS_CONN_STRING=redis://127.0.0.1:6379/0

# Cluster
CLUSTER_ENABLED=true
CLUSTER_NODE_ID=1
CLUSTER_NODE_NAME=node-cn
CLUSTER_NODE_ADDRESS=https://cn.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me

# Seed nodes (only needed for first-time discovery)
# First node: leave empty or use its own address
CLUSTER_SEEDS=https://cn.example.com,https://us.example.com,https://eu.example.com

# Cluster tuning (optional)
CLUSTER_DISCOVERY_INTERVAL=30
CLUSTER_DEAD_PING_INTERVAL=120
CLUSTER_MAX_PING_FAILURES=3
CLUSTER_PUSH_INTERVAL=3
CLUSTER_SYNC_LOGS=true
CLUSTER_BATCH_SIZE=50
```

### 节点 2 — US

```bash
# /opt/one-api-pro/node2/.env
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

### 节点 3 — Europe

```bash
# /opt/one-api-pro/node3/.env
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

### 配置速查

| 变量| Node 1 | Node 2 | Node 3 | 备注|
| --- | --- | --- | --- | --- |
| `PORT` | 3000 | 3001 | 3002 | 同机时必须不同|
| `SQL_DSN` | `...oneapi_node1` | `...oneapi_node2` | `...oneapi_node3` | 独立 MySQL|
| `REDIS_CONN_STRING` | `:6379/0` | `:6380/0` | `:6381/0` | 独立 Redis|
| `CLUSTER_NODE_ID` | 1 | 2 | 3 | = `auto_increment_offset` |
| `CLUSTER_NODE_NAME` | `node-cn` | `node-us` | `node-eu` | 显示名|
| `CLUSTER_NODE_ADDRESS` | `https://cn.example.com` | `https://us.example.com` | `https://eu.example.com` | 公网 URL |
| `CLUSTER_SECRET` | same | same | same | 全集群一致|
| `CLUSTER_SEEDS` | 自身或空| 任一可达节点| 任一可达节点| 仅用于首次发现|

## 启动

```bash
# Node 1
./one-api-pro --env /opt/one-api-pro/node1/.env --port 3000

# Node 2
./one-api-pro --env /opt/one-api-pro/node2/.env --port 3001

# Node 3
./one-api-pro --env /opt/one-api-pro/node3/.env --port 3002
```

启动顺序 / Order:

1. 启动第一个节点；其 `CLUSTER_SEEDS` 留空或填自身；
2. 等待 ~5–10s，直到日志打印 `cluster module initialized`；
3. 启动其他节点，`CLUSTER_SEEDS` 指向任一存活节点；
4. 新节点 ping seed，传递性发现其余节点；
5. 全部上线后，在管理后台 **Settings → Node Management** 校验节点列表与心跳。

## Nginx（ip_hash）/ Nginx with ip_hash

```nginx
upstream one_api_cluster {
    ip_hash;  # 同一客户端固定到同一节点 / pin client to one node
    server cn.example.com:3000;
    server us.example.com:3000;
    server eu.example.com:3000;
}

server {
    listen 443 ssl http2;
    server_name api.example.com;

    ssl_certificate     /etc/letsencrypt/live/api.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.example.com/privkey.pem;

    location / {
        proxy_pass         http://one_api_cluster;
        proxy_http_version 1.1;

        proxy_set_header Host              $host;
        proxy_set_header X-Real-IP         $remote_addr;
        proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        proxy_set_header Upgrade           $http_upgrade;
        proxy_set_header Connection        "upgrade";
        proxy_buffering off;
        proxy_cache    off;
        proxy_read_timeout 300s;
    }
}
```

> `ip_hash` 是关键 — 它把同一客户端绑定到同一节点，让套餐限流与 Redis 缓存命中一致。

## 监控

```bash
# 列出全部节点（含心跳
curl -H "Authorization: Bearer YOUR_ROOT_TOKEN" \
  https://cn.example.com/api/cluster_node/
```

或在管理后台 **Settings → Node Management** 查看：

- `status`：`1` 在线，`2` 失败/禁用
- `last_heartbeat`：最近一次存活时间
- `ping_failures`：自上次成功以来的失败计数

## 升级

集群升级推荐 **滚动** 进行，避免多节点同时重启造成事件丢失：

```text
1. 升级 Node A（保持 B/C 运行）
2. 等待 A 健康、A↔B、A↔C 心跳恢复
3. 升级 Node B
4. ...
```

去中心化集群同步是 **实时推送**，离线期间的变更不会自动回填；节点恢复后建议用 `mysqldump` 同步一次以收敛漂移。

## 灾备

| 场景| 操作|
| --- | --- |
| 节点短暂离线 | 自动恢复（`status=1`），离线期间事件由 Pusher 补推 |
| 节点长时间离线 → 数据漂移 | `mysqldump` 从存活节点导入，重启 |
| 整个机房故障 | 切换 DNS 到其他机房节点，重新 `mysqldump` 重建 |
| Secret 泄露 | 管理后台 `PUT /api/cluster_node/` 改 secret；下次 ping 自动同步 |

## 注意事项

- `CLUSTER_SECRET` 必须在所有节点保持一致（初始值一致即可；后续可独立轮换）；
- `CLUSTER_NODE_ID` 必须唯一且与 MySQL `auto_increment_offset` 一致；
- `CLUSTER_NODE_ADDRESS` 必须能被其他节点访问（含协议前缀，例如 `https://`）；
- 失败节点 **不会** 被自动删除，仅标记 `status=2`；恢复后自动上线。
- `CLUSTER_SEEDS` 仅用于首次发现；节点加入后不再被查询。
- `logs` 表数据量大时可关闭 `CLUSTER_SYNC_LOGS`；
- 离线期间的变更不会自动回填，恢复后建议 `mysqldump` 同步。

下一步 / Next: [节点管理](/zh/decentralization/node-management) · [节点健康](/zh/decentralization/node-health)。

