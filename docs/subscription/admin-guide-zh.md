---
title: 订阅管理（管理员）
description: "套餐 CRUD、强制开通订阅、调整 / 撤销订阅。"
category: subscription
order: 3
---

# 订阅管理（管理员）

> 管理员对套餐（plan）和用户订阅（user_plan）的全部操作。实现：`controller/plan.go`、`controller/subscription.go`。

## 套餐 / Plan

数据模型：`model.Plan`。字段见 [套餐定价](./plan) 与 `model/plan.go`。

### 接口 / Endpoints

| Endpoint | Method | Auth | 说明 |
|---|---|---|---|
| `/api/plan/` | `GET` | Admin | 分页拉取（`config.ItemsPerPage`） |
| `/api/plan/search?keyword=` | `GET` | Admin | `name LIKE keyword%` |
| `/api/plan/:id` | `GET` | Admin | 单条详情 |
| `/api/plan/` | `POST` | Admin | 新建 |
| `/api/plan/` | `PUT` | Admin | 更新 |
| `/api/plan/:id` | `DELETE` | Admin | 删除 |
| `/api/plan/public` | `GET` | Public | 仅返回 `status=PlanStatusEnabled` 的套餐列表（用户下单页用） |
| `/api/plan/public/:id` | `GET` | Public | 公开详情 |

### 校验

`plan.Insert` / `plan.Update` 之前服务端会调 `ValidateDefaultModel`：当 `default_model` 非空时，必须能在 `model_limits` 找到对应键；否则报错：

- `default_model 'xxx' is set but model_limits is empty`
- `default_model 'xxx' is not found in model_limits`

### `model_limits` 格式

存 `text` 列，值为 JSON 对象：

```json
{
  "gpt-4o":    { "period_h": 5, "request_period": 100, "request_week": 500, "request_month": 2000, "token_period": 50000,  "token_week": 250000, "token_month": 1000000 },
  "claude-3.5-sonnet": { "period_h": 5, "request_period": 50, "request_week": 200, "request_month": 800, "token_period": 30000, "token_week": 120000, "token_month": 480000 }
}
```

- `period_h` 用于计算 period 窗口索引（`CalcWindowIndex`）；缺省回退为 5
- `request_*` / `token_*` 三种窗口各一对上限；任一为 0 表示该窗口对该维度不限

### features

`Features StringSlice` 字段（`model/plan.go`）：JSON 数组 `["...", "..."]`，前端展示为权益列表。

## 用户订阅 / User subscription

### 强制开通（绕过支付）

`POST /api/subscription/`（`controller/subscription.go::AddSubscription`）请求体：

```json
{
  "user_id": 42,
  "plan_id": 7,
  "billing_type": "token",
  "duration_days": 0,
  "notes": "商务合作赠送",
  "pay_method": "free"
}
```

流程：

1. 校验 `user_id` / `plan_id` 非空；`pay_method` 默认 `free`；白名单：`wechat` / `alipay` / `bank` / `offline` / `free`
2. 校验 plan 存在且 `status=PlanStatusEnabled`；校验用户存在
3. 调 `model.CreatePlanOrder` 写入 `Order`（`source=admin`，作为审计行），notes 默认为 "管理员开通"
4. **立刻** 调 `model.ActivatePackageByOrder(order, OrderUpgradeModeStack)` 激活订阅；管理员 grant 无论 pay_method 是什么都立即生效
5. 返回新 `user_plan`（含 `Plan`）+ 订单行

### 调整

`PUT /api/subscription/` 请求体（`UpdateSubscriptionRequest`）：

```json
{ "id": 88, "end_time": 1767225600, "status": 1, "billing_type": "token", "notes": "延期 30 天" }
```

非空字段被更新；`updated_time` 自动刷新。成功后调 `model.CacheDeleteUserActivePlans(userId)` 失效 Redis 缓存。

### 撤销

`DELETE /api/subscription/:id` 是硬删除：先 `First` 再 `Delete`，并清缓存。撤销后该用户立即失去该 plan 的窗口配额（无补偿）。

## 公开套餐 / Public

`/api/plan/public` 与 `/api/plan/public/:id` 不需要鉴权，供用户侧订阅页与下单流程调用。`/api/plan/`（admin）会返回 `status=0` 的下架套餐，方便管理。

## 周期任务 / Scheduled task

`model.ExpireUserPlans` 由后台启动（具体入口见 `main.go`）周期性把 `status=1 AND end_time <= now` 的 `user_plans` 翻成 `status=0`，让 `CheckPlanQuota` 不再选择它。

## 前端操作指南 / Frontend Guide

- **套餐管理**：`/setting/pricing` → "套餐" Tab（`web/default-pro/src/views/setting/PlanSetting.vue`）。列表 + 弹窗编辑 `name` / `description` / `price` / `duration_days` / `duration_text` / `sort` / `status` / `recommended` / `features[]` / `model_limits` JSON / `default_model`
- **订阅管理**：`/subscription`（`web/default-pro/src/views/subscription/Subscription.vue`）。管理员视角显示 user 列与 "添加订阅" 按钮；用户视角只读
- **添加订阅弹窗**：用户下拉从 `GET /api/user/search` 拉取；套餐下拉从 `GET /api/plan/` 拉取；`pay_method` 选项固定为 `free` / `offline` / `wechat` / `alipay` / `bank`

## 实现位置 / Implementation Pointers

| 关注点 | 位置 |
|---|---|
| 套餐 CRUD | `controller/plan.go` |
| 套餐模型 | `model/plan.go::Plan` |
| 用户订阅 CRUD | `controller/subscription.go` |
| 订单 + 激活 | `model/order_payment.go::CreatePlanOrder` / `ActivatePackageByOrder` |
| 缓存失效 | `model/plan.go::CacheDeleteUserActivePlans` |
| 周期过期 | `model/plan.go::ExpireUserPlans` |
