---
title: Admin Dashboard API
description: "/api/admin/dashboard/* endpoints."
category: api
order: 15
---
# Admin Dashboard API

`/api/admin/dashboard/*` powers the homepage of the admin console, exposing site-wide aggregates. All endpoints require **Admin**; on failure they return `{success, message}`.

## Endpoint index

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/api/admin/dashboard/overview` | GET | Admin | KPI overview (users / resources / quota / revenue) |
| `/api/admin/dashboard/charts` | GET | Admin | Site-wide day × model aggregates (`[]LogStatistic`) |
| `/api/admin/dashboard/top-users` | GET | Admin | Active-user leaderboard |


## 1. KPI Overview

**Endpoint:** `GET /api/admin/dashboard/overview`

**Auth:** Admin

**Query parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| range | string | Time window. One of `today` / `7d` / `30d` / `all`. Default `7d`. |

**Response:**

```json
{
  "success": true,
  "message": "",
  "data": {
    "users": {
      "total": 1024,
      "enabled": 900,
      "disabled": 80,
      "deleted": 44,
      "new_today": 5,
      "new_7d": 30,
      "new_30d": 120,
      "active_7d": 210
    },
    "tokens": { "total": 3200, "enabled": 2800 },
    "channels": { "total": 18, "enabled": 15 },
    "plans": { "total": 5, "enabled": 4 },
    "redemptions": { "total": 500, "used": 200, "unused": 300 },
    "subscriptions": { "total": 220, "active": 180, "expired": 40 },
    "quota": { "today": 120000, "week": 900000, "month": 3600000, "total": 12000000 },
    "revenue": {
      "total": 19999.50,
      "topup": 8000.00,
      "subscription": 11999.50,
      "refund": 0.00
    },
    "range": "7d",
    "generated_at": 1718000000
  }
}
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| users.total | int64 | All users (including disabled/deleted) |
| users.enabled | int64 | Users with status=enabled |
| users.disabled | int64 | Manually disabled users |
| users.deleted | int64 | Soft-deleted users |
| users.new_today / new_7d / new_30d | int64 | New users in the today / 7d / 30d window |
| users.active_7d | int64 | Users who had a consume log in the past 7 days AND `request_count > 0` |
| tokens.total / enabled | int64 | Token count / enabled tokens |
| channels.total / enabled | int64 | Channel count / enabled channels |
| plans.total / enabled | int64 | Plan count / listed plans |
| redemptions.total / used / unused | int64 | Redemption code count / used / unused |
| subscriptions.total / active / expired | int64 | User-plan count / active / expired |
| quota.today / week / month / total | int64 | Sum of quota from consume logs in each window |
| revenue.total / topup / subscription | float64 | Total revenue / topup / subscription (CNY) |
| revenue.refund | float64 | Refund total (CNY) |
| range | string | Actual window used |
| generated_at | int64 | Server-side generation timestamp |

> Revenue is filtered by `orders.pay_time` in the `range` window with `status=1` (paid) and `status=3` (refunded).


## 2. Site-wide Chart Data

**Endpoint:** `GET /api/admin/dashboard/charts`

**Auth:** Admin

**Query parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| range | string | Time window. One of `today` / `7d` / `30d` / `all`. Default `7d`. `all` is capped to the last 30 days. |

**Response:**

```json
{
  "success": true,
  "message": "",
  "data": [
    {
      "day": "2026-09-10",
      "model_name": "gpt-4o",
      "request_count": 120,
      "quota": 250000,
      "prompt_tokens": 80000,
      "completion_tokens": 50000
    },
    {
      "day": "2026-09-10",
      "model_name": "claude-3-5-sonnet",
      "request_count": 35,
      "quota": 70000,
      "prompt_tokens": 22000,
      "completion_tokens": 14000
    }
  ]
}
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| day | string | Day bucket (YYYY-MM-DD, UTC day boundary) |
| model_name | string | Model name |
| request_count | int | Requests for that model on that day |
| quota | int | Total quota consumed for that model on that day |
| prompt_tokens | int | Sum of input tokens |
| completion_tokens | int | Sum of output tokens |

> The response shape matches `model.SearchLogsByDayAndModel` (i.e. `GET /api/user/dashboard`), so the frontend can reuse the same chart-building logic.


## 3. Active-User Leaderboard

**Endpoint:** `GET /api/admin/dashboard/top-users`

**Auth:** Admin

**Query parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| range | string | Time window. One of `today` / `7d` / `30d` / `all`. Default `7d`. |
| limit | int | Rows to return, 1-200, default 20 |

**Response:**

```json
{
  "success": true,
  "message": "",
  "data": {
    "range": "7d",
    "items": [
      {
        "id": 12,
        "username": "alice",
        "email": "alice@example.com",
        "request_count": 1230,
        "quota": 800000,
        "balance": 5000000,
        "current_plan_name": "Pro Plan"
      }
    ]
  }
}
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| id | int | User ID |
| username | string | Username |
| email | string | Email |
| request_count | int64 | Requests in the window |
| quota | int64 | Quota consumed in the window |
| balance | int64 | Current remaining quota on the account |
| current_plan_name | string | Name of the earliest-expiring active subscription, if any |

> Sort order: `request_count` DESC, then window quota DESC. Only users with `status != deleted` and `request_count > 0` are included.