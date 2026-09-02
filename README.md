<p align="center">
  <img src="docs/logo.png" width="150" height="150" alt="one-api-pro logo">
</p>

<p align="center">
  One Api Pro · 基于Go语言的企业级 AI API Gateway
</p>
<p align="center">
  本项目基于 <a href="https://github.com/songquanpeng/one-api">one-api</a> (by <a href="https://github.com/songquanpeng">JustSong</a>) 深度重构开发，感谢原作者的开源贡献。
</p>

<p align="center">
  <a href="http://demo.one-api.pro"><img src="https://img.shields.io/badge/🚀%20Live%20Demo-demo.one--api.pro-165DFF.svg?style=for-the-badge" alt="Live Demo"></a>
</p>

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="license"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/language-Go-00ADD8.svg?logo=go&logoColor=white" alt="language"></a>
  <a href="https://gin-gonic.com/"><img src="https://img.shields.io/badge/framework-Gin-008080.svg?logo=go&logoColor=white" alt="framework"></a>
  <a href="https://vuejs.org/"><img src="https://img.shields.io/badge/frontend-Vue%203-42B883.svg?logo=vue.js&logoColor=white" alt="frontend"></a>
  <a href="https://arco.design/vue"><img src="https://img.shields.io/badge/ui-Arco%20Design-165DFF.svg" alt="ui"></a>
  <a href="https://vitejs.dev/"><img src="https://img.shields.io/badge/build-Vite-646CFF.svg?logo=vite&logoColor=white" alt="build"></a>
  <a href="https://gorm.io/"><img src="https://img.shields.io/badge/database-MySQL%20%7C%20PostgreSQL%20%7C%20SQLite-4479A1.svg?logo=mysql&logoColor=white" alt="database"></a>
  <a href="https://github.com/modelbus/one-api-pro"><img src="https://img.shields.io/badge/cluster-decentralized-FF6B6B.svg" alt="cluster"></a>
</p>

<p align="center">
  <a href="README.md">简体中文</a>
  &nbsp;·&nbsp;
  <a href="readme/README.en.md">English</a>
  &nbsp;·&nbsp;
  <a href="readme/README.zh-TW.md">繁體中文</a>
  &nbsp;·&nbsp;
  <a href="readme/README.ja.md">日本語</a>
  &nbsp;·&nbsp;
  <a href="readme/README.ru.md">Русский</a>
  &nbsp;·&nbsp;
  <a href="readme/README.ko.md">한국어</a>
  &nbsp;·&nbsp;
  <a href="readme/README.ar.md">العربية</a>
  &nbsp;·&nbsp;
  <a href="readme/README.de.md">Deutsch</a>
</p>

---

## 📑 目录

