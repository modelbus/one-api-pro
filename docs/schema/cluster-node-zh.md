---
title: 集群节点
description: Cluster 模式下注册的节点及健康状态。
category: schema
order: 12
---

# 集群节点（ClusterNode）

## 这是什么

`ClusterNode` 是 **Cluster 多节点部署模式下注册到中心表的一行记录**。每个节点启动时会向其他节点注册自己（包括节点 ID、地址、密钥），并持续心跳。

如果你只部署单机，**完全不需要关心**本节。

## 在哪里看到

- **后台 → 集群设置**：节点列表、状态、最后心跳时间
- **监控**：每个节点都会通过 pusher 主动推送数据变更到其他节点

## 关键字段

| 字段 | 含义 | 设置后影响 |
|---|---|---|
| 节点 ID | 全集群唯一标识 | 修改后该节点被认为「换了身份」 |
| 节点名称 | 展示用 | 修改不影响通信 |
| 状态 | 启用 / 禁用 | 禁用后不再接收其他节点的推送 |
| 地址 | `host:port` | 节点间 HTTP 推送目标；写错则推送失败 |
| 密钥 | 节点间互信的 Bearer Token | 必须与集群内其他节点一致 |
| 最后心跳时间 | Unix 秒 | 长时间未更新 → 标红「失联」 |

## 多节点部署

如果你计划部署 N 个节点：

1. 每个节点的 `CLUSTER_NODE_ID` 必须唯一
2. 每个节点的 `CLUSTER_NODE_SECRET` 必须相同
3. 每个节点要知道其他所有节点的 `地址`
4. 节点间通过 HTTP 主动推送数据变更，无中心数据库

详见 [去中心化](/decentralization/overview) 与 [多节点部署](/decentralization/deployment)。

## 相关页面

- [Cluster 概览](/decentralization/overview)
- [节点管理](/decentralization/node-management)
- [配置同步](/decentralization/config-sync)
- [节点健康](/decentralization/node-health)
- [多节点部署](/decentralization/deployment)
- [集群设置](/decentralization/cluster-settings)

## 相关 API

- `GET /api/cluster_node/` — 列表
- `POST /api/cluster_node/` — 新增
- `PUT /api/cluster_node/` — 更新
- `POST /api/cluster_node/:id/enable` — 启用
- `GET /api/cluster_node/ping/:id` — Ping 测试