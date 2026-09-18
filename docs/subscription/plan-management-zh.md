---
title: 套餐管理
description: "套餐（Plan）的增删改查、上下架、特性与计费字段维护。"
category: subscription
order: 5
---

# 套餐管理

> 在 `/api/plan` 上做套餐（Plan）CRUD；上线 / 下架切换；前端组件 `web/default-pro/src/views/setting/PlanSetting.vue`。

## 数据模型 / Data Model

`model.Plan`（`plans` 表）：

| 字段 | 类型 | 说明 |
|---|---|---|
| `name` | `varchar(100)` | 套餐展示名 |
| `description` | `text` | 列表卡片下方描述 |
| `price` | `decimal(10,2)` | 售价（元） |
| `duration_days` | `int` | 有效期长度 |
| `duration_text` | `varchar(50)` | 自定义展示（如「月卡」「季卡」） |
| `status` | `int` | `1=上架` / `0=下架` |
| `recommended` | `bool` | 推荐标记 |
| `sort` | `int` | 列表排序（升序） |
| `features` | `text`（JSON） | 特性列表（`StringSlice`，前端以多行编辑） |
| `model_limits` | `text`（JSON） | `{ "gpt-4": { request_period, token_period, period_h, ... }, ... }` |
| `default_model` | `varchar(100)` | 新订阅默认模型；必须出现在 `model_limits` 内，否则保存会被后端校验拒绝 |

`Plan.ValidateDefaultModel()`：若 `default_model` 非空，必须能在 `model_limits` 键集中找到，否则报「`default_model 'X' is set but model_limits is empty`」或「`default_model 'X' is not found in model_limits`」。

## 接口一览 / Endpoints

| Endpoint | Method | 鉴权 | 说明 |
|---|---|---|---|
| `/api/plan/list` | `GET` | Public | 拉取所有上架套餐（用户前端用） |
| `/api/plan/detail/:id` | `GET` | Public | 单套餐详情（用户前端用） |
| `/api/plan/` | `GET` | Admin | 分页拉取全部套餐（默认 `config.ItemsPerPage`） |
| `/api/plan/search?keyword=` | `GET` | Admin | 关键字模糊匹配（`name LIKE keyword%`） |
| `/api/plan/:id` | `GET` | Admin | 单套餐详情 |
| `/api/plan/` | `POST` | **Root** | 新建套餐（必填 `name`） |
| `/api/plan/` | `PUT` | **Root** | 更新套餐；`ValidateDefaultModel()` 同步校验 |
| `/api/plan/:id` | `DELETE` | **Root** | 删除套餐 |
| `/api/plan/current` | `GET` | User | 拉取当前用户的活跃订阅（`user_plans`），并把内嵌 `Plan` 字段平铺到顶层 |

实现：`controller/plan.go`；前端 `@/api/plan.js`。

## 上下架 / Toggling Publish

`PlanSetting.vue` 提供行内「启用 / 禁用」按钮，对当前行就地调 `PUT /api/plan/` 把 `status` 在 `1` ↔ `0` 之间切换。
`status=0` 的套餐在 `GetPublicPlans` 会被过滤，所以不再展示给自助下单用户。

## Features / Model Limits 编辑

- `features`：`StringSlice` 类型，对前端暴露为 JSON 数组。`PlanSetting.vue` 用动态行 `<a-input>` 编辑器（`featuresFromRecord` / `sanitizeFeaturesList` 工具函数 `web/default-pro/src/utils/plan.js`）。
- `model_limits`：JSON 对象；保存为 raw JSON 字符串。结构示意：

```json
{
  "gpt-4o": {
    "period_h": 5,
    "request_period": 100,
    "request_week": 1000,
    "request_month": 5000,
    "token_period": 50000,
    "token_week": 500000,
    "token_month": 2000000
  },
  "claude-3.5-sonnet": {
    "period_h": 5,
    "request_period": 50,
    "token_period": 20000,
    "request_week": 500,
    "token_month": 2500
  }
}
```

字段语义：
- `period_h` 决定 `WindowTypePeriod`（小时级）窗口大小。
- `request_period/week/month` = 该窗口内的调用次数上限；`token_*` 同理。

## 前端操作指南 / Frontend Guide

- 顶部表格列：ID / 名称 / 价格 / 有效天数 / 默认模型 / 推荐 ★ / 状态 tag / 排序 / 操作。
- 「新增套餐」按钮打开 900px 弹窗（`arco-modal-plan-wide`）：
  1. 套餐名称、描述（多行 textarea）；
  2. 价格（`precision=2`）、`duration_days`、`duration_text`；
  3. 排序、状态（启用/禁用）、推荐开关；
  4. 特性说明（动态行编辑，回车追加新行，删除按钮）；
  5. 默认模型（必须先在 `model_limits` 内出现）；
  6. 模型限额（JSON 字符串 textarea，含占位符 `{"gpt-4": 1000, "gpt-3.5": 5000}`）。
- 行内「启用 / 禁用」点击即提交；「删除」走二次确认。
- `sort` 升序展示，多套餐排序靠它控制首页「推荐置顶」效果。

## 接口实现 / Implementation Pointers

| 关注点 | 位置 |
|---|---|
| CRUD handler | `controller/plan.go` |
| 公开拉取 | `controller/plan.go::GetPublicPlans` / `GetPublicPlanDetail` |
| 用户当前订阅 | `controller/plan.go::GetCurrentPlan` → `model.GetActiveUserPlansByUserId` |
| `StringSlice` 双向桥接 | `model/plan.go` |
| 限额结构 | `model/plan.go::ModelLimitRule` |
| 前端编辑工具 | `web/default-pro/src/utils/plan.js` |
| 路由 | `router/api.go` |