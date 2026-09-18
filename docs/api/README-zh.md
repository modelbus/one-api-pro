---
title: API 参考 · 总览
description: "鉴权机制、响应格式、权限等级与术语表"
category: api
order: 1
---

# API 参考 · 总览

本文档汇总 One API Pro 的所有 `/api/*` 与 `/v1/*` 接口，以及鉴权、权限、计费等公共约定。

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

