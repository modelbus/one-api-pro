---
title: 兑换码管理（管理员）
description: "批量生成、状态切换、编辑、删除与查询。"
category: redemption
order: 2
---

# 兑换码管理（管理员）

> 路径：`/redemption`（`web/default-pro/src/views/redemption/Redemption.vue`）。CRUD 实现在 `controller/redemption.go`。

## 接口 / Endpoints

| Endpoint | Method | Auth | 说明 |
|---|---|---|---|
| `/api/redemption/` | `GET` | Admin | 分页（`config.ItemsPerPage`） |
| `/api/redemption/search?keyword=` | `GET` | Admin | `id=?` 或 `name LIKE kw%` |
| `/api/redemption/:id` | `GET` | Admin | 单条详情 |
| `/api/redemption/` | `POST` | Admin | 批量生成（详见下文） |
| `/api/redemption/` | `PUT` | Admin | 编辑；带 `?status_only=true` 时仅改状态 |
| `/api/redemption/:id` | `DELETE` | Admin | 硬删除 |

## 批量生成 / `POST /api/redemption/`

请求体：

```json
{
  "name": "2026-Q1 推广",
  "count": 50,
  "quota": 10000
}
```

服务端约束（`controller/redemption.go::AddRedemption`）：

- `name` 长度 1–20 字符；不合法返回 `兑换码名称长度必须在1-20之间`
- `count ∈ [1, 100]`；超出返回 `一次兑换码批量生成的个数不能大于 100`
- 每条生成 UUID 作为 `key`，逐行 `Insert`；失败时整批返回错误

响应 `data` 是 `string[]`（即生成的 `key` 列表），前端用于"批量导出"对话框。

## 编辑 / `PUT /api/redemption/`

```json
{ "id": 88, "name": "新名字", "quota": 20000 }
```

服务端先 `GetRedemptionById`，再覆盖 `Name` 与 `Quota`。`gorm.Select("name", "status", "quota", "redeemed_time")` 保证只更新白名单字段。

## 状态切换 / `PUT /api/redemption/?status_only=true`

```json
{ "id": 88, "status": 2 }   // 2 = Disabled
```

仅改 `status`；`status=3 (Used)` 不可手动设置——只能由 `Redeem` 事务在兑换成功后翻转。

## 删除 / `DELETE /api/redemption/:id`

硬删除 `redemptions` 行。已使用 (`status=3`) 的兑换码也可以删除；删除不影响 `users.quota`（已经发放的额度不会回退）。

## 前端操作指南 / Frontend Guide

页面：`/redemption`。

- 顶部搜索栏：按 id 精确匹配或 `name LIKE kw%`
- 列表行：ID / name / status chip / quota / 创建时间 / 兑换时间（未兑换显示 `-`）/ 操作
- 行内操作：
  - **复制**：把 `key` 写入剪贴板，便于发码
  - **启用 / 禁用**：popconfirm 后调 `PUT ?status_only=true`
  - **编辑**：弹窗修改 `name` / `quota`
  - **删除**：二次确认
- 「生成」按钮打开弹窗，填 `name` / `count` / `quota`；提交后弹出"批量导出"对话框展示 `data.key[]`

## 实现位置 / Implementation Pointers

| 关注点 | 位置 |
|---|---|
| CRUD | `controller/redemption.go` |
| 数据模型 | `model/redemption.go::Redemption` |
| 状态常量 | `model/redemption.go` (`RedemptionCodeStatusEnabled/Disabled/Used`) |
| 用户侧兑换 | `controller/user.go::TopUp` → `model.Redeem` |
