---
title: 日志查询
description: "全站日志过滤（类型 / 模型 / 用户 / 渠道 / 时间窗）、聚合与清理。"
category: misc
order: 11
---

# 日志查询

> 在 `logs` 表上做全站过滤（admin）或个人过滤（user）、聚合消费 / 订阅维度 quota、定期清理历史。前端组件：`web/default-pro/src/views/log/Log.vue`（admin 与 user 复用，自动按角色显隐列）。

## 接口一览

| Endpoint | Method | 鉴权 | 说明 |
|---|---|---|---|
| `/api/log/` | `GET` | Admin | 分页拉取（默认 `config.ItemsPerPage`） |
| `/api/log/search?keyword=` | `GET` | Admin | 关键字模糊匹配 |
| `/api/log/stat` | `GET` | Admin | quota 聚合（普通 / 订阅 双维度） |
| `/api/log/?target_timestamp=` | `DELETE` | Admin | 删除 `created_at < target` 的所有日志（无 admin 鉴权 token 时会失败） |
| `/api/log/self` | `GET` | User | 当前用户日志（不分页管理员视角） |
| `/api/log/self/search?keyword=` | `GET` | User | 当前用户关键字搜索 |
| `/api/log/self/stat` | `GET` | User | 当前用户 quota 聚合 |

实现：`controller/log.go`。

## 通用过滤参数

| 参数 | 类型 | 说明 |
|---|---|---|
| `p` | `int` | 页码（0-based） |
| `type` | `int` | 日志类型；`LogTypeTopup=1` / `Consume=2` / `Manage=3` / `System=4` / `Test=5` |
| `start_timestamp` | `int64` | unix 秒（含），默认 `0` 表示无下界 |
| `end_timestamp` | `int64` | unix 秒（含），默认 `0` 表示无上界 |
| `username` | `string` | 精确匹配（admin 端点可传） |
| `token_name` | `string` | 模糊匹配 |
| `model_name` | `string` | 模糊匹配 |
| `channel` | `int` | `channel_id`（admin 端点可传） |

> `LogTypeConsume` 的 `quota` 才计入计费；其它类型只用于审计。

## 聚合

`/api/log/stat` 返回：

```

- `quota` = 全口径 `sum(quota)`；
- `normal_quota` = `billing_source=0` 的部分；
- `subscription_quota` = `billing_source=1` 的部分。

后端 `SumUsedQuotaByBillingSource(logType, start, end, model, username, token_name, channel, billing_source)`，对应 `model/log.go`。

## 清理

`DELETE /api/log/?target_timestamp=<unix-seconds>`：删除所有 `created_at < target_timestamp` 的日志行；返回 `{ data: <affected_rows> }`。
target_timestamp 必须显式传（`0` 会被直接拒绝）。生产环境建议搭配定时任务（`OperationSetting.vue` 提供日期选择器 + 按钮）。

## 类型说明

| 类型 | 说明 | 是否计入 quota |
|---|---|---|
| `LogTypeTopup=1` | 充值 / 兑换码到账 / 管理员手动加额度 | 否 |
| `LogTypeConsume=2` | API 调用消耗 | 是 |
| `LogTypeManage=3` | 管理员操作（删/禁/调额度/开通套餐/开通订阅等） | 否 |
| `LogTypeSystem=4` | 系统行为（新用户注册赠送、邀请赠送等） | 否 |
| `LogTypeTest=5` | 渠道测试（`POST /api/channel/test`） | 否 |

## 前端操作指南

- 顶部欢迎条 + 「类型」下拉（全部 / 充值 / 消耗 / 管理 / 系统 / 测试）。
- 搜索栏：
  - 普通用户：「查询输入」「令牌名」「模型名」「开始时间」「结束时间」；
  - 管理员额外：「用户名」「渠道 ID」。
- 「统计」按钮 toggle 出/隐 `quota / normal_quota / subscription_quota` 摘要卡。
- 列表列：
  - 时间（点击复制 `request_id`）/ 渠道（admin）/ 来源（subscription chip 蓝 / quota chip 灰）/ 套餐 / 类型 / 模型 / 用户（admin）/ 令牌名；
  - 若非测试类型：Prompt / Completion / Quota；
  - 详情：`content` + `elapsed_time`（ms，颜色按档位变化）/ `is_stream` / `system_prompt_reset` 标签。
- 分页：`[10, 20, 50]`，滚动到末尾自动追加下一页。
- 操作设置（`/setting/operation`）提供「日志清理」入口：选日期 → 调 `DELETE /api/log/?target_timestamp=...`。

## 接口实现

| 关注点 | 位置 |
|---|---|
| 日志 handler | `controller/log.go` |
| `LogStatistic` 类型 | `model/log.go` |
| 配额聚合 | `model/log.go::SumUsedQuota` / `SumUsedQuotaByBillingSource` |
| 删除 | `model/log.go::DeleteOldLog` |
| 清理入口（前端） | `web/default-pro/src/views/setting/OperationSetting.vue` |
| 路由 | `router/api.go` |
