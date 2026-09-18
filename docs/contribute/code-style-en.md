---
title: Code Style
description: "Naming, comments, layout and alignment with AGENTS.md."
category: contribute
order: 2
---

# Code Style

> Naming, comments, layout and alignment with AGENTS.md.
> 命名、注释、目录结构与 AGENTS.md 的一致性。

This page mirrors the root [`AGENTS.md`](https://github.com/modelbus/one-api-pro/blob/main/AGENTS.md); if anything diverges, AGENTS.md is authoritative.

本规范的源头是仓库根目录 [`AGENTS.md`](https://github.com/modelbus/one-api-pro/blob/main/AGENTS.md)。本文是面向贡献者的「导读」+ 对外补充惯例。

## 1. Layout & Frontend / 目录与前端

```text
frontend lives ONLY in web/default-pro/
```

`web/air/ web/default/ web/berry/` are deprecated historical themes. `web/THEMES` only lists `default-pro`. `web/build/<theme>/` is the build output and is overwritten on every build — do not edit.

See [AGENTS.md §1](/en/contribute/dev-setup#2-project-layout).

`web/air/ web/default/ web/berry/` 是历史主题（已废弃）；`web/THEMES` 仅含 `default-pro`；`web/build/<theme>/` 是 `pnpm build` 产物，每次构建被覆盖，禁止手改。

## 2. Naming / 命名

| Kind / 类型 | Convention / 规范 | Example / 例 |
| --- | --- | --- |
| DB constants | PascalCase with category prefix | `OrderTypeTopup`, `SystemSettingKeyTopupEnabled` |
| API path | `/api/<resource>/<action>` | `/api/topup/order`, `/api/setting/topup` |
| Frontend component | PascalCase, multi-word | `TopupModal.vue`, `Dashboard.vue` |
| Frontend API module | camelCase export | `orderApi`, `topupApi` |
| Frontend util | camelCase pure function + co-located test | `src/utils/topup.js` + `topup.test.mjs` |
| Go package | all-lowercase | `model`, `controller`, `middleware` |

> **`topup` is the global name**: internal code uses `topup` (consistent with `LogTypeTopup=1` / `AdminTopUp` / `TopUp`); user-facing strings use 「充值」. Don't mix in `recharge`.
> **`topup` 全局统一**：内部代码命名 `topup`；UI 文案使用「充值」。禁止 `recharge` 命名（除非显式新业务线）。

## 3. Comments / 注释

### 3.1 File Header (required for new files) / 文件头注释

Every new file's first block must include 4 lines (function / English / version / date / author):

新增文件首段 4 行（功能 / 英文 / 版本 / 日期 / 作者）：

```go
// topup.go User self-service top-up business logic
// 版本: v0.0.10
// 日期: 2026-09-06
// 作者: opencode
package model
```

```vue
<!--
  TopupModal.vue Top-up modal component
  Top-up modal component

  版本: v0.0.10
  日期: 2026-09-06
  作者: opencode
-->
```

```js
// topup.js Top-up utility functions
// 版本: v0.0.10
// 日期: 2026-09-06
// 作者: opencode
```

### 3.2 Exported Symbols / 导出符号注释

Every `export function` / `export const` / `func` / `type` must be preceded by a one-line comment:

每个 `export function` / `export const` / `func` / `type` 在前加一行说明：

```go
// CreateTopupOrder validates amount, applies exchange rate, and creates a TP-prefixed order.
// 版本: v0.0.10
func CreateTopupOrder(in CreateTopupOrderInput) (*Order, error) { ... }
```

### 3.3 Inline Comments / 行内 / 块内注释

- Complex business logic: explain "why", not "what".
  复杂业务：解释"为什么"，不止"做了什么"。
- Key validations: cite the source (requirement / post-mortem / business rule).
  关键校验：注明来源（用户需求 / 事故复盘 / 业务约束）。
- `TODO: <scenario> (waiting for <milestone>)` / `XXX: <note> (to be replaced)`.
- Do NOT comment: trivial getters/setters, self-explanatory constants, one-off variables.
  不要注释：简单 getter/setter、命名自解释常量、一次性变量。

### 3.4 Editing Existing Code / 修改现有代码

- Logic change: update `日期:`; breaking change also bump `版本:`.
  逻辑变更 → 更新 `日期:`；破坏性变更同时改 `版本:`。
- Minor fixes (typo / rename): comments may stay.
  小修小补（typo / 重命名）→ 注释可不更新。
- Delete dead comments together with deleted code.
  删除代码 → 一并删除上方注释，不要留无主注释。

## 4. Backend / 后端规范

### 4.1 Unified Response / 统一响应

```go
c.JSON(http.StatusOK, gin.H{
    "success": bool,
    "message": string,
    "data":    <any>,
})
```

Errors also return 200 with `success:false` (legacy convention, easier for the frontend).

错误也用 200 返回 + `success:false`。

### 4.2 Layering / 分层

- `model/`: pure data + business logic, NO HTTP / Gin dependency.
- `controller/`: parameter parsing, model calls, response assembly only.
- `middleware/`: auth / rate-limit / turnstile / language.
- `common/payment/`: payment channels, self-register via `init() → RegisterChannel`.

### 4.3 Payments & Orders / 支付与订单

- Reuse `controller/buildPayInfo(payMethod, orderNo, amount, subject)` for `pay_url / qr_code / note`.
- Order-number prefixes: `TB` new / `UP` upgrade / `TP` topup (`model.GenerateOrderNo(prefix)`).
- Notify dispatch: by `order.Type` in `controller/payment.go::processNotify`.

### 4.4 System Settings / 系统设置

- Generic KV in `model.SystemSetting` table, indexed by `category / key`.
- Category constants: `SystemSettingCategoryPayment / Plan / Topup / General`.
- Key naming: `"<category>.<sub>.<field>"`, e.g. `"topup.enabled"` / `"plan.upgrade_mode"`.
- Legacy migration: remove UI read/write but keep DB row (do NOT actively delete).

### 4.5 Testing / 测试

- Use `glebarez/sqlite` in-memory + `gorm.AutoMigrate`.
- Test names: `TestXxx_Yyy` (behavior_scenario).
- Must cover: setting save, amount math, order activation idempotency.

## 5. Frontend / 前端规范

### 5.1 Basics / 基础

- Vue 3 `<script setup>` + Composition API.
- Alias `@` → `src/` (`vite.config.js`).
- Components PascalCase and **multi-word** (no `Index.vue` / `Modal.vue`).

### 5.2 API Modules / API 模块

```js
// src/api/order.js
import api from './index'

export const orderApi = {
  myOrders: (type) => api.get('/api/order/self', { params: type ? { type } : {} }),
}

export default orderApi
```

Call convention / 调用约定:

```js
const { data } = await topupApi.createOrder(payload)
if (data.success && data.data) { /* ... */ }
```

### 5.3 arco `<a-modal>` uses `:visible`, not `:model-value`

Good / 正例:

```vue
<a-modal
  :visible="modelValue"
  @update:visible="(v) => emit('update:modelValue', v)"
  @cancel="close"
  :title="title"
  :footer="false"
>
```

Bad (modal never shows) / 反例（弹窗永远不显示）:

```vue
<a-modal :model-value="modelValue" @cancel="close">  <!-- ❌ -->
```

### 5.4 arco-spin needs `style="width: 100%"`

```vue
<a-spin :loading="loading" style="width: 100%">
```

Without it, the parent container width behaves incorrectly.

不加会导致父容器宽度异常。

### 5.5 Forms / 表单

- Inline / number: `<a-input-number :precision="2">`.
- Switch: `<a-switch v-model="x">`.
- Selector: `<a-select v-model="...">` + `<a-option value="...">label</a-option>`.
- Validation: call a validator before submit; toast via `Message.error`.

### 5.6 Utils & Unit Tests / 工具函数与单测

- Pure logic goes to `src/utils/<name>.js`. Do NOT put complex logic inside Vue components.
- Tests: co-located `<name>.test.mjs`, using Node's built-in `node:test` + `node:assert/strict`.

```js
import { test } from 'node:test'
import assert from 'node:assert/strict'
import { validateTopupPresets } from './topup.js'

test('rejects duplicate amount', () => {
  assert.match(validateTopupPresets([{ amount: 10 }, { amount: 10 }]), /重复/)
})
```

### 5.7 Permissions / 权限

- `useAuthStore().isAdmin` / `.isRoot` gates UI.
- Backend middleware enforces; do NOT rely on frontend hiding only.

## 6. DB & Migrations / 数据库与迁移

- After changing `model/*.go`, ensure `AutoMigrate` works (non-breaking).
- Do NOT write destructive DDL (`DROP` / `TRUNCATE`) in business code.
- Complex migrations go to one-shot `cmd/migrate_*` scripts.

## 7. Known Pitfalls / 已知踩坑

1. arco `<a-modal>` uses `:visible`, not `:model-value`.
2. arco `<a-spin>` needs `style="width: 100%"`.
3. `plan.allow_topup` migrated to `topup.*`; UI no longer reads/writes; DB row kept.
4. Don't hardcode plan-name whitelists in Vue components (`VALID_PLAN_NAMES` etc.) — the backend API already accepts any plan name.
5. Always run the trio before committing: `go build ./...` + `go test ./model/ ./controller/ ./middleware/` + `pnpm build`.
6. Commit body must describe detailed changes in both Chinese and English (see [Commit Convention §3.4](/en/contribute/commit-convention#34-commit-body-critical-commit-body-critical)).
7. Granularity: one commit per file; split multi-feature files.
8. Frontend dir: `web/default-pro/` only; the other web/* dirs are historical themes, do not edit.

Next: [Commit Convention](/en/contribute/commit-convention) · [Add a Provider](/en/contribute/add-provider).