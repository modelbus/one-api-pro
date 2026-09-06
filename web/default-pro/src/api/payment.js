import api from './index'

// Public user-facing payment endpoints.
export const paymentApi = {
  // Admin: mark an order paid (mock/dev only)
  mockPay: (data) => api.post('/api/payment/mock/notify', data),
  // Public: returns whether any payment channel is enabled, plus a list
  // of supported methods with their enabled flags.
  status: () => api.get('/api/payment/status'),
}

export default paymentApi
