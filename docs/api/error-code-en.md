---
title: Error Codes
description: "Quick lookup for error messages in the unified response format."
category: api
order: 20
---
# Error Codes

One API Pro does **not** publish a numeric `error_code` enum. Errors surface in one of two unified response shapes:

- Admin-plane `/api/*`: HTTP 200 + `{success:false, message:"..."}`
- Gateway-plane `/v1/*` (OpenAI / Anthropic compatible): HTTP `<status>` + `{error:{message, type, param, code}}`

This page enumerates every **known `message` string** so the client can localize, monitor, or fall back gracefully.

## 1. Auth / Permission Failures

| Trigger | message | HTTP | Source |
|---------|---------|------|--------|
| No session and no access token | `无权进行此操作，未登录且未提供 access token` | 401 | `middleware/auth.go` |
| Access token matches no user | `无权进行此操作，access token 无效` | 200 | `middleware/auth.go` |
| Current user disabled or blacklisted | `用户已被封禁` | 200 / 403 | `middleware/auth.go`, `TokenAuth` |
| Role below required tier (UserAuth/AdminAuth/RootAuth) | `无权进行此操作，权限不足` | 200 | `middleware/auth.go` |
| Bearer (`sk-...`) invalid / expired | `<err.Error() from model.ValidateUserToken>` | 401 | `middleware/auth.go::TokenAuth` |
| Bearer token subnet restriction violated | `该令牌只能在指定网段使用：<subnet>，当前 ip：<ip>` | 403 | `middleware/auth.go::TokenAuth` |
| Bearer-token user banned | `用户已被封禁` | 403 | `middleware/auth.go::TokenAuth` |
| Bearer token `models` whitelist violated | `该令牌无权使用模型：<model>` | 403 | `middleware/auth.go::TokenAuth` |
| Non-admin user appends `-<channel>` to API key | `普通用户不支持指定渠道` | 403 | `middleware/auth.go::TokenAuth` |
| Non-admin user hits `/v1/oneapi/proxy/:channelid/...` | `普通用户不支持指定渠道` | 403 | `middleware/auth.go::TokenAuth` |

## 2. Generic Parameter Errors

| Trigger | message | Source |
|---------|---------|--------|
| JSON / multipart parse failure | `无效的参数` | most controllers |
| Required ID missing | `ID不能为空` | `model_price` / `group_price` |
| Path ID not numeric / record not found | `订单不存在` / `套餐不存在` / `<err>` | most controllers |
| Empty `model_name` (model price) | `模型名称不能为空` | `controller/model_price.go` |
| Empty `group_name` (group price) | `分组名称不能为空` | `controller/group_price.go` |
| Empty `name` (redemption) | `名称不能为空` | `controller/redemption.go` |

## 3. User / Login

| Trigger | message | Source |
|---------|---------|--------|
| Username already taken | `用户名已存在` / `<err>` | `controller/user.go` |
| Wrong login password | `用户名或密码错误` | `controller/user.go` |
| Email already taken | `邮箱已存在` / `<err>` | `controller/user.go` |
| Wrong email verification code | `验证码错误` | `controller/user.go` |
| Root user tries to self-delete | `不能删除超级管理员` | `controller/user.go` |
| Affiliate code claimed by another user | `推广码已被其他用户使用` | `controller/user.go` |

## 4. Plan / Subscription / Upgrade

| Trigger | message | Source |
|---------|---------|--------|
| Plan unpublished | `套餐已下架` | `model/order_payment.go::CreatePlanOrder` |
| User already has an active plan with the same `sort` | `您已经订阅了同级别的套餐` | `model/order_payment.go::CreatePlanOrder` |
| Downgrade to a lower-`sort` plan | `不能降级到低级别套餐` | `model/order_payment.go::CreatePlanOrder` |
| `order.plan_id <= 0` reached activation | `order.plan_id 无效 (order_no=... plan_id=...)` | `model/order_payment.go::ActivatePackageByOrder` |
| `user_plan` snapshot / live plan missing | `套餐不存在` | `model/order_payment.go::ActivatePackageByOrder` |
| Cancel non-pending order | `只能取消待支付订单` | `model/order_payment.go::CancelOrder` |
| Pay a non-pending order | `订单当前不可支付` | `controller/order.go::PayMyOrder` |
| Access another user's order | `无权访问此订单` | `controller/order.go` |

