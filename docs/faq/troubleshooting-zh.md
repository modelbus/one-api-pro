---
title: 故障排查
description: "渠道 401、计费异常、订单不激活等问题的定位。"
category: faq
order: 2
---

# 故障排查

> 渠道 401、计费异常、订单不激活等问题的定位。
> Diagnose channel 401s, billing anomalies, stuck orders.

按现象归类，先看「诊断信息」，再走「排查路径」，最后给「修复方案」。

Categorize by symptom, gather diagnostics first, then trace the cause, then apply a fix.

## 1. 渠道相关 / Channel Issues

### 1.1 渠道测试 200 但 `/v1/*` 返回 401

**诊断 / Diagnostics**：

```bash
# 触发测试
curl -X POST http://localhost:3000/api/channel/test/<id> -b cookies.txt
# 查看响应里的 x-request-id 与上游回包
```

**排查 / Trace**：

| 现象 / Symptom | 原因 / Cause | 修复 / Fix |
| --- | --- | --- |
| 测试 OK，调用 401 | `base_url` 错（如写了 `/v1/` 重复路径） | 去掉 `base_url` 末尾 `/v1`；OpenAI Base URL 一般 `https://api.openai.com` |
| 偶发 401 后正常 | 上游触发风控 | 切换账号 / 调整 `RATE_LIMIT` / 用 `cooldown` 隔开 |
| 多 Provider 共用 Key 时 | `abilities` 表里 `group` 没匹配 | 在「用户分组」中加白 |

### 1.2 渠道一直「冷却」/ 在 cooldown 永久不退

- **诊断**：`channels.cooldown_until` 字段一直大于 `now`。
- **原因**：上游持续报错；`CHANNEL_DEFAULT_COOLDON_SECDS=60` 起步，每次失败翻倍直至 `CHANNEL_MAX_COOLDOWN_SECONDS=600`。
- **修复**：手动改回 `cooldown_until=0`；排查上游报错根因。

### 1.3 渠道 `weight` 改了不生效

`BATCH_UPDATE_ENABLED=true` 时权重聚合写入有窗口（默认 5s）。要么关闭批处理，要么重启接口 `SLOW=` 相关 cache key。

If `BATCH_UPDATE_ENABLED=true`, weights are written through a batched window (default 5s). Either disable batching, or restart the relevant cache key.

## 2. 计费与缓存 / Billing & Cache

### 2.1 反复 403 `insufficient_user_quota`

**v0.0.21 之前的常见根因**：Redis `user_quota:<id>` 缓存与 DB 不一致。已修复：

Common pre-v0.0.21 cause: Redis `user_quota:<id>` cache drifts from DB. Fixed by:

- `IncreaseUserQuota` 直接写库后同步刷 Redis；
  `IncreaseUserQuota` writes DB then refreshes Redis synchronously.
- `preConsumeQuota` 的「余额充足免预扣」判定前移到 Redis 扣减之前；
  The "trusted, skip pre-consume" guard moved before the Redis decrement.
- `userQuotaLowWaterMark = 50_000`（曾为 500），避免高单价模型预扣落进「缓存漂移死区」。

**遇到此问题请检查版本**：

```bash
./one-api-pro --version
```

升级到 ≥ v0.0.21；若仍出现请清掉 `user_quota:<id>` 缓存重启：

Upgrade to ≥ v0.0.21. If it persists, clear the `user_quota:<id>` cache and restart:

```bash
redis-cli DEL user_quota:1 user_quota:2 ...
```

### 2.2 令牌「永不过期」被误判为过期

**修复历史**：v0.0.21 起 sentinel 从「仅 `-1` 表示永不过期」统一为「`<=0` 均视为永不过期」（`model/token.go` / `controller/token.go`）；前端 `buildTokenExpiredTime` 始终上报 `-1`。

Historical fix: since v0.0.21 the sentinel is "`<=0` means never expire" (`model/token.go` / `controller/token.go`); the frontend `buildTokenExpiredTime` always submits `-1`.

升级后历史脏数据 `expired_time=0` 不会再被误报为已过期。

After upgrade, dirty data `expired_time=0` is no longer flagged as expired.

### 2.3 `tiktoken` 计数 panic

**v0.0.21 之前**：防火墙下 HTTP 下载 BPE 超时 → fallback 不完整 → `nil` 编码器上调用 `Encode()` 触发 panic。

**v0.0.21 修复**：4 个 BPE 编码文件（cl100k_base / o200k_base / p50k_base / r50k_base）通过 `//go:embed` 内嵌进二进制，启动零网络依赖。

```text
httpBpeLoader 三级查找：
  1. embedded BPE（默认，开箱即用）
  2. TIKTOKEN_CACHE_DIR（高级覆盖，可跳过内嵌与 HTTP）
  3. HTTP 下载（兜底，应对未来 URL 变更）
```

升级到 v0.0.21+ 即可解决；高级用户可设置 `TIKTOKEN_CACHE_DIR` 跳过内嵌。

