---
title: 提交规范
description: "Conventional Commits 与中英双语 commit body。"
category: contribute
order: 3
---

# 提交规范

> Conventional Commits 与中英双语 commit body。
> Conventional Commits with bilingual commit body.

## 1. 类型 / Types

`feat` / `fix` / `refactor` / `style` / `chore` / `docs` / `ci` / `build` / `perf` / `test`

| 类型 / Type | 场景 / Scenario | 例 / Example |
| --- | --- | --- |
| `feat` | 新功能 / New feature | 新增充值弹窗 |
| `fix` | 修复 bug / Bug fix | 修复用户配额缓存漂移 |
| `refactor` | 既不改功能也不改 bug / Pure restructuring | 抽离 i18n 工具函数 |
| `style` | 格式 / 空白 / Style only | prettier / gofmt |
| `chore` | 构建 / 依赖 / 杂项 / Build or misc | bump dependency |
| `docs` | 文档 / Documentation | update docs/ |
| `ci` | CI 配置 / CI config | add GitHub Actions step |
| `build` | 构建系统 / Build system | update Dockerfile |
| `perf` | 性能 / Performance | cache hot-path |
| `test` | 测试 / Tests | add regression test |

## 2. scope 约定 / Scopes

`web` | `backend` | `model` | `controller` | `router` | `relay` | `channel` | `middleware` | `common` | `docs` | `release`

例：`feat(web): 新增充值弹窗`、`fix(relay): 修复 OpenAI 适配器 timeout`。

Examples: `feat(web): add topup modal`, `fix(relay): fix OpenAI adaptor timeout`.

## 3. 粒度（强约束）/ Granularity (strict)

- **一个 commit = 一个文件 或 一个功能**。
  One commit = one file OR one feature.
- **同一文件包含多个功能**：必须拆成多个 commit，每个独立 message。
  If a single file contains multiple features, split into multiple commits.
- **同一功能跨多文件**：可以拆成多次单文件 commit（按依赖顺序）。
  If one feature spans multiple files, prefer multiple single-file commits in dependency order.
- **禁止**：「feat: 重构一坨」「fix: 改了点东西」「update」等无差异描述。
  Forbidden: "feat: refactor a bunch" / "fix: changed something" / "update".

## 4. Commit Body 写法（关键）/ Commit Body (critical)

**commit body 必须用中文+英文同时说明本次修改的详细内容。**

**Commit body MUST be written in both Chinese and English describing the detailed changes.**

格式 / Format：

```text
<type>(<scope>): <中文摘要一句话>

<英文摘要一句话 / English one-liner>

- 详细变更 1 / Detailed change 1
- 详细变更 2 / Detailed change 2
- 详细变更 3 / Detailed change 3
```

### 完整示例 / Full Example

```text
feat(web): 充值弹窗支持自定义金额 chip

TopupModal: add custom-amount chip revealed on click.

- 新增 CUSTOM_IDX = -1 哨兵索引表示选中「自定义」chip
  Add CUSTOM_IDX = -1 sentinel for custom-amount selection
- allow_custom=true 时，chip 列表末尾追加「自定义」chip
  Append 「自定义」 chip at end of list when allow_custom is true
- 自定义金额输入框改为 v-if="isCustomSelected"，仅选中时显示
  Custom-amount input is now v-if="isCustomSelected", only shown when selected
```

## 5. 提交前自检 / Pre-commit Checklist

```bash
git diff --stat                # 确认改动范围
                              # 一文件多功能 → git add -p 分块暂存
                              # Multi-feature in one file → use git add -p to stage hunks separately
go build ./...                 # 后端编译
go test ./model/ ./controller/ ./middleware/  # 关键包单测
cd web/default-pro && pnpm build           # 前端构建
```

## 6. 分块暂存 / Stage in Hunks

```bash
# 查看 hunk
git add -p file.go

# 选项
y - 暂存此 hunk
n - 不暂存此 hunk
s - 拆成更小的 hunk
e - 手动编辑
```

## 7. 不要做的事 / What NOT to Do

| 反例 / Anti-pattern | 正确做法 / Correct |
| --- | --- |
| `feat: 重构一坨` | 按文件 / 功能拆分 |
| `fix: 改了点东西` | 说明症状 + 根因 + 修复 |
| 一个 commit 包含 5 个文件 + 3 个功能 | 拆成 3+ 个 commit |
| body 全中文或全英文 | 中英双语 |
| `update` / `wip` / `tmp` | 走完自检后再提交 |
| 提交后 force push 已被别人拉的分支 | 等合并完成；或与 reviewer 约定 |

## 8. PR 标题与描述 / PR Title & Description

PR 标题同 commit title；描述补充：

PR title mirrors the commit title; description adds:

1. **What**：本次变更做了什么（1–2 句）。
   **What**: what the change does (1–2 sentences).
2. **Why**：为什么要做（背景、需求、事故复盘）。
   **Why**: why (background, requirement, post-mortem).
3. **How**：实现要点（接口 / 关键函数）。
   **How**: implementation highlights (interfaces / key functions).
4. **Test**：是否新增测试、是否跑过 `pnpm build`。
   **Test**: new tests added? ran `pnpm build`?
5. **Risk**：潜在影响 / 回滚方案。
   **Risk**: impact / rollback plan.

## 9. 工作流程建议 / Recommended Workflow

1. **先 Plan**：对复杂改动输出方案 → 等用户确认 → 退出 Plan 模式编码。
   Plan first: for complex changes, output the plan → wait for approval → exit plan mode to code.
2. **逐文件 commit**：每个 commit 独立 message。
   Commit file by file with independent messages.
3. **复杂改动先写测试**：模型层单测（Go）/ 工具函数单测（Node）。
   Write tests first for complex changes: model-layer tests (Go) / util tests (Node).
4. **完成后自检**：build + test 三件套。
   Self-verify: the build + test trio.
5. **不主动 commit**：除非用户明确指示。
   Do NOT auto-commit unless explicitly told.

下一步 / Next: [新增 Provider](/zh/contribute/add-provider) · [发版流程](/zh/contribute/release-process)。