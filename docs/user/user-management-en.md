---
title: User Management
description: How admins create, enable, disable, top-up, and delete users.
category: user
order: 2
---

# User Management

> Admin → Users. Visible only when logged in as an admin.

## What you can do

| Action | Where | Effect |
|---|---|---|
| **List** | List page | All users; search, paginate, sort by quota |
| **Add** | Top-right "Add" | Create a normal user; can't create one with higher privilege than yourself |
| **Edit** | Row action | Change display name / password / group / quota |
| **Enable / Disable** | Row action | Disabled users can't log in or call APIs |
| **Promote / Demote** | Row action | Promote to admin / demote to user; Root can't be demoted |
| **Top up** | Row action | Add quota directly to the user (with a remark) |
| **Delete** | Row action | Soft delete; status becomes "Deleted" |

> Some actions (Disable, Delete, Promote) are visible only to Root; regular admins can only Edit and Top up.

## What each row shows

- Username (hover for email)
- Display name
- Group
- Current active subscription (with expiry)
- Quota remaining / used / requests
- Role chip (User / Admin / Root)
- Status chip (Enabled / Disabled / Deleted)
- Registration date

Top search: fuzzy match on `username` / `display_name` / `email`.

## Top up a user

Most common action. Click "Top up" on a row:

- Amount (in internal quota units)
- Remark (recorded in the audit log)
- Confirm

The quota arrives instantly; a top-up log entry is written.

## Disable vs Delete

| | Disable | Delete |
|---|---|---|
| Can user log in? | No | No |
| User data kept? | All | All (soft delete) |
| Visible in list? | Yes (marked "Disabled") | No (filtered out) |
| When to use | Temporary ban (may recover) | Permanently retired |

> Soft delete ≠ physical delete. Orders / call logs are preserved. Root can restore.

## Roles

| Role | Value | Can do |
|---|---|---|
| User | `RoleCommonUser` (=1) | Use API, see own data |
| Admin | `RoleAdminUser` (=10) | + admin Users / Orders / Logs |
| Root | `RoleRootUser` (=100) | + System Settings, delete users, configure payments |

Root accounts can't be disabled / deleted / demoted.

## Bulk actions

Select multiple rows → "Bulk actions" dropdown appears at the bottom:

- Bulk delete
- Bulk disable

A toast tells you which succeeded and which failed (e.g. "can't operate on Root").

## FAQ

- **Adjusted quota but user says it didn't arrive**: the change flushes Redis immediately, the next request should see the new value; if not, check Admin → Logs → User Operations.
- **New user didn't get a verification email**: admin sets the password directly, no email needed unless you require first-login change.
- **Can't delete a user**: probably a Root, or registered on another cluster node. Ask a Root to handle it.

## Related

- [User Schema](../schema/user)
- [Roles & Status](../schema/user)
- [Access Token](./access-token)
- [System Settings (admin)](../misc/system-settings)