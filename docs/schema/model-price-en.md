---
title: Model Price
description: Per-model unit price (token or per-request), the single source of truth for billing.
category: schema
order: 5
---

# Model Price

## What it is

`ModelPrice` is the **per-model unit price used for billing inside One API Pro**. After a call hits a channel, the deduction amount is fully decided by `call_consumption × ModelPrice`.

Think of it as the system's price list.

## Where to find it

- **Admin → Model Prices**: list, add, bulk edit
- **Admin → Dashboard**: today's / this month's revenue per model (computed from this table)

## Operator-relevant fields

| Field | Meaning | Effect |
|---|---|---|
| Model name | e.g. `gpt-4o` / `claude-sonnet-4` | Must match what [Channels](/en/schema/channel) advertise. |
| Billing type | `token` (per token) or `per_request` | Decides the semantics of the 3 fields below. |
| Input price | ¥ per million tokens | Drives input-side deduction. |
| Output price | ¥ per million tokens | Drives output-side deduction. |
| Cached price | ¥ per million tokens | For cached token segments (0 = no cache). |
| Per-request price | ¥ per call | Used when billing type is `per_request`. |
| Enabled | Toggle | Disabled → the model cannot be called. |

## Relationship to Channels

- The model selector inside a channel pulls from enabled ModelPrice entries.
- Editing ModelPrice does not affect channel availability, only the deduction amount.

## Relationship to Group Price

Final deduction:

```
deducted = consumption × ModelPrice × GroupPrice(group, model)
```

`GroupPrice` defaults to 1.0 (no discount). It only applies when the user is in a group that has an override for the model.

## Why it matters

Without ModelPrice, calls still go through, but no quota is deducted — effectively free. The system must be bootstrapped with ModelPrice entries before any real usage.

## Related pages

- [Model Price (business concept)](/en/pricing/model-price)
- [Model Price Management (admin)](/en/pricing/model-price-management)
- [Group Price](/en/schema/group-price)

## Related API

- `GET /api/model_price/` — list
- `POST /api/model_price/` — add / update
- `GET /api/model_price/options` — dropdown for channel edit page