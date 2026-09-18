---
title: Channel API
description: "All /api/channel/* endpoints."
category: api
order: 13
---
## 5. Channel Management (Channel)

### 5.1 List All Channels

**Endpoint:** `GET /api/channel/`

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
      "type": 1,
      "key": "sk-****xxxx",
      "status": 1,
      "name": "OpenAI",
      "weight": 1,
      "created_time": 1718000000,
      "test_time": 0,
      "response_time": 0,
      "base_url": "https://api.openai.com",
      "other": "",
      "balance": 0,
      "balance_updated_time": 0,
      "models": "gpt-4o,gpt-4o-mini",
      "group": "default",
      "used_quota": 0,
      "model_mapping": "",
      "priority": 0,
      "config": "{}",
      "system_prompt": "",
      "max_concurrency": 0,
      "cooldown_seconds": 60,
      "rpm": 0,
      "last_error": "",
      "last_error_time": 0
    }
  ]
}
```

> **Note:** Non-Root viewers see the `key` field masked.


### 5.2 Search Channels

**Endpoint:** `GET /api/channel/search`

**Auth:** Admin

**Query parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| keyword | string | Search keyword (matches ID or name prefix) |

**Response:** Same shape as 5.1.


### 5.3 Get a Single Channel

**Endpoint:** `GET /api/channel/:id`

**Auth:** Admin

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| id | int | Channel ID |

**Response:** Same shape as a single item in 5.1.


### 5.4 Add a Channel

**Endpoint:** `POST /api/channel/`

**Auth:** Admin

**Request body:**

```json
{
  "type": 1,
  "key": "sk-xxxxxxxx",
  "name": "OpenAI",
  "base_url": "https://api.openai.com",
  "models": "gpt-4o,gpt-4o-mini",
  "group": "default",
  "weight": 1,
  "priority": 0,
  "model_mapping": "",
  "config": "{}",
  "system_prompt": "",
  "max_concurrency": 0,
  "cooldown_seconds": 60,
  "rpm": 0
}
```

**Request fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| type | int | yes | Channel type (1=OpenAI, 3=Azure, etc.) |
| key | string | yes | API key, separate multiple keys with newlines |
| name | string | yes | Channel name |
| base_url | string | no | API base URL |
| models | string | yes | Supported models, comma separated |
| group | string | no | Group, default "default" |
| weight | int | no | Weight, default 0 |
| priority | int | no | Priority, default 0 |
| model_mapping | string | no | Model mapping JSON |
| config | string | no | Channel config JSON |
| system_prompt | string | no | System prompt |
| max_concurrency | int | no | Max concurrency, 0=unlimited |
| cooldown_seconds | int | no | Cooldown in seconds, default 60 |
| rpm | int | no | Requests-per-minute cap |

**Response:**

```json
{
  "success": true,
  "message": ""
}
```


### 5.5 Update a Channel

**Endpoint:** `PUT /api/channel/`

**Auth:** Admin

**Request body:** Same shape as 5.4, must include the `id` field.


### 5.6 Delete a Channel

**Endpoint:** `DELETE /api/channel/:id`

**Auth:** Admin


### 5.7 Delete Disabled Channels

**Endpoint:** `DELETE /api/channel/disabled`

**Auth:** Admin

**Description:** Delete every channel with status set to disabled (manual disable + auto disable).


### 5.8 Test a Channel

**Endpoint:** `GET /api/channel/test/:id`

**Auth:** Admin

**Query parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| model | string | Optional, the model used for the test call |

**Endpoint:** `GET /api/channel/test`

**Auth:** Admin

**Query parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| scope | string | Optional, test scope: `all`, `limited`, `disabled` |


### 5.9 Update Channel Balance

**Endpoint:** `GET /api/channel/update_balance/:id`

**Auth:** Admin

**Description:** Refresh the balance for the specified channel.

**Endpoint:** `GET /api/channel/update_balance`

**Auth:** Admin

**Description:** Refresh the balance for every channel.
