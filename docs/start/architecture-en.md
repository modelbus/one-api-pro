---
title: Architecture
description: "System components, data flow and extension points."
category: start
order: 4
---

# Architecture

> System components, data flow and extension points.
> 系统组件、数据流与扩展点。

## Layers / 分层

| Directory / 目录 | Responsibility / 职责 | May depend on / 可依赖 |
| --- | --- | --- |
| `main.go` | Process entrypoint, config loading, embed-theme bootstrap | `common/`, `model/` |
| `cmd/` | One-shot CLI tools (migrations, maintenance) | `common/`, `model/` |
| `common/` | Config, logging, Redis, HTTP, crypto, i18n, render, verification, payment | stdlib + third-party |
| `common/payment/` | Payment channel implementations (`wechat.go` / `alipay.go` / `bank.go`), self-register via `RegisterChannel` | `model.SystemSetting` |
| `common/embed_theme.go` | `//go:embed web/build/default-pro` bundles the built frontend into the binary | none |
| `model/` | GORM models + business functions (NO HTTP / Gin dependency) | `common/`, GORM |
| `controller/` | Gin handlers — parameter parsing, model calls, response assembly | `model/`, `middleware/` |
| `middleware/` | auth, rate-limit, Turnstile, gzip, language, cache | `model/`, `common/` |
| `relay/` | LLM relay layer (provider sub-packages) | `model/`, `relay/registry` |
| `relay/adaptor/` | Providers self-register in `init()` and implement the unified `Adaptor` interface | `relay/registry` |
| `relay/billing/` | Pre-consume / refund quota | `model/` |
| `relay/registry/` | Central provider metadata index (dual-mapped by `ID` and legacy `type`) | `relay/adaptor` |
| `router/` | Route registration | `controller/`, `middleware/` |
| `monitor/` | Patrols (success rate / channel auto-disable / cluster health) | `model/` |
| `web/default-pro/` | Frontend project (Vue 3 + Arco + Vite); `npm run build` output is `//go:embed`-ed | none |

> **Frontend is maintained in `web/default-pro/` ONLY.** The other `web/*` directories are historical themes or build output — do not edit (see [Code Style](/en/contribute/code-style)).
> **前端仅在 `web/default-pro/` 维护**。`web/air/ web/default/ web/berry/ web/THEMES web/build/` 是历史主题或构建产物，禁改。

## Request Lifecycle / 请求生命周期

For `POST /v1/chat/completions`:

以「用户调用 `/v1/chat/completions`」为例：

```text
Client ──► Gin router
        ├─► middleware: recover · request-id · logger · language · cors · gzip
        ├─► middleware: rate-limit (GLOBAL_API_RATE_LIMIT)
        ├─► middleware: TokenAuth       // parse sk-xxx
        │       ├─► CacheGetTokenByKey  // Redis hit / DB miss fallback
        │       ├─► token status check (enabled / disabled / expired / exhausted)
        │       ├─► subnet / models restrictions
        │       └─► user quota cache
        │
        ├─► relay/handler: dispatch to Chat / Embeddings / Image / Audio
        │       ├─► meta.Meta assembly (model, user, token, channel, isStream)
        │       ├─► Channel routing (middleware/distributor): weight / concurrency / cooldown / sticky
        │       ├─► preConsumeQuota: user_quota cache → PreConsumeTokenQuota
        │       ├─► adaptor.DoRequest → upstream provider
        │       ├─► adaptor.DoResponse: stream / JSON conversion, extract Usage
        │       └─► billing.PostConsumeQuota: settle the delta + refund (roll back on failure)
        │
        └─► middleware: ginzap / logger writes to the logs table
```

Highlights / 要点：

