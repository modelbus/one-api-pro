---
title: User
description: "One API Pro account: identity, quota, group and role."
category: schema
order: 2
---

# User

## What it is

`User` is the system's only identity. Every login, every API call, every order is attached to one user.

## Where to find it

- **Admin → Users**: list, detail, adjust balance, enable/disable.
- **Personal Center**: the current account's email, username and group.
- **Call logs**: every `/v1/*` call records the issuing `user_id`.

## Operator-relevant fields

These are the fields that actually change behavior. Most issues only need one of these touched:

| Field | Meaning | Effect of editing |
|---|---|---|
| Username | Login account | Next login uses the new username. |
| Email | Password reset / notifications | Re-verification required. |
| Password | Login secret | User is forced offline. |
| Group | `default` / `vip` / `svip` | Drives [Group Price](/en/schema/group-price) lookup and allowed-model scope. |
| Status | Enabled / Disabled | Disabled users cannot call any API or log in. |
| Quota `quota` | Balance (internal units) | First source billed; 0 rejects calls. |
| Used quota | Cumulative spent | Recorded only, no behavior. |
| Role | User / Admin / Root | Root-only access to the admin console. |
| Invite code / invitee | Invite relationship | Both sides earn bonus quota at signup. |

Fields you usually should not touch: registration time, last login time — managed by the system.

## Related entities

- A [Token](/en/schema/token) is auto-created when a user is added.
- All [Order](/en/schema/order) / [Topup](/en/schema/topup) / [Subscription](/en/schema/subscription) hang off `user_id`.
- Deduction order: `User.quota` → [Token quota](/en/schema/token) → [Subscription plan quota](/en/schema/subscription). See [Billing Rules](/en/subscription/billing-rules).

## Related pages

- [Registration & Login](/en/user/registration) — how to sign up, log in, recover password
- [Profile](/en/user/profile) — what the user can edit
- [User Management (Admin)](/en/user/user-management) — admin-side CRUD

## Related API

- `GET /api/user/self` — current user info
- `POST /api/user` — admin create
- `PUT /api/user` — admin update
- `POST /api/user/manage` — adjust balance / status