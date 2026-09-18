---
title: Top-up Settings
description: "Top-up enable switch, presets, allow_custom flag, exchange rate and currency."
category: admin
order: 14
---

# Top-up Settings

> Global toggle, preset chips, allow-custom flag and exchange rate for the online top-up flow. UI: `web/default-pro/src/views/setting/TopupSetting.vue`.

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/setting/topup` | `GET` | Root | Returns `{ enabled, allow_custom, exchange_rate, presets }` in one shot |
| `/api/setting/topup` | `PUT` | Root | Whole-bundle save (same shape) |

Implementation: `controller/topup.go::GetTopupSettings` / `PutTopupSettings`.

## Fields

| Field | Type | Notes |
|---|---|---|
| `enabled` | `bool` | Master switch; `false` blocks `/api/topup/order` with "充值功能未开启" |
| `allow_custom` | `bool` | Allow users to type their own amount; otherwise they must pick a preset |
| `exchange_rate` | `int` | `1 CNY = X quota`; must be `> 0`; defaults to `1` (1:1) |
| `presets` | `[{ amount, bonus_quota }]` | Quick-amount chip list; every `amount` must be `> 0` and unique |

Validation lives in `model/topup.go::SaveTopupSettings`:
- `amount <= 0` → "第 N 项金额必须大于 0".
- `bonus_quota < 0` → "第 N 项额度不能为负数".
- Duplicate `amount` → "快捷金额重复：X.XX 元已存在".
- `exchange_rate <= 0` → "兑换比例必须大于 0".

## Persistence

Four independent rows in `system_settings` (category = `topup`):

| key | value shape |
|---|---|
| `topup.enabled` | `"true"` / `"false"` |
| `topup.allow_custom` | `"true"` / `"false"` |
| `topup.exchange_rate` | integer string |
| `topup.presets` | JSON array: `[{"amount":10,"bonus_quota":10000}, ...]` |

PUT is full-overwrite. Missing fields are not preserved — the frontend always sends all four.

## User-side Order Resolution

`POST /api/topup/order` (`controller/topup.go::CreateTopupOrder`) delegates to `model/topup.go::ResolveTopupAmount`:

1. If `preset_amount > 0`, look up `presets` for a matching `amount`; hit returns `{ amount, bonus_quota }`; miss → "快捷金额未配置".
2. Otherwise compute `bonus_quota = amount * exchange_rate`.
3. When `allow_custom=false`, the custom path is rejected.

## Preset Chips

The user-facing `PaymentTopupModal.vue` renders `presets` as chips. When `allow_custom=true`, a trailing "自定义" chip is appended and `-1` is the sentinel index for "custom" (see `web/default-pro/src/utils/topup.js::validateTopupPresets`).

## Frontend Guide

- Top section: three vertical `<a-form-item>` rows for `enabled`, `allow_custom`, and `exchange_rate` (precision 0, min 1).
- Below: a bordered `<a-table>` for the presets. Add / remove rows; edit `amount` (precision 2, min 0.01) and `bonus_quota` (precision 0, min 0, step 1000) inline.
- Submit triggers a frontend `validateTopupPresets` check, then `PUT /api/setting/topup`.

## Implementation Pointers

| Concern | Location |
|---|---|
| Handler | `controller/topup.go` |
| Settings read/write | `model/topup.go::GetTopupSettings` / `SaveTopupSettings` |
| Amount resolution | `model/topup.go::ResolveTopupAmount` |
| Order creation | `model/topup.go::CreateTopupOrder` |
| Frontend validator | `web/default-pro/src/utils/topup.js` |
| Routes | `router/api.go` |