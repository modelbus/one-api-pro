---
title: Group Price API
description: "All /api/group_price/* endpoints."
category: api
order: 12
---
## 2. Group Price Management (Group Price)

Manage discount multipliers per user group. Requires Root.

### 2.1 List All Group Prices

**Endpoint:** `GET /api/group_price/`

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
      "group_name": "default",
      "model_name": "",
      "discount": 1.0,
      "created_at": 1718000000,
      "updated_at": 1718000000
    },
    {
      "id": 2,
      "group_name": "vip",
      "model_name": "gpt-4o",
      "discount": 0.8,
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
| group_name | string | Group name (e.g. default, vip, svip) |
| model_name | string | Model name; empty string means "default discount for all models in this group" |
| discount | float | Discount multiplier, 1.0=no discount, 0.8=20% off |
| created_at | int64 | Created at (Unix timestamp) |
| updated_at | int64 | Updated at (Unix timestamp) |


### 2.2 Add a Group Price

**Endpoint:** `POST /api/group_price/`

**Auth:** Root

**Request body:**

```json
{
  "group_name": "vip",
  "model_name": "gpt-4o",
  "discount": 0.8
}
```

**Request fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| group_name | string | yes | Group name |
| model_name | string | no | Model name; leave empty for the group's default discount across all models |
| discount | float | no | Discount multiplier, default 1.0 |

**Response:**

```json
{
  "success": true,
  "message": ""
}
```

**Errors:**

- Empty `group_name`: `{"success": false, "message": "分组名称不能为空"}`


### 2.3 Update a Group Price

**Endpoint:** `PUT /api/group_price/`

**Auth:** Root

**Request body:**

```json
{
  "id": 2,
  "group_name": "vip",
  "model_name": "gpt-4o",
  "discount": 0.7
}
```

**Request fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| id | uint | yes | Record ID |
| group_name | string | no | Group name |
| model_name | string | no | Model name |
| discount | float | no | Discount multiplier |

**Response:**

```json
{
  "success": true,
  "message": ""
}
```

**Errors:**

- `id` is 0: `{"success": false, "message": "ID不能为空"}`


### 2.4 Delete a Group Price

**Endpoint:** `DELETE /api/group_price/:id`

**Auth:** Root

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| id | int | Group price record ID |

**Response:**

```json
{
  "success": true,
  "message": ""
}
```
