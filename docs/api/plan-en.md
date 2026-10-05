---
title: Plan API
description: "All /api/plan/* endpoints."
category: api
order: 18
---
## 10. Plan Management (Plan)

### 10.1 List All Plans

**Endpoint:** `GET /api/plan/`

**Auth:** Admin

**Query parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| p | int | Page number, default 0 |

**Response:**

```json
{
  "success": true,
  "message": "",
  "data": [
    {
      "id": 1,
      "name": "Basic Plan",
      "price": 99.00,
      "tokens": 500000000,
      "model_limits": "{\"gpt-4o\":{\"request_month\":1000,\"token_month\":50000000}}",
      "description": "Basic plan description",
"features": ["1000 API calls/month", "GPT-4o support"],
      "sort": 0,
      "status": 1,
      "duration_days": 30,
      "duration_text": "30 days",
      "recommended": false,
      "created_time": 1718000000
    }
  ]
}
```

**Response fields:**

| Field | Type | Description |
|-------|------|-------------|
| id | uint | Plan ID |
| name | string | Plan name |
| price | float64 | Price |
| tokens | int64 | Token quota |
| model_limits | string | Model limit config as JSON; key is the model name, value is a `ModelLimitRule` |
| description | string | Description |
| features | array&lt;string&gt; | Feature list, one item per line on the user-facing plan card |
| sort | int | Sort weight |
| status | int | Status: 1=listed, 0=unlisted |
| duration_days | int | Validity in days |
| duration_text | string | Validity display text |
| recommended | bool | Whether the plan is recommended |


### 10.2 Search Plans

**Endpoint:** `GET /api/plan/search`

**Auth:** Admin

**Query parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| keyword | string | Search keyword |


### 10.3 Get Plan Detail

**Endpoint:** `GET /api/plan/:id`

**Auth:** Admin


### 10.4 Create a Plan

**Endpoint:** `POST /api/plan/`

**Auth:** Root

**Request body:**

```json
{
  "name": "Basic Plan",
  "price": 99.00,
  "tokens": 500000000,
  "model_limits": "{\"gpt-4o\":{\"request_month\":1000,\"token_month\":50000000}}",
  "description": "Basic plan description",
  "features": "Feature description",
  "sort": 0,
  "status": 1,
  "duration_days": 30,
  "duration_text": "30 days",
  "recommended": false
}
```


### 10.5 Update a Plan

**Endpoint:** `PUT /api/plan/`

**Auth:** Root

Same shape as create; must include the `id` field.

**`model_limits` field format:**

Keys are model names; values are `ModelLimitRule` objects:

```json
{
  "gpt-4o": {
    "period_h": 5,
    "request_period": 100,
    "request_week": 500,
    "request_month": 2000,
    "token_period": 500000,
    "token_week": 2000000,
    "token_month": 10000000
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| period_h | int | Rolling window length in hours, default 5 |
| request_period | int64 | Max requests per rolling window |
| request_week | int64 | Max requests per week |
| request_month | int64 | Max requests per month |
| token_period | int64 | Max tokens per rolling window |
| token_week | int64 | Max tokens per week |
| token_month | int64 | Max tokens per month |

**Handling of uncovered models:**

- If the requested model is not in `model_limits`, the plan does not cover the request and is skipped.
- If the user also has balance, the request falls back to global balance pay-as-you-go billing; insufficient balance is rejected by the pay-as-you-go path.
- If the platform has no price configured for that model, the response is 422 `model_price_not_found`.
- An empty `model_limits` string means the plan does not restrict models — every model is covered.


### 10.6 Delete a Plan

**Endpoint:** `DELETE /api/plan/:id`

**Auth:** Root
