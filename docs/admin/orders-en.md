---
title: Orders Admin
description: "Order center from the admin view: filter, mark paid / refunded and delete."
category: admin
order: 7
---

# Orders Admin

> Site-wide order center on top of `/api/order`: list, filter, mark paid / refunded, delete. UI: `web/default-pro/src/views/admin/AdminOrders.vue`.

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/order/` | `GET` | Admin | Paginated list (`config.ItemsPerPage`) |
| `/api/order/search?keyword=` | `GET` | Admin | Matches `order_no` / `pay_trade_no` prefix |
| `/api/order/:id` | `GET` | Admin | Order detail (with embedded `User` brief) |
| `/api/order/:id` | `PUT` | Admin | **Mark paid / refunded** (see below) |
| `/api/order/:id` | `DELETE` | **Root** | Hard-delete an order |

Implementation: `controller/order.go`.

## Common Filter Query

`OrderAdminFilter` (`model/order.go`):

| Param | Type | Meaning of `0` / `""` | Notes |
|---|---|---|---|
| `type` | `int` | no filter | `1=plan` / `2=topup` |
| `status` | `int` | no filter | `0=pending` / `1=paid` / `2=canceled` / `3=refunded` |
| `source` | `int` | no filter | `1=user self-service` / `2=admin` |
| `user_id` | `int` | no filter | exact match |
| `plan_id` | `int` | no filter | exact match |
| `keyword` | `string` | no filter | fuzzy match on `order_no` or `pay_trade_no` (`LIKE 'kw%'`) |

`status` uses an empty string for "all" instead of `0` because `0` overlaps with `OrderStatusPending`.

## Mark Paid — `PUT /api/order/:id`

Body:

```json
{
  "status": 1,
  "pay_method": "offline",      // optional override
  "pay_trade_no": "TX-2026..."  // optional, admin-entered reference
}
```

Server behaviour branches on `req.Status`:

- `status=1`: persist `pay_method` / `pay_trade_no`, then call `model.ActivatePackageByOrder(o, OrderUpgradeModeStack)`.
  - Admin activation **always uses `stack`** (additive), never price-diff.
  - `pay_method` of `offline` / `bank` / `wechat` / `alipay` / `free` all activate immediately.
  - For topup orders (`OrderTypeTopup=2`), `model.ActivateTopupByOrder` is called (idempotent).
- `status=3`: call `model.MarkOrderRefunded(o)` — only flips the status bit. **Does NOT revoke quota that was already granted** (planned TODO).
- Other values: rejected.

Success messages: "订单已支付，套餐已激活" or "订单已标记为退款".

## Mark Refunded — `PUT /api/order/:id` (status=3)

Status flip only. Does not roll back the active subscription nor claw back quota. The order row stays as `status=3` for audit. The UI shows the **Refund** button only when `o.status === 1`.

## Delete — `DELETE /api/order/:id`

Visible only to Root; `o.Delete()` performs a hard row delete. There is no cascade to `user_plans`, but `user_plans.order_id` is a soft reference.

## Frontend Guide

- Welcome bar + search field, with three filter dropdowns on the left: order type (plan / topup), status, source.
- Top-right actions: reset filters, refresh.
- List row columns: ID / `order_no` (truncated with tooltip) / type chip / user summary (display_name + `#user_id`) / plan-or-amount (plan name parsed from `plan_info` for `type=1`) / amount / pay method / status chip / source / created-at / actions.
- Row actions:
  - **View**: opens a 560 px detail modal with all fields (`pay_time`, `pay_trade_no`, embedded `User`).
  - **Mark paid** (visible when `status=0`): modal asks for `pay_method` and optional `pay_trade_no`; submitting reloads the list.
  - **Refund** (visible when `status=1`): popconfirm → `PUT { status: 3 }`.
  - **Delete** (Root only, visible when `status !== 1`): popconfirm → `DELETE`.
- Pagination options `[10, 20, 50]`; the front-end pads `items.length + pageSize` to keep "next page" reachable.

## Implementation Pointers

| Concern | Location |
|---|---|
| Order CRUD | `controller/order.go` |
| Admin mark paid / refund | `controller/order.go::MarkOrderPaid` |
| Filter + user brief | `model/order.go::OrderAdminFilter` / `enrichOrdersWithUserBrief` |
| Activate plan | `model/order_payment.go::ActivatePackageByOrder` |
| Activate topup | `model/topup.go::ActivateTopupByOrder` |
| Mark refunded (status only) | `model/order_payment.go::MarkOrderRefunded` |
| Routes | `router/api.go` |