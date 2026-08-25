<template>
  <div class="plans-page">
    <h2 class="page-title">套餐</h2>

    <a-spin :loading="loading" style="width: 100%">
      <div v-if="currentPlan" class="current-banner">
        当前订阅：
        <strong>{{ currentPlan.plan ? currentPlan.plan.name : '套餐 #' + currentPlan.plan_id }}</strong>
        <span class="expire">到期时间：{{ formatTime(currentPlan.end_time) }}</span>
        <span class="days-left" v-if="daysLeft > 0">（剩余 {{ daysLeft }} 天）</span>
      </div>

      <div class="plans-grid">
        <div
          v-for="p in plans"
          :key="p.id"
          class="plan-card"
          :class="{ 'is-current': currentPlan && currentPlan.plan_id === p.id, 'is-downgrade': isDowngrade(p) }"
        >
          <a-tag v-if="p.recommended" color="arcoblue" class="badge">推荐</a-tag>
          <a-tag v-if="currentPlan && currentPlan.plan_id === p.id" color="green" class="badge">当前套餐</a-tag>
          <div class="plan-name">{{ p.name }}</div>
          <div class="plan-price">
            <span class="price-symbol">¥</span>
            <span class="price-num">{{ p.price }}</span>
            <span class="price-period">/ {{ p.duration_text || '月' }}</span>
          </div>
          <div v-if="getUpgradeDiff(p) !== null" class="plan-diff">
            差价：¥{{ getUpgradeDiff(p) }}
          </div>
          <div v-if="p.description" class="plan-desc">{{ p.description }}</div>
          <ul v-if="featureList(p).length" class="plan-features">
            <li v-for="(f, i) in featureList(p)" :key="i">✓ {{ f }}</li>
          </ul>
          <div class="plan-action">
            <a-button
              type="primary"
              :disabled="!canSubscribe(p)"
              long
              @click="openSubscribe(p)"
            >{{ subscribeButtonText(p) }}</a-button>
          </div>
        </div>
      </div>
    </a-spin>

    <a-modal
      v-model:visible="modalVisible"
      :title="modalTitle"
      :ok-loading="submitting"
      @ok="confirmSubscribe"
      width="480px"
    >
      <p v-if="selectedPlan">
        套餐：<strong>{{ selectedPlan.name }}</strong>
      </p>
      <p v-if="amount > 0">应付金额：<strong>¥{{ amount }}</strong></p>
      <p v-else>金额：<strong>免费</strong></p>
      <p v-if="mode === 'price_diff'">升级模式：差价升级</p>
      <p v-else>升级模式：新订阅</p>
      <a-radio-group v-model="payMethod" type="button">
        <a-radio value="wechat">微信支付</a-radio>
        <a-radio value="alipay">支付宝</a-radio>
        <a-radio value="bank">银行转账</a-radio>
      </a-radio-group>
      <a-alert v-if="payWarning" type="warning" style="margin-top:12px">{{ payWarning }}</a-alert>
      <div v-if="payQrCode" class="qr-area">
        <p>请使用 {{ payMethod === 'wechat' ? '微信' : payMethod === 'alipay' ? '支付宝' : '对应渠道' }} 扫描下方二维码完成支付：</p>
        <img v-if="payQrCode" :src="qrImage" alt="pay qr" class="qr-img" />
        <p class="qr-link">链接：<a :href="payQrCode" target="_blank">{{ payQrCode }}</a></p>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import planApi from '@/api/plan'
import orderApi from '@/api/order'

const router = useRouter()
const loading = ref(false)
const submitting = ref(false)
const plans = ref([])
const currentPlan = ref(null)
const modalVisible = ref(false)
const selectedPlan = ref(null)
const payMethod = ref('wechat')
const payWarning = ref('')
const payQrCode = ref('')
const qrImage = ref('')
const order = ref(null)
const amount = ref(0)
const mode = ref('')

const daysLeft = computed(() => {
  if (!currentPlan.value || !currentPlan.value.end_time) return 0
  const now = Math.floor(Date.now() / 1000)
  return Math.max(0, Math.ceil((currentPlan.value.end_time - now) / 86400))
})

const modalTitle = computed(() => {
  if (!selectedPlan.value) return '订阅套餐'
  return mode.value === 'price_diff' ? '升级套餐' : '订阅套餐'
})

function featureList(p) {
  try {
    const obj = JSON.parse(p.features || '{}')
    return Object.keys(obj)
  } catch {
    return []
  }
}

function canSubscribe(p) {
  if (!currentPlan.value) return true
  if (currentPlan.value.plan_id === p.id) return false
  if (p.sort <= (currentPlan.value.plan && currentPlan.value.plan.sort || 0)) {
    // downgrades are blocked
    return false
  }
  return true
}

