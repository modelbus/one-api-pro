---
title: 节点健康
description: 心跳、互相 Ping、节点失联检测与恢复。
category: decentralization
order: 4
---

# 节点健康

> 节点之间怎么知道「我还在」、「你还活着」、「你挂了」。

## 节点状态

每个节点有两种状态：

- **`status=1` alive**：能正常响应 ping
- **`status=2` dead**：连续 ping 失败超过阈值

状态变更通过 **每节点心跳 + 节点间互相 ping** 实现。

## 怎么工作

### 本地心跳

每个节点每 30 秒刷新自己的 `last_heartbeat` 时间戳。

### 互相 Ping

每个节点每 30 秒遍历所有「其他节点」，向每个发 ping：

```
节点 A ──ping──> 节点 B
       <──ok───

成功：A 收到响应 → 把 B 状态记为 alive
失败：累加 B 的 ping_failures，达阈值 → B 状态 = dead
```

死节点不会马上被踢出，**会按更慢的节奏继续探测**（默认 120 秒一次），恢复后自动回到 alive。

## 关键参数

| 变量 | 默认值 | 含义 |
|---|---|---|
| `CLUSTER_DISCOVERY_INTERVAL` | 30 秒 | 心跳 + 互 ping 周期 |
| `CLUSTER_DEAD_PING_INTERVAL` | 120 秒 | 死节点多久 ping 一次 |
| `CLUSTER_MAX_PING_FAILURES` | 3 | 连续失败几次标 dead |
| `CLUSTER_PUSH_INTERVAL` | 3 秒 | Pusher 推送节流（实际 100 ms 内部节流） |

进程内修改这些参数需重启。

## 死节点上的行为

- **Pusher 跳过**：不向死节点推送同步事件
- **手动 ping 仍可用**：后台「Ping」按钮可以强制 ping 一次
- **死节点不影响自己被 ping**：它仍能响应 ping，便于其他节点探测
- **恢复后自动补推**：死节点恢复时，本地 `sync_events` 堆积的事件由 Pusher 自动补推

## 怎么排查节点失联

1. 后台「集群设置 → 节点管理」看节点 `status` 和 `last_heartbeat` 时间
2. 点「Ping」手动 ping 一次，看返回
3. 检查网络：两个节点能否互相访问对方的 `address`
4. 检查 Secret：两边的 `secret_key` 是否一致
5. 看后端日志的 `[集群] ping 仍然失败` / `节点 X 已恢复` 等

## 故障转移

节点失联**不影响用户请求**：用户访问的是入口 Nginx/LB，仍然能命中其他节点。

只有当**所有**节点都挂掉时才会停止服务。所以 Cluster 模式天然是「高可用」架构。

## 常见问题

- **某节点一直 status=2**：检查网络、Secret；看后端启动日志
- **刚加的节点没被发现**：检查 `CLUSTER_SEEDS` 是否指向至少一个可达节点；首次启动后是否能 ping 通 seed
- **Secret 泄露怎么办**：在后台「节点管理 → 编辑」改 secret；下次 ping 自动传给其他节点

## 相关

- [Cluster 概览](./overview)
- [节点管理](./node-management)
- [多节点部署](./deployment)