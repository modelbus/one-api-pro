---
title: Plan Pricing
description: What a plan is, how to set it up, and how users buy it.
category: pricing
order: 3
---

# Plan Pricing

> What defines the "subscription plan" users see on the public subscribe page.

## What it is

A `Plan` is a "product" template. It defines:

- Name / price / validity
- Which models it covers and their quotas
- Whether it's the recommended option

Users see Plans on the public subscribe page.

## Where to find it

- **Admin → Plans**: add, publish, unpublish
- **Public → Subscribe**: only published Plans

## What to set when creating

| Field | Meaning | Effect |
|---|---|---|
| Name | Display name | Updates everywhere |
| Price | ¥ | What the user pays |
| Validity | Days | After expiry, status flips to "expired" — calls still work until quota is exhausted |
| Discount multiplier | Number | Discount applied during the subscription (default 1.0 = no discount) |
| Model quotas | JSON | Per-model window quotas (below) |
| Default model | Model name | Requests for non-listed models get rewritten to this; empty → 422 |
| Description | Rich text | Shown publicly |
| Recommended | Toggle | Shows "★ Recommended" badge |
| Status | Published / Unpublished | Unpublished blocks new subscriptions; existing ones are unaffected |

## Model quotas (`model_limits`) JSON

Each plan can set per-model three-window quotas:

```json
{
  "gpt-4o": {
    "period_h": 5,
    "request_period": 100,
    "request_week": 500,
    "request_month": 2000,
    "token_period": 50000,
    "token_week": 250000,
    "token_month": 1000000
  }
}
```

| Field | Meaning |
|---|---|
| `period_h` | Period window in hours (default 5) |
| `request_period` / `request_week` / `request_month` | Call-count cap per window (`0` = unlimited) |
| `token_period` / `token_week` / `token_month` | Token cap per window |

`model_limits = null` means the plan doesn't restrict that model (pay-as-you-go path).

See [Billing Rules](../en/subscription/billing-rules).

## User purchase flow

```
User picks a plan at /pricing
    ↓
POST /api/order/plan (order number prefix TB / UP)
    ↓
Payment → callback → ActivatePackageByOrder
    ↓
Subscription created (end_time = now + duration_days)
```

Admins can also "grant" manually to bypass payment (common for support compensation).

## Upgrade / Downgrade

Switching plans:

- **Price-diff (default)**: pay only the difference; remaining quota is pro-rated
- **Stack**: pay full price for a new subscription; old one keeps running

Admin picks the mode in [Plan Settings](../en/subscription/plan-settings).

## How to publish a new plan

1. Admin → Plans → Add
2. Name, price, validity, discount
3. Fill `model_limits` for every model the plan should cover
4. Pick a `default_model` (suggested: the plan's flagship model)
5. Set Status to Published
6. Save

Tip: test before publishing. Open a test subscription and verify calls, quotas, upgrade, downgrade.

## vs. user groups

A plan **isn't tied to a user group**. A VIP user buying a basic plan is billed by:

```
consumption × ModelPrice × plan.discount × (group.discount / group.default)
```

## FAQ

- **Changed model_limits — does it affect old subscriptions?** No. Old subscriptions use the snapshot at purchase time.
- **Is default_model required?** Required if you want to limit models; leave empty if the plan is unrestricted.
- **Plan deleted — old subscriptions still work?** Yes. Delete only blocks new subscriptions.

## Related

- [Plan Schema](/en/schema/plan)
- [Plan Management (admin)](../en/subscription/plan-management)
- [Upgrade & Downgrade](../en/subscription/upgrade-downgrade)
- [My Orders (user)](../en/user/orders)

## Related API

- `GET /api/plan/` — list (admin)
- `GET /api/plan/public` — list (public)
- `POST /api/plan/` — add (Root)
- `PUT /api/plan/` — update (Root)
- `DELETE /api/plan/:id` — delete (Root)