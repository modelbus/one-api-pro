---
title: Profile
description: Change password, view Access Token, invite link, delete account.
category: user
order: 3
---

# Profile

> Top-right avatar → Personal Center.

## Where

After login, click the avatar in the top-right → **Personal Center**.

## What you can change

By tab:

### Basic

- **Display name**: blank is fine, max 32 chars. Display only.
- **Email**: blank is fine, but to receive password-reset / redemption / order mails, enter a valid one.

### Change password

1. Enter new password twice (must match)
2. Save
3. Current session stays valid; other devices are forced offline

### Access Token

A UUID for calling `/api/*` admin endpoints. See [Access Token](./access-token).

### Invite link

Form: `https://your-host/?aff=<your-code>`. Share it:

- The friend registers via the link → both of you earn bonus quota
- Admin sets the reward amount in [System Settings](../misc/system-settings)

### Third-party binding

Link GitHub / Lark / email to the account. If admin enabled that login method, you can log in via the third party.

### Delete account (dangerous)

"Delete Account" button at the bottom:

1. After double confirmation, your account is disabled (status = "deleted")
2. **You can't log in or call any API**
3. Admin can restore it, but only by acting manually

> After deletion your [orders](../schema/order) and [call logs](../schema/log) remain in the system (for compliance), but the account itself is gone.

## Safety

- Password ≥ 12 chars, mix upper/lower + digits + symbols
- Even with third-party login enabled, set a master password
- Keep Access Token only on trusted machines; never paste it in public issues

## FAQ

- **Changed email but no reset mail arrived**: check the email for typos, also check spam
- **Invite link didn't work**: both inviter and invitee must be active; the new user can't register with an email that's already used
- **Access Token accidentally deleted**: regenerate (the old one stops working immediately)

## Related

- [Access Token](./access-token)
- [API Key](../api/token)
- [User Management (admin)](../user/user-management)