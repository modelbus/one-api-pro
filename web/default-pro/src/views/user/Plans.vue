<template>
  <div class="page-container">
    <!-- 升级提示横幅 -->
    <div v-if="currentPlan && !isExpired(currentPlan)" class="upgrade-banner">
      <span class="banner-icon">💡</span>
      <i18n-t keypath="plansPage.upgradeBanner" tag="span">
        <template #name><strong>{{ currentPlan.name }}</strong></template>
      </i18n-t>
    </div>

    <!-- 套餐卡片网格 -->
    <a-spin :loading="loading" dot :tip="$t('plansPage.loading')" class="plans-spin">
    <div v-if="!loading && planList.length === 0" class="empty-state">
      <div class="empty-illustration">
        <svg width="80" height="80" viewBox="0 0 80 80" fill="none" xmlns="http://www.w3.org/2000/svg">
          <rect x="14" y="22" width="52" height="44" rx="6" stroke="#D1D1D6" stroke-width="2" fill="#FFFFFF"/>
          <rect x="20" y="32" width="22" height="3" rx="1.5" fill="#D1D1D6"/>
          <rect x="20" y="40" width="40" height="3" rx="1.5" fill="#E5E5E7"/>
          <rect x="20" y="48" width="32" height="3" rx="1.5" fill="#E5E5E7"/>
          <circle cx="56" cy="56" r="10" fill="#F7F8FA" stroke="#D1D1D6" stroke-width="2"/>
          <line x1="50" y1="56" x2="62" y2="56" stroke="#D1D1D6" stroke-width="2" stroke-linecap="round"/>
        </svg>
      </div>
      <p class="empty-text">{{ $t('plansPage.emptyTitle') }}</p>
      <p class="empty-desc">{{ $t('plansPage.emptyDesc') }}</p>
    </div>
    <div v-else class="plans-grid">
      <div
        v-for="plan in planList"
        :key="plan.id"
        class="plan-card"
        :class="{
          recommended: plan.recommended,
          'is-current': isCurrentPlan(plan)
        }"
      >
        <div v-if="plan.recommended" class="plan-badge">{{ $t('plansPage.recommended') }}</div>
        <div v-if="isCurrentPlan(plan)" class="current-badge">{{ $t('plansPage.currentPlan') }}</div>

        <div class="plan-name">{{ plan.name }}</div>
        <div v-if="plan.multiplier" class="plan-multiplier">{{ plan.multiplier }}</div>
        <div class="plan-price">
          <span class="plan-price-symbol">¥</span>
          <span class="plan-price-num">{{ plan.price }}</span>
          <span class="plan-price-period">/{{ plan.duration_text || $t('plansPage.monthUnit') }}</span>
        </div>
        <div class="plan-desc">{{ plan.description }}</div>

        <hr class="plan-divider" />

        <ul class="plan-features">
          <li v-for="(feat, i) in plan.features" :key="i">{{ feat }}</li>
        </ul>

        <div class="plan-btn-wrap">
          <a-button
            :type="isCurrentPlan(plan) ? 'outline' : plan.recommended ? 'primary' : 'outline'"
            long
            size="large"
            class="plan-btn"
            :disabled="isCurrentPlan(plan) || isLowerPlan(plan)"
            @click="isCurrentPlan(plan) || isLowerPlan(plan) ? null : (currentPlan && !isExpired(currentPlan) ? onUpgradeClick(plan) : onPurchaseClick(plan))"
          >
            <template v-if="isCurrentPlan(plan)">{{ $t('plansPage.currentPlan') }}</template>
            <template v-else-if="isLowerPlan(plan)">{{ $t('plansPage.unavailable') }}</template>
            <template v-else-if="currentPlan && !isExpired(currentPlan)">{{ $t('plansPage.upgradeTo', { name: plan.name }) }}</template>
            <template v-else>{{ $t('plansPage.subscribeNow') }}</template>
          </a-button>
        </div>
      </div>
    </div>
    </a-spin>

    <!-- 升级确认弹窗 -->
    <a-modal
      v-model:visible="upgradeModalVisible"
      :footer="false"
      :mask-closable="true"
      :title="$t('plansPage.confirmUpgradeTitle')"
      class="upgrade-modal"
      width="520px"
    >
      <div v-if="targetPlan" class="upgrade-modal-body">
        <div class="upgrade-compare">
          <div class="compare-card left">
            <div class="compare-label">{{ $t('plansPage.currentPlan') }}</div>
            <div class="compare-name">{{ currentPlan?.name }}</div>
            <div v-if="currentPlan && !isExpired(currentPlan)" class="compare-days">
              <i18n-t keypath="plansPage.remainingDays" tag="span">
                <template #n><strong>{{ getRemainingDays(currentPlan) }}</strong></template>
              </i18n-t>
            </div>
            <div v-else class="compare-days expired">{{ $t('plansPage.expired') }}</div>
            <div class="compare-price">¥{{ currentPlan?.price || 0 }}/{{ currentPlan?.duration_text || $t('plansPage.monthUnit') }}</div>
          </div>
          <div class="compare-arrow">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
              <path d="M5 12h14M14 5l7 7-7 7" stroke="#007AFF" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <div class="compare-card right">
            <div class="compare-label">{{ $t('plansPage.upgradeToLabel') }}</div>
            <div class="compare-name">{{ targetPlan.name }}</div>
            <div class="compare-days">{{ targetPlan.duration_text || $t('plansPage.oneMonth') }}</div>
            <div class="compare-price highlight">¥{{ targetPlan.price }}/{{ targetPlan.duration_text || $t('plansPage.monthUnit') }}</div>
          </div>
        </div>

        <div class="upgrade-diff">
          <div class="diff-title">{{ $t('plansPage.diffTitle') }}</div>
          <div class="diff-row">
            <span>{{ $t('plansPage.newPlanPrice') }}</span>
            <span>¥{{ targetPlan.price }}</span>
          </div>
          <div class="diff-row">
            <span>{{ $t('plansPage.remainingValue') }}</span>
            <span class="deduct">-¥{{ getRemainingValue() }}</span>
          </div>
          <div class="diff-row total">
            <span>{{ $t('plansPage.payableDiff') }}</span>
            <span class="final-price">¥{{ getUpgradeDiff() }}</span>
          </div>
        </div>

        <div class="upgrade-note">
          {{ $t('plansPage.upgradeNote') }}
        </div>

        <div class="pay-picker">
          <div class="pay-picker-label">{{ $t('plansPage.payMethod') }}</div>
          <div class="pay-picker-list">
            <button
              type="button"
              class="pay-picker-item wechat"
              :class="{ active: selectedPayMethod === 'wechat' }"
              :disabled="!availablePayMethods.includes('wechat')"
              @click="selectedPayMethod = 'wechat'"
            >
              <icon-wechatpay :size="28" class="pay-picker-icon" />
              <span class="pay-picker-name">{{ $t('plansPage.wechatPay') }}</span>
              <span v-if="selectedPayMethod === 'wechat'" class="pay-picker-check">✓</span>
            </button>
            <button
              type="button"
              class="pay-picker-item alipay"
              :class="{ active: selectedPayMethod === 'alipay' }"
              :disabled="!availablePayMethods.includes('alipay')"
              @click="selectedPayMethod = 'alipay'"
            >
              <icon-alipay-circle :size="28" class="pay-picker-icon" />
              <span class="pay-picker-name">{{ $t('plansPage.alipayPay') }}</span>
              <span v-if="selectedPayMethod === 'alipay'" class="pay-picker-check">✓</span>
            </button>
          </div>
        </div>

        <div class="upgrade-footer">
          <a-button @click="upgradeModalVisible = false">{{ $t('plansPage.cancel') }}</a-button>
          <a-button type="primary" :loading="upgrading" @click="confirmUpgrade">{{ $t('plansPage.confirmUpgrade') }}</a-button>
        </div>
      </div>
    </a-modal>

    <!-- 购买确认弹窗 -->
    <a-modal
      v-model:visible="purchaseModalVisible"
      :footer="false"
      :mask-closable="true"
      :title="$t('plansPage.confirmPurchaseTitle')"
      width="420px"
    >
      <div v-if="selectedPlan" class="purchase-modal-body">
        <div class="purchase-summary">
          <div class="purchase-row">
            <span class="purchase-label">{{ $t('plansPage.planName') }}</span>
            <span class="purchase-value">{{ selectedPlan.name }}</span>
          </div>
          <div class="purchase-row">
            <span class="purchase-label">{{ $t('plansPage.planPrice') }}</span>
            <span class="purchase-price">¥{{ selectedPlan.price }}/{{ selectedPlan.duration_text || $t('plansPage.monthUnit') }}</span>
          </div>
          <div class="purchase-row">
            <span class="purchase-label">{{ $t('plansPage.scenario') }}</span>
            <span class="purchase-value">{{ selectedPlan.description }}</span>
          </div>
        </div>

        <div class="pay-picker">
          <div class="pay-picker-label">{{ $t('plansPage.payMethod') }}</div>
          <div class="pay-picker-list">
            <button
              type="button"
              class="pay-picker-item wechat"
              :class="{ active: selectedPayMethod === 'wechat' }"
              :disabled="!availablePayMethods.includes('wechat')"
              @click="selectedPayMethod = 'wechat'"
            >
              <icon-wechatpay :size="28" class="pay-picker-icon" />
              <span class="pay-picker-name">{{ $t('plansPage.wechatPay') }}</span>
              <span v-if="selectedPayMethod === 'wechat'" class="pay-picker-check">✓</span>
            </button>
            <button
              type="button"
              class="pay-picker-item alipay"
              :class="{ active: selectedPayMethod === 'alipay' }"
              :disabled="!availablePayMethods.includes('alipay')"
              @click="selectedPayMethod = 'alipay'"
            >
              <icon-alipay-circle :size="28" class="pay-picker-icon" />
              <span class="pay-picker-name">{{ $t('plansPage.alipayPay') }}</span>
              <span v-if="selectedPayMethod === 'alipay'" class="pay-picker-check">✓</span>
            </button>
          </div>
        </div>

        <div class="purchase-footer">
          <a-button @click="purchaseModalVisible = false">{{ $t('plansPage.cancel') }}</a-button>
          <a-button type="primary" :loading="purchasing" @click="confirmPurchase">{{ $t('plansPage.confirmPurchase') }}</a-button>
        </div>
      </div>
    </a-modal>

    <!-- 支付二维码弹窗 -->
    <a-modal
      v-model:visible="paymentModalVisible"
      :footer="false"
      :mask-closable="true"
      :title="paymentModalTitle"
      width="420px"
      class="payment-modal"
    >
      <div class="payment-modal-body">
        <div class="qrcode-wrap">
          <img v-if="qrcodeDataUrl" :src="qrcodeDataUrl" :alt="paymentModalTitle" class="qrcode-img" />
          <div v-else class="qrcode-loading">{{ $t('plansPage.generatingQr') }}</div>
        </div>
        <div class="payment-tip">{{ paymentTip }}</div>
        <div class="payment-tip-sub">{{ $t('plansPage.payWithin5min') }}</div>
        <div class="payment-order-info">
          <div class="payment-info-row">
            <span>{{ $t('plansPage.planLabel') }}</span>
            <span class="payment-info-value">{{ paymentPackageName }}</span>
          </div>
          <div class="payment-info-row">
            <span>{{ $t('plansPage.amountLabel') }}</span>
            <span class="payment-info-price">¥{{ paymentAmount }}</span>
          </div>
        </div>
        <div class="payment-actions">
          <a-button long @click="paymentModalVisible = false">{{ $t('plansPage.close') }}</a-button>
          <a-button type="primary" long @click="$router.push('/orders')">{{ $t('plansPage.goOrders') }}</a-button>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Message } from '@arco-design/web-vue'
