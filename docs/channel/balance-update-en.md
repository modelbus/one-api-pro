---
title: Balance Update
description: "Periodic balance refresh per channel, per-provider endpoints, and fallback behavior."
category: channel
order: 5
---

# Balance Update

> `channels.balance` is written by `updateChannelBalance` against each provider's upstream API. Implementation: `controller/channel-billing.go`.

## Field

`channels` table:

| Field | Type | Notes |
|---|---|---|
| `balance` | `float64` | Balance in USD (display unit depends on provider) |
| `balance_updated_time` | `int64` | Last refresh time (unix seconds); used by the UI to detect stale data |

## Scheduled refresh

`controller/channel-billing.go::AutomaticallyUpdateChannels(frequency int)` is an infinite loop:

```go
for {
  time.Sleep(time.Duration(frequency) * time.Minute)
  _ = updateAllChannelsBalance()
}
```

`main.go` does not currently launch this goroutine — drive it externally (ops cron, separate process, or a future PR), or call `go controller.AutomaticallyUpdateChannels(N)` from your own bootstrap.

## Per-provider endpoints

`updateChannelBalance(channel)` dispatches on `registry.IDByLegacyType(channel.Type)` (`controller/channel-billing.go`):

| Provider | Endpoint / formula |
|---|---|
| `openai` | `GET {base_url}/v1/dashboard/billing/subscription` + `usage`; `balance = HardLimitUSD - usage.TotalUsage / 100` |
| `azure` | Not implemented; returns `尚未实现` |
| `custom` (any OpenAI-compatible) | Same as `openai`, using `channel.base_url` |
| `closeai` | `GET {base_url}/dashboard/billing/credit_grants` |
| `openai-sb` | `GET https://api.openai-sb.com/sb-api/user/status?api_key=...`, parses `data.credit` |
| `aiproxy` | `GET https://aiproxy.io/api/report/getUserOverview`, parses `data.totalPoints` |
| `api2gpt` | `GET https://api.api2gpt.com/dashboard/billing/credit_grants` |
| `aigc2d` | `GET https://api.aigc2d.com/dashboard/billing/credit_grants` |
| `siliconflow` | `GET https://api.siliconflow.cn/v1/user/info`, parses `data.totalBalance` |
| `deepseek` | `GET https://api.deepseek.com/user/balance`, picks the `Currency=CNY` entry's `TotalBalance` |
| `openrouter` | `GET https://openrouter.ai/api/v1/credits`, `balance = total_credits - total_usage` |
| others | Returns `尚未实现` |

Every implementation eventually calls `channel.UpdateBalance(value)` to persist `channels.balance` and `channels.balance_updated_time`.

## Auto-disable on empty balance

`updateAllChannelsBalance` (`controller/channel-billing.go:410`) iterates enabled channels and calls `updateChannelBalance`. When `balance <= 0` (including the case where the upstream returned `err=nil` but balance is non-positive), it calls `monitor.DisableChannel(id, name, "余额不足")` and the status becomes `ChannelStatusAutoDisabled=3`.

`UpdateAllChannelsBalance` is currently a stub that returns `success=true` immediately — the real refresh is driven by the scheduled task.

## Manual refresh

| Endpoint | Method | Auth | Behavior |
|---|---|---|---|
| `/api/channel/update_balance/:id` | `GET` | Admin | Refreshes one channel and returns the new `balance` |
| `/api/channel/update_balance` | `GET` | Admin | Stub for now (returns `success=true` immediately) |

## Fallback

Providers without an implemented balance endpoint do not bubble errors up to the caller — `updateChannelBalance` returns an error from the unknown-branch path, and `updateAllChannelsBalance` simply `continue`s. The UI keeps the old value; combine with `balance_updated_time` to detect staleness.

## Implementation Pointers

| Concern | Location |
|---|---|
| Per-provider fetch | `controller/channel-billing.go::updateChannelXxxBalance` |
| Dispatcher | `controller/channel-billing.go::updateChannelBalance` |
| Scheduled task | `controller/channel-billing.go::AutomaticallyUpdateChannels` |
| Persist | `model/channel.go::UpdateBalance` |
| Auto-disable on empty | `controller/channel-billing.go::updateAllChannelsBalance` |
