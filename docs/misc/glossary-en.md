---
title: Glossary
description: "Glossary of One API Pro terms."
category: misc
order: 3
---

# Glossary

> Glossary of One API Pro terms.
> One API Pro 常用术语解释。

Grouped by topic; each entry has zh / en, meaning, source location, related code.

按主题分组；每条给出中英对照、含义、出现位置、相关代码。

## Roles & Permissions / 角色与权限

| Term | Meaning | Value | Constant | Source |
| --- | --- | --- | --- | --- |
| Guest / 访客 | Not signed in | 0 | `model.RoleGuestUser` | default |
| User / 普通用户 | Signed in, authenticated | 1 | (implicit) | default |
| Admin / 管理员 | Admin console access | 10 | `model.RoleAdminUser` | `middleware/auth.go::AdminAuth` |
| Root / 超级管理员 | System settings + admin management | 100 | `model.RoleRootUser` | `middleware/auth.go::RootAuth` |

The frontend uses `useAuthStore().isAdmin` / `isRoot` to gate UI; the backend middleware is the source of truth.

前端通过 `useAuthStore().isAdmin` / `isRoot` 控制 UI 显隐；后端中间件兜底。

## User Status / 用户状态

| Value | Meaning | Constant |
| --- | --- | --- |
| 0 | Default value (new user) | — |
| 1 | Enabled | `model.UserStatusEnabled` |
| 2 | Disabled | `model.UserStatusDisabled` |
| 3 | Deleted | `model.UserStatusDeleted` |

> Note: the constant values differ from the role values to avoid colliding with "Enabled" (which is 1).

## Token Status / 令牌状态

| Value | Meaning | Constant |
| --- | --- | --- |
| 1 | Enabled | `model.TokenStatusEnabled` |
| 2 | Disabled | `model.TokenStatusDisabled` |
| 3 | Expired | `model.TokenStatusExpired` |
| 4 | Quota exhausted | `model.TokenStatusExhausted` |

## Orders / 订单

| Term | Meaning | Value | Prefix |
| --- | --- | --- | --- |
| Plan Subscription / 套餐订阅 | New plan purchase | `OrderTypePlanSubscription` = 1 | `TB` |
| Topup / 充值 | Online top-up | `OrderTypeTopup` = 2 | `TP` |
| Upgrade Differential / 差价升级 | Plan upgrade differential | `OrderTypePlanSubscription` = 1 | `UP` |

| Status | Meaning | Constant |
| --- | --- | --- |
| 0 | Pending | `model.OrderStatusPending` |
| 1 | Paid | `model.OrderStatusPaid` |
| 2 | Canceled | `model.OrderStatusCanceled` |
| 3 | Refunded | `model.OrderStatusRefunded` |

Source / 来源:

| Enum | Meaning | Constant |
| --- | --- | --- |
| 1 | User self-service | `model.OrderSourceUserSelf` |
| 2 | Admin-placed | `model.OrderSourceAdmin` |

Pay-method constants / 支付方式常量:

```text
OrderPayMethodWechat  = "wechat"
OrderPayMethodAlipay  = "alipay"
OrderPayMethodBank    = "bank"      // reserved / 预留
OrderPayMethodOffline = "offline"   // offline collection
OrderPayMethodFree    = "free"      // admin free grant
```

## Log Type / 日志类型

`model/log.go`:

| Value | Meaning | Constant |
| --- | --- | --- |
| 0 | Unknown | `LogTypeUnknown` |
| 1 | Topup | `LogTypeTopup` |
| 2 | Consume | `LogTypeConsume` |
| 3 | Manage | `LogTypeManage` |
| 4 | System | `LogTypeSystem` |
| 5 | Test | `LogTypeTest` |

The Dashboard / Admin Dashboard only count `LogTypeConsume` by default.

「个人仪表盘」与「运营仪表盘」默认只统计 `LogTypeConsume`。

