---
title: Group Pricing
description: "`group_price` table, discount multipliers, per-model overrides, and default groups."
category: pricing
order: 2
---

# Group Pricing

> User-group multipliers on top of model pricing. `group_price` defines "this group on this model" at a decimal multiplier; default `1.0` (no discount). Implementation: `model/model_price.go::GroupPrice`.

## Data model

`model.GroupPrice` (`group_price` table):

| Field | Type | Notes |
|---|---|---|
| `group_name` | `varchar(32)` | Matches `users.group`; first column of the composite unique index |
| `model_name` | `varchar(100)` | Model name; empty string means "default multiplier for this group", second column of the unique index |
| `discount` | `decimal(10,4)` | Multiplier: `1.0` no discount; `0.8` charges 80% (vip 20% off) |
| `created_at` / `updated_at` | `bigint` | unix seconds |

`(group_name, model_name)` is the unique key; the same combination cannot be inserted twice.

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/group_price/` | `GET` | Admin | Full list, ordered by `group_name, model_name` |
| `/api/group_price/` | `POST` | Admin | Create; `group_name` required; `discount=0` falls back to `1.0` |
| `/api/group_price/` | `PUT` | Admin | Update |
| `/api/group_price/:id` | `DELETE` | Admin | Delete (rebuilds `groupPriceMap` after) |

## Lookup order

`model.GetGroupDiscount(groupName, modelName)`:

1. `groupPriceMap[groupName][modelName]` — exact hit
2. `groupPriceMap[groupName][""]` — group-wide default
3. Otherwise `1.0`

`relay/billing/ratio/model.go::GetGroupDiscount(groupName, modelName, fallbackNames...)` retries step 1 with each fallback name before falling back to `""` (e.g. `claude-3.5-sonnet-internet` → `claude-3.5-sonnet`); the fallback chain mirrors `GetModelPrice`.

## Default groups

`model_price.go::defaultGroupPrices` seeds three rows:

```go
{GroupName: "default", ModelName: "", Discount: 1.0}
{GroupName: "vip",     ModelName: "", Discount: 1.0}
{GroupName: "svip",    ModelName: "", Discount: 1.0}
```

Inserted when the table is empty. `Discount=1.0` is a placeholder; admins tune them in the UI.

## Cache

- Boot: `InitGroupPriceCache` loads the entire table into `groupPriceMap map[string]map[string]float64`
- After writes: every CRUD calls `InitGroupPriceCache`
- Scheduled: `SyncGroupPriceCache(frequency)` shares the frequency with `SyncModelPriceCache`
- Redis: `CacheGetGroupPrice(group, model)` second-tier cache (`ModelPriceCacheSeconds=300`); DB on miss

## Coverage interaction

- A row with `model_name=""` is the group's default for every model
- A per-model row (`model_name="gpt-4o"`) only takes precedence for that exact model; other models fall back to the group's `""` row

Examples:

| group | model | discount | Meaning |
|---|---|---|---|
| `vip` | `""` | `0.8` | vip users pay 80% on all models |
| `vip` | `gpt-4o` | `0.5` | vip users pay 50% on gpt-4o |
| `svip` | `""` | `1.0` | svip placeholder (no discount) |

## Frontend Guide

Route: `/setting/pricing` → **Group Pricing** tab (`web/default-pro/src/views/setting/PricingSetting.vue`).

- List: `group_name` / `model_name` / `discount`
- Edit modal: the three fields; `discount` defaults to `1.0`
- **Add** uses the same form

## Implementation Pointers

| Concern | Location |
|---|---|
| Data model | `model/model_price.go::GroupPrice` |
| Cache init | `model/model_price.go::InitGroupPriceCache` |
| Lookup (with fallback) | `model/model_price.go::GetGroupDiscount` |
| Redis cache | `model/model_price.go::CacheGetGroupPrice` |
| CRUD | `controller/model_price.go` |
| Default groups | `model/model_price.go::defaultGroupPrices` |
| Billing settlement | `relay/handler/helper.go::postConsumeQuota` |
