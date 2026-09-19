---
title: Add a Channel
description: "Field-by-field: what each setting is, where to get it, and what it changes."
category: channel
order: 2
---

# Add a Channel

> Each field on the New Channel form: what it is, where to get it, what it changes.

Open Admin → Channels → Add. The form has these fields.

## Required

### Provider type

- **What**: which upstream (OpenAI / Anthropic / Azure / Zhipu / DeepSeek / Ollama …).
- **Where**: pick from the dropdown.
- **Effect**: determines the request protocol (OpenAI-compatible / Anthropic Messages / Azure custom). **Wrong choice = nothing works.**

### Channel name

- **What**: a free-text label.
- **Where**: anything you like (e.g. `OpenAI Official`, `Azure East`, `Zhipu Production`).
- **Effect**: display only, no protocol impact.

### Base URL

- **What**: the upstream provider's endpoint URL.
- **Where**:
  - OpenAI: `https://api.openai.com/v1`
  - Azure OpenAI: from your deployment page, e.g. `https://<resource>.openai.azure.com/openai/deployments/<dep>`
  - Third-party proxy: the URL your vendor gives you
- **Effect**: wrong value = 404. Note whether `/v1` is required.

### API key / credential

- **What**: the upstream access token.
- **Where**: upstream console → API Keys / Access Tokens.
- **Effect**: wrong value = 401.
- **Tip**: copy whole, no quotes / spaces / newlines.

### Models

- **What**: which models this channel serves (multi-select).
- **Where**: options come from [Model Price](/en/pricing/model-price) entries you have enabled.
- **Effect**: unselected models are not routed here.
- **Tip**: set the Model Price first, otherwise calls succeed but no quota is deducted.

## Common

### Weight / Priority

- **What**: a number; higher = more likely to be picked.
- **Suggestion**: cheap-and-fast channels high, expensive-and-slow low.

### Concurrency cap

- **What**: max simultaneous in-flight calls.
- **Effect**: exceeding it trips the breaker / re-routes.
- **Suggestion**: match your upstream provider's RPM / TPM limits.

### Enabled

- **What**: a toggle.
- **Effect**: disabled = router skips.

## Advanced (rarely needed)

### Custom request headers

Only when the upstream requires a specific header (e.g. legacy Azure with `api-key`).

### Response JSONPath

For non-standard upstream responses; leave blank 99% of the time.

### Retry / timeout

Bump only if the upstream is slow or flaky.

## After save

1. Save.
2. Open the channel detail.
3. [Channel Test](./channel-test) — enter a model name; confirm success.
4. Once passing, enable the channel.

## Related

- [Channel Routing](./channel-routing) — how multi-channel selection works
- [Channel Test](./channel-test)
- [Balance Update](./balance-update)
- [Provider List](./provider-list)