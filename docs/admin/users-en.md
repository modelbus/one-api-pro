---
title: Users Admin
description: "User CRUD, status toggle, quota adjustment and role management."
category: admin
order: 2
---

# Users Admin

> Create / list / edit / disable users, adjust quota, promote or demote roles, plus batch operations by username.

Route: `/user`, served by `web/default-pro/src/views/user/User.vue`. All endpoints live under `/api/user` and the admin-prefixed routes are guarded by `AdminAuth`.

## Endpoints

| Endpoint | Method | Auth | Description |
|---|---|---|---|
| `/api/user/` | `GET` | Admin | Paginated list (page size = `config.ItemsPerPage`). `?p=` page, `?order=quota\|used_quota` |
| `/api/user/search?keyword=` | `GET` | Admin | Search by `username` / `display_name` / `email` |
| `/api/user/:id` | `GET` | Admin | Single-user detail |
| `/api/user/` | `POST` | Admin | Create a user; the target role must be `<` the caller's role |
| `/api/user/manage` | `POST` | Admin | Single-action manage (see below) |
| `/api/user/manage` | `POST` | **Root only** | Batch action: `{ action: "batch-delete"\|"batch-disable", usernames: [...] }` |
| `/api/user/` | `PUT` | Admin | Edit `display_name` / `password` / `group` / `quota`; role cannot be raised to `>=` the caller |
| `/api/user/:id` | `DELETE` | Admin | Delete user (cannot delete peers or higher) |
| `/api/user/topup` | `POST` | Admin | Manually grant quota to a user (legacy entrypoint) |

Implementation: `controller/user.go`; the frontend calls `api` module directly.

## Role & Status Constants

| Constant | Value | Meaning |
|---|---|---|
| `RoleCommonUser` | `1` | Regular user |
| `RoleAdminUser` | `10` | Admin |
| `RoleRootUser` | `100` | Super-admin (cannot be disabled / deleted / demoted) |
| `UserStatusEnabled` | `1` | Enabled |
| `UserStatusDisabled` | `2` | Disabled |
| `UserStatusDeleted` | `3` | Soft-deleted (lists filter on `status != UserStatusDeleted`) |

The router guard checks `user.role >= 10` for admin pages; rows with `role >= 100` have all action buttons disabled in the UI.

## POST `/api/user/manage` — Single Action

Body: `{ username, action }`. Supported actions:

| Action | Effect | Who | Audit |
|---|---|---|---|
| `enable` | `user.status = 1` | Admin (cannot target a user with `role >=` self; root bypass) | – |
| `disable` | `user.status = 2` | **Root** only | `LogTypeManage`: "管理员禁用用户" |
| `delete` | Soft delete | **Root** only | `LogTypeManage`: "管理员删除用户" |
| `promote` | `user.role = 10` | **Root** only; errors out if already an admin | – |
| `demote` | `user.role = 1`; root excluded | Admin | – |

Response: `{ success, message, data: { role, status } }`.

## POST `/api/user/manage` — Batch

Body: `{ action: "batch-delete" | "batch-disable", usernames: [...] }`.

Response:

```jsonc
{
  "success": true,
  "data": {
    "succeeded": ["alice"],
    "failed": { "bob": "无法操作超级管理员用户" },
    "total": 2,
    "succeeded_count": 1,
    "failed_count": 1
  }
}
```

Per-row failures are returned as `name: reason`; the frontend surfaces partial failures as `Message.warning` and keeps the partial-success toast. When **every** row fails the response is `success: false` but the `failed` map is still populated.

## Side Effects of `PUT /api/user/`

- When `origin.quota != updated.quota` an audit log of type `LogTypeManage` is written, and the Redis `user_quota` cache for the user is invalidated (`model.CacheUpdateUserQuota`) so the next request sees the new balance.
- A blank `password` is treated as "no change"; the sentinel string `$I_LOVE_U` is used to pass `common.Validate.Struct` and rolled back afterwards.
- Role promotion is gated twice: the original user's role must be `<` the caller, and the target role must also be `<` the caller — both relaxed for `RoleRootUser`.

## Manual Top-Up

`POST /api/user/topup` (legacy endpoint in `controller/user.go::AdminTopUp`):

```json
{ "user_id": 42, "quota": 100000, "remark": "support compensation" }
```

`model.IncreaseUserQuota` adds to the user's balance; a `LogTypeTopup` entry is written (default remark: `通过 API 充值 <LogQuota(quota)>`).

## Frontend Guide

- The search box performs a fuzzy search; switching to search mode replaces the list (no append, pagination resets).
- Sort dropdown: `default` / `quota` / `used_quota` — handled server-side by `model.GetAllUsers(order)`.
- Row actions: **Edit** (display_name / password / group / quota), **Enable / Disable**, **Promote / Demote**, **Delete**, each with a confirmation popover. All buttons are disabled for root users.
- Batch select: a checkbox column appears once any row is checked. The bottom toolbar shows a count and a "批量操作" dropdown with `批量删除` and `批量禁用`. Batch delete is a soft delete, matching the single-row action.
- Row columns: ID, username (tooltip = email), display name, group, current active plan (chip with name · billing · expiry), remaining / used / request counts, role chip, status chip, registered date, actions.

## Implementation Pointers

| Concern | Location |
|---|---|
| CRUD + manage handler | `controller/user.go` |
| Batch helper | `controller/user.go::manageUserBatch` |
| Role/status constants | `model/user.go` |
| Soft-delete behavior | `model/user.go::GetAllUsers` |
| Manual quota grant | `controller/user.go::AdminTopUp` → `model.IncreaseUserQuota` |
| Quota cache invalidation | `model.CacheUpdateUserQuota` (called inside `UpdateUser`) |
| Route table | `router/api.go` |