function isDowngrade(p) {
  return currentPlan.value && p.sort <= (currentPlan.value.plan && currentPlan.value.plan.sort || 0)
}

function subscribeButtonText(p) {
  if (currentPlan.value && currentPlan.value.plan_id === p.id) return '当前套餐'
  if (isDowngrade(p)) return '暂不可用'
  if (currentPlan.value) return '升级到 ' + p.name
  return '立即订阅'
}

function getUpgradeDiff(p) {
  if (!currentPlan.value || !currentPlan.value.plan) return null
  if (p.sort <= currentPlan.value.plan.sort) return null
  // backend returns the actual amount; we can't compute it client-side
  // without the plan prices, so return a placeholder
  return null
}

function formatTime(t) {
  if (!t) return '-'
  const d = new Date(t * 1000)
  return d.toLocaleString('zh-CN')
}

async function loadData() {
  loading.value = true
  try {
    const [listRes, currentRes] = await Promise.all([
      planApi.list(),
      planApi.current().catch(() => ({ data: { data: null } })),
    ])
    plans.value = (listRes.data && listRes.data.data) || []
    currentPlan.value = (currentRes.data && currentRes.data.data) || null
  } catch (e) {
    Message.error('加载套餐失败')
  } finally {
    loading.value = false
  }
}

function openSubscribe(p) {
  if (!canSubscribe(p)) return
  selectedPlan.value = p
  payMethod.value = 'wechat'
  payWarning.value = ''
  payQrCode.value = ''
  qrImage.value = ''
  order.value = null
  amount.value = p.price
  mode.value = currentPlan.value ? 'price_diff' : 'stack'
  modalVisible.value = true
}

async function confirmSubscribe() {
  if (!selectedPlan.value) return
  submitting.value = true
  payWarning.value = ''
  payQrCode.value = ''
  try {
    const { data } = await orderApi.createPlanOrder({
      plan_id: selectedPlan.value.id,
      pay_method: payMethod.value,
    })
    if (!data.success) {
      Message.error(data.message || '创建订单失败')
      return
    }
    order.value = data.order
    amount.value = data.amount
    mode.value = data.mode
    const pay = data.pay || {}
    payWarning.value = pay.warning || pay.note || ''
    payQrCode.value = pay.pay_url || pay.qr_code || ''
    // Render the QR code as data: URL via a tiny inline SVG-as-img trick:
    // we just put the URL into the modal text — a real client would
    // convert the URL to a QR image. For the spec we leave the URL
    // visible so the user can copy it.
    qrImage.value = ''
    Message.success('订单已创建，订单号 ' + data.order.order_no)
  } catch (e) {
    Message.error(e.response?.data?.message || '网络错误')
  } finally {
    submitting.value = false
  }
}

onMounted(loadData)
</script>

<style scoped>
.plans-page { padding: 16px; }
.page-title { font-size: 20px; font-weight: 600; margin-bottom: 16px; color: var(--color-text-1); }
.current-banner { padding: 12px 16px; background: var(--color-fill-2); border-radius: 6px; margin-bottom: 16px; }
.current-banner .expire { margin-left: 12px; color: var(--color-text-3); }
.current-banner .days-left { margin-left: 8px; color: #165dff; }
.plans-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 16px; }
.plan-card { position: relative; background: var(--color-bg-2); border: 1px solid var(--color-border-2); border-radius: 8px; padding: 20px 18px 18px; display: flex; flex-direction: column; gap: 8px; transition: border-color 0.15s, box-shadow 0.15s; }
.plan-card:hover { border-color: #165dff40; box-shadow: 0 6px 24px rgba(0, 0, 0, 0.06); }
.plan-card.is-current { border-color: #00b42a; box-shadow: 0 0 0 1px #00b42a inset; }
.plan-card .badge { position: absolute; top: 12px; right: 12px; }
.plan-name { font-size: 16px; font-weight: 600; color: var(--color-text-1); }
.plan-price { display: flex; align-items: baseline; gap: 4px; }
.plan-price .price-symbol { font-size: 14px; color: var(--color-text-3); }
.plan-price .price-num { font-size: 28px; font-weight: 700; color: #165dff; }
.plan-price .price-period { font-size: 13px; color: var(--color-text-3); }
.plan-diff { font-size: 13px; color: #ff7d00; }
.plan-desc { color: var(--color-text-2); font-size: 13px; line-height: 1.6; }
.plan-features { list-style: none; padding: 0; margin: 0; color: var(--color-text-2); font-size: 13px; line-height: 1.8; }
.plan-action { margin-top: auto; padding-top: 12px; }
.qr-area { margin-top: 12px; padding: 12px; background: var(--color-fill-1); border-radius: 6px; }
.qr-img { display: block; margin: 12px auto; max-width: 200px; }
.qr-link { font-size: 12px; word-break: break-all; color: var(--color-text-3); }
</style>
