---
title: 备份与恢复
description: "数据库与上传资源的备份恢复方案。"
category: install
order: 7
---

# 备份与恢复

> 数据库与上传资源的备份恢复方案。
> Backup and recovery for the database and uploads.

数据主要落地点 / Where data lives:

| 类型 / Type | 路径 / Path | 备注 / Notes |
| --- | --- | --- |
| 主数据库 / Primary DB | `SQLITE_PATH`（默认 `one-api-pro.db`）或外部 MySQL/PostgreSQL | 包含全部业务表与系统设置 |
| 日志数据库 / Logs DB | `LOG_SQL_DSN`（若设置） | 高频写入，独立备份即可 |
| 上传 / 上传缓存 / Uploads | `/app/data` 下的运行时缓存 | 临时文件可重建 |
| 日志 / Logs | `--log-dir`（默认 `./logs`） | 用于排障与审计 |
| 配置 / Config | `/app/config/.env`、数据库 `options` 表 | 与备份联动 |

## SQLite 备份 / SQLite backup

SQLite 数据库在事务内一致时可直接拷贝；推荐使用 `sqlite3 .backup`（在线、热备）。

The SQLite DB is consistent inside a transaction and can be safely copied; prefer `sqlite3 .backup` (online / hot backup).

### 在线热备 / Online hot backup

```bash
# 在容器外（或主机）执行
sqlite3 /opt/one-api-pro/data/one-api-pro.db \
  ".timeout 5000" \
  ".backup '/opt/backup/one-api-pro-$(date +%Y%m%d-%H%M%S).db'"
```

### 简单复制（停机或低写入）/ Simple copy (downtime / low write)

```bash
systemctl stop one-api-pro
cp /opt/one-api-pro/data/one-api-pro.db \
   /opt/backup/one-api-pro-$(date +%Y%m%d).db
cp -r /opt/one-api-pro/data/logs \
      /opt/backup/logs-$(date +%Y%m%d)
systemctl start one-api-pro
```

> Docker 部署时拷贝容器内文件：先 `docker cp` 或通过 `docker exec` 调用 `sqlite3`。
> For Docker deployments, copy the file out with `docker cp` or invoke `sqlite3` via `docker exec`.

### 恢复 / Restore

```bash
systemctl stop one-api-pro
cp /opt/backup/one-api-pro-20260101.db \
   /opt/one-api-pro/data/one-api-pro.db
systemctl start one-api-pro
```

## MySQL 备份 / MySQL backup

```bash
# 全量逻辑备份（推荐）
mysqldump -h <host> -u root -p \
  --single-transaction --routines --triggers \
  --databases oneapi \
  > oneapi-$(date +%Y%m%d-%H%M%S).sql

# 日志库独立备份（可选）
mysqldump -h <host> -u root -p \
  --single-transaction \
  --databases oneapi_logs \
  > oneapi_logs-$(date +%Y%m%d-%H%M%S).sql

# 恢复
mysql -h <host> -u root -p oneapi < oneapi-20260101.sql
```

`--single-transaction` 借助 InnoDB 事务一致性，对运行中的库仍可得到一致快照。

`--single-transaction` relies on InnoDB consistency and gives a coherent snapshot even while the service is running.

## PostgreSQL 备份 / PostgreSQL backup

```bash
# 自定义格式，支持并行恢复与单表抽取
pg_dump -h <host> -U postgres -Fc oneapi \
  > oneapi-$(date +%Y%m%d-%H%M%S).dump

# 恢复
pg_restore -h <host> -U postgres -d oneapi --clean --if-exists \
  oneapi-20260101.dump
```

## 建议的备份节奏 / Suggested schedule

| 数据库 / DB | 频率 / Frequency | 保留 / Retention | 工具 / Tool |
| --- | --- | --- | --- |
| SQLite | 每 6 小时 | 14 天 / 14 days | `sqlite3 .backup` |
| MySQL `oneapi` | 每日全量 + 每 6 小时增量 binlog | 30 天全量 + 14 天 binlog | `mysqldump` + `binlog` |
| MySQL `oneapi_logs` | 每日全量 | 7 天 | `mysqldump` |
| PostgreSQL | 每日全量 | 30 天 | `pg_dump -Fc` |
| `/app/config/.env` | 每次修改后 | 不限 / indefinite | 版本控制 / VCS |

> ⚠️ **去中心化集群**：每个节点持有独立 MySQL，建议每个节点各自独立备份；如需跨节点统一恢复，可选主从复制或物理快照。
> ⚠️ **Decentralized clusters**: every node owns its own MySQL — back each node up independently; for cluster-wide recovery use master/slave replication or storage-level snapshots.

## 验证与演练 / Validation & drill

备份完成后建议周期性恢复演练：

Periodically rehearse restores:

1. 在隔离环境恢复到空库；
   Restore to an empty database in an isolated environment.
2. 启动一份 `one-api-pro --env` 指向该库，确认 `/api/status` 正常；
   Boot `one-api-pro --env` against the restored DB and confirm `/api/status` works.
3. 用 root 账号登录后台，校验关键数据（用户、渠道、套餐）。
   Log in as root and spot-check key data (users, channels, plans).

## 灾难恢复 / Disaster recovery

| 场景 / Scenario | 恢复步骤 / Steps |
| --- | --- |
| 容器 / 主机故障 | 重新部署容器 + 挂载同一 `/app/data` 与 `/app/config` 目录即可恢复 |
| 数据库误删表 | 从最近的 `mysqldump` / `sqlite3 .backup` 全量恢复 |
| `options` 系统设置被错误清空 | 使用 `options` 表的备份恢复；如无备份，可手动通过 `/api/option/` 接口重新设置 |
| 集群节点长期离线后数据漂移 | 从一个存活节点 `mysqldump` 同步数据后重启该节点 |

下一步 / Next: [反向代理](/zh/install/reverse-proxy) · [Cluster 概览](/zh/decentralization/overview)。
