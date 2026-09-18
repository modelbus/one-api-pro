---
title: Misc Public API
description: "Group list, system options and other public endpoints."
category: api
order: 22
---
# Misc Public API



## 3. Group List (Group)

### 3.1 Get All Group Names

**Endpoint:** `GET /api/group/`

**Auth:** Admin

**Request parameters:** None

**Response:**

```json
{
  "success": true,
  "message": "",
  "data": ["default", "vip", "svip"]
}
```

**Description:** Returns every group name defined in the `group_prices` table.




## 4. System Options (Option)

### 4.1 Get All System Options

**Endpoint:** `GET /api/option/`

**Auth:** Root

**Request parameters:** None

**Response:**

```json
{
  "success": true,
  "message": "",
  "data": [
    { "key": "QuotaForNewUser", "value": "1000000" },
    { "key": "QuotaPerUnit", "value": "500000" },
    { "key": "TopUpLink", "value": "" },
    ...
  ]
}
```

**Description:** Sensitive options (tokens, SMTP passwords, etc.) are filtered or masked.


### 4.2 Update a System Option

**Endpoint:** `PUT /api/option/`

**Auth:** Root

**Request body:**

```json
{
  "key": "QuotaPerUnit",
  "value": "500000"
}
```

**Request fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| key | string | yes | Option name |
| value | string | yes | Option value |

**Available options:**

| Option name | Description |
|-------------|-------------|
| QuotaForNewUser | Quota granted to new users |
| QuotaForInviter | Reward quota for the inviter |
| QuotaForInvitee | Reward quota for the invitee |
| QuotaRemindThreshold | Quota reminder threshold |
| PreConsumedQuota | Pre-consumed quota |
| TopUpLink | Top-up link |
| ChatLink | Chat link |
| QuotaPerUnit | Quota per unit (how much internal quota equals ¥1) |
| DisplayInCurrencyEnabled | Whether to display amounts as currency (true/false) |
| DisplayTokenStatEnabled | Whether to display token stats (true/false) |
| ApproximateTokenEnabled | Whether to use approximate token counting (true/false) |
| RetryTimes | Retry count |
| LogConsumeEnabled | Whether to enable consumption logs (true/false) |
| AutomaticDisableChannelEnabled | Auto-disable channels (true/false) |
| AutomaticEnableChannelEnabled | Auto-enable channels (true/false) |
| ChannelDisableThreshold | Channel disable threshold |
| ChannelDefaultCooldownSeconds | Default cooldown (seconds) |
| ChannelMaxCooldownSeconds | Max cooldown (seconds) |
| ChannelConcurrencyEnabled | Enable per-channel concurrency limit (true/false) |
| ChannelStickySessionEnabled | Enable sticky session (true/false) |
| ErrorNext | Error response strategy JSON |

**Response:**

```json
{
  "success": true,
  "message": ""
}
```




## 14. Other Public Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/status` | GET | Get system status (no auth required) |
| `/api/notice` | GET | Get system notice |
| `/api/about` | GET | Get about page content |
| `/api/home_page_content` | GET | Get home page content |
| `/api/verification` | GET | Send email verification code |
| `/api/reset_password` | GET | Send password reset email |
| `/api/user/reset` | POST | Reset password |
| `/api/oauth/github` | GET | GitHub OAuth callback |
| `/api/oauth/oidc` | GET | OIDC OAuth callback |
| `/api/oauth/lark` | GET | Lark OAuth callback |
| `/api/oauth/state` | GET | Generate OAuth state code |
| `/api/oauth/wechat` | GET | WeChat OAuth callback |
| `/api/oauth/wechat/bind` | GET | Bind WeChat account (requires login) |
| `/api/oauth/email/bind` | GET | Bind email (requires login) |
| `/api/user/logout` | GET | Logout |
| `/api/user/aff` | GET | Get affiliate code |
| `/api/user/subscription` | GET | Get subscription info |


## Appendix B: Permission Levels

| Level | Value | Description |
|-------|-------|-------------|
| Guest | 0 | Unauthenticated visitor |
| User | 1 | Regular user |
| Admin | 10 | Admin |
| Root | 100 | Super admin |


## Appendix C: Billing Types

| billing_type | Description |
|-------------|-------------|
| `token` | Per-token billing, priced in ¥ per million tokens |
| `per_request` | Per-request billing, priced in ¥ per request |

**Token billing formula:**

```
quota = ceil((inputPrice × inputTokens + outputPrice × completionTokens + cachedPrice × cachedTokens) / 1,000,000 × groupDiscount × QuotaPerUnit)
```

**Per-request billing formula:**

```
quota = ceil(perRequestPrice × sizeRatio × N × groupDiscount × QuotaPerUnit)
```

**Group discount matching rules:**

1. `GroupName + ModelName` exact match → use that discount.
2. `GroupName + ""` (empty) → use that group's default discount.
3. No match → discount is 1.0 (no discount).

**`QuotaPerUnit`:** Default 500,000, meaning 500,000 internal quota equals ¥1.


## Appendix D: Channel Type Reference

| Type | Adaptor |
|------|---------|
| 1 | OpenAI |
| 2 | Azure OpenAI |
| 3 | Custom channel |
| 4 | Claude (Anthropic) |
| 5 | Google Gemini |
| 6 | Tongyi Qianwen (Ali) |
| 7 | iFlytek Spark (SparkDesk) |
| 8 | Baidu ERNIE (Baidu) |
| 9 | ByteDance Doubao (Doubao) |
| 10 | MiniMax |
| 11 | DeepSeek |
| 12 | Cohere |
| 13 | 360 Brain |
| 14 | Ollama |
| 15 | Moonshot |
| 16 | Zhipu AI (GLM) |
| 17 | Baichuan |
| 18 | Yi (01.AI) |
| 19 | DeepSeek (dedicated adaptor) |
| 20-40+ | Other adaptors |