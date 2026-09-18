---
title: 新增支付通道
description: "如何接入一个新的支付通道。"
category: contribute
order: 5
---

# 新增支付通道

> 如何接入一个新的支付通道。
> How to onboard a new payment channel.

支付通道统一实现 `common/payment.Channel` 接口，通过 `init() → RegisterChannel` 自注册。下单链路复用 `controller/buildPayInfo`，回调走 `processNotify`，零侵入。

Payment channels implement `common/payment.Channel` and self-register via `init() → RegisterChannel`. The order flow reuses `controller/buildPayInfo`; the callback goes through `processNotify`.

## 1. 文件位置 / File Location

```text
common/payment/
├── payment.go            # Channel 接口 + RegisterChannel + New(payMethod)
├── wechat.go             # 已实现：微信 Native
├── alipay.go             # 已实现：支付宝
└── <your_channel>.go    # 新增
```

## 2. 接口 / The Channel Interface

`common/payment/payment.go::Channel`：

```go
type Channel interface {
    Name() string                                             // 匹配 OrderPayMethod<X>
    IsEnabled() (bool, error)                                 // 从 system_settings 读 enabled 开关
    PrePay(orderNo string, amount float64, subject string) (*PrePayResult, error)
    VerifyNotify(payload []byte) (*NotifyResult, error)
}
```

返回结构：

```go
type PrePayResult struct {
    PayURL   string `json:"pay_url"`
    QRCode   string `json:"qr_code"`
    ExpireAt int64  `json:"expire_at"`  // Unix seconds; 0 = unknown
    TradeNo  string `json:"trade_no"`   // 渠道侧预支付 ID
}

type NotifyResult struct {
    OutTradeNo string   // = 我们的 order_no
    TradeNo    string   // 渠道侧交易号
    Amount     float64
    Paid       bool
}
```

## 3. 最小实现 / Minimal Example

参考 `common/payment/wechat.go`：

Reference `common/payment/wechat.go`:

```go
// common/payment/wechat.go
package payment

import (
    "github.com/modelbus/one-api-pro/model"
)

type wechatChannel struct{}

func init() { RegisterChannel(&wechatChannel{}) }

func (*wechatChannel) Name() string                          { return model.OrderPayMethodWechat }
func (*wechatChannel) IsEnabled() (bool, error)              { return SettingsBool(model.SystemSettingKeyWechatEnabled) }

type wechatConfig struct {
    AppID     string `json:"app_id"`
    MchID     string `json:"mch_id"`
    APIKey    string `json:"api_key"`
    NotifyURL string `json:"notify_url"`
}

func (c *wechatChannel) loadConfig() (*wechatConfig, error) {
    var cfg wechatConfig
    if err := SettingsJSON(model.SystemSettingKeyWechatConfig, &cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}

func (c *wechatChannel) PrePay(orderNo string, amount float64, subject string) (*PrePayResult, error) {
    enabled, _ := c.IsEnabled()
    if !enabled {
        return nil, errors.New("channel disabled")
    }
    cfg, err := c.loadConfig()
    if err != nil {
        return nil, err
    }
    // TODO: build provider-specific request body
    //       and return pay_url / qr_code / trade_no
    return &PrePayResult{PayURL: "...", QRCode: "...", TradeNo: "..."}, nil
}

func (c *wechatChannel) VerifyNotify(payload []byte) (*NotifyResult, error) {
    // TODO: parse provider async notify payload
    //       verify signature, return OutTradeNo / TradeNo / Amount / Paid
    return &NotifyResult{OutTradeNo: "...", Paid: true}, nil
}
```

## 4. 系统设置约定 / System Settings Convention

每个通道建议两个 key：

Each channel conventionally uses two keys:

| Key | 用途 / Purpose | 格式 / Format |
| --- | --- | --- |
| `payment.<name>.enabled` | 启停开关 | `{"enabled": bool}` |
| `payment.<name>.config` | 渠道参数 | 自定义 JSON |

常量放在 `model/system_setting.go`：

Place constants in `model/system_setting.go`:

