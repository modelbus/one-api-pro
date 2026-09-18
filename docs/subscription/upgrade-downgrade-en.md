---
title: Upgrade and Downgrade
description: "price_diff vs stack, the upgrade price formula, and same-tier / downgrade rejection."
category: subscription
order: 5
---

# Upgrade and Downgrade

> Two paths from one plan to another. Implementation: `model/order_payment.go::CreatePlanOrder`, `ActivatePackageByOrder`, `CalculateUpgradePrice`.

## Trigger

When the caller already has an active `user_plan` and submits another `POST /api/order/plan`, the server reads `plan.upgrade_mode` (system setting `model.SystemSettingKeyPlanUpgradeMode`, default `price_diff`) and branches:

- `price_diff` — pay the difference (default)
- `stack` — keep both plans

## Validation

`CreatePlanOrder` runs a hard check before any branching:

- Compare `currentPlan.Sort` to `newPlan.Sort`:
  - `newPlan.Sort < currentPlan.Sort` → `不能降级到低级别套餐`
  - `newPlan.Sort == currentPlan.Sort` → `您已经订阅了同级别的套餐`
- Without an active plan, the order is treated as a fresh subscription

`Plan.Sort` is set by the admin; higher numbers mean higher tier.

## `price_diff` mode

Order number prefix `UP`. Amount is computed by `CalculateUpgradePrice`:

```text
remaining_days = ceil((old_plan.end_time - now) / 86400)
old_daily      = old_plan.price / 30
new_daily      = new_plan.price / 30
upgrade_price  = max(0, (new_daily - old_daily) × remaining_days)
```

When `old_plan.end_time <= now` (the old plan has already expired) the function returns `new_plan.price` — i.e. full price for a fresh subscription.

`ActivatePackageByOrder(order, OrderUpgradeModePriceDiff)`:

1. Flip every `status=UserPlanStatusActive` `user_plans` row for the user to `status=Expired`
2. Insert the new `user_plan`: `start_time=now`, `end_time=now + plan.duration_days × 86400`, `billing_type=token`
3. Call `order.MarkOrderPaid` to mark the order paid
4. `CacheDeleteUserActivePlans(userId)`

## `stack` mode

Order number prefix `TB`. Amount = `new_plan.price` (no proration), independent of existing subscriptions.

`ActivatePackageByOrder(order, OrderUpgradeModeStack)`:

1. Does **not** touch existing active `user_plans`
2. Inserts the new `user_plan` (same start_time / end_time math as price_diff)
3. `MarkOrderPaid` + clear cache

Later `CheckPlanQuota` walks the plans by `end_time` ASC and uses the first non-exhausted one; multiple subscriptions naturally rotate by expiry.

## Admin grants

`POST /api/subscription/` (`controller/subscription.go::AddSubscription`) always uses the `OrderUpgradeModeStack` semantics: existing subscriptions are not closed; the new one is added alongside.

## Frontend Guide

- The user-side **Renew** button on `/subscription` directly calls `POST /api/order/plan { plan_id, pay_method }`; the backend decides `price_diff` vs `stack`
- After the order is created, render the QR / redirect from the returned `pay` object (`buildPayInfo`); admin grants skip this step
- The list's "Plan" column shows `plan.name`; ordering is driven by `plan.sort`
- Admins can extend, expire, or delete subscriptions inline from `/subscription`

## Implementation Pointers

| Concern | Location |
|---|---|
| Mode setting | `model/system_setting.go::SystemSettingKeyPlanUpgradeMode` |
| Order + mode decision | `model/order_payment.go::CreatePlanOrder` |
| Upgrade formula | `model/order_payment.go::CalculateUpgradePrice` |
| Activation (both modes) | `model/order_payment.go::ActivatePackageByOrder` |
| Admin grant | `controller/subscription.go::AddSubscription` |
| Pre-pay | `controller/order.go::buildPayInfo` |
