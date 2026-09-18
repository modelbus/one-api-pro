---
title: Order API
description: "/api/order/* endpoints (user self-service + admin)"
category: api
order: 8
---

# Order API

The Order API serves both the user self-service flow and the admin order center. Order types:

| Constant | Value | Meaning |
| --- | --- | --- |
| `OrderTypePlan` | 1 | Plan order (upgrade differential `UP` and new `TB`) |
| `OrderTypeTopup` | 2 | Top-up order (`TP`) |

## Conventions

- Auth: `/api/order/plan` and `/api/order/self/*` use Cookie Session or Access Token; `/api/order`, `/api/order/:id` and `/api/order/search` require admin (`AdminAuth` or root).
- All responses use the unified JSON envelope `{ success, message, data }`.
- Order status: `OrderStatusPending=0`, `OrderStatusPaid=1`, `OrderStatusCancelled=2`, `Refunded=3`. Pay status: `OrderPayStatusPending=0`, `OrderPayStatusPaid=1`, `OrderPayStatusRefunded=-1`.

## User Self-Service

### `POST /api/order/plan`

Create a plan subscription order. On upgrade the server computes the differential (prorated by remaining days) and the order number is prefixed `UP`; a new subscription uses `TB`.

Body (`CreatePlanOrderRequest`):

```json
{ "plan_id": 1, "mode": "price_diff | stack" }
```

- `mode=price_diff` — upgrade; subtracts the residual value of the existing subscription before charging (order number `UP`).
- `mode=stack` — stack a new subscription on top; charges the full plan price (order number `TB`).

Response: `{ success, message, data: { id, order_no, amount, ... } }`.

### `GET /api/order/self?type=1|2`

The current user's order list; `type=1` filters plan orders, `type=2` filters top-ups.

### `GET /api/order/self/:id`

Order detail, owner-checked against the current user.

### `POST /api/order/self/:id/cancel`

User-initiated cancellation of a still-pending (`status=0`) order.

### `POST /api/order/self/:id/pay`

Re-initiate payment for a pending self-order (reuses `buildPayInfo`); useful after a payment interruption.

## Admin

### `GET /api/order`

Admin order list (`AdminAuth`), paginated with multi-dimensional filters (type / status / source / user_id / plan_id / keyword).

### `GET /api/order/search?keyword=...`

Keyword search.

### `GET /api/order/:id`

Order detail (any user).

### `PUT /api/order/:id`

Admin marks an order as **paid** or **refunded**:

```json
{ "status": 1 | 3, "pay_method": "wechat|alipay|bank|offline|free", "pay_trade_no": "optional" }
```

- `status=1` — mark paid; for `offline` / `bank` channels the subscription activates immediately via `ActivatePackageByOrder`.
- `status=3` — mark refunded (status-only flip; quota is not reversed — the refund loop is a future enhancement).

### `DELETE /api/order/:id`

Hard-delete an order; root only.

## Order Number Prefixes

`GenerateOrderNo` is defined in `model/order_payment.go`:

| Prefix | Scenario |
| --- | --- |
| `TB` | New plan order (`CreatePlanOrder` with `mode=stack` or first purchase) |
| `UP` | Plan upgrade differential (`mode=price_diff`) |
| `TP` | Top-up (`CreateTopupOrder`) |

## Related

- Payment callback: `/api/payment/{wechat,alipay}/notify` (see `controller/payment.go::processNotify`) routes by `order.Type` to `ActivatePackageByOrder` or `ActivateTopupByOrder`.
- Top-up creation: `/api/topup/order` — see [Topup API](topup.md).
- Plan management: `/api/plan/*` — see [Plan API](plan.md).