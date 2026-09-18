---
title: Payment API
description: "/api/payment/* endpoints and callback notifications."
category: api
order: 10
---
# Payment API

The payment endpoints cover async callbacks for **WeChat Pay**, **Alipay**, and **bank transfer**, plus status queries, an admin Mock channel, and payment-channel configuration. Callback endpoints are public; authenticity is enforced via the channel signature.

## Endpoint index

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/api/payment/wechat/notify` | POST | Public | WeChat Pay async notification |
| `/api/payment/alipay/notify` | POST | Public | Alipay async notification |
| `/api/payment/mock/notify` | POST | Root | Admin Mock notification (test / manual confirmation) |
| `/api/payment/status` | GET | Public | Enabled status of every payment channel |
| `/api/setting/payment` | GET | Root | Read payment configuration |
| `/api/setting/payment/:method` | PUT | Root | Update payment configuration (`wechat` / `alipay` / `bank`) |


## Common: Async Notification Flow

All async callbacks (`wechat` / `alipay`) share the same handler (`controller/payment.go::processNotify`):

1. Read the raw request body (no JSON parsing).
2. Call the corresponding `payment.Channel.VerifyNotify` to verify the signature.
3. On success, extract `NotifyResult{OutTradeNo, TradeNo, Amount, Paid}`.
4. If `order.amount > 0` and `notify.amount > 0`, they must match.
5. Dispatch by order type:
   - `type = 1` (plan subscription) → `model.ActivatePackageByOrder`
   - `type = 2` (topup) → `model.ActivateTopupByOrder`
6. Channel-specific success payload:
   - WeChat: `<xml>` with `return_code=SUCCESS`
   - Alipay: literal `success`

> Activation is idempotent: re-delivering a paid-order notification will not double-credit quota or re-activate the plan.


## 1. WeChat Pay Callback

**Endpoint:** `POST /api/payment/wechat/notify`

**Auth:** Public (relies on WeChat signature verification)

**Request body:** WeChat's native form-urlencoded callback. Verified by `common/payment/wechat.go::VerifyNotify`.

**Response:**

| Scenario | Content-Type | Body |
|----------|--------------|------|
| Signature OK | `application/xml` | `<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>` |
| Signature / business error | `application/xml` | `<xml><return_code><![CDATA[FAIL]]></return_code><return_msg><![CDATA[<err>]]></return_msg></xml>` |

> Per the WeChat spec, returning `FAIL` on a business error lets WeChat retry the notification.


## 2. Alipay Callback

**Endpoint:** `POST /api/payment/alipay/notify`

**Auth:** Public (relies on Alipay RSA2 signature verification)

**Request body:** Alipay's form-urlencoded callback. Verified by `common/payment/alipay.go::VerifyNotify`.

**Response:**

| Scenario | Body |
|----------|------|
| Signature OK | literal `success` |
| Signature / business error | literal `fail` |


## 3. Mock Notification (admin)

**Endpoint:** `POST /api/payment/mock/notify`

**Auth:** Root

**Description:** Mark an order paid or refunded without going through a real payment channel. The handler dispatches by order type to `ActivateTopupByOrder` / `ActivatePackageByOrder` (default upgrade mode `stack`).

**Request body:**

```json
{
  "order_no": "TB20250912153000123456",
  "status": 1
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| order_no | string | yes | Target order number |
| status | int | yes | `1` = mark paid; `3` = mark refunded |

**Response:**

```json
{ "success": true, "message": "订单已支付，余额已到账" }
```

or

```json
{ "success": true, "message": "订单已支付，套餐已激活" }
```

or

```json
{ "success": true, "message": "订单已标记为退款" }
```

**Errors:**

| Scenario | message |
|----------|---------|
| `order_no` empty | `order_no 不能为空` |
| Order not found | `订单不存在` |
| `status` not 1 or 3 | `不支持的状态值（仅支持 1 或 3）` |
| Activation failed | `激活失败: <err>` |


## 4. Payment Channel Status

**Endpoint:** `GET /api/payment/status`

**Auth:** Public

**Description:** Reads the enabled status of every registered payment channel so the user-facing purchase page can decide whether to show the order UI.

**Response:**

```json
{
  "success": true,
  "message": "",
  "data": {
    "any_enabled": true,
    "methods": [
      { "name": "wechat", "label": "WeChat Pay", "enabled": true },
      { "name": "alipay", "label": "Alipay",     "enabled": true },
      { "name": "bank",   "label": "Bank Transfer", "enabled": false }
    ]
  }
}
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| any_enabled | bool | At least one channel is enabled |
| methods[].name | string | Channel id (`wechat` / `alipay` / `bank`) |
| methods[].label | string | Human-readable label |
| methods[].enabled | bool | Whether this channel is enabled |


## 5. Read Payment Configuration

**Endpoint:** `GET /api/setting/payment`

**Auth:** Root

**Response:**

```json
{
  "success": true,
  "message": "",
  "data": {
    "wechat_enabled":  { "enabled": true,  "config": { ... }, "description": "WeChat Pay toggle",   "updated_at": 1718000000 },
    "wechat_config":   { "enabled": true,  "config": { ... }, "description": "WeChat Pay settings","updated_at": 1718000000 },
    "alipay_enabled":  { "enabled": false, "config": {},      "description": "Alipay toggle",      "updated_at": 0 },
    "alipay_config":   { "enabled": false, "config": {},      "description": "Alipay settings",    "updated_at": 0 },
    "bank_enabled":    { "enabled": false, "config": {},      "description": "Bank transfer toggle","updated_at": 0 },
    "bank_config":     { "enabled": false, "config": {},      "description": "Bank account info",  "updated_at": 0 }
  }
}
```

> Each key corresponds to one row in `system_settings` (`payment.<method>.<enabled>` / `payment.<method>.config`).


## 6. Update Payment Configuration

**Endpoint:** `PUT /api/setting/payment/:method`

**Auth:** Root

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| method | string | `wechat` / `alipay` / `bank` |

**Request body:** Supports `application/json` or `multipart/form-data` (use multipart when uploading certificates).

JSON form:

```json
{
  "enabled": true,
  "config": {
    "app_id": "wx1234567890",
    "mch_id": "1900000000",
    "sign_type": "HMAC-SHA256"
  }
}
```

Multipart form:

| Field | Required | Description |
|-------|----------|-------------|
| `config` | yes | JSON string with the same shape as above |
| `cert_file` | WeChat optional | WeChat Pay API certificate; saved path written back to `config.cert_file` |
| `key_file` | WeChat optional | WeChat Pay API private key; saved path written back to `config.key_file` |
| `private_key_file` | Alipay optional | Merchant RSA private key; saved path written back to `config.private_key_file` |
| `public_key_file` | Alipay optional | Alipay public key; saved path written back to `config.public_key_file` |

> Uploaded certificates / keys are saved under `data/payment/<method>/<basename>`, and the absolute path is written back to the matching `xxx_file` field in `config`.

**Response:**

```json
{ "success": true, "message": "已保存" }
```

**Errors:**

| Scenario | message |
|----------|---------|
| Unknown `method` | `未知支付方式: <method>` |
| multipart missing `config` | `multipart 缺少 config 字段` |
| Invalid config JSON | `config 不是合法 JSON` / `解析配置失败: <err>` |
| File save failure | `创建目录失败: <err>` / `保存文件失败: <err>` |
| DB write failure | `<err>` |