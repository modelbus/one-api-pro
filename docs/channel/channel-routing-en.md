---
title: Channel Routing
description: How One API Pro picks one channel out of many when a request comes in.
category: channel
order: 3
---

# Channel Routing

> When a user hits `POST /v1/chat/completions`, the system must pick one channel among the candidates.
> This page explains the rules, why some channels never get picked, and how to debug it.

## The 4-step selection

1. **Filter**: keep only channels that are enabled, advertise the requested model, and allow the user's group.
2. **Health check**: drop channels in cooldown, auto-disabled, or hitting concurrency / RPM limits.
3. **Score**: weighted score from `weight`, `priority`, and `is_fallback`.
4. **Decide**: pick the top scorer; tie → random.

## Filter rules

| Filter | Where set | If failed |
|---|---|---|
| Status = Enabled | Channel detail toggle | Channel is skipped. |
| Requested model in `models` | Channel detail | Channel never serves that model. |
| User's group in `group` | Channel detail | Channel hidden from this user. |
| Not auto-disabled / cooldown | Automatic | System skips it. |

## Weight / Priority / Fallback

These three together determine the score:

- **Weight**: bigger = more likely. Set cheap-and-fast high.
- **Priority**: bigger rank wins (10 beats 5).
- **is_fallback**: when on, this channel is only tried after all non-fallback channels fail.

Adjust weight / priority to nudge routing without touching anything else.

## Health checks (automatic)

- **Concurrency**: in-flight calls >= `max_concurrency` → reject new.
- **RPM**: last-60s requests >= `rpm` → reject.
- **Cooldown**: N consecutive failures → wait a few minutes.
- **Auto-disable**: persistent severe failures → channel is auto-disabled.

These need **no manual config**. If a channel keeps getting auto-disabled, check its [logs](/en/schema/log) `error` field.

## Debug: why wasn't my channel picked?

1. **Check filter**: run [Channel Test](./channel-test) with the same model + group.
2. **Check weights**: list candidates' weights — the highest should win.
3. **Check auto-state**: red flag on the channel detail (auto-disabled / cooldown / cap).
4. **Check call logs**: the `channel` field shows which one was actually picked.

## Advanced: model name mapping

Some proxies rename native models (e.g. `gpt-4o` → `gpt-4o-2024`). Fill `model_mapping` in JSON:

```json
{ "gpt-4o": "gpt-4o-2024" }
```

A caller requesting `gpt-4o` is forwarded as `gpt-4o-2024` to this channel.

## Related

- [Channel Overview](./overview)
- [Add a Channel](./add-channel)
- [Channel Test](./channel-test)