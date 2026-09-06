<template>
  <div class="legal-layout">
    <header class="legal-header">
      <div class="legal-header-inner">
        <router-link to="/" class="legal-brand">
          <img :src="logoPng" alt="ONE-API-PRO" />
        </router-link>
        <router-link to="/" class="legal-back">返回首页</router-link>
      </div>
    </header>
    <main class="legal-main">
      <article class="legal-document">
        <header class="legal-document-header">
          <span class="legal-document-label">法律文件</span>
          <h1>{{ title }}</h1>
          <p>{{ description }}</p>
          <span v-if="updatedAt" class="legal-updated">{{ updatedAt }}</span>
        </header>
        <div class="legal-content">
          <slot />
        </div>
      </article>
    </main>
  </div>
</template>

<script setup>
import logoPng from '@/assets/logo.png'

defineProps({
  title: { type: String, required: true },
  description: { type: String, required: true },
  updatedAt: { type: String, default: '' },
})
</script>

<style scoped>
.legal-layout {
  min-height: 100vh;
  color: #1d2129;
  background: radial-gradient(circle at 10% 0%, rgba(22, 93, 255, 0.09), transparent 32%), radial-gradient(circle at 90% 100%, rgba(114, 46, 209, 0.08), transparent 34%), #f7f9fc;
}

.legal-header {
  background: rgba(255, 255, 255, 0.82);
  border-bottom: 1px solid rgba(29, 33, 41, 0.07);
  backdrop-filter: blur(16px);
}

.legal-header-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: min(1120px, calc(100% - 48px));
  height: 76px;
  margin: 0 auto;
}

.legal-brand {
  display: inline-flex;
  align-items: center;
}

.legal-brand img {
  display: block;
  width: 156px;
  height: auto;
  object-fit: contain;
}

.legal-back {
  color: #4e5969;
  font-size: 14px;
  font-weight: 500;
  transition: color 0.2s ease;
}

.legal-back:hover {
  color: rgb(var(--primary-6));
}

.legal-main {
  width: min(880px, calc(100% - 48px));
  margin: 0 auto;
  padding: 72px 0 88px;
}

.legal-document {
  padding: 56px 72px 68px;
  background: rgba(255, 255, 255, 0.94);
  border: 1px solid rgba(229, 233, 240, 0.9);
  border-radius: 22px;
  box-shadow: 0 18px 60px rgba(40, 62, 116, 0.09);
}

.legal-document-header {
  padding-bottom: 32px;
  border-bottom: 1px solid #eef0f3;
}

.legal-document-label {
  display: inline-block;
  margin-bottom: 14px;
  color: rgb(var(--primary-6));
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 0.08em;
}

.legal-document-header h1 {
  color: #1d2129;
  font-size: 34px;
  font-weight: 700;
  letter-spacing: -0.03em;
  line-height: 1.3;
}

.legal-document-header p {
  max-width: 680px;
  margin-top: 12px;
  color: #7b8699;
  font-size: 15px;
  line-height: 1.8;
}

.legal-updated {
  display: block;
  margin-top: 18px;
  color: #a1a9b5;
  font-size: 12px;
}

.legal-content {
  color: #4e5969;
  font-size: 14px;
  line-height: 1.9;
}

.legal-content :deep(.legal-intro) {
  margin: 28px 0 0;
}

.legal-content :deep(.legal-intro p),
.legal-content :deep(.legal-section p) {
  margin: 10px 0;
}

.legal-content :deep(.legal-section) {
  margin-top: 36px;
}

.legal-content :deep(.legal-section h2) {
  margin-bottom: 10px;
  color: #1d2129;
  font-size: 18px;
  font-weight: 650;
  line-height: 1.5;
}

.legal-content :deep(.legal-section ul) {
  margin: 10px 0 0;
  padding-left: 22px;
}

.legal-content :deep(.legal-section li) {
  margin: 7px 0;
}

@media (max-width: 720px) {
  .legal-header-inner {
    width: min(100% - 32px, 1120px);
    height: 68px;
  }

  .legal-brand img {
    width: 140px;
  }

  .legal-main {
    width: min(100% - 32px, 880px);
    padding: 32px 0 48px;
  }

  .legal-document {
    padding: 36px 24px 44px;
    border-radius: 18px;
  }

  .legal-document-header h1 {
    font-size: 28px;
  }
}
</style>
