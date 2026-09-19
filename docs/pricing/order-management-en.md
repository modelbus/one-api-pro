---
title: Order Management (Admin)
description: How admins query, mark paid/refund, and delete orders.
category: pricing
order: 7
---

# Order Management (Admin)

> Admin → Orders. Site-wide orders from the admin's perspective.

## Where

Admin → Orders.

## List shows

Per row:

- Order number
- Type chip (Plan / Top-up)
- User summary
- Amount / Plan name
- Payment method
- Status chip (Unpaid / Paid / Canceled / Refunded)
- Source (User / Admin)
- Created at

Top dropdowns: filter by type / status / source. Search box: prefix match on order number / pay-trade-no.

## Filters

| Dropdown | Options |
|---|---|
| Order type | All / Plan / Top-up |
| Status | All / Unpaid / Paid / Canceled / Refunded |
| Source | All / User / Admin |
| Search | Order number / pay-trade-no prefix |

## Mark paid

Use when: callback failed, bank transfer reconciled, testing.

1. Find the order → "Mark paid"
2. Pick `pay_method` (default `offline`)
3. Fill `pay_trade_no` (optional, external flow number)
4. Submit → order becomes Paid + activates plan / credits quota

> Admin "Mark paid" always uses **stack mode**; it won't overwrite an existing subscription.

## Mark refunded

Use when: customer requested refund, dispute resolved.

1. Only `status=1` (Paid) orders can be refunded
2. Row → "Refund" → double-confirm
3. Submit → order becomes Refunded

> **Refund does NOT auto-reverse already-credited quota or cancel subscriptions.** You'll need to:
> - Cancel the user's [subscription](../en/subscription/subscription-management)
> - Adjust the user's [balance](../en/schema/user)

## Delete (Root only)

Use when: test-order cleanup, accidental order creation.

1. Only `status !== 1` (not Paid) orders can be deleted
2. Root row → "Delete" → double-confirm
3. The order row is physically removed

> Hard delete is **irreversible**. Prefer "Refund" or canceling the subscription.

## FAQ

- **User says payment succeeded but the order stays "Unpaid"**: the async callback failed. Check `pay_time` / `pay_trade_no` in the order detail; mark paid manually.
- **Refund doesn't remove the user's balance**: correct — the system doesn't auto-reverse. Adjust manually in [User Management](../en/user/user-management).
- **Where to find the bank's flow number for reconciliation**: the `pay_trade_no` field in the order detail page.

## Related

- [My Orders (public)](../en/user/orders)
- [Order Schema](/en/schema/order)
- [Payment Channels](./payment)

## Related API

- `GET /api/order/` — list (Admin)
- `GET /api/order/search?keyword=` — keyword search (Admin)
- `GET /api/order/:id` — detail (Admin)
- `PUT /api/order/:id` — mark paid / refund (Admin)
- `DELETE /api/order/:id` — delete (Root)