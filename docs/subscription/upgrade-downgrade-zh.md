---
title: 套餐升降级
description: "差价升级（price_diff）与叠加（stack）、同级别与降级的拒绝规则。"
category: subscription
order: 5
---

# 套餐升降级

> 用户从当前套餐切到另一个套餐的两条路径。实现：`model/order_payment.go::CreatePlanOrder` / `ActivatePackageByOrder` / `CalculateUpgradePrice`。

## 触发

用户已经有 active `user_plan` 时再次发起 `POST /api/order/plan`，服务端按系统设置 `plan.upgrade_mode`（`model.SystemSettingKeyPlanUpgradeMode`，缺省 `price_diff`）分支：

- `price_diff` — 差价升级（默认）
- `stack` — 叠加

## 校验

不论走哪条分支，服务端先做硬校验（`CreatePlanOrder`）：

- `currentPlan.Sort` 与 `newPlan.Sort` 的比较：
  - `newPlan.Sort < currentPlan.Sort` → 返回 `不能降级到低级别套餐`
  - `newPlan.Sort == currentPlan.Sort` → 返回 `您已经订阅了同级别的套餐`
- 没有当前 active plan 时，直接走"全新订阅"逻辑

`Plan.Sort` 由管理员在套餐设置中手动维护，数字越大越高级。

## `price_diff`（差价升级）

订单号前缀 `UP`，金额按 `CalculateUpgradePrice`：

```

当 `old_plan.end_time <= now`（即原套餐已过期）时直接 `return new_plan.price`（视作全新订阅）。

激活路径 `ActivatePackageByOrder(order, OrderUpgradeModePriceDiff)`：

1. 把该用户全部 `status=UserPlanStatusActive` 的 `user_plans` 翻成 `status=Expired`
2. 插入新 `user_plan`，`start_time=now`，`end_time=now + plan.duration_days × 86400`，`billing_type=token`
3. 调 `order.MarkOrderPaid` 标记订单已支付
4. `CacheDeleteUserActivePlans(userId)`

## `stack`（叠加）

订单号前缀 `TB`，金额 = `new_plan.price`（无差价），与现有订阅无关。

激活路径 `ActivatePackageByOrder(order, OrderUpgradeModeStack)`：

1. **不**动现有 active `user_plans`
2. 插入新 `user_plan`（与 price_diff 同样的 start_time / end_time 计算）
3. `MarkOrderPaid` + 清缓存

后续 `CheckPlanQuota` 按 `end_time` 升序选择第一个未耗尽的 plan；多个订阅之间自然按到期日轮换使用。

## 管理员 grant

`POST /api/subscription/`（`controller/subscription.go::AddSubscription`）由管理员创建订阅，**始终**用 `OrderUpgradeModeStack` 语义：现有订阅不被关闭，新订阅并存。

## 前端操作指南

- 用户侧 `/subscription` 的「续费」按钮直接调 `POST /api/order/plan { plan_id, pay_method }`，由后端决定 `price_diff` / `stack`
- 订单创建成功后用返回的 `pay` 对象渲染二维码 / 跳转链接（`buildPayInfo`）；管理员 grant 不需要这一步
- 列表里的「套餐」列显示 `plan.name`，排序由 `plan.sort` 决定
- 管理员可在 `/subscription` 行内做「延期」/「过期」/「删除」

## 实现位置

| 关注点 | 位置 |
|---|---|
| 模式设置 | `model/system_setting.go::SystemSettingKeyPlanUpgradeMode` |
| 订单 + 模式判定 | `model/order_payment.go::CreatePlanOrder` |
| 差价公式 | `model/order_payment.go::CalculateUpgradePrice` |
| 激活（两种模式） | `model/order_payment.go::ActivatePackageByOrder` |
| 管理员 grant | `controller/subscription.go::AddSubscription` |
| 支付预下单 | `controller/order.go::buildPayInfo` |

