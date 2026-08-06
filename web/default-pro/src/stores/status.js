import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api'

export const useStatusStore = defineStore('status', () => {
  const status = ref({})
  const loaded = ref(false)

  async function fetchStatus() {
    try {
      const { data } = await api.get('/api/status')
      if (data.success) {
        status.value = data.data
        loaded.value = true
        document.title = data.data.system_name || 'One Api Pro'
      }
    } catch (e) {
      // ignore
    }
  }

  return { status, loaded, fetchStatus }
})
