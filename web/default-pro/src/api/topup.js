// 在线充值（topup）API 模块
// 版本: v0.0.10
// 日期: 2026-09-06
// 作者: opencode
//
// 说明：命名沿用 topup（与后端 LogTypeTopup/AdminTopUp 等保持一致），
//       UI 文案对外显示为"充值"。

import api from './index'

export const topupApi = {
  // 用户自助下单：amount 或 preset_amount 至少传一个
  // 返回结构与 CreatePlanOrder 对齐：{ order, amount, bonus_quota, pay }
  createOrder: (data) => api.post('/api/topup/order', data),
}

export default topupApi
