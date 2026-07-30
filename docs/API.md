# One Api Pro 接口文档

> 基础URL: `http://<host>:<port>/api`
> 
> 所有 `/api/*` 接口返回统一 JSON 格式：
> 
> **成功响应：**
> ```json
> { "success": true, "message": "", "data": ... }
> ```
> 
> **错误响应：**
> ```json
> { "success": false, "message": "错误描述" }
> ```

---

## 目录

- [1. 模型定价管理 (Model Price)](#1-模型定价管理-model-price)
- [2. 分组折扣管理 (Group Price)](#2-分组折扣管理-group-price)
- [3. 分组列表 (Group)](#3-分组列表-group)
- [4. 系统选项 (Option)](#4-系统选项-option)
- [5. 渠道管理 (Channel)](#5-渠道管理-channel)
- [6. 令牌管理 (Token)](#6-令牌管理-token)
- [7. 用户管理 (User)](#7-用户管理-user)
- [8. 日志 (Log)](#8-日志-log)
- [9. 兑换码 (Redemption)](#9-兑换码-redemption)
- [10. 套餐管理 (Plan)](#10-套餐管理-plan)
- [11. 订阅管理 (Subscription)](#11-订阅管理-subscription)
- [12. OpenAI 兼容接口 (v1)](#12-openai-兼容接口-v1)
- [13. 其他公共接口](#13-其他公共接口)
- [附录 A：鉴权机制](#附录-a鉴权机制)
- [附录 B：权限等级说明](#附录-b权限等级说明)
- [附录 C：计费类型说明](#附录-c计费类型说明)
- [附录 D：渠道类型对照表](#附录-d渠道类型对照表)

---

## 附录 A：鉴权机制

One Api Pro 支持两种鉴权方式：**Cookie Session** 和 **Access Token**。不同接口使用不同的鉴权方式。

### A.1 Cookie Session 鉴权

适用于 `/api/*` 管理接口。

用户登录后，服务端创建 Session 并通过 `Set-Cookie` 响应头返回 Session ID。后续请求浏览器会自动携带 Cookie，无需手动设置。

**登录方式：**

```bash
# 登录获取 Cookie
curl -X POST http://localhost:3000/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"username":"root","password":"123456"}' \
  -c cookies.txt

# 使用 Cookie 访问
curl http://localhost:3000/api/user/self -b cookies.txt
```

### A.2 Access Token 鉴权

适用于 `/api/*` 管理接口（当无 Cookie 时自动降级为 Token 鉴权）。

每个用户有一个固定的 Access Token（UUID 格式），可通过 `/api/user/token` 接口生成，也可在管理后台的用户详情页查看。

**使用方式：**

```bash
curl http://localhost:3000/api/user/self \
  -H "Authorization: <access_token>"
```

**鉴权流程：**

1. 请求到达时，中间件先检查 Cookie Session 中的 `username` 字段
2. 若 Session 有效，直接通过鉴权
3. 若 Session 无效或不存在，读取 `Authorization` 请求头
4. 去除 `Bearer ` 前缀后，在数据库 `users` 表中查找 `access_token` 匹配的记录
5. 找到则通过鉴权，否则返回 `401 Unauthorized`

**优先级：** Cookie Session > Access Token

### A.3 API Key 鉴权（Bearer Token）

适用于 `/v1/*` OpenAI 兼容接口。

用户通过 `/api/token/` 创建的令牌（格式：`sk-<random>`），用于调用 OpenAI 兼容的 AI 模型接口。

**使用方式：**

```bash
curl http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer sk-xxxxxxxx"
```

**鉴权流程：**

1. 读取 `Authorization` 请求头，去除 `Bearer ` 和 `sk-` 前缀
2. 以第一个 `-` 为分隔符取前半部分作为 Token Key
3. 在数据库 `tokens` 表中查找匹配的令牌
4. 验证令牌状态（启用/禁用/过期/耗尽）
5. 若令牌设置了 `subnet` 限制，验证客户端 IP 是否在允许的子网内
6. 若令牌设置了 `models` 限制，验证请求的模型是否在允许列表中
7. 全部通过后，将 `userId`、`tokenId`、`tokenName` 等信息写入请求上下文

**指定渠道（管理员专用）：**

管理员可在 API Key 后追加 `-<channel_id>` 来指定使用特定渠道：

```
sk-xxxxxxxx-5    # 使用渠道 ID 为 5 的渠道
```

普通用户不支持此功能，会返回 `403 Forbidden`。

**令牌状态码：**

| 状态值 | 含义 |
|--------|------|
| 1 | 启用 |
| 2 | 禁用 |
| 3 | 过期 |
| 4 | 额度耗尽 |

### A.4 权限等级与接口对应关系

| 权限等级 | 值 | 对应中间件 | 可访问接口 |
|----------|------|-----------|------------|
| Guest | 0 | 无需认证 | 登录、注册、公开信息接口 |
| User | 1 | `UserAuth()` | 自身令牌/日志/订阅、可用模型 |
| Admin | 10 | `AdminAuth()` | 所有渠道、所有用户、所有日志、分组列表 |
| Root | 100 | `RootAuth()` | 系统选项、模型定价、分组折扣、套餐管理、管理员充值 |

**鉴权失败响应：**

| 场景 | HTTP 状态码 | 响应体 |
|------|------------|--------|
| 未登录且无 Token | 401 | `{"success":false,"message":"无权进行此操作，未登录且未提供 access token"}` |
| Token 无效 | 200 | `{"success":false,"message":"无权进行此操作，access token 无效"}` |
| 权限不足 | 200 | `{"success":false,"message":"无权进行此操作，权限不足"}` |
| 用户被封禁 | 200 | `{"success":false,"message":"用户已被封禁"}` |
| 令牌无效/过期 | 401 | `{"success":false,"message":"<具体错误>"}` |
| 令牌子网限制 | 403 | `{"success":false,"message":"该令牌只能在指定网段使用：xxx，当前 ip：xxx"}` |
| 令牌模型限制 | 403 | `{"success":false,"message":"该令牌无权使用模型：xxx"}` |
| 普通用户指定渠道 | 403 | `{"success":false,"message":"普通用户不支持指定渠道"}` |

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

## 2. 分组折扣管理 (Group Price)

管理不同用户分组的折扣系数。需要 Root 权限。

### 2.1 获取所有分组折扣

**接口：** `GET /api/group_price/`

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
      "group_name": "default",
      "model_name": "",
      "discount": 1.0,
      "created_at": 1718000000,
      "updated_at": 1718000000
    },
    {
      "id": 2,
      "group_name": "vip",
      "model_name": "gpt-4o",
      "discount": 0.8,
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
| group_name | string | 分组名称（如 default, vip, svip） |
| model_name | string | 模型名称，空字符串表示该分组所有模型的默认折扣 |
| discount | float | 折扣系数，1.0=无折扣，0.8=八折 |
| created_at | int64 | 创建时间（Unix 时间戳） |
| updated_at | int64 | 更新时间（Unix 时间戳） |

---

### 2.2 添加分组折扣

**接口：** `POST /api/group_price/`

**权限：** Root

**请求体：**

```json
{
  "group_name": "vip",
  "model_name": "gpt-4o",
  "discount": 0.8
}
```

**请求字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| group_name | string | 是 | 分组名称 |
| model_name | string | 否 | 模型名称，留空表示该分组所有模型的默认折扣 |
| discount | float | 否 | 折扣系数，默认 1.0 |

**返回值：**

```json
{
  "success": true,
  "message": ""
}
```

**错误情况：**
- `group_name` 为空：`{"success": false, "message": "分组名称不能为空"}`

---

### 2.3 更新分组折扣

**接口：** `PUT /api/group_price/`

**权限：** Root

**请求体：**

```json
{
  "id": 2,
  "group_name": "vip",
  "model_name": "gpt-4o",
  "discount": 0.7
}
```

**请求字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | uint | 是 | 记录ID |
| group_name | string | 否 | 分组名称 |
| model_name | string | 否 | 模型名称 |
| discount | float | 否 | 折扣系数 |

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

### 2.4 删除分组折扣

**接口：** `DELETE /api/group_price/:id`

**权限：** Root

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 分组折扣记录ID |

**返回值：**

```json
{
  "success": true,
  "message": ""
}
```

---

## 3. 分组列表 (Group)

### 3.1 获取所有分组名称

**接口：** `GET /api/group/`

**权限：** Admin

**请求参数：** 无

**返回值：**

```json
{
  "success": true,
  "message": "",
  "data": ["default", "vip", "svip"]
}
```

**说明：** 返回所有在 `group_prices` 表中定义的分组名称列表。

---

## 4. 系统选项 (Option)

### 4.1 获取所有系统选项

**接口：** `GET /api/option/`

**权限：** Root

**请求参数：** 无

**返回值：**

```json
{
  "success": true,
  "message": "",
  "data": [
    { "key": "QuotaForNewUser", "value": "1000000" },
    { "key": "QuotaPerUnit", "value": "500000" },
    { "key": "TopUpLink", "value": "" },
    ...
  ]
}
```

**说明：** 敏感选项（如 Token、SMTP 密码等）的值会被过滤或脱敏。

---

### 4.2 更新系统选项

**接口：** `PUT /api/option/`

**权限：** Root

**请求体：**

```json
{
  "key": "QuotaPerUnit",
  "value": "500000"
}
```

**请求字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| key | string | 是 | 选项名称 |
| value | string | 是 | 选项值 |

**可用选项列表：**

| 选项名 | 说明 |
|--------|------|
| QuotaForNewUser | 新用户赠送额度 |
| QuotaForInviter | 邀请人奖励额度 |
| QuotaForInvitee | 被邀请人奖励额度 |
| QuotaRemindThreshold | 额度提醒阈值 |
| PreConsumedQuota | 预消耗额度 |
| TopUpLink | 充值链接 |
| ChatLink | 聊天链接 |
| QuotaPerUnit | 单位额度（1元=多少内部额度） |
| DisplayInCurrencyEnabled | 是否以货币显示（true/false） |
| DisplayTokenStatEnabled | 是否显示Token统计（true/false） |
| ApproximateTokenEnabled | 是否使用近似Token计算（true/false） |
| RetryTimes | 重试次数 |
| LogConsumeEnabled | 是否启用消费日志（true/false） |
| AutomaticDisableChannelEnabled | 自动禁用渠道（true/false） |
| AutomaticEnableChannelEnabled | 自动启用渠道（true/false） |
| ChannelDisableThreshold | 渠道禁用阈值 |
| ChannelDefaultCooldownSeconds | 默认冷却时间（秒） |
| ChannelMaxCooldownSeconds | 最大冷却时间（秒） |
| ChannelConcurrencyEnabled | 启用渠道并发限制（true/false） |
| ChannelStickySessionEnabled | 启用粘性会话（true/false） |
| ErrorNext | 错误响应策略 JSON |

**返回值：**

```json
{
  "success": true,
  "message": ""
}
```

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

## 6. 令牌管理 (Token)

### 6.1 获取所有令牌

**接口：** `GET /api/token/`

**权限：** User（仅返回当前用户的令牌）

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| p | int | 页码，默认 0 |
| order | string | 排序字段 |

**返回值：**

```json
{
  "success": true,
  "message": "",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "key": "sk-xxxxxxxx",
      "status": 1,
      "name": "my-token",
      "created_time": 1718000000,
      "accessed_time": 1718000000,
      "expired_time": -1,
      "remain_quota": 500000,
      "unlimited_quota": false,
      "used_quota": 100000,
      "models": null,
      "subnet": null,
      "updated_at": 1718000000
    }
  ]
}
```

**返回字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 令牌ID |
| user_id | int | 所属用户ID |
| key | string | 令牌密钥 |
| status | int | 状态：1=启用, 2=禁用, 3=过期, 4=耗尽 |
| name | string | 令牌名称 |
| created_time | int64 | 创建时间 |
| accessed_time | int64 | 最后访问时间 |
| expired_time | int64 | 过期时间，-1=永不过期 |
| remain_quota | int64 | 剩余额度 |
| unlimited_quota | bool | 是否无限额度 |
| used_quota | int64 | 已用额度 |
| models | string/null | 允许的模型（逗号分隔），null=全部 |
| subnet | string/null | 允许的子网 |

---

### 6.2 搜索令牌

**接口：** `GET /api/token/search`

**权限：** User

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| keyword | string | 搜索关键词 |

---

### 6.3 获取单个令牌

**接口：** `GET /api/token/:id`

**权限：** User

---

### 6.4 创建令牌

**接口：** `POST /api/token/`

**权限：** User

**请求体：**

```json
{
  "name": "my-token",
  "remain_quota": 500000,
  "expired_time": -1,
  "unlimited_quota": false,
  "models": null,
  "subnet": null
}
```

---

### 6.5 更新令牌

**接口：** `PUT /api/token/`

**权限：** User

与创建格式相同，加上 `id` 字段。查询参数 `status_only=1` 时仅更新状态。

---

### 6.6 删除令牌

**接口：** `DELETE /api/token/:id`

**权限：** User

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

---

### 8.2 获取用户自身日志

**接口：** `GET /api/log/self`

**权限：** User

---

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

---

### 8.4 日志统计

**接口：** `GET /api/log/stat`

**权限：** Admin

**接口：** `GET /api/log/self/stat`

**权限：** User

---

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

---

## 9. 兑换码 (Redemption)

### 9.1 获取所有兑换码

**接口：** `GET /api/redemption/`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| p | int | 页码，默认 0 |

---

### 9.2 搜索兑换码

**接口：** `GET /api/redemption/search`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| keyword | string | 搜索关键词 |

---

### 9.3 获取单个兑换码

**接口：** `GET /api/redemption/:id`

**权限：** Admin

---

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

---

### 9.5 更新兑换码

**接口：** `PUT /api/redemption/`

**权限：** Admin

查询参数 `status_only=1` 时仅更新状态。

---

### 9.6 删除兑换码

**接口：** `DELETE /api/redemption/:id`

**权限：** Admin

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
      "features": "功能特性描述",
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
| features | string | 功能特性描述 |
| sort | int | 排序权重 |
| status | int | 状态：1=上架, 0=下架 |
| duration_days | int | 有效天数 |
| duration_text | string | 有效期显示文本 |
| recommended | bool | 是否推荐 |

---

### 10.2 搜索套餐

**接口：** `GET /api/plan/search`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| keyword | string | 搜索关键词 |

---

### 10.3 获取套餐详情

**接口：** `GET /api/plan/:id`

**权限：** Admin

---

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

---

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

---

### 10.6 删除套餐

**接口：** `DELETE /api/plan/:id`

**权限：** Root

---

## 11. 订阅管理 (Subscription)

### 11.1 获取用户订阅信息

**接口：** `GET /api/subscription/self`

**权限：** User

**返回值：** 当前用户的活跃订阅列表，包含使用量详情。

---

### 11.2 获取所有订阅

**接口：** `GET /api/subscription/`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| p | int | 页码 |
| user_id | int | 按用户ID过滤 |
| status | int | 按状态过滤 |

---

### 11.3 搜索订阅

**接口：** `GET /api/subscription/search`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| keyword | string | 搜索关键词 |

---

### 11.4 获取订阅详情

**接口：** `GET /api/subscription/:id`

**权限：** Admin

---

### 11.5 获取订阅使用量

**接口：** `GET /api/subscription/:id/usage`

**权限：** User

---

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

---

### 11.7 更新订阅

**接口：** `PUT /api/subscription/`

**权限：** Admin

---

### 11.8 删除订阅

**接口：** `DELETE /api/subscription/:id`

**权限：** Admin

---

## 12. OpenAI 兼容接口 (v1)

以下接口与 OpenAI API 格式兼容，使用 Bearer Token 认证。

### 12.1 模型列表

**接口：** `GET /v1/models`

**权限：** Bearer Token

**返回值：**

```json
{
  "object": "list",
  "data": [
    {
      "id": "gpt-4o",
      "object": "model",
      "created": 1718000000,
      "owned_by": "one-api-pro"
    }
  ]
}
```

---

### 12.2 获取模型详情

**接口：** `GET /v1/models/:model`

**权限：** Bearer Token

---

### 12.3 Chat Completions

**接口：** `POST /v1/chat/completions`

**权限：** Bearer Token

**请求体：**

```json
{
  "model": "gpt-4o",
  "messages": [
    {"role": "system", "content": "You are a helpful assistant."},
    {"role": "user", "content": "Hello"}
  ],
  "temperature": 0.7,
  "max_tokens": 4096,
  "stream": true
}
```

---

### 12.4 Text Completions

**接口：** `POST /v1/completions`

---

### 12.5 Embeddings

**接口：** `POST /v1/embeddings`

**接口：** `POST /v1/engines/:model/embeddings`

---

### 12.6 图片生成

**接口：** `POST /v1/images/generations`

---

### 12.7 音频转写

**接口：** `POST /v1/audio/transcriptions`

**接口：** `POST /v1/audio/translations`

**接口：** `POST /v1/audio/speech`

---

### 12.8 内容审核

**接口：** `POST /v1/moderations`

---

### 12.9 计费查询（OpenAI 兼容）

**接口：** `GET /v1/dashboard/billing/subscription`

**接口：** `GET /v1/dashboard/billing/usage?start_date=2024-01-01&end_date=2024-12-31`

---

### 12.10 代理转发

**接口：** `ANY /v1/oneapi/proxy/:channelid/*target`

**说明：** 直接代理到指定渠道。

---

## 13. 其他公共接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/status` | GET | 获取系统状态（无需认证） |
| `/api/notice` | GET | 获取系统公告 |
| `/api/about` | GET | 获取关于页面内容 |
| `/api/home_page_content` | GET | 获取首页内容 |
| `/api/verification` | GET | 发送邮箱验证码 |
| `/api/reset_password` | GET | 发送密码重置邮件 |
| `/api/user/reset` | POST | 重置密码 |
| `/api/oauth/github` | GET | GitHub OAuth 回调 |
| `/api/oauth/oidc` | GET | OIDC OAuth 回调 |
| `/api/oauth/lark` | GET | 飞书 OAuth 回调 |
| `/api/oauth/state` | GET | 生成 OAuth 状态码 |
| `/api/oauth/wechat` | GET | 微信 OAuth 回调 |
| `/api/oauth/wechat/bind` | GET | 微信账号绑定（需登录） |
| `/api/oauth/email/bind` | GET | 邮箱绑定（需登录） |
| `/api/user/logout` | GET | 退出登录 |
| `/api/user/aff` | GET | 获取推广码 |
| `/api/user/subscription` | GET | 获取订阅信息 |

---

## 附录 B：权限等级说明

| 等级 | 值 | 说明 |
|------|------|------|
| Guest | 0 | 未登录用户 |
| User | 1 | 普通用户 |
| Admin | 10 | 管理员 |
| Root | 100 | 超级管理员 |

---

## 附录 C：计费类型说明

| billing_type | 说明 |
|-------------|------|
| `token` | 按 Token 计费，价格单位为 ¥/百万tokens |
| `per_request` | 按次计费，价格单位为 ¥/次 |

**Token 计费公式：**

```
quota = ceil((inputPrice × inputTokens + outputPrice × completionTokens + cachedPrice × cachedTokens) / 1,000,000 × groupDiscount × QuotaPerUnit)
```

**按次计费公式：**

```
quota = ceil(perRequestPrice × sizeRatio × N × groupDiscount × QuotaPerUnit)
```

**分组折扣匹配规则：**
1. `GroupName + ModelName` 精确匹配 → 使用该折扣
2. `GroupName + ""(空)` → 使用该分组默认折扣
3. 无匹配 → 折扣为 1.0（无折扣）

**QuotaPerUnit 说明：** 默认 500,000，表示 500,000 内部额度 = 1 元人民币。

---

## 附录 D：渠道类型对照表

| Type | 适配器 |
|------|--------|
| 1 | OpenAI |
| 2 | Azure OpenAI |
| 3 | 自定义渠道 |
| 4 | Claude (Anthropic) |
| 5 | Google Gemini |
| 6 | 通义千问 (Ali) |
| 7 | 讯飞星火 (SparkDesk) |
| 8 | 百度文心 (Baidu) |
| 9 | 字节豆包 (Doubao) |
| 10 | MiniMax |
| 11 | DeepSeek |
| 12 | Cohere |
| 13 | 360 智脑 |
| 14 | Ollama |
| 15 | 月之暗面 (Moonshot) |
| 16 | 智谱AI (GLM) |
| 17 | 百里 (Baichuan) |
| 18 | 零一万物 (Yi) |
| 19 | DeepSeek (独立适配器) |
| 20-40+ | 其他适配器 |

---

## 附录 E：集群管理 API

集群 API 使用两种认证方式：

- **集群密钥（Cluster Secret）**：节点间内部通信，Header 格式 `X-Cluster-Secret: <节点 secret>`（目标节点的 secret，从本地 DB 查）
- **Root 管理员**：管理后台调用，使用 `Authorization: Bearer <Root Access Token>` 或 Cookie Session

| 接口 | 方法 | 说明 | 认证 |
|------|------|------|------|
| `/api/cluster/ping` | POST | 节点发现与心跳（双向注册） | 集群密钥 |
| `/api/cluster/sync` | POST | 接收同步事件 | 集群密钥 |
| `/api/cluster_node/` | GET | 查看节点列表（含本机、已禁用节点） | Root 管理员 |
| `/api/cluster_node/:id` | GET | 查看单个节点 | Root 管理员 |
| `/api/cluster_node/` | POST | 添加节点（含 secret） | Root 管理员 |
| `/api/cluster_node/` | PUT | 更新节点（可更新 secret） | Root 管理员 |
| `/api/cluster_node/:id` | DELETE | 软删除（禁用）节点 | Root 管理员 |
| `/api/cluster_node/:id/enable` | POST | 重新启用已禁用的节点 | Root 管理员 |
| `/api/cluster_node/ping/:id` | GET | 手动 ping 节点 | Root 管理员 |