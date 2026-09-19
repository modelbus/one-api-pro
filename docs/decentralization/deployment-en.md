---
title: Multi-node Deployment
description: "Production multi-node setup: env vars, Nginx, disaster recovery."
category: decentralization
order: 5
---

# Multi-node Deployment

> Complete steps to deploy One API Pro as a 3-node cluster.

## Topology

```
                 ┌─────────────┐
                 │  Nginx/LB   │  ← single entry, ip_hash
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
        └──────── HTTP push sync ────────┘
```

Each node is fully independent: its own MySQL, Redis, address.

## Prerequisites

Per node:

- Linux server + Docker / binary
- Independent MySQL instance (**do not share**)
- Independent Redis
- Public domain (e.g. `cn.example.com`) + TLS
- Reverse proxy (Nginx / Caddy / Cloudflare)

## MySQL config

Each node has its own MySQL. Key `my.cnf`:

```ini
[mysqld]
server-id = <unique per node, matches CLUSTER_NODE_ID>
auto_increment_increment = 50     # supports up to 50 nodes
auto_increment_offset = <CLUSTER_NODE_ID>
log_bin = ON
binlog_format = ROW                # recommended for future replication
```

Pre-create the empty database:

```sql
CREATE DATABASE oneapi_node1 CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## Redis config

Each node has its own local Redis (cache + rate limiting only):

```bash
redis-server --port 6379 --appendonly yes
```

## Joining an existing cluster

A new node **must import a DB snapshot from a live peer first**, otherwise it only sees changes after it joins.

### Option 1: mysqldump (recommended)

```bash
# Export from a live node
mysqldump -h cn.example.com -u root -p \
  --single-transaction --routines --triggers \
  oneapi > /tmp/oneapi.sql

# Import into the new node
mysql -u root -p oneapi < /tmp/oneapi.sql
```

### Option 2: Snapshot API

```bash
curl -H "X-Cluster-Secret: <target-secret>" \
  "https://cn.example.com/api/cluster/snapshot?tables=users,tokens,channels,abilities,options,redemptions,plans,user_plans,plan_usages" \
  -o snapshot.json
```

## Startup order

1. **Start the first node** with `CLUSTER_SEEDS` empty or pointing to itself
2. Wait ~5-10 s for the log `cluster module initialized`
3. **Start the other nodes** with `CLUSTER_SEEDS` pointing to any live node
4. New nodes ping the seed, discover the rest transitively
5. Admin → Cluster → Node Management to verify the list and heartbeats

## Env vars (3-node example)

| Variable | Node 1 | Node 2 | Node 3 |
|---|---|---|---|
| `PORT` | 3000 | 3001 | 3002 |
| `SQL_DSN` | `...oneapi_node1` | `...oneapi_node2` | `...oneapi_node3` |
| `REDIS_CONN_STRING` | `:6379/0` | `:6380/0` | `:6381/0` |
| `CLUSTER_ENABLED` | true | true | true |
| `CLUSTER_NODE_ID` | **1** | **2** | **3** |
| `CLUSTER_NODE_NAME` | `node-cn` | `node-us` | `node-eu` |
| `CLUSTER_NODE_ADDRESS` | `https://cn.example.com` | `https://us.example.com` | `https://eu.example.com` |
| `CLUSTER_SECRET` | same | same | same |
| `CLUSTER_SEEDS` | self or empty | `https://cn.example.com` | `https://cn.example.com` |

**Critical constraints:**

- `CLUSTER_SECRET` identical across all nodes (initial; can rotate independently later)
- `CLUSTER_NODE_ID` unique and matches MySQL `auto_increment_offset`
- `CLUSTER_NODE_ADDRESS` reachable by other nodes (must include scheme `https://`)
- `CLUSTER_SEEDS` only used for first-time discovery; ignored after that

### Node 1 complete .env

```bash
PORT=3000
SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node1?charset=utf8mb4&parseTime=True&loc=Local
REDIS_CONN_STRING=redis://127.0.0.1:6379/0

CLUSTER_ENABLED=true
CLUSTER_NODE_ID=1
CLUSTER_NODE_NAME=node-cn
CLUSTER_NODE_ADDRESS=https://cn.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me

# Optional tuning
CLUSTER_DISCOVERY_INTERVAL=30
CLUSTER_DEAD_PING_INTERVAL=120
CLUSTER_MAX_PING_FAILURES=3
CLUSTER_PUSH_INTERVAL=3
CLUSTER_SYNC_LOGS=true
CLUSTER_BATCH_SIZE=50
```

Node 2 / 3 are similar — only ID / Name / Address / SQL_DSN differ.

## Nginx (ip_hash entry)

```nginx
upstream one_api_cluster {
    ip_hash;  # pin client to one node — quota limits and Redis cache stay consistent
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

`ip_hash` pins a client to one node so rate-limit and cache don't break across hops.

## Rolling upgrades

Recommended: upgrade one node at a time:

```
1. Upgrade Node A (keep B/C running)
2. Wait for A to be healthy and A↔B, A↔C heartbeats to recover
3. Upgrade Node B
4. Upgrade Node C
```

After each step check heartbeats and sync. Backlog events are caught up by Pusher.

## Disaster recovery

| Scenario | Action |
|---|---|
| Node briefly offline | Auto-recovers (`status=1`); backlog events are pushed when it comes back |
| Node offline long → data drift | `mysqldump` from a live node, restart |
| Whole DC fails | Switch DNS to another DC's nodes; `mysqldump` to rebuild |
| Secret leaked | Admin → Node Management → edit secret; next ping auto-propagates |

## Monitoring

Admin → Cluster → Node Management — each row shows:

- `status`: `1` alive, `2` failed/disabled
- `last_heartbeat`: last alive timestamp
- `ping_failures`: failure count since last success

Or via API:

```bash
curl -H "Authorization: Bearer <root-token>" \
  https://cn.example.com/api/cluster_node/
```

## Related

- [Cluster Overview](./overview)
- [Node Management](./node-management)
- [Node Health](./node-health)