import { IconWechatpay, IconAlipayCircle } from '@arco-design/web-vue/es/icon'
import QRCode from 'qrcode'
import planApi from '@/api/plan'
import orderApi from '@/api/order'
import paymentApi from '@/api/payment'

const router = useRouter()
const { t } = useI18n()

const currentPlan = ref(null)
const planList = ref([])
const loading = ref(false)

const upgradeModalVisible = ref(false)
const targetPlan = ref(null)

const purchaseModalVisible = ref(false)
const selectedPlan = ref(null)

const paymentModalVisible = ref(false)
const paymentCodeUrl = ref('')
const paymentAmount = ref(0)
const paymentPackageName = ref('')
const qrcodeDataUrl = ref('')
const upgrading = ref(false)
const purchasing = ref(false)

// Payment method picker state. Defaults to WeChat; reset to the
// first available method whenever a new modal is opened.
const selectedPayMethod = ref('wechat')
const availablePayMethods = ref([]) // names like 'wechat' | 'alipay'
const lastPayMethod = ref('wechat') // remembers the last choice for the QR modal title

// paymentModalTitle / paymentTip are overridden by openPaymentModal /
// openBankTransferModal so the modal can show different copy for bank
// transfer vs WeChat/Alipay.
const paymentModalTitle = ref('')
const paymentTip = ref('')
function refreshDefaultPaymentCopy() {
  paymentModalTitle.value = lastPayMethod.value === 'alipay' ? t('plansPage.alipayScanTitle') : t('plansPage.wechatScanTitle')
  paymentTip.value = lastPayMethod.value === 'alipay' ? t('plansPage.alipayScanTip') : t('plansPage.wechatScanTip')
}

