<template>
  <div class="topup-page">
    <a-row :gutter="24" class="topup-row">
      <a-col :xs="24" :md="12">
        <a-card :bordered="false" class="topup-card">
          <template #title>
            <span class="card-title">{{ $t('topup.selfTopup') }}</span>
          </template>

          <div class="quota-display">
            <a-statistic
              :title="$t('topup.currentQuota')"
              :value="quota"
              :value-style="{ color: '#165DFF' }"
            >
              <template #prefix>
                <span class="stat-prefix">$</span>
              </template>
            </a-statistic>
          </div>

          <a-divider style="margin: 20px 0" />

          <div class="code-redeem-section">
            <p class="section-label">Redeem Code</p>
            <a-input-group style="width: 100%">
              <a-input
                v-model="redeemCode"
                placeholder="Enter redemption code"
                allow-clear
                :max-length="128"
              />
              <a-button @click="pasteCode">
                <template #icon><icon-copy /></template>
                Paste
              </a-button>
            </a-input-group>
            <a-button
              type="primary"
              long
              :loading="redeeming"
              :disabled="!redeemCode.trim()"
              style="margin-top: 12px"
              @click="handleRedeem"
            >
              Redeem
            </a-button>
          </div>
        </a-card>
      </a-col>

      <a-col :xs="24" :md="12">
        <a-card :bordered="false" class="topup-card">
          <template #title>
            <span class="card-title">Buy Credits</span>
          </template>

          <div class="quota-display">
            <a-statistic
              :title="$t('topup.currentQuota')"
              :value="quota"
              :value-style="{ color: '#165DFF' }"
            >
              <template #prefix>
                <span class="stat-prefix">$</span>
              </template>
            </a-statistic>
          </div>

          <a-divider style="margin: 20px 0" />

          <template v-if="topUpLink">
            <div class="buy-section">
              <p class="section-label">External Purchase</p>
              <a-alert type="info" style="margin-bottom: 16px">
                Click the button below to purchase credits through our external payment page.
              </a-alert>
              <a-button type="outline" size="large" :href="topUpLink" target="_blank" long>
                Buy Credits
              </a-button>
            </div>
          </template>

          <template v-else>
            <div class="buy-section">
              <a-empty description="No top-up link available. Please contact the administrator." />
            </div>
          </template>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconCopy } from '@arco-design/web-vue/es/icon'
import { useStatusStore } from '@/stores/status'
import api from '@/api'

const statusStore = useStatusStore()

const loading = ref(true)
const quota = ref(0)
const redeemCode = ref('')
const redeeming = ref(false)
const topUpLink = ref('')

onMounted(async () => {
  await fetchUserData()
})

async function fetchUserData() {
  loading.value = true
  try {
    const { data } = await api.get('/api/user/self')
    if (data.success) {
      const user = data.data
      quota.value = user.quota ?? user.remain_quota ?? 0
    }
    topUpLink.value = statusStore.status?.top_up_link ?? ''
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || 'Failed to fetch user data')
  } finally {
    loading.value = false
  }
}

async function pasteCode() {
  try {
    const text = await navigator.clipboard.readText()
    if (text) {
      redeemCode.value = text.trim()
      Message.success('Code pasted from clipboard')
    }
  } catch {
    Message.warning('Unable to access clipboard')
  }
}

async function handleRedeem() {
  const code = redeemCode.value.trim()
  if (!code) {
    Message.warning('Please enter a redemption code')
    return
  }
  redeeming.value = true
  try {
    const { data } = await api.post('/api/user/topup', { key: code })
    if (data.success) {
      Message.success('Top-up successful')
      redeemCode.value = ''
      await fetchUserData()
    } else {
      Message.error(data.message || 'Top-up failed')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || 'Top-up failed')
  } finally {
    redeeming.value = false
  }
}
</script>

<style scoped>
.topup-page {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.topup-row {
  align-items: stretch;
}

.topup-card {
  height: 100%;
}

.card-title {
  font-size: 17px;
  font-weight: 600;
}

.quota-display {
  text-align: center;
  padding: 16px 0;
}

.stat-prefix {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-3);
}

.section-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-2);
  margin: 0 0 12px 0;
}

.code-redeem-section {
  padding: 0;
}

.buy-section {
  padding: 0;
}
</style>
