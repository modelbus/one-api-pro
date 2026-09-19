---
title: Redemption
description: Bulk-issued quota or plan redemption codes.
category: schema
order: 11
---

# Redemption

## What it is

`Redemption` is a **bulk-issued redeemable code**. Admins generate codes, end users redeem them at the public "Redemption Center" to receive quota or activate a plan.

## Where to find it

- **Admin → Redemption Codes**: list, CSV / TXT export
- **Public → Redemption**: user enters code
- **Admin → User detail**: per-user redemption history

## What you set per code

| Field | Where | Effect |
|---|---|---|
| Type | Admin → Add | "Quota" adds to `User.quota`; "Plan" creates a Subscription. |
| Quota amount | Admin → Add | Credited on redeem (quota type only). |
| Plan | Admin → Add | Activated on redeem (plan type only). |
| Expiry | Admin → Add | After expiry, code is unusable. |
| Max redemptions | Admin → Add | Defaults to 1. |

## Operator-relevant fields (per code)

| Field | Meaning | Effect |
|---|---|---|
| Code | The string | Changing it after issuance breaks distribution. |
| Type | Quota / Plan | Drives the redeem action. |
| Quota amount | Number | Credited (quota type only). |
| Plan FK | Plan | Activated (plan type only). |
| Total redemptions | Default 1 | Larger numbers mean "multi-redeem code". |
| Redeemed | Number | Increments automatically; maxed-out codes stop working. |
| Expire time | Unix seconds / 0 = never | After expiry, code is unusable. |
| Status | Unused / Exhausted / Expired | Drives whether the code still works. |

## Crediting on redeem

- Quota type: `User.quota += amount` immediately, no approval
- Plan type: create a [Subscription](/en/schema/subscription) with `expire_at = now + Plan.duration_days`

## Code format

24-char random string by default. For bulk distribution, export CSV / TXT (one code per line).

## Related pages

- [Redemption Overview](/en/redemption/overview)
- [Redemption Management (admin)](/en/redemption/redemption-management)
- [User Guide](/en/redemption/user-guide)
- [Quota Calculation](/en/redemption/quotas)

## Related API

- `POST /api/redemption/` — create
- `GET /api/redemption/` — list
- `POST /api/redemption/redeem` — user redeem
- `DELETE /api/redemption/:id` — delete