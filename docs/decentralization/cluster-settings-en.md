---
title: Cluster Settings
description: Cluster mode env vars, node management, and API.
category: decentralization
order: 6
---

# Cluster Settings

> Cluster mode config, admin UI, and API.

## Environment variables

Loaded at startup:

| Variable | Required | Default | Notes |
|---|---|---|---|
| `CLUSTER_ENABLED` | Yes | — | Set to `"true"` to enable cluster |
| `CLUSTER_NODE_ID` | Yes | — | Integer 1–49; out-of-range fails fast |
| `CLUSTER_NODE_NAME` | No | `node-<id>` | Node name |
| `CLUSTER_NODE_ADDRESS` | Yes | — | This node's publicly reachable URL |
| `CLUSTER_SECRET` | Yes | — | This node's initial secret |
| `CLUSTER_SEEDS` | No | empty | Seed addresses for first-time discovery (comma-separated) |
| `CLUSTER_DISCOVERY_INTERVAL` | No | 30 | Discovery period (seconds) |
| `CLUSTER_DEAD_PING_INTERVAL` | No | 120 | Dead-node ping period (seconds) |
| `CLUSTER_MAX_PING_FAILURES` | No | 3 | Failures before marking dead |
| `CLUSTER_PUSH_INTERVAL` | No | 3 | Pusher push period (seconds) |
| `CLUSTER_BATCH_SIZE` | No | 50 | Sync batch size |
| `CLUSTER_SYNC_LOGS` | No | true | Set `"false"` to disable log sync |

> If `CLUSTER_ENABLED != "true"`, all `/api/cluster_node/*` endpoints return `cluster mode not enabled`.

## Admin → Cluster

The node list lives in [Node Management](./node-management).

## Node fields

| Field | Type | Notes |
|---|---|---|
| `node_id` | int UNIQUE | Cluster-wide node number (1–49) |
| `node_name` | varchar(64) | Node name |
| `address` | varchar(256) | Node URL (with `http://` or `https://`) |
| `secret_key` | varchar(128) | Cross-node auth key |
| `status` | int | `1=alive` / `2=failed` |
| `disabled` | bool | Soft-disable flag |
| `last_heartbeat` | bigint | Last heartbeat unix seconds |
| `ping_failures` | int | Consecutive ping failures |

## API

| Endpoint | Method | Auth | Notes |
|---|---|---|---|
| `/api/cluster_node/` | GET | Root | All nodes (with `is_self` flag) |
| `/api/cluster_node/:id` | GET | Root | One node |
| `/api/cluster_node/` | POST | Root | Register a new node |
| `/api/cluster_node/` | PUT | Root | Edit name / address / secret |
| `/api/cluster_node/:id` | DELETE | Root | Soft-disable |
| `/api/cluster_node/:id/enable` | POST | Root | Re-enable |
| `/api/cluster_node/ping/:id` | GET | Root | Manual ping |

## Register a new node

```json
{
  "node_id": 1,
  "node_name": "node-cn",
  "address": "https://cn.example.com",
  "secret": "<peer's CLUSTER_SECRET>"
}
```

Constraints:

- `node_id` ∈ [1, 49]
- `address` non-empty
- `secret` non-empty

> Self-registration (auto-write on startup) covers most cases; manual add is for debugging / transitory setups.

## "This node" badge

The list endpoint adds `is_self: bool` (`n.NodeId == cluster.NodeID`). The UI uses this to highlight the local node.

## Delete vs Disable

| Action | API | Effect |
|---|---|---|
| Delete (soft) | `DELETE /api/cluster_node/:id` | `disabled=true` + `status=2` |
| Enable | `POST /api/cluster_node/:id/enable` | `disabled=false` + `status=1` + reset failure counter |
| Hard delete | Manual SQL | `DELETE FROM cluster_nodes WHERE node_id = ?` |

> The current node cannot be disabled (controller explicitly blocks it).

## Manual Ping

`GET /api/cluster_node/ping/:id` pings the target node:

- Success → `{ data: <peer response> }`
- Failure → `{ success: false, message: <reason> }`

Useful for diagnosing "why is this node dead".

## FAQ

- **All env vars correct but endpoints return "cluster mode not enabled"**: confirm `CLUSTER_ENABLED` literally equals the string `"true"`.
- **`CLUSTER_NODE_ID` out of range**: must be 1–49; otherwise the process refuses to start.

## Related

- [Cluster Overview](./overview)
- [Node Management](./node-management)
- [Multi-node Deployment](./deployment)