---
title: 架构总览
description: "系统组件、数据流与扩展点。"
category: start
order: 4
---

# 架构总览

> 系统组件、数据流与扩展点。
> System components, data flow and extension points.

## 分层 / Layers

| 目录 / Directory | 职责 / Responsibility | 可依赖 / May depend on |
| --- | --- | --- |
| `main.go` | 进程入口、配置加载、Embed 主题装载 | `common/`, `model/` |
| `cmd/` | 一次性 CLI 工具（迁移脚本、维护任务） | `common/`, `model/` |
| `common/` | 配置、日志、Redis、HTTP、加密、I18N、Render、验证码、支付通道 | 仅标准库 / 第三方 |
| `common/payment/` | 各支付通道实现（`wechat.go` / `alipay.go` / `bank.go`），通过 `RegisterChannel` 自注册 | `model.SystemSetting` |
| `common/embed_theme.go` | 通过 `//go:embed web/build/default-pro` 把前端产物打包进二进制 | 无 |
| `model/` | GORM 模型 + 业务函数（**不依赖** HTTP / Gin） | `common/`, GORM |
| `controller/` | Gin Handler；只做参数解析、调 model、组装统一响应 | `model/`, `middleware/` |
| `middleware/` | auth / rate-limit / Turnstile / gzip / language / cache | `model/`, `common/` |
| `relay/` | LLM 转发层（按 provider 划分子包） | `model/`, `relay/registry` |
| `relay/adaptor/` | Provider 自注册（`init()` → `registry.Register`），统一 `Adaptor` 接口 | `relay/registry` |
| `relay/billing/` | 预扣 / 退还配额 | `model/` |
| `relay/registry/` | Provider 元数据中央索引（按 `ID` 与 legacy `type` 双向） | `relay/adaptor` |
| `router/` | 路由注册 | `controller/`, `middleware/` |
| `monitor/` | 巡检（成功率 / 渠道下线 / 集群健康） | `model/` |
| `web/default-pro/` | 前端项目（Vue 3 + Arco + Vite），`npm run build` 产物被 `//go:embed` | 无 |

> **前端仅在 `web/default-pro/` 维护**。`web/air/ web/default/ web/berry/ web/THEMES web/build/` 是历史主题或构建产物，禁改（详见 [开发规范](/zh/contribute/code-style)）。
> **Frontend is maintained in `web/default-pro/` ONLY.** The other `web/*` directories are historical themes or build output — do not edit.

## 请求生命周期 / Request Lifecycle

以「用户调用 `/v1/chat/completions`」为例：

The lifecycle for `POST /v1/chat/completions`:

```text
Client ──► Gin router
        ├─► middleware: recover · request-id · logger · language · cors · gzip
        ├─► middleware: rate-limit (GLOBAL_API_RATE_LIMIT)
        ├─► middleware: TokenAuth       // 解析 sk-xxx
        │       ├─► CacheGetTokenByKey  // Redis 命中 / miss 回源
        │       ├─► token 状态校验（启用 / 禁用 / 过期 / 耗尽）
        │       ├─► subnet / models 限制
        │       └─► user quota cache
        │
        ├─► relay/handler: 分流到 Chat / Embeddings / Image / Audio
        │       ├─► meta.Meta 组装（model、user、token、channel、isStream）
        │       ├─► 渠道路由（middleware/distributor）：权重 / 并发 / 冷却 / 粘性
        │       ├─► preConsumeQuota：user_quota 缓存 → PreConsumeTokenQuota
        │       ├─► adaptor.DoRequest → 上游 Provider
        │       ├─► adaptor.DoResponse：转换流 / JSON，提取 Usage
        │       └─► billing.PostConsumeQuota：差额清算 + 退款（失败回滚）
        │
        └─► middleware: ginzap / logger 落库 logs 表
```

要点 / Highlights：

