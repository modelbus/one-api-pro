---
title: 节点管理
description: "节点的注册、查询、启用与删除。"
category: decentralization
order: 2
---

# 节点管理

> 节点的注册、查询、启用与删除。
> Register, query, enable and remove nodes.

`cluster_nodes` 表是节点信息与运行时状态的唯一来源。节点既会在 **本地启动时自动写入自身记录**，也支持 **通过管理后台添加远程节点**。

The `cluster_nodes` table is the single source of truth for node info and runtime state. A node writes its own record on startup, and admins can also add remote nodes manually from the admin UI.

## 节点字段 / Node fields

模型定义见 `model/cluster_node.go`：

The model lives in `model/cluster_node.go`:

| 字段 / Field | 类型 / Type | 说明 / Description |
| --- | --- | --- |
| `node_id` | int | 节点编号 `1-49`；与 `CLUSTER_NODE_ID`、`MySQL auto_increment_offset` 一致 |
| `node_name` | string | 显示名 / display name |
| `address` | string | 公网 URL（含协议前缀） |
| `secret_key` | string | 本节点对其他节点的认证 secret（详见下文 Per-node Secret） |
| `status` | int | `1` = 启用（alive），`2` = 失败（ping 多次失败），`3` = 被管理员禁用 |
| `disabled` | bool | `true` 表示管理员禁用，软删除标记 |
| `last_heartbeat` | int64 | 最近心跳时间戳 |
| `ping_failures` | int | 连续 ping 失败次数 |
| `last_ping_attempt` | int64 | 最近一次 ping 尝试 |

> 模型里同时存在 `status=3`（被管理员禁用），与 `disabled` 字段语义重叠。`status` 是当前 runtime 状态；`disabled` 是管理员意图的禁用。删除操作同时设置 `disabled=true` 与 `status=2`（参见 `controller/cluster_node.go::DeleteClusterNode`）。
> The model exposes both `status=3` ("disabled by admin") and a `disabled` boolean. `status` reflects runtime state; `disabled` reflects admin intent. The delete handler sets `disabled=true` and `status=2` (see `controller/cluster_node.go::DeleteClusterNode`).

## 自注册 / Self-registration

每个节点首次启动时会调用 `SaveLocalNode`：

Every node calls `SaveLocalNode` on first boot:

- 若本地 `cluster_nodes` 已有本机记录，仅刷新 `address`、`node_name`、`status`、`last_heartbeat`、`ping_failures=0`；
  If a record for `node_id == CLUSTER_NODE_ID` already exists locally, only `address`, `node_name`, `status`, `last_heartbeat` and `ping_failures=0` are refreshed.
- 否则以 `CLUSTER_SECRET` 作为初值创建本节点记录。
  Otherwise a record is created with `secret_key = CLUSTER_SECRET`.

这样设计的原因 / Why this is intentional:

1. **可观测性**：管理后台能看到本节点实时地址、心跳与状态，便于排障；
   **Admin visibility** — Settings → Node Management shows the local node's address, status, heartbeat.
2. **传递性发现**：B 收到 A 的 ping 后，A 会返回自己已知的全部节点列表，B 把这些节点 merge 进本地；C 即可通过 B 间接认识 A；
   **Transitive discovery** — the ping response includes the sender's full node list, which B merges locally.
3. **存活信号**：本地 `last_heartbeat` 由 `discoverOnce` 每周期刷新一次，反映本节点自身存活；
   **Liveness signal** — local `last_heartbeat` is refreshed every cycle by `discoverOnce`.

> 五层防护避免自循环：`GetAllRemoteNodes`、`GetAliveNodesForSync` 的 SQL 过滤 `node_id != ?`，`handlePing` 拒绝自 ping，`mergeDiscoveredNodes` 跳过自身，`ApplyEvents` 跳过本机发出的事件。完整说明见仓库 `README.md`。
> Five layers of guards prevent loops: SQL filters in `GetAllRemoteNodes` / `GetAliveNodesForSync`, self-ping rejection in `handlePing`, self-skip in `mergeDiscoveredNodes`, and self-skip in `ApplyEvents`. See the repo README for the full diagram.

## 远程节点的添加 / Adding remote nodes

管理后台 **Settings → Node Management**（`web/default-pro/src/views/setting/ClusterSetting.vue`）调用 `/api/cluster_node/`：

The admin UI **Settings → Node Management** (`web/default-pro/src/views/setting/ClusterSetting.vue`) calls `/api/cluster_node/`:

