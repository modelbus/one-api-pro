---
title: Channel
description: "Smallest unit for onboarding an upstream LLM provider: credentials, routing and stats."
category: schema
order: 4
---

# Channel

## What it is

A `Channel` is the smallest unit for onboarding an upstream provider (OpenAI, Anthropic, Azure, Ollama…). Each new channel becomes one selectable origin for [routing](/en/channel/channel-routing).

## Where to find it

- **Admin → Channels**: list, add, test, enable/disable, refresh balance
- **Admin → Dashboard → Top Users**: which channel each user's calls landed on
- **Call logs**: every call records the channel id it hit

## What you usually set per channel

These are the fields you fill in 99% of the time:

| Field | Where it comes from | Effect |
|---|---|---|
| Provider type | Pick from dropdown (OpenAI / Anthropic / Azure / Ollama …) | Determines request protocol. |
| Channel name | Anything you like | Display only. |
| Base URL | The provider's endpoint URL | Wrong value → 404. OpenAI: `https://api.openai.com/v1`. Azure: per deployment. |
| API key / credentials | Provider console → API Keys | Wrong value → 401. No quotes / spaces / newlines. |
| Models (multi-select) | Pick from [Model Price](/en/schema/model-price) | Only selected models can be routed here. |
| Priority / weight | Number 0–N | Higher = more likely to be picked. |
| Concurrency cap | Number | Calls beyond it trip the breaker / re-route. |
| Enabled | Toggle | Disabled channels are skipped. |

Most channels only need the above. The fields below are rarely used:

- **Custom request headers**: only if the provider requires a specific header (e.g. Azure `api-key`).
- **Response JSONPath**: parses provider responses; leave blank unless you know why.
- **Retry count / timeout**: bump only if the provider is slow or flaky.

## How billing works

Once a channel is hit, the deduction amount is decided by [Model Price](/en/schema/model-price). Channels store credentials and routing, not price.

## Auto balance refresh

For providers that expose an account-balance endpoint (OpenAI / Azure …), One API Pro can poll the upstream balance periodically and show it on the channel list — handy for spotting channels that are about to go empty. Enable via the "Auto update balance" toggle on the channel edit page.

## Related pages

- [Channel Overview](/en/channel/overview)
- [Add a Channel](/en/channel/add-channel)
- [Channel Test](/en/channel/channel-test)
- [Channel Routing](/en/channel/channel-routing)
- [Balance Update](/en/channel/balance-update)
- [Provider List](/en/channel/provider-list)

## Related API

- `GET /api/channel/` — list
- `POST /api/channel/` — create
- `PUT /api/channel/` — update
- `GET /api/channel/test/:id` — test connectivity