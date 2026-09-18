---
title: Redemption Quota Rules
description: "How `redemptions.quota` is defined, its relation to plan grants, expiry, and units."
category: redemption
order: 4
---

# Redemption Quota Rules

> The credited amount comes from `redemptions.quota` alone; unlike top-up orders there is no "amount × rate" semantic — codes carry only an integer quota.

## Field

`model.Redemption.Quota` (`quota`, `bigint`, default 100):

- Type `int64`; the same unit as `users.quota`
- On successful redeem: `users.quota = users.quota + redemption.quota` (inside the `Redeem` transaction)

The row never stores a money amount or exchange rate — codes are currency-agnostic.

## Conversion

There is no `amount × rate` path for redemption codes. Admins write the integer `quota` at generation time.

If the operator wants to distribute "X RMB worth of quota", the frontend can convert amount → quota using the same `exchange_rate` from `model.GetTopupSettings()`, but **the persisted value is still the integer quota**.

## Plan binding

Codes are orthogonal to plans: redeeming never creates, extends, or touches `user_plans`. If the operator needs plan-bound distribution, the frontend can group codes by `plan.name` for display, but every individual code remains general-purpose.

## Expiry

`redemptions` has no `expire_time` column. Lifecycle is controlled solely through `status`:

- `Enabled=1` → redeemable at any time
- `Disabled=2` → `Redeem` rejects
- `Used=3` → `Redeem` rejects

For "30-day expiry" semantics, batch-update `status` to `2` at the appropriate cutoff (`PUT ?status_only=true {status:2}`) or run a one-shot cron that flips the column.

## Quota vs balance

`users.quota` is the balance consumed by pay-as-you-go traffic. Redeem credits add to this column; subscription window traffic does not draw on it.

When a user holds both an active subscription and redeemed quota:

- Requests inside the plan window take the subscription path (`meta.PlanId > 0`) and do not spend `users.quota`
- After the plan is exhausted, requests take the pay-as-you-go path and consume `users.quota` — that is where redeemed quota gets spent

## Plan-bundled grants

Plans do not auto-grant quota: `Plan` carries no quota field, and `UserPlan` does not consume quota. To grant quota alongside a plan, chain an `OrderTypeTopup=2` order after the admin grant, or call `IncreaseUserQuota` directly outside `ActivatePackageByOrder`.

## Implementation Pointers

| Concern | Location |
|---|---|
| Field definition | `model/redemption.go::Redemption.Quota` |
| Adding quota | `model/redemption.go::Redeem` (`UPDATE users SET quota = quota + ?`) |
| Top-up rate | `model/topup.go::GetTopupSettings` / `SystemSettingKeyTopupExchangeRate` |
| Plan grant | `controller/subscription.go::AddSubscription` |
