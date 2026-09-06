import api from './index'

export const orderApi = {
  // User self-service
  createPlan: (data) => api.post('/api/order/plan', data),
  myOrders: (type) => api.get('/api/order/self', { params: type ? { type } : {} }),
  myOrder: (id) => api.get(`/api/order/self/${id}`),
  cancelMyOrder: (id) => api.post(`/api/order/self/${id}/cancel`),
  payMyOrder: (id, data) => api.post(`/api/order/self/${id}/pay`, data || {}),
  // Admin
  all: (p = 0, type) => api.get('/api/order', { params: { p, ...(type ? { type } : {}) } }),
  search: (keyword, type) => api.get('/api/order/search', { params: { keyword, ...(type ? { type } : {}) } }),
  get: (id) => api.get(`/api/order/${id}`),
  markPaid: (id, data) => api.put(`/api/order/${id}`, data),
  delete: (id) => api.delete(`/api/order/${id}`),
}

export default orderApi
