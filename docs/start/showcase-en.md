---
title: Showcase
description: "Typical scenarios and examples built with One API Pro."
category: start
order: 2
---

# Showcase

> Typical scenarios and examples built with One API Pro.
> 社区使用 One API Pro 构建的典型场景与示例。

Below are common deployment scenarios from the community. See the `assets/` folder of the `/docs` repo for screenshots and demo images.

下面是社区里较常见的落地方式。截图与镜像请见 `/docs` 仓库 `assets/` 目录。

## 1. Solo OpenAI Gateway / 个人 OpenAI 网关

**Scenario**: run a single instance on a home NAS, wrap OpenAI, DeepSeek and Gemini behind one `sk-` egress used by ChatGPT-Next-Web, Lobe-Chat, immersive-translate, etc.

**Setup**:

- Single SQLite instance, one channel with multiple fallback providers.
- Tokens unrestricted in models and subnet.
- Track balance via `/api/log/self`.

## 2. Internal Team Billing / 团队内部分账

**Scenario**: 10-person team, one user each, Plan + pay-as-you-go fallback.

**Setup**:

- "10 CNY / month" plan + metered `quota_per_unit` for overage.
- `token.models` restricts the model group; `token.subnet` restricts the company egress IP.
- `/api/log/self` produces per-user statements; export invoices from the admin `Logs` view.

## 3. SaaS Resale / SaaS 商业化分发

**Scenario**: expose One API Pro as a paid product, enable WeChat / Alipay, mix Plan + top-up.

**Setup**:

- Plans: `basic / pro / team`, orders of `OrderTypePlanSubscription` (`TB` prefix).
- Top-up: `topup.*` settings with custom amount + preset chips, `OrderTypeTopup` (`TP` prefix).
- Upgrade differential: `OrderUpgradeModePriceDiff` auto-computes the delta and creates `UP`-prefixed orders.
- Redemptions: bulk-export CSVs for offline acquisition.

## 4. Multi-region Multi-active / 多机房多活

**Scenario**: Beijing / Shanghai / Frankfurt regions, serve each from the nearest node.

**Setup**:

- `CLUSTER_ENABLED=true`, three independent MySQL + Redis stacks.
- GeoDNS routes the OpenAI-compatible endpoint to the nearest node.
- See [Cluster Overview](/en/decentralization/overview#sync-scope) for the sync matrix.

## 5. Custom Provider Onboarding / 自定义 Provider 接入

**Scenario**: an enterprise-internal LLM gateway needs to be wrapped as an OpenAI-compatible protocol.

**Setup**:

- See [Add a Provider](/en/contribute/add-provider).
- Create four files under `relay/adaptor/provider/<name>/`: `register.go` / `adaptor.go` / `main.go` / `constants.go`.
- Register with the registry; after restart the provider appears in the "New Channel" dropdown.

## Screenshots / 截图与素材

| Page / 页面 | Use / 用途 |
| --- | --- |
| Demo-Index.png | Dashboard home (hero metrics + trend charts) |
| Demo-Token.png | Token management: remaining / used, key column, available models |
| Demo-Subscribe.png | My subscriptions / renewal |
| Demo-Plan.png | Plan detail and upgrade differential |
| Demo-cluster.png | Cluster node status |

> Screenshots are continuously updated; PRs with your deployment screenshots are welcome.

Next: [Quick Tour](/en/start/quick-tour) · [Architecture](/en/start/architecture).