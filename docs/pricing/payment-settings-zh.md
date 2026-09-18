---
title: 支付通道设置
description: "微信 / 支付宝 / 银行转账三个通道的证书上传、回调地址与启用开关。"
category: pricing
order: 12
---

# 支付通道设置

> 在 `system_settings` 表上维护三种支付通道（微信、支付宝、银行转账）的启用开关、参数配置和证书文件路径。前端组件：`web/default-pro/src/views/setting/PaymentSetting.vue`。

## 接口一览

| Endpoint | Method | 鉴权 | 说明 |
|---|---|---|---|
| `/api/setting/payment` | `GET` | Root | 一次性返回三种通道的 `{ enabled, config, description, updated_at }` |
| `/api/setting/payment/:method` | `PUT` | Root | 保存单个通道；支持 `multipart/form-data` 上传证书 |

`:method` 取值 `wechat` / `alipay` / `bank`。

实现：`controller/setting_payment.go`。

## 请求体格式

PUT 请求采用 `multipart/form-data`：

| Form 字段 | 说明 |
|---|---|
| `config` | JSON 字符串：`{ "enabled": bool, "config": { ... } }` |
| `cert_file` | 仅 wechat：商户证书（`.pem`） |
| `key_file` | 仅 wechat：商户私钥（`.pem`） |
| `private_key_file` | 仅 alipay：应用私钥 |
| `public_key_file` | 仅 alipay：支付宝公钥 |

文件会被保存到 `data/payment/<method>/<basename>`，路径会被写回 `config.<field>` 键内。
注意：上传文件**不**会在 response 中返回路径，需要再次 GET 才能确认保存成功。

## 数据库

每种通道使用两条 `system_settings` 行：

| key | 用途 |
|---|---|
| `payment.<method>.enabled` | 仅 `{"enabled": bool}` 的最小 JSON |
| `payment.<method>.config` | 完整配置 JSON（含 `enabled` 与所有键） |

`category` 字段统一为 `payment`。GET 端点把两条行合并回 `enabled` + `config` 对象。

## 字段约定

| 方法 | 推荐字段 |
|---|---|
| `wechat` | `app_id` / `mch_id` / `api_key` / `notify_url` / `cert_file` / `key_file` |
| `alipay` | `app_id` / `gateway` / `notify_url` / `private_key` / `public_key` / `private_key_file` / `public_key_file` |
| `bank` | `account_name` / `account_no` / `bank_name` / `branch` / `notes` |

`app_id` / `mch_id` / `api_key` / `cert_file` / `key_file` 等键名由支付 SDK 直接读取；前后端字段名需严格一致。

## 前端操作指南

- 顶部三段式：「微信支付」/「支付宝」/「银行转账」；
- 每段顶部都是「启用开关」`a-switch`，切换时立即 PUT 仅 `enabled` 标志（其余 config 也一起带上，保证幂等）；
- 启用后才展开表单字段（按上面的字段约定）；
- 「上传证书」按钮（仅 wechat/alipay）以 `<a-upload custom-request>` 自定义请求方式直接 PUT 当前表单 + 文件；
- 「保存」按钮提交整张表单；银行通道没有证书，所以「保存」即可。

## 接口实现

| 关注点 | 位置 |
|---|---|
| Handler | `controller/setting_payment.go` |
| 设置读写 | `model/system_setting.go::GetSystemSetting` / `UpsertSystemSetting` |
| 证书落盘 | `controller/setting_payment.go::PutPaymentMethod`（`data/payment/<method>/`） |
| 支付通道实现 | `common/payment/`（`wechat.go` / `alipay.go` / `bank.go`） |
| 公开通道状态 | `controller/payment.go::GetPaymentStatus` |
| 路由 | `router/api.go` |
