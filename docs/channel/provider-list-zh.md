---
title: Provider 一览
description: "前端常量中的全部 Provider、tag 分类与 OpenAI 兼容协议覆盖。"
category: channel
order: 6
---

# Provider 一览

> 数据源：`web/default-pro/src/constants/providers.js`。该文件同时是前端的展示色板与渠道类型选项的真相源。

## 数据结构 / Data shape

```js
{
  name: 'DeepSeek',
  slug: 'deepseek',
  color: '#4D6BFE',
  tag: '国产'
}
```

- `name` — 显示名
- `slug` — 与后端 `relay/adaptor/provider/<slug>` 子包一一对应，可作渠道 `type` 候选
- `color` — 列表徽标颜色
- `tag` — 「国产」/「海外」分组

## 国产 / Domestic

| slug | name |
|---|---|
| `deepseek` | DeepSeek |
| `qwen` | 阿里通义千问 |
| `wenxin` | 百度文心一言 |
| `chatglm` | 智谱 ChatGLM |
| `doubao` | 字节豆包 |
| `hunyuan` | 腾讯混元 |
| `spark` | 讯飞星火 |
| `moonshot` | 月之暗面 Moonshot |
| `baichuan` | 百川智能 |
| `stepfun` | 阶跃星辰 |
| `zeroone` | 零一万物 |
| `minimax` | 商汤日日新 |
| `internlm` | InternLM |
| `siliconcloud` | 硅基流动 |
| `aihubmix` | AIHubMix 推理时代 |

## 海外 / Overseas

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

## OpenAI 兼容协议 / OpenAI-compatible coverage

任意 OpenAI 兼容的中转站 / 自部署网关不必新增 Provider；只需在新增渠道时把 `type` 设为 `openai=1`，填好 `base_url` 与 `models`，再通过 `model_mapping` 校正上游模型名即可走通所有路由与计费。`updateChannelBalance` 在 type 为 `openai` / `custom` 时同样使用该 `base_url` 拉取余额。

## 渠道类型映射 / Channel type map

前端导出 `CHANNEL_TYPE_MAP` 给新增渠道弹窗使用：

```js
{
  openai: 1, claude: 2, azure: 3, gemini: 4,
  baidu: 5, aliyun: 6, tencent: 7, xunfei: 8,
  zhipu: 9, deepseek: 10, midjourney: 11
}
```

后端 `relay.GetAdaptorByChannel(channel.Type)` 通过该数值路由到 `relay/adaptor/provider/<slug>` 包；OpenAI 兼容中转复用 `openai` adaptor 并覆盖 `base_url`。

## 实现位置 / Implementation Pointers

| 关注点 | 位置 |
|---|---|
| Provider 常量 | `web/default-pro/src/constants/providers.js` |
| 后端 adaptor | `relay/adaptor/provider/<slug>/` |
| Adaptor 注册 | `relay/adaptor/openai` 等 |
