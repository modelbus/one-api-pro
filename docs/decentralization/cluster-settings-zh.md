---
title: 集群节点管理
description: "注册集群节点、密钥分发、心跳状态、节点启用 / 禁用 / Ping。"
category: decentralization
order: 17
---

# 集群节点管理

> 多节点集群模式下，把每个 node 注册到主节点，便于日志同步、渠道分流与故障转移。前端组件：`web/default-pro/src/views/setting/ClusterSetting.vue`。

> 集群模式由环境变量开启；如果 `CLUSTER_ENABLED != "true"`，所有 `/api/cluster_node/*` 接口都会返回 `集群模式未启用`。

## 环境变量 / Cluster Env

启动参数（来自 `cluster/config.go::LoadConfig`）：

| 变量 | 必填 | 说明 |
|---|---|---|
| `CLUSTER_ENABLED` | 是 | `"true"` 才启用集群；其它值回退单节点 |
| `CLUSTER_NODE_ID` | 是 | 1–49 之间的整数，FatalLog 越界即拒启 |
| `CLUSTER_NODE_NAME` | 否 | 默认 `node-<id>` |
| `CLUSTER_NODE_ADDRESS` | 是 | 本节点对外可达的 URL |
| `CLUSTER_SECRET` | 是 | 本节点的初始 secret |
| `CLUSTER_SEEDS` | 否 | 种子节点地址列表（逗号分隔） |
| `CLUSTER_PUSH_INTERVAL` | 否 | push 周期秒，默认 3 |
| `CLUSTER_DISCOVERY_INTERVAL` | 否 | 节点发现周期秒，默认 30 |
| `CLUSTER_DEAD_PING_INTERVAL` | 否 | 死节点重试 ping 周期秒，默认 120 |
| `CLUSTER_MAX_PING_FAILURES` | 否 | 最大连续 ping 失败次数，默认 3 |
| `CLUSTER_SYNC_LOGS` | 否 | `"false"` 关闭日志同步；默认开启 |
| `CLUSTER_BATCH_SIZE` | 否 | 同步批次大小，默认 50 |

## 数据模型 / Data Model

`model.ClusterNode`（`cluster_nodes` 表）：

| 字段 | 类型 | 说明 |
|---|---|---|
| `node_id` | `int` UNIQUE | 集群内节点编号（1–49） |
| `node_name` | `varchar(64)` | 节点名 |
| `address` | `varchar(256)` | 节点 URL（HTTP） |
| `secret_key` | `varchar(128)` | 用于跨节点认证的密钥 |
| `status` | `int` | `1=alive` / `2=failed` |
| `disabled` | `bool` | 软禁用标记 |
| `last_heartbeat` / `last_ping_attempt` | `bigint` | unix 秒 |
| `ping_failures` | `int` | 连续 ping 失败次数 |
| `created_at` / `updated_at` | `bigint` | unix 秒 |

`IsAlive()` = `status == 1 && !disabled`。

## 接口一览 / Endpoints

| Endpoint | Method | 鉴权 | 说明 |
|---|---|---|---|
| `/api/cluster_node/` | `GET` | Root | 全量节点列表（带 `is_self` 标识） |
| `/api/cluster_node/:id` | `GET` | Root | 单节点详情 |
| `/api/cluster_node/` | `POST` | Root | 注册新节点（`node_id` 1–49） |
| `/api/cluster_node/` | `PUT` | Root | 修改节点名称 / 地址 / secret（传 secret 会重置 status + ping_failures） |
| `/api/cluster_node/:id` | `DELETE` | Root | **软禁用**（设 `disabled=true`，status=2）。物理删除需手动 SQL |
| `/api/cluster_node/:id/enable` | `POST` | Root | 重新启用（`disabled=false` + status=1 + 清零失败次数） |
| `/api/cluster_node/ping/:id` | `GET` | Root | 主动向目标节点发 ping，返回对方 ping 的结果对象 |

实现：`controller/cluster_node.go`。

## 注册新节点 / `POST /api/cluster_node/`

请求体：

```json
{
  "node_id": 2,
  "node_name": "node-shanghai",
  "address": "https://sh.example.com",
  "secret": "<32+ char shared secret>"
}
```

约束：
- `node_id` ∈ [1, 49]，否则返回 `节点编号必须在 1-49 之间`；
- `address` 不能为空；
- `secret` 不能为空（这是其他节点访问本节点的认证密钥）；

## 「本节点」徽章 / `is_self`

列表接口 `GetAllClusterNodes` 给每行加一个 `is_self: bool` 字段，值为 `n.NodeId == cluster.NodeID`。前端用此徽章标识「这是当前节点」。

## 删除 vs 禁用 / Disable vs Hard-Delete

- **删除** = 软禁用：`disabled=true`，`status=2`。UI 行的「禁用」按钮触发；提示文案「节点已禁用（物理删除需要手动 SQL: `DELETE FROM cluster_nodes WHERE node_id = ?`）」。
- **启用** = `POST /api/cluster_node/:id/enable`；把 `disabled=false` + status=1 + 清零失败次数 + 刷新 `last_heartbeat`。
- 当前节点不能被禁用（`controller` 内显式拦截：`nodeId == cluster.NodeID`）。

## Ping / 主动探测

`GET /api/cluster_node/ping/:id` 调用 `cluster.PingNode(cluster.GetDB(), &node)`，向目标节点发 ping；失败时返回 `success: false`，message 含失败原因（HTTP / 鉴权 / 超时）；成功时返回 `{ data: <对方响应> }`。

## 前端操作指南 / Frontend Guide

- 顶部表格列：节点 ID / 节点名 / address（截断 tooltip）/ 状态 chip（green / red）/ 上次心跳时间 / 操作。
- 行内操作：
  - **Ping**：立刻调 ping 端点，loading 圈到对应行；
  - **编辑**：打开 500px 弹窗（`node_id` 编辑时锁定）— 注意修改 secret 会立即重置 status；
  - **启用 / 禁用**：popconfirm；
  - **删除**：popconfirm → 软删除。
- 「新增节点」按钮打开相同弹窗，校验同上。

## 接口实现 / Implementation Pointers

| 关注点 | 位置 |
|---|---|
| Handler | `controller/cluster_node.go` |
| 节点模型 | `model/cluster_node.go` |
| 配置加载 | `cluster/config.go::LoadConfig` |
| 心跳 / Ping | `cluster/handler.go` / `cluster/cluster.go` |
| 路由 | `router/api.go` |