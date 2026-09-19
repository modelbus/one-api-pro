---
title: 套餐定价
description: 套餐是什么、怎么设置、用户怎么买。
category: pricing
order: 3
---

# 套餐定价

> 用户在前台订阅页看到的「包月套餐」是怎么定义的。

## 这是什么

`Plan` 是「可订阅商品」的模板。它定义了：

- 叫什么 / 多贵 / 多长有效期
- 哪些模型能用、能用多少
- 是不是推荐套餐

用户在前台看到的套餐就是它。

## 在哪里看到

- **后台 → 套餐管理**（管理员）：新增、上下架
- **前台 → 订阅**（用户）：列出当前上架的 Plan

## 创建套餐时需要填的

| 项 | 含义 | 设置后影响 |
|---|---|---|
| 名称 | 用户看到的名字 | 修改后所有位置同步更新 |
| 价格 | ¥ | 用户下单金额 |
| 有效期 | 天数 | 过期后用户进入「套餐已过期」状态，但仍可调用到额度耗尽 |
| 折扣倍率 | 数字 | 用户订阅期间内调用的最终折扣（默认 1.0 = 无折扣） |
| 模型配额 | JSON | 每个模型的窗口配额（见下） |
| 默认模型 | 模型名 | 不在配额里的请求被路由到该模型；空则 422 |
| 描述 | 富文本 | 前台展示 |
| 推荐 | 开关 | 前台套餐卡显示「★ 推荐」徽章 |
| 状态 | 上架 / 下架 | 下架后用户不能新订；已有订阅不受影响 |

## 模型配额（model_limits）JSON

每个套餐对每个模型可以设三档窗口配额：

```json
{
  "gpt-4o": {
    "period_h": 5,
    "request_period": 100,
    "request_week": 500,
    "request_month": 2000,
    "token_period": 50000,
    "token_week": 250000,
    "token_month": 1000000
  }
}
```

| 字段 | 含义 |
|---|---|
| `period_h` | period 窗口小时数（默认 5） |
| `request_period` / `request_week` / `request_month` | 三窗口调用次数上限（`0` = 不限） |
| `token_period` / `token_week` / `token_month` | 三窗口 token 上限 |

`model_limits` 为 null 时套餐对该模型不限（按量计费路径）。

详见 [计费规则](../subscription/billing-rules)。

## 用户购买流程

```
用户在前台 /pricing 选套餐
    ↓
POST /api/order/plan  (生成订单号 TB / UP)
    ↓
支付 → 回调 → ActivatePackageByOrder
    ↓
创建 Subscription (end_time = now + duration_days)
```

管理员也可以「手动开通」绕过支付，常用于客服补偿。

## 升降级

切换套餐时：

- **差价模式**（默认）：按差价扣费，旧订阅剩余额度折算到新订阅
- **叠加模式**：付全款开新订阅，旧订阅继续生效

管理员在 [套餐业务配置](../subscription/plan-settings) 切换。

## 怎么上架新套餐

1. 后台 → 套餐管理 → 新增
2. 填名称、价格、有效期、折扣
3. 填 model_limits（每个模型的三窗口配额）
4. 填 default_model（建议填套餐的主推模型）
5. 状态改为「上架」
6. 保存

建议填完先测一遍再上架：开通一个测试订阅，验证调用、配额、升级、降级。

## 与用户组的关系

套餐**不直接绑定用户组**。一个 VIP 用户买普通套餐，扣费规则：

```
调用消耗 × ModelPrice × 套餐折扣 × (用户组折扣 / 用户组默认折扣)
```

## 常见问题

- **改了 model_limits 老订阅会变吗**：不会。老订阅沿用订阅时的快照
- **默认模型必须填吗**：套餐要限制模型时必填，否则不限模型时为空
- **删除套餐后老订阅还能用吗**：能，删除只影响新订阅

## 相关页面

- [套餐（数据结构）](/schema/plan)
- [套餐管理（后台）](../subscription/plan-management)
- [套餐升降级](../subscription/upgrade-downgrade)
- [我的订单（用户）](../user/orders)

## 相关 API

- `GET /api/plan/` — 列表（管理员）
- `GET /api/plan/public` — 列表（前台公开）
- `POST /api/plan/` — 新增（Root）
- `PUT /api/plan/` — 更新（Root）
- `DELETE /api/plan/:id` — 删除（Root）