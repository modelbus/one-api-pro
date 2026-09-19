---
title: Log
description: Per-call record of every /v1/* request/response/billing — the first place to check when something goes wrong.
category: schema
order: 13
---

# Log

## What it is

`Log` records **every call to `/v1/*`**: caller, model, token consumption, channel hit, status code, error message. All traces are persisted for troubleshooting and auditing.

## Where to find it

- **Admin → Logs**: filter by time / user / model / status; paged
- **Personal Center → My Logs**: the current user's own logs (compact view)
- **Admin → Dashboard**: today's / this month's stats

## Operator-relevant fields

| Field | Meaning | Effect |
|---|---|---|
| Caller | User | "Who is calling". |
| Token | Token FK | "Which key is calling". |
| Channel | Channel FK | "Which upstream path was hit". |
| Model | e.g. `gpt-4o` | "Which model was used". |
| Input / output tokens | Number | Compute cost, locate overruns. |
| Status code | 200 / 401 / 429 / 5xx | First place to look on failure. |
| Error message | Text | Raw error from the upstream provider. |
| Latency | ms | Locate slow requests. |
| Quota charged | Number | Final deduction for this call. |

## Common troubleshooting paths

- **User reports "charged too much"**: check input/output tokens vs the upstream provider's own billing; ensure token counting matches.
- **User reports "always 429"**: look for "quota insufficient" in the error message; check `User.quota`.
- **User reports "model not working"**: see if `channel` is empty; check whether the model is selected in the channel.
- **Locating slow requests**: sort by latency descending; the slowest channel is the suspect.

## Cleanup / archival

Logs grow unbounded. Configure a retention window (90 days is typical in production). Admin → Logs → "Clear history" deletes logs older than the threshold.

## Related pages

- [Logs (admin)](/en/misc/logs)

## Related API

- `GET /api/log/` — list (admin)
- `GET /api/log/self` — current user's logs
- `GET /api/log/stat` — stats
- `DELETE /api/log/` — clear history