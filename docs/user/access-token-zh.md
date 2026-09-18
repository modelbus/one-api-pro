---
title: Access Token
description: "管理用户的系统级 Access Token。"
category: user
order: 4
---

# Access Token

> 管理用户的系统级 Access Token。

Access Token 是 **UUID 格式** 的字符串，用作 `/api/*` 管理接口的鉴权头；当浏览器 Cookie 失效时（如脚本调用、CI），使用 Access Token 是最稳妥的方式。

## 与 API Key 的区别

| 维度| Access Token | API Key（`sk-…`） |
| --- | --- | --- |
| 格式| UUID 字符串 | `sk-` + 随机串 |
| 用途| `/api/*` 管理接口 | `/v1/*` OpenAI 兼容接口 |
| 鉴权方式| `Authorization: <uuid>` | `Authorization: Bearer sk-…` |
| 颁发方| 用户自生成（个人中心） | 用户在「令牌」页自创建 |
| 重生| 任何时候都可重新生成，旧 Token 立即失效 | 删除重建，旧 Key 立即失效 |

## 生成

UI：个人中心 → Access Token → 「生成」。

```bash
curl http://localhost:3000/api/user/token -b cookies.txt
# { success:true, data:"<uuid>" }
```

仅显示一次，必须立刻复制保存；刷新页面后再访问就拿不到了。

## 使用

```bash
# 调用管理接口
curl http://localhost:3000/api/user/self \
  -H "Authorization: <your_access_token>"

# 调 list 类的 admin API（要求 admin）
curl http://localhost:3000/api/user/list \
  -H "Authorization: <your_access_token>"
```

后端中间件鉴权优先级：Cookie Session 优先 → Access Token 兜底（详见 [API 总览 · 鉴权](/zh/api/README#附录-a鉴权机制)）。

## 重新生成

再次点「生成」会生成新的 UUID，旧的立即失效。这是一种「紧急撤销」手段。

## 配合 SDK

脚本 / CI 里常见的两种用法：

```bash
# 1. 临时导出环境变量
export OAP_TOKEN=$(curl -s http://localhost:3000/api/user/token -b cookies.txt | jq -r .data)

# 2. 直接放 Secret Manager
gh secret set OAP_TOKEN < token.txt
```

## 安全实践

| 建议| 说明|
| --- | --- |
| **像密码一样保护** | 一旦泄露，攻击者可以读取 / 修改你的 `/api/*` 资源 |
| **避免提交到仓库** | CI 临时文件、`.env` 都加进 `.gitignore` |
| **定期重生** | 推荐每 90 天重生一次；外泄时立即重生 |
| **不同环境用不同 Token** | 生产 / 预发 / 测试不要共用 |
| **配合限流** | `GLOBAL_API_RATE_LIMIT` 兜底防爆破 |

## 常见问题

- **「生成」按钮没反应**：检查 `status` 接口返回的 `github_oauth` 等字段；正常情况下所有用户都能生成。
- **重生成后旧 Token 仍可调用**：缓存层延迟几秒；建议立即刷新页面并清掉本地凭证。
- **能否给其他人用**：Access Token 与 `users.id` 一一对应；他人使用 = 越权。

下一步 / Next: [API Key（令牌）](/zh/api/token) · [我的订单](/zh/user/orders) · [Chat Playground](/zh/user/chat)。
