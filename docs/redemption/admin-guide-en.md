---
title: Redemption Admin
description: "Batch generation, status toggle, edit, delete, and search."
category: redemption
order: 2
---

# Redemption Admin

> Route: `/redemption` (`web/default-pro/src/views/redemption/Redemption.vue`). CRUD implementation: `controller/redemption.go`.

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/redemption/` | `GET` | Admin | Paginated (`config.ItemsPerPage`) |
| `/api/redemption/search?keyword=` | `GET` | Admin | `id=?` or `name LIKE kw%` |
| `/api/redemption/:id` | `GET` | Admin | Single row |
| `/api/redemption/` | `POST` | Admin | Batch generate (see below) |
| `/api/redemption/` | `PUT` | Admin | Edit; with `?status_only=true` only changes `status` |
| `/api/redemption/:id` | `DELETE` | Admin | Hard delete |

## Batch generate — `POST /api/redemption/`

Request body:

```json
{
  "name": "2026-Q1 推广",
  "count": 50,
  "quota": 10000
}
```

Server-side rules (`controller/redemption.go::AddRedemption`):

- `name` length 1–20; otherwise `兑换码名称长度必须在1-20之间`
- `count ∈ [1, 100]`; otherwise `一次兑换码批量生成的个数不能大于 100`
- Each row gets a UUID for `key`; one failed `Insert` aborts the whole batch

`data` is `string[]` (the generated `key` list); the frontend uses this for the "Batch Export" dialog.

## Edit — `PUT /api/redemption/`

```json
{ "id": 88, "name": "新名字", "quota": 20000 }
```

The server `GetRedemptionById` first, then overwrites `Name` and `Quota`. `gorm.Select("name", "status", "quota", "redeemed_time")` keeps updates on the allow-list only.

## Status toggle — `PUT /api/redemption/?status_only=true`

```json
{ "id": 88, "status": 2 }   // 2 = Disabled
```

Only `status` is updated. `status=3 (Used)` cannot be set manually — only the `Redeem` transaction can flip it.

## Delete — `DELETE /api/redemption/:id`

Hard-deletes the `redemptions` row. Used codes can also be deleted; deletion does not touch `users.quota` (already-credited quota is not clawed back).

## Frontend Guide

Route: `/redemption`.

- Top search bar: exact id match or `name LIKE kw%`
- List rows: ID / name / status chip / quota / created_time / redeemed_time (`-` when unused) / actions
- Row actions:
  - **Copy**: writes `key` to the clipboard for distribution
  - **Enable / Disable**: popconfirm + `PUT ?status_only=true`
  - **Edit**: modal for `name` / `quota`
  - **Delete**: secondary confirmation
- **Generate** opens a modal for `name` / `count` / `quota`; the response opens a "Batch Export" dialog showing `data.key[]`

## Implementation Pointers

| Concern | Location |
|---|---|
| CRUD | `controller/redemption.go` |
| Data model | `model/redemption.go::Redemption` |
| Status constants | `model/redemption.go` (`RedemptionCodeStatusEnabled/Disabled/Used`) |
| User-side redeem | `controller/user.go::TopUp` → `model.Redeem` |
