---
title: 模型定价 API
description: 模型定价 API：模型定价管理 (Model Price)
category: api
order: 11
---
## 1. 模型定价管理 (Model Price)

管理模型的 Token 定价（¥/百万tokens）和按次定价。需要 Root 权限。

### 1.1 获取所有模型定价

**接口：** `GET /api/model_price/`

**权限：** Root

**请求参数：** 无

**返回值：**

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

**返回字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 记录ID |
| model_name | string | 模型名称（唯一） |
| input_price | float | 输入价格（¥/百万tokens） |
| output_price | float | 输出价格（¥/百万tokens） |
| cached_price | float | 缓存价格（¥/百万tokens），0 表示不支持 |
| per_request_price | float | 按次价格（¥/次） |
| billing_type | string | 计费类型：`token` 或 `per_request` |
| enabled | bool | 是否启用 |
| created_at | int64 | 创建时间（Unix 时间戳） |
| updated_at | int64 | 更新时间（Unix 时间戳） |

---

### 1.2 添加模型定价

**接口：** `POST /api/model_price/`

**权限：** Root

**请求体：**

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

**请求字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| model_name | string | 是 | 模型名称，不可重复 |
| input_price | float | 否 | 输入价格，默认 0 |
| output_price | float | 否 | 输出价格，默认 0 |
| cached_price | float | 否 | 缓存价格，默认 0 |
| per_request_price | float | 否 | 按次价格，默认 0 |
| billing_type | string | 否 | `token`（默认）或 `per_request` |
| enabled | bool | 否 | 默认 true |

**返回值：**

```json
{
  "success": true,
  "message": ""
}
```

**错误情况：**
- `model_name` 为空：`{"success": false, "message": "模型名称不能为空"}`
- `model_name` 已存在：数据库唯一约束报错

---

### 1.3 更新模型定价

**接口：** `PUT /api/model_price/`

**权限：** Root

**请求体：**

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

**请求字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | uint | 是 | 记录ID |
| model_name | string | 否 | 模型名称（更新时不可修改，用于标识） |
| input_price | float | 否 | 输入价格 |
| output_price | float | 否 | 输出价格 |
| cached_price | float | 否 | 缓存价格 |
| per_request_price | float | 否 | 按次价格 |
| billing_type | string | 否 | 计费类型 |
| enabled | bool | 否 | 是否启用 |

**返回值：**

```json
{
  "success": true,
  "message": ""
}
```

**错误情况：**
- `id` 为 0：`{"success": false, "message": "ID不能为空"}`

---

### 1.4 删除模型定价

**接口：** `DELETE /api/model_price/:id`

**权限：** Root

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 模型定价记录ID |

**返回值：**

```json
{
  "success": true,
  "message": ""
}
```

**错误情况：**
- `id` 无效：`{"success": false, "message": "无效的ID"}`

---
