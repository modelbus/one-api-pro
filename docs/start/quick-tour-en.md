---
title: Quick Tour
description: "End-to-end flow from install to your first model call."
category: start
order: 3
---

# Quick Tour

> End-to-end flow from install to your first model call.
> 从安装到调用第一个模型的端到端流程。

This guide assumes a Docker-ready machine (Linux / macOS / Windows). See [Requirements](/en/install/requirements) for the full setup.

本指南假设你已经有了一台机器（Linux / macOS / Windows）和 Docker。完整说明见 [Install & Upgrade](/en/install/requirements)。

## Step 1 — Launch / 第 1 步：启动

```bash
docker run -d --name one-api-pro \
  --restart always \
  -p 3000:3000 \
  -v $(pwd)/data:/app/data \
  -e TZ=Asia/Shanghai \
  ghcr.io/modelbus/one-api-pro:latest
```

After start, open `http://localhost:3000`. Default root: `root` / `123456` — change it immediately.

启动后访问 `http://localhost:3000`，看到登录页即成功。默认 root 账号：`root` / `123456`，**首次登录后立即改密**。

## Step 2 — Sign in & create an Access Token / 第 2 步：登录与创建 Access Token

1. Sign in as root, then go to "Personal Center" → "Access Token".
2. Click "Generate", copy the UUID (user-level token, used for `/api/*`).
3. Or just use the Cookie Session for admin endpoints.

```bash
# Login (cookie)
curl -X POST http://localhost:3000/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"username":"root","password":"123456"}' \
  -c cookies.txt

# Self
curl http://localhost:3000/api/user/self -b cookies.txt
```

> The language toggle lives in the top-right of auth pages; the choice is persisted in `localStorage.lang`.
> 登录页右上角可切换中 / English（`localStorage.lang` 持久化）。

## Step 3 — Add a Channel / 第 3 步：新建渠道

Navigate to "Channels" → "New Channel":

1. **Type**: OpenAI (or DeepSeek, Gemini, Anthropic, …).
2. **Name / Base URL**: prefilled; editable.
3. **Key**: upstream API key.
4. **Models**: empty = all models supported by that provider (from `relay/adaptor/<provider>/constants.go`).

After saving, click "Test"; a `200` means the channel is reachable.

新建完成后点「测试」，看到 200 即渠道可用。

## Step 4 — Create an API Key / 第 4 步：创建 API Key

Navigate to "Tokens" → "New Token":

- **Name**: any recognizable string.
- **Quota**: optional; `0` = unlimited.
- **Expires at**: optional; `0` / "Never expire" means never.
- **Models**: empty = unlimited.

Click copy — the key is shown **only once**. It is prefixed with `sk-` and works against any `/v1/*` compatible endpoint.

点复制，**只显示这一次**。该 Key 以 `sk-` 前缀，可用于调用 `/v1/*` 兼容接口。

## Step 5 — Call your first model / 第 5 步：调用第一个模型

```bash
curl http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer sk-xxxxxxxxxxxxxxxx" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o-mini",
    "messages": [{"role":"user","content":"Hello from One API Pro!"}]
  }'
```

A JSON reply means success. Streaming (`"stream": true`), tool calls (`tools`), multimodal (`image_url`) all pass through per the OpenAI protocol.

返回 JSON 即成功。流式（`"stream": true`）、工具调用（`tools`）、多模态（`image_url`）均按 OpenAI 协议透传。

## Verify / 验证

From the "Logs" page or `curl /api/log/self -b cookies.txt`:

- `prompt_tokens` / `completion_tokens` / `total_tokens` are exact counts from embedded tiktoken BPE — see [`TIKTOKEN_CACHE_DIR`](/en/install/config#tokenizer-cache--tokenizer-cache).
- `quota` is computed from the [price table](/en/pricing/model-price) and the [group ratio](/en/pricing/group-price).
- User-level `quota` and `token.remain_quota` are decremented atomically; the Redis cache is refreshed via `IncreaseUserQuota` to stay consistent.

进入「日志」或运行 `curl http://localhost:3000/api/log/self -b cookies.txt` 查看本次调用。

Next: [Architecture](/en/start/architecture) · [Subscriptions](/en/subscription/overview) · [Top-up](/en/pricing/topup).