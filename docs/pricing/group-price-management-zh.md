---
title: 分组定价管理
description: 管理员如何在后台给用户组设置模型折扣。
category: pricing
order: 4
---

# 分组定价管理

> 后台 → 分组定价。给不同用户组设置模型折扣。

## 在哪里

后台 → 分组定价。

## 列表看到什么

每行展示：

- 用户组名（如 `default` / `vip` / `svip`）
- 模型名（空 = 该组所有模型的默认折扣）
- 折扣倍率（`0.8` = 八折）
- 操作按钮

## 怎么新增

1. 列表右上角「新增」
2. 填：
   - 用户组（必填，下拉候选由系统维护）
   - 模型名（必填，下拉候选来自 [ModelPrice](/pricing/model-price) 中已启用的模型）
   - 折扣倍率（`0.8` = 八折，`1.2` = 加价 20%；留空用 `1.0`）
3. 保存

## 怎么编辑

行内「编辑」修改 `discount`。改后立即生效。

## 命中规则

1. 先查「(user_group, model) 完全匹配」
2. 没有则查「(user_group, 空模型)」 = 该用户组的默认折扣
3. 都没有 → 1.0

## 推荐配置

| group | model | discount | 效果 |
|---|---|---|---|
| `vip` | _(空)_ | `0.8` | vip 用户对所有模型 80% 收费 |
| `vip` | `gpt-4o` | `0.5` | vip 用户对 gpt-4o 50% 收费（叠加在 0.8 之上？不，分组模型优先级更高） |
| `svip` | _(空)_ | `0.7` | svip 用户对所有模型 70% 收费 |

## 注意事项

- 同一对 `(group, model)` 不能重复添加（会报唯一键冲突）
- 修改一个组的折扣，会即时影响所有用该组的用户
- 不要把所有组都设 1.0（默认就是 1.0）

## 怎么把组分配给用户

组在 [系统设置](../misc/system-settings) 配置，新用户默认 `default` 组。管理员在 [用户管理](../user/user-management) 里给具体用户改组。

## 常见问题

- **设了折扣但用户反映没生效**：检查用户当前的 `group` 字段（不是用户名）
- **改了折扣立即生效吗**：是。下次调用按新折扣扣
- **套餐折扣和组折扣什么关系**：套餐折扣优先于组折扣。具体见 [计费规则](/subscription/billing-rules)

## 相关页面

- [分组折扣（业务概念）](/pricing/group-price)
- [分组折扣（数据结构）](/schema/group-price)
- [用户管理（管理员）](../user/user-management)

## 相关 API

- `GET /api/group_price/` — 列表（Root）
- `POST /api/group_price/` — 新增（Root）
- `PUT /api/group_price/` — 更新（Root）
- `DELETE /api/group_price/:id` — 删除（Root）