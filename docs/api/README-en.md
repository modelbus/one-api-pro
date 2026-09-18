---
title: API Reference · Overview
description: "Auth mechanisms, response format, permission tiers and glossary."
category: api
order: 1
---

# API Reference · Overview

This document aggregates every `/api/*` and `/v1/*` endpoint in One API Pro, together with the shared conventions for auth, permissions and billing.

## Appendix A: Auth Mechanisms

One API Pro supports two auth methods: **Cookie Session** and **Access Token**. Different endpoints use different auth schemes.

### A.1 Cookie Session

Used for the `/api/*` admin endpoints.

After login, the server creates a session and returns the session ID through the `Set-Cookie` response header. The browser carries the cookie automatically on subsequent requests, no manual setup required.

**Login:**

```bash
# Login to obtain the cookie
curl -X POST http://localhost:3000/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"username":"root","password":"123456"}' \
  -c cookies.txt

# Use the cookie
curl http://localhost:3000/api/user/self -b cookies.txt
```

### A.2 Access Token

Used for the `/api/*` admin endpoints (falls back to token auth when no cookie is present).

Every user has a fixed Access Token (UUID format), generated through `/api/user/token` and visible in the admin user detail page.

**Usage:**

```bash
curl http://localhost:3000/api/user/self \
  -H "Authorization: <access_token>"
```

**Auth flow:**

1. When a request arrives, the middleware first checks the `username` field in the Cookie Session.
2. If the session is valid, auth passes immediately.
3. If the session is invalid or missing, the `Authorization` request header is read.
4. After stripping the `Bearer ` prefix, the `users` table is searched for a record with a matching `access_token`.
5. If found, auth passes; otherwise `401 Unauthorized` is returned.

**Priority:** Cookie Session > Access Token.

### A.3 API Key (Bearer Token)

Used for the `/v1/*` OpenAI-compatible endpoints.

Tokens created through `/api/token/` (format: `sk-<random>`) are used to call the OpenAI-compatible model endpoints.

**Usage:**

```bash
curl http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer sk-xxxxxxxx"
```

**Auth flow:**

1. Read the `Authorization` header, strip the `Bearer ` and `sk-` prefixes.
2. Take the part before the first `-` as the token key.
3. Search the `tokens` table for a matching token.
4. Validate the token status (enabled / disabled / expired / exhausted).
5. If the token has a `subnet` restriction, verify the client IP is inside it.
6. If the token has a `models` restriction, verify the requested model is in the allow list.
7. When all checks pass, write `userId`, `tokenId`, `tokenName`, etc. into the request context.

**Specifying a channel (admin only):**

Admins can append `-<channel_id>` after the API key to pin the request to a specific channel:

```
sk-xxxxxxxx-5    # Use the channel with id=5
```

Non-admin users get `403 Forbidden` if they try this.

**Token status codes:**

| Status | Meaning |
|--------|---------|
| 1 | Enabled |
| 2 | Disabled |
| 3 | Expired |
| 4 | Quota exhausted |

### A.4 Permission Levels and Endpoint Mapping

| Level | Value | Middleware | Accessible endpoints |
|-------|-------|------------|----------------------|
| Guest | 0 | No auth | Login, register, public info endpoints |
| User | 1 | `UserAuth()` | Own tokens / logs / subscription, available models |
| Admin | 10 | `AdminAuth()` | All channels, all users, all logs, group list |
| Root | 100 | `RootAuth()` | System options, model prices, group prices, plan admin, admin top-up |

**Auth failure responses:**

| Scenario | HTTP status | Body |
|----------|-------------|------|
| Not logged in and no token | 401 | `{"success":false,"message":"无权进行此操作，未登录且未提供 access token"}` |
| Token invalid | 200 | `{"success":false,"message":"无权进行此操作，access token 无效"}` |
| Insufficient permission | 200 | `{"success":false,"message":"无权进行此操作，权限不足"}` |
| User banned | 200 | `{"success":false,"message":"用户已被封禁"}` |
| Token invalid / expired | 401 | `{"success":false,"message":"<specific error>"}` |
| Token subnet restriction | 403 | `{"success":false,"message":"该令牌只能在指定网段使用：xxx，当前 ip：xxx"}` |
| Token model restriction | 403 | `{"success":false,"message":"该令牌无权使用模型：xxx"}` |
| Non-admin specifying channel | 403 | `{"success":false,"message":"普通用户不支持指定渠道"}` |
