# AGENTS.md — One API Pro 开发规范 / Development Conventions

> 面向所有 Agent 与贡献者 / For all agents and contributors
>
> 本文档汇总本项目不可从代码自动推断、但每次改动都必须遵守的约定。
> This document captures conventions that cannot be inferred from code alone but must be followed for every change.

---

## TL;DR — 30 秒必读 / 30-second Quick Start

- 前端维护目录：**仅** `web/default-pro/` / Frontend lives in **`web/default-pro/` only**
- 提交风格：Conventional Commits + 中文 commit body（本次修改的中文+英文说明） / Conventional Commits with Chinese commit body (Chinese + English explanation of the change)
- 粒度：默认一文件一 commit；同文件多功能要拆分 / One file per commit by default; split a file into multiple commits if it contains multiple features
- 版本号读根目录 `VERSION` / Version from root `VERSION` file
- 新增文件 4 行文件头注释（功能 / 版本 / 日期 / 作者） / New files: 4-line header comment
- 弹窗用 arco `<a-modal :visible>`，**不是** `:model-value` / Use `:visible`, not `:model-value`
- 提交前必跑：`go build ./...` + `go test ./model/ ./controller/ ./middleware/` + `npm run build` / Always verify before committing

---

## 1. 项目结构 / Project Layout

```
one-api-pro/
├── cmd/                  命令行工具 / CLI tools
├── common/               公共库：config、helper、payment 等
│   ├── config/           运行时配置
│   ├── payment/          微信/支付宝/银行支付通道实现
│   └── ...
├── controller/           HTTP handler（Gin）
├── middleware/           Gin 中间件（auth、ratelimit、turnstile）
├── model/                GORM 数据模型 + 业务函数
├── relay/                LLM 渠道转发（按 provider 划分子包）
├── router/               路由注册
├── monitor/              监控/巡检
├── web/
│   ├── default-pro/      ★ 当前唯一维护的前端项目
│   ├── THEMES            主题名清单（仅含 default-pro）
│   └── build/            npm run build 产物，被 //go:embed 嵌入二进制
├── docs/                 用户文档与截图
├── main.go               程序入口
├── VERSION               当前版本号（每次发版更新）
├── CHANGELOG             变更日志
├── go.mod                Go module 声明（go 1.25.0）
├── Dockerfile            多阶段构建
└── .github/workflows/    CI：release-docker / release
```

> **⚠️ 前端仅在 `web/default-pro/` 维护。**
> 禁止在 `web/air/`、`web/default/`、`web/berry/`、`web/THEMES`、`web/build/` 改动业务代码：
> - `web/air` `web/default` `web/berry` 是历史主题保留目录，已废弃。
> - `web/THEMES` 仅是构建脚本读取的主题名清单。
> - `web/build/<theme>/` 是 `npm run build` 产物，每次构建被覆盖，禁止手改。

> **⚠️ Frontend is maintained in `web/default-pro/` ONLY.**
> Do NOT edit business code under `web/air/`, `web/default/`, `web/berry/`, `web/THEMES`, or `web/build/`:
> - `web/air` `web/default` `web/berry` are historical theme directories, deprecated.
> - `web/THEMES` is just a theme name list read by the build script.
> - `web/build/<theme>/` is the output of `npm run build` and is overwritten on every build.

---

## 2. 开发环境 / Dev Environment

| 项目 / Item | 版本 / Version |
|---|---|
| Go | 1.25.0（`go.mod`） |
| Node | 22+ |
| 包管理 / Pkg manager | pnpm 9（`pnpm-lock.yaml`） |
| 前端框架 / Frontend | Vue 3 `<script setup>` + Vite 8 |
| UI 库 / UI kit | Arco Design Vue 2.58 |
| 状态 / State | Pinia |
| HTTP | axios |
| 路由 / Router | vue-router 4 |
| DB | GORM（MySQL / PostgreSQL / SQLite） |

### 启动 / Run

