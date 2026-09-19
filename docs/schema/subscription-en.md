---
title: Subscription
description: An active subscription instance for a user, based on a Plan.
category: schema
order: 8
---

# Subscription

## What it is

`Subscription` is an instance of a [Plan](/en/schema/plan). When a user subscribes in the public UI, the system creates a Subscription and applies the Plan's rules.

## Where to find it

- **Personal Center → My Subscriptions**: active plan, days left, used / total quota
- **Admin → Subscriptions**: every user's subscription with expiry and status
- **Call logs**: each call logs the subscription it matched

## Operator-relevant fields

| Field | Meaning | Effect |
|---|---|---|
| Parent | User FK | Cascade-delete with user. |
| Plan FK | Plan this is based on | Editing a Plan does NOT retroactively update existing subscriptions. |
| Start time | When it began | Usually set automatically at checkout. |
| Expire time | When it stops | After expiry → "plan expired" but calls keep working until quota is exhausted. |
| Used quota | Cumulative consumption | Recorded only. |
| Status | Active / Expired / Cancelled | Drives whether the subscription discount applies. |

## Effect on billing

The full deduction formula:

```
final = consumption × ModelPrice × GroupPrice(group, model)
            × Subscription.discount (if active subscription covers this model)
```

- No active subscription → no subscription discount
- Expired subscription → same (no discount), but the call still works
- Multiple subscriptions during upgrade transition → only the active one applies

## Subscription vs User quota

- User quota (`User.quota`): deducted first when billing
- Subscription quota: cap on what the subscription discount applies to; after exhaustion the discount stops
- They are independent deduction sources — see [Billing Rules](/en/subscription/billing-rules)

## Related pages

- [Subscription Overview](/en/subscription/overview)
- [User Guide](/en/subscription/user-guide)
- [Subscription Management (admin)](/en/subscription/subscription-management)

## Related API

- `GET /api/subscription/self` — current user's subscriptions
- `GET /api/subscription/` — list (admin)
- `POST /api/subscription/` — create (admin manual grant)
- `PUT /api/subscription/` — update
- `DELETE /api/subscription/:id` — delete