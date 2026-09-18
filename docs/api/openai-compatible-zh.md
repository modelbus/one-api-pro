---
title: OpenAI 兼容接口
description: "OpenAI 兼容接口：OpenAI 兼容接口 (v1)"
category: api
order: 20
---
## 12. OpenAI 兼容接口 (v1)

以下接口与 OpenAI API 格式兼容，使用 Bearer Token 认证。

### 12.1 模型列表

**接口：** `GET /v1/models`

**权限：** Bearer Token

**返回值：**

```json
{
  "object": "list",
  "data": [
    {
      "id": "gpt-4o",
      "object": "model",
      "created": 1718000000,
      "owned_by": "one-api-pro"
    }
  ]
}
```


### 12.2 获取模型详情

**接口：** `GET /v1/models/:model`

**权限：** Bearer Token


### 12.3 Chat Completions

**接口：** `POST /v1/chat/completions`

**权限：** Bearer Token

**请求体：**

```json
{
  "model": "gpt-4o",
  "messages": [
    {"role": "system", "content": "You are a helpful assistant."},
    {"role": "user", "content": "Hello"}
  ],
  "temperature": 0.7,
  "max_tokens": 4096,
  "stream": true
}
```


### 12.4 Text Completions

**接口：** `POST /v1/completions`


### 12.5 Embeddings

**接口：** `POST /v1/embeddings`

**接口：** `POST /v1/engines/:model/embeddings`


### 12.6 图片生成

**接口：** `POST /v1/images/generations`


### 12.7 音频转写

**接口：** `POST /v1/audio/transcriptions`

**接口：** `POST /v1/audio/translations`

**接口：** `POST /v1/audio/speech`


### 12.8 内容审核

**接口：** `POST /v1/moderations`


### 12.9 计费查询（OpenAI 兼容）

**接口：** `GET /v1/dashboard/billing/subscription`

**接口：** `GET /v1/dashboard/billing/usage?start_date=2024-01-01&end_date=2024-12-31`


### 12.10 代理转发

**接口：** `ANY /v1/oneapi/proxy/:channelid/*target`

**说明：** 直接代理到指定渠道。

