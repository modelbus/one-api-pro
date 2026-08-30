<p align="center">
  <img src="../docs/logo.png" width="150" height="150" alt="one-api-pro logo">
</p>

<p align="center">
  One Api Pro · 基於 Go 語言打造的企業級 AI API Gateway
</p>
<p align="center">
  本項目基於 <a href="https://github.com/songquanpeng/one-api">one-api</a> (by <a href="https://github.com/songquanpeng">JustSong</a>) 深度重構開發，感謝原作者的開源貢獻。
</p>

<p align="center">
  <a href="../LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="license"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/language-Go-00ADD8.svg?logo=go&logoColor=white" alt="language"></a>
  <a href="https://gin-gonic.com/"><img src="https://img.shields.io/badge/framework-Gin-008080.svg?logo=go&logoColor=white" alt="framework"></a>
  <a href="https://vuejs.org/"><img src="https://img.shields.io/badge/frontend-Vue%203-42B883.svg?logo=vue.js&logoColor=white" alt="frontend"></a>
  <a href="https://arco.design/vue"><img src="https://img.shields.io/badge/ui-Arco%20Design-165DFF.svg" alt="ui"></a>
  <a href="https://vitejs.dev/"><img src="https://img.shields.io/badge/build-Vite-646CFF.svg?logo=vite&logoColor=white" alt="build"></a>
  <a href="https://gorm.io/"><img src="https://img.shields.io/badge/database-MySQL%20%7C%20PostgreSQL%20%7C%20SQLite-4479A1.svg?logo=mysql&logoColor=white" alt="database"></a>
  <a href="https://github.com/modelbus/one-api-pro"><img src="https://img.shields.io/badge/cluster-decentralized-FF6B6B.svg" alt="cluster"></a>
</p>

<p align="center">
  <a href="../README.md">简体中文</a>
  &nbsp;·&nbsp;
  <a href="README.en.md">English</a>
  &nbsp;·&nbsp;
  <a href="README.zh-TW.md">繁體中文</a>
  &nbsp;·&nbsp;
  <a href="README.ja.md">日本語</a>
  &nbsp;·&nbsp;
  <a href="README.ru.md">Русский</a>
  &nbsp;·&nbsp;
  <a href="README.ko.md">한국어</a>
  &nbsp;·&nbsp;
  <a href="README.ar.md">العربية</a>
  &nbsp;·&nbsp;
  <a href="README.de.md">Deutsch</a>
</p>

---

## 📑 目錄