- **Three auth schemes**: Cookie Session (`/api/*`), user Access Token (fallback for `/api/*`), API Key (`/v1/*`). See [API Overview · Auth](/en/api/README#appendix-a-auth-mechanisms).
- **Channel routing**: weighted round-robin by default; queues when concurrency is full; skips cooldown channels; sticky to the same channel for the same user across requests.
- **Two-phase billing**: pre-consume `PreConsumedQuota + promptTokens`, then settle against the real `usage` (`PostConsumeQuota`); failures / exceptions roll back via `ReturnPreConsumedQuota`.
- **Cache consistency**: top-up paths (`IncreaseUserQuota` etc.) refresh the Redis `user_quota:<id>` key synchronously, avoiding the "DB has the truth, cache keeps returning stale values" → repeated 403 scenario.

## Data Model / 数据模型

Core tables (auto-created via `gorm.AutoMigrate`):

核心表（`model/*.go` 中 `AutoMigrate` 自动建表）：

| Table / 表 | Key fields / 关键字段 | Notes / 说明 |
| --- | --- | --- |
| `users` | `role` (0/10/100), `status` (1/2/3), `quota`, `group`, `access_token` | User accounts & roles |
| `tokens` | `key`, `status`, `remain_quota`, `expired_time`, `models`, `subnet` | API key table |
| `channels` | `type`, `key`, `base_url`, `weight`, `cooldown`, `models`, `group` | Provider channel configs |
| `abilities` | `(group, model_id, channel_id)` unique | Which group / model / channel combinations are allowed |
| `redemptions` | `code`, `quota`, `used_user_id`, `expired_time` | Redemption codes |
| `plans` | `name`, `quota`, `price`, `duration_days`, `status` | Subscription plans |
| `user_plans` | `(user_id, plan_id)`, `start_at`, `expire_at`, `remain_quota` | User subscriptions |
| `plan_usages` | Daily aggregated plan usage | Reporting |
| `orders` | `order_no`, `type` (1 plan / 2 topup), `status` (0/1/2/3), `pay_method`, `pay_trade_no`, `source` (1/2) | Orders |
| `logs` | `type` (consume/topup/manage/system/test), `model_name`, `prompt/completion_tokens`, `quota`, `channel_id` | Call logs |
| `options` / `system_settings` | `(key, value, category)` | Generic KV config |
| `cluster_nodes` | `node_id`, `address`, `status`, `last_heartbeat`, `ping_failures` | Cluster nodes |

## Provider Self-Registration / Provider 自注册

Providers register themselves at runtime via `init()`:

Provider 是运行时通过 `init()` 自注册到 `relay/registry` 的：

```go
// relay/adaptor/provider/deepseek/register.go
package deepseek

import (
    "github.com/modelbus/one-api-pro/relay/adaptor"
    "github.com/modelbus/one-api-pro/relay/registry"
)

func init() {
    registry.Register(registry.ChannelMeta{
        ID:             "deepseek",
        Name:           "DeepSeek",
        DefaultBaseURL: "https://api.deepseek.com",
        LegacyType:     36,
    }, func() adaptor.Adaptor { return &Adaptor{} })
}
```

For the full walkthrough see [Add a Provider](/en/contribute/add-provider).

新增 Provider 的步骤见 [新增 Provider](/zh/contribute/add-provider)。

## Payment Channels / 支付通道

Payment channels implement `common/payment.Channel` and self-register — adding a new one is zero-intrusion:

支付通道采用 `common/payment.Channel` 接口 + `RegisterChannel` 自注册，新增通道零侵入：

```go
// common/payment/wechat.go
func init() { RegisterChannel(&wechatChannel{}) }
```

The order flow funnels through `controller/buildPayInfo(payMethod, orderNo, amount, subject)`:

下单链路统一走 `controller/buildPayInfo(payMethod, orderNo, amount, subject)`：

1. Select the channel by `pay_method` via `payment.New(payMethod)`.
2. Read `system_settings` and check `IsEnabled()`.
3. Call `channel.PrePay(orderNo, amount, subject)` to obtain `pay_url / qr_code / trade_no`.
4. Async notify hits `controller/payment.go::processNotify` → `VerifyNotify` → `MarkOrderPaid` → dispatch by `order.Type` (activate plan or credit top-up).

For the full walkthrough see [Add a Payment Channel](/en/contribute/add-payment).

新增通道的步骤见 [新增支付通道](/zh/contribute/add-payment)。

## Cluster Mode / 集群模式

When `CLUSTER_ENABLED=true`, One API Pro runs in decentralized multi-active mode:

`CLUSTER_ENABLED=true` 时启动「去中心化多活」模式：

- Each node owns its own MySQL + Redis.
- GORM callbacks (`AfterCreate / AfterUpdate / AfterDelete`) capture business-table changes.
- A background goroutine pushes events to all alive nodes via HTTP POST.
- The receiver compares `updated_at` and only writes the newer version (last-writer-wins).

See [Cluster Overview](/en/decentralization/overview), [Config Sync](/en/decentralization/config-sync), [Node Health](/en/decentralization/node-health).

详见 [Cluster 概览](/zh/decentralization/overview)、[配置同步](/zh/decentralization/config-sync)、[节点健康](/zh/decentralization/node-health)。

## Extension Points / 扩展点

| What you want / 想做的事 | Where / 改哪里 |
| --- | --- |
| Add a new provider | `relay/adaptor/provider/<name>/` + implement the `Adaptor` interface |
| Add a new payment channel | `common/payment/<name>.go` + implement `Channel` and `init()` |
| Add a new admin menu | `web/default-pro/src/router/index.js` + `web/default-pro/src/views/` |
| Tweak the unified response | every `gin.H{"success":...}` in `controller/*.go` |
| Tweak the billing formula | `relay/billing/` + `model.PreConsumeTokenQuota / PostConsumeQuota` |
| Add a cache key | `model/cache.go` with the unified `model.cacheSeconds` TTL |

Next: [Dev Setup](/en/contribute/dev-setup) · [Code Style](/en/contribute/code-style).