---
title: Top-up Settings
description: Enable / disable online top-up, configure presets and exchange rate.
category: pricing
order: 9
---

# Top-up Settings

> Admin → Top-up Settings. Visible only to Root.

## Configurable fields

| Field | Meaning | Effect |
|---|---|---|
| `enabled` | Master switch | Off hides the public top-up entry |
| `allow_custom` | Allow custom amount | Off forces users to pick a preset |
| `exchange_rate` | ¥1 = X quota | Drives custom-amount conversion; must be > 0 |
| `presets` | Preset amounts | Each entry is `{amount, bonus_quota}` |

## Recommended setup

Enable + allow custom + a few presets:

| amount | bonus_quota |
|---|---|
| 10 | 100000 |
| 50 | 500000 |
| 100 | 1000000 |
| 500 | 5000000 |

`exchange_rate = 100000` (¥1 = 100k quota, matches default `QuotaPerUnit=500000`).

## How to change

Admin → Top-up Settings → edit fields → Save.

Validation:

- Each preset `amount > 0`, `bonus_quota >= 0`
- No duplicate amounts
- `exchange_rate > 0`

## User-facing impact

| Field | What users see |
|---|---|
| `enabled = false` | Top-up entry gone |
| `allow_custom = false` | Only preset chips visible; input field hidden |
| `presets` | Number and order of chips |
| `exchange_rate` | Custom-amount quota credit |

## How to test

1. After enabling, log in as a test account and visit the top-up page
2. Click a preset → credited quota = `bonus_quota`
3. Enter a custom amount → credited quota = `amount × exchange_rate`

## Notes

- Exchange rate changes only affect **new orders**; existing orders keep their snapshotted rate
- Preset changes only affect **new orders**
- Preset amounts must be unique (constraint)

## FAQ

- **All settings saved but user still can't see top-up**: verify `enabled=true` and at least one payment channel enabled
- **Changed exchange rate but no effect**: takes effect for new orders only
- **Can preset amounts be decimals?**: yes — `amount` supports two decimals (e.g. `9.99` ¥)

## Related

- [Top-up (concept)](./topup)
- [Top-up Management (admin)](./topup-management)
- [Payment Channels](./payment)

## Related API

- `GET /api/setting/topup` — read settings (Root)
- `PUT /api/setting/topup` — update settings (Root, **full replacement**)