- [🚀 快速開始](#-快速開始)
- [🔧 技術棧](#-技術棧)
  - [Go 後端](#go-後端)
  - [Vue 3 前端](#vue-3-前端)
- [✨ 功能亮點](#-功能亮點)
- [🔥 對比 one-api](#-對比-one-api)
- [📸 截圖展示](#-截圖展示)
- [⚙️ 配置](#%EF%B8%8F-配置)
  - [🔧 環境變數](#-環境變數)
  - [⌨️ 命令列參數](#%EF%B8%8F-命令列參數)
- [📖 接口文檔](#-接口文檔)
- [📦 部署](#-部署)
  - [🔨 手動部署](#-手動部署)
  - [🏢 多機部署](#-多機部署)
  - [🌐 集群部署（去中心化多活）](#-集群部署去中心化多活)
- [🗺️ 開發計劃](#%EF%B8%8F-開發計劃)
- [License](#license)

---

## 🚀 快速開始

### 1. 取得可執行檔

從 [GitHub Releases](https://github.com/modelbus/one-api-pro/releases/latest) 下載預編譯版本，或從原始碼編譯：

```bash
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro
```

### 2.（原始碼構建）構建 Vue 3 前端

```bash
cd web
sh build.sh        # 依序構建 web/THEMES 中的每個主題（預設 default-pro）
cd ..
```

### 3.（原始碼構建）構建後端

> 後端必須在前端構建完成之後再編譯，以便嵌入最新前端產物。

```bash
go build -ldflags "-s -w" -o one-api-pro
```

### 4.（可選）一鍵打包多平台

使用根目錄的 `release.sh` 腳本，可一鍵完成依賴下載、前端構建、多平台交叉編譯：

```bash
./release.sh                          # 使用 VERSION 檔案作為版本號
./release.sh v0.1.0                   # 指定版本號
./release.sh v0.1.0 --skip-frontend   # 跳過前端構建（重複使用既有 web/build）
```

> 前置依賴：`go`、`node`、`npm`。版本號來自根目錄 `VERSION` 檔案（自動相容有無 `v` 前綴）。

打包產物為**靜態編譯的裸可執行檔**（無需解壓縮，直接執行），輸出到 `dist/` 目錄：

```
dist/one-api-pro-linux-amd64
dist/one-api-pro-linux-arm64
dist/one-api-pro-windows-amd64.exe
dist/one-api-pro-darwin-amd64
dist/one-api-pro-darwin-arm64
```

> 其中 `linux-*` 為靜態連結，適用於 CentOS / Ubuntu。GitHub Releases 由 `.github/workflows/release.yml` 在推送 `v*` tag 時自動構建發布，與本機 `release.sh` 輸出邏輯一致。

### 5. 啟動

```bash
./one-api-pro --port 3000 --log-dir ./logs
```

造訪 `http://localhost:3000`，使用初始帳號 `root / 123456` 登入。

> 詳細部署方式請見 [📦 部署](#-部署)，接口文檔請見 [📖 接口文檔](#-接口文檔)。

---

## 🔧 技術棧

本項目基於以下開源技術構建，感謝所有開源項目作者。

### Go 後端

| 技術 | 用途 |
| --- | --- |
| [Gin](https://github.com/gin-gonic/gin) | HTTP Web 框架 |
| [GORM](https://gorm.io) | ORM 函式庫，支援 SQLite / MySQL / PostgreSQL |
| [go-redis/redis](https://github.com/go-redis/redis) | Redis 客戶端 |
| [golang-jwt/jwt](https://github.com/golang-jwt/jwt) | JWT 驗證 |
| [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2) | AWS Bedrock 整合 |
| [Google API Go Client](https://github.com/googleapis/google-api-go-client) | Google Gemini / PaLM2 整合 |
| [pkoukk/tiktoken-go](https://github.com/pkoukk/tiktoken-go) | Token 計數 |
| [gorilla/websocket](https://github.com/gorilla/websocket) | WebSocket 支援（如訊飛等通道） |
| [joho/godotenv](https://github.com/joho/godotenv) | `.env` 設定檔解析 |

### Vue 3 前端

| 技術 | 用途 |
| --- | --- |
| [Vue 3](https://vuejs.org) | 前端框架（組合式 API） |
| [Vite](https://vitejs.dev) | 構建工具 |
| [Arco Design Vue](https://arco.design/vue) | UI 元件庫 |
| [Pinia](https://pinia.vuejs.org) | 狀態管理 |
| [Vue Router 4](https://router.vuejs.org) | 路由管理 |
| [Axios](https://axios-http.com) | HTTP 客戶端 |
| [ECharts](https://echarts.apache.org) | 資料視覺化圖表 |
| [vue-i18n](https://vue-i18n.intlify.dev) | 國際化 |

---

## ✨ 功能亮點

One Api Pro 是**企業級 AI API 閘道**，以 Go 語言 + Vue 3 全新打造，在保留原版 one-api 全部功能的基礎上，進行了架構級重構與企業級增強。

### 🖥️ 視覺化儀表板

全新的 Vue 3 + Arco Design 管理後台，提供資料視覺化儀表板，核心指標、使用趨勢、模型用量分布一目了然。

| 核心指標卡 | 使用趨勢圖 |
|:---:|:---:|
| ![儀表板首頁](../docs/Demo-Index.png) | ![儀表板首頁](../docs/Demo-Index.png) |

### 🔑 精細的令牌管理

支援多維度令牌管控：可用模型白名單、IP 子網路限制、額度上限、過期時間、無限額度，權限粒度可細化到單一模型。

| 令牌管理 |
|:---:|
| ![令牌管理](../docs/Demo-Token.png) |

### 📦 套餐訂閱體系

內建完整的套餐與訂閱體系：按 Token / 按請求計費，週期限頻（小時 / 週 / 月），按模型精細管控，支援推薦套餐與價格配置。

| 套餐管理 | 訂閱管理 |
|:---:|:---:|
| ![套餐管理](../docs/Demo-Plan.png) | ![訂閱管理](../docs/Demo-Subscribe.png) |

### 💳 訂單與真實支付

每次套餐下單都會留下一條**完整的訂單稽核紀錄**（訂單編號、使用者、套餐快照 JSON、金額、支付方式、狀態、支付時間、通道流水號），支援套餐／加值兩種訂單類型，原生接入 **微信支付 Native**（PC 掃碼）與**支付寶面對面支付**（TradePrecreate），並預置銀行／線下／免費三種管理端通道。套餐升級差價按剩餘天數比例自動計算，疊加模式下新舊套餐並行生效，全部規則可在「營運 → 套餐營運」子 Tab 中熱切換。

| 訂單中心 | 支付設定 |
|:---:|:---:|
| ![訂單中心](../docs/Demo-Order.png) | ![支付設定](../docs/Demo-Payment.png) |

### 🌐 去中心化多活集群

支援去中心化多活集群部署，每個節點獨立 MySQL + Redis，透過應用層事件同步實現資料互信，無需共用資料庫，天然支援全球多地域就近存取。

| 集群節點管理 |
|:---:|
| ![集群節點管理](../docs/Demo-cluster.png) |

### 🧩 其他核心能力

- **30+ 模型平台接入**：OpenAI / Anthropic / Gemini / DeepSeek / 通義千問 / 文心一言 / 訊飛 / 智譜 等主流平台全覆蓋，統一為 OpenAI 相容接口
- **精確成本核算**：按 Token 或按次計費，Prompt / Completion / Cached 獨立定價，分組折扣疊加，週期用量追蹤
- **通道負載均衡**：按權重隨機分配、自動故障切換、冷卻／停用策略、通道併發與 RPM 限流
- **多級權限體系**：Guest / User / Admin / Root 四級權限，修復原版 API 權限漏洞，精細化管理員操作權限
- **企業級安全**：全鏈路 HTTPS、Token 驗證、子網路 IP 限制、稽核日誌即時追蹤

---

## 🔥 對比 one-api

| 對比維度 | one-api | one-api-pro |
| --- | --- | --- |
| 項目名稱 | one-api | one-api-pro |
| Adaptor 架構 | 集中式常數管理（channeltype/define.go 56 行 iota + url.go 平行陣列 + helper.go 雙層 switch），新增供應商必須修改 4 個框架檔案 | 自註冊機制（registry + register.go），新增供應商只需建立套件 + 註冊即可，框架程式碼零修改 |
| 權限精細化 | 管理員與一般使用者權限邊界模糊，任何人都可透過 API 操作設定項 | 分級權限體系，修復 API 權限漏洞，精細化管理員操作權限 |
| 訂閱模式 | 無套餐／訂閱體系 | 完整套餐訂閱 + 週期限頻 + 按模型管控 |
| 去中心化集群 | 無獨立集群支援，多機部署需共用 MySQL | 支援去中心化多活集群，每節點獨立 MySQL + Redis，透過應用層事件同步實現資料互信，無需共用資料庫 |
| 目錄結構 | relay/adaptor/ 平鋪 40 個目錄，基礎協定與供應商混在一起，relay/model/ 與根目錄 model/ 衝突 | adaptor/openai/、adaptor/anthropic/ 作為基礎協定獨立放置，adaptor/provider/ 統一收納 37 個供應商，relay/schema/ 消除命名衝突 |
| 管理後台 | 3 套前端主題（default/berry/air），基礎管理功能 | Vue 3 + Arco Design 全新管理後台，視覺化儀表板 |
| 持續更新 | 原項目已於 2024 年停止更新 | 持續維護更新，針對企業級場景最佳化 |

---

## 📸 截圖展示

### 🖥️ 儀表板
![儀表板首頁](../docs/Demo-Index.png)

### 🔑 令牌管理
![令牌管理](../docs/Demo-Token.png)

### 📦 套餐管理
![套餐管理](../docs/Demo-Plan.png)

### 🔄 訂閱管理
![訂閱管理](../docs/Demo-Subscribe.png)

### 🌐 集群節點管理
![集群節點管理](../docs/Demo-cluster.png)

---

## ⚙️ 配置

系統本身開箱即用。

你可以透過設定環境變數或命令列參數進行設定；啟動後，使用 `root` 使用者登入管理後台繼續設定。

> **提示**：如果你不知道某個設定項的含義，可以暫時刪掉值以查看進一步的提示文字。

### 🔧 環境變數

> One Api Pro 支援從 `.env` 檔案中讀取環境變數，請參照 `.env.example` 檔案，使用時請將其重新命名為 `.env`。也可透過 `--env` 參數指定設定檔路徑（支援相對路徑），詳見命令列參數一節。

1. `REDIS_CONN_STRING`：設定後將把 Redis 用作快取。
   + 範例：`REDIS_CONN_STRING=redis://default:redispw@localhost:49153`
   + 如果資料庫存取延遲很低，就沒有必要啟用 Redis，啟用後反而會出現資料滯後的問題。
   + 如果需要使用哨兵或集群模式：
     + 則需將該環境變數設定為節點列表，例如：`localhost:49153,localhost:49154,localhost:49155`。
     + 此外還需設定以下環境變數：
       + `REDIS_PASSWORD`：Redis 集群或哨兵模式下的密碼設定。
       + `REDIS_MASTER_NAME`：Redis 哨兵模式下主節點的名稱。
2. `SESSION_SECRET`：設定後將使用固定的工作階段金鑰，如此系統重新啟動後已登入使用者的 cookie 仍然有效。
   + 範例：`SESSION_SECRET=random_string`
3. `SQL_DSN`：設定後將使用指定的資料庫而非 SQLite，請使用 MySQL 或 PostgreSQL。
   + 範例：
     + MySQL：`SQL_DSN=root:123456@tcp(localhost:3306)/oneapi`
     + PostgreSQL：`SQL_DSN=postgres://postgres:123456@localhost:5432/oneapi`（適配中，歡迎回饋）
   + 注意需提前建立資料庫 `oneapi`，無需手動建表，程式將自動建表。
   + 如果使用雲端資料庫：如果雲端伺服器需要驗證身分，需要在連線參數中加入 `?tls=skip-verify`。
   + 請根據你的資料庫設定調整下列參數（或者保持預設值）：
     + `SQL_MAX_IDLE_CONNS`：最大閒置連線數，預設為 `100`。
     + `SQL_MAX_OPEN_CONNS`：最大開啟連線數，預設為 `1000`。
       + 如果出現 `Error 1040: Too many connections` 錯誤，請適當減小此值。
     + `SQL_CONN_MAX_LIFETIME`：連線的最大生命週期，預設為 `60`，單位為分鐘。
4. `LOG_SQL_DSN`：設定後將為 `logs` 表使用獨立的資料庫，請使用 MySQL 或 PostgreSQL。
5. `FRONTEND_BASE_URL`：設定後將把頁面請求重新導向到指定網址，僅限從伺服器設定。
   + 範例：`FRONTEND_BASE_URL=https://openai.justsong.cn`
6. `MEMORY_CACHE_ENABLED`：啟用記憶體快取，會使使用者額度的更新存在一定延遲，可選值為 `true` 和 `false`，未設定則預設為 `false`。
   + 範例：`MEMORY_CACHE_ENABLED=true`
7. `SYNC_FREQUENCY`：在啟用快取的情況下與資料庫同步設定的頻率，單位為秒，預設為 `600` 秒。
   + 範例：`SYNC_FREQUENCY=60`
8. `NODE_TYPE`：設定後將指定節點類型，可選值為 `master` 和 `slave`，未設定則預設為 `master`。
   + 範例：`NODE_TYPE=slave`
9. `CHANNEL_UPDATE_FREQUENCY`：設定後將定期更新通道餘額，單位為分鐘，未設定則不進行更新。
   + 範例：`CHANNEL_UPDATE_FREQUENCY=1440`
10. `CHANNEL_TEST_FREQUENCY`：設定後將定期檢查通道，單位為分鐘，未設定則不進行檢查。
    +範例：`CHANNEL_TEST_FREQUENCY=1440`
11. `POLLING_INTERVAL`：批量更新通道餘額與測試可用性時的請求間隔，單位為秒，預設無間隔。
    + 範例：`POLLING_INTERVAL=5`
12. `BATCH_UPDATE_ENABLED`：啟用資料庫批量更新聚合，會使使用者額度的更新存在一定延遲，可選值為 `true` 和 `false`，未設定則預設為 `false`。
    + 範例：`BATCH_UPDATE_ENABLED=true`
    + 如果你遇到資料庫連線數過多的問題，可以嘗試啟用該選項。
13. `BATCH_UPDATE_INTERVAL=5`：批量更新聚合的時間間隔，單位為秒，預設為 `5`。
    + 範例：`BATCH_UPDATE_INTERVAL=5`
14. 請求頻率限制：
    + `GLOBAL_API_RATE_LIMIT`：全域 API 速率限制（中繼請求除外），單一 IP 三分鐘內的最大請求數，預設為 `180`。
    + `GLOBAL_WEB_RATE_LIMIT`：全域 Web 速率限制，單一 IP 三分鐘內的最大請求數，預設為 `60`。
15. 編碼器快取設定：
    + `TIKTOKEN_CACHE_DIR`：程式啟動時會連網下載通用模型的詞元編碼（如 `gpt-3.5-turbo`、`gpt-4`、`gpt-4o`）。若網路受限或離線，下載逾時（約 30 秒）後會自動降級為近似 token 計數（約 `0.38 × 字元數`），服務仍可正常啟動。如需精確計費，可在連網環境預先下載編碼檔案至該目錄，再遷移到離線環境。
    + `DATA_GYM_CACHE_DIR`：目前此設定作用與 `TIKTOKEN_CACHE_DIR` 一致，但優先級沒有它高。
16. `RELAY_TIMEOUT`：中繼逾時設定，單位為秒，預設不設定逾時時間。
17. `RELAY_PROXY`：設定後使用此代理來請求 API。
18. `USER_CONTENT_REQUEST_TIMEOUT`：使用者上傳內容下載逾時時間，單位為秒。
19. `USER_CONTENT_REQUEST_PROXY`：設定後使用此代理來請求使用者上傳的內容，例如圖片。
20. `SQLITE_BUSY_TIMEOUT`：SQLite 鎖等待逾時設定，單位為毫秒，預設 `3000`。
21. `GEMINI_SAFETY_SETTING`：Gemini 的安全設定，預設 `BLOCK_NONE`。
22. `GEMINI_VERSION`：One Api Pro 所使用的 Gemini 版本，預設為 `v1`。
23. `THEME`：系統的主題設定，預設為 `default-pro`（Vue 3 管理後台），也可切換為 `default` / `berry` / `air`（舊 React 主題），具體可選值請參考[此處](../web/README.md)。
24. `ENABLE_METRIC`：是否依請求成功率停用通道，預設不啟用，可選值為 `true` 和 `false`。
25. `METRIC_QUEUE_SIZE`：請求成功率統計佇列大小，預設為 `10`。
26. `METRIC_SUCCESS_RATE_THRESHOLD`：請求成功率閾值，預設為 `0.8`。
27. `INITIAL_ROOT_TOKEN`：如果設定了此值，則在系統首次啟動時會自動建立一個值為此環境變數值的 root 使用者令牌。
28. `INITIAL_ROOT_ACCESS_TOKEN`：如果設定了此值，則在系統首次啟動時會自動建立一個值為此環境變數的 root 使用者系統管理令牌。
29. `ENFORCE_INCLUDE_USAGE`：是否強制在 stream 模型下回傳 usage，預設不啟用，可選值為 `true` 和 `false`。
30. `TEST_PROMPT`：測試模型時的使用者 prompt，預設為 `Print your model name exactly and do not output without any other text.`。

#### 🌐 集群設定（去中心化多活部署）

> 不設定以下環境變數時，系統以單節點模式運行，無任何副作用。

1. `CLUSTER_ENABLED`：是否啟用集群模式，預設不啟用。
   + 範例：`CLUSTER_ENABLED=true`
2. `CLUSTER_NODE_ID`：節點編號（1-49），必須與 MySQL 的 `auto_increment_offset` 一致，不同節點不能重複。
   + 範例：`CLUSTER_NODE_ID=1`
3. `CLUSTER_NODE_NAME`：節點名稱，便於識別，預設為 `node-{NODE_ID}`。
   + 範例：`CLUSTER_NODE_NAME=node-cn`
4. `CLUSTER_NODE_ADDRESS`：本節點的公開網路存取網址（需包含協定前綴），其他節點透過此網址推送資料。
   + 範例：`CLUSTER_NODE_ADDRESS=https://cn.example.com`
5. `CLUSTER_SECRET`：本節點的初始 secret，**每個節點獨立**。首次啟動時作為初始 secret 寫入資料庫，之後可由 admin 修改。
   + 範例：`CLUSTER_SECRET=MyClusterSecret123abc`
6. `CLUSTER_SEEDS`：種子節點網址（逗號分隔），新節點啟動時向種子節點註冊以取得集群資訊，只需設定一個可達的節點即可。第一個節點可以不設定或設定自己的網址。
   + 範例：`CLUSTER_SEEDS=https://cn.example.com`
   + 多個種子：`CLUSTER_SEEDS=https://cn.example.com,https://us.example.com`
7. `CLUSTER_PUSH_INTERVAL`：同步事件推送間隔，單位為秒，預設 `3`。
8. `CLUSTER_DISCOVERY_INTERVAL`：節點發現間隔，單位為秒，存活節點每週期互相 ping，預設 `30`。
9. `CLUSTER_DEAD_PING_INTERVAL`：失敗節點 ping 間隔，單位為秒，比存活間隔長以減少無效請求，預設 `120`。
10. `CLUSTER_MAX_PING_FAILURES`：連續 ping 失敗次數，達到後標記節點為失敗狀態，預設 `3`。
11. `CLUSTER_SYNC_LOGS`：是否同步日誌表，日誌資料量較大可按需關閉，預設 `true`。
     + 範例：`CLUSTER_SYNC_LOGS=false`
12. `CLUSTER_BATCH_SIZE`：每次推送最大事件數，預設 `50`。

### ⌨️ 命令列參數

1. `--port <port_number>`: 指定伺服器監聽的連接埠號碼，預設為 `3000`。
   + 範例：`--port 3000`
2. `--log-dir <log_dir>`: 指定日誌資料夾，如果未設定，預設儲存至工作目錄的 `logs` 資料夾下。
   + 範例：`--log-dir ./logs`
3. `--env <env_file_path>`: 指定設定檔路徑，支援相對路徑和絕對路徑。未指定時自動載入目前目錄的 `.env` 檔案。
   + 範例：`--env ./config.env`
   + 範例：`--env /etc/one-api-pro/production.env`
   + 多實例部署範例：
     ```bash
     ./one-api-pro --env ./instances/instance1.env --port 3001 &
     ./one-api-pro --env ./instances/instance2.env --port 3002 &
     ```
   + 設定優先級：命令列參數 > 系統環境變數 > `--env` 指定的設定檔 > 預設值
4. `--version`: 列印系統版本號並退出。
   + 範例：`./one-api-pro --version`
   + 版本號來源（優先級由高到低）：
     1. 目前工作目錄或可執行檔同目錄下的 `VERSION` 檔案（自動相容有無 `v` 前綴，如 `0.0.2` 或 `v0.0.2`）；
     2. 編譯時透過 `-ldflags "-X .../common.Version=..."` 注入的版本號（`release.sh` 與 CI 均會自動注入）；
     3. 原始碼中的預設值 `common/constants.go`。
   + 因此只需維護根目錄的 `VERSION` 檔案一處，即可讓 `--version`、啟動日誌、`/api/status` 接口與前端儀表板顯示的版本號保持一致。
5. `--help`: 檢視命令的使用說明與參數說明。
   + 範例：`./one-api-pro --help`

---

## 📖 接口文檔

完整的接口文檔已獨立維護在 [docs/API.md](../docs/API.md)，涵蓋：

- **驗證機制**：Cookie Session / Access Token / API Key（Bearer Token）三種驗證方式
- **管理接口**：模型定價、分組折扣、通道、令牌、使用者、日誌、兌換碼、套餐、訂閱等完整 CRUD
- **OpenAI 相容接口**：`/v1/models`、`/v1/chat/completions`、`/v1/embeddings`、圖片、音訊、內容稽核等
- **集群管理 API**：節點發現、心跳、資料同步等去中心化集群接口

👉 [檢視完整接口文檔 →](../docs/API.md)

---

## 📦 部署

### 🔨 手動部署

#### 1. 取得可執行檔

任選以下方式之一：

**方式一：下載預編譯版本（推薦）**

從 [GitHub Releases](https://github.com/modelbus/one-api-pro/releases/latest) 下載對應平台的裸可執行檔（Linux / macOS / Windows），無需解壓縮即可直接執行。

**方式二：使用 release.sh 一鍵打包**

```shell
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro
./release.sh            # 多平台打包，產物輸出到 dist/ 目錄
```

**方式三：從原始碼編譯**

```shell
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro

# 構建前端（Vue 3 管理後台，依序構建 web/THEMES）
cd web
sh build.sh

# 構建後端（注意：必須在構建前端之後執行，以便嵌入最新前端產物）
cd ..
go build -ldflags "-s -w" -o one-api-pro
```

#### 2. 執行

```shell
chmod u+x one-api-pro
./one-api-pro --port 3000 --log-dir ./logs
```

#### 3. 存取

存取 [http://localhost:3000/](http://localhost:3000/) 並登入。初始帳號使用者名為 `root`，密碼為 `123456`。

### 🏢 多機部署
1. 所有伺服器將 `SESSION_SECRET` 設為相同的值。
2. 必須設定 `SQL_DSN`，使用 MySQL 資料庫而非 SQLite，所有伺服器連線到同一個資料庫。
3. 所有從屬伺服器必須將 `NODE_TYPE` 設為 `slave`，不設定則預設為主伺服器。
4. 設定 `SYNC_FREQUENCY` 後伺服器將定期從資料庫同步設定，在使用遠端資料庫的情況下，建議設定此項並啟用 Redis，無論主從。
5. 從屬伺服器可以選擇設定 `FRONTEND_BASE_URL`，以將頁面請求重新導向到主伺服器。
6. 在從屬伺服器上**分別**裝好 Redis，設定好 `REDIS_CONN_STRING`，這樣可以在快取未過期的情況下做到資料庫零存取，可減少延遲（Redis 集群或哨兵模式的支援請參考環境變數說明）。
7. 如果主伺服器存取資料庫的延遲也較高，則也需要啟用 Redis，並設定 `SYNC_FREQUENCY`，以定期從資料庫同步設定。

環境變數的具體使用方法請詳見[此處](#-環境變數)。

### 🌐 集群部署（去中心化多活）

集群模式允許多個節點各自部署獨立的 One Api Pro + MySQL，透過應用層事件同步實現資料互信，無需共用資料庫。

> **適用場景**：全球多地域部署、就近存取降低延遲、高可用容災、多節點負載均衡。

#### 🗺️ 架構概覽

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

- **去中心化**：所有節點地位平等，無主從之分，任何節點資料變更後主動推送至所有存活節點
- **零侵入**：透過 GORM 回呼捕捉資料變更，不修改現有業務程式碼
- **非同步推送**：資料同步不阻塞主流程，透過背景 goroutine 批量推送
- **衝突解決**：基於 `updated_at` 時間戳比較，只有更新的資料才會寫入
- **限流同步**：通道併發與 RPM 限流計數器透過資料庫表實現跨節點同步
- **單節點相容**：不設定集群環境變數時，系統完全以單節點模式運行

#### 📊 同步範圍

| 資料表 | 是否同步 | 說明 |
|--------|---------|------|
| users | ✅ | 使用者資訊 |
| tokens | ✅ | API 令牌 |
| channels | ✅ | 通道設定 |
| abilities | ✅ | 通道能力 |
| options | ✅ | 系統設定 |
| redemptions | ✅ | 兌換碼 |
| plans | ✅ | 訂閱計劃 |
| user_plans | ✅ | 使用者訂閱 |
| plan_usages | ✅ | 計劃用量 |
| channel_counters | ✅ | 通道限流計數器 |
| cluster_nodes | 🔄 Discovery | 集群節點資訊（由發現機制維護，不走資料同步） |
| logs | ⚠️ 可選 | 日誌資料量較大，透過 `CLUSTER_SYNC_LOGS` 控制 |

#### 🚀 部署步驟

**1. MySQL 設定（每個節點必須使用獨立的 MySQL 實例）**

每個節點都需要一個**獨立的 MySQL 實例**（不能在同一 MySQL 實例中建立多個資料庫來部署多個節點，因為 `auto_increment_offset` 是實例級變數）。

```ini
# 節點 1 的 my.cnf
[mysqld]
server-id = 1
auto_increment_increment = 50
auto_increment_offset = 1
log_bin = mysql-bin
binlog_format = ROW

# 節點 2 的 my.cnf
[mysqld]
server-id = 2
auto_increment_increment = 50
auto_increment_offset = 2
log_bin = mysql-bin
binlog_format = ROW

# 節點 3 的 my.cnf
[mysqld]
server-id = 3
auto_increment_increment = 50
auto_increment_offset = 3
log_bin = mysql-bin
binlog_format = ROW
```

> `auto_increment_increment` 設為 50，最多支援 50 個節點。每個節點的 `offset` 必須與 `CLUSTER_NODE_ID` 一致且互不相同。

> **重要說明：** `auto_increment_increment` 和 `auto_increment_offset` 是 MySQL 的**系統級變數**，對實例內所有資料庫生效，無法為不同資料庫設定不同的值，也無法在表級別設定（MySQL 表選項僅支援 `AUTO_INCREMENT` 起始值，不支援步長）。因此每個節點**必須使用獨立的 MySQL 實例**，不能在同一個 MySQL 實例中透過建立不同資料庫來部署多個節點。如需在同一台機器上執行多個 MySQL 實例，可以使用不同連接埠啟動多個 mysqld 行程，或使用 Docker 執行多個獨立的 MySQL 容器。

> **關於 `server-id` 與 binlog：** `server-id` 在同一集群的所有 MySQL 實例中必須互不相同。強烈建議啟用 `log_bin` 與 `binlog_format=ROW`——它們用於未來的主從複製擴充與 point-in-time recovery。集群資料同步本身不依賴 binlog（透過 GORM 回呼在應用層實現），但 binlog 提供了額外的可靠性保障。

**2. Redis 設定（每個節點必須使用獨立的 Redis 實例）**

每個節點也需要**獨立的 Redis 實例**（連接埠不同或在不同機器上）。Redis 在本集群架構中不用於節點間通訊，只用於本節點的快取、限流等業務用途。

**3. 新節點初始化資料**

新節點上線時，需要先取得既有節點的資料快照：

```bash
# 方式一：從既有節點匯出並匯入
mysqldump -h existing-node -u root -p oneapi > backup.sql
mysql -u root -p oneapi < backup.sql

# 方式二：透過 API 取得快照（需先啟動服務）
curl -H "X-Cluster-Secret: your-secret" \
  "https://existing-node/api/cluster/snapshot?tables=users,tokens,channels,abilities,options,redemptions,plans,user_plans,plan_usages" \
  -o snapshot.json
```

**4. 環境變數設定（完整案例）**

以下是 3 節點集群的完整 `.env` 設定範例。每個節點都使用獨立的 MySQL 與 Redis 實例，連接埠與路徑各不相同。

**節點 1 — 中國節點（`/opt/one-api-pro/node1/.env`）：**
```bash
# ========================
# 基礎設定
# ========================
PORT=3000
SYSTEM_NAME=One Api Pro Cluster

# ========================
# 資料庫（獨立 MySQL 實例）
# ========================
SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node1?charset=utf8mb4&parseTime=True&loc=Local

# ========================
# Redis（獨立 Redis 實例）
# ========================
REDIS_CONN_STRING=redis://127.0.0.1:6379/0

# ========================
# 集群設定
# ========================
CLUSTER_ENABLED=true
CLUSTER_NODE_ID=1
CLUSTER_NODE_NAME=node-cn
CLUSTER_NODE_ADDRESS=https://cn.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me

# 種子節點（首次啟動時引導發現其他節點）
# 第一個節點：填入自己的網址或留空
# 後續節點：填入任意一個已存活節點的網址
CLUSTER_SEEDS=https://cn.example.com,https://us.example.com,https://eu.example.com

# ========================
# 集群調校（可選）
# ========================
CLUSTER_DISCOVERY_INTERVAL=30
CLUSTER_DEAD_PING_INTERVAL=120
CLUSTER_MAX_PING_FAILURES=3
CLUSTER_PUSH_INTERVAL=3
CLUSTER_SYNC_LOGS=true
CLUSTER_BATCH_SIZE=50
```

**節點 2 — 美國節點（`/opt/one-api-pro/node2/.env`）：**
```bash
# 基礎設定
PORT=3001
SYSTEM_NAME=One Api Pro Cluster

# 資料庫（獨立 MySQL 實例，連接埠或機器與節點 1 不同）
SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node2?charset=utf8mb4&parseTime=True&loc=Local

# Redis（獨立 Redis 實例）
REDIS_CONN_STRING=redis://127.0.0.1:6380/0

# 集群設定
CLUSTER_ENABLED=true
CLUSTER_NODE_ID=2
CLUSTER_NODE_NAME=node-us
CLUSTER_NODE_ADDRESS=https://us.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me   # 必須與節點 1 完全一致

# 填入任意一個已存活節點的網址
CLUSTER_SEEDS=https://cn.example.com
```

**節點 3 — 歐洲節點（`/opt/one-api-pro/node3/.env`）：**
```bash
# 基礎設定
PORT=3002
SYSTEM_NAME=One Api Pro Cluster

# 資料庫
SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node3?charset=utf8mb4&parseTime=True&loc=Local

# Redis
REDIS_CONN_STRING=redis://127.0.0.1:6381/0

# 集群設定
CLUSTER_ENABLED=true
CLUSTER_NODE_ID=3
CLUSTER_NODE_NAME=node-eu
CLUSTER_NODE_ADDRESS=https://eu.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me   # 必須與所有節點一致

# 填入任意一個已存活節點的網址
CLUSTER_SEEDS=https://cn.example.com
```

**設定參數對照表：**

| 環境變數 | 節點 1 | 節點 2 | 節點 3 | 說明 |
|---|---|---|---|---|
| `PORT` | 3000 | 3001 | 3002 | 監聽連接埠（同一機器需不同） |
| `SQL_DSN` | ...oneapi_node1 | ...oneapi_node2 | ...oneapi_node3 | 獨立 MySQL 實例 |
| `REDIS_CONN_STRING` | :6379/0 | :6380/0 | :6381/0 | 獨立 Redis 實例 |
| `CLUSTER_NODE_ID` | 1 | 2 | 3 | 節點編號，對應 MySQL `auto_increment_offset` |
| `CLUSTER_NODE_NAME` | node-cn | node-us | node-eu | 節點名稱，便於識別 |
| `CLUSTER_NODE_ADDRESS` | https://cn.example.com | https://us.example.com | https://eu.example.com | 節點公開網路網址（其他節點透過此網址存取） |
| `CLUSTER_SECRET` | 同一個值 | 同一個值 | 同一個值 | **所有節點必須完全一致** |
| `CLUSTER_SEEDS` | 自己的網址或留空 | 任意存活節點 | 任意存活節點 | 首次啟動引導，後續自動發現 |

**5. 啟動命令**

每個節點使用 `--env` 參數載入自己的設定檔：

```bash
# 節點 1
./one-api-pro --env /opt/one-api-pro/node1/.env --port 3000

# 節點 2
./one-api-pro --env /opt/one-api-pro/node2/.env --port 3001

# 節點 3
./one-api-pro --env /opt/one-api-pro/node3/.env --port 3002
```

**6. 啟動順序**

1. 啟動第一個節點（Node A），`CLUSTER_SEEDS` 留空或填入自己的網址
2. 等待 Node A 完全啟動（約 5-10 秒，看到「集群模組初始化完成」日誌）
3. 啟動後續節點，`CLUSTER_SEEDS` 填入任意一個已存活節點的網址
4. 後續節點啟動後會自動 ping 種子節點，傳遞性發現所有其他節點
5. 所有節點啟動後，可透過任一節點的管理後台「設定 → 節點管理」頁面檢視節點狀態

**7. Nginx 負載均衡設定範例（可選）**

```nginx
upstream one_api_cluster {
    ip_hash;  # 基於 IP 雜湊，同一使用者的請求固定到同一節點，保證 session/cache 命中
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

> **使用 `ip_hash` 是關鍵**：保證同一使用者的請求始終到同一節點，避免 plan 限頻、Redis 快取等狀態在不同節點間遺失。

**8. 驗證集群狀態**

部署完成後，可透過以下方式驗證：

```bash
# 檢視節點列表（在任一節點上呼叫）
curl -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  https://cn.example.com/api/cluster_node/

# 應回傳所有節點的列表，包含 status、last_heartbeat、ping_failures 等欄位
```

或在管理後台：**設定 → 節點管理** 頁面檢視節點列表、狀態、最後心跳時間等。

> 💡 集群管理 API 詳見 [docs/API.md 附錄 E：集群管理 API](../docs/API.md#附錄-e集群管理-api)

#### ⚠️ 注意事項

- 每個節點必須有獨立的 MySQL 實例與 Redis 實例，不共用資料庫
- `CLUSTER_SECRET` 在所有節點間必須一致，請使用強密碼並妥善保管
- `CLUSTER_NODE_ID` 在所有節點間必須互不相同，且與 MySQL `auto_increment_offset` 一致
- `CLUSTER_NODE_ADDRESS` 必須是其他節點可存取的公開網路網址（包含協定前綴如 `https://`）
- 新節點上線前的資料初始化需手動完成（從線上節點拉取快照）
- 日誌表（logs）資料量較大，可透過 `CLUSTER_SYNC_LOGS=false` 關閉日誌同步
- MySQL 的 `auto_increment_increment` 和 `auto_increment_offset` 必須與 `CLUSTER_NODE_ID` 設定一致
- 節點發現採用 ping 雙向註冊機制，失敗節點不會被刪除，只標記為 status=2，網路恢復後自動復活
- `CLUSTER_SEEDS` 只是首次啟動的引導；節點一旦透過 ping 發現其他節點，就不再依賴 SEEDS
- 節點離線期間其他節點產生的變更**不會自動補傳**，離線節點重新上線後需拉取快照補齊資料

#### 📝 關於「本機節點」自我註冊

每個節點啟動時會在自己的 `cluster_nodes` 表中寫入一筆本機記錄（`node_id` 等於本機設定的 `CLUSTER_NODE_ID`）。這是**刻意的設計**，原因如下：

1. **管理後台展示**：在「設定 → 節點管理」頁面，管理者需要看到本機資訊（網址、狀態、心跳時間等），以便排查問題
2. **節點發現傳遞性**：當節點 B 收到節點 A 的 ping 請求時，A 在回應中回傳完整的節點列表（包含 A 自身）。B 收到後將其合併到本機表中。這樣 C 透過 B 的回應也能學習到 A 的存在
3. **存活判斷依據**：本機記錄的 `last_heartbeat` 由本機每 30 秒自動更新一次（在 `discoverOnce` 函式中），反映本機正常運行的狀態

**自我註冊不會導致循環同步資料**。系統在 5 個層面做了防護：

| 防護點 | 作用 |
|---|---|
| ① `GetAllRemoteNodes` SQL 過濾 | 發現時 SQL 加 `node_id != ?` 排除本機 |
| ② `GetAliveNodesForSync` SQL 過濾 | 推送時 SQL 加 `node_id != ?` 排除本機 |
| ③ `handlePing` 拒絕自 ping | 明確拒絕 `req.NodeId == NodeID` |
| ④ `mergeDiscoveredNodes` 跳過本機 | 合併發現節點時跳過本機 |
| ⑤ `ApplyEvents` 跳過本機事件 | 應用事件時跳過本機產生的事件 |

資料流是單向的：從本機推到遠端，從遠端拉過來應用到本機，**永遠不會有迴路**。

管理後台會在本機節點名稱旁顯示「本機」藍色徽章，並對本機停用「刪除」與「手動 Ping」操作（這兩個操作對本機無意義）。

#### 🔐 關於「每節點獨立 secret」

每個節點有**自己的 secret**，不再使用全域共用的 secret。設計原因：

1. **安全性**：一個節點洩露 secret 不會影響其他節點
2. **管理靈活**：每個節點可以獨立輪換自己的 secret
3. **自動發現**：節點間 ping 時自動攜帶自己的 secret 供對方保存

**Secret 生命週期**：
- 節點首次啟動：用 `CLUSTER_SECRET` 環境變數作為初始值，寫入 `cluster_nodes.secret_key` 欄位
- 後續啟動：從 `cluster_nodes.secret_key` 讀取
- Admin 可以在「節點管理」頁面修改其他節點的 secret
- ping 時 `X-Cluster-Secret` 標頭 = **目標節點**的 secret（從本地 DB 查詢）

**新增節點流程**：
1. 在節點 A 上新增 B 節點記錄，填入 B 的 `CLUSTER_SECRET` 值
2. 在節點 B 上新增 A 節點記錄，填入 A 的 `CLUSTER_SECRET` 值
3. A ping B：用 B 的 secret；B 接收：驗證 B 自己的 secret ✓
4. B 回應中攜帶 A、B 各自的 secret，A 更新本地保存

#### 🗑️ 關於「軟刪除節點」

Admin 刪除節點時**不實際刪除**記錄，而是設定 `disabled = true`：

- 防止被刪除節點「自動長回來」（ping 機制會重新註冊）
- 已停用的節點仍然會回應 ping（讓對方知道本節點線上），但不會取得本節點資訊
- 實際刪除需要手動 SQL：`DELETE FROM cluster_nodes WHERE node_id = ?`

#### 🔄 關於「資料同步機制」（重要）

**集群資料同步**完全依賴 **GORM 事件 + HTTP 主動推送**機制：
- 任何業務表的 INSERT/UPDATE/DELETE 操作 → GORM 回呼捕捉 → 寫入 `sync_events` 表 → Pusher goroutine 推送到所有存活節點
- 接收方用 `WithSkipHook` 寫入本地資料庫（不會回環）
- 接收方跳過 `event.NodeId == 本機 NodeID` 的事件（雙重保險）

**架構權衡**：本設計**不實現跨節點主動拉取**，原因如下：
1. **侵入業務**：跨節點拉取需要知道每張表的業務唯一欄位，會侵入業務程式碼
2. **主鍵衝突**：跨節點自增 ID 不連續（不同的 `auto_increment_offset`），使用源節點 id 會破壞 offset 設計
3. **複雜度高**：維護成本高，可靠性提升有限
4. **主動推送夠用**：95% 的場景（節點線上時的常規同步）完全由推送涵蓋

**已知限制與運維要求**：
- 節點離線期間其他節點產生的資料變更 → **永久遺失**（推送是即時的）
- 節點重新上線後無法自動補齊離線期間的資料
- 新節點加入後只能接收到加入之後的資料變更，無歷史資料
- **運維補救**：使用 `mysqldump` 從其他節點匯出後匯入

**典型部署場景對照**：

| 場景 | 是否需要拉取 | 處理方式 |
|---|---|---|
| 節點永久線上 | ❌ | 推送完全夠用 |
| 節點偶爾重啟（分鐘級） | ⚠️ | 短時離線資料遺失，運維可接受 |
| 節點頻繁維護 | ❌ | 推送繼續，重啟後立即恢復 |
| 新節點加入集群 | ❌ | DBA 手動 `mysqldump` 初始化 |
| 節點長期離線後恢復 | ❌ | DBA 手動 `mysqldump` 補齊 |

如果部署後存取出現空白頁面，請詳見 [#97](https://github.com/modelbus/one-api-pro/issues/97)。

---

## 🗺️ 開發計劃

### ✅ 已完成

- [x] **架構級重構**：Adaptor 自註冊機制，新增供應商零框架修改
- [x] **Vue 3 全新管理後台**：Arco Design + 視覺化儀表板 + 30+ 模型平台圖示
- [x] **套餐訂閱體系**：按 Token / 按次計費，週期限頻，按模型精細管控
- [x] **去中心化多活集群**：GORM 事件驅動 + HTTP 主動推送同步，無需共用資料庫
- [x] **精確成本核算**：Prompt / Completion / Cached 獨立定價，分組折扣疊加
- [x] **多級權限體系**：Guest / User / Admin / Root 四級，修復原版 API 權限漏洞
- [x] **OpenAI 相容接口**：完整支援 models / chat / completions / embeddings / images / audio / moderations
- [x] **套餐下單與升級流程**：原生 POST `/api/order/plan` 建立套餐訂閱訂單，支援 `stack`（疊加）與 `price_diff`（差價升級）兩種模式，差價按剩餘天數比例自動計算，含同級與降級校驗
- [x] **訂單稽核與訂單中心**：新增 `orders` 表（type/source/order_no/plan_info/amount/status/pay_status/pay_method/pay_time/pay_trade_no）持久化所有支付／管理開通流水，前端使用者側 `/plans` 與 `/orders` 頁面完整呈現
- [x] **真實支付整合（gopay）**：原生接入微信支付 Native（PC 掃碼）與支付寶面對面支付（TradePrecreate），支付回呼走 `/api/payment/{wechat,alipay}/notify` 完成驗籤 + 訂單啟用閉環
- [x] **支付／套餐營運設定**：「營運設定」下新增「套餐營運」（差價升級 vs 疊加）與「支付」（微信／支付寶／銀行 三通道獨立開關 + 憑證上傳 + 通知 URL 設定），按需顯示表單

### 🔄 進行中

- [ ] **更豐富的通道診斷與智慧路由最佳化**：已具備自動冷卻（`CooldownFilter`）、託底降級（`FallbackFilter`）與低成功率自動停用（`monitor`），下一步補全獨立診斷面板／節點級 ping 與人工複核流程
- [ ] 更豐富的用量分析報表與匯出
- [ ] 多語言國際化（i18n）完善

### 🔭 規劃中

- [ ] **支付通道擴充**：Apple Pay、銀聯、Stripe 等；支援異步退款 API + 自動化退款流水
- [ ] **餘額（quota）線上加值**：使用者可自助在「個人」區域為帳戶加值額度，與訂閱套餐按需互不干擾
- [ ] **與常見平台財務對接**：對接主流財務／對帳平台，自動同步加值、消費、退款等財務流水
- [ ] **Token 餘量預警機制**：帳戶／令牌 Token 餘量低時自動預警，支援多通道通知
- [ ] **日誌稽核與稽核報表**：完整的操作稽核日誌與視覺化稽核報表，滿足合規要求
- [ ] **AI 智慧分析**：基於大模型對用量、成本、通道健康度進行智慧分析與建議
- [ ] 外掛化擴充機制
- [ ] 企業級 SSO / LDAP 對接
- [ ] 用量告警與通知通道擴充（釘釘／飛書／企業微信等）
- [ ] 更多模型平台的持續接入

> 💡 歡迎提交 PR 或 Issue 參與共建，請詳見 [Issues](https://github.com/modelbus/one-api-pro/issues)。

---

## License

[MIT License](../LICENSE)
