---
title: Channel Overview
description: What a channel is, why you need multiple, and its states.
category: channel
order: 1
---

# Channel Overview

## What is a channel

A **channel is one upstream-provider endpoint**. Examples:

- For OpenAI you may have an official channel, an Azure self-hosted channel, a third-party proxy channel.
- For DeepSeek you may have an official channel and a reseller channel.

Each channel **independently stores credentials, base URL, allowed models, and routing weight**. When a request comes in, the [router](./channel-routing) picks one channel to forward to.

## Why have multiple channels

- **Redundancy** — one channel fails (out of credit / rate limited / down), traffic auto-routes elsewhere.
- **Cost** — cheaper channels first, expensive ones as fallback.
- **Region / compliance** — different regions / vendors.
- **Multiple accounts on one Provider** — spread quota, reduce risk.

## Channel states

| State | Set by | Meaning |
|---|---|---|
| Enabled | Admin | Routable. |
| Manually disabled | Admin | Admin clicked disable — router skips. |
| Auto disabled | System | Too many errors / out of credit / timed out — router skips; admin must re-enable. |

Auto-disabled is a **temporary breaker**, not a permanent shutdown. Re-enable it from Admin → Channels or run [Channel Test](./channel-test) — passing auto-clears the flag.

## What lives in a channel

The full list lives on [Add a Channel](./add-channel). In short:

- **Required**: provider type, base URL, API key, allowed models
- **Common**: weight / priority, concurrency cap
- **Rare**: custom headers, response JSONPath

Field-by-field explanation on [Add a Channel](./add-channel).

## Which providers are supported

Full list on [Provider List](./provider-list). Any **OpenAI-compatible** proxy works as the `openai` type: set `base_url` + `models`, then use `model_mapping` to fix model names.

## Next

- Create your first channel → [Add a Channel](./add-channel)
- Understand multi-channel selection → [Channel Routing](./channel-routing)
- Troubleshoot "why isn't this channel used" → [Channel Test](./channel-test)