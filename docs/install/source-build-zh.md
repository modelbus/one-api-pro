---
title: 源码编译
description: "本地编译 Go 后端与前端主题。"
category: install
order: 4
---

# 源码编译

> 本地编译 Go 后端与前端主题。
> Build the Go backend and web theme from source.

## 先决条件 / Prerequisites

| 工具 / Tool | 版本 / Version | 说明 / Notes |
| --- | --- | --- |
| Go | 1.25.0 | `go.mod` 中声明；GOTOOLCHAIN=auto 自动拉取 |
| Node.js | 22+ | 仅构建前端时需要 |
| pnpm | 9 | 推荐；仓库提交 `pnpm-lock.yaml` |

> 仓库根目录的 `VERSION` 文件会被二进制自动读取（`--version` 与 `/api/status` 都用此值）。
> The root `VERSION` file is auto-detected by the binary (`--version` and `/api/status` both use it).

## 构建顺序（关键） / Order of operations (critical)

前端产物通过 `//go:embed web/build/default-pro` 嵌入二进制；**必须先构建前端，再构建后端**。

The front-end bundle is embedded via `//go:embed web/build/default-pro`; you **must build the front-end first**.

### 1. 构建前端 / Build the front-end

```bash
cd web
sh build.sh
```

`build.sh` 读取 `web/THEMES` 中的主题名清单，逐个进入主题目录执行 `npm install && npm run build`。当前唯一维护主题是 `default-pro`（Vue 3 + Arco Design）。

`build.sh` reads the theme list in `web/THEMES`, enters each theme directory and runs `npm install && npm run build`. The only actively maintained theme today is `default-pro` (Vue 3 + Arco Design).

构建完成后产物落在 / Output:

```text
web/build/default-pro/   # 静态资源，会被 //go:embed 嵌入二进制
```

### 2. 构建 Go 后端 / Build the Go backend

```bash
cd ..
go build -ldflags "-s -w" -o one-api-pro
```

`-s -w` 去除符号表与调试信息，最终二进制约 40 MB。
`-s -w` strips the symbol table and DWARF debug info, yielding a ~40 MB binary.

### 3. 启动 / Run

```bash
chmod u+x one-api-pro
./one-api-pro --port 3000 --log-dir ./logs
```

打开 [http://localhost:3000/](http://localhost:3000/) ，默认账号 `root` / `123456`。
Open [http://localhost:3000/](http://localhost:3000/) and log in with the default account `root` / `123456`.

## 多平台发布 / Multi-platform release

仓库自带 `release.sh`，可一次性构建 `linux/amd64`、`linux/arm64`、`windows/amd64`、`darwin/amd64`、`darwin/arm64`。

The repo ships `release.sh`, which builds `linux/amd64`, `linux/arm64`, `windows/amd64`, `darwin/amd64` and `darwin/arm64` in one shot.

```bash
./release.sh                # 版本号取 VERSION 文件或 git tag
./release.sh v0.0.22        # 显式指定
./release.sh v0.0.22 --skip-frontend   # 复用已有 web/build，跳过前端构建
```

产物 / Output:

```text
dist/
  one-api-pro-linux-amd64
  one-api-pro-linux-arm64
  one-api-pro-windows-amd64.exe
  one-api-pro-darwin-amd64
  one-api-pro-darwin-arm64
```

`release.sh` 内部依次：读取 `VERSION` → `go mod download` → `(cd web && sh build.sh)` → 5 个目标交叉编译（`CGO_ENABLED=0`）。
`release.sh` reads `VERSION`, runs `go mod download`, runs `(cd web && sh build.sh)` and then cross-compiles the 5 targets with `CGO_ENABLED=0`.

## 前端单独构建 / Front-end only

```bash
cd web/default-pro
npm install
npm run build
```

构建产物输出到 `web/build/default-pro/`；清理旧产物：

The bundle lands in `web/build/default-pro/`. To clean:

```bash
rm -rf web/build/default-pro
```

## 校验 / Verification

构建完成后建议运行最小自检：

After a successful build, run a minimal self-check:

```bash
./one-api-pro --version
./one-api-pro --help | head -20
./one-api-pro --port 3001 --log-dir ./logs &
sleep 3
curl -fsS http://localhost:3001/api/status
```

`/api/status` 应返回 JSON，包含 `version`、`build_time` 等字段（具体字段以源码为准）。

`/api/status` should return JSON including `version`, `build_time` and similar fields (exact schema follows the source).

下一步 / Next: [配置项与环境变量](/zh/install/config) · [Docker 部署](/zh/install/docker-deploy)。
