import axios from 'axios'

const api = axios.create({
  baseURL: '',
  withCredentials: true,
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    const { status } = error.response || {}
    if (status === 401) {
      window.location.href = '/login?expired=true'
    } else if (status === 429) {
      // rate limited - could show toast
    }
    return Promise.reject(error)
  }
)

export default api
