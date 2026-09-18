---
title: Payment Channels
description: "WeChat / Alipay / Bank-transfer configuration: cert upload, notify URL and enable toggle."
category: admin
order: 12
---

# Payment Channels

> Maintain enable switches, parameters and certificate file paths for the three payment channels on top of `system_settings`. UI: `web/default-pro/src/views/setting/PaymentSetting.vue`.

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/setting/payment` | `GET` | Root | Single response containing all three channels' `{ enabled, config, description, updated_at }` |
| `/api/setting/payment/:method` | `PUT` | Root | Save one channel; accepts `multipart/form-data` for cert uploads |

`:method` is one of `wechat` / `alipay` / `bank`.

Implementation: `controller/setting_payment.go`.

## Payload Shape

PUT requests use `multipart/form-data`:

| Form field | Notes |
|---|---|
| `config` | JSON string: `{ "enabled": bool, "config": { … } }` |
| `cert_file` | wechat only — merchant certificate (`.pem`) |
| `key_file` | wechat only — merchant private key (`.pem`) |
| `private_key_file` | alipay only — app private key |
| `public_key_file` | alipay only — Alipay public key |

Files are saved under `data/payment/<method>/<basename>`. The saved path replaces the `xxx_file` key inside the config map. The PUT response does **not** echo the file path — call GET again to confirm persistence.

## Persistence

Each channel occupies two `system_settings` rows:

| key | Purpose |
|---|---|
| `payment.<method>.enabled` | Minimal JSON containing only `{"enabled": bool}` |
| `payment.<method>.config` | Full configuration JSON (including `enabled` and all keys) |

`category` is fixed to `payment`. The GET endpoint merges the two rows back into a single `{ enabled, config }` object.

## Field Conventions

| Method | Recommended keys |
|---|---|
| `wechat` | `app_id` / `mch_id` / `api_key` / `notify_url` / `cert_file` / `key_file` |
| `alipay` | `app_id` / `gateway` / `notify_url` / `private_key` / `public_key` / `private_key_file` / `public_key_file` |
| `bank` | `account_name` / `account_no` / `bank_name` / `branch` / `notes` |

Keys like `app_id` / `mch_id` / `api_key` / `cert_file` / `key_file` are read directly by the payment SDK — keep the names consistent across backend and frontend.

## Frontend Guide

- Three sections: **WeChat**, **Alipay**, **Bank**.
- Each section's top item is the `enabled` switch; toggling it immediately PUTs the minimal `{ enabled, config }` payload (config is sent along too, to stay idempotent).
- Form fields are only revealed after enabling.
- The **Upload cert** / **Upload private key** buttons (wechat / alipay) use `<a-upload custom-request>` to PUT the current form together with the file.
- The **Save** button submits the entire section; bank has no cert, so just click Save.

## Implementation Pointers

| Concern | Location |
|---|---|
| Handler | `controller/setting_payment.go` |
| Settings CRUD | `model/system_setting.go::GetSystemSetting` / `UpsertSystemSetting` |
| Cert storage | `controller/etting_payment.go::PutPaymentMethod` (`data/payment/<method>/`) |
| Channel implementations | `common/payment/` (`wechat.go`, `alipay.go`, `bank.go`) |
| Public channel status | `controller/payment.go::GetPaymentStatus` |
| Routes | `router/api.go` |