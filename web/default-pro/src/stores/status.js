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
        // 站点信息不参与 Tab 标题：标题由路由 meta 驱动（见 router/index.js），
        // 否则刷新后会被后端 system_name 的默认值（如 one-api）覆盖。
        // system_name must not drive document.title; the router owns the tab title.
      }
    } catch (e) {
      // ignore
    }
  }

  return { status, loaded, fetchStatus }
})
