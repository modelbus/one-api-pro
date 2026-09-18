---
title: 支付通道
description: "微信 Native、支付宝当面付、银行转账、线下、免费通道；证书上传与回调地址配置。"
category: pricing
order: 5
---

# 支付通道

> 支付通道抽象为 `common/payment.Channel` 接口，按 `pay_method` 注册到全局 `registry`。每个通道从 `system_settings` 读自己的 `enabled` / `config` JSON，调用方通过 `payment.New(pay_method)` 拿到实例。

## 注册的通道 / Registered channels

`common/payment/payment.go::AnyChannelEnabled` 列举 `wechat` / `alipay` / `bank` 三个候选（`offline` / `free` 始终视为可用）；启动时由各文件 `init()` 注册：

| 通道 | pay_method | 实现 | 主要凭证 |
|---|---|---|---|
| 微信 Native | `OrderPayMethodWechat="wechat"` | `common/payment/wechat.go` | `app_id` / `mch_id` / `api_key` / `notify_url`；可选 `cert_file` / `key_file`（退款用） |
| 支付宝 当面付 | `OrderPayMethodAlipay="alipay"` | `common/payment/alipay.go` | `app_id` / `private_key` / `public_key` / `private_key_file` / `public_key_file` / `notify_url` / `gateway` |
| 银行转账 | `OrderPayMethodBank="bank"` | `common/payment/bank.go::bankChannel` | `account_name` / `account_no` / `bank_name` / `branch` / `notes`（`PrePay` 不可用，等管理员标记收款） |
| 线下支付 | `OrderPayMethodOffline="offline"` | `common/payment/bank.go::offlineChannel` | 同上；`IsEnabled()` 始终返回 true |
| 免费 / 管理员赠送 | `OrderPayMethodFree="free"` | `common/payment/bank.go::freeChannel` | 始终视为启用；`PrePay` 返回空对象 |

## 系统设置 / Settings

支付配置存 `system_settings` 表 `payment.*` 键（`model/system_setting.go`）：

| Key | 含义 |
|---|---|
| `payment.wechat.enabled` / `payment.wechat.config` | 微信总开关 / 凭证 JSON |
| `payment.alipay.enabled` / `payment.alipay.config` | 支付宝总开关 / 凭证 JSON |
| `payment.bank.enabled` / `payment.bank.config` | 银行转账总开关 / 凭证 JSON |

`IsEnabled()` 调 `payment.SettingsBool(key)` 读 `{"enabled": true|false}` 解析。

## 微信 Native / WeChat Native

`wechatChannel.PrePay`：

1. `IsEnabled()` 拒绝未启用
2. 加载 `wechatConfig`：`app_id` / `mch_id` / `api_key` 必填，缺一返回 `微信支付参数不完整`
3. `UnifiedOrder(ctx, BodyMap{body, out_trade_no, total_fee, spbill_create_ip, notify_url, trade_type=Native})`
4. `total_fee` = 元 × 100（`int64(amount*100 + 0.5)`）
5. 返回 `PrePayResult{PayURL=CodeURL, QRCode=CodeURL, ExpireAt=0, TradeNo=PrepayId}`

`VerifyNotify`：

- 解析 XML 为 `wechatNotifyXML`（`return_code` / `result_code` / `out_trade_no` / `transaction_id` / `total_fee` / `sign`）
- 按 ASCII 升序拼接非空字段 + `&key=<API_KEY>` 做 MD5 → 比对 `sign`
- 校验 `return_code == SUCCESS` 且 `result_code == SUCCESS`
- 回调金额 = `TotalFee / 100`（元）

`WechatNotify` 端点（`POST /api/payment/wechat/notify`，免鉴权）：

- 成功 → 返回 `<xml><return_code>SUCCESS</return_code>...</xml>`
- 失败 → 返回 `<xml><return_code>FAIL</return_code><return_msg>...</return_msg></xml>`（让微信重试）

## 支付宝 当面付 / Alipay Face-to-Face

`alipayChannel.PrePay`：

- `app_id` 必填；优先使用 `private_key_file` / `public_key_file`，否则读 inline `private_key` / `public_key`
- `gateway` 默认 `https://openapi.alipay.com/gateway.do`（生产），可填沙箱
- `TradePrecreate(ctx, "当面付", out_trade_no, total_amount)` 生成二维码字符串

