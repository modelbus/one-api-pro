---
title: 源码编译
description: 从源代码本地构建 Go 后端与前端主题。
category: install
order: 4
---

# 源码编译

> 适用：二次开发、贡献代码、需要在 PR 中调试。

## 前置条件

| 工具 | 版本 |
|---|---|
| Go | 1.25+ |
| Node.js | 22+ |
| pnpm | 9+ |
| Git | 任何近版本 |

## 后端

```bash
git clone https://github.com/modelbus/one-api-pro
cd one-api-pro
go build -o one-api-pro .
```

启动：

```bash
./one-api-pro
```

第一次启动会自动创建 SQLite 数据库并输出默认管理员账号。

## 前端

```bash
cd web/default-pro
pnpm install
pnpm build       # 产物会被 //go:embed 进二进制
```

如果不构建前端直接跑二进制，访问页面会看到「前端资源未嵌入」的提示。

## 调试循环

后端改完直接 `./one-api-pro` 重启。

前端改完：

```bash
cd web/default-pro
pnpm build
# 重新 build 后端二进制（因为 //go:embed）
cd ../..
go build -o one-api-pro .
```

或者更快的：前端 dev server 单独跑，配置代理到后端 3000 端口（适合 UI 频繁修改）。

## IDE / 编辑器

推荐 GoLand / VS Code + Go 扩展。代码风格见 [编码规范](../contribute/code-style)。

## 相关文档

- [开发环境搭建](../contribute/dev-setup)
- [编码规范](../contribute/code-style)
- [提交规范](../contribute/commit-convention)