---
title: Cluster Overview
description: "Design goals, constraints and suitable use cases of the decentralized multi-active cluster."
category: decentralization
order: 1
---

# Cluster Overview

> Design goals, constraints and suitable use cases of the decentralized multi-active cluster.

## Goals

One API Pro's **Cluster mode** provides a decentralized, multi-active multi-node deployment. Its core goals are:

- **No shared database**: every node owns its own MySQL and Redis; nodes actively push sync events to each other over HTTP.
- **Zero business intrusion**: data changes are captured via GORM callbacks — no business-code changes needed.
- **Multi-active local access**: deploy across regions / data centers, serve users from the nearest node.
- **Convergent conflicts**: last-writer-wins by `updated_at` ensures eventual consistency.

## When to use

| Scenario | Recommendation |
| --- | --- |
| Single-region, modest traffic | A single instance is enough — no Cluster needed |
| Multi-region / multi-AZ / cross-region DR | ✅ Cluster mode |
| Latency-sensitive, want local access | ✅ Cluster mode |
| K8s multi-replica + shared DB | Use the multi-instance-shared-DB pattern (see `install/docker-compose`); no Cluster needed |
| Strict consistency / distributed transactions | ❌ Not suitable; One API Pro does not implement cross-node transactions |

## Architecture at a glance

```text
              ┌─────────────┐
              │  Nginx/LB   │   (single entry, ip_hash LB)
              └──────┬──────┘
                     │
       ┌─────────────┼─────────────┐
       │             │             │
 ┌─────┴─────┐ ┌─────┴─────┐ ┌─────┴─────┐
 │  Node A   │ │  Node B   │ │  Node C   │
 │ one-api   │ │ one-api   │ │ one-api   │
 │ + MySQL   │ │ + MySQL   │ │ + MySQL   │
 │ + Redis   │ │ + Redis   │ │ + Redis   │
 └─────┬─────┘ └─────┬─────┘ └─────┬─────┘
       │             │             │
       └────── HTTP push of sync events ──────┘
```

Every node is equal: any data change on a node is actively pushed to all alive nodes.

## Core characteristics

- **Decentralized** — peers, no central coordinator.
- **Zero-invasion** — `INSERT` / `UPDATE` / `DELETE` are captured automatically by GORM callbacks.
- **Async push** — sync runs in a background goroutine, never blocking the main flow.
- **Conflict resolution** — the receiver compares `updated_at`; only the newer version is written (last-writer-wins).
- **Rate-limit sync** — `channel_counters` (channel concurrency / RPM) sync per node, so global rate-limit state can be aggregated across nodes.
- **Single-node compatible** — without `CLUSTER_*` env vars, the system runs in plain single-node mode with no side effects.

## Sync scope

| Table | Synced? | Notes |
| --- | --- | --- |
| `users` | ✅ | user accounts |
| `tokens` | ✅ | API tokens |
| `channels` | ✅ | provider channels |
| `abilities` | ✅ | channel abilities |
| `options` | ✅ | system settings |
| `redemptions` | ✅ | redemption codes |
| `plans` | ✅ | subscription plans |
| `user_plans` | ✅ | user subscriptions |
| `plan_usages` | ✅ | plan usage |
| `channel_counters` | ✅ | channel rate-limit counters |
| `cluster_nodes` | 🔄 | maintained by the discovery mechanism, not data sync |
| `logs` | ⚠️ | controlled by `CLUSTER_SYNC_LOGS` |

## Design trade-offs

### Push-only, no active pull

Data sync relies entirely on GORM callbacks + HTTP push; **active cross-node pull is not implemented**.

Why:

1. **Business intrusion** — pull would require knowing each table's business-unique field, polluting business code.
2. **Primary-key conflicts** — auto-increment IDs differ across nodes (different `auto_increment_offset`); using the source's ID would break the offset design.
3. **Complexity** — high maintenance cost for limited reliability gain.
4. **Push is enough** — covers ~95% of normal scenarios (alive nodes, normal traffic).

### Known limits

- **Changes made while a node is offline are not back-filled**; the node must be re-seeded via `mysqldump` from a live node after coming back.
- New nodes only see changes from the moment they join — no history.
- Disable `logs` sync (`CLUSTER_SYNC_LOGS=false`) if the table grows too large.

## Cluster vs shared-DB

| Aspect | Cluster (decentralized) | Multi-instance shared DB (classic) |
| --- | --- | --- |
| Database | per-node independent MySQL | shared MySQL |
| Redis | per-node | shared |
| Consistency | eventual (`updated_at` LWW) | strong (shared DB) |
| Cross-region latency | low (local node) | high (DB round-trip) |
| Offline tolerance | data lost during offline, manual re-seed | any node may go down without loss |
| Fit | cross-region / multi-AZ | single-AZ multi-replica |

Next: [Node Management](/en/decentralization/node-management) · [Config Sync](/en/decentralization/config-sync) · [Node Health](/en/decentralization/node-health).
