---
title: Docker 单实例部署
description: 用官方镜像跑一条 One API Pro 实例（最简最快）。
category: install
order: 2
---

# Docker 单实例部署

> 适用：试用、个人、小团队。
> 部署时长：5 分钟。

## 一条命令跑起来

```bash
docker run -d --name one-api-pro --restart always \
  -p 3000:3000 \
  -v $(pwd)/data:/app/data \
  -e TZ=Asia/Shanghai \
  ghcr.io/modelbus/one-api-pro:latest
```

启动后访问 `http://localhost:3000`。默认管理员账号：

- 用户名：`root`
- 密码：`123456`

**第一件事**：登录后立即改密码（个人中心）。

## 关键参数解释

| 项 | 含义 |
|---|---|
| `-d` | 后台运行 |
| `--restart always` | 进程崩溃 / Docker 重启时自动拉起 |
| `-p 3000:3000` | 把宿主机的 3000 端口映射到容器 3000 |
| `-v $(pwd)/data:/app/data` | 把数据库和上传文件持久化到宿主机 `./data` 目录 |
| `-e TZ=Asia/Shanghai` | 设置时区（影响日志时间戳、CRON 任务） |

## 后续要做什么

按顺序：

1. **改默认密码**：登录后 → 个人中心 → 修改密码
2. **添加渠道**：后台 → 渠道 → 新增，参考 [新增渠道](../channel/add-channel)
3. **配置模型定价**：后台 → 模型定价，参考 [模型定价](../pricing/model-price-management)
5. **（可选）启用反向代理**：直接暴露 3000 端口不安全，生产建议用 [反向代理](./reverse-proxy)
6. **（可选）配置备份**：参考 [备份与恢复](./backup-restore)

## 数据存在哪

容器内的 `/app/data` 目录：

- `one-api.db` — SQLite 主库
- `logs/` — 调用日志导出
- `uploads/` — 用户上传文件

通过 `-v` 挂载到宿主机，**重启 / 升级容器数据不丢**。

## 查看日志

```bash
docker logs -f one-api-pro
```

只显示最近 200 行加 `-n 200`。问题排查时常用 [容器内日志 + 后台日志查询] 配合定位。

## 升级版本

```bash
docker pull ghcr.io/modelbus/one-api-pro:latest
docker stop one-api-pro
docker rm one-api-pro
# 用同样参数重新 run
```

数据在挂载的 `./data` 里，不会丢。详细流程参考 [版本升级](./upgrade)。

## 常见问题

- **忘记密码**：用 `docker exec -it one-api-pro one-api-pro reset-password root 新密码`（v0.0.20+）
- **容器起不来**：看 `docker logs one-api-pro`，常见原因是 3000 端口被占用
- **想用 MySQL 而非 SQLite**：加 `-e SQL_DSN='user:pass@tcp(host:3306)/db'` 切换

## 相关文档

- [系统要求](./requirements)
- [docker-compose 部署](./docker-compose) — 适合生产
- [配置项与环境变量](./config)
- [备份与恢复](./backup-restore)