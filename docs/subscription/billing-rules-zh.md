---
title: 计费规则
description: "三种窗口（period/week/month）、加权扣减公式、batchUpdate 路径与扣减顺序。"
category: subscription
order: 4
---

# 计费规则

> 套餐的窗口配额与请求时扣减。实现：`model/plan.go::CalcWindowIndex`、`model/plan_quota.go::WeightedUsage`、中间件 `middleware/plan_quota.go::PlanQuotaCheck`。

## 三种窗口

| windowType | 长度 | 起点 |
|---|---|---|
| `period` | `rule.period_h` 小时（默认 5） | `start_time` |
| `week` | 7 天 | `start_time` |
| `month` | 30 天 | `start_time` |

每个窗口索引由 `model.CalcWindowIndex(now, startTime, windowType, periodH)` 计算：

```

下次重置时间 `CalcNextResetTime = start_time + (windowIndex+1) * windowDuration`，因此 period 窗口的"每天 5 小时滚动"按 `start_time` 自然推进，而不是按自然日对齐。

## 加权使用率

每个 plan 内多个模型有自己的限额。`WeightedUsage` 把所有已用 `(consumed, limit)` 对折算到一个 0..100 的统一使用率：

```

任一窗口的 `weighted >= 100` 即视为该窗口耗尽。

具体规则（`model/plan_quota.go::WeightedUsage`）：

1. 遍历 `PlanUsage` 行，按 `(model, window_type, window_index)` 匹配当前窗口
2. 按 `billing_type` 决定使用哪个限额与已用量：
   - `BillingTypeRequest="request"`：`consumed = requests`，`limit = rule.request_period|week|month`
   - `BillingTypeToken="token"`：`consumed = prompt_tokens + completion_tokens`，`limit = rule.token_period|week|month`
   - 其它 / 缺省：先尝试 request 维度；若 `limit <= 0` 则回退到 token 维度
3. `consumed × 100 / limit` 累加；`limit <= 0` 的模型不参与
4. 只累加当前窗口（`windowIndex` 一致）的 usage 行

## 模型解析

`FindLimit(limits, model, defaultModel)` 优先级：

1. 精确命中 `limits[model]` → 用该规则
2. 否则回退到 `limits[defaultModel]` → 用 default_model 的规则
3. 都没有 → 返回 `not found`，相关 usage 不参与加权

`CheckPlanQuota` 还会把命中的 default_model 透出给 middleware，由 `PlanQuotaCheck` 把请求的 `RequestModel` 改写到 `default_model`，从而把未知模型路由到套餐默认模型。

## 扣减顺序

请求进入 relay 中间件链：

1. `PlanQuotaCheck` → `CheckPlanQuota` 遍历 `CacheGetUserActivePlans(userId)`（`end_time` 升序）。对每个 `user_plan`：
   - 若 `end_time <= now`，把状态翻成 expired 并跳过
   - 取 `plan.model_limits`；若 nil 直接视为 usable
   - 计算 period / week / month 三个 weighted；任一 `>= 100` 视为该 plan 耗尽，进入下一个
   - 第一个 usable 的 plan 被放行，`plan_id` 写入 `meta.PlanId`，`billing_type` 写入 `meta.BillingType`
2. 没有匹配 plan：放行；`meta.PlanId=0` 走按量计费
3. `postConsumeQuota`（`relay/handler/helper.go`）：
   - `PlanId > 0`：`IncrementPlanUsage(plan_id, model, window_type, window_index, 1, prompt_tokens, completion_tokens, cached_tokens)` 三个窗口都写一次；并把预扣的 quota 退回 token
   - `PlanId == 0`：按 `priceResult.BillingType` 计算 `quota`：`PerRequest` 用 `per_request_price`，`Token` 用 `(input_price × prompt + output_price × completion + cached_price × cached) × group_discount`；`PostConsumeTokenQuota` 结算 `users.quota` 与 `tokens.quota`

## 缓存

- `user_plans:<id>`：Redis 缓存当前用户的活跃订阅（`UserPlanCacheSeconds=300`），`AddSubscription` / `UpdateSubscription` / `DeleteSubscription` 写后调 `CacheDeleteUserActivePlans` 失效
- `model_price:<name>` / `group_price:<group>:<model>`：模型与分组的折扣（`ModelPriceCacheSeconds=300`）
- `model_price` 与 `group_price` 表由 `SyncModelPriceCache` / `SyncGroupPriceCache` 周期同步

## batchUpdate 路径

当 `BATCH_UPDATE_ENABLED=true` 时：

- `model.UpdateChannelUsedQuota` 把渠道已用额度累加写到内存，由后台批量任务统一刷库（`BATCH_UPDATE_INTERVAL`，默认 5s）
- `model.IncreaseUserQuota` 也走 batch；其它写路径仍同步落库

启用后会减少 DB 写放大，但需要在 `main.go` 显式 `model.InitBatchUpdater()` 并保证进程内单点。

## 实现位置

| 关注点 | 位置 |
|---|---|
| 窗口计算 | `model/plan.go::CalcWindowIndex` / `GetWindowDurationSeconds` |
| 加权扣减 | `model/plan_quota.go::WeightedUsage` / `CalculateWeightedUsage` |
| 模型解析 | `model/plan_quota.go::FindLimit` |
| 中间件判定 | `middleware/plan_quota.go::PlanQuotaCheck` |
| 写 plan_usage | `model/plan.go::IncrementPlanUsage` |
| 计费结算 | `relay/handler/helper.go::postConsumeQuota` |
| 缓存 | `model/plan.go::CacheGetUserActivePlans` / `CacheDeleteUserActivePlans` |

