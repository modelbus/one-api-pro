---
title: Backup & Restore
description: Database and uploads backup strategy.
category: install
order: 7
---

# Backup & Restore

> Daily automatic backup + quick recovery is the two legs of data safety.

## Where data lives

| Type | Path | Importance |
|---|---|---|
| SQLite main DB | `./data/one-api.db` (container `/app/data/`) | **Critical** |
| MySQL data volume | `./mysql-data/` | **Critical** |
| Redis data volume | `./redis-data/` | Important (inter-node state) |
| User uploads | `./data/uploads/` | Depends on usage |
| Exported logs | `./data/logs/` | Low |

## Recommended strategy

### 1. Daily database full backup (required)

**SQLite**:

```bash
docker exec one-api-pro \
  sqlite3 /app/data/one-api.db ".backup '/app/data/backup-$(date +%Y%m%d).db'"

# or stop-then-copy:
docker stop one-api-pro
cp ./data/one-api.db ./backup/one-api-$(date +%Y%m%d).db
docker start one-api-pro
```

**MySQL**:

```bash
docker exec one-api-pro-mysql \
  mysqldump -uoneapi -poneapi-pass oneapi | gzip > ./backup/mysql-$(date +%Y%m%d).sql.gz
```

**Automate**: add to crontab, run at midnight, retain 30 days.

### 2. Weekly uploads backup (if used)

```bash
tar czf ./backup/uploads-$(date +%Y%m%d).tgz ./data/uploads/
```

### 3. Off-site backup (required in production)

Sync `./backup/` to object storage:

- AWS S3 / Aliyun OSS / Tencent COS
- rsync to another machine
- 7 days → 30 days → 1 year tiered

Tools: `restic` / `borgbackup` / cloud-provider CLIs.

## Restore

```bash
# 1. Stop
docker compose down

# 2. Replace DB
docker compose up -d mysql
gunzip < ./backup/mysql-YYYYMMDD.sql.gz | \
  docker exec -i one-api-pro-mysql mysql -uoneapi -poneapi-pass oneapi

# 3. Restart all
docker compose up -d
```

## Verify after restore

- [ ] Log in to the admin console
- [ ] Admin → Users — count matches
- [ ] Admin → Channels — configurations intact
- [ ] Admin → Plans — plans still active
- [ ] Send a test call — log records it

## Related

- [Docker Single-instance](./docker-deploy)
- [Docker Compose](./docker-compose)
- [System Requirements](./requirements)