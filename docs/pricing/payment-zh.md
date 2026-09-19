---
title: 支付通道
description: 微信 / 支付宝 / 银行转账的接入方式与后台配置。
category: pricing
order: 5
---

# 支付通道

> 用户付款走哪个通道？怎么配置？怎么处理回调？

## 支持的通道

| 通道 | pay_method | 适用场景 |
|---|---|---|
| **微信 Native** | `wechat` | 用户用微信扫一扫付款 |
| **支付宝 当面付** | `alipay` | 用户用支付宝扫一扫付款 |
| **银行转账** | `bank` | 用户用银行 / 公司网银转账，管理员手工确认 |
| **线下支付** | `offline` | 不走线上系统，管理员手工 |
| **免费 / 赠送** | `free` | 仅用于管理员手动开通订阅 |

默认开启前三种里至少一种，否则用户在结算页看不到支付按钮。

## 在哪里看到

- **后台 → 支付设置**：管理员配置每个通道
- **用户侧**：套餐下单、充值页面的「支付方式」下拉

## 微信 Native（推荐 C 端）

### 申请

1. 微信支付商户平台 → 产品中心 → 申请「Native 支付」
2. 拿到 `app_id` / `mch_id` / `api_key`（v2）或「商户 API 证书」（v3）
3. 配置回调地址：`https://your-domain.com/api/payment/wechat/notify`

### 后台配置

后台 → 支付设置 → 微信：

| 项 | 填什么 |
|---|---|
| 启用 | 开关 |
| `app_id` | 商户平台 → 账号中心 |
| `mch_id` | 商户平台 → 账号中心 |
| `api_key` | 商户平台 → API 安全 → v2 密钥 |
| `notify_url` | `https://your-domain.com/api/payment/wechat/notify` |
| 证书 / 私钥（PEM） | 退款用；选填 |

## 支付宝 当面付（推荐 C 端）

### 申请

1. 支付宝开放平台 → 创建应用 → 签约「当面付」
2. 拿到 `app_id` / 应用公钥 / 应用私钥
3. 配置回调地址：`https://your-domain.com/api/payment/alipay/notify`

### 后台配置

后台 → 支付设置 → 支付宝：

| 项 | 填什么 |
|---|---|
| 启用 | 开关 |
| `app_id` | 开放平台 → 我的应用 |
| `private_key` / `public_key` | 应用公钥 / 私钥（PEM 内容） |
| 或 `private_key_file` / `public_key_file` | 文件路径 |
| `gateway` | 默认 `https://openapi.alipay.com/gateway.do`（生产） |
| `notify_url` | `https://your-domain.com/api/payment/alipay/notify` |

## 银行转账（B 端 / 公司用户）

适合：金额大、对账要求高、不适合线上扫码的公司客户。

### 流程

1. 用户下单 → 订单状态 = 待支付
2. 用户用银行 / 公司网银转账到你提供的账户
3. 管理员在后台「订单」详情页手动「标记已支付」
4. 系统激活套餐 / 充值

### 后台配置

后台 → 支付设置 → 银行转账：

| 项 | 填什么 |
|---|---|
| 启用 | 开关 |
| `account_name` | 收款账户名（公司名） |
| `account_no` | 银行账号 |
| `bank_name` | 开户行 |
| `branch` | 支行 |
| `notes` | 备注（让用户备注订单号） |

## 用户侧支付流程

```
用户在套餐页下单
    ↓
弹窗显示可用支付方式（已启用的）
    ↓
用户选微信 → 显示二维码 → 扫码支付
或选银行 → 显示账号 + 备注订单号 → 用户转账
    ↓
微信/支付宝自动回调 → 订单状态 = 已支付 → 激活套餐
银行转账 → 管理员后台手动标记已支付 → 激活套餐
```

## 异步回调

第三方支付完成后会异步通知 One API Pro，验证签名 → 更新订单 → 激活套餐：

- `POST /api/payment/wechat/notify`（微信）
- `POST /api/payment/alipay/notify`（支付宝）

这两个端点**免鉴权**，由支付平台直接调用。要确保能从公网访问。

## 手动激活（调试 / 对账）

`POST /api/payment/mock/notify`（Root）：传 `{order_no, status}`：

- `status=1` → 强制激活（按订单类型加 quota 或开订阅）
- `status=3` → 标记为已退款

用于：测试、对账修正、回调失败补救。

## 常见问题

- **支付方式按钮不显示**：检查后台是否至少启用了一种通道
- **回调一直不来**：检查 `notify_url` 是否能从公网访问；检查防火墙 / 反代是否拦截
- **金额校验失败**：检查订单的 `amount` 与 `pay_method` 返回的金额是否一致（避免汇率 / 优惠造成的不一致）
- **银行转账收款后怎么操作**：管理员去「订单」详情页手动点「标记已支付」

## 相关页面

- [支付配置（后台）](./payment-settings)
- [订单管理（后台）](./order-management)
- [我的订单（用户）](../user/orders)

## 相关 API

- `GET /api/payment/status` — 当前已启用的支付方式（Public）
- `POST /api/payment/wechat/notify` — 微信回调
- `POST /api/payment/alipay/notify` — 支付宝回调
- `POST /api/payment/mock/notify` — 手动激活（Root）
- `GET /api/setting/payment` — 读取支付配置（Root）
- `PUT /api/setting/payment/:method` — 更新单个通道配置（Root）