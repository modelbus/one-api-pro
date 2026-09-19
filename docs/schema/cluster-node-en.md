---
title: Cluster Node
description: A node registered in Cluster mode and its health state.
category: schema
order: 12
---

# Cluster Node

## What it is

`ClusterNode` is the row registered in the central cluster table for **each node in a multi-node Cluster deployment**. Every node registers itself to the others on startup and continues to heartbeat.

If you run a single instance, you can **skip this section entirely**.

## Where to find it

- **Admin → Cluster**: node list, status, last heartbeat
- **Monitoring**: every node actively pushes data changes to the others

## Operator-relevant fields

| Field | Meaning | Effect |
|---|---|---|
| Node ID | Globally unique ID | Changing it is treated as a new identity. |
| Node name | Display only | No effect on communication. |
| Status | Enabled / Disabled | Disabled nodes stop receiving pushes. |
| Address | `host:port` | Push target; wrong value = push fails. |
| Secret | Bearer Token for inter-node trust | Must match across the cluster. |
| Last heartbeat | Unix seconds | Stale → flagged "lost". |

## Multi-node deployment

For N nodes:

1. Each node has a unique `CLUSTER_NODE_ID`
2. All nodes share the same `CLUSTER_NODE_SECRET`
3. Every node knows every other node's address
4. Nodes push data changes to each other over HTTP — no central DB

See [Decentralization](/en/decentralization/overview) and [Multi-node Deployment](/en/decentralization/deployment).

## Related pages

- [Cluster Overview](/en/decentralization/overview)
- [Node Management](/en/decentralization/node-management)
- [Config Sync](/en/decentralization/config-sync)
- [Node Health](/en/decentralization/node-health)
- [Multi-node Deployment](/en/decentralization/deployment)
- [Cluster Settings](/en/decentralization/cluster-settings)

## Related API

- `GET /api/cluster_node/` — list
- `POST /api/cluster_node/` — create
- `PUT /api/cluster_node/` — update
- `POST /api/cluster_node/:id/enable` — enable
- `GET /api/cluster_node/ping/:id` — ping test