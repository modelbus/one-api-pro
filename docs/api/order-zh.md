---
title: 订单 API
description: "/api/order/* 端点（用户自助与管理端）"
category: api
order: 8
---

# 订单 API

订单 API 同时覆盖用户自助与管理端后台。订单类型：

| 常量 | 值 | 含义 |
| --- | --- | --- |
| `OrderTypePlan` | 1 | 套餐订单（含升级差价 `UP` 与新建 `TB`） |
| `OrderTypeTopup` | 2 | 充值订单（`TP`） |

## 公共约定

- 鉴权：`/api/order/plan`、`/api/order/self/*` 走 Cookie Session 或 Access Token；`/api/order`、`/api/order/:id`、`/api/order/search` 仅管理员（RootAuth）。
- 所有响应使用统一 JSON 格式 `{ success, message, data }`。
- 订单状态：`OrderStatusPending=0`、`OrderStatusPaid=1`、`OrderStatusCancelled=2`、`Refunded=3`。支付状态：`OrderPayStatusPending=0`、`OrderPayStatusPaid=1`、`OrderPayStatusRefunded=-1`。

## 用户自助

### `POST /api/order/plan`

新建套餐订阅订单。差价升级场景下 server 计算差价（按剩余天数比例），订单号 `UP` 前缀；普通新建为 `TB`。

请求体（`CreatePlanOrderRequest`）：

```

- `mode=price_diff`：差价升级；扣除已有订阅剩余天数比例后计算应付金额（订单号 `UP`）。
- `mode=stack`：叠加订阅；按所选套餐全额扣减（订单号 `TB`）。

响应：`{ success, message, data: { id, order_no, amount, ... } }`。

### `GET /api/order/self?type=1|2`

当前用户订单列表，可按类型过滤。`type=1` 套餐、`type=2` 充值。

### `GET /api/order/self/:id`

订单详情，校验当前用户为订单 owner。

### `POST /api/order/self/:id/cancel`

用户主动取消未支付订单（仅 `status=0`）。

### `POST /api/order/self/:id/pay`

用户对自有未支付订单重新发起支付（复用 `buildPayInfo`），常用于支付中断后继续。

## 管理端

### `GET /api/order`

管理端订单列表（`AdminAuth`），分页 + 多维过滤（type / status / source / user_id / plan_id / keyword）。

### `GET /api/order/search?keyword=...`

关键字搜索。

### `GET /api/order/:id`

订单详情（任意用户）。

### `PUT /api/order/:id`

管理端标记订单为「已支付」或「已退款」：

```json
{ "status": 1 | 3, "pay_method": "wechat|alipay|bank|offline|free", "pay_trade_no": "可选" }
```

- `status=1`：标记已支付，离线/银行通道会立即激活套餐（`ActivatePackageByOrder`）。
- `status=3`：标记已退款（仅翻状态，quota 不回退；退款闭环待后续迭代）。

### `DELETE /api/order/:id`

物理删除订单，仅 Root 可操作。

## 订单号规则

`GenerateOrderNo` 在 `model/order_payment.go` 定义：

| 前缀 | 场景 |
| --- | --- |
| `TB` | 套餐新建（`CreatePlanOrder` 中 `mode=stack` 或首次购买） |
| `UP` | 套餐差价升级（`mode=price_diff`） |
| `TP` | 充值（`CreateTopupOrder`） |

## 相关页

- 支付回调：`/api/payment/{wechat,alipay}/notify` 走 `controller/payment.go::processNotify`，按 `order.Type` 分发到 `ActivatePackageByOrder` 或 `ActivateTopupByOrder`。
- 充值下单：`/api/topup/order` 详见 [Topup API](topup.md)。
- 套餐管理：`/api/plan/*` 详见 [Plan API](plan.md)。