- **三段鉴权**：Cookie Session（`/api/*`）、用户 Access Token（`/api/*` 兜底）、API Key（`/v1/*`）。详见 [API 总览 · 鉴权](/zh/api/README#附录-a鉴权机制)。
- **渠道路由**：默认按权重轮询；并发满则排队；冷却期内跳过；同 `user` 的连续请求粘到同一渠道。
- **计费分两段**：先预扣 `PreConsumedQuota + promptTokens`，结束时按真实 `usage` 清算（`PostConsumeQuota`），失败 / 异常通过 `ReturnPreConsumedQuota` 回滚。
- **缓存一致性**：`IncreaseUserQuota` 等加额链路同步刷 Redis `user_quota:<id>`，避免「库里有但缓存返回旧值」导致的持续 403。

## 数据模型 / Data Model

核心表（`model/*.go` 中 `AutoMigrate` 自动建表）：

Core tables (auto-created via `gorm.AutoMigrate` in `model/*.go`):

| 表 / Table | 关键字段 / Key fields | 说明 / Notes |
| --- | --- | --- |
| `users` | `role` (0/10/100)、`status` (1/2/3)、`quota`、`group`、`access_token` | 用户与角色 |
| `tokens` | `key`、`status`、`remain_quota`、`expired_time`、`models`、`subnet` | API Key 库 |
| `channels` | `type`、`key`、`base_url`、`weight`、`cooldown`、`models`、`group` | 渠道配置 |
| `abilities` | `(group, model_id, channel_id)` 唯一 | 「哪个组 / 哪个模型 / 哪个渠道」组合允许 |
| `redemptions` | `code`、`quota`、`used_user_id`、`expired_time` | 兑换码 |
| `plans` | `name`、`quota`、`price`、`duration_days`、`status` | 套餐 |
| `user_plans` | `(user_id, plan_id)`、`start_at`、`expire_at`、`remain_quota` | 用户订阅 |
| `plan_usages` | 套餐用量快照（按天聚合） | 报表 |
| `orders` | `order_no`、`type` (1 plan / 2 topup)、`status` (0/1/2/3)、`pay_method`、`pay_trade_no`、`source` (1/2) | 订单 |
| `logs` | `type` (consume/topup/manage/system/test)、`model_name`、`prompt/completion_tokens`、`quota`、`channel_id` | 调用日志 |
| `options` / `system_settings` | `(key, value, category)` | 通用 KV 配置 |
| `cluster_nodes` | `node_id`、`address`、`status`、`last_heartbeat`、`ping_failures` | Cluster 节点 |

## Provider 自注册 / Provider Self-Registration

Provider 是运行时通过 `init()` 自注册到 `relay/registry` 的：

Providers register themselves at runtime via `init()`:

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

新增 Provider 的步骤见 [新增 Provider](/zh/contribute/add-provider)。

See [Add a Provider](/en/contribute/add-provider) for the full walkthrough.

## 支付通道 / Payment Channels

支付通道采用 `common/payment.Channel` 接口 + `RegisterChannel` 自注册，新增通道零侵入：

Payment channels implement `common/payment.Channel` and self-register:

```go
// common/payment/wechat.go
func init() { RegisterChannel(&wechatChannel{}) }
```

下单链路统一走 `controller/buildPayInfo(payMethod, orderNo, amount, subject)`：

The order flow funnels through `controller/buildPayInfo(payMethod, orderNo, amount, subject)`:

1. 按 `pay_method` 选 `payment.New(payMethod)`；
2. 读 `system_settings` 检查 `IsEnabled()`；
3. 调 `channel.PrePay(orderNo, amount, subject)` 拿到 `pay_url / qr_code / trade_no`；
4. 异步回调 `controller/payment.go::processNotify`：`VerifyNotify` → `MarkOrderPaid` → 按 `order.Type` 触发套餐激活 / 充值入账。

新增通道的步骤见 [新增支付通道](/zh/contribute/add-payment)。

See [Add a Payment Channel](/en/contribute/add-payment).

## 集群模式 / Cluster Mode

`CLUSTER_ENABLED=true` 时启动「去中心化多活」模式：

When `CLUSTER_ENABLED=true`, One API Pro runs in decentralized multi-active mode:

- 每个节点独立 MySQL + Redis；
- GORM callbacks（`AfterCreate / AfterUpdate / AfterDelete`）捕获业务表变更；
- 背景 goroutine 把事件以 HTTP POST 推给所有存活节点；
- 接收端按 `updated_at` 决定是否覆盖（last-writer-wins）。

详见 [Cluster 概览](/zh/decentralization/overview)、[配置同步](/zh/decentralization/config-sync)、[节点健康](/zh/decentralization/node-health)。

See [Cluster Overview](/en/decentralization/overview), [Config Sync](/en/decentralization/config-sync), [Node Health](/en/decentralization/node-health).

## 扩展点 / Extension Points

| 想做的事 / What you want | 改哪里 / Where |
| --- | --- |
| 增加新 Provider | `relay/adaptor/provider/<name>/` + 实现 `Adaptor` 接口 |
| 增加新支付通道 | `common/payment/<name>.go` + 实现 `Channel` 接口 + `init()` 注册 |
| 增加新后台菜单 | `web/default-pro/src/router/index.js` + `web/default-pro/src/views/` |
| 改统一响应 | `controller/*.go` 中所有 `gin.H{"success":...}` |
| 改计费公式 | `relay/billing/` + `model.PreConsumeTokenQuota / PostConsumeQuota` |
| 加缓存键 | `model/cache.go` 中按 `model.cacheSeconds` 统一 TTL |

下一步 / Next: [开发环境搭建](/zh/contribute/dev-setup) · [贡献指南](/zh/contribute/code-style)。