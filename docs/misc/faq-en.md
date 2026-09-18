---
title: FAQ
description: "Common questions for users and admins."
category: misc
order: 1
---

# FAQ

> Common questions for users and admins.
> 用户与管理员的常见问答。

## User / 用户

### Q1: Where is my API Key?

The "New Token" page creates an **OpenAI-compatible key** (`sk-` prefix, used for `/v1/*`). The UUID-style **Access Token** lives in Personal Center, used for `/api/*`. They are not interchangeable.

`/token` 页新建的 Key 是 **OpenAI 兼容 Key**（`sk-` 前缀），用于 `/v1/*` 调用。UUID 格式的 **Access Token** 在个人中心生成，用于 `/api/*` 管理接口。两者不可混用。

### Q2: Why does my balance not change?

Admins may have granted only Token quota (`tokens.remain_quota`), not user quota (`users.quota`). Calls first deduct from the Token quota, then from the user quota. See [Architecture · cache consistency](/en/start/architecture#request-lifecycle-request-lifecycle).

管理员可能只给了 Token 额度（`tokens.remain_quota`），没有给用户额度（`users.quota`）。调用会优先扣 Token 的 `remain_quota`；Token 用尽后再扣用户。

### Q3: Can I keep using the API after my plan expires?

You can keep calling until tokens are exhausted, but the plan's discount ratio is no longer applied, and the `daily_quota` progress bar drops to 0.

可以继续调用到令牌耗尽，但不再享受套餐折扣倍率；且 `daily_quota` 进度条变为 0。

### Q4: My invite link isn't giving any reward?

`/api/user/aff` returns the invite URL. After the invitee registers, `IncreaseUserQuota` credits the bonus. Since v0.0.21 the Redis cache is refreshed synchronously — **the next call sees the new balance immediately**.

`/api/user/aff` 返回的 URL 即邀请链接；被邀请人成功注册后会调用 `IncreaseUserQuota` 增加奖励。v0.0.21 后此链路已同步刷 Redis 缓存，**下个请求立即可见**。

### Q5: I never received the password-reset email.

Currently emails are only logged via `logger.SysLog`; integrate SMTP / a third-party provider in production (see [Profile · Password Reset](/en/user/profile#password-reset--password-reset)).

当前邮件发送以 `logger.SysLog` 形式输出，**生产环境请接 SMTP / 第三方**。

## Admin / 管理员

### Q1: Channel test passes (200) but `/v1/*` returns 401?

Usually a **Base URL / Key mismatch**, or the upstream account is banned or geo-restricted. Check the `x-request-id` from the "Test" page response when opening a support ticket.

通常是 **Base URL / Key 不匹配**，或上游账号失效、地区限制。检查渠道「测试」页响应头里的 `x-request-id` 回到上游工单。

### Q2: Can I restrict a user group to specific models?

Set the token's "Models" allow-list, or configure `models` on the channel (model + user-group matrix). Enforced by `relay/middleware/should-check-model`.

设置 Token 时填写「可用模型」白名单，或在渠道侧配置 `models` 限制（按模型 + 用户组双维度）。中间件 `relay/middleware/should-check-model` 强制校验。

### Q3: How do I configure upgrade differential?

System option `plan.upgrade_mode` values:

系统选项 `plan.upgrade_mode` 取值：

| Value / 值 | Meaning / 含义 |
| --- | --- |
| `price_diff` (default) | Pay the differential; `UP` prefix orders |
| `stack` | Full-price stack; the old subscription keeps running |

### Q4: Plan vs top-up — which is consumed first?

**Plan first**: deduction order is `token.remain_quota` → `users.quota` (top-up). Insufficient balance returns 403. When a plan has `daily_quota`, the daily cap plus plan ratio are applied first against `users.quota`.

**套餐优先**：调用扣减顺序 `token.remain_quota` → `users.quota`（充值）→ 不够则 403。套餐有 `daily_quota` 时，按日限额 + 套餐倍率再扣 `users.quota`。

### Q5: How does an admin top-up another user?

In the admin "User Detail" → "Adjust quota" page, or call `POST /api/user/<id>/quota` directly. Both paths refresh Redis synchronously since v0.0.21, so the next call sees the new balance immediately.

管理后台「用户详情」→ 「调整额度」，或 `POST /api/user/<id>/quota` 直接调；这两条链路在 v0.0.21 后都已同步刷 Redis，下一次请求立即生效。

## Deployment / 部署

### Q1: Port 3000 is taken after Docker start?

Use `-p 8080:3000`; if you reverse-proxy via Nginx, update `proxy_pass http://127.0.0.1:8080;`.

修改 `-p 8080:3000`；同时若用 Nginx 反代，`proxy_pass http://127.0.0.1:8080;`。

### Q2: Can I run without Redis?

Yes — without `REDIS_CONN_STRING`, `RedisEnabled=false` and all caches hit DB directly. Under heavy concurrency this noticeably slows down `users.quota` deductions; Redis is recommended.

可以，未配置 `REDIS_CONN_STRING` 时 `RedisEnabled=false`，所有缓存直接走 DB；并发较高时会显著拉低 `users.quota` 扣减性能，建议挂 Redis。

### Q3: Can SQLite be used in production?

Yes but with caveats: single-writer; growth of the `logs` table slows things down. For production prefer MySQL / PostgreSQL; `LOG_SQL_DSN` may be a separate database.

可以但有限制：单写者；日志表膨胀会拖慢。生产建议 MySQL / PostgreSQL，`LOG_SQL_DSN` 可单独配置。

### Q4: How do I enable cluster mode?

See [Cluster Overview](/en/decentralization/overview). Minimal: 3 nodes, each with its own MySQL / Redis, `CLUSTER_ENABLED=true` and a shared `CLUSTER_SECRET`.

详见 [Cluster 概览](/zh/decentralization/overview)。最简模式：3 节点，每节点独立 MySQL/Redis，`CLUSTER_ENABLED=true` + 共享 `CLUSTER_SECRET`。

Next: [Troubleshooting](/en/faq/troubleshooting) · [Glossary](/en/faq/glossary).