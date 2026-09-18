---
title: Chat Playground
description: "内置的对话调试与多模型对比。"
category: user
order: 6
---

# Chat Playground

> 内置的对话调试与多模型对比。
> Built-in chat playground for multi-model comparison.

入口：`/chat`（`web/default-pro/src/views/chat/Chat.vue`）。Chat Playground 是一个 **嵌入式 iframe**，由系统选项 `ChatLink` 配置指向：

Entry: `/chat` (`web/default-pro/src/views/chat/Chat.vue`). The Chat Playground is an **embedded iframe** whose URL is configured via the system option `ChatLink`:

```text
ChatLink (环境变量 CHAT_LINK) ─► iframe src
```

## 适用场景 / Use Cases

- **多模型对比**：在一个 UI 里同时跑 GPT-4o / Claude / DeepSeek / Qwen 比拼回复。
  Multi-model comparison: run GPT-4o / Claude / DeepSeek / Qwen side-by-side in one UI.
- **调试提示词**：长 prompt 反复试错，不用每次都改 CLI。
  Prompt iteration: iterate long prompts without re-running CLIs.
- **教学 / Demo**：给团队 / 客户演示 LLM 行为。
  Demo / training: showcase LLM behavior to the team / customers.

## 配置 / Configure

环境变量 / Env var：

```bash
CHAT_LINK=https://chat.example.com/
```

或 / or

```yaml
CHAT_LINK: "https://chat-oneapi.lobehub.com/"
```

前端获取 `status.chat_link` 后直接渲染：

The frontend reads `status.chat_link` and renders it directly:

```js
chatLink.value = statusStore.status?.chat_link || ''
```

留空则 `/chat` 页显示空状态（提示管理员配置 `ChatLink`）。

When empty, `/chat` shows an empty state prompting the admin to set `ChatLink`.

## 鉴权 / Auth

推荐配合第三方 Chat UI：

Recommended pairings with third-party chat UIs:

| UI | 推荐配置 / Recommended setup | 备注 / Notes |
| --- | --- | --- |
| **Lobe Chat** | `CHAT_LINK=https://chat-oneapi.lobehub.com/` | 直接用 One API Pro 当模型网关 |
| **NextChat (ChatGPT-Next-Web)** | 自部署 NextChat，`API_BASE_URL` 指向 One API Pro | 注意 `/v1/*` 跨域（CORS） |
| **Open WebUI** | 同上 | 支持流式、tools、多模态 |

> 第三方 UI 调用的是 One API Pro 的 `sk-` Key（不是 Access Token），需要在它们的「自定义 Base URL」处填 `http://<your-host>:3000`。
> Third-party UIs use One API Pro's `sk-` key (not the Access Token); fill the "Custom Base URL" field with `http://<your-host>:3000`.

## 排错 / Troubleshooting

| 现象 / Symptom | 排查 / Check |
| --- | --- |
| `/chat` 显示空状态 | `CHAT_LINK` 是否设置；`/api/status` 返回的 `chat_link` 是否非空 |
| iframe 加载但「无模型」 | 第三方 UI 侧的 `API Base URL` / `API Key` 是否正确指向 One API Pro |
| iframe 报 CORS | 反代处加 `Access-Control-Allow-Origin: <origin>`；或把三方 UI 与 One API Pro 同源部署 |
| 调用返回 401 | 在 One API Pro 个人中心新建「永不过期 + 不限模型」的 API Key 写入第三方 UI |

## 自托管第三方 UI / Self-hosted Third-party UI

如果 `ChatLink` 指向自托管的 Lobe / NextChat / Open WebUI：

If `ChatLink` points to a self-hosted Lobe / NextChat / Open WebUI:

1. 给 One API Pro 起一个稳定的域名（避免 IP 直连 + 反代 443）。
   Give One API Pro a stable domain (avoid bare-IP + reverse-proxy 443).
2. 在第三方 UI 「模型提供方」处选择 OpenAI 兼容，Base URL 填 `https://<one-api-domain>`，API Key 填 `sk-…`。
   In the third-party UI's "Model Provider" choose OpenAI-compatible, set Base URL to `https://<one-api-domain>` and API Key to `sk-…`.
3. 在 One API Pro 的「令牌」处允许第三方 UI 调用所需的模型（不限模型可省去此步）。
   In One API Pro's "Tokens", allow the models the third-party UI calls (skip if unlimited).

## 不实现内嵌聊天 UI 的原因 / Why no in-house chat UI

保持 One API Pro 的前端「重在网关 / 计费 / 管理」三块，聊天 UI 由社区更专业的项目承担（更新更快、协议支持更广）。这种「网关 + 可插拔 UI」的拆分也是 One API Pro 与上游 One-API 在产品定位上的主要差异之一。

One API Pro's frontend stays focused on gateway / billing / admin; the chat UI is delegated to community-maintained projects that move faster and support more protocols. This split — "gateway + pluggable UI" — is also the main product-positioning difference between One API Pro and upstream One-API.

下一步 / Next: [接入第三方 Chat UI](https://github.com/lobehub/lobe-chat) · [新增 Provider](/zh/contribute/add-provider)。