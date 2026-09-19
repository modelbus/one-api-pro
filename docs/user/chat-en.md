---
title: Chat Playground
description: Chat with models directly inside One API Pro.
category: user
order: 6
---

# Chat Playground

> Want to chat with models directly inside One API Pro to compare them? That's what this is.

## What it is

A built-in chat page at `/chat`. You can:

- Switch models inside the same conversation to compare answers
- Iterate on prompts without re-running CLI commands
- Demo LLM behavior to teammates / customers

## How to open

No extra setup — just visit `/chat` after logging in. If the page is blank, it means the admin hasn't set [ChatLink](../misc/system-settings) (it's empty by default).

## Usage

Same as ChatGPT / Claude:

1. Pick a model (dropdown lists every model you can use)
2. Type your prompt
3. See the response

Tip: open two tabs with different models for the same prompt to compare side-by-side.

## vs. third-party Chat UIs

The Chat Playground is a simple built-in page. For a richer experience (multiple conversations, file uploads, plugins, etc.), use a third-party Chat UI connected via your [API Key](../api/token).

See "Third-party Chat UI" below.

## Third-party Chat UI

The recommended path: use Lobe Chat / NextChat / Open WebUI and connect via One API Pro's `sk-` key.

| Third-party UI | How to configure |
|---|---|
| **Lobe Chat** | Pick "OpenAI compatible"; Base URL = `http://<your-host>:3000`; API Key = `sk-…` |
| **NextChat** | Same |
| **Open WebUI** | Same; supports streaming, tools, multimodal |

Steps:

1. Personal Center → Tokens → create an API Key (recommended: "never expires + unrestricted models")
2. In the third-party UI, choose OpenAI-compatible; Base URL = your One API Pro (`http://localhost:3000` or your domain); API Key = the `sk-…`
3. Models configured in [Model Price](../schema/model-price) will appear in the UI

## FAQ

- **`/chat` is blank**: admin hasn't set `ChatLink`. Ask the admin to configure [System Settings](../misc/system-settings), or just use a third-party UI.
- **Calls return 401**: create a "never expires + unrestricted" API Key and try again.
- **Third-party UI shows "no models"**: check its Base URL points to One API Pro and the API Key is valid.
- **CORS error**: deploy the UI and One API Pro on the same origin (one reverse proxy), or add `Access-Control-Allow-Origin` in the reverse proxy.

## Related

- [Access Token](./access-token)
- [API Key](../api/token)
- [Model Price](../schema/model-price)