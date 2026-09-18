---
title: 套餐运营设置
description: "套餐业务规则：升级模式（叠加 vs 差价）、到期策略等。"
category: subscription
order: 13
---

# 套餐运营设置

> 套餐的全局业务规则——目前只包含 `upgrade_mode`（升级时采用「叠加」还是「差价」）。前端组件：`web/default-pro/src/views/setting/OperationSetting.vue` 末段「套餐」区。

## 接口一览

| Endpoint | Method | 鉴权 | 说明 |
|---|---|---|---|
| `/api/setting/plan` | `GET` | Root | 拉取套餐运营配置；缺省值 `price_diff` |
| `/api/setting/plan` | `PUT` | Root | 保存套餐运营配置；只接受 `price_diff` / `stack` |

实现：`controller/setting_payment.go::GetPlanSettings` / `PutPlanSettings`。

## 升级模式

设置 key：`plan.upgrade_mode`，取值常量：

| 取值 | 常量 | 含义 |
|---|---|---|
| `price_diff` | `OrderUpgradeModePriceDiff` | **差价升级**：当用户已有活跃订阅并升级到更贵的套餐时，按差价补足（仅补差额，不叠加旧套餐） |
| `stack` | `OrderUpgradeModeStack` | **叠加升级**：在已有激活之上叠加新套餐的有效期（`end_time += duration_days`），常用作「VIP 时长续费」 |

> 默认值为 `price_diff`；PUT 时如果值不在上表内会返回 `upgrade_mode 必须是 price_diff 或 stack`。

## 与订单激活的协作

- 用户自助下单（`POST /api/order/plan`）：根据 `upgrade_mode` 计算升级时是否需要补差价（差价模式下若新套餐价 ≤ 旧套餐实付，订单金额被 clamp 到 0）。
- 管理员手动开通（`POST /api/subscription`）：**始终**走 `stack`（`OrderUpgradeModeStack`），不受此设置影响。

## 数据存储

`system_settings` 表，key `plan.upgrade_mode`，category `plan`。
仅一行字符串值；GET 端点把它包装成 `{ data: { upgrade_mode: "price_diff" } }`。

> 历史 `plan.allow_topup` 字段已迁移到 `topup.enabled`，DB 行保留但 UI 不再读写（见 `model/system_setting.go` 注释）。

## 前端操作指南

- 在 `/setting/operation` 页面最底部「套餐」区块，提供两个单选按钮：「差价升级 (`price_diff`)」/「叠加升级 (`stack`)」。
- 修改后点击「保存」按钮（位于本区块内，独立于其它区段）。
- 默认值显示为「差价升级」。

## 接口实现

| 关注点 | 位置 |
|---|---|
| 设置 handler | `controller/setting_payment.go` (`GetPlanSettings` / `PutPlanSettings`) |
| 常量定义 | `model/order.go::OrderUpgradeModePriceDiff` / `OrderUpgradeModeStack` |
| 用户自助下单差异 | `model/order_payment.go::CreatePlanOrder` |
| 管理员开通（始终 stack） | `controller/subscription.go::AddSubscription` |
| 路由 | `router/api.go` |
