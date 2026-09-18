---
title: 系统设置
description: "站点外观、登录 / 注册策略、OAuth / Turnstile / SMTP / 公告 / 主页内容。"
category: misc
order: 16
---

# 系统设置

> 站点级全局配置：服务器地址、登录与注册开关、GitHub / Lark / 微信 OAuth、Turnstile、SMTP、外观、公告。前端组件：`web/default-pro/src/views/setting/SystemSetting.vue`。

## 接口一览

| Endpoint | Method | 鉴权 | 说明 |
|---|---|---|---|
| `/api/option/` | `GET` | Root | 拉取所有选项（自动屏蔽 `*Token` / `*Secret`） |
| `/api/option/` | `PUT` | Root | 单 key 持久化 |

实现：`controller/option.go`；底层走 `options` 表 + `config.OptionMap`。

## 分区

| 分区 | 字段 |
|---|---|
| Server Address | `ServerAddress`（外部 API base） |
| 登录 / 注册 | `PasswordLoginEnabled` / `PasswordRegisterEnabled` / `RegisterEnabled` / `EmailVerificationEnabled` / `GitHubOAuthEnabled` / `LarkOAuthEnabled` / `WeChatAuthEnabled` / `TurnstileCheckEnabled` |
| GitHub OAuth | `GitHubClientId` / `GitHubClientSecret` |
| Lark OAuth | `LarkClientId` / `LarkClientSecret` |
| 微信扫码登录 | `WeChatServerAddress` / `WeChatServerToken` / `WeChatAccountQRCodeImageURL` |
| Turnstile | `TurnstileSiteKey` / `TurnstileSecretKey` |
| SMTP | `SMTPServer` / `SMTPPort` / `SMTPAccount` / `SMTPFrom` / `SMTPToken` |
| 外观 | `SystemName` / `Logo` / `Theme` |
| 内容 | `Notice` / `HomePageContent` |

## 字段语义

| 字段 | 说明 |
|---|---|
| `ServerAddress` | 用户在「令牌」页面看到的 base URL；通常填 `https://api.example.com` |
| `PasswordLoginEnabled` | 关闭后用户只能走 OAuth|
| `PasswordRegisterEnabled` | 关闭后只能邀请或 OAuth 注册 |
| `RegisterEnabled` | 全局注册开关；false 时 `POST /api/user/register` 直接拒绝 |
| `EmailVerificationEnabled` | 注册时强制要求邮箱 + 验证码 |
| `EmailDomainRestrictionEnabled` + `EmailDomainWhitelist` | 邮箱域名白名单（在 `/setting/operation` 维护） |
| `GitHubOAuthEnabled` / `LarkOAuthEnabled` / `WeChatAuthEnabled` | 启用对应的 OAuth 登录入口；启用前必须填齐 Client/Secret 配置 |
| `TurnstileCheckEnabled` | 注册 / 重置密码等关键路径强制 Cloudflare Turnstile 人机校验 |
| `Notice` | 顶部公告（公开接口 `GET /api/notice` 读取） |
| `HomePageContent` | 首页 hero 文案（公开接口 `GET /api/home_page_content` 读取） |
| `SystemName` / `Logo` / `Theme` | 站点名称|

## 联动校验

`controller/option.go::UpdateOption` 在启用部分开关前会校验前置字段（详见 [operation-setting](operation-setting)）。

## OAuth 登录接入流程

1. 在对应 OAuth 提供方创建应用（GitHub / Lark / 微信开放平台），拿到 Client ID 与 Client Secret。
2. 回调 URL 填写 `{ServerAddress}/api/oauth/{github\|lark\|wechat}`。
3. 把 Client ID / Secret 填入本页面。
4. 启用对应 `*Enabled` 开关；保存即可。

## 邮件（SMTP）

填入 SMTP 服务器地址、端口（默认 587）、账号、发件人地址与 Token（推荐用应用专用密码）。
启用后：
- 注册时的邮箱验证码；
- 密码重置邮件；
- 邀请通知邮件；
都会通过该 SMTP 发出。

## Turnstile 人机校验

1. 在 Cloudflare 控制台创建 Turnstile widget，拿到 Site Key（前端）+ Secret Key（后端）。
2. 把两者填入本页面，启用 `TurnstileCheckEnabled`。
3. 注册 / 重置密码 / 兑换码页面会自动嵌入 Turnstile 校验。

## 前端操作指南

- 顶部：「Server Address」单输入框 + 保存按钮（只改一项）。
- 「登录 / 注册开关」8 个开关横排，change 立即生效。
- 各 OAuth / Turnstile / SMTP 区段：填好字段后点击对应区段的「保存」按钮（独立 submit）。
- 「外观」3 字段一起保存。
- 「内容」2 个 textarea 一起保存。

## 接口实现

| 关注点 | 位置 |
|---|---|
| Handler + 校验 | `controller/option.go` |
| 启动时初始化 / 后台同步 | `model/option.go::InitOptionMap` / `SyncOptions` |
| 公开公告 / 主页 | `controller/misc.go` (`GetNotice` / `GetHomePageContent`) |
| GitHub OAuth | `controller/auth/github.go` |
| Lark OAuth | `controller/auth/lark.go` |
| 微信 OAuth | `controller/auth/wechat.go` |
| 路由 | `router/api.go` |
