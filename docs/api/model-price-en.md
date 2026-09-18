---
title: Model Price API
description: "All /api/model_price/* endpoints."
category: api
order: 11
---
## 1. Model Price Management (Model Price)

Manage per-model token pricing (¥ per million tokens) and per-request pricing. Requires Root.

### 1.1 List All Model Prices

**Endpoint:** `GET /api/model_price/`

**Auth:** Root

**Request parameters:** None

**Response:**

```json
{
  "success": true,
  "message": "",
  "data": [
    {
      "id": 1,
      "model_name": "gpt-4o",
      "input_price": 2.5,
      "output_price": 10.0,
      "cached_price": 1.25,
      "per_request_price": 0,
      "billing_type": "token",
      "enabled": true,
      "created_at": 1718000000,
      "updated_at": 1718000000
    }
  ]
}
```

**Response fields:**

| Field | Type | Description |
|-------|------|-------------|
| id | uint | Record ID |
| model_name | string | Model name (unique) |
| input_price | float | Input price (¥ per million tokens) |
| output_price | float | Output price (¥ per million tokens) |
| cached_price | float | Cached price (¥ per million tokens), 0 means unsupported |
| per_request_price | float | Per-request price (¥ per request) |
| billing_type | string | Billing type: `token` or `per_request` |
| enabled | bool | Whether enabled |
| created_at | int64 | Created at (Unix timestamp) |
| updated_at | int64 | Updated at (Unix timestamp) |


### 1.2 Add a Model Price

**Endpoint:** `POST /api/model_price/`

**Auth:** Root

**Request body:**

```json
{
  "model_name": "gpt-4o",
  "input_price": 2.5,
  "output_price": 10.0,
  "cached_price": 1.25,
  "per_request_price": 0,
  "billing_type": "token",
  "enabled": true
}
```

**Request fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| model_name | string | yes | Model name, must be unique |
| input_price | float | no | Input price, default 0 |
| output_price | float | no | Output price, default 0 |
| cached_price | float | no | Cached price, default 0 |
| per_request_price | float | no | Per-request price, default 0 |
| billing_type | string | no | `token` (default) or `per_request` |
| enabled | bool | no | Default true |

**Response:**

```json
{
  "success": true,
  "message": ""
}
```

**Errors:**

- Empty `model_name`: `{"success": false, "message": "模型名称不能为空"}`
- Duplicate `model_name`: DB unique-constraint violation.


### 1.3 Update a Model Price

**Endpoint:** `PUT /api/model_price/`

**Auth:** Root

**Request body:**

```json
{
  "id": 1,
  "model_name": "gpt-4o",
  "input_price": 3.0,
  "output_price": 12.0,
  "cached_price": 1.5,
  "per_request_price": 0,
  "billing_type": "token",
  "enabled": true
}
```

**Request fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| id | uint | yes | Record ID |
| model_name | string | no | Model name (used as the lookup key during update, not changed) |
| input_price | float | no | Input price |
| output_price | float | no | Output price |
| cached_price | float | no | Cached price |
| per_request_price | float | no | Per-request price |
| billing_type | string | no | Billing type |
| enabled | bool | no | Whether enabled |

**Response:**

```json
{
  "success": true,
  "message": ""
}
```

**Errors:**

- `id` is 0: `{"success": false, "message": "ID不能为空"}`


### 1.4 Delete a Model Price

**Endpoint:** `DELETE /api/model_price/:id`

**Auth:** Root

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| id | int | Model price record ID |

**Response:**

```json
{
  "success": true,
  "message": ""
}
```

**Errors:**

- Invalid `id`: `{"success": false, "message": "无效的ID"}`
