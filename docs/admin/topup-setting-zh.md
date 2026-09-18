---
title: 充值设置
description: "在线充值的总开关、预设金额、自定义金额 chip、兑换比例与币种。"
category: admin
order: 14
---

# 充值设置

> 控制在线充值（topup）业务的全局开关、自定义金额、预设金额 chip 与兑换比例。前端组件：`web/default-pro/src/views/setting/TopupSetting.vue`。

## 接口一览 / Endpoints

| Endpoint | Method | 鉴权 | 说明 |
|---|---|---|---|
| `/api/setting/topup` | `GET` | Root | 一次性返回 `{ enabled, allow_custom, exchange_rate, presets }` |
| `/api/setting/topup` | `PUT` | Root | 整体覆盖保存（同上结构） |

实现：`controller/topup.go::GetTopupSettings` / `PutTopupSettings`。

## 字段 / Fields

| 字段 | 类型 | 说明 |
|---|---|---|
| `enabled` | `bool` | 总开关。`false` 时用户自助下单直接报「充值功能未开启」 |
| `allow_custom` | `bool` | 是否允许用户输入自定义金额；`false` 时只能从 `presets` 里选 |
| `exchange_rate` | `int` | 自定义金额 1 元 = X quota；必须 > 0；缺省 = 1（即 1:1） |
| `presets` | `[{ amount, bonus_quota }]` | 快捷金额 chip 列表；每个 `amount` 必须 > 0 且全局唯一 |

校验由 `model/topup.go::SaveTopupSettings` 完成：
- `amount <= 0` → 「第 N 项金额必须大于 0」；
- `bonus_quota < 0` → 「第 N 项额度不能为负数」；
- 重复金额 → 「快捷金额重复：X.XX 元已存在」；
- `exchange_rate <= 0` → 「兑换比例必须大于 0」。

## 数据库 / Persistence

四个独立 `system_settings` 行（category=`topup`）：

| key | value 形态 |
|---|---|
| `topup.enabled` | `"true"` / `"false"` |
| `topup.allow_custom` | `"true"` / `"false"` |
| `topup.exchange_rate` | 整数字符串 |
| `topup.presets` | JSON 数组：`[{"amount":10,"bonus_quota":10000}, ...]` |

PUT 是整体覆盖——不存在的字段会被保留（前端总会把所有字段都发上来）。

## 用户侧下单 / How It Works

`POST /api/topup/order`（`controller/topup.go::CreateTopupOrder`）的解析由 `model/topup.go::ResolveTopupAmount` 完成：

1. `preset_amount > 0` 时：到 `presets` 中查匹配金额，命中即返回该条；不命中 → 「快捷金额未配置」；
2. 否则按 `amount × exchange_rate` 计算 `bonus_quota`；
3. `allow_custom=false` 时不允许走自定义分支。

## 预设 chip / Frontend Chips

`PaymentTopupModal.vue`（用户侧）按 `presets` 渲染金额 chip，并在 `allow_custom=true` 时在末尾追加一个「自定义」chip。
选中 chip 的索引 `-1` 作为「自定义」哨兵（参考 `web/default-pro/src/utils/topup.js` 里的 `validateTopupPresets`）。

## 前端操作指南 / Frontend Guide

- 顶部「总开关」「允许自定义金额」「兑换比例（1 元 = X quota）」三组竖排 `<a-form-item>`。
- 下面「快捷金额」表格（bordered cell）支持新增、删除、编辑 `amount` / `bonus_quota`；编辑直接在 `<a-input-number>` 里改：
  - `amount`：precision=2，min=0.01；
  - `bonus_quota`：precision=0，min=0，step=1000。
- 提交前会跑 `validateTopupPresets` 给出前端预校验错误。
- 点击「保存」调用 `PUT /api/setting/topup`；成功后 `Message.success`。

## 接口实现 / Implementation Pointers

| 关注点 | 位置 |
|---|---|
| Controller | `controller/topup.go` |
| 设置读写 | `model/topup.go::GetTopupSettings` / `SaveTopupSettings` |
| 金额解析 | `model/topup.go::ResolveTopupAmount` |
| 创建订单 | `model/topup.go::CreateTopupOrder` |
| 前端校验工具 | `web/default-pro/src/utils/topup.js` |
| 路由 | `router/api.go` |