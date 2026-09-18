---
title: 集群 API
description: "集群 API：集群管理 API"
category: api
order: 21
---
## 13. 集群管理 API

集群 API 用于去中心化多活集群的节点管理、数据同步和节点发现。使用两种认证方式：

- **集群密钥（Cluster Secret）**：节点间内部通信，Header 格式 `X-Cluster-Secret: <节点 secret>`（目标节点的 secret，从本地 DB 查）
- **Root 管理员**：管理后台调用，使用 `Authorization: Bearer <Root Access Token>` 或 Cookie Session

### 13.1 节点发现与心跳（内部）

**接口：** `POST /api/cluster/ping`

**认证：** 集群密钥

**说明：** 节点间双向 ping，用于集群节点发现和存活检测。请求方在收到响应后会保存对方返回的完整节点列表，实现传递性发现。

**请求体：**

```

**请求字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| node_id | int | 是 | 发起 ping 的节点 ID（与 `CLUSTER_NODE_ID` 一致） |
| node_name | string | 是 | 发起 ping 的节点名称 |
| address | string | 是 | 发起 ping 的节点公网地址（包含协议前缀） |
| secret_key | string | 是 | 发起 ping 的节点 secret |

**响应（包含全部已知节点列表）：**

```

### 13.2 数据同步（内部）

**接口：** `POST /api/cluster/sync`

**认证：** 集群密钥

**说明：** 接收其他节点推送的数据变更事件。接收方会跳过 `event.NodeId == 本机 NodeID` 的事件，避免回环。

**请求体：**

```

**请求字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| source_node_id | int | 是 | 事件来源节点 ID |
| events | array | 是 | 事件列表，每批默认最多 50 条 |

**事件字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 事件唯一 ID |
| table_name | string | 表名（users|
| operation | string | 操作类型：`INSERT` / `UPDATE` / `DELETE` |
| primary_key | uint | 主键值 |
| data | object | 变更数据（DELETE 时为空） |
| event_time | int64 | 事件时间戳（秒） |

**响应：**

```

### 13.3 获取全部节点列表

**接口：** `GET /api/cluster_node/`

**认证：** Root 管理员

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| p | int | 页码，默认 0 |

**返回值：**

```

**返回字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 记录 ID |
| node_id | int | 节点编号（与 `CLUSTER_NODE_ID` 对应） |
| node_name | string | 节点名称 |
| address | string | 节点公网地址 |
| status | int | 状态：1=存活, 2=失败 |
| last_heartbeat | int64 | 最后心跳时间（Unix 时间戳） |
| ping_failures | int | 连续 ping 失败次数 |
| disabled | bool | 是否已被管理员禁用 |
| secret_key | string | 节点的访问密钥 |

### 13.4 获取单个节点

**接口：** `GET /api/cluster_node/:id`

**认证：** Root 管理员

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 节点记录 ID（非 node_id） |

### 13.5 添加节点

**接口：** `POST /api/cluster_node/`

**认证：** Root 管理员

**请求体：**

```

**请求字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| node_id | int | 是 | 节点编号（1-49），集群内唯一 |
| node_name | string | 是 | 节点名称 |
| address | string | 是 | 节点公网地址（包含协议前缀如 `https://`） |
| secret_key | string | 是 | 节点初始 secret，应与目标节点 `CLUSTER_SECRET` 一致 |

### 13.6 更新节点

**接口：** `PUT /api/cluster_node/`

**认证：** Root 管理员

**请求体：**

```

**请求字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | uint | 是 | 节点记录 ID |
| node_id | int | 否 | 节点编号 |
| node_name | string | 否 | 节点名称 |
| address | string | 否 | 节点公网地址 |
| secret_key | string | 否 | 更新该节点的 secret，更新后其他节点访问本节点需要使用新值 |
| disabled | bool | 否 | 禁用状态 |

> **secret_key 更新机制：** 节点 secret 在 ping 时由 `X-Cluster-Secret` 头部携带，目标节点使用自己的 secret 校验。当管理员更新某节点的 secret 后，其他节点在下一次 ping 时会自动学习到新值。

### 13.7 软删除节点

**接口：** `DELETE /api/cluster_node/:id`

**认证：** Root 管理员

**说明：** 不物理删除记录，而是设置 `disabled = true`。被禁用的节点仍然会响应 ping（让对方知道其在线），但其他节点不会再向其推送事件。如需物理删除需手动执行 SQL：`DELETE FROM cluster_nodes WHERE node_id = ?`。

### 13.8 重新启用已禁用节点

**接口：** `POST /api/cluster_node/:id/enable`

**认证：** Root 管理员

**说明：** 将 `disabled` 重置为 `false`，恢复节点参与集群通信的能力。

### 13.9 手动 Ping 节点

**接口：** `GET /api/cluster_node/ping/:id`

**认证：** Root 管理员

**说明：** 主动发起一次 ping 请求，常用于排查节点连通性问题。