```go
const (
    SystemSettingKey<Name>Enabled = "payment.<name>.enabled"
    SystemSettingKey<Name>Config  = "payment.<name>.config"
)
```

`payment.go` 已提供 `SettingsBool` / `SettingsJSON` 辅助函数，直接用。

`payment.go` already provides `SettingsBool` / `SettingsJSON` helpers — use them directly.

## 5. 注册常量 / Register the Pay Method Constant

`model/order.go`：

```go
const OrderPayMethod<Name> = "<name>"
```

> **`Name()` 必须返回这个常量字符串**；`payment.New(payMethod)` 按字符串查注册表，名字不匹配会报「未注册的支付方式」。
> **`Name()` MUST return this constant string**; `payment.New(payMethod)` looks up the registry by name; mismatches produce "未注册的支付方式".

## 6. 下单链路 / Order Flow

```text
controller/topup.go::CreateTopupOrder
    └─► model.CreateTopupOrder
            └─► order = TP-prefixed order row
    └─► buildPayInfo(req.PayMethod, order.OrderNo, order.Amount, "余额充值")
            ├─► if payMethod == "bank": 返回占位 + note
            ├─► ch, _ := payment.New(payMethod)
            ├─► enabled, _ := ch.IsEnabled()
            └─► r, _ := ch.PrePay(orderNo, amount, subject)
```

回调链路：

Callback flow:

```text
POST /api/payment/notify/<name>     ← 上游回调
    └─► controller/payment.go::<Name>Notify
            └─► processNotify(c, model.OrderPayMethod<Name>)
                    ├─► ch.VerifyNotify(rawBody)
                    ├─► orders.GetByOrderNo(OutTradeNo)
                    └─► order.MarkOrderPaid(payMethod, tradeNo)
                            └─► order.Type 分发：套餐激活 / 充值入账
```

模板见 `controller/payment.go::WechatNotify`。

See `controller/payment.go::WechatNotify` for the template.

## 7. 前端配置 / Frontend Configuration

「支付配置」页（`PaymentSetting.vue`）按 `payment.<name>.enabled` / `payment.<name>.config` 自动渲染表单字段；新增字段请同步：

The "Payment Settings" page (`PaymentSetting.vue`) renders form fields based on `payment.<name>.enabled` / `payment.<name>.config` keys. When adding new fields, sync:

- `web/default-pro/src/views/setting/PaymentSetting.vue` 的 i18n 字段
- i18n keys under `web/default-pro/src/i18n/pages/setting.js` (namespace `settingPage.payment`)

## 8. 测试 / Tests

- 单测：抽离「签名校验」为可纯函数，覆正常 / 错误签名 / 过期 3 类场景。
- 集成：在「支付配置」填 sandbox 参数，「新建充值订单」测端到端。
- 异步回调：可用 ngrok / localtunnel 起一个公网地址，渠道后台填到 `notify_url`。

Unit tests: extract signature verification as a pure function; cover normal / wrong-signature / expired scenarios.

Integration: fill sandbox params in "Payment Settings", create a top-up order end-to-end.

Async notify: expose a public URL via ngrok / localtunnel, and put it into the provider console's `notify_url`.

## 9. 提交 / Commit

按 [提交规范 §3](/zh/contribute/commit-convention) 一文件一 commit，body 写明：

Per [Commit Convention §3](/en/contribute/commit-convention), one commit per file, body includes:

```text
feat(payment): 新增 <name> 支付通道

Add <name> payment channel.

- 实现 payment.Channel 接口：Name / IsEnabled / VerifyNotify
- PrePay 调 <name> SDK，返回 PayURL/QRCode/TradeNo
- 配置 key：payment.<name>.enabled / payment.<name>.config
- model.OrderPayMethod<Name> 常量 + 注册
- controller/payment.go::<Name>Notify 路由
- 单测覆盖：正常回调 / 错误签名 / 过期
```

下一步 / Next: [发版流程](/zh/contribute/release-process) · [新增 Provider](/zh/contribute/add-provider)。