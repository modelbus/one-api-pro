---
title: 兑换码 API
description: "兑换码 API：兑换码 (Redemption)"
category: api
order: 17
---
## 9. 兑换码 (Redemption)

### 9.1 获取所有兑换码

**接口：** `GET /api/redemption/`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| p | int | 页码，默认 0 |


### 9.2 搜索兑换码

**接口：** `GET /api/redemption/search`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| keyword | string | 搜索关键词 |


### 9.3 获取单个兑换码

**接口：** `GET /api/redemption/:id`

**权限：** Admin


### 9.4 创建兑换码（批量）

**接口：** `POST /api/redemption/`

**权限：** Admin

**请求体：**

```json
{
  "name": "兑换码名称",
  "quota": 500000,
  "count": 10
}
```

**请求字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 兑换码名称 |
| quota | int64 | 是 | 兑换额度 |
| count | int | 否 | 批量创建数量（1-100），默认 1 |

**返回值：** 返回创建的兑换码列表，包含自动生成的 `key`。


### 9.5 更新兑换码

**接口：** `PUT /api/redemption/`

**权限：** Admin

查询参数 `status_only=1` 时仅更新状态。


### 9.6 删除兑换码

**接口：** `DELETE /api/redemption/:id`

**权限：** Admin

