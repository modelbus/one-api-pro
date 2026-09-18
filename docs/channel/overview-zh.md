---
title: 渠道概览
description: "渠道概念、生命周期与支持的 Provider 范围。"
category: channel
order: 1
---

# 渠道概览

> 渠道（Channel）是对一个上游 LLM Provider 接入实例的封装，保存鉴权、地址、模型与路由策略，并参与请求时的过滤与冷却决策。

## 概念

渠道是 One API Pro 中"上游凭证 + 路由策略"的最小单位。一条渠道对应一个 Provider 的一种接入方式（如 DeepSeek 官方、Azure OpenAI 自部署、OpenAI 兼容的中转站）。每条渠道独立保存：

- 凭证与基础 URL（`key` / `base_url` / `config`）
- 模型白名单（`models`，逗号分隔）
- 用户组白名单（`group`，逗号分隔）
- 模型名映射（`model_mapping`，JSON 对象）
- 路由参数：权重（`weight`）、优先级（`priority`）、并发上限（`max_concurrency`）、RPM 上限（`rpm`）、冷却秒数（`cooldown_seconds`）
- 业务开关：是否纯 fallback（`is_fallback` / `fallback_priority`）、系统提示词（`system_prompt`）
- 运行时指标：`status`（启用/手动禁用/自动禁用）、`balance`（美元）、`response_time`、`last_error`

数据模型定义见 `model/channel.go::Channel`；插入/更新时会同步生成 `abilities` 表行（每条 (channel_id, model) 一行）供路由层快速过滤。

## 状态机

| 值 | 常量 | 含义 |
|---|---|---|
| 0 | `ChannelStatusUnknown` | 默认值；新渠道不会出现此值（默认 1） |
| 1 | `ChannelStatusEnabled` | 启用；路由层可选 |
| 2 | `ChannelStatusManuallyDisabled` | 管理员手动禁用 |
| 3 | `ChannelStatusAutoDisabled` | 系统因余额不足 / 错误率过高 / 超时自动禁用 |

自动禁用后必须由管理员手动或在批量测试通过后回到 `Enabled`。

## 支持的 Provider

完整 Provider 列表见 [Provider 一览](./provider-list)。所有走 OpenAI 兼容协议的中转站都可以复用 `openai` 类型：填好 `base_url` 与 `models`，再通过 `model_mapping` 校正模型名即可。

## 相关文档

- [新增渠道](./add-channel) — 字段含义与新增流程
- [渠道路由](./channel-routing) — 过滤、冷却、并发、RPM、Sticky、fallback
- [渠道测试](./channel-test) — 单渠道与全渠道测试
- [余额刷新](./balance-update) — 自动与手动拉取上游余额

