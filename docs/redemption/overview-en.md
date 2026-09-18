---
title: Redemption Overview
description: "Redemption code types, use cases, and typical scenarios."
category: redemption
order: 1
---

# Redemption Overview

> Redemption codes are One API Pro's way of distributing `quota` by code. Admins batch-generate them; external channels (referrals, partnerships, events) hand them out; users enter them on `/redeem` to top up `users.quota`.

## Types

`model.Redemption` (`redemptions` table) has `key` as a 32-char UUID with a unique index.

| Status | Constant | Meaning |
|---|---|---|
| 1 | `RedemptionCodeStatusEnabled` | Active, can be redeemed |
| 2 | `RedemptionCodeStatusDisabled` | Disabled by an admin |
| 3 | `RedemptionCodeStatusUsed` | Used (flipped inside the redemption transaction) |

> Status values deliberately skip `0` — the default value must not collide with "Used".

`count` (`gorm:"-:all"`, not persisted) is only used at batch generation time: the server loops `count` times and inserts one row per UUID.

## Use cases

| Scenario | Usage |
|---|---|
| Promotion / referral | One code per quota, scoped by `name`; users redeem themselves |
| Compensation | Admin issues a code to a user after a complaint or incident |
| Marketing campaign | Batch generate (`count ≤ 100 / batch`), distribute via channels |

> Redemption codes only add to `users.quota`. They do **not** create or affect `user_plans` and are orthogonal to plans.

## One-time vs multi-use

Each row in `redemptions` corresponds to a single UUID — once redeemed it flips to `Used=3`. There is no "multi-use" semantic; for distribution to multiple users, generate multiple rows (same `name`, different `key`).

## Side effects on redeem

`model.Redeem(ctx, key, userId)` (`model/redemption.go:55`) runs in a transaction:

1. `SELECT … FOR UPDATE` locks the `redemptions.key` row
2. Validates `status == Enabled`; otherwise returns `该兑换码已被使用`
3. `UPDATE users SET quota = quota + redemptions.quota`
4. Sets `RedeemedTime = now`, `status = Used`, and `Save`
5. Writes a `LogTypeTopup` log: `通过兑换码充值 <LogQuota(redemption.Quota)>`

The credit is only visible after commit; any step failing rolls back the whole transaction, preventing concurrent redemption and double credit.

## Relation to top-up orders

Redemption credits land via `users.quota += redemptions.quota` (in `Redeem`); the `OrderTypeTopup=2` activation path (`model.ActivateTopupByOrder` → `IncreaseUserQuota`) ends up at the same column. They are independent; orders carry the amount / exchange-rate info needed for reconciliation, redemption codes carry just `quota`.

## Implementation Pointers

| Concern | Location |
|---|---|
| Data model | `model/redemption.go::Redemption` |
| Redemption transaction | `model/redemption.go::Redeem` |
| CRUD | `controller/redemption.go` |
| User-side entry point | `controller/user.go::TopUp` (`POST /api/user/topup`) |
