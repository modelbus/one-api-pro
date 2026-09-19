---
title: Subscription Management (Admin)
description: How admins create plans, manually grant subscriptions, and revoke them.
category: subscription
order: 3
---

# Subscription Management (Admin)

> Admin → Subscriptions. Visible only to admin accounts.

## Plan management

Admin → Plans. You can:

- **Create a plan**: name, price, quota (by token or by request), validity, discount multiplier, description
- **Publish / Unpublish**: published plans show on the public subscribe page; unpublished plans stop new subscriptions but existing ones are unaffected
- **Edit**: change any field
- **Delete**: unpublish first, then delete (recommended to avoid accidents)

See [Plan Management (admin)](./plan-management).

## User subscription management

Admin → Subscriptions. Lists every user's active subscriptions.

### Grant a subscription manually (admin compensation)

1. Find the user
2. Click "Grant"
3. Pick a plan
4. Pick start time (now / N days later)
5. Confirm

The user gets the plan immediately, no payment needed.

### Revoke a subscription

1. Find the subscription
2. Click "Revoke"
3. Double-confirm

Subscription stops immediately. Already-used quota is not refunded.

### Adjust subscription expiry

1. Find the subscription
2. Click "Adjust expiry"
3. Modify the date
4. Confirm

Common for: customer support compensation, trial extensions.

## Business rules

[Plan Settings](./plan-settings) controls:

- **Upgrade mode**: price-diff vs stack (default: price-diff)
- More global rules will be added later

## FAQ

- **Will the user see a manually granted subscription?** Yes. It appears in their "My Subscriptions" immediately.
- **Can they re-subscribe after a revoke?** Yes. They can buy it themselves, or you grant it again.
- **Do plan changes apply to old subscriptions?** No. Old subscriptions use the plan as it was at purchase time; new subscriptions use the new settings.

## Related

- [Subscription Schema](../schema/subscription)
- [Plan Schema](../schema/plan)
- [Plan Management](./plan-management)
- [User Management (admin)](../user/user-management)