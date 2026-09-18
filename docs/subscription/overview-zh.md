---
title: 订阅概览
description: "订阅、按量计费与充值额度三种计费方式的区别与扣减顺序。"
category: subscription
order: 1
---

# 订阅概览

> One API Pro 同时支持三种计费通道：订阅（plan，按限额与窗口扣减）、按量计费（pay-as-you-go，从用户 `quota` 扣减）、充值额度（topup，也是用户 `quota`）。三者由中间件统一判定后扣减。

## 三种计费

| 模式 | 来源 | 落库 | 扣减点 |
|---|---|---|---|
| 订阅 | `model.Plan` + `model.UserPlan` | `user_plans` 行 + `plan_usages` 行 | 中间件 `middleware/plan_quota.go::PlanQuotaCheck` 按窗口配额放行；`relay/handler/helper.go::postConsumeQuota` 在用完请求后写 `plan_usages` |
| 按量计费 | 用户直接调用 `/v1/chat/completions` 且无有效订阅 | `tokens.quota` + `users.quota` | `model.PreConsumeTokenQuota` 预扣 → 上游返回 → `PostConsumeTokenQuota` 结算 |
| 充值 | 兑换码 / 充值订单 / 管理员手动加 | `users.quota` 直接累加 | `model.IncreaseUserQuota`（不参与请求时的扣减判定，仅提供余额） |

## 扣减顺序

请求到达 relay 中间件链时：

1. `PlanQuotaCheck` 调 `model.CheckPlanQuota(userId, model)` 遍历 `CacheGetUserActivePlans(userId)` 返回的活跃订阅（按 `end_time` 升序），按 `(model, window_type)` 匹配到第一个未耗尽的 plan 后放行，并把 `plan_id` 写进 `meta.PlanId`。没有匹配则按"无订阅"路径放行。
2. 计费阶段 `postConsumeQuota`：
   - 若 `meta.PlanId > 0`：把 `requests`、`prompt_tokens`、`completion_tokens`、`cached_tokens` 累加到该 plan 的 `plan_usages`，并把预扣的 quota 退回 token（订阅用户不消耗 `quota`）。
   - 否则：按 `priceResult.billing_type` 计算应扣 `quota` 并 `PostConsumeTokenQuota` 结算；非订阅用户从 `users.quota` 中扣减。

## 订阅与 topup 的关系

- 订阅有自己的窗口配额：每个 plan 在 `plan.model_limits` 中定义模型级 `request_period/week/month` / `token_period/week/month`。**不消耗** `users.quota`。
- `topup`（充值额度）增加的是 `users.quota`，是按量计费时的真实余额；管理员手动加 quota（`/api/user/topup` legacy）或管理员 grant 走订单中心（`OrderTypeTopup=2`）也是同一条路径。
- 一个用户可以同时有多个订阅（叠加模式 `OrderUpgradeModeStack`），订阅之间按 `end_time` 升序依次使用；topup 与订阅独立，订阅不消耗 topup 余额。

## 计划升级

当用户已有活跃订阅时再次下单时，根据系统设置 `plan.upgrade_mode` 走两条路径：

- `price_diff`（默认）：新订阅比当前更高级时按"差价"扣费，订单号前缀 `UP`，`ActivatePackageByOrder` 会先把所有当前活跃 `user_plans` 置为 expired 再插入新行；同级别或更低级别拒绝。
- `stack`：新订阅与当前共存，互不影响，订单号前缀 `TB`，订单金额 = 新套餐原价。

详见 [套餐升降级](./upgrade-downgrade)。

## 过期

`model.ExpireUserPlans`（周期任务）把 `status=1 AND end_time <= now` 的行翻成 `UserPlanStatusExpired=0`。到期后该订阅不再出现在 `CheckPlanQuota` 结果中，用户请求走按量计费。

## 实现位置

| 关注点 | 位置 |
|---|---|
| 数据模型 | `model/plan.go::Plan` / `UserPlan` / `PlanUsage` |
| 配额判定 | `model/plan_quota.go::CheckPlanQuota` |
| 中间件 | `middleware/plan_quota.go::PlanQuotaCheck` |
| 计费结算 | `relay/handler/helper.go::postConsumeQuota` |
| 升级 / 叠加 | `model/order_payment.go::CreatePlanOrder` / `ActivatePackageByOrder` |
| 用户充值 | `model/topup.go::CreateTopupOrder` / `ActivateTopupByOrder` |

