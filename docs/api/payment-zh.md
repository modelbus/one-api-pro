---
title: 支付 API
description: "/api/payment/* 与回调通知。"
category: api
order: 10
---
# 支付 API

支付相关端点涵盖**微信支付**、**支付宝**、**银行转账**三类的异步回调、状态查询与管理员 Mock 通道；另含支付通道配置接口。所有回调端点均无需鉴权，靠渠道签名校验合法性。

## 端点一览

| 接口 | 方法 | 权限 | 说明 |
|------|------|------|------|
| `/api/payment/wechat/notify` | POST | 公开 | 微信支付异步回调 |
| `/api/payment/alipay/notify` | POST | 公开 | 支付宝异步回调 |
| `/api/payment/mock/notify` | POST | Root | 管理员 Mock 支付回调（用于测试/手工确认） |
| `/api/payment/status` | GET | 公开 | 各支付通道的启用状态 |
| `/api/setting/payment` | GET | Root | 读取支付配置 |
| `/api/setting/payment/:method` | PUT | Root | 更新支付配置（`wechat` / `alipay` / `bank`） |


## 公共约定：异步回调流程

所有异步回调（`wechat` / `alipay`）走同一处理流程（见 `controller/payment.go::processNotify`）：

1. 读取请求原始 body（不解析 JSON）。
2. 调对应 `payment.Channel.VerifyNotify` 校验签名。
3. 校验通过后取出 `NotifyResult{OutTradeNo, TradeNo, Amount, Paid}`。
4. 若 `order.amount > 0` 且 `notify.amount > 0`，必须金额相等，否则丢弃。
5. 按订单类型分发到激活逻辑：
   - `type = 1`（套餐订阅）→ `model.ActivatePackageByOrder`
   - `type = 2`（在线充值）→ `model.ActivateTopupByOrder`
6. 渠道专用的成功响应：
   - 微信：`<xml>...</xml>` 中 `return_code=SUCCESS`
   - 支付宝：明文 `success`

> 订单激活是幂等的：已支付的订单重复通知不会重复发额度/激活套餐。


## 1. 微信支付回调

**接口：** `POST /api/payment/wechat/notify`

**权限：** 公开（依赖微信签名校验）

**请求体：** 微信原生 form-urlencoded 回调（`xml` 解析前的内容），由 `common/payment/wechat.go::VerifyNotify` 校验。

**响应：**

| 场景 | Content-Type | Body |
|------|--------------|------|
| 校验通过 | `application/xml` | `<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>` |
| 校验/业务失败 | `application/xml` | `<xml><return_code><![CDATA[FAIL]]></return_code><return_msg><![CDATA[<err>]]></return_msg></xml>` |

> 按微信规范，业务失败时返回 `FAIL`，让微信侧触发重试。


## 2. 支付宝回调

**接口：** `POST /api/payment/alipay/notify`

**权限：** 公开（依赖支付宝 RSA2 签名校验）

**请求体：** 支付宝 form-urlencoded 回调，由 `common/payment/alipay.go::VerifyNotify` 校验。

**响应：**

| 场景 | Body |
|------|------|
| 校验通过 | 明文 `success` |
| 校验/业务失败 | 明文 `fail` |


## 3. Mock 支付回调（管理员）

**接口：** `POST /api/payment/mock/notify`

**权限：** Root

**说明：** 用于测试 / 手工操作，无需经过真实支付通道即可标记订单为已支付或已退款。订单按 type 自动分发到 `ActivateTopupByOrder` / `ActivatePackageByOrder`（默认升级模式 `stack`）。

**请求体：**

