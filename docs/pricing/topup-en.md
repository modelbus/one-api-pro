---
title: Top-up Settings
description: "`topup.*` system settings, master switch, presets, custom amount and rate, supported payment channels."
category: pricing
order: 4
---

# Top-up Settings

> Top-up is the user-facing flow that turns a paid amount into `users.quota`. Settings live in the `system_settings` table under `topup.*`; the user flow goes through `OrderTypeTopup=2`. Implementation: `model/topup.go`, `controller/topup.go`, `web/default-pro/src/views/setting/TopupSetting.vue`.

## System settings

Defined in `model/system_setting.go`:

| Key | Type | Notes |
|---|---|---|
| `topup.enabled` | bool string | Master switch; default `false` |
| `topup.allow_custom` | bool string | Allow user-entered custom amount; default `false` |
| `topup.presets` | JSON | Preset amount list: `[{"amount": 10, "bonus_quota": 10000}, ...]` |
| `topup.exchange_rate` | int64 string | Custom amount: 1 yuan = X quota (default `1`, i.e. 1:1) |

Category constant: `SystemSettingCategoryTopup="topup"`.

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/setting/topup` | `GET` | Root | Read all four settings (`TopupSettingsResponse`) |
| `/api/setting/topup` | `PUT` | Root | Save (whole-object replace) |
| `/api/topup/order` | `POST` | User | Create top-up order (see below) |

`SaveTopupSettings` validates:

- Every preset: `amount > 0`, `bonus_quota >= 0`; `amount` must be unique (`快捷金额重复：X.XX 元已存在`)
- `exchange_rate > 0` (`兑换比例必须大于 0`)
- `presets == nil` persists as `[]`

## Presets

`model.TopupPreset`:

```json
{
  "amount": 10,
  "bonus_quota": 10000
}
```

- `amount` — the yuan amount the user pays
- `bonus_quota` — the integer quota credited (e.g. 10 yuan at 1:1000 rate → 10000 quota)

`ResolveTopupAmount`:

- `PresetAmount > 0` → look up the matching preset; not found → `快捷金额未配置`
- Otherwise fall back to the custom path; requires `allow_custom == true`

## Custom amount and rate

When `PresetAmount == 0`:

```text
bonus_quota = int(amount × exchange_rate)
```

Default `exchange_rate = 1` (1 yuan = 1 quota). Production usually raises it to 1:1000 or similar so the UI doesn't show sub-yuan amounts.

## Top-up order

`POST /api/topup/order` body (`CreateTopupOrderRequest`):

```json
{
  "amount": 100,
  "preset_amount": 0,
  "pay_method": "wechat"
}
```

Server validation:

1. `amount <= 0 && preset_amount <= 0` → `amount 或 preset_amount 至少传一个`
2. `pay_method` must be in the `OrderPayMethod*` whitelist (`wechat` / `alipay` / `bank` / `offline` / `free`)
3. Self-service accepts only `wechat` / `alipay` / `bank` (`自助充值仅支持 wechat / alipay / bank`)
4. At least one payment channel must be enabled (`payment.AnyChannelEnabled()`); otherwise `系统尚未开通任何支付通道，请设置后开启支付`
5. `topup.enabled == true`; otherwise `充值功能未开启`
6. Resolve amount / quota via `ResolveTopupAmount`

`model.CreateTopupOrder` writes:

- Order number prefix `TP` (`GenerateOrderNo("TP")`)
- `Type=OrderTypeTopup=2`, `Source=OrderSourceUserSelf=1`, `Status=OrderStatusPending=0`, `PayStatus=OrderPayStatusPending=0`
- `PlanInfo` stores `TopupOrderPlanInfo` JSON (`amount` / `preset_amount` / `bonus_quota` / `exchange_rate` snapshot) for reconciliation

Then `buildPayInfo(pay_method, order_no, amount, "余额充值")` returns pre-pay params to the frontend.

## Activation

`controller/payment.go::processNotify` dispatches by `order.Type`:

- `OrderTypeTopup=2` → `model.ActivateTopupByOrder(order)`:
  - Idempotent: `status == Paid` returns immediately
  - Parses `PlanInfo` for `bonus_quota`; on failure, falls back to `amount × exchange_rate`
  - `IncreaseUserQuota(order.UserId, bonus_quota)` adds to `users.quota`
  - `MarkOrderPaid` writes `pay_status=1` / `pay_time` / `pay_trade_no`

## Migration note

Legacy `plan.allow_topup` (DB row preserved) was migrated to `topup.enabled`; the UI no longer reads or writes it. To inspect the legacy value: `SELECT * FROM system_settings WHERE 'key' = 'plan.allow_topup'`.

## Frontend Guide

Page: `/setting/operation` → Top-up section (`web/default-pro/src/views/setting/TopupSetting.vue`).

- Switches: `enabled` / `allow_custom`
- Rate: `exchange_rate` (integer, min 1)
- Preset table: inline edit `amount` + `bonus_quota`; add / remove rows; saving replaces the whole set
- User-side `/pricing` or `/topup`: preset chips + custom amount input (only when `allow_custom=true`)

## Implementation Pointers

| Concern | Location |
|---|---|
| Setting keys | `model/system_setting.go::SystemSettingKeyTopup*` |
| Read / save settings | `model/topup.go::GetTopupSettings` / `SaveTopupSettings` |
| Amount resolution | `model/topup.go::ResolveTopupAmount` |
| Order creation | `model/topup.go::CreateTopupOrder` |
| Activation | `model/topup.go::ActivateTopupByOrder` |
| Settings endpoints | `controller/topup.go::GetTopupSettings` / `PutTopupSettings` |
| User order endpoint | `controller/topup.go::CreateTopupOrder` |
| Payment notify dispatch | `controller/payment.go::processNotify` |
