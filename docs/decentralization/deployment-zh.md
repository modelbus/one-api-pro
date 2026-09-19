---
title: 多节点部署
description: 在生产环境启用 Cluster 模式：环境变量、Nginx、灾备。
category: decentralization
order: 5
---

# 多节点部署

> 把 One API Pro 部署成 3 节点 Cluster 的完整步骤。

## 总体拓扑

```
                 ┌─────────────┐
                 │  Nginx/LB   │  ← 单一入口，ip_hash
                 └──────┬──────┘
                        │
        ┌───────────────┼───────────────┐
        │               │               │
   ┌────┴────┐    ┌────┴────┐    ┌────┴────┐
   │ Node 1  │    │ Node 2  │    │ Node 3  │
   │ one-api │    │ one-api │    │ one-api │
   │ MySQL   │    │ MySQL   │    │ MySQL   │
   │ Redis   │    │ Redis   │    │ Redis   │
   └────┬────┘    └────┬────┘    └────┬────┘
        │               │               │
        └──────── HTTP push 同步 ────────┘
```

每个节点完全独立：各自的 MySQL、Redis、地址。

## 前置条件

每个节点需要：

- Linux 服务器 + Docker / 二进制
- 独立的 MySQL 实例（**不要共享数据库**）
- 独立的 Redis 实例
- 公网域名（如 `cn.example.com`）+ TLS 证书
- 反向代理（Nginx / Caddy / Cloudflare）

## MySQL 配置

每节点一个独立的 MySQL 实例。`my.cnf` 关键项：

```ini
[mysqld]
server-id = <与 CLUSTER_NODE_ID 唯一>
auto_increment_increment = 50     # 支持最多 50 个节点
auto_increment_offset = <CLUSTER_NODE_ID>
log_bin = ON
binlog_format = ROW                # 推荐，未来主从复制用
```

预创建空库：

```sql
CREATE DATABASE oneapi_node1 CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## Redis 配置

每节点一个本地 Redis，仅做缓存与限流。无需集群：

```bash
redis-server --port 6379 --appendonly yes
```

## 加入现有节点

新节点**必须从存活节点同步一次数据库快照**，否则只能看到加入后的变更。

### 方法 1：mysqldump（推荐）

```bash
# 从存活节点导出
mysqldump -h cn.example.com -u root -p \
  --single-transaction --routines --triggers \
  oneapi > /tmp/oneapi.sql

# 导入到新节点
mysql -u root -p oneapi < /tmp/oneapi.sql
```

### 方法 2：快照 API

```bash
curl -H "X-Cluster-Secret: <target-secret>" \
  "https://cn.example.com/api/cluster/snapshot?tables=users,tokens,channels,abilities,options,redemptions,plans,user_plans,plan_usages" \
  -o snapshot.json
```

## 启动顺序

1. **先启第一个节点**：`CLUSTER_SEEDS` 留空或填自身
2. 等 ~5-10s，看到日志 `cluster module initialized`
3. **再启其他节点**：`CLUSTER_SEEDS` 指向任一存活节点
4. 新节点 ping seed，传递性发现其余节点
5. 后台 → 集群设置 → 节点管理 校验列表与心跳

## 环境变量（3 节点示例）

| 变量 | Node 1 | Node 2 | Node 3 |
|---|---|---|---|
| `PORT` | 3000 | 3001 | 3002 |
| `SQL_DSN` | `...oneapi_node1` | `...oneapi_node2` | `...oneapi_node3` |
| `REDIS_CONN_STRING` | `:6379/0` | `:6380/0` | `:6381/0` |
| `CLUSTER_ENABLED` | true | true | true |
| `CLUSTER_NODE_ID` | **1** | **2** | **3** |
| `CLUSTER_NODE_NAME` | `node-cn` | `node-us` | `node-eu` |
| `CLUSTER_NODE_ADDRESS` | `https://cn.example.com` | `https://us.example.com` | `https://eu.example.com` |
| `CLUSTER_SECRET` | 同一个值 | 同一个值 | 同一个值 |
| `CLUSTER_SEEDS` | 自身或空 | `https://cn.example.com` | `https://cn.example.com` |

**关键约束：**

- `CLUSTER_SECRET` 所有节点必须一致（首次启动时；后续可独立轮换）
- `CLUSTER_NODE_ID` 必须唯一，与 MySQL `auto_increment_offset` 一致
- `CLUSTER_NODE_ADDRESS` 必须是其他节点能访问的公网 URL（含 `https://`）
- `CLUSTER_SEEDS` 仅首次发现时使用；加入后不再查询

### Node 1 完整 .env

```bash
PORT=3000
SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node1?charset=utf8mb4&parseTime=True&loc=Local
REDIS_CONN_STRING=redis://127.0.0.1:6379/0

CLUSTER_ENABLED=true
CLUSTER_NODE_ID=1
CLUSTER_NODE_NAME=node-cn
CLUSTER_NODE_ADDRESS=https://cn.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me

# 可选调优
CLUSTER_DISCOVERY_INTERVAL=30
CLUSTER_DEAD_PING_INTERVAL=120
CLUSTER_MAX_PING_FAILURES=3
CLUSTER_PUSH_INTERVAL=3
CLUSTER_SYNC_LOGS=true
CLUSTER_BATCH_SIZE=50
```

Node 2 / 3 类似，只改 ID / Name / Address / SQL_DSN。

## Nginx（ip_hash 入口）

```nginx
upstream one_api_cluster {
    ip_hash;  # 同一客户端固定到同一节点，套餐限流 + Redis 缓存命中一致
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

`ip_hash` 让同一客户端始终落到同一节点，缓存与限流不会因为切换节点而失效。

## 滚动升级

集群升级推荐 **滚动** 进行（一次只重启一个节点）：

```
1. 升级 Node A（保持 B/C 运行）
2. 等 A 健康、A↔B、A↔C 心跳恢复
3. 升级 Node B
4. 升级 Node C
```

每步观察心跳和同步状态。离线期间事件由 Pusher 补推。

## 灾备

| 场景 | 操作 |
|---|---|
| 节点短暂离线 | 自动恢复（`status=1`），离线期间事件由 Pusher 补推 |
| 节点长时间离线导致漂移 | `mysqldump` 从存活节点导入，重启 |
| 整个机房故障 | 切换 DNS 到其他机房节点，重新 `mysqldump` 重建 |
| Secret 泄露 | 后台 → 节点管理 → 编辑 secret；下次 ping 自动同步 |

## 监控

后台 → 集群设置 → 节点管理 看每行：

- `status`：`1` 在线，`2` 失败/禁用
- `last_heartbeat`：最近一次存活时间
- `ping_failures`：自上次成功以来的失败计数

或调 API：

```bash
curl -H "Authorization: Bearer <root-token>" \
  https://cn.example.com/api/cluster_node/
```

## 相关

- [Cluster 概览](./overview)
- [节点管理](./node-management)
- [节点健康](./node-health)