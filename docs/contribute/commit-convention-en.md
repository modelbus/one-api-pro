---
title: Commit Convention
description: "Conventional Commits with bilingual commit body."
category: contribute
order: 3
---

# Commit Convention

> Conventional Commits with bilingual commit body.
> Conventional Commits 与中英双语 commit body。

## 1. Types / 类型

`feat` / `fix` / `refactor` / `style` / `chore` / `docs` / `ci` / `build` / `perf` / `test`

| Type | Scenario / 场景 | Example |
| --- | --- | --- |
| `feat` | New feature / 新功能 | Add topup modal |
| `fix` | Bug fix / 修复 bug | Fix user-quota cache drift |
| `refactor` | Pure restructuring / 既不改功能也不改 bug | Extract i18n utils |
| `style` | Format / whitespace only / 格式 / 空白 | prettier / gofmt |
| `chore` | Build / deps / misc / 构建 / 依赖 / 杂项 | bump dependency |
| `docs` | Documentation / 文档 | update docs/ |
| `ci` | CI config / CI 配置 | add GitHub Actions step |
| `build` | Build system / 构建系统 | update Dockerfile |
| `perf` | Performance / 性能 | cache hot-path |
| `test` | Tests / 测试 | add regression test |

## 2. Scopes / scope 约定

`web` | `backend` | `model` | `controller` | `router` | `relay` | `channel` | `middleware` | `common` | `docs` | `release`

Examples: `feat(web): add topup modal`, `fix(relay): fix OpenAI adaptor timeout`.

例：`feat(web): 新增充值弹窗`、`fix(relay): 修复 OpenAI 适配器 timeout`。

## 3. Granularity (strict) / 粒度（强约束）

- **One commit = one file OR one feature**.
- If a single file contains multiple features, split into multiple commits.
- If one feature spans multiple files, prefer multiple single-file commits in dependency order.
- **Forbidden**: "feat: refactor a bunch" / "fix: changed something" / "update".
  禁止："feat: 重构一坨" / "fix: 改了点东西" / "update".

## 4. Commit Body (critical) / Commit Body 写法（关键）

**Commit body MUST be written in both Chinese and English describing the detailed changes.**

**commit body 必须用中文+英文同时说明本次修改的详细内容。**

Format / 格式:

```text
<type>(<scope>): <中文摘要一句话>

<English one-liner>

- 详细变更 1 / Detailed change 1
- 详细变更 2 / Detailed change 2
- 详细变更 3 / Detailed change 3
```

### Full Example / 完整示例

```text
feat(web): TopupModal supports custom-amount chip

TopupModal: add custom-amount chip revealed on click.

- Add CUSTOM_IDX = -1 sentinel for custom-amount selection
  新增 CUSTOM_IDX = -1 哨兵索引表示选中「自定义」chip
- Append 「自定义」 chip at end of list when allow_custom is true
  allow_custom=true 时，chip 列表末尾追加「自定义」chip
- Custom-amount input is now v-if="isCustomSelected", only shown when selected
  自定义金额输入框改为 v-if="isCustomSelected"，仅选中时显示
```

## 5. Pre-commit Checklist / 提交前自检

```bash
git diff --stat                # review the diff
                              # Multi-feature in one file → use git add -p to stage hunks separately
                              # 一文件多功能 → git add -p 分块暂存
go build ./...                 # backend build
go test ./model/ ./controller/ ./middleware/  # key-package tests
cd web/default-pro && pnpm build           # frontend build
```

## 6. Stage in Hunks / 分块暂存

```bash
# Review hunks
git add -p file.go

# Options
y - stage this hunk
n - skip this hunk
s - split into smaller hunks
e - edit manually
```

## 7. What NOT to Do / 不要做的事

| Anti-pattern / 反例 | Correct / 正确做法 |
| --- | --- |
| `feat: refactor a bunch` | Split per file / feature |
| `fix: changed something` | State symptom + root cause + fix |
| One commit with 5 files + 3 features | Split into 3+ commits |
| Body in one language only | Bilingual zh + en |
| `update` / `wip` / `tmp` | Finish self-check first |
| Force-push a branch others pulled from | Wait for merge; coordinate with reviewer |

## 8. PR Title & Description / PR 标题与描述

PR title mirrors the commit title; description adds:

PR 标题同 commit title；描述补充：

1. **What**: what the change does (1–2 sentences).
   本次变更做了什么（1–2 句）。
2. **Why**: why (background, requirement, post-mortem).
   为什么要做（背景、需求、事故复盘）。
3. **How**: implementation highlights (interfaces / key functions).
   实现要点（接口 / 关键函数）。
4. **Test**: new tests added? ran `pnpm build`?
   是否新增测试、是否跑过 `pnpm build`。
5. **Risk**: impact / rollback plan.
   潜在影响 / 回滚方案。

## 9. Recommended Workflow / 工作流程建议

1. **Plan first**: for complex changes, output the plan → wait for approval → exit plan mode to code.
   先 Plan：对复杂改动输出方案 → 等用户确认 → 退出 Plan 模式编码。
2. **Commit file by file** with independent messages.
   逐文件 commit：每个 commit 独立 message。
3. **Write tests first for complex changes**: model-layer tests (Go) / util tests (Node).
   复杂改动先写测试：模型层单测（Go）/ 工具函数单测（Node）。
4. **Self-verify**: the build + test trio.
   完成后自检：build + test 三件套。
5. **Do NOT auto-commit** unless explicitly told.
   不主动 commit：除非用户明确指示。

Next: [Add a Provider](/en/contribute/add-provider) · [Release Process](/en/contribute/release-process).