---
title: Log API
description: "All /api/log/* endpoints."
category: api
order: 16
---
## 8. Logs (Log)

### 8.1 List All Logs

**Endpoint:** `GET /api/log/`

**Auth:** Admin

**Query parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| p | int | Page number, default 0 |
| type | int | Log type: 0=all, 1=topup, 2=consume, 3=admin, 4=system |
| start_timestamp | int64 | Start timestamp |
| end_timestamp | int64 | End timestamp |
| model_name | string | Filter by model name |
| username | string | Filter by username |
| token_name | string | Filter by token name |
| channel | int | Filter by channel ID |
| group | string | Filter by group |

**Response fields:**

| Field | Type | Description |
|-------|------|-------------|
| id | int | Log ID |
| user_id | int | User ID |
| created_at | int64 | Created at |
| type | int | Log type |
| content | string | Log content |
| username | string | Username |
| token_name | string | Token name |
| model_name | string | Model name |
| quota | int | Quota consumed |
| prompt_tokens | int | Input tokens |
| completion_tokens | int | Output tokens |
| cached_tokens | int | Cached tokens |
| channel_id | int | Channel ID |
| request_id | string | Request ID |
| elapsed_time | int64 | Elapsed time (ms) |
| is_stream | bool | Whether the response is streaming |
| billing_source | int | Billing source: 0=normal, 1=subscription |
| plan_id | int | Plan ID |
| session_key | string | Session key |


### 8.2 Get Self Logs

**Endpoint:** `GET /api/log/self`

**Auth:** User


### 8.3 Search Logs

**Endpoint:** `GET /api/log/search`

**Auth:** Admin

**Endpoint:** `GET /api/log/self/search`

**Auth:** User

**Query parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| keyword | string | Search keyword |
| type | int | Log type (optional) |


### 8.4 Log Statistics

**Endpoint:** `GET /api/log/stat`

**Auth:** Admin

**Endpoint:** `GET /api/log/self/stat`

**Auth:** User


### 8.5 Delete Historical Logs

**Endpoint:** `DELETE /api/log/`

**Auth:** Admin

**Query parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| target_timestamp | int64 | Delete logs older than this timestamp |

**Response:**

```json
{
  "success": true,
  "message": "",
  "data": 100
}
```

> `data` is the number of deleted log rows.
