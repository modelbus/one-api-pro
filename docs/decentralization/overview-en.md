---
title: Cluster Overview
description: What Cluster mode is, when to use it, and how it differs from multi-instance + shared DB.
category: decentralization
order: 1
---

# Cluster Overview

> What is One API Pro's multi-node decentralized deployment, when to use it, and how it differs from ordinary multi-replica.

## What Cluster is

**Cluster mode = multiple One API Pro instances, each with its own MySQL + Redis, with HTTP-based push synchronization between nodes.**

No central node, no shared database. Every node is equal: whichever node you hit serves you.

## vs. "Multi-instance + shared DB"

| Dimension | Cluster | Multi-instance + shared DB (traditional) |
|---|---|---|
| Database | Each node owns its own MySQL | All instances share one MySQL |
| Redis | Each node owns its own | Shared |
| Consistency | Eventual (based on timestamp) | Strong |
| Cross-region latency | Low (local node serves) | High (cross-region DB roundtrip) |
| Offline tolerance | Changes during offline aren't backfilled | Nodes can go down freely |
| Best for | Multi-region / cross-DC | Single-DC multi-replica |

In short: **Cluster = HA + multi-region**; **multi-replica = load sharing inside one DC**.

## When to use Cluster

✅ **Use it when:**

- Multi-DC / multi-region / cross-region DR
- Latency-sensitive users; serve from the nearest node
- Don't want all traffic backhauling to a central DC

❌ **Don't use it when:**

- Single-DC, small-medium traffic: single instance + docker-compose is enough
- You already run K8s multi-replica with shared DB: that's a different pattern, don't mix
- You need strong consistency / distributed transactions: One API Pro doesn't implement cross-node TX

## What it does

- **Decentralized**: nodes are equal, no central coordinator
- **Auto-sync**: any record changed on one node is pushed to all others
- **Conflict converges**: last-writer-wins by `updated_at`, eventual consistency
- **Rate-limit aggregation**: per-channel concurrency / RPM counts are synchronized; global state lives across nodes
- **Zero code change**: triggers capture changes automatically; no business-code edits

## Known limitations

- **Changes made while a node is offline aren't backfilled**: after recovery, restore that node's DB from a live peer (e.g. `mysqldump`).
- **A new node only sees changes made after it joins**: pre-join data must be imported manually.
- **Call logs grow fast — disable their sync**: set `CLUSTER_SYNC_LOGS=false` on each node.

## What's synced

Accounts, tokens, channels, plans, subscriptions, redemption codes, system settings, call logs (optional).

Not synced: the cluster-node table itself (maintained by the discovery mechanism).

## Architecture

```
                 ┌─────────────┐
                 │  Nginx/LB   │  ← single entry, ip_hash
                 └──────┬──────┘
                        │
        ┌───────────────┼───────────────┐
        │               │               │
   ┌────┴────┐    ┌────┴────┐    ┌────┴────┐
   │ Node A  │    │ Node B  │    │ Node C  │
   │ one-api │    │ one-api │    │ one-api │
   │ MySQL   │    │ MySQL   │    │ MySQL   │
   │ Redis   │    │ Redis   │    │ Redis   │
   └────┬────┘    └────┬────┘    └────┬────┘
        │               │               │
        └──────── HTTP push sync ────────┘
```

Every node's data changes are pushed to all live nodes.

## Next

- Enable Cluster → [Multi-node Deployment](./deployment)
- Day-to-day node operations → [Node Management](./node-management)
- Sync troubleshooting → [Node Health](./node-health)