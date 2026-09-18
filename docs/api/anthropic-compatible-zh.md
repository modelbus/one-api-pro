---
title: Anthropic 兼容接口
description: "/v1/messages 兼容接口。"
category: api
order: 19
---
# Anthropic 兼容接口

One API Pro 通过 `relay/adaptor/anthropic` 将 OpenAI 风格的 Chat Completions 请求（`/v1/chat/completions`）转换为 Anthropic Messages 协议，并兼容处理 Anthropic SSE 流式事件。下文描述 Anthropic 兼容适配器的请求/响应字段与流式事件约定，便于客户端按 Anthropic SDK 直连 Anthropic 风格模型。

> Anthropic Messages 的端点路径（`/v1/messages`）当前由 OpenAI Chat Completions 端点承担：在 `/v1/chat/completions` 使用 `claude-*` 模型时，请求体会先按 Anthropic Messages 协议发送给上游渠道，响应则归一化为 OpenAI Chat Completions 形态返回。本节描述归一化前的 Anthropic 形态（供客户端直接调用 Anthropic SDK + One API Pro 网关的混合场景参考）。

## 端点一览

| 接口 | 方法 | 权限 | 说明 |
|------|------|------|------|
| `/v1/chat/completions` (model=`claude-*`) | POST | Bearer Token | OpenAI Chat Completions 入口，命中 anthropic 适配器 |
| `/v1/models` | GET | Bearer Token | 模型列表（与 OpenAI 兼容） |


## 1. 请求格式（Anthropic Messages）

> 适配器内部将 OpenAI `Chat Completions` 请求转换为下列 Anthropic Messages 形态（参见 `relay/adaptor/anthropic/main.go::ConvertRequest`）。

**请求体：**

```json
{
  "model": "claude-3-5-sonnet",
  "messages": [
    { "role": "user", "content": "Hello" }
  ],
  "system": "You are a helpful assistant.",
  "max_tokens": 4096,
  "temperature": 0.7,
  "top_p": 0.9,
  "top_k": 40,
  "stop_sequences": ["\n\nHuman:"],
  "stream": false,
  "tools": [
    {
      "name": "get_weather",
      "description": "Get current weather",
      "input_schema": {
        "type": "object",
        "properties": { "city": { "type": "string" } },
        "required": ["city"]
      }
    }
  ],
  "tool_choice": { "type": "auto" }
}
```

**字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| model | string | 是 | 模型名；OpenAI 端常用 `claude-3-5-sonnet` / `claude-3-opus` 等 |
| messages | array | 是 | 对话历史，每项 `Message`，见下 |
| system | string | 否 | 系统提示词（也接受 `messages` 中的 `role=system` 消息，会被合并到此字段） |
| max_tokens | int | 否 | 最大输出 token；缺省时填充 4096 |
| temperature | float | 否 | 采样温度 |
| top_p | float | 否 | 核采样 |
| top_k | int | 否 | Top-K 采样 |
| stop_sequences | array&lt;string&gt; | 否 | 自定义停止序列 |
| stream | bool | 否 | 是否流式响应 |
| tools | array | 否 | 工具列表，每项 `Tool` |
| tool_choice | object/string | 否 | 工具选择；适配 OpenAI 的 `function` 形式或字符串 `auto`/`any` |

**`Message`：**

```json
{ "role": "user", "content": [{ "type": "text", "text": "Hello" }] }
```

| 字段 | 类型 | 说明 |
|------|------|------|
| role | string | `user` / `assistant`；OpenAI 的 `tool` 角色会被映射为 `user` + `tool_result` |
| content | array | 内容块数组，元素为 `Content` |

**`Content`：**

| 字段 | 类型 | 说明 |
|------|------|------|
| type | string | `text` / `image` / `tool_use` / `tool_result` |
| text | string | 文本内容（`type=text` / `type=tool_result`） |
| source | object | 图片源（`type=image`），`{type: "base64", media_type: "image/png", data: "..."}` |
| id | string | `tool_use` 的 id |
| name | string | `tool_use` 的工具名 |
| input | object | `tool_use` 的入参 |
| tool_use_id | string | `tool_result` 对应的 `tool_use.id` |
| content | string | `tool_result` 的文本结果 |