## 5. Topup / Payment

| Trigger | message | Source |
|---------|---------|--------|
| No payment channel enabled | `系统尚未开通任何支付通道，请设置后开启支付` | `controller/payment.go::noPaymentEnabledMsg` |
| Callback signature / business check failed | `<err.Error()>` (WeChat returns XML `FAIL`, Alipay returns `fail`) | `controller/payment.go::processNotify` |
| Amount mismatch on callback | `amount mismatch` | `controller/payment.go::processNotify` |
| Callback not flagged as paid | `payment not marked paid by channel` | `controller/payment.go::processNotify` |
| Callback order missing | `<err from model.GetOrderByOrderNo>` | `controller/payment.go::processNotify` |
| Mock: empty `order_no` | `order_no 不能为空` | `controller/payment.go::MockPay` |
| Mock: order missing | `订单不存在` | `controller/payment.go::MockPay` |
| Mock: bad `status` value | `不支持的状态值（仅支持 1 或 3）` | `controller/payment.go::MockPay` |
| Mock: activation failed | `激活失败: <err>` | `controller/payment.go::MockPay` |
| Self-service: invalid amount | `amount 或 preset_amount 至少传一个` | `controller/topup.go::CreateTopupOrder` |
| Self-service: pay method not allowed | `自助充值仅支持 wechat / alipay / bank` / `自助下单仅支持 wechat / alipay / bank` | `controller/topup.go`, `controller/order.go` |
| Topup disabled | `充值功能未开启` | `model/topup.go::CreateTopupOrder` |
| Custom amount disabled | `未开启自定义金额` | `model/topup.go::ResolveTopupAmount` |
| Custom amount must be > 0 | `充值金额必须大于 0` | `model/topup.go::ResolveTopupAmount` |
| Preset amount not configured | `快捷金额未配置` | `model/topup.go::ResolveTopupAmount` |
| Save topup settings: invalid preset | `第 N 项金额必须大于 0` / `第 N 项额度不能为负数` | `model/topup.go::SaveTopupSettings` |
| Save topup settings: duplicate preset | `快捷金额重复：X.XX 元已存在` | `model/topup.go::SaveTopupSettings` |
| Save topup settings: bad exchange rate | `兑换比例必须大于 0` | `model/topup.go::SaveTopupSettings` |
| Save plan settings: bad upgrade mode | `upgrade_mode 必须是 price_diff 或 stack` | `controller/setting_payment.go::PutPlanSettings` |

## 6. Channel / Relay

| Trigger | message | Source |
|---------|---------|--------|
| No channel configured for the user's group | `分组 <group> 下未配置任何渠道` / `no channels available` | `controller/relay.go` |
| Upstream returns 401 / 403 | `<upstream error message>` passed through | `relay/handler` |
| Channel auto-disabled | `渠道已被自动禁用` / `<cooldown message>` | `relay/interceptor` |
| `/v1/images/edits`, `/v1/files`, `/v1/assistants/*`, `/v1/threads/*`, etc. not implemented | `{"error":{"message":"not_implemented","type":"one_api_error"}}` | `controller/relay.go::RelayNotImplemented` |
| Channel lookup failed | `channel not found` | `controller/relay.go::routeChannel` |

## 7. Relay Gateway Error Body (`/v1/*`)

```json
{
  "error": {
    "message": "<msg> (request id: <uuid>)",
    "type": "one_api_error",
    "param": "",
    "code": "<upstream error type or relay_internal_*>"
  }
}
```

`code` is typically the upstream channel's error type. Internal-only codes produced by the gateway include `read_response_body_failed`, `unmarshal_response_body_failed`, `marshal_response_body_failed`, `close_response_body_failed`, `do_request_failed`, `unknown_error` (see `relay/handler/error.go`).

## 8. Error Monitoring Tips

- Auth failures: exact-message match against §1 above.
- Business failures (`success:false` + message): group by message prefix:
  - `充值` / `订单` / `套餐` / `订阅` / `兑换` → business
  - `无权` / `用户已被封禁` / `普通用户不支持` → auth / permission
- Gateway failures (`/v1/*`): inspect `error.code` first, then `error.type == "one_api_error"`, then look for `(request id: ...)` at the end of `error.message` for triage.
- All error responses use either HTTP 200 (admin-plane) or HTTP `<upstream>` (gateway-plane); never infer business success/failure from the HTTP status code alone.