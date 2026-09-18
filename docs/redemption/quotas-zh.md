---
title: 兑换码额度计算
description: "兑换码 `quota` 字段与套餐绑定、过期、范围与单位。"
category: redemption
order: 4
---

# 兑换码额度计算

> 兑换码的到账金额以 `redemptions.quota` 为单一来源；与套餐 / topup 订单的"金额"语义不同：兑换码不存金额，只存 quota 整数。

## 字段 / Field

`model.Redemption.Quota`（`quota`，`bigint`，默认 100）：

- 类型 `int64`；单位 = `users.quota` 同样的 quota 单位
- 兑换成功后 `users.quota = users.quota + redemption.quota`（`Redeem` 事务）

不在数据库里保存任何金额或汇率；兑换码与币种、汇率无关。

## 换算 / Conversion

兑换码没有"amount × rate"换算路径——管理员在生成时直接填写 `quota` 整数。

如果运营上希望"以人民币金额分发"，可以在前端做金额 → quota 的换算（与 `model.GetTopupSettings()` 返回的 `exchange_rate` 一致），但**落库值仍是 quota 整数**。

## 与套餐的绑定 / Plan binding

兑换码与套餐正交：兑换码不会创建 / 续期 / 影响 `user_plans`。管理员若想"绑定到某套餐"，需要在前端按 `plan.name` 分桶展示码池，但每张码本身仍是通用 quota。

## 有效期 / Expiry

`redemptions` 表没有 `expire_time` 字段。过期策略只在 `status` 维度：

- 启用中（`Enabled=1`） → 任意时刻可兑换
- 已停用（`Disabled=2`） → `Redeem` 拒绝
- 已使用（`Used=3`） → `Redeem` 拒绝

如需"30 天后失效"，由管理员在合适时间点批量调 `PUT ?status_only=true {status:2}`，或写一个一次性 cron 任务直接更新 `status`。

## 兑换额度与余额的差异 / Quota vs balance

`users.quota` 是按量计费（pay-as-you-go）的真实余额；兑换码到账后该余额增加，但不参与套餐窗口配额。

当用户同时有活跃订阅 + 兑换码到账余额时：

- 套餐窗口内的请求走订阅路径（`meta.PlanId > 0`），不扣 `users.quota`
- 订阅耗尽后的请求走按量计费，从 `users.quota` 扣减；这部分会被兑换码到账的余额"充值"

## 套餐内赠额 / Plan grants

套餐本身没有"自动赠额"语义：`Plan` 不携带 quota，`UserPlan` 也不消耗 quota。若需要"开套餐即赠 N quota"，可在管理员 grant 流程后串一个 `OrderTypeTopup=2` 的订单，或在 `ActivatePackageByOrder` 之外直接 `IncreaseUserQuota`。

## 实现位置 / Implementation Pointers

| 关注点 | 位置 |
|---|---|
| 字段定义 | `model/redemption.go::Redemption.Quota` |
| 加 quota | `model/redemption.go::Redeem` (`UPDATE users SET quota = quota + ?`) |
| topup 汇率 | `model/topup.go::GetTopupSettings` / `SystemSettingKeyTopupExchangeRate` |
| 套餐 grant | `controller/subscription.go::AddSubscription` |
