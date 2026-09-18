---
title: Config Sync
description: "Event dispatch, watcher, applier and conflict resolution."
category: decentralization
order: 3
---

# Config Sync

> Event dispatch, watcher, applier and conflict resolution.

Data flow: `business SQL → GORM callback → sync_events → pusher goroutine → remote /sync → applier`.

## Synced tables

Defined in `cluster/event.go::syncTableList`:

```text
users, tokens, channels, abilities,
options, redemptions, plans,
user_plans, plan_usages,
channel_counters,
logs (gated by CLUSTER_SYNC_LOGS)
```

Internal tables (`cluster_nodes`, `sync_events`, ...) are **not** synced.

## Watcher: capture changes

`cluster/watcher.go` registers four GORM callbacks that fire only on non-internal tables:

| Callback | Trigger | Event written |
| --- | --- | --- |
| `cluster:after_create` | after GORM `Create` | `insert` |
| `cluster:after_update` | after GORM `Update` | `update` |
| `cluster:before_delete` | before GORM `Delete` | capture old row JSON |
| `cluster:after_delete` | after GORM `Delete` | `delete` |

Each event writes one row to `sync_events`:

```text
sync_events (
  id, table_name, row_id, row_key, action,
  data TEXT, node_id, created_at, pushed
)
```

- `row_id`: `id` for single-PK tables; empty for composite-PK tables (`channel_counters`, `abilities`) where the full row lives in `data`.
- `row_key`: only populated for the `options` table.
- `node_id`: stamped with the current node so the receiver can deduplicate.
- `pushed=0`: waiting for the pusher; deleted in bulk after a successful push.

> The `sync_events` insert is performed via `WithSkipHook(db).Session(NewDB: true).Exec(...)`, **completely bypassing GORM callbacks** — this avoids the `createSyncEvent → DB.Create → afterCreate → createSyncEvent` infinite recursion.

## Pusher goroutine

`cluster/pusher.go` runs a goroutine listening on `syncNotifyChan` (capacity 256):

```text
watcher writes event → NotifySyncEvent() → syncNotifyChan (non-blocking, drop on full)
                                ↓
                       StartPusher picks up signal
                                ↓
                  100 ms throttle + drainNotifyChan
                                ↓
                  pushEvents: read pushed=0 rows (ORDER BY created_at, LIMIT BatchSize)
                                ↓
           for each event, POST to every alive node's /api/cluster/sync (WithSkipHook)
                                ↓
            all success → bulk DELETE FROM sync_events WHERE id IN (...)
```

Highlights:

- **Throttling**: 100 ms sleep + drain merges multiple changes in the same window into one push.
- **Batching**: `BatchSize` (default 50) caps each request's payload.
- **Alive nodes only**: `GetAliveNodesForSync` filters to `status=1 AND disabled=false AND node_id != ?`.
- **No delete on partial failure**: if any node fails, the event stays so the next cycle retries.
- **Catch-up**: on startup, count rows with `pushed=0` and trigger an immediate catch-up push.

Headers:

- `X-Cluster-Secret`: target node's secret (looked up locally)
- `X-Cluster-Node-Id`: current `CLUSTER_NODE_ID`
- 5-second timeout

## Receiver: Applier

`cluster/handler.go::handleSync` calls `ApplyEvents` inside a goroutine and returns 200 immediately:

```go
go func() {
    defer func() { recover() ... }()   // defensive: panic recovery
    ApplyEvents(req.Events)
}()
c.JSON(http.StatusOK, gin.H{"success": true})
```

`ApplyEvents` dispatches by `action`:

| `action` | Handler | Notes |
| --- | --- | --- |
| `insert` | `applyInsert` | falls back to update if PK already exists (idempotent) |
| `update` | `applyUpdate` | compares `updated_at`; only newer rows are written |
| `delete` | `applyDelete` | composite-PK tables use composite predicates; falls back to `data.id` if `row_id` is missing |

Special tables:

- `options`: upsert by `key`.
- `abilities`: composite PK `(group, model, channel_id)`.
- `channel_counters`: composite PK `(channel_id, node_id)`.

All writes go through `WithSkipHook(db)` — **the receiver never re-triggers the watcher**.

## Conflict resolution

`applyUpdate` compares the incoming `updated_at` with the local one:

```text
incoming.updated_at <= local.updated_at   → skip (local is newer or equal)
incoming.updated_at >  local.updated_at   → overwrite
```

Last-Writer-Wins.

> All nodes should keep their system clocks roughly in sync (NTP recommended). Large skew can cause newer updates to be rejected.

## Cache invalidation

`applier.go::invalidateCache` evicts Redis cache after applying an event (only when Redis is enabled):

| Table | Keys evicted |
| --- | --- |
| `users` | `user_group:<id>`, `user_quota:<id>`, `user_enabled:<id>`, `user_plans:<id>` |
| `tokens` | `token:<key>` |
| `channels` | reload channel cache; delete all `group_models:*` |
| `options` | reload `OptionMap` |
| `user_plans` | `user_plans:<user_id>` |

## Event cleanup

`cluster/pusher.go::StartEventCleanup` deletes `sync_events` rows older than 7 days every hour to keep the table bounded.

Next: [Node Health](/en/decentralization/node-health) · [Multi-node Deployment](/en/decentralization/deployment).
