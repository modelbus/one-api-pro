---
title: 充值业务配置
description: 启用 / 关闭在线充值、配置预设金额和汇率。
category: pricing
order: 9
---

# 充值业务配置

> 后台 → 充值业务配置。Root 可见。

## 当前可配项

| 项 | 含义 | 设置后影响 |
|---|---|---|
| `enabled` | 总开关 | 关闭后用户在前台看不到充值入口 |
| `allow_custom` | 允许自定义金额 | 关闭后用户只能选预设金额 |
| `exchange_rate` | 1 元 = X quota | 影响自定义金额的换算；必须 > 0 |
| `presets` | 预设金额列表 | 每条 `{amount, bonus_quota}` |

## 推荐配置

启用 + 允许自定义 + 几个预设：

| amount | bonus_quota |
|---|---|
| 10 | 100000 |
| 50 | 500000 |
| 100 | 1000000 |
| 500 | 5000000 |

`exchange_rate = 100000`（1 元 = 10 万 quota，与 `QuotaPerUnit=500000` 默认值兼容）。

## 怎么改

后台 → 充值业务配置 → 改字段 → 保存。

校验：
- 每个预设 `amount > 0`，`bonus_quota >= 0`
- 不能有重复金额
- `exchange_rate > 0`

## 用户侧影响

| 项 | 用户看到什么 |
|---|---|
| `enabled = false` | 充值入口消失 |
| `allow_custom = false` | 只能点预设金额；输入框隐藏 |
| 预设金额列表 | 决定 chip 数量和顺序 |
| `exchange_rate` | 自定义金额的最终到账额度 |

## 怎么测

1. 启用配置后，让测试账号在前台访问充值页
2. 点预设金额 → 检查到账额度 = `bonus_quota`
3. 输自定义金额 → 检查到账额度 = `amount × exchange_rate`

## 注意事项

- 改汇率只影响**新订单**；已创建的订单沿用下单时的快照汇率
- 改预设金额只影响**新订单**；已有订单不影响
- 预设金额不能重复（unique constraint）

## 常见问题

- **配置都改了但用户看不到充值入口**：检查是否启用了 `enabled`；并且至少一个支付通道已启用
- **汇率改完没生效**：检查是否保存成功；保存后立即对**新订单**生效
- **怎么让预设金额带小数**：`amount` 支持两位小数（如 `9.99` 元）

## 相关页面

- [充值（业务概念）](./topup)
- [充值管理（后台）](./topup-management)
- [支付通道](./payment)

## 相关 API

- `GET /api/setting/topup` — 读取配置（Root）
- `PUT /api/setting/topup` — 更新配置（Root，**整体覆盖**）