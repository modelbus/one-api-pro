---
title: Upgrade & Downgrade
description: "Two paths when switching plans: price-diff and stack."
category: subscription
order: 5
---

# Upgrade & Downgrade

> You already have a subscription. What happens if you buy a different plan?

## Two paths

Admin chooses in [Plan Settings](./plan-settings):

### Price-diff (default)

- Upgrading to a more expensive plan → pay only the **difference**
- Remaining quota on the old plan is pro-rated into the new subscription
- Order number prefix: `UP`

Example: Plan A has 20 days left, 30% used. You upgrade to Plan B (B is ¥30 more than A). You pay about `¥30 × 80% = ¥24` (pro-rated by remaining days).

### Stack

- Pay full price for the new plan
- Old subscription **keeps running** until its original expiry
- Order number prefix: `TB`

Good for: temporarily trying a higher plan, or holding two plans at once.

## When orders get rejected

- **Same plan or lower plan**: price-diff rejects it ("already on the same or higher tier"); stack allows it (two subscriptions side by side)
- **Current subscription is expired**: treated as a fresh purchase
- **Target plan is unpublished**: rejected with "plan not available"

## How to do it

1. User: pick a new plan in `/pricing` → place order
2. System picks the right path automatically
3. After payment it takes effect immediately

When an admin grants a new subscription manually (bypasses payment), the same path is used.

## FAQ

- **Differential is negative**: downgrading from an expensive plan to a cheap one — old plan's remaining value can exceed the new price. Confirm with admin.
- **Upgraded but old plan still charges**: only happens in stack mode. Price-diff mode replaces, doesn't double-charge.
- **Can I cancel an upgrade**: before payment yes; after payment ask the admin to revoke.

## Related

- [Subscription Overview](./overview)
- [Plan Settings](./plan-settings)
- [Plan Management (admin)](./plan-management)
- [My Orders](../user/orders)