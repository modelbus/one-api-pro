---
title: 配置同步
description: 数据如何在节点间同步、什么时候会冲突、冲突如何收敛。
category: decentralization
order: 3
---

# 配置同步

> 节点 A 改了一条记录，节点 B、C、D 多久能看见？答：通常 < 1 秒。

## 同步流程

```
节点 A 写入数据库
    ↓
GORM 回调自动捕获变更
    ↓
写入 sync_events 临时表
    ↓
Pusher 协程批量推送到每个存活节点
    ↓
节点 B/C/D 的 Applier 收到事件并应用
```

任一节点改数据，其他节点通常 **< 1 秒**内可见。

## 同步哪些表

- `users` / `tokens`：账号与 Token
- `channels` / `abilities`：渠道
- `options`：系统设置
- `redemptions`：兑换码
- `plans` / `user_plans` / `plan_usages`：套餐 / 订阅 / 用量
- `channel_counters`：渠道限流计数
- `logs`：调用日志（可通过 `CLUSTER_SYNC_LOGS=false` 关掉）

**不同步**：`cluster_nodes` / `sync_events`（内部表，由发现机制维护）。

## 冲突怎么处理

如果两个节点**几乎同时**改同一条记录，系统按「最后写入胜出」：

```
A 改记录的 updated_at = 10:00:01.100
B 改记录的 updated_at = 10:00:01.050
→ A 的版本最终保留
```

规则：

- 接收方收到事件，比较入站记录的 `updated_at` 与本地
- 入站更新才覆盖；本地更新则丢弃入站

> 所有节点的**系统时钟必须基本一致**（建议启用 NTP）。漂移过大会导致本应被采纳的更新被错误丢弃。

## 节点离线期间产生的变更

**不会自动回填**。比如节点 B 离线 1 小时，期间节点 A 改了 100 条记录；B 恢复后只能看到恢复后的新变更，看不到离线期间的。

修复方法：从存活节点手工 `mysqldump` 同步一次。

## 节点间推送的请求头

每次推送会带：

- `X-Cluster-Secret`：目标节点的 secret（校验通过才接收）
- `X-Cluster-Node-Id`：发送方节点 ID
- 超时 5 秒

## 同步出问题了怎么排查

1. **A 改了 B 没收到**：
   - 看 B 后端日志是否有 sync 事件
   - 在 A 端手动 Ping B，看是否通
   - 检查两边的 `secret_key` 是否一致
2. **数据冲突，丢了一条更新**：
   - 检查两节点时钟是否同步（开启 NTP）
   - 看数据库 `sync_events` 表是否有积压（`pushed=0` 太多 = 推送阻塞）
3. **同步延迟很大（> 几秒）**：
   - 看 `sync_events` 积压
   - 看 Pusher 协程是否在运行（看进程监控）
   - 看节点间网络延迟

## 哪些场景关掉同步更快

- 日志表巨大：在节点环境变量设 `CLUSTER_SYNC_LOGS=false`，日志只存在本地节点
- 节点只读：把当前节点的 secret 设为随机值（拒绝接收）

## 相关

- [Cluster 概览](./overview)
- [节点管理](./node-management)
- [节点健康](./node-health)