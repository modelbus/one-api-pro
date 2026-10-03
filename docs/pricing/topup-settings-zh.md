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
| `exchange_rate` | 1 元 = X quota | 影响自定义金额的换算；必须 ≥ `QuotaPerUnit` |
| `presets` | 预设金额列表 | 每条 `{amount, bonus_quota}` |

## 推荐配置

启用 + 允许自定义 + 几个预设（以默认 `QuotaPerUnit=500000` 即 1 元 = 1 元额度为例）：

| amount | bonus_quota |
|---|---|
| 10 | 5000000 |
| 50 | 25000000 |
| 100 | 50000000 |
| 500 | 250000000 |

`exchange_rate = 500000`（1 元 = `QuotaPerUnit` 额度，充 1 元到账 1 元）。

> `QuotaPerUnit` 是系统级换算基准（默认 500000），可在系统设置中修改；
> `exchange_rate` 不得低于它，否则会出现「充 1 元到账不足 1 元」。
> 若要赠送，请按倍率设置（如 `1000000` = 2 倍）。

## 怎么改

后台 → 充值业务配置 → 改字段 → 保存。

校验：
- 每个预设 `amount > 0`，`bonus_quota >= 0`
- 不能有重复金额
- `exchange_rate >= QuotaPerUnit`（低于基准值的存量数据会在读取时自动归一）

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
- 低于基准值的存量 `exchange_rate`（如历史 `1` / `100000`）会在读取时自动归一为 `QuotaPerUnit`，无需数据迁移
- 改预设金额只影响**新订单**；已有订单不影响
- 预设金额不能重复（unique constraint）

## 常见问题

- **配置都改了但用户看不到充值入口**：检查是否启用了 `enabled`；并且至少一个支付通道已启用
- **汇率改完没生效**：检查是否保存成功；保存后立即对**新订单**生效
- **为什么保存时报「兑换比例不能低于 500000」**：`exchange_rate` 低于系统 `QuotaPerUnit` 会被拒绝；请按基准值或更高的赠送倍率填写
- **怎么让预设金额带小数**：`amount` 支持两位小数（如 `9.99` 元）

## 相关页面

- [充值（业务概念）](./topup)
- [充值管理（后台）](./topup-management)
- [支付通道](./payment)

## 相关 API

- `GET /api/setting/topup` — 读取配置（Root）
- `PUT /api/setting/topup` — 更新配置（Root，**整体覆盖**）