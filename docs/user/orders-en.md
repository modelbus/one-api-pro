---
title: My Orders
description: "All your orders: plan purchases, upgrade differentials, top-ups."
category: user
order: 5
---

# My Orders

> Everything you've ever paid for: plans, plan upgrade differentials, top-ups.

## Where

Visit `/orders` after logging in.

## Three order types

| Type | What | When |
|---|---|---|
| **Plan order** | Buy / upgrade a plan | Subscribe page, plan upgrade |
| **Top-up order** | Add credit to your balance | Top-up page |
| **Admin-created** | Admin creates it for you | Admin backend manual / free grant |

Order-number prefix: `TB` (new) / `UP` (upgrade diff) / `TP` (top-up).

## Order states

| Status | Meaning | How it gets here |
|---|---|---|
| Unpaid | Order created but not paid | You just clicked pay |
| Paid | Callback received / admin marked | WeChat / Alipay confirmed / admin manual |
| Canceled | You canceled | Click cancel (only when unpaid) |
| Refunded | Admin refunded | Admin action in backend |

## How to pay

If an order is "Unpaid", the list shows a "Pay" button. Click it — the page opens the previously selected payment method (WeChat / Alipay / bank).

Payment completes and the page returns automatically — no need to refresh.

## Upgrade differentials

When upgrading to a more expensive plan:

- Default: pay only the **difference**. `newPlan.price - oldPlan.remaining_value`.
- Order-number prefix: `UP`.
- Remaining quota is pro-rated into the new plan.

Alternative "stack" mode: pay full price for the new plan; the old one keeps running until its original expiry. Admin can toggle this in [Plan Settings](../subscription/plan-settings).

## Filtering

Three tabs at the top of the page:

- **All**: every order
- **Plan**: only `TB` / `UP`
- **Top-up**: only `TP`

## FAQ

- **Order stuck on "Unpaid"**: callback hasn't arrived. Ask an admin to mark it paid in Admin → Orders, or check the payment channel config
- **Paid but order not activated**: usually activates within 1 minute. If it's been > 5 min, see [Troubleshooting](../misc/troubleshooting#payment-succeeded-but-order-not-activated)
- **Negative upgrade differential**: downgrading from an expensive plan to a cheap one — the old plan's remaining value can exceed the new price. Confirm with the admin.

## Related

- [Subscription (Token Plan)](../subscription/overview)
- [Upgrade & Downgrade](../subscription/upgrade-downgrade)
- [Top-up](../pricing/topup-settings)
- [Payment Settings (admin)](../pricing/payment-settings)