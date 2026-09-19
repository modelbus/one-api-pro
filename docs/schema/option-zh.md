---
title: 系统设置
description: 系统级 KV 设置：站点名、注册策略、Turnstile、邮件签名等。
category: schema
order: 14
---

# 系统设置（Option）

## 这是什么

`Option` 是一个 KV 存储，存的是「系统级开关与配置」。与 [User.quota](/schema/user) / [Channel](/schema/channel) 这种「业务实体」不同，Option 是「运行时参数」。

## 在哪里看到

- **后台 → 系统设置**：通用设置（站点名、Logo、注册策略、Turnstile、邮件签名）
- **后台 → 运营设置**：公告、邮件通知模板

## 关键字段（运维视角）

| 设置项 | 含义 | 设置后影响 |
|---|---|---|
| 站点名 | 显示在浏览器标题、邮件、PDF | 立刻生效 |
| Logo URL | 顶部 logo | 立刻生效 |
| 默认新用户额度 | `QuotaForNewUser` | 新注册用户初始 `quota` |
| 注册开关 | 允许新用户注册 | 关闭后注册页面返回 403 |
| 邀请奖励 | 注册邀请双方各得多少 | 下一次新用户注册生效 |
| Turnstile / Captcha | Cloudflare 验证码 key | 防止机器人注册；留空则禁用 |
| 邮件签名 | SMTP 邮件结尾署名 | 下一次发邮件生效 |
| 服务条款 / 隐私政策 URL | 文案链接 | 前台页面替换 |

## 不在 Option 中存的内容

- 用户资料 → [User](/schema/user)
- 模型单价 → [ModelPrice](/schema/model-price)
- 套餐 → [Plan](/schema/plan)

Option 仅用于「没有专属 UI、没有专属数据模型、但需要随时调整」的全局参数。

## 相关页面

- [系统设置（后台）](/misc/system-settings)
- [运营设置（后台）](/misc/operation-settings)

## 相关 API

- `GET /api/option/` — 获取所有设置（敏感值会被脱敏）
- `PUT /api/option/` — 更新单条（Root）