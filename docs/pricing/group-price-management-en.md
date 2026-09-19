---
title: Group Price Management
description: How admins set per-group model discounts in the admin console.
category: pricing
order: 4
---

# Group Price Management

> Admin → Group Prices. Set per-group model discounts.

## Where

Admin → Group Prices.

## List shows

Per row:

- Group name (`default` / `vip` / `svip` etc.)
- Model name (empty = the group's default for all models)
- Discount multiplier (`0.8` = 20% off)
- Action buttons

## How to add

1. Top-right "Add"
2. Fill:
   - Group (required; dropdown from system-maintained list)
   - Model name (required; dropdown from enabled [Model Prices](/en/pricing/model-price))
   - Discount multiplier (`0.8` = 20% off, `1.2` = 20% markup; blank = `1.0`)
3. Save

## How to edit

Row → "Edit" → change `discount` → Save. Takes effect immediately.

## Match order

1. Exact `(user_group, model)` match
2. Otherwise `(user_group, "")` — the group's default
3. Otherwise 1.0

## Suggested setup

| group | model | discount | Effect |
|---|---|---|---|
| `vip` | _(empty)_ | `0.8` | VIP pays 80% across all models |
| `vip` | `gpt-4o` | `0.5` | VIP pays 50% on gpt-4o (overrides the group's 0.8 for that model) |
| `svip` | _(empty)_ | `0.7` | SVIP pays 70% across all models |

## Notes

- The same `(group, model)` pair cannot be added twice (unique constraint)
- Changing a group's discount affects every user in that group, immediately
- Don't set every group to 1.0 (already the default)

## How to assign groups to users

Groups are configured in [System Settings](../en/misc/system-settings); new users default to `default`. Admins change a user's group in [User Management](../en/user/user-management).

## FAQ

- **Discount set but user says it didn't apply**: verify the user's actual `group` field (not username)
- **Effective immediately?**: yes, next call uses the new discount
- **Plan discount vs group discount?**: plan discount wins first. See [Billing Rules](/en/subscription/billing-rules)

## Related

- [Group Price (concept)](/en/pricing/group-price)
- [Group Price Schema](/en/schema/group-price)
- [User Management (admin)](/en/user/user-management)

## Related API

- `GET /api/group_price/` — list (Root)
- `POST /api/group_price/` — add (Root)
- `PUT /api/group_price/` — update (Root)
- `DELETE /api/group_price/:id` — delete (Root)