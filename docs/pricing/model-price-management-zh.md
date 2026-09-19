---
title: 模型定价管理
description: 管理员如何在后台维护模型单价。
category: pricing
order: 3
---

# 模型定价管理

> 后台 → 模型定价。管理员在这里维护每个模型的单价。

## 在哪里

后台 → 模型定价。

## 列表看到什么

每行展示：

- 模型名
- 输入价 / 输出价 / 缓存价 / 单次请求价
- 计费类型（token / per_request）
- 状态（启用 / 禁用）
- 操作按钮

## 怎么新增

1. 列表右上角「新增」
2. 填：
   - 模型名（必填，唯一）
   - 输入价 / 输出价 / 缓存价 / 单次请求价
   - 计费类型（`token` 默认，或 `per_request`）
   - 启用（默认开）
3. 保存

## 怎么编辑

行内「编辑」打开弹窗，修改任意字段保存即可。改后立即生效。

## 怎么禁用

行内「禁用」：

- 禁用后该模型不可被调用（不会出现在「渠道编辑」下拉里）
- 不影响已开始的调用

通常用于：上游 Provider 下架某个模型时。

## 怎么删除

行内「删除」（二次确认）。删除后该模型彻底不可用。

通常不需要删除：默认状态改成「禁用」即可。

## 注意事项

- **价格单位是元 / 百万 token**。例：`0.002` 表示每百万 token 收 0.002 元
- **缓存价通常 = 输入价的一半**。0 表示不区分缓存
- **新增模型后** 还要在 [渠道](/channel/add-channel) 里的「可用模型」勾选上，才能被路由到

## 系统自带默认价格

系统启动时会种入一些主流模型的默认价格（GPT-4o / Claude / DeepSeek 等）。可在后台调整覆盖。

## 常见问题

- **价格改完后调用没按新价扣**：检查缓存重建是否完成；通常改后立即生效
- **调用返回「model price not set」**：模型不在 [ModelPrice](/pricing/model-price) 表里，或 `enabled=false`
- **想给某 Provider 全系列统一价**：可在渠道编辑时用 `model_mapping`，但每条 ModelPrice 仍要逐个加

## 相关页面

- [模型定价（业务概念）](/pricing/model-price)
- [模型定价（数据结构）](/schema/model-price)
- [渠道管理](/channel/add-channel)

## 相关 API

- `GET /api/model_price/` — 列表（Root）
- `POST /api/model_price/` — 新增（Root）
- `PUT /api/model_price/` — 更新（Root）
- `DELETE /api/model_price/:id` — 删除（Root）