const NO_PAYMENT_MSG = computed(() => t('plansPage.noPaymentMsg'))

async function loadPlans() {
  loading.value = true
  try {
    const [allPlans, current] = await Promise.all([
      planApi.list(),
      planApi.current().catch(() => null),
    ])
    // Backend returns { success, message, data: [...] } for list,
    // { success, message, data: <userplan> } for current.
    // axios exposes the parsed body as response.data.
    currentPlan.value = current?.data?.data || null
    const list = Array.isArray(allPlans?.data?.data) ? allPlans.data.data : []
    planList.value = list.map((item, index) => ({
      id: item.id,
      name: item.name,
      price: item.price,
      description: item.description,
      recommended: item.recommended,
      features: parseFeatures(item.features),
      multiplier: item.multiplier,
      sort: item.sort ?? index + 1,
      duration_days: item.duration_days ?? 30,
      duration_text: item.duration_text || t('plansPage.oneMonth'),
    }))
  } catch (e) {
    Message.error(t('plansPage.loadPlansFailed'))
  } finally {
    loading.value = false
  }
}

function parseFeatures(raw) {
  if (!raw) return []
  if (Array.isArray(raw)) return raw
  try {
    const obj = typeof raw === 'string' ? JSON.parse(raw) : raw
    if (Array.isArray(obj)) return obj
    if (typeof obj === 'object' && obj !== null) {
      // features stored as { "API调用": true, ... } -> render as strings
      return Object.entries(obj).map(([k, v]) => v === true ? k : `${k}: ${v}`)
    }
  } catch {}
  return []
}

