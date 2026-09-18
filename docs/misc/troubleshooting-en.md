---
title: Troubleshooting
description: "Diagnose channel 401s, billing anomalies, stuck orders."
category: misc
order: 2
---

# Troubleshooting

> Diagnose channel 401s, billing anomalies, stuck orders.
> 渠道 401、计费异常、订单不激活等问题的定位。

Categorize by symptom, gather diagnostics first, then trace the cause, then apply a fix.

按现象归类，先看「诊断信息」，再走「排查路径」，最后给「修复方案」。

## 1. Channel Issues / 渠道相关

### 1.1 Channel test passes (200) but `/v1/*` returns 401

**Diagnostics / 诊断**:

```bash
# Trigger test
curl -X POST http://localhost:3000/api/channel/test/<id> -b cookies.txt
# Inspect x-request-id + upstream response in the body
```

**Trace / 排查**:

| Symptom | Cause | Fix |
| --- | --- | --- |
| Test OK, call 401 | Wrong `base_url` (e.g. duplicated `/v1/`) | Drop the trailing `/v1`; OpenAI base URL is usually `https://api.openai.com` |
| Sporadic 401 then OK | Upstream rate-limit / risk control | Rotate the key, adjust rate-limit, increase `cooldown` |
| Shared key across providers | `abilities.group` mismatch | Whitelist in the user-group table |

### 1.2 Channel stuck in cooldown forever

- **Diagnostic**: `channels.cooldown_until` keeps being > `now`.
- **Cause**: Upstream keeps erroring; `CHANNEL_DEFAULT_COOLDOWN_SECONDS=60` doubles on each failure up to `CHANNEL_MAX_COOLDOWN_SECONDS=600`.
- **Fix**: Manually reset `cooldown_until=0`; trace the upstream root cause.

### 1.3 Channel `weight` change does not take effect

With `BATCH_UPDATE_ENABLED=true` weight writes are batched (default 5s). Either disable batching, or restart the relevant cache key.

`BATCH_UPDATE_ENABLED=true` 时权重聚合写入有窗口（默认 5s）。要么关闭批处理，要么重启相关 cache key。

## 2. Billing & Cache / 计费与缓存

### 2.1 Repeated 403 `insufficient_user_quota`

**Common pre-v0.0.21 cause**: Redis `user_quota:<id>` cache drifts from DB. Fixed by:

**v0.0.21 之前的常见根因**：Redis `user_quota:<id>` 缓存与 DB 不一致。已修复：

- `IncreaseUserQuota` writes DB then refreshes Redis synchronously.
- The "trusted, skip pre-consume" guard moved before the Redis decrement.
- `userQuotaLowWaterMark = 50_000` (was 500), removing the "cache drift dead zone" for high-cost models.

**Check your version**:

```bash
./one-api-pro --version
```

Upgrade to ≥ v0.0.21. If the symptom persists, clear the `user_quota:<id>` cache and restart:

升级到 ≥ v0.0.21；若仍出现请清掉 `user_quota:<id>` 缓存重启：

```bash
redis-cli DEL user_quota:1 user_quota:2 ...
```

### 2.2 "Never expire" token flagged as expired

**Historical fix**: since v0.0.21 the sentinel is "`<=0` means never expire" (`model/token.go` / `controller/token.go`); the frontend `buildTokenExpiredTime` always submits `-1`.

**修复历史**：v0.0.21 起 sentinel 从「仅 `-1` 表示永不过期」统一为「`<=0` 均视为永不过期」。

After upgrade, dirty data `expired_time=0` is no longer flagged as expired.

### 2.3 `tiktoken` counting panic

**Pre-v0.0.21**: under a firewall the HTTP BPE download times out → incomplete fallback → `Encode()` is called on a `nil` encoder → panic.

**v0.0.21 fix**: 4 BPE files (`cl100k_base` / `o200k_base` / `p50k_base` / `r50k_base`) are `//go:embed`-ed into the binary; startup has no network dependency.

```text
httpBpeLoader three-tier lookup:
  1. embedded BPE (default, works out of the box)
  2. TIKTOKEN_CACHE_DIR (advanced override; skips embedded + HTTP)
  3. HTTP download (last-resort fallback for future URL changes)
```

Upgrade to v0.0.21+ to resolve; advanced users can set `TIKTOKEN_CACHE_DIR` to skip the embedded copy.

## 3. Orders & Async Notify / 订单与异步通知

