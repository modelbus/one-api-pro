---
title: 其他公共 API
description: 分组列表、系统选项与其他公共端点。
category: api
order: 22
---
# 其他公共 API



## 3. 分组列表 (Group)

### 3.1 获取所有分组名称

**接口：** `GET /api/group/`

**权限：** Admin

**请求参数：** 无

**返回值：**

```json
{
  "success": true,
  "message": "",
  "data": ["default", "vip", "svip"]
}
```

**说明：** 返回所有在 `group_prices` 表中定义的分组名称列表。

---


## 4. 系统选项 (Option)

### 4.1 获取所有系统选项

**接口：** `GET /api/option/`

**权限：** Root

**请求参数：** 无

**返回值：**

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

**说明：** 敏感选项（如 Token、SMTP 密码等）的值会被过滤或脱敏。

---

### 4.2 更新系统选项

**接口：** `PUT /api/option/`

**权限：** Root

**请求体：**

```json
{
  "key": "QuotaPerUnit",
  "value": "500000"
}
```

**请求字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| key | string | 是 | 选项名称 |
| value | string | 是 | 选项值 |

**可用选项列表：**

| 选项名 | 说明 |
|--------|------|
| QuotaForNewUser | 新用户赠送额度 |
| QuotaForInviter | 邀请人奖励额度 |
| QuotaForInvitee | 被邀请人奖励额度 |
| QuotaRemindThreshold | 额度提醒阈值 |
| PreConsumedQuota | 预消耗额度 |
| TopUpLink | 充值链接 |
| ChatLink | 聊天链接 |
| QuotaPerUnit | 单位额度（1元=多少内部额度） |
| DisplayInCurrencyEnabled | 是否以货币显示（true/false） |
| DisplayTokenStatEnabled | 是否显示Token统计（true/false） |
| ApproximateTokenEnabled | 是否使用近似Token计算（true/false） |
| RetryTimes | 重试次数 |
| LogConsumeEnabled | 是否启用消费日志（true/false） |
| AutomaticDisableChannelEnabled | 自动禁用渠道（true/false） |
| AutomaticEnableChannelEnabled | 自动启用渠道（true/false） |
| ChannelDisableThreshold | 渠道禁用阈值 |
| ChannelDefaultCooldownSeconds | 默认冷却时间（秒） |
| ChannelMaxCooldownSeconds | 最大冷却时间（秒） |
| ChannelConcurrencyEnabled | 启用渠道并发限制（true/false） |
| ChannelStickySessionEnabled | 启用粘性会话（true/false） |
| ErrorNext | 错误响应策略 JSON |

**返回值：**

```json
{
  "success": true,
  "message": ""
}
```

---


## 14. 其他公共接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/status` | GET | 获取系统状态（无需认证） |
| `/api/notice` | GET | 获取系统公告 |
| `/api/about` | GET | 获取关于页面内容 |
| `/api/home_page_content` | GET | 获取首页内容 |
| `/api/verification` | GET | 发送邮箱验证码 |
| `/api/reset_password` | GET | 发送密码重置邮件 |
| `/api/user/reset` | POST | 重置密码 |
| `/api/oauth/github` | GET | GitHub OAuth 回调 |
| `/api/oauth/oidc` | GET | OIDC OAuth 回调 |
| `/api/oauth/lark` | GET | 飞书 OAuth 回调 |
| `/api/oauth/state` | GET | 生成 OAuth 状态码 |
| `/api/oauth/wechat` | GET | 微信 OAuth 回调 |
| `/api/oauth/wechat/bind` | GET | 微信账号绑定（需登录） |
| `/api/oauth/email/bind` | GET | 邮箱绑定（需登录） |
| `/api/user/logout` | GET | 退出登录 |
| `/api/user/aff` | GET | 获取推广码 |
| `/api/user/subscription` | GET | 获取订阅信息 |

---

## 附录 B：权限等级说明

| 等级 | 值 | 说明 |
|------|------|------|
| Guest | 0 | 未登录用户 |
| User | 1 | 普通用户 |
| Admin | 10 | 管理员 |
| Root | 100 | 超级管理员 |

---

## 附录 C：计费类型说明

| billing_type | 说明 |
|-------------|------|
| `token` | 按 Token 计费，价格单位为 ¥/百万tokens |
| `per_request` | 按次计费，价格单位为 ¥/次 |

**Token 计费公式：**

```
quota = ceil((inputPrice × inputTokens + outputPrice × completionTokens + cachedPrice × cachedTokens) / 1,000,000 × groupDiscount × QuotaPerUnit)
```

**按次计费公式：**

```
quota = ceil(perRequestPrice × sizeRatio × N × groupDiscount × QuotaPerUnit)
```

**分组折扣匹配规则：**
1. `GroupName + ModelName` 精确匹配 → 使用该折扣
2. `GroupName + ""(空)` → 使用该分组默认折扣
3. 无匹配 → 折扣为 1.0（无折扣）

**QuotaPerUnit 说明：** 默认 500,000，表示 500,000 内部额度 = 1 元人民币。

---

## 附录 D：渠道类型对照表

| Type | 适配器 |
|------|--------|
| 1 | OpenAI |
| 2 | Azure OpenAI |
| 3 | 自定义渠道 |
| 4 | Claude (Anthropic) |
| 5 | Google Gemini |
| 6 | 通义千问 (Ali) |
| 7 | 讯飞星火 (SparkDesk) |
| 8 | 百度文心 (Baidu) |
| 9 | 字节豆包 (Doubao) |
| 10 | MiniMax |
| 11 | DeepSeek |
| 12 | Cohere |
| 13 | 360 智脑 |
| 14 | Ollama |
| 15 | 月之暗面 (Moonshot) |
| 16 | 智谱AI (GLM) |
| 17 | 百里 (Baichuan) |
| 18 | 零一万物 (Yi) |
| 19 | DeepSeek (独立适配器) |
| 20-40+ | 其他适配器 |

