---
title: Upgrade
description: "Upgrade procedure and data-migration notes."
category: install
order: 6
---

# Upgrade

> Upgrade procedure and data-migration notes.

## Before upgrading

1. Back up the database and config (see [Backup & Restore](/en/install/backup-restore)).
2. Skim the latest [CHANGELOG](https://github.com/modelbus/one-api-pro/blob/main/CHANGELOG) for breaking changes that need manual action.
3. Verify whether any new mandatory environment variables have been introduced (in particular namespace shifts between v0.0.x releases).

## Docker upgrade

The `/app/config` and `/app/data` volumes are reused, so the SQLite DB and `.env` are preserved.

```bash
# 1. Pull the new image
docker pull ghcr.io/modelbus/one-api-pro:latest

# 2. Stop the old container
docker stop one-api-pro

# 3. Start a new container with the same volumes & port
docker run -d \
  --name one-api-pro \
  --restart unless-stopped \
  -p 3000:3000 \
  -v /opt/one-api-pro/config:/app/config \
  -v /opt/one-api-pro/data:/app/data \
  -e TZ=Asia/Shanghai \
  ghcr.io/modelbus/one-api-pro:latest

# 4. Watch the logs
docker logs -f one-api-pro
```

`AutoMigrate` runs on first boot to add any missing tables; it never drops columns. To roll back:

```bash
docker stop one-api-pro && docker rm one-api-pro
docker run -d --name one-api-pro \
  -p 3000:3000 \
  -v /opt/one-api-pro/config:/app/config \
  -v /opt/one-api-pro/data:/app/data \
  ghcr.io/modelbus/one-api-pro:v0.0.20   # pin to a known-good tag
```

## Binary upgrade

```bash
# 1. Stop the old process
systemctl stop one-api-pro   # or kill / supervisorctl stop

# 2. Back up binary + config + data
cp one-api-pro one-api-pro.bak
cp -r config config.bak
cp -r data data.bak

# 3. Replace the binary
cp new-one-api-pro ./one-api-pro
chmod u+x one-api-pro

# 4. Start
systemctl start one-api-pro
```

`AutoMigrate` again handles table creation; the version is taken from the `VERSION` file next to the binary or in CWD (`--version`).

## docker-compose upgrade

```bash
# Pull the latest image
docker compose pull one-api-pro

# Roll over
docker compose up -d
```

Volumes declared in the top-level `volumes:` block are preserved.

## Database migration

- **Non-breaking**: new columns and tables are handled by `AutoMigrate`.
- **Breaking changes** are called out in CHANGELOG and shipped as one-shot scripts under `cmd/migrate_*`.
- **Cross-database switch** (SQLite → MySQL, MySQL → PostgreSQL) requires logical dumps via `mysqldump` / `pg_dump` / `sqlite3 .dump` plus any companion migration script.

## Cluster upgrade

To avoid dropping events while every node restarts at once, prefer **rolling upgrades**:

```text
1. Upgrade Node A (keep B / C running)
2. Wait until A is healthy and A↔B / A↔C heartbeats recover
3. Upgrade Node B
4. ...
```

In decentralized cluster mode, writes made while a node is offline are not auto-replayed. After it comes back, run `mysqldump` from a live node to re-seed (see [Multi-node Deployment](/en/decentralization/deployment)).

Next: [Backup & Restore](/en/install/backup-restore) · [Cluster Overview](/en/decentralization/overview).
