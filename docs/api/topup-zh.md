---
title: 充值 API
description: "/api/topup/* 端点。"
category: api
order: 9
---
# 充值 API

`/api/topup/*` 提供**在线支付充值**（余额充值）端点：用户自助下单、查询支付参数，以及管理员维护充值设置（总开关、自定义金额、快捷金额）。订单 type 固定为 `2`（`OrderTypeTopup`）；回调通知走 `/api/payment/*`。

## 端点一览

| 接口 | 方法 | 权限 | 说明 |
|------|------|------|------|
| `/api/topup/order` | POST | User | 用户自助创建充值订单 |
| `/api/setting/topup` | GET | Root | 读取充值设置 |
| `/api/setting/topup` | PUT | Root | 保存充值设置 |

## 公共约定

- **额度基准为系统常量**：1 元 = `QuotaPerUnit` = `1,000,000` quota（1 quota = 1e-6 元），语义同微信支付的「分」。该常量仅参与计算，不入库、不可修改；`/api/status` 以只读字段 `quota_per_unit` 下发，供前端折算展示。
- 充值订单号前缀：`TP`（参见 `model.GenerateOrderNo("TP")`）。
- 任意支付通道已启用（`payment.AnyChannelEnabled().any_enabled == true`）且 `topup.enabled == true` 才允许下单，否则拒绝。
- 金额解析规则（见 `model.ResolveTopupAmount`）：
  - `preset_amount > 0`：命中预设，从配置中找 `amount` 相等的 preset，使用其 `bonus_quota`。
  - 否则使用 `amount` 自定义金额，需 `topup.allow_custom == true`；**恒为 1:1**，`bonus_quota = round(amount × QuotaPerUnit)`。
- **不再提供兑换比例**：`topup.exchange_rate` 已移除，自定义金额恒为「支付多少到账多少」。营销赠送（如「充 10 得 15」）通过快捷金额的「支付金额 / 到账金额」体现。
- 订单激活通过 `model.ActivateTopupByOrder`（异步回调时触发）：给用户 `IncreaseUserQuota(bonus_quota)` 并把订单置为 `status=1`；订单快照中记录 `credit_amount`（到账金额，元）与 `bonus_quota`（到账额度）。

## 1. 创建充值订单

**接口：** `POST /api/topup/order`

**权限：** User

**说明：** 校验金额 → 调 `model.CreateTopupOrder` 持久化订单（type=2, status=0）→ 调 `controller.buildPayInfo` 生成支付参数。

**前置条件：**

- 已登录（`c.GetInt("id") != 0`）
- 任意支付通道已启用
- `topup.enabled == true`

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| amount | float64 | 二选一 | 自定义金额（元）；`preset_amount > 0` 时忽略。到账额度 = `round(amount × 1000000)` |
| preset_amount | float64 | 二选一 | 命中预设金额（元）；命中后使用对应 preset 的 `bonus_quota` |
| pay_method | string | 是 | `wechat` / `alipay` / `bank` |

```json
{
  "preset_amount": 10,
  "pay_method": "wechat"
}
```

**返回示例：**

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

> 上例为「充 10 得 15」的预设：支付 10 元，到账 15 元（`bonus_quota = 15000000`）。

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
    "presets": [
      { "amount": 10, "bonus_quota": 15000000 },
      { "amount": 50, "bonus_quota": 50000000 }
    ]
  }
}
```

**字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| enabled | bool | `topup.enabled`，总开关 |
| allow_custom | bool | `topup.allow_custom`，是否允许用户输入自定义金额（恒 1:1） |
| presets | array | `topup.presets`（JSON 反序列化结果），每项 `TopupPreset` |

`TopupPreset`：

| 字段 | 类型 | 说明 |
|------|------|------|
| amount | float64 | 支付金额（元） |
| bonus_quota | int64 | 到账额度（quota），必须 >= `round(amount × 1000000)` |

## 3. 保存充值设置

**接口：** `PUT /api/setting/topup`

**权限：** Root

**请求体：**

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

**字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| enabled | bool | 否 | 总开关，默认 false |
| allow_custom | bool | 否 | 是否允许自定义金额，默认 false |
| presets | array | 否 | 快捷金额列表，会整体覆盖原有列表 |

**校验规则**（`model.SaveTopupSettings`）：

- `presets[i].amount > 0`
- `presets` 内 `amount` 不可重复
- `presets[i].bonus_quota >= round(presets[i].amount × 1000000)`（允许赠送，不允许缩水）

**返回：**

```json
{ "success": true, "message": "已保存" }
```

**错误情况：**

| 场景 | message |
|------|---------|
| JSON 解析失败 | `无效的参数` |
| 第 N 项金额 <= 0 | `第 N 项金额必须大于 0` |
| 快捷金额重复 X 元 | `快捷金额重复：X.XX 元已存在` |
| 到账额度低于支付金额折算额度 | `第 N 项到账额度不能低于支付金额折算额度（X.XX 元 = Y 额度）` |