| 操作 / Action | 方法 / Method | 路径 / Path | 说明 / Notes |
| --- | --- | --- | --- |
| 列表 / List | `GET` | `/api/cluster_node/` | 全部未禁用节点；本节点 `is_self=true` |
| 详情 / Detail | `GET` | `/api/cluster_node/:id` | 单个节点详情 |
| 新增 / Add | `POST` | `/api/cluster_node/` | `node_id` 在 `1-49`，`secret` 必填（对方 secret） |
| 更新 / Update | `PUT` | `/api/cluster_node/` | 改 name / address / secret；改 secret 自动恢复 `status=1` |
| 禁用（软删）/ Disable | `DELETE` | `/api/cluster_node/:id` | 设置 `disabled=true` |
| 重新启用 / Enable | `POST` | `/api/cluster_node/:id/enable` | 重置 `disabled=false` 并刷新心跳 |
| 手动 Ping / Manual ping | `GET` | `/api/cluster_node/ping/:id` | 立即发起一次 ping，便于排障 |

所有接口需要 **Root 权限**（`middleware.RootAuth`）。

All endpoints require **Root privileges** (`middleware.RootAuth`).

## Per-node Secret

每个节点持 **自己的** secret，由本地 `cluster_nodes.secret_key` 持久化，取代早期的全局共享密钥设计：

Each node carries **its own** secret, persisted in `cluster_nodes.secret_key`, replacing the earlier global-shared-secret design:

- **安全性**：一个节点 secret 泄露不影响其他节点；
  **Security** — one node's secret leaking does not affect others.
- **灵活性**：每个节点可独立轮换 secret；
  **Flexibility** — every node can rotate its own secret independently.
- **自动传递**：ping 响应会携带全部节点 secret，实现跨节点同步；
  **Auto-discovery** — the ping response includes every node's secret, so peers learn each other's secrets.

生命周期 / Lifecycle:

1. **首次启动**：使用 `CLUSTER_SECRET` 环境变量作为 `secret_key` 初值；
   **First boot** — uses the `CLUSTER_SECRET` env var as the initial value.
2. **后续启动**：从 `cluster_nodes.secret_key` 读取，不再依赖环境变量；
   **Subsequent boots** — read from `cluster_nodes.secret_key`; the env var is no longer consulted.
3. **轮换**：通过管理后台 `PUT /api/cluster_node/` 修改 `secret`；下一次 ping 会把新值同步给其他节点；
   **Rotation** — change the `secret` field via the admin UI; the next ping propagates the new value to peers.
4. **校验**：ping/sync 请求头 `X-Cluster-Secret` 等于 **目标节点的** secret（从本地 DB 查）。
   **Verification** — the `X-Cluster-Secret` header equals the **target node's** secret (looked up locally).

新增节点步骤 / Adding a new node:

1. 在 A 上添加 B 的记录（输入 B 的 `CLUSTER_SECRET`）；
   Add Node B's record on Node A (input B's `CLUSTER_SECRET`).
2. 在 B 上添加 A 的记录（输入 A 的 `CLUSTER_SECRET`）；
   Add Node A's record on Node B (input A's `CLUSTER_SECRET`).
3. A → B 的 ping 用 B 的 secret，B 用自己的 secret 校验通过；
   A → B ping uses B's secret, which B verifies with its own secret.
4. B 的响应同时携带 A 与 B 的 secret，A 更新本地副本。
   B's response carries both A's and B's secrets; A updates its local copy.

## 软删除 / Soft delete

`DELETE /api/cluster_node/:id` **不物理删除**记录，而是设置 `disabled = true`：

`DELETE /api/cluster_node/:id` does **not** hard-delete — it sets `disabled = true`:

- 防止被删节点通过 ping 重新注册；
  Prevents a deleted node from "regrowing" via ping-driven re-registration.
- 被禁用节点仍会响应 ping（让对方知道它在线），但不会向其推送事件；
  Disabled nodes still respond to pings (so peers know they're online) but won't fetch this node's info.
- 物理删除需要手动 SQL：`DELETE FROM cluster_nodes WHERE node_id = ?;`
  Hard delete requires manual SQL: `DELETE FROM cluster_nodes WHERE node_id = ?;`

启用通过 `POST /api/cluster_node/:id/enable`，会自动重置 `ping_failures=0` 与 `last_heartbeat`。

Re-enable via `POST /api/cluster_node/:id/enable`; this also resets `ping_failures=0` and refreshes `last_heartbeat`.

## API 速查 / API cheatsheet

| 路径 / Path | 方法 / Method | 用途 / Purpose |
| --- | --- | --- |
| `/api/cluster/ping` | `POST` | 节点间心跳（内部） |
| `/api/cluster/sync` | `POST` | 节点间事件推送（内部） |
| `/api/cluster_node/` | `GET` / `POST` / `PUT` | 节点列表 / 新增 / 更新 |
| `/api/cluster_node/:id` | `GET` / `DELETE` | 节点详情 / 禁用 |
| `/api/cluster_node/:id/enable` | `POST` | 重新启用 |
| `/api/cluster_node/ping/:id` | `GET` | 手动 ping（管理用） |

接口契约详见 [集群 API](/zh/api/cluster)。
Full contract lives in [Cluster API](/en/api/cluster).

下一步 / Next: [配置同步](/zh/decentralization/config-sync) · [节点健康](/zh/decentralization/node-health)。
