---
title: Overview
description: "Positioning, core features and use cases of One API Pro."
category: start
order: 1
---

# Overview

> Positioning, core features and use cases of One API Pro.
> One API Pro 的定位、核心特性与适用场景。

## One-liner / 一句话

One API Pro unifies 40+ LLM providers behind a single **OpenAI-compatible** HTTP gateway, with channel routing, subscription & metered billing, WeChat / Alipay / bank top-ups, redemption codes, decentralized multi-node clustering and an operations dashboard — all production-ready out of the box.

把 40+ LLM 提供商统一为一个 **OpenAI 兼容** 的 HTTP 网关，内置渠道路由、订阅/按量计费、微信/支付宝/银行收款、兑换码、去中心化多节点集群与运营仪表盘，开箱即用。

## Core Features / 核心特性

| Module / 模块 | Capability / 能力 |
| --- | --- |
| **Multi-provider aggregation** | OpenAI, Anthropic, Gemini, Qwen, DeepSeek, Doubao, Mistral, Cohere, xAI, VertexAI, AWS Claude and more — all behind the OpenAI `/v1/chat/completions` protocol |
| **Smart channel routing** | Weighted, concurrency-aware, cooldown, sticky sessions, automatic failover, auto-disable on low success rate |
| **Subscription + pay-as-you-go** | Token Plan (recurring), metered billing, group discount ratio, plan upgrade differential orders (`UP` prefix) |
| **Top-up & payments** | WeChat Native, Alipay, bank transfer; custom amounts, exchange rate, preset chips (`TP` prefix orders) |
| **Redemption codes** | Bulk generation, export, revoke; for activation, promotions and gifts |
| **Multi-node cluster** | Per-node DB + Redis, GORM callbacks + HTTP push-based sync (see [Cluster Overview](/en/decentralization/overview)) |
| **Operations dashboard** | KPIs, trend charts, model Top-N distribution, user leaderboard, order analytics (`/api/admin/dashboard/*`) |
| **Three auth schemes** | Cookie Session, user-level Access Token, `sk-`-style API Key — each serving a different surface |
| **i18n** | Full zh / en internationalization (admin, auth, landing, legal) with Arco components localized |

## Use Cases / 适用场景

- **Solo / small teams**: self-host an OpenAI-compatible proxy, consolidate billing across providers.
  个人 / 小团队：自建 OpenAI 兼容代理，统一管理多家供应商账单。
- **AI app vendors**: per-customer channels with subscription / token quotas for internal billing.
  AI 应用开发商：为多个客户/项目分别配置渠道，按订阅/Token 配额做内部计费。
- **Enterprise gateway**: aggregate vendors with RBAC, audit, rate-limit.
  企业内部平台：聚合多家大模型做内部网关，配合 RBAC、审计、限流。
- **Commercial resale**: WeChat / Alipay enabled, mixed Token Plan + top-up billing for resale.
  商业化分发平台：开启微信/支付宝收款，做 Token Plan + 充值混合计费，对外销售额度。

## Relationship with upstream One-API / 与上游 One-API 的关系

One API Pro is built on top of the upstream One-API and extends it:

- **Kept**: model protocol compatibility, core auth & billing flow, single-binary deploy.
  保留：模型协议兼容、核心鉴权与计费流程、单二进制部署。
- **Replaced**: routing, middleware, caching, monitoring, i18n, payment, clustering.
  替换：路由 / 中间件 / 缓存策略 / 监控 / i18n / 支付 / 集群。
- **Added**: plan / subscription system, top-up, redemption, ops dashboard, decentralized cluster, `v1/chat/completions` streaming + non-stream dual-protocol support.

Next: [Showcase](/en/start/showcase) · [Quick Tour](/en/start/quick-tour) · [Architecture](/en/start/architecture).