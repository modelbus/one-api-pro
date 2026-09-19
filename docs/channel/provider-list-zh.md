---
title: Provider 全清单
description: One API Pro 已适配的 LLM Provider 列表与适配说明。
category: channel
order: 6
---

# Provider 全清单

> 这里列出全部内置适配的 Provider 类型。新建渠道时从下拉里选对应类型即可。
>
> 如果你的上游是 **OpenAI 兼容** 的中转站，统一选 `OpenAI`，然后填 `base_url` + `models` 即可。

## 海外

| Provider 类型 | 适用 | 备注 |
|---|---|---|
| `openai` | OpenAI / 任意 OpenAI 兼容中转 | 默认类型；填 `base_url` 即可 |
| `anthropic` | Anthropic Claude | Messages API |
| `azure` | Azure OpenAI | 需填部署名 / API Key |
| `gemini` | Google Gemini | 原生协议 |
| `openrouter` | OpenRouter（聚合多个模型） | 走 OpenAI 协议 |
| `groq` | Groq | 走 OpenAI 协议 |
| `mistral` | Mistral AI | 走 OpenAI 协议 |
| `cohere` | Cohere | 走 OpenAI 协议 |
| `togetherai` | Together AI | 走 OpenAI 协议 |
| `replicate` | Replicate | 自定义协议 |
| `aws` | AWS Bedrock | 自定义协议 |
| `vertexai` | Google Vertex AI | 需 GCP 凭证 |
| `palm` | Google PaLM（旧） | 即将停用 |
| `xai` | x.AI (Grok) | 走 OpenAI 协议 |

## 国内

| Provider 类型 | 适用 | 备注 |
|---|---|---|
| `deepseek` | DeepSeek | 走 OpenAI 协议 |
| `moonshot` | 月之暗面 Kimi | 走 OpenAI 协议 |
| `zhipu` | 智谱 GLM | 自定义协议 |
| `qwen` / `ali` / `alibailian` | 阿里通义千问 / 百炼 | 多端点适配 |
| `doubao` | 字节豆包 | 自定义协议 |
| `baichuan` | 百川 | 自定义协议 |
| `baidu` / `baiduv2` | 文心一言 | 多版本适配 |
| `tencent` | 腾讯混元 | 自定义协议 |
| `lingyiwanwu` | 智谱零一万物 | 自定义协议 |
| `stepfun` | 阶跃星辰 | 自定义协议 |
| `minimax` | MiniMax | 自定义协议 |
| `novita` | Novita AI（GPU 推理聚合） | 自定义协议 |

## 本地 / 自托管

| Provider 类型 | 适用 | 备注 |
|---|---|---|
| `ollama` | Ollama（本地推理） | 走 OpenAI 协议 |
| `proxy` | 通用 HTTP 代理 | 当作纯透传 |

## 通用 OpenAI 兼容中转站

任何走 OpenAI 兼容协议的第三方聚合 / 中转 / 转发服务，都选 `openai` 类型，配置项：

- `Base URL`：服务商提供的接入地址
- `API Key`：服务商给的密钥
- `Models`：服务商支持的模型列表
- `Model Mapping`（可选）：把官方模型名映射成服务商内部的模型名

例如：

```json
{ "gpt-4o": "openai-gpt-4o-vip", "claude-sonnet-4": "anthropic-claude-3.5" }
```

## 怎么确认我的 Provider 在不在

- 直接看 [新增渠道] 的 Provider 下拉列表
- 也可以用 [渠道连通性测试](./channel-test) 实测

## 适配新 Provider

如果你的 Provider 不在清单里，且协议非 OpenAI 兼容，需要做适配：

- 写一个 `relay/adaptor/provider/<name>/` 子包
- 实现 `Adaptor` 接口 + `init()` 注册
- 在 [新增渠道] 的下拉里自动出现

详见 [新增 Provider](/contribute/add-provider)。

## 相关文档

- [渠道概览](./overview)
- [新增渠道](./add-channel)
- [渠道路由策略](./channel-routing)