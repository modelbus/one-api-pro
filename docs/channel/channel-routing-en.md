---
title: Channel Routing
description: "Filter chain, cooldown, concurrency, RPM, sticky sessions, fallback, and auto-disable."
category: channel
order: 3
---

# Channel Routing

> The routing layer shrinks the "available channels" set through a chain of filters, then a selector picks one. Implementation lives in `channelrouter/`.

## Pipeline

`channelrouter.ChannelRouter` (`channelrouter/router.go`) owns four stateful objects and a `[]ChannelFilter` chain. Each request walks the chain in order:

```
candidates
  → StatusFilter          // status == Enabled
  → FallbackFilter        // drop is_fallback=true
  → CooldownFilter        // drop channels in cooldown
  → ConcurrencyFilter     // drop channels at max_concurrency
  → RPMFilter             // drop channels at rpm cap
  → StickySessionFilter   // pin to the previous channel for this token, if still alive
  → PriorityFilter        // when IgnoreFirstPriority is set, drop the top priority tier
```

After all filters pass, `PriorityRandomSelector` (`channelrouter/selector.go`) picks one randomly within the top priority tier.

## Field semantics

| Field | Filter | Behavior |
|---|---|---|
| `status` | `StatusFilter` | Only `ChannelStatusEnabled=1` survives |
| `is_fallback` | `FallbackFilter` | When `true`, hidden from normal routing; only used by the fallback path |
| `fallback_priority` | fallback path | Ascending order among fallback-only channels |
| `priority` | `PriorityFilter` + Selector | Higher value ranks earlier; random within a tier |
| `weight` | — | Reserved; routing currently keys off `priority` |
| `max_concurrency` | `ConcurrencyFilter` | `<=0` means unlimited. Non-cluster mode uses an in-memory atomic counter; cluster mode sums across `channel_counters` rows |
| `rpm` | `RPMFilter` | `<=0` means unlimited. 60-second sliding window (`channelrouter/rpm.go`) |
| `cooldown_seconds` | default 60 | Cooldown written by the upstream error interceptor (`relay/interceptor/channel_action.go`); capped by `ChannelMaxCooldownSeconds` (default 600) |

## Cooldown

`channelrouter/cooldown.go::CooldownManager` is a `sync.Map` keyed by channel id:

- `SetCooldown(channelId, seconds, reason, statusCode)` writes an entry; a background goroutine prunes expired entries every 30 seconds
- `IsInCooldown(channelId)` is consulted by `CooldownFilter`; expired entries are removed lazily
- Trigger: `relay/interceptor/channel_action.go::ChannelActionHandler`, based on upstream HTTP status and error type (`insufficient_quota` / `authentication_error` / `permission_error` / `invalid_api_key` / keywords like `credit` / `balance` / `已欠费`)

`CooldownFilter` excludes any channel still in cooldown when the request returns, giving natural back-off for 429 / 5xx / auth failures.

## Auto-disable

`monitor.DisableChannel` flips status to `ChannelStatusAutoDisabled=3` and notifies root via email / message pusher. Triggers:

1. **Metric daemon** — `monitor/metric.go` keeps the most recent `MetricQueueSize` (default 10) outcomes per channel; once the buffer fills, success rate < `MetricSuccessRateThreshold` (default 0.8) calls `MetricDisableChannel`. Switch: `config.EnableMetric`.
2. **Balance exhausted** — `controller/channel-billing.go::updateAllChannelsBalance` disables a channel whose balance reaches `<= 0`.
3. **Slow response** — During batch tests, response time > `ChannelDisableThreshold` (default 5 s) with `AutomaticDisableChannelEnabled=true` triggers disable.
4. **Error policy** — `ChannelActionHandler` disables on hard upstream errors (401 / 403 / etc.) when `AutomaticDisableChannelEnabled` is on.

Auto-disabled channels only return to `Enabled` after an admin enables them or a later batch test passes.

## Fallback path

The main path retries within normal channels; when all of them fail, the router enters the fallback path. Only `is_fallback=true` channels are considered, ordered by `fallback_priority` ascending, then `priority`. `FallbackFilter` guarantees that normal traffic never picks a fallback channel.

## Sticky session

`channelrouter/sticky.go` keys on `MakeSessionKey(userId, model)` and caches the last channel id used for that token. `StickySessionFilter` shrinks the candidates to a single channel when the cached one is still available; otherwise falls back to random selection. After the relay completes, the upper layer calls `SetStickySession(sessionKey, channelId)`.

Switch: `ChannelStickySessionEnabled` (default `false`). On boot, when enabled, `LoadFromLogDB` rebuilds the map from the last 24 hours of `logs` rows where `type=2`.

## Implementation Pointers

| Concern | Location |
|---|---|
| Pipeline | `channelrouter/router.go` |
| Filters | `channelrouter/filter.go` |
| Cooldown | `channelrouter/cooldown.go` |
| Concurrency counter | `channelrouter/concurrency.go` |
| RPM counter | `channelrouter/rpm.go` |
| Sticky session | `channelrouter/sticky.go` |
| Upstream error interceptor | `relay/interceptor/channel_action.go` |
| Disable notification | `monitor/channel.go` |
| Success-rate metric | `monitor/metric.go` |