### 3.1 Async notify verification failed

| Channel | Verification | Field source |
| --- | --- | --- |
| WeChat | MD5 signature + `return_code=SUCCESS` | `payment.wechat.config.api_key` |
| Alipay | RSA2 sign verification | `payment.alipay.config.private_key/public_key` |
| Bank | Manual reconciliation (admin "Mark Paid") | via admin UI |

**Diagnostics / 诊断**:

```bash
# 1. Tail notify logs
tail -f logs/$(date +%Y-%m-%d).log | grep -i notify

# 2. Query order status
curl http://localhost:3000/api/order/<id> -b cookies.txt
```

**Common pitfalls / 典型问题**:

- WeChat API Key rotated but the "Payment Settings" page not saved → callback signature fails.
- Alipay `notify_url` reachable but the domain isn't ICP-filed → callback rejected.
- "Bank Transfer" has no callback; admin must click "Mark Paid" in the order detail.

### 3.2 Order "Paid" but plan not activated

- Check `orders.status == 1`; if so but `user_plans` not increased:
  检查 `orders.status` 是否 = 1；若是但 `user_plans` 未增加：
- Pre-v0.0.21 stale data may come from cache races; run:
  v0.0.21 之前的脏数据可能源于缓存竞态；执行：

  ```sql
  SELECT * FROM orders WHERE status=1 AND paid_at > NOW() - INTERVAL 1 DAY;
  -- find anomalies and re-trigger activation via /api/order/<id>/reapply
  ```

## 4. OAuth & Callbacks / OAuth 与回调

### 4.1 GitHub callback lands on a white screen

- Are `GITHUB_CLIENT_ID` / `GITHUB_CLIENT_SECRET` set?
- Does `FRONTEND_BASE_URL` share the origin with the callback URL?
- Browser blocking third-party cookies? Fall back to local :3000 for debugging.

### 4.2 Lark callback 400

Add the callback URL to "Security → Redirect URLs" in the Lark developer console; protocol and case must match.

回调域名要在「飞书开发者后台 → 安全设置 → 重定向 URL」里加白；URL 区分协议与大小写。

## 5. Cluster / 集群

### 5.1 A node keeps showing offline

| Field | Meaning | Check |
| --- | --- | --- |
| `status` | 1 alive / 2 abnormal | Try `curl /api/cluster_node/<id>/ping` first |
| `disabled` | Disabled flag | Restore manually in Cluster Settings |
| `last_heartbeat` | Last heartbeat | Considered dead when delta > `CLUSTER_DEAD_PING_INTERVAL (120)` |
| `ping_failures` | Consecutive ping failures | Auto-evicted after > `CLUSTER_MAX_PING_FAILURES (3)` |

### 5.2 New node can't see historical changes

By design: Cluster mode does **not** perform active pull; new nodes only see events from join-time onward. See [Cluster Overview · design trade-offs](/en/decentralization/overview#design-trade-offs-design-trade-offs).

设计取舍：Cluster 模式**不做主动拉取**，新节点只看加入之后的同步事件。

**Fix**: `mysqldump` a baseline from a live node, or temporarily disable `CLUSTER_ENABLED` and let the new node attach to the shared DB.

### 5.3 Data drift between nodes

- Mismatched MySQL timezones (set `TZ`).
- NTP drift > 1s (install `chrony`).
- `CLUSTER_PUSH_INTERVAL (3s)` too short for the volume; bump to 10.

节点 `MySQL` 时区不一致（设置 `TZ`）；NTP 漂移 > 1s（建议 `chrony`）；`CLUSTER_PUSH_INTERVAL (3s)` 太短被队列挤压，调高到 10。

## 6. Collecting Diagnostics / 收集诊断信息

When filing an issue, attach:

提交 issue 时附上：

1. **Version / 版本**: `./one-api-pro --version` (known issues fixed in ≥ v0.0.21).
2. **Deployment / 部署方式**: single instance / Cluster / Docker.
3. **Key env vars / 关键环境变量**: `SQL_DSN` type, whether `REDIS_CONN_STRING` is set, `CHANNEL_*` config.
4. **Repro / 复现步骤**: the call + expected vs actual response.
5. **Log snippet / 日志片段**: `grep -i "<keyword>" logs/$(date +%Y-%m-%d).log`.
6. **Screenshots / 截图** + browser / OS.

Next: [FAQ](/en/faq/faq) · [Glossary](/en/faq/glossary).