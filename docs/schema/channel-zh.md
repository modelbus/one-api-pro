---
title: 渠道
description: 接入一个上游大模型 Provider 的最小单位：凭证、路由与统计。
category: schema
order: 4
---

# 渠道（Channel）

## 这是什么

`Channel` 是接入一个上游 Provider（如 OpenAI、Anthropic、Azure）的最小单位。每新建一个渠道，One API Pro 就多一个可选的上游来源；用户发起调用时，系统按 [路由策略](/channel/channel-routing) 从多个渠道里挑一个转发。

## 在哪里看到

- **后台 → 渠道**：列表、新增、连通性测试、启用/禁用、刷新余额
- **后台 → 运营仪表盘 → Top Users**：每个用户的调用最终打到哪个渠道
- **调用日志**：每条调用记录命中的渠道 ID

## 你需要为每个渠道设置什么

新增渠道时，**通常只需要填这些字段**：

| 项 | 含义 | 设置后影响 |
|---|---|---|
| `type` | Provider 类型 | 决定请求协议（OpenAI 兼容 / Anthropic / Azure 等） |
| `name` | 渠道名称 | 仅用于辨识多个渠道 |
| `base_url` | 上游 Provider 的接入地址 | 错就 404 / 拒连；OpenAI 通常 `https://api.openai.com/v1`；Azure 按部署填 |
| `key` | API Key / 凭证 | 错就 401；不要带引号、空格、换行 |
| `models` | 模型列表（多选） | 从 [模型定价](/schema/model-price) 中已有模型里勾选；不勾的模型不会被路由到该渠道 |
| `weight` / `priority` | 路由权重 / 优先级 | 决定被路由命中的概率 |
| `status` | 启用 / 禁用 | 禁用后路由跳过 |

下面这些是**大多数情况不需要动**的字段，除非遇到问题：

- `system_prompt`：自定义请求头 / 响应转换 JSONPath
- `max_concurrency` / `cooldown_seconds`：并发 / 冷却
- `group`：用户组白名单

## 它如何参与计费

调用命中渠道后，扣费金额由 [模型定价](/schema/model-price) 决定。渠道本身不存价格，只存凭证与路由。

## 余额自动更新

对支持「账户余额查询接口」的 Provider（OpenAI / Azure 等），One API Pro 可以定时拉取上游余额展示在渠道列表里，便于发现「渠道快要没钱」。开启方式：渠道编辑页底部「自动更新余额」开关。

## 相关页面

- [渠道概览](/channel/overview)
- [新增渠道](/channel/add-channel)
- [渠道连通性测试](/channel/channel-test)
- [渠道路由策略](/channel/channel-routing)
- [余额自动更新](/channel/balance-update)
- [Provider 全清单](/channel/provider-list)

## 相关 API

- `GET /api/channel/` — 渠道列表
- `POST /api/channel/` — 新建渠道
- `PUT /api/channel/` — 更新渠道
- `GET /api/channel/test/:id` — 测试连通性