---
title: Cluster 概览
description: "去中心化多活集群的设计目标、约束与适用场景。"
category: decentralization
order: 1
---

# Cluster 概览

> 去中心化多活集群的设计目标、约束与适用场景。
> Design goals, constraints and suitable use cases of the decentralized multi-active cluster.

## 目标 / Goals

One API Pro 的 **Cluster 模式** 提供去中心化的多节点多活部署，核心目标是：

The **Cluster mode** in One API Pro provides a decentralized, multi-active multi-node deployment. Its core goals are:

- **不共享数据库**：每个节点持有独立的 MySQL 与 Redis；节点间通过 HTTP 主动推送同步；
  **No shared database**: every node owns its own MySQL and Redis; nodes actively push sync events to each other over HTTP.
- **零侵入业务**：通过 GORM 回调捕获业务表变更，无需改业务代码；
  **Zero business intrusion**: data changes are captured via GORM callbacks — no business-code changes needed.
- **多活就近接入**：跨地域 / 跨机房部署，本地节点就近服务，降低延迟；
  **Multi-active local access**: deploy across regions / data centers, serve users from the nearest node.
- **冲突可收敛**：基于 `updated_at` 的最后写入胜出，保证最终一致；
  **Convergent conflicts**: last-writer-wins by `updated_at` ensures eventual consistency.

## 何时使用 / When to use

| 场景 / Scenario | 推荐 / Recommendation |
| --- | --- |
| 单机房中小流量 | 单实例即可，无需 Cluster |
| 多机房 / 多区域 / 跨地域容灾 | ✅ Cluster 模式 |
| 业务对延迟敏感，希望就近接入 | ✅ Cluster 模式 |
| 已有 K8s 多副本 + 共享 DB | 仍用多实例共享 DB 方案（参见 `install/docker-compose`），不需要 Cluster |
| 强一致 / 分布式事务 | ❌ 不适合；One API Pro 不实现跨节点事务 |

## 架构一览 / Architecture at a glance

```text
              ┌─────────────┐
              │  Nginx/LB   │   (单一入口, ip_hash 负载均衡 / single entry, ip_hash LB)
              └──────┬──────┘
                     │
       ┌─────────────┼─────────────┐
       │             │             │
 ┌─────┴─────┐ ┌─────┴─────┐ ─────┴─────┐
 │  Node A   │ │  Node B   │ │  Node C   │
 │ one-api   │ │ one-api   │ │ one-api   │
 │ + MySQL   │ │ + MySQL   │ │ + MySQL   │
 │ + Redis   │ │ + Redis   │ │ + Redis   │
 └─────┬─────┘ └─────┬─────┘ └─────┬─────┘
       │             │             │
       └────── HTTP push of sync events ──────┘
```

每个节点都是平等的：任何节点的数据变更都会被主动推送到所有存活节点。
Every node is equal: any data change on a node is actively pushed to all alive nodes.

## 核心特性 / Core characteristics

- **去中心化 / Decentralized** — 节点间对等，无中心协调器；no central coordinator.
- **零侵入 / Zero-invasion** — 通过 GORM callbacks 自动捕获 `INSERT` / `UPDATE` / `DELETE`。
  Zero invasion: `INSERT` / `UPDATE` / `DELETE` are captured automatically by GORM callbacks.
- **异步推送 / Async push** — 同步发生在后台 goroutine，不阻塞主流程；不会拖累请求路径。
  Async push: sync runs in a background goroutine, never blocking the main flow.
- **冲突解决 / Conflict resolution** — 接收方比较 `updated_at`，仅写入更新的版本；最后写入胜出。
  Conflict resolution: the receiver compares `updated_at`; only the newer version is written (last-writer-wins).
- **限流同步 / Rate-limit sync** — `channel_counters`（渠道并发 / RPM）按节点维度同步，全局限流状态可跨节点聚合。
  Rate-limit sync: `channel_counters` (channel concurrency / RPM) sync per node, so global rate-limit state can be aggregated across nodes.
