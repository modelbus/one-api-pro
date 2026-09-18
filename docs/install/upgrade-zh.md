---
title: 版本升级
description: "版本升级流程与数据迁移注意事项。"
category: install
order: 6
---

# 版本升级

> 版本升级流程与数据迁移注意事项。
> Upgrade procedure and data-migration notes.

## 升级前 / Before upgrading

1. 备份数据库与配置（详见 [备份与恢复](/zh/install/backup-restore)）；
   Back up the database and config (see [Backup & Restore](/en/install/backup-restore)).
2. 查阅最新 [CHANGELOG](https://github.com/modelbus/one-api-pro/blob/main/CHANGELOG)，确认是否有需要手工介入的破坏性变更；
   Check the latest [CHANGELOG](https://github.com/modelbus/one-api-pro/blob/main/CHANGELOG) for any breaking changes that require manual action.
3. 检查是否有新增的必填环境变量（特别是 v0.0.x 之间的命名空间调整）。
   Verify whether any new mandatory environment variables have been introduced.

## Docker 升级 / Docker upgrade

数据卷 `/app/config` 与 `/app/data` 在升级中保持原样挂载，数据库与 `.env` 都不会被覆盖。

The `/app/config` and `/app/data` volumes are reused, so the SQLite DB and `.env` are preserved.

```bash
# 1. 拉取新版镜像
docker pull ghcr.io/modelbus/one-api-pro:latest

# 2. 停掉旧容器
docker stop one-api-pro

# 3. 用同样的卷与端口启动新容器
docker run -d \
  --name one-api-pro \
  --restart unless-stopped \
  -p 3000:3000 \
  -v /opt/one-api-pro/config:/app/config \
  -v /opt/one-api-pro/data:/app/data \
  -e TZ=Asia/Shanghai \
  ghcr.io/modelbus/one-api-pro:latest

# 4. 观察日志
docker logs -f one-api-pro
```

首次启动时 `AutoMigrate` 会自动补齐缺失表；不会主动删除列。失败回滚：

`AutoMigrate` runs on first boot to add any missing tables; it never drops columns. To roll back:

```bash
docker stop one-api-pro && docker rm one-api-pro
docker run -d --name one-api-pro \
  -p 3000:3000 \
  -v /opt/one-api-pro/config:/app/config \
  -v /opt/one-api-pro/data:/app/data \
  ghcr.io/modelbus/one-api-pro:v0.0.20   # 显式锁回上一个 tag
```

## 二进制升级 / Binary upgrade

```bash
# 1. 停旧进程
systemctl stop one-api-pro   # 或 kill / supervisorctl stop

# 2. 备份可执行文件 + 配置 + 数据
cp one-api-pro one-api-pro.bak
cp -r config config.bak
cp -r data data.bak

# 3. 替换二进制
cp new-one-api-pro ./one-api-pro
chmod u+x one-api-pro

# 4. 启动
systemctl start one-api-pro
```

`AutoMigrate` 同样会自动建表；版本号取自当前目录或二进制同目录的 `VERSION` 文件（`--version` 输出）。

`AutoMigrate` again handles table creation; the version is taken from the `VERSION` file next to the binary or in CWD.

## docker-compose 升级 / docker-compose upgrade

```bash
# 拉取最新镜像
docker compose pull one-api-pro

# 滚动重启
docker compose up -d
```

数据卷声明在 `volumes:` 顶层，不会被替换。

Volumes declared at the top-level `volumes:` block are preserved.

## 数据库迁移 / Database migration

- **不破坏**：常规新增字段、新增表由 `AutoMigrate` 自动处理，无需手工介入；
  **Non-breaking**: new columns/tables are handled by `AutoMigrate`.
- **破坏性变更**：CHANGELOG 会明确标注；通常以 `cmd/migrate_*` 一次性脚本形式发布；
  **Breaking changes** are called out in CHANGELOG and shipped as one-shot scripts under `cmd/migrate_*`.
- **跨数据库切换**（SQLite → MySQL、MySQL → PostgreSQL）：需自行使用 `mysqldump` / `pg_dump` / `sqlite3 .dump` 做逻辑导出与导入，必要时运行配套迁移脚本。
  **Cross-DB switch** (SQLite → MySQL, MySQL → PostgreSQL) requires logical dumps via `mysqldump` / `pg_dump` / `sqlite3 .dump` plus any companion migration script.

## 集群升级 / Cluster upgrade

为避免多节点同时重启造成事件丢失，建议 **逐节点滚动升级**（Rolling Upgrade）。

To avoid dropping events while every node restarts at once, prefer **rolling upgrades**.

```text
1. 升级节点 A（保持 B / C 运行）
2. 等待 A 健康、A↔B、A↔C 心跳恢复
3. 升级节点 B
4. ...
```

去中心化集群模式下，单节点短暂离线期间产生的写入不会自动回填；该节点恢复后建议执行 `mysqldump` 同步一次（详见 [多节点部署](/zh/decentralization/deployment)）。

In decentralized cluster mode, writes made while a node is offline are not auto-replayed; after it comes back, run `mysqldump` from a live node to re-seed (see [Multi-node Deployment](/en/decentralization/deployment)).

下一步 / Next: [备份与恢复](/zh/install/backup-restore) · [Cluster 概览](/zh/decentralization/overview)。
