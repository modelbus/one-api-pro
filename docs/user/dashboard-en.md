---
title: Dashboard
description: "Usage statistics, subscription status and order overview."
category: user
order: 2
---

# Dashboard

> Usage statistics, subscription status and order overview.
> 消费统计、订阅状态、订单概览。

Path: `/dashboard` (`web/default-pro/src/views/dashboard/Dashboard.vue`).

## Page Sections / 页面分区

| Section / 区域 | Field / 字段 | Source / 数据来源 | Notes / 说明 |
| --- | --- | --- | --- |
| Welcome strip | username / role / version | `useAuthStore` + `/api/status` | Role chip: User / Admin / Root |
| Hero metrics (4 cards) | Total Tokens / Total Requests / Total Quota / Current Plan | `/api/log/self` aggregate + `/api/user/self` + `/api/subscription/self` | Card colors: blue / green / orange / purple |
| Usage progress | Today / Last 7 days (by threshold) | `/api/log/self` + `model.SumUsedQuota(LogTypeConsume, ...)` | Turns red above 80% |
| Trend charts | Requests / Quota / Tokens three-line (echarts) | `/api/log/self` | Default last 7 days |
| Model Top-N | Bar chart of top N by call volume | `/api/log/self` aggregate | `n=10` |
| Usage details | Date / model / request count / quota / tokens table | `/api/log/self` paginated | Default `pageSize=8` |
| API Key overview | Masked key + creation time | `/api/token/self` | Click to copy |
| Quick actions | Manage tokens / redemption / usage log | Route navigation | Admins see an extra "Admin Dashboard" button |

> Consumption stats only count `logs.type = LogTypeConsume`. Top-ups (`LogTypeTopup`) and admin grants (`LogTypeManage`) are excluded.
> 用量统计基于 `logs.type = LogTypeConsume` 过滤；充值、管理员加额等不计入消费。

## Data Fetching / 数据获取

```text
onMounted ──► Promise.all([
  api.get('/api/user/self'),                       // user info
  api.get('/api/log/self', { params: { p: 0 } }),  // first-screen logs (8 rows)
  api.get('/api/log/self', { params: { p: 0, page_size: 50, type: 1 } }),  // trends/distribution
  api.get('/api/token/self'),                      // token overview
  api.get('/api/subscription/self'),               // current subscription
  api.get('/api/order/self', { params: { type: 2 } }),  // latest top-up order
])
```

Concurrent fetch avoids waterfalls; gated by `a-spin :loading="loading" style="width:100%"`. (`arco-spin` pitfall: `AGENTS.md` §10.4).

## Key Calculations / 关键计算

| Name / 名称 | Formula / 公式 | Notes / 说明 |
| --- | --- | --- |
| `todayTokens` | `SUM(tokens_used)` today `/ < 0, now>` | Resets at 00:00 local (`TZ`-dependent) |
| `sevendayTokens` | `SUM(tokens_used)` last 7 days | Rolling 7-day window |
| `todayPercent` | `todayTokens / dailyQuota` | `dailyQuota = plan.daily_quota`; 0 when no plan |
| `planFoot` | Subscribed: expiry date; Expired: "Expired — renew"; None: empty | See `statItems` |
| `quota` format | `< 10000` raw; `>= 10000` → `xxx.xx w` | Displayed on the "Total Quota" card |
| `tokens` format | `< 1000` raw; `>= 1000` → `xx.x K`; `>= 1e6` → `xx.x M` | In the detail table |

> `daily_quota` lives in `plans.daily_quota`; if unset, `dailyPercent` is 0 and the progress bar is hidden.
> 若套餐未设置每日上限，进度条不展示。

## Time & Timezone / 时间与时区

- Server uses `helper.GetTimestamp()` (`time.Now().Unix()`); set `TZ=Asia/Shanghai` in CN deployments.
  服务端 `helper.GetTimestamp()` 使用 `time.Now().Unix()`，未配置 `TZ` 时按 UTC。
- "Today" is computed client-side via `new Date()`; switching timezones mid-day can shift the boundary by ±1h.
  仪表盘当天范围由前端按本地时区计算，跨时区切换会有 ±1h 偏差。

## Customize / 自定义

- **Show/hide quick actions**: edit the `quickActions` computed in `Dashboard.vue`.
- **Adjust Top-N**: edit `slice(0, N)` inside `barOption`.
- **Reorder cards**: edit `statItems`, keeping `value / foot / icon` aligned.

## FAQ / 常见问题

- **Data may lag by minutes**: log writes are async, but quota updates are real-time after v0.0.21 (`IncreaseUserQuota` and friends now refresh Redis synchronously).
  数据延迟几分钟：日志写入是异步的；v0.0.21 之后 `IncreaseUserQuota` 等关键加额链路已同步刷 Redis。
- **"Expired" but still works**: check `expire_at` vs client time; wrong `TZ` is a common cause.
  显示「已过期」但还能用：检查 `expire_at` 与客户端时间；`TZ` 错误是常见原因。

Next: [Access Token](/en/user/access-token) · [Profile](/en/user/profile) · [My Orders](/en/user/orders).