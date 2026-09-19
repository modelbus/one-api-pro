---
title: 管理员手动充值
description: 管理员如何直接给用户加额度（绕过支付通道）。
category: pricing
order: 8
---

# 管理员手动充值

> 客服补偿、活动赠送、对账修复 —— 看这里。

## 两种方式

### 方式 1：通过订单中心（推荐，有审计）

1. 后台 → 订单管理 → 新建订单（选「充值」类型）
2. 线下转账（银行转账）或管理员手动确认
3. 行内「标记已付」→ 系统创建充值订单 + 加额度 + 写审计日志

**优点**：每笔充值都有 `Order` 记录，便于对账。

### 方式 2：直接加额度（兼容旧版）

仅用于特殊情况：内部测试 / 紧急补偿 / 修复对账。

```
POST /api/user/topup
{
  "user_id": 42,
  "quota": 100000,
  "remark": "客服补偿"
}
```

`remark` 是审计说明，必填。空时会自动填 `通过 API 充值 <LogQuota(quota)>`。

**缺点**：没有 `Order` 行，事后审计困难。

## 推荐用方式 1

默认走订单中心：每笔充值都能追溯到具体来源（哪个银行流水 / 哪个客服工单）。

## 怎么查充值订单

后台 → 订单管理 → 类型下拉选「充值」：

- `order_no`：`TP` 前缀
- `user_id`：充值目标用户
- `amount`：充值金额（元）
- `plan_info` 是 JSON，含 `bonus_quota`（实际到账额度）
- `pay_method`：微信 / 支付宝 / 银行 / 线下 / 免费
- `status`：待支付 / 已支付 / 已取消 / 已退款
- `pay_time`：支付完成时间

## 怎么标记已付

详见 [订单管理](./order-management)。

## 常见问题

- **用户反映付款了但余额没到**：看订单详情 `pay_time` / `pay_trade_no`；手动标记已付
- **加了多次重复**：系统幂等，已支付的订单不会被重复激活
- **退款后用户余额还在**：是的，系统不自动回退。需要手动 [调整用户余额](../user/user-management)

## 相关页面

- [充值（业务概念）](./topup)
- [充值业务配置](./topup-settings)
- [订单管理](./order-management)
- [用户管理（管理员）](../user/user-management)

## 相关 API

- `POST /api/user/topup` — 直接加额度（旧路径，兼容保留）
- `POST /api/order` — 创建充值订单（推荐）
- `PUT /api/order/:id` — 标记已付（Admin）
- `GET /api/order/?type=2` — 查充值订单