---
title: 版本升级
description: "版本升级流程与数据迁移注意事项。"
category: install
order: 6
---

# 版本升级

> 版本升级流程与数据迁移注意事项。

## 升级前

1. 备份数据库与配置（详见 [备份与恢复](/zh/install/backup-restore)）；
2. 查阅最新 [CHANGELOG](https://github.com/modelbus/one-api-pro/blob/main/CHANGELOG)，确认是否有需要手工介入的破坏性变更；
3. 检查是否有新增的必填环境变量（特别是 v0.0.x 之间的命名空间调整）。

## Docker 升级

数据卷 `/app/config` 与 `/app/data` 在升级中保持原样挂载，数据库与 `.env` 都不会被覆盖。

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

```bash
docker stop one-api-pro && docker rm one-api-pro
docker run -d --name one-api-pro \
  -p 3000:3000 \
  -v /opt/one-api-pro/config:/app/config \
  -v /opt/one-api-pro/data:/app/data \
  ghcr.io/modelbus/one-api-pro:v0.0.20   # 显式锁回上一个 tag
```

## 二进制升级

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

## docker-compose 升级

```bash
# 拉取最新镜像
docker compose pull one-api-pro

# 滚动重启
docker compose up -d
```

数据卷声明在 `volumes:` 顶层，不会被替换。

## 数据库迁移

- **不破坏**：常规新增字段、新增表由 `AutoMigrate` 自动处理，无需手工介入；
- **破坏性变更**：CHANGELOG 会明确标注；通常以 `cmd/migrate_*` 一次性脚本形式发布；
- **跨数据库切换**（SQLite → MySQL、MySQL → PostgreSQL）：需自行使用 `mysqldump` / `pg_dump` / `sqlite3 .dump` 做逻辑导出与导入，必要时运行配套迁移脚本。

## 集群升级

为避免多节点同时重启造成事件丢失，建议 **逐节点滚动升级**（Rolling Upgrade）。

```text
1. 升级节点 A（保持 B / C 运行）
2. 等待 A 健康、A↔B、A↔C 心跳恢复
3. 升级节点 B
4. ...
```

去中心化集群模式下，单节点短暂离线期间产生的写入不会自动回填；该节点恢复后建议执行 `mysqldump` 同步一次（详见 [多节点部署](/zh/decentralization/deployment)）。

下一步 / Next: [备份与恢复](/zh/install/backup-restore) · [Cluster 概览](/zh/decentralization/overview)。

