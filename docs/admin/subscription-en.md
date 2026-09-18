---
title: Subscriptions Admin
description: "Grant, query, edit and revoke user subscriptions (user_plans)."
category: admin
order: 6
---

# Subscriptions Admin

> Grant, query, adjust and revoke user subscriptions on the `user_plans` table; admins can bypass payment and activate immediately.

Route: `/subscription`. The page is shared between admins and users (`web/default-pro/src/views/subscription/Subscription.vue`); the "add subscription" button is gated by `authStore.isAdmin`.

## Data Model

`model.UserPlan` (table `user_plans`):

| Field | Type | Notes |
|---|---|---|
| `user_id` | `int` (indexed) | User id |
| `plan_id` | `int` | Plan id (required; `Insert()` rejects `plan_id=0` to prevent orphan rows) |
| `order_id` | `int` | Linked order id (`0` for admin grants) |
| `start_time` / `end_time` | `int64` | unix seconds; `end_time <= now` is considered expired |
| `status` | `int` | `UserPlanStatusActive=1` / `UserPlanStatusExpired=0` |
| `billing_type` | `varchar(20)` | `token` or `request`, mirrors the plan |
| `notes` | `text` | Free-form notes ("管理员开通") |

The background task `ExpireUserPlans()` flips `status=1 AND end_time <= now` rows to `expired`.

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/subscription/self` | `GET` | User | The current user's subscriptions (any status) |
| `/api/subscription/` | `GET` | Admin | Paginated list; supports `?user_id=` and `?status=` filters |
| `/api/subscription/search?keyword=` | `GET` | Admin | `username LIKE keyword%` |
| `/api/subscription/:id` | `GET` | Admin | Subscription detail (with embedded `Plan`) |
| `/api/subscription/:id/usage` | `GET` | User | Current subscription's limits / usage / weighted progress / next reset |
| `/api/subscription/` | `POST` | Admin | Admin grant (see below) |
| `/api/subscription/` | `PUT` | Admin | Adjust `end_time` / `status` / `billing_type` / `notes` |
| `/api/subscription/:id` | `DELETE` | Admin | Revoke |

Implementation: `controller/subscription.go`.

## Admin Grant

```json
{
  "user_id": 42,
  "plan_id": 7,
  "billing_type": "token",
  "duration_days": 0,        // optional; 0 = inherit plan default
  "notes": "商务合作赠送",
  "pay_method": "free"       // wechat | alipay | bank | offline | free
}
```

Server-side:

1. Validates `pay_method` (default `free`).
2. Validates the plan exists and `status=1`.
3. Calls `model.CreatePlanOrder` to insert an `Order` with `source=admin` as the audit trail.
4. **Immediately** calls `model.ActivatePackageByOrder(order, OrderUpgradeModeStack)` — admin grants take effect right away regardless of `pay_method` (`free` / `offline` / `bank` / `wechat` / `alipay`).
5. Returns the new `user_plan` (with the embedded `Plan`) together with the order.

Admin grants always use `OrderUpgradeModeStack` — price_diff logic does not apply.

## Adjust Subscription

```json
{
  "id": 88,
  "end_time": 1767225600,
  "status": 1,
  "billing_type": "token",
  "notes": "延期 30 天"
}
```

After success `model.CacheDeleteUserActivePlans(userId)` invalidates the Redis cache (`user_plans:<id>`); the next request reloads fresh data.

## Revoke Subscription

`DELETE /api/subscription/:id` is a hard delete on the `user_plan` row and clears the user's cache. The plan's quota caps / usage stop applying immediately — there is no automatic reset.

## `/api/subscription/:id/usage`

Response:

```jsonc
{
  "subscription": { /* user_plan + plan */ },
  "usage": [ /* PlanUsage rows: per (model, window_type, window_index) */ ],
  "weighted": { "period": 0.42, "week": 0.15, "month": 0.08 },
  "limits":   { "gpt-4o": { request_period: 100, token_period: 50000, period_h: 5, … }, … },
  "model_usage": { /* per model, per window: used vs limit */ },
  "next_reset": { "period": 1767225600, "week": 1767225600, "month": 1767225600 },
  "now": 1735660800,
  "start_time": 1733059200,
  "billing_type": "token"
}
```

`weighted` is the result of `model.CalculateWeightedUsage` which folds per-request and per-token consumption into a single ratio, suitable for a usage progress bar.

## Frontend Guide

- Columns: ID / user (admin only) / plan name / billing type / start / end / status / actions.
- Row actions:
  - **Usage**: opens the usage view / modal showing the `weighted` progress.
  - **Renew** (admin): a modal that lets you push `end_time` forward by `duration_days`.
  - **Expire** (admin): confirmation popover, then `PUT` with `status=0`.
  - **Delete**: hard delete (irreversible).
- The **Add Subscription** button is admin-only. Modal fields: `user_id`, `plan_id` (dropdown from `/api/plan/`), `billing_type` (token / request), `pay_method` (free / offline / wechat / alipay / bank), notes.

## Implementation Pointers

| Concern | Location |
|---|---|
| Subscription CRUD | `controller/subscription.go` |
| Limits / usage | `model/plan.go::GetUserSubscriptionInfo` / `CalculateWeightedUsage` / `CalcModelUsageDetails` / `CalcNextResetTime` |
| Auto-expiry | `model/plan.go::ExpireUserPlans` |
| Cache invalidation | `model/plan.go::CacheDeleteUserActivePlans` |
| Routes | `router/api.go` |