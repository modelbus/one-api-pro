---
title: 备份与恢复
description: 数据库与上传资源的备份恢复方案。
category: install
order: 7
---

# 备份与恢复

> 每天自动备份 + 出事能快速恢复，是数据安全的两条腿。

## 数据存在哪里

| 类型 | 路径 | 重要性 |
|---|---|---|
| SQLite 主库 | `./data/one-api.db`（容器内 `/app/data/`） | **核心** |
| MySQL 数据卷 | `./mysql-data/` | **核心** |
| Redis 数据卷 | `./redis-data/` | 重要（节点间状态） |
| 用户上传 | `./data/uploads/` | 视业务而定 |
| 导出日志 | `./data/logs/` | 一般 |

## 推荐方案

### 1. 每天数据库全量备份（必须）

**SQLite**：

```bash
docker exec one-api-pro \
  sqlite3 /app/data/one-api.db ".backup '/app/data/backup-$(date +%Y%m%d).db'"

# 或停服后复制
docker stop one-api-pro
cp ./data/one-api.db ./backup/one-api-$(date +%Y%m%d).db
docker start one-api-pro
```

**MySQL**：

```bash
docker exec one-api-pro-mysql \
  mysqldump -uoneapi -poneapi-pass oneapi | gzip > ./backup/mysql-$(date +%Y%m%d).sql.gz
```

**自动**：写到 crontab 每天凌晨跑一次，保留 30 天。

### 2. 上传资源每周备份（视业务）

```bash
tar czf ./backup/uploads-$(date +%Y%m%d).tgz ./data/uploads/
```

### 3. 异地备份（生产必做）

把 `./backup/` 同步到对象存储：

- AWS S3 / 阿里 OSS / 腾讯 COS
- rsync 到另一台机器
- 7 天 → 30 天 → 1 年分层

工具：`restic` / `borgbackup` / 各云厂商 CLI。

## 恢复

```bash
# 1. 停服
docker compose down

# 2. 替换数据库
docker compose up -d mysql   # 先起 MySQL
gunzip < ./backup/mysql-YYYYMMDD.sql.gz | \
  docker exec -i one-api-pro-mysql mysql -uoneapi -poneapi-pass oneapi

# 3. 重启全部
docker compose up -d
```

## 验证

恢复后立刻做：

- [ ] 登录后台
- [ ] 后台 → 用户管理，确认用户数与备份前一致
- [ ] 后台 → 渠道管理，确认渠道配置正确
- [ ] 后台 → 套餐，确认套餐有效
- [ ] 发一条测试调用，确认日志记录成功

## 相关文档

- [Docker 单实例部署](./docker-deploy)
- [docker-compose 部署](./docker-compose)
- [系统要求](./requirements)