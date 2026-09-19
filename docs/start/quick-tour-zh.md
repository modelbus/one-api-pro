---
title: 5 分钟快速体验
description: "从安装到调用第一个模型的端到端流程。"
category: start
order: 3
---

# 5 分钟快速体验

> 从安装到调用第一个模型的端到端流程。

本指南假设你已经有了一台机器（Linux / macOS / Windows）和 Docker。完整说明见 [安装与更新](/zh/install/requirements)。

## 第 1 步：启动

```bash
docker run -d --name one-api-pro \
  --restart always \
  -p 3000:3000 \
  -v $(pwd)/data:/app/data \
  -e TZ=Asia/Shanghai \
  ghcr.io/modelbus/one-api-pro:latest
```

启动后访问 `http://localhost:3000`，看到登录页即成功。默认 root 账号：`root` / `123456`，**首次登录后立即改密**。

## 第 2 步：登录与创建 Access Token

1. 用 root 登录后进入「个人中心」→「Access Token」。
2. 点击「生成」，复制得到的 UUID（用户级 Token，用于调用 `/api/*`）。
3. 也可以直接用 Cookie Session 调用管理接口。

```bash
# 登录拿 Cookie
curl -X POST http://localhost:3000/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"username":"root","password":"123456"}' \
  -c cookies.txt

# 自查
curl http://localhost:3000/api/user/self -b cookies.txt
```

> 登录页右上角可切换中 / English（`localStorage.lang` 持久化）。

## 第 3 步：新建渠道

进入「渠道」→「新建渠道」，填写：

1. **类型**：选 OpenAI（或 DeepSeek、Gemini、Anthropic 等）；
2. **名称 / Base URL**：默认填好，可改；
3. **密钥**：上游 API Key；
4. **模型**：留空 = 支持该 Provider 全部模型（来自 `relay/adaptor/<provider>/constants.go`）。

新建完成后点「测试」，看到 200 即渠道可用。

## 第 4 步：创建 API Key

进入「令牌」→「新建令牌」：

- **名称**：任意可识别字符串；
- **额度**：可选；填 0 = 不限；
- **过期时间**：可选；填 0 / 勾选「永不过期」即不过期；
- **可用模型**：留空 = 不限。

点复制，**只显示这一次**。该 Key 以 `sk-` 前缀，可用于调用 `/v1/*` 兼容接口。

## 第 5 步：调用第一个模型

```bash
curl http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer sk-xxxxxxxxxxxxxxxx" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o-mini",
    "messages": [{"role":"user","content":"Hello from One API Pro!"}]
  }'
```

返回 JSON 即成功。流式（`"stream": true`）、工具调用（`tools`）、多模态（`image_url`）均按 OpenAI 协议透传。

## 验证

进入「日志」或运行 `curl http://localhost:3000/api/log/self -b cookies.txt` 查看本次调用：

- `prompt_tokens` / `completion_tokens` / `total_tokens` 由 tiktoken 内嵌 BPE 精确计数（见 [配置项：`TIKTOKEN_CACHE_DIR`](/zh/install/config#tokenizer-缓存--tokenizer-cache)）。
- `quota` 字段按 [定价](/zh/pricing/model-price) 与 [分组折扣](/zh/pricing/group-price) 计算。
- 用户级 `quota` / `token.remain_quota` 同步扣减，Redis 缓存通过 `IncreaseUserQuota` 一致性更新。

Inspect via `Logs` page or `curl /api/log/self -b cookies.txt`. Token counts are accurate thanks to embedded BPE; quota is computed from the price table and group ratio. Redis cache is kept consistent via `IncreaseUserQuota`.

下一步 / Next: [架构总览](/zh/start/architecture) · [套餐订阅](/zh/subscription/overview) · [充值](/zh/pricing/topup-settings)。
