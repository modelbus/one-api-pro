---
title: Billing Rules
description: "Three windows (period/week/month), the weighted-usage formula, the batchUpdate path, and the deduction order."
category: subscription
order: 4
---

# Billing Rules

> Window quotas and per-request deduction. Implementation: `model/plan.go::CalcWindowIndex`, `model/plan_quota.go::WeightedUsage`, middleware `middleware/plan_quota.go::PlanQuotaCheck`.

## Three windows

| windowType | Length | Anchor |
|---|---|---|
| `period` | `rule.period_h` hours (default 5) | `start_time` |
| `week` | 7 days | `start_time` |
| `month` | 30 days | `start_time` |

`model.CalcWindowIndex(now, startTime, windowType, periodH)`:

```go
elapsed := now - startTime
period   = elapsed / (period_h * 3600)
week     = elapsed / (7 * 86400)
month    = elapsed / (30 * 86400)
```

`CalcNextResetTime = start_time + (windowIndex+1) * windowDuration`, so a 5-hour period window rolls based on `start_time` rather than the calendar day.

## Weighted usage

Multiple models within a plan each have their own limit. `WeightedUsage` collapses all `(consumed, limit)` pairs into a single 0..100 ratio:

```text
weighted = Σ consumed_i × QuotaPoolCapacity / limit_i
        = Σ consumed_i × 100 / limit_i
```

Any window whose `weighted >= 100` is considered exhausted.

`model/plan_quota.go::WeightedUsage` rules:

1. Walk `PlanUsage` rows; match the current `(model, window_type, window_index)`
2. Pick `(limit, consumed)` by `billing_type`:
   - `BillingTypeRequest="request"`: `consumed = requests`, `limit = rule.request_period|week|month`
   - `BillingTypeToken="token"`: `consumed = prompt_tokens + completion_tokens`, `limit = rule.token_period|week|month`
   - Otherwise / default: try request first; if `limit <= 0`, fall back to token
3. `consumed × 100 / limit` is summed; models with `limit <= 0` are skipped
4. Only rows whose `windowIndex` matches the current window contribute

## Finding the limit

`FindLimit(limits, model, defaultModel)` priority:

1. Exact hit `limits[model]` → use it
2. Fall back to `limits[defaultModel]`
3. Otherwise `not found`; that usage row is excluded from the weighted sum

`CheckPlanQuota` also exposes the resolved `default_model` to the middleware. `PlanQuotaCheck` rewrites `RequestModel` to `default_model` so an off-list model is funneled into the plan default.

## Deduction order

A request walks the relay middleware chain as follows:

1. `PlanQuotaCheck` → `CheckPlanQuota` walks `CacheGetUserActivePlans(userId)` (sorted by `end_time` ASC). For each `user_plan`:
   - If `end_time <= now`, flip status to expired and skip
   - Get `plan.model_limits`; if nil, treat the plan as usable
   - Compute period / week / month weighted values; any `>= 100` exhausts the plan and we move on
   - The first usable plan is admitted; `plan_id` is written to `meta.PlanId`, `billing_type` to `meta.BillingType`
2. No plan match → admit anyway; `meta.PlanId=0` falls through to pay-as-you-go
3. `postConsumeQuota` (`relay/handler/helper.go`):
   - `PlanId > 0`: call `IncrementPlanUsage(plan_id, model, window_type, window_index, 1, prompt_tokens, completion_tokens, cached_tokens)` for each of the three windows; refund any pre-deducted quota back to the token
   - `PlanId == 0`: compute `quota` from `priceResult.BillingType` — `PerRequest` uses `per_request_price`; `Token` uses `(input_price × prompt + output_price × completion + cached_price × cached) × group_discount`; `PostConsumeTokenQuota` settles `users.quota` and `tokens.quota`

## Caches

- `user_plans:<id>`: Redis cache of the user's active plans (`UserPlanCacheSeconds=300`). `AddSubscription` / `UpdateSubscription` / `DeleteSubscription` call `CacheDeleteUserActivePlans` after writes
- `model_price:<name>` / `group_price:<group>:<model>`: model and group discounts (`ModelPriceCacheSeconds=300`)
- The `model_price` and `group_price` tables are kept fresh by `SyncModelPriceCache` / `SyncGroupPriceCache`

## batchUpdate path

When `BATCH_UPDATE_ENABLED=true`:

- `model.UpdateChannelUsedQuota` queues channel used-quota increments in memory; a background task flushes them every `BATCH_UPDATE_INTERVAL` seconds (default 5s)
- `model.IncreaseUserQuota` also goes through batch; other writes stay synchronous

This reduces DB write amplification but requires `main.go` to call `model.InitBatchUpdater()` and run as a single instance.

## Implementation Pointers

| Concern | Location |
|---|---|
| Window math | `model/plan.go::CalcWindowIndex` / `GetWindowDurationSeconds` |
| Weighted deduction | `model/plan_quota.go::WeightedUsage` / `CalculateWeightedUsage` |
| Limit lookup | `model/plan_quota.go::FindLimit` |
| Middleware check | `middleware/plan_quota.go::PlanQuotaCheck` |
| plan_usage writer | `model/plan.go::IncrementPlanUsage` |
| Billing settlement | `relay/handler/helper.go::postConsumeQuota` |
| Cache | `model/plan.go::CacheGetUserActivePlans` / `CacheDeleteUserActivePlans` |
