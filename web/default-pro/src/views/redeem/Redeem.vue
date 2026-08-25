<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">{{ $t('redeem.title') }}</h1>
      <p class="page-subtitle">{{ $t('redeem.subtitle') }}</p>
    </div>

    <!-- 当前额度 -->
    <a-card :bordered="false" class="quota-card">
      <a-spin :loading="quotaLoading">
        <a-row :gutter="16" align="center">
          <a-col :span="12">
            <a-statistic
              :title="$t('redeem.currentQuota')"
              :value="quota"
              :value-style="{ color: '#165DFF', fontWeight: 700 }"
            >
              <template #prefix><span class="stat-prefix">$</span></template>
            </a-statistic>
          </a-col>
          <a-col :span="12" v-if="quotaUsed > 0 || quotaTotal > 0">
            <a-statistic
              :title="$t('redeem.usage')"
              :value="quotaUsed"
              :precision="0"
            >
              <template #suffix>/ {{ quotaTotal || '∞' }}</template>
            </a-statistic>
          </a-col>
        </a-row>
      </a-spin>
    </a-card>

    <!-- 兑换表单（参考 tbus-web Redeem.vue 风格） -->
    <div class="redeem-card">
      <div class="redeem-icon">🎁</div>
      <div class="redeem-form">
        <a-input
          v-model="code"
          :placeholder="$t('redeem.placeholder')"
          size="large"
          class="redeem-input"
          :status="error ? 'error' : ''"
          allow-clear
          :max-length="128"
          @press-enter="handleRedeem"
        />
        <div v-if="error" class="redeem-error">{{ error }}</div>
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

    <!-- 兑换成功弹窗 -->
    <div v-if="success" class="modal-overlay" @click.self="success = false">
      <div class="modal success-modal">
        <div class="modal-icon">✅</div>
        <div class="modal-title">{{ $t('redeem.successTitle') }}</div>
        <div class="modal-desc">{{ $t('redeem.successDesc', { amount: redeemedAmount }) }}</div>
        <div class="modal-footer">
          <a-button @click="success = false">{{ $t('common.close') }}</a-button>
          <a-button type="primary" @click="goPlans">{{ $t('redeem.viewPlans') }}</a-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import api from '@/api'

const router = useRouter()

const code = ref('')
const loading = ref(false)
const error = ref('')
const success = ref(false)

const quotaLoading = ref(true)
const quota = ref(0)
const quotaUsed = ref(0)
const quotaTotal = ref(0)

const redeemedAmount = ref(0)

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
      quotaTotal.value = (quota.value + quotaUsed.value) || 0
    }
  } catch (e) {
    // silent
  } finally {
    quotaLoading.value = false
  }
}

async function handleRedeem() {
  error.value = validateCode(code.value)
  if (error.value) return

  loading.value = true
  error.value = ''
  try {
    const { data } = await api.post('/api/user/topup', { key: code.value.trim() })
    if (data?.success) {
      redeemedAmount.value = data.data ?? 0
      success.value = true
      code.value = ''
      await fetchQuota()
      Message.success('兑换成功')
    } else {
      error.value = data?.message || '兑换失败'
    }
  } catch (e) {
    error.value = e.response?.data?.message || e.message || '兑换失败'
  } finally {
    loading.value = false
  }
}

function goPlans() {
  success.value = false
  router.push('/plans')
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

.quota-card {
  background: var(--color-bg-2);
  border-radius: 12px;
  margin-bottom: 20px;
  padding: 4px 0;
}
.stat-prefix {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-3);
  margin-right: 4px;
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
.redeem-btn {
  width: 100%;
  border-radius: 8px;
}
.redeem-hint {
  font-size: 12px;
  color: var(--color-text-4);
  margin-top: 4px;
}

/* 成功弹窗 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}
.modal {
  background: #fff;
  border-radius: 16px;
  padding: 32px;
  width: 360px;
  max-width: 90vw;
}
.success-modal {
  text-align: center;
}
.modal-icon {
  font-size: 56px;
  margin-bottom: 16px;
}
.modal-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--color-text-1);
  margin-bottom: 8px;
}
.modal-desc {
  font-size: 14px;
  color: var(--color-text-3);
  margin-bottom: 24px;
}
.modal-footer {
  display: flex;
  justify-content: center;
  gap: 12px;
}
</style>