- [🚀 快速开始](#-快速开始)
- [🔧 技术栈](#-技术栈)
  - [Go 后端](#go-后端)
  - [Vue 3 前端](#vue-3-前端)
- [✨ 功能亮点](#-功能亮点)
- [🔥 对比 one-api](#-对比-one-api)
- [📸 截图展示](#-截图展示)
- [⚙️ 配置](#%EF%B8%8F-配置)
  - [🔧 环境变量](#-环境变量)
  - [⌨️ 命令行参数](#%EF%B8%8F-命令行参数)
- [📖 接口文档](#-接口文档)
- [📦 部署](#-部署)
  - [🔨 手动部署](#-手动部署)
  - [🏢 多机部署](#-多机部署)
  - [🌐 集群部署（去中心化多活）](#-集群部署去中心化多活)
- [🗺️ 开发计划](#%EF%B8%8F-开发计划)
- [License](#license)

---

## 🚀 快速开始

### 1. 获取可执行文件

从 [GitHub Releases](https://github.com/modelbus/one-api-pro/releases/latest) 下载预编译版本，或从源码编译：

```bash
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro
```

### 2.（源码构建）构建 Vue 3 前端

```bash
cd web
sh build.sh        # 按 web/THEMES 依次构建主题（默认 default-pro）
cd ..
```

### 3.（源码构建）构建后端

> 后端必须在前端构建完成之后再编译，以嵌入最新前端产物。

```bash
go build -ldflags "-s -w" -o one-api-pro
```

### 4.（可选）一键打包多平台

使用根目录的 `release.sh` 脚本，可一键完成依赖下载、前端构建、多平台交叉编译：

```bash
./release.sh                          # 使用 VERSION 文件作为版本号
./release.sh v0.1.0                   # 指定版本号
./release.sh v0.1.0 --skip-frontend   # 跳过前端构建（复用已有 web/build）
```

> 前置依赖：`go`、`node`、`npm`。版本号来自根目录 `VERSION` 文件（自动兼容有无 `v` 前缀）。

打包产物为**静态编译的裸可执行文件**（无需解压，直接运行），输出到 `dist/` 目录：

```
dist/one-api-pro-linux-amd64
dist/one-api-pro-linux-arm64
dist/one-api-pro-windows-amd64.exe
dist/one-api-pro-darwin-amd64
dist/one-api-pro-darwin-arm64
```

> 其中 `linux-*` 为静态链接，CentOS / Ubuntu 通用。GitHub Releases 由 `.github/workflows/release.yml` 在推送 `v*` tag 时自动构建发布，与本地 `release.sh` 输出逻辑一致。

### 5. 启动

```bash
./one-api-pro --port 3000 --log-dir ./logs
```

访问 `http://localhost:3000`，使用初始账号 `root / 123456` 登录。

> 详细部署方式见 [📦 部署](#-部署)，接口文档见 [📖 接口文档](#-接口文档)。

---

## 🔧 技术栈

本项目基于以下开源技术构建，感谢所有开源项目作者。

### Go 后端

| 技术 | 用途 |
| --- | --- |
| [Gin](https://github.com/gin-gonic/gin) | HTTP Web 框架 |
| [GORM](https://gorm.io) | ORM 库，支持 SQLite / MySQL / PostgreSQL |
| [go-redis/redis](https://github.com/go-redis/redis) | Redis 客户端 |
| [golang-jwt/jwt](https://github.com/golang-jwt/jwt) | JWT 鉴权 |
| [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2) | AWS Bedrock 集成 |
| [Google API Go Client](https://github.com/googleapis/google-api-go-client) | Google Gemini / PaLM2 集成 |
| [pkoukk/tiktoken-go](https://github.com/pkoukk/tiktoken-go) | Token 计数 |
| [gorilla/websocket](https://github.com/gorilla/websocket) | WebSocket 支持（讯飞等渠道） |
| [joho/godotenv](https://github.com/joho/godotenv) | .env 配置文件解析 |

### Vue 3 前端

| 技术 | 用途 |
| --- | --- |
| [Vue 3](https://vuejs.org) | 前端框架（组合式 API） |
| [Vite](https://vitejs.dev) | 构建工具 |
| [Arco Design Vue](https://arco.design/vue) | UI 组件库 |
| [Pinia](https://pinia.vuejs.org) | 状态管理 |
| [Vue Router 4](https://router.vuejs.org) | 路由管理 |
| [Axios](https://axios-http.com) | HTTP 客户端 |
| [ECharts](https://echarts.apache.org) | 数据可视化图表 |
| [vue-i18n](https://vue-i18n.intlify.dev) | 国际化 |

---

## ✨ 功能亮点

One Api Pro 是一个**企业级 AI API 网关**，基于 Go 语言 + Vue 3 全新打造，在保留原版 one-api 全部功能的基础上，进行了架构级重构与企业级增强。

### 🖥️ 可视化仪表盘

全新的 Vue 3 + Arco Design 管理后台，提供数据可视化仪表盘，核心指标、使用趋势、模型用量分布一目了然。

| 核心指标卡 | 使用趋势图 |
|:---:|:---:|
| ![仪表盘首页](docs/Demo-Index.png) | ![仪表盘首页](docs/Demo-Index.png) |

### 🔑 精细的令牌管理

支持多维度令牌管控：可用模型白名单、IP 子网限制、额度上限、过期时间、无限额度，权限粒度细化到单个模型。

| 令牌管理 |
|:---:|
| ![令牌管理](docs/Demo-Token.png) |

### 📦 套餐订阅体系

内置完整的套餐与订阅体系：按 Token / 按请求计费，周期限频（小时 / 周 / 月），按模型精细管控，支持推荐套餐与价格配置。

| 套餐管理 | 订阅管理 |
|:---:|:---:|
| ![套餐管理](docs/Demo-Plan.png) | ![订阅管理](docs/Demo-Subscribe.png) |

### 💳 订单与真实支付

每个套餐下单都会留下一条**完整的订单审计记录**（订单号、用户、套餐快照 JSON、金额、支付方式、状态、支付时间、渠道流水号），支持套餐 / 充值两种订单类型，原生接入 **微信支付 Native**（PC 扫码）与**支付宝当面付**（TradePrecreate），并预置银行 / 线下 / 免费三种管理端通道。套餐升级差价按剩余天数比例自动计算，叠加模式下新旧套餐并行生效，全部规则可在「运营 → 套餐运营」子 Tab 中热切换。

| 订单中心 | 支付配置 |
|:---:|:---:|
| ![订单中心](docs/Demo-Order.png) | ![支付配置](docs/Demo-Payment.png) |

### 🌐 去中心化多活集群

支持去中心化多活集群部署，每个节点独立 MySQL + Redis，通过应用层事件同步实现数据互信，无需共享数据库，天然支持全球多地域就近访问。

| 集群节点管理 |
|:---:|
| ![集群节点管理](docs/Demo-cluster.png) |

### 🧩 其他核心能力

- **30+ 模型平台接入**：OpenAI / Anthropic / Gemini / DeepSeek / 通义千问 / 文心一言 / 讯飞 / 智谱 等主流平台全覆盖，统一 OpenAI 兼容接口
- **精确成本核算**：按 Token 或按次计费，Prompt / Completion / Cached 独立定价，分组折扣叠加，周期用量追踪
- **渠道负载均衡**：按权重随机分配、自动故障切换、冷却 / 禁用策略、渠道并发与 RPM 限流
- **多级权限体系**：Guest / User / Admin / Root 四级权限，修复原版 API 权限漏洞，精细化管理员操作权限
- **企业级安全**：全链路 HTTPS、Token 鉴权、子网 IP 限制、审计日志实时追踪

---

## 🔥 对比 one-api

| 对比维度 | one-api | one-api-pro |
| --- | --- | --- |
| 项目名称 | one-api | one-api-pro |
| Adaptor 架构 | 集中式常量管理（channeltype/define.go 56 行 iota + url.go 平行数组 + helper.go 双层 switch），新增提供商必须修改 4 个框架文件 | 自注册机制（registry + register.go），新增提供商只需创建包 + 注册即可，框架代码零修改 |
| 权限精细化 | 管理员与普通用户权限边界模糊，任何人可通过 API 操作设置项 | 分级权限体系，修复 API 权限漏洞，精细化管理员操作权限 |
| 订阅模式 | 无套餐/订阅体系 | 完整套餐订阅 + 周期限频 + 按模型管控 |
| 去中心化集群 | 无独立集群支持，多机部署需共享 MySQL | 支持去中心化多活集群，每节点独立 MySQL + Redis，通过应用层事件同步实现数据互信，无需共享数据库 |
| 目录结构 | relay/adaptor/ 平铺 40 个目录，基础协议与供应商混在一起，relay/model/ 与根 model/ 冲突 | adaptor/openai/、adaptor/anthropic/ 作为基础协议独立放置，adaptor/provider/ 统一收纳 37 个供应商，relay/schema/ 消除命名冲突 |
| 管理后台 | 3 套前端主题（default/berry/air），基础管理功能 | Vue 3 + Arco Design 全新管理后台，可视化仪表盘 |
| 持续更新 | 原项目已于 2024 年停止更新 | 持续维护更新，针对企业级场景优化 |

---

## 📸 截图展示

### 🖥️ 仪表盘
![仪表盘首页](docs/Demo-Index.png)

### 🔑 令牌管理
![令牌管理](docs/Demo-Token.png)

### 📦 套餐管理
![套餐管理](docs/Demo-Plan.png)

### 🔄 订阅管理
![订阅管理](docs/Demo-Subscribe.png)

### 🌐 集群节点管理
![集群节点管理](docs/Demo-cluster.png)

---

## ⚙️ 配置

系统本身开箱即用。

你可以通过设置环境变量或者命令行参数进行配置；启动后，使用 `root` 用户登录管理后台继续配置。

> **提示**：如果你不知道某个配置项的含义，可以临时删掉值以看到进一步的提示文字。

### 🔧 环境变量

> One Api Pro 支持从 `.env` 文件中读取环境变量，请参照 `.env.example` 文件，使用时请将其重命名为 `.env`。也可通过 `--env` 参数指定配置文件路径（支持相对路径），详见命令行参数一节。

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
    + `TIKTOKEN_CACHE_DIR`：程序启动时会联网下载通用模型的词元编码（如 `gpt-3.5-turbo`、`gpt-4`、`gpt-4o`）。若网络受限或离线，下载超时（约 30 秒）后会自动降级为近似 token 计数（约 `0.38 × 字符数`），服务仍可正常启动。如需精确计费，可在联网环境预先下载编码文件至该目录，再迁移到离线环境。
    + `DATA_GYM_CACHE_DIR`：目前该配置作用与 `TIKTOKEN_CACHE_DIR` 一致，但是优先级没有它高。
16. `RELAY_TIMEOUT`：中继超时设置，单位为秒，默认不设置超时时间。
17. `RELAY_PROXY`：设置后使用该代理来请求 API。
18. `USER_CONTENT_REQUEST_TIMEOUT`：用户上传内容下载超时时间，单位为秒。
19. `USER_CONTENT_REQUEST_PROXY`：设置后使用该代理来请求用户上传的内容，例如图片。
20. `SQLITE_BUSY_TIMEOUT`：SQLite 锁等待超时设置，单位为毫秒，默认 `3000`。
21. `GEMINI_SAFETY_SETTING`：Gemini 的安全设置，默认 `BLOCK_NONE`。
22. `GEMINI_VERSION`：One Api Pro 所使用的 Gemini 版本，默认为 `v1`。
23. `THEME`：系统的主题设置，默认为 `default-pro`（Vue 3 管理后台），也可切换为 `default` / `berry` / `air`（旧 React 主题），具体可选值参考[此处](./web/README.md)。
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
   + 例子：`./one-api-pro --version`
   + 版本号来源（优先级从高到低）：
     1. 当前工作目录或可执行文件同目录下的 `VERSION` 文件（自动兼容有无 `v` 前缀，如 `0.0.2` 或 `v0.0.2`）；
     2. 编译时通过 `-ldflags "-X .../common.Version=..."` 注入的版本号（`release.sh` 与 CI 均会自动注入）；
     3. 源码中的默认值 `common/constants.go`。
   + 因此只需维护根目录 `VERSION` 文件一处，即可让 `--version`、启动日志、`/api/status` 接口与前端仪表盘展示的版本号保持一致。
5. `--help`: 查看命令的使用帮助和参数说明。
   + 例子：`./one-api-pro --help`

---

## 📖 接口文档

完整的接口文档已独立维护在 [docs/API.md](docs/API.md)，涵盖：

- **鉴权机制**：Cookie Session / Access Token / API Key（Bearer Token）三种鉴权方式
- **管理接口**：模型定价、分组折扣、渠道、令牌、用户、日志、兑换码、套餐、订阅等完整 CRUD
- **OpenAI 兼容接口**：`/v1/models`、`/v1/chat/completions`、`/v1/embeddings`、图片、音频、内容审核等
- **集群管理 API**：节点发现、心跳、数据同步等去中心化集群接口

👉 [查看完整接口文档 →](docs/API.md)

---

## 📦 部署

### 🔨 手动部署

#### 1. 获取可执行文件

任选以下方式之一：

**方式一：下载预编译版本（推荐）**

从 [GitHub Releases](https://github.com/modelbus/one-api-pro/releases/latest) 下载对应平台的裸可执行文件（Linux / macOS / Windows），无需解压即可直接运行。

**方式二：使用 release.sh 一键打包**

```shell
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro
./release.sh            # 多平台打包，产物输出到 dist/ 目录
```

**方式三：从源码编译**

```shell
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro

# 构建前端（Vue 3 管理后台，按 web/THEMES 依次构建）
cd web
sh build.sh

# 构建后端（注意：必须在构建前端之后执行，以便嵌入最新前端产物）
cd ..
go build -ldflags "-s -w" -o one-api-pro
```

#### 2. 运行

```shell
chmod u+x one-api-pro
./one-api-pro --port 3000 --log-dir ./logs
```

#### 3. 访问

访问 [http://localhost:3000/](http://localhost:3000/) 并登录。初始账号用户名为 `root`，密码为 `123456`。

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
CLUSTER_DISCOVERY_INTERVAL=30
CLUSTER_DEAD_PING_INTERVAL=120
CLUSTER_MAX_PING_FAILURES=3
CLUSTER_PUSH_INTERVAL=3
CLUSTER_SYNC_LOGS=true
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

> 💡 集群管理 API 详见 [docs/API.md 附录 E：集群管理 API](docs/API.md#附录-e集群管理-api)

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

如果部署后访问出现空白页面，详见 [#97](https://github.com/modelbus/one-api-pro/issues/97)。

---

## 🗺️ 开发计划

### ✅ 已完成

- [x] **架构级重构**：Adaptor 自注册机制，新增供应商零框架修改
- [x] **Vue 3 全新管理后台**：Arco Design + 可视化仪表盘 + 30+ 模型平台图标
- [x] **套餐订阅体系**：按 Token / 按次计费，周期限频，按模型精细管控
- [x] **去中心化多活集群**：GORM 事件驱动 + HTTP 主动推送同步，无需共享数据库
- [x] **精确成本核算**：Prompt / Completion / Cached 独立定价，分组折扣叠加
- [x] **多级权限体系**：Guest / User / Admin / Root 四级，修复原版 API 权限漏洞
- [x] **OpenAI 兼容接口**：完整支持 models / chat / completions / embeddings / images / audio / moderations
- [x] **套餐下单与升级流程**：原生 POST `/api/order/plan` 创建套餐订阅订单，支持 `stack`（叠加）与 `price_diff`（差价升级）两种模式，差价按剩余天数比例自动计算，含同级与降级校验
- [x] **订单审计与订单中心**：新增 `orders` 表（type/source/order_no/plan_info/amount/status/pay_status/pay_method/pay_time/pay_trade_no）持久化所有支付/管理开通流水，前端用户侧 `/plans` 与 `/orders` 页面完整呈现
- [x] **真实支付集成（gopay）**：原生接入微信支付 Native（PC 扫码）与支付宝当面付（TradePrecreate），支付回调走 `/api/payment/{wechat,alipay}/notify` 完成验签 + 订单激活闭环
- [x] **支付 / 套餐运营设置**：「运营设置」下新增「套餐运营」（差价升级 vs 叠加）与「支付」（微信 / 支付宝 / 银行 三通道独立开关 + 证书上传 + 通知 URL 配置），按需显示表单

### 🔄 进行中

- [ ] **更丰富的渠道诊断与智能路由优化**：已具备自动冷却（`CooldownFilter`）、托底降级（`FallbackFilter`）与低成功率自动禁用（`monitor`），下一步补全独立诊断面板 / 节点级 ping 与人工复核流程
- [ ] 更丰富的用量分析报表与导出
- [ ] 多语言国际化（i18n）完善

### 🔭 规划中

- [ ] **支付通道扩展**：Apple Pay、银联、Stripe 等；支持异步退款 API + 自动化退款流水
- [ ] **余额（quota）在线充值**：用户可自助在「个人」区域为账户充值额度，与订阅套餐按需互不干扰
- [ ] **与常见平台财务对接**：对接主流财务 / 对账平台，自动同步充值、消费、退款等财务流水
- [ ] **Token 余量预警机制**：账户 / 令牌 Token 余量低时自动预警，支持多通道通知
- [ ] **日志审计及审计报表**：完整的操作审计日志与可视化审计报表，满足合规要求
- [ ] **AI 智能分析**：基于大模型对用量、成本、渠道健康度进行智能分析与建议
- [ ] 插件化扩展机制
- [ ] 企业级 SSO / LDAP 对接
- [ ] 用量告警与通知渠道扩展（钉钉 / 飞书 / 企业微信等）
- [ ] 更多模型平台的持续接入

> 💡 欢迎提交 PR 或 Issue 参与共建，详见 [Issues](https://github.com/modelbus/one-api-pro/issues)。

---

## License

[MIT License](LICENSE)
