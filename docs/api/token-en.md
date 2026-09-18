---
title: Token API
description: "All /api/token/* endpoints."
category: api
order: 14
---
## 6. Token Management (Token)

### 6.1 List All Tokens

**Endpoint:** `GET /api/token/`

**Auth:** User (returns only the caller's own tokens)

**Query parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| p | int | Page number, default 0 |
| order | string | Sort field |

**Response:**

```json
{
  "success": true,
  "message": "",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "key": "sk-xxxxxxxx",
      "status": 1,
      "name": "my-token",
      "created_time": 1718000000,
      "accessed_time": 1718000000,
      "expired_time": -1,
      "remain_quota": 500000,
      "unlimited_quota": false,
      "used_quota": 100000,
      "models": null,
      "subnet": null,
      "updated_at": 1718000000
    }
  ]
}
```

**Response fields:**

| Field | Type | Description |
|-------|------|-------------|
| id | int | Token ID |
| user_id | int | Owner user ID |
| key | string | Token key |
| status | int | Status: 1=enabled, 2=disabled, 3=expired, 4=exhausted |
| name | string | Token name |
| created_time | int64 | Created time |
| accessed_time | int64 | Last access time |
| expired_time | int64 | Expiration time, -1=never |
| remain_quota | int64 | Remaining quota |
| unlimited_quota | bool | Whether quota is unlimited |
| used_quota | int64 | Quota consumed |
| models | string/null | Allowed models (comma separated), null=all |
| subnet | string/null | Allowed subnet |


### 6.2 Search Tokens

**Endpoint:** `GET /api/token/search`

**Auth:** User

**Query parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| keyword | string | Search keyword |


### 6.3 Get a Single Token

**Endpoint:** `GET /api/token/:id`

**Auth:** User


### 6.4 Create a Token

**Endpoint:** `POST /api/token/`

**Auth:** User

**Request body:**

```json
{
  "name": "my-token",
  "remain_quota": 500000,
  "expired_time": -1,
  "unlimited_quota": false,
  "models": null,
  "subnet": null
}
```


### 6.5 Update a Token

**Endpoint:** `PUT /api/token/`

**Auth:** User

Same shape as create, plus the `id` field. When the query parameter `status_only=1` is set, only the status field is updated.


### 6.6 Delete a Token

**Endpoint:** `DELETE /api/token/:id`

**Auth:** User
