---
title: 节点管理
description: 集群节点的注册、查询、启用与删除。
category: decentralization
order: 2
---

# 节点管理

> Cluster 模式下怎么新增 / 查看 / 启用 / 禁用节点。

## 在哪里

后台 → 集群设置 → 节点管理。

## 列表看到什么

每行展示：

- 节点 ID / 名称
- 地址（`https://node-b.example.com`）
- 状态（启用 / 失败 / 禁用）
- 最近心跳时间
- 连续 ping 失败次数
- 操作按钮

标红的「失联」表示连续多次 ping 失败。

## 节点字段

| 字段 | 含义 |
|---|---|
| `node_id` | 节点编号 1–49；与 `CLUSTER_NODE_ID` / MySQL `auto_increment_offset` 一致 |
| `node_name` | 展示名 |
| `address` | 公网 URL（含协议前缀） |
| `secret_key` | 本节点对其他节点的认证 secret |
| `status` | `1`=启用 / `2`=失败 / `3`=被管理员禁用 |
| `disabled` | 管理员意图的软删除标记 |
| `last_heartbeat` | 最近心跳时间戳 |
| `ping_failures` | 连续 ping 失败次数 |

## 自注册（推荐）

每个节点**首次启动时自动注册自己**。无需在后台手动添加：

- 节点启动 → 写自身记录到 `cluster_nodes` 表
- 周期任务刷新心跳、地址
- 其他节点通过 ping 也能间接发现它

这是推荐方式：节点只需关心自己的 `CLUSTER_NODE_ID` 和 `CLUSTER_NODE_SECRET`。

## 手动新增（仅特殊场景）

有些场景需要手动登记：

1. 后台 → 节点管理 → 「新增节点」
2. 填：
   - 节点 ID（1–49）
   - 节点名称（任意）
   - 地址（对方的公网 URL，如 `https://node-b.example.com`）
   - Secret（对方的 `CLUSTER_SECRET` 值）
3. 保存

> 自注册已覆盖绝大多数场景；只有调试 / 临时让节点互认时才需要手动。

## 启用 / 禁用 / 删除

| 操作 | 在哪里 | 效果 |
|---|---|---|
| 启用 | 行内「启用」 | `disabled=false` + 重置失败计数 |
| 禁用（软删除） | 行内「禁用」 | `disabled=true`；不再推送事件，但仍响应 ping |
| 物理删除 | 手动 SQL | `DELETE FROM cluster_nodes WHERE node_id = ?` |

通常只需要「禁用」即可（可恢复）。物理删除只用于节点彻底下线的场景。

## Secret 是干什么的

每个节点对其他节点的请求头 `X-Cluster-Secret` 必须等于**目标节点**的 secret。Secret 泄露不影响其他节点（每个节点独立）。

### 首次启动

使用 `CLUSTER_SECRET` 环境变量作为 `secret_key` 初值。

### 后续启动

从数据库 `cluster_nodes.secret_key` 读取，不再依赖环境变量。

### 轮换

通过后台修改 → 下次 ping 会把新值同步给其他节点。

## 怎么排查节点失联

1. 后台看节点「最后心跳时间」是否很久没更新
2. 点「Ping」手动 ping 一次，看返回
3. 检查网络：两个节点能不能互相访问对方的 `address`
4. 检查 Secret：两边的 `secret_key` 是否一致
5. 看 [节点健康](./node-health) 的故障转移文档

## 常见问题

- **节点列表里看不到自己**：检查 `CLUSTER_NODE_ID` 环境变量是否设置；看后端日志有没有 `cluster_nodes` 表的写入失败
- **新增后两个节点互不认识**：两边的 `address` 是否互相能访问；`secret_key` 是否填对
- **重启用节点不工作**：重启节点进程，让它重写心跳

## 相关

- [Cluster 概览](./overview)
- [节点健康](./node-health)
- [多节点部署](./deployment)