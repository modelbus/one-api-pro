import api from './index'

// User-facing payment endpoints (notify URLs are public; the rest are
// admin-only).
export const paymentApi = {
  // Admin: mark an order paid (mock/dev only)
  mockPay: (data) => api.post('/api/payment/mock/notify', data),
}

export default paymentApi