function isCurrentPlan(plan) {
  if (!currentPlan.value) return false
  if (isExpired(currentPlan.value)) return false
  return (currentPlan.value.name || '').toLowerCase() === (plan.name || '').toLowerCase()
}

function isLowerPlan(plan) {
  if (!currentPlan.value || isExpired(currentPlan.value)) return false
  return plan.sort < (currentPlan.value.sort ?? 0)
}

function isExpired(plan) {
  if (!plan) return false
  if (plan.is_expired !== undefined) return !!plan.is_expired
  if (plan.expire_time) return plan.expire_time * 1000 < Date.now()
  return false
}

function getRemainingDays(plan) {
  if (!plan) return 0
  if (plan.remaining_days !== undefined) return plan.remaining_days
  if (!plan.expire_time) return 0
  const now = Date.now()
  const end = plan.expire_time * 1000
  return Math.max(0, Math.ceil((end - now) / (1000 * 60 * 60 * 24)))
}

function getRemainingValue() {
  if (!currentPlan.value || isExpired(currentPlan.value)) return 0
  const remaining = getRemainingDays(currentPlan.value)
  const totalDays = currentPlan.value.duration_days || 30
  const currentPrice = Number(currentPlan.value.price) || 0
  return Math.round((remaining / totalDays) * currentPrice)
}

