---
title: Backup & Restore
description: "Backup and recovery for the database and uploads."
category: install
order: 7
---

# Backup & Restore

> Backup and recovery for the database and uploads.

Where data lives:

| Type | Path | Notes |
| --- | --- | --- |
| Primary DB | `SQLITE_PATH` (default `one-api-pro.db`) or external MySQL/PostgreSQL | All business tables + system settings |
| Logs DB | `LOG_SQL_DSN` (if set) | High-frequency writes; back up independently |
| Uploads / cache | under `/app/data` | Temporary; rebuildable |
| Logs | `--log-dir` (default `./logs`) | Used for troubleshooting & audit |
| Config | `/app/config/.env`, `options` table | Back up alongside the DB |

## SQLite backup

The SQLite DB is consistent inside a transaction and can be safely copied; prefer `sqlite3 .backup` (online / hot backup).

### Online hot backup

```bash
# Run outside (or inside) the container
sqlite3 /opt/one-api-pro/data/one-api-pro.db \
  ".timeout 5000" \
  ".backup '/opt/backup/one-api-pro-$(date +%Y%m%d-%H%M%S).db'"
```

### Simple copy (downtime / low write)

```bash
systemctl stop one-api-pro
cp /opt/one-api-pro/data/one-api-pro.db \
   /opt/backup/one-api-pro-$(date +%Y%m%d).db
cp -r /opt/one-api-pro/data/logs \
      /opt/backup/logs-$(date +%Y%m%d)
systemctl start one-api-pro
```

> For Docker deployments, copy the file out with `docker cp` or invoke `sqlite3` via `docker exec`.

### Restore

```bash
systemctl stop one-api-pro
cp /opt/backup/one-api-pro-20260101.db \
   /opt/one-api-pro/data/one-api-pro.db
systemctl start one-api-pro
```

## MySQL backup

```bash
# Full logical backup (recommended)
mysqldump -h <host> -u root -p \
  --single-transaction --routines --triggers \
  --databases oneapi \
  > oneapi-$(date +%Y%m%d-%H%M%S).sql

# Logs DB independent backup (optional)
mysqldump -h <host> -u root -p \
  --single-transaction \
  --databases oneapi_logs \
  > oneapi_logs-$(date +%Y%m%d-%H%M%S).sql

# Restore
mysql -h <host> -u root -p oneapi < oneapi-20260101.sql
```

`--single-transaction` relies on InnoDB consistency and gives a coherent snapshot even while the service is running.

## PostgreSQL backup

```bash
# Custom format — supports parallel restore and selective table restore
pg_dump -h <host> -U postgres -Fc oneapi \
  > oneapi-$(date +%Y%m%d-%H%M%S).dump

# Restore
pg_restore -h <host> -U postgres -d oneapi --clean --if-exists \
  oneapi-20260101.dump
```

## Suggested schedule

| DB | Frequency | Retention | Tool |
| --- | --- | --- | --- |
| SQLite | every 6 hours | 14 days | `sqlite3 .backup` |
| MySQL `oneapi` | daily full + 6-hour binlog increment | 30-day full + 14-day binlog | `mysqldump` + `binlog` |
| MySQL `oneapi_logs` | daily full | 7 days | `mysqldump` |
| PostgreSQL | daily full | 30 days | `pg_dump -Fc` |
| `/app/config/.env` | after each change | indefinite | VCS |

> ⚠️ **Decentralized clusters**: every node owns its own MySQL — back each node up independently. For cluster-wide recovery, use master/slave replication or storage-level snapshots.

## Validation & drill

Periodically rehearse restores:

1. Restore to an empty database in an isolated environment.
2. Boot `one-api-pro --env` against the restored DB and confirm `/api/status` works.
3. Log in as root and spot-check key data (users, channels, plans).

## Disaster recovery

| Scenario | Steps |
| --- | --- |
| Container / host failure | Redeploy the container with the same `/app/data` and `/app/config` mounts |
| Accidental DROP TABLE | Restore from the latest `mysqldump` / `sqlite3 .backup` |
| `options` table wiped | Restore from a backup of `options`; otherwise re-set via `/api/option/` |
| Cluster node offline too long → drift | `mysqldump` from a live node, re-import, restart |

Next: [Reverse Proxy](/en/install/reverse-proxy) · [Cluster Overview](/en/decentralization/overview).
