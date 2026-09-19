---
title: Group Price
description: Group discount multiplier on specific models; defaults to 1.0 when not set.
category: schema
order: 6
---

# Group Price

## What it is

`GroupPrice` lets a user group apply a discount multiplier on a model (e.g. 0.8 = 20% off). It composes with [Model Price](/en/schema/model-price):

```
deducted = consumption × ModelPrice × GroupPrice(group, model)
```

## Where to find it

- **Admin → Group Prices**: list, add, bulk edit
- The logged-in user's group and the matched GroupPrice are shown in Personal Center as "Current discount".

## Operator-relevant fields

| Field | Meaning | Effect |
|---|---|---|
| Group | `default` / `vip` / `svip` | Matches the user's `group` field. |
| Model | e.g. `gpt-4o` | Empty string = default multiplier for **all models** in this group. |
| Multiplier | 1.0 = no change; 0.8 = 20% off; 1.2 = 20% markup | Multiplies the deduction. |

## Match order

1. Exact `(group, model)` match
2. Otherwise `(group, "")` — the group's default multiplier
3. Otherwise 1.0

## Common recipes

- **Flat VIP discount across all models**

  | Group | Model | Multiplier |
  |---|---|---|
  | vip | _(empty)_ | 0.8 |

- **VIP + extra discount on a premium model**

  | Group | Model | Multiplier |
  |---|---|---|
  | vip | _(empty)_ | 0.8 |
  | vip | gpt-4o | 0.6 |

## Don't

- Don't set every group to 1.0 (that is the default).
- Don't set two rows with the same `(group, model)` — unique constraint violation.

## Related pages

- [Group Price (business concept)](/en/pricing/group-price)
- [Group Price Management (admin)](/en/pricing/group-price-management)

## Related API

- `GET /api/group_price/` — list
- `POST /api/group_price/` — add
- `PUT /api/group_price/` — update
- `DELETE /api/group_price/:id` — delete