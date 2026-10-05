---
title: Plan
description: "Subscription plan template: quota, validity, billing."
category: schema
order: 7
---

# Plan

## What it is

`Plan` is the "product" template a user can subscribe to. After checkout, the system creates a [Subscription](/en/schema/subscription) instance from one Plan.

## Where to find it

- **Admin → Plans**: list, add, publish / unpublish
- **Public → Subscribe**: only published plans are visible
- **Admin → Plan Settings**: global upgrade / expiry policy

## Operator-relevant fields

| Field | Meaning | Effect |
|---|---|---|
| Name | Shown to the user | Updating it updates everywhere. |
| Price | ¥ | The amount the user is charged. |
| Billing type | `token` or `request` | Which dimension the per-model limits are counted on. |
| Plan quota | CNY | Shared virtual balance for the models listed in the plan; when exhausted the whole plan stops working. `0` = unlimited. |
| Model limits | JSON | Declares which models the plan covers and their rolling-window (5h / week / month) limits. Malformed JSON is rejected on save. |
| Validity | Days | After expiry → "plan expired" status, but calls still work until quota runs out. |
| Description | Rich text | Shown on the subscribe page. |
| Status | Published / unpublished | Unpublished stops new subscriptions; existing ones are unaffected. |
| Recommended | Toggle | Shows the "Recommended" badge on the public card. |

## Plan vs Subscription

- **Plan**: the template — name, price, quota, validity.
- **Subscription**: an instance — a user subscribing to a Plan at a specific moment. One user can hold multiple over time.

## Upgrades and differentials

When switching from Plan A to a higher-priced Plan B:

- If `B.price > A.price`, a differential order is created (prefix `UP`)
- Remaining days / quota are pro-rated into the new plan

Upgrade rules are maintained in [Plan Settings](/en/subscription/plan-settings).

## Plan quota and invalidation rules

The models listed in `model_limits` share the `virtual_amount` pool:

- Every subscription call adds `model price × group discount` to `user_plans.used_amount`
- `used_amount >= virtual_amount` → the whole plan becomes unusable and requests fall back to the account balance
- Any rolling window (5h / week / month) reaching 100% weighted usage → same fallback
- Models not listed in `model_limits` are never covered by the plan
- Subscription usage counts into `users.used_quota` / `request_count` but never deducts `users.quota`

> `virtual_amount = 0` means unlimited (legacy-compatible); only the window limits apply then.

## Related pages

- [Plan Pricing](/en/pricing/plan)
- [Plan Management (admin)](/en/subscription/plan-management)
- [Plan Settings](/en/subscription/plan-settings)
- [Upgrade & Downgrade](/en/subscription/upgrade-downgrade)

## Related API

- `GET /api/plan/` — list (admin)
- `GET /api/plan/list` — list (public)
- `POST /api/plan/` — add
- `PUT /api/plan/` — update
- `DELETE /api/plan/:id` — delete