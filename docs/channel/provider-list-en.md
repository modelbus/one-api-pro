---
title: Provider List
description: All LLM providers adapted by One API Pro.
category: channel
order: 6
---

# Provider List

> Every provider has a built-in adapter. Pick the matching type when creating a channel.
>
> For **OpenAI-compatible** proxies (any third-party aggregator), always pick `OpenAI` — just set `base_url` + `models`.

## International

| Provider type | Used for | Notes |
|---|---|---|
| `openai` | OpenAI / any OpenAI-compatible proxy | Default type; set `base_url` only. |
| `anthropic` | Anthropic Claude | Messages API. |
| `azure` | Azure OpenAI | Needs deployment name / API key. |
| `gemini` | Google Gemini | Native protocol. |
| `openrouter` | OpenRouter (multi-model aggregator) | OpenAI protocol. |
| `groq` | Groq | OpenAI protocol. |
| `mistral` | Mistral AI | OpenAI protocol. |
| `cohere` | Cohere | OpenAI protocol. |
| `togetherai` | Together AI | OpenAI protocol. |
| `replicate` | Replicate | Custom protocol. |
| `aws` | AWS Bedrock | Custom protocol. |
| `vertexai` | Google Vertex AI | GCP credentials required. |
| `palm` | Google PaLM (legacy) | Will be deprecated. |
| `xai` | x.AI (Grok) | OpenAI protocol. |

## China

| Provider type | Used for | Notes |
|---|---|---|
| `deepseek` | DeepSeek | OpenAI protocol. |
| `moonshot` | Moonshot Kimi | OpenAI protocol. |
| `zhipu` | Zhipu GLM | Custom protocol. |
| `qwen` / `ali` / `alibailian` | Alibaba Qwen / Bailian | Multiple endpoints. |
| `doubao` | ByteDance Doubao | Custom protocol. |
| `baichuan` | Baichuan | Custom protocol. |
| `baidu` / `baiduv2` | ERNIE | Multi-version. |
| `tencent` | Tencent Hunyuan | Custom protocol. |
| `lingyiwanwu` | Zhipu Lingyi | Custom protocol. |
| `stepfun` | StepFun | Custom protocol. |
| `minimax` | MiniMax | Custom protocol. |
| `novita` | Novita AI (GPU aggregator) | Custom protocol. |

## Local / Self-hosted

| Provider type | Used for | Notes |
|---|---|---|
| `ollama` | Ollama (local inference) | OpenAI protocol. |
| `proxy` | Generic HTTP proxy | Pure passthrough. |

## Generic OpenAI-compatible proxies

Any third-party OpenAI-compatible service uses the `openai` type with:

- `Base URL`: the URL the proxy gives you
- `API Key`: the key the proxy issues
- `Models`: which models the proxy supports
- `Model Mapping` (optional): map official names to internal names

Example:

```json
{ "gpt-4o": "openai-gpt-4o-vip", "claude-sonnet-4": "anthropic-claude-3.5" }
```

## How to know whether your provider is supported

- Just look at the dropdown on the New Channel form.
- Or use [Channel Test](./channel-test) to verify.

## Adapting a new Provider

If your provider isn't listed AND the protocol isn't OpenAI-compatible, you need to adapt it:

- Write a `relay/adaptor/provider/<name>/` subpackage
- Implement the `Adaptor` interface + register in `init()`
- It appears in the dropdown automatically

See [Add a Provider](/en/contribute/add-provider).

## Related

- [Channel Overview](./overview)
- [Add a Channel](./add-channel)
- [Channel Routing](./channel-routing)