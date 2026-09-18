---
title: 常见问题
description: "用户与管理员的常见问答。"
category: misc
order: 1
---

# 常见问题

> 用户与管理员的常见问答。

## 用户

### Q1：我的 API Key 在哪里？

`/token` 页新建的 Key 是 **OpenAI 兼容 Key**（`sk-` 前缀），用于 `/v1/*` 调用。UUID 格式的 **Access Token** 在个人中心生成，用于 `/api/*` 管理接口。两者不可混用。

### Q2：为什么账户额度一直不变？

管理员可能只给了 Token 额度（`tokens.remain_quota`），没有给用户额度（`users.quota`）。调用会优先扣 Token 的 `remain_quota`；Token 用尽后再扣用户。详见 [架构总览 · 缓存一致性](/zh/start/architecture#请求生命周期-request-lifecycle)。

### Q3：套餐过期后还能继续用吗？

可以继续调用到令牌耗尽，但不再享受套餐折扣倍率；且 `daily_quota` 进度条变为 0。

### Q4：邀请链接没看到奖励？

`/api/user/aff` 返回的 URL 即邀请链接；被邀请人成功注册后会调用 `IncreaseUserQuota` 增加奖励。v0.0.21 后此链路已同步刷 Redis 缓存，**下个请求立即可见**。

### Q5：重置密码邮件没收到？

当前邮件发送以 `logger.SysLog` 形式输出，**生产环境请接 SMTP / 第三方**（详见 [个人资料 · 找回密码](/zh/user/profile#找回密码--password-reset)）。

## 管理员

### Q1：渠道测试 200 但 `/v1/*` 调用 401？

通常是 **Base URL / Key 不匹配**，或上游账号失效、地区限制。检查渠道「测试」页响应头里的 `x-request-id` 回到上游工单。

### Q2：能不能让某个用户组只能用特定模型？

设置 Token 时填写「可用模型」白名单，或在渠道侧配置 `models` 限制（按模型 + 用户组双维度）。中间件 `relay/middleware/should-check-model` 强制校验。

### Q3：升级差价如何配置？

系统选项 `plan.upgrade_mode` 取值：

| 值| 含义|
| --- | --- |
| `price_diff`（默认） | 按差价下单 `UP` 前缀订单 |
| `stack` | 全额叠加，旧订阅继续生效 |

### Q4：充值额度与套餐额度谁优先？

**套餐优先**：调用扣减顺序 `token.remain_quota` → `users.quota`（充值）→ 不够则 403。套餐有 `daily_quota` 时，按日限额 + 套餐倍率再扣 `users.quota`。

### Q5：管理员怎么给别人加额度？

管理后台「用户详情」→ 「调整额度」，或 `POST /api/user/<id>/quota` 直接调；这两条链路在 v0.0.21 后都已同步刷 Redis，下一次请求立即生效。

## 部署

### Q1：Docker 启动后 3000 端口被占用？

修改 `-p 8080:3000`；同时若用 Nginx 反代，`proxy_pass http://127.0.0.1:8080;`。

### Q2：能不能不挂 Redis？

可以，未配置 `REDIS_CONN_STRING` 时 `RedisEnabled=false`，所有缓存直接走 DB；并发较高时会显著拉低 `users.quota` 扣减性能，建议挂 Redis。

### Q3：能不能用 SQLite 跑生产？

可以但有限制：单写者；日志表膨胀会拖慢。生产建议 MySQL / PostgreSQL，`LOG_SQL_DSN` 可单独配置。

### Q4：Cluster 模式如何接入？

详见 [Cluster 概览](/zh/decentralization/overview)。最简模式：3 节点，每节点独立 MySQL/Redis，`CLUSTER_ENABLED=true` + 共享 `CLUSTER_SECRET`。

下一步 / Next: [故障排查](/zh/faq/troubleshooting) · [术语表](/zh/faq/glossary)。
