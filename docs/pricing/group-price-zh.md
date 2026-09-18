---
title: 分组折扣
description: "`group_price` 表、折扣倍率、模型级覆盖与默认分组。"
category: pricing
order: 2
---

# 分组折扣

> 用户组（`users.group`）对模型计费的乘数；`group_price` 表定义"某用户组对某模型"的折扣倍率，默认 `1.0`（不折扣）。实现：`model/model_price.go::GroupPrice`。

## 数据模型

`model.GroupPrice`（`group_price` 表）：

| 字段 | 类型 | 说明 |
|---|---|---|
| `group_name` | `varchar(32)` | 用户组名（与 `users.group` 对应），联合唯一索引第一列 |
| `model_name` | `varchar(100)` | 模型名；空字符串表示"该用户组的所有模型默认倍率"，联合唯一索引第二列 |
| `discount` | `decimal(10,4)` | 倍率：`1.0` 不折扣；`0.8` 即 80% 收费（vip 八折） |
| `created_at` / `updated_at` | `bigint` | unix 秒 |

`(group_name, model_name)` 联合唯一，因此同一组合不能重复插入。

## 接口

| Endpoint | Method | Auth | 说明 |
|---|---|---|---|
| `/api/group_price/` | `GET` | Admin | 全量，按 `group_name, model_name` 升序 |
| `/api/group_price/` | `POST` | Admin | 新建；`group_name` 必填；`discount=0` 自动回退 `1.0` |
| `/api/group_price/` | `PUT` | Admin | 更新 |
| `/api/group_price/:id` | `DELETE` | Admin | 删除（写后重建 `groupPriceMap`） |

## 查找顺序

`model.GetGroupDiscount(groupName, modelName)`：

1. `groupPriceMap[groupName][modelName]` → 命中返回
2. 否则 `groupPriceMap[groupName][""]` → 该用户组的全局倍率
3. 都没有 → `1.0`

`relay/billing/ratio/model.go::GetGroupDiscount(groupName, modelName, fallbackNames...)` 在 (1) 失败时按顺序尝试 `fallbackNames`（例如 `claude-3.5-sonnet-internet` → `claude-3.5-sonnet`），与 `GetModelPrice` 的 fallback 链对齐。

## 默认分组

`model_price.go::defaultGroupPrices` 初始化三种默认组：

```

启动时若表为空则写入。`Discount=1.0` 是占位，管理员按需在 UI 上调小。

## 缓存

- 启动：`InitGroupPriceCache` 把全表读入 `groupPriceMap map[string]map[string]float64`
- 写后：每个 CRUD 都调 `InitGroupPriceCache` 重建
- 周期：`SyncGroupPriceCache(frequency)` 与 `SyncModelPriceCache` 共用频率
- Redis：`CacheGetGroupPrice(group, model)` 二级缓存（`ModelPriceCacheSeconds=300`），未命中走 DB

## 与模型覆盖的关系

- 一个 group 的 `(model_name="")` 行是该用户组所有模型的默认倍率
- 单模型覆盖（`model_name="gpt-4o"`）只在 `GetGroupDiscount(group, "gpt-4o")` 时优先命中；其它模型仍走 `""` 行

例如：

| group | model | discount | 含义 |
|---|---|---|---|
| `vip` | `""` | `0.8` | vip 用户对所有模型 80% 收费 |
| `vip` | `gpt-4o` | `0.5` | vip 用户对 gpt-4o 50% 收费 |
| `svip` | `""` | `1.0` | svip 不打折（占位） |

## 前端操作指南

页面：`/setting/pricing` → "分组折扣" Tab（`web/default-pro/src/views/setting/PricingSetting.vue`）。

- 列表：group_name / model_name / discount
- 编辑弹窗：三个字段；`discount` 默认 `1.0`
- 「新增」按相同表单

## 实现位置

| 关注点 | 位置 |
|---|---|
| 数据模型 | `model/model_price.go::GroupPrice` |
| 缓存初始化 | `model/model_price.go::InitGroupPriceCache` |
| 查找（含 fallback） | `model/model_price.go::GetGroupDiscount` |
| Redis 缓存 | `model/model_price.go::CacheGetGroupPrice` |
| CRUD | `controller/model_price.go` |
| 默认分组 | `model/model_price.go::defaultGroupPrices` |
| 计费结算 | `relay/handler/helper.go::postConsumeQuota` |

