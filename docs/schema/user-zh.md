---
title: 用户
description: One API Pro 的账号体系：身份、额度、用户组、配额。
category: schema
order: 2
---

# 用户（User）

## 这是什么

`User` 是系统的唯一身份。每个登录、每次 API 调用、每条订单最终都关联到一个用户。

## 在哪里看到

- **管理后台 → 用户管理**：列表、详情、调整余额、启用/禁用
- **前台 → 个人中心**：当前登录账号的邮箱、用户名、用户组
- **调用日志**：每条 `/v1/*` 调用都记录了发起调用的 `user_id`

## 关键字段（运维视角）

下面是真正影响行为的字段，遇到问题时通常只动这些：

| 字段 | 含义 | 设置后影响 |
|---|---|---|
| `username` | 登录账号 | 修改后下次登录用新用户名 |
| `email` | 找回密码 / 通知 | 改邮箱需要重新走验证 |
| `password` | 登录凭证 | 重置后用户被踢下线 |
| `group` | `default` / `vip` / `svip` 等 | 影响 [分组折扣](/schema/group-price) 命中与可访问模型范围 |
| `status` | 启用 / 禁用 | 禁用后无法调用任何 API、无法登录 |
| `quota` | 余额（按内部单位计） | 扣调用费时优先扣这里；扣到 0 直接拒绝 |
| `used_quota` | 累计已扣 | 仅记录，不影响行为 |
| `role` | 普通用户 / 管理员 / Root | Root 才能进后台；普通用户只能看到自己 |
| `aff_code` / `inviter_id` | 邀请关系 | 注册时填写邀请码，双方获奖励额度 |

不推荐手动改的字段：注册时间、最后登录时间（由系统自动维护）。

## 相关数据结构

- 用户创建时，会自动创建一个 [Token](/schema/token)，用于调用管理后台 API
- 用户的所有 [Order](/schema/order) / [Topup](/schema/topup) / [Subscription](/schema/subscription) 都挂在 `user_id` 上
- 用户的实际扣费顺序：先用 `quota`，再用 [Token 配额](/schema/token)，再用 [Subscription 套餐额度](/schema/subscription)。具体规则见 [套餐与计费](/subscription/billing-rules)

## 相关功能页

- [注册与登录](/user/registration) — 如何注册、登录、找回密码
- [个人资料](/user/profile) — 用户能改什么
- [用户管理（管理员）](/user/user-management) — 管理员如何调整账户

## 相关 API

- `GET /api/user/self` — 查看当前用户信息
- `POST /api/user` — 管理员创建用户
- `PUT /api/user` — 管理员更新用户
- `POST /api/user/manage` — 管理员调整余额、启用 / 禁用