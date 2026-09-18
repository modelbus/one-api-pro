---
title: 渠道路由
description: "路由链 Filter、冷却、并发、RPM、Sticky Session、fallback 与自动禁用。"
category: channel
order: 3
---

# 渠道路由

> 路由层把"可用渠道"按 Filter 链收缩成候选集，再由 Selector 选出一条最终渠道。实现见 `channelrouter/`。

## 路由链 / Pipeline

`channelrouter.ChannelRouter`（`channelrouter/router.go`）持有四类状态对象和一组 `ChannelFilter`，按以下顺序对候选集逐步收缩：

```
candidates
  → StatusFilter          // 仅 status=Enabled
  → FallbackFilter        // 排除 is_fallback=true
  → CooldownFilter        // 排除处于冷却中的渠道
  → ConcurrencyFilter     // 排除已达 max_concurrency 的渠道
  → RPMFilter             // 排除已达 rpm 的渠道
  → StickySessionFilter   // 若同 token 之前用过某渠道且仍可用，直接锁定
  → PriorityFilter        // IgnoreFirstPriority 时切除优先级最高的桶
```

通过全部 Filter 后，`PriorityRandomSelector`（`channelrouter/selector.go`）在剩下的桶内随机挑一条。

## 字段含义 / Field semantics

| 字段 | Filter | 行为 |
|---|---|---|
| `status` | `StatusFilter` | 仅 `ChannelStatusEnabled=1` 通过；其它状态被剔除 |
| `is_fallback` | `FallbackFilter` | 设为 true 后只在 fallback 路径中被选中；正常请求看不到 |
| `fallback_priority` | Fallback 路径 | 同为 fallback 时按升序选 |
| `priority` | `PriorityFilter` + Selector | 越大越靠前；同优先级桶内随机抽一条 |
| `weight` | — | 保留字段，当前以 `priority` 为主 |
| `max_concurrency` | `ConcurrencyFilter` | `<=0` 不限制；非集群模式下是节点内 in-memory 计数，集群模式下累加 `channel_counters` 表 |
| `rpm` | `RPMFilter` | `<=0` 不限制；按 60 秒滑窗计数（`channelrouter/rpm.go`） |
| `cooldown_seconds` | 触发冷却后的秒数 | 上游错误拦截器 `relay/interceptor/channel_action.go` 写入 `CooldownManager`，冷却时长被 `ChannelMaxCooldownSeconds`（默认 600）截断 |

## 冷却 / Cooldown

`channelrouter/cooldown.go::CooldownManager` 是 `sync.Map` 形态的冷却字典：

- `SetCooldown(channelId, seconds, reason, statusCode)` 写入条目；过期清理由后台 goroutine 每 30 秒跑一次
- `IsInCooldown(channelId)` 在 Filter 中查询；过期条目 lazy 删除
- 触发点：`relay/interceptor/channel_action.go::ChannelActionHandler`，根据上游 HTTP 状态码与错误类型（`insufficient_quota` / `authentication_error` / `permission_error` / `invalid_api_key` / 关键字 `credit`、`balance`、`已欠费` 等）写入冷却或触发自动禁用

`CooldownFilter` 不会剔除回包时还没到期的渠道，因此对 429 / 5xx / 鉴权失败有自然的退避效果。

## 自动禁用 / Auto-disable

`monitor.DisableChannel` 把状态置为 `ChannelStatusAutoDisabled=3`，并向 root 邮箱 / 消息推送发送通知。触发源：

1. **指标守护** — `monitor/metric.go` 维护每个渠道最近 `MetricQueueSize`（默认 10）次请求的成功率；累计满 10 次后若成功率 < `MetricSuccessRateThreshold`（默认 0.8）则调用 `MetricDisableChannel`。开关：`config.EnableMetric`。
2. **余额归零** — `controller/channel-billing.go::updateAllChannelsBalance` 在刷新余额时，若余额 `<= 0` 自动禁用。
3. **超时** — 批量测试时若响应时间 > `ChannelDisableThreshold`（默认 5s）且 `AutomaticDisableChannelEnabled=true` 则禁用。
4. **错误策略** — `ChannelActionHandler` 根据上游错误内容（如 401 / 403）触发禁用，受 `AutomaticDisableChannelEnabled` 控制。

被自动禁用的渠道只有管理员手动启用，或下一次批量测试通过后才会回到 `Enabled`。

## Fallback 路径 / Fallback path

主路径在所有正常渠道都返回失败后才进入 fallback：仅 `is_fallback=true` 的渠道参与，按 `fallback_priority` 升序、`priority` 次序依次重试。`FallbackFilter` 保证正常请求不会把 fallback 渠道选走。

## Sticky Session / 会话粘性

`channelrouter/sticky.go` 以 `MakeSessionKey(userId, model)` 作为键，缓存"用户 → 上次渠道"。请求时 `StickySessionFilter` 把候选集收缩成该渠道一条（若仍可用），否则回退到随机选择。请求完成后由路由上层调用 `SetStickySession(sessionKey, channelId)` 写入。

启用开关：`ChannelStickySessionEnabled`（默认 `false`）。启动时若开启，会从 `logs` 表加载近 24 小时 type=2 的会话记录重建映射（`LoadFromLogDB`）。

## 实现位置 / Implementation Pointers

| 关注点 | 位置 |
|---|---|
| 路由管线 | `channelrouter/router.go` |
| Filter 实现 | `channelrouter/filter.go` |
| 冷却管理 | `channelrouter/cooldown.go` |
| 并发计数 | `channelrouter/concurrency.go` |
| RPM 计数 | `channelrouter/rpm.go` |
| Sticky Session | `channelrouter/sticky.go` |
| 上游错误拦截 | `relay/interceptor/channel_action.go` |
| 自动禁用通知 | `monitor/channel.go` |
| 成功率指标 | `monitor/metric.go` |
