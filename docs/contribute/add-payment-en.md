---
title: Add a Payment Channel
description: "How to onboard a new payment channel."
category: contribute
order: 5
---

# Add a Payment Channel

> How to onboard a new payment channel.
> 如何接入一个新的支付通道。

Payment channels implement `common/payment.Channel` and self-register via `init() → RegisterChannel`. The order flow reuses `controller/buildPayInfo`; the callback goes through `processNotify`.

支付通道统一实现 `common/payment.Channel` 接口，通过 `init() → RegisterChannel` 自注册。下单链路复用 `controller/buildPayInfo`，回调走 `processNotify`，零侵入。

## 1. File Location / 文件位置

```text
common/payment/
├── payment.go            # Channel interface + RegisterChannel + New(payMethod)
├── wechat.go             # already implemented: WeChat Native
├── alipay.go             # already implemented: Alipay
└── <your_channel>.go    # new
```

## 2. The Channel Interface / 接口

`common/payment/payment.go::Channel`:

```go
type Channel interface {
    Name() string                                             // matches OrderPayMethod<X>
    IsEnabled() (bool, error)                                 // reads the enabled flag from system_settings
    PrePay(orderNo string, amount float64, subject string) (*PrePayResult, error)
    VerifyNotify(payload []byte) (*NotifyResult, error)
}
```

Return shapes / 返回结构:

```go
type PrePayResult struct {
    PayURL   string `json:"pay_url"`
    QRCode   string `json:"qr_code"`
    ExpireAt int64  `json:"expire_at"`  // Unix seconds; 0 = unknown
    TradeNo  string `json:"trade_no"`   // provider-side prepay id
}

type NotifyResult struct {
    OutTradeNo string   // = our order_no
    TradeNo    string   // provider-side transaction id
    Amount     float64
    Paid       bool
}
```

## 3. Minimal Example / 最小实现

Reference `common/payment/wechat.go`:

参考 `common/payment/wechat.go`：

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

## 4. System Settings Convention / 系统设置约定

Each channel conventionally uses two keys:

每个通道建议两个 key：

| Key | Purpose / 用途 | Format / 格式 |
| --- | --- | --- |
| `payment.<name>.enabled` | on/off switch | `{"enabled": bool}` |
| `payment.<name>.config` | channel params | custom JSON |

Place constants in `model/system_setting.go`:

常量放在 `model/system_setting.go`：

```go
const (
    SystemSettingKey<Name>Enabled = "payment.<name>.enabled"
    SystemSettingKey<Name>Config  = "payment.<name>.config"
)
```

`payment.go` already provides `SettingsBool` / `SettingsJSON` helpers — use them directly.

`payment.go` 已提供 `SettingsBool` / `SettingsJSON` 辅助函数，直接用。

## 5. Register the Pay Method Constant / 注册常量

In `model/order.go`:

`model/order.go`：

```go
const OrderPayMethod<Name> = "<name>"
```

> **`Name()` MUST return this constant string**; `payment.New(payMethod)` looks up the registry by name; mismatches produce "未注册的支付方式".
> **`Name()` 必须返回这个常量字符串**；名字不匹配会报「未注册的支付方式」。

## 6. Order Flow / 下单链路

```text
controller/topup.go::CreateTopupOrder
    └─► model.CreateTopupOrder
            └─► order = TP-prefixed order row
    └─► buildPayInfo(req.PayMethod, order.OrderNo, order.Amount, "余额充值")
            ├─► if payMethod == "bank": return placeholder + note
            ├─► ch, _ := payment.New(payMethod)
            ├─► enabled, _ := ch.IsEnabled()
            └─► r, _ := ch.PrePay(orderNo, amount, subject)
```

Callback flow / 回调链路:

```text
POST /api/payment/notify/<name>     ← upstream callback
    └─► controller/payment.go::<Name>Notify
            └─► processNotify(c, model.OrderPayMethod<Name>)
                    ├─► ch.VerifyNotify(rawBody)
                    ├─► orders.GetByOrderNo(OutTradeNo)
                    └─► order.MarkOrderPaid(payMethod, tradeNo)
                            └─► dispatch by order.Type: activate plan / credit top-up
```

See `controller/payment.go::WechatNotify` for the template.

模板见 `controller/payment.go::WechatNotify`。

## 7. Frontend Configuration / 前端配置

The "Payment Settings" page (`PaymentSetting.vue`) renders form fields based on `payment.<name>.enabled` / `payment.<name>.config` keys. When adding new fields, sync:

「支付配置」页（`PaymentSetting.vue`）按 `payment.<name>.enabled` / `payment.<name>.config` 自动渲染表单字段；新增字段请同步：

- i18n fields in `web/default-pro/src/views/setting/PaymentSetting.vue`
- i18n keys under `web/default-pro/src/i18n/pages/setting.js` (namespace `settingPage.payment`)

## 8. Tests / 测试

- Unit: extract signature verification as a pure function; cover normal / wrong-signature / expired scenarios.
- Integration: fill sandbox params in "Payment Settings", create a top-up order end-to-end.
- Async notify: expose a public URL via ngrok / localtunnel, and put it into the provider console's `notify_url`.

单测：抽离「签名校验」为可纯函数，覆正常 / 错误签名 / 过期 3 类场景。
集成：在「支付配置」填 sandbox 参数，「新建充值订单」测端到端。
异步回调：可用 ngrok / localtunnel 起一个公网地址，渠道后台填到 `notify_url`。

## 9. Commit / 提交

Per [Commit Convention §3](/en/contribute/commit-convention), one commit per file, body includes:

按 [提交规范 §3](/zh/contribute/commit-convention) 一文件一 commit，body 写明：

```text
feat(payment): add <name> payment channel

Add <name> payment channel.

- Implement payment.Channel: Name / IsEnabled / VerifyNotify
- PrePay calls <name> SDK, returns PayURL/QRCode/TradeNo
- Config keys: payment.<name>.enabled / payment.<name>.config
- model.OrderPayMethod<Name> constant + registration
- controller/payment.go::<Name>Notify route
- Unit tests: normal callback / wrong signature / expired
```

Next: [Release Process](/en/contribute/release-process) · [Add a Provider](/en/contribute/add-provider).