---
title: Top-up API
description: "/api/topup/* endpoints."
category: api
order: 9
---
# Top-up API

`/api/topup/*` exposes **online top-up** (balance credit) endpoints: users create their own orders, and admins maintain top-up settings (master switch, custom amount, quick amounts). Order type is always `2` (`OrderTypeTopup`); async notifications go to `/api/payment/*`.

## Endpoints

| Endpoint | Method | Auth | Description |
|------|------|------|------|
| `/api/topup/order` | POST | User | Create a top-up order as a user |
| `/api/setting/topup` | GET | Root | Read top-up settings |
| `/api/setting/topup` | PUT | Root | Save top-up settings |

## Conventions

- **The quota base is a system constant**: ¥1 = `QuotaPerUnit` = `1,000,000` quota (1 quota = 1e-6 CNY), conceptually like WeChat Pay's "cent". It takes part in computation only, is never stored in the DB and cannot be modified; `/api/status` returns it read-only as `quota_per_unit` for frontend rendering.
- Top-up order-number prefix: `TP` (see `model.GenerateOrderNo("TP")`).
- An order is only accepted when at least one payment channel is enabled (`payment.AnyChannelEnabled().any_enabled == true`) and `topup.enabled == true`.
- Amount resolution (see `model.ResolveTopupAmount`):
  - `preset_amount > 0`: match a preset by its `amount` and use its `bonus_quota`.
  - Otherwise use the custom `amount`, which requires `topup.allow_custom == true`; it is **always 1:1**, `bonus_quota = round(amount × QuotaPerUnit)`.
- **No exchange rate any more**: `topup.exchange_rate` has been removed, so a custom amount credits exactly what you pay. Promotions (e.g. "pay 10 get 15") are expressed by the pay/credit columns of quick amounts.
- Activation goes through `model.ActivateTopupByOrder` (triggered by the async callback): it calls `IncreaseUserQuota(bonus_quota)` and marks the order `status=1`. The order snapshot stores `credit_amount` (credited CNY) and `bonus_quota` (credited quota).

## 1. Create a top-up order

**Endpoint:** `POST /api/topup/order`

**Auth:** User

**Behaviour:** validate the amount → `model.CreateTopupOrder` persists the order (type=2, status=0) → `controller.buildPayInfo` builds the payment info.

**Preconditions:**

- Signed in (`c.GetInt("id") != 0`)
- At least one payment channel enabled
- `topup.enabled == true`

**Request body:**

| Field | Type | Required | Description |
|------|------|------|------|
| amount | float64 | one of two | Custom amount (CNY); ignored when `preset_amount > 0`. Credits = `round(amount × 1000000)` |
| preset_amount | float64 | one of two | Preset amount (CNY); uses the matched preset's `bonus_quota` |
| pay_method | string | yes | `wechat` / `alipay` / `bank` |

```json
{
  "preset_amount": 10,
  "pay_method": "wechat"
}
```

**Response example:**

```json
{
  "success": true,
  "message": "",
  "order": {
    "id": 1,
    "order_no": "TP20261003120000001",
    "user_id": 1,
    "type": 2,
    "amount": 10,
    "status": 0,
    "pay_status": 0,
    "pay_method": "wechat",
    "plan_info": "{\"amount\":10,\"preset_amount\":10,\"credit_amount\":15,\"bonus_quota\":15000000}"
  },
  "amount": 10,
  "bonus_quota": 15000000,
  "pay": {
    "status": "success",
    "pay_url": "weixin://wxpay/bizpayurl?pr=xxxx",
    "qr_code": "weixin://wxpay/bizpayurl?pr=xxxx",
    "trade_no": "4200001234202610031234567890"
  }
}
```

> The example above is a "pay 10 get 15" preset: ¥10 paid, ¥15 credited (`bonus_quota = 15000000`).

**`pay` fields:**

| Field | Type | Description |
|------|------|------|
| status | string | `success` (payable) / `warning` (channel disabled or pre-order failed) |
| pay_url | string | Payment redirect URL (WeChat/Alipay) |
| qr_code | string | QR payload (same as `pay_url` for Native channels) |
| expire_at | int64 | Expiry timestamp; 0 means unknown |
| trade_no | string | Upstream pre-order id (not equal to `order_no`) |
| note | string | Returned for the `bank` channel only: transfer instructions |
| warning | string | Returned only when `status="warning"`: failure reason |

**Errors:**

| Scenario | message |
|------|---------|
| Not signed in | `未登录` |
| JSON parse failure | `无效的参数` |
| `amount <= 0 && preset_amount <= 0` | `amount 或 preset_amount 至少传一个` |
| Empty `pay_method` | `pay_method 不能为空` |
| Unsupported pay method | `不支持的支付方式` |
| Pay method not `wechat/alipay/bank` | `自助充值仅支持 wechat / alipay / bank` |
| No payment channel configured | `系统尚未开通任何支付通道，请设置后开启支付` |
| Top-up disabled | `充值功能未开启` |
| Preset not configured | `快捷金额未配置` |
| Custom amount disabled | `未开启自定义金额` |
| Amount must be > 0 | `充值金额必须大于 0` |
| Order-number generation failed | `生成订单号失败，请重试` |

## 2. Read top-up settings

**Endpoint:** `GET /api/setting/topup`

**Auth:** Root

**Response example:**

```json
{
  "success": true,
  "message": "",
  "data": {
    "enabled": true,
    "allow_custom": true,
    "presets": [
      { "amount": 10, "bonus_quota": 15000000 },
      { "amount": 50, "bonus_quota": 50000000 }
    ]
  }
}
```

**Fields:**

| Field | Type | Description |
|------|------|------|
| enabled | bool | `topup.enabled`, master switch |
| allow_custom | bool | `topup.allow_custom`, whether custom amounts are allowed (always 1:1) |
| presets | array | `topup.presets` (deserialized JSON), each item is a `TopupPreset` |

`TopupPreset`:

| Field | Type | Description |
|------|------|------|
| amount | float64 | Pay amount (CNY) |
| bonus_quota | int64 | Credited quota; must be >= `round(amount × 1000000)` |

## 3. Save top-up settings

**Endpoint:** `PUT /api/setting/topup`

**Auth:** Root

**Request body:**

```json
{
  "enabled": true,
  "allow_custom": true,
  "presets": [
    { "amount": 10, "bonus_quota": 15000000 },
    { "amount": 50, "bonus_quota": 50000000 }
  ]
}
```

**Fields:**

| Field | Type | Required | Description |
|------|------|------|------|
| enabled | bool | no | Master switch, default false |
| allow_custom | bool | no | Allow custom amounts, default false |
| presets | array | no | Quick amounts; the whole list is replaced |

**Validation** (`model.SaveTopupSettings`):

- `presets[i].amount > 0`
- `amount` must not repeat within `presets`
- `presets[i].bonus_quota >= round(presets[i].amount × 1000000)` (bonuses allowed, under-crediting not)

**Response:**

```json
{ "success": true, "message": "已保存" }
```

**Errors:**

| Scenario | message |
|------|---------|
| JSON parse failure | `无效的参数` |
| Row N amount <= 0 | `第 N 项金额必须大于 0` |
| Duplicate amount X | `快捷金额重复：X.XX 元已存在` |
| Credit below pay | `第 N 项到账额度不能低于支付金额折算额度（X.XX 元 = Y 额度）` |
