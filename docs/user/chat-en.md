---
title: Chat Playground
description: "Built-in chat playground for multi-model comparison."
category: user
order: 6
---

# Chat Playground

> Built-in chat playground for multi-model comparison.
> 内置的对话调试与多模型对比。

Entry: `/chat` (`web/default-pro/src/views/chat/Chat.vue`). The Chat Playground is an **embedded iframe** whose URL is configured via the system option `ChatLink`:

入口：`/chat`（`web/default-pro/src/views/chat/Chat.vue`）。Chat Playground 是一个 **嵌入式 iframe**，由系统选项 `ChatLink` 配置指向：

```text
ChatLink (env var CHAT_LINK) ─► iframe src
```

## Use Cases / 适用场景

- **Multi-model comparison**: run GPT-4o / Claude / DeepSeek / Qwen side-by-side in one UI.
  多模型对比：在一个 UI 里同时跑 GPT-4o / Claude / DeepSeek / Qwen 比拼回复。
- **Prompt iteration**: iterate long prompts without re-running CLIs.
  调试提示词：长 prompt 反复试错，不用每次都改 CLI。
- **Demo / training**: showcase LLM behavior to the team / customers.
  教学 / Demo：给团队 / 客户演示 LLM 行为。

## Configure / 配置

Env var:

环境变量：

```bash
CHAT_LINK=https://chat.example.com/
```

or / 或

```yaml
CHAT_LINK: "https://chat-oneapi.lobehub.com/"
```

The frontend reads `status.chat_link` and renders it directly:

前端获取 `status.chat_link` 后直接渲染：

```js
chatLink.value = statusStore.status?.chat_link || ''
```

When empty, `/chat` shows an empty state prompting the admin to set `ChatLink`.

留空则 `/chat` 页显示空状态（提示管理员配置 `ChatLink`）。

## Auth / 鉴权

Recommended pairings with third-party chat UIs:

推荐配合第三方 Chat UI：

| UI | Recommended setup / 推荐配置 | Notes / 备注 |
| --- | --- | --- |
| **Lobe Chat** | `CHAT_LINK=https://chat-oneapi.lobehub.com/` | Use One API Pro as the model gateway directly |
| **NextChat (ChatGPT-Next-Web)** | Self-host NextChat; set `API_BASE_URL` to One API Pro | Mind the CORS for `/v1/*` |
| **Open WebUI** | same as above | Streaming, tools, multimodal all supported |

> Third-party UIs use One API Pro's `sk-` key (not the Access Token); fill the "Custom Base URL" field with `http://<your-host>:3000`.
> 第三方 UI 调用的是 One API Pro 的 `sk-` Key（不是 Access Token），需要在它们的「自定义 Base URL」处填 `http://<your-host>:3000`。

## Troubleshooting / 排错

| Symptom / 现象 | Check / 排查 |
| --- | --- |
| `/chat` empty | Is `CHAT_LINK` set? Is `/api/status` returning `chat_link`? |
| iframe loaded but no models | The third-party UI's `API Base URL` / `API Key` pointing to One API Pro? |
| iframe CORS error | Add `Access-Control-Allow-Origin: <origin>` in the reverse proxy, or co-locate the third-party UI and One API Pro under the same origin |
| 401 from the upstream | Create a "never expire + unlimited models" API Key in One API Pro and paste it into the third-party UI |

## Self-hosted Third-party UI / 自托管第三方 UI

If `ChatLink` points to a self-hosted Lobe / NextChat / Open WebUI:

如果 `ChatLink` 指向自托管的 Lobe / NextChat / Open WebUI：

1. Give One API Pro a stable domain (avoid bare-IP + reverse-proxy 443).
   给 One API Pro 起一个稳定的域名（避免 IP 直连 + 反代 443）。
2. In the third-party UI's "Model Provider" choose OpenAI-compatible, set Base URL to `https://<one-api-domain>` and API Key to `sk-…`.
   在第三方 UI 「模型提供方」处选择 OpenAI 兼容，Base URL 填 `https://<one-api-domain>`，API Key 填 `sk-…`。
3. In One API Pro's "Tokens", allow the models the third-party UI calls (skip if unlimited).
   在 One API Pro 的「令牌」处允许第三方 UI 调用所需的模型（不限模型可省去此步）。

## Why no in-house chat UI / 不实现内嵌聊天 UI 的原因

One API Pro's frontend stays focused on gateway / billing / admin; the chat UI is delegated to community-maintained projects that move faster and support more protocols. This split — "gateway + pluggable UI" — is also the main product-positioning difference between One API Pro and upstream One-API.

保持 One API Pro 的前端「重在网关 / 计费 / 管理」三块，聊天 UI 由社区更专业的项目承担（更新更快、协议支持更广）。这种「网关 + 可插拔 UI」的拆分也是 One API Pro 与上游 One-API 在产品定位上的主要差异之一。

Next: [Lobe Chat](https://github.com/lobehub/lobe-chat) · [Add a Provider](/en/contribute/add-provider).