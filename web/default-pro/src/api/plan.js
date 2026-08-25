import api from './index'

export const planApi = {
  // Public
  list: () => api.get('/api/plan/list'),
  detail: (id) => api.get(`/api/plan/detail/${id}`),
  // User
  current: () => api.get('/api/plan/current'),
  // Admin
  all: (p = 0) => api.get('/api/plan/', { params: { p } }),
  search: (keyword) => api.get('/api/plan/search', { params: { keyword } }),
  get: (id) => api.get(`/api/plan/${id}`),
  add: (data) => api.post('/api/plan/', data),
  update: (data) => api.put('/api/plan/', data),
  delete: (id) => api.delete(`/api/plan/${id}`),
}

export default planApi
