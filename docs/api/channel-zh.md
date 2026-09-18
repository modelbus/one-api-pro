---
title: 渠道 API
description: 渠道 API：渠道管理 (Channel)
category: api
order: 13
---
## 5. 渠道管理 (Channel)

### 5.1 获取所有渠道

**接口：** `GET /api/channel/`

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

> **注意：** 非 Root 用户查看时，`key` 字段会被脱敏处理。

---

### 5.2 搜索渠道

**接口：** `GET /api/channel/search`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| keyword | string | 搜索关键词（匹配 ID 或名称前缀） |

**返回值：** 与 5.1 相同格式

---

### 5.3 获取单个渠道

**接口：** `GET /api/channel/:id`

**权限：** Admin

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 渠道ID |

**返回值：** 与 5.1 中单条数据格式相同

---

### 5.4 添加渠道

**接口：** `POST /api/channel/`

**权限：** Admin

**请求体：**

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

**请求字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | int | 是 | 渠道类型（1=OpenAI, 3=Azure, 等） |
| key | string | 是 | API密钥，多个用换行分隔 |
| name | string | 是 | 渠道名称 |
| base_url | string | 否 | API基础URL |
| models | string | 是 | 支持的模型，逗号分隔 |
| group | string | 否 | 分组，默认 "default" |
| weight | int | 否 | 权重，默认 0 |
| priority | int | 否 | 优先级，默认 0 |
| model_mapping | string | 否 | 模型映射JSON |
| config | string | 否 | 渠道配置JSON |
| system_prompt | string | 否 | 系统提示词 |
| max_concurrency | int | 否 | 最大并发数，0=不限 |
| cooldown_seconds | int | 否 | 冷却时间（秒），默认60 |
| rpm | int | 否 | 每分钟请求数限制 |

**返回值：**

```json
{
  "success": true,
  "message": ""
}
```

---

### 5.5 更新渠道

**接口：** `PUT /api/channel/`

**权限：** Admin

**请求体：** 与 5.4 相同格式，需包含 `id` 字段。

---

### 5.6 删除渠道

**接口：** `DELETE /api/channel/:id`

**权限：** Admin

---

### 5.7 删除已禁用渠道

**接口：** `DELETE /api/channel/disabled`

**权限：** Admin

**说明：** 删除所有状态为禁用（手动禁用+自动禁用）的渠道。

---

### 5.8 测试渠道

**接口：** `GET /api/channel/test/:id`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| model | string | 可选，指定测试使用的模型 |

**接口：** `GET /api/channel/test`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| scope | string | 可选，测试范围：`all`、`limited`、`disabled` |

---

### 5.9 更新渠道余额

**接口：** `GET /api/channel/update_balance/:id`

**权限：** Admin

**说明：** 更新指定渠道的余额。

**接口：** `GET /api/channel/update_balance`

**权限：** Admin

**说明：** 更新所有渠道的余额。

---
