---
title: 令牌（Token）管理
description: "用户级 API 令牌 CRUD、剩余额度、状态与子网白名单。"
category: api
order: 10
---

# 令牌（Token）管理

> 在 `tokens` 表上为每个用户管理 API 令牌（`sk-...`），是 OpenAI 兼容接口 `/v1/...` 的鉴权凭据。前端组件：`web/default-pro/src/views/token/Token.vue`。

> 该页面在普通用户视角是「我的令牌」；管理员可直接访问同一接口拉取当前用户的 token（不带 user_id 查询参数）；跨用户查询请走 `/api/order` 或 `GET /api/user/self` 等替代路径。

## 数据模型

`model.Token`（`tokens` 表，关键字段）：

| 字段 | 类型 | 说明 |
|---|---|---|
| `user_id` | `int` | 所属用户 |
| `name` | `varchar(50)` | 令牌名称 |
| `key` | `varchar` UNIQUE | 密钥本体（`sk-` 前缀在请求头中拼接） |
| `status` | `int` | `TokenStatusEnabled=1` / `Disabled=2` / `Expired=3` / `Exhausted=4` |
| `expired_time` | `int64` | unix 秒；`-1` 表示永不过期 |
| `remain_quota` | `int64` | 剩余额度（quota） |
| `unlimited_quota` | `bool` | 无限额度开关 |
| `models` | `text`（逗号分隔） | 允许使用的模型白名单；空字符串表示全部 |
| `subnet` | `text` | 客户端子网白名单（如 `10.0.0.0/8,192.168.0.0/16`） |
| `created_time` / `accessed_time` | `int64` | 创建 / 最近访问时间 |

## 接口一览

| Endpoint | Method | 鉴权 | 说明 |
|---|---|---|---|
| `/api/token/` | `GET` | User | 当前用户令牌列表（`?p=` `?order=`） |
| `/api/token/search?keyword=` | `GET` | User | 按 `name` 模糊搜索 |
| `/api/token/:id` | `GET` | User | 单条详情（必须属于当前用户） |
| `/api/token/` | `POST` | User | 创建令牌（自动用 `random.GenerateKey` 生成密钥） |
| `/api/token/` | `PUT` | User | 编辑（`?status_only=true` 仅改状态） |
| `/api/token/:id` | `DELETE` | User | 删除 |
| `/api/token/status` | `GET` | 内部（鉴权由 relay 中间件完成） | 返回 OpenAI 兼容的 `credit_summary` |
| `/api/token/self/token` | `GET` | User | **生成 / 重新生成用户级 AccessToken**（与 `token.key` 不同，是 dashboard 用的 access_token） |

实现：`controller/token.go`。

## 创建令牌

请求体：

```

- `name` 不能超过 30 字符；
- `subnet` 非空时会用 `network.IsValidSubnets` 校验 CIDR 列表格式；
- `key` 在后端由 `random.GenerateKey()` 生成，前端永远拿不到生成前的样子（创建接口的 response 是完整 `Token`）；
- 新建令牌状态默认 `TokenStatusEnabled=1`。

## 启用前置条件

`UpdateToken` 在 `status = Enabled` 时会二次校验：

- 如果原 token 已过期（`expired_time <= now`）：返回 "令牌已过期，无法启用，请先修改令牌过期时间，或者设置为永不过期"；
- 如果原 token 已耗尽（`remain_quota <= 0` 且非 `unlimited_quota`）：返回 "令牌可用额度已用尽..."。

## OpenAI 兼容

`GetTokenStatus` 由 relay 调用，返回：

```json
{
  "object": "credit_summary",
  "total_granted": <remain_quota>,
  "total_used": 0,
  "total_available": <remain_quota>,
  "expires_at": <expired_time * 1000>   // -1 时返回 0
}
```

> 注意 `total_used` 当前固定返回 `0`，因为 token 没有单独的 `used_quota` 字段；如需按 token 维度用量，需后续单独聚合。

## 配额消耗路径

调用 `/v1/chat/completions` 时由 `relay` 中间件读取 `token` 信息：

- `unlimited_quota=true` → 跳过扣减；
- `models` 不为空 → 仅当请求 model ∈ 白名单才放行；
- `subnet` 不为空 → 校验客户端 IP 命中；
- `expired_time > 0 AND expired_time <= now` → 视为 expired；
- `remain_quota` 不足 → 视为 exhausted。

## 前端操作指南

- 顶部 Base URL 卡片：`{{ baseUrl }}/v1`，复制即用；附带「使用指南」按钮。
- 「新增令牌」弹窗字段：`name`、`expired_time`（unix 时间或 `-1`）、`remain_quota`、`unlimited_quota` 开关、`models`（多选/逗号分隔）、`subnet`（CIDR）。
- 行内：
  - **查看完整 Key**：弹窗展示明文（创建后只能这次看到）；
  - **编辑**：更新除 `key` 外的所有字段；
  - **启用/禁用**：popconfirm；
  - **删除**：二次确认。
- 模型列点击「眼睛」图标弹窗展示 `models` 完整列表。

## 接口实现

| 关注点 | 位置 |
|---|---|
| CRUD handler | `controller/token.go` |
| AccessToken（用户级） | `controller/user.go::GenerateAccessToken` |
| OpenAI 兼容 credit_summary | `controller/token.go::GetTokenStatus` |
| 中间件（鉴权 / 配额） | `middleware/auth.go` / `relay/` |
| 路由 | `router/api.go` |
