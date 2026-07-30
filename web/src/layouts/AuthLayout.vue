<template>
  <div class="auth-layout">
    <div class="auth-bg"></div>
    <div class="auth-container">
      <div class="auth-card">
        <div class="auth-brand">
          <img v-if="statusStore.status?.logo" :src="statusStore.status.logo" class="brand-logo" />
          <h1 class="brand-name">{{ statusStore.status?.system_name || 'One Api Pro' }}</h1>
          <p class="brand-desc">企业级 AI API 网关</p>
        </div>
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useStatusStore } from '@/stores/status'

const statusStore = useStatusStore()

onMounted(async () => {
  if (!statusStore.loaded) await statusStore.fetchStatus()
})
</script>

<style scoped>
.auth-layout {
  min-height: 100vh;
  display: flex;
  position: relative;
}

.auth-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  opacity: 0.06;
}

.auth-container {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  position: relative;
  z-index: 1;
}

.auth-card {
  width: 420px;
  max-width: 100%;
  background: var(--color-bg-2);
  border-radius: 8px;
  padding: 48px 40px;
  box-shadow: 0 4px 24px rgb(0 0 0 / 8%);
}

.auth-brand {
  text-align: center;
  margin-bottom: 40px;
}

.brand-logo {
  width: 56px;
  height: 56px;
  margin-bottom: 16px;
}

.brand-name {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text-1);
  margin: 0 0 8px;
}

.brand-desc {
  font-size: 14px;
  color: var(--color-text-4);
  margin: 0;
}
</style>
