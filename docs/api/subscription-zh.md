---
title: 订阅 API
description: "订阅 API：订阅管理 (Subscription)"
category: api
order: 19
---
## 11. 订阅管理 (Subscription)

### 11.1 获取用户订阅信息

**接口：** `GET /api/subscription/self`

**权限：** User

**返回值：** 当前用户的活跃订阅列表，包含使用量详情。

### 11.2 获取所有订阅

**接口：** `GET /api/subscription/`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| p | int | 页码 |
| user_id | int | 按用户ID过滤 |
| status | int | 按状态过滤 |

### 11.3 搜索订阅

**接口：** `GET /api/subscription/search`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| keyword | string | 搜索关键词 |

### 11.4 获取订阅详情

**接口：** `GET /api/subscription/:id`

**权限：** Admin

### 11.5 获取订阅使用量

**接口：** `GET /api/subscription/:id/usage`

**权限：** User

### 11.6 创建订阅

**接口：** `POST /api/subscription/`

**权限：** Admin

**请求体：**

```json
{
  "user_id": 2,
  "plan_id": 1,
  "billing_type": "token",
  "duration_days": 30,
  "notes": "管理员备注"
}
```

### 11.7 更新订阅

**接口：** `PUT /api/subscription/`

**权限：** Admin

### 11.8 删除订阅

**接口：** `DELETE /api/subscription/:id`

**权限：** Admin

