---
title: Dev Setup
description: "Local setup for Go 1.25, Node 22 and pnpm 9."
category: contribute
order: 1
---

# Dev Setup

> Local setup for Go 1.25, Node 22 and pnpm 9.
> Go 1.25、Node 22、pnpm 9 的本地准备。

## 1. Toolchain / 必备工具

| Item / 项目 | Version / 版本 | Notes |
| --- | --- | --- |
| Go | 1.25.0 (pinned in `go.mod`) | <https://go.dev/dl/> |
| Node.js | 22+ | <https://nodejs.org/> |
| pnpm | 9 | `npm i -g pnpm`; locked by `pnpm-lock.yaml` |
| Git | 2.30+ | Commit messages follow Conventional Commits |
| SQLite (optional) | 3 | Only when you don't run MySQL/PG locally |

> Windows: prefer WSL2; macOS: `brew install go node`; Linux: package manager or `nvm`.

## 2. Clone & Init / 克隆与初始化

```bash
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro

# Backend
go mod download

# Frontend
cd web/default-pro
pnpm install
cd ../..
```

## 3. Run the Backend / 后端运行

```bash
go run .
```

Default: listens on `:3000`, loads SQLite `one-api-pro.db`, runs `gorm.AutoMigrate` on first start.

默认监听 `:3000`，加载 SQLite `one-api-pro.db`，首次启动自动 `gorm.AutoMigrate`。

Useful env vars / 常用环境变量:

```bash
SQL_DSN=                              # empty → SQLite; non-empty for MySQL/PG DSN
REDIS_CONN_STRING=redis://127.0.0.1:6379
DEBUG=true                            # Gin debug logs
DEBUG_SQL=true                         # GORM SQL logs
SESSION_SECRET=$(uuidgen)             # shared by multi-instance
```

## 4. Run the Frontend / 前端运行

```bash
cd web/default-pro
pnpm dev    # http://localhost:3001
```

`vite.config.js` proxies `/api/*` and `/v1/*` to `:3000` automatically; no CORS setup needed.

`vite.config.js` 已配置自动把 `/api/*` 与 `/v1/*` 代理到 `:3000`，不用手动配 CORS。

Theme build (output is `//go:embed`-ed into the binary):

主题构建（产物被 `//go:embed` 进二进制）：

```bash
cd web/default-pro
pnpm build    # output → web/build/default-pro/
```

## 5. Database Initialization / 数据库初始化

`model.InitDB()` creates tables and the root account (`root` / `123456`) on first start — **change the password immediately**.

启动时 `model.InitDB()` 会自动建表 + 创建 root 账号，**首次启动后立即改密**。

## 6. Tests / 测试

```bash
# Backend: glebarez/sqlite in-memory + AutoMigrate
go test ./...

# Frontend: Node's built-in test runner, no framework dependency
node web/default-pro/src/utils/*.test.mjs
```

Common test locations / 常用单测位置:

| Path / 路径 | Coverage / 覆盖 |
| --- | --- |
| `model/cache_test.go` | Redis cache + low-water fallback |
| `model/order_payment_test.go` | Order-number generation, activation idempotency |
| `model/token_test.go` | "Never expire" sentinel |
| `controller/topup_test.go` | Top-up amount math + exchange rate |
| `relay/billing/quota_preconsume_test.go` | Pre-consume formula + Plan short-circuit |
| `web/default-pro/src/utils/*.test.mjs` | Frontend amount / token / top-up utils |

## 7. Project Layout / 项目结构

```text
one-api-pro/
├── cmd/                  one-shot CLI tools
├── common/               config / helper / payment / render / i18n
├── controller/           Gin handlers
├── middleware/           auth / rate-limit / turnstile / language
├── model/                GORM models + business functions (no HTTP dependency)
├── relay/                LLM relay layer
│   ├── adaptor/          provider self-registration (init → registry.Register)
│   ├── billing/          pre-consume / refund quota
│   ├── registry/         central index (by ID and LegacyType)
│   └── handler/          relay entry (Chat / Embeddings / Image)
├── router/               route registration
├── monitor/              monitoring / patrol
├── web/
│   ├── default-pro/      ★ the only maintained frontend project
│   ├── air/ default/ berry/   historical themes (do not edit)
│   ├── THEMES            theme-name list (only default-pro)
│   └── build/            output of `pnpm build` (do not edit)
├── docs/                 user docs
├── main.go                entrypoint
├── VERSION               current version
├── CHANGELOG/            per-version .md
├── go.mod                Go module (1.25.0)
└── .github/workflows/    release / release-docker
```

## 8. Recommended Editor / 推荐编辑器

- **VSCode**: extensions Go (`golang.go`), Vue (`Vue.volar`), ESLint, Prettier.
  VSCode：插件 Go + Vue (Volar) + ESLint + Prettier。
- **GoLand**: built-in `go fmt` / `goimports` + Go toolchain.
  GoLand：内置 `go fmt` / `goimports` + 内置 Go 工具链。

`.editorconfig` and `.prettierrc` are committed; formatting runs automatically before commit.

`.editorconfig` 与 `.prettierrc` 已包含，提交前会自动按风格格式化。

## 9. Pre-commit Checklist / 提交前自检

Per [Commit Convention §3.5](/en/contribute/commit-convention#pre-commit-checklist), always run the trio:

按 [提交规范 §3.5](/zh/contribute/commit-convention#提交前自检) 必跑三件套：

```bash
go build ./...
go test ./model/ ./controller/ ./middleware/
cd web/default-pro && pnpm build
```

## 10. Getting Help / 遇到问题

- Docs first: [Architecture](/en/start/architecture) and [Troubleshooting](/en/faq/troubleshooting).
- Issues: <https://github.com/modelbus/one-api-pro/issues>; attach diagnostics (see [Troubleshooting §6](/en/faq/troubleshooting#collecting-diagnostics)).
- Discussions: open-ended RFC / Q&A.

Next: [Code Style](/en/contribute/code-style) · [Commit Convention](/en/contribute/commit-convention) · [Add a Provider](/en/contribute/add-provider).