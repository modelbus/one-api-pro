---
title: 模型定价
description: "`model_price` 表、单位、字段含义与 `ListModelPriceOptions` 下拉数据源。"
category: pricing
order: 1
---

# 模型定价

> 每个模型在 `model_price` 表里有一条独立定价记录，决定按 token / 按次计费的金额。实现：`model/model_price.go::ModelPrice`、`controller/model_price.go`。

## 数据模型 / Data model

`model.ModelPrice`（`model_price` 表）：

| 字段 | 类型 | 说明 |
|---|---|---|
| `model_name` | `varchar(100)` UNIQUE | 模型名；按它做匹配（请求 model 命中走 `GetModelPrice`） |
| `input_price` | `decimal(16,6)` | 输入单价（¥/百万 token） |
| `output_price` | `decimal(16,6)` | 输出单价（¥/百万 token） |
| `cached_price` | `decimal(16,6)` | 缓存命中部分单价（¥/百万 token），通常 = `input_price` 的 1/2 |
| `per_request_price` | `decimal(16,6)` | 按次计费金额（DALL·E / TTS / 嵌入等非 token 模型） |
| `billing_type` | `varchar(20)` | `token` / `per_request` |
| `enabled` | `bool` | 是否启用；下拉只列出 `enabled=true` 的行 |

价格单位是"元 / 百万 token"。`relay/billing/ratio` 计算时再乘上 group discount。

## 接口 / Endpoints

| Endpoint | Method | Auth | 说明 |
|---|---|---|---|
| `/api/model_price/` | `GET` | Root | 全量（含 `enabled=false`） |
| `/api/model_price/options` | `GET` | Admin | 仅返回 `enabled=true` 的 `model_name` 列表，渠道编辑下拉专用 |
| `/api/model_price/` | `POST` | Admin | 新建（`model_name` 必填；`billing_type` 缺省 `token`） |
| `/api/model_price/` | `PUT` | Admin | 更新（写后重建 `modelPriceMap`） |
| `/api/model_price/:id` | `DELETE` | Admin | 删除（写后重建缓存） |

`controller/model_price.go` 中 `GetAllModelPrices` 标记为 Deprecated（仅 Root），原因：它会泄露 `input_price` 等计费字段到非必要场景；新增前端场景请优先使用 `ListModelPriceOptions`。

## 缓存 / Cache

`InitModelPriceCache` 启动时把 `enabled=true` 行写入进程内 `modelPriceMap`（`map[string]*ModelPrice`），由读写锁保护：

- `GetModelPrice(name)` 是 RLock 下读 map
- `FindModelPriceByPattern(name)` 是兼容：用 `name` 精确 + 子串匹配；适配 `-internet` 等后缀
- 所有写操作（POST / PUT / DELETE）后调 `InitModelPriceCache` 重建
- 周期任务 `SyncModelPriceCache(frequency)` 每 N 秒重建一次（`SyncFrequency` 来自 `config.SyncFrequency`）

`CacheGetModelPrice(name)` / `CacheGetGroupPrice(group, model)` 是 Redis 二级缓存（`ModelPriceCacheSeconds=300`），先查 Redis，未命中走 DB 并回填。

## 计费结算路径 / Billing flow

`relay/billing/ratio/model.go::GetModelPrice(name, fallbackNames...)`：

1. 先用 `name` 查 `model.GetModelPrice`
2. 未命中则依次尝试 `fallbackNames`（如 `gpt-4o-internet` → `gpt-4o`）
3. 仍未命中返回 `errors.New("model price not set")`

随后 `postConsumeQuota` 按 `BillingType` 计算 quota：

```go
if BillingType == PerRequest {
  quota = per_request_price * 1 * group_discount
} else {
  quota = (input_price * prompt_tokens
         + output_price * completion_tokens
         + cached_price * cached_tokens)
         * group_discount / 1_000_000   // 元 → quota
}
```

## 默认定价 / Default prices

`model_price.go::defaultModelPrices` 列出系统初始化时种入的常用模型（`gpt-4o` / `claude-3.5-sonnet` / `deepseek-chat` / `qwen-max` / `gemini-1.5-pro` 等），仅在 `model_price` 表为空时由 `InitDefaultPrices` 写入。`OnConflict{DoNothing: true}` 保证重跑安全。

## 前端操作指南 / Frontend Guide

页面：`/setting/pricing` → "模型定价" Tab（`web/default-pro/src/views/setting/PricingSetting.vue`）。

- 列表：model_name / input_price / output_price / cached_price / per_request_price / billing_type / enabled
- 编辑弹窗：六个数值字段 + 计费类型下拉（`token` / `per_request`）+ enabled 开关
- 「新增」按相同表单

## 实现位置 / Implementation Pointers

| 关注点 | 位置 |
|---|---|
| 数据模型 | `model/model_price.go::ModelPrice` |
| 缓存初始化 | `model/model_price.go::InitModelPriceCache` |
| 查找（精确 / 模糊） | `model/model_price.go::GetModelPrice` / `FindModelPriceByPattern` |
| Redis 缓存 | `model/model_price.go::CacheGetModelPrice` |
| 下拉数据源 | `controller/model_price.go::ListModelPriceOptions` |
| CRUD | `controller/model_price.go` |
| 默认定价 | `model/model_price.go::InitDefaultPrices` |
| 计费结算 | `relay/handler/helper.go::postConsumeQuota` |