---
title: 错误码表
description: "统一响应格式下的错误信息速查。"
category: api
order: 20
---
# 错误码表

One API Pro 没有显式的 `error_code` 数字枚举。错误通过两种统一响应格式传递：

- 管理面 `/api/*`：HTTP 200 + `{success:false, message:"..."}`
- 网关面 `/v1/*`（含 OpenAI / Anthropic 兼容）：HTTP `<status>` + `{error:{message, type, param, code}}`

本文汇总所有**已知 message 文本**，便于前端做 i18n、监控打点或文案兜底。

## 1. 鉴权

| 触发场景 | message 原文 | HTTP | 触发位置 |
|----------|-------------|------|----------|
| 既无 Session 也无 Access Token | `无权进行此操作，未登录且未提供 access token` | 401 | `middleware/auth.go` |
| Access Token 查不到用户 | `无权进行此操作，access token 无效` | 200 | `middleware/auth.go` |
| 当前用户被禁用或黑名单 | `用户已被封禁` | 200 / 403 | `middleware/auth.go`、`TokenAuth` |
| 当前角色不足（UserAuth/AdminAuth/RootAuth） | `无权进行此操作，权限不足` | 200 | `middleware/auth.go` |
| Bearer Token (`sk-...`) 整体无效 / 过期 | `<model.ValidateUserToken 返回的 err.Error()>` | 401 | `middleware/auth.go::TokenAuth` |
| Bearer Token 绑定子网，IP 不在子网内 | `该令牌只能在指定网段使用：<subnet>，当前 ip：<ip>` | 403 | `middleware/auth.go::TokenAuth` |
| 用户被封禁（Bearer 路径） | `用户已被封禁` | 403 | `middleware/auth.go::TokenAuth` |
| Bearer Token 指定了模型，但不在 `models` 白名单 | `该令牌无权使用模型：<model>` | 403 | `middleware/auth.go::TokenAuth` |
| 普通用户在 API Key 后追加 `-<channel>` 指定渠道 | `普通用户不支持指定渠道` | 403 | `middleware/auth.go::TokenAuth` |
| 普通用户在 `/v1/oneapi/proxy/:channelid/...` 指定渠道 | `普通用户不支持指定渠道` | 403 | `middleware/auth.go::TokenAuth` |

## 2. 通用参数错误

| 触发场景 | message 原文 | 触发位置 |
|----------|-------------|----------|
| JSON| `无效的参数` | 大量 controller |
| 必填 ID 缺失 | `ID不能为空` | `model_price` / `group_price` |
| 路径上 ID 非数字 / 找不到记录 | `订单不存在` / `套餐不存在` / `<err>` | 多数 controller |
| `model_name` 为空（模型定价） | `模型名称不能为空` | `controller/model_price.go` |
| `group_name` 为空（分组折扣） | `分组名称不能为空` | `controller/group_price.go` |
| `key` 为空（兑换码） | `名称不能为空` | `controller/redemption.go` |

## 3. 用户

| 触发场景 | message 原文 | 触发位置 |
|----------|-------------|----------|
| 用户名已被注册 | `用户名已存在` / `<err>` | `controller/user.go` |
| 登录密码错误 | `用户名或密码错误` | `controller/user.go` |
| 邮箱已被注册 | `邮箱已存在` / `<err>` | `controller/user.go` |
| 邮箱验证码错误 | `验证码错误` | `controller/user.go` |
| Root 用户尝试自删 | `不能删除超级管理员` | `controller/user.go` |
| 推广码已被其他用户占用 | `推广码已被其他用户使用` | `controller/user.go` |

## 4. 套餐

| 触发场景 | message 原文 | 触发位置 |
|----------|-------------|----------|
| 套餐已下架 | `套餐已下架` | `model/order_payment.go::CreatePlanOrder` |
| 用户已有同 sort 的激活套餐 | `您已经订阅了同级别的套餐` | `model/order_payment.go::CreatePlanOrder` |
| 升级到更低 sort 的套餐 | `不能降级到低级别套餐` | `model/order_payment.go::CreatePlanOrder` |
| `order.plan_id <= 0` 触发激活 | `order.plan_id 无效 (order_no=... plan_id=...)` | `model/order_payment.go::ActivatePackageByOrder` |
| `user_plan` 找不到对应 plan 快照 / 套餐 | `套餐不存在` | `model/order_payment.go::ActivatePackageByOrder` |
| 取消非 pending 订单 | `只能取消待支付订单` | `model/order_payment.go::CancelOrder` |
| 支付非 pending 订单 | `订单当前不可支付` | `controller/order.go::PayMyOrder` |
| 访问他人订单 | `无权访问此订单` | `controller/order.go` |

