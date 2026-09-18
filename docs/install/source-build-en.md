---
title: Source Build
description: "Build the Go backend and web theme from source."
category: install
order: 4
---

# Source Build

> Build the Go backend and web theme from source.

## Prerequisites

| Tool | Version | Notes |
| --- | --- | --- |
| Go | 1.25.0 | Declared in `go.mod`; auto-fetched via `GOTOOLCHAIN=auto` |
| Node.js | 22+ | Required only to build the front-end |
| pnpm | 9 | Recommended; the repo ships `pnpm-lock.yaml` |

> The root `VERSION` file is auto-detected by the binary (`--version` and `/api/status` both use it).

## Order of operations (critical)

The front-end bundle is embedded via `//go:embed web/build/default-pro`; you **must build the front-end first**.

### 1. Build the front-end

```bash
cd web
sh build.sh
```

`build.sh` reads the theme list in `web/THEMES`, enters each theme directory and runs `npm install && npm run build`. The only actively maintained theme today is `default-pro` (Vue 3 + Arco Design).

Output:

```text
web/build/default-pro/   # static assets, embedded via //go:embed
```

### 2. Build the Go backend

```bash
cd ..
go build -ldflags "-s -w" -o one-api-pro
```

`-s -w` strips the symbol table and DWARF debug info, yielding a ~40 MB binary.

### 3. Run

```bash
chmod u+x one-api-pro
./one-api-pro --port 3000 --log-dir ./logs
```

Open [http://localhost:3000/](http://localhost:3000/) and log in with the default account `root` / `123456`.

## Multi-platform release

The repo ships `release.sh`, which builds `linux/amd64`, `linux/arm64`, `windows/amd64`, `darwin/amd64` and `darwin/arm64` in one shot.

```bash
./release.sh                # version from VERSION file or git tag
./release.sh v0.0.22        # explicit version
./release.sh v0.0.22 --skip-frontend   # reuse existing web/build, skip the front-end
```

Output:

```text
dist/
  one-api-pro-linux-amd64
  one-api-pro-linux-arm64
  one-api-pro-windows-amd64.exe
  one-api-pro-darwin-amd64
  one-api-pro-darwin-arm64
```

`release.sh` reads `VERSION`, runs `go mod download`, runs `(cd web && sh build.sh)` and then cross-compiles the 5 targets with `CGO_ENABLED=0`.

## Front-end only

```bash
cd web/default-pro
npm install
npm run build
```

The bundle lands in `web/build/default-pro/`. To clean:

```bash
rm -rf web/build/default-pro
```

## Verification

After a successful build, run a minimal self-check:

```bash
./one-api-pro --version
./one-api-pro --help | head -20
./one-api-pro --port 3001 --log-dir ./logs &
sleep 3
curl -fsS http://localhost:3001/api/status
```

`/api/status` should return JSON including `version`, `build_time` and similar fields (exact schema follows the source).

Next: [Configuration](/en/install/config) · [Docker Deploy](/en/install/docker-deploy).
