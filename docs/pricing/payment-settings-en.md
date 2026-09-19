---
title: Payment Channel Settings
description: Enable and configure each payment channel (WeChat / Alipay / Bank Transfer) in the admin console.
category: pricing
order: 10
---

# Payment Channel Settings

> Admin → Payment Settings. Visible only to Root.

## Three channels

Each channel has its own toggle and credentials. **Enable at least one** or users can't pay.

### WeChat Pay

| Field | What | Where |
|---|---|---|
| `app_id` | AppID | WeChat Pay merchant console → Account center |
| `mch_id` | Merchant ID | Same |
| `api_key` | v2 key | Merchant console → API security → APIv2 key |
| `notify_url` | Callback URL | `https://your-domain.com/api/payment/wechat/notify` |
| Cert / Private key | PEM files (for refunds) | Merchant console → API security → API certificates |

### Alipay Face-to-Face

| Field | What | Where |
|---|---|---|
| `app_id` | App ID | Alipay Open Platform → My apps |
| `private_key` / `public_key` | App private key / Alipay public key | Same |
| or `private_key_file` / `public_key_file` | Paths to the above PEM files | Auto-saved on upload |
| `gateway` | Gateway | Default `https://openapi.alipay.com/gateway.do` (production) |
| `notify_url` | Callback URL | `https://your-domain.com/api/payment/alipay/notify` |

### Bank transfer

| Field | What |
|---|---|
| `account_name` | Beneficiary name (company name) |
| `account_no` | Bank account number |
| `bank_name` | Bank |
| `branch` | Branch |
| `notes` | Prompt shown to users to include the order number in remarks |

## How to enable / configure

Admin → Payment Settings:

1. Pick a channel → turn on "Enabled"
2. Form fields expand → fill in
3. For WeChat / Alipay, upload the PEM file(s)
4. Click Save

The toggle is effective immediately.

## Certificate / file upload

PEM files for WeChat and Alipay:

- Saved to `data/payment/<method>/<basename>` on disk
- The path is written back into the `config` field (e.g. `cert_file: /app/data/payment/wechat/apiclient_cert.pem`)
- **Not returned in the response** — you need to GET again to confirm

## Example configs

### WeChat production

```json
{
  "enabled": true,
  "config": {
    "app_id": "wx0123456789abcdef",
    "mch_id": "1900000001",
    "api_key": "your-strong-api-key",
    "notify_url": "https://api.example.com/api/payment/wechat/notify"
  }
}
```

### Alipay production

```json
{
  "enabled": true,
  "config": {
    "app_id": "2021000123456789",
    "gateway": "https://openapi.alipay.com/gateway.do",
    "notify_url": "https://api.example.com/api/payment/alipay/notify"
  }
}
```

For private key, upload the PEM file rather than pasting.

### Bank transfer

```json
{
  "enabled": true,
  "config": {
    "account_name": "XX Tech Ltd.",
    "account_no": "6225 1234 5678 9012",
    "bank_name": "CMB",
    "branch": "Shanghai Branch",
    "notes": "Please include the order number in your transfer remarks."
  }
}
```

## How to verify

1. Admin → Payment Settings → confirm toggle is on
2. Public endpoint `/api/payment/status` returns the enabled channels (no login)
3. Place a test order → user can pick the channel → mock the payment to verify the full chain

## Notes

- Changes take effect immediately, but don't affect orders already created
- Credentials are **sensitive**: don't commit them to git / share unnecessarily
- Switching channels: existing orders stay on the old channel; new orders use the new channel

## FAQ

- **Users see no payment methods**: no channel enabled, or all are mis-configured
- **Payment succeeded but order stays unpaid**: verify `notify_url` is publicly reachable; check backend logs
- **Certificate upload fails**: check PEM format; private key must not have a passphrase
- **Switched channel — what about old orders**: existing orders use the old channel; new orders use the new one

## Related

- [Payment Channels](./payment)
- [Top-up Settings](./topup-settings)
- [Order Management](./order-management)

## Related API

- `GET /api/setting/payment` — read all channel configs (Root)
- `PUT /api/setting/payment/:method` — save one channel (Root, supports file upload)
- `GET /api/payment/status` — public enabled-channels query (no login)