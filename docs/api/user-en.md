---
title: User API
description: "All /api/user/* endpoints."
category: api
order: 15
---
## 7. User Management (User)

### 7.1 Register

**Endpoint:** `POST /api/user/register`

**Auth:** None (public)

**Request body:**

```json
{
  "username": "newuser",
  "password": "password123",
  "email": "user@example.com",
  "verification_code": "123456"
}
```


### 7.2 Login

**Endpoint:** `POST /api/user/login`

**Auth:** None (public)

**Request body:**

```json
{
  "username": "root",
  "password": "123456"
}
```

**Response:**

```json
{
  "success": true,
  "message": "",
  "data": {
    "id": 1,
    "username": "root",
    "display_name": "Root User",
    "role": 100,
    "status": 1,
    "email": "",
    "quota": 500000000000000,
    "access_token": "uuid-token-string",
    "group": "default",
    "aff_code": "abcdef"
  }
}
```


### 7.3 Get Current User Info

**Endpoint:** `GET /api/user/self`

**Auth:** User


### 7.4 Update Current User Info

**Endpoint:** `PUT /api/user/self`

**Auth:** User

**Request body:**

```json
{
  "display_name": "New Name",
  "password": "newpassword"
}
```


### 7.5 Delete Current User

**Endpoint:** `DELETE /api/user/self`

**Auth:** User

**Note:** The root user cannot self-delete.


### 7.6 Get User Dashboard

**Endpoint:** `GET /api/user/dashboard`

**Auth:** User

**Response:** Usage statistics for the past 7 days.


### 7.7 Generate Access Token

**Endpoint:** `GET /api/user/token`

**Auth:** User

**Note:** Generates a new UUID-format access token.


### 7.8 Get Affiliate Code

**Endpoint:** `GET /api/user/aff`

**Auth:** User


### 7.9 Redeem (User Redemption Code)

**Endpoint:** `POST /api/user/topup`

**Auth:** User

**Request body:**

```json
{
  "key": "redemption-code"
}
```


### 7.10 Get Available Models

**Endpoint:** `GET /api/user/available_models`

**Auth:** User

**Response:** Models available to the caller's user group.


### 7.11 Admin Operations

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/user/` | GET | List all users (paginated with `?p=N`) |
| `/api/user/search` | GET | Search users |
| `/api/user/:id` | GET | Get a specific user |
| `/api/user/` | POST | Create a user |
| `/api/user/manage` | POST | Manage a user (disable / enable / delete / promote / demote) |
| `/api/user/` | PUT | Update a user (admin) |
| `/api/user/:id` | DELETE | Delete a user |
| `/api/topup` | POST | Admin top-up |

**Admin top-up request body:**

```json
{
  "user_id": 2,
  "quota": 500000,
  "remark": "Topup remark"
}
```

**Manage user request body:**

```json
{
  "username": "testuser",
  "action": "disable"
}
```

**Allowed `action` values:** `disable`, `enable`, `delete`, `promote`, `demote`
