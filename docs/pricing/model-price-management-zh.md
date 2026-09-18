---
title: 模型定价
description: "模型价格 CRUD，供计费系统实时查询；为渠道编辑弹窗提供下拉数据源。"
category: pricing
order: 3
---

# 模型定价

> 给每个模型维护输入 / 输出 / 缓存 / 单次请求单价与计费类型；供 `relay` 计费链路实时查询。

入口路由：管理员设置 → **Pricing**（`/setting/pricing`）。前端组件：`web/default-pro/src/views/setting/PricingSetting.vue`。
该页签是 Root-only；普通管理员在「下拉候选」接口里也能拉取启用列表。

## 数据模型 / Data Model

`model.ModelPrice`（`model_price` 表）：

| 字段 | 类型 | 说明 |
|---|---|---|
| `model_name` | `varchar(100)` UNIQUE | 模型主键 |
| `input_price` | `decimal(16,6)` | 输入 token 单价（USD / 1K token） |
| `output_price` | `decimal(16,6)` | 输出 token 单价 |
| `cached_price` | `decimal(16,6)` | 缓存 token 单价 |
| `per_request_price` | `decimal(16,6)` | 单次请求固定价（如 dall-e / whisper） |
| `billing_type` | `varchar(20)` | `token` 或 `per_request` |
| `enabled` | `bool` | 是否启用；`false` 不参与计费、也不出现在下拉候选中 |
| `created_at` / `updated_at` | `bigint` | unix 秒 |

启动时会由 `InitDefaultPrices()` 写入一份主流模型默认表（gpt-4o、claude-3.5、deepseek、qwen-plus 等）；用户可自由增删覆盖。

## 接口一览 / Endpoints

| Endpoint | Method | 鉴权 | 说明 |
|---|---|---|---|
| `/api/model_price/` | `GET` | Root | 全量返回 `model_price` 表所有行（含 `enabled=false`） |
| `/api/model_price/options` | `GET` | Admin | **仅返回 `enabled=true` 的 `model_name` 列表**，供渠道编辑弹窗用 |
| `/api/model_price/` | `POST` | Root | 新建一行；`Insert()` 会把 `id` 强制置零以防脏数据撞主键 |
| `/api/model_price/` | `PUT` | Root | 更新指定 `id` 的字段（`input_price/output_price/cached_price/per_request_price/billing_type/enabled`） |
| `/api/model_price/:id` | `DELETE` | Root | 删除一行 |

实现：`controller/model_price.go`；前端调用走 `@/api`（GET/POST/PUT/DELETE）。

## 缓存行为 / Caching

每次 Add / Update / Delete 后都会同步调用 `model.InitModelPriceCache()`，把内存 cache `modelPriceMap` 整表重建。
后台还有 `SyncModelPriceCache(frequency)` 定时协程，按配置频率（默认 300 秒）从 DB 拉取同步，保证多实例部署的最终一致性。
计费热路径走 `model.CacheGetModelPrice(modelName)`（Redis 优先，无 Redis 时回 DB），键 `model_price:<name>`，TTL `ModelPriceCacheSeconds = 300`。

## 下拉候选 / Channel Dropdown Source

`GET /api/model_price/options` 是路由 `/api/channel` 表单里"可用模型"下拉的真实来源（admin 即可访问，避免把计费细节泄漏给普通渠道编辑流程）。
返回值是 `string[]`，按 `model_name asc` 排序；空字符串模型会被去重忽略。

## 计费模式 / Billing Modes

- `billing_type = "token"`（默认）：按 `(prompt_tokens + completion_tokens)` 计费；命中缓存的 token 按 `cached_price` 单价。
- `billing_type = "per_request"`：按调用次数计费，忽略 token（适用 dall-e / whisper / tts 等）。

`FindModelPriceByPattern` 提供子串兜底匹配：未命中精确键时按 `strings.Contains(modelName, pattern)` 兜底（适用于按系列统一价）。

## 前端操作指南 / Frontend Guide

- 页签顶部切换「模型定价」/「分组定价」。
- 列表列：模型名 / 输入价 / 输出价 / 缓存价 / 单次请求价 / 计费类型（token / per_request tag）/ 操作。
- 点击「新增」打开 640px 弹窗：模型名（必填）、输入价/输出价/缓存价/单次请求价（`precision=6`）、计费类型下拉。
- 编辑时模型名允许修改但**主键重复**会被 GORM 唯一索引拦截；保存后立即刷新整张表（包含缓存重建）。
- 删除走二次确认弹窗。

## 接口实现 / Implementation Pointers

| 关注点 | 位置 |
|---|---|
| CRUD handler | `controller/model_price.go` |
| 下拉候选 handler | `controller/model_price.go::ListModelPriceOptions` |
| 默认价格表 | `model/model_price.go::defaultModelPrices` |
| 缓存初始化 / 同步 | `model/model_price.go::InitModelPriceCache` / `SyncModelPriceCache` |
| 计费热路径 | `model/model_price.go::CacheGetModelPrice` |
| 路由 | `router/api.go` |