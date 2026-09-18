---
title: 展示案例
description: "社区使用 One API Pro 构建的典型场景与示例。"
category: start
order: 2
---

# 展示案例

> 社区使用 One API Pro 构建的典型场景与示例。
> Typical scenarios and examples built with One API Pro.

下面是社区里较常见的落地方式。截图与镜像请见 `/docs` 仓库 `assets/` 目录。

Below are common deployment scenarios from the community. See the `assets/` folder of the `/docs` repo for screenshots and demo images.

## 1. 个人 OpenAI 网关 / Solo OpenAI Gateway

**场景 / Scenario**：家里 NAS 上跑一个实例，把 OpenAI、DeepSeek、Gemini 全部包成一个 `sk-` 出口，给 ChatGPT-Next-Web、Lobe-Chat、沉浸式翻译都用。

**做法 / Setup**：

- 单实例 SQLite，单 `Channel` 绑多个 Provider 备用；
- Token 不限模型、不限子网；
- 余额变动靠 `/api/log/self` 观察。

## 2. 团队内部分账 / Internal Team Billing

**场景 / Scenario**：10 人小团队，每人一个 user，套餐（`plans` + `user_plans`）+ 按量兜底。

**做法 / Setup**：

- 「月度 10 元套餐」+ 超出按量 `quota_per_unit` 扣费；
- `token.models` 限定模型组，`token.subnet` 限定公司出口 IP；
- `/api/log/self` 出账单，发票从管理后台 `Logs` 导。

## 3. SaaS 商业化分发 / SaaS Resale

**场景 / Scenario**：把 One API Pro 当成对外的产品，开通微信 / 支付宝收款、做套餐 + 充值混合。

**做法 / Setup**：

- 套餐：`basic / pro / team` 三档，`OrderTypePlanSubscription`（`TB` 前缀）；
- 充值：开启 `topup.*`，自定义金额 + 预设 chip，`OrderTypeTopup`（`TP` 前缀）；
- 升级差价：`OrderUpgradeModePriceDiff` 自动算差价，下 `UP` 前缀订单；
- 兑换码：地推拉新、`Redemption` 批量导出 CSV。

## 4. 多机房多活 / Multi-region Multi-active

**场景 / Scenario**：北京 / 上海 / 法兰克福三地机房，本地机房就近服务。

**做法 / Setup**：

- `CLUSTER_ENABLED=true`，三节点独立 MySQL + Redis；
- 域名 NS 三地 GeoDNS，OpenAI 兼容接口按地区落到最近节点；
- 同步范围见 [Cluster 概览](/zh/decentralization/overview#同步范围-sync-scope)。

## 5. 自定义 Provider 接入 / Custom Provider Onboarding

**场景 / Scenario**：企业内私有部署的 LLM 网关，需要把内部 API 包装成 OpenAI 兼容协议。

**做法 / Setup**：

- 参考 [新增 Provider](/zh/contribute/add-provider)；
- 在 `relay/adaptor/provider/<name>/` 下新建 4 个文件：`register.go` / `adaptor.go` / `main.go` / `constants.go`；
- 注册进 `registry`，重启即出现在「新建渠道」的下拉里。

## 截图与素材 / Screenshots

| 页面 / Page | 用途 / Use |
| --- | --- |
| Demo-Index.png | 控制台首页（视觉指标 + 趋势图） |
| Demo-Token.png | 令牌管理：剩余/已用、密钥列、可用模型 |
| Demo-Subscribe.png | 我的订阅 / 续费 |
| Demo-Plan.png | 套餐详情与差价升级 |
| Demo-cluster.png | Cluster 节点状态 |

> 截图持续更新中，欢迎 PR 提交你的部署截图。

下一步 / Next: [5 分钟快速体验](/zh/start/quick-tour) · [架构总览](/zh/start/architecture)。