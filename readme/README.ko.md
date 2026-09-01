<p align="center">
  <img src="../docs/logo.png" width="150" height="150" alt="one-api-pro logo">
</p>

<p align="center">
  One Api Pro · Go로 구축된 엔터프라이즈급 AI API Gateway
</p>
<p align="center">
  본 프로젝트는 <a href="https://github.com/songquanpeng/one-api">one-api</a> (by <a href="https://github.com/songquanpeng">JustSong</a>)을 기반으로 심층 재구축했으며, 원저자의 오픈소스 기여에 감사를 표합니다.
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

## 📑 목차

- [🚀 빠른 시작](#-빠른-시작)
- [🔧 기술 스택](#-기술-스택)
  - [Go 백엔드](#go-백엔드)
  - [Vue 3 프론트엔드](#vue-3-프론트엔드)
- [✨ 기능 하이라이트](#-기능-하이라이트)
- [🔥 one-api와의 비교](#-one-api와의-비교)
- [📸 스크린샷](#-스크린샷)
- [⚙️ 구성](#%EF%B8%8F-구성)
  - [🔧 환경 변수](#-환경-변수)
  - [⌨️ 명령줄 인수](#%EF%B8%8F-명령줄-인수)
- [📖 API 문서](#-api-문서)
- [📦 배포](#-배포)
  - [🔨 수동 배포](#-수동-배포)
  - [🏢 다중 호스트 배포](#-다중-호스트-배포)
  - [🌐 클러스터 배포 (분산 다중 활성)](#-클러스터-배포-분산-다중-활성)
- [🗺️ 개발 로드맵](#%EF%B8%8F-개발-로드맵)
- [License](#license)

---

## 🚀 빠른 시작

### 1. 실행 파일 얻기

[GitHub Releases](https://github.com/modelbus/one-api-pro/releases/latest)에서 사전 빌드된 버전을 다운로드하거나 소스에서 빌드합니다:

```bash
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro
```

### 2. (소스 빌드) Vue 3 프론트엔드 빌드

```bash
cd web
sh build.sh        # web/THEMES에 따라 테마를 순서대로 빌드 (기본 default-pro)
cd ..
```

### 3. (소스 빌드) 백엔드 빌드

> 백엔드는 최신 프론트엔드 산출물을 포함하도록 프론트엔드 빌드가 완료된 후에 컴파일해야 합니다.

```bash
go build -ldflags "-s -w" -o one-api-pro
```

### 4. (선택) 멀티플랫폼 원클릭 패키징

루트 디렉토리의 `release.sh` 스크립트를 사용하면 의존성 다운로드, 프론트엔드 빌드, 멀티플랫폼 크로스 컴파일을 한 번에 수행할 수 있습니다:

```bash
./release.sh                          # VERSION 파일을 버전 번호로 사용
./release.sh v0.1.0                   # 버전 번호 지정
./release.sh v0.1.0 --skip-frontend   # 프론트엔드 빌드 건너뛰기 (기존 web/build 재사용)
```

> 사전 요구사항: `go`, `node`, `npm`. 버전 번호는 루트 `VERSION` 파일에서 가져옵니다 (`v` 접두사 유무를 자동으로 지원).

패키징 산출물은 **정적 링크된 순수 실행 파일**이며 (압축 해제 없이 직접 실행), `dist/` 디렉토리에 출력됩니다:

```
dist/one-api-pro-linux-amd64
dist/one-api-pro-linux-arm64
dist/one-api-pro-windows-amd64.exe
dist/one-api-pro-darwin-amd64
dist/one-api-pro-darwin-arm64
```

> 여기서 `linux-*`는 정적 링크 방식으로 CentOS / Ubuntu에서 모두 사용할 수 있습니다. GitHub Releases는 `.github/workflows/release.yml`에 의해 `v*` 태그가 푸시될 때 자동으로 빌드·배포되며, 로컬 `release.sh`의 출력 로직과 동일합니다.

### 5. 실행

```bash
./one-api-pro --port 3000 --log-dir ./logs
```

`http://localhost:3000`에 접속하여 초기 계정 `root / 123456`으로 로그인합니다.

> 자세한 배포 방법은 [📦 배포](#-배포), 인터페이스 문서는 [📖 API 문서](#-api-문서)를 참고하세요.

---

## 🔧 기술 스택

본 프로젝트는 다음 오픈소스 기술로 구축되었으며, 모든 오픈소스 프로젝트 저자에게 감사를 표합니다.

### Go 백엔드

| 기술 | 용도 |
| --- | --- |
| [Gin](https://github.com/gin-gonic/gin) | HTTP 웹 프레임워크 |
| [GORM](https://gorm.io) | ORM 라이브러리, SQLite / MySQL / PostgreSQL 지원 |
| [go-redis/redis](https://github.com/go-redis/redis) | Redis 클라이언트 |
| [golang-jwt/jwt](https://github.com/golang-jwt/jwt) | JWT 인증 |
| [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2) | AWS Bedrock 통합 |
| [Google API Go Client](https://github.com/googleapis/google-api-go-client) | Google Gemini / PaLM2 통합 |
| [pkoukk/tiktoken-go](https://github.com/pkoukk/tiktoken-go) | Token 계산 |
| [gorilla/websocket](https://github.com/gorilla/websocket) | WebSocket 지원 (iFlytek 등 채널) |
| [joho/godotenv](https://github.com/joho/godotenv) | `.env` 구성 파일 파싱 |

### Vue 3 프론트엔드

| 기술 | 용도 |
| --- | --- |
| [Vue 3](https://vuejs.org) | 프론트엔드 프레임워크 (Composition API) |
| [Vite](https://vitejs.dev) | 빌드 도구 |
| [Arco Design Vue](https://arco.design/vue) | UI 컴포넌트 라이브러리 |
| [Pinia](https://pinia.vuejs.org) | 상태 관리 |
| [Vue Router 4](https://router.vuejs.org) | 라우팅 관리 |
| [Axios](https://axios-http.com) | HTTP 클라이언트 |
| [ECharts](https://echarts.apache.org) | 데이터 시각화 차트 |
| [vue-i18n](https://vue-i18n.intlify.dev) | 국제화 |

---

## ✨ 기능 하이라이트

One Api Pro는 **엔터프라이즈급 AI API Gateway**로, Go 언어 + Vue 3로 새롭게 제작되었으며, 기존 one-api의 모든 기능을 유지하면서 아키텍처 수준의 재구축과 엔터프라이즈급 강화를 수행했습니다.

### 🖥️ 시각화 대시보드

새로운 Vue 3 + Arco Design 관리 백엔드로, 데이터 시각화 대시보드를 제공하며 핵심 지표, 사용 추세, 모델 사용량 분포를 한눈에 파악할 수 있습니다.

| 핵심 지표 카드 | 사용 추세 차트 |
|:---:|:---:|
| ![대시보드 홈](../docs/Demo-Index.png) | ![대시보드 홈](../docs/Demo-Index.png) |

### 🔑 정밀한 토큰 관리

다차원 토큰 통제 지원: 사용 가능 모델 화이트리스트, IP 서브넷 제한, 할당량 상한, 만료 시간, 무제한 할당량. 권한 세분화는 단일 모델까지 가능합니다.

| 토큰 관리 |
|:---:|
| ![토큰 관리](../docs/Demo-Token.png) |

### 📦 플랜 구독 체계

완전한 플랜 및 구독 체계 내장: Token 기준 / 요청 기준 과금, 주기별 요율 제한 (시간 / 주 / 월), 모델별 정밀 통제, 추천 플랜 및 가격 구성 지원.

| 플랜 관리 | 구독 관리 |
|:---:|:---:|
| ![플랜 관리](../docs/Demo-Plan.png) | ![구독 관리](../docs/Demo-Subscribe.png) |

### 💳 주문 및 실제 결제

플랜을 주문할 때마다 **완전한 주문 감사 기록** (주문 번호, 사용자, 플랜 스냅샷 JSON, 금액, 결제 수단, 상태, 결제 시간, 채널 거래 번호)이 남습니다. 플랜 / 충전 두 가지 주문 유형을 지원하며, **WeChat Pay Native** (PC QR)와 **Alipay Face-to-Face** (TradePrecreate)를 네이티브로 연동하고, 은행 / 오프라인 / 무료 세 가지 관리자 채널을 미리 구성해 두었습니다. 플랜 업그레이드 차액은 남은 일수 비율로 자동 계산되며, 중첩 모드에서는 신규·기존 플랜이 병행 적용됩니다. 모든 규칙은 「운영 → 플랜 운영」 하위 탭에서 핫스왑으로 전환할 수 있습니다.

| 주문 센터 | 결제 구성 |
|:---:|:---:|
| ![주문 센터](../docs/Demo-Order.png) | ![결제 구성](../docs/Demo-Payment.png) |

### 🌐 분산 다중 활성 클러스터

분산 다중 활성 클러스터 배포를 지원합니다. 각 노드는 독립적인 MySQL + Redis를 사용하며, 애플리케이션 계층 이벤트 동기화를 통해 데이터 상호 신뢰를 실현합니다. 공유 데이터베이스가 필요 없어 전 세계 멀티 리전 근거리 접속을 자연스럽게 지원합니다.

| 클러스터 노드 관리 |
|:---:|
| ![클러스터 노드 관리](../docs/Demo-cluster.png) |

### 🧩 기타 핵심 기능

- **30+ 모델 플랫폼 연동**: OpenAI / Anthropic / Gemini / DeepSeek / Qwen / ERNIE Bot / iFlytek / Zhipu 등 주요 플랫폼 전부 지원, 통합 OpenAI 호환 인터페이스
- **정밀 원가 계산**: Token 기준 또는 호출당 과금, Prompt / Completion / Cached 독립 가격, 그룹 할인 중첩, 주기 사용량 추적
- **채널 로드 밸런싱**: 가중치 기반 랜덤 분배, 자동 장애 전환, 쿨다운 / 비활성 정책, 채널 동시성 및 RPM 요율 제한
- **다단계 권한 체계**: Guest / User / Admin / Root 4단계 권한, 기존 API 권한 취약점 수정, 세분화된 관리자 작업 권한
- **엔터프라이즈급 보안**: 전 구간 HTTPS, Token 인증, 서브넷 IP 제한, 실시간 감사 로그

---

## 🔥 one-api와의 비교

| 비교 항목 | one-api | one-api-pro |
| --- | --- | --- |
| 프로젝트 이름 | one-api | one-api-pro |
| Adaptor 아키텍처 | 집중식 상수 관리 (channeltype/define.go 56줄 iota + url.go 병렬 배열 + helper.go 이중 switch), 공급자 추가 시 4개 프레임워크 파일 수정 필요 | 자체 등록 메커니즘 (registry + register.go), 공급자 추가 시 패키지 생성 + 등록만 하면 되고 프레임워크 코드 수정 불필요 |
| 권한 세분화 | 관리자와 일반 사용자 권한 경계가 모호하며, 누구든 API로 설정을 변경 가능 | 단계별 권한 체계로 API 권한 취약점 수정, 세분화된 관리자 권한 |
| 구독 모드 | 플랜 / 구독 체계 없음 | 완전한 플랜 구독 + 주기 요율 제한 + 모델별 통제 |
| 분산 클러스터 | 독립 클러스터 미지원, 다중 호스트 배포 시 MySQL 공유 필요 | 분산 다중 활성 클러스터 지원, 각 노드 독립적 MySQL + Redis, 애플리케이션 계층 이벤트 동기화로 데이터 상호 신뢰, 공유 데이터베이스 불필요 |
| 디렉터리 구조 | relay/adaptor/ 에 40개 디렉터리 평평하게 배치, 기본 프로토콜과 공급자 혼재, relay/model/ 과 루트 model/ 충돌 | adaptor/openai/, adaptor/anthropic/ 을 기본 프로토콜로 독립 배치, adaptor/provider/ 에 37개 공급자 통합, relay/schema/ 로 명명 충돌 제거 |
| 관리 백엔드 | 3가지 프론트엔드 테마 (default/berry/air), 기본 관리 기능 | Vue 3 + Arco Design 전면 신규 관리 백엔드, 시각화 대시보드 |
| 지속 업데이트 | 원본 프로젝트는 2024년에 업데이트 중단 | 지속적인 유지보수 및 업데이트, 엔터프라이즈급 시나리오에 최적화 |

---

## 📸 스크린샷

### 🖥️ 대시보드
![대시보드 홈](../docs/Demo-Index.png)

### 🔑 토큰 관리
![토큰 관리](../docs/Demo-Token.png)

### 📦 플랜 관리
![플랜 관리](../docs/Demo-Plan.png)

### 🔄 구독 관리
![구독 관리](../docs/Demo-Subscribe.png)

### 🌐 클러스터 노드 관리
![클러스터 노드 관리](../docs/Demo-cluster.png)

---

## ⚙️ 구성

시스템은 기본적으로 바로 사용할 수 있습니다.

환경 변수 또는 명령줄 인수를 통해 구성할 수 있으며, 시작 후 `root` 사용자로 로그인하여 관리 백엔드에서 계속 설정할 수 있습니다.

> **힌트**: 특정 구성 항목의 의미를 모를 경우, 값을 임시로 삭제하면 추가 힌트 텍스트가 표시됩니다.

### 🔧 환경 변수

> One Api Pro는 `.env` 파일에서 환경 변수를 읽을 수 있습니다. `.env.example` 파일을 참조하여 사용 시 이를 `.env`로 이름을 변경하세요. `--env` 인자로 구성 파일 경로를 지정할 수도 있으며 (상대 경로 지원), 자세한 내용은 명령줄 인수 섹션을 참고하세요.

1. `REDIS_CONN_STRING`: 설정 시 Redis를 캐시로 사용합니다.
   + 예시: `REDIS_CONN_STRING=redis://default:redispw@localhost:49153`
   + 데이터베이스 접근 지연이 매우 낮다면 Redis를 활성화할 필요가 없으며, 활성화 시 오히려 데이터 지연 문제가 발생할 수 있습니다.
   + Sentinel 또는 클러스터 모드를 사용해야 하는 경우:
     + 해당 환경 변수를 노드 목록으로 설정해야 합니다. 예: `localhost:49153,localhost:49154,localhost:49155`.
     + 이 외에도 다음 환경 변수를 설정해야 합니다:
       + `REDIS_PASSWORD`: Redis 클러스터 또는 Sentinel 모드에서의 비밀번호 설정.
       + `REDIS_MASTER_NAME`: Redis Sentinel 모드에서 마스터 노드의 이름.
2. `SESSION_SECRET`: 설정 시 고정된 세션 키를 사용하므로, 시스템 재시작 후에도 로그인된 사용자의 쿠키가 유효합니다.
   + 예시: `SESSION_SECRET=random_string`
3. `SQL_DSN`: 설정 시 SQLite 대신 지정된 데이터베이스를 사용하며, MySQL 또는 PostgreSQL을 사용하세요.
   + 예시:
     + MySQL: `SQL_DSN=root:123456@tcp(localhost:3306)/oneapi`
     + PostgreSQL: `SQL_DSN=postgres://postgres:123456@localhost:5432/oneapi` (적응 중, 피드백 환영)
   + 데이터베이스 `oneapi`를 미리 생성해야 하며, 테이블 생성은 필요하지 않습니다. 프로그램이 자동으로 생성합니다.
   + 클라우드 데이터베이스 사용 시: 클라우드 서버가 신원 인증을 요구한다면 연결 매개변수에 `?tls=skip-verify`를 추가해야 합니다.
   + 데이터베이스 구성에 따라 다음 매개변수를 수정하세요 (또는 기본값 유지):
     + `SQL_MAX_IDLE_CONNS`: 최대 유휴 연결 수, 기본값 `100`.
     + `SQL_MAX_OPEN_CONNS`: 최대 열린 연결 수, 기본값 `1000`.
       + `Error 1040: Too many connections` 오류가 발생하면 이 값을 적절히 줄이세요.
     + `SQL_CONN_MAX_LIFETIME`: 연결의 최대 수명, 기본값 `60`, 단위는 분.
4. `LOG_SQL_DSN`: 설정 시 `logs` 테이블에 독립적인 데이터베이스를 사용하며, MySQL 또는 PostgreSQL을 사용하세요.
5. `FRONTEND_BASE_URL`: 설정 시 페이지 요청을 지정된 주소로 리디렉션하며, 서버에서만 설정할 수 있습니다.
   + 예시: `FRONTEND_BASE_URL=https://openai.justsong.cn`
6. `MEMORY_CACHE_ENABLED`: 메모리 캐시를 활성화하며, 사용자 할당량 업데이트에 일정한 지연이 발생할 수 있습니다. 값은 `true`와 `false`, 미설정 시 기본값 `false`.
   + 예시: `MEMORY_CACHE_ENABLED=true`
7. `SYNC_FREQUENCY`: 캐시 활성화 시 데이터베이스와 구성을 동기화하는 빈도이며, 단위는 초, 기본값 `600`초.
   + 예시: `SYNC_FREQUENCY=60`
8. `NODE_TYPE`: 설정 시 노드 유형을 지정하며, 값은 `master`와 `slave`, 미설정 시 기본값 `master`.
   + 예시: `NODE_TYPE=slave`
9. `CHANNEL_UPDATE_FREQUENCY`: 설정 시 채널 잔액을 주기적으로 갱신하며, 단위는 분, 미설정 시 갱신하지 않습니다.
   + 예시: `CHANNEL_UPDATE_FREQUENCY=1440`
10. `CHANNEL_TEST_FREQUENCY`: 설정 시 채널을 주기적으로 점검하며, 단위는 분, 미설정 시 점검하지 않습니다.
    + 예시: `CHANNEL_TEST_FREQUENCY=1440`
11. `POLLING_INTERVAL`: 채널 잔액 일괄 갱신 및 가용성 테스트 시 요청 간격, 단위는 초, 기본값은 간격 없음.
    + 예시: `POLLING_INTERVAL=5`
12. `BATCH_UPDATE_ENABLED`: 데이터베이스 일괄 갱신 집계를 활성화하며, 사용자 할당량 업데이트에 일정한 지연이 발생할 수 있습니다. 값은 `true`와 `false`, 미설정 시 기본값 `false`.
    + 예시: `BATCH_UPDATE_ENABLED=true`
    + 데이터베이스 연결 수가 과도한 문제가 발생한다면 이 옵션을 활성화해 보세요.
13. `BATCH_UPDATE_INTERVAL=5`: 일괄 갱신 집계의 시간 간격, 단위는 초, 기본값 `5`.
    + 예시: `BATCH_UPDATE_INTERVAL=5`
14. 요청 속도 제한:
    + `GLOBAL_API_RATE_LIMIT`: 전역 API 속도 제한 (중계 요청 제외), 단일 ip 3분 내 최대 요청 수, 기본값 `180`.
    + `GLOBAL_WEB_RATE_LIMIT`: 전역 Web 속도 제한, 단일 ip 3분 내 최대 요청 수, 기본값 `60`.
15. 인코더 캐시 설정:
    + `TIKTOKEN_CACHE_DIR`: 프로그램 시작 시 일반 모델의 토큰 인코딩 (예: `gpt-3.5-turbo`, `gpt-4`, `gpt-4o`)을 온라인으로 다운로드합니다. 네트워크 제한 또는 오프라인 상태라면 다운로드 타임아웃 (약 30초) 후 자동으로 근사 token 계산 (약 `0.38 × 문자 수`)으로 폴백하며 서비스는 정상 시작됩니다. 정확한 과금이 필요하다면 온라인 환경에서 인코딩 파일을 이 디렉토리에 미리 다운로드한 뒤 오프라인 환경으로 이전하세요.
    + `DATA_GYM_CACHE_DIR`: 현재 이 구성은 `TIKTOKEN_CACHE_DIR`과 동일한 역할을 하며, 우선순위는 더 낮습니다.
16. `RELAY_TIMEOUT`: 중계 타임아웃 설정, 단위는 초, 기본값은 타임아웃 시간 미설정.
17. `RELAY_PROXY`: 설정 시 이 프록시를 사용하여 API를 요청합니다.
18. `USER_CONTENT_REQUEST_TIMEOUT`: 사용자가 업로드한 콘텐츠 다운로드 타임아웃 시간, 단위는 초.
19. `USER_CONTENT_REQUEST_PROXY`: 설정 시 이 프록시를 사용하여 사용자가 업로드한 콘텐츠 (예: 이미지)를 요청합니다.
20. `SQLITE_BUSY_TIMEOUT`: SQLite 잠금 대기 타임아웃 설정, 단위는 밀리초, 기본값 `3000`.
21. `GEMINI_SAFETY_SETTING`: Gemini 안전 설정, 기본값 `BLOCK_NONE`.
22. `GEMINI_VERSION`: One Api Pro에서 사용하는 Gemini 버전, 기본값 `v1`.
23. `THEME`: 시스템 테마 설정, 기본값 `default-pro` (Vue 3 관리 백엔드), `default` / `berry` / `air` (기존 React 테마)로 전환 가능, 구체적인 값은 [여기](../web/README.md) 참조.
24. `ENABLE_METRIC`: 요청 성공률에 따라 채널을 비활성화할지 여부, 기본값은 비활성화, 값은 `true`와 `false`.
25. `METRIC_QUEUE_SIZE`: 요청 성공률 통계 큐 크기, 기본값 `10`.
26. `METRIC_SUCCESS_RATE_THRESHOLD`: 요청 성공률 임계값, 기본값 `0.8`.
27. `INITIAL_ROOT_TOKEN`: 이 값이 설정되면 시스템 최초 시작 시 해당 환경 변수 값을 가진 root 사용자 토큰이 자동 생성됩니다.
28. `INITIAL_ROOT_ACCESS_TOKEN`: 이 값이 설정되면 시스템 최초 시작 시 해당 환경 변수 값을 가진 root 사용자 시스템 관리 토큰이 자동 생성됩니다.
29. `ENFORCE_INCLUDE_USAGE`: stream 모델 하에서 usage 반환을 강제할지 여부, 기본값은 비활성화, 값은 `true`와 `false`.
30. `TEST_PROMPT`: 모델 테스트 시 사용자 prompt, 기본값은 `Print your model name exactly and do not output without any other text.`.

#### 🌐 클러스터 구성 (분산 다중 활성 배포)

> 아래 환경 변수를 구성하지 않으면 시스템은 단일 노드 모드로 실행되며 부작용이 없습니다.

1. `CLUSTER_ENABLED`: 클러스터 모드 활성화 여부, 기본값 비활성화.
   + 예시: `CLUSTER_ENABLED=true`
2. `CLUSTER_NODE_ID`: 노드 번호 (1-49), MySQL의 `auto_increment_offset`과 일치해야 하며, 각 노드가 중복될 수 없습니다.
   + 예시: `CLUSTER_NODE_ID=1`
3. `CLUSTER_NODE_NAME`: 노드 이름, 식별 편의를 위한 것으로 기본값은 `node-{NODE_ID}`.
   + 예시: `CLUSTER_NODE_NAME=node-cn`
4. `CLUSTER_NODE_ADDRESS`: 이 노드의 공인 접근 주소 (프로토콜 접두사 포함), 다른 노드가 이 주소로 데이터를 푸시합니다.
   + 예시: `CLUSTER_NODE_ADDRESS=https://cn.example.com`
5. `CLUSTER_SECRET`: 이 노드의 초기 secret, **각 노드 독립적**. 최초 시작 시 초기 secret으로 데이터베이스에 기록되며 이후 admin이 수정할 수 있습니다.
   + 예시: `CLUSTER_SECRET=MyClusterSecret123abc`
6. `CLUSTER_SEEDS`: 시드 노드 주소 (쉼표로 구분), 새 노드 시작 시 시드 노드에 등록하여 클러스터 정보를 획득하며, 도달 가능한 노드 하나만 구성하면 됩니다. 첫 번째 노드는 비우거나 자신의 주소를 구성할 수 있습니다.
   + 예시: `CLUSTER_SEEDS=https://cn.example.com`
   + 다중 시드: `CLUSTER_SEEDS=https://cn.example.com,https://us.example.com`
7. `CLUSTER_PUSH_INTERVAL`: 동기화 이벤트 푸시 간격, 단위는 초, 기본값 `3`.
8. `CLUSTER_DISCOVERY_INTERVAL`: 노드 발견 간격, 단위는 초, 생존 노드는 매 주기 서로 ping, 기본값 `30`.
9. `CLUSTER_DEAD_PING_INTERVAL`: 실패 노드 ping 간격, 단위는 초, 생존 간격보다 길게 설정해 무효 요청을 줄이며, 기본값 `120`.
10. `CLUSTER_MAX_PING_FAILURES`: 연속 ping 실패 횟수, 도달 시 노드를 실패 상태로 표시, 기본값 `3`.
11. `CLUSTER_SYNC_LOGS`: 로그 테이블 동기화 여부, 로그 데이터량이 크므로 필요에 따라 끌 수 있으며, 기본값 `true`.
    + 예시: `CLUSTER_SYNC_LOGS=false`
12. `CLUSTER_BATCH_SIZE`: 푸시당 최대 이벤트 수, 기본값 `50`.

### ⌨️ 명령줄 인수

1. `--port <port_number>`: 서버가 수신할 포트 번호를 지정, 기본값 `3000`.
   + 예시: `--port 3000`
2. `--log-dir <log_dir>`: 로그 폴더를 지정, 미설정 시 작업 디렉토리의 `logs` 폴더에 저장.
   + 예시: `--log-dir ./logs`
3. `--env <env_file_path>`: 구성 파일 경로를 지정, 상대 경로와 절대 경로를 지원. 미지정 시 현재 디렉토리의 `.env` 파일을 자동 로드.
   + 예시: `--env ./config.env`
   + 예시: `--env /etc/one-api-pro/production.env`
   + 다중 인스턴스 배포 예시:
     ```bash
     ./one-api-pro --env ./instances/instance1.env --port 3001 &
     ./one-api-pro --env ./instances/instance2.env --port 3002 &
     ```
   + 구성 우선순위: 명령줄 인수 > 시스템 환경 변수 > `--env` 지정 구성 파일 > 기본값
4. `--version`: 시스템 버전 번호를 출력하고 종료.
   + 예시: `./one-api-pro --version`
   + 버전 번호 출처 (우선순위 높은 순):
     1. 현재 작업 디렉토리 또는 실행 파일과 같은 디렉토리의 `VERSION` 파일 (`v` 접두사 유무를 자동 지원, 예: `0.0.2` 또는 `v0.0.2`);
     2. 컴파일 시 `-ldflags "-X .../common.Version=..."`로 주입된 버전 번호 (`release.sh`와 CI 모두 자동 주입);
     3. 소스 내 기본값 `common/constants.go`.
   + 따라서 루트 디렉토리의 `VERSION` 파일 한 곳만 관리하면 `--version`, 시작 로그, `/api/status` 인터페이스와 프론트엔드 대시보드에 표시되는 버전 번호를 일치시킬 수 있습니다.
5. `--help`: 명령 사용법 및 인수 설명을 확인.
   + 예시: `./one-api-pro --help`

---

## 📖 API 문서

완전한 인터페이스 문서는 [docs/API.md](../docs/API.md)에 독립적으로 관리됩니다:

- **인증 메커니즘**: Cookie Session / Access Token / API Key (Bearer Token) 세 가지 인증 방식
- **관리 인터페이스**: 모델 가격, 그룹 할인, 채널, 토큰, 사용자, 로그, 교환 코드, 플랜, 구독 등 완전한 CRUD
- **OpenAI 호환 인터페이스**: `/v1/models`, `/v1/chat/completions`, `/v1/embeddings`, 이미지, 오디오, 콘텐츠 심사 등
- **클러스터 관리 API**: 노드 발견, 하트비트, 데이터 동기화 등 분산 클러스터 인터페이스

👉 [전체 인터페이스 문서 보기 →](../docs/API.md)

---

## 📦 배포

### 🔨 수동 배포

#### 1. 실행 파일 얻기

다음 방법 중 하나를 선택하세요:

**방법 1: 사전 빌드된 버전 다운로드 (권장)**

[GitHub Releases](https://github.com/modelbus/one-api-pro/releases/latest)에서 해당 플랫폼의 순수 실행 파일 (Linux / macOS / Windows)을 다운로드하며, 압축 해제 없이 바로 실행할 수 있습니다.

**방법 2: release.sh로 원클릭 패키징**

```shell
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro
./release.sh            # 멀티플랫폼 패키징, 산출물은 dist/ 디렉토리로 출력
```

**방법 3: 소스에서 컴파일**

```shell
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro

# 프론트엔드 빌드 (Vue 3 관리 백엔드, web/THEMES에 따라 순서대로 빌드)
cd web
sh build.sh

# 백엔드 빌드 (주의: 최신 프론트엔드 산출물을 포함하도록 반드시 프론트엔드 빌드 이후에 실행)
cd ..
go build -ldflags "-s -w" -o one-api-pro
```

#### 2. 실행

```shell
chmod u+x one-api-pro
./one-api-pro --port 3000 --log-dir ./logs
```

#### 3. 접속

[http://localhost:3000/](http://localhost:3000/)에 접속하여 로그인합니다. 초기 계정 사용자명은 `root`, 비밀번호는 `123456`입니다.

### 🏢 다중 호스트 배포
1. 모든 서버의 `SESSION_SECRET`을 동일한 값으로 설정합니다.
2. 반드시 `SQL_DSN`을 설정하고 SQLite가 아닌 MySQL 데이터베이스를 사용하며, 모든 서버가 동일한 데이터베이스에 연결합니다.
3. 모든 종속 서버는 `NODE_TYPE`을 `slave`로 설정해야 하며, 미설정 시 기본적으로 마스터 서버가 됩니다.
4. `SYNC_FREQUENCY`를 설정하면 서버가 주기적으로 데이터베이스에서 구성을 동기화합니다. 원격 데이터베이스 사용 시 마스터/슬레이브 여부와 관계없이 이 항목을 설정하고 Redis를 활성화하는 것이 좋습니다.
5. 종속 서버는 선택적으로 `FRONTEND_BASE_URL`을 설정하여 페이지 요청을 마스터 서버로 리디렉션할 수 있습니다.
6. 종속 서버에 **각각** Redis를 설치하고 `REDIS_CONN_STRING`을 설정하면, 캐시가 만료되지 않는 동안 데이터베이스 접근을 0회로 줄여 지연을 감소시킬 수 있습니다 (Redis 클러스터 또는 Sentinel 모드 지원은 환경 변수 설명 참조).
7. 마스터 서버의 데이터베이스 접근 지연이 높은 경우에도 Redis를 활성화하고 `SYNC_FREQUENCY`를 설정하여 주기적으로 데이터베이스에서 구성을 동기화해야 합니다.

환경 변수의 구체적인 사용 방법은 [여기](#환경-변수)를 참조하세요.

### 🌐 클러스터 배포 (분산 다중 활성)

클러스터 모드는 여러 노드가 각자 독립적인 One Api Pro + MySQL을 배포하고, 애플리케이션 계층 이벤트 동기화를 통해 데이터 상호 신뢰를 실현하므로 공유 데이터베이스가 필요 없습니다.

> **적용 시나리오**: 전 세계 멀티 리전 배포, 근거리 접속으로 지연 감소, 고가용성 재해 복구, 멀티 노드 로드 밸런싱.

#### 🗺️ 아키텍처 개요

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

#### ⭐ 핵심 특징

- **분산**: 모든 노드가 동등하며 마스터/슬레이브 구분이 없고, 어떤 노드든 데이터 변경 후 모든 생존 노드에 능동적으로 푸시
- **제로 침투**: GORM 콜백으로 데이터 변경을 포착하며, 기존 비즈니스 코드를 수정하지 않음
- **비동기 푸시**: 데이터 동기화가 메인 흐름을 차단하지 않고, 백그라운드 goroutine이 일괄 푸시
- **충돌 해결**: `updated_at` 타임스탬프 비교 기반, 더 새로운 데이터만 기록
- **속도 제한 동기화**: 채널 동시성 및 RPM 요율 제한 카운터를 데이터베이스 테이블로 크로스 노드 동기화
- **단일 노드 호환**: 클러스터 환경 변수를 구성하지 않으면 시스템은 완전히 단일 노드 모드로 실행

#### 📊 동기화 범위

| 데이터 테이블 | 동기화 여부 | 설명 |
|--------|---------|------|
| users | ✅ | 사용자 정보 |
| tokens | ✅ | API 토큰 |
| channels | ✅ | 채널 구성 |
| abilities | ✅ | 채널 능력 |
| options | ✅ | 시스템 설정 |
| redemptions | ✅ | 교환 코드 |
| plans | ✅ | 구독 플랜 |
| user_plans | ✅ | 사용자 구독 |
| plan_usages | ✅ | 플랜 사용량 |
| channel_counters | ✅ | 채널 속도 제한 카운터 |
| cluster_nodes | 🔄 Discovery | 클러스터 노드 정보 (발견 메커니즘이 유지, 데이터 동기화 없음) |
| logs | ⚠️ 선택 | 로그 데이터량이 크며, `CLUSTER_SYNC_LOGS`로 제어 |

#### 🚀 배포 단계

**1. MySQL 구성 (각 노드는 반드시 독립적인 MySQL 인스턴스 사용)**

각 노드마다 **독립적인 MySQL 인스턴스**가 필요합니다 (`auto_increment_offset`은 인스턴스 수준 변수이므로, 하나의 MySQL 인스턴스에 여러 데이터베이스를 만들어 여러 노드를 배포할 수 없습니다).

```ini
# 노드 1의 my.cnf
[mysqld]
server-id = 1
auto_increment_increment = 50
auto_increment_offset = 1
log_bin = mysql-bin
binlog_format = ROW

# 노드 2의 my.cnf
[mysqld]
server-id = 2
auto_increment_increment = 50
auto_increment_offset = 2
log_bin = mysql-bin
binlog_format = ROW

# 노드 3의 my.cnf
[mysqld]
server-id = 3
auto_increment_increment = 50
auto_increment_offset = 3
log_bin = mysql-bin
binlog_format = ROW
```

> `auto_increment_increment`를 50으로 설정하면 최대 50개 노드를 지원합니다. 각 노드의 `offset`은 `CLUSTER_NODE_ID`와 일치해야 하며 서로 달라야 합니다.

> **중요:** `auto_increment_increment`와 `auto_increment_offset`은 MySQL의 **시스템 수준 변수**로, 인스턴스 내 모든 데이터베이스에 적용됩니다. 데이터베이스별로 다른 값을 설정할 수 없고 테이블 수준에서도 설정할 수 없습니다 (MySQL 테이블 옵션은 `AUTO_INCREMENT` 시작 값만 지원하며 증분은 지원하지 않습니다). 따라서 각 노드**는 반드시 독립적인 MySQL 인스턴스**를 사용해야 하며, 하나의 MySQL 인스턴스 내에서 서로 다른 데이터베이스를 만들어 여러 노드를 배포할 수 없습니다. 한 머신에서 여러 MySQL 인스턴스를 실행해야 한다면, 다른 포트로 여러 mysqld 프로세스를 시작하거나 Docker로 여러 독립적인 MySQL 컨테이너를 실행할 수 있습니다.

> **`server-id` 및 binlog 정보:** `server-id`는 동일 클러스터 내 모든 MySQL 인스턴스에서 서로 달라야 합니다. `log_bin`과 `binlog_format=ROW`는 적극 권장합니다 — 향후 마스터-슬레이브 복제 확장과 point-in-time recovery에 사용됩니다. 클러스터 데이터 동기화 자체는 binlog에 의존하지 않으며 (GORM 콜백을 통해 애플리케이션 계층에서 구현), binlog는 추가적인 신뢰성 보장을 제공합니다.

**2. Redis 구성 (각 노드는 반드시 독립적인 Redis 인스턴스 사용)**

각 노드마다 **독립적인 Redis 인스턴스**도 필요합니다 (포트가 다르거나 다른 머신). Redis는 이 클러스터 아키텍처에서 노드 간 통신에 사용되지 않으며, 이 노드의 캐시, 속도 제한 등 비즈니스 용도로만 사용됩니다.

**3. 새 노드 초기화 데이터**

새 노드가 온라인 상태가 되면 먼저 기존 노드의 데이터 스냅샷을 가져와야 합니다:

```bash
# 방법 1: 기존 노드에서 내보내기 및 가져오기
mysqldump -h existing-node -u root -p oneapi > backup.sql
mysql -u root -p oneapi < backup.sql

# 방법 2: API를 통해 스냅샷 가져오기 (서비스 먼저 시작 필요)
curl -H "X-Cluster-Secret: your-secret" \
  "https://existing-node/api/cluster/snapshot?tables=users,tokens,channels,abilities,options,redemptions,plans,user_plans,plan_usages" \
  -o snapshot.json
```

**4. 환경 변수 구성 (전체 예시)**

다음은 3노드 클러스터의 완전한 `.env` 구성 예시입니다. 각 노드는 독립적인 MySQL 및 Redis 인스턴스를 사용하며 포트와 경로가 서로 다릅니다.

**노드 1 — 중국 노드 (`/opt/one-api-pro/node1/.env`):**
```bash
# ========================
# 기본 구성
# ========================
PORT=3000
SYSTEM_NAME=One Api Pro Cluster

# ========================
# 데이터베이스 (독립 MySQL 인스턴스)
# ========================
SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node1?charset=utf8mb4&parseTime=True&loc=Local

# ========================
# Redis (독립 Redis 인스턴스)
# ========================
REDIS_CONN_STRING=redis://127.0.0.1:6379/0

# ========================
# 클러스터 구성
# ========================
CLUSTER_ENABLED=true
CLUSTER_NODE_ID=1
CLUSTER_NODE_NAME=node-cn
CLUSTER_NODE_ADDRESS=https://cn.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me

# 시드 노드 (최초 시작 시 다른 노드 발견을 안내)
# 첫 번째 노드: 자신의 주소를 입력하거나 비워 둠
# 후속 노드: 임의의 생존 노드 주소를 입력
CLUSTER_SEEDS=https://cn.example.com,https://us.example.com,https://eu.example.com

# ========================
# 클러스터 튜닝 (선택)
# ========================
CLUSTER_DISCOVERY_INTERVAL=30
CLUSTER_DEAD_PING_INTERVAL=120
CLUSTER_MAX_PING_FAILURES=3
CLUSTER_PUSH_INTERVAL=3
CLUSTER_SYNC_LOGS=true
CLUSTER_BATCH_SIZE=50
```

**노드 2 — 미국 노드 (`/opt/one-api-pro/node2/.env`):**
```bash
# 기본 구성
PORT=3001
SYSTEM_NAME=One Api Pro Cluster

# 데이터베이스 (독립 MySQL 인스턴스, 노드 1과 포트 또는 머신 다름)
SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node2?charset=utf8mb4&parseTime=True&loc=Local

# Redis (독립 Redis 인스턴스)
REDIS_CONN_STRING=redis://127.0.0.1:6380/0

# 클러스터 구성
CLUSTER_ENABLED=true
CLUSTER_NODE_ID=2
CLUSTER_NODE_NAME=node-us
CLUSTER_NODE_ADDRESS=https://us.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me   # 노드 1과 완전히 동일해야 함

# 임의의 생존 노드 주소 입력
CLUSTER_SEEDS=https://cn.example.com
```

**노드 3 — 유럽 노드 (`/opt/one-api-pro/node3/.env`):**
```bash
# 기본 구성
PORT=3002
SYSTEM_NAME=One Api Pro Cluster

# 데이터베이스
SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node3?charset=utf8mb4&parseTime=True&loc=Local

# Redis
REDIS_CONN_STRING=redis://127.0.0.1:6381/0

# 클러스터 구성
CLUSTER_ENABLED=true
CLUSTER_NODE_ID=3
CLUSTER_NODE_NAME=node-eu
CLUSTER_NODE_ADDRESS=https://eu.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me   # 모든 노드와 일치해야 함

# 임의의 생존 노드 주소 입력
CLUSTER_SEEDS=https://cn.example.com
```

**구성 매개변수 대조표:**

| 환경 변수 | 노드 1 | 노드 2 | 노드 3 | 설명 |
|---|---|---|---|---|
| `PORT` | 3000 | 3001 | 3002 | 수신 포트 (같은 머신에서는 달라야 함) |
| `SQL_DSN` | ...oneapi_node1 | ...oneapi_node2 | ...oneapi_node3 | 독립 MySQL 인스턴스 |
| `REDIS_CONN_STRING` | :6379/0 | :6380/0 | :6381/0 | 독립 Redis 인스턴스 |
| `CLUSTER_NODE_ID` | 1 | 2 | 3 | 노드 번호, MySQL `auto_increment_offset`에 대응 |
| `CLUSTER_NODE_NAME` | node-cn | node-us | node-eu | 노드 이름, 식별 편의 |
| `CLUSTER_NODE_ADDRESS` | https://cn.example.com | https://us.example.com | https://eu.example.com | 노드 공인 주소 (다른 노드가 이 주소로 접근) |
| `CLUSTER_SECRET` | 동일한 값 | 동일한 값 | 동일한 값 | **모든 노드가 완전히 일치해야 함** |
| `CLUSTER_SEEDS` | 자신의 주소 또는 비움 | 임의 생존 노드 | 임의 생존 노드 | 최초 시작 안내, 이후 자동 발견 |

**5. 시작 명령**

각 노드는 `--env` 인자로 자신의 구성 파일을 로드합니다:

```bash
# 노드 1
./one-api-pro --env /opt/one-api-pro/node1/.env --port 3000

# 노드 2
./one-api-pro --env /opt/one-api-pro/node2/.env --port 3001

# 노드 3
./one-api-pro --env /opt/one-api-pro/node3/.env --port 3002
```

**6. 시작 순서**

1. 첫 번째 노드(Node A)를 시작하고, `CLUSTER_SEEDS`는 비우거나 자신의 주소를 입력합니다
2. Node A가 완전히 시작될 때까지 대기합니다 (약 5-10초, "클러스터 모듈 초기화 완료" 로그 확인)
3. 후속 노드를 시작하고, `CLUSTER_SEEDS`에 임의의 생존 노드 주소를 입력합니다
4. 후속 노드 시작 후 자동으로 시드 노드를 ping하여 전이적으로 모든 다른 노드를 발견합니다
5. 모든 노드가 시작되면 임의 노드의 관리 백엔드 "설정 → 노드 관리" 페이지에서 노드 상태를 확인합니다

**7. Nginx 로드 밸런싱 구성 예시 (선택)**

```nginx
upstream one_api_cluster {
    ip_hash;  # IP 해시 기반, 동일 사용자 요청을 같은 노드에 고정해 session/cache 적중 보장
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

> **`ip_hash` 사용이 핵심입니다**: 동일 사용자의 요청이 항상 같은 노드로 가도록 보장하여, 플랜 요율 제한, Redis 캐시 등 상태가 서로 다른 노드 사이에서 유실되지 않게 합니다.

**8. 클러스터 상태 검증**

배포 완료 후 다음 방법으로 검증할 수 있습니다:

```bash
# 노드 목록 보기 (임의 노드에서 호출)
curl -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  https://cn.example.com/api/cluster_node/

# 모든 노드 목록이 status, last_heartbeat, ping_failures 등 필드와 함께 반환되어야 함
```

또는 관리 백엔드에서: **설정 → 노드 관리** 페이지에서 노드 목록, 상태, 마지막 하트비트 시간 등을 확인할 수 있습니다.

> 💡 클러스터 관리 API에 대한 자세한 내용은 [docs/API.md 부록 E: 클러스터 관리 API](../docs/API.md#附录-e集群管理-api)

#### ⚠️ 주의사항

- 각 노드는 반드시 독립적인 MySQL 인스턴스와 Redis 인스턴스를 가져야 하며 데이터베이스를 공유하지 않습니다
- `CLUSTER_SECRET`은 모든 노드 간에 일치해야 하며, 강력한 비밀번호를 사용하고 잘 보관하세요
- `CLUSTER_NODE_ID`는 모든 노드 간에 서로 달라야 하며 MySQL `auto_increment_offset`과 일치해야 합니다
- `CLUSTER_NODE_ADDRESS`는 다른 노드가 접근 가능한 공인 주소여야 합니다 (프로토콜 접두사 `https://` 포함)
- 새 노드 온라인 전 데이터 초기화는 수동으로 완료해야 합니다 (온라인 노드에서 스냅샷 가져오기)
- 로그 테이블(logs)은 데이터량이 크므로 `CLUSTER_SYNC_LOGS=false`로 로그 동기화를 끌 수 있습니다
- MySQL의 `auto_increment_increment`와 `auto_increment_offset`은 `CLUSTER_NODE_ID` 구성과 일치해야 합니다
- 노드 발견은 ping 양방향 등록 메커니즘을 사용하며, 실패 노드는 삭제되지 않고 status=2로만 표시되며 네트워크 복구 후 자동으로 부활합니다
- `CLUSTER_SEEDS`는 최초 시작 시의 안내일 뿐이며, 노드가 ping으로 다른 노드를 발견하면 더 이상 SEEDS에 의존하지 않습니다
- 노드 오프라인 동안 다른 노드에서 발생한 변경사항은 **자동으로 재전송되지 않으며**, 오프라인 노드가 다시 온라인되면 스냅샷을 가져와 데이터를 보충해야 합니다

#### 📝 "로컬 노드" 자기 등록 정보

각 노드가 시작될 때 자신의 `cluster_nodes` 테이블에 로컬 레코드를 하나 기록합니다 (`node_id`는 로컬에 구성된 `CLUSTER_NODE_ID`와 동일). 이는 **의도된 설계**이며 이유는 다음과 같습니다:

1. **관리 백엔드 표시**: "설정 → 노드 관리" 페이지에서 관리자가 로컬 정보(주소, 상태, 하트비트 시간 등)를 확인하여 문제를 해결할 수 있도록
2. **노드 발견 전이성**: 노드 B가 노드 A의 ping 요청을 받을 때, A는 응답에 전체 노드 목록(A 자신 포함)을 반환합니다. B는 이를 수신하여 로컬 테이블에 병합합니다. 이렇게 C는 B의 응답을 통해 A의 존재를 학습할 수 있습니다
3. **생존 판단 근거**: 로컬 레코드의 `last_heartbeat`는 로컬에서 30초마다 자동 갱신되며 (`discoverOnce` 함수), 로컬의 정상 동작 상태를 반영합니다

**자기 등록은 순환 동기화 데이터를 유발하지 않습니다.** 시스템은 5개 계층에서 방어합니다:

| 방어 지점 | 역할 |
|---|---|
| ① `GetAllRemoteNodes` SQL 필터 | 발견 시 SQL에 `node_id != ?`를 추가해 로컬 제외 |
| ② `GetAliveNodesForSync` SQL 필터 | 푸시 시 SQL에 `node_id != ?`를 추가해 로컬 제외 |
| ③ `handlePing` 자기 ping 거부 | `req.NodeId == NodeID`를 명시적으로 거부 |
| ④ `mergeDiscoveredNodes` 로컬 건너뜀 | 발견 노드 병합 시 로컬 건너뜀 |
| ⑤ `ApplyEvents` 로컬 이벤트 건너뜀 | 이벤트 적용 시 로컬 생성 이벤트 건너뜀 |

데이터 흐름은 단방향입니다: 로컬에서 원격으로 푸시, 원격에서 로컬로 가져와 적용, **절대 순환이 없습니다**.

관리 백엔드는 로컬 노드 이름 옆에 "로컬" 파란색 배지를 표시하고, 로컬 노드에 대해 "삭제"와 "수동 Ping"을 비활성화합니다 (이 두 작업은 로컬 노드에 의미가 없습니다).

#### 🔐 "노드별 독립 secret" 정보

각 노드는 **자체 secret**을 가지며, 더 이상 전역 공유 secret을 사용하지 않습니다. 설계 이유:

1. **보안성**: 한 노드의 secret이 유출되어도 다른 노드에 영향을 주지 않습니다
2. **관리 유연성**: 각 노드가 자체 secret을 독립적으로 교체할 수 있습니다
3. **자동 발견**: 노드 간 ping 시 자체 secret을 자동으로 전달하여 상대방이 보관하도록 합니다

**Secret 수명주기**:
- 노드 최초 시작: `CLUSTER_SECRET` 환경 변수를 초기 값으로 사용하여 `cluster_nodes.secret_key` 필드에 기록
- 이후 시작: `cluster_nodes.secret_key`에서 읽기
- Admin이 "노드 관리" 페이지에서 다른 노드의 secret을 수정할 수 있음
- ping 시 `X-Cluster-Secret` 헤더 = **대상 노드**의 secret (로컬 DB에서 조회)

**새 노드 추가 절차**:
1. 노드 A에서 B 노드 레코드를 추가하고 B의 `CLUSTER_SECRET` 값을 입력
2. 노드 B에서 A 노드 레코드를 추가하고 A의 `CLUSTER_SECRET` 값을 입력
3. A가 B를 ping: B의 secret 사용; B가 수신: 자체 secret 검증 ✓
4. B의 응답에 A, B 각각의 secret이 포함되고 A가 로컬 저장소를 갱신

#### 🗑️ "소프트 삭제 노드" 정보

Admin이 노드를 삭제할 때 **물리적으로 삭제하지 않고** `disabled = true`로 설정합니다:

- 삭제된 노드가 "자동으로 다시 등장"하는 것을 방지 (ping 메커니즘이 재등록함)
- 비활성된 노드는 여전히 ping에 응답합니다 (상대방에게 이 노드가 온라인임을 알림) 그러나 이 노드의 정보는 가져오지 않음
- 물리적 삭제는 수동 SQL 필요: `DELETE FROM cluster_nodes WHERE node_id = ?`

#### 🔄 "데이터 동기화 메커니즘" 정보 (중요)

**클러스터 데이터 동기화**는 전적으로 **GORM 이벤트 + HTTP 능동 푸시** 메커니즘에 의존합니다:
- 모든 비즈니스 테이블의 INSERT/UPDATE/DELETE 작업 → GORM 콜백이 포착 → `sync_events` 테이블에 기록 → Pusher goroutine이 모든 생존 노드에 푸시
- 수신측은 `WithSkipHook`으로 로컬 데이터베이스에 기록 (순환 없음)
- 수신측은 `event.NodeId == 로컬 NodeID` 이벤트를 건너뜁니다 (이중 안전장치)

**아키텍처 트레이드오프**: 본 설계는 **크로스 노드 능동 풀(pull)을 구현하지 않습니다**, 이유는 다음과 같습니다:
1. **비즈니스 침투**: 크로스 노드 풀은 각 테이블의 비즈니스 고유 필드를 알아야 하며, 비즈니스 코드를 침해합니다
2. **기본 키 충돌**: 크로스 노드 자동 증가 ID가 연속되지 않으며 (다른 `auto_increment_offset`), 소스 노드 id를 사용하면 offset 설계가 깨집니다
3. **복잡도 높음**: 유지보수 비용이 높고 신뢰성 개선은 제한적입니다
4. **능동 푸시로 충분**: 95%의 시나리오 (노드 온라인 시 일반 동기화)는 푸시로 완전히 해결됩니다

**알려진 제한 및 운영 요구사항**:
- 노드 오프라인 동안 다른 노드에서 발생한 데이터 변경 → **영구 손실** (푸시는 실시간)
- 노드가 다시 온라인되어도 오프라인 기간의 데이터를 자동 보충할 수 없음
- 새 노드가 가입한 이후의 데이터 변경만 수신할 수 있으며, 과거 데이터는 없음
- **운영 해결책**: `mysqldump`로 다른 노드에서 내보낸 후 가져오기

**일반적인 배포 시나리오 대조**:

| 시나리오 | 풀 필요 여부 | 처리 방식 |
|---|---|---|
| 노드 상시 온라인 | ❌ | 푸시로 충분 |
| 노드 간헐적 재시작 (분 단위) | ⚠️ | 짧은 오프라인 데이터 유실, 운영상 허용 가능 |
| 노드 빈번한 유지보수 | ❌ | 푸시 계속, 재시작 후 즉시 복구 |
| 새 노드 클러스터 가입 | ❌ | DBA가 수동으로 `mysqldump` 초기화 |
| 노드 장기 오프라인 후 복구 | ❌ | DBA가 수동으로 `mysqldump` 보충 |

배포 후 접속 시 빈 화면이 나타나면 [#97](https://github.com/modelbus/one-api-pro/issues/97)을 참조하세요.

---

## 🗺️ 개발 로드맵

### ✅ 완료

- [x] **아키텍처 수준 재구축**: Adaptor 자체 등록 메커니즘, 공급자 추가 시 프레임워크 수정 없음
- [x] **Vue 3 전면 신규 관리 백엔드**: Arco Design + 시각화 대시보드 + 30+ 모델 플랫폼 아이콘
- [x] **플랜 구독 체계**: Token / 호출당 과금, 주기 요율 제한, 모델별 정밀 통제
- [x] **분산 다중 활성 클러스터**: GORM 이벤트 기반 + HTTP 능동 푸시 동기화, 공유 데이터베이스 불필요
- [x] **정밀 원가 계산**: Prompt / Completion / Cached 독립 가격, 그룹 할인 중첩
- [x] **다단계 권한 체계**: Guest / User / Admin / Root 4단계, 기존 API 권한 취약점 수정
- [x] **OpenAI 호환 인터페이스**: models / chat / completions / embeddings / images / audio / moderations 완전 지원
- [x] **플랜 주문 및 업그레이드 절차**: 네이티브 `POST /api/order/plan`으로 플랜 구독 주문 생성, `stack`(중첩)과 `price_diff`(차액 업그레이드) 두 가지 모드 지원, 차액은 남은 일수 비율로 자동 계산, 동일 등급 및 다운그레이드 검증 포함
- [x] **주문 감사 및 주문 센터**: 신규 `orders` 테이블 (type/source/order_no/plan_info/amount/status/pay_status/pay_method/pay_time/pay_trade_no)로 모든 결제/관리 승인 흐름을 영속화, 프론트엔드 사용자측 `/plans` 및 `/orders` 페이지에 완전히 표시
- [x] **실제 결제 통합 (gopay)**: WeChat Pay Native (PC QR)와 Alipay Face-to-Face (TradePrecreate) 네이티브 연동, 결제 콜백은 `/api/payment/{wechat,alipay}/notify`로 서명 검증 + 주문 활성화 폐루프 완성
- [x] **결제 / 플랜 운영 설정**: 「운영 설정」 아래에 「플랜 운영」(차액 업그레이드 vs 중첩) 및 「결제」(위챗 / 알리페이 / 은행 3채널 독립 스위치 + 인증서 업로드 + 알림 URL 구성), 필요에 따라 폼 표시

### 🔄 진행 중

- [ ] **더 풍부한 채널 진단 및 지능형 라우팅 최적화**: 자동 쿨다운(`CooldownFilter`), 폴백 다운그레이드(`FallbackFilter`) 및 저성공률 자동 비활성화(`monitor`)는 이미 준비됨, 다음 단계로 독립 진단 패널 / 노드 수준 ping 및 사람 검토 절차 완성
- [ ] 더 풍부한 사용량 분석 리포트 및 내보내기
- [ ] 다국어 국제화(i18n) 완성

### 🔭 계획 중

- [ ] **결제 채널 확장**: Apple Pay, UnionPay, Stripe 등; 비동기 환불 API + 자동화 환불 흐름
- [ ] **잔액(quota) 온라인 충전**: 사용자가 「개인」 영역에서 계정 할당량을 직접 충전 가능, 구독 플랜과 필요에 따라 서로 간섭 없음
- [ ] **일반 플랫폼 재무 연동**: 주요 재무 / 대조 플랫폼 연동, 충전·소비·환불 등 재무 흐름 자동 동기화
- [ ] **Token 잔량 경고 메커니즘**: 계정 / 토큰 Token 잔량이 낮을 때 자동 경고, 다중 채널 알림 지원
- [ ] **로그 감사 및 감사 리포트**: 완전한 작업 감사 로그와 시각화 감사 리포트, 규정 요구 충족
- [ ] **AI 지능 분석**: 대규모 모델 기반으로 사용량, 비용, 채널 건강도에 대한 지능형 분석 및 제안
- [ ] 플러그인화 확장 메커니즘
- [ ] 엔터프라이즈급 SSO / LDAP 연동
- [ ] 사용량 경보 및 알림 채널 확장 (DingTalk / Feishu / WeCom 등)
- [ ] 더 많은 모델 플랫폼 지속 연동

> 💡 PR 또는 Issue 제출을 환영합니다. 자세한 내용은 [Issues](https://github.com/modelbus/one-api-pro/issues)를 참조하세요.

---

## License

본 프로젝트는 [MIT License](../LICENSE)에 따라 배포됩니다.
