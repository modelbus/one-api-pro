---
title: Source Build
description: Build the Go backend and the web theme from source.
category: install
order: 4
---

# Source Build

> For hacking on it, contributing, or debugging a PR.

## Prerequisites

| Tool | Version |
|---|---|
| Go | 1.25+ |
| Node.js | 22+ |
| pnpm | 9+ |
| Git | any recent |

## Backend

```bash
git clone https://github.com/modelbus/one-api-pro
cd one-api-pro
go build -o one-api-pro .
```

Run:

```bash
./one-api-pro
```

First start auto-creates the SQLite DB and prints the default admin.

## Frontend

```bash
cd web/default-pro
pnpm install
pnpm build       # output is embedded into the binary via //go:embed
```

If you skip the frontend build, pages show "frontend assets not embedded".

## Dev loop

Backend changes: just restart `./one-api-pro`.

Frontend changes:

```bash
cd web/default-pro
pnpm build
# rebuild the backend (because of //go:embed)
cd ../..
go build -o one-api-pro .
```

For frequent UI changes, run the frontend dev server with a proxy to backend:3000.

## Editor

GoLand or VS Code + Go extension. Code style: [Code Style](../contribute/code-style).

## Related

- [Dev Setup](../contribute/dev-setup)
- [Code Style](../contribute/code-style)
- [Commit Convention](../contribute/commit-convention)