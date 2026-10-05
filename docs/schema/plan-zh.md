---
title: 套餐
description: 可购买的订阅模板：额度、有效期、计费配置。
category: schema
order: 7
---

# 套餐（Plan）

## 这是什么

`Plan` 是「可订阅商品」的模板。用户在前台下单后，会基于某个 Plan 创建一条 [Subscription](/schema/subscription)。

## 在哪里看到

- **后台 → 套餐管理**：列表、新增、上下架
- **前台 → 订阅**：列出当前上架的 Plan
- **后台 → 套餐业务配置**：升级模式、过期策略等全局规则

## 你需要为每个套餐设置什么

| 项 | 含义 | 设置后影响 |
|---|---|---|
| `name` | 名称 | 用户看到的名字 |
| `price` | 价格 | ¥ |
| `billing_type` | 计费维度 | `token`（按 Token）或 `request`（按请求次数）；即套餐内模型的限额按哪个维度统计 |
| `virtual_amount` | 套餐额度 | 元；套餐内模型共享的虚拟余额，消耗完整个套餐不可用；`0` 表示不限额度 |
| `model_limits` | 模型限额 | JSON；声明套餐覆盖哪些模型及其滚动窗口（5 小时 / 周 / 月）限额；格式写错会被拒绝保存 |
| `duration_days` | 有效期 | 天数；过期后用户进入「套餐已过期」状态，但仍可调用到额度耗尽 |
| `description` | 描述 | 富文本 / 多行；用户在订阅页能看到 |
| `status` | 状态 | 上架 / 下架；下架后用户不能新订 |
| `recommended` | 推荐 | 开关；前台套餐卡显示「推荐」徽章 |

## 套餐与订阅的区别

- **Plan**：模板，定义了「这种订阅叫什么、多少钱、多少额度、有效期多长」
- **Subscription**：实例化，某个用户在某个时刻订阅了某个 Plan 后生成的「实例」。一个用户可以多次订阅同一个或不同的 Plan

## 套餐升级 / 差价

当用户从 Plan A 切换到 Plan B（更高价格）时：

- 如果 B 价格 > A 价格，按差价生成一条订单（订单号前缀 `UP`）
- 剩余天数 / 额度按比例折算到新套餐

升级规则在 [套餐业务配置](/subscription/plan-settings) 中维护。

## 套餐额度与失效规则

套餐覆盖范围由 `model_limits` 里的模型名决定，这些模型共享 `virtual_amount` 这笔虚拟余额：

- 每笔订阅调用的费用按「模型定价 × 分组折扣」累加到 `user_plans.used_amount`
- `used_amount >= virtual_amount` → 套餐整体不可用，请求回落用户全局余额按量计费
- 任一滚动窗口（5 小时 / 周 / 月）的加权用量达到 100% → 同样整体不可用
- 未在 `model_limits` 中列出的模型不参与套餐，直接按全局余额计费
- 订阅消费会计入 `users.used_quota` / `request_count`，但**不扣减** `users.quota`

> `virtual_amount = 0` 表示不限额度（兼容存量套餐），此时只由窗口限额控制。

## 相关页面

- [套餐定价](/pricing/plan)
- [套餐管理（后台）](/subscription/plan-management)
- [套餐业务配置](/subscription/plan-settings)
- [套餐升降级](/subscription/upgrade-downgrade)

## 相关 API

- `GET /api/plan/` — 列表（管理员）
- `GET /api/plan/list` — 列表（前台公开）
- `POST /api/plan/` — 新增
- `PUT /api/plan/` — 更新
- `DELETE /api/plan/:id` — 删除