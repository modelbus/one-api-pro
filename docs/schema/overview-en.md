---
title: Schema Overview
description: Core persisted data models in One API Pro and how they fit together.
category: schema
order: 1
---

# Schema Overview

> This category explains the data One API Pro persists, answering:
>
> - What is each entity, who creates it, and where can I find it?
> - Where do I change it from?
> - What changes when I edit it?

For HTTP endpoints, see the [API Reference](/en/api/README).

## Reading order

First time here? Read in this order:

1. [User](/en/schema/user) — identity and quota hub of the whole system
2. [Channel](/en/schema/channel) — smallest unit for an upstream provider
3. [Plan](/en/schema/plan) and [Subscription](/en/schema/subscription) — monetization model
4. [Order](/en/schema/order) — what users actually paid for

Debugging an issue? Jump to:

- Quota problems → [User](/en/schema/user) → [Access Token](/en/schema/token) → [Log](/en/schema/log)
- Routing misbehaves → [Channel](/en/schema/channel) → [Model Price](/en/schema/model-price) → [Group Price](/en/schema/group-price)
- Billing/subscription policy → [Plan](/en/schema/plan) → [Subscription](/en/schema/subscription) → [Order](/en/schema/order) → [Topup](/en/schema/topup)

## Entity map

| Entity | Created by | Consumed by | Manage in |
|---|---|---|---|
| User | Self-registration / admin | All features | Admin → Users |
| Token | User | Admin API | Admin → API Tokens |
| Channel | Admin | Channel router | Admin → Channels |
| ModelPrice | Admin | Billing | Admin → Model Prices |
| GroupPrice | Admin | Billing | Admin → Group Prices |
| Plan | Admin | Checkout | Admin → Plans |
| Subscription | System (after checkout) | Billing / rate limit | Admin → Subscriptions |
| Order | User | Finance | Admin → Orders |
| Topup | User | Billing | Admin → Topups |
| Redemption | Admin | User redeem | Admin → Redemption Codes |
| ClusterNode | Node on startup | Cluster sync | Admin → Cluster |
| Log | System (per call) | Admin debugging | Admin → Logs |
| Option | Admin | Global behavior | Admin → System |