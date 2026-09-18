---
title: 诊断 API
description: "/api/diag/* 诊断端点。"
category: api
order: 16
---
# 诊断 API

`/api/diag/*` 提供管理员专用的诊断端点，用于排查生产环境的数据一致性问题（典型场景：`user_plans.plan_id` 总是显示为 0）。所有端点均要求 **Admin** 权限。

## 端点一览

| 接口 | 方法 | 权限 | 说明 |
|------|------|------|------|
| `/api/diag/subscriptions` | GET | Admin | 比对 `user_plans` 与 `orders`，定位脏数据 |

## 1. 订阅诊断

**接口：** `GET /api/diag/subscriptions`

**权限：** Admin

**说明：** 扫描前 100 条 `user_plans`，把 `user_plan` 行 + `orders` 行 + `users` 简档拼成单行返回，便于排查 `plan_id <= 0` 等异常。本接口只读，不会修改任何数据。

**返回示例：**

```

**字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| total_user_plans | int | 本次扫描到的 `user_plans` 总数（最多 100 行） |
| missing_plan_id | int | `plan_id <= 0` 的行数（疑似脏数据） |
| rows | array | 诊断行，每行包括 `user_plan` / `order_no` / `order_plan_id` / `order_type` / `user` |

**每行字段：**

| 字段 | 类型 | 说明 |
|------|------|------|
| user_plan | object | 原始 `UserPlan` 行 |
| order_no | string | 来源订单号；若 `order_id <= 0` 则为空 |
| order_plan_id | int | 订单上的 `plan_id` |
| order_type | int | 订单 type（1=套餐订阅, 2=在线充值） |
| user | object | `UserBrief`：精简的用户信息（id|

> 仅当 `user_plan.order_id > 0` 时才会尝试加载对应的 `orders` 行；其余字段为 0 / 空字符串。

**典型用法：**

```bash
# 1. 列出最近 100 条 user_plan，定位脏数据
curl -s 'http://localhost:3000/api/diag/subscriptions' -b cookies.txt | jq

# 2. 在 `missing_plan_id > 0` 时，按 user_plan.id 反查触发路径
#    - plan_id = 0：写入时 order.plan_id 缺失/被改写
#    - order_id = 0：孤儿 user_plan（订单被硬删除后未级联清理）
```
