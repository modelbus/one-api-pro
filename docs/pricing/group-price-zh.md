---
title: 分组折扣
description: 让某个用户组对某些模型有折扣。
category: pricing
order: 2
---

# 分组折扣

> 给 VIP 打个八折、给 SVIP 某些模型五折 —— 看这里。

## 这是什么

`GroupPrice` 让「某个用户组」对「某些模型」有折扣倍率。配合 [ModelPrice](/pricing/model-price) 使用：

```
最终扣费 = 调用消耗 × ModelPrice × GroupPrice(group, model)
```

`GroupPrice` 默认 1.0（无折扣）。用户在某个用户组，且该用户组对当前模型有折扣时生效。

## 在哪里看到

后台 → 分组折扣。

## 需要设置什么

每条填：

| 项 | 含义 | 设置后影响 |
|---|---|---|
| `group_name` | 用户组名 | 与用户的 `group` 字段匹配才生效 |
| `model_name` | 模型名 | 空 = 该用户组**所有模型**的默认折扣；填具体模型 = 只针对该模型 |
| `discount` | 折扣倍率 | `1.0`=无折扣；`0.8`=八折；`1.2`=加价 20% |

## 命中规则

1. 先查「(user_group, model) 完全匹配」
2. 没有则查「(user_group, 空模型)」 = 该用户组的默认折扣
3. 都没有 → 1.0

## 常见用法

### VIP 群体整体八折

| group | model | discount |
|---|---|---|
| `vip` | _(空)_ | `0.8` |

### VIP 对高级模型额外优惠

| group | model | discount |
|---|---|---|
| `vip` | _(空)_ | `0.8` |
| `vip` | `gpt-4o` | `0.5` |

### 多个用户组

系统启动时会种入 `default` / `vip` / `svip` 三个组，每个 `discount=1.0`。按需在 UI 上调。

## 怎么改

后台 → 分组折扣 → 新增：

1. 选 `group_name`（如 `vip`）
2. 填 `model_name`（空 = 该组所有模型的默认折扣）
3. 填 `discount`（如 `0.8` = 八折）
4. 保存

## 不要

- 把所有组都设成 1.0（默认就是 1.0）
- 同一对 `(group, model)` 创建多条记录（会报唯一键冲突）

## 与套餐的关系

[订阅套餐](/subscription/overview) 的折扣是另一条独立的折扣路径（`subscription.discount`）。两条折扣同时存在：

- 调用命中订阅 → 用 `subscription.discount`（套餐折扣）
- 调用未命中订阅 → 用 `GroupPrice`（用户组折扣）

## 常见问题

- **设了 VIP 八折但用户反映没生效**：检查用户当前是否在 `vip` 组（个人中心或后台查看）
- **改了折扣立即生效吗**：是。下次调用按新折扣计算
- **一个用户能享受多个折扣吗**：VIP 用户组折扣 × 套餐折扣（如果有）= 实际扣费

## 相关页面

- [分组折扣管理（后台）](./group-price-management)
- [分组折扣（数据结构）](/schema/group-price)
- [ModelPrice](/pricing/model-price)

## 相关 API

- `GET /api/group_price/` — 列表
- `POST /api/group_price/` — 新增
- `PUT /api/group_price/` — 更新
- `DELETE /api/group_price/:id` — 删除