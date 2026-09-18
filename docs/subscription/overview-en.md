---
title: Subscription Overview
description: "How subscriptions, pay-as-you-go, and top-up quota interact, and the deduction order across the three."
category: subscription
order: 1
---

# Subscription Overview

> One API Pro supports three billing paths side by side: subscriptions (plans with window-based quotas), pay-as-you-go (drawn from `users.quota`), and top-up (also credited to `users.quota`). Middleware picks the right path per request.

## Three billing modes

| Mode | Source | Persists to | Deduction point |
|---|---|---|---|
| Subscription | `model.Plan` + `model.UserPlan` | `user_plans` row + `plan_usages` rows | `middleware/plan_quota.go::PlanQuotaCheck` admits based on window quotas; `relay/handler/helper.go::postConsumeQuota` writes `plan_usages` after the call |
| Pay-as-you-go | User calls `/v1/chat/completions` with no usable subscription | `tokens.quota` + `users.quota` | `model.PreConsumeTokenQuota` pre-deducts → upstream returns → `PostConsumeTokenQuota` settles |
| Top-up | Redemption code / top-up order / admin manual grant | `users.quota` direct addition | `model.IncreaseUserQuota` (does not participate in admission; only funds the balance) |

## Deduction order

When a request hits the relay middleware chain:

1. `PlanQuotaCheck` calls `model.CheckPlanQuota(userId, model)`, which walks `CacheGetUserActivePlans(userId)` ordered by `end_time` ASC and admits the first plan whose weighted usage is below `QuotaPoolCapacity` (100). On hit, `plan_id` is stashed in `meta.PlanId`. No hit falls through to the no-plan path.
2. During billing (`postConsumeQuota`):
   - When `meta.PlanId > 0`: increments `plan_usages` (requests, prompt_tokens, completion_tokens, cached_tokens) for the chosen plan and refunds any pre-deducted quota back to the token (subscription users do not consume `quota`).
   - Otherwise: computes the `quota` delta based on `priceResult.billing_type`, then `PostConsumeTokenQuota` settles against `users.quota`.

## Subscription vs top-up

- A subscription carries its own window quotas: every plan declares per-model `request_period/week/month` and `token_period/week/month` in `plan.model_limits`. **It does not consume `users.quota`.**
- Top-up credits `users.quota`; that balance is what pay-as-you-go burns through. Admin manual grant (`/api/user/topup`, legacy) and admin grants via the order center (`OrderTypeTopup=2`) end up on the same path.
- A user can hold multiple subscriptions (the `OrderUpgradeModeStack` mode); they are consumed in `end_time` ASC order. Top-up and subscriptions are independent — subscriptions never spend top-up balance.

## Upgrade paths

When the user already has an active subscription and orders again, the system reads `plan.upgrade_mode` and branches:

- `price_diff` (default): when the new plan is higher-tier, charge the difference; order number prefix `UP`; `ActivatePackageByOrder` marks all current `user_plans` expired before inserting the new one. Same-tier or lower-tier orders are rejected.
- `stack`: keep existing plans and create a new one alongside; order number prefix `TB`; `amount = new_plan.price`.

See [Upgrade and Downgrade](./upgrade-downgrade).

## Expiry

`model.ExpireUserPlans` (scheduled task) flips `status=1 AND end_time <= now` rows to `UserPlanStatusExpired=0`. Once expired, the subscription no longer appears in `CheckPlanQuota`, and subsequent requests fall through to pay-as-you-go.

## Implementation Pointers

| Concern | Location |
|---|---|
| Data model | `model/plan.go::Plan` / `UserPlan` / `PlanUsage` |
| Quota admission | `model/plan_quota.go::CheckPlanQuota` |
| Middleware | `middleware/plan_quota.go::PlanQuotaCheck` |
| Billing settlement | `relay/handler/helper.go::postConsumeQuota` |
| Upgrade / stack | `model/order_payment.go::CreatePlanOrder` / `ActivatePackageByOrder` |
| User top-up | `model/topup.go::CreateTopupOrder` / `ActivateTopupByOrder` |
