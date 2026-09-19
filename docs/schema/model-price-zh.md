---
title: 模型定价
description: 每个模型按 token 或按调用次数维护的单价，是计费的唯一依据。
category: schema
order: 5
---

# 模型定价（ModelPrice）

## 这是什么

`ModelPrice` 定义**一个模型在 One API Pro 内部计费的单位价格**。每次调用命中后，扣费金额完全由「调用消耗 × 该模型的 ModelPrice」决定。

你可以把它理解为「系统的价目表」。

## 在哪里看到

- **后台 → 模型定价**：列表、新增、批量调整
- **后台 → 运营仪表盘**：当日 / 当月按模型的收入预估（基于此表）

## 关键字段（运维视角）

| 字段 | 含义 | 设置后影响 |
|---|---|---|
| 模型名 | `gpt-4o` / `claude-sonnet-4` 等 | 必须与 [渠道](/schema/channel) 模型列表中的值匹配 |
| 计费类型 | `token`（按 token 计费）或 `per_request`（按次计费） | 决定后续 3 个字段的语义 |
| 输入单价 | ¥/百万 token | 影响每次调用的输入扣费 |
| 输出单价 | ¥/百万 token | 影响每次调用的输出扣费 |
| 缓存单价 | ¥/百万 token | 影响命中缓存的 token 段（0 表示不缓存） |
| 按次价格 | ¥/次 | 当计费类型为 `per_request` 时使用 |
| 启用 | 开关 | 禁用后，模型不可被调用 |

## 与「渠道」的关系

- 渠道里勾选「可用模型」时，可选项来自 ModelPrice 中已启用的模型。
- 修改 ModelPrice 不影响渠道的可用性，只影响计费。

## 与「分组折扣」的关系

最终扣费金额公式：

```
最终扣费 = 调用消耗 × ModelPrice × GroupPrice(group, model)
```

`GroupPrice` 默认 1.0（无折扣）。当用户在某个用户组，且该用户组对当前模型有折扣时生效。

## 为什么要配置

不配置的话，调用也走通，但内部额度不会扣减，相当于「白嫖」。所以系统上线**必须**先把常用模型的价格维护进来。

## 相关页面

- [模型定价（业务概念）](/pricing/model-price)
- [模型定价管理（后台）](/pricing/model-price-management)
- [分组折扣](/schema/group-price)

## 相关 API

- `GET /api/model_price/` — 列表
- `POST /api/model_price/` — 新增 / 更新
- `GET /api/model_price/options` — 渠道编辑时使用的下拉选项