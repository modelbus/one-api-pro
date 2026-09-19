---
title: Access Token
description: 个人中心生成的 Access Token 用来调用管理后台 API。
category: user
order: 4
---

# Access Token

> 浏览器登录失效了，但你的脚本 / CI 还想调管理后台 API？用 Access Token。

## Access Token 是什么

个人中心里生成的一串 UUID。它让你不用登录网页，也能调用 One API Pro 的 `/api/*` 管理接口（比如查自己、改资料、查日志等）。

适合这种场景：

- 写脚本 / 自动化任务
- 在 CI / 服务器上调用
- 浏览器 Cookie 不可用的场景

## 和 API Key 的区别

| 项 | Access Token | API Key（`sk-…`） |
|---|---|---|
| 是什么 | UUID 字符串 | `sk-` 开头的随机串 |
| 用来调哪 | `/api/*` 管理接口 | `/v1/*` OpenAI 兼容接口 |
| 鉴权头 | `Authorization: <uuid>` | `Authorization: Bearer sk-…` |
| 谁生成 | 你（个人中心） | 你（令牌页面） |
| 失效 | 重新生成后旧 Token 立即失效 | 删除旧 Key 后失效 |

简而言之：**Access Token 管理 One API Pro 自己**，**API Key 让别人用 One API Pro 调模型**。

## 怎么生成

**方式一：网页**

1. 登录后台
2. 右上角头像 → **个人中心**
3. 找到 **Access Token** 一栏
4. 点「生成」
5. **立刻复制保存** —— 刷新页面就再也看不到了

**方式二：用现有 Cookie 调一次 API**

如果你已经在别处登录了，可以直接：

```bash
curl http://localhost:3000/api/user/token -b cookies.txt
```

返回值就是新生成的 UUID。

## 怎么用

把 UUID 放到 `Authorization` 请求头里（**没有** `Bearer` 前缀）：

```bash
curl http://localhost:3000/api/user/self \
  -H "Authorization: <your_access_token>"
```

请求成功后会返回当前登录用户的信息。

## 重新生成（撤销旧 Token）

回到个人中心 → Access Token → 再次点「生成」。会得到新 UUID，**旧 UUID 立即失效**。

适用于：怀疑泄露、紧急撤销、定期轮换。

## 安全建议

| 建议 | 原因 |
|---|---|
| 像密码一样保管 | 拿到它的人能调 `/api/*` 全部权限 |
| 别提交到仓库 | 一旦提交，git 历史里永久可见 |
| 不同环境用不同的 Token | 减少爆炸半径 |
| 怀疑泄露就立即重生成 | 重生成后旧 Token 立即失效 |
| CI 里用专门的 Token | 别用你个人登录那个 |

## 常见问题

- **「生成」按钮点不动**：刷新页面后再试；或检查是否被弹窗拦截
- **重生成后旧 Token 还能用**：系统有几秒缓存，等几秒再试
- **能不能给别人用**：不行。Access Token 和你的账号一一对应，给别人 = 把账号给他

## 相关

- [个人资料](/user/profile) — 改密码、邮箱
- [我的订单](/user/orders)
- [API 文档](/api/README)