---
title: 配置同步
description: "事件分发、Watcher、Applier 与冲突解决。"
category: decentralization
order: 3
---

# 配置同步

> 事件分发、Watcher、Applier 与冲突解决。
> Event dispatch, watcher, applier and conflict resolution.

数据同步链路：`业务 SQL → GORM 回调 → sync_events → 推送协程 → 远端 /sync → Applier`。

Data flow: `business SQL → GORM callback → sync_events → pusher goroutine → remote /sync → applier`.

## 同步的表 / Synced tables

由 `cluster/event.go::syncTableList` 控制：

Defined in `cluster/event.go::syncTableList`:

```text
users, tokens, channels, abilities,
options, redemptions, plans,
user_plans, plan_usages,
channel_counters,
logs（受 CLUSTER_SYNC_LOGS 控制）
```

`cluster_nodes`、`sync_events` 等内部表**不会**被同步。

Internal tables (`cluster_nodes`, `sync_events`, ...) are **not** synced.

## Watcher：捕获变更 / Watcher: capture changes

`cluster/watcher.go` 注册了四个 GORM 回调，仅对非内部表生效：

`cluster/watcher.go` registers four GORM callbacks that fire only on non-internal tables:

| 回调 / Callback | 触发时机 / Trigger | 写入事件 / Event |
| --- | --- | --- |
| `cluster:after_create` | GORM `Create` 之后 | `insert` |
| `cluster:after_update` | GORM `Update` 之后 | `update` |
| `cluster:before_delete` | GORM `Delete` 之前 | 抓取旧行 JSON 备用 |
| `cluster:after_delete` | GORM `Delete` 之后 | `delete` |

每个事件都会向 `sync_events` 写入一行：

Each event writes one row to `sync_events`:

```text
sync_events (
  id, table_name, row_id, row_key, action,
  data TEXT, node_id, created_at, pushed
)
```

- `row_id`：单主键表的 `id`；复合主键表（如 `channel_counters`、`abilities`）为空，由 `data` JSON 携带完整字段；
  `row_id`: `id` for single-PK tables; empty for composite-PK tables (`channel_counters`, `abilities`) where the full row lives in `data`.
- `row_key`：仅 `options` 表存 `key`；
  `row_key`: only populated for the `options` table.
- `node_id`：写入时打上当前节点编号，用于接收方识别与去重；
  `node_id`: stamped with the current node so the receiver can deduplicate.
- `pushed=0`：等待 Pusher 推送；推送成功后由 Pusher 批量删除。
  `pushed=0`: waiting for the pusher; deleted in bulk after a successful push.

> 关键：写入 `sync_events` 用 `WithSkipHook(db).Session(NewDB: true).Exec(...)` 直接执行 SQL，**完全绕开 GORM 回调**，避免 `createSyncEvent → DB.Create → afterCreate → createSyncEvent` 的无限递归。
> The `sync_events` insert is performed via `WithSkipHook(db).Session(NewDB: true).Exec(...)`, **completely bypassing GORM callbacks** — this avoids the `createSyncEvent → DB.Create → afterCreate → createSyncEvent` infinite recursion.

## Pusher：推送协程 / Pusher goroutine

`cluster/pusher.go` 启动一个 goroutine 监听 `syncNotifyChan`（带 256 容量）：

`cluster/pusher.go` runs a goroutine listening on `syncNotifyChan` (capacity 256):

```text
watcher 写入事件 → NotifySyncEvent() → syncNotifyChan (非阻塞满则丢弃)
                                ↓
                       StartPusher 收到信号
                                ↓
                  100ms 节流 + drainNotifyChan
                                ↓
                  pushEvents: 读取 pushed=0 的事件 (按 created_at, Limit=BatchSize)
                                ↓
           对每个事件用 WithSkipHook 直接 HTTP POST 到每个存活节点 /api/cluster/sync
                                ↓
            全部成功 → 批量 DELETE FROM sync_events WHERE id IN (...)
```

要点 / Highlights:

- **节流**：每次通知后 sleep 100ms 并 drain 通道，把同一窗口内的多次变更合并推送；
  **Throttling**: 100 ms sleep + drain merges multiple changes in the same window into one push.
- **批量**：`BatchSize`（默认 50）限制单次推送事件数，防止请求过大；
  **Batching**: `BatchSize` (default 50) caps each request's payload.
