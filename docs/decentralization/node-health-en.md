---
title: Node Health
description: "Heartbeat, ping and failover strategy."
category: decentralization
order: 4
---

# Node Health

> Heartbeat, ping and failover strategy.

Node health is maintained by two independent mechanisms:

1. **Local heartbeat**: every node refreshes its own `cluster_nodes.last_heartbeat` every cycle.
2. **Mutual ping**: peers ping each other and update the remote row's `status` and `ping_failures` based on the response.

## Key tunables

| Variable | Default | Effect |
| --- | --- | --- |
| `CLUSTER_DISCOVERY_INTERVAL` | `30` | `discoverOnce` interval (seconds) — heartbeat + mutual-ping period |
| `CLUSTER_DEAD_PING_INTERVAL` | `120` | Ping interval for failed nodes (seconds) |
| `CLUSTER_MAX_PING_FAILURES` | `3` | Consecutive failures before a node is marked `status=2` |
| `CLUSTER_PUSH_INTERVAL` | `3` | Pusher notification throttle (seconds; actual internal throttle is 100 ms) |

> All parameters are loaded in `cluster/config.go::LoadConfig`; changes require a restart.

## discoverOnce cycle

`cluster/node.go::discoverOnce` is the core entry point for node health:

```text
┌─────────────────────────────────────────────────────────┐
│                      every 30 seconds                    │
└─────────────────────────────────────────────────────────┘
                           │
              1) refresh local last_heartbeat = now
                           │
              2) iterate all non-self nodes:
                 ┌────────────────┬───────────────────┐
                 │  status=1 (alive)│  status=2 (dead) │
                 ├────────────────┼───────────────────┤
                 │ pingAliveNode  │ pingDeadNode       │
                 └────────────────┴───────────────────┘
```

`StartDiscovery` is a simple `for { sleep; discoverOnce }` loop:

```go
func StartDiscovery(db *gorm.DB) {
    for {
        time.Sleep(time.Duration(DiscoveryInterval) * time.Second)
        discoverOnce(db)
    }
}
```

## pingAliveNode: health check on alive nodes

```text
HTTP POST {node.Address}/api/cluster/ping  (5-second timeout)
            │
   ┌────────┴─────────┐
   │ success          │ failure / !Success
   ├──────────────────┼──────────────────
   │ status=1         │ ping_failures++
   │ ping_failures=0  │
   │ last_heartbeat=now│  ┌─ ≥ MaxPingFailures?
   │ merge resp.Nodes │  │  ├ yes → status=2 + SysError log
   └──────────────────┘  │  └ no  → bump failure counter
                        └───────────────
```

`mergeDiscoveredNodes` merges the node list carried by the ping response into the local DB — this is the core of **transitive discovery**: D may learn about A transitively through B.

## pingDeadNode: recovery probe on failed nodes

To reduce wasted requests on failed nodes, dead nodes are pinged less often (default 120 seconds):

```text
now - node.last_ping_attempt < DeadPingInterval   → skip
otherwise:
   POST /api/cluster/ping
        ├─ failure → only refresh last_ping_attempt
        └─ success → status=1, ping_failures=0, last_heartbeat=now
                     SysLogf("node %s recovered")
```

`last_ping_attempt` is maintained jointly by both ping handlers.

## Status state machine

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
        │ admin: DELETE /api/cluster_node/:id
        ▼
   disabled=true, status=2
        ▲
        │ admin: POST /api/cluster_node/:id/enable
        └───── disabled=false, status=1, ping_failures=0
```

> Failed nodes are **not** auto-deleted — they only flip to `status=2` and auto-resurrect when the network is back.

## Behaviour toward failed nodes

- **Pusher skips**: `GetAliveNodesForSync` filters to `status=1 AND disabled=false`, so failed nodes don't receive events.
- **Manual ping is always available**: admins can force a ping via `GET /api/cluster_node/ping/:id`.
- **Failed nodes still respond to pings** so peers can detect recovery.
- **Pusher queues events** for failed nodes (`pushed=0`); on recovery the Pusher resumes, and a startup catch-up handles rows queued while the node was down.

## Troubleshooting

| Symptom | Check |
| --- | --- |
| Node stuck at `status=2` | Network + `X-Cluster-Secret` (ping uses the peer's secret); check startup logs for `[集群] ping 仍然失败` |
| `ping_failures` keeps growing | Are `cluster_discover_interval` and `cluster_dead_ping_interval` reasonable? |
| Newly added node never discovered | `CLUSTER_SEEDS` must point at one reachable node; verify first boot can reach the seed |
| Secret leaked | Update via the admin UI `PUT /api/cluster_node/`; next ping auto-propagates |

Next: [Multi-node Deployment](/en/decentralization/deployment) · [Config Sync](/en/decentralization/config-sync).