**`Tool`：**

| 字段 | 类型 | 说明 |
|------|------|------|
| name | string | 工具名 |
| description | string | 工具描述 |
| input_schema | object | 工具入参的 JSON Schema，含 `type` / `properties` / `required` |


## 2. 响应格式（非流式）

```json
{
  "id": "msg_01XYZ",
  "type": "message",
  "role": "assistant",
  "content": [
    { "type": "text", "text": "Hello! How can I help you today?" }
  ],
  "model": "claude-3-5-sonnet",
  "stop_reason": "end_turn",
  "stop_sequence": null,
  "usage": {
    "input_tokens": 12,
    "output_tokens": 18,
    "cache_read_input_tokens": 0,
    "cache_creation_input_tokens": 0
  },
  "error": { "type": "", "message": "" }
}
```

**字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 消息 ID |
| type | string | 固定 `message` |
| role | string | 固定 `assistant` |
| content | array | 内容块数组（同上） |
| model | string | 实际使用的模型 |
| stop_reason | string | `end_turn` / `stop_sequence` / `max_tokens` / `tool_use` |
| stop_sequence | string/null | 命中的停止序列 |
| usage.input_tokens | int | 输入 token |
| usage.output_tokens | int | 输出 token |
| usage.cache_read_input_tokens | int | 缓存读取 token |
| usage.cache_creation_input_tokens | int | 缓存写入 token |
| error | object | 错误对象，正常响应为空 |

> 适配器收到上游响应后，会把 `stop_reason` 映射为 OpenAI 的 `finish_reason`：`end_turn`/`stop_sequence` → `stop`，`max_tokens` → `length`，`tool_use` → `tool_calls`。响应在 `/v1/chat/completions` 上以 OpenAI `chat.completion` 形态返回。


## 3. 流式响应（SSE）

SSE 事件顺序与 Anthropic Messages 流式规范一致；客户端按 `event` 前缀的 `event` 类型分发：

```
event: message_start
data: {"type":"message_start","message":{...full message...}}

event: content_block_start
data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello"}}

event: content_block_stop
data: {"type":"content_block_stop","index":0}

event: message_delta
data: {"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":18}}

event: message_stop
data: {"type":"message_stop"}
```

| 事件 | 关键字段 | 说明 |
|------|----------|------|
| `message_start` | `message.{id,model,role,usage}` | 整条消息的元信息；`usage` 在此已是累计值 |
| `content_block_start` | `content_block.{type,text/id/name/input}` | 起始一个内容块（`text` 或 `tool_use`） |
| `content_block_delta` | `delta.{type,text/partial_json}` | 增量内容；`input_json_delta` 表示工具入参增量 JSON |
| `content_block_stop` | — | 当前内容块结束 |
| `message_delta` | `delta.{stop_reason,stop_sequence}`、`usage` | 累计的 usage（取 max 合并，避免重复计数） |
| `message_stop` | — | 整条消息结束 |

> 在 `/v1/chat/completions` 端点上，One API Pro 会把上述事件归一化为 OpenAI 的 `chat.completion.chunk` 形态，便于使用 OpenAI SDK 的客户端接入。


## 4. 鉴权与计费

- 鉴权：Bearer Token（即 `/api/token/` 创建的 `sk-xxxxxxxx`），与 OpenAI 兼容接口一致。
- 计费：按 OpenAI 侧的 `prompt_tokens` + `completion_tokens` 计算；`cache_read_input_tokens` / `cache_creation_input_tokens` 在 `model.ClaudeUsage2OpenAI` 中折算到 `prompt_tokens`。
- 重试 / Fallback：与 OpenAI 兼容接口一致（`config.RetryTimes`、`ErrorNext`），适用 `relay/handler` 通用重试策略。


## 5. 错误响应（来自上游）

```json
{
  "type": "error",
  "error": {
    "type": "invalid_request_error",
    "message": "messages: must be non-empty"
  }
}
```

经 `/v1/chat/completions` 中转后，最终返回 OpenAI 风格的错误体：

```json
{
  "error": {
    "message": "<msg> (request id: <uuid>)",
    "type": "one_api_error",
    "param": "",
    "code": "<upstream error type>"
  }
}
```