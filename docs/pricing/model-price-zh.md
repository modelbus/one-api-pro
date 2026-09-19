---
title: 模型定价
description: 每个模型的单价，是系统计费的唯一依据。
category: pricing
order: 1
---

# 模型定价

> 用户调用模型时扣多少钱，看这里。

## 这是什么

`ModelPrice` 是「系统的价目表」。每个模型在表里有一条记录，决定：

- 按 token 计费：输入 / 输出各 ¥ 多少 / 百万 token
- 按次计费：调用一次 ¥ 多少（用于 DALL·E / 图像生成等）

每次调用扣费金额 = 调用消耗 × 该模型的 ModelPrice × 用户组折扣。

## 在哪里看到

- **后台 → 模型定价**：列表、新增、批量调整
- **后台 → 运营仪表盘**：当日 / 当月按模型的收入预估（基于此表）

## 需要设置什么

后台 → 模型定价 → 新增，每条填：

| 项 | 含义 | 设置后影响 |
|---|---|---|
| `model_name` | 模型名 | 必须与 [渠道](/channel/add-channel) 模型列表中的值匹配；错则该模型无法被调用 |
| `billing_type` | `token`（按 token）或 `per_request`（按次） | 决定下面 4 个字段的语义 |
| `input_price` | 输入单价（¥/百万 token） | 影响每次调用的输入扣费 |
| `output_price` | 输出单价（¥/百万 token） | 影响每次调用的输出扣费 |
| `cached_price` | 缓存命中部分单价（¥/百万 token） | 通常是 `input_price` 的一半；0 表示不缓存 |
| `per_request_price` | 按次价格（¥/次） | 仅 `billing_type=per_request` 时使用 |
| `enabled` | 开关 | 禁用后模型不可被调用 |

> 系统初始化时会种入一些常用模型的默认价格（GPT-4o / Claude / DeepSeek 等）。可在后台调整。

## 与「渠道」的关系

- 渠道里勾选「可用模型」时，可选项来自 `ModelPrice` 中 `enabled=true` 的模型
- 修改 `ModelPrice` 不影响渠道可用性，只影响计费

## 什么时候需要改

- 系统上线后**必须**先把常用模型价格维护进来，否则调用不扣费（白嫖）
- 上游 Provider 调价时
- 新增模型时

## 常见问题

- **新增模型但渠道选不到**：先在「模型定价」里把该模型设为 `enabled=true`
- **调用成功但余额没扣**：检查模型是否在 ModelPrice 里；价格是不是 0
- **改价后老用户的扣费怎么算**：调用按当时模型的 ModelPrice 算；改价只影响之后的调用

## 相关页面

- [模型定价管理（后台）](./model-price-management)
- [模型定价（数据结构）](/schema/model-price)
- [分组折扣](/schema/group-price)

## 相关 API

- `GET /api/model_price/` — 列表（Root）
- `GET /api/model_price/options` — 仅 `enabled=true` 的模型名（Admin，用于渠道编辑下拉）
- `POST /api/model_price/` — 新增
- `PUT /api/model_price/` — 更新
- `DELETE /api/model_price/:id` — 删除