import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(JSON.parse(localStorage.getItem('user')) || null)

  const isLoggedIn = computed(() => !!user.value)
  const isAdmin = computed(() => user.value && user.value.role >= 10)
  const isRoot = computed(() => user.value && user.value.role >= 100)

  async function login(username, password) {
    const { data } = await api.post('/api/user/login', { username, password })
    if (data.success) {
      user.value = data.data
      localStorage.setItem('user', JSON.stringify(data.data))
      return data.data
    }
    throw new Error(data.message)
  }

  async function register(form) {
    const { data } = await api.post('/api/user/register', form)
    if (data.success) return data
    throw new Error(data.message)
  }

  async function logout() {
    try {
      await api.get('/api/user/logout')
    } catch (e) {
      // ignore
    }
    user.value = null
    localStorage.removeItem('user')
  }

  function loadUser() {
    const stored = localStorage.getItem('user')
    if (stored) {
      try {
        user.value = JSON.parse(stored)
      } catch (e) {
        user.value = null
        localStorage.removeItem('user')
      }
    }
  }

  return { user, isLoggedIn, isAdmin, isRoot, login, register, logout, loadUser }
})
