---
title: Docker Single-instance Deploy
description: Run one One API Pro instance from the official image — fastest path.
category: install
order: 2
---

# Docker Single-instance Deploy

> For trial / personal / small team. About 5 minutes.

## One-command start

```bash
docker run -d --name one-api-pro --restart always \
  -p 3000:3000 \
  -v $(pwd)/data:/app/data \
  -e TZ=Asia/Shanghai \
  ghcr.io/modelbus/one-api-pro:latest
```

Open `http://localhost:3000`. Default admin:

- Username: `root`
- Password: `123456`

**First thing**: change the password after login (Personal Center).

## Key flags explained

| Flag | Meaning |
|---|---|
| `-d` | Run in background |
| `--restart always` | Auto-restart on crash / Docker restart |
| `-p 3000:3000` | Map host port 3000 → container 3000 |
| `-v $(pwd)/data:/app/data` | Persist DB and uploads to host `./data` |
| `-e TZ=Asia/Shanghai` | Time zone (affects log timestamps, cron) |

## Next steps

In order:

1. **Change default password** — log in → Personal Center → change password.
2. **Add a channel** — Admin → Channels → Add, see [Add a Channel](../channel/add-channel).
3. **Configure model prices** — Admin → Model Prices, see [Model Price Management](../pricing/model-price-management).
4. **(Optional) Reverse proxy** — exposing port 3000 directly is unsafe for production; use [Reverse Proxy](./reverse-proxy).
5. **(Optional) Configure backups** — see [Backup & Restore](./backup-restore).

## Where data lives

Inside the container at `/app/data`:

- `one-api.db` — main SQLite DB
- `logs/` — exported call logs
- `uploads/` — uploaded files

Mounted to the host via `-v` — **restarts / image upgrades do NOT lose data**.

## Tail logs

```bash
docker logs -f one-api-pro
```

Add `-n 200` for the last 200 lines. Pair with Admin → Logs when troubleshooting.

## Upgrade

```bash
docker pull ghcr.io/modelbus/one-api-pro:latest
docker stop one-api-pro
docker rm one-api-pro
# re-run with the same flags
```

Data stays in the mounted `./data`. Full flow: [Upgrade](./upgrade).

## FAQ

- **Forgot password**: `docker exec -it one-api-pro one-api-pro reset-password root newpassword` (v0.0.20+)
- **Container won't start**: `docker logs one-api-pro` — usually port 3000 already in use.
- **Want MySQL instead of SQLite**: add `-e SQL_DSN='user:pass@tcp(host:3306)/db'`.

## Related

- [System Requirements](./requirements)
- [Docker Compose](./docker-compose) — for production
- [Configuration](./config)
- [Backup & Restore](./backup-restore)