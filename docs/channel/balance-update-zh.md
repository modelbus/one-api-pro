---
title: 余额刷新
description: "渠道余额的周期性拉取、Provider 实现差异与兜底计数。"
category: channel
order: 5
---

# 余额刷新

> `channels.balance` 字段由 `updateChannelBalance` 在周期任务里按 Provider 调用各自上游接口写入。实现：`controller/channel-billing.go`。

## 字段 / Field

`channels` 表：

| 字段 | 类型 | 说明 |
|---|---|---|
| `balance` | `float64` | 余额（美元，按 Provider 自行换算的展示值） |
| `balance_updated_time` | `int64` | 上次刷新时间（unix 秒）；列表里可据此判定是否过期 |

## 周期任务 / Scheduled refresh

`controller/channel-billing.go::AutomaticallyUpdateChannels(frequency int)` 是一个死循环：

```go
for {
  time.Sleep(time.Duration(frequency) * time.Minute)
  _ = updateAllChannelsBalance()
}
```

注意：当前仓库的 `main.go` 没有显式启动该 goroutine；需要外部（运维 cron / 第三方进程 / 后续 PR）以分钟级频率调用 `updateAllChannelsBalance()`，或自行 `go controller.AutomaticallyUpdateChannels(N)`。

## Provider 实现差异 / Per-provider endpoints

`updateChannelBalance(channel)` 按 `registry.IDByLegacyType(channel.Type)` 路由到对应实现（`controller/channel-billing.go`）：

| Provider | 端点 / 公式 |
|---|---|
| `openai` | `GET {base_url}/v1/dashboard/billing/subscription` + `usage`；`balance = HardLimitUSD - usage.TotalUsage / 100` |
| `azure` | 暂未实现，返回错误 `尚未实现` |
| `custom`（任意 OpenAI 兼容） | 同 `openai`，使用 `channel.base_url` |
| `closeai` | `GET {base_url}/dashboard/billing/credit_grants` |
| `openai-sb` | `GET https://api.openai-sb.com/sb-api/user/status?api_key=...`，解析 `data.credit` |
| `aiproxy` | `GET https://aiproxy.io/api/report/getUserOverview`，解析 `data.totalPoints` |
| `api2gpt` | `GET https://api.api2gpt.com/dashboard/billing/credit_grants` |
| `aigc2d` | `GET https://api.aigc2d.com/dashboard/billing/credit_grants` |
| `siliconflow` | `GET https://api.siliconflow.cn/v1/user/info`，解析 `data.totalBalance` |
| `deepseek` | `GET https://api.deepseek.com/user/balance`，取 `Currency=CNY` 的 `TotalBalance` |
| `openrouter` | `GET https://openrouter.ai/api/v1/credits`，`balance = total_credits - total_usage` |
| 其它 | 返回 `尚未实现` |

所有实现最终都会调 `channel.UpdateBalance(value)` 写回 `channels.balance` 与 `channels.balance_updated_time`。

## 余额归零自动禁用 / Auto-disable on empty

`updateAllChannelsBalance`（`controller/channel-billing.go:410`）对每个已启用渠道调用 `updateChannelBalance`；若 `balance <= 0`（含上游返回 `err=nil` 但余额非正的情况），调用 `monitor.DisableChannel(id, name, "余额不足")`，状态变为 `ChannelStatusAutoDisabled=3`。

`UpdateAllChannelsBalance` 当前路由是占位实现：直接返回 `success=true`，实际刷新由周期任务驱动。

## 手动刷新 / Manual refresh

| Endpoint | Method | Auth | 行为 |
|---|---|---|---|
| `/api/channel/update_balance/:id` | `GET` | Admin | 单条拉取一次并返回最新 `balance` |
| `/api/channel/update_balance` | `GET` | Admin | 当前为占位实现（立即返回 `success=true`，等待周期任务 / 后续接入） |

## 兜底 / Fallback

未实现余额接口的 Provider 不会抛错给上层调用方——`updateChannelBalance` 在分支未命中时直接返回 error，但 `updateAllChannelsBalance` 只 `continue`，跳过失败渠道。展示侧会一直保留旧值；管理员可结合 `balance_updated_time` 判定是否长时间未更新。

## 实现位置 / Implementation Pointers

| 关注点 | 位置 |
|---|---|
| 单 Provider 余额拉取 | `controller/channel-billing.go::updateChannelXxxBalance` |
| 入口路由 | `controller/channel-billing.go::updateChannelBalance` |
| 周期任务 | `controller/channel-billing.go::AutomaticallyUpdateChannels` |
| 写回 | `model/channel.go::UpdateBalance` |
| 余额归零自动禁用 | `controller/channel-billing.go::updateAllChannelsBalance` |
