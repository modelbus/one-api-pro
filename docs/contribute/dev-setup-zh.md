---
title: 开发环境搭建
description: "Go 1.25、Node 22、pnpm 9 的本地准备。"
category: contribute
order: 1
---

# 开发环境搭建

> Go 1.25、Node 22、pnpm 9 的本地准备。
> Local setup for Go 1.25, Node 22 and pnpm 9.

## 1. 必备工具 / Toolchain

| 项目 / Item | 版本 / Version | 备注 / Notes |
| --- | --- | --- |
| Go | 1.25.0（`go.mod` 锁定） | <https://go.dev/dl/> |
| Node.js | 22+ | <https://nodejs.org/> |
| pnpm | 9 | `npm i -g pnpm`；`pnpm-lock.yaml` 锁定 |
| Git | 2.30+ | commit message 走 Conventional Commits |
| SQLite（可选） | 3 | 仅本机无 MySQL/PG 时用 |

> Windows 推荐用 WSL2；macOS 用 `brew install go node`；Linux 用包管理器或 `nvm`。

## 2. 克隆与初始化 / Clone & Init

```bash
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro

# 后端
go mod download

# 前端
cd web/default-pro
pnpm install
cd ../..
```

## 3. 后端运行 / Run the Backend

```bash
go run .
```

默认监听 `:3000`，加载 SQLite `one-api-pro.db`，首次启动自动 `gorm.AutoMigrate`。

By default it listens on `:3000`, loads SQLite `one-api-pro.db`, and runs `gorm.AutoMigrate` on first start.

常用环境变量 / Useful env vars:

```bash
SQL_DSN=                              # 空 → SQLite；非空时填 MySQL/PG DSN
REDIS_CONN_STRING=redis://127.0.0.1:6379
DEBUG=true                            # Gin debug 日志
DEBUG_SQL=true                         # GORM SQL 日志
SESSION_SECRET=$(uuidgen)             # 多实例共享 Session
```

## 4. 前端运行 / Run the Frontend

```bash
cd web/default-pro
pnpm dev    # http://localhost:3001
```

`vite.config.js` 已配置自动把 `/api/*` 与 `/v1/*` 代理到 `:3000`，不用手动配 CORS。

`vite.config.js` proxies `/api/*` and `/v1/*` to `:3000` automatically; no CORS setup needed.

主题构建（产物被 `//go:embed` 进二进制）：

Theme build (output is `//go:embed`-ed into the binary):

```bash
cd web/default-pro
pnpm build    # 产物 → web/build/default-pro/
```

## 5. 数据库初始化 / Database Initialization

启动时 `model.InitDB()` 会自动建表 + 创建 root 账号（`root` / `123456`），**首次启动后立即改密**。

`model.InitDB()` creates tables and the root account (`root` / `123456`) on first start — **change the password immediately**.

## 6. 测试 / Tests

```bash
# 后端：glebarez/sqlite 内存 + AutoMigrate
go test ./...

# 前端：Node 自带 test runner，无框架依赖
node web/default-pro/src/utils/*.test.mjs
```

常用单测位置 / Common test locations:

| 路径 / Path | 覆盖 / Coverage |
| --- | --- |
| `model/cache_test.go` | Redis 缓存与低水位回源 |
| `model/order_payment_test.go` | 订单号生成、激活幂等 |
| `model/token_test.go` | 「永不过期」sentinel |
| `controller/topup_test.go` | 充值金额计算与汇率 |
| `relay/billing/quota_preconsume_test.go` | 预扣公式与 Plan 短路 |
| `web/default-pro/src/utils/*.test.mjs` | 前端金额 / token / 充值工具 |

## 7. 项目结构 / Project Layout

```text
one-api-pro/
├── cmd/                  一次性 CLI 工具
├── common/               config / helper / payment / render / i18n
├── controller/           Gin handler
├── middleware/           auth / rate-limit / turnstile / language
├── model/                GORM 模型 + 业务函数（不依赖 HTTP）
├── relay/                LLM 渠道转发
│   ├── adaptor/          Provider 自注册（init → registry.Register）
│   ├── billing/          预扣 / 退还配额
│   ├── registry/         中央索引（按 ID 与 LegacyType）
│   └── handler/          转发入口（Chat / Embeddings / Image）
├── router/               路由注册
├── monitor/              监控 / 巡检
├── web/
│   ├── default-pro/      ★ 唯一维护的前端项目
│   ├── air/ default/ berry/   历史主题（禁改）
│   ├── THEMES            主题名清单（仅含 default-pro）
│   └── build/            npm run build 产物（禁手改）
├── docs/                 用户文档
├── main.go               入口
├── VERSION               当前版本号
├── CHANGELOG/            变更日志（每个版本独立 .md）
├── go.mod                Go module（1.25.0）
└── .github/workflows/    release / release-docker
```

## 8. 推荐编辑器 / Editor

- **VSCode**：插件 Go (`golang.go`) + Vue (`Vue.volar`) + ESLint + Prettier。
  VSCode: Go, Vue (Volar), ESLint, Prettier.
- **GoLand**：内置 `go fmt` / `goimports` + 内置 Go 工具链。
  GoLand: built-in Go toolchain.

`.editorconfig` 与 `.prettierrc` 已包含，提交前会自动按风格格式化。

`.editorconfig` and `.prettierrc` are committed; formatting runs automatically before commit.

## 9. 提交前自检 / Pre-commit Checklist

按 [提交规范 §3.5](/zh/contribute/commit-convention#提交前自检) 必跑三件套：

Per [Commit Convention §3.5](/en/contribute/commit-convention#pre-commit-checklist), always run the trio:

```bash
go build ./...
go test ./model/ ./controller/ ./middleware/
cd web/default-pro && pnpm build
```

## 10. 遇到问题 / Getting Help

- **文档**：先看 [架构总览](/zh/start/architecture) 与 [故障排查](/zh/faq/troubleshooting)。
- **Issue**：<https://github.com/modelbus/one-api-pro/issues>，按版本附诊断信息（见 [故障排查 §6](/zh/faq/troubleshooting#收集诊断信息)）。
- **Discussions**：开放性讨论 / RFC。

下一步 / Next: [编码规范](/zh/contribute/code-style) · [提交规范](/zh/contribute/commit-convention) · [新增 Provider](/zh/contribute/add-provider)。