---
title: 订阅管理
description: "给指定用户开通 / 调整 / 撤销套餐订阅（user_plans）。"
category: subscription
order: 6
---

# 订阅管理

> 在 `user_plans` 表上为指定用户开通、查询、调整、撤销套餐订阅；管理员可绕过支付直接生效。

入口路由：`/subscription`（同一页面，admin 与 user 复用视图；前端 `web/default-pro/src/views/subscription/Subscription.vue`，按 `authStore.isAdmin` 切换"新建订阅"按钮可见性）。

## 数据模型 / Data Model

`model.UserPlan`（`user_plans` 表）：

| 字段 | 类型 | 说明 |
|---|---|---|
| `user_id` | `int`（索引） | 用户 ID |
| `plan_id` | `int` | 套餐 ID（必填；`Insert()` 拒绝 `plan_id=0`，防止脏行） |
| `order_id` | `int` | 关联订单 ID（管理员开通时为 0） |
| `start_time` / `end_time` | `int64` | unix 秒；`end_time <= now` 视为已过期 |
| `status` | `int` | `UserPlanStatusActive=1` / `UserPlanStatusExpired=0` |
| `billing_type` | `varchar(20)` | `token` 或 `request`，与套餐保持一致 |
| `notes` | `text` | 备注（如「管理员开通」） |

后台定时任务 `ExpireUserPlans()` 把 `status=1 AND end_time <= now` 的行批量改为 `expired`。

## 接口一览 / Endpoints

| Endpoint | Method | 鉴权 | 说明 |
|---|---|---|---|
| `/api/subscription/self` | `GET` | User | 当前登录用户全部订阅（不限状态） |
| `/api/subscription/` | `GET` | Admin | 分页拉取所有订阅，支持 `?user_id=` 与 `?status=` 过滤 |
| `/api/subscription/search?keyword=` | `GET` | Admin | 按 username 前缀模糊搜索 |
| `/api/subscription/:id` | `GET` | Admin | 订阅详情（附 `Plan` 字段） |
| `/api/subscription/:id/usage` | `GET` | User | 拉取当前订阅的限额 / 实际用量 / 加权使用 / 下次重置时间 |
| `/api/subscription/` | `POST` | Admin | 管理员开通套餐（详见下方） |
| `/api/subscription/` | `PUT` | Admin | 调整 `end_time` / `status` / `billing_type` / `notes` |
| `/api/subscription/:id` | `DELETE` | Admin | 撤销订阅 |

实现：`controller/subscription.go`。

## 管理员开通 / `POST /api/subscription/`

请求体：

```json
{
  "user_id": 42,
  "plan_id": 7,
  "billing_type": "token",
  "duration_days": 0,            // 可选：0 表示沿用套餐默认值
  "notes": "商务合作赠送",
  "pay_method": "free"           // wechat | alipay | bank | offline | free
}
```

服务端会：

1. 校验 `pay_method` 必须是合法值，缺省 `free`；
2. 校验套餐存在且 `status=1`；
3. 走 `model.CreatePlanOrder` 创建一条 `source=admin` 的 `Order`，作为审计凭证；
4. **立即**调 `model.ActivatePackageByOrder(order, OrderUpgradeModeStack)` 激活订阅（不论 `pay_method` 是 `free` / `offline` / `bank` / `wechat` / `alipay`）。管理员手动开通始终即时生效；
5. 把新生成的 `user_plan` 行（含 `plan` 字段）随 `order` 一起返回。

`ActivatePackageByOrder` 的升级模式始终是 `stack`（管理员主动开通不应用差价逻辑）。

## 调整订阅 / `PUT /api/subscription/`

```json
{
  "id": 88,
  "end_time": 1767225600,
  "status": 1,
  "billing_type": "token",
  "notes": "延期 30 天"
}
```

更新成功后调用 `model.CacheDeleteUserActivePlans(userId)` 清掉 Redis 缓存（`user_plans:<id>`），下一次请求会重新加载新数据。

## 撤销订阅 / `DELETE /api/subscription/:id`

物理删除（`user_plan` 表），同时清掉用户缓存。
撤销后该用户的订阅额度 / 限额随之失效，不会自动重置。

## 限额查询 / `/api/subscription/:id/usage`

返回体：

```jsonc
{
  "subscription": { /* user_plan + plan */ },
  "usage": [ /* PlanUsage rows: per (model, window_type, window_index) */ ],
  "weighted": { "period": 0.42, "week": 0.15, "month": 0.08 },  // 加权使用率
  "limits":   { "gpt-4o": { request_period: 100, token_period: 50000, period_h: 5, … }, … },
  "model_usage": { /* 每模型按窗口拆分的 used vs limit */ },
  "next_reset": { "period": 1767225600, "week": 1767225600, "month": 1767225600 },
  "now": 1735660800,
  "start_time": 1733059200,
  "billing_type": "token"
}
```

`weighted` 是 `model.CalculateWeightedUsage` 把"按请求 / 按 token"两种计费维度的消耗折算成一个综合比率；用于在 UI 上展示"使用进度条"。

## 前端操作指南 / Frontend Guide

- 顶部表格列：ID / 用户（admin 可见）/ 套餐名 / 计费类型 / 起 / 止 / 状态 / 操作。
- 行内操作：
  - **用量**：跳到 `/subscription/usage?id=...` 或弹窗展示 `weighted` 进度；
  - **续期**：admin 打开弹窗，直接在 `end_time` 上加 `duration_days`；
  - **立即过期**：二次确认后调用 `PUT` 把 `status` 置 `0`；
  - **删除**：直接删行（不可恢复）。
- 顶部「新增订阅」按钮仅 admin 可见，弹窗字段：`user_id`、`plan_id`（下拉由 `/api/plan/` 提供）、`billing_type`（token / request）、`pay_method`（free / offline / wechat / alipay / bank）、备注。

## 接口实现 / Implementation Pointers

| 关注点 | 位置 |
|---|---|
| Subscription CRUD handler | `controller/subscription.go` |
| 限额 / 用量 | `model/plan.go::GetUserSubscriptionInfo` / `CalculateWeightedUsage` / `CalcModelUsageDetails` / `CalcNextResetTime` |
| 自动过期 | `model/plan.go::ExpireUserPlans` |
| 缓存失效 | `model/plan.go::CacheDeleteUserActivePlans` |
| 路由 | `router/api.go` |