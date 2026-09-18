---
title: 充值 API
description: "/api/topup/* 端点。"
category: api
order: 9
---
# 充值 API

`/api/topup/*` 提供**在线支付充值**（余额充值）端点：用户自助下单、查询支付参数，以及管理员维护充值设置（总开关、自定义金额、快捷金额、汇率）。订单 type 固定为 `2`（`OrderTypeTopup`）；回调通知走 `/api/payment/*`。

## 端点一览

| 接口 | 方法 | 权限 | 说明 |
|------|------|------|------|
| `/api/topup/order` | POST | User | 用户自助创建充值订单 |
| `/api/setting/topup` | GET | Root | 读取充值设置 |
| `/api/setting/topup` | PUT | Root | 保存充值设置 |


## 公共约定

- 充值订单号前缀：`TP`（参见 `model.GenerateOrderNo("TP")`）。
- 任意支付通道已启用（`payment.AnyChannelEnabled().any_enabled == true`）且 `topup.enabled == true` 才允许下单，否则拒绝。
- 金额解析规则（见 `model.ResolveTopupAmount`）：
  - `preset_amount > 0`：命中预设，从配置中找 `amount` 相等的 preset，使用其 `bonus_quota`。
  - 否则使用 `amount` 自定义金额，需 `topup.allow_custom == true`，`bonus_quota = amount × exchange_rate`。
- 兑换比例默认 `1:1`（1 元 = 1 quota），仅作用于自定义金额。
- 订单激活通过 `model.ActivateTopupByOrder`（异步回调时触发）：给用户 `IncreaseUserQuota(bonus_quota)` 并把订单置为 `status=1`。


## 1. 创建充值订单

**接口：** `POST /api/topup/order`

**权限：** User

**说明：** 校验金额 → 调 `model.CreateTopupOrder` 持久化订单（type=2, status=0）→ 调 `controller.buildPayInfo` 生成支付参数。

**前置条件：**

- 已登录（`c.GetInt("id") != 0`）
- 任意支付通道已启用
- `topup.enabled == true`

**请求体：**

```json
{
  "amount": 100.0,
  "preset_amount": 0,
  "pay_method": "wechat"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| amount | float64 | 二选一 | 自定义金额（元）；`preset_amount > 0` 时忽略 |
| preset_amount | float64 | 二选一 | 命中预设金额（元）；命中后使用对应 preset 的 `bonus_quota` |
| pay_method | string | 是 | `wechat` / `alipay` / `bank` |

**返回示例：**

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

**`pay` 字段：**

| 字段 | 类型 | 说明 |
|------|------|------|
| status | string | `success`（可支付）/ `warning`（通道未启用或预下单失败） |
| pay_url | string | 支付跳转 URL（微信/支付宝） |
| qr_code | string | 二维码内容（Native 通道与 `pay_url` 一致） |
| expire_at | int64 | 过期时间戳，0 表示未知 |
| trade_no | string | 支付通道预下单 id（不等于 `order_no`） |
| note | string | 仅 `bank` 通道返回：转账须知 |
| warning | string | 仅 `status="warning"` 时返回：失败原因 |

**错误情况：**

| 场景 | message |
|------|---------|
| 未登录 | `未登录` |
| JSON 解析失败 | `无效的参数` |
| `amount <= 0 && preset_amount <= 0` | `amount 或 preset_amount 至少传一个` |
| `pay_method` 为空 | `pay_method 不能为空` |
| 不支持的支付方式 | `不支持的支付方式` |
| 支付方式非 `wechat/alipay/bank` | `自助充值仅支持 wechat / alipay / bank` |
| 未配置任何支付通道 | `系统尚未开通任何支付通道，请设置后开启支付` |
| 充值功能未开启 | `充值功能未开启` |
| 快捷金额未配置 | `快捷金额未配置` |
| 未开启自定义金额 | `未开启自定义金额` |
| 充值金额必须大于 0 | `充值金额必须大于 0` |
| 订单号生成失败 | `生成订单号失败，请重试` |


## 2. 读取充值设置

**接口：** `GET /api/setting/topup`

**权限：** Root

**返回示例：**

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

**字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| enabled | bool | `topup.enabled`，总开关 |
| allow_custom | bool | `topup.allow_custom`，是否允许用户输入自定义金额 |
| exchange_rate | int64 | `topup.exchange_rate`，自定义金额 1 元 = X quota，默认 1 |
| presets | array | `topup.presets`（JSON 反序列化结果），每项 `TopupPreset` |

**`TopupPreset`：**

| 字段 | 类型 | 说明 |
|------|------|------|
| amount | float64 | 充值金额（元） |
| bonus_quota | int64 | 实际入账 quota |


## 3. 保存充值设置

**接口：** `PUT /api/setting/topup`

**权限：** Root

**请求体：**

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

**字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| enabled | bool | 否 | 总开关，默认 false |
| allow_custom | bool | 否 | 是否允许自定义金额，默认 false |
| exchange_rate | int64 | 否 | 自定义金额 1 元 = X quota，必须大于 0 |
| presets | array | 否 | 快捷金额列表，会整体覆盖原有列表 |

**校验规则**（`model.SaveTopupSettings`）：

- `presets[i].amount > 0`
- `presets[i].bonus_quota >= 0`
- `presets` 内 `amount` 不可重复
- `exchange_rate > 0`

**返回：**

```json
{ "success": true, "message": "已保存" }
```

**错误情况：**

| 场景 | message |
|------|---------|
| JSON 解析失败 | `无效的参数` |
| 第 N 项金额 <= 0 | `第 N 项金额必须大于 0` |
| 第 N 项额度 < 0 | `第 N 项额度不能为负数` |
| 快捷金额重复 X 元 | `快捷金额重复：X.XX 元已存在` |
| 兑换比例 <= 0 | `兑换比例必须大于 0` |