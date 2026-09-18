---
title: Cluster Nodes
description: "Register cluster nodes, distribute secrets, monitor heartbeat, enable / disable / ping."
category: decentralization
order: 17
---

# Cluster Nodes

> Register every node with the master so logs, channel routing and failover work across the cluster. UI: `web/default-pro/src/views/setting/ClusterSetting.vue`.

> Cluster mode is gated by env vars. When `CLUSTER_ENABLED != "true"` all `/api/cluster_node/*` endpoints return `集群模式未启用`.

## Cluster Environment Variables

Boot parameters (`cluster/config.go::LoadConfig`):

| Variable | Required | Notes |
|---|---|---|
| `CLUSTER_ENABLED` | yes | `"true"` enables clustering; anything else falls back to single-node |
| `CLUSTER_NODE_ID` | yes | Integer in `[1, 49]`; out-of-range calls `FatalLog` and refuses to start |
| `CLUSTER_NODE_NAME` | no | Default `node-<id>` |
| `CLUSTER_NODE_ADDRESS` | yes | Externally reachable URL of this node |
| `CLUSTER_SECRET` | yes | Initial secret for this node |
| `CLUSTER_SEEDS` | no | Comma-separated seed-node URLs |
| `CLUSTER_PUSH_INTERVAL` | no | Push interval (s); default 3 |
| `CLUSTER_DISCOVERY_INTERVAL` | no | Discovery interval (s); default 30 |
| `CLUSTER_DEAD_PING_INTERVAL` | no | Dead-node re-ping interval (s); default 120 |
| `CLUSTER_MAX_PING_FAILURES` | no | Max consecutive ping failures; default 3 |
| `CLUSTER_SYNC_LOGS` | no | `"false"` disables log sync; default on |
| `CLUSTER_BATCH_SIZE` | no | Sync batch size; default 50 |

## Data Model

`model.ClusterNode` (table `cluster_nodes`):

| Field | Type | Notes |
|---|---|---|
| `node_id` | `int` UNIQUE | Cluster-local id (1–49) |
| `node_name` | `varchar(64)` | Display name |
| `address` | `varchar(256)` | HTTP URL |
| `secret_key` | `varchar(128)` | Shared secret used for cross-node auth |
| `status` | `int` | `1=alive` / `2=failed` |
| `disabled` | `bool` | Soft-disable flag |
| `last_heartbeat` / `last_ping_attempt` | `bigint` | unix seconds |
| `ping_failures` | `int` | Consecutive failure count |
| `created_at` / `updated_at` | `bigint` | unix seconds |

`IsAlive()` = `status == 1 && !disabled`.

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/cluster_node/` | `GET` | Root | Full list with `is_self` flag |
| `/api/cluster_node/:id` | `GET` | Root | Single node detail |
| `/api/cluster_node/` | `POST` | Root | Register a new node (`node_id` 1–49) |
| `/api/cluster_node/` | `PUT` | Root | Update name / address / secret (passing a new secret resets `status` and `ping_failures`) |
| `/api/cluster_node/:id` | `DELETE` | Root | **Soft-disable** (`disabled=true`, `status=2`). Physical delete requires manual SQL |
| `/api/cluster_node/:id/enable` | `POST` | Root | Re-enable (`disabled=false` + status=1 + reset failures + refresh heartbeat) |
| `/api/cluster_node/ping/:id` | `GET` | Root | Active ping; returns the target node's response on success |

Implementation: `controller/cluster_node.go`.

## Register a New Node — `POST /api/cluster_node/`

```json
{
  "node_id": 2,
  "node_name": "node-shanghai",
  "address": "https://sh.example.com",
  "secret": "<32+ char shared secret>"
}
```

Constraints:
- `node_id ∈ [1, 49]` — otherwise `节点编号必须在 1-49 之间`.
- `address` required.
- `secret` required (shared secret other nodes use to authenticate against this node).

## "This Node" Marker

`GetAllClusterNodes` decorates the list with `is_self: bool = (n.NodeId == cluster.NodeID)`. The UI uses this to badge the current node.

## Disable vs Hard-Delete

- **Delete** = soft-disable: sets `disabled=true`, `status=2`. The UI's disable button does exactly this and the response reminds you that physical deletion requires `DELETE FROM cluster_nodes WHERE node_id = ?`.
- **Enable** = `POST /api/cluster_node/:id/enable` — resets `disabled`, status, `ping_failures` and `last_heartbeat`.
- The current node (`nodeId == cluster.NodeID`) cannot be disabled — the controller explicitly rejects it.

## Active Ping

`GET /api/cluster_node/ping/:id` invokes `cluster.PingNode(cluster.GetDB(), &node)`. Failures return `success: false` with a descriptive `message` (HTTP error, auth error, timeout). Success returns `{ data: <peerResponse> }`.

## Frontend Guide

- Table columns: node_id / node_name / address (truncated tooltip) / status chip (green/red) / last heartbeat / actions.
- Row actions:
  - **Ping** — calls the ping endpoint immediately, with a per-row loading spinner.
  - **Edit** — opens a 500 px modal (`node_id` is locked when editing); changing `secret` resets status immediately.
  - **Enable / Disable** — popconfirm.
  - **Delete** — popconfirm → soft-delete.
- **Add node** opens the same modal; the same validations apply.

## Implementation Pointers

| Concern | Location |
|---|---|
| Handler | `controller/cluster_node.go` |
| Node model | `model/cluster_node.go` |
| Config loader | `cluster/config.go::LoadConfig` |
| Heartbeat / ping | `cluster/handler.go`, `cluster/cluster.go` |
| Routes | `router/api.go` |