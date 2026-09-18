---
title: Group Price
description: "Per (group × model) discount multipliers that compose with model_price."
category: pricing
order: 4
---

# Group Price

> Discount multiplier per `(group_name, model_name)` pair; combined with `model_price` to produce the final unit price for a user.

Route: Admin → **Setting → Pricing** (`/setting/pricing`) → **Group Price** tab. Served by `web/default-pro/src/views/setting/PricingSetting.vue`.
Group-price management is Root-only.

## Data Model

`model.GroupPrice` (table `group_price`):

| Field | Type | Notes |
|---|---|---|
| `group_name` | `varchar(32)`, part of composite unique index with `model_name` | User group (`default`, `vip`, `svip`, …) |
| `model_name` | `varchar(100)`, composite unique index | `""` means "all models for this group share one discount" |
| `discount` | `decimal(10,4)` | Multiplier. `1.0` = no change; `< 1.0` = discount, `> 1.0` = surcharge |

`InitDefaultPrices()` seeds three defaults (`default`, `vip`, `svip` with `discount=1.0`) on first run.

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/group_price/` | `GET` | Root | Full table dump |
| `/api/group_price/` | `POST` | Root | Insert; defaults `discount=1.0`; `Insert()` zeros the id defensively |
| `/api/group_price/` | `PUT` | Root | Update `discount` for a given id |
| `/api/group_price/:id` | `DELETE` | Root | Delete a row |

Implementation: `controller/model_price.go` (shared file with model-price handlers).

## Lookup Logic

`model.GetGroupDiscount(groupName, modelName)` falls back in this order:

1. Exact match `(groupName, modelName)` → that row's `discount`.
2. Fallback match `(groupName, "")` → the "all models in this group" discount.
3. Neither found → `1.0` (no change).

`GetGroupNames()` returns the set of distinct `group_name`s currently in cache — used as the user-edit group's dropdown data source.

## Caching Behaviour

`GroupPriceCacheSeconds = 300` (default). Every Add/Update/Delete calls `model.InitGroupPriceCache()` to rebuild the in-memory `groupPriceMap`. The billing hot path uses `model.CacheGetGroupPrice(groupName, modelName)` — Redis-first with key `group_price:<group>:<model>`, DB fallback.

## How It Composes

Final unit price = `model_price.input_price` × `group_price.discount` (resolved per model per group).
The discount also applies to `per_request` rows.

## Frontend Guide

- Columns: group_name / model_name / discount / actions.
- **Add** opens a 520 px modal:
  - `group_name` (required; dropdown from `GET /api/group/`).
  - `model_name` (required; dropdown from `/api/model_price/options`).
  - `discount` (precision 4, min 0; blank defaults to 1.0).
- Row actions: edit, delete (with confirmation).
- Changing a group's discount immediately affects every user in that group for the matching model; the periodic `SyncGroupPriceCache` background task is the safety net.

## Implementation Pointers

| Concern | Location |
|---|---|
| CRUD handler | `controller/model_price.go` (`AddGroupPrice` / `UpdateGroupPrice` / `DeleteGroupPrice`) |
| Default group rows | `model/model_price.go::defaultGroupPrices` |
| Cache init / sync | `model/model_price.go::InitGroupPriceCache` / `SyncGroupPriceCache` |
| Billing hot path | `model/model_price.go::CacheGetGroupPrice` |
| Group dropdown source | `controller/group.go::GetGroups` → `model.GetGroupNames` |
| Routes | `router/api.go` |