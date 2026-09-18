---
title: Model Pricing
description: "`model_price` table, units, field semantics, and the `ListModelPriceOptions` dropdown source."
category: pricing
order: 1
---

# Model Pricing

> Each model has a row in `model_price` that drives per-token / per-request billing. Implementation: `model/model_price.go::ModelPrice`, `controller/model_price.go`.

## Data model

`model.ModelPrice` (`model_price` table):

| Field | Type | Notes |
|---|---|---|
| `model_name` | `varchar(100)` UNIQUE | Model name; matched against the request model via `GetModelPrice` |
| `input_price` | `decimal(16,6)` | Input unit price (¥ / 1M tokens) |
| `output_price` | `decimal(16,6)` | Output unit price (¥ / 1M tokens) |
| `cached_price` | `decimal(16,6)` | Cached-input unit price (¥ / 1M tokens), typically half of `input_price` |
| `per_request_price` | `decimal(16,6)` | Per-request price for non-token models (DALL·E / TTS / embeddings) |
| `billing_type` | `varchar(20)` | `token` / `per_request` |
| `enabled` | `bool` | Whether the row is active; the dropdown only shows `enabled=true` rows |

Prices are in ¥ per 1M tokens. The actual quota deduction multiplies by the group discount (see [Group Pricing](./group-price)).

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/model_price/` | `GET` | Root | Full list (including `enabled=false`) |
| `/api/model_price/options` | `GET` | Admin | Returns the `model_name` list of `enabled=true` rows only — for the channel-edit dropdown |
| `/api/model_price/` | `POST` | Admin | Create (`model_name` required; `billing_type` defaults to `token`) |
| `/api/model_price/` | `PUT` | Admin | Update (rebuilds `modelPriceMap` after success) |
| `/api/model_price/:id` | `DELETE` | Admin | Delete (rebuilds cache after success) |

`GetAllModelPrices` in `controller/model_price.go` is marked Deprecated (Root only) because it leaks the billing fields into contexts that don't need them; prefer `ListModelPriceOptions` for new UI surfaces.

## Cache

`InitModelPriceCache` seeds an in-process `modelPriceMap` (`map[string]*ModelPrice`) of `enabled=true` rows, guarded by a RWMutex:

- `GetModelPrice(name)` reads the map under RLock
- `FindModelPriceByPattern(name)` first tries the exact match, then substring match (handles `-internet` suffixes)
- Every write (POST / PUT / DELETE) rebuilds the map via `InitModelPriceCache`
- The scheduled task `SyncModelPriceCache(frequency)` rebuilds it every N seconds (`config.SyncFrequency`)

`CacheGetModelPrice(name)` / `CacheGetGroupPrice(group, model)` are the Redis tier (`ModelPriceCacheSeconds=300`); Redis miss falls through to DB and warms back.

## Billing flow

`relay/billing/ratio/model.go::GetModelPrice(name, fallbackNames...)`:

1. Try `name` against `model.GetModelPrice`
2. Fall back through `fallbackNames` (e.g. `gpt-4o-internet` → `gpt-4o`)
3. Otherwise return `errors.New("model price not set")`

Then `postConsumeQuota` computes the quota based on `BillingType`:

```go
if BillingType == PerRequest {
  quota = per_request_price * 1 * group_discount
} else {
  quota = (input_price * prompt_tokens
         + output_price * completion_tokens
         + cached_price * cached_tokens)
         * group_discount / 1_000_000   // ¥ → quota
}
```

## Default prices

`model_price.go::defaultModelPrices` lists common models seeded at first boot (`gpt-4o` / `claude-3.5-sonnet` / `deepseek-chat` / `qwen-max` / `gemini-1.5-pro` etc.). `InitDefaultPrices` only runs when the table is empty, using `OnConflict{DoNothing: true}` so re-runs are safe.

## Frontend Guide

Route: `/setting/pricing` → **Model Pricing** tab (`web/default-pro/src/views/setting/PricingSetting.vue`).

- List columns: `model_name` / `input_price` / `output_price` / `cached_price` / `per_request_price` / `billing_type` / `enabled`
- Edit modal: the six numeric fields, a billing-type dropdown (`token` / `per_request`), and the enabled switch
- **Add** uses the same form

## Implementation Pointers

| Concern | Location |
|---|---|
| Data model | `model/model_price.go::ModelPrice` |
| Cache init | `model/model_price.go::InitModelPriceCache` |
| Lookup (exact / fuzzy) | `model/model_price.go::GetModelPrice` / `FindModelPriceByPattern` |
| Redis cache | `model/model_price.go::CacheGetModelPrice` |
| Dropdown source | `controller/model_price.go::ListModelPriceOptions` |
| CRUD | `controller/model_price.go` |
| Default prices | `model/model_price.go::InitDefaultPrices` |
| Billing settlement | `relay/handler/helper.go::postConsumeQuota` |
