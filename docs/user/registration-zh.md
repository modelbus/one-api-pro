---
title: 注册与登录
description: "邮箱、GitHub、OIDC、飞书、微信登录方式。"
category: user
order: 1
---

# 注册与登录

> 邮箱、GitHub、OIDC、飞书、微信登录方式。
> Email, GitHub, OIDC, Lark and WeChat login flows.

入口：`/register`、`/login`、`/reset`、`/oauth/github`、`/oauth/lark`。代码位于 `web/default-pro/src/views/auth/`。

Entry points: `/register`, `/login`, `/reset`, `/oauth/github`, `/oauth/lark`. Source: `web/default-pro/src/views/auth/`.

## 登录方式 / Sign-in Methods

| 方式 / Method | 启用条件 / Prerequisite | 入口 / Entry | 备注 / Notes |
| --- | --- | --- | --- |
| **用户名 + 密码** | 始终可用 / always on | `/login` | 默认 root：`root` / `123456`，首次登录后必须改密 |
| **邮箱注册** | 启用 `EMAIL_VERIFICATION=true` 后注册需邮件验证码 | `/register` | 验证码走 `common/verification.go`，进程内 Map，**集群多节点下不互通**（建议接外部 SMS / 邮件服务） |
| **GitHub OAuth** | `GITHUB_CLIENT_ID` + `GITHUB_CLIENT_SECRET` 非空 | `/oauth/github` | 回调 URL：`{origin}/oauth/github`；首次会要求邮箱匹配 |
| **飞书 OAuth** | `LARK_CLIENT_ID` + `LARK_CLIENT_SECRET` 非空 | `/oauth/lark` | 飞书「网页登录」应用 |
| **微信扫码登录** | `WECHAT_AUTH_ENABLED=true` + 公众号 / 网站应用配置 | 登录页二维码 | 走 `controllers/oauth/wechat.go` |
| **OIDC 通用登录** | `OIDC_*` 系列环境变量 | 注册 / 登录页底部按钮 | 兼容 Keycloak / Auth0 / Authing 等 |

> 这些开关在前端通过 `useStatusStore().status` 暴露，控制 UI 显隐（`status.github_client_id`、`status.lark_client_id`、`status.email_verification`、`status.oidc`）。后端中间件兜底，不要只靠前端隐藏。
> These flags surface through `useStatusStore().status`; backend middleware still enforces — don't rely on frontend-only hiding.

## 注册流程 / Registration

```text
1. 用户访问 /register
   User visits /register
2. 填写 username / password / password2 / email（可选）+ invite code（可选）+ 同意协议
   Fill username / password / password2 / optional email / optional invite code / accept terms
3. 若 EMAIL_VERIFICATION=true，邮箱必填 + 验证码
   If EMAIL_VERIFICATION=true, email + verification code required
4. POST /api/user/register（form 表单）
   POST /api/user/register
5. 服务端校验：密码强度、用户名唯一性、邮箱唯一性、邀请码有效性、邮箱验证码
   Server validates password strength, username/email uniqueness, invite code, email code
6. 注册成功 → 跳转 /login
   On success → jump to /login
```

**关键约束 / Constraints**：

- **协议必勾选**：注册页 `<a-checkbox v-model="agreedToTerms">`，未勾选禁止提交。
  Terms checkbox is enforced; submission is blocked when unchecked.
- **邀请码**：URL 带 `?aff=XXX` 时自动填入，注册成功后邀请人获得奖励（`aff` 字段见 `model/user.go`）。
  Invite code from `?aff=XXX` is pre-filled; inviter gets rewarded upon success (`aff` field, see `model/user.go`).
- **关闭注册**：设置系统选项 `REGISTER_ENABLED=false` 后 `/register` 路由返回 403。
  Setting `REGISTER_ENABLED=false` makes `/register` return 403.
- **默认用户组**：注册得到的 `users.group` 默认 `default`，可在「运营设置」改默认组（`user.default_group`）。
  Default group is `default`; configurable via the operations setting `user.default_group`.

