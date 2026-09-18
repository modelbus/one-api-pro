---
title: 令牌 API
description: 令牌 API：令牌管理 (Token)
category: api
order: 14
---
## 6. 令牌管理 (Token)

### 6.1 获取所有令牌

**接口：** `GET /api/token/`

**权限：** User（仅返回当前用户的令牌）

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| p | int | 页码，默认 0 |
| order | string | 排序字段 |

**返回值：**

```json
{
  "success": true,
  "message": "",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "key": "sk-xxxxxxxx",
      "status": 1,
      "name": "my-token",
      "created_time": 1718000000,
      "accessed_time": 1718000000,
      "expired_time": -1,
      "remain_quota": 500000,
      "unlimited_quota": false,
      "used_quota": 100000,
      "models": null,
      "subnet": null,
      "updated_at": 1718000000
    }
  ]
}
```

**返回字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 令牌ID |
| user_id | int | 所属用户ID |
| key | string | 令牌密钥 |
| status | int | 状态：1=启用, 2=禁用, 3=过期, 4=耗尽 |
| name | string | 令牌名称 |
| created_time | int64 | 创建时间 |
| accessed_time | int64 | 最后访问时间 |
| expired_time | int64 | 过期时间，-1=永不过期 |
| remain_quota | int64 | 剩余额度 |
| unlimited_quota | bool | 是否无限额度 |
| used_quota | int64 | 已用额度 |
| models | string/null | 允许的模型（逗号分隔），null=全部 |
| subnet | string/null | 允许的子网 |

---

### 6.2 搜索令牌

**接口：** `GET /api/token/search`

**权限：** User

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| keyword | string | 搜索关键词 |

---

### 6.3 获取单个令牌

**接口：** `GET /api/token/:id`

**权限：** User

---

### 6.4 创建令牌

**接口：** `POST /api/token/`

**权限：** User

**请求体：**

```json
{
  "name": "my-token",
  "remain_quota": 500000,
  "expired_time": -1,
  "unlimited_quota": false,
  "models": null,
  "subnet": null
}
```

---

### 6.5 更新令牌

**接口：** `PUT /api/token/`

**权限：** User

与创建格式相同，加上 `id` 字段。查询参数 `status_only=1` 时仅更新状态。

---

### 6.6 删除令牌

**接口：** `DELETE /api/token/:id`

**权限：** User

---
