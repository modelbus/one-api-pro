<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">{{ $t('redeem.title') }}</h1>
      <p class="page-subtitle">{{ $t('redeem.subtitle') }}</p>
    </div>

    <!-- 兑换表单（参考 tbus-web Redeem.vue 风格） -->
    <div class="redeem-card">
      <div class="redeem-icon">🎁</div>
      <div class="redeem-form">
        <a-input
          v-model="code"
          :placeholder="$t('redeem.placeholder')"
          size="large"
          class="redeem-input"
          :status="error ? 'error' : (success ? 'success' : '')"
          allow-clear
          :max-length="128"
          @press-enter="handleRedeem"
        />
        <div v-if="error" class="redeem-error">{{ error }}</div>
        <div v-else-if="success" class="redeem-success">{{ success }}</div>
        <a-button
          type="primary"
          size="large"
          class="redeem-btn"
          :loading="loading"
          @click="handleRedeem"
        >
          {{ $t('redeem.submit') }}
        </a-button>
      </div>
      <div class="redeem-hint">{{ $t('redeem.hint') }}</div>
    </div>

    <!-- 额度信息（弱化展示，让兑换区域成为视觉重点） -->
    <div class="quota-info">
      <a-spin :loading="quotaLoading" style="width: 100%;">
        <div class="quota-info-row">
          <div class="quota-info-label">{{ $t('redeem.currentQuota') }}</div>
          <div :class="['quota-info-value', 'quota-info-value-current', { flash: quotaFlash }]" :key="quotaFlashKey">{{ formatNumber(quota) }}</div>
        </div>
        <div class="quota-info-row" v-if="quotaUsed > 0 || quotaTotal > 0">
          <div class="quota-info-label">{{ $t('redeem.usage') }}</div>
          <div class="quota-info-value">
            {{ formatNumber(quotaUsed) }} <span class="quota-divider">/</span> {{ formatNumber(quotaTotal) }}
          </div>
        </div>
      </a-spin>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'

const code = ref('')
const loading = ref(false)
const error = ref('')
const success = ref('')

const quotaLoading = ref(true)
const quota = ref(0)
const quotaUsed = ref(0)
const quotaTotal = ref(0)
const quotaFlash = ref(false)
const quotaFlashKey = ref(0)

const numberFormatter = new Intl.NumberFormat('en-US')

function formatNumber(n) {
  if (n == null || isNaN(n)) return '0'
  return numberFormatter.format(Math.floor(Number(n)))
}

function validateCode(c) {
  const v = (c || '').trim()
  if (!v) return '请输入兑换码'
  if (v.length < 8) return '兑换码格式不正确'
  return ''
}

async function fetchQuota() {
  quotaLoading.value = true
  try {
    const { data } = await api.get('/api/user/self')
    if (data?.success && data.data) {
      quota.value = data.data.quota ?? data.data.remain_quota ?? 0
      quotaUsed.value = data.data.used_quota ?? 0
      quotaTotal.value = (Number(quota.value) + Number(quotaUsed.value)) || 0
    }
  } catch (e) {
    // silent
  } finally {
    quotaLoading.value = false
  }
}

function flashQuota() {
  quotaFlashKey.value++
  quotaFlash.value = true
  setTimeout(() => { quotaFlash.value = false }, 700)
}

async function handleRedeem() {
  error.value = ''
  success.value = ''
  const v = validateCode(code.value)
  if (v) {
    error.value = v
    return
  }

  loading.value = true
  try {
    const { data } = await api.post('/api/user/topup', { key: code.value.trim() })
    if (data?.success) {
      const amount = data.data ?? 0
      success.value = `兑换成功，已到账 ${formatNumber(amount)} 额度`
      code.value = ''
      // Re-fetch quota so the "当前额度/已使用" panel reflects the
      // new balance immediately and the number visibly increases.
      await fetchQuota()
      flashQuota()
    } else {
      error.value = data?.message || '兑换失败'
    }
  } catch (e) {
    error.value = e.response?.data?.message || e.message || '兑换失败'
  } finally {
    loading.value = false
  }
}

onMounted(fetchQuota)
</script>

<style scoped>
.page-container {
  padding: 24px;
  max-width: 720px;
  margin: 0 auto;
}
.page-header { margin-bottom: 20px; }
.page-title {
  font-size: 22px;
  font-weight: 700;
  color: var(--color-text-1);
  margin: 0 0 4px;
}
.page-subtitle {
  font-size: 14px;
  color: var(--color-text-3);
  margin: 0;
}

/* 兑换卡片（参考 tbus-web Redeem.vue） */
.redeem-card {
  background: #fff;
  border: 1px solid var(--color-border-2);
  border-radius: 12px;
  padding: 40px 24px;
  text-align: center;
}
.redeem-icon {
  font-size: 48px;
  margin-bottom: 20px;
}
.redeem-form {
  max-width: 400px;
  margin: 0 auto 16px;
}
.redeem-input {
  margin-bottom: 12px;
  border-radius: 8px;
}
.redeem-error {
  color: #f53f3f;
  font-size: 13px;
  margin-bottom: 12px;
  text-align: left;
}
.redeem-success {
  color: #00b42a;
  font-size: 13px;
  margin-bottom: 12px;
  text-align: left;
}
.redeem-btn {
  width: 100%;
  border-radius: 8px;
}
.redeem-hint {
  font-size: 12px;
  color: var(--color-text-4);
  margin-top: 4px;
}

/* 额度信息：弱化展示（小字号、低对比度、放在兑换区域下方） */
.quota-info {
  max-width: 400px;
  margin: 16px auto 0;
  padding: 12px 16px;
  background: transparent;
}
.quota-info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 0;
}
.quota-info-row + .quota-info-row {
  border-top: 1px dashed var(--color-border-2);
}
.quota-info-label {
  font-size: 13px;
  color: var(--color-text-4);
}
.quota-info-value {
  font-size: 13px;
  color: var(--color-text-3);
  font-variant-numeric: tabular-nums;
  font-feature-settings: 'tnum';
  transition: color 0.2s;
}
.quota-info-value-current {
  font-size: 15px;
  color: #00b42a;
  font-weight: 600;
}
.quota-info-value-current.flash {
  animation: quota-flash 0.6s ease-out;
}
@keyframes quota-flash {
  0%   { transform: scale(1); }
  50%  { transform: scale(1.18); }
  100% { transform: scale(1); }
}
.quota-divider {
  margin: 0 6px;
  color: var(--color-text-4);
}
</style>
