---
title: Cluster API
description: "Cluster management endpoints for the decentralized multi-active cluster."
category: api
order: 21
---
## 13. Cluster Management API

The cluster API manages nodes, syncs data and supports node discovery for the decentralized multi-active cluster. Two auth methods are used:

- **Cluster Secret** for inter-node traffic; header format `X-Cluster-Secret: <node secret>` (the target node's secret, looked up from the local DB).
- **Root admin** for the admin UI; uses `Authorization: Bearer <Root Access Token>` or cookie session.

### 13.1 Node Discovery and Heartbeat (internal)

**Endpoint:** `POST /api/cluster/ping`

**Auth:** Cluster secret

**Description:** Two-way ping between nodes for cluster discovery and liveness checks. The responder returns its full node list so the requester can do transitive discovery.

**Request body:**

```json
{
  "node_id": 1,
  "node_name": "node-cn",
  "address": "https://cn.example.com",
  "secret_key": "node-1-secret"
}
```

**Request fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| node_id | int | yes | ID of the node sending the ping (matches `CLUSTER_NODE_ID`) |
| node_name | string | yes | Name of the node sending the ping |
| address | string | yes | Public address of the node sending the ping (including the scheme prefix) |
| secret_key | string | yes | Secret of the node sending the ping |

**Response (includes the full known node list):**

```json
{
  "success": true,
  "data": {
    "nodes": [
      {
        "id": 1,
        "node_id": 1,
        "node_name": "node-cn",
        "address": "https://cn.example.com",
        "status": 1,
        "last_heartbeat": 1718000000,
        "ping_failures": 0,
        "disabled": false,
        "created_at": 1718000000,
        "updated_at": 1718000000
      }
    ]
  }
}
```


### 13.2 Data Sync (internal)

**Endpoint:** `POST /api/cluster/sync`

**Auth:** Cluster secret

**Description:** Receives data-change events pushed from other nodes. The receiver skips events whose `event.NodeId` equals its own node ID to avoid loops.

**Request body:**

```json
{
  "source_node_id": 1,
  "events": [
    {
      "id": 1001,
      "table_name": "channels",
      "operation": "UPDATE",
      "primary_key": 5,
      "data": {
        "id": 5,
        "name": "OpenAI",
        "status": 2
      },
      "event_time": 1718000000
    }
  ]
}
```

**Request fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| source_node_id | int | yes | ID of the node that emitted the events |
| events | array | yes | Event list, max 50 per batch by default |

**Event fields:**

| Field | Type | Description |
|-------|------|-------------|
| id | int64 | Unique event ID |
| table_name | string | Table name (users / tokens / channels / abilities / options / plans / user_plans / redemptions, etc.) |
| operation | string | Operation type: `INSERT` / `UPDATE` / `DELETE` |
| primary_key | uint | Primary key value |
| data | object | Change payload (empty on `DELETE`) |
| event_time | int64 | Event timestamp (seconds) |

**Response:**

```json
{
  "success": true,
  "data": {
    "applied": 1,
    "skipped": 0
  }
}
```


### 13.3 List All Nodes

**Endpoint:** `GET /api/cluster_node/`

**Auth:** Root admin

**Query parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| p | int | Page number, default 0 |

**Response:**

```json
{
  "success": true,
  "message": "",
  "data": [
    {
      "id": 1,
      "node_id": 1,
      "node_name": "node-cn",
      "address": "https://cn.example.com",
      "status": 1,
      "last_heartbeat": 1718000000,
      "ping_failures": 0,
      "disabled": false,
      "secret_key": "node-1-secret",
      "created_at": 1718000000,
      "updated_at": 1718000000
    }
  ]
}
```

**Response fields:**

| Field | Type | Description |
|-------|------|-------------|
| id | uint | Record ID |
| node_id | int | Node number (matches `CLUSTER_NODE_ID`) |
| node_name | string | Node name |
| address | string | Node public address |
| status | int | Status: 1=alive, 2=failed |
| last_heartbeat | int64 | Last heartbeat time (Unix timestamp) |
| ping_failures | int | Consecutive ping failures |
| disabled | bool | Whether the admin has disabled this node |
| secret_key | string | The node's access secret |


### 13.4 Get a Single Node

**Endpoint:** `GET /api/cluster_node/:id`

**Auth:** Root admin

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| id | int | Node record ID (not `node_id`) |


### 13.5 Add a Node

**Endpoint:** `POST /api/cluster_node/`

**Auth:** Root admin

**Request body:**

```json
{
  "node_id": 2,
  "node_name": "node-us",
  "address": "https://us.example.com",
  "secret_key": "node-2-secret"
}
```

**Request fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| node_id | int | yes | Node number (1-49), unique within the cluster |
| node_name | string | yes | Node name |
| address | string | yes | Node public address (must include the scheme prefix, e.g. `https://`) |
| secret_key | string | yes | Initial secret; must match the target node's `CLUSTER_SECRET` |


### 13.6 Update a Node

**Endpoint:** `PUT /api/cluster_node/`

**Auth:** Root admin

**Request body:**

```json
{
  "id": 2,
  "node_id": 2,
  "node_name": "node-us",
  "address": "https://us.example.com",
  "secret_key": "new-secret-value",
  "disabled": false
}
```

**Request fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| id | uint | yes | Node record ID |
| node_id | int | no | Node number |
| node_name | string | no | Node name |
| address | string | no | Node public address |
| secret_key | string | no | Update the node's secret. After the update, other nodes must use the new value when contacting this node. |
| disabled | bool | no | Disabled state |

> **Secret rotation:** A node's secret is carried by the `X-Cluster-Secret` header during ping and verified by the target against its own secret. When the admin updates a node's secret, other nodes learn the new value automatically on the next ping.


### 13.7 Soft Delete a Node

**Endpoint:** `DELETE /api/cluster_node/:id`

**Auth:** Root admin

**Description:** Does not hard-delete the record; sets `disabled = true` instead. Disabled nodes still respond to pings (so peers know they are online), but no other node will push events to them. Hard delete requires running SQL manually: `DELETE FROM cluster_nodes WHERE node_id = ?`.


### 13.8 Re-enable a Disabled Node

**Endpoint:** `POST /api/cluster_node/:id/enable`

**Auth:** Root admin

**Description:** Resets `disabled` to `false`, restoring the node's participation in cluster communication.


### 13.9 Manually Ping a Node

**Endpoint:** `GET /api/cluster_node/ping/:id`

**Auth:** Root admin

**Description:** Triggers a single ping, typically used to diagnose connectivity.
