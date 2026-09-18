---
title: Redemption Codes
description: "Bulk-generate redemption codes, toggle status, copy keys and search."
category: redemption
order: 9
---

# Redemption Codes

> Bulk-generate redemption codes (each carrying a quota), toggle their status, search by keyword. UI: `web/default-pro/src/views/redemption/Redemption.vue`.

## Data Model

`model.Redemption` (table `redemptions`):

| Field | Type | Notes |
|---|---|---|
| `id` | `int` PK | Internal id |
| `user_id` | `int` | Creator (admin id) |
| `key` | `char(32)` UNIQUE | The code itself (UUID without dashes) |
| `name` | `varchar` (indexed) | Code batch name |
| `quota` | `bigint` | Quota granted on redemption |
| `status` | `int` | `RedemptionCodeStatusEnabled=1` / `Disabled=2` / `Used=3` (default `1`) |
| `count` | `int` (not persisted) | Batch size, used only on creation; serialized with `gorm:"-:all"` |
| `created_time` / `redeemed_time` | `bigint` | unix seconds; `redeemed_time=0` means unused |

Status values intentionally avoid `0` so a default-zero row never looks "used".

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/redemption/` | `GET` | Admin | Paginated list (`config.ItemsPerPage`) |
| `/api/redemption/search?keyword=` | `GET` | Admin | `id=?` OR `name LIKE kw%` |
| `/api/redemption/:id` | `GET` | Admin | Single detail |
| `/api/redemption/` | `POST` | Admin | Generate a batch (`count` 1–100) |
| `/api/redemption/` | `PUT` | Admin | Edit (`?status_only=true` only changes status) |
| `/api/redemption/:id` | `DELETE` | Admin | Delete one |

Implementation: `controller/redemption.go`.

## Bulk Generate — `POST /api/redemption/`

```json
{ "name": "2026-Q1 Promo", "count": 50, "quota": 10000 }
```

Constraints:
- `name` length 1–20 chars;
- `count` ∈ [1, 100] — returns `一次兑换码批量生成的个数不能大于 100` otherwise;
- Each code receives a UUID as `key`; the response data is the array of `key`s for export.

## Status Toggle — `PUT /api/redemption/?status_only=true`

```json
{ "id": 88, "status": 2 }
```

With `?status_only=true`, only `status` is updated; otherwise `name` / `quota` are also written.
`status=3` (Used) cannot be set manually — it is flipped inside the user-side redeem transaction.

## User-side Redemption (transactional)

`model.Redeem(ctx, key, userId)` is what `POST /api/user/topup` calls:

1. `SELECT … FOR UPDATE` on `redemptions.key`;
2. Assert `status = RedemptionCodeStatusEnabled`;
3. `UPDATE users SET quota = quota + redemptions.quota WHERE id = user_id`;
4. Set `redeemed_time = now()`, `status = Used`;
5. Record a `LogTypeTopup` log (`通过兑换码充值 <LogQuota>`).

The whole flow runs inside `gorm.Transaction` to make concurrent redemption safe.

## Frontend Guide

- Search bar (id or `name` prefix).
- Row columns: ID / name / status chip / quota / created time / redeemed time (`-` if unused) / actions.
- Row actions:
  - **Copy**: copies `key` to the clipboard.
  - **Enable / Disable**: popconfirm → `PUT ?status_only=true`.
  - **Edit**: modal for `name` / `quota`.
  - **Delete**: confirmation popover.
- **Generate** opens a modal asking for `name`, `count`, `quota`. After submission the generated `key[]` opens the bulk-export dialog.

## Implementation Pointers

| Concern | Location |
|---|---|
| CRUD handler | `controller/redemption.go` |
| Redeem transaction | `model/redemption.go::Redeem` |
| Status constants | `model/redemption.go` |
| User redeem entrypoint | `controller/user.go::TopUp` |
| Routes | `router/api.go` |