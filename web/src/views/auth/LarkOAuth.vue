<template>
  <div class="oauth-page">
    <a-spin :loading="true" tip="飞书登录中..." :style="{ width: '100%' }">
      <div style="padding: 60px 0"></div>
    </a-spin>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

onMounted(async () => {
  const { code, state } = route.query
  try {
    const { data } = await api.get('/api/oauth/lark', { params: { code, state } })
    if (data.success) {
      authStore.user = data.data
      localStorage.setItem('user', JSON.stringify(data.data))
      router.push(state && state.startsWith('bind') ? '/setting' : '/')
    }
  } catch (e) {
    router.push('/login?error=oauth')
  }
})
</script>

<style scoped>
.oauth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-fill-2);
}
</style>
