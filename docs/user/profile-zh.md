---
title: 个人资料
description: "修改密码、OAuth 绑定、通知偏好。"
category: user
order: 3
---

# 个人资料

> 修改密码、OAuth 绑定、通知偏好。
> Change password, OAuth bindings, notification preferences.

入口：顶栏右上角头像 → 个人中心，路径 `/setting/personal`。代码 `web/default-pro/src/views/setting/PersonalSetting.vue`。

Entry: top-right avatar → "Personal Center", path `/setting/personal`. Source: `web/default-pro/src/views/setting/PersonalSetting.vue`.

## 功能区 / Sections

| 区 / Section | 字段 / Field | API | 备注 / Notes |
| --- | --- | --- | --- |
| **资料** | 显示名 | `PUT /api/user/self` | 显示名可空；最多 32 字符 |
| **修改密码** | 新密码 + 确认密码 | `PUT /api/user/self` | 二次输入必须一致；留空表示不改 |
| **Access Token** | 生成 / 显示 / 复制 | `GET /api/user/token` | UUID 格式，只显示一次 |
| **邀请链接** | `?aff=XXX` 链接 | `GET /api/user/aff` | 邀请人获得邀请奖励 |
| **第三方绑定** | GitHub / Lark / 邮箱 | `GET /api/oauth/github/bind`、`/api/oauth/lark/bind`、`GET /api/oauth/email/bind` | 见下 |
| **危险区** | 注销账号 | `DELETE /api/user/self` | 二次确认；状态置为 `UserStatusDeleted (3)` |

> API Token（`sk-`）与 Access Token（UUID）不同：API Token 用于 `/v1/*` 调用，Access Token 用于 `/api/*` 管理接口。
> "API Token" (`sk-…`) and "Access Token" (UUID) are different — the former calls `/v1/*`, the latter calls `/api/*`.

## 字段校验 / Validation

- **显示名** (`display_name`)：1–32 字符，前后 `strings.TrimSpace`。
  Display name: 1–32 chars, trimmed.
- **新密码**：6–32 字符，必须含字母 + 数字（前后端各校验一次）。
  New password: 6–32 chars, must contain both letters and digits (validated both sides).
- **确认密码**：必须与新密码一致；前端用 `===` 即时校验。
  Confirmation: must equal the new password; checked instantly with `===`.
- **邀请链接**：邀请人只能为自己生成一个 `aff` code；多次调用返回同一个。
  Invite link: one `aff` code per inviter, returned by repeated calls.

## 第三方绑定流程 / Third-party Binding Flows

```text
GitHub:  window.location.href = `https://github.com/login/oauth/authorize?client_id=${id}&scope=user:email`
         └─► 回调 /oauth/github → /api/oauth/github/callback → users.oauth_provider='github'
Lark  :  window.location.href = `https://open.feishu.cn/open-apis/authen/v1/authorize?app_id=${id}&redirect_uri=…`
         └─► 回调 /oauth/lark → /api/oauth/lark/callback → users.oauth_provider='lark'
Email :  弹窗 input email + 验证码
         └─► GET /api/verification?email=… 发送
         └─► GET /api/oauth/email/bind?email=…&code=… 绑定
```

只有当 `status.github_client_id` / `status.lark_client_id` 非空时，对应按钮才显示；绑定的具体开关在「系统设置」中。

The corresponding buttons render only when `status.github_client_id` / `status.lark_client_id` is non-empty; the toggle itself lives under "System Settings".

## 注销账号 / Account Deletion

- **不可恢复**：`DELETE /api/user/self` 把 `users.status` 置为 `UserStatusDeleted (3)`，从 `GET /api/user/search` 等列表接口中过滤掉；订单、订阅、日志保留以满足审计。
  **Irreversible**: `users.status` is set to `UserStatusDeleted (3)`. The user is filtered out of list queries; orders, subscriptions and logs are kept for audit.
- **退登**：删除成功后 `useAuthStore.logout()` 触发，跳回 `/`。
  Auto-logout: `useAuthStore.logout()` runs and redirects to `/`.
- **管理员降级**：管理员不能删除自己（`auth_helper` 中 `RoleRootUser` 不允许自删），需要另一个 Root 协助。
  Admin safety: a Root cannot delete themselves; another Root must assist.

## 密码策略 / Password Policy

```text
- 长度 6–32
- 必须包含字母（a-zA-Z）
- 必须包含数字（0-9）
- 不允许纯空格
- 与旧密码相同 → 拒绝
```

后端校验失败返回 `{ success:false, message }`；前端在 `Message.error` 中展示。

Backend rejects with `{ success:false, message }`; frontend surfaces it via `Message.error`.

## i18n 键值 / i18n Keys

| Key | zh | en |
| --- | --- | --- |
| `settingPage.personal.profile` | 资料 | Profile |
| `settingPage.personal.displayName` | 显示名 | Display name |
| `settingPage.personal.newPassword` | 新密码 | New password |
| `settingPage.personal.confirmPassword` | 确认密码 | Confirm password |
| `settingPage.personal.saveChanges` | 保存 | Save |
| `settingPage.personal.accessToken` | Access Token | Access Token |
| `settingPage.personal.affLink` | 邀请链接 | Invite link |
| `settingPage.personal.thirdPartyBinding` | 第三方绑定 | Third-party binding |
| `settingPage.personal.dangerZone` | 危险区 | Danger zone |
| `settingPage.personal.deleteAccount` | 注销账号 | Delete account |

新加字段 / 区块时，记得同步更新 `src/i18n/pages/setting.js` 的 `personal` 命名空间。

When adding fields / sections, remember to extend the `personal` namespace in `src/i18n/pages/setting.js`.

下一步 / Next: [Access Token](/zh/user/access-token) · [Chat Playground](/zh/user/chat)。