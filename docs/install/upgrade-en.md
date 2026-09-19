---
title: Upgrade
description: How to upgrade safely to a new version.
category: install
order: 6
---

# Upgrade

> Standard pre-flight, mid-flight, and post-flight procedure.

## Before

1. **Back up the database** — see [Backup & Restore](./backup-restore).
2. **Read the [CHANGELOG](https://github.com/modelbus/one-api-pro/blob/main/CHANGELOG)** for the "Upgrade Notes" section of the target version.
3. **Check breaking changes** — incompatible fields, endpoints, config keys.

## Upgrade paths

### Docker / docker-compose

```bash
docker compose pull one-api-pro
docker compose up -d
```

Only the one-api-pro container restarts. Volumes (`./data`, `./mysql-data`) are preserved.

### Source build

```bash
git pull
cd web/default-pro && pnpm install && pnpm build && cd ../..
go build -o one-api-pro .
systemctl stop one-api-pro
systemctl start one-api-pro
```

### Pre-built binary

Replace the binary, restart the service. Data files are unchanged.

## During

`AutoMigrate` runs at startup — it adds tables / columns as needed. Watch the log:

```
[migrate] adding column foo.bar
[migrate] creating index idx_xxx
```

If it fails, look at the specific error (usually a column-type conflict).

## After

1. **Re-read [CHANGELOG](https://github.com/modelbus/one-api-pro/blob/main/CHANGELOG) "Upgrade Notes"**.
2. **Check Admin → Dashboard** — compare success rates before/after.
3. **Sample a few call logs** — confirm behavior is as expected.
4. **New features**: some versions require manual config in Admin → System Settings (the CHANGELOG will mention it).

## Rollback

If something goes badly wrong:

```bash
# Docker
docker compose down
mv ./mysql-data ./mysql-data.broken
cp -r ./backup-YYYYMMDD/mysql-data ./
docker compose up -d
```

If the schema didn't change between versions, you can simply downgrade to the previous binary.

## Related

- [Backup & Restore](./backup-restore)
- [CHANGELOG](../changelog/index)