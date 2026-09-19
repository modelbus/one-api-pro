---
title: Plan Settings
description: "Global plan rules: upgrade mode, etc."
category: subscription
order: 13
---

# Plan Settings

> Admin → Plan Settings.

## Currently configurable

### Upgrade mode

When a user with an active subscription buys a new plan, which mode is used?

- **Price-diff** (default): pay only the difference; remaining quota is pro-rated.
- **Stack**: pay full price for a new subscription; the old one keeps running.

See [Upgrade & Downgrade](./upgrade-downgrade).

## How to change

Admin → Plan Settings → pick the mode → save.

Takes effect immediately for all subsequent upgrades; in-flight upgrades are unaffected.

## When to switch

| Scenario | Recommended mode |
|---|---|
| Standard SaaS | Price-diff (default) |
| Frequent upgrades / trial use | Stack |
| Want upgrades to feel like "renewal", not "switch" | Price-diff |
| Want stacked subscriptions | Stack |

## FAQ

- **Changed mode — what about old subscriptions?** Not affected. Old subscriptions follow the rules at their purchase time.
- **Can I cancel a price-diff upgrade?** Not by the user. After payment ask admin to revoke.

## Related

- [Upgrade & Downgrade](./upgrade-downgrade)
- [Subscription Overview](./overview)
- [Plan Management (admin)](./plan-management)