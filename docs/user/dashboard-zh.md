---
title: 个人仪表盘
description: "消费统计、订阅状态、订单概览。"
category: user
order: 2
---

# 个人仪表盘

> 消费统计、订阅状态、订单概览。

路径 / Path：`/dashboard`（`web/default-pro/src/views/dashboard/Dashboard.vue`）。

## 页面分区

| 区域| 字段| 数据来源| 说明|
| --- | --- | --- | --- |
| 顶部欢迎条 | 用户名 / 角色 / 版本 | `useAuthStore` + `/api/status` | 角色 chip：用户 / 管理员|
| 核心指标（4 卡） | 总 Tokens / 总请求 / 总 Quota / 当前套餐 | `/api/log/self` 聚合 + `/api/user/self` + `/api/subscription/self` | 卡颜色：蓝 / 绿 / 橙 / 紫 |
| 用量进度 | 当日 / 当周（按设定阈值） | `/api/log/self` + `model.SumUsedQuota(LogTypeConsume, ...)` | 超过 80% 转红 |
| 趋势图 | 请求| `/api/log/self` | 默认显示最近 7 天 |
| 模型分布 Top-N | 调用量前 N 模型柱状图 | `/api/log/self` 聚合 | `n=10` |
| 使用明细 | 时间 / 模型 / 请求数| `/api/log/self` 分页 | 默认 `pageSize=8` |
| API Key 概览 | 最近一个令牌的脱敏 key + 创建时间 | `/api/token/self` | 点击复制 |
| 快捷入口 | 管理 Token / 兑换码 / 用量日志 | 路由跳转 | 管理员视角下额外显示「运营仪表盘」按钮 |

> 用量统计基于 `logs.type = LogTypeConsume` 过滤；充值（`LogTypeTopup`）、管理员加额（`LogTypeManage`）等不计入消费。

## 数据获取

```text
onMounted ──► Promise.all([
  api.get('/api/user/self'),                       // 用户信息
  api.get('/api/log/self', { params: { p: 0 } }),  // 日志（首屏 8 条）
  api.get('/api/log/self', { params: { p: 0, page_size: 50, type: 1 } }),  // 用于趋势/分布
  api.get('/api/token/self'),                      // 令牌概览
  api.get('/api/subscription/self'),               // 当前订阅
  api.get('/api/order/self', { params: { type: 2 } }),  // 最近充值订单（用于下次访问跳转）
])
```

并发请求避免瀑布流；`Promise.all` 全部 resolve 后再渲染，loading 由 `a-spin :loading="loading" style="width:100%"` 控制。

Concurrent fetch avoids waterfalls; gated by `a-spin :loading="loading" style="width:100%"`. (`arco-spin` 踩坑：`AGENTS.md` §10.4)

## 关键计算

| 名称| 公式| 说明|
| --- | --- | --- |
| `todayTokens` | `SUM(tokens_used)` 当日 `/ < 0, now>` | 每日 0 点重置（受 `TZ` 影响） |
| `sevendayTokens` | `SUM(tokens_used)` 最近 7 天 | 滚动 7 天窗口 |
| `todayPercent` | `todayTokens| `dailyQuota = plan.daily_quota`，未订阅为 0 |
| `planFoot` | 已订阅：到期日；已过期：「已过期，去续费」；未订阅：空 | 见 `statItems` 计算 |
| `quota` 格式 | `< 10000` 原样；`>= 10000` → `xxx.xx w` | 显示在「总 Quota」卡 |
| `tokens` 格式 | `< 1000` 原样；`>= 1000` → `xx.x K`；`>= 1e6` → `xx.x M` | 显示在用量明细 |

> `daily_quota` 是 `plans.daily_quota` 字段；若套餐未设置每日上限，`dailyPercent` 为 0，不展示进度条。

## 时间与时区

- 服务端 `helper.GetTimestamp()` 使用 `time.Now().Unix()`，未配置 `TZ` 时按 UTC；推荐设置 `TZ=Asia/Shanghai`。
- 仪表盘当天范围由前端按本地时区计算（`new Date()`），跨时区切换会有 ±1h 偏差。

## 自定义

- **隐藏/显示快捷入口**：编辑 `Dashboard.vue` 的 `quickActions` 计算属性。
- **调整 Top-N**：修改 `barOption` 中的 `slice(0, N)`。
- **修改卡片顺序**：编辑 `statItems` 数组；保持 `value / foot / icon` 一一对应。

## 常见问题

- **数据延迟几分钟**：日志写入是异步的（`middleware/logger.go`），但 v0.0.21 之后 `IncreaseUserQuota` 等关键加额链路已同步刷 Redis，不会有分钟级延迟。
- **显示「已过期」但还能用**：检查 `/api/subscription/self` 的 `expire_at` 与客户端时间；服务器 `TZ` 配置错误会触发此问题。

下一步 / Next: [Access Token](/zh/user/access-token) · [个人资料](/zh/user/profile) · [我的订单](/zh/user/orders)。
