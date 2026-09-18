---
title: 用户订阅使用指南
description: "用户侧 `/subscription` 与 `/subscription/:id/usage` 看到的内容：进度条、剩余时间、模型明细。"
category: subscription
order: 2
---

# 用户订阅使用指南

> 用户视角仅看到自己的活跃订阅 + 每窗口用量；本页与 [订阅概览](./overview) 互补：概览讲机制，本页讲 UI 字段。

## 接口

| Endpoint | Method | Auth | 说明 |
|---|---|---|---|
| `/api/subscription/self` | `GET` | User | 当前用户所有状态的 `user_plans`（最多 1000 条） |
| `/api/subscription/info` | `GET` | User | 按 `model.GetUserSubscriptionInfo` 聚合：每个 plan × 每个 model × 每个 windowType 一条 usage |
| `/api/subscription/current` | `GET` | User | 当前生效的 plan 详情（`end_time` 最近且未过期），并把 `Plan` 字段拍平 |
| `/api/subscription/:id/usage` | `GET` | User | 单条订阅的详细用量：原始 `plan_usages`、`weighted` 加权使用率、`model_usage` 模型级明细、`next_reset` 下次重置时间 |

实现：`controller/subscription.go`。

## `/api/subscription/current` 字段

服务端把 `Plan` 嵌入对象拍平到顶层：

```json
{
  "id": 88,
  "user_id": 42,
  "plan_id": 7,
  "order_id": 1735,
  "start_time": 1733059200,
  "end_time": 1767225600,
  "expire_time": 1767225600,
  "expire_date": "2026-01-01",
  "status": 1,
  "billing_type": "token",
  "notes": "管理员开通",
  "created_time": 1733059200,
  "updated_time": 1733059200,
  "is_expired": false,
  "remaining_days": 90,
  "name": "Pro 月卡",
  "price": 99.00,
  "duration_days": 30,
  "duration_text": "30 天",
  "sort": 50,
  "recommended": true,
  "description": "...",
  "features": ["GPT-4o 200K tokens/天", "..."],
  "model_limits": { "gpt-4o": { ... } },
  "default_model": "gpt-4o"
}
```

- `is_expired` = `end_time <= now`
- `remaining_days` = `ceil((end_time - now) / 86400)`（已过期则为 0）
- `features` 是 `StringSlice`，落库是 JSON 数组

## `/api/subscription/:id/usage` 字段

返回结构：

```

- `weighted.*` 是 `model.CalculateWeightedUsage` 输出；以 `QuotaPoolCapacity=100` 为满分；任一窗口值 `>= 100` 时 `CheckPlanQuota` 视为耗尽。
- `model_usage[*]` 是 `model.CalcModelUsageDetails`：每个 (model, windowType) 给出当前窗口已用 requests / tokens 与对应百分比。
- `next_reset` 由 `model.CalcNextResetTime` 计算：period 用 `plan.model_limits[<default>].period_h`（默认 5），week=7d，month=30d。

## 前端操作指南

页面：`/subscription`（`web/default-pro/src/views/subscription/Subscription.vue`）。表格列：

- ID / 套餐名 / billing_type / start_time / end_time / status / 操作
- 行内操作：
  - **用量**：调用 `/api/subscription/:id/usage`，渲染加权进度条 + 各模型明细
  - **续费**：用户自助发起新的 `POST /api/order/plan`，根据 `plan.upgrade_mode` 走差价或叠加
  - **取消**：仅管理员可见
- 空状态文案：管理员与普通用户不同
- 「添加订阅」按钮仅 `authStore.isAdmin` 可见

## 实现位置

| 关注点 | 位置 |
|---|---|
| 用户列表 | `controller/subscription.go::GetUserSubscriptions` |
| 当前 plan | `controller/subscription.go::GetCurrentPlan` |
| 用户聚合信息 | `model/plan.go::GetUserSubscriptionInfo` |
| 详细用量 | `controller/subscription.go::GetSubscriptionUsage` |
| 加权计算 | `model/plan_quota.go::CalculateWeightedUsage` |
| 下次重置时间 | `model/plan_quota.go::CalcNextResetTime` |

