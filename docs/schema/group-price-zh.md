---
title: 分组折扣
description: 用户组对特定模型的折扣倍率，未命中时按 1.0 计费。
category: schema
order: 6
---

# 分组折扣（GroupPrice）

## 这是什么

`GroupPrice` 让「某个用户组」对「某个模型」拥有折扣倍率（如 0.8 = 八折）。它与 [ModelPrice](/schema/model-price) 组合使用，最终扣费金额：

```
最终扣费 = 调用消耗 × ModelPrice × GroupPrice(group, model)
```

## 在哪里看到

- **后台 → 分组折扣**：列表、新增、批量调整
- 用户登录后，其用户组 + 命中的 GroupPrice 会显示在个人中心的「当前折扣」

## 关键字段（运维视角）

| 字段 | 含义 | 设置后影响 |
|---|---|---|
| `group_name` | 用户组名 `default` / `vip` / `svip` | 与用户的 `group` 字段匹配才生效 |
| `model_name` | 模型名 `gpt-4o` / `claude-sonnet-4` 等；空字符串表示该用户组**所有模型**的默认折扣 | 空时为全组默认折扣 |
| `discount` | 折扣倍率，1.0 = 原价 | 0.8 = 八折；1.2 = 加价 20% |

## 命中规则

1. 先查「(user_group, model) 完全匹配」的记录
2. 没匹配上，再查「(user_group, 空模型)」记录（该用户组的默认折扣）
3. 都没匹配上，按 1.0 计费

## 常见用法

- **VIP 群体整体八折**：

  | `group_name` | `model_name` | `discount` |
  |---|---|---|
  | `vip` | _(空)_ | 0.8 |

- **VIP 对某款高级模型额外优惠**：

  | `group_name` | `model_name` | `discount` |
  |---|---|---|
  | `vip` | _(空)_ | 0.8 |
  | `vip` | `gpt-4o` | 0.6 |

## 不推荐用法

- 不要把所有用户组都设置 1.0（默认就是 1.0，无意义）
- 不要给同一对 `(group_name, model_name)` 设置多条记录（重复会报唯一键错误）

## 相关页面

- [分组折扣（业务概念）](/pricing/group-price)
- [分组折扣管理（后台）](/pricing/group-price-management)

## 相关 API

- `GET /api/group_price/` — 列表
- `POST /api/group_price/` — 新增
- `PUT /api/group_price/` — 更新
- `DELETE /api/group_price/:id` — 删除