---
title: Option
description: "System-wide KV settings: site name, signup policy, Turnstile, mail signature, etc."
category: schema
order: 14
---

# Option

## What it is

`Option` is a KV store for system-level switches and config — the "runtime parameters" of the system, as opposed to business entities like [User](/en/schema/user) or [Channel](/en/schema/channel).

## Where to find it

- **Admin → System Settings**: site name, logo, signup policy, Turnstile, mail signature
- **Admin → Operation Settings**: announcements, mail templates

## Operator-relevant fields

| Setting | Meaning | Effect |
|---|---|---|
| Site name | Browser title, mail, PDF | Immediate. |
| Logo URL | Header logo | Immediate. |
| Default new-user quota | `QuotaForNewUser` | Initial `User.quota` for newly registered users. |
| Signup switch | Allow new registrations | Off → register page returns 403. |
| Invite reward | Bonus for both inviter and invitee | Applies on the next signup. |
| Turnstile / Captcha | Cloudflare key | Anti-bot signup; leave blank to disable. |
| Mail signature | SMTP mail footer | Applies on the next mail. |
| Terms / Privacy URL | Link in copy | Frontend pages swap in. |

## What does NOT live in Option

- User profile → [User](/en/schema/user)
- Model unit price → [Model Price](/en/schema/model-price)
- Plan → [Plan](/en/schema/plan)

Option is reserved for "no dedicated UI, no dedicated data model, but might need tuning at any time" globals.

## Related pages

- [System Settings (admin)](/en/misc/system-settings)
- [Operation Settings (admin)](/en/misc/operation-settings)

## Related API

- `GET /api/option/` — get all (sensitive values are masked)
- `PUT /api/option/` — update one (Root)