```bash
# 后端 / Backend
go run .

# 前端（开发模式，自动代理 /api /v1 → :3000）/ Frontend dev
cd web/default-pro
npm run dev   # http://localhost:3001

# 主题构建（生成 web/build/default-pro，被嵌入 Go 二进制）/ Theme build
cd web/default-pro
npm run build
```

### 测试 / Tests

```bash
# Go 测试：使用 glebarez/sqlite 内存库 + gorm.AutoMigrate
# Go tests use glebarez/sqlite in-memory + gorm.AutoMigrate
go test ./...

# 前端工具函数单测（无框架依赖，Node 自带 test runner）
# Frontend pure-function unit tests (no framework, uses Node's built-in test runner)
node web/default-pro/src/utils/*.test.mjs
```

---

## 3. 提交规范 / Commit Conventions

### 3.1 类型 / Types

`feat` / `fix` / `refactor` / `style` / `chore` / `docs` / `ci` / `build` / `perf` / `test`

### 3.2 scope 约定 / Scopes

`web` | `backend` | `model` | `controller` | `router` | `relay` | `channel` | `middleware` | `common` | `docs` | `release`

### 3.3 粒度（强约束）/ Granularity (strict)

- **一个 commit = 一个文件 或 一个功能** / One commit = one file OR one feature
- **同一文件包含多个功能**：必须拆成多个 commit，每个独立 message
  / If a single file contains multiple features, split into multiple commits
- **同一功能跨多文件**：可以拆成多次单文件 commit（按依赖顺序）
  / If one feature spans multiple files, prefer multiple single-file commits in dependency order
- 禁止"feat: 重构一坨" / "fix: 改了点东西" / "update" 等无差异描述
  / Vague messages like "feat: refactor a bunch" / "fix: changed something" are forbidden

### 3.4 Commit Body 写法（关键）/ Commit Body (critical)

**commit body 必须用中文+英文同时说明本次修改的详细内容。**
**Commit body MUST be written in both Chinese and English describing the detailed changes.**

格式 / Format：

```
<type>(<scope>): <中文摘要一句话>

<英文摘要一句话 / English one-liner>

- 详细变更 1 / Detailed change 1
- 详细变更 2 / Detailed change 2
- 详细变更 3 / Detailed change 3
```

**正例 / Good example**：

```
feat(web): 充值弹窗支持自定义金额 chip

TopupModal: add custom-amount chip revealed on click.

- 新增 CUSTOM_IDX = -1 哨兵索引表示选中「自定义」chip
  Add CUSTOM_IDX = -1 sentinel for custom-amount selection
- allow_custom=true 时，chip 列表末尾追加「自定义」chip
  Append 「自定义」 chip at end of list when allow_custom is true
- 自定义金额输入框改为 v-if="isCustomSelected"，仅选中时显示
  Custom-amount input is now v-if="isCustomSelected", only shown when selected
```

### 3.5 提交前自检 / Pre-commit Checklist

```bash
git diff --stat                # 确认改动范围
# 一文件多功能 → git add -p 分块暂存
# Multi-feature in one file → use git add -p to stage hunks separately
go build ./...                 # 后端编译
go test ./model/ ./controller/ ./middleware/  # 关键包单测
cd web/default-pro && npm run build           # 前端构建
```

---

## 4. 版本号约定 / Versioning

- 根目录 `VERSION` 文件：当前版本号（例：`0.0.10`）
  / Root `VERSION` file: current version (e.g. `0.0.10`)
- 新增/修改文件时，文件头注释的 `版本:` 取当前 `VERSION` 内容
  / When creating/modifying files, set header comment `版本:` to current `VERSION` content
- 发版流程：更新 `VERSION` → 更新 `CHANGELOG` → 打 tag
  / Release flow: bump `VERSION` → update `CHANGELOG` → tag

---

## 5. 命名约定 / Naming Conventions

### 5.1 `topup` 全局统一 / `topup` is the global name

- 内部代码命名使用 `topup`（与历史 `LogTypeTopup=1` / `AdminTopUp` / `TopUp` 保持一致）
  / Internal code uses `topup` (consistent with legacy `LogTypeTopup=1` / `AdminTopUp` / `TopUp`)
