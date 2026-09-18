---
title: 兑换码使用（用户）
description: "`/redeem` 页面、调用接口、成功与失败提示。"
category: redemption
order: 3
---

# 兑换码使用（用户）

> 用户在 `/redeem` 页输入 32 位 `key` 后调 `POST /api/user/topup`，由 `model.Redeem` 在事务内加 quota 并落 `LogTypeTopup` 日志。

## 接口

| 项 | 值 |
|---|---|
| Method | `POST` |
| Path | `/api/user/topup`（亦用作管理员手动加 quota，见 [管理员视图](#管理员视图)） |
| Auth | User（自助兑换）/ Admin（手动补偿） |
| Body | `{ "key": "..." }` |

成功响应：

```

`data` 为本次到账的 quota 数值（= `redemption.quota`）。

## 失败模式

`model.Redeem` 在事务内会返回 `errors.New(...)`，由 controller 包装成 `success:false, message`：

| 触发 | message |
|---|---|
| `key` 为空 | `未提供兑换码` |
| `userId` 为 0 | `无效的 user id` |
| `key` 不存在 | `无效的兑换码` |
| `status != Enabled`（已用 / 已停用） | `该兑换码已被使用` |
| 任意 SQL 错误 | `兑换失败，<原因>` |

事务保证：并发请求同一 key 时只有一个事务能进入关键段，其它事务在 `FOR UPDATE` 上等待后看到 `status=Used` 并立即失败。

## 管理员视图

管理员同样可以调 `POST /api/user/topup`（`controller/user.go::AdminTopUp`），不传 `key` 而是直接传 `{ "user_id": 42, "quota": 10000, "remark": "..." }`。这条路径不走 `Redeem`，直接 `IncreaseUserQuota(user_id, quota)`，并通过 `RecordTopupLog` 写 `LogTypeTopup` 日志；与兑换码语义独立，但对外是同一个端点。

> 如果仅做"加额度"操作，建议直接走管理员 grant 订单中心（`OrderTypeTopup=2` → `ActivateTopupByOrder`），可保留完整审计订单行；`AdminTopUp` 是兼容旧版遗留。

## 前端操作指南

页面：`/redeem`（由前端自定义组件承载）。

- 一个 `<a-input>` + 「兑换」按钮
- 提交后根据 `success` 字段弹 `Message.success('兑换成功，到账 X')` 或 `Message.error(message)`
- 成功后通常把 `users.quota` 重新拉取或直接刷新顶部余额

## 实现位置

| 关注点 | 位置 |
|---|---|
| 用户兑换入口 | `controller/user.go::TopUp` |
| 兑换事务 | `model/redemption.go::Redeem` |
| 管理员手动加 | `controller/user.go::AdminTopUp` |
| 日志 | `model/log.go::RecordLog` (`LogTypeTopup`) |

