---
title: 系统要求
description: 部署 One API Pro 需要什么样的硬件、操作系统、端口与数据库。
category: install
order: 1
---

# 系统要求

> 部署前要确认的环境清单。

## 硬件

| 规模 | 推荐 |
|---|---|
| 试用 / 个人 | 1 vCPU / 1 GB RAM / 5 GB 磁盘 |
| 小团队（< 100 用户） | 2 vCPU / 4 GB RAM / 20 GB SSD |
| 生产（数百到数千用户） | 4+ vCPU / 8+ GB RAM / 50 GB+ SSD（取决于日志保留时长） |

One API Pro 是单一 Go 二进制。CPU 主要吃在请求转发上；内存主要吃在日志缓存与并发请求。

## 操作系统

Linux（推荐 Ubuntu 22.04 / Debian 12 / CentOS Stream 9）、macOS 12+、Windows 10+。Docker 部署忽略此条。

## 端口

| 端口 | 用途 | 是否必须对外开放 |
|---|---|---|
| `3000` | HTTP 主端口（管理后台 + `/api/*` + `/v1/*`） | **是** |
| 集群节点间 | `CLUSTER_NODE_PORT`（默认 `3000`） | 仅多节点部署需要 |

## 数据库

两种选择，二选一：

### 内嵌 SQLite（默认）

零配置，开箱即用。适合单实例、< 1000 用户、QPS 较低（< 50）。

- 数据落盘在容器内 `/app/data/one-api.db`
- 备份只需复制这个文件

### 外部 MySQL（生产推荐）

适合多实例、高 QPS、大数据量场景。

- 最低 MySQL 5.7 / 8.0
- 需要提前创建空数据库 + 用户
- 连接串通过环境变量 `SQL_DSN` 配置

> 部署前决定好：上线后从 SQLite 迁到 MySQL 是可以的，但需要停机做数据迁移。

## Redis（可选但推荐）

集群多节点模式**必须** Redis：

- 最低 Redis 6.0
- 用于节点间状态共享与限流
- 单实例部署可不用 Redis

## 反向代理（生产推荐）

直接暴露容器端口到公网不安全。建议在前面套一层：

- Nginx / Caddy / Cloudflare
- 用于 TLS 终止 + HSTS + 真实 IP 透传
- 见 [反向代理](./reverse-proxy)

## 其他

- **域名**：建议准备一个正式域名 + TLS 证书
- **SMTP**：可选，用于发送密码重置 / 兑换码 / 订单通知邮件
- **支付通道**：若要启用充值 / 套餐支付，需先在 [支付配置](../pricing/payment-settings) 里配好微信 / 支付宝 / 银行

## 准备清单（部署前确认）

- [ ] 服务器 ≥ 2 vCPU / 4 GB RAM
- [ ] 操作系统 Linux / macOS / Windows
- [ ] 80 / 443（反向代理）或 3000（直连）可达
- [ ] 数据库就绪：默认 SQLite 即可，或 MySQL ≥ 5.7
- [ ] Redis ≥ 6.0（仅多节点）
- [ ] 反向代理（Nginx / Caddy / Cloudflare 之一）
- [ ] 域名 + TLS 证书
- [ ] （可选）SMTP 邮件账号
- [ ] （可选）支付通道凭证

## 下一步

- 最快的部署 → [Docker 单实例部署](./docker-deploy)
- 生产环境 → [docker-compose 部署](./docker-compose)
- 想自己编译 → [源码编译](./source-build)