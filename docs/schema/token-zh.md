---
title: 访问令牌
description: 用户在管理后台生成的系统级 Token，用于调用 /api/* 接口。
category: schema
order: 3
---

# 访问令牌（Token）

## 这是什么

`Token` 是用户在 **前台个人中心** 生成的系统级访问令牌，对应 OpenAI 兼容 Key 之外提供的「管理 API Key」（UUID 格式）。用于调用 `/api/*` 管理后台类接口。

> 与 [渠道 API Key](/schema/channel) 是两个完全不同的东西：
>
> - 用户 Token → 调用 One API Pro 的 `/api/*`
> - 渠道 API Key → 由 One API Pro 调用上游 Provider 时使用

## 在哪里看到

- **前台 → 个人中心 → Access Token**：查看、复制、重置
- **后台 → API Token 管理**：列出所有用户的 Token、调整额度、禁用、删除
- **调用日志**：每条 `/api/*` 调用都记录了使用的 `token_id`

## 关键字段（运维视角）

| 字段 | 含义 | 设置后影响 |
|---|---|---|
| Key | UUID 字符串 | 修改后旧 Key 立即失效 |
| 名称 | 备注用 | 不影响行为 |
| 所属用户 | 关联到 User | 删除用户会级联删除 Token |
| 状态 | 启用 / 禁用 / 过期 / 额度耗尽 | 禁用后用此 Key 调用 API 直接被拒绝 |
| 剩余额度 `remain_quota` | Token 独立额度 | 扣到 0 即拒绝；不影响用户 `quota` |
| 额度上限 `quota` | 本次发放的额度 | 与 `remain_quota` 不同，仅记录初始额度 |
| 过期时间 | Unix 秒 / `-1` 表示永不过期 | 过期后状态变为「已过期」 |
| 允许的模型 `models` | 空数组表示不限制 | 只允许调用列出的模型；非空时其他模型直接 403 |
| 允许的子网 `subnet` | CIDR 列表 | 客户端 IP 不在白名单则 403 |

## 与「用户额度」的关系

扣费顺序：

1. 优先扣 `Token.remain_quota`
2. 耗尽后再扣 [User.quota](/schema/user)
3. 仍有 [Subscription](/schema/subscription) 时按套餐折扣扣减

所以一个用户可以同时拥有：

- 个人账户余额
- 多个 Token，每个有独立额度
- 一个活跃订阅

这便于做项目组、团队子账号等场景。

## 相关 API

- `GET /api/user/token` — 获取当前用户的 Token
- `POST /api/token` — 创建 Token
- `PUT /api/token` — 更新 Token（调整额度、过期时间、模型白名单）
- `DELETE /api/token/:id` — 删除 Token