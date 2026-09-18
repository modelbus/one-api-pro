---
title: 套餐定价
description: "`Plan` 数据模型、有效期、模型覆盖、分组绑定与升降级语义。"
category: pricing
order: 3
---

# 套餐定价

> 套餐（Plan）是订阅的定价与配额载体。实现：`model/plan.go::Plan`、`controller/plan.go`、`web/default-pro/src/views/setting/PlanSetting.vue`。

## 数据模型 / Data model

`model.Plan`（`plans` 表）字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| `name` | `varchar(100)` | 套餐名（必须；`Insert` 时为空则报错） |
| `description` | `text` | 套餐描述 |
| `price` | `decimal(10,2)` | 价格（元） |
| `duration_days` | `int` | 有效天数（默认 30）；创建订阅时 `end_time = start_time + duration_days * 86400` |
| `duration_text` | `varchar(50)` | 展示文本（如 "30 天"、"季度"） |
| `status` | `int` | `PlanStatusEnabled=1` / `PlanStatusDisabled=0`；`/api/plan/public` 仅返回 `Enabled` |
| `recommended` | `bool` | 前端 "推荐" 标签 |
| `sort` | `int` | 升降级比较键；值越大越高级 |
| `features` | `StringSlice` (text JSON) | 权益文案列表 |
| `model_limits` | `text` JSON | 各模型窗口配额（见下） |
| `default_model` | `varchar(100)` | 不在 limits 的请求被路由到该模型；`ValidateDefaultModel` 校验必须在 `model_limits` 存在 |

## `model_limits` 格式 / `model_limits` shape

`model/model_plan.go::ModelLimitRule`：

```json
{
  "gpt-4o":            { "period_h": 5, "request_period": 100, "request_week": 500, "request_month": 2000, "token_period": 50000,  "token_week": 250000, "token_month": 1000000 },
  "claude-3.5-sonnet": { "period_h": 5, "request_period": 50,  "request_week": 200, "request_month": 800,  "token_period": 30000,  "token_week": 120000, "token_month": 480000 }
}
```

- `period_h` = period 窗口小时数（缺省 5）
- `request_*` / `token_*` = 三个窗口的请求次数 / token 上限；`0` 表示不限
- `model_limits` 为 null 时 `CheckPlanQuota` 直接视为 usable（无窗口限制）

详见 [计费规则](../subscription/billing-rules)。

## features / 权益列表

`StringSlice` 是 `plan.go` 自定义的 JSON / 文本双向桥接：

- 写库：`MarshalJSON` 输出数组；`Value()` 序列化为 JSON 字符串
- 读库：`Scan` 优先按 JSON 数组解析，失败则按 `\n` 拆分（兼容历史纯文本）
- nil 时 `MarshalJSON` 输出 `[]` 而非 `null`

## 接口 / Endpoints

| Endpoint | Method | Auth | 说明 |
|---|---|---|---|
| `/api/plan/` | `GET` | Admin | 分页 |
| `/api/plan/search?keyword=` | `GET` | Admin | `name LIKE kw%` |
| `/api/plan/:id` | `GET` | Admin | 详情 |
| `/api/plan/` | `POST` | Admin | 新建（`name` 必填；`ValidateDefaultModel`） |
| `/api/plan/` | `PUT` | Admin | 更新（同上） |
| `/api/plan/:id` | `DELETE` | Admin | 删除 |
| `/api/plan/public` | `GET` | Public | 仅 `status=Enabled`，用户侧套餐列表 |
| `/api/plan/public/:id` | `GET` | Public | 公开详情 |

## 分组绑定 / Group binding

套餐自身不直接绑定到具体用户组；用户组判定由渠道的 `group` 字段控制。套餐的"高级感"主要通过：

- `recommended=true` — 列表标记推荐
- `sort` — 升降级比较（同 `sort` 拒绝；`sort` 小的不能升到 `sort` 大的之上）
- `default_model` — 当用户请求一个不在 `model_limits` 里的模型时，重写到该模型；为空则直接 422（`PlanQuotaCheck`）

## 与订单 / 订阅的关系 / Plan ↔ Order ↔ Subscription

- 用户自助下单：`POST /api/order/plan` → `model.CreatePlanOrder` → `buildPayInfo` 返回预支付参数 → 支付回调 → `ActivatePackageByOrder(order, mode)`
- 管理员 grant：`POST /api/subscription/` → 立即 `ActivatePackageByOrder(order, OrderUpgradeModeStack)`
- 升降级模式由系统设置 `plan.upgrade_mode` 决定（默认 `price_diff`）

详见 [套餐升降级](../subscription/upgrade-downgrade) 与 [订阅管理（管理员）](../subscription/admin-guide)。

## 前端操作指南 / Frontend Guide

页面：`/setting/pricing` → "套餐" Tab（`web/default-pro/src/views/setting/PlanSetting.vue`）。

- 列表列：name / price / duration_days / duration_text / sort / status / recommended
- 编辑弹窗：name / description / price / duration_days / duration_text / sort / status / recommended / features[] / model_limits JSON / default_model
- 「推荐」开关（★）：`recommended=true`
- 「启用 / 禁用」按钮切换 `status`
- 删除带二次确认；删除套餐不会回滚已存在的 `user_plans`（历史订阅仍按当时快照走）

## 实现位置 / Implementation Pointers

| 关注点 | 位置 |
|---|---|
| 数据模型 | `model/plan.go::Plan` / `ModelLimitRule` |
| 校验 | `model/plan.go::ValidateDefaultModel` |
| CRUD | `controller/plan.go` |
| 公开接口 | `controller/plan.go::GetPublicPlans` / `GetPublicPlanDetail` |
| 创建订单 | `model/order_payment.go::CreatePlanOrder` |
| 激活 | `model/order_payment.go::ActivatePackageByOrder` |
| 升降级 | `model/order_payment.go::CalculateUpgradePrice` |
