---
title: 管理员手动充值
description: "管理员直接给用户加额度；以及管理员视角的充值订单查询。"
category: pricing
order: 8
---

# 管理员手动充值

> 绕过支付通道，直接给指定用户加 `quota`（走 `model.IncreaseUserQuota`），同时落 `LogTypeTopup` 审计日志。

入口路由：`/admin/orders`（管理员订单中心）内的「标记已付」流程会创建 `OrderTypeTopup` 订单并即时激活。
另外有一个独立入口（兼容旧版） `POST /api/user/topup`，由 `controller/user.go::AdminTopUp` 提供；新代码建议走订单中心的「标记已付」流程以保留审计行。

## 接口 / Endpoints

| Endpoint | Method | 鉴权 | 说明 |
|---|---|---|---|
| `/api/user/topup` | `POST` | Admin | 直接累加 user.quota（旧路径，保留兼容） |
| `/api/topup/order` | `POST` | User | 用户自助充值下单（需要 `topup.enabled=true`） |
| `/api/order/:id` | `PUT` | Admin | 标记已付（管理员支付渠道，详见 orders 文档） |
| `/api/order/` | `GET` | Admin | 在订单中心按 `type=2` 过滤充值订单 |

> 用户自助入口请参考用户侧文档；本页只覆盖管理员视角。

## POST `/api/user/topup`（兼容路径）

```json
{ "user_id": 42, "quota": 100000, "remark": "客服补偿" }
```

- 调 `model.IncreaseUserQuota(user_id, quota)` 直接累加（带 Redis 缓存失效）。
- `remark` 为空时自动生成 `通过 API 充值 <LogQuota(quota)>`。
- 写一条 `LogTypeTopup` 审计日志（`model.RecordTopupLog`）。

> 没有 `Order` 行；若需要审计追溯，请改用 `POST /api/order` + 标记已付的流程（创建 `Order` 行再激活）。

## 通过「订单中心」手动开通 / Activate Topup via Order Center

推荐路径：

1. 在 `/admin/orders` 顶部点「新建订单」或在前端「充值」表单走 `POST /api/order` 创建 `OrderTypeTopup=2` 订单。
2. 用户完成线下转账（`pay_method=bank` / `offline`），管理员在订单中心点「标记已付」。
3. 后端 `PUT /api/order/:id { status:1, pay_method, pay_trade_no }` 调 `model.ActivateTopupByOrder`：
   - 解析 `plan_info` 拿到下单时快照的 `bonus_quota`；
   - 调 `IncreaseUserQuota` 累加；
   - 把 `status` / `pay_status` / `pay_time` / `pay_trade_no` 落库。

`ActivateTopupByOrder` 幂等：已 `status=1` 的订单直接返回 nil，不会重复加 quota。

## 充值订单查询 / Querying Topup Orders

`GET /api/order/?type=2`（或 admin 页面下拉筛选 `type=Topup`）即可仅看充值订单。
列表字段与套餐订单一致：`order_no` / `user_id` / `plan_id`（充值时为 0） / `amount` / `pay_method` / `status` / `source`。

`plan_info` 字段是 JSON 字符串，结构：

```json
{ "amount": 100.00, "preset_amount": 100.00, "bonus_quota": 100000, "exchange_rate": 1000 }
```

充值订单的 `plan_id` 始终为 `0`，但 `plan_info.bonus_quota` 才是最终到账额度（避免汇率调整造成的二次解释）。

## 前端操作指南 / Frontend Guide

- 手动加额度（兼容路径）：管理员旧 UI 在「用户」行内提供「+ 额度」入口（如有），提交 `{ user_id, quota, remark }`。
- 标记已付：订单中心订单状态 `pending` 时，行尾出现「标记已付」按钮。
- 退款：已支付订单出现「退款」按钮（仅翻状态，不回退 quota）。
- 删除：Root 可对非 paid 订单执行删除。

## 接口实现 / Implementation Pointers

| 关注点 | 位置 |
|---|---|
| 兼容路径手动加额度 | `controller/user.go::AdminTopUp` |
| 用户自助下单 | `controller/topup.go::CreateTopupOrder` |
| 订单中心标记已付 | `controller/order.go::MarkOrderPaid` |
| 充值激活（幂等） | `model/topup.go::ActivateTopupByOrder` |
| 充值下单 | `model/topup.go::CreateTopupOrder` |
| 路由 | `router/api.go` |