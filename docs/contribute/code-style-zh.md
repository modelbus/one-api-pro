---
title: 编码规范
description: "命名、注释、目录结构与 AGENTS.md 的一致性。"
category: contribute
order: 2
---

# 编码规范

> 命名、注释、目录结构与 AGENTS.md 的一致性。
> Naming, comments, layout and alignment with AGENTS.md.

本规范的源头是仓库根目录 [`AGENTS.md`](https://github.com/modelbus/one-api-pro/blob/main/AGENTS.md)。本文是面向贡献者的「导读」+ 对外补充惯例。

This page mirrors the root [`AGENTS.md`](https://github.com/modelbus/one-api-pro/blob/main/AGENTS.md); if anything diverges, AGENTS.md is authoritative.

## 1. 目录与前端 / Layout & Frontend

```text
frontend lives ONLY in web/default-pro/
```

`web/air/ web/default/ web/berry/` 是历史主题（已废弃）；
`web/THEMES` 仅含 `default-pro`；
`web/build/<theme>/` 是 `pnpm build` 产物，每次构建被覆盖，禁止手改。

详见 [AGENTS.md §1](/zh/contribute/dev-setup#2-项目结构--project-layout)。

`web/air/ web/default/ web/berry/` are deprecated historical themes. `web/THEMES` only lists `default-pro`. `web/build/<theme>/` is the build output and is overwritten on every build — do not edit.

See [AGENTS.md §1](/en/contribute/dev-setup#2-project-layout).

## 2. 命名约定 / Naming

| 类型 / Kind | 规范 / Convention | 例 / Example |
| --- | --- | --- |
| 数据库常量 | PascalCase + 类别前缀 | `OrderTypeTopup`、`SystemSettingKeyTopupEnabled` |
| API 路径 | `/api/<resource>/<action>` | `/api/topup/order`、`/api/setting/topup` |
| 前端组件 | PascalCase 多词 | `TopupModal.vue`、`Dashboard.vue` |
| 前端 API 模块 | camelCase export | `orderApi`、`topupApi` |
| 前端工具函数 | camelCase 纯函数 + 同目录单测 | `src/utils/topup.js` + `topup.test.mjs` |
| Go 包名 | 全小写单词 | `model`, `controller`, `middleware` |

> **`topup` 全局统一**：内部代码命名 `topup`（与历史 `LogTypeTopup=1` / `AdminTopUp` / `TopUp` 一致）；UI 文案使用「充值」。禁止 `recharge` 命名（除非显式新业务线）。
> **`topup` is the global name**: internal code uses `topup` (consistent with `LogTypeTopup=1` / `AdminTopUp` / `TopUp`); user-facing strings use 「充值」. Don't mix in `recharge`.

## 3. 注释约定 / Comments

### 3.1 文件头注释（新增文件必填）

新增文件首段 4 行（功能 / 英文 / 版本 / 日期 / 作者）：

Every new file's first block must include 4 lines (function / English / version / date / author):

```go
// topup.go 用户自助充值业务
// User self-service top-up business logic
// 版本: v0.0.10
// 日期: 2026-09-06
// 作者: opencode
package model
```

```vue
<!--
  TopupModal.vue 充值弹窗组件
  Top-up modal component

  版本: v0.0.10
  日期: 2026-09-06
  作者: opencode
-->
```

```js
// topup.js 充值工具函数
// Top-up utility functions
// 版本: v0.0.10
// 日期: 2026-09-06
// 作者: opencode
```

### 3.2 导出符号注释

每个 `export function` / `export const` / `func` / `type` 在前加一行说明：

Every `export function` / `export const` / `func` / `type` must be preceded by a one-line comment:

```go
// CreateTopupOrder 校验金额、汇率换算、生成 TP 订单号
// 版本: v0.0.10
func CreateTopupOrder(in CreateTopupOrderInput) (*Order, error) { ... }
```

### 3.3 行内 / 块内注释（按需）

- **复杂业务**：解释"为什么"，不止"做了什么"。
  Complex business logic: explain "why", not "what".
- **关键校验**：注明来源（用户需求 / 事故复盘 / 业务约束）。
  Key validations: cite the source (requirement / post-mortem / business rule).
- `TODO: <场景>（待 <里程碑>）` / `XXX: <说明>（待替换）`。
- **不要注释**：简单 getter/setter、命名自解释常量、一次性变量。
  Do NOT comment: trivial getters/setters, self-explanatory constants, one-off variables.

### 3.4 修改现有代码

- **逻辑变更** → 更新 `日期:`；破坏性变更同时改 `版本:`。
  Logic change: update `日期:`; breaking change also bump `版本:`.
- **小修小补**（typo / 重命名）→ 注释可不更新。
  Minor fixes (typo / rename): comments may stay.
- **删除代码** → 一并删除上方注释，不要留无主注释。
  Delete dead comments together with deleted code.

## 4. 后端规范 / Backend

### 4.1 统一响应

```go
c.JSON(http.StatusOK, gin.H{
    "success": bool,
    "message": string,
    "data":    <any>,
})
```

错误也用 200 返回 + `success:false`（项目历史约定，便于前端统一处理）。

Errors also return 200 with `success:false` (legacy convention, easier for the frontend).

### 4.2 分层

- `model/`：纯数据 + 业务函数，**不要**引入 HTTP / Gin 依赖。
  `model/`: pure data + business logic, NO HTTP / Gin dependency.
- `controller/`：只做参数解析、调用 model、组装响应。
  `controller/`: parameter parsing, model calls, response assembly only.
- `middleware/`：auth / rate-limit / turnstile / language。
- `common/payment/`：各支付通道（wechat / alipay / bank），通过 `init() → RegisterChannel` 自注册。

### 4.3 支付与订单

- 复用 `controller/buildPayInfo(payMethod, orderNo, amount, subject)` 生成支付参数（`pay_url / qr_code / note`）。
- 订单号前缀：`TB` 新订 / `UP` 差价升级 / `TP` 充值（`model.GenerateOrderNo(prefix)`）。
- 通知分发：`controller/payment.go::processNotify` 按 `order.Type` 分发。

### 4.4 系统设置

- 通用 KV 存 `model.SystemSetting` 表，`category / key` 维度。
- 类别常量：`SystemSettingCategoryPayment / Plan / Topup / General`。
- key 命名：`"<category>.<sub>.<field>"`，例：`"topup.enabled"` / `"plan.upgrade_mode"`。
- **旧字段迁移**：UI 删除读写但 DB 行保留（不要主动删）。

### 4.5 测试

- 用 `glebarez/sqlite` 内存库 + `gorm.AutoMigrate` 建表。
- 测试函数命名：`TestXxx_Yyy`（行为_场景）。
- 关键路径必须覆盖：设置保存、金额计算、订单激活幂等。

## 5. 前端规范 / Frontend

### 5.1 基础

- Vue 3 `<script setup>` + Composition API
- 路径别名 `@` → `src/`（`vite.config.js`）
- 组件命名 PascalCase，**多词**（不要 `Index.vue` / `Modal.vue` 这类单泛词）

### 5.2 API 模块

```js
// src/api/order.js
import api from './index'

export const orderApi = {
  myOrders: (type) => api.get('/api/order/self', { params: type ? { type } : {} }),
}

export default orderApi
```

调用约定：

```js
const { data } = await topupApi.createOrder(payload)
if (data.success && data.data) { /* ... */ }
```

### 5.3 arco `<a-modal>` 必须 `:visible`，不是 `:model-value`

正例：

```vue
<a-modal
  :visible="modelValue"
  @update:visible="(v) => emit('update:modelValue', v)"
  @cancel="close"
  :title="title"
  :footer="false"
>
```

反例（弹窗永远不显示）：

```vue
<a-modal :model-value="modelValue" @cancel="close">  <!-- ❌ -->
```

### 5.4 arco-spin 加 `style="width: 100%"`

```vue
<a-spin :loading="loading" style="width: 100%">
```

不加会导致父容器宽度异常。

### 5.5 表单

- 行内 / 数字：`<a-input-number :precision="2">`
- 开关：`<a-switch v-model="x">`
- 选择器：`<a-select v-model="...">` + `<a-option value="...">label</a-option>`
- 表单校验：函数式校验在提交前调用，弹 `Message.error`。

### 5.6 工具函数与单测

- 纯逻辑抽到 `src/utils/<name>.js`，**不要**在组件里写复杂内联逻辑。
- 单测：同目录 `<name>.test.mjs`，用 Node 自带 `node:test` + `node:assert/strict`。

```js
import { test } from 'node:test'
import assert from 'node:assert/strict'
import { validateTopupPresets } from './topup.js'

test('rejects duplicate amount', () => {
  assert.match(validateTopupPresets([{ amount: 10 }, { amount: 10 }]), /重复/)
})
```

### 5.7 权限

- `useAuthStore().isAdmin` / `.isRoot` 控制 UI 显隐；
- 后端中间件兜底（不要只靠前端隐藏）。

## 6. 数据库与迁移 / DB & Migrations

- 改 `model/*.go` 后必须保证 `AutoMigrate` 能正确建表（不破坏性）。
- 不在业务代码中写破坏性 DDL（`DROP` / `TRUNCATE` 等）。
- 复杂迁移写 `cmd/migrate_*` 一次性脚本。

## 7. 已知踩坑 / Known Pitfalls

1. **arco `<a-modal>` 用 `:visible`，不是 `:model-value`**。
2. **arco `<a-spin>` 加 `style="width: 100%"`**。
3. **`plan.allow_topup` 已迁移到 `topup.*`**，UI 不再读写，DB 行保留。
4. **不要在 Vue 组件里硬编码套餐名白名单**（`VALID_PLAN_NAMES` 等）—— 后端 API 已支持任意套餐名。
5. **提交前必跑三件套**：`go build ./...` + `go test ./model/ ./controller/ ./middleware/` + `pnpm build`。
6. **Commit body 必须中英双语说明详细变更**（[提交规范 §3.4](/zh/contribute/commit-convention#34-commit-body-写法关键-commit-body-critical)）。
7. **粒度**：一文件一 commit；同文件多功能必须拆 commit。
8. **前端目录**：仅 `web/default-pro/`，其它 web 子目录是历史主题，禁改。

下一步 / Next: [提交规范](/zh/contribute/commit-convention) · [新增 Provider](/zh/contribute/add-provider)。