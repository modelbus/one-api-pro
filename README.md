

<p align="center">
  <a href="https://github.com/songquanpeng/one-api"><img src="https://raw.githubusercontent.com/songquanpeng/one-api/main/web/default/public/logo.png" width="150" height="150" alt="one-api logo"></a>
</p>

<div align="center">

# One Api Pro

_✨ 针对企业级的 AI API Gateway ✨_

> 本项目基于 [one-api](https://github.com/songquanpeng/one-api) (by [JustSong](https://github.com/songquanpeng)) 二次开发，感谢原作者的开源贡献。

</div>

<p align="center">
  <a href="https://raw.githubusercontent.com/songquanpeng/one-api/main/LICENSE">
    <img src="https://img.shields.io/github/license/songquanpeng/one-api?color=brightgreen" alt="license">
  </a>
  <a href="https://github.com/songquanpeng/one-api/releases/latest">
    <img src="https://img.shields.io/github/v/release/songquanpeng/one-api?color=brightgreen&include_prereleases" alt="release">
  </a>
  <a href="https://hub.docker.com/repository/docker/justsong/one-api">
    <img src="https://img.shields.io/docker/pulls/justsong/one-api?color=brightgreen" alt="docker pull">
  </a>
  <a href="https://github.com/songquanpeng/one-api/releases/latest">
    <img src="https://img.shields.io/github/downloads/songquanpeng/one-api/total?color=brightgreen&include_prereleases" alt="release">
  </a>
  <a href="https://goreportcard.com/report/github.com/songquanpeng/one-api">
    <img src="https://goreportcard.com/badge/github.com/songquanpeng/one-api" alt="GoReportCard">
  </a>
</p>

<p align="center">
  <a href="https://github.com/songquanpeng/one-api#部署">部署教程</a>
  ·
  <a href="https://github.com/songquanpeng/one-api#使用方法">使用方法</a>
  ·
  <a href="https://github.com/songquanpeng/one-api/issues">意见反馈</a>
  ·
  <a href="https://github.com/songquanpeng/one-api#截图展示">截图展示</a>
  ·
  <a href="https://openai.justsong.cn/">在线演示</a>
  ·
  <a href="https://github.com/songquanpeng/one-api#常见问题">常见问题</a>
  ·
  <a href="https://github.com/songquanpeng/one-api#相关项目">相关项目</a>
  ·
  <a href="https://iamazing.cn/page/reward">赞赏支持</a>

---

## 📑 目录

- [🔥 one-api-pro vs one-api](#-one-api-pro-vs-one-api)
- [📦 部署](#-部署)
- [⚙️ 配置](#%EF%B8%8F-配置)
- [📖 使用方法](#-使用方法)
- [🎬 演示](#-演示)
- [❓ 常见问题](#-常见问题)
- [🔧 第三方技术](#-第三方技术)

---

## 🔥 one-api-pro vs one-api

| 对比维度 | one-api | one-api-pro |
|----------|---------|-------------|
| 项目名称 | one-api | one-api-pro |
| Adaptor 架构 | 集中式常量管理（channeltype/define.go 56 行 iota + url.go 平行数组 + helper.go 双层 switch），新增提供商必须修改 4 个框架文件 | 自注册机制（registry + register.go），新增提供商只需创建包 + 注册即可，框架代码零修改 |
| 订阅模式 | 无套餐/订阅体系 | 完整套餐订阅 + 周期限频 + 按模型管控 |
| 去中心化集群 | 无独立集群支持，多机部署需共享 MySQL | 支持去中心化多活集群，每节点独立 MySQL + Redis，通过应用层事件同步实现数据互信，无需共享数据库 |
| 目录结构 | relay/adaptor/ 平铺 40 个目录，基础协议与供应商混在一起，relay/model/ 与根 model/ 冲突 | adaptor/openai/、adaptor/anthropic/ 作为基础协议独立放置，adaptor/provider/ 统一收纳 37 个供应商，relay/schema/ 消除命名冲突 |
| 管理后台 | 3 套前端主题（default/berry/air），基础管理功能 | Vue 3 + Arco Design 全新管理后台，可视化仪表盘 |
| 持续更新 | 原项目已于 2024 年停止更新 | 持续维护更新，针对企业级场景优化 |

## 📦 部署
### 🔨 手动部署
1. 从 [GitHub Releases](https://github.com/songquanpeng/one-api/releases/latest) 下载可执行文件或者从源码编译：
   ```shell
   git clone https://github.com/Leon-PanPan/one-api-pro.git

   # 构建前端
   cd one-api-pro/web/default
   npm install
   npm run build

   # 构建后端
   cd ../..
   go mod download
   go build -ldflags "-s -w" -o one-api-pro
   ````
2. 运行：
   ```shell
   chmod u+x one-api-pro
   ./one-api-pro --port 3000 --log-dir ./logs
   ```
3. 访问 [http://localhost:3000/](http://localhost:3000/) 并登录。初始账号用户名为 `root`，密码为 `123456`。

更加详细的部署教程[参见此处](https://iamazing.cn/page/how-to-deploy-a-website)。

### 🏢 多机部署
1. 所有服务器 `SESSION_SECRET` 设置一样的值。
2. 必须设置 `SQL_DSN`，使用 MySQL 数据库而非 SQLite，所有服务器连接同一个数据库。
3. 所有从服务器必须设置 `NODE_TYPE` 为 `slave`，不设置则默认为主服务器。
4. 设置 `SYNC_FREQUENCY` 后服务器将定期从数据库同步配置，在使用远程数据库的情况下，推荐设置该项并启用 Redis，无论主从。
5. 从服务器可以选择设置 `FRONTEND_BASE_URL`，以重定向页面请求到主服务器。
6. 从服务器上**分别**装好 Redis，设置好 `REDIS_CONN_STRING`，这样可以做到在缓存未过期的情况下数据库零访问，可以减少延迟（Redis 集群或者哨兵模式的支持请参考环境变量说明）。
7. 如果主服务器访问数据库延迟也比较高，则也需要启用 Redis，并设置 `SYNC_FREQUENCY`，以定期从数据库同步配置。

环境变量的具体使用方法详见[此处](#环境变量)。

### 🌐 集群部署（去中心化多活）

集群模式允许多个节点各自部署独立的 One Api Pro + MySQL，通过应用层事件同步实现数据互信，无需共享数据库。

> **适用场景**：全球多地域部署、就近访问降低延迟、高可用容灾、多节点负载均衡。

#### 🗺️ 架构概览

```
                    ┌─────────────┐
                    │  Nginx/LB   │  （统一入口，ip_hash 负载均衡）
                    └──────┬──────┘
                           │
            ┌──────────────┼──────────────┐
            │              │              │
     ┌──────┴──────┐ ┌────┴───────┐ ┌───┴────────┐
     │  Node A     │ │  Node B     │ │  Node C     │
     │ (one-api-pro)   │ │ (one-api-pro)   │ │ (one-api-pro)   │
     │ + MySQL     │ │ + MySQL     │ │ + MySQL     │
     │ + Redis     │ │ + Redis     │ │ + Redis     │
     └──────┬──────┘ └─────┬──────┘ └────┬────────┘
            │              │              │
            └────── HTTP 推送同步事件 ──────┘
```

#### ⭐ 核心特性

- **去中心化**：所有节点地位平等，无主从之分，任何节点数据变更后主动推送至所有存活节点
- **零侵入**：通过 GORM 回调捕获数据变更，不修改现有业务代码
- **异步推送**：数据同步不阻塞主流程，通过后台 goroutine 批量推送
- **冲突解决**：基于 `updated_at` 时间戳比较，只有更新的数据才写入
- **限流同步**：渠道并发和 RPM 限流计数器通过数据库表实现跨节点同步
- **单节点兼容**：不配置集群环境变量时，系统完全以单节点模式运行

#### 📊 同步范围

| 数据表 | 是否同步 | 说明 |
|--------|---------|------|
| users | ✅ | 用户信息 |
| tokens | ✅ | API 令牌 |
| channels | ✅ | 渠道配置 |
| abilities | ✅ | 渠道能力 |
| options | ✅ | 系统设置 |
| redemptions | ✅ | 兑换码 |
| plans | ✅ | 订阅计划 |
| user_plans | ✅ | 用户订阅 |
| plan_usages | ✅ | 计划用量 |
| channel_counters | ✅ | 渠道限流计数器 |
| cluster_nodes | 🔄 Discovery | 集群节点信息（由发现机制维护，不走数据同步） |
| logs | ⚠️ 可选 | 日志数据量较大，通过 `CLUSTER_SYNC_LOGS` 控制 |

#### 🚀 部署步骤

**1. MySQL 配置（每个节点必须使用独立的 MySQL 实例）**

每个节点都需要一个**独立的 MySQL 实例**（不能在同一 MySQL 实例中创建多个数据库来部署多个节点，因为 `auto_increment_offset` 是实例级变量）。

```ini
# 节点 1 的 my.cnf
[mysqld]
server-id = 1
auto_increment_increment = 50
auto_increment_offset = 1
log_bin = mysql-bin
binlog_format = ROW

# 节点 2 的 my.cnf
[mysqld]
server-id = 2
auto_increment_increment = 50
auto_increment_offset = 2
log_bin = mysql-bin
binlog_format = ROW

# 节点 3 的 my.cnf
[mysqld]
server-id = 3
auto_increment_increment = 50
auto_increment_offset = 3
log_bin = mysql-bin
binlog_format = ROW
```

> `auto_increment_increment` 设为 50，最多支持 50 个节点。每个节点的 `offset` 必须与 `CLUSTER_NODE_ID` 一致且互不相同。

> **重要说明：** `auto_increment_increment` 和 `auto_increment_offset` 是 MySQL 的**系统级变量**，对实例内所有数据库生效，无法为不同数据库设置不同的值，也无法在表级别设置（MySQL 表选项仅支持 `AUTO_INCREMENT` 起始值，不支持步长）。因此每个节点**必须使用独立的 MySQL 实例**，不能在同一个 MySQL 实例中通过创建不同数据库来部署多个节点。如需在同一台机器上运行多个 MySQL 实例，可以使用不同端口启动多个 mysqld 进程，或使用 Docker 运行多个独立的 MySQL 容器。

> **关于 `server-id` 和 binlog：** `server-id` 在同一集群的所有 MySQL 实例中必须互不相同。`log_bin` 和 `binlog_format=ROW` 强烈建议启用——它们用于未来的主从复制扩展和 point-in-time recovery。集群数据同步本身不依赖 binlog（通过 GORM 回调在应用层实现），但 binlog 提供了额外的可靠性保障。

**2. Redis 配置（每个节点必须使用独立的 Redis 实例）**

每个节点也需要**独立的 Redis 实例**（端口不同或在不同机器上）。Redis 在本集群架构中不用于节点间通信，只用于本节点的缓存、限流等业务用途。

**3. 新节点初始化数据**

新节点上线时，需要先获取已有节点的数据快照：

```bash
# 方式一：从已有节点导出并导入
mysqldump -h existing-node -u root -p oneapi > backup.sql
mysql -u root -p oneapi < backup.sql

# 方式二：通过 API 获取快照（需先启动服务）
curl -H "X-Cluster-Secret: your-secret" \
  "https://existing-node/api/cluster/snapshot?tables=users,tokens,channels,abilities,options,redemptions,plans,user_plans,plan_usages" \
  -o snapshot.json
```

**4. 环境变量配置（完整案例）**

以下是 3 节点集群的完整 `.env` 配置示例。每个节点都使用独立的 MySQL 和 Redis 实例，端口和路径各不相同。

**节点 1 — 中国节点（`/opt/one-api-pro/node1/.env`）：**
```bash
# ========================
# 基础配置
# ========================
PORT=3000
SYSTEM_NAME=One Api Pro Cluster

# ========================
# 数据库（独立 MySQL 实例）
# ========================
SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node1?charset=utf8mb4&parseTime=True&loc=Local

# ========================
# Redis（独立 Redis 实例）
# ========================
REDIS_CONN_STRING=redis://127.0.0.1:6379/0

# ========================
# 集群配置
# ========================
CLUSTER_ENABLED=true
CLUSTER_NODE_ID=1
CLUSTER_NODE_NAME=node-cn
CLUSTER_NODE_ADDRESS=https://cn.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me

# 种子节点（首次启动时引导发现其他节点）
# 第一个节点：填自己的地址或留空
# 后续节点：填任意一个已存活节点的地址
CLUSTER_SEEDS=https://cn.example.com,https://us.example.com,https://eu.example.com

# ========================
# 集群调优（可选）
# ========================
# 节点发现间隔（秒），存活节点每周期互相 ping，默认 30
CLUSTER_DISCOVERY_INTERVAL=30
# 失败节点 ping 间隔（秒），比存活间隔长以减少无效请求，默认 120
CLUSTER_DEAD_PING_INTERVAL=120
# 连续 ping 失败次数，达到后标记节点为失败状态（status=2），默认 3
CLUSTER_MAX_PING_FAILURES=3
# 同步事件推送间隔（秒），默认 3
CLUSTER_PUSH_INTERVAL=3
# 是否同步日志表（logs 数据量较大，按需关闭），默认 true
CLUSTER_SYNC_LOGS=true
# 每次推送最大事件数，默认 50
CLUSTER_BATCH_SIZE=50
```

**节点 2 — 美国节点（`/opt/one-api-pro/node2/.env`）：**
```bash
# 基础配置
PORT=3001
SYSTEM_NAME=One Api Pro Cluster

# 数据库（独立 MySQL 实例，端口或机器与节点 1 不同）
SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node2?charset=utf8mb4&parseTime=True&loc=Local

# Redis（独立 Redis 实例）
REDIS_CONN_STRING=redis://127.0.0.1:6380/0

# 集群配置
CLUSTER_ENABLED=true
CLUSTER_NODE_ID=2
CLUSTER_NODE_NAME=node-us
CLUSTER_NODE_ADDRESS=https://us.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me   # 必须与节点 1 完全一致

# 填任意一个已存活节点的地址
CLUSTER_SEEDS=https://cn.example.com
```

**节点 3 — 欧洲节点（`/opt/one-api-pro/node3/.env`）：**
```bash
# 基础配置
PORT=3002
SYSTEM_NAME=One Api Pro Cluster

# 数据库
SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node3?charset=utf8mb4&parseTime=True&loc=Local

# Redis
REDIS_CONN_STRING=redis://127.0.0.1:6381/0

# 集群配置
CLUSTER_ENABLED=true
CLUSTER_NODE_ID=3
CLUSTER_NODE_NAME=node-eu
CLUSTER_NODE_ADDRESS=https://eu.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me   # 必须与所有节点一致

# 填任意一个已存活节点的地址
CLUSTER_SEEDS=https://cn.example.com
```

**配置参数对照表：**

| 环境变量 | 节点 1 | 节点 2 | 节点 3 | 说明 |
|---|---|---|---|---|
| `PORT` | 3000 | 3001 | 3002 | 监听端口（同一机器需要不同） |
| `SQL_DSN` | ...oneapi_node1 | ...oneapi_node2 | ...oneapi_node3 | 独立 MySQL 实例 |
| `REDIS_CONN_STRING` | :6379/0 | :6380/0 | :6381/0 | 独立 Redis 实例 |
| `CLUSTER_NODE_ID` | 1 | 2 | 3 | 节点编号，对应 MySQL `auto_increment_offset` |
| `CLUSTER_NODE_NAME` | node-cn | node-us | node-eu | 节点名称，便于识别 |
| `CLUSTER_NODE_ADDRESS` | https://cn.example.com | https://us.example.com | https://eu.example.com | 节点公网地址（其他节点通过此地址访问） |
| `CLUSTER_SECRET` | 同一个值 | 同一个值 | 同一个值 | **所有节点必须完全一致** |
| `CLUSTER_SEEDS` | 自己的地址或留空 | 任意存活节点 | 任意存活节点 | 首次启动引导，后续自动发现 |

**5. 启动命令**

每个节点使用 `--env` 参数加载自己的配置文件：

```bash
# 节点 1
./one-api-pro --env /opt/one-api-pro/node1/.env --port 3000

# 节点 2
./one-api-pro --env /opt/one-api-pro/node2/.env --port 3001

# 节点 3
./one-api-pro --env /opt/one-api-pro/node3/.env --port 3002
```

**6. 启动顺序**

1. 启动第一个节点（Node A），`CLUSTER_SEEDS` 留空或填自己的地址
2. 等待 Node A 完全启动（约 5-10 秒，看到"集群模块初始化完成"日志）
3. 启动后续节点，`CLUSTER_SEEDS` 填写任意一个已存活节点的地址
4. 后续节点启动后会自动 ping 种子节点，传递性发现所有其他节点
5. 所有节点启动后，可通过任一节点的管理后台"设置 → 节点管理"页面查看节点状态

**7. Nginx 负载均衡配置示例（可选）**

```nginx
upstream one_api_cluster {
    ip_hash;  # 基于 IP 哈希，同一用户请求固定到同一节点，保证 session/cache 命中
    server cn.example.com:3000;
    server us.example.com:3000;
    server eu.example.com:3000;
}

server {
    listen 443 ssl;
    server_name api.example.com;

    location / {
        proxy_pass http://one_api_cluster;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 300s;
    }
}
```

> **使用 `ip_hash` 是关键**：保证同一用户的请求始终到同一节点，避免 plan 限频、Redis 缓存等状态在不同节点间丢失。

**8. 验证集群状态**

部署完成后，可以通过以下方式验证：

```bash
# 查看节点列表（任一节点上调用）
curl -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  https://cn.example.com/api/cluster_node/

# 应返回所有节点的列表，包含 status、last_heartbeat、ping_failures 等字段
```

或在管理后台：**设置 → 节点管理** 页面查看节点列表、状态、最后心跳时间等。

#### 📡 集群管理 API

| 接口 | 方法 | 说明 | 认证 |
|------|------|------|------|
| `/api/cluster/ping` | POST | 节点发现与心跳（双向注册） | 集群密钥 |
| `/api/cluster/sync` | POST | 接收同步事件 | 集群密钥 |
| `/api/cluster_node/` | GET | 查看节点列表（含本机、已禁用节点） | Root 管理员 |
| `/api/cluster_node/:id` | GET | 查看单个节点 | Root 管理员 |
| `/api/cluster_node/` | POST | 添加节点（含 secret） | Root 管理员 |
| `/api/cluster_node/` | PUT | 更新节点（可更新 secret） | Root 管理员 |
| `/api/cluster_node/:id` | DELETE | 软删除（禁用）节点 | Root 管理员 |
| `/api/cluster_node/:id/enable` | POST | 重新启用已禁用的节点 | Root 管理员 |
| `/api/cluster_node/ping/:id` | GET | 手动 ping 节点 | Root 管理员 |

#### ⚠️ 注意事项

- 每个节点必须有独立的 MySQL 实例和 Redis 实例，不共享数据库
- `CLUSTER_SECRET` 在所有节点间必须一致，请使用强密码并妥善保管
- `CLUSTER_NODE_ID` 在所有节点间必须互不相同，且与 MySQL `auto_increment_offset` 一致
- `CLUSTER_NODE_ADDRESS` 必须是其他节点可访问的公网地址（包含协议前缀如 `https://`）
- 新节点上线前的数据初始化需要手动完成（从在线节点拉取快照）
- 日志表（logs）数据量较大，可通过 `CLUSTER_SYNC_LOGS=false` 关闭日志同步
- MySQL 的 `auto_increment_increment` 和 `auto_increment_offset` 必须与 `CLUSTER_NODE_ID` 配置一致
- 节点发现采用 ping 双向注册机制，失败节点不会被删除，只标记为 status=2，网络恢复后自动复活
- `CLUSTER_SEEDS` 只是首次启动的引导；节点一旦通过 ping 发现其他节点，就不再依赖 SEEDS
- 节点离线期间其他节点产生的变更**不会自动补传**，离线节点重新上线后需拉取快照补齐数据

#### 📝 关于"本机节点"自我注册

每个节点启动时会在自己的 `cluster_nodes` 表中写入一条本机记录（`node_id` 等于本机配置的 `CLUSTER_NODE_ID`）。这是**有意的设计**，原因如下：

1. **管理后台展示**：在"设置 → 节点管理"页面，管理员需要看到本机信息（地址、状态、心跳时间等），以便排查问题
2. **节点发现传递性**：当节点 B 收到节点 A 的 ping 请求时，A 在响应中返回完整的节点列表（包含 A 自身）。B 收到后将其合并到本地表中。这样 C 通过 B 的响应也能学习到 A 的存在
3. **存活判断依据**：本机记录的 `last_heartbeat` 由本机每 30 秒自动更新一次（`discoverOnce` 函数中），反映本机正常运行的状态

**自我注册不会导致循环同步数据**。系统在 5 个层面做了防护：

| 防护点 | 作用 |
|---|---|
| ① `GetAllRemoteNodes` SQL 过滤 | 发现时 SQL 加 `node_id != ?` 排除本机 |
| ② `GetAliveNodesForSync` SQL 过滤 | 推送时 SQL 加 `node_id != ?` 排除本机 |
| ③ `handlePing` 拒绝自 ping | 显式拒绝 `req.NodeId == NodeID` |
| ④ `mergeDiscoveredNodes` 跳过本机 | 合并发现节点时跳过本机 |
| ⑤ `ApplyEvents` 跳过本机事件 | 应用事件时跳过本机产生的事件 |

数据流是单向的：从本机推到远程，从远程拉过来应用到本机，**永远不会有回路**。

管理后台会在本机节点名称旁显示"本机"蓝色徽章，并对本机禁用"删除"和"手动 Ping"操作（这两个操作对本机无意义）。

#### 🔐 关于"每节点独立 secret"

每个节点有**自己的 secret**，不再使用全局共享 secret。设计原因：

1. **安全性**：一个节点泄露 secret 不会影响其他节点
2. **管理灵活**：每个节点可以独立轮换自己的 secret
3. **自动发现**：节点间 ping 时自动携带自己的 secret 供对方保存

**Secret 生命周期**：
- 节点首次启动：用 `CLUSTER_SECRET` 环境变量作为初始值，写入 `cluster_nodes.secret_key` 字段
- 后续启动：从 `cluster_nodes.secret_key` 读取
- Admin 可以在"节点管理"页面修改其他节点的 secret
- ping 时 `X-Cluster-Secret` 头部 = **目标节点**的 secret（从本地 DB 查）

**添加新节点流程**：
1. 在节点 A 上添加 B 节点记录，填入 B 的 `CLUSTER_SECRET` 值
2. 在节点 B 上添加 A 节点记录，填入 A 的 `CLUSTER_SECRET` 值
3. A ping B：用 B 的 secret；B 接收：验证 B 自己的 secret ✓
4. B 响应中携带 A、B 各自的 secret，A 更新本地保存

#### 🗑️ 关于"软删除节点"

Admin 删除节点时**不物理删除**记录，而是设置 `disabled = true`：

- 防止被删除节点"自动长回来"（ping 机制会重新注册）
- 已禁用的节点仍然会响应 ping（让对方知道本节点在线），但不会获取本节点信息
- 物理删除需要手动 SQL：`DELETE FROM cluster_nodes WHERE node_id = ?`

#### 🔄 关于"数据同步机制"（重要）

**集群数据同步**完全依赖 **GORM 事件 + HTTP 主动推送**机制：
- 任何业务表的 INSERT/UPDATE/DELETE 操作 → GORM 回调捕获 → 写入 `sync_events` 表 → Pusher goroutine 推送到所有存活节点
- 接收方用 `WithSkipHook` 写本地数据库（不会回环）
- 接收方跳过 `event.NodeId == 本机 NodeID` 的事件（双重保险）

**架构权衡**：本设计**不实现跨节点主动拉取**，原因如下：
1. **侵入业务**：跨节点拉取需要知道每张表的业务唯一字段，会侵入业务代码
2. **主键冲突**：跨节点自增 ID 不连续（不同 `auto_increment_offset`），使用源节点 id 会破坏 offset 设计
3. **复杂度高**：维护成本高，可靠性提升有限
4. **主动推送够用**：95% 的场景（节点在线时的常规同步）完全由推送覆盖

**已知限制与运维要求**：
- 节点离线期间其他节点产生的数据变更 → **永久丢失**（推送是实时的）
- 节点重新上线后无法自动补齐离线期间的数据
- 新节点加入后只能接收到加入之后的数据变更，无历史数据
- **运维补救**：使用 `mysqldump` 从其他节点导出后导入

**典型部署场景对照**：

| 场景 | 是否需要拉取 | 处理方式 |
|---|---|---|
| 节点永久在线 | ❌ | 推送完全够用 |
| 节点偶尔重启（分钟级） | ⚠️ | 短时离线数据丢失，运维可接受 |
| 节点频繁维护 | ❌ | 推送继续，重启后立即恢复 |
| 新节点加入集群 | ❌ | DBA 手动 `mysqldump` 初始化 |
| 节点长期离线后恢复 | ❌ | DBA 手动 `mysqldump` 补齐 |

如果部署后访问出现空白页面，详见 [#97](https://github.com/songquanpeng/one-api/issues/97)。

### 🔌 部署第三方服务配合 One Api Pro 使用
> 欢迎 PR 添加更多示例。

#### 💬 ChatGPT Next Web
项目主页：https://github.com/Yidadaa/ChatGPT-Next-Web

```bash
docker run --name chat-next-web -d -p 3001:3000 yidadaa/chatgpt-next-web
```

注意修改端口号，之后在页面上设置接口地址（例如：https://openai.justsong.cn/ ）和 API Key 即可。

#### 💬 ChatGPT Web
项目主页：https://github.com/Chanzhaoyu/chatgpt-web

```bash
docker run --name chatgpt-web -d -p 3002:3002 -e OPENAI_API_BASE_URL=https://openai.justsong.cn -e OPENAI_API_KEY=sk-xxx chenzhaoyu94/chatgpt-web
```

注意修改端口号、`OPENAI_API_BASE_URL` 和 `OPENAI_API_KEY`。

#### 🤖 QChatGPT - QQ 机器人
项目主页：https://github.com/RockChinQ/QChatGPT

根据[文档](https://qchatgpt.rockchin.top)完成部署后，在 `data/provider.json`设置`requester.openai-chat-completions.base-url`为 One Api Pro 实例地址，并填写 API Key 到 `keys.openai` 组中，设置 `model` 为要使用的模型名称。

运行期间可以通过`!model`命令查看、切换可用模型。

### ☁️ 部署到第三方平台
<details>
<summary><strong>部署到 Sealos </strong></summary>
<div>

> Sealos 的服务器在国外，不需要额外处理网络问题，支持高并发 & 动态伸缩。

点击以下按钮一键部署（部署后访问出现 404 请等待 3~5 分钟）：

[![Deploy-on-Sealos.svg](https://raw.githubusercontent.com/labring-actions/templates/main/Deploy-on-Sealos.svg)](https://cloud.sealos.io/?openapp=system-fastdeploy?templateName=one-api)

</div>
</details>

<details>
<summary><strong>部署到 Zeabur</strong></summary>
<div>

> Zeabur 的服务器在国外，自动解决了网络的问题，同时免费的额度也足够个人使用

[![Deploy on Zeabur](https://zeabur.com/button.svg)](https://zeabur.com/templates/7Q0KO3)

1. 首先 fork 一份代码。
2. 进入 [Zeabur](https://zeabur.com?referralCode=songquanpeng)，登录，进入控制台。
3. 新建一个 Project，在 Service -> Add Service 选择 Marketplace，选择 MySQL，并记下连接参数（用户名、密码、地址、端口）。
4. 复制链接参数，运行 ```create database `one-api-pro` ``` 创建数据库。
5. 然后在 Service -> Add Service，选择 Git（第一次使用需要先授权），选择你 fork 的仓库。
6. Deploy 会自动开始，先取消。进入下方 Variable，添加一个 `PORT`，值为 `3000`，再添加一个 `SQL_DSN`，值为 `<username>:<password>@tcp(<addr>:<port>)/one-api` ，然后保存。 注意如果不填写 `SQL_DSN`，数据将无法持久化，重新部署后数据会丢失。
7. 选择 Redeploy。
8. 进入下方 Domains，选择一个合适的域名前缀，如 "my-one-api"，最终域名为 "my-one-api.zeabur.app"，也可以 CNAME 自己的域名。
9. 等待部署完成，点击生成的域名进入 One Api Pro。

</div>
</details>

<details>
<summary><strong>部署到 Render</strong></summary>
<div>

> Render 提供免费额度，绑卡后可以进一步提升额度

Render 可以直接部署 docker 镜像，不需要 fork 仓库：https://dashboard.render.com

</div>
</details>

## ⚙️ 配置
系统本身开箱即用。

你可以通过设置环境变量或者命令行参数进行配置。

等到系统启动后，使用 `root` 用户登录系统并做进一步的配置。

**Note**：如果你不知道某个配置项的含义，可以临时删掉值以看到进一步的提示文字。

## 📖 使用方法
在`渠道`页面中添加你的 API Key，之后在`令牌`页面中新增访问令牌。

之后就可以使用你的令牌访问 One Api Pro 了，使用方式与 [OpenAI API](https://platform.openai.com/docs/api-reference/introduction) 一致。

你需要在各种用到 OpenAI API 的地方设置 API Base 为你的 One Api Pro 的部署地址，例如：`https://openai.justsong.cn`，API Key 则为你在 One API 中生成的令牌。

注意，具体的 API Base 的格式取决于你所使用的客户端。

例如对于 OpenAI 的官方库：
```bash
OPENAI_API_KEY="sk-xxxxxx"
OPENAI_API_BASE="https://<HOST>:<PORT>/v1"
```

```mermaid
graph LR
    A(用户)
    A --->|使用 One API 分发的 key 进行请求| B(One API)
    B -->|中继请求| C(OpenAI)
    B -->|中继请求| D(Azure)
    B -->|中继请求| E(其他 OpenAI API 格式下游渠道)
    B -->|中继并修改请求体和返回体| F(非 OpenAI API 格式下游渠道)
```

可以通过在令牌后面添加渠道 ID 的方式指定使用哪一个渠道处理本次请求，例如：`Authorization: Bearer ONE_API_PRO_KEY-CHANNEL_ID`。
注意，需要是管理员用户创建的令牌才能指定渠道 ID。

不加的话将会使用负载均衡的方式使用多个渠道。

### 🔧 环境变量
> One API 支持从 `.env` 文件中读取环境变量，请参照 `.env.example` 文件，使用时请将其重命名为 `.env`。也可通过 `--env` 参数指定配置文件路径（支持相对路径），详见命令行参数一节。
1. `REDIS_CONN_STRING`：设置之后将使用 Redis 作为缓存使用。
   + 例子：`REDIS_CONN_STRING=redis://default:redispw@localhost:49153`
   + 如果数据库访问延迟很低，没有必要启用 Redis，启用后反而会出现数据滞后的问题。
   + 如果需要使用哨兵或者集群模式：
     + 则需要把该环境变量设置为节点列表，例如：`localhost:49153,localhost:49154,localhost:49155`。
     + 除此之外还需要设置以下环境变量：
       + `REDIS_PASSWORD`：Redis 集群或者哨兵模式下的密码设置。
       + `REDIS_MASTER_NAME`：Redis 哨兵模式下主节点的名称。
2. `SESSION_SECRET`：设置之后将使用固定的会话密钥，这样系统重新启动后已登录用户的 cookie 将依旧有效。
   + 例子：`SESSION_SECRET=random_string`
3. `SQL_DSN`：设置之后将使用指定数据库而非 SQLite，请使用 MySQL 或 PostgreSQL。
   + 例子：
     + MySQL：`SQL_DSN=root:123456@tcp(localhost:3306)/oneapi`
     + PostgreSQL：`SQL_DSN=postgres://postgres:123456@localhost:5432/oneapi`（适配中，欢迎反馈）
   + 注意需要提前建立数据库 `oneapi`，无需手动建表，程序将自动建表。
   + 如果使用本地数据库：部署命令可添加 `--network="host"` 以使得容器内的程序可以访问到宿主机上的 MySQL。
   + 如果使用云数据库：如果云服务器需要验证身份，需要在连接参数中添加 `?tls=skip-verify`。
   + 请根据你的数据库配置修改下列参数（或者保持默认值）：
     + `SQL_MAX_IDLE_CONNS`：最大空闲连接数，默认为 `100`。
     + `SQL_MAX_OPEN_CONNS`：最大打开连接数，默认为 `1000`。
       + 如果报错 `Error 1040: Too many connections`，请适当减小该值。
     + `SQL_CONN_MAX_LIFETIME`：连接的最大生命周期，默认为 `60`，单位分钟。
4. `LOG_SQL_DSN`：设置之后将为 `logs` 表使用独立的数据库，请使用 MySQL 或 PostgreSQL。
5. `FRONTEND_BASE_URL`：设置之后将重定向页面请求到指定的地址，仅限从服务器设置。
   + 例子：`FRONTEND_BASE_URL=https://openai.justsong.cn`
6. `MEMORY_CACHE_ENABLED`：启用内存缓存，会导致用户额度的更新存在一定的延迟，可选值为 `true` 和 `false`，未设置则默认为 `false`。
   + 例子：`MEMORY_CACHE_ENABLED=true`
7. `SYNC_FREQUENCY`：在启用缓存的情况下与数据库同步配置的频率，单位为秒，默认为 `600` 秒。
   + 例子：`SYNC_FREQUENCY=60`
8. `NODE_TYPE`：设置之后将指定节点类型，可选值为 `master` 和 `slave`，未设置则默认为 `master`。
   + 例子：`NODE_TYPE=slave`
9. `CHANNEL_UPDATE_FREQUENCY`：设置之后将定期更新渠道余额，单位为分钟，未设置则不进行更新。
   + 例子：`CHANNEL_UPDATE_FREQUENCY=1440`
10. `CHANNEL_TEST_FREQUENCY`：设置之后将定期检查渠道，单位为分钟，未设置则不进行检查。 
   +例子：`CHANNEL_TEST_FREQUENCY=1440`
11. `POLLING_INTERVAL`：批量更新渠道余额以及测试可用性时的请求间隔，单位为秒，默认无间隔。
    + 例子：`POLLING_INTERVAL=5`
12. `BATCH_UPDATE_ENABLED`：启用数据库批量更新聚合，会导致用户额度的更新存在一定的延迟可选值为 `true` 和 `false`，未设置则默认为 `false`。
    + 例子：`BATCH_UPDATE_ENABLED=true`
    + 如果你遇到了数据库连接数过多的问题，可以尝试启用该选项。
13. `BATCH_UPDATE_INTERVAL=5`：批量更新聚合的时间间隔，单位为秒，默认为 `5`。
    + 例子：`BATCH_UPDATE_INTERVAL=5`
14. 请求频率限制：
    + `GLOBAL_API_RATE_LIMIT`：全局 API 速率限制（除中继请求外），单 ip 三分钟内的最大请求数，默认为 `180`。
    + `GLOBAL_WEB_RATE_LIMIT`：全局 Web 速率限制，单 ip 三分钟内的最大请求数，默认为 `60`。
15. 编码器缓存设置：
    + `TIKTOKEN_CACHE_DIR`：默认程序启动时会联网下载一些通用的词元的编码，如：`gpt-3.5-turbo`，在一些网络环境不稳定，或者离线情况，可能会导致启动有问题，可以配置此目录缓存数据，可迁移到离线环境。
    + `DATA_GYM_CACHE_DIR`：目前该配置作用与 `TIKTOKEN_CACHE_DIR` 一致，但是优先级没有它高。
16. `RELAY_TIMEOUT`：中继超时设置，单位为秒，默认不设置超时时间。
17. `RELAY_PROXY`：设置后使用该代理来请求 API。
18. `USER_CONTENT_REQUEST_TIMEOUT`：用户上传内容下载超时时间，单位为秒。
19. `USER_CONTENT_REQUEST_PROXY`：设置后使用该代理来请求用户上传的内容，例如图片。
20. `SQLITE_BUSY_TIMEOUT`：SQLite 锁等待超时设置，单位为毫秒，默认 `3000`。
21. `GEMINI_SAFETY_SETTING`：Gemini 的安全设置，默认 `BLOCK_NONE`。
22. `GEMINI_VERSION`：One Api Pro 所使用的 Gemini 版本，默认为 `v1`。
23. `THEME`：系统的主题设置，默认为 `default`，具体可选值参考[此处](./web/README.md)。
24. `ENABLE_METRIC`：是否根据请求成功率禁用渠道，默认不开启，可选值为 `true` 和 `false`。
25. `METRIC_QUEUE_SIZE`：请求成功率统计队列大小，默认为 `10`。
26. `METRIC_SUCCESS_RATE_THRESHOLD`：请求成功率阈值，默认为 `0.8`。
27. `INITIAL_ROOT_TOKEN`：如果设置了该值，则在系统首次启动时会自动创建一个值为该环境变量值的 root 用户令牌。
28. `INITIAL_ROOT_ACCESS_TOKEN`：如果设置了该值，则在系统首次启动时会自动创建一个值为该环境变量的 root 用户创建系统管理令牌。
29. `ENFORCE_INCLUDE_USAGE`：是否强制在 stream 模型下返回 usage，默认不开启，可选值为 `true` 和 `false`。
30. `TEST_PROMPT`：测试模型时的用户 prompt，默认为 `Print your model name exactly and do not output without any other text.`。

#### 🌐 集群配置（去中心化多活部署）
> 不配置以下环境变量时，系统以单节点模式运行，无任何副作用。

1. `CLUSTER_ENABLED`：是否启用集群模式，默认不启用。
   + 例子：`CLUSTER_ENABLED=true`
2. `CLUSTER_NODE_ID`：节点编号（1-49），必须与 MySQL 的 `auto_increment_offset` 一致，不同节点不能重复。
   + 例子：`CLUSTER_NODE_ID=1`
3. `CLUSTER_NODE_NAME`：节点名称，便于识别，默认为 `node-{NODE_ID}`。
   + 例子：`CLUSTER_NODE_NAME=node-cn`
4. `CLUSTER_NODE_ADDRESS`：本节点的公网访问地址（需包含协议前缀），其他节点通过此地址推送数据。
   + 例子：`CLUSTER_NODE_ADDRESS=https://cn.example.com`
5. `CLUSTER_SECRET`：本节点的初始 secret，**每个节点独立**。首次启动时作为初始 secret 写入数据库，之后可由 admin 修改。
   + 例子：`CLUSTER_SECRET=MyClusterSecret123abc`
6. `CLUSTER_SEEDS`：种子节点地址（逗号分隔），新节点启动时向种子节点注册获取集群信息，只需配置一个可达节点即可。第一个节点可以不配置或配置自己的地址。
   + 例子：`CLUSTER_SEEDS=https://cn.example.com`
   + 多个种子：`CLUSTER_SEEDS=https://cn.example.com,https://us.example.com`
7. `CLUSTER_PUSH_INTERVAL`：同步事件推送间隔，单位为秒，默认 `3`。
8. `CLUSTER_DISCOVERY_INTERVAL`：节点发现间隔，单位为秒，存活节点每周期互相 ping，默认 `30`。
9. `CLUSTER_DEAD_PING_INTERVAL`：失败节点 ping 间隔，单位为秒，比存活间隔长以减少无效请求，默认 `120`。
10. `CLUSTER_MAX_PING_FAILURES`：连续 ping 失败次数，达到后标记节点为失败状态，默认 `3`。
11. `CLUSTER_SYNC_LOGS`：是否同步日志表，日志数据量较大可按需关闭，默认 `true`。
     + 例子：`CLUSTER_SYNC_LOGS=false`
12. `CLUSTER_BATCH_SIZE`：每次推送最大事件数，默认 `50`。

### ⌨️ 命令行参数
1. `--port <port_number>`: 指定服务器监听的端口号，默认为 `3000`。
   + 例子：`--port 3000`
2. `--log-dir <log_dir>`: 指定日志文件夹，如果没有设置，默认保存至工作目录的 `logs` 文件夹下。
   + 例子：`--log-dir ./logs`
3. `--env <env_file_path>`: 指定配置文件路径，支持相对路径和绝对路径。未指定时自动加载当前目录的 `.env` 文件。
   + 例子：`--env ./config.env`
   + 例子：`--env /etc/one-api-pro/production.env`
   + 多实例部署示例：
     ```bash
     ./one-api-pro --env ./instances/instance1.env --port 3001 &
     ./one-api-pro --env ./instances/instance2.env --port 3002 &
     ```
   + 配置优先级：命令行参数 > 系统环境变量 > `--env` 指定的配置文件 > 默认值
4. `--version`: 打印系统版本号并退出。
5. `--help`: 查看命令的使用帮助和参数说明。

## 🎬 演示
### 🔗 在线演示
注意，该演示站不提供对外服务：
https://openai.justsong.cn

### 📸 截图展示
![channel](https://user-images.githubusercontent.com/39998050/233837954-ae6683aa-5c4f-429f-a949-6645a83c9490.png)
![token](https://user-images.githubusercontent.com/39998050/233837971-dab488b7-6d96-43af-b640-a168e8d1c9bf.png)

## ❓ 常见问题
1. 额度是什么？怎么计算的？One Api Pro 的额度计算有问题？
   + 额度 = 分组倍率 * 模型倍率 * （提示 token 数 + 补全 token 数 * 补全倍率）
   + 其中补全倍率对于 GPT3.5 固定为 1.33，GPT4 为 2，与官方保持一致。
   + 如果是非流模式，官方接口会返回消耗的总 token，但是你要注意提示和补全的消耗倍率不一样。
   + 注意，One Api Pro 的默认倍率就是官方倍率，是已经调整过的。
2. 账户额度足够为什么提示额度不足？
   + 请检查你的令牌额度是否足够，这个和账户额度是分开的。
   + 令牌额度仅供用户设置最大使用量，用户可自由设置。
3. 提示无可用渠道？
   + 请检查的用户分组和渠道分组设置。
   + 以及渠道的模型设置。
4. 渠道测试报错：`invalid character '<' looking for beginning of value`
   + 这是因为返回值不是合法的 JSON，而是一个 HTML 页面。
   + 大概率是你的部署站的 IP 或代理的节点被 CloudFlare 封禁了。
5. ChatGPT Next Web 报错：`Failed to fetch`
   + 部署的时候不要设置 `BASE_URL`。
   + 检查你的接口地址和 API Key 有没有填对。
   + 检查是否启用了 HTTPS，浏览器会拦截 HTTPS 域名下的 HTTP 请求。
6. 报错：`当前分组负载已饱和，请稍后再试`
   + 上游渠道 429 了。
7. 升级之后我的数据会丢失吗？
   + 如果使用 MySQL，不会。
   + 如果使用 SQLite，需要按照我所给的部署命令挂载 volume 持久化 one-api-pro.db 数据库文件，否则容器重启后数据会丢失。
8. 升级之前数据库需要做变更吗？
   + 一般情况下不需要，系统将在初始化的时候自动调整。
   + 如果需要的话，我会在更新日志中说明，并给出脚本。
9. 手动修改数据库后报错：`数据库一致性已被破坏，请联系管理员`？
   + 这是检测到 ability 表里有些记录的渠道 id 是不存在的，这大概率是因为你删了 channel 表里的记录但是没有同步在 ability 表里清理无效的渠道。
   + 对于每一个渠道，其所支持的模型都需要有一个专门的 ability 表的记录，表示该渠道支持该模型。

## 🔧 第三方技术

本项目在后端和前端都使用了大量的开源第三方库，下面列出主要依赖及其用途。

### 🔙 后端（Go）

#### 🌐 Web 框架与中间件

| 技术 | 用途 |
| --- | --- |
| [Gin](https://github.com/gin-gonic/gin) | HTTP Web 框架，负责路由注册、请求处理、中间件链和响应渲染，是整个后端服务的基础 |
| [gin-contrib/cors](https://github.com/gin-contrib/cors) | 跨域资源共享（CORS）中间件，允许前端从不同域访问 API |
| [gin-contrib/gzip](https://github.com/gin-contrib/gzip) | Gzip 压缩中间件，对 HTTP 响应进行压缩以减少传输体积 |
| [gin-contrib/sessions](https://github.com/gin-contrib/sessions) | 基于 Cookie 的 Session 管理，保存用户登录状态 |
| [gin-contrib/static](https://github.com/gin-contrib/static) | 静态文件服务中间件，用于托管前端构建产物 |
| [go-playground/validator](https://github.com/go-playground/validator) | Struct tag 驱动的请求参数校验 |
| [gorilla/websocket](https://github.com/gorilla/websocket) | WebSocket 协议实现，用于支持流式响应（打字机效果）等长连接场景 |

#### 🗄️ 数据库与 ORM

| 技术 | 用途 |
| --- | --- |
| [GORM](https://gorm.io) | Go 主流 ORM 库，负责业务模型定义、数据库迁移、查询和写入 |
| [gorm.io/driver/sqlite](https://gorm.io/driver/sqlite) | GORM 的 SQLite 驱动，作为默认数据库，方便单机零配置部署 |
| [gorm.io/driver/mysql](https://gorm.io/driver/mysql) | GORM 的 MySQL 驱动，支持高并发场景下的 MySQL 部署 |
| [gorm.io/driver/postgres](https://gorm.io/driver/postgres) | GORM 的 PostgreSQL 驱动，提供 PostgreSQL 数据库支持 |
| [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) | MySQL 协议驱动，被 GORM MySQL driver 间接依赖 |
| [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) | SQLite 的 CGo 实现，被 GORM SQLite driver 间接依赖 |
| [jackc/pgx](https://github.com/jackc/pgx) | PostgreSQL 驱动，被 GORM PostgreSQL driver 间接依赖 |

#### ⚡ 缓存、限流与同步原语

| 技术 | 用途 |
| --- | --- |
| [go-redis/redis](https://github.com/go-redis/redis) | Redis 客户端，用于缓存、分布式限流计数器、集群模式下的节点限流共享等 |
| [patrickmn/go-cache](https://github.com/patrickmn/go-cache) | 进程内的内存缓存（`MEMORY_CACHE_ENABLED`），用于减轻数据库压力 |
| [golang/sync](https://pkg.go.dev/golang.org/x/sync) （`errgroup` 等） | 提供同步原语与错误聚合，简化并发任务编排 |
| [cloud.google.com/go/iam](https://cloud.google.com/go/iam) | 间接用于 Google API 客户端中的 IAM 凭证处理（Gemini 等渠道） |

#### 🔐 认证、安全与配置

| 技术 | 用途 |
| --- | --- |
| [golang-jwt/jwt](https://github.com/golang-jwt/jwt) | JWT（JSON Web Token）的生成与解析，用于系统访问令牌等场景 |
| [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) | Go 官方加密扩展库，提供 bcrypt 等密码哈希算法用于用户密码加密 |
| [joho/godotenv](https://github.com/joho/godotenv) | 解析 `.env` 配置文件，将其中的键值对注入为环境变量 |

#### 🤖 大模型 / AI 集成

| 技术 | 用途 |
| --- | --- |
| [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2) 及其 `credentials`、`bedrockruntime` 子模块 | 访问 AWS Bedrock 上的 Claude 系列模型 |
| [Google API Go Client](https://github.com/googleapis/google-api-go-client) （`google.golang.org/api`） | 调用 Google Gemini / PaLM2 系列模型 |
| [pkoukk/tiktoken-go](https://github.com/pkoukk/tiktoken-go) | OpenAI `tiktoken` 的 Go 移植版，用于对 prompt/completion 进行精确的 token 计数与计费 |
| [jinzhu/copier](https://github.com/jinzhu/copier) | 结构体之间按字段名快速复制，用于在不同请求/响应模型之间做转换 |

#### 🛠️ 通用工具

| 技术 | 用途 |
| --- | --- |
| [google/uuid](https://github.com/google/uuid) | 生成 UUID，用于兑换码、请求标识等 |
| [pkg/errors](https://github.com/pkg/errors) | 带堆栈信息的错误包装与处理，便于排障 |
| [golang.org/x/image](https://pkg.go.dev/golang.org/x/image) | Go 官方图像处理扩展库，用于图片相关辅助处理 |
| [smartystreets/goconvey](https://github.com/smartystreets/goconvey) | 测试与断言辅助库，配套 `stretchr/testify` 在测试中使用 |
| [stretchr/testify](https://github.com/stretchr/testify) | 常用的断言、mock、测试套件库，用于编写单元测试 |

### 🎨 前端（React，默认主题 `web/default`）

#### 🏗️ 核心框架与构建

| 技术 | 用途 |
| --- | --- |
| [React](https://react.dev) 18 | 负责整个管理后台 UI 的组件化渲染 |
| [React DOM](https://react.dev) | 在浏览器中挂载 React 组件树 |
| [react-router-dom](https://reactrouter.com) | 客户端路由，支撑多页面（渠道、令牌、用户管理等）切换 |
| [history](https://github.com/remix-run/history) | 被 react-router 用于管理浏览器历史栈 |
| [react-scripts](https://create-react-app.dev) | Create React App 的脚本集合，提供开发服务器、打包、测试等构建流程 |
| [Prettier](https://prettier.io) | 代码格式化工具，统一前端代码风格 |

#### UI 与交互

| 技术 | 用途 |
| --- | --- |
| [Semantic UI React](https://react.semantic-ui.com) 及其 [Semantic UI CSS](https://semantic-ui.com) | 管理后台主要的 UI 组件库（按钮、表单、表格、菜单等） |
| [react-toastify](https://github.com/fkhadra/react-toastify) | 全局通知/Toast 组件，用于操作结果提示 |
| [react-dropzone](https://react-dropzone.js.org) | 文件拖拽上传组件，用于头像等资源上传 |
| [react-turnstile](https://github.com/Le0developer/react-turnstile) | 集成 Cloudflare Turnstile 人机验证 |

#### 数据可视化与工具

| 技术 | 用途 |
| --- | --- |
| [Recharts](https://recharts.org) | 基于 React 的图表库，用于仪表盘等数据可视化 |
| [axios](https://axios-http.com) | HTTP 客户端，统一封装与管理后端 API 的交互 |
| [moment](https://momentjs.com) | 日期与时间的格式化、解析、相对时间展示 |
| [marked](https://github.com/markedjs/marked) | Markdown 渲染器，用于渲染公告、关于页等 Markdown 内容 |
| [i18next](https://www.i18next.com) 及其 `browser-languagedetector`、`react-i18next` | 国际化方案，支持多语言界面切换 |

### 部署与运维相关

| 技术 | 用途 |
| --- | --- |
| [Docker](https://www.docker.com) | 容器化打包与运行，便于跨平台一致部署 |
| [Nginx](https://nginx.org) | 文档中给出的反向代理与 HTTPS 终止示例 |
| [Let's Encrypt / certbot](https://certbot.eff.org) | 文档中推荐的免费 HTTPS 证书签发与自动续期方案 |
| [宝塔面板](https://www.bt.cn) | 文档中提到的可视化运维面板，可一键安装 One API |
| [Watchtower](https://github.com/containrrr/watchtower) | 自动监控并更新 Docker 镜像 |
| [MySQL](https://www.mysql.com) / [PostgreSQL](https://www.postgresql.org) | 生产环境推荐的关系型数据库 |
| [Redis](https://redis.io) | 可选的缓存与分布式限流后端 |
| [Cloudflare AI Gateway](https://developers.cloudflare.com/ai-gateway/) / [Cloudflare Turnstile](https://www.cloudflare.com/products/turnstile/) | 渠道代理与人机验证服务（功能层面集成） |
| [Feishu](https://open.feishu.cn) / [GitHub OAuth](https://github.com/settings/applications/new) | 第三方登录功能所对接的开放平台 |

## 适配器架构

One API 通过一套**适配器（Adaptor）接口**来对接不同的大模型厂商，所有渠道的 API 请求/响应处理都走统一的中继流程。

### 设计模式

```
                    Adaptor 接口（10 methods）
                    ┌───────┴──────────┐
                    │                  │
            GetRequestURL(),       DoResponse()
            SetupRequestHeader(),   │
            ConvertRequest(),       ├─ Passthrough → StreamPipeline.Run(openai配置)
            DoRequest()             ├─ Converting  → StreamPipeline.Run(gemini配置)
                                    ├─ EventDriven → AnthropicPipeline.Run(...)
                                    └─ 自定义      → xunfei/aws 独立实现
```

| 设计模式 | 用在哪 | 说明 |
|---------|-------|------|
| **Template Method** | `openai.Adaptor` 作为基类，子类覆写差异方法 | 20+ 个 OpenAI 兼容渠道继承 `openai.Adaptor`，仅覆写 `GetRequestURL` / `SetupRequestHeader` / `GetModelList` 等差异步骤 |
| **Strategy** | `StreamPipeline` 和 `AnthropicPipeline` | 将流式响应处理拆成三种策略：透传（Passthrough）、逐行转换（Converting）、事件驱动（EventDriven），每个 adaptor 选择其一 |
| **Hook（内嵌 Strategy）** | `StreamPipeline` 的 `ProcessLine` / `Render` / `MergeUsage` | 对逐行转换型渠道（gemini / baidu / cohere 等），通过注入拆帧函数、转换函数、usage 提取函数来参数化处理管线 |

### 目录结构

```
relay/adaptor/
├── interface.go              # Adaptor 接口契约
├── common.go                 # DoRequestHelper, SetupCommonRequestHeader
├── pipeline.go               # StreamPipeline — 标准 SSE/NDJSON 流处理
│
├── openai/                   # 基类：纯粹的 OpenAI 协议实现（无渠道特化分支）
│   └── adaptor.go            # 提供完整默认实现，子类通过 struct embedding 继承
│
├── groq/adaptor.go           # 子类：仅覆写模型列表（~10 行）
├── doubao/adaptor.go         # 子类：覆写 GetRequestURL + 模型列表（~25 行）
├── azure/adaptor.go          # 子类：覆写 GetRequestURL + SetupRequestHeader（~30 行）
│
├── anthropic/                # 独立渠道：AnthropicPipeline（事件驱动流处理）
├── gemini/                   # 独立渠道：StreamPipeline（逐行转换型）
├── xunfei/                   # 独立渠道：WebSocket 自定义传输
└── aws/                      # 独立渠道：AWS SDK 自定义传输
```

### 新增一个渠道需要写多少代码

**纯 OpenAI 兼容（如 Groq）—— ~10 行**：
```go
package groq
import oa "github.com/songquanpeng/one-api/relay/adaptor/openai"

type Adaptor struct { *oa.Adaptor }                       // 继承全部默认行为
func (a *Adaptor) GetModelList() []string { return ModelList }
func (a *Adaptor) GetChannelName() string { return "groq" }

var ModelList = []string{"llama-3.1-8b-instant", ...}
```

**URL 有差异（如 Doubao）—— ~25 行**，多覆写一个 `GetRequestURL` 方法。

**完全不同协议（如新厂商）—— ~80 行**，实现全部 Adaptor 方法 + 选择对应的 Pipeline 引擎。

### Dispatcher 派发机制

`relay/adaptor.go:GetAdaptorByChannel(channelType)` 接收渠道类型，返回对应的适配器：

```
channeltype.Doubao  → &doubao.Adaptor{Adaptor: &openai.Adaptor{}}  (子类)
channeltype.Anthropic  → &anthropic.Adaptor{}                      (独立渠道)
channeltype.Unknown    → &openai.Adaptor{}                         (默认兜底)
```

### Usage 提取

所有渠道的 `DoResponse` 最终返回 `*model.Usage`（共三个关键字段）：

| 字段 | 用途 |
|------|------|
| `PromptTokens` | 输入 token 总数 |
| `CompletionTokens` | 输出 token 总数 |
| `PromptTokensDetails.CachedTokens` | 缓存命中 token（按 `CachedPrice` 计价） |

计费公式：`非缓存Input × InputPrice + 缓存Input × CachedPrice + Output × OutputPrice`。各渠道在各自的 Pipeline 配置或 Hook 中完成原始响应 → Usage 的映射，对上游计费管线完全透明。

## 相关项目
* [FastGPT](https://github.com/labring/FastGPT): 基于 LLM 大语言模型的知识库问答系统
* [ChatGPT Next Web](https://github.com/Yidadaa/ChatGPT-Next-Web):  一键拥有你自己的跨平台 ChatGPT 应用
* [VChart](https://github.com/VisActor/VChart):  不只是开箱即用的多端图表库，更是生动灵活的数据故事讲述者。
* [VMind](https://github.com/VisActor/VMind):  不仅自动，还很智能。开源智能可视化解决方案。
* [CherryStudio](https://github.com/CherryHQ/cherry-studio):  全平台支持的AI客户端, 多服务商集成管理、本地知识库支持。

## 注意

本项目使用 MIT 协议进行开源，**在此基础上**，必须在页面底部保留署名以及指向本项目的链接。如果不想保留署名，必须首先获得授权。

同样适用于基于本项目的二开项目。

依据 MIT 协议，使用者需自行承担使用本项目的风险与责任，本开源项目开发者与此无关。
