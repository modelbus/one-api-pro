---
title: Anthropic Compatible API
description: "/v1/messages compatible endpoints."
category: api
order: 19
---
# Anthropic Compatible API

One API Pro uses the `relay/adaptor/anthropic` package to translate OpenAI-style Chat Completions requests (`/v1/chat/completions`) into Anthropic Messages, and to consume Anthropic SSE streaming events. This page describes the Anthropic request / response and streaming-event conventions so clients can talk to Claude-style models through the gateway.

> The native `/v1/messages` endpoint is currently served by `/v1/chat/completions`: when a `claude-*` model is requested, the body is sent upstream in Anthropic Messages form, and the response is normalized into OpenAI Chat Completions shape. The Anthropic wire format documented below is the form the adapter speaks to the upstream.

## Endpoint index

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/v1/chat/completions` (model=`claude-*`) | POST | Bearer Token | OpenAI Chat Completions entry, hits the anthropic adapter |
| `/v1/models` | GET | Bearer Token | Model list (OpenAI-compatible) |


## 1. Request Format (Anthropic Messages)

> Internally the adapter converts OpenAI `Chat Completions` requests into the Anthropic Messages shape below (`relay/adaptor/anthropic/main.go::ConvertRequest`).

**Request body:**

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

**Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| model | string | yes | Model name. Common: `claude-3-5-sonnet`, `claude-3-opus`, etc. |
| messages | array | yes | Conversation history, see `Message` below |
| system | string | no | System prompt. A `role=system` entry in `messages` is folded into this field. |
| max_tokens | int | no | Max output tokens. Filled with 4096 when omitted. |
| temperature | float | no | Sampling temperature |
| top_p | float | no | Nucleus sampling |
| top_k | int | no | Top-K sampling |
| stop_sequences | array&lt;string&gt; | no | Custom stop sequences |
| stream | bool | no | Whether to stream the response |
| tools | array | no | Tool list, each a `Tool` |
| tool_choice | object/string | no | Tool choice; OpenAI `function` shape or string `auto` / `any` |

**`Message`:**

```json
{ "role": "user", "content": [{ "type": "text", "text": "Hello" }] }
```

| Field | Type | Description |
|-------|------|-------------|
| role | string | `user` / `assistant`. OpenAI's `tool` role is mapped to `user` + `tool_result`. |
| content | array | Content blocks (`Content`) |

**`Content`:**

| Field | Type | Description |
|-------|------|-------------|
| type | string | `text` / `image` / `tool_use` / `tool_result` |
| text | string | Text payload (`type=text` or `type=tool_result`) |
| source | object | Image source (`type=image`); `{type:"base64", media_type:"image/png", data:"..."}` |
| id | string | `tool_use` id |
| name | string | `tool_use` tool name |
| input | object | `tool_use` arguments |
| tool_use_id | string | Matching `tool_use.id` for a `tool_result` |
| content | string | Text payload of a `tool_result` |

**`Tool`:**

| Field | Type | Description |
|-------|------|-------------|
| name | string | Tool name |
| description | string | Tool description |
| input_schema | object | JSON Schema for the tool arguments; contains `type` / `properties` / `required` |


## 2. Response Format (non-streaming)

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

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| id | string | Message ID |
| type | string | Always `message` |
| role | string | Always `assistant` |
| content | array | Content blocks (same shape as above) |
| model | string | Actual model used |
| stop_reason | string | `end_turn` / `stop_sequence` / `max_tokens` / `tool_use` |
| stop_sequence | string/null | Stop sequence hit, if any |
| usage.input_tokens | int | Input tokens |
| usage.output_tokens | int | Output tokens |
| usage.cache_read_input_tokens | int | Cache read tokens |
| usage.cache_creation_input_tokens | int | Cache write tokens |
| error | object | Error object; empty on success |

> The adapter maps `stop_reason` into OpenAI `finish_reason`: `end_turn` / `stop_sequence` → `stop`, `max_tokens` → `length`, `tool_use` → `tool_calls`. The response is returned on `/v1/chat/completions` in OpenAI `chat.completion` shape.


## 3. Streaming Response (SSE)

SSE event order follows the Anthropic Messages streaming spec; clients dispatch by the `event` prefix:

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

| Event | Key fields | Description |
|-------|------------|-------------|
| `message_start` | `message.{id,model,role,usage}` | Message metadata; `usage` is cumulative. |
| `content_block_start` | `content_block.{type,text/id/name/input}` | Opens a content block (`text` or `tool_use`). |
| `content_block_delta` | `delta.{type,text/partial_json}` | Incremental content; `input_json_delta` carries partial JSON for tool args. |
| `content_block_stop` | — | Closes the current content block. |
| `message_delta` | `delta.{stop_reason,stop_sequence}`, `usage` | Cumulative usage (max-merged to avoid double counting). |
| `message_stop` | — | Closes the whole message. |

> At `/v1/chat/completions`, One API Pro normalizes these events into OpenAI `chat.completion.chunk` shape so OpenAI-SDK clients can consume them directly.


## 4. Auth and Billing

- **Auth:** Bearer Token (i.e. `sk-xxxxxxxx` issued by `/api/token/`), identical to the OpenAI-compatible endpoints.
- **Billing:** Computed from OpenAI-side `prompt_tokens` + `completion_tokens`; `cache_read_input_tokens` / `cache_creation_input_tokens` are folded into `prompt_tokens` by `model.ClaudeUsage2OpenAI`.
- **Retry / Fallback:** Same as the OpenAI-compatible endpoints (`config.RetryTimes`, `ErrorNext`), driven by the generic `relay/handler` retry chain.


## 5. Error Response (from upstream)

```json
{
  "type": "error",
  "error": {
    "type": "invalid_request_error",
    "message": "messages: must be non-empty"
  }
}
```

After relaying through `/v1/chat/completions`, the final body uses the standard OpenAI error shape:

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