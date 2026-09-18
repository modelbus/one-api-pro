---
title: Subscription API
description: "All /api/subscription/* endpoints."
category: api
order: 19
---
## 11. Subscription Management (Subscription)

### 11.1 Get User Subscription Info

**Endpoint:** `GET /api/subscription/self`

**Auth:** User

**Response:** The caller's active subscription list with usage details.


### 11.2 List All Subscriptions

**Endpoint:** `GET /api/subscription/`

**Auth:** Admin

**Query parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| p | int | Page number |
| user_id | int | Filter by user ID |
| status | int | Filter by status |


### 11.3 Search Subscriptions

**Endpoint:** `GET /api/subscription/search`

**Auth:** Admin

**Query parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| keyword | string | Search keyword |


### 11.4 Get Subscription Detail

**Endpoint:** `GET /api/subscription/:id`

**Auth:** Admin


### 11.5 Get Subscription Usage

**Endpoint:** `GET /api/subscription/:id/usage`

**Auth:** User


### 11.6 Create a Subscription

**Endpoint:** `POST /api/subscription/`

**Auth:** Admin

**Request body:**

```json
{
  "user_id": 2,
  "plan_id": 1,
  "billing_type": "token",
  "duration_days": 30,
  "notes": "Admin notes"
}
```


### 11.7 Update a Subscription

**Endpoint:** `PUT /api/subscription/`

**Auth:** Admin


### 11.8 Delete a Subscription

**Endpoint:** `DELETE /api/subscription/:id`

**Auth:** Admin
