<p align="center">
  <img src="../docs/logo.png" width="150" height="150" alt="one-api-pro logo">
</p>

<p align="center">
  One Api Pro · Go 言語で構築されたエンタープライズ向け AI API Gateway
</p>
<p align="center">
  <a href="https://github.com/songquanpeng/one-api">one-api</a> (by <a href="https://github.com/songquanpeng">JustSong</a>) を深くリファクタリングして開発されました。原作者のオープンソースへの貢献に感謝します。
</p>

<p align="center">
  👉 <strong>オンライン Demo を見る</strong>：<a href="http://demo.one-api.pro">http://demo.one-api.pro</a>
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

## 📑 目次

- [🚀 クイックスタート](#-クイックスタート)
- [🔧 技術スタック](#-技術スタック)
  - [Go バックエンド](#go-バックエンド)
  - [Vue 3 フロントエンド](#vue-3-フロントエンド)
- [✨ 機能ハイライト](#-機能ハイライト)
- [🔥 one-api との比較](#-one-api-との比較)
- [📸 スクリーンショット](#-スクリーンショット)
- [⚙️ 設定](#%EF%B8%8F-設定)
  - [🔧 環境変数](#-環境変数)
  - [⌨️ コマンドライン引数](#%EF%B8%8F-コマンドライン引数)
- [📖 API ドキュメント](#-api-ドキュメント)
- [📦 デプロイ](#-デプロイ)
  - [🔨 手動デプロイ](#-手動デプロイ)
  - [🏢 マルチホストデプロイ](#-マルチホストデプロイ)
  - [🌐 分散クラスタデプロイ](#-分散クラスタデプロイ)
- [🗺️ ロードマップ](#%EF%B8%8F-ロードマップ)
- [License](#license)

---

## 🚀 クイックスタート

### 1. 実行可能ファイルを取得

[GitHub Releases](https://github.com/modelbus/one-api-pro/releases/latest) からプリコンパイル済みのバージョンをダウンロードするか、ソースからビルドします：

```bash
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro
```

### 2.（ソースビルド）Vue 3 フロントエンドをビルド

```bash
cd web
sh build.sh        # web/THEMES に従って各テーマをビルド（デフォルトは default-pro）
cd ..
```

### 3.（ソースビルド）バックエンドをビルド

> バックエンドはフロントエンドのビルド完了後にコンパイルする必要があります。最新のフロントエンド成果物を組み込むためです。

```bash
go build -ldflags "-s -w" -o one-api-pro
```

### 4.（オプション）ワンクリックでマルチプラットフォームをパッケージング

ルートディレクトリの `release.sh` スクリプトを使用すると、依存関係のダウンロード、フロントエンドビルド、マルチプラットフォームのクロスコンパイルをワンクリックで実行できます：

```bash
./release.sh                          # VERSION ファイルをバージョン番号として使用
./release.sh v0.1.0                   # バージョン番号を指定
./release.sh v0.1.0 --skip-frontend   # フロントエンドのビルドをスキップ（既存の web/build を再利用）
```

> 前提要件：`go`、`node`、`npm`。バージョン番号はルートディレクトリの `VERSION` ファイルから取得されます（`v` プレフィックスの有無を自動判別）。

パッケージング成果物は**静的リンクされた裸の実行可能ファイル**です（解凍不要、直接実行可能）、`dist/` ディレクトリに出力されます：

```
dist/one-api-pro-linux-amd64
dist/one-api-pro-linux-arm64
dist/one-api-pro-windows-amd64.exe
dist/one-api-pro-darwin-amd64
dist/one-api-pro-darwin-arm64
```

> このうち `linux-*` は静的リンクされており、CentOS / Ubuntu で共通に使用できます。GitHub Releases は `.github/workflows/release.yml` により `v*` タグがプッシュされた際に自動でビルド・公開されます。ローカルの `release.sh` の出力ロジックと一致します。

### 5. 起動

```bash
./one-api-pro --port 3000 --log-dir ./logs
```

`http://localhost:3000` にアクセスし、初期アカウント `root / 123456` でログインします。

> 詳細なデプロイ方法は [📦 デプロイ](#-デプロイ)、API ドキュメントは [📖 API ドキュメント](#-api-ドキュメント) を参照してください。

---

## 🔧 技術スタック

本プロジェクトは以下のオープンソース技術を基に構築されています。すべてのオープンソースプロジェクトの作者に感謝します。

### Go バックエンド

| 技術 | 用途 |
| --- | --- |
| [Gin](https://github.com/gin-gonic/gin) | HTTP Web フレームワーク |
| [GORM](https://gorm.io) | ORM ライブラリ、SQLite / MySQL / PostgreSQL をサポート |
| [go-redis/redis](https://github.com/go-redis/redis) | Redis クライアント |
| [golang-jwt/jwt](https://github.com/golang-jwt/jwt) | JWT 認証 |
| [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2) | AWS Bedrock 連携 |
| [Google API Go Client](https://github.com/googleapis/google-api-go-client) | Google Gemini / PaLM2 連携 |
| [pkoukk/tiktoken-go](https://github.com/pkoukk/tiktoken-go) | Token カウント |
| [gorilla/websocket](https://github.com/gorilla/websocket) | WebSocket 対応（訊飛などのチャネル） |
| [joho/godotenv](https://github.com/joho/godotenv) | .env 設定ファイルの読み込み |

### Vue 3 フロントエンド

| 技術 | 用途 |
| --- | --- |
| [Vue 3](https://vuejs.org) | フロントエンドフレームワーク（Composition API） |
| [Vite](https://vitejs.dev) | ビルドツール |
| [Arco Design Vue](https://arco.design/vue) | UI コンポーネントライブラリ |
| [Pinia](https://pinia.vuejs.org) | 状態管理 |
| [Vue Router 4](https://router.vuejs.org) | ルーティング管理 |
| [Axios](https://axios-http.com) | HTTP クライアント |
| [ECharts](https://echarts.apache.org) | データビジュアライゼーショングラフ |
| [vue-i18n](https://vue-i18n.intlify.dev) | 国際化 |

---

## ✨ 機能ハイライト

One Api Pro は**エンタープライズ向け AI API ゲートウェイ**であり、Go 言語 + Vue 3 で新たに構築されました。オリジナルの one-api の全機能を維持した上で、アーキテクチャレベルのリファクタリングとエンタープライズ向けの強化が施されています。

### 🖥️ ビジュアルダッシュボード

新しい Vue 3 + Arco Design 管理画面は、データビジュアライゼーションのダッシュボードを提供し、主要指標、利用トレンド、モデル別の使用量分布を一目で把握できます。

| 主要指標カード | 利用トレンドグラフ |
|:---:|:---:|
| ![ダッシュボードホーム](../docs/Demo-Index.png) | ![ダッシュボードホーム](../docs/Demo-Index.png) |

### 🔑 きめ細かいトークン管理

多次元のトークン管理をサポート：使用可能モデルのホワイトリスト、IP サブネット制限、クォータ上限、有効期限、無制限クォータ。権限の粒度は単一モデルまで細分化できます。

| トークン管理 |
|:---:|
| ![トークン管理](../docs/Demo-Token.png) |

### 📦 プラン・サブスクリプション体系

プランとサブスクリプションの完全な体系を内蔵：Token / リクエスト単位の課金、周期レート制限（時間 / 週 / 月）、モデル別のきめ細かい管理、おすすめプランと価格設定をサポートします。

| プラン管理 | サブスクリプション管理 |
|:---:|:---:|
| ![プラン管理](../docs/Demo-Plan.png) | ![サブスクリプション管理](../docs/Demo-Subscribe.png) |

### 💳 注文と実決済

各プランの注文は**完全な注文監査記録**（注文番号、ユーザー、プランスナップショット JSON、金額、支払い方法、ステータス、支払い時刻、チャネル取引番号）を残します。プラン / チャージの 2 種類の注文タイプをサポートし、**微信支付 Native**（PC スキャン）と**支付宝当面付**（TradePrecreate）をネイティブに連携し、銀行 / オフライン / 無料の 3 つの管理側チャネルも用意されています。プランアップグレードの差額は残り日数の割合で自動計算され、積み上げモードでは新旧プランが並行して有効になり、すべてのルールは「運営 → プラン運営」サブタブでホットスイッチできます。

| 注文センター | 支払い設定 |
|:---:|:---:|
| ![注文センター](../docs/Demo-Order.png) | ![支払い設定](../docs/Demo-Payment.png) |

### 🌐 分散型アクティブ・アクティブクラスタ

分散型アクティブ・アクティブクラスタのデプロイをサポートします。各ノードは独立した MySQL + Redis を持ち、アプリケーション層のイベント同期でデータ相互信頼を実現します。データベース共有は不要で、グローバルな複数地域への最寄りアクセスを自然にサポートします。

| クラスタノード管理 |
|:---:|
| ![クラスタノード管理](../docs/Demo-cluster.png) |

### 🧩 その他のコア機能

- **30+ のモデルプラットフォーム連携**：OpenAI / Anthropic / Gemini / DeepSeek / 通義千問 / 文心一言 / 訊飛 / 智譜 などの主要プラットフォームを網羅し、OpenAI 互換インターフェースに統一
- **正確なコスト計算**：Token 単位またはリクエスト単位の課金、Prompt / Completion / Cached の独立した価格設定、グループ割引の積み上げ、周期使用量の追跡
- **チャネル負荷分散**：重み付きランダム割り当て、自動フェイルオーバー、クールダウン / 無効化ポリシー、チャネル並行性と RPM レート制限
- **多段階権限システム**：Guest / User / Admin / Root の 4 段階、オリジナルの API 権限の脆弱性を修正、管理者操作権限をきめ細かく制御
- **エンタープライズ向けセキュリティ**：全経路 HTTPS、Token 認証、サブネット IP 制限、監査ログのリアルタイム追跡

---

## 🔥 one-api との比較

| 比較観点 | one-api | one-api-pro |
| --- | --- | --- |
| プロジェクト名 | one-api | one-api-pro |
| Adaptor アーキテクチャ | 集中定数管理（channeltype/define.go の 56 行の iota + url.go の平行配列 + helper.go の二重 switch）、プロバイダーを追加するには 4 つのフレームワークファイルを変更する必要がある | 自己登録メカニズム（registry + register.go）、プロバイダーを追加するにはパッケージを作成して登録するだけでよく、フレームワークコードの変更はゼロ |
| 権限のきめ細かさ | 管理者と一般ユーザーの権限境界が曖昧で、誰でも API で設定項目を操作できる | 階層型権限システム、API 権限の脆弱性を修正、管理者操作権限をきめ細かく制御 |
| サブスクリプションモード | プラン/サブスクリプション体系なし | 完全なプラン・サブスクリプション + 周期レート制限 + モデル別管理 |
| 分散クラスタ | 独立したクラスタサポートなし、マルチホストデプロイでは MySQL を共有する必要がある | 分散型アクティブ・アクティブクラスタをサポート、各ノードが独立した MySQL + Redis を持ち、アプリケーション層のイベント同期でデータ相互信頼を実現、データベース共有不要 |
| ディレクトリ構造 | relay/adaptor/ に 40 個のディレクトリをフラットに配置、基本プロトコルとプロバイダーが混在、relay/model/ がルートの model/ と衝突 | adaptor/openai/、adaptor/anthropic/ を基本プロトコルとして独立配置、adaptor/provider/ に 37 社のプロバイダーを統合、relay/schema/ で名前の衝突を解消 |
| 管理画面 | 3 種類のフロントエンドテーマ（default/berry/air）、基本管理機能 | Vue 3 + Arco Design の新しい管理画面、ビジュアルダッシュボード |
| 継続更新 | 元プロジェクトは 2024 年に更新を停止した | 継続的に保守・更新、エンタープライズ向けシナリオに最適化 |

---

## 📸 スクリーンショット

### 🖥️ ダッシュボード
![ダッシュボードホーム](../docs/Demo-Index.png)

### 🔑 トークン管理
![トークン管理](../docs/Demo-Token.png)

### 📦 プラン管理
![プラン管理](../docs/Demo-Plan.png)

### 🔄 サブスクリプション管理
![サブスクリプション管理](../docs/Demo-Subscribe.png)

### 🌐 クラスタノード管理
![クラスタノード管理](../docs/Demo-cluster.png)

---

## ⚙️ 設定

システム自体はインストール直後から使用可能です。

環境変数やコマンドライン引数で設定できます。起動後、`root` ユーザーで管理画面にログインして設定を続行できます。

> **ヒント**：設定項目の意味がわからない場合は、一時的に値（バリュー）を削除すると詳細なヒントテキストが表示されます。

### 🔧 環境変数

> One Api Pro は `.env` ファイルから環境変数を読み取ることをサポートします。`.env.example` ファイルを参照し、使用時には `.env` に名前を変更してください。`--env` 引数で設定ファイルのパス（相対パス対応）を指定することもできます。詳細はコマンドライン引数の節を参照してください。

1. `REDIS_CONN_STRING`：設定後、Redis をキャッシュとして使用します。
   - 例：`REDIS_CONN_STRING=redis://default:redispw@localhost:49153`
   - データベースアクセスの遅延が非常に低い場合、Redis を有効にする必要はありません。有効にすると逆にデータ遅延の問題が発生します。
   - センチネルまたはクラスタモードを使用する場合：
     - その環境変数をノードリストに設定する必要があります。例：`localhost:49153,localhost:49154,localhost:49155`。
     - その他に以下の環境変数も設定する必要があります：
       - `REDIS_PASSWORD`：Redis クラスタまたはセンチネルモードでのパスワード設定。
       - `REDIS_MASTER_NAME`：Redis センチネルモードでのマスターノードの名前。
2. `SESSION_SECRET`：設定後、固定のセッションキーを使用します。これによりシステム再起動後もログイン済みユーザーの cookie が引き続き有効になります。
   - 例：`SESSION_SECRET=random_string`
3. `SQL_DSN`：設定後、SQLite ではなく指定のデータベースを使用します。MySQL または PostgreSQL を使用してください。
   - 例：
     - MySQL：`SQL_DSN=root:123456@tcp(localhost:3306)/oneapi`
     - PostgreSQL：`SQL_DSN=postgres://postgres:123456@localhost:5432/oneapi`（対応中、フィードバック歓迎）
   - 事前にデータベース `oneapi` を作成しておく必要があります。テーブルの手動作成は不要で、プログラムが自動で作成します。
   - クラウドデータベースを使用する場合：クラウドサーバーが認証を要求する場合は、接続パラメータに `?tls=skip-verify` を追加する必要があります。
   - データベース設定に応じて以下のパラメータを変更してください（またはデフォルト値を維持）：
     - `SQL_MAX_IDLE_CONNS`：最大空き接続数、デフォルトは `100`。
     - `SQL_MAX_OPEN_CONNS`：最大開いた接続数、デフォルトは `1000`。
       - `Error 1040: Too many connections` エラーが出る場合は、この値を適宜減らしてください。
     - `SQL_CONN_MAX_LIFETIME`：接続の最大ライフタイム、デフォルトは `60`、単位は分。
4. `LOG_SQL_DSN`：設定後、`logs` テーブル用に独立したデータベースを使用します。MySQL または PostgreSQL を使用してください。
5. `FRONTEND_BASE_URL`：設定後、ページリクエストを指定されたアドレスにリダイレクトします。サーバーからのみ設定可能です。
   - 例：`FRONTEND_BASE_URL=https://openai.justsong.cn`
6. `MEMORY_CACHE_ENABLED`：メモリキャッシュを有効にします。ユーザー額度の更新に一定の遅延が生じます。`true` と `false` を選択でき、未設定の場合は `false` になります。
   - 例：`MEMORY_CACHE_ENABLED=true`
7. `SYNC_FREQUENCY`：キャッシュを有効にした場合にデータベースと設定を同期する頻度、単位は秒、デフォルトは `600` 秒。
   - 例：`SYNC_FREQUENCY=60`
8. `NODE_TYPE`：設定後、ノードタイプを指定します。`master` と `slave` から選択でき、未設定の場合は `master` がデフォルト。
   - 例：`NODE_TYPE=slave`
9. `CHANNEL_UPDATE_FREQUENCY`：設定後、チャネル残高を定期的に更新します。単位は分、未設定の場合は更新しません。
   - 例：`CHANNEL_UPDATE_FREQUENCY=1440`
10. `CHANNEL_TEST_FREQUENCY`：設定後、チャネルを定期的にチェックします。単位は分、未設定の場合はチェックしません。
    - 例：`CHANNEL_TEST_FREQUENCY=1440`
11. `POLLING_INTERVAL`：チャネル残高の一括更新および可用性テスト時のリクエスト間隔、単位は秒、デフォルトは間隔なし。
    - 例：`POLLING_INTERVAL=5`
12. `BATCH_UPDATE_ENABLED`：データベースのバッチ更新集約を有効にします。ユーザー額度の更新に一定の遅延が生じます。`true` と `false` を選択でき、未設定の場合は `false` になります。
    - 例：`BATCH_UPDATE_ENABLED=true`
    - データベース接続数が多すぎる問題が発生した場合は、このオプションの有効化を試してください。
13. `BATCH_UPDATE_INTERVAL=5`：バッチ更新集約の時間間隔、単位は秒、デフォルトは `5`。
    - 例：`BATCH_UPDATE_INTERVAL=5`
14. リクエスト頻度制限：
    - `GLOBAL_API_RATE_LIMIT`：グローバル API レート制限（リレーリクエストを除く）、単一 IP の 3 分以内の最大リクエスト数、デフォルトは `180`。
    - `GLOBAL_WEB_RATE_LIMIT`：グローバル Web レート制限、単一 IP の 3 分以内の最大リクエスト数、デフォルトは `60`。
15. エンコーダキャッシュ設定：
    - `TIKTOKEN_CACHE_DIR`：プログラム起動時にオンラインで汎用モデルのトークンエンコーダ（例：`gpt-3.5-turbo`、`gpt-4`、`gpt-4o`）をダウンロードします。ネットワークが制限されているかオフラインの場合、ダウンロードがタイムアウト（約 30 秒）すると自動的に近似 token カウント（約 `0.38 × 文字数`）にフォールバックし、サービスは正常に起動できます。正確な課金が必要な場合は、オンライン環境で事前にエンコーダファイルをこのディレクトリにダウンロードし、オフライン環境へ移行してください。
    - `DATA_GYM_CACHE_DIR`：現時点ではこの設定の役割は `TIKTOKEN_CACHE_DIR` と同じですが、それより優先度は低くなります。
16. `RELAY_TIMEOUT`：リレータイムアウト設定、単位は秒、デフォルトではタイムアウト未設定。
17. `RELAY_PROXY`：設定後、そのプロキシを使用して API をリクエストします。
18. `USER_CONTENT_REQUEST_TIMEOUT`：ユーザーコンテンツのダウンロードタイムアウト時間、単位は秒。
19. `USER_CONTENT_REQUEST_PROXY`：設定後、そのプロキシを使用してユーザーがアップロードしたコンテンツ（例：画像）をリクエストします。
20. `SQLITE_BUSY_TIMEOUT`：SQLite ロック待ちタイムアウト設定、単位はミリ秒、デフォルトは `3000`。
21. `GEMINI_SAFETY_SETTING`：Gemini の安全設定、デフォルトは `BLOCK_NONE`。
22. `GEMINI_VERSION`：One Api Pro が使用する Gemini のバージョン、デフォルトは `v1`。
23. `THEME`：システムのテーマ設定、デフォルトは `default-pro`（Vue 3 管理画面）、`default` / `berry` / `air`（旧 React テーマ）にも切り替え可能です。具体的な選択肢は[こちら](../web/README.md)を参照してください。
24. `ENABLE_METRIC`：リクエスト成功率に基づいてチャネルを無効化するかどうか、デフォルトでは無効、`true` と `false` を選択できます。
25. `METRIC_QUEUE_SIZE`：リクエスト成功率統計のキューサイズ、デフォルトは `10`。
26. `METRIC_SUCCESS_RATE_THRESHOLD`：リクエスト成功率のしきい値、デフォルトは `0.8`。
27. `INITIAL_ROOT_TOKEN`：設定した場合、システム初回起動時にこの環境変数値を値とする root ユーザートークンを自動作成します。
28. `INITIAL_ROOT_ACCESS_TOKEN`：設定した場合、システム初回起動時にこの環境変数値を値とする root ユーザーのシステム管理用アクセストークンを自動作成します。
29. `ENFORCE_INCLUDE_USAGE`：stream モデルでの usage 返却在を強制するかどうか、デフォルトでは無効、`true` と `false` を選択できます。
30. `TEST_PROMPT`：モデルテスト時のユーザー prompt、デフォルトは `Print your model name exactly and do not output without any other text.`。

#### 🌐 クラスタ設定（分散型アクティブ・アクティブデプロイ）

> 以下の環境変数を設定しない場合、システムは単一ノードモードで動作し、副作用はありません。

1. `CLUSTER_ENABLED`：クラスタモードを有効にするかどうか、デフォルトでは無効。
   - 例：`CLUSTER_ENABLED=true`
2. `CLUSTER_NODE_ID`：ノード番号（1-49）、MySQL の `auto_increment_offset` と一致する必要があり、異なるノードで重複できません。
   - 例：`CLUSTER_NODE_ID=1`
3. `CLUSTER_NODE_NAME`：ノード名、識別しやすくするため、デフォルトは `node-{NODE_ID}`。
   - 例：`CLUSTER_NODE_NAME=node-cn`
4. `CLUSTER_NODE_ADDRESS`：本ノードの公開アクセスアドレス（プロトコルプレフィックスを含める必要があります）、他のノードはこのアドレスにデータをプッシュします。
   - 例：`CLUSTER_NODE_ADDRESS=https://cn.example.com`
5. `CLUSTER_SECRET`：本ノードの初期 secret、**各ノード独立**。初回起動時に初期 secret としてデータベースに書き込まれ、その後 admin が変更可能。
   - 例：`CLUSTER_SECRET=MyClusterSecret123abc`
6. `CLUSTER_SEEDS`：シードノードアドレス（カンマ区切り）、新規ノード起動時にシードノードに登録してクラスタ情報を取得します。到達可能なノードを 1 つ設定すれば十分です。最初のノードは設定しないか、自分のアドレスを設定できます。
   - 例：`CLUSTER_SEEDS=https://cn.example.com`
   - 複数のシード：`CLUSTER_SEEDS=https://cn.example.com,https://us.example.com`
7. `CLUSTER_PUSH_INTERVAL`：同期イベントプッシュ間隔、単位は秒、デフォルトは `3`。
8. `CLUSTER_DISCOVERY_INTERVAL`：ノード発見間隔、単位は秒、存続ノードは毎周期お互いに ping し合い、デフォルトは `30`。
9. `CLUSTER_DEAD_PING_INTERVAL`：失敗ノードの ping 間隔、単位は秒、存続間隔より長くして無効なリクエストを減らします、デフォルトは `120`。
10. `CLUSTER_MAX_PING_FAILURES`：連続 ping 失敗回数、これに達するとノードを失敗状態としてマーク、デフォルトは `3`。
11. `CLUSTER_SYNC_LOGS`：ログテーブルを同期するかどうか、ログデータ量が大きい場合は必要に応じて無効化、デフォルトは `true`。
    - 例：`CLUSTER_SYNC_LOGS=false`
12. `CLUSTER_BATCH_SIZE`：毎回のプッシュ最大イベント数、デフォルトは `50`。

### ⌨️ コマンドライン引数

1. `--port <port_number>`: サーバーがリッスンするポート番号を指定、デフォルトは `3000`。
   - 例：`--port 3000`
2. `--log-dir <log_dir>`: ログフォルダを指定、未設定の場合は作業ディレクトリの `logs` フォルダに保存。
   - 例：`--log-dir ./logs`
3. `--env <env_file_path>`: 設定ファイルのパスを指定、相対パスと絶対パスに対応。未指定時はカレントディレクトリの `.env` ファイルを自動ロード。
   - 例：`--env ./config.env`
   - 例：`--env /etc/one-api-pro/production.env`
   - 複数インスタンスデプロイ例：
     ```bash
     ./one-api-pro --env ./instances/instance1.env --port 3001 &
     ./one-api-pro --env ./instances/instance2.env --port 3002 &
     ```
   - 設定の優先順位：コマンドライン引数 > システム環境変数 > `--env` で指定した設定ファイル > デフォルト値
4. `--version`: システムバージョン番号を表示して終了。
   - 例：`./one-api-pro --version`
   - バージョン番号の由来（優先度が高い順）：
     1. カレント作業ディレクトリまたは実行可能ファイルと同ディレクトリの `VERSION` ファイル（`v` プレフィックスの有無を自動判別、例：`0.0.2` または `v0.0.2`）；
     2. ビルド時に `-ldflags "-X .../common.Version=..."` で注入されたバージョン番号（`release.sh` と CI の両方が自動注入）；
     3. ソース内のデフォルト値 `common/constants.go`。
   - したがって、ルートディレクトリの `VERSION` ファイルを 1 箇所だけ管理すれば、`--version`、起動ログ、`/api/status` インターフェース、フロントエンドダッシュボードに表示されるバージョン番号を一致させられます。
5. `--help`: コマンドの使用ヘルプと引数の説明を表示。
   - 例：`./one-api-pro --help`

---

## 📖 API ドキュメント

完全な API ドキュメントは [docs/API.md](../docs/API.md) に独立して管理されており、以下を網羅しています：

- **認証メカニズム**：Cookie Session / Access Token / API Key（Bearer Token）の 3 つの認証方式
- **管理インターフェース**：モデル価格、グループ割引、チャネル、トークン、ユーザー、ログ、利用コード、プラン、サブスクリプションなどの完全な CRUD
- **OpenAI 互換インターフェース**：`/v1/models`、`/v1/chat/completions`、`/v1/embeddings`、画像、音声、モデレーションなど
- **クラスタ管理 API**：ノード発見、ハートビート、データ同期などの分散型クラスタインターフェース

👉 [完全な API ドキュメントを見る →](../docs/API.md)

---

## 📦 デプロイ

### 🔨 手動デプロイ

#### 1. 実行可能ファイルを取得

以下のいずれかの方法を選びます：

**方法一：プリコンパイル済みバージョンをダウンロード（推奨）**

[GitHub Releases](https://github.com/modelbus/one-api-pro/releases/latest) から対応プラットフォーム（Linux / macOS / Windows）の裸の実行可能ファイルをダウンロードします。解凍不要でそのまま実行できます。

**方法二：release.sh でワンクリックパッケージング**

```shell
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro
./release.sh            # マルチプラットフォームのパッケージング、成果物は dist/ に出力
```

**方法三：ソースからコンパイル**

```shell
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro

# フロントエンドをビルド（Vue 3 管理画面、web/THEMES に従って順次ビルド）
cd web
sh build.sh

# バックエンドをビルド（注意：最新のフロントエンド成果物を組み込むため、必ずフロントエンドのビルド後に実行）
cd ..
go build -ldflags "-s -w" -o one-api-pro
```

#### 2. 実行

```shell
chmod u+x one-api-pro
./one-api-pro --port 3000 --log-dir ./logs
```

#### 3. アクセス

[http://localhost:3000/](http://localhost:3000/) にアクセスしてログインします。初期アカウントのユーザー名は `root`、パスワードは `123456` です。

### 🏢 マルチホストデプロイ
1. すべてのサーバーで `SESSION_SECRET` を同じ値に設定します。
2. `SQL_DSN` を必ず設定し、SQLite ではなく MySQL データベースを使用し、すべてのサーバーが同じデータベースに接続します。
3. すべての従属サーバーで `NODE_TYPE` を `slave` に設定する必要があります。設定しない場合はメインサーバーがデフォルトになります。
4. `SYNC_FREQUENCY` を設定すると、サーバーはデータベースから設定を定期的に同期します。リモートデータベースを使用する場合は、主従を問わずこの項目の設定と Redis の有効化を推奨します。
5. 従属サーバーは必要に応じて `FRONTEND_BASE_URL` を設定し、ページリクエストをメインサーバーにリダイレクトできます。
6. 従属サーバーには**それぞれ** Redis をインストールし、`REDIS_CONN_STRING` を設定します。これにより、キャッシュが未失効の間はデータベースへのアクセスをゼロにでき、遅延を削減できます（Redis クラスタまたはセンチネルモードの対応は環境変数の説明を参照）。
7. メインサーバーのデータベースアクセス遅延も高い場合は、Redis の有効化と `SYNC_FREQUENCY` の設定が必要です。データベースから設定を定期的に同期するためです。

環境変数の具体的な使用方法は[こちら](#-環境変数)を参照してください。

### 🌐 分散クラスタデプロイ

クラスタモードでは、複数のノードがそれぞれ独立した One Api Pro + MySQL をデプロイし、アプリケーション層のイベント同期でデータ相互信頼を実現できます。データベース共有は不要です。

> **適用シナリオ**：グローバルな複数地域デプロイ、最寄りアクセスによる遅延削減、高可用性・災害復旧、複数ノードの負荷分散。

#### 🗺️ アーキテクチャ概要

```
                    ┌─────────────┐
                    │  Nginx/LB   │  （統一入口，ip_hash 負載均衡）
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

#### ⭐ コア特性

- **分散型**：全ノードが対等で主従の区別なし、どのノードでもデータ変更後は全存続ノードへ能動的にプッシュ
- **ゼロ侵入**：GORM コールバックでデータ変更を捕捉し、既存のビジネスコードを変更しない
- **非同期プッシュ**：データ同期はメインプロセスをブロックせず、バックグラウンドの goroutine で一括プッシュ
- **競合解決**：`updated_at` タイムスタンプの比較に基づき、より新しいデータのみ書き込む
- **レート制限同期**：チャネル並行性と RPM レート制限カウンタをデータベーステーブル経由でクロスノード同期
- **単一ノード互換**：クラスタ環境変数を設定しない場合、システムは完全に単一ノードモードで動作

#### 📊 同期範囲

| データテーブル | 同期するか | 説明 |
|------|------|------|
| users | ✅ | ユーザー情報 |
| tokens | ✅ | API トークン |
| channels | ✅ | チャネル設定 |
| abilities | ✅ | チャネル機能 |
| options | ✅ | システム設定 |
| redemptions | ✅ | 利用コード |
| plans | ✅ | サブスクリプションプラン |
| user_plans | ✅ | ユーザーサブスクリプション |
| plan_usages | ✅ | プラン使用量 |
| channel_counters | ✅ | チャネルレート制限カウンタ |
| cluster_nodes | 🔄 Discovery | クラスタノード情報（発見メカニズムで管理され、データ同期は行われない） |
| logs | ⚠️ オプション | ログデータ量が大きいため、`CLUSTER_SYNC_LOGS` で制御 |

#### 🚀 デプロイ手順

**1. MySQL 設定（各ノードは独立した MySQL インスタンスを使用する必要があります）**

各ノードには**独立した MySQL インスタンス**が必要です（同一の MySQL インスタンス内に複数のデータベースを作成して複数のノードをデプロイすることはできません。`auto_increment_offset` はインスタンスレベルの変数だからです）。

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

> `auto_increment_increment` を 50 に設定すると、最大 50 ノードをサポートできます。各ノードの `offset` は `CLUSTER_NODE_ID` と一致し、かつ互いに異なる必要があります。

> **重要：** `auto_increment_increment` と `auto_increment_offset` は MySQL の**システムレベル変数**であり、インスタンス内のすべてのデータベースに影響します。データベースごとに異なる値を設定することも、テーブルレベルの設定もできません（MySQL のテーブルオプションは `AUTO_INCREMENT` の開始値のみをサポートし、増分はサポートしません）。したがって、各ノードは**独立した MySQL インスタンス**を使用する必要があり、同一の MySQL インスタンス内で異なるデータベースを作成して複数のノードをデプロイすることはできません。同一マシン上で複数の MySQL インスタンスを実行する必要がある場合は、異なるポートで複数の mysqld プロセスを開始するか、Docker で複数の独立した MySQL コンテナを実行してください。

> **`server-id` と binlog について：** `server-id` は同一クラスタのすべての MySQL インスタンスで互いに異なる必要があります。`log_bin` と `binlog_format=ROW` は強く推奨します。これらは将来のマスタースレーブレプリケーション拡張と point-in-time recovery のためのものです。クラスタデータ同期自体は binlog に依存しません（GORM コールバックでアプリケーション層に実装）が、binlog は追加の信頼性保証を提供します。

**2. Redis 設定（各ノードは独立した Redis インスタンスを使用する必要があります）**

各ノードにも**独立した Redis インスタンス**が必要です（ポートが異なるか、異なるマシン上）。Redis はこのクラスタアーキテクチャではノード間通信には使用されず、本ノードのキャッシュ、レート制限などのビジネス用途のみに使用されます。

**3. 新規ノードの初期化データ**

新規ノードをオンラインにする際は、まず既存ノードのデータスナップショットを取得する必要があります：

```bash
# 方法一：既存ノードからエクスポートしてインポート
mysqldump -h existing-node -u root -p oneapi > backup.sql
mysql -u root -p oneapi < backup.sql

# 方法二：API でスナップショットを取得（サービスの起動が必要）
curl -H "X-Cluster-Secret: your-secret" \
  "https://existing-node/api/cluster/snapshot?tables=users,tokens,channels,abilities,options,redemptions,plans,user_plans,plan_usages" \
  -o snapshot.json
```

**4. 環境変数設定（完全な例）**

以下は 3 ノードクラスタの完全な `.env` 設定例です。各ノードは独立した MySQL と Redis インスタンスを使用し、ポートとパスはそれぞれ異なります。

**ノード 1 — 中国ノード（`/opt/one-api-pro/node1/.env`）：**
```bash
# ========================
# 基本設定
# ========================
PORT=3000
SYSTEM_NAME=One Api Pro Cluster

# ========================
# データベース（独立した MySQL インスタンス）
# ========================
SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node1?charset=utf8mb4&parseTime=True&loc=Local

# ========================
# Redis（独立した Redis インスタンス）
# ========================
REDIS_CONN_STRING=redis://127.0.0.1:6379/0

# ========================
# クラスタ設定
# ========================
CLUSTER_ENABLED=true
CLUSTER_NODE_ID=1
CLUSTER_NODE_NAME=node-cn
CLUSTER_NODE_ADDRESS=https://cn.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me

# シードノード（初回起動時に他のノードを発見するためのガイド）
# 最初のノード：自分のアドレスを記入するか空のまま
# 後続のノード：稼働中の任意のノードのアドレスを記入
CLUSTER_SEEDS=https://cn.example.com,https://us.example.com,https://eu.example.com

# ========================
# クラスタ調整（オプション）
# ========================
CLUSTER_DISCOVERY_INTERVAL=30
CLUSTER_DEAD_PING_INTERVAL=120
CLUSTER_MAX_PING_FAILURES=3
CLUSTER_PUSH_INTERVAL=3
CLUSTER_SYNC_LOGS=true
CLUSTER_BATCH_SIZE=50
```

**ノード 2 — 米国ノード（`/opt/one-api-pro/node2/.env`）：**
```bash
# 基本設定
PORT=3001
SYSTEM_NAME=One Api Pro Cluster

# データベース（独立した MySQL インスタンス。ポートまたはマシンがノード 1 と異なる）
SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node2?charset=utf8mb4&parseTime=True&loc=Local

# Redis（独立した Redis インスタンス）
REDIS_CONN_STRING=redis://127.0.0.1:6380/0

# クラスタ設定
CLUSTER_ENABLED=true
CLUSTER_NODE_ID=2
CLUSTER_NODE_NAME=node-us
CLUSTER_NODE_ADDRESS=https://us.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me   # ノード 1 と完全に一致している必要があります

# 稼働中の任意のノードのアドレスを記入
CLUSTER_SEEDS=https://cn.example.com
```

**ノード 3 — 欧州ノード（`/opt/one-api-pro/node3/.env`）：**
```bash
# 基本設定
PORT=3002
SYSTEM_NAME=One Api Pro Cluster

# データベース
SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node3?charset=utf8mb4&parseTime=True&loc=Local

# Redis
REDIS_CONN_STRING=redis://127.0.0.1:6381/0

# クラスタ設定
CLUSTER_ENABLED=true
CLUSTER_NODE_ID=3
CLUSTER_NODE_NAME=node-eu
CLUSTER_NODE_ADDRESS=https://eu.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me   # 全ノードで一致している必要があります

# 稼働中の任意のノードのアドレスを記入
CLUSTER_SEEDS=https://cn.example.com
```

**設定パラメータ対照表：**

| 環境変数 | ノード 1 | ノード 2 | ノード 3 | 説明 |
|---|---|---|---|---|
| `PORT` | 3000 | 3001 | 3002 | リッスンポート（同一マシンでは異なる必要がある） |
| `SQL_DSN` | ...oneapi_node1 | ...oneapi_node2 | ...oneapi_node3 | 独立した MySQL インスタンス |
| `REDIS_CONN_STRING` | :6379/0 | :6380/0 | :6381/0 | 独立した Redis インスタンス |
| `CLUSTER_NODE_ID` | 1 | 2 | 3 | ノード番号、MySQL の `auto_increment_offset` に対応 |
| `CLUSTER_NODE_NAME` | node-cn | node-us | node-eu | ノード名、識別しやすくする |
| `CLUSTER_NODE_ADDRESS` | https://cn.example.com | https://us.example.com | https://eu.example.com | ノードの公開アドレス（他のノードはこのアドレスでアクセス） |
| `CLUSTER_SECRET` | 同じ値 | 同じ値 | 同じ値 | **全ノードで完全に一致させる必要がある** |
| `CLUSTER_SEEDS` | 自分のアドレスまたは空 | 任意の存続ノード | 任意の存続ノード | 初回起動時のガイド、以降は自動発見 |

**5. 起動コマンド**

各ノードは `--env` 引数で自分の設定ファイルを読み込みます：

```bash
# ノード 1
./one-api-pro --env /opt/one-api-pro/node1/.env --port 3000

# ノード 2
./one-api-pro --env /opt/one-api-pro/node2/.env --port 3001

# ノード 3
./one-api-pro --env /opt/one-api-pro/node3/.env --port 3002
```

**6. 起動順序**

1. 最初のノード（Node A）を起動、`CLUSTER_SEEDS` は空にするか自分のアドレスを入力
2. Node A が完全に起動するのを待ちます（約 5-10 秒、「クラスタモジュール初期化完了」のログを確認）
3. 後続ノードを起動、`CLUSTER_SEEDS` に稼働中の任意のノードのアドレスを入力
4. 後続ノードは起動後、シードノードに自動で ping し、推移的にすべての他のノードを発見
5. 全ノード起動後、任意のノードの管理画面「設定 → ノード管理」ページでノード状態を確認

**7. Nginx 負荷分散設定例（オプション）**

```nginx
upstream one_api_cluster {
    ip_hash;  # IP ハッシュに基づき、同一ユーザーのリクエストを同じノードに固定し、session/cache のヒット率を保証
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

> **`ip_hash` を使うことが重要**：同一ユーザーのリクエストを常に同じノードへ固定し、プランレート制限、Redis キャッシュなどの状態が異なるノード間で失われないようにします。

**8. クラスタ状態の検証**

デプロイ完了後、以下の方法で検証できます：

```bash
# ノード一覧を表示（任意のノード上で呼び出し）
curl -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  https://cn.example.com/api/cluster_node/

# 全ノードのリスト（status、last_heartbeat、ping_failures などのフィールド）が返されるはずです
```

または管理画面：**設定 → ノード管理** ページでノード一覧、状態、最終ハートビート時刻などを確認できます。

> 💡 クラスタ管理 API の詳細は [docs/API.md 付録 E：クラスタ管理 API](../docs/API.md#附录-e集群管理-api) を参照してください

#### ⚠️ 注意事項

- 各ノードは独立した MySQL インスタンスと Redis インスタンスが必要で、データベースを共有しません
- `CLUSTER_SECRET` は全ノードで一致させる必要があり、強力なパスワードを使用して適切に保管してください
- `CLUSTER_NODE_ID` は全ノードで互いに異なる必要があり、MySQL の `auto_increment_offset` と一致させる必要があります
- `CLUSTER_NODE_ADDRESS` は他のノードからアクセス可能な公開アドレスである必要があります（`https://` などのプロトコルプレフィックスを含む）
- 新規ノードのオンライン前のデータ初期化は手動で行う必要があります（オンラインノードからスナップショットを取得）
- ログテーブル（logs）はデータ量が大きいため、`CLUSTER_SYNC_LOGS=false` でログ同期を無効化できます
- MySQL の `auto_increment_increment` と `auto_increment_offset` は `CLUSTER_NODE_ID` 設定と一致させる必要があります
- ノード発見は ping の双方向登録メカニズムを採用しており、失敗ノードは削除されず、status=2 とマークされるだけです。ネットワーク復旧後は自動的に復活します
- `CLUSTER_SEEDS` は初回起動時のガイドにすぎません。ノードが ping で他のノードを発見した後は SEEDS に依存しません
- ノードがオフラインの間に他のノードで発生した変更は**自動で再送されません**。オフラインノードが再オンラインした後はスナップショットを取得してデータを補完する必要があります

#### 📝 「ローカルノード」の自己登録について

各ノードは起動時に自分の `cluster_nodes` テーブルに 1 件のローカルレコードを書き込みます（`node_id` は本機設定の `CLUSTER_NODE_ID` と等しい）。これは**意図的な設計**であり、理由は以下の通りです：

1. **管理画面での表示**：「設定 → ノード管理」ページで、管理者が本機情報（アドレス、状態、ハートビート時刻など）を確認し、問題を調査できるようにするため
2. **ノード発見の推移性**：ノード B がノード A の ping リクエストを受けたとき、A は応答で完全なノードリスト（A 自身を含む）を返します。B はそれを受信してローカルテーブルにマージします。これにより C は B の応答を通じて A の存在も学習できます
3. **存続判断の根拠**：本機レコードの `last_heartbeat` は本機によって 30 秒ごとに自動更新され（`discoverOnce` 関数内）、本機が正常に稼働している状態を反映します

**自己登録により循環同期データが発生することはありません**。システムは 5 つのレベルで保護しています：

| 防護ポイント | 作用 |
|---|---|
| ① `GetAllRemoteNodes` SQL フィルタ | 発見時、SQL に `node_id != ?` を付加して本機を除外 |
| ② `GetAliveNodesForSync` SQL フィルタ | プッシュ時、SQL に `node_id != ?` を付加して本機を除外 |
| ③ `handlePing` は自己 ping を拒否 | `req.NodeId == NodeID` を明示的に拒否 |
| ④ `mergeDiscoveredNodes` は本機をスキップ | 発見ノードのマージ時に本機をスキップ |
| ⑤ `ApplyEvents` は本機イベントをスキップ | イベント適用時に本機が生成したイベントをスキップ |

データフローは一方向です：本機からリモートへプッシュ、リモートから取得して本機に適用。**永遠にループはありません**。

管理画面は本機のノード名の横に「本機」の青いバッジを表示し、本機に対する「削除」と「手動 Ping」操作を無効化します（これらは本機にとって意味がありません）。

#### 🔐 「ノードごとの独立した secret」について

各ノードは**自分の secret** を持ち、グローバル共有 secret は使用しません。設計理由：

1. **セキュリティ**：1 つのノードの secret が漏えいしても他のノードには影響しません
2. **管理の柔軟性**：各ノードは自分の secret を独立してローテートできます
3. **自動発見**：ノード間 ping 時、相手に保存させるため自分の secret を自動的に携行します

**Secret のライフサイクル**：
- ノード初回起動：`CLUSTER_SECRET` 環境変数を初期値として使用し、`cluster_nodes.secret_key` フィールドに書き込み
- 以降の起動：`cluster_nodes.secret_key` から読み取る
- Admin は「ノード管理」ページで他のノードの secret を変更できます
- ping 時の `X-Cluster-Secret` ヘッダー = **ターゲットノード**の secret（ローカル DB から検索）

**新規ノード追加フロー**：
1. ノード A で B ノードのレコードを追加し、B の `CLUSTER_SECRET` 値を入力
2. ノード B で A ノードのレコードを追加し、A の `CLUSTER_SECRET` 値を入力
3. A が B に ping：B の secret を使用；B が受信：B 自身の secret を検証 ✓
4. B の応答に A、B 各々の secret が含まれ、A がローカル保存分を更新

#### 🗑️ 「ノードのソフト削除」について

Admin がノードを削除するとき、**物理的に削除せず** `disabled = true` を設定します：

- 削除されたノードが「自動的に復活する」のを防ぐ（ping メカニズムが再登録するため）
- 無効化されたノードは依然として ping に応答します（相手に本ノードがオンラインであることを知らせる）が、本ノードの情報は取得しません
- 物理削除には手動 SQL が必要：`DELETE FROM cluster_nodes WHERE node_id = ?`

#### 🔄 「データ同期メカニズム」について（重要）

**クラスタデータ同期**は完全に **GORM イベント + HTTP 能動プッシュ** メカニズムに依存しています：
- 任意のビジネステーブルの INSERT/UPDATE/DELETE 操作 → GORM コールバックで捕捉 → `sync_events` テーブルに書き込み → Pusher goroutine が全存続ノードへプッシュ
- 受信側は `WithSkipHook` でローカルデータベースに書き込みます（ループバックなし）
- 受信側は `event.NodeId == 本機 NodeID` のイベントをスキップします（二重の保険）

**アーキテクチャのトレードオフ**：本設計は**クロスノードの能動プルを実装しません**。理由は以下の通りです：
1. **ビジネスへの侵入**：クロスノードプルでは各テーブルのビジネス固有フィールドを知る必要があり、ビジネスコードを汚染します
2. **主キー競合**：クロスノードの自動増分 ID は連続しません（異なる `auto_increment_offset`）、ソースノードの id を使用すると offset 設計を破壊します
3. **複雑さが高い**：保守コストが高く、信頼性の向上は限定的
4. **能動プッシュで十分**：95% のシナリオ（ノードがオンライン時の通常同期）は完全にプッシュでカバーされます

**既知の制限と運用要件**：
- ノードがオフラインの間に他のノードで発生したデータ変更 → **恒久的に消失**（プッシュはリアルタイム）
- ノードが再オンラインした後、オフライン期間のデータを自動で補完できません
- 新規ノードは参加後のデータ変更のみを受信でき、履歴データはありません
- **運用上の対策**：`mysqldump` で他のノードからエクスポートしてからインポート

**典型的なデプロイシナリオ対照**：

| シナリオ | プルが必要か | 処理方法 |
|---|---|---|
| ノードが恒久的にオンライン | ❌ | プッシュで十分 |
| ノードが時々再起動（分単位） | ⚠️ | 短時間のオフラインデータ消失、運用上許容可能 |
| ノードが頻繁にメンテナンス | ❌ | プッシュ継続、再起動後すぐに復旧 |
| 新規ノードがクラスタに参加 | ❌ | DBA が手動で `mysqldump` 初期化 |
| ノードが長期オフライン後に復旧 | ❌ | DBA が手動で `mysqldump` を補完 |

デプロイ後にアクセスして空白ページが表示される場合は、[#97](https://github.com/modelbus/one-api-pro/issues/97) を参照してください。

---

## 🗺️ ロードマップ

### ✅ 完了

- [x] **アーキテクチャレベルのリファクタリング**：Adaptor 自己登録メカニズム、新規プロバイダー追加でフレームワーク変更ゼロ
- [x] **Vue 3 の新しい管理画面**：Arco Design + ビジュアルダッシュボード + 30+ モデルプラットフォームのアイコン
- [x] **プラン・サブスクリプション体系**：Token / リクエスト単位の課金、周期レート制限、モデル別のきめ細かい管理
- [x] **分散型アクティブ・アクティブクラスタ**：GORM イベント駆動 + HTTP 能動プッシュ同期、データベース共有不要
- [x] **正確なコスト計算**：Prompt / Completion / Cached の独立した価格設定、グループ割引の積み上げ
- [x] **多段階権限システム**：Guest / User / Admin / Root の 4 段階、オリジナルの API 権限の脆弱性を修正
- [x] **OpenAI 互換インターフェース**：models / chat / completions / embeddings / images / audio / moderations を完全サポート
- [x] **プラン注文とアップグレードフロー**：ネイティブ `POST /api/order/plan` でサブスクリプション注文を作成、`stack`（積み上げ）と `price_diff`（差額アップグレード）の 2 モードをサポート、差額は残り日数の割合で自動計算、同レベル・ダウングレードの検証を含む
- [x] **注文監査と注文センター**：新規 `orders` テーブル（type/source/order_no/plan_info/amount/status/pay_status/pay_method/pay_time/pay_trade_no）で全決済・管理側開通のフローを永続化、フロントエンドのユーザー側 `/plans` と `/orders` ページで完全に表示
- [x] **実決済連携（gopay）**：ネイティブに微信支付 Native（PC スキャン）と支付宝当面付（TradePrecreate）を連携、決済コールバックは `/api/payment/{wechat,alipay}/notify` で検証 + 注文アクティブ化のループを完了
- [x] **決済 / プラン運営設定**：「運営設定」の下に「プラン運営」（差額アップグレード vs 積み上げ）と「決済」（微信 / 支付宝 / 銀行 の 3 チャネル独立スイッチ + 証明書アップロード + 通知 URL 設定）を新設、必要に応じてフォームを表示

### 🔄 進行中

- [ ] **より豊富なチャネル診断とインテリジェントルーティング最適化**：自動クールダウン（`CooldownFilter`）、フォールバック（`FallbackFilter`）、低成功率自動無効化（`monitor`）は既に実装済み。次のステップとしては独立した診断パネル / ノードレベルの ping と手動レビューフローを補完
- [ ] より豊富な使用量分析レポートとエクスポート
- [ ] 多言語国際化（i18n）の充実

### 🔭 計画中

- [ ] **決済チャネルの拡張**：Apple Pay、銀聯、Stripe など；非同期払い戻し API + 自動化払い戻しフローをサポート
- [ ] **残高（quota）のオンラインチャージ**：ユーザーが「個人」エリアでアカウントへクォータをチャージ可能、サブスクリプションプランとは必要に応じて独立
- [ ] **一般的なプラットフォームとの財務連携**：主流の財務 / 決済照合プラットフォームと連携し、チャージ、消費、払い戻しなどの財務フローを自動同期
- [ ] **Token 残量予警メカニズム**：アカウント / トークンの Token 残量が少ないときに自動予警、マルチチャネル通知をサポート
- [ ] **ログ監査と監査レポート**：完全な操作監査ログとビジュアル監査レポート、コンプライアンス要件を満たす
- [ ] **AI スマート分析**：大規模モデルに基づき使用量、コスト、チャネル健全性をスマートに分析し提案
- [ ] プラグイン拡張メカニズム
- [ ] エンタープライズ向け SSO / LDAP 連携
- [ ] 使用量アラートと通知チャネルの拡張（DingTalk / 飛書 / 企業微信など）
- [ ] より多くのモデルプラットフォームの継続的な連携

> 💡 PR や Issue の提出を歓迎します。詳しくは [Issues](https://github.com/modelbus/one-api-pro/issues) を参照してください。

---

## License

[MIT License](../LICENSE)
