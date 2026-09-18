---
title: 运营仪表盘 API
description: "/api/admin/dashboard/* 端点。"
category: api
order: 15
---
# 运营仪表盘 API

`/api/admin/dashboard/*` 提供给管理后台首页使用的全站聚合数据。所有端点均要求 **Admin** 权限，且 `success: false` 时返回 `{success, message}`。

## 端点一览

| 接口 | 方法 | 权限 | 说明 |
|------|------|------|------|
| `/api/admin/dashboard/overview` | GET | Admin | KPI 总览（用户/资源/额度/收入） |
| `/api/admin/dashboard/charts` | GET | Admin | 站点级按日 × 模型聚合（`[]LogStatistic`） |
| `/api/admin/dashboard/top-users` | GET | Admin | 活跃用户排行榜 |


## 1. KPI 总览

**接口：** `GET /api/admin/dashboard/overview`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| range | string | 时间窗口，可选 `today` / `7d` / `30d` / `all`，默认 `7d` |

**返回示例：**

```json
{
  "success": true,
  "message": "",
  "data": {
    "users": {
      "total": 1024,
      "enabled": 900,
      "disabled": 80,
      "deleted": 44,
      "new_today": 5,
      "new_7d": 30,
      "new_30d": 120,
      "active_7d": 210
    },
    "tokens": { "total": 3200, "enabled": 2800 },
    "channels": { "total": 18, "enabled": 15 },
    "plans": { "total": 5, "enabled": 4 },
    "redemptions": { "total": 500, "used": 200, "unused": 300 },
    "subscriptions": { "total": 220, "active": 180, "expired": 40 },
    "quota": { "today": 120000, "week": 900000, "month": 3600000, "total": 12000000 },
    "revenue": {
      "total": 19999.50,
      "topup": 8000.00,
      "subscription": 11999.50,
      "refund": 0.00
    },
    "range": "7d",
    "generated_at": 1718000000
  }
}
```

**字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| users.total | int64 | 用户总数（含禁用/已删除） |
| users.enabled | int64 | 启用中的用户数 |
| users.disabled | int64 | 手动禁用的用户数 |
| users.deleted | int64 | 软删除的用户数 |
| users.new_today / new_7d / new_30d | int64 | 今日 / 近 7 日 / 近 30 日新增 |
| users.active_7d | int64 | 近 7 日有消费日志且 `request_count > 0` 的用户数 |
| tokens.total / enabled | int64 | 令牌总数 / 启用数 |
| channels.total / enabled | int64 | 渠道总数 / 启用数 |
| plans.total / enabled | int64 | 套餐总数 / 上架数 |
| redemptions.total / used / unused | int64 | 兑换码总数 / 已使用 / 未使用 |
| subscriptions.total / active / expired | int64 | 用户订阅总数 / 激活中 / 已过期 |
| quota.today / week / month / total | int64 | 各窗口消费日志的 quota 合计 |
| revenue.total / topup / subscription | float64 | 收入合计 / 充值 / 套餐（单位：元） |
| revenue.refund | float64 | 退款合计（单位：元） |
| range | string | 实际生效的窗口名 |
| generated_at | int64 | 服务端生成时间戳 |

> 收入按 `orders.pay_time` 在 `range` 区间内过滤 `status=1`（已支付）和 `status=3`（已退款）的金额。


## 2. 全站图表数据

**接口：** `GET /api/admin/dashboard/charts`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| range | string | 时间窗口，可选 `today` / `7d` / `30d` / `all`，默认 `7d`；`all` 收敛为最近 30 天 |

**返回示例：**

```json
{
  "success": true,
  "message": "",
  "data": [
    {
      "day": "2026-09-10",
      "model_name": "gpt-4o",
      "request_count": 120,
      "quota": 250000,
      "prompt_tokens": 80000,
      "completion_tokens": 50000
    },
    {
      "day": "2026-09-10",
      "model_name": "claude-3-5-sonnet",
      "request_count": 35,
      "quota": 70000,
      "prompt_tokens": 22000,
      "completion_tokens": 14000
    }
  ]
}
```

**字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| day | string | 日期（YYYY-MM-DD，UTC 截断到天） |
| model_name | string | 模型名 |
| request_count | int | 当天该模型的请求数 |
| quota | int | 当天该模型消费的 quota 合计 |
| prompt_tokens | int | 输入 Token 合计 |
| completion_tokens | int | 输出 Token 合计 |

> 返回结构与 `model.SearchLogsByDayAndModel`（即 `GET /api/user/dashboard`）完全一致，前端可以复用同一套图表绘制逻辑。


## 3. 活跃用户排行榜

**接口：** `GET /api/admin/dashboard/top-users`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| range | string | 时间窗口，可选 `today` / `7d` / `30d` / `all`，默认 `7d` |
| limit | int | 返回行数，1-200，默认 20 |

**返回示例：**

```json
{
  "success": true,
  "message": "",
  "data": {
    "range": "7d",
    "items": [
      {
        "id": 12,
        "username": "alice",
        "email": "alice@example.com",
        "request_count": 1230,
        "quota": 800000,
        "balance": 5000000,
        "current_plan_name": "Pro Plan"
      }
    ]
  }
}
```

**字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 用户 ID |
| username | string | 用户名 |
| email | string | 邮箱 |
| request_count | int64 | 窗口内的请求数 |
| quota | int64 | 窗口内消费的 quota 合计 |
| balance | int64 | 当前账户剩余 quota |
| current_plan_name | string | 当前激活的最早到期订阅的套餐名 |

> 排序：先按 `request_count` 降序，再按窗口内消费 quota 降序。仅统计 `users.status != deleted` 且 `users.request_count > 0` 的用户。