## 3. 订单与异步通知 / Orders & Async Notify

### 3.1 异步通知校验失败

| 渠道 / Channel | 校验 / Verification | 字段来源 / Field source |
| --- | --- | --- |
| WeChat | MD5 签名 + `return_code=SUCCESS` | `payment.wechat.config.api_key` |
| Alipay | RSA2 验签 | `payment.alipay.config.private_key/public_key` |
| Bank | 手工对账（管理员「标记已支付」） | 走后台人工流程 |

**诊断 / Diagnostics**：

```bash
# 1. 看回调日志
tail -f logs/$(date +%Y-%m-%d).log | grep -i notify

# 2. 主动查订单状态
curl http://localhost:3000/api/order/<id> -b cookies.txt
```

**典型问题 / Common pitfalls**：

- WeChat API Key 改了但「支付配置」页没保存 → 回调签名失败。
- Alipay `notify_url` 与回调地址公网可达，但**不在 ICP 备案域**也会被拒。
- 「银行转账」无回调，必须管理员在订单详情手动「标记已支付」。

### 3.2 订单「已支付」但套餐没激活

- 检查 `orders.status` 是否 = 1；若是但 `user_plans` 未增加：
  Check `orders.status == 1`; if so but `user_plans` not increased:
- v0.0.21 之前的脏数据可能源于缓存竞态；执行：
  Pre-v0.0.21 stale data may come from cache races; run:

  ```sql
  SELECT * FROM orders WHERE status=1 AND paid_at > NOW() - INTERVAL 1 DAY;
  -- 找出异常后用控制器 /api/order/<id>/reapply 重新触发激活
  ```

## 4. OAuth 与回调 / OAuth & Callbacks

### 4.1 GitHub 回调后白屏

- `GITHUB_CLIENT_ID` / `GITHUB_CLIENT_SECRET` 是否设置；
- `FRONTEND_BASE_URL` 是否与回调 URL 同源；
- 浏览器第三方 Cookie 拦截；改用本地 3000 调试。

### 4.2 飞书回调 400

回调域名要在「飞书开发者后台 → 安全设置 → 重定向 URL」里加白；URL 区分协议与大小写。

Add the callback URL to "Security → Redirect URLs" in the Lark developer console; protocol and case must match.

## 5. 集群 / Cluster

### 5.1 节点一直显示「离线」

| 字段 / Field | 含义 / Meaning | 排查 / Check |
| --- | --- | --- |
| `status` | 1 在线 / 2 异常 | 异常节点先 `curl /api/cluster_node/<id>/ping` |
| `disabled` | 是否被禁用 | 集群设置里手动恢复 |
| `last_heartbeat` | 上次心跳 | 间隔 > `CLUSTER_DEAD_PING_INTERVAL (120)` 即视为死亡 |
| `ping_failures` | 连续 ping 失败计数 | > `CLUSTER_MAX_PING_FAILURES (3)` 后自动剔除 |

### 5.2 新节点看不到历史变更

设计取舍：Cluster 模式**不做主动拉取**，新节点只看加入之后的同步事件。详见 [Cluster 概览 · 设计取舍](/zh/decentralization/overview#设计取舍-design-trade-offs)。

Design choice: Cluster mode does **not** do active pull; new nodes only see events from join-time onward. See [Cluster Overview · design trade-offs](/en/decentralization/overview#design-trade-offs-design-trade-offs).

**修复**：从存活节点 `mysqldump` 拉一次 baseline；或临时关闭 `CLUSTER_ENABLED`，让新节点直接接 DB。

### 5.3 节点间数据漂移

- 节点 `MySQL` 时区不一致（设置 `TZ`）；
  Mismatched MySQL timezones (set `TZ`).
- 节点 NTP 漂移 > 1s（建议 `chrony`）；
  NTP drift > 1s (install `chrony`).
- `CLUSTER_PUSH_INTERVAL (3s)` 太短被队列挤压，调高到 10。
  `CLUSTER_PUSH_INTERVAL (3s)` is too short for the queue; bump to 10.

## 6. 收集诊断信息 / Collecting Diagnostics

提交 issue 时附上：

When filing an issue, attach:

1. **版本 / Version**：`./one-api-pro --version`（≥ v0.0.21 已知问题已修复）；
2. **部署方式 / Deployment**：单实例 / Cluster / Docker；
3. **关键环境变量 / Key env vars**：`SQL_DSN` 类型、`REDIS_CONN_STRING` 是否设置、`CHANNEL_*` 配置；
4. **复现步骤 / Repro**：调用 + 期望 vs 实际响应；
5. **日志片段 / Log snippet**：`grep -i "<关键字>" logs/$(date +%Y-%m-%d).log`；
6. **截图 / Screenshots** + 浏览器 / OS。

下一步 / Next: [常见问题](/zh/faq/faq) · [术语表](/zh/faq/glossary)。