---
title: 兑换码
description: 可批量生成并发放给用户的额度 / 套餐兑换码。
category: schema
order: 11
---

# 兑换码（Redemption）

## 这是什么

`Redemption` 是**批量发放的兑换码**。管理员生成后，用户在「兑换码中心」输入兑换码即可获得额度或激活套餐。

## 在哪里看到

- **后台 → 兑换码管理**：列表、批量导出 CSV / TXT
- **前台 → 兑换码**：用户输入兑换码
- **后台 → 用户管理 → 详情**：用户已兑换过的记录

## 你需要为「兑换码功能」设置什么

| 项 | 在哪里 | 设置后影响 |
|---|---|---|
| `type` | 兑换码管理 → 新增 | 「额度」直接加 User.quota；「套餐」创建 Subscription |
| `quota` | 兑换码管理 → 新增 | 兑换后写入额度 |
| `plan_id` | 兑换码管理 → 新增 | 仅「套餐」类型需要；兑换后激活 |
| `expired_time` | 兑换码管理 → 新增 | 过期后该码不可再兑换 |
| `max_redemptions` | 兑换码管理 → 新增 | 默认 1 |

## 关键字段（每张码）

| 字段 | 含义 | 设置后影响 |
|---|---|---|
| `key` | 兑换码字符串 | 用户输入的凭证；改了就发不出去了 |
| `type` | 类型 | 额度 / 套餐；决定兑换后的入账动作 |
| `quota` | 额度值 | 兑换后写入额度（仅「额度」类型） |
| `plan_id` | 关联套餐 | 兑换后激活的套餐（仅「套餐」类型） |
| `max_redemptions` | 总次数 | 默认 1；大于 1 时支持「一码多兑」 |
| `redeemed_count` | 已兑次数 | 自动递增；满后失效 |
| `expired_time` | 过期时间 | Unix 秒 / 0 = 永不过期 |
| `status` | 状态 | 未使用 / 已用完 / 已过期；决定是否还能兑换 |

## 兑换后的入账

- 额度类型：`User.quota += 兑换额度`，无需审批
- 套餐类型：创建一条 [Subscription](/schema/subscription)，状态为「生效中」，到期时间 = 现在 + Plan 的 `duration_days`

## 兑换码格式

默认生成 24 位随机字符串。批量发放建议使用 CSV 导出，后台也支持 TXT 一码一行。

## 相关页面

- [兑换码概览](/redemption/overview)
- [兑换码管理（后台）](/redemption/redemption-management)
- [用户兑换流程](/redemption/user-guide)
- [兑换额度计算](/redemption/quotas)

## 相关 API

- `POST /api/redemption/` — 创建兑换码
- `GET /api/redemption/` — 列表
- `POST /api/redemption/redeem` — 用户兑换（前台）
- `DELETE /api/redemption/:id` — 删除