---
title: Payment Channels
description: "WeChat Native, Alipay Face-to-Face, bank transfer, offline, and free channels; certificate upload and notify URL configuration."
category: pricing
order: 5
---

# Payment Channels

> Payment channels implement `common/payment.Channel`. Each is registered against a `pay_method` and reads its `enabled` / `config` JSON from `system_settings`. Callers obtain an instance via `payment.New(pay_method)`.

## Registered channels

`common/payment/payment.go::AnyChannelEnabled` enumerates `wechat` / `alipay` / `bank` (offline / free are always considered enabled). Registration happens in each file's `init()`:

| Channel | pay_method | Implementation | Credentials |
|---|---|---|---|
| WeChat Native | `OrderPayMethodWechat="wechat"` | `common/payment/wechat.go` | `app_id` / `mch_id` / `api_key` / `notify_url`; optional `cert_file` / `key_file` (refund) |
| Alipay Face-to-Face | `OrderPayMethodAlipay="alipay"` | `common/payment/alipay.go` | `app_id` / `private_key` / `public_key` / `private_key_file` / `public_key_file` / `notify_url` / `gateway` |
| Bank transfer | `OrderPayMethodBank="bank"` | `common/payment/bank.go::bankChannel` | `account_name` / `account_no` / `bank_name` / `branch` / `notes` (`PrePay` not implemented; admin marks received) |
| Offline | `OrderPayMethodOffline="offline"` | `common/payment/bank.go::offlineChannel` | Same as bank; `IsEnabled()` always returns true |
| Free / admin grant | `OrderPayMethodFree="free"` | `common/payment/bank.go::freeChannel` | Always enabled; `PrePay` returns an empty object |

## Settings

Payment config lives in the `system_settings` table under `payment.*` (`model/system_setting.go`):

| Key | Meaning |
|---|---|
| `payment.wechat.enabled` / `payment.wechat.config` | WeChat master switch / credential JSON |
| `payment.alipay.enabled` / `payment.alipay.config` | Alipay master switch / credential JSON |
| `payment.bank.enabled` / `payment.bank.config` | Bank master switch / credential JSON |

`IsEnabled()` calls `payment.SettingsBool(key)` which parses `{"enabled": true|false}`.

## WeChat Native

`wechatChannel.PrePay`:

1. Reject when `IsEnabled()` is false
2. Load `wechatConfig`: `app_id` / `mch_id` / `api_key` are required; otherwise `微信支付参数不完整`
3. Call `UnifiedOrder(ctx, BodyMap{body, out_trade_no, total_fee, spbill_create_ip, notify_url, trade_type=Native})`
4. `total_fee` = yuan × 100 (`int64(amount*100 + 0.5)`)
5. Returns `PrePayResult{PayURL=CodeURL, QRCode=CodeURL, ExpireAt=0, TradeNo=PrepayId}`

`VerifyNotify`:

- Parses the XML into `wechatNotifyXML` (`return_code` / `result_code` / `out_trade_no` / `transaction_id` / `total_fee` / `sign`)
- Builds the canonical string by sorting non-empty fields in ASCII order and appending `&key=<API_KEY>`, MD5, compared case-insensitively against `sign`
- Validates `return_code == SUCCESS` and `result_code == SUCCESS`
- Amount = `TotalFee / 100` (yuan)

`WechatNotify` endpoint (`POST /api/payment/wechat/notify`, public):

- Success → `<xml><return_code>SUCCESS</return_code>...</xml>`
- Failure → `<xml><return_code>FAIL</return_code><return_msg>...</return_msg></xml>` (WeChat retries)

## Alipay Face-to-Face

`alipayChannel.PrePay`:

- `app_id` required; prefers `private_key_file` / `public_key_file`, otherwise reads inline `private_key` / `public_key`
- `gateway` defaults to `https://openapi.alipay.com/gateway.do` (production); can be set to the sandbox URL
- Calls `TradePrecreate(ctx, "当面付", out_trade_no, total_amount)` to get the QR string

