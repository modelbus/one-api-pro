---
title: Top-up Settings
description: Enable / disable online top-up and configure pay/credit amounts for presets.
category: pricing
order: 9
---

# Top-up Settings

> Admin → Top-up Settings. Visible only to Root.

## Configurable fields

| Field | Meaning | Effect |
|---|---|---|
| `enabled` | Master switch | Off hides the public top-up entry |
| `allow_custom` | Allow custom amount | Off forces users to pick a preset; on, custom amounts are always 1:1 (credit equals pay) |
| `presets` | Preset amounts | Each entry is `{amount, bonus_quota}`; edited in the UI as a "Pay Amount (CNY) / Credit Amount (CNY)" pair |

> The quota base is a system constant: **¥1 = 1,000,000 quota** (1 quota = 1e-6 CNY), a "cent"-like unit used in computation only — never stored and never modifiable. `topup.exchange_rate` has been removed.

## Recommended setup

Enable + allow custom + a few presets (¥1 = 1,000,000 quota):

| Pay Amount (CNY) | Credit Amount (CNY) | bonus_quota |
|---|---|---|
| 10 | 10 | 10000000 |
| 50 | 50 | 50000000 |
| 100 | 100 | 100000000 |
| 500 | 500 | 500000000 |

Promotional example "pay 10 get 15":

| Pay Amount (CNY) | Credit Amount (CNY) | bonus_quota |
|---|---|---|
| 10 | 15 | 15000000 |

> Credit may exceed pay (a bonus) but must **never** be lower than pay (no under-crediting).

## How to change

Admin → Top-up Settings → edit fields → Save.

Validation:

- Each preset `amount > 0`
- No duplicate pay amounts
- `bonus_quota >= round(amount × 1000000)` (credit not below pay)

## User-facing impact

| Field | What users see |
|---|---|
| `enabled = false` | Top-up entry gone |
| `allow_custom = false` | Only preset chips visible; input field hidden |
| `presets` | Number and order of chips; each chip shows both pay and credit amounts |
| Custom amount | Always 1:1 — credit equals pay |

## How to test

1. After enabling, log in as a test account and visit the top-up page
2. Click a preset → credited quota = `bonus_quota`
3. Enter a custom amount → credited quota = `round(amount × 1000000)`

## Notes

- Preset changes only affect **new orders**; existing orders keep their snapshot (`bonus_quota` / `credit_amount` in `plan_info`)
- Custom amounts no longer have a rate setting — they are always 1:1
- Preset amounts must be unique (constraint)

## FAQ

- **All settings saved but user still can't see top-up**: verify `enabled=true` and at least one payment channel enabled
- **Why does saving fail with "到账额度不能低于支付金额折算额度"?**: the credit amount is lower than the pay amount; raise the credit or lower the pay
- **Can preset amounts be decimals?**: yes — both pay and credit support two decimals (e.g. `9.99` CNY)

## Related

- [Top-up (concept)](./topup)
- [Top-up Management (admin)](./topup-management)
- [Payment Channels](./payment)

## Related API

- `GET /api/setting/topup` — read settings (Root)
- `PUT /api/setting/topup` — update settings (Root, **full replacement**)
