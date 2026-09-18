---
title: Add a Channel
description: "Channel fields, the create endpoint, and multi-key batch insert."
category: channel
order: 2
---

# Add a Channel

> `POST /api/channel/` — implementation: `controller/channel.go::AddChannel`.

## Endpoint

| Field | Value |
|---|---|
| Method | `POST` |
| Path | `/api/channel/` |
| Auth | Admin |
| Body | `Channel` JSON; the `key` field accepts `\n`-separated values for batch insert |

`AddChannel` (`controller/channel.go:80`) splits `key` on `\n` — one line becomes one channel row — and `model.BatchInsertChannels` then calls `AddAbilities()` to mirror every model in `models` into the `abilities` table.

## Fields

| Field | JSON type | Notes |
|---|---|---|
| `type` | `int` | Provider type enum (mirrors `CHANNEL_TYPE_MAP`: openai=1, claude=2, azure=3, gemini=4, baidu=5, aliyun=6, tencent=7, xunfei=8, zhipu=9, deepseek=10, midjourney=11, …) |
| `key` | `string` | Credential; sent as `Authorization: Bearer <key>`. Accepts `\n`-separated multi-key batch |
| `name` | `string` | Display name (indexed) |
| `base_url` | `*string` | Upstream API root. Empty falls back to provider default; required for OpenAI-compatible relays |
| `models` | `string` | Model allow-list, comma-separated. When testing, a model not in the list falls back to the first listed one |
| `group` | `string` | User-group allow-list, comma-separated. `ContainsGroup` does exact match; empty group is visible to everyone |
| `model_mapping` | `*string` | JSON object: `{ "source model": "upstream actual model" }`. Request model is rewritten before forwarding |
| `system_prompt` | `*string` | Prepended to the system message before forwarding (used by certain relay flows) |
| `weight` | `*uint` | Weighted-round-robin weight; the router currently keys off `priority`, `weight` is reserved |
| `priority` | `*int64` | Higher value sorts earlier within the same priority tier; the selector picks randomly within a tier |
| `max_concurrency` | `*int` | Per-node (or whole-cluster when clustering is on) cap; `<=0` means unlimited. Honored by `ConcurrencyFilter` |
| `cooldown_seconds` | `int` | Cooldown applied after an upstream error (default 60) |
| `rpm` | `*int` | Requests-per-minute cap; honored by `RPMFilter`. `<=0` means unlimited |
| `is_fallback` | `*bool` | When `true`, the channel is reserved for the fallback path and excluded from normal routing |
| `fallback_priority` | `*int64` | Order among fallback-only channels (ascending) |
| `config` | `string` | Provider-specific JSON (`ChannelConfig`: region / sk / ak / user_id / api_version / library_id / plugin / vertex_ai_*) |
| `status` | `int` | Default 1; see [Channel Routing](./channel-routing) |

> The default list endpoint `GetAllChannels` runs `Omit("key")`, so the frontend never sees the real credential. Only `GetChannel(..., selectAll=true)` returns it.

## Multi-key batch insert

Put one key per line in the `key` field:

```json
{
  "type": 1,
  "name": "OpenAI-bulk",
  "base_url": "https://api.openai.com",
  "models": "gpt-4o,gpt-4o-mini",
  "group": "default,vip",
  "key": "sk-AAA...\nsk-BBB...\nsk-CCC..."
}
```

Each line becomes an independent channel row; the batch fails as a whole on any insert error.

## `PUT /api/channel/`

`UpdateChannel` parses the raw JSON and updates only the keys that actually appear in the payload, so a partial update like `{id, status}` cannot wipe `name` / `models` / `group` / `config` / `balance`. See `controller/channel.go:151`.

## Frontend Guide

Route: `/channel` → **Add Channel** (`web/default-pro/src/views/channel/Channel.vue`).

- **Type** dropdown uses `CHANNEL_TYPE_MAP` + Provider list
- **Base URL** is required for OpenAI-compatible relays; official providers fall back to defaults when blank
- **Models** is populated from `GET /api/model_price/options` (AdminAuth) — only enabled `model_price` rows show up
- **Group** is a multi-select, defined in user management
- **Model Mapping** is edited as key/value pairs and serialized to JSON before submit
- **Max Concurrency / RPM / Cooldown** with value `<=0` mean "unlimited"

## Implementation Pointers

| Concern | Location |
|---|---|
| CRUD handlers | `controller/channel.go` |
| Data model | `model/channel.go::Channel` |
| Ability sync | `model/ability.go::AddAbilities` / `UpdateAbilities` |
| Partial-update safety | `controller/channel.go::UpdateChannel` |
