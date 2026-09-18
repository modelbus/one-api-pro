---
title: 套餐 API
description: "套餐 API：套餐管理 (Plan)"
category: api
order: 18
---
## 10. 套餐管理 (Plan)

### 10.1 获取所有套餐

**接口：** `GET /api/plan/`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| p | int | 页码，默认 0 |

**返回值：**

```json
{
  "success": true,
  "message": "",
  "data": [
    {
      "id": 1,
      "name": "基础套餐",
      "price": 99.00,
      "tokens": 500000000,
      "model_limits": "{\"gpt-4o\":{\"request_month\":1000,\"token_month\":50000000}}",
      "default_model": "gpt-4o",
      "description": "基础套餐描述",
"features": ["API 调用 1000 次/月", "支持 GPT-4o"],
      "sort": 0,
      "status": 1,
      "duration_days": 30,
      "duration_text": "30天",
      "recommended": false,
      "created_time": 1718000000
    }
  ]
}
```

**返回字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 套餐ID |
| name | string | 套餐名称 |
| price | float64 | 价格 |
| tokens | int64 | Token配额 |
| model_limits | string | 模型限额配置JSON，key为模型名称，value为ModelLimitRule |
| default_model | string | 默认模型名称，不在model_limits中的请求模型将转发至此模型计费；为空则不转发，未配置模型返回422 |
| description | string | 描述 |
| features | `array<string>` | 功能特性列表，每项一行展示在用户端套餐卡 |
| sort | int | 排序权重 |
| status | int | 状态：1=上架, 0=下架 |
| duration_days | int | 有效天数 |
| duration_text | string | 有效期显示文本 |
| recommended | bool | 是否推荐 |


### 10.2 搜索套餐

**接口：** `GET /api/plan/search`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| keyword | string | 搜索关键词 |


### 10.3 获取套餐详情

**接口：** `GET /api/plan/:id`

**权限：** Admin


### 10.4 创建套餐

**接口：** `POST /api/plan/`

**权限：** Root

**请求体：**

```json
{
  "name": "基础套餐",
  "price": 99.00,
  "tokens": 500000000,
  "model_limits": "{\"gpt-4o\":{\"request_month\":1000,\"token_month\":50000000}}",
  "default_model": "gpt-4o",
  "description": "基础套餐描述",
  "features": "功能特性描述",
  "sort": 0,
  "status": 1,
  "duration_days": 30,
  "duration_text": "30天",
  "recommended": false
}
```


### 10.5 更新套餐

**接口：** `PUT /api/plan/`

**权限：** Root

与创建格式相同，需包含 `id` 字段。

**`model_limits` 字段格式说明：**

key 为模型名称，value 为 `ModelLimitRule` 对象：

```json
{
  "gpt-4o": {
    "period_h": 5,
    "request_period": 100,
    "request_week": 500,
    "request_month": 2000,
    "token_period": 500000,
    "token_week": 2000000,
    "token_month": 10000000
  }
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| period_h | int | 滚动周期时长（小时），默认5 |
| request_period | int64 | 周期内最大请求数 |
| request_week | int64 | 周内最大请求数 |
| request_month | int64 | 月内最大请求数 |
| token_period | int64 | 周期内最大Token数 |
| token_week | int64 | 周内最大Token数 |
| token_month | int64 | 月内最大Token数 |

**`default_model` 字段说明：**

- 当用户请求的模型不在 `model_limits` 中时，系统会将请求转发至 `default_model` 指定的模型
- `default_model` 必须是 `model_limits` 中已配置的模型名称，否则创建/更新套餐时会报错
- 如果 `default_model` 为空且用户请求的模型不在 `model_limits` 中，将返回 422 错误
- 转发后，实际请求将使用 `default_model` 发送到上游，日志中记录的模型名称也是 `default_model`


### 10.6 删除套餐

**接口：** `DELETE /api/plan/:id`

**权限：** Root