## 登录流程 / Sign-in

```text
1. POST /api/user/login → Cookie Session 写入
   Login → set Cookie Session
2. 前端 useAuthStore.login() → 写入 localStorage.user
   useAuthStore.login() persists localStorage.user
3. 跳转 /dashboard
   Redirect to /dashboard
```

- 登录失败：`{ success:false, message }`，前端通过 `Message.error` 提示。
  Login failures: `{ success:false, message }`, surfaced via `Message.error`.
- 多次失败会被 `GLOBAL_WEB_RATE_LIMIT` 兜底限流（`/install/config#会话与限流-ratelimits`）。
  Repeated failures are throttled by `GLOBAL_WEB_RATE_LIMIT`.
- Session 加密密钥见 `SESSION_SECRET`；多实例必须一致，否则会话不互通。
  Session key is `SESSION_SECRET`; multi-instance deployments must share it.

## 找回密码 / Password Reset

```text
1. /reset 输入邮箱
   /reset — enter email
2. POST /api/password/reset → 后端写入 `common/verification.go`（purpose="r"）
   Backend stores code via `common/verification.go` with purpose "r"
3. 邮件 / 控制台输出验证码（开发期走 stdout）
   Send code via email / console in dev
4. /reset/:token 输入新密码 → POST /api/password/reset_confirm
   /reset/:token — enter new password and POST /api/password/reset_confirm
5. 后端校验验证码 + 写入新密码（bcrypt）
   Backend validates code and rehashes via bcrypt
```

> 当前实现的邮件发送以**控制台日志**形式输出（`logger.SysLog`），生产环境请接 SMTP / 第三方邮件。
> Current email sending is logged via `logger.SysLog`; integrate SMTP / a third-party provider in production.

## 第三方绑定 / Third-party Binding

在「个人中心 → 第三方绑定」可绑定邮箱 / GitHub / 飞书（独立于登录用 OAuth）：

In "Personal Center → Third-party Binding" you can bind email / GitHub / Lark (independent of sign-in OAuth):

- **绑定 GitHub**：跳转 `https://github.com/login/oauth/authorize?client_id=...&scope=user:email`；
  Bind GitHub: redirect to `https://github.com/login/oauth/authorize?client_id=...&scope=user:email`.
- **绑定飞书**：跳转飞书 `authen/v1/authorize`，回调 `/oauth/lark`；
  Bind Lark: redirect to Lark's `authen/v1/authorize`, callback at `/oauth/lark`.
- **绑定邮箱**：弹窗输入邮箱 + 验证码。
  Bind email: modal asking for email + verification code.

绑定后 `users.oauth_provider` / `users.oauth_id` 字段更新；下次可用对应方式登录。

After binding, `users.oauth_provider` / `users.oauth_id` are updated; the bound provider can be used for next login.

## 常见踩坑 / Pitfalls

- **OAuth 回调 URL 必须与 `FRONTEND_BASE_URL` 同源**，否则 cookie 写不进去。
  The OAuth callback URL must share the same origin as `FRONTEND_BASE_URL`; otherwise the cookie won't be set.
- **多实例部署必须保证 `SESSION_SECRET` 一致**，否则用户在 A 节点登录后访问 B 节点会被踢回登录。
  Multi-instance deployments must share `SESSION_SECRET`; otherwise users bounce between nodes.
- **飞书回调域名需在飞书后台「安全设置 → 重定向 URL」里加白**。
  Add the Lark callback URL to the app's "Security → Redirect URLs" allow-list.
- **GitHub OAuth 用户无公开邮箱**：回调 `scope=user:email` 必带；否则拿不到邮箱，会落到 401。
  GitHub OAuth users with no public email need `scope=user:email`, otherwise the lookup errors out.

下一步 / Next: [个人仪表盘](/zh/user/dashboard) · [Access Token](/zh/user/access-token) · [个人资料](/zh/user/profile)。