---
title: Manual Top-up (Admin)
description: How admins directly add quota to a user (bypassing payment channels).
category: pricing
order: 8
---

# Manual Top-up (Admin)

> Customer service compensation, event gifts, reconciliation fix — see here.

## Two ways

### Option 1: Through Order Management (recommended, has audit trail)

1. Admin → Orders → New Order (pick "Top-up")
2. Offline bank transfer or admin manual confirmation
3. Row → "Mark paid" → system creates top-up order + credits quota + writes audit log

**Pro**: every top-up has an `Order` record for reconciliation.

### Option 2: Direct quota credit (legacy)

Only for special cases: internal testing, urgent compensation, reconciliation fix.

```
POST /api/user/topup
{
  "user_id": 42,
  "quota": 100000,
  "remark": "Customer service compensation"
}
```

`remark` is the audit note (required). Empty → auto-filled with `API topup <LogQuota(quota)>`.

**Con**: no `Order` row; hard to audit afterwards.

## Prefer option 1

Default to Order Management: every top-up is traceable to a source (bank flow / support ticket).

## How to query top-up orders

Admin → Orders → Type dropdown = "Top-up":

- `order_no`: `TP` prefix
- `user_id`: target user
- `amount`: ¥ paid
- `plan_info` is JSON containing `bonus_quota` (actual credited)
- `pay_method`: wechat / alipay / bank / offline / free
- `status`: unpaid / paid / canceled / refunded
- `pay_time`: payment completion time

## How to mark paid

See [Order Management](./order-management).

## FAQ

- **User says they paid but balance didn't arrive**: check order detail's `pay_time` / `pay_trade_no`; mark paid manually
- **Added the same quota twice**: system is idempotent — already-paid orders won't be re-activated
- **Refunded but balance still there**: correct — system doesn't auto-reverse. Manually adjust in [User Management](../en/user/user-management)

## Related

- [Top-up (concept)](./topup)
- [Top-up Settings](./topup-settings)
- [Order Management](./order-management)
- [User Management (admin)](/en/user/user-management)

## Related API

- `POST /api/user/topup` — direct quota credit (legacy, retained)
- `POST /api/order` — create top-up order (recommended)
- `PUT /api/order/:id` — mark paid (Admin)
- `GET /api/order/?type=2` — list top-up orders