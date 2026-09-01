<template>
  <div class="auth-layout">
    <div class="auth-background" aria-hidden="true">
      <span class="auth-orb auth-orb-primary"></span>
      <span class="auth-orb auth-orb-secondary"></span>
      <span class="auth-orb auth-orb-tertiary"></span>
    </div>
    <div class="auth-container">
      <div class="auth-card">
        <div class="auth-brand">
          <img :src="logoPng" class="brand-logo" alt="ONE-API-PRO" />
        </div>
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useStatusStore } from '@/stores/status'
import logoPng from '@/assets/logo.png'

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
  overflow: hidden;
  isolation: isolate;
}

.auth-background {
  position: absolute;
  inset: 0;
  z-index: 0;
  overflow: hidden;
  background: linear-gradient(135deg, #f5f9ff 0%, #f4f7ff 48%, #fbf8ff 100%);
}

.auth-background::before {
  position: absolute;
  inset: 0;
  content: '';
  background-image: linear-gradient(rgba(110, 139, 205, 0.06) 1px, transparent 1px), linear-gradient(90deg, rgba(110, 139, 205, 0.06) 1px, transparent 1px);
  background-size: 48px 48px;
  mask-image: linear-gradient(to bottom, transparent, #000 24%, #000 76%, transparent);
}

.auth-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(2px);
  opacity: 0.7;
}

.auth-orb-primary {
  top: -180px;
  right: -140px;
  width: 520px;
  height: 520px;
  background: radial-gradient(circle, rgba(72, 125, 255, 0.2), rgba(72, 125, 255, 0) 70%);
}

.auth-orb-secondary {
  bottom: -240px;
  left: -180px;
  width: 620px;
  height: 620px;
  background: radial-gradient(circle, rgba(145, 92, 255, 0.16), rgba(145, 92, 255, 0) 70%);
}

.auth-orb-tertiary {
  top: 42%;
  left: 20%;
  width: 180px;
  height: 180px;
  background: radial-gradient(circle, rgba(46, 212, 191, 0.1), rgba(46, 212, 191, 0) 70%);
}

.auth-container {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
  position: relative;
  z-index: 1;
}

.auth-card {
  width: 520px;
  max-width: 100%;
  padding: 36px 40px 32px;
  background: rgba(255, 255, 255, 0.94);
  border: 1px solid rgba(255, 255, 255, 0.9);
  border-radius: 24px;
  box-shadow: 0 24px 80px rgba(40, 62, 116, 0.14), 0 4px 18px rgba(40, 62, 116, 0.06);
  backdrop-filter: blur(20px);
}

.auth-brand {
  display: flex;
  justify-content: center;
  margin-bottom: 28px;
}

.brand-logo {
  display: block;
  width: 184px;
  max-width: 100%;
  height: auto;
  object-fit: contain;
}

.auth-card :deep(.auth-heading) {
  margin-bottom: 28px;
  text-align: center;
}

.auth-card :deep(.auth-heading h1) {
  color: #1d2129;
  font-size: 26px;
  font-weight: 700;
  letter-spacing: -0.02em;
  line-height: 1.35;
}

.auth-card :deep(.auth-heading p) {
  margin-top: 8px;
  color: #86909c;
  font-size: 14px;
  line-height: 1.6;
}

.auth-card :deep(.auth-form .arco-form-item) {
  margin-bottom: 20px;
}

.auth-card :deep(.auth-form .arco-input-wrapper) {
  min-height: 48px;
  background: #f8faff;
  border-color: #e3e9f3;
  border-radius: 12px;
  transition: background 0.2s ease, border-color 0.2s ease, box-shadow 0.2s ease;
}

.auth-card :deep(.auth-form .arco-input-wrapper:hover) {
  background: #fff;
  border-color: #9db9ff;
}

.auth-card :deep(.auth-form .arco-input-wrapper-focus) {
  background: #fff;
  border-color: #165dff;
  box-shadow: 0 0 0 3px rgba(22, 93, 255, 0.1);
}

.auth-card :deep(.auth-form .arco-input) {
  color: #1d2129;
  background: transparent;
  font-size: 14px;
}

.auth-card :deep(.auth-form .arco-input::placeholder) {
  color: #a7afbd;
}

.auth-card :deep(.auth-form .arco-input-prefix) {
  margin-right: 10px;
  color: #8d98aa;
}

.auth-card :deep(.auth-form .arco-input-prefix .arco-icon) {
  font-size: 17px;
}

.auth-card :deep(.auth-form .arco-btn) {
  border-radius: 12px;
  font-weight: 600;
}

.auth-card :deep(.auth-form .arco-btn-primary) {
  min-height: 48px;
  box-shadow: 0 8px 18px rgba(22, 93, 255, 0.2);
}

.auth-card :deep(.auth-form .arco-alert) {
  border-radius: 12px;
}

@media (max-width: 560px) {
  .auth-container {
    padding: 24px 16px;
  }

  .auth-card {
    padding: 28px 20px 24px;
    border-radius: 20px;
  }

  .brand-logo {
    width: 168px;
  }

  .auth-card :deep(.auth-heading h1) {
    font-size: 24px;
  }
}
</style>
