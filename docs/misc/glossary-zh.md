---
title: 术语表
description: "One API Pro 常用术语解释。"
category: misc
order: 3
---

# 术语表

> One API Pro 常用术语解释。

按主题分组；每条给出中英对照、含义、出现位置、相关代码。

## 角色与权限

| 术语| 含义| 值| 常量| 出现|
| --- | --- | --- | --- | --- |
| Guest / 访客 | 未登录用户 | 0 | `model.RoleGuestUser` | 默认 |
| User / 普通用户 | 已登录、已认证 | 1 | （隐式） | 默认 |
| Admin / 管理员 | 后台管理权限 | 10 | `model.RoleAdminUser` | `middleware/auth.go::AdminAuth` |
| Root / 超级管理员 | 系统设置 + 管理员管理 | 100 | `model.RoleRootUser` | `middleware/auth.go::RootAuth` |

前端通过 `useAuthStore().isAdmin` / `isRoot` 控制 UI 显隐；后端中间件兜底。

## 用户状态

| 值| 含义| 常量|
| --- | --- | --- |
| 0 | 默认值（新建用户） | — |
| 1 | 启用| `model.UserStatusEnabled` |
| 2 | 封禁| `model.UserStatusDisabled` |
| 3 | 已注销| `model.UserStatusDeleted` |

> **注**：常量值 ≠ 角色值，因为「未指定」默认是 0，要避免与「启用」撞车。

## 令牌状态

| 值| 含义| 常量|
| --- | --- | --- |
| 1 | 启用| `model.TokenStatusEnabled` |
| 2 | 禁用| `model.TokenStatusDisabled` |
| 3 | 过期| `model.TokenStatusExpired` |
| 4 | 额度耗尽| `model.TokenStatusExhausted` |

## 订单

| 术语| 含义| 值| 订单号前缀|
| --- | --- | --- | --- |
| 套餐订阅| 新订套餐 | `OrderTypePlanSubscription` = 1 | `TB` |
| 充值| 在线支付充值 | `OrderTypeTopup` = 2 | `TP` |
| 差价升级| 同套餐升级 | `OrderTypePlanSubscription` = 1 | `UP` |

| 状态值| 含义| 常量|
| --- | --- | --- |
| 0 | 待支付| `model.OrderStatusPending` |
| 1 | 已支付| `model.OrderStatusPaid` |
| 2 | 已取消| `model.OrderStatusCanceled` |
| 3 | 已退款| `model.OrderStatusRefunded` |

来源 / Source / 支付方式 / Pay method：

| 枚举| 含义| 常量|
| --- | --- | --- |
| 1 | 用户自助 | `model.OrderSourceUserSelf` |
| 2 | 管理员下单 | `model.OrderSourceAdmin` |

支付渠道常量 / Pay method constants：

```text
OrderPayMethodWechat  = "wechat"
OrderPayMethodAlipay  = "alipay"
OrderPayMethodBank    = "bank"      // 预留 / reserved
OrderPayMethodOffline = "offline"   // 管理员线下收款 / offline collection
OrderPayMethodFree    = "free"      // 管理员免费赠送 / admin free grant
```

## 日志类型

`model/log.go` 定义：

| 值| 含义| 常量|
| --- | --- | --- |
| 0 | 未知| `LogTypeUnknown` |
| 1 | 充值| `LogTypeTopup` |
| 2 | 消费| `LogTypeConsume` |
| 3 | 管理| `LogTypeManage` |
| 4 | 系统| `LogTypeSystem` |
| 5 | 测试| `LogTypeTest` |

「个人仪表盘」与「运营仪表盘」默认只统计 `LogTypeConsume`。

## 系统设置

`system_settings` 表是 `key / value / category` 的 KV 库；常见 key：

| Key 常量| 默认| 类别| 备注|
| --- | --- | --- | --- |
| `payment.wechat.enabled` | false | payment | `{"enabled": bool}` |
| `payment.wechat.config` | — | payment | `{app_id, mch_id, api_key, notify_url, cert_file, key_file}` |
| `payment.alipay.enabled` | false | payment | 同上 |
| `payment.alipay.config` | — | payment | `{app_id, private_key, public_key, notify_url}` |
| `payment.bank.enabled` | false | payment | 预留 |
| `payment.bank.config` | — | payment | `{account_name, account_no, bank_name, branch, note}` |
| `plan.upgrade_mode` | `price_diff` | plan | `price_diff` / `stack` |
| `topup.enabled` | true | topup | 充值总开关 |
| `topup.allow_custom` | true | topup | 是否允许自定义金额 |
| `topup.presets` | — | topup | JSON: `[{amount, bonus_quota}, ...]` |
| `topup.exchange_rate` | 1 | topup | 1 元 = X quota |

## 套餐升级模式

| 模式| 含义| 常量|
| --- | --- | --- |
| `price_diff`（默认） | 差价升级：新套餐价 − 旧套餐剩余价值 | `OrderUpgradeModePriceDiff` |
| `stack` | 叠加：新套餐全价，旧订阅继续生效至过期 | `OrderUpgradeModeStack` |

详见 [套餐升降级](/zh/subscription/upgrade-downgrade)。

See [套餐升降级](/en/subscription/upgrade-downgrade).

## 配额

| 术语| 含义| 字段|
| --- | --- | --- |
| Quota | 内部计费单位；与「元」换算由 `QuotaPerUnit = 500_000` 决定（每元 = 500_000 quota） | — |
| Token quota | 某 API Key 自带的额度 | `tokens.remain_quota` |
| User quota | 用户级余额（含充值、套餐未消耗部分） | `users.quota` |
| Plan quota | 套餐分配的额度（`daily_quota` + 总配额） | `plans.quota|

## 缓存

| 键| 含义| TTL |
| --- | --- | --- |
| `token:<key>` | Token JSON 序列化 | `SyncFrequency` |
| `user_group:<id>` | 用户分组（决定可见模型） | `SyncFrequency` |
| `user_quota:<id>` | 用户余额 | `SyncFrequency` |
| `user_status:<id>` | 用户状态 | `SyncFrequency` |
| `group_models:<group>` | 用户组可用模型映射 | `SyncFrequency` |

`SyncFrequency` 默认 600s；`userQuotaLowWaterMark = 50_000` 决定何时回源。

## 集群

| 术语| 含义| 字段|
| --- | --- | --- |
| 节点| 一台部署的 One API Pro 实例 | `cluster_nodes.node_id` |
| 心跳| 节点间 `MarkHeartbeat()` 写入 `last_heartbeat` | `cluster_nodes.last_heartbeat` |
| Ping 失败| 探测失败累计 | `cluster_nodes.ping_failures` |
| 主动推送| GORM callbacks 捕获表变更 → HTTP POST 同步 | — |

## 鉴权

| 术语| 含义| 头部|
| --- | --- | --- |
| Cookie Session | 浏览器登录后的会话 Cookie | `Cookie` |
| Access Token | UUID，用户级管理接口 Token | `Authorization: <uuid>` |
| API Key | `sk-…` 形式，OpenAI 兼容 Key | `Authorization: Bearer sk-…` |

详见 [API 总览 · 鉴权机制](/zh/api/README#附录-a鉴权机制)。

下一步 / Next: [故障排查](/zh/faq/troubleshooting) · [贡献指南](/zh/contribute/dev-setup)。
