---
title: 订单中心（管理员视图）
description: "管理员订单中心：查询 / 标记已付 / 标记退款 / 删除。"
category: pricing
order: 7
---

# 订单中心（管理员视图）

> 在 `/api/order` 上对全站订单做查询、过滤、标记已付、标记退款与删除。前端组件：`web/default-pro/src/views/admin/AdminOrders.vue`。

## 接口一览 / Endpoints

| Endpoint | Method | 鉴权 | 说明 |
|---|---|---|---|
| `/api/order/` | `GET` | Admin | 分页拉取订单列表（默认 `config.ItemsPerPage`） |
| `/api/order/search?keyword=` | `GET` | Admin | 关键字匹配 `order_no` / `pay_trade_no` 前缀 |
| `/api/order/:id` | `GET` | Admin | 订单详情（含 admin 嵌入的 `User` 简档） |
| `/api/order/:id` | `PUT` | Admin | **标记已付 / 标记退款**（详见下方） |
| `/api/order/:id` | `DELETE` | **Root** | 物理删除一条订单 |

实现：`controller/order.go`。

## 通用过滤参数 / Filter Query

`OrderAdminFilter`（`model/order.go`）：

| 参数 | 类型 | 0 / `""` 含义 | 说明 |
|---|---|---|---|
| `type` | `int` | 不过滤 | `1=套餐` / `2=充值` |
| `status` | `int` | 不过滤 | `0=待支付` / `1=已支付` / `2=已取消` / `3=已退款` |
| `source` | `int` | 不过滤 | `1=用户自助` / `2=管理员` |
| `user_id` | `int` | 不过滤 | 精确匹配 |
| `plan_id` | `int` | 不过滤 | 精确匹配 |
| `keyword` | `string` | 不过滤 | 模糊匹配 `order_no` 或 `pay_trade_no`（`LIKE 'kw%'`） |

注意 `status` 用空字符串表示「全部」而非 `0`，因为 `0` 与 `OrderStatusPending` 重叠。

## 标记已付 / `PUT /api/order/:id`

请求体：

```json
{
  "status": 1,
  "pay_method": "offline",      // 可选覆盖
  "pay_trade_no": "TX-2026..."  // 可选，管理员手动录入
}
```

服务端按 `req.Status` 分支：

- `status=1`：写 `pay_method` / `pay_trade_no`，然后调 `model.ActivatePackageByOrder(o, OrderUpgradeModeStack)` 激活套餐；
  - **关键**：管理员手动激活始终走 `stack`（叠加）模式，不走差价；
  - `pay_method` 在 `offline` / `bank` / `wechat` / `alipay` / `free` 下都会立即激活；
  - 充值订单（`OrderTypeTopup=2`）走 `model.ActivateTopupByOrder`，幂等；
- `status=3`：调 `model.MarkOrderRefunded(o)`，仅翻状态位，**不回退已发放的 quota**（计划中 TODO）；
- 其它值：拒绝。

成功消息：`订单已支付，套餐已激活` 或 `订单已标记为退款`。

## 标记退款 / `PUT /api/order/:id` (status=3)

仅状态翻转，不会触发套餐撤销或 quota 回退。
充值订单退款后 user 余额仍按原值发放，仅 order 表的 `status=3` 留作审计。
前端 UI：仅当 `o.status === 1` 时展示「退款」按钮（红色 popconfirm）。

## 删除 / `DELETE /api/order/:id`

仅 Root 可见；调用 `o.Delete()` 直接物理删除行（不级联到 `user_plans`，但 `user_plans.order_id` 是弱引用）。

## 前端操作指南 / Frontend Guide

- 顶部欢迎条 + 搜索框；左侧下拉：`订单类型`（套餐 / 充值）、`状态`、`来源`。
- 顶部右侧：「重置筛选」「刷新」。
- 列表行：ID / `order_no`（截断 tooltip） / 类型 chip / 用户摘要（display_name + `#user_id`）/ 套餐或金额（type=1 时解析 `plan_info` 取套餐名）/ 金额 / 支付方式 / 状态 chip / 来源 / 下单时间 / 操作。
- 行操作：
  - **查看**：打开详情弹窗，列出全部字段（含 `pay_time` / `pay_trade_no` / 嵌入的 `User`）；
  - **标记已付**（仅 `status=0` 可见）：弹窗选 `pay_method`、填 `pay_trade_no`（可空），提交后立即激活并刷新；
  - **退款**（仅 `status=1` 可见）：popconfirm 后调 `PUT { status: 3 }`；
  - **删除**（Root only，`status !== 1` 时可见）：popconfirm 后 `DELETE`。
- 列表分页：`[10, 20, 50]`；前端用 `offset + pageSize` 兜底渲染至少还有下一页。

## 接口实现 / Implementation Pointers

| 关注点 | 位置 |
|---|---|
| 订单 CRUD handler | `controller/order.go` |
| 管理员标记已付 / 退款 | `controller/order.go::MarkOrderPaid` |
| 过滤 + 嵌入 user | `model/order.go::OrderAdminFilter` / `enrichOrdersWithUserBrief` |
| 激活套餐 | `model/order_payment.go::ActivatePackageByOrder` |
| 激活充值 | `model/topup.go::ActivateTopupByOrder` |
| 标记退款（仅状态位） | `model/order_payment.go::MarkOrderRefunded` |
| 路由 | `router/api.go` |