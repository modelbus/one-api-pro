---
title: Model Price Management
description: How admins maintain per-model prices in the admin console.
category: pricing
order: 3
---

# Model Price Management

> Admin → Model Prices. Where each model's unit price lives.

## Where

Admin → Model Prices.

## List shows

Per row:

- Model name
- Input / output / cached / per-request prices
- Billing type (`token` / `per_request`)
- Status (Enabled / Disabled)
- Action buttons

## How to add

1. Top-right "Add"
2. Fill:
   - Model name (required, unique)
   - Input / output / cached / per-request price
   - Billing type (`token` default, or `per_request`)
   - Enabled (default on)
3. Save

## How to edit

Row → "Edit" → modify any field → Save. Takes effect immediately.

## How to disable

Row → "Disable":

- Disabled models can't be called (won't appear in channel-edit dropdown)
- In-flight calls are unaffected

Typical when: the upstream provider sunsets a model.

## How to delete

Row → "Delete" (double-confirm). Once deleted, the model is permanently unavailable.

Usually unnecessary — flipping status to "Disabled" is enough.

## Things to know

- **Unit is ¥ per million tokens.** Example: `0.002` means ¥0.002 per million tokens.
- **Cached price is usually half of input price.** `0` means no caching distinction.
- **After adding a model**, you also need to enable it in the relevant [channels](/channel/add-channel) so it can be routed.

## Default prices

On first start the system seeds default prices for common models (GPT-4o / Claude / DeepSeek etc.). Edit freely.

## FAQ

- **Changed price but calls aren't billed at the new rate**: the in-memory cache rebuilds within seconds; usually immediate.
- **"model price not set" error**: the model isn't in [Model Price](/en/pricing/model-price), or `enabled=false`.
- **One price for a whole provider family**: use `model_mapping` on the channel, but you still need a ModelPrice row per exact name.

## Related

- [Model Price (concept)](/en/pricing/model-price)
- [Model Price Schema](/en/schema/model-price)
- [Channel Management](/en/channel/add-channel)

## Related API

- `GET /api/model_price/` — list (Root)
- `POST /api/model_price/` — add (Root)
- `PUT /api/model_price/` — update (Root)
- `DELETE /api/model_price/:id` — delete (Root)