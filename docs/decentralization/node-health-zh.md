---
title: 节点健康
description: "心跳、Ping 与故障转移策略。"
category: decentralization
order: 4
---

# 节点健康

> 心跳、Ping 与故障转移策略。

节点健康由两个独立机制保障：

1. **本地心跳**：每个节点每周期把自己 `cluster_nodes.last_heartbeat` 刷成当前时间；
2. **互 ping**：节点间互相 ping，根据响应情况更新对方 `status` 与 `ping_failures`。

## 关键参数

| 变量| 默认| 作用|
| --- | --- | --- |
| `CLUSTER_DISCOVERY_INTERVAL` | `30` | `discoverOnce` 调用间隔（秒），即心跳 + 互 ping 周期 |
| `CLUSTER_DEAD_PING_INTERVAL` | `120` | 已失败节点多久 ping 一次（秒） |
| `CLUSTER_MAX_PING_FAILURES` | `3` | 连续失败多少次后标记节点 `status=2` |
| `CLUSTER_PUSH_INTERVAL` | `3` | Pusher 协程的通知节流间隔（秒；实际为 100 ms 内部节流） |

> 所有参数在 `cluster/config.go::LoadConfig` 加载；进程内变更需重启。

## discoverOnce 周期

`cluster/node.go::discoverOnce` 是节点健康的核心入口：

```text
┌─────────────────────────────────────────────────────────┐
│                       每 30 秒触发                       │
└─────────────────────────────────────────────────────────┘
                           │
              1) 刷新本机 last_heartbeat = now
                           │
              2) 遍历所有非本机节点：
                 ┌────────────────┬───────────────────┐
                 │  status=1 (alive)│  status=2 (dead) │
                 ├────────────────┼───────────────────┤
                 │ pingAliveNode  │ pingDeadNode       │
                 └────────────────┴───────────────────┘
```

`StartDiscovery` 是个无脑 `for { sleep; discoverOnce }` 循环：

```go
func StartDiscovery(db *gorm.DB) {
    for {
        time.Sleep(time.Duration(DiscoveryInterval) * time.Second)
        discoverOnce(db)
    }
}
```

## pingAliveNode：存活节点的健康检测

```text
HTTP POST {node.Address}/api/cluster/ping  (5 秒超时)
            │
   ┌────────┴─────────┐
   │ 成功             │ 失败 / 非 Success
   ├──────────────────┼──────────────────
   │ status=1         │ ping_failures++
   │ ping_failures=0  │
   │ last_heartbeat=now│  ┌─ ≥ MaxPingFailures?
   │ merge resp.Nodes │  │  ├ 是 → status=2 + SysError 日志
   └──────────────────┘  │  └ 否 → 仅累加失败计数
                        └───────────────
```

`mergeDiscoveredNodes` 把 ping 响应里携带的节点列表合并进本地 — 这是 **传递性发现** 的核心：节点 D 可能通过 B → A → D 间接发现。

## pingDeadNode：失败节点的恢复探测

为减少对已失败节点的无谓请求，dead 节点的 ping 间隔更长（默认 120 秒）：

```text
now - node.last_ping_attempt < DeadPingInterval   → 跳过
否则：
   POST /api/cluster/ping
        ├─ 失败 → 仅刷新 last_ping_attempt
        └─ 成功 → status=1, ping_failures=0, last_heartbeat=now
                    SysLogf("节点 %s 已恢复")
```

`last_ping_attempt` 字段由 `pingAliveNode` / `pingDeadNode` 共同维护。

## 状态机

```text
       ┌─────────────────────────────────────┐
       │              status=1 (alive)       │
       │   pingFailures=0..MaxPingFailures-1 │
       └────────┬───────────────┬────────────┘
                │ ping OK       │ ping fail × MaxPingFailures
                ▼               ▼
       last_heartbeat=now   status=2 (dead)
                                ▲
                                │ ping OK
        ┌───────────────────────┘
        │
        │ 手动 disable
        ▼
   disabled=true, status=2
        ▲
        │ admin: POST /api/cluster_node/:id/enable
        └───── disabled=false, status=1, ping_failures=0
```

> 失败的节点 **不会** 被自动删除 — 仅 `status=2`，恢复后自动回到 `status=1`。

## 失败节点上的行为

- **Pusher 跳过**：`GetAliveNodesForSync` 只取 `status=1 AND disabled=false`，不向失败节点推送；
- **手动 ping 仍可用**：管理员可通过 `GET /api/cluster_node/ping/:id` 强制 ping 一次；
- **不影响 `/api/cluster/ping`**：失败节点本身仍能响应 ping（便于对方知道它在线）；
- **Pusher 通知**：失败节点在事件表里堆积（`pushed=0`），恢复后由 Pusher 自动补推 + 启动 catch-up 兜底；

## 排障

| 现象| 排查|
| --- | --- |
| 节点一直 `status=2` | 检查网络与 `X-Cluster-Secret`（ping 用的是对方 secret）；看启动日志的 `[集群] ping 仍然失败` |
| 节点 ping_failures 累加 | `cluster_discover_interval` 与 `cluster_dead_ping_interval` 配置是否合理 |
| 新加节点未被发现 | 是否配置了 `CLUSTER_SEEDS` 指向至少一个可达节点；首次启动后是否能 ping 通 seed |
| Secret 泄露 | 管理员后台 `PUT /api/cluster_node/` 改 secret；下次 ping 自动传播 |

下一步 / Next: [多节点部署](/zh/decentralization/deployment) · [配置同步](/zh/decentralization/config-sync)。

