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
| `exchange_rate` | ¥1 = X quota | Drives custom-amount conversion; must be >= `QuotaPerUnit` |
| `presets` | Preset amounts | Each entry is `{amount, bonus_quota}` |

## Recommended setup

Enable + allow custom + a few presets (assuming the default `QuotaPerUnit=500000`, i.e. ¥1 credits ¥1):

| amount | bonus_quota |
|---|---|
| 10 | 5000000 |
| 50 | 25000000 |
| 100 | 50000000 |
| 500 | 250000000 |

`exchange_rate = 500000` (¥1 = `QuotaPerUnit` quota, so ¥1 in credits ¥1).

> `QuotaPerUnit` is the system-wide conversion base (default 500000) and can be changed in system settings.
> `exchange_rate` must not be lower than it, otherwise ¥1 would credit less than ¥1.
> To grant a bonus, use a multiplier (e.g. `1000000` = 2x).

## How to change

Admin → Top-up Settings → edit fields → Save.

Validation:

- Each preset `amount > 0`, `bonus_quota >= 0`
- No duplicate amounts
- `exchange_rate >= QuotaPerUnit` (legacy sub-base values are normalized on read)

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
- Legacy sub-base `exchange_rate` values (e.g. `1` / `100000`) are normalized to `QuotaPerUnit` on read — no data migration needed
- Preset changes only affect **new orders**
- Preset amounts must be unique (constraint)

## FAQ

- **All settings saved but user still can't see top-up**: verify `enabled=true` and at least one payment channel enabled
- **Changed exchange rate but no effect**: takes effect for new orders only
- **Why does saving fail with "兑换比例不能低于 500000"?**: an `exchange_rate` below the system `QuotaPerUnit` is rejected; use the base value or a higher bonus multiplier
- **Can preset amounts be decimals?**: yes — `amount` supports two decimals (e.g. `9.99` ¥)

## Related

- [Top-up (concept)](./topup)
- [Top-up Management (admin)](./topup-management)
- [Payment Channels](./payment)

## Related API

- `GET /api/setting/topup` — read settings (Root)
- `PUT /api/setting/topup` — update settings (Root, **full replacement**)