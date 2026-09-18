---
title: Node Management
description: "Register, query, enable and remove nodes."
category: decentralization
order: 2
---

# Node Management

> Register, query, enable and remove nodes.

The `cluster_nodes` table is the single source of truth for node info and runtime state. A node writes its own record on startup, and admins can also add remote nodes manually from the admin UI.

## Node fields

The model lives in `model/cluster_node.go`:

| Field | Type | Description |
| --- | --- | --- |
| `node_id` | int | Node ID `1-49`; matches `CLUSTER_NODE_ID` and MySQL `auto_increment_offset` |
| `node_name` | string | Display name |
| `address` | string | Public URL (with scheme) |
| `secret_key` | string | Per-node secret other nodes use to authenticate (see Per-node Secret below) |
| `status` | int | `1` = alive, `2` = failed (consecutive ping failures), `3` = disabled by admin |
| `disabled` | bool | `true` when disabled by an admin (soft-delete marker) |
| `last_heartbeat` | int64 | Last heartbeat timestamp |
| `ping_failures` | int | Consecutive ping failures |
| `last_ping_attempt` | int64 | Last ping attempt timestamp |

> The model exposes both `status=3` ("disabled by admin") and a `disabled` boolean. `status` reflects runtime state; `disabled` reflects admin intent. The delete handler sets `disabled=true` and `status=2` (see `controller/cluster_node.go::DeleteClusterNode`).

## Self-registration

Every node calls `SaveLocalNode` on first boot:

- If a record for `node_id == CLUSTER_NODE_ID` already exists locally, only `address`, `node_name`, `status`, `last_heartbeat` and `ping_failures=0` are refreshed.
- Otherwise a record is created with `secret_key = CLUSTER_SECRET`.

Why this is intentional:

1. **Admin visibility** — Settings → Node Management shows the local node's address, status and heartbeat for troubleshooting.
2. **Transitive discovery** — the ping response includes the sender's full node list, which the receiver merges locally. So C can learn about A through B.
3. **Liveness signal** — local `last_heartbeat` is refreshed every cycle by `discoverOnce`, reflecting this node's own health.

> Five layers of guards prevent loops: SQL filters in `GetAllRemoteNodes` / `GetAliveNodesForSync`, self-ping rejection in `handlePing`, self-skip in `mergeDiscoveredNodes`, and self-skip in `ApplyEvents`. See the repo README for the full diagram.

## Adding remote nodes

The admin UI **Settings → Node Management** (`web/default-pro/src/views/setting/ClusterSetting.vue`) calls `/api/cluster_node/`:

| Action | Method | Path | Notes |
| --- | --- | --- | --- |
| List | `GET` | `/api/cluster_node/` | All non-disabled nodes; local node carries `is_self=true` |
| Detail | `GET` | `/api/cluster_node/:id` | Single node detail |
| Add | `POST` | `/api/cluster_node/` | `node_id` in `1-49`; `secret` is required (the peer's secret) |
| Update | `PUT` | `/api/cluster_node/` | Update name / address / secret; changing secret auto-recovers `status=1` |
| Disable (soft delete) | `DELETE` | `/api/cluster_node/:id` | Sets `disabled=true` |
| Re-enable | `POST` | `/api/cluster_node/:id/enable` | Resets `disabled=false` and refreshes the heartbeat |
| Manual ping | `GET` | `/api/cluster_node/ping/:id` | One-shot ping for troubleshooting |

All endpoints require **Root privileges** (`middleware.RootAuth`).

## Per-node Secret

Each node carries **its own** secret, persisted in `cluster_nodes.secret_key`, replacing the earlier global-shared-secret design:

- **Security** — one node's secret leaking does not affect others.
- **Flexibility** — every node can rotate its own secret independently.
- **Auto-discovery** — the ping response includes every node's secret, so peers learn each other's secrets.

Lifecycle:

1. **First boot** — uses the `CLUSTER_SECRET` env var as the initial `secret_key`.
2. **Subsequent boots** — read from `cluster_nodes.secret_key`; the env var is no longer consulted.
3. **Rotation** — change the `secret` field via the admin UI; the next ping propagates the new value to peers.
4. **Verification** — the `X-Cluster-Secret` header equals the **target node's** secret (looked up locally).

Adding a new node:

1. Add Node B's record on Node A (input B's `CLUSTER_SECRET`).
2. Add Node A's record on Node B (input A's `CLUSTER_SECRET`).
3. A → B ping uses B's secret, which B verifies with its own secret.
4. B's response carries both A's and B's secrets; A updates its local copy.

## Soft delete

`DELETE /api/cluster_node/:id` does **not** hard-delete — it sets `disabled = true`:

- Prevents a deleted node from "regrowing" via ping-driven re-registration.
- Disabled nodes still respond to pings (so peers know they're online) but won't fetch this node's info.
- Hard delete requires manual SQL: `DELETE FROM cluster_nodes WHERE node_id = ?;`

Re-enable via `POST /api/cluster_node/:id/enable`; this also resets `ping_failures=0` and refreshes `last_heartbeat`.

## API cheatsheet

| Path | Method | Purpose |
| --- | --- | --- |
| `/api/cluster/ping` | `POST` | Inter-node heartbeat (internal) |
| `/api/cluster/sync` | `POST` | Inter-node event push (internal) |
| `/api/cluster_node/` | `GET` / `POST` / `PUT` | List / add / update nodes |
| `/api/cluster_node/:id` | `GET` / `DELETE` | Detail / disable |
| `/api/cluster_node/:id/enable` | `POST` | Re-enable |
| `/api/cluster_node/ping/:id` | `GET` | Manual ping (admin) |

The full contract lives in [Cluster API](/en/api/cluster).

Next: [Config Sync](/en/decentralization/config-sync) · [Node Health](/en/decentralization/node-health).
