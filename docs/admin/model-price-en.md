---
title: Model Price
description: "Admin CRUD for model prices; also the read-only options endpoint used by the channel-edit dropdown."
category: admin
order: 3
---

# Model Price

> Per-model input / output / cached / per-request prices plus a billing-type switch, queried live by the relay billing path.

Route: Admin → **Setting → Pricing** (`/setting/pricing`). Served by `web/default-pro/src/views/setting/PricingSetting.vue`.
The CRUD tabs are Root-only; the options endpoint (used by the channel-edit dropdown) requires only Admin.

## Data Model

`model.ModelPrice` (table `model_price`):

| Field | Type | Notes |
|---|---|---|
| `model_name` | `varchar(100)` UNIQUE | Primary identifier |
| `input_price` | `decimal(16,6)` | Input token price |
| `output_price` | `decimal(16,6)` | Output token price |
| `cached_price` | `decimal(16,6)` | Cached token price |
| `per_request_price` | `decimal(16,6)` | Per-request flat price (e.g. dall-e, whisper) |
| `billing_type` | `varchar(20)` | `token` or `per_request` |
| `enabled` | `bool` | When `false`, the row is excluded from billing and the dropdown |
| `created_at` / `updated_at` | `bigint` | unix seconds |

On startup `InitDefaultPrices()` seeds a default catalogue (gpt-4o, claude-3.5, deepseek, qwen-plus, …); admins can freely add / edit / delete rows.

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/model_price/` | `GET` | Root | Full table dump (includes `enabled=false` rows) |
| `/api/model_price/options` | `GET` | Admin | **Only `enabled=true` `model_name`s**; consumed by the channel-edit dropdown |
| `/api/model_price/` | `POST` | Root | Insert a new row; `Insert()` zero-fills the id defensively to avoid PK collisions |
| `/api/model_price/` | `PUT` | Root | Update `input_price/output_price/cached_price/per_request_price/billing_type/enabled` for a given id |
| `/api/model_price/:id` | `DELETE` | Root | Delete a row |

Implementation: `controller/model_price.go`; the frontend uses the `@/api` module.

## Cache Behaviour

After Add / Update / Delete, `model.InitModelPriceCache()` rebuilds the in-memory `modelPriceMap`.
A background goroutine `SyncModelPriceCache(frequency)` re-syncs from the DB every `frequency` seconds (default 300) so multi-instance deployments stay eventually consistent.
The billing hot path uses `model.CacheGetModelPrice(modelName)` (Redis first, DB fallback), key `model_price:<name>`, TTL `ModelPriceCacheSeconds = 300`.

## Channel Dropdown Source

`GET /api/model_price/options` is what feeds the "available models" dropdown in the channel-edit form. Only `enabled=true` rows are returned, sorted by `model_name asc`; empty model_name rows are deduped and dropped.

## Billing Modes

- `billing_type = "token"` (default): charges on token usage, with cached tokens billed at `cached_price`.
- `billing_type = "per_request"`: charges per call, ignoring tokens (dall-e, whisper, tts, …).

`FindModelPriceByPattern` falls back to `strings.Contains(modelName, pattern)` when no exact key matches — useful when a series of variants should share one price row.

## Frontend Guide

- The pricing page is split into two tabs: **Model Price** and **Group Price**.
- Columns: model_name / input_price / output_price / cached_price / per_request_price / billing_type tag / actions.
- **Add** opens a 640 px modal: model_name (required), four price inputs with `precision=6`, billing-type select.
- Editing model_name is allowed but unique-index collisions will be rejected by GORM; saving triggers a full table refresh (and cache rebuild).
- Deletes use a confirmation popover.

## Implementation Pointers

| Concern | Location |
|---|---|
| CRUD handler | `controller/model_price.go` |
| Options handler | `controller/model_price.go::ListModelPriceOptions` |
| Default catalogue | `model/model_price.go::defaultModelPrices` |
| Cache init / sync | `model/model_price.go::InitModelPriceCache` / `SyncModelPriceCache` |
| Billing hot path | `model/model_price.go::CacheGetModelPrice` |
| Routes | `router/api.go` |