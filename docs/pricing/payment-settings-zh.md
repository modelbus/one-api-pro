---
title: 支付通道设置
description: 在后台启用 / 配置各支付通道（微信 / 支付宝 / 银行转账）。
category: pricing
order: 10
---

# 支付通道设置

> 后台 → 支付设置。Root 可见。

## 三个通道

每种通道配置独立的开关与凭证。**至少启用一种**，否则用户无法付款。

### 微信支付

| 字段 | 填什么 | 哪里取 |
|---|---|---|
| `app_id` | 商户平台 AppID | 微信支付商户平台 → 账号中心 |
| `mch_id` | 商户号 | 同上 |
| `api_key` | v2 密钥 | 商户平台 → API 安全 → APIv2 密钥 |
| `notify_url` | 回调地址 | `https://your-domain.com/api/payment/wechat/notify` |
| 证书 / 私钥 | PEM 文件（退款用） | 商户平台 → API 安全 → API 证书 |

### 支付宝 当面付

| 字段 | 填什么 | 哪里取 |
|---|---|---|
| `app_id` | 应用 AppID | 支付宝开放平台 → 我的应用 |
| `private_key` / `public_key` | 应用私钥 / 支付宝公钥 | 同上 |
| 或 `private_key_file` / `public_key_file` | 上面两个的 PEM 文件路径 | 上传后系统自动保存路径 |
| `gateway` | 网关 | 默认 `https://openapi.alipay.com/gateway.do`（生产） |
| `notify_url` | 回调地址 | `https://your-domain.com/api/payment/alipay/notify` |

### 银行转账

| 字段 | 填什么 |
|---|---|
| `account_name` | 收款账户名（公司名） |
| `account_no` | 银行账号 |
| `bank_name` | 开户行 |
| `branch` | 支行 |
| `notes` | 让用户在转账备注里填订单号的提示文本 |

## 怎么启用 / 配置

后台 → 支付设置：

1. 选一个通道 → 打开「启用」开关
2. 表单字段展开 → 填各项
3. 微信 / 支付宝可上传证书（PEM 文件）
4. 点「保存」

切换开关即时生效。

## 证书 / 文件上传

微信和支付宝的 PEM 文件：

- 上传后系统存到 `data/payment/<method>/<basename>`
- 路径写回 `config` 字段（如 `cert_file: /app/data/payment/wechat/apiclient_cert.pem`）
- **不返回**在 response 里，需要重新 GET 才能确认

## 配置示例

### 微信生产

```json
{
  "enabled": true,
  "config": {
    "app_id": "wx0123456789abcdef",
    "mch_id": "1900000001",
    "api_key": "your-strong-api-key",
    "notify_url": "https://api.example.com/api/payment/wechat/notify"
  }
}
```

### 支付宝生产

```json
{
  "enabled": true,
  "config": {
    "app_id": "2021000123456789",
    "gateway": "https://openapi.alipay.com/gateway.do",
    "notify_url": "https://api.example.com/api/payment/alipay/notify"
  }
}
```

私钥建议上传 PEM 文件而不是粘贴。

### 银行转账

```json
{
  "enabled": true,
  "config": {
    "account_name": "XX 科技有限公司",
    "account_no": "6225 1234 5678 9012",
    "bank_name": "招商银行",
    "branch": "上海分行",
    "notes": "请在备注里填写订单号便于对账"
  }
}
```

## 怎么验证配置成功

1. 后台 → 支付设置 → 看「启用」开关是打开状态
2. 公开接口 `/api/payment/status` 应返回已启用的通道列表（无需登录）
3. 创建测试订单 → 用户能选到该通道 → 模拟支付验证完整链路

## 注意事项

- 修改后立即生效，但已下单的订单不受影响
- 凭证是**敏感数据**：避免提交到 git / 泄露给无关人员
- 切换支付通道时，新订单走新通道；老订单走原通道

## 常见问题

- **用户看不到任何支付方式**：三个通道都没启用，或全部未正确配置
- **支付成功但订单一直待支付**：检查 `notify_url` 是否能被公网访问；查后端日志
- **证书上传后提示错误**：检查 PEM 格式；私钥不要带密码（passphrase）
- **切换通道后老订单怎么办**：老订单走原通道；新订单走新通道

## 相关页面

- [支付通道](./payment)
- [充值业务配置](./topup-settings)
- [订单管理](./order-management)

## 相关 API

- `GET /api/setting/payment` — 读取所有通道配置（Root）
- `PUT /api/setting/payment/:method` — 保存单个通道（Root，支持文件上传）
- `GET /api/payment/status` — 公开查询已启用通道（无需登录）