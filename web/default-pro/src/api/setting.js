import api from './index'

// Admin: payment configuration + plan-operations configuration.
export const settingApi = {
  // Payment: GET the bundle, PUT per-method.
  getPayment: () => api.get('/api/setting/payment'),
  putPaymentMethod: (method, formData) =>
    api.put(`/api/setting/payment/${method}`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }),
  // Plan operations
  getPlan: () => api.get('/api/setting/plan'),
  putPlan: (data) => api.put('/api/setting/plan', data),
}

export default settingApi