- 对外 UI 文案使用「充值」/ User-facing UI text uses "充值"
- 订单 type 常量沿用 `OrderTypeTopup = 2` / Order type constant stays `OrderTypeTopup = 2`
- 设置 key 沿用 `topup.enabled` / `topup.allow_custom` / `topup.presets` / `topup.exchange_rate`
  / Setting keys: `topup.enabled` / `topup.allow_custom` / `topup.presets` / `topup.exchange_rate`
- **禁止混用** `recharge` 命名（除非显式是新业务线）
  / Do NOT mix in `recharge` (unless explicitly a new business line)

### 5.2 其它命名 / Other Naming

- 数据库常量：PascalCase + 类别前缀（`OrderTypeTopup`、`SystemSettingKeyTopupEnabled`）
  / DB constants: PascalCase with category prefix
- API 路径：`/api/<resource>/<action>`（例：`/api/topup/order`、`/api/setting/topup`）
  / API paths: `/api/<resource>/<action>`
- 前端组件：PascalCase（`TopupModal.vue`、`Dashboard.vue`）
  / Frontend components: PascalCase
- 前端 API 模块：camelCase export（`orderApi`、`topupApi`）
  / Frontend API modules: camelCase export
- 前端工具函数：camelCase 纯函数，放 `src/utils/<name>.js`
  / Frontend utils: camelCase pure functions in `src/utils/<name>.js`

---

## 6. 注释约定 / Comment Conventions

### 6.1 文件级注释（新增文件必须）/ File-level Comments (required for new files)

新增文件首段（Go 在 `package` 前；Vue 在 `<template>` 内最前；JS/TS 在文件最顶）写 4 行：

For new files, the first block (Go: before `package`; Vue: top of `<template>`; JS/TS: very top) should be 4 lines:

```go
// <文件名> <功能一句话描述>
// <功能英文描述一行>
// 版本: vX.Y.Z
// 日期: YYYY-MM-DD
// 作者: <opencode | your name>
package model
```

```vue
<!--
  <ComponentName> <功能一句话描述>
  <One-line English description>

  版本: vX.Y.Z
  日期: YYYY-MM-DD
  作者: <opencode | your name>
-->
```

```js
// <文件名> <功能一句话描述>
// <One-line English description>
// 版本: vX.Y.Z
// 日期: YYYY-MM-DD
// 作者: <opencode | your name>
```

### 6.2 导出符号注释（必须）/ Exported Symbol Comments (required)

每个 `export function` / `export const` / `func` / `type` 在前加：

Every `export function` / `export const` / `func` / `type` must be preceded by:

```go
// <一句话说明行为 / one-line behavior description>
// 版本: vX.Y.Z
// 日期: YYYY-MM-DD
func CreateTopupOrder(in CreateTopupOrderInput) (*Order, error) { ... }
```

### 6.3 行内 / 块内注释（按需）/ Inline Comments (as needed)

- **复杂业务逻辑**：解释"为什么这样做"，不止是"做了什么"
  / Complex business logic: explain "why", not just "what"
- **关键校验**：注明来源（用户需求 / 事故复盘 / 业务约束）
  / Key validations: cite the source (requirement / post-mortem / business rule)
- **TODO**：`// TODO: <场景>（待 <里程碑>）` / `// TODO: <scenario> (waiting for <milestone>)`
- **临时占位**：`// XXX: <说明>（待替换）`

### 6.4 不写注释的情况 / When NOT to Comment

- 简单 getter/setter
- 命名自解释的常量
- 一次性临时变量

### 6.5 修改现有代码的注释更新 / Updating Comments on Edits

- **逻辑变更** → 更新 `日期:`；破坏性变更同时改 `版本:`
  / Logic change: update `日期:`; breaking change also bump `版本:`
- **小修小补**（typo / 重命名）→ 注释可不更新
  / Minor fixes (typo/rename): comments may stay
- **删除代码** → 一并删除其上方注释，不要留无主注释
  / Delete dead comments together with deleted code

