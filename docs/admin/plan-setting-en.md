---
title: Plan Business Rules
description: "Plan business rules: upgrade mode (stack vs price-difference), expiry policy."
category: admin
order: 13
---

# Plan Business Rules

> Site-wide plan policies — currently only `upgrade_mode` (additive `stack` vs price-difference `price_diff`). UI: `web/default-pro/src/views/setting/OperationSetting.vue`, **Plan** section at the bottom.

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/setting/plan` | `GET` | Root | Read the plan-operations bundle; defaults to `price_diff` |
| `/api/setting/plan` | `PUT` | Root | Save the bundle; only `price_diff` / `stack` are accepted |

Implementation: `controller/setting_payment.go::GetPlanSettings` / `PutPlanSettings`.

## Upgrade Mode

Setting key: `plan.upgrade_mode`. Allowed values:

| Value | Constant | Meaning |
|---|---|---|
| `price_diff` | `OrderUpgradeModePriceDiff` | **Price difference**: when a user upgrades to a more expensive plan, they pay only the difference between the two plans; the old subscription is not stacked |
| `stack` | `OrderUpgradeModeStack` | **Stack**: extend `end_time` by the new plan's `duration_days`; useful for "top-up your VIP time" |

Default is `price_diff`. PUT rejects any other value with `upgrade_mode 必须是 price_diff 或 stack`.

## How It Combines with Order Activation

- Self-service plan purchase (`POST /api/order/plan`): uses this setting to decide whether the upgrade charges the price difference (in `price_diff` mode, the amount is clamped to `0` when the new plan is not strictly more expensive than the active one).
- Admin grant (`POST /api/subscription`): **always** uses `stack` (`OrderUpgradeModeStack`), independent of this setting.

## Persistence

Stored as a single row in `system_settings`, key `plan.upgrade_mode`, category `plan`.
The GET endpoint wraps the string into `{ data: { upgrade_mode: "price_diff" } }`.

> The legacy `plan.allow_topup` row was migrated to `topup.enabled`; the row is preserved but the UI no longer reads/writes it (see the comment in `model/system_setting.go`).

## Frontend Guide

- In `/setting/operation`, the bottom-most **Plan** section exposes two radio buttons: **Price difference** (`price_diff`) and **Stack** (`stack`).
- Click **Save** (a button specific to this block) after changing the value.
- The default UI selection is "Price difference".

## Implementation Pointers

| Concern | Location |
|---|---|
| Handler | `controller/setting_payment.go` (`GetPlanSettings` / `PutPlanSettings`) |
| Constants | `model/order.go::OrderUpgradeModePriceDiff` / `OrderUpgradeModeStack` |
| Self-service difference logic | `model/order_payment.go::CreatePlanOrder` |
| Admin grant (always stack) | `controller/subscription.go::AddSubscription` |
| Routes | `router/api.go` |