---
title: 数据结构总览
description: One API Pro 持久化的核心数据模型及它们之间的关系。
category: schema
order: 1
---

# 数据结构总览

> 本目录解释 One API Pro 在数据库中持久化的核心数据，回答这些问题：
>
> - 这些数据是什么、由谁创建、在哪里能看到？
> - 我该在哪里修改它们？
> - 修改后会影响哪些行为？

如果你正在查找 HTTP 接口，请前往 [API 参考](/api/README)。

## 阅读建议

如果你是 **第一次接触**，按顺序读：

1. [用户](/schema/user) — 整个系统的身份与额度中心
2. [渠道](/schema/channel) — 接入上游大模型的最小单位
3. [套餐与订阅](/schema/plan) 与 [订阅](/schema/subscription) — 商业化模型
4. [订单](/schema/order) — 用户实际付费的记录

如果你是 **开发者 / 排查问题**，按需查阅：

- 用户额度异常 → [用户](/schema/user) → [访问令牌](/schema/token) → [调用日志](/schema/log)
- 路由是否生效 → [渠道](/schema/channel) → [模型定价](/schema/model-price) → [分组折扣](/schema/group-price)
- 订阅/扣费策略 → [套餐](/schema/plan) → [订阅](/schema/subscription) → [订单](/schema/order) → [充值](/schema/topup)

## 关系一览

| 数据 | 谁创建 | 谁消费 | 在哪里管理 |
|---|---|---|---|
| User | 注册 / 管理员 | 所有功能 | 后台「用户管理」 |
| Token | 用户 | 管理后台 API | 后台「API Token 管理」 |
| Channel | 管理员 | 渠道路由器 | 后台「渠道」 |
| ModelPrice | 管理员 | 计费 | 后台「模型定价」 |
| GroupPrice | 管理员 | 计费 | 后台「分组折扣」 |
| Plan | 管理员 | 用户下单 | 后台「套餐管理」 |
| Subscription | 系统（下单后） | 计费 / 限速 | 后台「订阅管理」 |
| Order | 用户 | 财务对账 | 后台「订单」 |
| Topup | 用户 | 计费 | 后台「充值管理」 |
| Redemption | 管理员 | 用户兑换 | 后台「兑换码管理」 |
| ClusterNode | 节点启动时注册 | 集群同步 | 后台「集群设置」 |
| Log | 系统（每次调用） | 管理员排查 | 后台「日志」 |
| Option | 管理员 | 全局行为 | 后台「系统设置」 |