function getUpgradeDiff() {
  const diff = (Number(targetPlan.value?.price) || 0) - getRemainingValue()
  return Math.max(0, diff)
}

// Pre-flight check: refuse to open the purchase/upgrade modal when the
// admin has not enabled any payment channel. Returns true when it is
// safe to continue. Side-effect: refreshes `availablePayMethods` so
// the picker only shows enabled methods.
async function ensurePaymentEnabled() {
  try {
    const { data } = await paymentApi.status()
    const d = data?.data || {}
    const anyEnabled = !!d.any_enabled
    if (!anyEnabled) {
      Message.error(NO_PAYMENT_MSG.value)
      return false
    }
    availablePayMethods.value = (d.methods || [])
      .filter(m => m.enabled && (m.name === 'wechat' || m.name === 'alipay'))
      .map(m => m.name)
    // Default-select the first available method if the current
    // selection is no longer enabled (e.g. admin disabled WeChat).
    if (!availablePayMethods.value.includes(selectedPayMethod.value)) {
      selectedPayMethod.value = availablePayMethods.value[0] || 'wechat'
    }
    return true
  } catch (e) {
    Message.error(t('plansPage.payStatusFailed'))
    return false
  }
}

function handleUpgrade(plan) {
  const currentSort = currentPlan.value?.sort ?? 0
  if (currentSort >= plan.sort) {
    Message.warning({ content: t('plansPage.cannotDowngrade'), duration: 3000 })
    return
  }
  targetPlan.value = plan
  upgradeModalVisible.value = true
}

function handlePurchase(plan) {
  selectedPlan.value = plan
  purchaseModalVisible.value = true
}

async function onUpgradeClick(plan) {
  if (!(await ensurePaymentEnabled())) return
  handleUpgrade(plan)
}

async function onPurchaseClick(plan) {
  if (!(await ensurePaymentEnabled())) return
  handlePurchase(plan)
}

async function generateQRCode(url) {
  try {
    qrcodeDataUrl.value = await QRCode.toDataURL(url, {
      width: 220,
      margin: 2,
      color: { dark: '#000000', light: '#ffffff' },
    })
  } catch (e) {
    console.error(t('plansPage.qrGenerateFailed'), e)
  }
}

async function openPaymentModal(codeUrl, amount, packageName, payMethod) {
  paymentCodeUrl.value = codeUrl
  paymentAmount.value = amount
  paymentPackageName.value = packageName
  lastPayMethod.value = payMethod || selectedPayMethod.value
  refreshDefaultPaymentCopy()
  await generateQRCode(codeUrl)
  paymentModalVisible.value = true
}

// Show the payment modal in "transfer note" mode (no QR code) — used
// when the user picked bank transfer and the backend returned a note
// instead of a pre-payment URL.
function openBankTransferModal(amount, packageName, note, payMethod) {
  paymentAmount.value = amount
  paymentPackageName.value = packageName
  lastPayMethod.value = payMethod || selectedPayMethod.value
  qrcodeDataUrl.value = ''
  paymentTip.value = note || t('plansPage.bankTransferTip')
  paymentModalTitle.value = t('plansPage.transferInfo')
  paymentModalVisible.value = true
}

