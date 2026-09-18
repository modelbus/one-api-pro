---
title: Admin Dashboard
description: "Aggregate KPIs, charts and user leaderboard."
category: admin
order: 1
---

# Admin Dashboard

> Whole-site KPI overview, day × model charts, active-user leaderboard and revenue panels.

Route: `/admin/dashboard`, served by `web/default-pro/src/views/admin/AdminDashboard.vue`.
The page is visible to every admin (`role >= 10`); regular users are redirected to `/dashboard` by the router guard.

## Endpoints

All three endpoints live under `/api/admin/dashboard` and are protected by `AdminAuth` middleware.

| Endpoint | Query | Description |
|---|---|---|
| `GET /api/admin/dashboard/overview` | `range=today\|7d\|30d\|all` (default `7d`) | Single-call KPI aggregate (no chart rows) |
| `GET /api/admin/dashboard/charts`   | `range=today\|7d\|30d\|all` (`all` capped to the last 30 days) | `[]LogStatistic` aggregated by `day × model_name` |
| `GET /api/admin/dashboard/top-users` | `range=today\|7d\|30d\|all`, `limit=1..200` (default 20) | `[{ id, username, email, request_count, quota, balance, current_plan_name }]` |

Implementation: `controller/admin/dashboard.go` + `model/admin_dashboard.go`.

## Overview Payload

`overview` returns `AdminDashboardOverview`:

| Group | Key fields | Notes |
|---|---|---|
| `users` | `total / enabled / disabled / deleted / new_today / new_7d / new_30d / active_7d` | Bucketed by `status`; "new_X" windows are `created_at >= start_of_(day/week/month)`; `active_7d = DISTINCT user_id FROM logs WHERE type=2 AND created_at >= 7d AND users.request_count > 0` |
| `tokens` / `channels` / `plans` | `total / enabled` | Resource totals; `enabled` = `status=1` |
| `redemptions` | `total / used / unused` | Three redemption status counts |
| `subscriptions` | `total / active / expired` | Aggregated from `user_plans` |
| `quota` | `today / week / month / total` | Computed via `SumUsedQuota(LogTypeConsume, …)` with day/week/month/all windows |
| `revenue` | `total / topup / subscription / refund` | `orders` table: `status=1` split by `type` (`OrderTypeTopup=2` / `OrderTypePlanSubscription=1`); refunds are summed where `status=3` |
| `range` / `generated_at` | metadata | `range` is the canonical window key; `generated_at` is a unix-second timestamp |

## Charts Payload

`charts` mirrors `model.SearchLogsByDayAndModel` and returns `[]LogStatistic`:

```ts
interface LogStatistic {
  day: string              // YYYY-MM-DD, truncated to start of day
  model_name: string
  request_count: number    // only LogTypeConsume rows
  quota: number            // sum(quota)
  prompt_tokens: number
  completion_tokens: number
}
```

The shape is identical to `GET /api/user/dashboard`; the frontend reuses the same chart-building logic (`AdminDashboard.vue::chartData`). SQL is written to work on MySQL / PostgreSQL / SQLite via dialect-specific date formatting.

## Top Users

Scope: users with at least one `type=2` log in the chosen window AND `users.request_count > 0`, ordered by `request_count DESC, quota DESC`. `limit` accepts 1–200 and falls back to 20 on overflow. The "current plan name" is filled in by `batchTopActivePlanNames(userIds)` in one batched query — no N+1 — and only the earliest-expiring active subscription is shown.

## Frontend Panels

`AdminDashboard.vue` shows a top welcome bar with a range radio (`today / 7d / 30d / all`) and a left 16 / right 8 grid:

- **Left column (16/24)** — stacked in order:
  - **Users KPI**: 8 cards (total / new today / new 7d / new 30d / active 7d / disabled / deleted / today ratio).
  - **Revenue KPI**: 4 cards (total / topup / subscription / refund) in `¥`.
  - **Resources KPI**: 5 cards (tokens / channels / plans / redemptions / subscriptions); the "plans" and "redemptions" cards are clickable shortcuts into the management pages.
  - **Quota + 3 charts**: 4 KPI cards (today / 7d / 30d / total) plus 3 line charts (request count / quota / tokens).
  - **Model distribution**: stacked bar chart of the Top 8 models by `quota`.
  - **Usage details**: flat `day × model` table, sorted by date DESC then `quota` DESC.
- **Right column (8/24)** — promo placeholder / leaderboard card (with `today / 7d / 30d` switch, shows rank / username / email / request count / quota / current plan) / announcements / changelog / contact links.

## Tips

- Switching `range` reloads all three endpoints; the chip in the top-right shows the last successful `overview` refresh time.
- All three line charts (requests / quota / tokens) and the model distribution share the same `[]LogStatistic` payload — no extra round-trips.
- `users.total` is `COUNT(*)` over the users table; soft-deleted rows are bucketed separately in `users.deleted`.

## Implementation Pointers

| Concern | Location |
|---|---|
| Overview aggregate SQL | `model/admin_dashboard.go::GetAdminDashboardOverview` |
| Revenue aggregation | `model/admin_dashboard.go::aggregateRevenue` |
| Day × model aggregation | `model/admin_dashboard.go::SearchAdminLogsByDayAndModel` |
| Range parsing | `model/admin_dashboard.go::ParseAdminDashboardRange` / `ParseAdminChartsRange` |
| Top users + active plan | `model/admin_dashboard.go::GetAdminTopUsers` / `batchTopActivePlanNames` |
| HTTP handler | `controller/admin/dashboard.go` |
| Route registration | `router/api.go` (`/api/admin/dashboard/*`) |