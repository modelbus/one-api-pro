---
title: 日志 API
description: "日志 API：日志 (Log)"
category: api
order: 16
---
## 8. 日志 (Log)

### 8.1 获取所有日志

**接口：** `GET /api/log/`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| p | int | 页码，默认 0 |
| type | int | 日志类型：0=全部, 1=充值, 2=消费, 3=管理, 4=系统 |
| start_timestamp | int64 | 起始时间戳 |
| end_timestamp | int64 | 结束时间戳 |
| model_name | string | 模型名过滤 |
| username | string | 用户名过滤 |
| token_name | string | 令牌名过滤 |
| channel | int | 渠道ID过滤 |
| group | string | 分组过滤 |

**返回字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 日志ID |
| user_id | int | 用户ID |
| created_at | int64 | 创建时间 |
| type | int | 日志类型 |
| content | string | 日志内容 |
| username | string | 用户名 |
| token_name | string | 令牌名 |
| model_name | string | 模型名 |
| quota | int | 消耗额度 |
| prompt_tokens | int | 输入Token数 |
| completion_tokens | int | 输出Token数 |
| cached_tokens | int | 缓存Token数 |
| channel_id | int | 渠道ID |
| request_id | string | 请求ID |
| elapsed_time | int64 | 耗时（ms） |
| is_stream | bool | 是否流式 |
| billing_source | int | 计费来源：0=普通, 1=订阅 |
| plan_id | int | 套餐ID |
| session_key | string | 会话Key |


### 8.2 获取用户自身日志

**接口：** `GET /api/log/self`

**权限：** User


### 8.3 搜索日志

**接口：** `GET /api/log/search`

**权限：** Admin

**接口：** `GET /api/log/self/search`

**权限：** User

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| keyword | string | 搜索关键词 |
| type | int | 日志类型（可选） |


### 8.4 日志统计

**接口：** `GET /api/log/stat`

**权限：** Admin

**接口：** `GET /api/log/self/stat`

**权限：** User


### 8.5 删除历史日志

**接口：** `DELETE /api/log/`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| target_timestamp | int64 | 删除此时间戳之前的日志 |

**返回值：**

```json
{
  "success": true,
  "message": "",
  "data": 100
}
```

> `data` 为删除的日志条数。

