---
title: 兑换码管理
description: "兑换码批量生成、状态切换、复制密钥与查询。"
category: redemption
order: 9
---

# 兑换码管理

> 在 `redemptions` 表上批量生成兑换码（每张含 `quota`），开关状态、按关键字查询。前端组件：`web/default-pro/src/views/redemption/Redemption.vue`。

## 数据模型

`model.Redemption`（`redemptions` 表）：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | `int` PK | 内部 ID |
| `user_id` | `int` | 创建者（admin id） |
| `key` | `char(32)` UNIQUE | 兑换码本身（UUID，去掉横线） |
| `name` | `varchar`（索引） | 兑换码名（用于批量管理） |
| `quota` | `bigint` | 兑换后到账的额度（quota） |
| `status` | `int` | `RedemptionCodeStatusEnabled=1` / `Disabled=2` / `Used=3`（默认 1） |
| `count` | `int`（不落库） | 仅用于生成时的批量大小，序列化时 `g:"-all"` |
| `created_time` / `redeemed_time` | `bigint` | unix 秒；`redeemed_time=0` 表示未使用 |

> 状态值刻意避开 `0`，避免未初始化默认值撞"已使用"语义。

## 接口一览

| Endpoint | Method | 鉴权 | 说明 |
|---|---|---|---|
| `/api/redemption/` | `GET` | Admin | 分页拉取（`config.ItemsPerPage`） |
| `/api/redemption/search?keyword=` | `GET` | Admin | `id=?` 或 `name LIKE kw%` |
| `/api/redemption/:id` | `GET` | Admin | 单条详情 |
| `/api/redemption/` | `POST` | Admin | 批量生成（`count` 限制 1–100） |
| `/api/redemption/` | `PUT` | Admin | 编辑（`?status_only=true` 时仅改状态） |
| `/api/redemption/:id` | `DELETE` | Admin | 删除一条 |

实现：`controller/redemption.go`。

## 批量生成

请求体：

```json
{
  "name": "2026-Q1 推广",
  "count": 50,
  "quota": 10000
}
```

约束：
- `name` 长度 1–20 字符；
- `count` ∈ [1, 100]；超过 100 时返回 `一次兑换码批量生成的个数不能大于 100`；
- 服务端为每条生成 UUID 作为 `key`，写入后把 `key[]` 作为 `data` 返回，前端可用于「批量导出」。

## 状态切换

```json
{ "id": 88, "status": 2 }   // 2 = disabled
```

带 `?status_only=true` 时仅修改 `status` 字段；否则会一并更新 `name` / `quota`。
`status=3`（已使用）不可手动设置；该状态只能由用户调用 `/api/user/topup` 时在事务中翻转。

## 用户兑换

`model.Redeem(ctx, key, userId)` 是用户侧 `POST /api/user/topup` 的真实调用：

1. `SELECT … FOR UPDATE` 锁定 `redemptions.key`；
2. 校验 `status=RedemptionCodeStatusEnabled`；
3. `UPDATE users SET quota = quota + redemptions.quota WHERE id = user_id`；
4. 把 `redeemed_time = now()`、`status = Used`；
5. 落 `LogTypeTopup` 日志（`通过兑换码充值 <LogQuota>`）。
整个流程在 `gorm.Transaction` 内执行，避免并发兑换。

## 前端操作指南

- 顶部搜索栏（按 id 或 name 前缀）。
- 列表行：ID / name / status chip / quota / 创建时间 / 兑换时间（未兑换显示 `-`）/ 操作。
- 行内操作：
  - **复制**：把 `key` 复制到剪贴板；
  - **启用 / 禁用**：popconfirm 后调 `PUT ?status_only=true`；
  - **编辑**：打开弹窗修改 `name` / `quota`；
  - **删除**：二次确认。
- 「生成」按钮打开弹窗，填 `name`、`count`、`quota`，提交后弹出"批量导出"对话框（`data.key[]`）。

## 接口实现

| 关注点 | 位置 |
|---|---|
| CRUD handler | `controller/redemption.go` |
| 兑换事务 | `model/redemption.go::Redeem` |
| 状态常量 | `model/redemption.go`（`RedemptionCodeStatusEnabled/Disabled/Used`） |
| 用户侧兑换入口 | `controller/user.go::TopUp` |
| 路由 | `router/api.go` |
