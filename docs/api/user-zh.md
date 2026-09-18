---
title: 用户 API
description: 用户 API：用户管理 (User)
category: api
order: 15
---
## 7. 用户管理 (User)

### 7.1 用户注册

**接口：** `POST /api/user/register`

**权限：** 无（公开）

**请求体：**

```json
{
  "username": "newuser",
  "password": "password123",
  "email": "user@example.com",
  "verification_code": "123456"
}
```

---

### 7.2 用户登录

**接口：** `POST /api/user/login`

**权限：** 无（公开）

**请求体：**

```json
{
  "username": "root",
  "password": "123456"
}
```

**返回值：**

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

---

### 7.3 获取当前用户信息

**接口：** `GET /api/user/self`

**权限：** User

---

### 7.4 更新当前用户信息

**接口：** `PUT /api/user/self`

**权限：** User

**请求体：**

```json
{
  "display_name": "New Name",
  "password": "newpassword"
}
```

---

### 7.5 删除当前用户

**接口：** `DELETE /api/user/self`

**权限：** User

**说明：** Root 用户不可自删。

---

### 7.6 获取用户仪表盘

**接口：** `GET /api/user/dashboard`

**权限：** User

**返回值：** 7天内的使用统计。

---

### 7.7 生成访问令牌

**接口：** `GET /api/user/token`

**权限：** User

**说明：** 生成一个新的 UUID 格式访问令牌。

---

### 7.8 获取推广码

**接口：** `GET /api/user/aff`

**权限：** User

---

### 7.9 充值（用户兑换码）

**接口：** `POST /api/user/topup`

**权限：** User

**请求体：**

```json
{
  "key": "redemption-code"
}
```

---

### 7.10 获取可用模型

**接口：** `GET /api/user/available_models`

**权限：** User

**返回值：** 返回当前用户分组可用的模型列表。

---

### 7.11 管理员操作

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/user/` | GET | 获取所有用户（分页 `?p=N`） |
| `/api/user/search` | GET | 搜索用户 |
| `/api/user/:id` | GET | 获取指定用户 |
| `/api/user/` | POST | 创建用户 |
| `/api/user/manage` | POST | 管理用户（禁用/启用/删除/提升/降级） |
| `/api/user/` | PUT | 更新用户（管理员） |
| `/api/user/:id` | DELETE | 删除用户 |
| `/api/topup` | POST | 管理员充值 |

**管理员充值请求体：**

```json
{
  "user_id": 2,
  "quota": 500000,
  "remark": "充值备注"
}
```

**管理用户请求体：**

```json
{
  "username": "testuser",
  "action": "disable"
}
```

**action 可选值：** `disable`、`enable`、`delete`、`promote`、`demote`

---
