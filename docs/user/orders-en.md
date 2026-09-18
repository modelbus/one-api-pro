---
title: My Orders
description: "Plan, subscription upgrade and topup orders."
category: user
order: 5
---

# My Orders

> Plan, subscription upgrade and topup orders.
> 套餐、订阅升级、充值订单。

Entry: `/orders` (`web/default-pro/src/views/user/Orders.vue`). Covers three kinds of business orders:

入口：`/orders`（`web/default-pro/src/views/user/Orders.vue`）。涵盖三类业务订单：

## Order Types / 订单类型

| `type` | Name / 名称 | Triggered by / 触发场景 | Prefix / 订单号前缀 |
| --- | --- | --- | --- |
| 1 | Plan subscription / upgrade | User buys or upgrades a plan | `TB` new / `UP` differential |
| 2 | Top-up | Online top-up | `TP` |
| - | Admin-placed | Admin creates for a user / free grant | same as above |

Order number format: `{prefix} + yyyyMMddHHmmss + 6-digit random` (22 chars), produced by `model.GenerateOrderNo(prefix)`.

订单号格式：`{prefix} + yyyyMMddHHmmss + 6 位随机数`（22 字符），由 `model.GenerateOrderNo(prefix)` 生成。

## Status State Machine / 状态机

```text
        ┌───────────────┐
        │ 0 pending      │
        └──────┬────────┘
               │ callback verified / admin marked paid
               ▼
        ┌───────────────┐
        │ 1 paid         │
        └──────┬────────┘
               │ admin refunded
               ▼
        ┌───────────────┐
        │ 3 refunded     │
        └───────────────┘

  also: user self-cancel ─► 2 canceled (terminal)
```

| Status | Meaning | Triggered by |
| --- | --- | --- |
| 0 Pending | Order created, unpaid | New order |
| 1 Paid | Callback received / admin marked | `processNotify` / `MarkOrderPaid` |
| 2 Canceled | User cancelled | `POST /api/order/self/:id/cancel` |
| 3 Refunded | Admin refunded | `MarkOrderRefunded` |

## List Filters / 列表过滤

UI top tabs:

UI 顶部 Tab：

- **All / 全部**: `/api/order/self` without `type`.
- **Plan / 套餐**: `?type=1`.
- **Topup / 充值**: `?type=2`.

## Actions / 操作

| Action | When | API |
| --- | --- | --- |
| Pay | status = 0 | `POST /api/order/self/:id/pay` — reuses `controller/buildPayInfo` for `pay_url / qr_code` |
| View | any | `GET /api/order/self/:id` |
| Cancel | status = 0 | `POST /api/order/self/:id/cancel`, idempotent (errors when already paid) |

## Detail Fields / 详情字段

```text
orderNo           order number
type              1=plan, 2=topup
source            1=self-service, 2=admin-placed
planName          plan name (empty for topup)
amount            CNY
payMethod         wechat / alipay / bank / offline / free
payTradeNo        provider trade number
status            0/1/2/3
createdAt         created at
paidAt            paid at
refundedAt        refunded at
note              admin note
```

## Upgrade Differential / 升级差价

When `OrderUpgradeModePriceDiff` (default) is active and the user already has an active subscription:

当 `OrderUpgradeModePriceDiff`（默认）启用且用户已有有效订阅时：

```text
diff_amount = newPlan.price - max(remaining_value_of_current_plan, 0)
orderNo    = "UP" + timestamp + random
amount     = diff_amount
```

`OrderUpgradeModeStack` (stack): pay the full `newPlan.price` for a new subscription; the old one keeps running until expiry.

`OrderUpgradeModeStack`（叠加）：按 `newPlan.price` 全额开新订阅，旧订阅在过期前继续生效。

The exact logic lives in `model/order_payment.go::CreatePlanOrder`.

具体策略在 `model/order_payment.go::CreatePlanOrder` 中实现。

## Payment Availability / 支付渠道可用性

The "Pay" button branches on `payInfo.status`:

「支付」按钮按 `payInfo.status` 分级处理：

| `payInfo.status` | UI behavior / UI 行为 |
| --- | --- |
| `success` | Show QR / redirect; modal embeds `qr_code` or `pay_url` |
| `warning` | "Channel not enabled" or channel error toast; button stays but is disabled |
| `error` | Reject; prompt the user to switch payment method |

## FAQ / 常见问题

- **Order stays "Pending"**: async notify hasn't arrived; admin can mark paid manually, or check `payment.wechat.config` etc.
  订单一直显示「待支付」：异步回调未到；管理员后台「订单管理」手动 `MarkOrderPaid`，或检查支付渠道配置。
- **Negative upgrade differential**: the current plan's remaining value is larger than the new plan; UI still allows the order but should prompt the user.
  差价升级金额是负数：通常是模式为 `price_diff` 且旧套餐剩余价值高于新套餐全价；UI 仍允许下单，但请提示用户确认。
- **Payment succeeded but order stays inactive**: see [Troubleshooting · notify verification failed](/en/faq/troubleshooting#async-notify-verification-failed).

Next: [Chat Playground](/en/user/chat) · [Subscriptions](/en/subscription/overview) · [Top-up](/en/pricing/topup).