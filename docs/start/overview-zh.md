---
title: 项目简介
description: "One API Pro 的定位、核心特性与适用场景。"
category: start
order: 1
---

# 项目简介

> One API Pro 的定位、核心特性与适用场景。
> Positioning, core features and use cases of One API Pro.

## 一句话 / One-liner

把 40+ LLM 提供商统一为一个 **OpenAI 兼容** 的 HTTP 网关，内置渠道路由、订阅/按量计费、微信/支付宝/银行收款、兑换码、去中心化多节点集群与运营仪表盘，开箱即用。

One API Pro unifies 40+ LLM providers behind a single **OpenAI-compatible** HTTP gateway, with channel routing, subscription & metered billing, WeChat / Alipay / bank top-ups, redemption codes, decentralized multi-node clustering and an operations dashboard — all production-ready out of the box.

## 核心特性 / Core Features

| 模块 / Module | 能力 / Capability |
| --- | --- |
| **多模型聚合** | OpenAI、Anthropic、Gemini、Qwen、DeepSeek、Doubao、Mistral、Cohere、xAI、VertexAI、AWS Claude 等；全部统一为 OpenAI `/v1/chat/completions` 协议 |
| **智能渠道路由** | 权重、并发、冷却时间、自动粘性会话（Sticky）、失败转移、自动下线低成功率渠道 |
| **订阅 + 按量** | Token Plan（套餐按期）、按量计费、分组折扣倍率、套餐升降级差价订单（`UP` 前缀） |
| **充值与支付** | 微信 Native、支付宝、银行转账，支持自定义金额、汇率换算、预设金额 chip（`TP` 前缀订单） |
| **兑换码** | 批量生成、导出、回收；用于拉新、运营活动、赠送 |
| **多节点集群** | 每节点独立 DB + Redis，通过 GORM 回调 + HTTP 主动推送同步（详见 [Cluster 概览](/zh/decentralization/overview)） |
| **运营仪表盘** | KPI、趋势图、模型分布 Top-N、用户排行、订单统计（`/api/admin/dashboard/*`） |
| **鉴权三件套** | Cookie Session、用户级 Access Token、`sk-` 形式的 API Key，分别服务不同场景 |
| **国际化** | 全量中英文 i18n（管理后台、登录页、落地页、条款隐私），Arco 组件本地化 |

## 适用场景 / Use Cases

- **个人 / 小团队**：自建 OpenAI 兼容代理，统一管理多家供应商账单。
  Solo / small teams: self-host an OpenAI-compatible proxy, consolidate billing across providers.
- **AI 应用开发商**：为多个客户/项目分别配置渠道，按订阅/Token 配额做内部计费。
  AI app vendors: per-customer channels with subscription / token quotas for internal billing.
- **企业内部平台**：聚合多家大模型做内部网关，配合 RBAC、审计、限流。
  Enterprise gateway: aggregate vendors with RBAC, audit, rate-limit.
- **商业化分发平台**：开启微信/支付宝收款，做 Token Plan + 充值混合计费，对外销售额度。
  Commercial resale: WeChat / Alipay enabled, mixed Token Plan + top-up billing for resale.

## 与上游 One-API 的关系 / Relationship with upstream One-API

One API Pro 在上游 One-API 之上重构与扩展：

One API Pro is built on top of the upstream One-API and extends it:

- **保留 / Kept**：模型协议兼容、核心鉴权与计费流程、单二进制部署。
  Kept: model protocol compatibility, core auth & billing flow, single-binary deploy.
- **替换 / Replaced**：路由 / 中间件 / 缓存策略 / 监控 / i18n / 支付 / 集群。
  Replaced: routing, middleware, caching, monitoring, i18n, payment, clustering.
- **新增 / Added**：套餐系统（Plan / Subscription）、充值、兑换码、运营仪表盘、去中心化集群、`v1/chat/completions` 流式 + non-stream 双协议兼容。

下一步 / Next: [展示案例](/zh/start/showcase) · [5 分钟快速体验](/zh/start/quick-tour) · [架构总览](/zh/start/architecture)。