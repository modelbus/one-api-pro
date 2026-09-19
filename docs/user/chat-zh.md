---
title: Chat Playground
description: 在 One API Pro 里直接对话，调试模型与提示词。
category: user
order: 6
---

# Chat Playground

> 想直接在 One API Pro 里对话、对比几个模型？Chat Playground 就是这个。

## 这是什么

One API Pro 自带的网页聊天界面，地址是 `/chat`。你可以：

- 在同一个对话里切换不同模型，看谁的回复更好
- 调提示词不用每次改 CLI
- 给团队 / 客户演示 LLM 行为

## 怎么开

不需要额外配置 —— 登录后访问 `/chat` 就行。如果显示空白页，说明管理员没设置 [ChatLink](../misc/system-settings)（默认是空的）。

## 用法

跟用 ChatGPT / Claude 一样：

1. 选择模型（下拉框列出了你能用的所有模型）
2. 输入提示词
3. 看回复

提示：可以同时开两个 tab，分别选不同模型，对比同一个 prompt 的回复。

## 它和「用第三方 Chat UI」是什么关系

Chat Playground 只是 One API Pro 内置的一个简易聊天页。如果你需要更完整的聊天体验（多会话、文件上传、插件等），建议用第三方 Chat UI，通过 [API Key](./access-token) 或 [令牌](../api/token) 接 One API Pro。

详见下方的「接入第三方 Chat UI」一节。

## 接入第三方 Chat UI

更推荐的做法：用 Lobe Chat / NextChat / Open WebUI 等社区 UI，通过 One API Pro 的 `sk-` Key 接入。

| 第三方 UI | 推荐配置 |
|---|---|
| **Lobe Chat** | 在「模型服务」选 OpenAI 兼容；Base URL 填 `http://<your-host>:3000`；API Key 填 `sk-…` |
| **NextChat** | 同上 |
| **Open WebUI** | 同上；支持流式、tools、多模态 |

具体步骤：

1. 在 One API Pro 个人中心 → 令牌 → 新建一个 API Key（建议「永不过期 + 不限模型」）
2. 在第三方 UI 设置里选 OpenAI 兼容，Base URL 指向 One API Pro（如 `http://localhost:3000` 或你的域名）
3. API Key 填上一步新建的 `sk-…`
4. 在第三方 UI 模型列表里应能看到所有已配置 [模型定价](../schema/model-price) 的模型

## 常见问题

- **`/chat` 页面空白**：管理员没设置 `ChatLink`，联系管理员在 [系统设置](../misc/system-settings) 配一下，或者直接用第三方 Chat UI
- **调用返回 401**：新建一个「永不过期 + 不限模型」的 API Key 试试
- **第三方 UI 显示「无模型」**：检查它的 Base URL 是否正确指向了 One API Pro，且 API Key 在 One API Pro 这边有效
- **CORS 报错**：把第三方 UI 和 One API Pro 用同源部署（同一域名反向代理），或反代处加 `Access-Control-Allow-Origin`

## 相关

- [Access Token](./access-token)
- [API Key（令牌）](../api/token)
- [模型定价](../schema/model-price)