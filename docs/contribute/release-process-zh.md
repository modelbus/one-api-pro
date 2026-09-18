---
title: 发版流程
description: "VERSION、CHANGELOG、tag 与 CI 流程。"
category: contribute
order: 6
---

# 发版流程

> VERSION、CHANGELOG、tag 与 CI 流程。

One API Pro 的发版走 **GitHub Actions + SemVer tag**。完整流程如下：

```text
1. 准备 release 分支
   Prepare a release branch
2. 更新 VERSION 文件
   Bump VERSION
3. 写 CHANGELOG/<vX.Y.Z>.md
   Write CHANGELOG/<vX.Y.Z>.md
4. 提交 PR → 合并 → 打 tag
   PR → merge → tag
5. CI 自动构建二进制 + Docker 镜像
   CI auto-builds binaries + Docker image
```

## 1. 准备

- **代码冻结**：所有本期功能 PR 合并；bug fix 优先 backport。
- **版本号决定**：遵循 [SemVer 2.0](https://semver.org/)：
    - `MAJOR.MINOR.PATCH` 三位；
    - `MAJOR` 大改：`/v1/*` 协议破坏 / Provider 协议破坏；
    - `MINOR` 新功能 / 新 Provider / 新支付通道；
    - `PATCH` bug 修复 / 文档 / 性能。
- **CHANGELOG 文件命名**：`CHANGELOG/vX.Y.Z.md`（**包含 `v` 前缀**）。

## 2. 必填材料

### 2.1 `VERSION`

仓库根的 `VERSION` 文件，仅一行：

```bash
echo "0.0.22" > VERSION
```

**注意**：

- **不包含** `v` 前缀（`v` 仅出现在 tag 与 CHANGELOG 文件名）。
- 由 CI 编译时通过 `-ldflags "-X ...common.Version=${TAG}"` 注入二进制。

### 2.2 `CHANGELOG/vX.Y.Z.md`

格式可参考 `CHANGELOG/_template.md` 或历史 `v0.0.21.md`。建议结构：

```markdown
# vX.Y.Z — <一句话主题>

> 中文简介。

## 中文

### ✨ 新增功能

- ...

### 🐛 问题修复

- ...

### 🔧 重构

- ...

### ⚠️ 升级注意事项

- 零数据库迁移 / 字段迁移 / 行为变更
```

CI 会把 `CHANGELOG/${TAG}.md` 当作 GitHub Release body。

### 2.3 tag

```bash
git tag v0.0.22
git push origin v0.0.22
```

Tag 必须严格匹配 `^v[0-9]+\.[0-9]+\.[0-9]+([.-].+)?$`，否则 CI 在 normalize 步骤失败。

## 3. CI 流水线

仓库 `.github/workflows/`：

| Workflow | 触发| 产物|
| --- | --- | --- |
| `release.yml` | push tag `v*` 或 `workflow_dispatch` | linux/amd64、linux/arm64、windows/amd64、darwin/amd64、darwin/arm64 二进制；GitHub Release |
| `release-docker.yml` | push tag `v*.*.*` 或 `workflow_dispatch` | `ghcr.io/<owner>/one-api-pro:<tag>` 多架构镜像（amd64|

主要步骤：

1. **Normalize tag**：自动补 `v` 前缀；校验 semver。
3. **Verify CHANGELOG exists**：`CHANGELOG/${TAG}.md` 必须存在；缺失则 CI 拒绝构建。
   SRC="CHANGELOG/${TAG}.md"
   if [ ! -f "$SRC" ]; then echo "CHANGELOG file not found: $SRC"; exit 1; fi
   ```
4. **Build frontend**：`cd web && sh build.sh`，产物被 `//go:embed` 进二进制。
5. **Build binaries**：
   ```bash
   LDFLAGS="-s -w -X github.com/modelbus/one-api-pro/common.Version=${TAG}"
   CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o dist/one-api-pro-linux-amd64 .
   ```
6. **Release**：通过 `softprops/action-gh-release@v2` 上传二进制，body 取 `CHANGELOG/${TAG}.md`。
7. **Docker**：多阶段 `Dockerfile` → `docker buildx` 推 `ghcr.io`。

## 4. 校验清单

CI 通过后：

1. 访问 `https://github.com/<owner>/one-api-pro/releases/tag/vX.Y.Z`，确认 body 与附件齐全。
2. 拉镜像测试：
   ```bash
   docker run --rm -p 3000:3000 \
     -e TZ=Asia/Shanghai \
     -v $(pwd)/data:/app/data \
     ghcr.io/<owner>/one-api-pro:vX.Y.Z \
     --version
   ```
3. 跑一遍 [故障排查 §6](/zh/faq/troubleshooting#收集诊断信息) 提到的诊断命令，确认无回归。

## 5. 文档同步

发版后立即：

1. 把 `CHANGELOG/vX.Y.Z.md` 复制到 `docs/changelog/vX.Y.Z-{zh,en}.md`（仓库侧文档站独立副本）。
2. 必要时更新 `README.md` 中的「最新版本」「演示图」「关键链接」。
3. 如果有破坏性变更：在 `docs/start/overview.md` / `docs/install/upgrade.md` 顶部加「升级提示」Banner。
4. 运行 `node docs/scripts/sync-docs.mjs`（仓库侧维护人手动；CI 不动文档站）。

## 6. 热修流程

```text
1. 从上一个 release tag 拉 hotfix/<bug> 分支
   Branch from the previous tag as hotfix/<bug>
2. 修复 + 增加回归测试
   Fix + add regression test
3. 更新 VERSION (PATCH + 1) + 写 CHANGELOG
   Bump VERSION (PATCH + 1) and write CHANGELOG
4. PR → merge → tag vX.Y.(Z+1)
   PR → merge → tag vX.Y.(Z+1)
5. CI 自动发版
   CI auto-releases
```

## 7. Pre-release

```bash
git tag v0.0.22-rc.1
git push origin v0.0.22-rc.1
```

CI 会：

- 跳过 `latest` tag（避免把 RC 推到 latest）；
- 走同一个 release 流程，附件里追加 `-rc.N` 后缀。

`latest` is skipped when the tag includes a pre-release suffix (`-rc` / `-alpha` / `-beta`); otherwise the same flow applies.

## 8. 注意事项

- **不要修改已发布的 tag**（即使只是一行 typo）；重新打 `v0.0.22-fixed` 让用户迁移。
- **CHANGELOG 与 VERSION 不一致**：CI 不会检查两者是否对齐；人工 review 时务必确认。
- **先 push tag 后再发现 CHANGELOG 缺失**：可以补救：在 `CHANGELOG/vX.Y.Z.md` 写好内容，重新 push commit（不改 tag 哈希），CI 不会重跑；这时用 `workflow_dispatch` 手动触发一次即可。

下一步 / Next: [提交规范](/zh/contribute/commit-convention) · [故障排查](/zh/faq/troubleshooting)。
