---
title: 订阅
description: 用户当前生效的订阅实例（基于某个 Plan）。
category: schema
order: 8
---

# 订阅（Subscription）

## 这是什么

`Subscription` 是 [Plan](/schema/plan) 的实例化记录。用户在前台订阅某个 Plan 后，系统就会生成一条 Subscription，并按 Plan 的规则扣费。

## 在哪里看到

- **前台 → 个人中心 → 我的订阅**：当前生效的订阅、剩余天数、已用 / 总额度
- **后台 → 订阅管理**：所有用户的订阅、到期时间、当前状态
- **调用日志**：调用命中的订阅记录在每条 log 里

## 关键字段（运维视角）

| 字段 | 含义 | 设置后影响 |
|---|---|---|
| `user_id` | 所属用户 | 关联到 [User](/schema/user) |
| `plan_id` | 关联套餐 | 关联到 [Plan](/schema/plan)；改了 Plan 不会回溯到现有订阅 |
| `start_time` | 订阅开始时间 | 通常由下单时间自动写入 |
| `end_time` | 订阅失效时间 | 过期后用户进入「套餐过期」状态，但**仍可调用**到 `used_quota == quota` |
| `used_quota` | 已用额度 | 本订阅累计消耗；仅记录 |
| `status` | 状态 | 生效中 / 已过期 / 已取消 |

## 订阅对扣费的影响

调用扣费的完整公式：

```
final = consumption × ModelPrice × GroupPrice(group, model)
            × Subscription.discount (if active subscription covers this model)
```

- 没有活跃订阅 → 不享受订阅折扣
- 订阅过期 → 同样不享受折扣（但不影响用户调用）
- 多个订阅同时存在（升级过渡期）→ 仅当前生效的那条生效

## 订阅 vs 用户额度

- 用户额度（`User.quota`）：扣调用费时优先扣这里
- 订阅额度（`Plan.quota` / `Subscription` 配额）：订阅期间内的「总量」，扣到 0 后降级为原价
- 订阅和用户额度是**两条独立的扣费源**，具体见 [计费规则](/subscription/billing-rules)

## 相关页面

- [订阅概览](/subscription/overview)
- [用户使用指南](/subscription/user-guide)
- [订阅管理（后台）](/subscription/subscription-management)

## 相关 API

- `GET /api/subscription/self` — 当前用户的订阅
- `GET /api/subscription/` — 列表（管理员）
- `POST /api/subscription/` — 新建（管理员手动分配）
- `PUT /api/subscription/` — 更新
- `DELETE /api/subscription/:id` — 删除