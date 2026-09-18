---
title: 运营仪表盘
description: "KPI、图表、用户排行聚合视图。"
category: misc
order: 1
---

# 运营仪表盘

> 全站 KPI 概览、按日 × 模型的图表、活跃用户排行榜与营收看板。

入口路由：`/admin/dashboard`，对应前端组件 `web/default-pro/src/views/admin/AdminDashboard.vue`。
该页面对所有管理员（`role >= 10`）可见；普通用户会被 router guard 重定向到 `/dashboard`。

## 接口一览

仪表盘数据由三个独立的 Admin 端点提供，全部挂在 `/api/admin/dashboard` 下，使用 `AdminAuth` 中间件保护：

| Endpoint | Query | 返回 |
|---|---|---|
| `GET /api/admin/dashboard/overview` | `range=today\|7d\|30d\|all`，默认 `7d` | 单次聚合 KPI（不含图表数据） |
| `GET /api/admin/dashboard/charts`   | `range=today\|7d\|30d\|all`（`all` 收敛为最近 30 天） | `[]LogStatistic`，按 `day × model_name` 聚合 |
| `GET /api/admin/dashboard/top-users` | `range=today\|7d\|30d\|all`，`limit=1..200`（默认 20） | `[{ id, username, email, request_count, quota, balance, current_plan_name }]` |

后端实现见 `controller/admin/dashboard.go` 与 `model/admin_dashboard.go`。

## Overview 字段

`overview` 接口返回 `AdminDashboardOverview`（KPI-only）：

| 分组 | 关键字段 | 说明 |
|---|---|---|
| `users` | `total| 用户总数按 `status` 枚举分桶；新增按 `created_at >= start_of_(day/week/month)` 聚合；`active_7d = DISTINCT user_id FROM logs WHERE type=2 AND created_at >= 7d AND request_count > 0` |
| `tokens` / `channels` / `plans` | `total| 资源计数（按 `status=1` 计 enabled） |
| `redemptions` | `total| 兑换码三类状态计数 |
| `subscriptions` | `total| `user_plans` 表的状态聚合 |
| `quota` | `today| 通过 `SumUsedQuota(LogTypeConsume, …)` 汇总，窗口由后端按天/周/月/全量展开 |
| `revenue` | `total| 来自 `orders` 表：`status=1` 时按 `type` 拆分（`OrderTypeTopup=2` / `OrderTypePlanSubscription=1`），退款金额来自 `status=3` |
| `range` / `generated_at` | 元数据 | `range` 为规范化后的窗口键；`generated_at` 为服务端生成时间戳（秒） |

## Charts

`charts` 接口直接复用 `model.SearchLogsByDayAndModel` 的输出结构 `[]LogStatistic`：

```ts
interface LogStatistic {
  day: string              // YYYY-MM-DD（按 created_at 时区归零）
  model_name: string       // 模型名
  request_count: number    // 命中条数（仅 type=2 consume 日志）
  quota: number            // sum(quota)
  prompt_tokens: number    // sum(prompt_tokens)
  completion_tokens: number// sum(completion_tokens)
}
```

> 与 `GET /api/user/dashboard` 的 shape 完全一致；前端可共用同一份图表构建逻辑（`AdminDashboard.vue::chartData`）。
> 注意 SQL 同时兼容 MySQL / PostgreSQL / SQLite（`DATE_FORMAT` / `TO_CHAR(date_trunc…)` / `strftime` 三选一）。

`range=30d` 与 `range=all` 都收敛为最近 30 天，避免超长历史拖垮响应。

## Top Users

口径：近 `range` 时间内有 `type=2` 日志、且 `users.request_count > 0` 的用户，按 `request_count DESC, quota DESC` 排序。
`limit` 接受 1–200；越界默认回落为 20。
"当前订阅名称"通过 `batchTopActivePlanNames(userIds)` 批量补齐，避免 N+1；多张激活订阅时只展示最早到期的那一条。

## 前端面板

`AdminDashboard.vue` 顶部欢迎条带日期范围切换（`today / 7d / 30d / all`），下方按 16:8 分栏布局：

- **左 16/24 区块**（按顺序堆叠）
  - **用户 KPI**：8 张卡片（总数 / 今日新增 / 7 日新增 / 30 日新增 / 7 日活跃 / 禁用 / 软删除 / 今日占比）。
  - **营收 KPI**：4 张卡片（总营收 / 充值营收 / 套餐营收 / 退款），单位 `¥`。
  - **资源 KPI**：5 张卡片（令牌 / 渠道 / 套餐 / 兑换码 / 订阅），点击套餐或兑换码卡片可直接跳转到对应管理页。
  - **Quota + 三图**：四张 KPI 卡片（今日 / 7 日 / 30 日 / 总消耗）+ 三张折线图（请求量 / 额度 / Token）。
  - **模型分布**：基于 `chartData` Top 8 模型堆叠柱图，按天堆叠 `quota`。
  - **使用明细**：`day × model` 扁平表格，按日期降序 + `quota` 降序排序。
- **右 8/24 区块**
  - 广告位占位 / 排行榜（卡片式，按 `today / 7d / 30d` 切换，显示 rank / 用户名 / 邮箱 / 请求数 / 消耗 / 当前套餐）/ 系统公告 / 更新日志 / 资源链接。

## 操作提示

- `range` 切换会同时刷新 `overview` / `charts` / `top-users` 三个端点；右上的"上次刷新" chip 显示最后一次成功拉取 `overview` 的本地时间。
- 图表数据统一来自 `/api/admin/dashboard/charts`，三张折线图（请求量 / 额度 / Token）与模型分布共用同一份 `[]LogStatistic`，避免 N+1。
- 所有 KPI 字段均为「已软删除用户不再计入 `users.total`」的视角：`deleted` 单独展示但 `total` 仍按 `COUNT(*)` 计算。

## 接口实现

| 关注点 | 位置 |
|---|---|
| Overview 聚合 SQL | `model/admin_dashboard.go::GetAdminDashboardOverview` |
| Revenue 聚合 | `model/admin_dashboard.go::aggregateRevenue` |
| Charts 按日聚合 | `model/admin_dashboard.go::SearchAdminLogsByDayAndModel` |
| Range 解析 | `model/admin_dashboard.go::ParseAdminDashboardRange` / `ParseAdminChartsRange` |
| Top users + 当前套餐 | `model/admin_dashboard.go::GetAdminTopUsers` / `batchTopActivePlanNames` |
| HTTP handler | `controller/admin/dashboard.go` |
| 路由 | `router/api.go`（`/api/admin/dashboard/*`） |