`VerifyNotify` uses `alipay.VerifySign` (RSA2 with Alipay's public key) over the POST fields and returns `NotifyResult{OutTradeNo, TradeNo, Amount, Paid}`.

`AlipayNotify` endpoint (`POST /api/payment/alipay/notify`, public) returns the literal `success` on success and `fail` on error.

## Bank / Offline / Free

- `bankChannel`: `PrePay` returns `bank 支付未实现：等待管理员在后台标记收款`; no async notify; orders stay `pending` until an admin manually marks them paid via `PUT /api/order/:id`
- `offlineChannel`: same as bank, but `IsEnabled()` always returns true (always available)
- `freeChannel`: `PrePay` returns an empty object; `VerifyNotify` returns `Paid: true`. Used only by admin grants (`controller/subscription.go::AddSubscription`)

## User flow

`payment.AnyChannelEnabled()` decides whether `CreatePlanOrder` / `CreateTopupOrder` / `PayMyOrder` are allowed; on `false` they all return `系统尚未开通任何支付通道，请设置后开启支付` and the order is not persisted.

`buildPayInfo(pay_method, order_no, amount, "TBUS-"+package_name)` returns:

```json
{
  "status": "success",          // or "warning"
  "pay_url": "weixin://wxpay/bizpayurl?pr=...",
  "qr_code": "weixin://wxpay/bizpayurl?pr=...",
  "expire_at": 0,
  "trade_no": "prepay_id_xxx",
  "note": "..."                 // bank only
  "warning": "..."              // warning only
}
```

`status=warning` indicates the channel is not registered, disabled, or the SDK call failed; the order is still persisted so an admin can handle it manually.

## Async notify dispatch

`controller/payment.go::processNotify`:

1. `payment.New(pay_method).VerifyNotify(body)` validates the signature
2. `model.GetOrderByOrderNo(notif.OutTradeNo)` fetches the order
3. Amount mismatch: `order.Amount > 0 && notif.Amount > 0 && notif.Amount != order.Amount` → `amount mismatch`
4. Dispatch by `order.Type`:
   - `OrderTypeTopup=2` → `model.ActivateTopupByOrder` adds quota
   - Otherwise (plan order) → `model.ActivatePackageByOrder(order, OrderUpgradeModeStack)` activates the subscription

## Manual activation / test channel

`POST /api/payment/mock/notify` (Root, body `{order_no, status}`):

- `status=1` → dispatch by order type to `ActivateTopupByOrder` or `ActivatePackageByOrder`
- `status=3` → `model.MarkOrderRefunded` (status flip only; `users.quota` / `user_plans` untouched)

Used for tests or manual activation after reconciliation.

## Frontend Guide

Route: `/setting/payment` (`web/default-pro/src/views/setting/PaymentSetting.vue`).

- Each channel has a switch and a config form:
  - **WeChat**: `app_id` / `mch_id` / `api_key` / `notify_url` + PEM upload buttons for cert / key
  - **Alipay**: `app_id` / `gateway` / `notify_url` + PEM upload buttons for public / private key
  - **Bank**: `account_name` / `account_no` / `bank_name` / `notes`
- The switch fires `PUT /api/setting/payment/:method` immediately; form fields go through the same endpoint via the `config` field
- User-side `/pricing` or `/topup` modal: `GET /api/payment/status` decides which buttons to render (only enabled ones)

## Implementation Pointers

| Concern | Location |
|---|---|
| Interface | `common/payment/payment.go::Channel` |
| WeChat impl | `common/payment/wechat.go` |
| Alipay impl | `common/payment/alipay.go` |
| Bank / Offline / Free | `common/payment/bank.go` |
| Setting keys | `model/system_setting.go::SystemSettingKeyWechat* / Alipay* / Bank*` |
| Settings endpoints | `controller/setting_payment.go` (see frontend fields) |
| Notify dispatch | `controller/payment.go::processNotify` |
| Manual activation | `controller/payment.go::MockPay` |