- **仅存活节点**：`GetAliveNodesForSync` 只取 `status=1 AND disabled=false AND node_id != ?` 的节点；
  **Alive nodes only**: `GetAliveNodesForSync` filters to `status=1 AND disabled=false AND node_id != ?`.
- **失败不删除**：单个节点推送失败不会把该事件标为已推送，等下一周期重试；
  **No delete on partial failure**: if any node fails, the event stays so the next cycle retries.
- **重传兜底**：启动时统计 `pushed=0` 的事件数并触发一次补推。
  **Catch-up**: on startup, count rows with `pushed=0` and trigger an immediate catch-up push.

请求头 / Headers:

- `X-Cluster-Secret`: 目标节点的 secret（本地 DB 查）
- `X-Cluster-Node-Id`: 当前节点的 `CLUSTER_NODE_ID`
- 超时 5 秒
- 5-second timeout

## 接收方：Applier / Receiver: Applier

`cluster/handler.go::handleSync` 接收到事件后，**在 goroutine 内**调用 `ApplyEvents`，立即返回 200：

`cluster/handler.go::handleSync` calls `ApplyEvents` inside a goroutine and returns 200 immediately:

```go
go func() {
    defer func() { recover() ... }()   // 防御性：panic 不拖垮进程
    ApplyEvents(req.Events)
}()
c.JSON(http.StatusOK, gin.H{"success": true})
```

`ApplyEvents` 按 `action` 分发：

`ApplyEvents` dispatches by `action`:

| `action` | 处理 / Handler | 说明 / Notes |
| --- | --- | --- |
| `insert` | `applyInsert` | 若主键已存在则转为 update；保证幂等 |
| `update` | `applyUpdate` | 比较 `updated_at`，仅当入站更新才写入 |
| `delete` | `applyDelete` | 复合主键表走复合条件；缺 `row_id` 时尝试从 `data` JSON 兜底 |

特殊表 / Special tables:

- `options`：按 `key` upsert；
  `options`: upsert by `key`.
- `abilities`：复合主键 `(group, model, channel_id)`；
  `abilities`: composite PK `(group, model, channel_id)`.
- `channel_counters`：复合主键 `(channel_id, node_id)`；
  `channel_counters`: composite PK `(channel_id, node_id)`.

所有写入通过 `WithSkipHook(db)`，**避免接收方再次触发 Watcher 产生新事件**。

All writes go through `WithSkipHook(db)` — **the receiver never re-triggers the watcher**.

## 冲突解决 / Conflict resolution

`applyUpdate` 比较 `incoming.updated_at` 与本地 `updated_at`：

`applyUpdate` compares the incoming `updated_at` with the local one:

```text
incoming.updated_at <= local.updated_at   → 跳过（本地更新或相同）
incoming.updated_at >  local.updated_at   → 覆盖
```

最后写入胜出（Last-Writer-Wins）。

Last-Writer-Wins.

> 注意：所有节点需保证系统时钟基本一致（建议启用 NTP）。漂移过大会导致本应被采纳的更新被错误丢弃。
> All nodes should keep their system clocks roughly in sync (NTP recommended). Large skew can cause newer updates to be rejected.

## 缓存失效 / Cache invalidation

`applier.go::invalidateCache` 在事件应用后清理 Redis 缓存（仅当 Redis 启用时生效）：

`applier.go::invalidateCache` evicts Redis cache after applying an event (only when Redis is enabled):

| 表 / Table | 失效键 / Key |
| --- | --- |
| `users` | `user_group:<id>`, `user_quota:<id>`, `user_enabled:<id>`, `user_plans:<id>` |
| `tokens` | `token:<key>` |
| `channels` | 重载渠道缓存；删除全部 `group_models:*` |
| `options` | 重新加载 `OptionMap` |
| `user_plans` | `user_plans:<user_id>` |

## 日志清理 / Event cleanup

`cluster/pusher.go::StartEventCleanup` 每小时清理 `created_at < now - 7d` 的 `sync_events` 行，避免表膨胀。

`cluster/pusher.go::StartEventCleanup` deletes `sync_events` rows older than 7 days every hour.

下一步 / Next: [节点健康](/zh/decentralization/node-health) · [多节点部署](/zh/decentralization/deployment)。
