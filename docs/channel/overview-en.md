---
title: Channel Overview
description: "Channel concept, lifecycle, and supported providers."
category: channel
order: 1
---

# Channel Overview

> A Channel is the smallest unit that wraps an upstream LLM provider credential plus its routing policy. It is what the request router selects, gates, and rate-limits.

## Concept

A Channel is the minimum unit of "upstream credential + routing policy" in One API Pro. It typically maps to one logical upstream account: DeepSeek official, Azure OpenAI deployment, an OpenAI-compatible relay, etc. Each channel independently stores:

- Credential and base URL: `key` / `base_url` / `config`
- Model allow-list (`models`, comma-separated)
- User-group allow-list (`group`, comma-separated)
- Model name mapping (`model_mapping`, JSON object)
- Routing knobs: weight (`weight`), priority (`priority`), max concurrency (`max_concurrency`), RPM cap (`rpm`), cooldown seconds (`cooldown_seconds`)
- Behavioral flags: fallback-only (`is_fallback` / `fallback_priority`), system prompt prefix (`system_prompt`)
- Runtime metrics: `status`, `balance` (USD), `response_time`, `last_error`

The model is defined in `model/channel.go::Channel`. On insert / update, every (channel, model) pair is mirrored to the `abilities` table so the router can filter without re-reading the channel row.

## Status

| Value | Constant | Meaning |
|---|---|---|
| 0 | `ChannelStatusUnknown` | Default value; never persisted (DB default is 1) |
| 1 | `ChannelStatusEnabled` | Eligible for routing |
| 2 | `ChannelStatusManuallyDisabled` | Disabled by an admin |
| 3 | `ChannelStatusAutoDisabled` | Auto-disabled by monitor (low balance / low success rate / timeout) |

Auto-disabled channels only re-enter routing after an admin enables them or a batch test passes.

## Supported providers

The complete list is on [Provider List](./provider-list). Any OpenAI-compatible relay can be added as the `openai` type — fill `base_url` and `models`, then use `model_mapping` to rename models if needed.

## See also

- [Add a Channel](./add-channel) — every field, end to end
- [Channel Routing](./channel-routing) — filters, cooldown, concurrency, RPM, sticky, fallback
- [Channel Test](./channel-test) — single and batch tests
- [Balance Update](./balance-update) — automatic and on-demand balance refresh