// Inspect the response of /api/order/plan (and /api/order/self/:id/pay)
// and route the UI to either the QR/redirect modal, the bank-transfer
// note modal, or a toast. Crucially, we never auto-navigate to /orders
// when the backend reports a warning — that was the old (broken)
// behaviour the new pay.status field exists to fix.
function handlePaymentResult(data, plan, payMethod) {
  const pay = data?.pay || {}
  const status = pay.status
  const url = pay.pay_url || pay.qr_code
  const amount = data?.amount || plan?.price
  const name = data?.plan_name || plan?.name

  if (status === 'success' && url) {
    openPaymentModal(url, amount, name, payMethod)
    return
  }
  if (status === 'success' && pay.note) {
    openBankTransferModal(amount, name, pay.note, payMethod)
    return
  }
  if (status === 'warning') {
    // Channel problem — let the user know, but don't drag them away
    // from the plans page.
    Message.error(pay.warning || t('plansPage.getPayInfoFailed'))
    return
  }
  // Unknown / missing status — fall back to legacy behaviour so old
  // backends keep working, but still avoid an unsolicited redirect to
  // /orders.
  if (url) {
    openPaymentModal(url, amount, name, payMethod)
  } else if (pay.note) {
    openBankTransferModal(amount, name, pay.note, payMethod)
  } else {
    Message.error(pay.warning || t('plansPage.getPayInfoFailed'))
  }
}

async function confirmUpgrade() {
  if (!targetPlan.value) return
  upgrading.value = true
  try {
    const res = await orderApi.createPlan({
      plan_id: targetPlan.value.id,
      pay_method: selectedPayMethod.value,
    })
    const data = res?.data
    if (!data?.success) {
      // Surface the backend error (e.g. "no payment channel enabled").
      Message.error(data?.message || NO_PAYMENT_MSG.value)
      return
    }
    upgradeModalVisible.value = false
    handlePaymentResult(data, targetPlan.value, selectedPayMethod.value)
  } catch (e) {
    Message.error(e.response?.data?.message || t('plansPage.upgradeFailed'))
  } finally {
    upgrading.value = false
  }
}

async function confirmPurchase() {
  if (!selectedPlan.value) return
  purchasing.value = true
  try {
    const res = await orderApi.createPlan({
      plan_id: selectedPlan.value.id,
      pay_method: selectedPayMethod.value,
    })
    const data = res?.data
    if (!data?.success) {
      // Surface the backend error (e.g. "no payment channel enabled").
      Message.error(data?.message || NO_PAYMENT_MSG.value)
      return
    }
    purchaseModalVisible.value = false
    handlePaymentResult(data, selectedPlan.value, selectedPayMethod.value)
  } catch (e) {
    Message.error(e.response?.data?.message || t('plansPage.createOrderFailed'))
  } finally {
    purchasing.value = false
  }
}

onMounted(() => { loadPlans() })
</script>

<style scoped>
.page-container {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

/* 升级提示横幅 */
.upgrade-banner {
  background: linear-gradient(135deg, #EFF6FF 0%, #DBEAFE 100%);
  border: 1px solid #BFDBFE;
  border-radius: 10px;
  padding: 14px 20px;
  margin-bottom: 24px;
  font-size: 14px;
  color: #1D40AE;
  display: flex;
  align-items: center;
  gap: 10px;
}
.upgrade-banner strong { font-weight: 700; }
.banner-icon { font-size: 18px; }

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 80px 40px;
  color: #86868B;
}
.empty-illustration { margin-bottom: 16px; display: flex; justify-content: center; }
.empty-text {
  font-size: 15px;
  color: #1D1D1F;
  font-weight: 600;
  margin-bottom: 6px;
}
.empty-desc {
  font-size: 13px;
  color: #86868B;
  margin: 0;
}

/* 套餐卡片网格 */
.plans-spin { display: block; width: 100%; min-height: 300px; }
.plans-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  align-items: stretch;
}
.plan-card {
  background: #fff;
  border: 1px solid #E5E5E7;
  border-radius: 14px;
  padding: 24px 20px;
  position: relative;
  transition: all 0.25s;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
}
.plan-card:hover {
  border-color: #007AFF;
  box-shadow: 0 8px 24px rgba(0, 122, 255, 0.1);
}
.plan-card.recommended {
  border-color: #007AFF;
  box-shadow: 0 0 0 1px #007AFF;
}
.plan-card.is-current {
  border-color: #10B981;
  box-shadow: 0 0 0 1px #10B981;
}

