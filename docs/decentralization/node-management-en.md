---
title: Node Management
description: Cluster node registration, query, enable, and removal.
category: decentralization
order: 2
---

# Node Management

> How to register, view, enable, and disable cluster nodes.

## Where

Admin → Cluster → Node Management.

## List shows

Per row:

- Node ID / Name
- Address (`https://node-b.example.com`)
- Status (Enabled / Failed / Disabled)
- Last heartbeat
- Consecutive ping failures
- Action buttons

A red "Lost" tag means repeated ping failures.

## Fields

| Field | Meaning |
|---|---|
| `node_id` | Node number 1–49; must match `CLUSTER_NODE_ID` and MySQL `auto_increment_offset` |
| `node_name` | Display name |
| `address` | Public URL (with scheme) |
| `secret_key` | This node's authentication secret for inbound requests |
| `status` | `1`=enabled / `2`=failed / `3`=admin-disabled |
| `disabled` | Soft-delete flag |
| `last_heartbeat` | Latest heartbeat timestamp |
| `ping_failures` | Consecutive ping failures |

## Self-registration (recommended)

Each node **auto-registers on first startup**. No manual entry needed:

- Node starts → writes its own row into `cluster_nodes`
- A periodic task refreshes heartbeat and address
- Other nodes discover it via ping

This is the recommended path. Each node only needs to know its own `CLUSTER_NODE_ID` and `CLUSTER_NODE_SECRET`.

## Manual add (rare)

For debugging or transient setups:

1. Admin → Node Management → "Add"
2. Fill:
   - Node ID (1–49)
   - Name (any)
   - Address (peer public URL, e.g. `https://node-b.example.com`)
   - Secret (peer's `CLUSTER_SECRET`)
3. Save

> Self-registration covers almost everything; manual add is for niche cases.

## Enable / Disable / Delete

| Action | Where | Effect |
|---|---|---|
| Enable | Row → "Enable" | `disabled=false` + reset failure counter |
| Disable (soft) | Row → "Disable" | `disabled=true`; no events pushed, still responds to ping |
| Hard delete | Manual SQL | `DELETE FROM cluster_nodes WHERE node_id = ?` |

"Disable" is usually enough (reversible). Hard delete only when a node is permanently retired.

## What Secret is for

Each node requires the request header `X-Cluster-Secret` to equal the **target node's** secret. A secret leak is contained to one node.

### First start

Uses `CLUSTER_SECRET` env var as initial `secret_key`.

### Later starts

Reads from `cluster_nodes.secret_key`, no longer needs the env var.

### Rotation

Change via admin UI → next ping propagates the new value to other nodes.

## How to debug a lost node

1. Check "Last heartbeat" — is it stale?
2. Click "Ping" for a manual ping — see the response
3. Check network — can the two nodes reach each other's `address`?
4. Check Secret — are the `secret_key` values identical on both sides?
5. Read [Node Health](./node-health) for failover troubleshooting

## FAQ

- **My own node doesn't show in the list**: ensure `CLUSTER_NODE_ID` is set; check backend logs for `cluster_nodes` write errors.
- **Two nodes don't see each other**: verify both `address` values are reachable; verify `secret_key` matches.
- **Re-enable didn't help**: restart the node so it rewrites its heartbeat.

## Related

- [Cluster Overview](./overview)
- [Node Health](./node-health)
- [Multi-node Deployment](./deployment)