---
title: Provider List
description: "All providers declared in the frontend constant, tag groupings, and the OpenAI-compat protocol coverage."
category: channel
order: 6
---

# Provider List

> Source of truth: `web/default-pro/src/constants/providers.js`. This file is also the source for the channel-type picker and the brand colors.

## Data shape

```js
{
  name: 'DeepSeek',
  slug: 'deepseek',
  color: '#4D6BFE',
  tag: '国产'
}
```

- `name` — display name
- `slug` — matches `relay/adaptor/provider/<slug>` on the backend
- `color` — badge color in the channel list
- `tag` — `国产` / `海外` grouping

## Domestic

| slug | name |
|---|---|
| `deepseek` | DeepSeek |
| `qwen` | Alibaba Qwen |
| `wenxin` | Baidu ERNIE |
| `chatglm` | Zhipu ChatGLM |
| `doubao` | ByteDance Doubao |
| `hunyuan` | Tencent Hunyuan |
| `spark` | iFlytek Spark |
| `moonshot` | Moonshot AI |
| `baichuan` | Baichuan |
| `stepfun` | Stepfun |
| `zeroone` | Lingyiwanwu |
| `minimax` | SenseTime |
| `internlm` | InternLM |
| `siliconcloud` | SiliconFlow |
| `aihubmix` | AIHubMix |

## Overseas

| slug | name |
|---|---|
| `openai` | OpenAI |
| `claude` | Anthropic Claude |
| `gemini` | Google Gemini |
| `mistral` | Mistral AI |
| `palm` | Meta LLaMA |
| `gemma` | Google Gemma |
| `grok` | xAI Grok |
| `cohere` | Cohere |
| `groq` | Groq |
| `ollama` | Ollama |
| `openrouter` | OpenRouter |
| `together` | Together AI |
| `novita` | Novita AI |
| `cloudflare` | Cloudflare |
| `ai21` | AI21 Labs |
| `stability` | Stability AI |
| `midjourney` | Midjourney |
| `dalle` | DALL·E |
| `replicate` | Replicate |
| `runway` | Runway |
| `suno` | Suno |
| `perplexity` | Perplexity |
| `phind` | Phind |
| `devin` | Devin |
| `huggingface` | Hugging Face |

## OpenAI-compatible coverage

Any OpenAI-compatible relay or self-hosted gateway does not need a new Provider entry — create a channel with `type=openai=1`, fill `base_url` + `models`, and use `model_mapping` to rewrite upstream model names. Routing, billing, and balance updates (in `updateChannelBalance` for type `openai` / `custom`) all reuse that `base_url`.

## Channel type map

The frontend exports `CHANNEL_TYPE_MAP` for the channel-create dialog:

```js
{
  openai: 1, claude: 2, azure: 3, gemini: 4,
  baidu: 5, aliyun: 6, tencent: 7, xunfei: 8,
  zhipu: 9, deepseek: 10, midjourney: 11
}
```

Backend `relay.GetAdaptorByChannel(channel.Type)` dispatches to `relay/adaptor/provider/<slug>`. OpenAI-compatible relays reuse the `openai` adaptor and override `base_url`.

## Implementation Pointers

| Concern | Location |
|---|---|
| Provider constants | `web/default-pro/src/constants/providers.js` |
| Backend adaptors | `relay/adaptor/provider/<slug>/` |
| Adaptor registration | `relay/adaptor/openai` etc. |