.plan-badge {
  position: absolute;
  top: -11px;
  left: 50%;
  transform: translateX(-50%);
  background: #007AFF;
  color: #fff;
  font-size: 11px;
  font-weight: 600;
  padding: 3px 14px;
  border-radius: 100px;
  white-space: nowrap;
}
.current-badge {
  position: absolute;
  top: -11px;
  left: 50%;
  transform: translateX(-50%);
  background: #10B981;
  color: #fff;
  font-size: 11px;
  font-weight: 600;
  padding: 3px 14px;
  border-radius: 100px;
  white-space: nowrap;
}

.plan-name {
  font-size: 18px;
  font-weight: 700;
  color: #1D1D1F;
  margin-bottom: 10px;
  text-align: center;
}
.plan-multiplier {
  font-size: 12px;
  color: #007AFF;
  text-align: center;
  margin-bottom: 4px;
}
.plan-price {
  display: flex;
  align-items: baseline;
  justify-content: center;
  gap: 2px;
  margin-bottom: 6px;
}
.plan-price-symbol { font-size: 18px; font-weight: 600; color: #86868B; }
.plan-price-num { font-size: 40px; font-weight: 800; color: #1D1D1F; line-height: 1; }
.plan-price-period { font-size: 13px; color: #86868B; margin-left: 2px; }
.plan-desc {
  font-size: 13px;
  color: #86868B;
  text-align: center;
  margin-bottom: 0;
  line-height: 1.5;
  min-height: 40px;
}
.plan-divider {
  border: none;
  border-top: 1px solid #F2F2F7;
  margin: 16px 0;
}
.plan-features {
  list-style: none;
  padding: 0;
  margin: 0 0 20px;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.plan-features li {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 14px;
  color: #1D1D1F;
}
.plan-features li::before {
  content: '✓';
  color: #007AFF;
  font-weight: 700;
  flex-shrink: 0;
  margin-top: 1px;
}
.plan-btn-wrap { margin-top: auto; }
.plan-btn {
  width: 100%;
  height: 42px;
  border-radius: 6px;
  font-weight: 600;
}

/* 升级弹窗 */
.upgrade-modal-body { padding: 8px 0; }
.upgrade-compare { display: flex; align-items: stretch; gap: 0; }
.compare-card {
  flex: 1;
  border-radius: 12px;
  padding: 20px 16px;
  text-align: center;
}
.compare-card.left {
  background: linear-gradient(135deg, #F0FDF4, #DCFCE7);
  border: 1px solid #86EFAC;
}
.compare-card.right {
  background: linear-gradient(135deg, #EFF6FF, #DBEAFE);
  border: 1px solid #93C5FD;
}
.compare-label {
  font-size: 12px;
  font-weight: 600;
  color: #6B7280;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 8px;
}
.compare-arrow {
  width: 48px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}
.compare-name {
  font-size: 18px;
  font-weight: 700;
  margin-bottom: 8px;
}
.compare-days { font-size: 14px; color: #10B981; }
.compare-days.expired { color: #EF4444; }
.compare-price { font-size: 13px; color: #6B7280; }
.compare-price.highlight { color: #007AFF; font-weight: 600; }

.upgrade-diff {
  margin-top: 16px;
  background: #FFFBEB;
  border: 1px solid #FDE68A;
  border-radius: 10px;
  padding: 14px 16px;
}
.diff-title {
  font-size: 12px;
  font-weight: 600;
  color: #6B7280;
  text-transform: uppercase;
  margin-bottom: 10px;
}
.diff-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 6px;
  font-size: 14px;
}
.diff-row.total {
  border-top: 1px solid #E5E7EB;
  padding-top: 8px;
  margin-top: 8px;
}
.diff-row.total span:first-child { font-weight: 600; }
.final-price { font-size: 20px; font-weight: 800; color: #EF4444; }
.deduct { color: #10B981; }

.upgrade-note {
  background: #FEF3C7;
  border: 1px solid #FDE68A;
  border-radius: 8px;
  padding: 12px 14px;
  font-size: 13px;
  color: #92400E;
  margin-top: 16px;
  margin-bottom: 20px;
}
.upgrade-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 4px;
}

/* 购买弹窗 */
.purchase-modal-body { padding: 8px 0; }
.purchase-summary {
  display: flex;
  flex-direction: column;
  gap: 14px;
  margin-bottom: 20px;
}
.purchase-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
}
.purchase-label { font-size: 14px; color: #86868B; flex-shrink: 0; }
.purchase-value { font-size: 14px; color: #1D1D1F; text-align: right; }
.purchase-price { font-size: 16px; font-weight: 700; color: #007AFF; }
.purchase-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 4px;
}

/* 支付方式选择（购买 / 升级弹窗通用） */
.pay-picker {
  margin: 16px 0 20px;
}
.pay-picker-label {
  font-size: 13px;
  font-weight: 600;
  color: #6B7280;
  margin-bottom: 10px;
  letter-spacing: 0.3px;
}
.pay-picker-list {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}
.pay-picker-item {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px 12px;
  border: 1.5px solid #E5E5E7;
  border-radius: 10px;
  background: #fff;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  color: #1D1D1F;
  transition: all 0.15s;
}
.pay-picker-item:hover:not(:disabled) {
  border-color: #C7C7CC;
}
.pay-picker-item.active {
  border-width: 1.5px;
}
.pay-picker-item:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
.pay-picker-icon {
  flex-shrink: 0;
  font-size: 28px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}
/* WeChat 官方色：#07C160 */
.pay-picker-item.wechat .pay-picker-icon {
  color: #07C160;
}
.pay-picker-item.wechat.active {
  border-color: #07C160;
  background: rgba(7, 193, 96, 0.04);
  box-shadow: 0 0 0 2px rgba(7, 193, 96, 0.12);
}
.pay-picker-item.wechat.active .pay-picker-check {
  color: #07C160;
}
/* Alipay 官方色：#1677FF */
.pay-picker-item.alipay .pay-picker-icon {
  color: #1677FF;
}
.pay-picker-item.alipay.active {
  border-color: #1677FF;
  background: rgba(22, 119, 255, 0.04);
  box-shadow: 0 0 0 2px rgba(22, 119, 255, 0.12);
}
.pay-picker-item.alipay.active .pay-picker-check {
  color: #1677FF;
}
.pay-picker-name {
  flex: 1;
  text-align: left;
  display: inline-flex;
  align-items: center;
}
.pay-picker-check {
  flex-shrink: 0;
  font-size: 14px;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}

/* 响应式 */
@media (max-width: 1023px) {
  .plans-grid { grid-template-columns: repeat(2, 1fr); }
}
@media (max-width: 767px) {
  .plans-grid { grid-template-columns: 1fr; }
}

/* 支付二维码弹窗 */
.payment-modal-body {
  padding: 8px 0;
  text-align: center;
}
.payment-order-info {
  background: #F9FAFB;
  border-radius: 10px;
  padding: 14px 16px;
  margin-bottom: 16px;
  text-align: left;
}
.payment-info-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 6px;
  font-size: 14px;
}
.payment-info-value { color: #1D1D1F; font-weight: 600; }
.payment-info-price { color: #007AFF; font-weight: 700; font-size: 16px; }
.qrcode-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
  min-height: 220px;
}
.qrcode-img {
  width: 220px;
  height: 220px;
  border: 1px solid #E5E5E7;
  border-radius: 8px;
  display: block;
}
.qrcode-loading { color: #86868B; font-size: 14px; }
.payment-tip { font-size: 15px; font-weight: 600; color: #1D1D1F; margin-bottom: 4px; }
.payment-tip-sub { font-size: 12px; color: #86868B; margin-bottom: 20px; }
.payment-actions { display: flex; gap: 12px; }
</style>
