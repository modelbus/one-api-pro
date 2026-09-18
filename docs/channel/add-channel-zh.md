---
title: 新增渠道
description: "新增渠道的字段含义、保存路径与多 Key 批量录入。"
category: channel
order: 2
---

# 新增渠道

> 新增渠道请求：POST `/api/channel/`，实现见 `controller/channel.go::AddChannel`。

## 接口

| 项 | 值 |
|---|---|
| Method | `POST` |
| Path | `/api/channel/` |
| Auth | Admin |
| Body | `Channel` JSON；`key` 支持 `\n` 分隔的多行批量录入 |

`AddChannel` 在 `controller/channel.go:80` 会把 `key` 字段按 `\n` 拆分：每行一个 key 在事务内生成一条渠道记录，然后由 `model.BatchInsertChannels` 调 `AddAbilities()` 把每个 `models` 字段里的模型写入 `abilities` 表。

## 字段详解

| 字段 | JSON 类型 | 说明 |
|---|---|---|
| `type` | `int` | Provider 类型枚举（与前端常量 `CHANNEL_TYPE_MAP` 对应：openai=1, claude=2, azure=3, gemini=4, baidu=5, aliyun=6, tencent=7, xunfei=8, zhipu=9, deepseek=10, midjourney=11 …） |
| `key` | `string` | 鉴权凭证；服务端从 `Authorization: Bearer <key>` 头读取。允许 `"\n"` 分隔一次提交多条 |
| `name` | `string` | 显示名（索引列），便于列表检索 |
| `base_url` | `*string` | 上游 API 根地址。留空时使用 Provider 默认地址；OpenAI 兼容中转必填 |
| `models` | `string` | 模型白名单，逗号分隔。测试时若请求的模型不在列表里，会自动回退到列表的第一个 |
| `group` | `string` | 允许使用该渠道的用户组，逗号分隔。`ContainsGroup` 做精确匹配；空 group 表示对所有用户可见 |
| `model_mapping` | `*string` | JSON 对象：`{ "源模型": "上游实际模型名" }`。请求时命中后会把入参模型名改写到上游实际名 |
| `system_prompt` | `*string` | 转发到上游前拼接在系统消息前的提示词（部分 relay 流程会使用） |
| `weight` | `*uint` | 加权轮询权重；当前路由以 `priority` 为主，`weight` 保留字段 |
| `priority` | `*int64` | 同 priority 的渠道在同优先级内随机；值越大越靠前 |
| `max_concurrency` | `*int` | 单节点 / 全集群的最大并发；`<=0` 表示不限。`ConcurrencyFilter` 启用 |
| `cooldown_seconds` | `int` | 单次上游错误后该渠道进入冷却的秒数（默认 60） |
| `rpm` | `*int` | 每分钟请求上限；`RPMFilter` 启用。`<=0` 表示不限 |
| `is_fallback` | `*bool` | 设为 true 后仅在所有正常渠道耗尽时由 fallback 路径选中 |
| `fallback_priority` | `*int64` | 同为 fallback 时按此值升序选中 |
| `config` | `string` | Provider 特定 JSON（`ChannelConfig`：region|
| `status` | `int` | 默认 1；详见 [渠道路由](./channel-routing) |

> 列表接口 `GetAllChannels`（默认 scope）会 `Omit("key")`，前端永远拿不到真实 key；只有 `GetChannel` 不带 `id` 参数或 `selectAll=true` 时才会回填。

## 多 Key 批量

把多个 key 用换行写在 `key` 字段里即可一次创建多条渠道。例如：

```

服务端逐条 `Insert` + `AddAbilities`，失败会回滚整批。

## 更新

`UpdateChannel` 解析原始 JSON 后只对 payload 里实际出现的 key 做 `Updates`，避免 `{id, status}` 这类部分更新把其他字段清空。详见 `controller/channel.go:151`。

## 前端操作指南

路径：`/channel` → 「新增渠道」按钮（`web/default-pro/src/views/channel/Channel.vue`）。

- 「类型」下拉展示前端 `CHANNEL_TYPE_MAP` 与 Provider 列表
- 「Base URL」对 OpenAI 兼容中转是必填，对官方地址是覆盖
- 「模型」从 `GET /api/model_price/options` 拉取已启用模型清单（AdminAuth）
- 「分组」多选，用户组在用户管理里维护
- 「Model Mapping」以键值对编辑，提交时序列化为 JSON 字符串
- 「最大并发 / RPM / 冷却秒数」`<=0` 表示不限

## 实现位置

| 关注点 | 位置 |
|---|---|
| CRUD | `controller/channel.go` |
| 模型定义 | `model/channel.go::Channel` |
| 能力同步 | `model/ability.go::AddAbilities` / `UpdateAbilities` |
| 部分更新安全 | `controller/channel.go::UpdateChannel` |

