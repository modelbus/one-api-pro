---
title: Redemption API
description: "All /api/redemption/* endpoints."
category: api
order: 17
---
## 9. Redemption

### 9.1 List All Redemption Codes

**Endpoint:** `GET /api/redemption/`

**Auth:** Admin

**Query parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| p | int | Page number, default 0 |


### 9.2 Search Redemption Codes

**Endpoint:** `GET /api/redemption/search`

**Auth:** Admin

**Query parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| keyword | string | Search keyword |


### 9.3 Get a Single Redemption Code

**Endpoint:** `GET /api/redemption/:id`

**Auth:** Admin


### 9.4 Create Redemption Codes (Batch)

**Endpoint:** `POST /api/redemption/`

**Auth:** Admin

**Request body:**

```json
{
  "name": "Redemption name",
  "quota": 500000,
  "count": 10
}
```

**Request fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | yes | Redemption name |
| quota | int64 | yes | Redemption quota |
| count | int | no | Batch size (1-100), default 1 |

**Response:** Returns the list of created redemption codes, including the auto-generated `key`.


### 9.5 Update a Redemption Code

**Endpoint:** `PUT /api/redemption/`

**Auth:** Admin

When the query parameter `status_only=1` is set, only the status field is updated.


### 9.6 Delete a Redemption Code

**Endpoint:** `DELETE /api/redemption/:id`

**Auth:** Admin
