---
title: Top-up Admin
description: "Manual top-up for users and admin-side queries of top-up orders."
category: admin
order: 8
---

# Top-up Admin

> Bypass payment channels and add `quota` directly to a user (`model.IncreaseUserQuota`) with a `LogTypeTopup` audit row.

Route: `/admin/orders` admin order center (the "Mark as paid" flow creates an `OrderTypeTopup` order and activates it immediately). The legacy endpoint `POST /api/user/topup` (`controller/user.go::AdminTopUp`) is kept for compatibility.

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/user/topup` | `POST` | Admin | Directly add quota (legacy path) |
| `/api/topup/order` | `POST` | User | User self-service top-up (requires `topup.enabled=true`) |
| `/api/order/:id` | `PUT` | Admin | Mark paid (covers admin pay methods, see Orders page) |
| `/api/order/` | `GET` | Admin | Query top-up orders with `?type=2` |

User-facing flows are documented in the user docs; this page covers the admin view only.

## POST `/api/user/topup` (legacy)

```json
{ "user_id": 42, "quota": 100000, "remark": "support compensation" }
```

- Calls `model.IncreaseUserQuota(user_id, quota)` directly (Redis cache invalidated).
- When `remark` is empty, defaults to `通过 API 充值 <LogQuota(quota)>`.
- Writes a `LogTypeTopup` audit row via `model.RecordTopupLog`.

> No `Order` row is created; switch to the order-center path if you need an auditable order record.

## Activate Top-up via Order Center (recommended)

1. In `/admin/orders`, click **New order** (or use the user-side top-up form with `POST /api/order`) to create an `OrderTypeTopup=2` order.
2. After the user pays offline (`pay_method=bank` / `offline`), click **Mark paid** in the order list.
3. `PUT /api/order/:id { status:1, pay_method, pay_trade_no }` triggers `model.ActivateTopupByOrder`:
   - Parses `plan_info` to recover the snapshotted `bonus_quota`.
   - Calls `IncreaseUserQuota` to add to the balance.
   - Persists `status`, `pay_status`, `pay_time`, `pay_trade_no`.

`ActivateTopupByOrder` is idempotent — orders already at `status=1` are returned with `nil` immediately.

## Top-up Order Query

`GET /api/order/?type=2` (or pick `Type=Topup` in the admin filter dropdown) to see only top-up orders.
Columns are shared with plan orders: `order_no` / `user_id` / `plan_id` (always `0` for top-ups) / `amount` / `pay_method` / `status` / `source`.

`plan_info` is JSON with this shape:

```json
{ "amount": 100.00, "preset_amount": 100.00, "bonus_quota": 100000, "exchange_rate": 1000 }
```

`plan_info.bonus_quota` is the authoritative granted amount — prefer it over `amount` to be safe against later exchange-rate changes.

## Frontend Guide

- Manual quota grant (legacy): some admin user rows have a "+ quota" action that submits `{ user_id, quota, remark }`.
- Mark paid: pending orders in the order center get a "Mark paid" button at the end of the row.
- Refund: paid orders get a "Refund" button (status flip only, no quota claw-back).
- Delete: Root can hard-delete non-paid orders.

## Implementation Pointers

| Concern | Location |
|---|---|
| Legacy manual grant | `controller/user.go::AdminTopUp` |
| User self-service order | `controller/topup.go::CreateTopupOrder` |
| Order-center mark paid | `controller/order.go::MarkOrderPaid` |
| Idempotent activation | `model/topup.go::ActivateTopupByOrder` |
| Order creation | `model/topup.go::CreateTopupOrder` |
| Routes | `router/api.go` |