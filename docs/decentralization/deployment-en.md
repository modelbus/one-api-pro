---
title: Multi-node Deployment
description: "Production best practices for multi-node deployment."
category: decentralization
order: 5
---

# Multi-node Deployment

> Production best practices for multi-node deployment.

## Overall topology

```text
                ┌─────────────┐
                │  Nginx/LB   │   (single entry, ip_hash)
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

Each node is **fully independent**:

- Each runs its own MySQL instance.
- Each runs its own Redis instance.
- Each has its own `CLUSTER_NODE_ID`, `CLUSTER_NODE_ADDRESS` and independent secret.

## MySQL: one instance per node

> **Why**: `auto_increment_offset` is a MySQL **instance-level** variable; multiple databases on one instance cannot have different offsets.

Sample `my.cnf` (three-node example):

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

Highlights:

- `auto_increment_increment = 50` supports up to 50 nodes.
- Each node's `offset` equals `CLUSTER_NODE_ID`, unique across the cluster.
- `server-id` must be unique across all MySQL instances.
- `log_bin` + `binlog_format=ROW` recommended — enables future master/slave replication and point-in-time recovery.
- Cluster data sync does not depend on binlog (it uses GORM callbacks); binlog is just an extra safety net.

Pre-create the empty database:

```sql
CREATE DATABASE oneapi CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## Redis: one instance per node

Redis only serves **local cache + rate-limit state** here; cluster traffic does not flow through Redis. One instance per node is enough.

```bash
# Per-node startup
redis-server --port 6379 --appendonly yes
```

## Bootstrapping a new node

A newly added node must pull a snapshot from an existing node, otherwise it sees only changes from the moment it joins.

### Option 1: mysqldump import (recommended)

```bash
# Export from a live node
mysqldump -h <alive-node> -u root -p \
  --single-transaction --routines --triggers \
  --databases oneapi \
  > /opt/one-api-pro/seed/oneapi.sql

# Import into the new node (database must already exist)
mysql -u root -p oneapi < /opt/one-api-pro/seed/oneapi.sql
```

### Option 2: snapshot API

```bash
curl -H "X-Cluster-Secret: <target-node-secret>" \
  "https://<alive-node>/api/cluster/snapshot?tables=users,tokens,channels,abilities,options,redemptions,plans,user_plans,plan_usages" \
  -o /opt/one-api-pro/seed/snapshot.json
```

> Startup order: start the first node → wait for `cluster module initialized` in the logs → start the remaining nodes with `CLUSTER_SEEDS` pointing at the first one.

## Full .env examples

> All nodes share the same initial `CLUSTER_SECRET`; on first boot it is written to each node's `cluster_nodes.secret_key`. Subsequent rotation happens per-node via the admin UI or API.

### Node 1 — China

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

### Node 2 — US

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

### Node 3 — Europe

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

### Configuration cheat-sheet

| Variable | Node 1 | Node 2 | Node 3 | Notes |
| --- | --- | --- | --- | --- |
| `PORT` | 3000 | 3001 | 3002 | distinct on one host |
| `SQL_DSN` | `...oneapi_node1` | `...oneapi_node2` | `...oneapi_node3` | independent MySQL |
| `REDIS_CONN_STRING` | `:6379/0` | `:6380/0` | `:6381/0` | independent Redis |
| `CLUSTER_NODE_ID` | 1 | 2 | 3 | = `auto_increment_offset` |
| `CLUSTER_NODE_NAME` | `node-cn` | `node-us` | `node-eu` | display name |
| `CLUSTER_NODE_ADDRESS` | `https://cn.example.com` | `https://us.example.com` | `https://eu.example.com` | public URL |
| `CLUSTER_SECRET` | same | same | same | identical across cluster |
| `CLUSTER_SEEDS` | own or empty | any alive node | any alive node | bootstrap only |

## Startup

```bash
# Node 1
./one-api-pro --env /opt/one-api-pro/node1/.env --port 3000

# Node 2
./one-api-pro --env /opt/one-api-pro/node2/.env --port 3001

# Node 3
./one-api-pro --env /opt/one-api-pro/node3/.env --port 3002
```

Order:

1. Start the first node; leave `CLUSTER_SEEDS` empty or use its own address.
2. Wait ~5–10 seconds for `cluster module initialized` in the logs.
3. Start the remaining nodes with `CLUSTER_SEEDS` pointing at any alive node.
4. New nodes ping the seed and transitively learn the rest.
5. Once every node is up, verify the node list and heartbeats under **Settings → Node Management**.

## Nginx with ip_hash

```nginx
upstream one_api_cluster {
    ip_hash;  # pin client to one node
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

> `ip_hash` is critical — it pins a client to a single node so plan rate-limits and Redis cache stay consistent.

## Monitoring

```bash
# List all nodes (with heartbeat / failure counts)
curl -H "Authorization: Bearer YOUR_ROOT_TOKEN" \
  https://cn.example.com/api/cluster_node/
```

Or open the admin UI **Settings → Node Management** to inspect:

- `status`: `1` alive, `2` failed/disabled
- `last_heartbeat`: last liveness timestamp
- `ping_failures`: consecutive failure count

## Upgrade

Rolling upgrade is recommended — avoid restarting every node at once.

```text
1. Upgrade Node A (keep B / C running)
2. Wait until A is healthy and A↔B / A↔C heartbeats recover
3. Upgrade Node B
4. ...
```

Decentralized cluster sync is **realtime push**; changes made while a node is offline are not auto-replayed. After a node comes back, run `mysqldump` from a live node to re-seed and converge.

## Disaster recovery

| Scenario | Action |
| --- | --- |
| Brief node offline | Auto-recovery (`status=1`); Pusher replays queued events |
| Long offline → drift | `mysqldump` from a live node, re-import, restart |
| Whole datacenter down | Switch DNS to another datacenter; rebuild from `mysqldump` |
| Secret leaked | Update via admin UI `PUT /api/cluster_node/`; next ping auto-propagates |

## Operational notes

- `CLUSTER_SECRET` must be identical across all nodes on first boot; afterwards each node rotates its own `secret_key` independently.
- `CLUSTER_NODE_ID` must be unique and match MySQL `auto_increment_offset`.
- `CLUSTER_NODE_ADDRESS` must be reachable by other nodes (include the scheme, e.g. `https://`).
- Failed nodes are **not** auto-deleted — only marked `status=2`; they auto-resurrect when the network is back.
- `CLUSTER_SEEDS` is bootstrap-only; it is not consulted once peers are discovered.
- Disable `CLUSTER_SYNC_LOGS` if the `logs` table grows too large.
- Offline changes are not auto-replayed — re-seed via `mysqldump` after recovery.

Next: [Node Management](/en/decentralization/node-management) · [Node Health](/en/decentralization/node-health).
