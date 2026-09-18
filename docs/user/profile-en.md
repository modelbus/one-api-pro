---
title: Profile
description: "Change password, OAuth bindings, notification preferences."
category: user
order: 3
---

# Profile

> Change password, OAuth bindings, notification preferences.
> 修改密码、OAuth 绑定、通知偏好。

Entry: top-right avatar → "Personal Center", path `/setting/personal`. Source: `web/default-pro/src/views/setting/PersonalSetting.vue`.

入口：顶栏右上角头像 → 个人中心，路径 `/setting/personal`。代码 `web/default-pro/src/views/setting/PersonalSetting.vue`。

## Sections / 功能区

| Section / 区 | Field / 字段 | API | Notes / 备注 |
| --- | --- | --- | --- |
| **Profile** | Display name | `PUT /api/user/self` | Up to 32 chars, trim whitespace |
| **Change password** | New + confirm | `PUT /api/user/self` | Empty = no change; confirmation must match |
| **Access Token** | Generate / show / copy | `GET /api/user/token` | UUID format; shown only once |
| **Invite link** | `?aff=XXX` URL | `GET /api/user/aff` | Inviter earns on each signup |
| **Third-party binding** | GitHub / Lark / Email | `GET /api/oauth/github/bind`, `/api/oauth/lark/bind`, `/api/oauth/email/bind` | See below |
| **Danger zone** | Delete account | `DELETE /api/user/self` | Second confirmation; sets `UserStatusDeleted (3)` |

> "API Token" (`sk-…`) and "Access Token" (UUID) are different — the former calls `/v1/*`, the latter calls `/api/*`.
> API Token（`sk-`）与 Access Token（UUID）不同：API Token 用于 `/v1/*` 调用，Access Token 用于 `/api/*` 管理接口。

## Validation / 字段校验

- **Display name** (`display_name`): 1–32 chars, trimmed via `strings.TrimSpace`.
- **New password**: 6–32 chars, must contain both letters and digits (validated both sides).
- **Confirm password**: must equal new password; checked instantly with `===` on the frontend.
- **Invite link**: one `aff` code per inviter; repeated calls return the same code.

## Third-party Binding Flows / 第三方绑定流程

```text
GitHub:  window.location.href = `https://github.com/login/oauth/authorize?client_id=${id}&scope=user:email`
         └─► callback /oauth/github → /api/oauth/github/callback → users.oauth_provider='github'
Lark  :  window.location.href = `https://open.feishu.cn/open-apis/authen/v1/authorize?app_id=${id}&redirect_uri=…`
         └─► callback /oauth/lark → /api/oauth/lark/callback → users.oauth_provider='lark'
Email :  modal: enter email + verification code
         └─► GET /api/verification?email=…  send code
         └─► GET /api/oauth/email/bind?email=…&code=…  bind
```

Buttons render only when `status.github_client_id` / `status.lark_client_id` is non-empty; the toggle itself lives under "System Settings".

只有当 `status.github_client_id` / `status.lark_client_id` 非空时，对应按钮才显示；绑定的具体开关在「系统设置」中。

## Account Deletion / 注销账号

- **Irreversible**: `users.status` is set to `UserStatusDeleted (3)`. The user is filtered out of list queries; orders, subscriptions and logs are kept for audit.
  **不可恢复**：`users.status` 置为 `UserStatusDeleted (3)`，订单、订阅、日志保留以满足审计。
- **Auto-logout**: `useAuthStore.logout()` runs and redirects to `/`.
  退登：`useAuthStore.logout()` 触发，跳回 `/`。
- **Admin safety**: a Root cannot delete themselves; another Root must assist.
  管理员不能删除自己，需要另一个 Root 协助。

## Password Policy / 密码策略

```text
- 6–32 chars
- Must contain a letter (a-zA-Z)
- Must contain a digit (0-9)
- No whitespace-only passwords
- Same as the current password → reject
```

Backend rejects with `{ success:false, message }`; frontend surfaces it via `Message.error`.

## i18n Keys / i18n 键值

| Key | en | zh |
| --- | --- | --- |
| `settingPage.personal.profile` | Profile | 资料 |
| `settingPage.personal.displayName` | Display name | 显示名 |
| `settingPage.personal.newPassword` | New password | 新密码 |
| `settingPage.personal.confirmPassword` | Confirm password | 确认密码 |
| `settingPage.personal.saveChanges` | Save | 保存 |
| `settingPage.personal.accessToken` | Access Token | Access Token |
| `settingPage.personal.affLink` | Invite link | 邀请链接 |
| `settingPage.personal.thirdPartyBinding` | Third-party binding | 第三方绑定 |
| `settingPage.personal.dangerZone` | Danger zone | 危险区 |
| `settingPage.personal.deleteAccount` | Delete account | 注销账号 |

When adding fields / sections, remember to extend the `personal` namespace in `src/i18n/pages/setting.js`.

新加字段 / 区块时，记得同步更新 `src/i18n/pages/setting.js` 的 `personal` 命名空间。

Next: [Access Token](/en/user/access-token) · [Chat Playground](/en/user/chat).