```json
{
  "order_no": "TB20250912153000123456",
  "status": 1
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| order_no | string | 是 | 目标订单号 |
| status | int | 是 | `1` = 支付成功；`3` = 已退款 |

**返回：**

```json
{ "success": true, "message": "订单已支付，余额已到账" }
```

或

```json
{ "success": true, "message": "订单已支付，套餐已激活" }
```

或

```json
{ "success": true, "message": "订单已标记为退款" }
```

**错误情况：**

| 场景 | message |
|------|---------|
| `order_no` 为空 | `order_no 不能为空` |
| 订单不存在 | `订单不存在` |
| `status` 不是 1 或 3 | `不支持的状态值（仅支持 1 或 3）` |
| 激活失败 | `激活失败: <err>` |


## 4. 支付通道状态

**接口：** `GET /api/payment/status`

**权限：** 公开

**说明：** 读取所有已注册支付通道的启用状态，供前端购买页决定是否显示下单 UI。

**返回示例：**

```json
{
  "success": true,
  "message": "",
  "data": {
    "any_enabled": true,
    "methods": [
      { "name": "wechat", "label": "微信支付", "enabled": true },
      { "name": "alipay", "label": "支付宝",   "enabled": true },
      { "name": "bank",   "label": "银行转账", "enabled": false }
    ]
  }
}
```

**字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| any_enabled | bool | 至少有一个通道已启用 |
| methods[].name | string | 通道 id（`wechat` / `alipay` / `bank`） |
| methods[].label | string | 人类可读名称 |
| methods[].enabled | bool | 该通道当前是否已启用 |


## 5. 读取支付配置

**接口：** `GET /api/setting/payment`

**权限：** Root

**返回示例：**

```json
{
  "success": true,
  "message": "",
  "data": {
    "wechat_enabled":  { "enabled": true,  "config": { ... }, "description": "微信支付开关",   "updated_at": 1718000000 },
    "wechat_config":   { "enabled": true,  "config": { ... }, "description": "微信支付参数",   "updated_at": 1718000000 },
    "alipay_enabled":  { "enabled": false, "config": {},      "description": "支付宝开关",     "updated_at": 0 },
    "alipay_config":   { "enabled": false, "config": {},      "description": "支付宝参数",     "updated_at": 0 },
    "bank_enabled":    { "enabled": false, "config": {},      "description": "银行转账开关",   "updated_at": 0 },
    "bank_config":     { "enabled": false, "config": {},      "description": "银行账户信息",   "updated_at": 0 }
  }
}
```

> 每个键对应 `system_settings` 表的一行（`payment.<method>.<enabled>` / `payment.<method>.config`）。


## 6. 更新支付配置

**接口：** `PUT /api/setting/payment/:method`

**权限：** Root

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| method | string | `wechat` / `alipay` / `bank` |

**请求体：** 支持 `application/json` 或 `multipart/form-data`（上传证书时必须用 multipart）。

JSON 形式：

```json
{
  "enabled": true,
  "config": {
    "app_id": "wx1234567890",
    "mch_id": "1900000000",
    "sign_type": "HMAC-SHA256"
  }
}
```

multipart 形式（字段说明）：

| 字段 | 必填 | 说明 |
|------|------|------|
| `config` | 是 | JSON 字符串，内容同上 |
| `cert_file` | 微信可选 | 微信支付 API 证书；保存路径回写到 `config.cert_file` |
| `key_file` | 微信可选 | 微信支付 API 私钥；保存路径回写到 `config.key_file` |
| `private_key_file` | 支付宝可选 | 商户 RSA 私钥；保存路径回写到 `config.private_key_file` |
| `public_key_file` | 支付宝可选 | 支付宝公钥；保存路径回写到 `config.public_key_file` |

> 上传的证书 / 密钥文件落地在 `data/payment/<method>/<basename>`，原文件路径将写回配置。

**返回：**

```json
{ "success": true, "message": "已保存" }
```

**错误情况：**

| 场景 | message |
|------|---------|
| method 非法 | `未知支付方式: <method>` |
| multipart 缺 `config` | `multipart 缺少 config 字段` |
| 配置 JSON 非法 | `config 不是合法 JSON` / `解析配置失败: <err>` |
| 创建/保存文件失败 | `创建目录失败: <err>` / `保存文件失败: <err>` |
| DB 写入失败 | `<err>` |