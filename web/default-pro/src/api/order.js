// order 订单 API 模块
// Admin orders API module
// 版本: v0.0.13
// 日期: 2026-09-08
import api from './index'

// 仅追加 0/空值时透传的字段，避免污染 query。
// Only append filter params when truthy to keep query clean.
function buildAdminParams(extras = {}) {
  const out = {}
  for (const [k, v] of Object.entries(extras)) {
    if (v === 0 || (v !== '' && v !== null && v !== undefined)) {
      out[k] = v
    }
  }
  return out
}

export const orderApi = {
  // 用户自助
  // User self-service
  createPlan: (data) => api.post('/api/order/plan', data),
  myOrders: (type) => api.get('/api/order/self', { params: type ? { type } : {} }),
  myOrder: (id) => api.get(`/api/order/self/${id}`),
  cancelMyOrder: (id) => api.post(`/api/order/self/${id}/cancel`),
  payMyOrder: (id, data) => api.post(`/api/order/self/${id}/pay`, data || {}),

  // 管理员
  // Admin
  // list(params) 接收 { p, page_size, type, status, source, user_id, plan_id, keyword }
  list: (params = {}) => api.get('/api/order', { params: { p: 0, ...params } }),
  search: (keyword, type, status, source, userId, planId) =>
    api.get('/api/order/search', {
      params: {
        keyword,
        ...buildAdminParams({ type, status, source, user_id: userId, plan_id: planId }),
      },
    }),
  get: (id) => api.get(`/api/order/${id}`),
  markPaid: (id, data) => api.put(`/api/order/${id}`, data),
  delete: (id) => api.delete(`/api/order/${id}`),
}

export default orderApi