## 5. 充值

| 触发场景 | message 原文 | 触发位置 |
|----------|-------------|----------|
| 任意支付通道未启用 | `系统尚未开通任何支付通道，请设置后开启支付` | `controller/payment.go::noPaymentEnabledMsg` |
| 支付回调：渠道签名校验失败 / 业务校验失败 | `<err.Error()>`（微信返回 XML `FAIL`，支付宝返回 `fail`） | `controller/payment.go::processNotify` |
| 支付回调：金额与订单不一致 | `amount mismatch` | `controller/payment.go::processNotify` |
| 支付回调：未标记 paid | `payment not marked paid by channel` | `controller/payment.go::processNotify` |
| 支付回调：订单找不到 | `<model.GetOrderByOrderNo 返回 err>` | `controller/payment.go::processNotify` |
| Mock 通知：`order_no` 为空 | `order_no 不能为空` | `controller/payment.go::MockPay` |
| Mock 通知：订单不存在 | `订单不存在` | `controller/payment.go::MockPay` |
| Mock 通知：状态值非法 | `不支持的状态值（仅支持 1 或 3）` | `controller/payment.go::MockPay` |
| Mock 通知：激活失败 | `激活失败: <err>` | `controller/payment.go::MockPay` |
| 自助下单：金额非法 | `amount 或 preset_amount 至少传一个` | `controller/topup.go::CreateTopupOrder` |
| 自助下单：支付方式不允许 | `自助充值仅支持 wechat| `controller/topup.go`、`controller/order.go` |
| 充值功能未开启 | `充值功能未开启` | `model/topup.go::CreateTopupOrder` |
| 自定义金额未开启 | `未开启自定义金额` | `model/topup.go::ResolveTopupAmount` |
| 自定义金额必须 > 0 | `充值金额必须大于 0` | `model/topup.go::ResolveTopupAmount` |
| 快捷金额未配置 | `快捷金额未配置` | `model/topup.go::ResolveTopupAmount` |
| 充值设置：金额非法 | `第 N 项金额必须大于 0` / `第 N 项额度不能为负数` | `model/topup.go::SaveTopupSettings` |
| 充值设置：快捷金额重复 | `快捷金额重复：X.XX 元已存在` | `model/topup.go::SaveTopupSettings` |
| 充值设置：兑换比例非法 | `兑换比例必须大于 0` | `model/topup.go::SaveTopupSettings` |
| 套餐设置：升级模式非法 | `upgrade_mode 必须是 price_diff 或 stack` | `controller/setting_payment.go::PutPlanSettings` |

## 6. 渠道

| 触发场景 | message 原文 | 触发位置 |
|----------|-------------|----------|
| 用户分组下找不到任何渠道 | `分组 <group> 下未配置任何渠道` / `no channels available` | `controller/relay.go` |
| 渠道上游返回 401 / 403 | `<upstream error message>` 透传 | `relay/handler` |
| 渠道因错误被自动禁用 | `渠道已被自动禁用` / `<cooldown message>` | `relay/interceptor` |
| `/v1/images/edits`、`/v1/files`、`/v1/assistants/*`、`/v1/threads/*` 等未实现 | `{"error":{"message":"not_implemented","type":"one_api_error"}}` | `controller/relay.go::RelayNotImplemented` |
| 渠道找不到 | `channel not found` | `controller/relay.go::routeChannel` |

## 7. Relay 网关错误体（`/v1/*`）

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

`code` 通常是上游渠道返回的错误类型；网关自身产生的代码包括 `read_response_body_failed` / `unmarshal_response_body_failed` / `marshal_response_body_failed` / `close_response_body_failed` / `do_request_failed` / `unknown_error`（见 `relay/handler/error.go`）。

## 8. 错误监控建议

- 鉴权失败：message 完全匹配上表第 1 节即可统计归类。
- 业务失败（`success:false` + message）：按 message 前缀做分组，例如：
  - `充值` / `订单` / `套餐` / `订阅` / `兑换` → 业务侧
  - `无权` / `用户已被封禁` / `普通用户不支持` → 鉴权 / 权限
- 网关失败（`/v1/*`）：优先看 `error.code`，其次 `error.type = "one_api_error"`，最后看 `error.message` 末尾是否带 `(request id: ...)` 以便排障。
- 所有错误响应都使用 HTTP 200（管理面）或 HTTP `<upstream>`（网关面），不要假设 HTTP status 反映业务成败。