`VerifyNotify` 用 `alipay.VerifySign`（RSA2 + alipay public key）校验 POST 字段后返回 `NotifyResult{OutTradeNo, TradeNo, Amount, Paid}`。

`AlipayNotify` 端点（`POST /api/payment/alipay/notify`，免鉴权）成功返回字面量 `success`，失败返回 `fail`。

## 银行 / 线下 / 免费 / Bank / Offline / Free

- `bankChannel`：`PrePay` 返回错误 `bank 支付未实现：等待管理员在后台标记收款`；无回调；订单保持 `pending`，由管理员通过 `PUT /api/order/:id` 手动标记收款
- `offlineChannel`：与 bank 类似，但 `IsEnabled()` 始终返回 true（默认开放）
- `freeChannel`：`PrePay` 返回空对象；`VerifyNotify` 直接返回 `Paid: true`。仅用于管理员 grant（`controller/subscription.go::AddSubscription`）

## 用户侧流程 / User flow

`payment.AnyChannelEnabled()` 决定 `CreatePlanOrder` / `CreateTopupOrder` / `PayMyOrder` 是否允许下单；返回 false 时统一返回 `系统尚未开通任何支付通道，请设置后开启支付`，且不落订单。

`buildPayInfo(pay_method, order_no, amount, "TBUS-"+package_name)` 返回给前端的对象：

```json
{
  "status": "success",          // or "warning"
  "pay_url": "weixin://wxpay/bizpayurl?pr=...",
  "qr_code": "weixin://wxpay/bizpayurl?pr=...",
  "expire_at": 0,
  "trade_no": "prepay_id_xxx",
  "note": "..."                 // 仅 bank
  "warning": "..."              // 仅 warning
}
```

`status=warning` 时表示通道未注册 / 未启用 / SDK 调用失败；订单仍然落库，管理员可手动处理。

## 异步回调分发 / Async notify dispatch

`controller/payment.go::processNotify`：

1. `payment.New(pay_method).VerifyNotify(body)` 校验签名
2. `model.GetOrderByOrderNo(notif.OutTradeNo)` 拉订单
3. `order.Amount > 0 && notif.Amount > 0 && notif.Amount != order.Amount` → `amount mismatch`
4. 按 `order.Type` 分发：
   - `OrderTypeTopup=2` → `model.ActivateTopupByOrder` 加 quota
   - 其它（套餐订单） → `model.ActivatePackageByOrder(order, OrderUpgradeModeStack)` 激活订阅

## 手动激活 / 测试通道 / Manual activation

`POST /api/payment/mock/notify`（Root，请求体 `{order_no, status}`）：

- `status=1` → 按订单类型调 `ActivateTopupByOrder` 或 `ActivatePackageByOrder`
- `status=3` → `model.MarkOrderRefunded`（状态翻转；不修改 `users.quota` / `user_plans`）

用于测试或对账后手动激活。

## 前端操作指南 / Frontend Guide

页面：`/setting/payment`（`web/default-pro/src/views/setting/PaymentSetting.vue`）。

- 每个通道一个开关 + 配置表单：
  - **微信**：`app_id` / `mch_id` / `api_key` / `notify_url` + 证书 / 私钥 PEM 上传按钮
  - **支付宝**：`app_id` / `gateway` / `notify_url` + 公钥 / 私钥 PEM 上传按钮
  - **银行**：`account_name` / `account_no` / `bank_name` / `notes`
- 切换开关即时调 `PUT /api/setting/payment/:method`；表单字段通过同一端点的 `config` 字段一并保存
- 用户侧 `/pricing` 或 `/topup` 弹窗：`GET /api/payment/status` 决定渲染哪些支付按钮（已启用的）

## 实现位置 / Implementation Pointers

| 关注点 | 位置 |
|---|---|
| 接口定义 | `common/payment/payment.go::Channel` |
| 微信实现 | `common/payment/wechat.go` |
| 支付宝实现 | `common/payment/alipay.go` |
| Bank / Offline / Free | `common/payment/bank.go` |
| 系统设置 key | `model/system_setting.go::SystemSettingKeyWechat* / Alipay* / Bank*` |
| 设置接口 | `controller/setting_payment.go`（参考前端字段） |
| 回调分发 | `controller/payment.go::processNotify` |
| 手动激活 | `controller/payment.go::MockPay` |
