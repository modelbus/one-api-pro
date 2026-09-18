---
title: Subscription Admin
description: "Plan CRUD, admin grants, adjust / revoke subscriptions."
category: subscription
order: 3
---

# Subscription Admin

> All admin operations on plans (`plans` table) and user subscriptions (`user_plans` table). Implementation: `controller/plan.go`, `controller/subscription.go`.

## Plans

Data model: `model.Plan`. Fields are documented under [Plan Pricing](./plan) and `model/plan.go`.

### Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/plan/` | `GET` | Admin | Paginated (`config.ItemsPerPage`) |
| `/api/plan/search?keyword=` | `GET` | Admin | `name LIKE keyword%` |
| `/api/plan/:id` | `GET` | Admin | Single row |
| `/api/plan/` | `POST` | Admin | Create |
| `/api/plan/` | `PUT` | Admin | Update |
| `/api/plan/:id` | `DELETE` | Admin | Delete |
| `/api/plan/public` | `GET` | Public | Only `status=PlanStatusEnabled` rows (used by the user-facing subscribe page) |
| `/api/plan/public/:id` | `GET` | Public | Public detail |

### Validation

Before `Insert` / `Update`, the server calls `ValidateDefaultModel`: when `default_model` is non-empty, it must exist as a key in `model_limits`; otherwise it returns:

- `default_model 'xxx' is set but model_limits is empty`
- `default_model 'xxx' is not found in model_limits`

### `model_limits` shape

Persisted as a `text` column holding a JSON object:

```json
{
  "gpt-4o":    { "period_h": 5, "request_period": 100, "request_week": 500, "request_month": 2000, "token_period": 50000,  "token_week": 250000, "token_month": 1000000 },
  "claude-3.5-sonnet": { "period_h": 5, "request_period": 50, "request_week": 200, "request_month": 800, "token_period": 30000, "token_week": 120000, "token_month": 480000 }
}
```

- `period_h` drives the period window index (`CalcWindowIndex`); defaults to 5 when missing
- `request_*` / `token_*` are per-window caps; `0` means unlimited for that window × dimension

### features

`Features StringSlice` field (`model/plan.go`): a JSON array of strings, surfaced as a feature list in the UI.

## User subscriptions

### Force grant (bypasses payment)

`POST /api/subscription/` (`controller/subscription.go::AddSubscription`) request body:

```json
{
  "user_id": 42,
  "plan_id": 7,
  "billing_type": "token",
  "duration_days": 0,
  "notes": "商务合作赠送",
  "pay_method": "free"
}
```

Flow:

1. Validate `user_id` / `plan_id` non-empty; `pay_method` defaults to `free`; whitelist: `wechat` / `alipay` / `bank` / `offline` / `free`
2. Validate the plan exists and `status=PlanStatusEnabled`; validate the user exists
3. Call `model.CreatePlanOrder` to insert an `Order` row (`source=admin`, audit trail); `notes` defaults to "管理员开通"
4. **Immediately** call `model.ActivatePackageByOrder(order, OrderUpgradeModeStack)` to grant the subscription — admin grants always take effect right away regardless of `pay_method`
5. Returns the new `user_plan` (with embedded `Plan`) plus the order row

### Adjust

`PUT /api/subscription/` body (`UpdateSubscriptionRequest`):

```json
{ "id": 88, "end_time": 1767225600, "status": 1, "billing_type": "token", "notes": "延期 30 天" }
```

Non-empty fields are updated; `updated_time` is refreshed. On success, `model.CacheDeleteUserActivePlans(userId)` invalidates the Redis cache.

### Revoke

`DELETE /api/subscription/:id` is a hard delete: it `First`s the row, then deletes it, and clears the cache. The user immediately loses that plan's window quotas (no compensation).

## Public plans

`/api/plan/public` and `/api/plan/public/:id` require no auth and are used by the user-facing subscribe flow. `/api/plan/` (admin) returns rows of any status including disabled ones.

## Scheduled task

`model.ExpireUserPlans` runs periodically (see `main.go`) and flips `status=1 AND end_time <= now` rows to `status=0`, removing them from `CheckPlanQuota`.

## Frontend Guide

- **Plan management**: `/setting/pricing` → "套餐" tab (`web/default-pro/src/views/setting/PlanSetting.vue`). List + modal for editing `name` / `description` / `price` / `duration_days` / `duration_text` / `sort` / `status` / `recommended` / `features[]` / `model_limits` JSON / `default_model`
- **Subscription management**: `/subscription` (`web/default-pro/src/views/subscription/Subscription.vue`). Admin view shows the user column and the "Add Subscription" button; the user view is read-only
- **Add-subscription modal**: user dropdown from `GET /api/user/search`; plan dropdown from `GET /api/plan/`; `pay_method` is a fixed enum: `free` / `offline` / `wechat` / `alipay` / `bank`

## Implementation Pointers

| Concern | Location |
|---|---|
| Plan CRUD | `controller/plan.go` |
| Plan model | `model/plan.go::Plan` |
| Subscription CRUD | `controller/subscription.go` |
| Order + activation | `model/order_payment.go::CreatePlanOrder` / `ActivatePackageByOrder` |
| Cache invalidation | `model/plan.go::CacheDeleteUserActivePlans` |
| Scheduled expiry | `model/plan.go::ExpireUserPlans` |
