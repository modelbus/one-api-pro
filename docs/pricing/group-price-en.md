---
title: Group Price
description: Discounts for user groups on specific models.
category: pricing
order: 2
---

# Group Price

> VIP at 20% off, SVIP at 50% off specific models — this is where you set it.

## What it is

`GroupPrice` lets a user group apply a discount multiplier on specific models. It composes with [Model Price](/en/pricing/model-price):

```
deducted = consumption × ModelPrice × GroupPrice(group, model)
```

`GroupPrice` defaults to 1.0 (no discount).

## Where to find it

Admin → Group Prices.

## What to set

Per row:

| Field | Meaning | Effect |
|---|---|---|
| `group_name` | Group name | Matches the user's `group` field |
| `model_name` | Model name | Empty = default for the group; specific = that model only |
| `discount` | Multiplier | `1.0`=no change; `0.8`=20% off; `1.2`=20% markup |

## Match order

1. Exact `(user_group, model)` match
2. Otherwise `(user_group, "")` — the group's default multiplier
3. Otherwise 1.0

## Common recipes

### Flat VIP discount

| group | model | discount |
|---|---|---|
| `vip` | _(empty)_ | `0.8` |

### VIP + extra on a premium model

| group | model | discount |
|---|---|---|
| `vip` | _(empty)_ | `0.8` |
| `vip` | `gpt-4o` | `0.5` |

### Multiple groups

On first start the system seeds `default` / `vip` / `svip` at `discount=1.0`. Tune them in the UI as needed.

## How to change

Admin → Group Prices → Add:

1. Pick `group_name` (e.g. `vip`)
2. Fill `model_name` (empty = group default)
3. Fill `discount` (e.g. `0.8`)
4. Save

## Don't

- Don't set every group to 1.0 (it's already the default)
- Don't create two rows with the same `(group, model)` (unique constraint)

## vs. subscription discounts

[Subscription plans](/en/subscription/overview) have their own discount (`subscription.discount`). The two compose:

- Call hits a subscription → use `subscription.discount`
- Call misses → use `GroupPrice`

## FAQ

- **Set VIP 80% but user says it didn't apply**: confirm the user is actually in the `vip` group (Personal Center or admin)
- **Effective immediately?**: Yes. Next call uses the new discount.
- **Can one user stack discounts?**: Yes — `GroupPrice` × subscription discount (if active).

## Related

- [Group Price Management (admin)](./group-price-management)
- [Group Price Schema](/en/schema/group-price)
- [Model Price](/en/pricing/model-price)

## Related API

- `GET /api/group_price/` — list
- `POST /api/group_price/` — add
- `PUT /api/group_price/` — update
- `DELETE /api/group_price/:id` — delete