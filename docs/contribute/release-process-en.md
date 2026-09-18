---
title: Release Process
description: "VERSION, CHANGELOG, tag and CI pipeline."
category: contribute
order: 6
---

# Release Process

> VERSION, CHANGELOG, tag and CI pipeline.
> VERSION、CHANGELOG、tag 与 CI 流程。

Releases ship via **GitHub Actions + SemVer tags**. The complete flow:

One API Pro 的发版走 **GitHub Actions + SemVer tag**。完整流程如下：

```text
1. Prepare a release branch
2. Bump VERSION
3. Write CHANGELOG/<vX.Y.Z>.md
4. PR → merge → tag
5. CI auto-builds binaries + Docker image
```

## 1. Prepare / 准备

- Code freeze: all feature PRs merged; backport bug fixes first.
  代码冻结：所有本期功能 PR 合并；bug fix 优先 backport。
- Follow [SemVer 2.0](https://semver.org/): `MAJOR.MINOR.PATCH`.
    - `MAJOR` for breaking changes to `/v1/*` or provider protocols.
    - `MINOR` for new features / providers / payment channels.
    - `PATCH` for bug fixes / docs / perf.
  版本号决定：遵循 SemVer；`MAJOR` 大改、`MINOR` 新功能、`PATCH` bug 修复。
- CHANGELOG file: `CHANGELOG/vX.Y.Z.md` (with `v` prefix).
  CHANGELOG 文件命名：`CHANGELOG/vX.Y.Z.md`（**包含 `v` 前缀**）。

## 2. Required Artifacts / 必填材料

### 2.1 `VERSION`

A single-line file at the repo root:

仓库根的 `VERSION` 文件，仅一行：

```bash
echo "0.0.22" > VERSION
```

Notes:

- No `v` prefix (the `v` only appears in tag and CHANGELOG filename).
  **不包含** `v` 前缀（`v` 仅出现在 tag 与 CHANGELOG 文件名）。
- Injected into the binary by CI via `-ldflags "-X ...common.Version=${TAG}"`.
  由 CI 编译时通过 `-ldflags` 注入二进制。

### 2.2 `CHANGELOG/vX.Y.Z.md`

Reference `CHANGELOG/_template.md` or any historical `v0.0.21.md`. Suggested structure:

格式可参考 `CHANGELOG/_template.md` 或历史 `v0.0.21.md`。建议结构：

```markdown
# vX.Y.Z — <one-line theme>

> English one-liner.
> 中文简介。

## English

### ✨ New features

- ...

### 🐛 Bug fixes

- ...

### 🔧 Refactors

- ...

### ⚠️ Upgrade notes

- Zero DB migration / column migration / behavior change
```

CI uses `CHANGELOG/${TAG}.md` as the GitHub Release body.

CI 会把 `CHANGELOG/${TAG}.md` 当作 GitHub Release body。

### 2.3 tag

```bash
git tag v0.0.22
git push origin v0.0.22
```

Tags must match `^v[0-9]+\.[0-9]+\.[0-9]+([.-].+)?$`, otherwise the CI normalize step fails.

Tag 必须严格匹配 SemVer 正则，否则 CI 在 normalize 步骤失败。

## 3. CI Pipeline / CI 流水线

The repo's `.github/workflows/`:

仓库 `.github/workflows/`：

| Workflow | Trigger / 触发 | Output / 产物 |
| --- | --- | --- |
| `release.yml` | push tag `v*` or `workflow_dispatch` | linux/amd64, linux/arm64, windows/amd64, darwin/amd64, darwin/arm64 binaries; GitHub Release |
| `release-docker.yml` | push tag `v*.*.*` or `workflow_dispatch` | `ghcr.io/<owner>/one-api-pro:<tag>` multi-arch image (amd64 / arm64) |

Key steps / 主要步骤:

1. **Normalize tag**: auto-prefix `v`; validate semver.
3. **Verify CHANGELOG exists**: `CHANGELOG/${TAG}.md` must exist; otherwise CI refuses to build.
   ```bash
   SRC="CHANGELOG/${TAG}.md"
   if [ ! -f "$SRC" ]; then echo "CHANGELOG file not found: $SRC"; exit 1; fi
   ```
4. **Build frontend**: `cd web && sh build.sh`, the output is `//go:embed`-ed.
   构建前端，产物被 `//go:embed` 进二进制。
5. **Build binaries**:
   ```bash
   LDFLAGS="-s -w -X github.com/modelbus/one-api-pro/common.Version=${TAG}"
   CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o dist/one-api-pro-linux-amd64 .
   ```
6. **Release**: use `softprops/action-gh-release@v2`, body = `CHANGELOG/${TAG}.md`.
7. **Docker**: multi-stage `Dockerfile` → `docker buildx` push to `ghcr.io`.

## 4. Verification / 校验清单

After CI passes:

CI 通过后：

1. Visit `https://github.com/<owner>/one-api-pro/releases/tag/vX.Y.Z`, confirm body + assets.
   访问 GitHub Releases 页面，确认 body 与附件齐全。
2. Pull the image and smoke-test:
   拉镜像测试：
   ```bash
   docker run --rm -p 3000:3000 \
     -e TZ=Asia/Shanghai \
     -v $(pwd)/data:/app/data \
     ghcr.io/<owner>/one-api-pro:vX.Y.Z \
     --version
   ```
3. Run the diagnostics from [Troubleshooting §6](/en/faq/troubleshooting#collecting-diagnostics).
   跑一遍 [故障排查 §6](/zh/faq/troubleshooting#收集诊断信息) 提到的诊断命令。

## 5. Doc Sync / 文档同步

Right after release:

发版后立即：

1. Copy `CHANGELOG/vX.Y.Z.md` to `docs/changelog/vX.Y.Z-{zh,en}.md` (the doc site keeps its own copy).
   把 `CHANGELOG/vX.Y.Z.md` 复制到 `docs/changelog/vX.Y.Z-{zh,en}.md`。
2. Update the "Latest version" / demo images / key links in `README.md` as needed.
   必要时更新 `README.md` 中的「最新版本」「演示图」「关键链接」。
3. If breaking: add an "Upgrade banner" at the top of `docs/start/overview.md` / `docs/install/upgrade.md`.
   如果有破坏性变更：在 `docs/start/overview.md` / `docs/install/upgrade.md` 顶部加「升级提示」Banner。
4. Run `node docs/scripts/sync-docs.mjs` (maintained manually; CI doesn't touch the doc site).
   运行 `node docs/scripts/sync-docs.mjs`（仓库侧维护人手动；CI 不动文档站）。

## 6. Hotfix Flow / 热修流程

```text
1. Branch from the previous tag as hotfix/<bug>
2. Fix + add regression test
3. Bump VERSION (PATCH + 1) and write CHANGELOG
4. PR → merge → tag vX.Y.(Z+1)
5. CI auto-releases
```

```text
1. 从上一个 release tag 拉 hotfix/<bug> 分支
2. 修复 + 增加回归测试
3. 更新 VERSION (PATCH + 1) + 写 CHANGELOG
4. PR → merge → tag vX.Y.(Z+1)
5. CI 自动发版
```

## 7. Pre-release / 预发布

```bash
git tag v0.0.22-rc.1
git push origin v0.0.22-rc.1
```

`latest` is skipped when the tag includes a pre-release suffix (`-rc` / `-alpha` / `-beta`); otherwise the same flow applies.

CI 会：跳过 `latest` tag（避免把 RC 推到 latest）；走同一个 release 流程，附件里追加 `-rc.N` 后缀。

## 8. Pitfalls / 注意事项

- Don't amend published tags; cut a new `v0.0.22-fixed` and migrate users.
  不要修改已发布的 tag（即使只是一行 typo）；重新打 `v0.0.22-fixed` 让用户迁移。
- CI doesn't check CHANGELOG vs VERSION consistency; double-check manually.
  CHANGELOG 与 VERSION 不一致：CI 不会检查两者是否对齐；人工 review 时务必确认。
- If CHANGELOG is missing after the tag push: commit the missing file (don't re-tag), then run `workflow_dispatch` with the same tag.
  先 push tag 后再发现 CHANGELOG 缺失：可以补救：在 `CHANGELOG/vX.Y.Z.md` 写好内容，重新 push commit（不改 tag 哈希），CI 不会重跑；这时用 `workflow_dispatch` 手动触发一次即可。

Next: [Commit Convention](/en/contribute/commit-convention) · [Troubleshooting](/en/faq/troubleshooting).