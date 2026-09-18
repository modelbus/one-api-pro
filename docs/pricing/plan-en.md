---
title: Plan Pricing
description: "Plan data model, validity, model coverage, group bindings, and upgrade semantics."
category: pricing
order: 3
---

# Plan Pricing

> A Plan carries the subscription price, validity, model quotas, and sort key for upgrade decisions. Implementation: `model/plan.go::Plan`, `controller/plan.go`, `web/default-pro/src/views/setting/PlanSetting.vue`.

## Data model

`model.Plan` (`plans` table):

| Field | Type | Notes |
|---|---|---|
| `name` | `varchar(100)` | Plan name (required; `Insert` errors if empty) |
| `description` | `text` | Plan description |
| `price` | `decimal(10,2)` | Price in yuan |
| `duration_days` | `int` | Validity in days (default 30); `end_time = start_time + duration_days * 86400` |
| `duration_text` | `varchar(50)` | Display text (e.g. `30 天`, `季度`) |
| `status` | `int` | `PlanStatusEnabled=1` / `PlanStatusDisabled=0`; `/api/plan/public` returns only `Enabled` |
| `recommended` | `bool` | Frontend "recommended" tag |
| `sort` | `int` | Upgrade / downgrade comparison key; higher is higher-tier |
| `features` | `StringSlice` (text JSON) | Feature list |
| `model_limits` | `text` JSON | Per-model window quotas (see below) |
| `default_model` | `varchar(100)` | Requests for models not in `model_limits` are routed here; `ValidateDefaultModel` requires it to exist in `model_limits` |

## `model_limits` shape

`model/model_plan.go::ModelLimitRule`:

```json
{
  "gpt-4o":            { "period_h": 5, "request_period": 100, "request_week": 500, "request_month": 2000, "token_period": 50000,  "token_week": 250000, "token_month": 1000000 },
  "claude-3.5-sonnet": { "period_h": 5, "request_period": 50,  "request_week": 200, "request_month": 800,  "token_period": 30000,  "token_week": 120000, "token_month": 480000 }
}
```

- `period_h` — period window length in hours (default 5)
- `request_*` / `token_*` — request / token caps per window; `0` means unlimited
- When `model_limits` is null, `CheckPlanQuota` treats the plan as usable (no window limit)

See [Billing Rules](../subscription/billing-rules).

## features

`StringSlice` is a custom JSON/text bridge in `plan.go`:

- Writes: `MarshalJSON` outputs a JSON array; `Value()` serializes to a JSON string
- Reads: `Scan` first tries JSON array parsing, then `\n` fallback (legacy plain-text)
- nil marshals to `[]` instead of `null`

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/plan/` | `GET` | Admin | Paginated |
| `/api/plan/search?keyword=` | `GET` | Admin | `name LIKE kw%` |
| `/api/plan/:id` | `GET` | Admin | Detail |
| `/api/plan/` | `POST` | Admin | Create (name required; `ValidateDefaultModel` runs) |
| `/api/plan/` | `PUT` | Admin | Update (same validation) |
| `/api/plan/:id` | `DELETE` | Admin | Delete |
| `/api/plan/public` | `GET` | Public | Only `status=Enabled`; user-side plan list |
| `/api/plan/public/:id` | `GET` | Public | Public detail |

## Group binding

Plans do not bind directly to user groups; group routing is driven by `channel.group`. The "tier" feel of a plan comes from:

- `recommended=true` — flagged as recommended in the list
- `sort` — upgrade/downgrade comparison (same `sort` is rejected; lower `sort` cannot upgrade to higher `sort` beyond)
- `default_model` — if the user requests a model not in `model_limits`, the request is rewritten to it; empty + non-listed model → 422 (`PlanQuotaCheck`)

## Plan ↔ Order ↔ Subscription

- User self-service: `POST /api/order/plan` → `model.CreatePlanOrder` → `buildPayInfo` returns pre-pay params → payment notify → `ActivatePackageByOrder(order, mode)`
- Admin grant: `POST /api/subscription/` → immediately `ActivatePackageByOrder(order, OrderUpgradeModeStack)`
- Mode is read from the system setting `plan.upgrade_mode` (default `price_diff`)

See [Upgrade and Downgrade](../subscription/upgrade-downgrade) and [Subscription Admin](../subscription/admin-guide).

## Frontend Guide

Route: `/setting/pricing` → **Plans** tab (`web/default-pro/src/views/setting/PlanSetting.vue`).

- Columns: name / price / duration_days / duration_text / sort / status / recommended
- Edit modal: name / description / price / duration_days / duration_text / sort / status / recommended / features[] / model_limits JSON / default_model
- Recommended star (★): `recommended=true`
- **Enable / Disable** toggle changes `status`
- Delete is a secondary confirmation; deleting a plan does not roll back existing `user_plans` (they keep the snapshotted limits at activation time)

## Implementation Pointers

| Concern | Location |
|---|---|
| Data model | `model/plan.go::Plan` / `ModelLimitRule` |
| Validation | `model/plan.go::ValidateDefaultModel` |
| CRUD | `controller/plan.go` |
| Public endpoints | `controller/plan.go::GetPublicPlans` / `GetPublicPlanDetail` |
| Order creation | `model/order_payment.go::CreatePlanOrder` |
| Activation | `model/order_payment.go::ActivatePackageByOrder` |
| Upgrade math | `model/order_payment.go::CalculateUpgradePrice` |
