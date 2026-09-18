---
title: 运营设置
description: "配额策略、监控阈值、错误处理策略、邀请 / 注册 / 新用户额度。"
category: misc
order: 15
---

# 运营设置

> 全站运营策略：注册 / 邀请 / 新用户额度、监控阈值、渠道路由策略、错误处理策略。前端组件：`web/default-pro/src/views/setting/OperationSetting.vue`。

## 接口一览

| Endpoint | Method | 鉴权 | 说明 |
|---|---|---|---|
| `/api/option/` | `GET` | Root | 返回 `config.OptionMap` 中所有非 `Token` / `Secret` 键值对 |
| `/api/option/` | `PUT` | Root | `{ key, value }` 单条保存；按 key 触发对应校验 |

实现：`controller/option.go`；底层存储 `options` 表 + 内存 `config.OptionMap` 双写。

> 与 `system_settings` 表是平行的两套机制：`options` 是 key=value 字符串表，启动时加载到 `OptionMap`，对每个 `*Enabled` 开关同步影响 `config.*` 运行时变量；`system_settings` 是新版结构化设置（按 category 分组），由 `model.UpsertSystemSetting` 维护。本页只覆盖 `options` 表。

## 分区

页面按业务语义划分为若干 section，每个 section 内一组开关 + 数字字段 + 一个独立的「保存」按钮（部分开关在 change 事件里立即保存）：

| 分区 | 字段 |
|---|---|
| 额度策略 (Quota) | `QuotaForNewUser` / `PreConsumedQuota` / `QuotaForInviter` / `QuotaForInvitee` |
| 监控 (Monitor) | `ChannelDisableThreshold` / `QuotaRemindThreshold` + `AutomaticDisableChannelEnabled` / `AutomaticEnableChannelEnabled` / `LogConsumeEnabled` |
| 日志清理 | 日期选择器 + 「清理日志」按钮（调 `DELETE /api/log/?target_timestamp=...`） |
| 通用 (General) | `TopUpLink` / `ChatLink` / `QuotaPerUnit` / `RetryTimes` + `DisplayInCurrencyEnabled` / `DisplayTokenStatEnabled` / `ApproximateTokenEnabled` |
| 渠道路由 | `ChannelDefaultCooldownSeconds` / `ChannelMaxCooldownSeconds` + `ChannelConcurrencyEnabled` / `ChannelStickySessionEnabled` |
| 错误处理策略 | `ErrorNext` JSON：`{ passthrough, retry, disable, cooldown }` |
| 套餐 (Plan) | `plan.upgrade_mode`（详见 [plan-setting](plan-setting)） |

## 字段语义

| 字段 | 默认 | 用途 |
|---|---|---|
| `QuotaForNewUser` | env 启动值 | 注册赠送额度（quota） |
| `PreConsumedQuota` | 0 | 每次请求前预扣的 quota（请求失败时回退） |
| `QuotaForInviter` | 0 | 邀请人单次奖励 |
| `QuotaForInvitee` | 0 | 被邀请人注册奖励 |
| `ChannelDisableThreshold` | 0 | 单渠道连续失败 N 次自动禁用（需 `AutomaticDisableChannelEnabled=true`） |
| `QuotaRemindThreshold` | 0 | 余额低于此值时弹出提醒 |
| `QuotaPerUnit` | env 启动值 | UI 货币换算（每 `QuotaPerUnit` quota = 1 元） |
| `RetryTimes` | env 启动值 | 渠道级重试次数 |
| `ChannelDefaultCooldownSeconds` / `ChannelMaxCooldownSeconds` | env 启动值 | 渠道冷却区间 |
| `ErrorNext` | `{passthrough:true,retry:true,disable:true,cooldown:true}` | relay 错误链处理顺序（详见下） |

`TopUpLink` / `ChatLink` 用于「客户端主页」按钮的跳转链接。

## ErrorNext

存储为 JSON 对象 `{"passthrough": bool, "retry": bool, "disable": bool, "cooldown": bool}`：
- `passthrough`：透传错误到客户端（不会重试 / 切渠道）；
- `retry`：在同渠道内重试 `RetryTimes` 次；
- `disable`：连续失败达到阈值时禁用渠道；
- `cooldown`：失败后按 `ChannelDefaultCooldownSeconds` 冷却。

实际执行顺序由 relay 引擎读取；前端只负责把布尔集合保存为 JSON。

## 联动开关的额外校验

`PUT /api/option/` 在以下键启用前会要求前置配置已存在（`controller/option.go`）：

| Key | 前置要求 |
|---|---|
| `Theme` | 必须是 `config.ValidThemes` 内的合法主题名 |
| `GitHubOAuthEnabled` | `GitHubClientId` 与 `GitHubClientSecret` 均非空 |
| `EmailDomainRestrictionEnabled` | `EmailDomainWhitelist` 非空 |
| `WeChatAuthEnabled` | `WeChatServerAddress` 非空 |
| `TurnstileCheckEnabled` | `TurnstileSiteKey` 非空 |

不满足前置条件时返回明确文案，便于运营人员排错。

## 前端操作指南

- 顶部欢迎条 + 加载 spinner；进入页面即 `GET /api/option/`；
- 数字字段通过 `<a-input-number size="large">` 编辑；每个分区的「保存」按钮独立 PUT 该分区的全部字段；
- 开关字段大多用 `@change="saveSwitch(key)"` 立即单条 PUT；
- 「日志清理」选择器调用 `DELETE /api/log/?target_timestamp=<unix-seconds>`；target=0 时直接被后端拒绝；
- 「错误处理策略」分区整体保存为 JSON 字符串。

## 接口实现

| 关注点 | 位置 |
|---|---|
| Handler + 校验 | `controller/option.go` |
| 内存 OptionMap 初始化 | `model/option.go::InitOptionMap` / `SyncOptions` |
| 运行时同步更新 | `model/option.go::UpdateOption` → `updateOptionMap` |
| 与 `system_settings` 的关系 | 见 `model/system_setting.go` 头注释 |
| 路由 | `router/api.go` |
