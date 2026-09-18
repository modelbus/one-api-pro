---
title: Logs Admin
description: "Site-wide log query (type / model / user / channel / time window), aggregate and prune."
category: misc
order: 11
---

# Logs Admin

> Filter the `logs` table (admin-wide or per-user), aggregate quota by source, and prune historical rows. UI: `web/default-pro/src/views/log/Log.vue` (shared between admin and user; columns are hidden based on role).

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/log/` | `GET` | Admin | Paginated (`config.ItemsPerPage`) |
| `/api/log/search?keyword=` | `GET` | Admin | Keyword search |
| `/api/log/stat` | `GET` | Admin | Quota aggregate (normal + subscription) |
| `/api/log/?target_timestamp=` | `DELETE` | Admin | Delete all rows with `created_at < target` |
| `/api/log/self` | `GET` | User | Caller's own logs |
| `/api/log/self/search?keyword=` | `GET` | User | Caller's own keyword search |
| `/api/log/self/stat` | `GET` | User | Caller's own quota aggregate |

Implementation: `controller/log.go`.

## Common Filter

| Param | Type | Notes |
|---|---|---|
| `p` | `int` | Page index (0-based) |
| `type` | `int` | `LogTypeTopup=1` / `Consume=2` / `Manage=3` / `System=4` / `Test=5` |
| `start_timestamp` | `int64` | unix seconds (inclusive); `0` = no lower bound |
| `end_timestamp` | `int64` | unix seconds (inclusive); `0` = no upper bound |
| `username` | `string` | exact match (admin endpoints) |
| `token_name` | `string` | fuzzy match |
| `model_name` | `string` | fuzzy match |
| `channel` | `int` | `channel_id` (admin endpoints) |

Only `LogTypeConsume` rows contribute to billed quota; other types are audit only.

## Aggregate — `/api/log/stat`

```json
{
  "quota": 1234567,
  "normal_quota": 1000000,
  "subscription_quota": 234567
}
```

- `quota` — total `sum(quota)` over the window.
- `normal_quota` — sum where `billing_source=0`.
- `subscription_quota` — sum where `billing_source=1`.

Implemented in `model/log.go::SumUsedQuota` and `SumUsedQuotaByBillingSource`.

## Prune — `DELETE`

`DELETE /api/log/?target_timestamp=<unix-seconds>` removes all rows with `created_at < target_timestamp`. Returns `{ data: <affected_rows> }`. `target_timestamp=0` is rejected. Pair with the date-picker UI in Operation Setting for an interactive prune.

## Log Types

| Type | Meaning | Counts toward quota? |
|---|---|---|
| `LogTypeTopup=1` | Top-up / redemption / admin manual grant | No |
| `LogTypeConsume=2` | API call consumption | Yes |
| `LogTypeManage=3` | Admin actions (delete/disable/quota adjustment/grant subscription) | No |
| `LogTypeSystem=4` | System (new-user bonus, invite bonus, …) | No |
| `LogTypeTest=5` | Channel tests (`POST /api/channel/test`) | No |

## Frontend Guide

- Welcome bar + a `type` select (all / topup / consume / manage / system / test).
- Search bar:
  - Regular users: keyword, token_name, model_name, start/end timestamp.
  - Admins: + username, channel id.
- The **Statistics** button toggles the `quota / normal_quota / subscription_quota` summary card.
- Row columns: time (click to copy `request_id`) / channel (admin) / source chip (subscription blue / quota grey) / plan / type / model / user (admin) / token name. For non-test rows: Prompt / Completion / Quota. Detail column shows `content`, `elapsed_time` (ms; coloured by latency bucket), `is_stream`, `system_prompt_reset` tags.
- Pagination `[10, 20, 50]`; append-on-scroll.
- Operation Setting (`/setting/operation`) exposes a **Log Pruning** action — pick a date → `DELETE /api/log/?target_timestamp=...`.

## Implementation Pointers

| Concern | Location |
|---|---|
| Log handler | `controller/log.go` |
| `LogStatistic` type | `model/log.go` |
| Quota aggregate | `model/log.go::SumUsedQuota` / `SumUsedQuotaByBillingSource` |
| Delete | `model/log.go::DeleteOldLog` |
| Prune UI | `web/default-pro/src/views/setting/OperationSetting.vue` |
| Routes | `router/api.go` |