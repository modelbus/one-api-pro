---
title: Payment Channels
description: How to configure WeChat, Alipay, and bank transfer.
category: pricing
order: 5
---

# Payment Channels

> Which channels can users pay through? How to configure them? How are callbacks handled?

## Supported channels

| Channel | pay_method | Use case |
|---|---|---|
| **WeChat Native** | `wechat` | Users scan a WeChat QR |
| **Alipay Face-to-Face** | `alipay` | Users scan an Alipay QR |
| **Bank transfer** | `bank` | Users wire from a company account; admin confirms manually |
| **Offline** | `offline` | No online flow; admin handles manually |
| **Free / Grant** | `free` | For admin-granted subscriptions |

Enable at least one of the first three, otherwise users won't see any pay button at checkout.

## Where

- **Admin → Payment Settings**: configure each channel
- **User side**: pay-method dropdown at checkout / top-up

## WeChat Native (good for B2C)

### Apply

1. WeChat Pay merchant console → Products → apply for **Native Pay**
2. Get `app_id` / `mch_id` / `api_key` (v2) or "Merchant API Certificate" (v3)
3. Set the callback URL: `https://your-domain.com/api/payment/wechat/notify`

### Admin config

Admin → Payment Settings → WeChat:

| Field | Value |
|---|---|
| Enabled | toggle |
| `app_id` | Merchant console → Account center |
| `mch_id` | Merchant console → Account center |
| `api_key` | Merchant console → API security → v2 key |
| `notify_url` | `https://your-domain.com/api/payment/wechat/notify` |
| Cert / Private key (PEM) | For refunds; optional |

## Alipay Face-to-Face (good for B2C)

### Apply

1. Alipay Open Platform → create app → sign "Face-to-Face"
2. Get `app_id` / public key / private key
3. Set the callback URL: `https://your-domain.com/api/payment/alipay/notify`

### Admin config

Admin → Payment Settings → Alipay:

| Field | Value |
|---|---|
| Enabled | toggle |
| `app_id` | Open Platform → My apps |
| `private_key` / `public_key` | App public/private key (PEM content) |
| or `private_key_file` / `public_key_file` | File paths |
| `gateway` | Default `https://openapi.alipay.com/gateway.do` (production) |
| `notify_url` | `https://your-domain.com/api/payment/alipay/notify` |

## Bank transfer (good for B2B)

Best for: large amounts, reconciliation requirements, corporate customers.

### Flow

1. User places an order → status = unpaid
2. User wires money to your account (with order number in remarks)
3. Admin goes to Admin → Orders → "Mark as paid"
4. System activates the plan / credits the top-up

### Admin config

Admin → Payment Settings → Bank transfer:

| Field | Value |
|---|---|
| Enabled | toggle |
| `account_name` | Beneficiary name (company name) |
| `account_no` | Bank account number |
| `bank_name` | Bank |
| `branch` | Branch |
| `notes` | Ask users to include the order number in remarks |

## User-side checkout flow

```
User clicks Subscribe / Top-up
    ↓
Popup shows available payment methods (enabled ones only)
    ↓
User picks WeChat → QR code → scan to pay
User picks Alipay → QR code → scan to pay
User picks Bank transfer → show account + ask to include order number
    ↓
WeChat/Alipay callback → order = paid → activate plan
Bank transfer → admin manually marks paid → activate plan
```

## Async callbacks

After payment, the provider posts to:

- `POST /api/payment/wechat/notify` (WeChat)
- `POST /api/payment/alipay/notify` (Alipay)

These endpoints are **unauthenticated** (called by the payment platform). Make sure they're reachable from the public internet.

## Manual activation (debug / reconciliation)

`POST /api/payment/mock/notify` (Root): pass `{order_no, status}`:

- `status=1` → force activate (credits quota or creates subscription per the order type)
- `status=3` → mark as refunded

Use when: testing, fixing reconciliation, recovering from missed callbacks.

## FAQ

- **No payment buttons show**: at least one channel must be enabled
- **Callback never arrives**: check that `notify_url` is publicly reachable; check the reverse proxy / firewall
- **Amount mismatch error**: make sure the order `amount` matches what the channel returns — don't apply FX or discounts that change the amount
- **Bank transfer received — what now?**: admin goes to Admin → Orders → "Mark paid"

## Related

- [Payment Settings (admin)](./payment-settings)
- [Order Management (admin)](./order-management)
- [My Orders (user)](../en/user/orders)

## Related API

- `GET /api/payment/status` — currently enabled methods (Public)
- `POST /api/payment/wechat/notify` — WeChat callback
- `POST /api/payment/alipay/notify` — Alipay callback
- `POST /api/payment/mock/notify` — manual activation (Root)
- `GET /api/setting/payment` — read payment config (Root)
- `PUT /api/setting/payment/:method` — update one channel (Root)