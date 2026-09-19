---
title: Model Price
description: Per-model unit price — the single source of truth for billing.
category: pricing
order: 1
---

# Model Price

> What does it cost the user per call? Look here.

## What it is

`ModelPrice` is the system's "price list". Each model has one row, deciding:

- **Token billing**: ¥ per million input / output tokens
- **Per-request billing**: ¥ per call (used for DALL·E / image generation etc.)

Per-call deduction = consumption × ModelPrice × group discount.

## Where to find it

- **Admin → Model Prices**: list, add, bulk edit
- **Admin → Dashboard**: today's / this month's revenue per model

## What to set

Admin → Model Prices → Add:

| Field | Meaning | Effect |
|---|---|---|
| `model_name` | Model name | Must match a value in [Channel → Models](/channel/add-channel); mismatched names can't be called |
| `billing_type` | `token` (per token) or `per_request` (per call) | Decides the four fields below |
| `input_price` | ¥ per million input tokens | Drives input-side deduction |
| `output_price` | ¥ per million output tokens | Drives output-side deduction |
| `cached_price` | ¥ per million cached tokens | Usually half of `input_price`; 0 = no cache |
| `per_request_price` | ¥ per call | Only used when `billing_type=per_request` |
| `enabled` | Toggle | Disabled → model can't be called |

> The system seeds default prices for common models (GPT-4o / Claude / DeepSeek etc.) on first start. Edit them freely.

## FAQ

- **New model but not in the channel dropdown**: add the model here with `enabled=true` first.
- **Call succeeded but no quota deducted**: check the model is in Model Price and the price isn't zero.
- **Changed price — how is it billed going forward**: calls use the model's current ModelPrice at the time of the call.

## Related

- [Model Price Management (admin)](./model-price-management)
- [Model Price Schema](/en/schema/model-price)
- [Group Price](/en/schema/group-price)

## Related API

- `GET /api/model_price/` — list (Root)
- `GET /api/model_price/options` — enabled model names only (Admin; for channel edit dropdown)
- `POST /api/model_price/` — add
- `PUT /api/model_price/` — update
- `DELETE /api/model_price/:id` — delete