- **单节点兼容 / Single-node compatible** — 不配置 `CLUSTER_*` 时以单节点模式运行，无副作用。
  Single-node compatible: without `CLUSTER_*` env vars, the system runs in plain single-node mode with no side effects.

## 同步范围 / Sync scope

| 表 / Table | 是否同步 / Synced? | 备注 / Notes |
| --- | --- | --- |
| `users` | ✅ | 用户账号 / user accounts |
| `tokens` | ✅ | API tokens |
| `channels` | ✅ | 渠道 / provider channels |
| `abilities` | ✅ | 渠道能力 / channel abilities |
| `options` | ✅ | 系统设置 / system settings |
| `redemptions` | ✅ | 兑换码 / redemption codes |
| `plans` | ✅ | 套餐 / subscription plans |
| `user_plans` | ✅ | 用户订阅 / user subscriptions |
| `plan_usages` | ✅ | 套餐用量 / plan usage |
| `channel_counters` | ✅ | 渠道限流计数 / channel rate-limit counters |
| `cluster_nodes` | 🔄 由发现机制维护 | maintained by the discovery mechanism, not data sync |
| `logs` | ⚠️ 受 `CLUSTER_SYNC_LOGS` 控制 | controlled by `CLUSTER_SYNC_LOGS` |

## 设计取舍 / Design trade-offs

### 主动推送，不做主动拉取 / Push-only, no active pull

数据同步完全依赖 GORM 回调 + HTTP 主动推送，**不实现跨节点主动拉取**。

Data sync relies entirely on GORM callbacks + HTTP push; **active cross-node pull is not implemented**.

原因 / Why:

1. **业务侵入**：拉取需知道每张表的业务唯一字段，会污染业务代码；
   **Business intrusion** — pull would require knowing each table's business-unique field, polluting business code.
2. **主键冲突**：跨节点 `auto_increment` 不同（不同 `auto_increment_offset`），用源 ID 会破坏 offset 设计；
   **Primary-key conflicts** — auto-increment IDs differ across nodes (different `auto_increment_offset`); using the source's ID would break the offset design.
3. **复杂度**：维护成本高，可靠性收益有限；
   **Complexity** — high maintenance cost for limited reliability gain.
4. **推送足够**：覆盖约 95% 的正常场景（节点在线、流量正常）；
   **Push is enough** — covers ~95% of normal scenarios (alive nodes, normal traffic).

### 已知限制 / Known limits

- **节点离线期间产生的变更不会被回填**；节点恢复后需要从存活节点 `mysqldump` 同步一次。
  Changes made while a node is offline are **not back-filled**; the node must be re-seeded via `mysqldump` from a live node after coming back.
- 新加入节点只能看到加入之后的变更历史；
  New nodes only see changes from the moment they join — no history.
- `logs` 表数据量较大时建议关闭（`CLUSTER_SYNC_LOGS=false`）；
  Disable `logs` sync (`CLUSTER_SYNC_LOGS=false`) if the table grows too large.

## 与多实例共享 DB 方案的区别 / Cluster vs shared-DB

| 维度 / Aspect | Cluster（去中心化） | 多实例共享 DB（传统） |
| --- | --- | --- |
| 数据库 / DB | 每节点独立 MySQL | 共享 MySQL |
| Redis | 每节点独立 | 共享 |
| 一致性 / Consistency | 最终一致（基于 `updated_at`） | 强一致（共享 DB） |
| 跨地域延迟 | 低（本地节点服务） | 高（跨地域往返 DB） |
| 节点离线容忍 | 离线期间数据丢失，需手动补回 | 节点随时可下线，数据不丢 |
| 适合规模 / Fit | 跨地域、多机房 | 单机房多副本 |

下一步 / Next: [节点管理](/zh/decentralization/node-management) · [配置同步](/zh/decentralization/config-sync) · [节点健康](/zh/decentralization/node-health)。
