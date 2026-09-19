---
title: Topup
description: Online top-up of user quota, with presets, custom amount, exchange rate and crediting.
category: schema
order: 10
---

# Topup

## What it is

`Topup` is the user's online quota purchase. It is not tied to a subscription or plan — it simply adds credit to [User.quota](/en/schema/user).

Order number prefix: `TP`.

## Where to find it

- **Public → Top-up**: presets, custom amount input, current rate
- **Admin → Top-ups**: admin's order list
- **Admin → Top-up Settings**: master switch, preset list, rate, custom-amount policy

## What you set to enable top-up

Admin → Top-up Settings:

| Field | Meaning | Effect |
|---|---|---|
| Master switch | Enabled / not | Off hides the public entry. |
| Preset amounts | e.g. 10 / 50 / 100 | Quick-pick options in the UI. |
| Custom amount | Toggle | Off forces users to pick a preset. |
| Exchange rate | Quota per ¥ | A ¥10 top-up adds `10 × rate` to the user's quota. |
| Min / max amount | Number | Bounds on what users can enter. |

A sane starter:

- Enabled
- Presets: 10 / 50 / 100 / 500
- Custom allowed, min ¥1, max ¥10 000
- Rate `100000` (matches default `QuotaPerUnit=500000`)

## Operator-relevant fields (per order)

| Field | Meaning | Effect |
|---|---|---|
| Order number | `TP` prefix | Changing it breaks the callback. |
| Parent | User FK | Cascade-delete with user. |
| Amount | ¥ | Paid amount. |
| Quota | amount × rate | Credited on payment. |
| Status | Unpaid / Paid / Cancelled | Drives crediting. |
| Pay method | WeChat / Alipay / Bank / Mock | Which callback URL. |

## Crediting

On "Paid", `User.quota` is incremented by `amount × rate` immediately. No manual intervention.

## Top-up vs Subscriptions

- Top-up → credited to `User.quota`, never expires, no plan discount
- Subscription → drives the plan discount for a fixed period

A user can hold both: top-up balance (no discount) + active subscription (discount applied first).

## Related pages

- [Topup (concept)](/en/pricing/topup)
- [Topup Management (admin)](/en/pricing/topup-management)
- [Topup Settings](/en/pricing/topup-settings)

## Related API

- `POST /api/topup/order` — user creates top-up order
- `GET /api/setting/topup` — get settings
- `PUT /api/setting/topup` — update settings (Root)