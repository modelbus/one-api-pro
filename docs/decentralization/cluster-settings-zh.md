---
title: 集群设置
description: 集群模式的环境变量、节点管理与 API。
category: decentralization
order: 6
---

# 集群设置

> Cluster 模式的配置、管理界面与 API。

## 环境变量

启动时读取：

| 变量 | 必填 | 默认 | 说明 |
|---|---|---|---|
| `CLUSTER_ENABLED` | 是 | — | `"true"` 才启用集群 |
| `CLUSTER_NODE_ID` | 是 | — | 1–49 之间的整数，越界即拒启 |
| `CLUSTER_NODE_NAME` | 否 | `node-<id>` | 节点名 |
| `CLUSTER_NODE_ADDRESS` | 是 | — | 本节点对外可达的 URL |
| `CLUSTER_SECRET` | 是 | — | 本节点的初始 secret |
| `CLUSTER_SEEDS` | 否 | 空 | 种子节点地址列表（逗号分隔），仅首次发现使用 |
| `CLUSTER_DISCOVERY_INTERVAL` | 否 | 30 | 节点发现周期（秒） |
| `CLUSTER_DEAD_PING_INTERVAL` | 否 | 120 | 死节点重试 ping 周期（秒） |
| `CLUSTER_MAX_PING_FAILURES` | 否 | 3 | 连续失败多少次标 dead |
| `CLUSTER_PUSH_INTERVAL` | 否 | 3 | Pusher 推送周期（秒） |
| `CLUSTER_BATCH_SIZE` | 否 | 50 | 同步批次大小 |
| `CLUSTER_SYNC_LOGS` | 否 | true | `"false"` 关闭日志同步 |

> 如果 `CLUSTER_ENABLED != "true"`，所有 `/api/cluster_node/*` 接口都返回 `集群模式未启用`。

## 后台 → 集群设置

管理员视角的节点列表。详见 [节点管理](./node-management)。

## 节点字段

| 字段 | 类型 | 说明 |
|---|---|---|
| `node_id` | int UNIQUE | 集群内节点编号（1–49） |
| `node_name` | varchar(64) | 节点名 |
| `address` | varchar(256) | 节点 URL（含 http:// 或 https://） |
| `secret_key` | varchar(128) | 跨节点认证密钥 |
| `status` | int | `1=alive` / `2=failed` |
| `disabled` | bool | 软禁用标记 |
| `last_heartbeat` | bigint | 最近心跳 unix 秒 |
| `ping_failures` | int | 连续 ping 失败次数 |

## API 一览

| Endpoint | Method | 鉴权 | 说明 |
|---|---|---|---|
| `/api/cluster_node/` | GET | Root | 全量节点列表（含 `is_self` 标识本节点） |
| `/api/cluster_node/:id` | GET | Root | 单节点详情 |
| `/api/cluster_node/` | POST | Root | 注册新节点 |
| `/api/cluster_node/` | PUT | Root | 改节点名 / 地址 / secret |
| `/api/cluster_node/:id` | DELETE | Root | 软禁用（设 `disabled=true`） |
| `/api/cluster_node/:id/enable` | POST | Root | 重新启用 |
| `/api/cluster_node/ping/:id` | GET | Root | 主动 ping |

## 注册新节点

请求体：

```json
{
  "node_id": 1,
  "node_name": "node-cn",
  "address": "https://cn.example.com",
  "secret": "<对端的 CLUSTER_SECRET>"
}
```

约束：

- `node_id` ∈ [1, 49]
- `address` 不能为空
- `secret` 不能为空

> 自注册（启动时自动写自身记录）已覆盖大多数场景，手动新增主要用于调试 / 临时让节点互认。

## 「本节点」徽章

列表接口的 `is_self: bool` 字段标识当前节点（`n.NodeId == cluster.NodeID`）。前端用此徽章高亮「这是当前节点」。

## 删除 vs 禁用

| 操作 | API | 效果 |
|---|---|---|
| 删除（软） | `DELETE /api/cluster_node/:id` | `disabled=true` + `status=2` |
| 启用 | `POST /api/cluster_node/:id/enable` | `disabled=false` + `status=1` + 清零失败次数 |
| 物理删除 | 手动 SQL | `DELETE FROM cluster_nodes WHERE node_id = ?` |

> 当前节点不能被禁用（`controller` 显式拦截）。

## 主动 Ping

`GET /api/cluster_node/ping/:id` 向目标节点发起 ping：

- 成功 → 返回 `{ data: <对方响应> }`
- 失败 → `{ success: false, message: <失败原因> }`

常用于排查「为什么某节点状态是 dead」。

## 常见问题

- **环境变量都对了但接口返回「集群模式未启用」**：检查 `CLUSTER_ENABLED` 是否真的等于字符串 `"true"`
- **`CLUSTER_NODE_ID` 越界**：必须是 1-49 的整数；超出直接拒启

## 相关

- [Cluster 概览](./overview)
- [节点管理](./node-management)
- [多节点部署](./deployment)