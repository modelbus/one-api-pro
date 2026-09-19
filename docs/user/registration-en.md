---
title: Registration & Login
description: Email, username/password, GitHub, OIDC, Lark, WeChat login options.
category: user
order: 1
---

# Registration & Login

> How to create a new account, what login methods exist, and how to recover your password.

## Login methods

The `/login` page shows only the methods the admin has enabled:

| Method | When | How |
|---|---|---|
| **Username + password** | Always | Default admin: `root` / `123456` (must change on first login) |
| **Email registration** | When email verification is enabled | `/register`, fill email, get a code |
| **GitHub** | When GitHub OAuth is configured | Click the GitHub icon |
| **Lark** | When Lark OAuth is configured | Click the Lark icon |
| **WeChat scan** | When WeChat login is enabled | Scan the QR code |
| **OIDC / Keycloak** | When OIDC is configured | Click the button at the bottom |

Methods the admin hasn't enabled simply don't appear.

## Registration

1. Visit `/register`
2. Fill username, password (twice), optional email, optional invite code
3. **Must check** "I agree to the terms" to submit
4. If email verification is enabled, enter the code from your email
5. Submit → on success you're taken to `/login`

### Invite code

When registering via a link like `https://your-host/?aff=XXXX`, the invite code auto-fills. After successful registration, both you and the inviter earn bonus quota.

### Default group

New accounts default to the `default` group. Admins can change this in [Operation Settings](../misc/operation-settings).

## Login flow

1. `/login` — username + password (or third-party)
2. On success → `/dashboard`
3. On failure → red banner "Wrong username or password"

Repeated failures trigger rate limiting to prevent brute force.

## Password recovery

1. `/reset` — enter email
2. Backend emails a verification code
3. `/reset/:token` — enter new password
4. Done — log in with the new password

> The verification code is logged to stdout in dev. In production, configure SMTP so codes go by email.

## Multi-instance notes

For Cluster multi-node deployments:

- All nodes must share the same `SESSION_SECRET` (otherwise sessions don't share across nodes).
- Verification codes are stored in-process by default — not shared across nodes. In production, use an external SMS / email service.

## FAQ

- **GitHub login fails**: admin needs to set `GITHUB_CLIENT_ID` / `GITHUB_CLIENT_SECRET`; the OAuth callback URL must match the One API Pro origin.
- **Lark login fails**: in the Lark console, add One API Pro's domain to "Security settings → Redirect URLs".
- **No verification email**: check spam; verify admin set up SMTP.
- **`/register` returns 403**: admin disabled registration (`REGISTER_ENABLED=false`).

## Related

- [Dashboard](./dashboard)
- [Access Token](./access-token)
- [Profile](./profile)
- [System Settings (admin)](../misc/system-settings)