---

## 7. 后端规范 / Backend Conventions

### 7.1 统一响应格式 / Response Format

```go
c.JSON(http.StatusOK, gin.H{
    "success": bool,
    "message": string,
    "data":    <任意类型>,
})
```

错误也用 200 返回 + `success:false`（项目历史约定，便于前端统一处理）
/ Errors also return 200 with `success:false` (legacy convention, easier for frontend)

### 7.2 分层 / Layering

- `model/`：纯数据 + 业务函数，**不要**引入 HTTP / Gin 依赖
  / `model/`: pure data + business logic, NO HTTP / Gin dependency
- `controller/`：只做参数解析、调用 model、组装响应
  / `controller/`: parameter parsing, model calls, response assembly only
- `middleware/`：auth / rate-limit / turnstile
- `common/payment/`：各支付通道（wechat / alipay / bank）

### 7.3 支付与订单 / Payment & Orders

- 复用 `controller/buildPayInfo` 生成支付参数（pay_url / qr_code / note）
  / Reuse `controller/buildPayInfo`
- 订单号前缀：`TB` 新订 / `UP` 差价升级 / `TP` 充值
  / Order number prefixes: `TB` new / `UP` upgrade / `TP` topup
- 通知分发：`controller/payment.go::processNotify` 按 `order.Type` 分发
  / Notify dispatch: by `order.Type` in `processNotify`

### 7.4 系统设置 / System Settings

- 通用 kv 存 `model.SystemSetting` 表，key/category 维度
  / Generic KV stored in `model.SystemSetting` table
- 类别常量：`SystemSettingCategoryPayment` / `Plan` / `Topup` / `General`
- key 命名：`"<category>.<sub>.<field>"`，例：`"topup.enabled"` / `"plan.upgrade_mode"`
- 旧字段迁移：UI 删除读写但 DB 行保留（不要主动删）
  / Legacy field migration: remove UI read/write but keep DB row (do NOT actively delete)

### 7.5 测试 / Testing

- 用 `glebarez/sqlite` 内存库 + `gorm.AutoMigrate` 建表
  / Use `glebarez/sqlite` in-memory + `gorm.AutoMigrate`
- 测试函数命名：`TestXxx_Yyy`（行为_场景）
  / Test names: `TestXxx_Yyy`
- 关键路径（设置保存、金额计算、订单激活幂等）必须覆盖
  / Must cover: setting save, amount calculation, order activation idempotency

---

## 8. 前端规范 / Frontend Conventions

### 8.1 基础 / Basics

- Vue 3 `<script setup>` + Composition API
- 路径别名 `@` → `src/`（`vite.config.js`）
- 组件命名：PascalCase，多词（不要用 `Index.vue` / `Modal.vue` 这类单泛词）
  / Component names: PascalCase multi-word

### 8.2 API 模块 / API Modules

位置 / Location: `src/api/<name>.js`

```js
import api from './index'

export const orderApi = {
  myOrders: (type) => api.get('/api/order/self', { params: type ? { type } : {} }),
}

export default orderApi
```

调用约定 / Call convention:

```js
const { data } = await topupApi.createOrder(payload)
// data === { success, message, data: <实际数据> }
if (data.success && data.data) { /* ... */ }
```

### 8.3 弹窗（必踩坑）/ Modals (common pitfall)

arco `<a-modal>` 的可见性属性是 `visible`（**不是** `model-value`）。
/ Arco's `<a-modal>` visibility prop is `visible` (NOT `model-value`).

**正例 / Good**:

```vue
<a-modal
  :visible="modelValue"
  @update:visible="(v) => emit('update:modelValue', v)"
  @cancel="close"
  :title="title"
  :footer="false"
>
```

**反例（弹窗永远不显示）/ Bad (modal never shows)**:

```vue
<a-modal :model-value="modelValue" @cancel="close">  <!-- ❌ -->
```

### 8.4 arco-spin 必须设置 width / arco-spin Width

