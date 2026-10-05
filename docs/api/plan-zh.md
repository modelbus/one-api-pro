---
title: 套餐 API
description: "套餐 API：套餐管理 (Plan)"
category: api
order: 18
---
## 10. 套餐管理 (Plan)

### 10.1 获取所有套餐

**接口：** `GET /api/plan/`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| p | int | 页码，默认 0 |

**返回值：**

```json
{
  "success": true,
  "message": "",
  "data": [
    {
      "id": 1,
      "name": "基础套餐",
      "price": 99.00,
      "billing_type": "token",
      "virtual_amount": 100,
      "model_limits": "{\"gpt-4o\":{\"request_month\":1000,\"token_month\":50000000}}",
      "description": "基础套餐描述",
"features": ["API 调用 1000 次/月", "支持 GPT-4o"],
      "sort": 0,
      "status": 1,
      "duration_days": 30,
      "duration_text": "30天",
      "recommended": false,
      "created_time": 1718000000
    }
  ]
}
```

**返回字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 套餐ID |
| name | string | 套餐名称 |
| price | float64 | 价格（元） |
| billing_type | string | 计费维度：`token`（按 Token）或 `request`（按请求次数） |
| virtual_amount | float64 | 套餐虚拟余额总池（元）；`0` 表示不限额度 |
| model_limits | string | 模型限额配置JSON，key为模型名称，value为ModelLimitRule |
| description | string | 描述 |
| features | `array<string>` | 功能特性列表，每项一行展示在用户端套餐卡 |
| sort | int | 排序权重 |
| status | int | 状态：1=上架, 0=下架 |
| duration_days | int | 有效天数 |
| duration_text | string | 有效期显示文本 |
| recommended | bool | 是否推荐 |

> 单位口径：`virtual_amount` 存储层为「微元」（1 元 = 1_000_000），API 出入参统一为「元」，前端不需要做任何换算。

### 10.2 搜索套餐

**接口：** `GET /api/plan/search`

**权限：** Admin

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| keyword | string | 搜索关键词 |

### 10.3 获取套餐详情

**接口：** `GET /api/plan/:id`

**权限：** Admin

### 10.4 创建套餐

**接口：** `POST /api/plan/`

**权限：** Root

**请求体：**

```json
{
  "name": "基础套餐",
  "price": 99.00,
  "billing_type": "token",
  "virtual_amount": 100,
  "model_limits": "{\"gpt-4o\":{\"request_month\":1000,\"token_month\":50000000}}",
  "description": "基础套餐描述",
  "features": ["API 调用 1000 次/月", "支持 GPT-4o"],
  "sort": 0,
  "status": 1,
  "duration_days": 30,
  "duration_text": "30天",
  "recommended": false
}
```

**保存校验（不通过则返回 `success:false` 并拒绝落库）：**

1. `billing_type` 只能是 `token` 或 `request`；
2. `virtual_amount` 不能为负；
3. `model_limits` 必须是非空 JSON 对象，且可解析为 `map<string, ModelLimitRule>`；
4. 每条规则至少要配置一个与 `billing_type` 匹配的限额（`token` → `token_*`；`request` → `request_*`），否则视为无效配置。

> 校验的目的是堵住「配置写错 = 套餐免费不限量」：历史上 `model_limits` 写错或为空时，运行时会把套餐当成「不限制模型」直接放行。

### 10.5 更新套餐

**接口：** `PUT /api/plan/`

**权限：** Root

与创建格式相同，需包含 `id` 字段。

**`model_limits` 字段格式说明：**

key 为模型名称，value 为 `ModelLimitRule` 对象：

```

| 字段 | 类型 | 说明 |
|------|------|------|
| period_h | int | 滚动周期时长（小时），默认5 |
| request_period | int64 | 周期内最大请求数 |
| request_week | int64 | 周内最大请求数 |
| request_month | int64 | 月内最大请求数 |
| token_period | int64 | 周期内最大Token数 |
| token_week | int64 | 周内最大Token数 |
| token_month | int64 | 月内最大Token数 |

**未覆盖模型的处理：**

- 请求的模型不在 `model_limits` 中时，该套餐不适用于此请求，系统会跳过该套餐
- 若用户同时有余额，则回落至全局余额按量计费；余额不足由按量计费链路拒绝
- 平台未配置该模型的价格时，返回 422 `model_price_not_found`

**运行时失效条件（任一满足则该套餐整体不可用，请求回落全局余额）：**

- `model_limits` 为空或 JSON 非法（不再视为「不限制模型」）
- 套餐虚拟余额已耗尽：`user_plans.used_amount >= plans.virtual_amount`（`virtual_amount = 0` 表示不限额度）
- 任一计费窗口的加权用量达到 100%（`Σ 已用 × 100 / 限额 ≥ 100`）
- 套餐已过期（`end_time <= now`）

> 套餐消费会计入 `users.used_quota` / `request_count`，但**不扣减** `users.quota`；额度消耗记在 `user_plans.used_amount` 上。

### 10.6 删除套餐

**接口：** `DELETE /api/plan/:id`

**权限：** Root

