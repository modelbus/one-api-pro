---
title: 分组折扣 API
description: "分组折扣 API：分组折扣管理 (Group Price)"
category: api
order: 12
---
## 2. 分组折扣管理 (Group Price)

管理不同用户分组的折扣系数。需要 Root 权限。

### 2.1 获取所有分组折扣

**接口：** `GET /api/group_price/`

**权限：** Root

**请求参数：** 无

**返回值：**

```

**返回字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 记录ID |
| group_name | string | 分组名称（如 default, vip, svip） |
| model_name | string | 模型名称，空字符串表示该分组所有模型的默认折扣 |
| discount | float | 折扣系数，1.0=无折扣，0.8=八折 |
| created_at | int64 | 创建时间（Unix 时间戳） |
| updated_at | int64 | 更新时间（Unix 时间戳） |

### 2.2 添加分组折扣

**接口：** `POST /api/group_price/`

**权限：** Root

**请求体：**

```

**请求字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| group_name | string | 是 | 分组名称 |
| model_name | string | 否 | 模型名称，留空表示该分组所有模型的默认折扣 |
| discount | float | 否 | 折扣系数，默认 1.0 |

**返回值：**

```

**错误情况：**
- `group_name` 为空：`{"success": false, "message": "分组名称不能为空"}`

### 2.3 更新分组折扣

**接口：** `PUT /api/group_price/`

**权限：** Root

**请求体：**

```

**请求字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | uint | 是 | 记录ID |
| group_name | string | 否 | 分组名称 |
| model_name | string | 否 | 模型名称 |
| discount | float | 否 | 折扣系数 |

**返回值：**

```

**错误情况：**
- `id` 为 0：`{"success": false, "message": "ID不能为空"}`

### 2.4 删除分组折扣

**接口：** `DELETE /api/group_price/:id`

**权限：** Root

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 分组折扣记录ID |

**返回值：**

```

