---
title: Config Sync
description: How data syncs between nodes, when conflicts happen, and how they converge.
category: decentralization
order: 3
---

# Config Sync

> Node A changed a row. When do nodes B, C, D see it? Answer: usually < 1 second.

## Sync flow

```
Node A writes to its DB
    ↓
GORM callbacks auto-capture the change
    ↓
Row written to sync_events (temp table)
    ↓
Pusher goroutine batch-pushes to each live node
    ↓
Receiving node's Applier applies the event
```

Any node's change becomes visible on the others within ~1 second.

## What's synced

- `users` / `tokens`: accounts and tokens
- `channels` / `abilities`: channels
- `options`: system settings
- `redemptions`: redemption codes
- `plans` / `user_plans` / `plan_usages`: plans and subscriptions
- `channel_counters`: per-channel rate-limit counts
- `logs`: call logs (toggle via `CLUSTER_SYNC_LOGS=false`)

**Not synced:** `cluster_nodes` / `sync_events` (internal, maintained by discovery).

## How conflicts are resolved

If two nodes change the same row nearly simultaneously, the **last-writer-wins**:

```
Node A updated_at = 10:00:01.100
Node B updated_at = 10:00:01.050
→ A's version wins
```

Rule on the receiving side:

- Compare incoming row's `updated_at` with local
- If incoming is newer → overwrite; otherwise → discard

> **All nodes' clocks must be roughly in sync (enable NTP).** Large clock skew can cause the system to discard updates that should have been adopted.

## What happens during a node's offline window

**Changes are not backfilled automatically.** Example: node B is offline for 1 hour; node A makes 100 changes during that hour; when B comes back, it sees only new changes going forward.

Fix: restore B's DB from a live peer via `mysqldump` or similar.

## Request headers between nodes

Every push carries:

- `X-Cluster-Secret`: target node's secret (must match)
- `X-Cluster-Node-Id`: sender node ID
- 5-second timeout

## Troubleshooting sync issues

1. **A changed but B didn't receive**:
   - Check B's backend logs for sync records
   - From A, manually Ping B
   - Verify both sides' `secret_key` match
2. **Conflict — lost an update**:
   - Check both nodes' clocks (NTP enabled?)
   - Check `sync_events` table for backlog (`pushed=0` count)
3. **Sync latency > a few seconds**:
   - Check `sync_events` backlog
   - Check Pusher goroutine is alive (process monitoring)
   - Check inter-node network latency

## When to disable sync for speed

- Logs are huge: set `CLUSTER_SYNC_LOGS=false` on each node — logs stay only on the local node
- Node is read-only: set the local `secret_key` to a random value (rejects inbound events)

## Related

- [Cluster Overview](./overview)
- [Node Management](./node-management)
- [Node Health](./node-health)