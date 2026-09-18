---
title: OpenAI Compatible API
description: "/v1/chat/completions and other compatible endpoints."
category: api
order: 20
---
## 12. OpenAI Compatible API (v1)

The endpoints below are compatible with the OpenAI API format. They use Bearer Token auth.

### 12.1 List Models

**Endpoint:** `GET /v1/models`

**Auth:** Bearer Token

**Response:**

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


### 12.2 Get Model Detail

**Endpoint:** `GET /v1/models/:model`

**Auth:** Bearer Token


### 12.3 Chat Completions

**Endpoint:** `POST /v1/chat/completions`

**Auth:** Bearer Token

**Request body:**

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

**Endpoint:** `POST /v1/completions`


### 12.5 Embeddings

**Endpoint:** `POST /v1/embeddings`

**Endpoint:** `POST /v1/engines/:model/embeddings`


### 12.6 Image Generation

**Endpoint:** `POST /v1/images/generations`


### 12.7 Audio

**Endpoint:** `POST /v1/audio/transcriptions`

**Endpoint:** `POST /v1/audio/translations`

**Endpoint:** `POST /v1/audio/speech`


### 12.8 Moderations

**Endpoint:** `POST /v1/moderations`


### 12.9 Billing Query (OpenAI-compatible)

**Endpoint:** `GET /v1/dashboard/billing/subscription`

**Endpoint:** `GET /v1/dashboard/billing/usage?start_date=2024-01-01&end_date=2024-12-31`


### 12.10 Proxy Forwarding

**Endpoint:** `ANY /v1/oneapi/proxy/:channelid/*target`

**Description:** Proxies directly to the specified channel.