```vue
<a-spin :loading="loading" style="width: 100%">  <!-- ✅ -->
```

不加会导致父容器宽度异常。不加会引发布局错乱。
/ Without it, the parent container width behaves incorrectly.

### 8.5 表单 / Forms

- 行内编辑 / 数字输入：`<a-input-number :precision="2">`
- 开关：`<a-switch v-model="x">`
- 选择器：`<a-select v-model="x"><a-option value="...">label</a-option></a-select>`
- 表单校验：函数式校验在提交前调用，弹 `Message.error`

### 8.6 工具函数与单测 / Utils & Unit Tests

- 纯逻辑（金额计算、校验）抽到 `src/utils/<name>.js`，**不要**在 Vue 组件里写复杂内联逻辑
  / Extract pure logic to `src/utils/<name>.js`; do NOT put complex logic inside Vue components
- 单测：同目录 `<name>.test.mjs`，用 Node 自带 `node:test` + `node:assert/strict`
  / Tests: same dir `<name>.test.mjs`, use Node's built-in `node:test`

```js
import { test } from 'node:test'
import assert from 'node:assert/strict'
import { validateTopupPresets } from './topup.js'

test('rejects duplicate amount', () => {
  assert.match(validateTopupPresets([{ amount: 10 }, { amount: 10 }]), /重复/)
})
```

### 8.7 权限 / Permissions

- `useAuthStore().isAdmin` / `.isRoot` 控制 UI 显隐
- 后端中间件兜底（不要只靠前端隐藏）
  / Backend middleware enforces; do NOT rely on frontend hiding only

---

## 9. 数据库与迁移 / Database & Migration

- 改 `model/*.go` 后必须保证 `AutoMigrate` 能正确建表（不破坏性）
  / After changing models, ensure `AutoMigrate` works (non-breaking)
- 不在业务代码中写破坏性 DDL（`DROP` / `TRUNCATE` 等）
  / Do NOT write destructive DDL in business code
- 复杂迁移写 `cmd/migrate_*` 一次性脚本
  / Complex migrations: one-shot scripts in `cmd/migrate_*`

---

## 10. 常见踩坑清单 / Known Pitfalls

1. **arco `<a-modal>` 用 `:visible`（不是 `:model-value`）** / arco modal uses `:visible`
2. **arco `<a-spin>` 加 `style="width: 100%"`** / Add width to arco-spin
3. **`plan.allow_topup` 已迁移到 `topup.*`**，UI 不再读写，DB 行保留
   / `plan.allow_topup` migrated to `topup.*`; UI no longer reads/writes, DB row kept
4. **不要在 Vue 组件里写"VALID_PLAN_NAMES"这类硬编码白名单**——后端 API 已支持任意套餐名
   / Do NOT hardcode name whitelists; backend supports all plan names
5. **提交前必跑** `go build ./...` + `go test ./model/ ./controller/ ./middleware/` + `npm run build`
   / Always run all three builds before committing
6. **Commit body 必须中英双语说明详细变更**（详见 §3.4）
   / Commit body must describe detailed changes in both Chinese and English
7. **粒度**：一文件一 commit；同文件多功能必须拆 commit
   / One file per commit; split multi-feature files
8. **前端目录**：仅 `web/default-pro/`，其它 web 子目录是历史主题，禁改
   / Frontend dir: `web/default-pro/` only

---

## 11. 工作流程建议 / Recommended Workflow

1. **先 Plan 模式输出方案** → 等用户确认 → 退出 plan 模式编码
   / Output plan first → wait for approval → exit plan mode to code
2. **逐文件 commit**：每个 commit 独立 message
   / Commit file by file with independent messages
3. **复杂改动先写测试**：模型层单测（Go）/ 工具函数单测（Node）
   / Write tests first for complex changes
4. **完成任务后自检**：build + test 三件套
   / Self-verify: all three build + test commands
5. **不主动 commit**——除非用户明确指示
   / Do NOT auto-commit unless explicitly told

---

> 最后更新 / Last updated: 2026-09-07
> 版本 / Version: v0.0.10
