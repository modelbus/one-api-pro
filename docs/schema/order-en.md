---
title: Order
description: Unified record for plan purchases, upgrades, and differentials; the source of truth for finance reconciliation.
category: schema
order: 9
---

# Order

## What it is

`Order` is the **unified record for every paid action** — new plan, plan upgrade, top-up. The order number prefix encodes the type:

- `TB` — new plan
- `UP` — upgrade differential
- `TP` — top-up (see [Topup](/en/schema/topup))

## Where to find it

- **Public → My Orders**: the current user's orders and statuses
- **Admin → Orders**: admin view with filters
- **Payment callbacks**: third-party payment callbacks identify the order by number

## Operator-relevant fields

| Field | Meaning | Effect |
|---|---|---|
| Order number | `TB` / `UP` / `TP` prefix + timestamp + random | Changing it breaks payment callbacks. |
| Parent | User FK | Cascade-delete with user. |
| Type | New / Upgrade / Top-up | Drives the downstream fulfillment. |
| Plan FK | Target plan | For upgrades only. |
| Subscription FK | Activated subscription | Set when fulfillment completes. |
| Amount | ¥ | What was actually paid. |
| Status | Unpaid / Paid / Cancelled / Refunded | Drives fulfillment and refunds. |
| Pay method | WeChat / Alipay / Bank / Mock | Which callback URL is hit. |

## Lifecycle

```
Unpaid → (user pays) → Paid → (system activates subscription / credits) → Done
       ↘ (cancel / timeout) → Cancelled
```

## Activation

On "Paid", the system:

1. Creates or extends a [Subscription](/en/schema/subscription)
2. For upgrade orders, pro-rates the remaining quota from the old subscription

Activation is **asynchronous** — it happens within seconds of the callback.

## Refunds

Triggered from the order detail page. Refunds:

1. Mark the order as "Refunded"
2. Cancel the linked subscription
3. Do not refund quota already consumed

## Related pages

- [My Orders (user)](/en/user/orders)
- [Order Management (admin)](/en/pricing/order-management)

## Related API

- `POST /api/order/plan` — user creates a plan order
- `GET /api/order/self` — current user's orders
- `GET /api/order/` — list (admin)
- `PUT /api/order/:id` — mark paid (admin)
- `DELETE /api/order/:id` — delete (Root)