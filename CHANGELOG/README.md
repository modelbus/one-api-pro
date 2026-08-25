# 发布日志目录（CHANGELOG）

本目录按 **版本号（tag）** 维护每次发布的 Markdown 发布日志，文件名与 Git tag 完全一致，例如：

```
CHANGELOG/v1.2.0.md
CHANGELOG/v1.1.0.md
```

CI（`.github/workflows/release.yml`）会自动读取与当前 tag 同名的 Markdown 文件作为 **GitHub Release 的发布说明（body）**。如果该文件不存在，则回退到 GitHub 自动生成的 release notes。

---

## 文件约定

每个发布日志文件应同时提供 **中文** 与 **英文** 两段内容，便于不同语种的用户查看。建议从 `_template.md` 复制后填写：

- **中文段**：使用主要标题（`## 中文`），讲述面向中文用户的功能/修复/变更亮点。
- **English section**：使用 `## English` 标题，提供等价英文版本。
- 中英两段建议保持条目数量与顺序一致，便于跨语言查阅。

---

## 如何新增一个版本的发布日志

由 Release Manager 在打 tag **之前** 完成（推荐在 `master` 分支上操作）：

1. 找到上一个 tag，确认本次版本号（例如 `v1.3.0`）。
2. 复制模板：
   ```bash
   cp CHANGELOG/_template.md CHANGELOG/v1.3.0.md
   ```
3. 收集自上一个 tag 以来的 commit：
   ```bash
   git log v1.2.0..HEAD --oneline
   ```
4. 用 AI agent 解析 commit messages，提取每个 commit 的语义、归类为下列章节：
   - ✨ **新增功能 / New Features**
   - 🐛 **问题修复 / Bug Fixes**
   - 🚀 **性能优化 / Performance**
   - 🔧 **重构 / Refactor**
   - 📚 **文档 / Documentation**
   - 🔒 **安全 / Security**
5. 撰写中文段落、英文段落，并补充 **升级注意事项 / Upgrade Notes**。
6. 提交 PR：
   ```
   docs(changelog): add v1.3.0 release notes
   ```
   合并后，由维护者推送 tag 触发发布。

### 让 AI agent 自动生成草稿（推荐流程）

> 调用你团队使用的 AI agent，给它以下 prompt：

```
请基于以下 git commit messages，生成一个 CHANGELOG/v<NEW_TAG>.md 文件。
要求：
1. 同时输出中文与英文两个段（标题分别为 "## 中文" 与 "## English"）。
2. 按照以下章节分类：新增功能、问题修复、性能优化、重构、文档、安全。
3. 每个 bullet 用动宾结构描述用户视角的价值，避免直接照抄 commit subject。
4. 在结尾追加 "升级注意事项 / Upgrade Notes" 段，说明需要用户操作（如有）。
5. 最后使用这里给的 markdown 模板：
   ...（把 _template.md 内容贴给 agent）
```

把 agent 输出复制进 `CHANGELOG/v<NEW_TAG>.md`，人工 review 后提 PR。

---

## CI 行为说明

- 触发条件：推送形如 `v*` 的 tag，或手动触发 workflow（`workflow_dispatch`）。
- 发布脚本会读取 `CHANGELOG/${TAG}.md`：
  - 文件存在 → 内容作为 release body；
  - 文件不存在 → 使用 `generate_release_notes: true`，由 GitHub 自动按 PR 列表生成。

> 💡 如果你不希望 CI 退回到自动生成，可以把 `release.yml` 里的 `generate_release_notes: true` 改掉。
