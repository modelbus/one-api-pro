---
title: 渠道测试
description: "单渠道测试与批量测试、测试请求体构造、自动禁用联动。"
category: channel
order: 4
---

# 渠道测试

> 测试入口：`GET /api/channel/test/:id?model=`（单渠道）、`POST /api/channel/test?scope=all`（批量）。实现：`controller/channel-test.go`。

## 单渠道测试 / Single channel

请求：

| 项 | 值 |
|---|---|
| Method | `GET` |
| Path | `/api/channel/test/:id` |
| Query | `model`（可选；不传时取该渠道 `models` 列表的第一个） |
| Auth | Admin |

服务端行为（`controller/channel-test.go::TestChannel`）：

1. `model.GetChannelById(id, true)` 拉取完整渠道（含 `key`）
2. `buildTestRequest(model)` 构造 OpenAI ChatCompletion 请求体，仅含 `model` + 一条 `user` 消息（`config.TestPrompt`，默认 `Output only your specific model name with no additional text.`）
3. 若请求模型不在渠道 `models` 列表里，回退到列表第一个；然后用 `model_mapping` 改写到上游实际模型名
4. `relay.GetAdaptorByChannel(channel.Type)` 拿到 Provider adaptor，调 `ConvertRequest` → `DoRequest` → `DoResponse`
5. 解析响应取首条 `choices[0].content`；测试成功后写一条 `Log`（`RecordTestLog`）
6. 把耗时（毫秒）写回 `channels.response_time`

返回结构：

```json
{
  "success": true,
  "message": "gpt-4o",
  "time": 0.842,
  "modelName": "gpt-4o"
}
```

失败时 `success=false`，`message` 携带上游 HTTP 状态码与错误内容。

## 批量测试 / Batch test

请求：

| 项 | 值 |
|---|---|
| Method | `POST` |
| Path | `/api/channel/test` |
| Query | `scope`：`all`（默认）/ `disabled` |
| Auth | Admin |

实现 `testChannels` 用 `testAllChannelsLock` 保证全进程内只有一个测试在跑。流程：

1. 拉取全部渠道（`scope=all` 时含已禁用）；按顺序、间隔 `config.RequestInterval` 测试
2. 每个渠道测试结束后：
   - 若本已启用且响应时间 > `config.ChannelDisableThreshold * 1000` ms（默认 5 s），按 `AutomaticDisableChannelEnabled` 选择自动禁用或仅发通知
   - 若本已启用且 `monitor.ShouldDisableChannel(openaiErr, -1)` 判定为应禁用，调 `monitor.DisableChannel`
   - 若本已禁用但本轮成功，调 `monitor.EnableChannel` 重新启用
3. 全部结束后若 `notify=true`，向 root 发邮件 / 消息「渠道测试完成」

返回 `success=true` 表示测试已启动（实际结果异步落到 `channels.response_time` 与 `channels.last_error`）。

## 周期自动测试 / Scheduled test

`main.go:88` 在环境变量 `CHANNEL_TEST_FREQUENCY`（分钟）非空时启动 `controller.AutomaticallyTestChannels`，每 N 分钟跑一次 `testChannels(ctx, false, "all")`。

## 测试 Prompt / Test prompt

全局可配：`config.TestPrompt`（环境变量 `TEST_PROMPT`）。生产环境建议改成能稳定区分模型的短句，例如：

```
TEST_PROMPT="Output only your specific model name with no additional text."
```

## 实现位置 / Implementation Pointers

| 关注点 | 位置 |
|---|---|
| 单渠道测试 | `controller/channel-test.go::TestChannel` |
| 批量测试 | `controller/channel-test.go::testChannels` |
| 周期任务 | `controller/channel-test.go::AutomaticallyTestChannels` |
| 启用 / 禁用判定 | `monitor/manage.go::ShouldDisableChannel` / `ShouldEnableChannel` |
| 测试请求构造 | `controller/channel-test.go::buildTestRequest` |
| 测试日志 | `model.RecordTestLog` |
