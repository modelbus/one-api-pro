---
title: 我的订单
description: "套餐、订阅升级、充值订单。"
category: user
order: 5
---

# 我的订单

> 套餐、订阅升级、充值订单。

入口：`/orders`（`web/default-pro/src/views/user/Orders.vue`）。涵盖三类业务订单：

## 订单类型

| `type` | 名称| 触发场景| 订单号前缀|
| --- | --- | --- | --- |
| 1 | 套餐订阅 / 升级 | 用户购买套餐 / 升级套餐 | `TB` 新订 / `UP` 差价升级 |
| 2 | 充值 | 在线支付充值 | `TP` |
| - | 管理员后台下单 | 管理员代客下单 / 免费赠送 | 同上 |

订单号格式：`{prefix} + yyyyMMddHHmmss + 6 位随机数`（22 字符），由 `model.GenerateOrderNo(prefix)` 生成。

## 状态机

```text
        ┌───────────────┐
        │ 0 待支付       │   pending
        └──────┬────────┘
               │ callback verified / admin marked paid
               ▼
        ┌───────────────┐
        │ 1 已支付       │   paid
        └──────┬────────┘
               │ admin refunded
               ▼
        ┌───────────────┐
        │ 3 已退款       │   refunded
        └───────────────┘

  同时：用户自助取消 ─► 2 已取消（终态）
        also: user self-cancel ─► 2 canceled (terminal)
```

| 状态值| 含义| 触发|
| --- | --- | --- |
| 0 待支付| 下单未付款 | 用户新建订单 |
| 1 已支付| 收到回调 / 管理员标记 | `processNotify` / `MarkOrderPaid` |
| 2 已取消| 用户主动取消 | `POST /api/order/self/:id/cancel` |
| 3 已退款| 管理员退款 | `MarkOrderRefunded` |

## 列表过滤

UI 顶部 Tab：

- **全部 / All**：`/api/order/self` 不传 `type`；
- **套餐订单 / Plan**：`?type=1`；
- **充值订单 / Topup**：`?type=2`。

## 操作

| 操作| 触发条件| API |
| --- | --- | --- |
| 支付| 状态 = 0 | `POST /api/order/self/:id/pay`，复用 `controller/buildPayInfo` 拿 `pay_url|
| 查看| 任何状态 | `GET /api/order/self/:id` |
| 取消| 状态 = 0 | `POST /api/order/self/:id/cancel`，幂等（已支付返回错误） |

## 详情字段

```text
orderNo           订单号 / order number
type              1=套餐 2=充值
source            1=用户自助 2=管理员下单
planName          套餐名称（充值订单为空）
amount            金额（元）
payMethod         wechat / alipay / bank / offline / free
payTradeNo        支付渠道流水号
status            0/1/2/3
createdAt         下单时间
paidAt            支付完成时间
refundedAt        退款时间
note              管理员备注
```

## 升级差价

当 `OrderUpgradeModePriceDiff`（默认）启用且用户已有有效订阅时：

```text
diff_amount = newPlan.price - max(remaining_value_of_current_plan, 0)
orderNo    = "UP" + timestamp + random
amount     = diff_amount
```

`OrderUpgradeModeStack`（叠加）：按 `newPlan.price` 全额开新订阅，旧订阅在过期前继续生效。

具体策略在 `model/order_payment.go::CreatePlanOrder` 中实现。

## 支付渠道可用性

「支付」按钮按 `payInfo.status` 分级处理：

| `payInfo.status` | UI 行为|
| --- | --- |
| `success` | 显示二维码 / 跳转链接；弹窗内嵌 `qr_code` 或 `pay_url` |
| `warning` | 弹窗提示「该支付方式尚未启用」或渠道错误；按钮保留但不可点击 |
| `error` | 直接拒绝；提示「请更换支付方式」 |

## 常见问题

- **订单一直显示「待支付」**：异步回调未到；登录管理员后台在「订单管理」手动 `MarkOrderPaid`，或检查支付渠道配置（`payment.wechat.config` 等）。
- **差价升级金额是负数**：通常是模式为 `price_diff` 且旧套餐剩余价值高于新套餐全价；UI 仍允许下单，但请提示用户确认。
- **支付成功但订单一直未激活**：见 [故障排查 · 异步通知校验失败](/zh/faq/troubleshooting#异步通知校验失败)。

下一步 / Next: [Chat Playground](/zh/user/chat) · [套餐订阅](/zh/subscription/overview) · [充值](/zh/pricing/topup)。
