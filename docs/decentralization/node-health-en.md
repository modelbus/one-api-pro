---
title: Node Health
description: Heartbeat, mutual ping, dead-node detection, and recovery.
category: decentralization
order: 4
---

# Node Health

> How nodes know "I'm still here", "you're still alive", and "you're dead".

## Node state

Each node has two states:

- **`status=1` alive**: responding to ping normally
- **`status=2` dead**: failed consecutive pings beyond the threshold

State changes are driven by **per-node heartbeat + mutual pinging**.

## How it works

### Local heartbeat

Every node refreshes its own `last_heartbeat` every 30 seconds.

### Mutual ping

Every 30 seconds, each node iterates "the other nodes" and pings each:

```
Node A ──ping──> Node B
       <──ok───

Success: A receives response → mark B as alive
Failure: increment B's ping_failures; when it hits the threshold → status = dead
```

A dead node isn't kicked out. It's probed at a slower rate (default 120 s) and auto-recovers when it comes back.

## Key parameters

| Variable | Default | Meaning |
|---|---|---|
| `CLUSTER_DISCOVERY_INTERVAL` | 30 s | Heartbeat + ping cycle |
| `CLUSTER_DEAD_PING_INTERVAL` | 120 s | How often to ping a dead node |
| `CLUSTER_MAX_PING_FAILURES` | 3 | Failures before marking dead |
| `CLUSTER_PUSH_INTERVAL` | 3 s | Pusher throttle (actual 100 ms internal) |

Process-level changes to these require a restart.

## What happens on a dead node

- **Pusher skips it**: no events pushed to dead nodes
- **Manual ping still works**: Admin → "Ping" button forces a probe
- **Dead nodes still respond to ping**: so other nodes can detect recovery
- **Catch-up on recovery**: events that piled up locally (`sync_events` with `pushed=0`) are auto-pushed after recovery

## Troubleshooting a lost node

1. Admin → Cluster → Node Management: check `status` and `last_heartbeat`
2. Click "Ping" for a manual probe — see the response
3. Network: can the two nodes reach each other's `address`?
4. Secret: do both sides' `secret_key` match?
5. Backend logs: look for `[cluster] ping still failing` / `node X recovered`

## Failover

A lost node doesn't break user requests: traffic goes through Nginx/LB and is served by remaining nodes.

Service only stops when **all** nodes are down. That's the natural HA property of Cluster mode.

## FAQ

- **Node stuck at status=2**: check network and Secret; see backend startup logs
- **New node never discovered**: confirm `CLUSTER_SEEDS` points to at least one reachable node; first start should ping the seed successfully
- **Secret leaked**: edit it in Admin → Node Management; next ping auto-propagates

## Related

- [Cluster Overview](./overview)
- [Node Management](./node-management)
- [Multi-node Deployment](./deployment)