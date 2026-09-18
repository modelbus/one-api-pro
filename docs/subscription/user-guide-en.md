---
title: User Subscription Guide
description: "What users see on /subscription and /subscription/:id/usage — progress bars, remaining time, per-model detail."
category: subscription
order: 2
---

# User Subscription Guide

> This page documents the user-facing fields; the mechanics live in [Subscription Overview](./overview).

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/subscription/self` | `GET` | User | All `user_plans` for the caller (any status, up to 1000) |
| `/api/subscription/info` | `GET` | User | Aggregated view via `model.GetUserSubscriptionInfo`: one entry per (plan, model, windowType) |
| `/api/subscription/current` | `GET` | User | Currently active plan (nearest non-expired `end_time`); flattens the embedded `Plan` |
| `/api/subscription/:id/usage` | `GET` | User | Per-subscription detail: raw `plan_usages`, `weighted` usage, `model_usage` per-model rows, `next_reset` |

Implementation: `controller/subscription.go`.

## `/api/subscription/current` fields

The server flattens the embedded `Plan`:

```json
{
  "id": 88,
  "user_id": 42,
  "plan_id": 7,
  "order_id": 1735,
  "start_time": 1733059200,
  "end_time": 1767225600,
  "expire_time": 1767225600,
  "expire_date": "2026-01-01",
  "status": 1,
  "billing_type": "token",
  "notes": "管理员开通",
  "created_time": 1733059200,
  "updated_time": 1733059200,
  "is_expired": false,
  "remaining_days": 90,
  "name": "Pro 月卡",
  "price": 99.00,
  "duration_days": 30,
  "duration_text": "30 天",
  "sort": 50,
  "recommended": true,
  "description": "...",
  "features": ["GPT-4o 200K tokens/天", "..."],
  "model_limits": { "gpt-4o": { ... } },
  "default_model": "gpt-4o"
}
```

- `is_expired` = `end_time <= now`
- `remaining_days` = `ceil((end_time - now) / 86400)` (0 when expired)
- `features` is a `StringSlice`, persisted as a JSON array

## `/api/subscription/:id/usage` fields

```jsonc
{
  "subscription": { /* user_plan + embedded plan */ },
  "usage":         [ /* PlanUsage rows */ ],
  "weighted":      { "period": 42.0, "week": 15.0, "month": 8.0 },   // percentage
  "limits":        { "gpt-4o": { "request_period": 100, "token_period": 50000, "period_h": 5, ... } },
  "model_usage":   { "period": [ { "model": "gpt-4o", "requests": 42, "token_percent": 21.5, "request_percent": 42.0 } ], ... },
  "next_reset":    { "period": 1767225600, "week": 1767225600, "month": 1767225600 },
  "now":           1735660800,
  "start_time":    1733059200,
  "billing_type":  "token"
}
```

- `weighted.*` is the output of `model.CalculateWeightedUsage`; values are percentages against `QuotaPoolCapacity=100`. Any value `>= 100` makes `CheckPlanQuota` treat the plan as exhausted.
- `model_usage[*]` is `model.CalcModelUsageDetails`: per (model, windowType) rows showing current-window usage and percentage.
- `next_reset` is computed by `model.CalcNextResetTime`: period uses `plan.model_limits[<default>].period_h` (default 5), week=7d, month=30d.

## Frontend Guide

Route: `/subscription` (`web/default-pro/src/views/subscription/Subscription.vue`). Columns:

- ID / plan name / billing_type / start_time / end_time / status / actions
- Row actions:
  - **Usage**: calls `/api/subscription/:id/usage` and renders the weighted progress bar plus per-model breakdown
  - **Renew**: triggers `POST /api/order/plan`; the server picks `price_diff` or `stack` based on `plan.upgrade_mode`
  - **Cancel**: admin only
- Empty-state copy differs for admins vs users
- The **Add Subscription** button is gated by `authStore.isAdmin`

## Implementation Pointers

| Concern | Location |
|---|---|
| User list | `controller/subscription.go::GetUserSubscriptions` |
| Current plan | `controller/subscription.go::GetCurrentPlan` |
| User aggregated info | `model/plan.go::GetUserSubscriptionInfo` |
| Detailed usage | `controller/subscription.go::GetSubscriptionUsage` |
| Weighted calculation | `model/plan_quota.go::CalculateWeightedUsage` |
| Next reset | `model/plan_quota.go::CalcNextResetTime` |
