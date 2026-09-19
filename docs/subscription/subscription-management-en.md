---
title: Subscription Management
description: How admins view, adjust, and revoke user subscriptions.
category: subscription
order: 7
---

# Subscription Management

> Admin → Subscriptions.

## List

Shows all active user subscriptions. Each row:

- User
- Plan
- Status (Active / Expired / Revoked)
- Start / End time
- Used / Total quota
- Action buttons

Filter by user, plan, or status.

## Actions

### Grant a subscription manually

1. Top-right "Grant" → pick user, plan, start time
2. Confirm → user gets it immediately, no payment needed

Common for: customer support compensation, event gifts, internal testing.

### Adjust expiry

1. Find the subscription → "Adjust expiry"
2. Modify the date → confirm

Lengthen or shorten the validity; doesn't affect used quota.

### Revoke

1. Find the subscription → "Revoke"
2. Double-confirm → subscription stops immediately

Common for: violation handling, refund-linked revocation.

## Status

| Status | Meaning |
|---|---|
| Active | Currently usable |
| Expired | Past `end_time`, auto-set by periodic job |
| Revoked | Admin-revoked |

## Bulk actions

Bulk revoke is supported (e.g. mass-ban multiple users).

## FAQ

- **No subscription shown for a user**: maybe they never subscribed; remove the "Active" filter.
- **Revoke — does their top-up quota disappear?** No. Revoke only ends the subscription. Top-up quota still works for pay-as-you-go.
- **Edited a plan — does it change active subscriptions?** No. Plans use a snapshot at purchase time.

## Related

- [Subscription Schema](../schema/subscription)
- [Plan Management (admin)](./plan-management)
- [Subscription Management (admin)](./admin-guide)