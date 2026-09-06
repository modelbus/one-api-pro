<p align="center">
  <img src="../docs/logo.png" width="150" height="150" alt="one-api-pro logo">
</p>

<p align="center">
  One Api Pro · Ein unternehmensweiter AI-API-Gateway auf Basis von Go
</p>
<p align="center">
  Dieses Projekt ist eine umfassende Weiterentwicklung von <a href="https://github.com/songquanpeng/one-api">one-api</a> (von <a href="https://github.com/songquanpeng">JustSong</a>) – vielen Dank an den ursprünglichen Autor für seinen Open-Source-Beitrag.
</p>

<p align="center">
  👉 <strong>Live-Demo ansehen</strong>: <a href="http://demo.one-api.pro">http://demo.one-api.pro</a>
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

## 📑 Inhaltsverzeichnis

- [🚀 Schnellstart](#-schnellstart)
- [🔧 Technologie-Stack](#-technologie-stack)
  - [Go-Backend](#go-backend)
  - [Vue-3-Frontend](#vue-3-frontend)
- [✨ Funktions-Highlights](#-funktions-highlights)
- [🔥 Vergleich mit one-api](#-vergleich-mit-one-api)
- [📸 Screenshots](#-screenshots)
- [⚙️ Konfiguration](#%EF%B8%8F-konfiguration)
  - [🔧 Umgebungsvariablen](#-umgebungsvariablen)
  - [⌨️ Befehlszeilenargumente](#%EF%B8%8F-befehlszeilenargumente)
- [📖 API-Dokumentation](#-api-dokumentation)
- [📦 Bereitstellung](#-bereitstellung)
  - [🔨 Manuelle Bereitstellung](#-manuelle-bereitstellung)
  - [🏢 Multi-Host-Bereitstellung](#-multi-host-bereitstellung)
  - [🌐 Dezentrale Cluster-Bereitstellung](#-dezentrale-cluster-bereitstellung)
- [🗺️ Entwicklungsplan](#%EF%B8%8F-entwicklungsplan)
- [Lizenz](#lizenz)

---

## 🚀 Schnellstart

### 1. Die ausführbare Datei beschaffen

Lade die vorkompilierte Version von den [GitHub Releases](https://github.com/modelbus/one-api-pro/releases/latest) herunter oder kompiliere sie aus dem Quellcode:

```bash
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro
```

### 2. (Quellcode-Build) Vue-3-Frontend erstellen

```bash
cd web
sh build.sh        # erstellt jede Theme aus web/THEMES (Standard: default-pro)
cd ..
```

### 3. (Quellcode-Build) Backend erstellen

> Das Backend muss nach dem Frontend kompiliert werden, damit die neuesten Frontend-Artefakte eingebettet werden.

```bash
go build -ldflags "-s -w" -o one-api-pro
```

### 4. (Optional) One-Click-Paketierung für mehrere Plattformen

Verwende das Skript `release.sh` im Projektstamm, um in einem Schritt die Abhängigkeiten herunterzuladen, das Frontend zu erstellen und die Kreuzkompilierung für alle Plattformen durchzuführen:

```bash
./release.sh                          # verwendet die VERSION-Datei als Versionsnummer
./release.sh v0.1.0                   # gibt eine Versionsnummer an
./release.sh v0.1.0 --skip-frontend   # überspringt den Frontend-Build (nutzt vorhandenes web/build)
```

> Voraussetzungen: `go`, `node`, `npm`. Die Versionsnummer stammt aus der Datei `VERSION` im Projektstamm (kompatibel mit oder ohne `v`-Präfix).

Die Ausgabe sind **statisch verknüpfte, eigenständige ausführbare Dateien** (kein Entpacken nötig, direkt ausführbar), die in das Verzeichnis `dist/` geschrieben werden:

```
dist/one-api-pro-linux-amd64
dist/one-api-pro-linux-arm64
dist/one-api-pro-windows-amd64.exe
dist/one-api-pro-darwin-amd64
dist/one-api-pro-darwin-arm64
```

> Die `linux-*`-Binaries sind statisch verknüpft und funktionieren auf CentOS / Ubuntu. GitHub Releases werden von `.github/workflows/release.yml` automatisch erstellt und veröffentlicht, wenn ein `v*`-Tag gepusht wird – identisch zur lokalen Ausgabe von `release.sh`.

### 5. Starten

```bash
./one-api-pro --port 3000 --log-dir ./logs
```

Rufe `http://localhost:3000` auf und melde dich mit dem Standardkonto `root / 123456` an.

> Details zur Bereitstellung findest du unter [📦 Bereitstellung](#-bereitstellung), die API-Dokumentation unter [📖 API-Dokumentation](#-api-dokumentation).

---

## 🔧 Technologie-Stack

Dieses Projekt wurde mit den folgenden Open-Source-Technologien aufgebaut – allen Autoren der Open-Source-Projekte sei gedankt.

### Go-Backend

| Technologie | Verwendung |
| --- | --- |
| [Gin](https://github.com/gin-gonic/gin) | HTTP-Web-Framework |
| [GORM](https://gorm.io) | ORM-Bibliothek, unterstützt SQLite / MySQL / PostgreSQL |
| [go-redis/redis](https://github.com/go-redis/redis) | Redis-Client |
| [golang-jwt/jwt](https://github.com/golang-jwt/jwt) | JWT-Authentifizierung |
| [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2) | AWS-Bedrock-Integration |
| [Google API Go Client](https://github.com/googleapis/google-api-go-client) | Google-Gemini / PaLM2-Integration |
| [pkoukk/tiktoken-go](https://github.com/pkoukk/tiktoken-go) | Token-Zählung |
| [gorilla/websocket](https://github.com/gorilla/websocket) | WebSocket-Unterstützung (z. B. iFlytek) |
| [joho/godotenv](https://github.com/joho/godotenv) | Analyse der `.env`-Konfigurationsdatei |

### Vue-3-Frontend

| Technologie | Verwendung |
| --- | --- |
| [Vue 3](https://vuejs.org) | Frontend-Framework (Composition API) |
| [Vite](https://vitejs.dev) | Build-Tool |
| [Arco Design Vue](https://arco.design/vue) | UI-Komponentenbibliothek |
| [Pinia](https://pinia.vuejs.org) | Zustandsverwaltung |
| [Vue Router 4](https://router.vuejs.org) | Routing |
| [Axios](https://axios-http.com) | HTTP-Client |
| [ECharts](https://echarts.apache.org) | Datenvisualisierung (Diagramme) |
| [vue-i18n](https://vue-i18n.intlify.dev) | Internationalisierung |

---

## ✨ Funktions-Highlights

One Api Pro ist ein **unternehmensweiter KI-API-Gateway**, neu aufgebaut mit Go + Vue 3. Es behält alle Funktionen des Originals und bietet zugleich eine architektonische Neustrukturierung sowie unternehmensweite Verbesserungen.

### 🖥️ Visualisiertes Dashboard

Ein brandneues Admin-Interface aus Vue 3 + Arco Design bietet ein visualisiertes Dashboard, auf dem Kennzahlen, Nutzungstrends und die Verteilung der Modellnutzung auf einen Blick erkennbar sind.

| Kernkennzahl-Karten | Nutzungstrend-Diagramm |
|:---:|:---:|
| ![Dashboard](../docs/Demo-Index.png) | ![Dashboard](../docs/Demo-Index.png) |

### 🔑 Feingranulare Token-Verwaltung

Unterstützt eine mehrdimensionale Token-Kontrolle: Modell-Whitelist, IP-Subnetz-Beschränkungen, Kontingentobergrenzen, Ablaufzeit und unbegrenztes Kontingent – die Rechtegranularität reicht bis hinunter zu einem einzelnen Modell.

| Token-Verwaltung |
|:---:|
| ![Token-Verwaltung](../docs/Demo-Token.png) |

### 📦 Paket- und Abonnementsystem

Enthält ein vollständiges Paket- und Abonnementsystem: Abrechnung pro Token / pro Anfrage, periodische Ratenbegrenzung (stündlich / wöchentlich / monatlich), präzise Modellkontrolle sowie Unterstützung für empfohlene Pakete und Preisgestaltung.

| Paketverwaltung | Abonnementverwaltung |
|:---:|:---:|
| ![Paketverwaltung](../docs/Demo-Plan.png) | ![Abonnementverwaltung](../docs/Demo-Subscribe.png) |

### 💳 Bestellungen und echte Zahlungen

Jede Paketbestellung hinterlässt einen **vollständigen Bestell-Auditdatensatz** (Bestellnummer, Benutzer, Paket-Snapshot als JSON, Betrag, Zahlungsmethode, Status, Zahlungszeit, Kanalseriennummer). Unterstützt werden die beiden Bestelltypen Paket / Aufladung. Nativ integriert sind **WeChat Pay Native** (PC-QR-Code) und **Alipay Face-to-Face Pay** (TradePrecreate); zusätzlich sind die drei Kanäle Bank / Offline / Frei für die Verwaltung vorkonfiguriert. Der Aufpreis beim Paket-Upgrade wird automatisch proportional zu den verbleibenden Tagen berechnet; im Additionsmodus laufen altes und neues Paket parallel. Alle Regeln lassen sich im Untertab „Paketbetrieb“ unter „Betrieb“ zur Laufzeit umschalten.

| Bestellzentrum | Zahlungskonfiguration |
|:---:|:---:|
| ![Bestellzentrum](../docs/Demo-Order.png) | ![Zahlungskonfiguration](../docs/Demo-Payment.png) |

### 🌐 Dezentraler Active-Active-Cluster

Unterstützt die Bereitstellung dezentraler Active-Active-Cluster. Jeder Knoten betreibt ein eigenes MySQL + Redis und synchronisiert Daten über Ereignisse auf Anwendungsebene – ohne gemeinsame Datenbank, ideal für latenzarme Zugriffe in mehreren Regionen weltweit.

| Cluster-Knotenverwaltung |
|:---:|
| ![Cluster-Knotenverwaltung](../docs/Demo-cluster.png) |

### 🧩 Weitere Kernfunktionen

- **Über 30 Modellplattformen**: OpenAI / Anthropic / Gemini / DeepSeek / Qwen / Wenxin / iFlytek / Zhipu u. a. – vollständig abgedeckt und einheitlich hinter einer OpenAI-kompatiblen Schnittstelle
- **Präzise Kostenrechnung**: Abrechnung pro Token oder pro Anfrage, unabhängige Preisgestaltung für Prompt / Completion / Cached, gestapelte Gruppenrabatte, Verfolgung der Nutzung im Zeitraum
- **Kanal-Lastausgleich**: gewichtete Zufallsverteilung, automatisches Failover, Kühl-/Deaktivierungsstrategien, Kanal-Konkurrenz und RPM-Ratenbegrenzung
- **Mehrstufiges Rechtesystem**: die vier Ebenen Guest / User / Admin / Root, behebt die API-Berechtigungslücken des Originals und verfeinert die Admin-Berechtigungen
- **Unternehmensweite Sicherheit**: durchgängiges HTTPS, Token-Authentifizierung, IP-Subnetz-Beschränkungen, Echtzeit-Audit-Logs

---

## 🔥 Vergleich mit one-api

| Vergleichsdimension | one-api | one-api-pro |
| --- | --- | --- |
| Projektname | one-api | one-api-pro |
| Adaptor-Architektur | Zentrale Konstantenverwaltung (channeltype/define.go mit 56 Zeilen iota + parallele Arrays in url.go + zweistufiges switch in helper.go). Ein neuer Anbieter erfordert Änderungen an 4 Framework-Dateien | Selbstregistrierungsmechanismus (registry + register.go). Ein neuer Anbieter benötigt nur ein neues Paket plus eine Registrierung – ohne Änderungen am Framework-Code |
| Rechtegranularität | Die Grenzen zwischen Admin- und Normalnutzer-Rechten sind verschwommen; jeder kann Einstellungen über die API ändern | Gestuftes Rechtesystem, behebt API-Berechtigungslücken, verfeinert die Admin-Rechte |
| Abonnementmodus | Kein Paket-/Abonnementsystem | Vollständiges Paket-Abonnement + periodische Ratenbegrenzung + Modellkontrolle |
| Dezentraler Cluster | Keine unabhängige Cluster-Unterstützung; Mehrfach-Bereitstellung erfordert gemeinsame MySQL | Unterstützt dezentralen Active-Active-Cluster, jeder Knoten mit eigener MySQL + Redis, Datenvertrauen durch Ereignissynchronisierung auf Anwendungsebene, ohne gemeinsame Datenbank |
| Verzeichnisstruktur | relay/adaptor/ flach mit 40 Verzeichnissen; Basisprotokolle und Anbieter vermischt; relay/model/ kollidiert mit dem Wurzelverzeichnis model/ | adaptor/openai/ und adaptor/anthropic/ als Basisprotokolle separat; adaptor/provider/ sammelt 37 Anbieter; relay/schema/ beseitigt Namenskonflikte |
| Admin-Interface | 3 Frontend-Themes (default / berry / air), grundlegende Verwaltungsfunktionen | Brandneues Vue 3 + Arco Design Admin-Interface mit visualisiertem Dashboard |
| Laufende Aktualisierung | Das Originalprojekt wurde 2024 eingestellt | Wird kontinuierlich gepflegt und für unternehmensweite Szenarien optimiert |

---

## 📸 Screenshots

### 🖥️ Dashboard
![Dashboard](../docs/Demo-Index.png)

### 🔑 Token-Verwaltung
![Token-Verwaltung](../docs/Demo-Token.png)

### 📦 Paketverwaltung
![Paketverwaltung](../docs/Demo-Plan.png)

### 🔄 Abonnementverwaltung
![Abonnementverwaltung](../docs/Demo-Subscribe.png)

### 🌐 Cluster-Knotenverwaltung
![Cluster-Knotenverwaltung](../docs/Demo-cluster.png)

---

## ⚙️ Konfiguration

Das System ist sofort einsatzbereit.

Du kannst es über Umgebungsvariablen oder Befehlszeilenargumente konfigurieren; nach dem Start melde dich als `root` an und fahre mit der Konfiguration im Admin-Interface fort.

> **Hinweis:** Wenn du die Bedeutung einer Einstellung nicht kennst, lösche ihren Wert vorübergehend, um einen Hinweistext zu sehen.

### 🔧 Umgebungsvariablen

> One Api Pro liest Umgebungsvariablen aus einer `.env`-Datei. Beachte die Datei `.env.example`; benenne sie vor der Verwendung in `.env` um. Du kannst den Pfad auch über den Parameter `--env` angeben (relative und absolute Pfade werden unterstützt) – siehe Abschnitt Befehlszeilenargumente.

1. `REDIS_CONN_STRING`: Nach der Aktivierung wird Redis als Cache verwendet.
   + Beispiel: `REDIS_CONN_STRING=redis://default:redispw@localhost:49153`
   + Wenn die Datenbankzugriffe bereits sehr niedrige Latenz haben, ist Redis nicht notwendig; die Aktivierung kann sogar zu Daten-Hinterherhinken führen.
   + Falls du den Sentinel- oder Cluster-Modus benötigst:
     + Dann setze diese Variable auf eine Knotenliste, z. B. `localhost:49153,localhost:49154,localhost:49155`.
     + Zusätzlich sind noch folgende Umgebungsvariablen zu setzen:
       + `REDIS_PASSWORD`: Passwort für den Redis-Cluster- oder Sentinel-Modus.
       + `REDIS_MASTER_NAME`: Name des Master-Knotens im Redis-Sentinel-Modus.
2. `SESSION_SECRET`: Setzt einen festen Sitzungsschlüssel, sodass die Cookies angemeldeter Benutzer nach einem Systemneustart gültig bleiben.
   + Beispiel: `SESSION_SECRET=random_string`
3. `SQL_DSN`: Nach der Aktivierung wird die angegebene Datenbank statt SQLite verwendet – bitte MySQL oder PostgreSQL verwenden.
   + Beispiele:
     + MySQL: `SQL_DSN=root:123456@tcp(localhost:3306)/oneapi`
     + PostgreSQL: `SQL_DSN=postgres://postgres:123456@localhost:5432/oneapi` (in Arbeit, Feedback willkommen)
   + Beachte, dass die Datenbank `oneapi` vorher angelegt werden muss; Tabellen werden automatisch erstellt.
   + Bei Cloud-Datenbanken: Wenn der Cloud-Server eine Identitätsprüfung verlangt, füge `?tls=skip-verify` zu den Verbindungsparametern hinzu.
   + Passe die folgenden Parameter an deine Datenbankkonfiguration an (oder lasse die Standardwerte):
     + `SQL_MAX_IDLE_CONNS`: Maximale Anzahl inaktiver Verbindungen, Standard `100`.
     + `SQL_MAX_OPEN_CONNS`: Maximale Anzahl geöffneter Verbindungen, Standard `1000`.
       + Bei der Fehlermeldung `Error 1040: Too many connections` reduziere diesen Wert angemessen.
     + `SQL_CONN_MAX_LIFETIME`: Maximale Lebensdauer einer Verbindung, Standard `60`, Einheit Minuten.
4. `LOG_SQL_DSN`: Nach der Aktivierung wird für die `logs`-Tabelle eine eigene Datenbank verwendet – bitte MySQL oder PostgreSQL verwenden.
5. `FRONTEND_BASE_URL`: Nach der Aktivierung werden Seitenanfragen an die angegebene Adresse umgeleitet – nur auf dem Server festlegen.
   + Beispiel: `FRONTEND_BASE_URL=https://openai.justsong.cn`
6. `MEMORY_CACHE_ENABLED`: Aktiviert den Speicher-Cache, was zu einer gewissen Verzögerung bei Kontingentaktualisierungen führt. Werte: `true` und `false`; nicht gesetzt ist der Standard `false`.
   + Beispiel: `MEMORY_CACHE_ENABLED=true`
7. `SYNC_FREQUENCY`: Häufigkeit (in Sekunden), mit der bei aktiviertem Cache die Konfiguration mit der Datenbank synchronisiert wird, Standard `600` Sekunden.
   + Beispiel: `SYNC_FREQUENCY=60`
8. `NODE_TYPE`: Nach der Aktivierung wird der Knotentyp festgelegt. Werte: `master` und `slave`; nicht gesetzt ist der Standard `master`.
   + Beispiel: `NODE_TYPE=slave`
9. `CHANNEL_UPDATE_FREQUENCY`: Nach der Aktivierung wird das Kanalguthaben regelmäßig aktualisiert (Einheit Minuten); nicht gesetzt ist keine Aktualisierung.
   + Beispiel: `CHANNEL_UPDATE_FREQUENCY=1440`
10. `CHANNEL_TEST_FREQUENCY`: Nach der Aktivierung werden die Kanäle regelmäßig geprüft (Einheit Minuten); nicht gesetzt ist keine Prüfung.
    + Beispiel: `CHANNEL_TEST_FREQUENCY=1440`
11. `POLLING_INTERVAL`: Intervall (in Sekunden) zwischen Anfragen zur massenhaften Kontingentaktualisierung und Verfügbarkeitstests, standardmäßig kein Intervall.
    + Beispiel: `POLLING_INTERVAL=5`
12. `BATCH_UPDATE_ENABLED`: Aktiviert die Aggregation von Datenbank-Batch-Updates, was zu einer gewissen Verzögerung bei Kontingentaktualisierungen führt. Werte: `true` und `false`; nicht gesetzt ist der Standard `false`.
    + Beispiel: `BATCH_UPDATE_ENABLED=true`
    + Wenn du zu viele Datenbankverbindungen hast, versuche diese Option zu aktivieren.
13. `BATCH_UPDATE_INTERVAL`: Zeitintervall der Batch-Update-Aggregation, Einheit Sekunden, Standard `5`.
    + Beispiel: `BATCH_UPDATE_INTERVAL=5`
14. Anfrageratenbegrenzung:
    + `GLOBAL_API_RATE_LIMIT`: Globale API-Ratenbegrenzung (außer Relay-Anfragen), maximale Anzahl Anfragen pro IP in drei Minuten, Standard `180`.
    + `GLOBAL_WEB_RATE_LIMIT`: Globale Web-Ratenbegrenzung, maximale Anzahl Anfragen pro IP in drei Minuten, Standard `60`.
15. Cache-Einstellungen für den Encoder:
    + `TIKTOKEN_CACHE_DIR`: Beim Programmstart werden die Token-Encodings gängiger Modelle (z. B. `gpt-3.5-turbo`, `gpt-4`, `gpt-4o`) von einer Online-Quelle heruntergeladen. Bei eingeschränktem Netzwerk oder offline greift nach einem Timeout (ca. 30 Sekunden) automatisch eine Näherungs-Zählung (ca. `0.38 × Zeichenzahl`), der Dienst startet weiterhin normal. Für exakte Abrechnung lade die Encoding-Dateien in einer Online-Umgebung vorab in dieses Verzeichnis und übertrage sie in die Offline-Umgebung.
    + `DATA_GYM_CACHE_DIR`: Diese Konfiguration wirkt derzeit wie `TIKTOKEN_CACHE_DIR`, hat jedoch eine geringere Priorität.
16. `RELAY_TIMEOUT`: Relay-Timeout in Sekunden, standardmäßig kein Timeout.
17. `RELAY_PROXY`: Nach der Aktivierung wird dieser Proxy für API-Anfragen verwendet.
18. `USER_CONTENT_REQUEST_TIMEOUT`: Timeout für das Herunterladen benutzerdefinierter Inhalte, Einheit Sekunden.
19. `USER_CONTENT_REQUEST_PROXY`: Nach der Aktivierung wird dieser Proxy für benutzerdefinierte Inhalte, z. B. Bilder, verwendet.
20. `SQLITE_BUSY_TIMEOUT`: Timeout für die SQLite-Sperrwarterei, Einheit Millisekunden, Standard `3000`.
21. `GEMINI_SAFETY_SETTING`: Sicherheitseinstellung von Gemini, Standard `BLOCK_NONE`.
22. `GEMINI_VERSION`: Gemini-Version, die One Api Pro verwendet, Standard `v1`.
23. `THEME`: Theme-Einstellung des Systems, Standard `default-pro` (Vue-3-Admin-Interface); alternativ `default` / `berry` / `air` (alte React-Themes). Weitere Werte siehe [hier](../web/README.md).
24. `ENABLE_METRIC`: Deaktiviert Kanäle basierend auf der Anfrage-Erfolgsrate, standardmäßig aus. Werte: `true` und `false`.
25. `METRIC_QUEUE_SIZE`: Größe der Warteschlange für die Erfolgsraten-Statistik, Standard `10`.
26. `METRIC_SUCCESS_RATE_THRESHOLD`: Schwellenwert für die Anfrage-Erfolgsrate, Standard `0.8`.
27. `INITIAL_ROOT_TOKEN`: Falls gesetzt, wird beim ersten Systemstart automatisch ein Root-Benutzer-Token mit diesem Wert erstellt.
28. `INITIAL_ROOT_ACCESS_TOKEN`: Falls gesetzt, wird beim ersten Systemstart automatisch ein Systemverwaltungs-Token für den Root-Benutzer mit diesem Wert erstellt.
29. `ENFORCE_INCLUDE_USAGE`: Erzwingt die Rückgabe von `usage` im Stream-Modus, standardmäßig aus. Werte: `true` und `false`.
30. `TEST_PROMPT`: Benutzer-Prompt beim Testen von Modellen, Standard `Print your model name exactly and do not output without any other text.`.

#### 🌐 Cluster-Konfiguration (dezentrale Active-Active-Bereitstellung)

> Sind die folgenden Umgebungsvariablen nicht gesetzt, läuft das System im Einzelknoten-Modus ohne Nebenwirkungen.

1. `CLUSTER_ENABLED`: Aktiviert den Cluster-Modus, standardmäßig deaktiviert.
   + Beispiel: `CLUSTER_ENABLED=true`
2. `CLUSTER_NODE_ID`: Knotennummer (1–49), muss mit `auto_increment_offset` von MySQL übereinstimmen; verschiedene Knoten dürfen sich nicht wiederholen.
   + Beispiel: `CLUSTER_NODE_ID=1`
3. `CLUSTER_NODE_NAME`: Knotenname zur Identifikation, Standard `node-{NODE_ID}`.
   + Beispiel: `CLUSTER_NODE_NAME=node-cn`
4. `CLUSTER_NODE_ADDRESS`: Öffentliche Adresse dieses Knotens (muss ein Protokollpräfix enthalten); andere Knoten pushen Daten an diese Adresse.
   + Beispiel: `CLUSTER_NODE_ADDRESS=https://cn.example.com`
5. `CLUSTER_SECRET`: Anfangs-Secret dieses Knotens, **pro Knoten unabhängig**. Es wird beim ersten Start als Startwert in die Datenbank geschrieben und kann später vom Admin geändert werden.
   + Beispiel: `CLUSTER_SECRET=MyClusterSecret123abc`
6. `CLUSTER_SEEDS`: Adressen der Seed-Knoten (kommagetrennt); ein neuer Knoten registriert sich beim Start bei den Seeds, um Cluster-Informationen zu erhalten. Ein erreichbarer Knoten genügt. Der erste Knoten kann leer bleiben oder die eigene Adresse enthalten.
   + Beispiel: `CLUSTER_SEEDS=https://cn.example.com`
   + Mehrere Seeds: `CLUSTER_SEEDS=https://cn.example.com,https://us.example.com`
7. `CLUSTER_PUSH_INTERVAL`: Push-Intervall für Synchronisationsereignisse, Einheit Sekunden, Standard `3`.
8. `CLUSTER_DISCOVERY_INTERVAL`: Intervall der Knoten-Erkennung, Einheit Sekunden; lebende Knoten pingen sich in jedem Zyklus gegenseitig, Standard `30`.
9. `CLUSTER_DEAD_PING_INTERVAL`: Ping-Intervall für ausgefallene Knoten, Einheit Sekunden; länger als das Lebensintervall, um unnötige Anfragen zu reduzieren, Standard `120`.
10. `CLUSTER_MAX_PING_FAILURES`: Anzahl aufeinanderfolgender Ping-Fehler; danach wird der Knoten als ausgefallen markiert, Standard `3`.
11. `CLUSTER_SYNC_LOGS`: Synchronisiert die Log-Tabelle; bei großer Datenmenge bedarfsweise deaktivierbar, Standard `true`.
    + Beispiel: `CLUSTER_SYNC_LOGS=false`
12. `CLUSTER_BATCH_SIZE`: Maximale Anzahl Ereignisse pro Push, Standard `50`.

### ⌨️ Befehlszeilenargumente

1. `--port <port_number>`: Legt den Port fest, auf dem der Server lauscht, Standard `3000`.
   + Beispiel: `--port 3000`
2. `--log-dir <log_dir>`: Legt das Log-Verzeichnis fest; ohne Angabe wird standardmäßig im `logs`-Ordner des Arbeitsverzeichnisses gespeichert.
   + Beispiel: `--log-dir ./logs`
3. `--env <env_file_path>`: Legt den Pfad zur Konfigurationsdatei fest; relative und absolute Pfade werden unterstützt. Ohne Angabe wird die `.env`-Datei im aktuellen Verzeichnis geladen.
   + Beispiel: `--env ./config.env`
   + Beispiel: `--env /etc/one-api-pro/production.env`
   + Beispiel für mehrere Instanzen:
     ```bash
     ./one-api-pro --env ./instances/instance1.env --port 3001 &
     ./one-api-pro --env ./instances/instance2.env --port 3002 &
     ```
   + Konfigurationspriorität: Befehlszeilenargumente > System-Umgebungsvariablen > `--env`-Konfigurationsdatei > Standardwerte
4. `--version`: Gibt die Systemversion aus und beendet das Programm.
   + Beispiel: `./one-api-pro --version`
   + Versionsquellen (Priorität absteigend):
     1. Die `VERSION`-Datei im aktuellen Arbeitsverzeichnis oder im Verzeichnis der ausführbaren Datei (kompatibel mit oder ohne `v`-Präfix, z. B. `0.0.2` oder `v0.0.2`);
     2. Die zur Build-Zeit über `-ldflags "-X .../common.Version=..."` injizierte Version (`release.sh` und die CI injizieren sie automatisch);
     3. Der Standardwert im Quellcode: `common/constants.go`.
   + Daher genügt es, die `VERSION`-Datei im Projektstamm zu pflegen, damit `--version`, das Start-Log, der Endpunkt `/api/status` und die im Dashboard angezeigte Version übereinstimmen.
5. `--help`: Zeigt die Verwendungshilfe und die Parameterbeschreibung an.
   + Beispiel: `./one-api-pro --help`

---

## 📖 API-Dokumentation

Die vollständige API-Dokumentation wird separat in [docs/API.md](../docs/API.md) gepflegt und umfasst:

- **Authentifizierung**: drei Verfahren – Cookie-Session / Access Token / API Key (Bearer Token)
- **Admin-Endpunkte**: vollständiges CRUD für Modellpreise, Gruppenrabatte, Kanäle, Token, Benutzer, Logs, Einlösecodes, Pakete, Abonnements usw.
- **OpenAI-kompatible Endpunkte**: `/v1/models`, `/v1/chat/completions`, `/v1/embeddings`, Bilder, Audio, Inhaltsmoderation usw.
- **Cluster-API**: dezentrale Cluster-Endpunkte für Knoten-Erkennung, Heartbeat und Datensynchronisation

👉 [Vollständige API-Dokumentation ansehen →](../docs/API.md)

---

## 📦 Bereitstellung

### 🔨 Manuelle Bereitstellung

#### 1. Die ausführbare Datei beschaffen

Wähle eine der folgenden Varianten:

**Variante 1: Vorkompilierte Version herunterladen (empfohlen)**

Lade die eigenständige ausführbare Datei für deine Plattform (Linux / macOS / Windows) von den [GitHub Releases](https://github.com/modelbus/one-api-pro/releases/latest) herunter – ohne Entpacken direkt ausführbar.

**Variante 2: One-Click-Paketierung mit release.sh**

```shell
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro
./release.sh            # Paketierung für mehrere Plattformen, Ausgabe im dist/-Verzeichnis
```

**Variante 3: Aus dem Quellcode kompilieren**

```shell
git clone https://github.com/modelbus/one-api-pro.git
cd one-api-pro

# Frontend erstellen (Vue-3-Admin-Interface, gemäß den Themes aus web/THEMES)
cd web
sh build.sh

# Backend erstellen (Achtung: muss NACH dem Frontend ausgeführt werden, damit die neuesten Artefakte eingebettet werden)
cd ..
go build -ldflags "-s -w" -o one-api-pro
```

#### 2. Ausführen

```shell
chmod u+x one-api-pro
./one-api-pro --port 3000 --log-dir ./logs
```

#### 3. Zugriff

Öffne [http://localhost:3000/](http://localhost:3000/) und melde dich an. Das Standardkonto ist `root`, das Passwort `123456`.

### 🏢 Multi-Host-Bereitstellung
1. Setze `SESSION_SECRET` auf allen Servern auf denselben Wert.
2. Setze zwingend `SQL_DSN`, verwende MySQL statt SQLite – alle Server verbinden sich mit derselben Datenbank.
3. Alle Slave-Server müssen `NODE_TYPE` auf `slave` setzen; ohne Angabe gilt der Master-Server.
4. Setze `SYNC_FREQUENCY`, damit die Server die Konfiguration regelmäßig aus der Datenbank synchronisieren. Bei einer entfernten Datenbank wird dies unabhängig von Master/Slave empfohlen; aktiviere dabei Redis.
5. Slave-Server können optional `FRONTEND_BASE_URL` setzen, um Seitenanfragen zum Master umzuleiten.
6. Richte auf den Slave-Servern **jeweils** ein eigenes Redis ein und setze `REDIS_CONN_STRING` entsprechend. So kann die Datenbank bei warmem Cache gar nicht erreicht werden, was die Latenz senkt (zur Sentinel-/Cluster-Unterstützung siehe die Umgebungsvariablen).
7. Wenn auch der Master eine hohe Datenbanklatenz hat, aktiviere auch dort Redis und setze `SYNC_FREQUENCY`, um die Konfiguration regelmäßig aus der Datenbank zu synchronisieren.

Einzelheiten zur Verwendung der Umgebungsvariablen findest du [hier](#-umgebungsvariablen).

### 🌐 Dezentrale Cluster-Bereitstellung

Der Cluster-Modus erlaubt es mehreren Knoten, jeweils ein eigenes One Api Pro + MySQL bereitzustellen und Daten über Ereignisse auf Anwendungsebene zu synchronisieren – ohne gemeinsame Datenbank.

> **Anwendungsfälle**: globale Bereitstellung in mehreren Regionen, latenzarme lokale Zugriffe, Hochverfügbarkeit/Notfallwiederherstellung, Lastausgleich über mehrere Knoten.

#### 🗺️ Architektur

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

#### ⭐ Kernmerkmale

- **Dezentral**: Alle Knoten sind gleichberechtigt, ohne Master/Slave; jede Datenänderung wird aktiv an alle lebenden Knoten gepusht
- **Nicht-invasiv**: GORM-Callbacks erfassen Datenänderungen, ohne bestehenden Geschäftslogik-Code zu verändern
- **Asynchroner Push**: Die Datensynchronisation blockiert den Hauptablauf nicht; sie erfolgt über einen Hintergrund-Goroutine
- **Konfliktlösung**: Vergleich anhand des `updated_at`-Zeitstempels; nur neuere Daten werden geschrieben
- **Ratenbegrenzungssynchronisation**: Kanal-Konkurrenz- und RPM-Zähler werden über Datenbanktabellen knotenübergreifend synchronisiert
- **Einzelknoten-kompatibel**: Ohne Cluster-Umgebungsvariablen läuft das System vollständig im Einzelknoten-Modus

#### 📊 Synchronisationsumfang

| Datentabelle | Synchronisiert? | Beschreibung |
| --- | --- | --- |
| users | ✅ | Benutzerinformationen |
| tokens | ✅ | API-Token |
| channels | ✅ | Kanal-Konfigurationen |
| abilities | ✅ | Kanal-Fähigkeiten |
| options | ✅ | Systemeinstellungen |
| redemptions | ✅ | Einlösecodes |
| plans | ✅ | Abonnement-Pakete |
| user_plans | ✅ | Benutzer-Abonnements |
| plan_usages | ✅ | Paket-Nutzung |
| channel_counters | ✅ | Ratenbegrenzungszähler der Kanäle |
| cluster_nodes | 🔄 Discovery | Knoteninformationen (werden vom Erkennungsmechanismus gepflegt, nicht über Datensynchronisation) |
| logs | ⚠️ Optional | Log-Daten sind groß; via `CLUSTER_SYNC_LOGS` steuerbar |

#### 🚀 Bereitstellungsschritte

**1. MySQL-Konfiguration (jeder Knoten muss eine eigene MySQL-Instanz verwenden)**

Jeder Knoten benötigt eine **eigene MySQL-Instanz** (es ist nicht möglich, mehrere Knoten durch das Anlegen mehrerer Datenbanken in derselben Instanz zu betreiben, da `auto_increment_offset` eine instanzweite Variable ist).

```ini
# Node 1 my.cnf
[mysqld]
server-id = 1
auto_increment_increment = 50
auto_increment_offset = 1
log_bin = mysql-bin
binlog_format = ROW

# Node 2 my.cnf
[mysqld]
server-id = 2
auto_increment_increment = 50
auto_increment_offset = 2
log_bin = mysql-bin
binlog_format = ROW

# Node 3 my.cnf
[mysqld]
server-id = 3
auto_increment_increment = 50
auto_increment_offset = 3
log_bin = mysql-bin
binlog_format = ROW
```

> `auto_increment_increment` ist auf 50 gesetzt – es unterstützt maximal 50 Knoten. Der `offset` jedes Knotens muss mit `CLUSTER_NODE_ID` übereinstimmen und unterschiedlich sein.

> **Wichtig:** `auto_increment_increment` und `auto_increment_offset` sind **Instanz-weite Variablen** von MySQL und gelten für alle Datenbanken der Instanz. Sie lassen sich weder pro Datenbank noch auf Tabellenebene unterschiedlich setzen (die Tabellenoption von MySQL unterstützt nur einen `AUTO_INCREMENT`-Startwert, keine Schrittweite). Daher **muss** jeder Knoten eine **eigene MySQL-Instanz** verwenden. Es ist nicht möglich, mehrere Knoten durch das Anlegen verschiedener Datenbanken in derselben MySQL-Instanz zu betreiben. Um mehrere Instanzen auf derselben Maschine zu betreiben, starte mehrere `mysqld`-Prozesse auf unterschiedlichen Ports oder verwende mehrere Docker-Container.

> **Zu `server-id` und binlog:** `server-id` muss in allen MySQL-Instanzen des Clusters verschieden sein. `log_bin` und `binlog_format=ROW` werden dringend empfohlen – sie dienen künftiger Master-Slave-Replikation und der Point-in-Time-Recovery. Die Cluster-Datensynchronisation selbst hängt nicht vom binlog ab (sie erfolgt über GORM-Callbacks auf Anwendungsebene), doch das binlog bietet zusätzliche Zuverlässigkeit.

**2. Redis-Konfiguration (jeder Knoten muss eine eigene Redis-Instanz verwenden)**

Jeder Knoten benötigt ebenfalls eine **eigene Redis-Instanz** (anderer Port oder andere Maschine). Redis dient in dieser Cluster-Architektur nicht der Kommunikation zwischen Knoten, sondern nur dem lokalen Cache, der Ratenbegrenzung und anderen Geschäftszwecken.

**3. Neuen Knoten initialisieren**

Beim Hinzufügen eines neuen Knotens muss zuerst ein Daten-Snapshot von einem bestehenden Knoten bezogen werden:

```bash
# Variante 1: Aus einem bestehenden Knoten exportieren und importieren
mysqldump -h existing-node -u root -p oneapi > backup.sql
mysql -u root -p oneapi < backup.sql

# Variante 2: Snapshot über die API holen (Dienst muss zuerst laufen)
curl -H "X-Cluster-Secret: your-secret" \
  "https://existing-node/api/cluster/snapshot?tables=users,tokens,channels,abilities,options,redemptions,plans,user_plans,plan_usages" \
  -o snapshot.json
```

**4. Umgebungsvariablen-Konfiguration (vollständiges Beispiel)**

Nachfolgend ein vollständiges `.env`-Konfigurationsbeispiel für einen 3-Knoten-Cluster. Jeder Knoten nutzt eine eigene MySQL- und Redis-Instanz, mit unterschiedlichen Ports und Pfaden.

**Knoten 1 – China-Knoten (`/opt/one-api-pro/node1/.env`):**
```bash
# ========================
# Grundkonfiguration
# ========================
PORT=3000
SYSTEM_NAME=One Api Pro Cluster

# ========================
# Datenbank (eigene MySQL-Instanz)
# ========================
SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node1?charset=utf8mb4&parseTime=True&loc=Local

# ========================
# Redis (eigene Redis-Instanz)
# ========================
REDIS_CONN_STRING=redis://127.0.0.1:6379/0

# ========================
# Cluster-Konfiguration
# ========================
CLUSTER_ENABLED=true
CLUSTER_NODE_ID=1
CLUSTER_NODE_NAME=node-cn
CLUSTER_NODE_ADDRESS=https://cn.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me

# Seed-Knoten (nur für die anfängliche Erkennung anderer Knoten)
# Erster Knoten: eigene Adresse oder leer
# Weitere Knoten: Adresse eines beliebigen lebenden Knotens
CLUSTER_SEEDS=https://cn.example.com,https://us.example.com,https://eu.example.com

# ========================
# Cluster-Tuning (optional)
# ========================
CLUSTER_DISCOVERY_INTERVAL=30
CLUSTER_DEAD_PING_INTERVAL=120
CLUSTER_MAX_PING_FAILURES=3
CLUSTER_PUSH_INTERVAL=3
CLUSTER_SYNC_LOGS=true
CLUSTER_BATCH_SIZE=50
```

**Knoten 2 – US-Knoten (`/opt/one-api-pro/node2/.env`):**
```bash
# Grundkonfiguration
PORT=3001
SYSTEM_NAME=One Api Pro Cluster

# Datenbank (eigene MySQL-Instanz, Port oder Maschine abweichend von Knoten 1)
SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node2?charset=utf8mb4&parseTime=True&loc=Local

# Redis (eigene Redis-Instanz)
REDIS_CONN_STRING=redis://127.0.0.1:6380/0

# Cluster-Konfiguration
CLUSTER_ENABLED=true
CLUSTER_NODE_ID=2
CLUSTER_NODE_NAME=node-us
CLUSTER_NODE_ADDRESS=https://us.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me   # muss exakt mit Knoten 1 übereinstimmen

# Adresse eines beliebigen lebenden Knotens
CLUSTER_SEEDS=https://cn.example.com
```

**Knoten 3 – Europa-Knoten (`/opt/one-api-pro/node3/.env`):**
```bash
# Grundkonfiguration
PORT=3002
SYSTEM_NAME=One Api Pro Cluster

# Datenbank
SQL_DSN=root:password@tcp(127.0.0.1:3306)/oneapi_node3?charset=utf8mb4&parseTime=True&loc=Local

# Redis
REDIS_CONN_STRING=redis://127.0.0.1:6381/0

# Cluster-Konfiguration
CLUSTER_ENABLED=true
CLUSTER_NODE_ID=3
CLUSTER_NODE_NAME=node-eu
CLUSTER_NODE_ADDRESS=https://eu.example.com
CLUSTER_SECRET=your-strong-shared-secret-key-change-me   # muss mit allen Knoten übereinstimmen

# Adresse eines beliebigen lebenden Knotens
CLUSTER_SEEDS=https://cn.example.com
```

**Gegenüberstellung der Konfigurationsparameter:**

| Umgebungsvariable | Knoten 1 | Knoten 2 | Knoten 3 | Beschreibung |
|---|---|---|---|---|
| `PORT` | 3000 | 3001 | 3002 | Lauschport (auf derselben Maschine unterschiedlich) |
| `SQL_DSN` | ...oneapi_node1 | ...oneapi_node2 | ...oneapi_node3 | eigene MySQL-Instanz |
| `REDIS_CONN_STRING` | :6379/0 | :6380/0 | :6381/0 | eigene Redis-Instanz |
| `CLUSTER_NODE_ID` | 1 | 2 | 3 | Knotennummer, entspricht MySQL `auto_increment_offset` |
| `CLUSTER_NODE_NAME` | node-cn | node-us | node-eu | Knotenname zur Identifikation |
| `CLUSTER_NODE_ADDRESS` | https://cn.example.com | https://us.example.com | https://eu.example.com | öffentliche Adresse (andere Knoten greifen darüber zu) |
| `CLUSTER_SECRET` | derselbe Wert | derselbe Wert | derselbe Wert | **Alle Knoten müssen exakt übereinstimmen** |
| `CLUSTER_SEEDS` | eigene Adresse oder leer | beliebiger lebender Knoten | beliebiger lebender Knoten | anfängliche Erkennung, danach automatische Entdeckung |

**5. Startbefehle**

Jeder Knoten lädt seine Konfigurationsdatei über den Parameter `--env`:

```bash
# Knoten 1
./one-api-pro --env /opt/one-api-pro/node1/.env --port 3000

# Knoten 2
./one-api-pro --env /opt/one-api-pro/node2/.env --port 3001

# Knoten 3
./one-api-pro --env /opt/one-api-pro/node3/.env --port 3002
```

**6. Startreihenfolge**

1. Starte den ersten Knoten (Knoten A); lasse `CLUSTER_SEEDS` leer oder verwende die eigene Adresse.
2. Warte, bis Knoten A vollständig gestartet ist (ca. 5–10 Sekunden; achte auf das Log „Cluster-Modul initialisiert“).
3. Starte die weiteren Knoten und trage bei `CLUSTER_SEEDS` die Adresse eines beliebigen lebenden Knotens ein.
4. Die neuen Knoten pingen nach dem Start automatisch die Seed-Knoten und entdecken transitiv alle anderen Knoten.
5. Sobald alle Knoten laufen, kannst du unter „Einstellungen → Knotenverwaltung“ im Admin-Interface eines beliebigen Knotens den Knotenstatus einsehen.

**7. Nginx-Lastausgleich-Beispiel (optional)**

```nginx
upstream one_api_cluster {
    ip_hash;  # gleicher Client landet immer auf demselben Knoten – bewahrt Session & Cache
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

> **`ip_hash` ist entscheidend:** Es stellt sicher, dass Anfragen desselben Benutzers immer auf demselben Knoten landen, damit Zustände wie Plan-Limits und Redis-Cache nicht zwischen den Knoten verloren gehen.

**8. Cluster-Status prüfen**

Nach der Bereitstellung kannst du wie folgt prüfen:

```bash
# Liste der Knoten abrufen (auf einem beliebigen Knoten aufrufen)
curl -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  https://cn.example.com/api/cluster_node/

# Sollte die Liste aller Knoten mit Feldern wie status, last_heartbeat, ping_failures zurückgeben
```

Oder rufe im Admin-Interface die Seite **Einstellungen → Knotenverwaltung** auf, um Knotenliste, Status und letzte Herzschlagzeit zu sehen.

> 💡 Einzelheiten zur Cluster-API: [docs/API.md, Anhang E – Cluster-API](../docs/API.md#附录-e集群管理-api)

#### ⚠️ Hinweise

- Jeder Knoten muss eine eigene MySQL- und Redis-Instanz besitzen; keine gemeinsame Datenbank.
- `CLUSTER_SECRET` muss in allen Knoten identisch sein. Verwende ein starkes Passwort und bewahre es sicher auf.
- `CLUSTER_NODE_ID` muss in allen Knoten unterschiedlich sein und mit MySQL `auto_increment_offset` übereinstimmen.
- `CLUSTER_NODE_ADDRESS` muss eine für andere Knoten erreichbare öffentliche Adresse sein (inkl. Protokollpräfix wie `https://`).
- Die Dateninitialisierung neuer Knoten muss manuell erfolgen (Snapshot von einem Online-Knoten ziehen).
- Die Log-Tabelle (`logs`) ist groß; die Log-Synchronisation kann mit `CLUSTER_SYNC_LOGS=false` deaktiviert werden.
- MySQL `auto_increment_increment` und `auto_increment_offset` müssen zur `CLUSTER_NODE_ID`-Konfiguration passen.
- Die Knoten-Erkennung nutzt einen bidirektionalen Ping-Mechanismus; ausgefallene Knoten werden nicht gelöscht, sondern nur als status=2 markiert und leben nach Netzerwiederherstellung automatisch wieder auf.
- `CLUSTER_SEEDS` ist nur die anfängliche Erkennung; sobald ein Knoten andere per Ping entdeckt hat, verlässt er sich nicht mehr auf SEEDS.
- Änderungen anderer Knoten während der Offline-Phase werden **nicht automatisch nachgeliefert**. Nach dem Wiederauftauchen muss der Offline-Knoten einen Snapshot ziehen, um die Daten aufzufüllen.

#### 📝 Zur Selbstregistrierung des „lokalen Knotens“

Jeder Knoten schreibt beim Start in der eigenen Tabelle `cluster_nodes` einen Datensatz über sich selbst (`node_id` gleich dem lokal konfigurierten `CLUSTER_NODE_ID`). Das ist **gewollt**, aus folgenden Gründen:

1. **Anzeige im Admin-Interface**: Auf der Seite „Einstellungen → Knotenverwaltung“ muss der Admin lokale Angaben (Adresse, Status, Herzschlagzeit) sehen, um Probleme zu beheben.
2. **Transitive Erkennung**: Wenn Knoten B einen Ping von Knoten A erhält, gibt A in der Antwort die vollständige Knotenliste zurück (inklusive A selbst). B führt sie in der lokalen Tabelle zusammen. So kann C über die Antwort von B auch die Existenz von A erfahren.
3. **Basis für die Lebenderkennung**: Das lokale `last_heartbeat` wird alle 30 Sekunden (in der Funktion `discoverOnce`) automatisch aktualisiert und spiegelt den normalen Betriebsstatus dieses Knotens wider.

**Die Selbstregistrierung verursacht keine Loop-Synchronisation.** Das System schützt sich auf fünf Ebenen:

| Schutzmaßnahme | Wirkung |
|---|---|
| ① `GetAllRemoteNodes` SQL-Filter | Bei der Erkennung schließt die SQL-Bedingung `node_id != ?` den lokalen Knoten aus |
| ② `GetAliveNodesForSync` SQL-Filter | Beim Push schließt die SQL-Bedingung `node_id != ?` den lokalen Knoten aus |
| ③ `handlePing` verweigert Selbst-Ping | `req.NodeId == NodeID` wird explizit abgelehnt |
| ④ `mergeDiscoveredNodes` überspringt den lokalen Knoten | Fügt beim Zusammenführen der erkannten Knoten den lokalen Knoten nicht ein |
| ⑤ `ApplyEvents` überspringt lokale Ereignisse | Bearbeitet lokale Ereignisse beim Anwenden nicht |

Der Datenfluss ist unidirektional: vom lokalen zum entfernten Knoten und vom entfernten zurück zum lokalen – **es gibt nie eine Schleife**.

Das Admin-Interface zeigt neben dem lokalen Knotennamen ein blaues Abzeichen „lokaler Knoten“ und deaktiviert für den lokalen Knoten die Aktionen „Löschen“ und „Manuelles Pingen“ (beide sind für den lokalen Knoten bedeutungslos).

#### 🔐 Zu „pro Knoten unabhängigen Secrets“

Jeder Knoten besitzt **sein eigenes Secret**; ein global geteiltes Secret wird nicht mehr verwendet. Gründe:

1. **Sicherheit**: Ein Leak des Secrets eines Knotens wirkt sich nicht auf andere Knoten aus.
2. **Flexible Verwaltung**: Jeder Knoten kann sein Secret unabhängig rotieren.
3. **Automatische Erkennung**: Beim Pingen übermitteln die Knoten automatisch ihr eigenes Secret, damit die Gegenseite es speichert.

**Lebenszyklus des Secrets:**
- Beim ersten Start: `CLUSTER_SECRET` wird als Startwert in das Feld `cluster_nodes.secret_key` geschrieben.
- Bei späteren Starts: aus `cluster_nodes.secret_key` gelesen.
- Admin können das Secret anderer Knoten auf der Seite „Knotenverwaltung“ ändern.
- Beim Pingen enthält der Header `X-Cluster-Secret` das Secret des **Zielknotens** (aus der lokalen DB gelesen).

**Ablauf beim Hinzufügen eines neuen Knotens:**
1. Füge auf Knoten A den Datensatz für Knoten B hinzu und trage den Wert von B's `CLUSTER_SECRET` ein.
2. Füge auf Knoten B den Datensatz für Knoten A hinzu und trage den Wert von A's `CLUSTER_SECRET` ein.
3. A pingt B: mit dem Secret von B; B empfängt: verifiziert das eigene Secret ✓.
4. Die Antwort von B enthält die Secrets von A und B; A aktualisiert die lokal gespeicherten Werte.

#### 🗑️ Zum „Soft-Delete“ von Knoten

Beim Löschen eines Knotens **löscht** der Admin den Datensatz **nicht physisch**, sondern setzt `disabled = true`:

- Verhindert, dass der gelöschte Knoten „automatisch wieder auftaucht“ (der Ping-Mechanismus würde ihn erneut registrieren).
- Deaktivierte Knoten beantworten weiterhin Pings (damit die Gegenseite den Online-Status kennt), holen aber keine Informationen dieses Knotens ab.
- Das physische Löschen erfordert manuelles SQL: `DELETE FROM cluster_nodes WHERE node_id = ?`.

#### 🔄 Zum „Datensynchronisationsmechanismus“ (wichtig)

**Die Cluster-Datensynchronisation** beruht vollständig auf **GORM-Ereignissen + aktivem HTTP-Push**:
- Beliebige INSERT/UPDATE/DELETE-Operationen auf Geschäftstabellen → von GORM-Callbacks erfasst → in die Tabelle `sync_events` geschrieben → vom Pusher-Goroutine an alle lebenden Knoten gepusht.
- Der Empfänger schreibt über `WithSkipHook` in die lokale Datenbank (kein Loop-Back).
- Der Empfänger überspringt zusätzlich Ereignisse mit `event.NodeId == lokale NodeID` (doppelte Absicherung).

**Architekturabwägung:** Dieses Design **implementiert keinen aktiven Pull zwischen Knoten**, aus folgenden Gründen:
1. **Eingriff in die Geschäftslogik**: Der knotenübergreifende Pull müsste die geschäftlich eindeutigen Felder jeder Tabelle kennen und würde die Geschäftslogik beeinträchtigen.
2. **Primärschlüssel-Konflikte**: Die Auto-Increment-IDs sind knotenübergreifend nicht kontinuierlich (verschiedene `auto_increment_offset`); die Verwendung der Quell-ID würde das Offset-Design zerstören.
3. **Hohe Komplexität**: hoher Wartungsaufwand, begrenzter Zuverlässigkeitsgewinn.
4. **Aktiver Push genügt**: 95 % der Szenarien (normale Synchronisation bei Online-Knoten) sind vollständig durch Push abgedeckt.

**Bekannte Grenzen und Betriebsanforderungen:**
- Während der Offline-Phase erzeugte Datenänderungen → **dauerhaft verloren** (Push erfolgt nur in Echtzeit).
- Nach dem Wiederauftauchen können die in der Offline-Phase entstandenen Daten nicht automatisch ergänzt werden.
- Ein neuer Knoten erhält nur die Änderungen nach seinem Beitritt; es gibt keine historischen Daten.
- **Abhilfe durch den Betrieb**: mit `mysqldump` aus einem anderen Knoten exportieren und importieren.

**Typische Bereitstellungsszenarien im Vergleich:**

| Szenario | Pull erforderlich? | Vorgehen |
|---|---|---|
| Knoten dauerhaft online | ❌ | Push reicht völlig aus |
| Knoten startet gelegentlich neu (im Minutenbereich) | ⚠️ | Kurzer Datenverlust bei Offline-Zeit, betrieblich akzeptabel |
| Knoten wird häufig gewartet | ❌ | Push läuft weiter; nach Neustart sofort wiederhergestellt |
| Neuer Knoten tritt dem Cluster bei | ❌ | DBA initialisiert manuell mit `mysqldump` |
| Knoten nach langer Offline-Phase wieder da | ❌ | DBA ergänzt manuell mit `mysqldump` |

Wenn nach der Bereitstellung eine leere Seite erscheint, siehe [#97](https://github.com/modelbus/one-api-pro/issues/97).

---

## 🗺️ Entwicklungsplan

### ✅ Abgeschlossen

- [x] **Architektur-Umbau**: Selbstregistrierungs-Mechanismus der Adaptoren; neue Anbieter erfordern keine Framework-Änderungen.
- [x] **Brandneues Vue-3-Admin-Interface**: Arco Design + visualisiertes Dashboard + über 30 Modellplattform-Icons.
- [x] **Paket- und Abonnementsystem**: Abrechnung pro Token / pro Anfrage, periodische Ratenbegrenzung, präzise Modellkontrolle.
- [x] **Dezentraler Active-Active-Cluster**: durch GORM-Ereignisse + aktiven HTTP-Push synchronisiert, ohne gemeinsame Datenbank.
- [x] **Präzise Kostenrechnung**: unabhängige Preisgestaltung für Prompt / Completion / Cached, gestapelte Gruppenrabatte.
- [x] **Mehrstufiges Rechtesystem**: die Ebenen Guest / User / Admin / Root; behebt die API-Berechtigungslücken des Originals.
- [x] **OpenAI-kompatible Schnittstelle**: vollständige Unterstützung für models / chat / completions / embeddings / images / audio / moderations.
- [x] **Bestell- und Upgrade-Ablauf für Pakete**: natives `POST /api/order/plan` erstellt die Bestellungen für das Paket-Abonnement; unterstützt `stack` (additiv) und `price_diff` (Aufpreis-Upgrade) – der Aufpreis wird automatisch über die verbleibenden Tage berechnet, inklusive Prüfung auf gleiche/absteigende Stufe.
- [x] **Bestell-Audit und Bestellzentrum**: neue Tabelle `orders` (type/source/order_no/plan_info/amount/status/pay_status/pay_method/pay_time/pay_trade_no) speichert alle Zahlungs-/Admin-Freischaltungen; die Frontend-Seiten `/plans` und `/orders` bilden sie vollständig ab.
- [x] **Reale Zahlungsintegration (gopay)**: nativ WeChat Pay Native (PC-QR-Code) und Alipay Face-to-Face Pay (TradePrecreate); die Zahlungs-Callbacks laufen über `/api/payment/{wechat,alipay}/notify` und schließen Verifikation + Bestellaktivierung.
- [x] **Zahlungs-/Paketbetrieb-Einstellungen**: unter „Betriebseinstellungen“ neu hinzugekommen: „Paketbetrieb“ (Aufpreis-Upgrade vs. Additiv) und „Zahlung“ (WeChat / Alipay / Bank als drei unabhängige Schalter + Zertifikats-Upload + Benachrichtigungs-URL); Formulare erscheinen bedarfsweise.

### 🔄 In Umsetzung

- [ ] **Bessere Kanaldiagnose und intelligentere Routenoptimierung**: automatische Kühlung (`CooldownFilter`), Fallback-Degradierung (`FallbackFilter`) und automatische Deaktivierung bei niedriger Erfolgsrate (`monitor`) sind vorhanden; als Nächstes folgen ein eigenständiges Diagnose-Panel, Knoten-Ping und ein manueller Prüfablauf.
- [ ] Umfangreichere Nutzungsanalyse-Berichte und Exporte.
- [ ] Vervollständigung der Mehrsprachigkeit (i18n).

### 🔭 Geplant

- [ ] **Erweiterung der Zahlungskanäle**: Apple Pay, UnionPay, Stripe usw.; Unterstützung einer asynchronen Erstattungs-API + automatisierter Erstattungsbuchungen.
- [ ] **Online-Aufladung des Guthabens (Quota)**: Benutzer können ihr Konto im Bereich „Persönlich“ selbst aufladen; unabhängig von Abonnement-Paketen.
- [ ] **Finanzanbindung an gängige Plattformen**: Anbindung an gängige Finanz-/Abstimmungsplattformen; automatische Synchronisation von Aufladungen, Verbrauch und Erstattungen.
- [ ] **Warnmechanismus für niedrigen Token-Restbestand**: automatische Warnung bei niedrigem Token-Restbestand von Konto/Token; unterstützt Mehrkanal-Benachrichtigungen.
- [ ] **Log-Audit und Audit-Berichte**: vollständige Betriebs-Audit-Logs und visualisierte Audit-Berichte zur Erfüllung von Compliance-Anforderungen.
- [ ] **KI-Analyse**: intelligente Analyse und Empfehlungen zu Nutzung, Kosten und Kanalgesundheit auf Basis großer Sprachmodelle.
- [ ] Plug-in-Erweiterungsmechanismus.
- [ ] Unternehmensweites SSO / LDAP.
- [ ] Verbrauchswarnungen und Erweiterung der Benachrichtigungskanäle (DingTalk / Feishu / WeCom usw.).
- [ ] Kontinuierliche Aufnahme weiterer Modellplattformen.

> 💡 PRs und Issues sind willkommen – siehe [Issues](https://github.com/modelbus/one-api-pro/issues).

---

## Lizenz

[MIT-Lizenz](../LICENSE)
