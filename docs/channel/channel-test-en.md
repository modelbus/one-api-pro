---
title: Channel Test
description: "Single-channel and batch tests, the test request shape, and how the auto-disable path reacts."
category: channel
order: 4
---

# Channel Test

> Two endpoints: `GET /api/channel/test/:id?model=` for a single channel and `POST /api/channel/test?scope=all` for batch. Implementation: `controller/channel-test.go`.

## Single-channel test

| Field | Value |
|---|---|
| Method | `GET` |
| Path | `/api/channel/test/:id` |
| Query | `model` (optional; falls back to the first entry in `channel.models`) |
| Auth | Admin |

Server behavior (`controller/channel-test.go::TestChannel`):

1. `model.GetChannelById(id, true)` fetches the channel including `key`
2. `buildTestRequest(model)` constructs an OpenAI ChatCompletion body with one `user` message using `config.TestPrompt` (default `Output only your specific model name with no additional text.`)
3. If the requested model is not in `channel.models`, the first listed model is used; `model_mapping` then rewrites it to the upstream model name
4. `relay.GetAdaptorByChannel(channel.Type)` returns the provider adaptor; the test calls `ConvertRequest` → `DoRequest` → `DoResponse`
5. Parses the response and extracts `choices[0].content`; on success writes a `Log` via `RecordTestLog`
6. Stores the elapsed milliseconds on `channels.response_time`

Response on success:

```json
{
  "success": true,
  "message": "gpt-4o",
  "time": 0.842,
  "modelName": "gpt-4o"
}
```

On failure `success=false` and `message` carries the upstream status and error body.

## Batch test

| Field | Value |
|---|---|
| Method | `POST` |
| Path | `/api/channel/test` |
| Query | `scope` — `all` (default) / `disabled` |
| Auth | Admin |

`testChannels` uses a process-wide `testAllChannelsLock` so only one batch run is alive at any time. Flow:

1. Loads every channel (including disabled ones for `scope=all`); walks them with `config.RequestInterval` between each
2. After each channel's test:
   - If it was enabled and response time > `config.ChannelDisableThreshold * 1000` ms (default 5 s), auto-disable (when `AutomaticDisableChannelEnabled=true`) or send a notification
   - If it was enabled and `monitor.ShouldDisableChannel(openaiErr, -1)` decides it should be disabled, call `monitor.DisableChannel`
   - If it was disabled but this run succeeded, call `monitor.EnableChannel` to bring it back
3. When the run finishes and `notify=true`, root receives a "channel test completed" message

Response `success=true` means the test was started. Results land asynchronously on `channels.response_time` and `channels.last_error`.

## Scheduled test

`main.go:88` starts `controller.AutomaticallyTestChannels` whenever the `CHANNEL_TEST_FREQUENCY` env var (minutes) is set; it runs `testChannels(ctx, false, "all")` every N minutes.

## Test prompt

Globally configurable: `config.TestPrompt` (env `TEST_PROMPT`). In production use a short prompt that the model will answer cleanly, e.g.:

```
TEST_PROMPT="Output only your specific model name with no additional text."
```

## Implementation Pointers

| Concern | Location |
|---|---|
| Single-channel test | `controller/channel-test.go::TestChannel` |
| Batch test | `controller/channel-test.go::testChannels` |
| Scheduled runner | `controller/channel-test.go::AutomaticallyTestChannels` |
| Enable / disable decision | `monitor/manage.go::ShouldDisableChannel` / `ShouldEnableChannel` |
| Test request builder | `controller/channel-test.go::buildTestRequest` |
| Test log | `model.RecordTestLog` |
