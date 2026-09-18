---
title: Diagnostic API
description: "/api/diag/* diagnostic endpoints."
category: api
order: 16
---
# Diagnostic API

`/api/diag/*` exposes admin-only diagnostic endpoints used to debug production data inconsistencies (typical scenario: `user_plans.plan_id` always shows 0). All endpoints require **Admin**.

## Endpoint index

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/api/diag/subscriptions` | GET | Admin | Compare `user_plans` against `orders` to locate dirty rows |


## 1. Subscription Diagnostics

**Endpoint:** `GET /api/diag/subscriptions`

**Auth:** Admin

**Description:** Scans the first 100 `user_plans` rows and joins each with its `orders` row (when `order_id > 0`) and a brief `users` profile, so you can triage anomalies like `plan_id <= 0`. The endpoint is read-only and never mutates data.

**Response:**

```json
{
  "success": true,
  "message": "",
  "data": {
    "total_user_plans": 100,
    "missing_plan_id": 2,
    "rows": [
      {
        "user_plan": {
          "id": 42,
          "user_id": 7,
          "plan_id": 1,
          "order_id": 88,
          "start_time": 1718000000,
          "end_time": 1720592000,
          "status": 1,
          "billing_type": "token",
          "created_time": 1718000000,
          "updated_time": 1718000000
        },
        "order_no": "TB20250912153000123456",
        "order_plan_id": 1,
        "order_type": 1,
        "user": {
          "id": 7,
          "username": "alice",
          "display_name": "Alice",
          "email": "alice@example.com",
          "role": 1,
          "status": 1
        }
      },
      {
        "user_plan": {
          "id": 43,
          "user_id": 9,
          "plan_id": 0,
          "order_id": 0,
          "start_time": 0,
          "end_time": 0,
          "status": 1,
          "billing_type": "token",
          "created_time": 1718000000,
          "updated_time": 1718000000
        },
        "order_no": "",
        "order_plan_id": 0,
        "order_type": 0,
        "user": {
          "id": 9,
          "username": "bob",
          "display_name": "Bob",
          "email": "bob@example.com",
          "role": 1,
          "status": 1
        }
      }
    ]
  }
}
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| total_user_plans | int | Number of `user_plans` rows returned (capped at 100) |
| missing_plan_id | int | Count of rows with `plan_id <= 0` (suspected dirty data) |
| rows | array | One row per scanned `user_plan`, joined with order + user brief |

**Per-row fields:**

| Field | Type | Description |
|-------|------|-------------|
| user_plan | object | Raw `UserPlan` row |
| order_no | string | Source order number; empty when `order_id <= 0` |
| order_plan_id | int | `plan_id` on the source order |
| order_type | int | Order type (1=plan subscription, 2=topup) |
| user | object | `UserBrief`: id / username / display_name / email / role / status |

> The order row is only loaded when `user_plan.order_id > 0`; otherwise the order fields are zero / empty.

**Typical usage:**

```bash
# 1. List the first 100 user_plans to spot dirty rows
curl -s 'http://localhost:3000/api/diag/subscriptions' -b cookies.txt | jq

# 2. When missing_plan_id > 0, look up the trigger path by user_plan.id:
#    - plan_id = 0   : order.plan_id was missing/overwritten at activation time
#    - order_id = 0  : orphan user_plan (order hard-deleted without cascade cleanup)
```