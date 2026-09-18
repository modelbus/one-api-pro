---
title: 分组定价
description: "按 用户分组 × 模型 维护折扣倍率；缺省折扣为 1.0。"
category: admin
order: 4
---

# 分组定价

> 给「用户分组 × 模型」组合设置折扣倍率（`discount`），与 `model_price` 相乘得到最终用户实付单价。

入口路由：管理员设置 → **Pricing**（`/setting/pricing`）的「分组定价」页签。前端组件：`web/default-pro/src/views/setting/PricingSetting.vue`。
分组定价是 Root-only。

## 数据模型 / Data Model

`model.GroupPrice`（`group_price` 表）：

| 字段 | 类型 | 说明 |
|---|---|---|
| `group_name` | `varchar(32)`，与 `model_name` 组成复合唯一索引 | 用户分组名（如 `default` / `vip` / `svip`） |
| `model_name` | `varchar(100)`，复合唯一索引 | 模型名；空字符串 `""` 表示「该分组的所有模型统一折扣」 |
| `discount` | `decimal(10,4)` | 倍率（`1.0` = 不打折）；小于 `1.0` 折扣，大于 `1.0` 加价 |

启动时由 `InitDefaultPrices()` 写入三条默认分组（`default` / `vip` / `svip`，折扣均为 `1.0`），首次部署即用。

## 接口一览 / Endpoints

| Endpoint | Method | 鉴权 | 说明 |
|---|---|---|---|
| `/api/group_price/` | `GET` | Root | 全量返回 `group_price` 行 |
| `/api/group_price/` | `POST` | Root | 新建一行；缺省 `discount=1.0`；`Insert()` 强制清零 `id` |
| `/api/group_price/` | `PUT` | Root | 更新 `discount` 字段 |
| `/api/group_price/:id` | `DELETE` | Root | 删除一行 |

实现：`controller/model_price.go`（与模型定价共用一个文件，handler 共享鉴权风格）。

## 查询逻辑 / Lookup

`model.GetGroupDiscount(groupName, modelName)` 的兜底顺序：

1. 精确匹配 `(groupName, modelName)`：返回该条 `discount`。
2. 兜底匹配 `(groupName, "")`：返回「该分组所有模型统一折扣」。
3. 都没有：返回 `1.0`（不打折）。

`GetGroupNames()` 返回当前缓存中所有出现过的 `group_name`，供用户管理/编辑下拉使用。

## 缓存行为 / Caching

`model.GroupPriceCacheSeconds` 默认 300s；每次 CRUD 都会调 `model.InitGroupPriceCache()` 全量重建内存 `groupPriceMap`。
计费热路径走 `model.CacheGetGroupPrice(groupName, modelName)`（Redis 优先，键 `group_price:<group>:<model>`，回 DB 兜底）。

## 实际计费 / How It Combines

最终单价 = `model_price.input_price` × `group_price.discount`（按模型名 × 分组名查）。
对 `per_request` 计费的模型，`discount` 同样作用于 `per_request_price`。

## 前端操作指南 / Frontend Guide

- 「分组定价」页签：表列 = 分组名 / 模型名 / 折扣 / 操作。
- 「新增」弹窗（520px）：
  - `group_name`（必填，下拉候选由 `GET /api/group/` 提供）
  - `model_name`（必填，下拉候选来自 `/api/model_price/options`）
  - `discount`（`precision=4`，最小 `0`，留空使用 `1.0`）
- 行内操作：编辑、删除（二次确认）。
- 注意：修改一个组的折扣，会即时影响所有用该 group 调用对应模型的用户实付额；后台有 `SyncGroupPriceCache` 定时兜底。

## 接口实现 / Implementation Pointers

| 关注点 | 位置 |
|---|---| 
| CRUD handler | `controller/model_price.go` (`AddGroupPrice` / `UpdateGroupPrice` / `DeleteGroupPrice`) |
| 默认分组表 | `model/model_price.go::defaultGroupPrices` |
| 缓存初始化 / 同步 | `model/model_price.go::InitGroupPriceCache` / `SyncGroupPriceCache` |
| 计费热路径 | `model/model_price.go::CacheGetGroupPrice` |
| 分组下拉数据 | `controller/group.go::GetGroups` → `model.GetGroupNames` |
| 路由 | `router/api.go` |