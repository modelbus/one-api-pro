---
title: Topup API
description: "/api/topup/* endpoints."
category: api
order: 9
---
# Topup API

`/api/topup/*` provides online topup (balance recharge) endpoints: user self-service ordering, payment-parameter generation, and admin maintenance of topup settings (master switch, custom amount, presets, exchange rate). Order type is fixed to `2` (`OrderTypeTopup`); payment notifications still go through `/api/payment/*`.

## Endpoint index

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/api/topup/order` | POST | User | Create a topup order (self-service) |
| `/api/setting/topup` | GET | Root | Read topup settings |
| `/api/setting/topup` | PUT | Root | Save topup settings |


## Common Conventions

- Order number prefix: `TP` (see `model.GenerateOrderNo("TP")`).
- Ordering is only allowed when at least one payment channel is enabled (`payment.AnyChannelEnabled().any_enabled == true`) AND `topup.enabled == true`; otherwise the request is rejected.
- Amount resolution (`model.ResolveTopupAmount`):
  - `preset_amount > 0`: hit a preset; locate the preset whose `amount` matches and use its `bonus_quota`.
  - Otherwise use the custom `amount`, which requires `topup.allow_custom == true` and yields `bonus_quota = amount × exchange_rate`.
- Default exchange rate `1:1` (1 CNY = 1 quota), only applied to the custom amount path.
- Activation goes through `model.ActivateTopupByOrder` (triggered by the async callback): calls `IncreaseUserQuota(bonus_quota)` and sets the order to `status=1`.


## 1. Create a Topup Order

**Endpoint:** `POST /api/topup/order`

**Auth:** User

**Description:** Validate the body → call `model.CreateTopupOrder` to persist the order (type=2, status=0) → call `controller.buildPayInfo` to produce payment parameters.

**Prerequisites:**

- Logged in (`c.GetInt("id") != 0`).
- At least one payment channel is enabled.
- `topup.enabled == true`.

**Request body:**

```json
{
  "amount": 100.0,
  "preset_amount": 0,
  "pay_method": "wechat"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| amount | float64 | one of | Custom amount (CNY). Ignored when `preset_amount > 0`. |
| preset_amount | float64 | one of | Hits a preset (CNY). When > 0, the matched preset's `bonus_quota` is used. |
| pay_method | string | yes | `wechat` / `alipay` / `bank` |

**Response:**

```json
{
  "success": true,
  "message": "",
  "order": {
    "id": 12,
    "type": 2,
    "source": 1,
    "order_no": "TP20250912153000123456",
    "user_id": 7,
    "plan_id": 0,
    "plan_info": "{\"amount\":100,\"bonus_quota\":100000,\"exchange_rate\":1000}",
    "amount": 100.00,
    "status": 0,
    "pay_status": 0,
    "pay_method": "wechat",
    "create_time": 1718000000,
    "update_time": 1718000000
  },
  "amount": 100.00,
  "bonus_quota": 100000,
  "pay": {
    "status": "success",
    "pay_url": "weixin://wxpay/bizpayurl?pr=xxxxxxxx",
    "qr_code": "weixin://wxpay/bizpayurl?pr=xxxxxxxx",
    "expire_at": 1718003600,
    "trade_no": "PFX20250912153000123456"
  }
}
```

**`pay` fields:**

| Field | Type | Description |
|-------|------|-------------|
| status | string | `success` (ready to pay) / `warning` (channel not enabled or pre-pay failed) |
| pay_url | string | Payment redirect URL (WeChat / Alipay) |
| qr_code | string | QR-code payload (Native channels share this with `pay_url`) |
| expire_at | int64 | Expiration timestamp, 0 if unknown |
| trade_no | string | Channel pre-payment id (different from `order_no`) |
| note | string | Only present for the `bank` channel; transfer instructions |
| warning | string | Only when `status="warning"`; explains why pre-pay failed |

**Errors:**

| Scenario | message |
|----------|---------|
| Not logged in | `未登录` |
| Invalid JSON | `无效的参数` |
| `amount <= 0 && preset_amount <= 0` | `amount 或 preset_amount 至少传一个` |
| Empty `pay_method` | `pay_method 不能为空` |
| Unsupported pay method | `不支持的支付方式` |
| Pay method not `wechat/alipay/bank` | `自助充值仅支持 wechat / alipay / bank` |
| No payment channel enabled | `系统尚未开通任何支付通道，请设置后开启支付` |
| Topup disabled | `充值功能未开启` |
| Preset amount not configured | `快捷金额未配置` |
| Custom amount disabled | `未开启自定义金额` |
| Amount must be > 0 | `充值金额必须大于 0` |
| Order number generation failed | `生成订单号失败，请重试` |


## 2. Read Topup Settings

**Endpoint:** `GET /api/setting/topup`

**Auth:** Root

**Response:**

```json
{
  "success": true,
  "message": "",
  "data": {
    "enabled": true,
    "allow_custom": true,
    "exchange_rate": 1000,
    "presets": [
      { "amount": 50.0,  "bonus_quota": 50000 },
      { "amount": 100.0, "bonus_quota": 110000 },
      { "amount": 500.0, "bonus_quota": 600000 }
    ]
  }
}
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| enabled | bool | `topup.enabled`, master switch |
| allow_custom | bool | `topup.allow_custom`, allow user-supplied custom amounts |
| exchange_rate | int64 | `topup.exchange_rate`, CNY → quota ratio for custom amounts, default 1 |
| presets | array | `topup.presets` (decoded from JSON). Each entry is a `TopupPreset`. |

**`TopupPreset`:**

| Field | Type | Description |
|-------|------|-------------|
| amount | float64 | Topup amount (CNY) |
| bonus_quota | int64 | Quota actually credited |


## 3. Save Topup Settings

**Endpoint:** `PUT /api/setting/topup`

**Auth:** Root

**Request body:**

```json
{
  "enabled": true,
  "allow_custom": true,
  "exchange_rate": 1000,
  "presets": [
    { "amount": 50.0,  "bonus_quota": 50000 },
    { "amount": 100.0, "bonus_quota": 110000 },
    { "amount": 500.0, "bonus_quota": 600000 }
  ]
}
```

**Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| enabled | bool | no | Master switch, default false |
| allow_custom | bool | no | Allow custom amounts, default false |
| exchange_rate | int64 | no | CNY → quota ratio for custom amounts, must be > 0 |
| presets | array | no | Preset list; the whole list replaces the previous one |

**Validation rules** (`model.SaveTopupSettings`):

- `presets[i].amount > 0`
- `presets[i].bonus_quota >= 0`
- `presets` must not contain duplicate `amount` values
- `exchange_rate > 0`

**Response:**

```json
{ "success": true, "message": "已保存" }
```

**Errors:**

| Scenario | message |
|----------|---------|
| Invalid JSON | `无效的参数` |
| Item N has `amount <= 0` | `第 N 项金额必须大于 0` |
| Item N has `bonus_quota < 0` | `第 N 项额度不能为负数` |
| Duplicate preset amount | `快捷金额重复：X.XX 元已存在` |
| `exchange_rate <= 0` | `兑换比例必须大于 0` |