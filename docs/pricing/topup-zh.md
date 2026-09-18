---
title: 充值设置
description: "`topup.*` 系统设置、总开关、快捷金额、自定义金额与汇率、支持的支付通道。"
category: pricing
order: 4
---

# 充值设置

> 充值（topup）= 用户在线自助支付把金额兑换成 `users.quota`。设置存储在 `system_settings` 表的 `topup.*` 键；用户侧流程走 `OrderTypeTopup=2`。实现：`model/topup.go`、`controller/topup.go`、`web/default-pro/src/views/setting/TopupSetting.vue`。

## 系统设置 / Settings

`model/system_setting.go` 定义：

| Key | 类型 | 说明 |
|---|---|---|
| `topup.enabled` | bool 字符串 | 总开关；缺省 `false` |
| `topup.allow_custom` | bool 字符串 | 是否允许用户输入自定义金额；缺省 `false` |
| `topup.presets` | JSON | 快捷金额列表 `[{"amount": 10, "bonus_quota": 10000}, ...]` |
| `topup.exchange_rate` | int64 字符串 | 自定义金额 1 元 = X quota（缺省 `1`，即 1:1） |

类别常量：`SystemSettingCategoryTopup="topup"`。

## 接口 / Endpoints

| Endpoint | Method | Auth | 说明 |
|---|---|---|---|
| `/api/setting/topup` | `GET` | Root | 读取四项设置（`TopupSettingsResponse`） |
| `/api/setting/topup` | `PUT` | Root | 整体覆盖保存 |
| `/api/topup/order` | `POST` | User | 创建充值订单（详见下文） |

`SaveTopupSettings` 校验：

- 每个 preset `amount > 0`、`bonus_quota >= 0`；`amount` 不可重复（`快捷金额重复：X.XX 元已存在`）
- `exchange_rate > 0`（`兑换比例必须大于 0`）
- `presets == nil` 时落库为 `[]`

## 快捷金额 / Presets

`model.TopupPreset`：

```json
{
  "amount": 10,
  "bonus_quota": 10000
}
```

- `amount` = 用户支付金额（元，浮点）
- `bonus_quota` = 到账 quota 整数（= 10 元 × 默认汇率 = 10000 quota 时即 1:1000 汇率）

`ResolveTopupAmount`：

- `PresetAmount > 0` → 在 presets 里查 `amount == PresetAmount` 命中返回；未命中返回 `快捷金额未配置`
- 否则走自定义金额路径，要求 `allow_custom == true`

## 自定义金额与汇率 / Custom amount

当 `PresetAmount == 0` 时走自定义路径：

```text
bonus_quota = int(amount × exchange_rate)
```

汇率缺省 `1`（即 1 元 = 1 quota）。生产环境通常调高以体现"1 元 = 1000 quota"或类似的换算，让前端不再做金额展示小数。

## 充值订单 / Top-up order

`POST /api/topup/order` 请求体（`CreateTopupOrderRequest`）：

```json
{
  "amount": 100,
  "preset_amount": 0,
  "pay_method": "wechat"
}
```

服务端校验：

1. `amount <= 0 && preset_amount <= 0` → `amount 或 preset_amount 至少传一个`
2. `pay_method` 必须在 `OrderPayMethod*` 白名单内（`wechat` / `alipay` / `bank` / `offline` / `free`）
3. 自助路径只接受 `wechat` / `alipay` / `bank`（`自助充值仅支持 wechat / alipay / bank`）
4. 任意支付通道必须启用（`payment.AnyChannelEnabled()`），否则返回 `系统尚未开通任何支付通道，请设置后开启支付`
5. `topup.enabled == true`，否则 `充值功能未开启`
6. 解析金额 / 配额（`ResolveTopupAmount`）

`model.CreateTopupOrder` 创建订单：

- 订单号前缀 `TP`（`GenerateOrderNo("TP")`）
- `Type=OrderTypeTopup=2`、`Source=OrderSourceUserSelf=1`、`Status=OrderStatusPending=0`、`PayStatus=OrderPayStatusPending=0`
- `PlanInfo` 写入 `TopupOrderPlanInfo` JSON（amount / preset_amount / bonus_quota / exchange_rate 快照），用于对账

随后 `buildPayInfo(pay_method, order_no, amount, "余额充值")` 生成预支付参数返回给前端。

## 激活 / Activation

支付回调（`controller/payment.go::processNotify`）按 `order.Type` 分发：

- `OrderTypeTopup=2` → `model.ActivateTopupByOrder(order)`：
  - 幂等：`status == Paid` 直接返回
  - 解析 `PlanInfo` 取出 `bonus_quota`；缺失时按 `amount × exchange_rate` 兜底
  - `IncreaseUserQuota(order.UserId, bonus_quota)` 直接加到 `users.quota`
  - `MarkOrderPaid` 落 `pay_status=1` / `pay_time` / `pay_trade_no`

## 历史兼容 / Migration note

旧字段 `plan.allow_topup`（DB 行保留）已迁移到 `topup.enabled`，UI 不再读写；如需查看历史值，直接 `SELECT * FROM system_settings WHERE 'key' = 'plan.allow_topup'`。

## 前端操作指南 / Frontend Guide

页面：`/setting/operation` 中的"充值"区块（`web/default-pro/src/views/setting/TopupSetting.vue`）。

- 三个开关：`enabled` / `allow_custom`
- 汇率：`exchange_rate`（整数，最小 1）
- 快捷金额表：行内编辑 amount + bonus_quota；新增 / 删除行；提交时整组覆盖保存
- 用户侧 `/pricing` 或 `/topup`：preset chip + 「自定义金额」入口（仅当 `allow_custom=true`）

## 实现位置 / Implementation Pointers

| 关注点 | 位置 |
|---|---|
| 系统设置常量 | `model/system_setting.go::SystemSettingKeyTopup*` |
| 设置读写 | `model/topup.go::GetTopupSettings` / `SaveTopupSettings` |
| 金额解析 | `model/topup.go::ResolveTopupAmount` |
| 创建订单 | `model/topup.go::CreateTopupOrder` |
| 激活 | `model/topup.go::ActivateTopupByOrder` |
| 设置接口 | `controller/topup.go::GetTopupSettings` / `PutTopupSettings` |
| 用户下单接口 | `controller/topup.go::CreateTopupOrder` |
| 支付回调分发 | `controller/payment.go::processNotify` |
