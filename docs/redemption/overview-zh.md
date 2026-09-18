---
title: 兑换码概览
description: "兑换码的类型、用途与典型场景。"
category: redemption
order: 1
---

# 兑换码概览

> 兑换码是 One API Pro 中"按码发放额度"的载体：管理员批量生成，外部渠道（推广、合作、活动）发放给终端用户，用户在 `/redeem` 输入后到账 `quota`。

## 类型

数据模型 `model.Redemption`（`redemptions` 表）字段 `key` 是 32 位 UUID 字符串，唯一索引。

| 状态 | 常量 | 含义 |
|---|---|---|
| 1 | `RedemptionCodeStatusEnabled` | 已启用，可被兑换 |
| 2 | `RedemptionCodeStatusDisabled` | 已停用（管理员手动） |
| 3 | `RedemptionCodeStatusUsed` | 已使用（用户兑换后由事务翻转） |

> 状态值刻意避开 `0`：默认值不能落在 "已使用" 语义上。

`is_count` 字段（`Count`，落库时 `gorm:"-:all"` 不持久化）只在批量生成时使用：服务端按 `count` 循环生成对应行，每条独立 UUID。

## 用途

| 场景 | 用法 |
|---|---|
| 推广 / 邀请 | 一码一额，限定批次名（`name`）；用户自助兑换 |
| 充值补偿 | 管理员手动发码给投诉 / 故障用户 |
| 营销活动 | 批量生成（`count ≤ 100/批`），按渠道发放 |

> 兑换码只是把 `users.quota` 加上对应额度，**不会** 创建或影响 `user_plans`，与套餐正交。

## 单次 vs 多次

`redemptions` 表本身每条对应一个唯一 UUID：单条被兑换即置为 `Used=3`，是天然的一次性码。多次使用需要管理员生成多条（同 `name` 不同 `key`）。没有"一码多次"语义；如需分发额度到多用户，按数量生成即可。

## 兑换成功后的副作用

`model.Redeem(ctx, key, userId)`（`model/redemption.go:55`）在事务内：

1. `SELECT … FOR UPDATE` 锁住 `redemptions.key` 行
2. 校验 `status == Enabled`；否则返回 `该兑换码已被使用`
3. `UPDATE users SET quota = quota + redemptions.quota`
4. `redemption.RedeemedTime = now`，`status = Used`，`Save`
5. 写一条 `LogTypeTopup` 日志：`通过兑换码充值 <LogQuota(redemption.Quota)>`

事务提交后才会对用户可见；任一步骤失败整笔回滚，避免并发兑换与重复到账。

## 与充值订单的关系

兑换码到账额度走的是 `users.quota += redemptions.quota`（`Redeem`），与 `OrderTypeTopup=2` 的订单激活路径（`model.ActivateTopupByOrder` → `IncreaseUserQuota`）最终都把额度累加到 `users.quota`。两者可独立使用；订单可生成对账所需的金额 / 兑换率信息，兑换码仅含纯 `quota`。

## 实现位置

| 关注点 | 位置 |
|---|---|
| 数据模型 | `model/redemption.go::Redemption` |
| 兑换事务 | `model/redemption.go::Redeem` |
| CRUD | `controller/redemption.go` |
| 用户侧入口 | `controller/user.go::TopUp`（`POST /api/user/topup`） |