## System Settings / 系统设置

The `system_settings` table is a KV store with `(key, value, category)`. Common keys:

`system_settings` 表是 `key / value / category` 的 KV 库；常见 key：

| Constant | Default | Category | Notes |
| --- | --- | --- | --- |
| `payment.wechat.enabled` | false | payment | `{"enabled": bool}` |
| `payment.wechat.config` | — | payment | `{app_id, mch_id, api_key, notify_url, cert_file, key_file}` |
| `payment.alipay.enabled` | false | payment | same |
| `payment.alipay.config` | — | payment | `{app_id, private_key, public_key, notify_url}` |
| `payment.bank.enabled` | false | payment | reserved |
| `payment.bank.config` | — | payment | `{account_name, account_no, bank_name, branch, note}` |
| `plan.upgrade_mode` | `price_diff` | plan | `price_diff` / `stack` |
| `topup.enabled` | true | topup | Master switch for top-up |
| `topup.allow_custom` | true | topup | Allow custom amount |
| `topup.presets` | — | topup | JSON: `[{amount, bonus_quota}, ...]` |
| `topup.exchange_rate` | 1 | topup | 1 CNY = X quota |

## Plan Upgrade Mode / 套餐升级模式

| Mode | Meaning | Constant |
| --- | --- | --- |
| `price_diff` (default) | Differential: new price minus remaining value of the current plan | `OrderUpgradeModePriceDiff` |
| `stack` | Stack: pay full price for the new plan; the old plan keeps running until expiry | `OrderUpgradeModeStack` |

See [Subscription Upgrade & Downgrade](/en/subscription/upgrade-downgrade).

详见 [套餐升降级](/zh/subscription/upgrade-downgrade)。

## Quota / 配额

| Term | Meaning | Field |
| --- | --- | --- |
| Quota | Internal billing unit; the CNY ratio is `QuotaPerUnit = 500_000` (1 CNY = 500_000 quota) | — |
| Token quota | Per-API-key budget | `tokens.remain_quota` |
| User quota | User-level balance (top-up + unspent plan balance) | `users.quota` |
| Plan quota | Plan-assigned budget (`daily_quota` + total) | `plans.quota / plans.daily_quota` |

## Cache / 缓存

| Key | Meaning | TTL |
| --- | --- | --- |
| `token:<key>` | Token JSON | `SyncFrequency` |
| `user_group:<id>` | User group (determines visible models) | `SyncFrequency` |
| `user_quota:<id>` | User balance | `SyncFrequency` |
| `user_status:<id>` | User status | `SyncFrequency` |
| `group_models:<group>` | Group → models mapping | `SyncFrequency` |

`SyncFrequency` defaults to 600s; `userQuotaLowWaterMark = 50_000` decides when to fall back to DB.

## Cluster / 集群

| Term | Meaning | Field |
| --- | --- | --- |
| Node / 节点 | A deployed One API Pro instance | `cluster_nodes.node_id` |
| Heartbeat / 心跳 | `MarkHeartbeat()` writes `last_heartbeat` | `cluster_nodes.last_heartbeat` |
| Ping Failure / Ping 失败 | Cumulative ping failures | `cluster_nodes.ping_failures` |
| Push / 主动推送 | GORM callbacks capture table changes → HTTP POST sync | — |

## Auth / 鉴权

| Term | Meaning | Header |
| --- | --- | --- |
| Cookie Session | Browser-side session cookie | `Cookie` |
| Access Token | UUID, user-level admin token | `Authorization: <uuid>` |
| API Key | `sk-…`-style, OpenAI-compatible key | `Authorization: Bearer sk-…` |

See [API Overview · Auth](/en/api/README#appendix-a-auth-mechanisms).

详见 [API 总览 · 鉴权机制](/zh/api/README#附录-a鉴权机制)。

Next: [Troubleshooting](/en/faq/troubleshooting) · [Contribute](/en/contribute/dev-setup).