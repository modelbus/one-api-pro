---
title: Access Token
description: "Manage your system-level Access Token."
category: user
order: 4
---

# Access Token

> Manage your system-level Access Token.
> 管理用户的系统级 Access Token。

The Access Token is a **UUID** used as the auth header for `/api/*` admin endpoints. When the browser cookie is unavailable (scripts, CI), the Access Token is the canonical way to authenticate.

Access Token 是 **UUID 格式** 的字符串，用作 `/api/*` 管理接口的鉴权头；当浏览器 Cookie 失效时（如脚本调用、CI），使用 Access Token 是最稳妥的方式。

## API Key vs Access Token / 与 API Key 的区别

| Aspect / 维度 | Access Token | API Key (`sk-…`) |
| --- | --- | --- |
| Format / 格式 | UUID string | `sk-` + random |
| Used for / 用途 | `/api/*` admin endpoints | `/v1/*` OpenAI-compatible endpoints |
| Auth / 鉴权 | `Authorization: <uuid>` | `Authorization: Bearer sk-…` |
| Issued by / 颁发方 | User self-generates in Personal Center | User creates in the Tokens page |
| Rotate / 重生 | Regenerate anytime, old token invalidated | Delete + recreate, old key invalidated |

## Generate / 生成

UI: Personal Center → Access Token → "Generate".

UI：个人中心 → Access Token → 「生成」。

```bash
curl http://localhost:3000/api/user/token -b cookies.txt
# { success:true, data:"<uuid>" }
```

Shown only once — copy it immediately; refreshing the page loses it.

仅显示一次，必须立刻复制保存；刷新页面后再访问就拿不到了。

## Usage / 使用

```bash
# Admin endpoints
curl http://localhost:3000/api/user/self \
  -H "Authorization: <your_access_token>"

# Admin list (requires admin role)
curl http://localhost:3000/api/user/list \
  -H "Authorization: <your_access_token>"
```

Backend middleware priority: Cookie Session first, Access Token fallback (see [API Overview · Auth](/en/api/README#appendix-a-auth-mechanisms)).

后端中间件鉴权优先级：Cookie Session 优先 → Access Token 兜底。

## Rotation / 重新生成

Click "Generate" again to mint a new UUID; the old one is invalidated immediately — a useful emergency-rotation lever.

再次点「生成」会生成新的 UUID，旧的立即失效。这是一种「紧急撤销」手段。

## With SDKs / 配合 SDK

Two common usages in scripts / CI:

脚本 / CI 里常见的两种用法：

```bash
# 1. Export as env var
export OAP_TOKEN=$(curl -s http://localhost:3000/api/user/token -b cookies.txt | jq -r .data)

# 2. Store in a Secret Manager
gh secret set OAP_TOKEN < token.txt
```

## Security Practices / 安全实践

| Recommendation / 建议 | Notes / 说明 |
| --- | --- |
| **Treat it like a password** | A leak lets attackers read or modify your `/api/*` resources |
| **Never commit it** | `.env` and CI temp files should be `.gitignore`d |
| **Rotate periodically** | Every 90 days is reasonable; rotate immediately on exposure |
| **One token per environment** | Don't share between prod, staging and test |
| **Pair with rate-limit** | `GLOBAL_API_RATE_LIMIT` is a final backstop |

## FAQ / 常见问题

- "Generate" button does nothing: verify the user can call `/api/user/token` (any role works).
  「生成」按钮没反应：检查用户是否可调用 `/api/user/token`（任意角色都可以）。
- Old token still works after rotation: cache propagation may take a few seconds; refresh and clear local credentials immediately.
  重生成后旧 Token 仍可调用：缓存层延迟几秒；建议立即刷新页面并清掉本地凭证。
- Can I share with someone else? No — the token is tied to your `users.id`; sharing = impersonation.
  能否给其他人用：Access Token 与 `users.id` 一一对应；他人使用 = 越权。

Next: [API Key (Tokens)](/en/api/token) · [My Orders](/en/user/orders) · [Chat Playground](/en/user/chat).