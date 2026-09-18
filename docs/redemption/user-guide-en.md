---
title: User Redemption Guide
description: "/redeem page, the endpoint, success and failure messages."
category: redemption
order: 3
---

# User Redemption Guide

> On `/redeem` the user enters the 32-char `key` and the page calls `POST /api/user/topup`. `model.Redeem` adds quota and writes a `LogTypeTopup` audit row inside a transaction.

## Endpoint

| Field | Value |
|---|---|
| Method | `POST` |
| Path | `/api/user/topup` (also used by admin manual grant, see [Admin view](#admin-view)) |
| Auth | User (self-service) / Admin (manual compensation) |
| Body | `{ "key": "..." }` |

Success response:

```json
{ "success": true, "data": 10000 }
```

`data` is the credited quota (equal to `redemption.quota`).

## Failure modes

`model.Redeem` returns `errors.New(...)` from inside the transaction; the controller wraps it as `success:false, message`:

| Trigger | message |
|---|---|
| `key` is empty | `未提供兑换码` |
| `userId` is 0 | `无效的 user id` |
| `key` not found | `无效的兑换码` |
| `status != Enabled` (used / disabled) | `该兑换码已被使用` |
| Any SQL error | `兑换失败，<reason>` |

The transaction guarantees: concurrent requests for the same key serialize on `FOR UPDATE` — only one wins, others see `status=Used` and bail out.

## Admin view

Admins can call `POST /api/user/topup` (`controller/user.go::AdminTopUp`) too, but pass `{ "user_id": 42, "quota": 10000, "remark": "..." }` instead of a `key`. This path bypasses `Redeem` and directly calls `IncreaseUserQuota(user_id, quota)`, then writes a `LogTypeTopup` log via `RecordTopupLog`. It is independent of the redemption semantic but lives behind the same endpoint.

> For "add quota" operations, prefer the order-center admin grant (`OrderTypeTopup=2` → `ActivateTopupByOrder`) — it keeps the audit order row. `AdminTopUp` is the legacy compatibility path.

## Frontend Guide

Page: `/redeem` (component-level integration).

- One `<a-input>` + **Redeem** button
- After submit, use the `success` field to show `Message.success('兑换成功，到账 X')` or `Message.error(message)`
- On success, refresh the top balance display (or re-fetch `users.quota`)

## Implementation Pointers

| Concern | Location |
|---|---|
| User redeem entry | `controller/user.go::TopUp` |
| Redemption transaction | `model/redemption.go::Redeem` |
| Admin manual grant | `controller/user.go::AdminTopUp` |
| Log | `model/log.go::RecordLog` (`LogTypeTopup`) |
