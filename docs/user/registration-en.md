---
title: Registration & Login
description: "Email, GitHub, OIDC, Lark and WeChat login flows."
category: user
order: 1
---

# Registration & Login

> Email, GitHub, OIDC, Lark and WeChat login flows.
> 邮箱、GitHub、OIDC、飞书、微信登录方式。

Entry points: `/register`, `/login`, `/reset`, `/oauth/github`, `/oauth/lark`. Source: `web/default-pro/src/views/auth/`.

入口：`/register`、`/login`、`/reset`、`/oauth/github`、`/oauth/lark`。代码位于 `web/default-pro/src/views/auth/`。

## Sign-in Methods / 登录方式

| Method / 方式 | Prerequisite / 启用条件 | Entry / 入口 | Notes / 备注 |
| --- | --- | --- | --- |
| **Username + password** | Always on / 始终可用 | `/login` | Default root: `root` / `123456` — must change on first login |
| **Email registration** | `EMAIL_VERIFICATION=true` requires a code | `/register` | Codes live in `common/verification.go` (in-process Map; **not shared across cluster nodes** — wire to an external service for multi-node) |
| **GitHub OAuth** | `GITHUB_CLIENT_ID` + `GITHUB_CLIENT_SECRET` non-empty | `/oauth/github` | Callback URL: `{origin}/oauth/github`; first-time logins need matching email |
| **Lark OAuth** | `LARK_CLIENT_ID` + `LARK_CLIENT_SECRET` non-empty | `/oauth/lark` | Lark "Web app" type |
| **WeChat QR-code login** | `WECHAT_AUTH_ENABLED=true` + official-account / web-app credentials | QR on login page | Implemented in `controllers/oauth/wechat.go` |
| **Generic OIDC** | `OIDC_*` env vars | Button at the bottom of register/login | Works with Keycloak / Auth0 / Authing, etc. |

> These flags surface through `useStatusStore().status` (`status.github_client_id`, `status.lark_client_id`, `status.email_verification`, `status.oidc`); backend middleware still enforces — don't rely on frontend hiding only.
> 这些开关在前端通过 `useStatusStore().status` 暴露，控制 UI 显隐。后端中间件兜底，不要只靠前端隐藏。

## Registration / 注册流程

```text
1. User visits /register
2. Fill username / password / password2 / optional email / optional invite code / accept terms
3. If EMAIL_VERIFICATION=true, email + verification code are required
4. POST /api/user/register (form payload)
5. Server validates password strength, username/email uniqueness, invite code, email code
6. On success → redirect to /login
```

**Constraints / 关键约束**:

- **Terms checkbox required**: `<a-checkbox v-model="agreedToTerms">` blocks submission.
  协议必勾选：未勾选禁止提交。
- **Invite code**: `?aff=XXX` pre-fills the field; the inviter is rewarded on success (see `model/user.go`).
  邀请码：URL 带 `?aff=XXX` 时自动填入，注册成功后邀请人获得奖励。
- **Disable registration**: setting `REGISTER_ENABLED=false` makes `/register` return 403.
  关闭注册：设置系统选项 `REGISTER_ENABLED=false` 后 `/register` 路由返回 403。
- **Default group**: new users land in `default`; configurable via the operations setting `user.default_group`.
  默认用户组：注册得到的 `users.group` 默认 `default`，可在「运营设置」改默认组。

## Sign-in / 登录流程

```text
1. POST /api/user/login → Cookie Session is set
2. Frontend useAuthStore.login() → persists localStorage.user
3. Redirect to /dashboard
```

- Login failures: `{ success:false, message }`, surfaced via `Message.error`.
  登录失败：`{ success:false, message }`，前端通过 `Message.error` 提示。
- Repeated failures are throttled by `GLOBAL_WEB_RATE_LIMIT` (see `/install/config#session--rate-limits-ratelimits`).
  多次失败会被 `GLOBAL_WEB_RATE_LIMIT` 兜底限流。
- Session key is `SESSION_SECRET`; multi-instance deployments must share it.
  Session 加密密钥见 `SESSION_SECRET`；多实例必须一致，否则会话不互通。

## Password Reset / 找回密码

```text
1. /reset — enter email
2. POST /api/password/reset → backend stores in `common/verification.go` with purpose "r"
3. Code is sent via email / printed to console in dev
4. /reset/:token — enter new password and POST /api/password/reset_confirm
5. Backend validates the code and rehashes via bcrypt
```

> Email sending currently logs to `logger.SysLog`; integrate SMTP / a third-party provider in production.
> 当前实现的邮件发送以**控制台日志**形式输出（`logger.SysLog`），生产环境请接 SMTP / 第三方邮件。

## Third-party Binding / 第三方绑定

In "Personal Center → Third-party Binding" you can bind email / GitHub / Lark independently of sign-in OAuth:

在「个人中心 → 第三方绑定」可绑定邮箱 / GitHub / 飞书（独立于登录用 OAuth）：

- **GitHub**: redirect to `https://github.com/login/oauth/authorize?client_id=...&scope=user:email`.
  GitHub：跳转 `https://github.com/login/oauth/authorize?client_id=...&scope=user:email`。
- **Lark**: redirect to Lark's `authen/v1/authorize`, callback at `/oauth/lark`.
  飞书：跳转飞书 `authen/v1/authorize`，回调 `/oauth/lark`。
- **Email**: modal asking for email + verification code.
  邮箱：弹窗输入邮箱 + 验证码。

After binding, `users.oauth_provider` / `users.oauth_id` are updated; the bound provider can be used for next login.

绑定后 `users.oauth_provider` / `users.oauth_id` 字段更新；下次可用对应方式登录。

## Pitfalls / 常见踩坑

- The OAuth callback URL must share the same origin as `FRONTEND_BASE_URL`; otherwise the cookie won't be set.
  OAuth 回调 URL 必须与 `FRONTEND_BASE_URL` 同源，否则 cookie 写不进去。
- Multi-instance deployments must share `SESSION_SECRET`; otherwise users bounce between nodes.
  多实例部署必须保证 `SESSION_SECRET` 一致，否则用户在 A 节点登录后访问 B 节点会被踢回登录。
- Add the Lark callback URL to the app's "Security → Redirect URLs" allow-list.
  飞书回调域名需在飞书后台「安全设置 → 重定向 URL」里加白。
- GitHub OAuth users with no public email need `scope=user:email`; otherwise the lookup errors out.
  GitHub OAuth 用户无公开邮箱：回调 `scope=user:email` 必带；否则拿不到邮箱，会落到 401。

Next: [Dashboard](/en/user/dashboard) · [Access Token](/en/user/access-token) · [Profile](/en/user/profile).