---
title: 用户管理
description: "用户 CRUD、状态切换、额度调整与角色管理。"
category: user
order: 2
---

# 用户管理

> 用户的增删改查、状态切换、额度修改、角色提升/降级与按 username 批量操作。

入口路由：`/user`，对应前端组件 `web/default-pro/src/views/user/User.vue`。
所有接口都挂在 `/api/user` 下，admin 路由前缀走 `AdminAuth` 中间件。

## 接口一览 / Endpoints

| Endpoint | Method | 鉴权 | 说明 |
|---|---|---|---|
| `/api/user/` | `GET` | Admin | 分页拉取用户列表（每页 `config.ItemsPerPage`）。`?p=` 页码，`?order=quota\|used_quota` 排序 |
| `/api/user/search?keyword=` | `GET` | Admin | 模糊搜索；命中 `username` / `display_name` / `email` |
| `/api/user/:id` | `GET` | Admin | 拉取单个用户详情 |
| `/api/user/` | `POST` | Admin | 创建普通用户（不能把对方创建成 ≥ 自己权限的角色） |
| `/api/user/manage` | `POST` | Admin | 见下方「单条 manage」说明 |
| `/api/user/manage` | `POST` | Admin（root only） | 批量：`action=batch-delete\|batch-disable`，body 为 `{ usernames: [...] }` |
| `/api/user/` | `PUT` | Admin | 编辑用户（display_name / password / group / quota）；role 不能升级到 ≥ 自己 |
| `/api/user/:id` | `DELETE` | Admin | 删除用户（不能删同级或更高级） |
| `/api/user/topup` | `POST` | Admin | 管理员手动给用户加额度（兼容旧入口，等价于下方管理入口） |

实现：`controller/user.go`；前端 API 调用直接走 `api` 模块（`@/api`）。

## 角色与状态常量 / Roles & Status

| 常量 | 值 | 说明 |
|---|---|---|
| `RoleCommonUser` | `1` | 普通用户 |
| `RoleAdminUser` | `10` | 普通管理员 |
| `RoleRootUser` | `100` | 超级管理员（不可被禁/删/降级） |
| `UserStatusEnabled` | `1` | 启用 |
| `UserStatusDisabled` | `2` | 禁用 |
| `UserStatusDeleted` | `3` | 软删除（列表查询会自动 `WHERE status != deleted`） |

> 前端 router 通过 `user.role >= 10` 判断是否可访问 admin 页面；root (`role >= 100`) 在 UI 中所有按钮都会被 disable。

## POST `/api/user/manage` 单条动作 / Single Actions

请求体：`{ username, action }`，`action` 接受：

| Action | 行为 | 谁可触发 | 审计日志 |
|---|---|---|---|
| `enable` | `user.status = 1` | Admin（不能改 ≥ 自己的用户，root 例外） | 不写 |
| `disable` | `user.status = 2` | **Root** only | `LogTypeManage`：「管理员禁用用户」 |
| `delete` | 软删除 | **Root** only | `LogTypeManage`：「管理员删除用户」 |
| `promote` | `user.role = 10` | **Root** only；已是管理员则报错 | 不写 |
| `demote` | `user.role = 1`；不能对 root 降级 | Admin | 不写 |

返回：`{ success, message, data: { role, status } }`。

## POST `/api/user/manage` 批量动作 / Batch Actions

请求体：`{ action: "batch-delete" | "batch-disable", usernames: [...] }`。

返回：

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

失败项会以 `name: reason` 的形式返回，前端列表层会合并展示为带 `Message.warning` 的 toast。
所有失败时整体视为失败（`success:false`），但 data 仍包含 `failed` 明细，便于排查。

## 编辑用户的副作用 / Side Effects of `PUT /api/user/`

- 当 `origin.quota != updated.quota` 时，会写一条 `LogTypeManage` 审计日志，并把 Redis 中的 `user_quota` 缓存失效（`model.CacheUpdateUserQuota`），下一次请求就会读到更新后的余额。
- `password` 字段为空时不会修改；用哨兵串 `$I_LOVE_U` 跳过 `common.Validate.Struct` 的非空校验。
- 角色提升路径在两层校验：原用户角色不能 ≥ 自己；目标角色也不能 ≥ 自己。两者都通过 root 特例放行。

## 手动加额度 / Manual Top-Up

`POST /api/user/topup` 由 `controller/user.go::AdminTopUp` 处理，请求体：

```json
{ "user_id": 42, "quota": 100000, "remark": "客服补偿" }
```

成功后调用 `model.IncreaseUserQuota` 直接累加，并写一条 `LogTypeTopup` 日志。
没有 `remark` 时会自动生成 `通过 API 充值 <LogQuota(quota)>`。

## 前端操作指南 / Frontend Guide

- 顶部搜索栏支持关键字模糊搜索；切到搜索模式后列表会被替换、禁用追加加载、分页同步重置。
- 排序下拉：`默认` / `quota` / `used_quota`，由后端 `model.GetAllUsers(order)` 处理。
- 行内操作：
  - **编辑**：修改 display_name / password / group / quota。
  - **启用/禁用** / **提升/降级** / **删除**：单条二次确认弹窗；root 用户的对应按钮会被 disable。
- 批量：勾选多行后底部出现「批量操作」下拉，支持「批量删除」「批量禁用」。批量删除是软删除（与单条一致）。
- 用户列展示：用户名（悬停提示邮箱）、显示名、所在分组、当前活跃套餐（计划到期日）、剩余/已用/请求次数、角色 chip、状态 chip、注册日期、操作按钮。

## 接口实现 / Implementation Pointers

| 关注点 | 位置 |
|---|---|
| User CRUD / manage | `controller/user.go` |
| 批量 manage | `controller/user.go::manageUserBatch` |
| 角色/状态常量 | `model/user.go` |
| 软删除行为 | `model/user.go::GetAllUsers`（`WHERE status != UserStatusDeleted`） |
| 手动加额度 | `controller/user.go::AdminTopUp` → `model.IncreaseUserQuota` |
| 配额缓存失效 | `model.CacheUpdateUserQuota`（`UpdateUser` 内调用） |
| 路由 | `router/api.go` |