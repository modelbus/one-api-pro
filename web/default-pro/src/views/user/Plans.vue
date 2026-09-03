<template>
  <div class="page-container">
    <!-- 升级提示横幅 -->
    <div v-if="currentPlan && !isExpired(currentPlan)" class="upgrade-banner">
      <span class="banner-icon">💡</span>
      <span>您当前正在使用 <strong>{{ currentPlan.name }}</strong>，升级将补齐差价</span>
    </div>

    <!-- 套餐卡片网格 -->
    <a-spin :loading="loading" dot tip="加载中..." class="plans-spin">
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
      <p class="empty-text">暂无可用套餐</p>
      <p class="empty-desc">请联系管理员配置套餐后再来查看</p>
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
        <div v-if="plan.recommended" class="plan-badge">推荐</div>
        <div v-if="isCurrentPlan(plan)" class="current-badge">当前套餐</div>

        <div class="plan-name">{{ plan.name }}</div>
        <div v-if="plan.multiplier" class="plan-multiplier">{{ plan.multiplier }}</div>
        <div class="plan-price">
          <span class="plan-price-symbol">¥</span>
          <span class="plan-price-num">{{ plan.price }}</span>
          <span class="plan-price-period">/{{ plan.duration_text || '月' }}</span>
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
            <template v-if="isCurrentPlan(plan)">当前套餐</template>
            <template v-else-if="isLowerPlan(plan)">暂不可用</template>
            <template v-else-if="currentPlan && !isExpired(currentPlan)">升级到 {{ plan.name }}</template>
            <template v-else>立即订阅</template>
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
      title="确认升级"
      class="upgrade-modal"
      width="520px"
    >
      <div v-if="targetPlan" class="upgrade-modal-body">
        <div class="upgrade-compare">
          <div class="compare-card left">
            <div class="compare-label">当前套餐</div>
            <div class="compare-name">{{ currentPlan?.name }}</div>
            <div v-if="currentPlan && !isExpired(currentPlan)" class="compare-days">
              剩余 <strong>{{ getRemainingDays(currentPlan) }}</strong> 天
            </div>
            <div v-else class="compare-days expired">已过期</div>
            <div class="compare-price">¥{{ currentPlan?.price || 0 }}/{{ currentPlan?.duration_text || '月' }}</div>
          </div>
          <div class="compare-arrow">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
              <path d="M5 12h14M14 5l7 7-7 7" stroke="#007AFF" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <div class="compare-card right">
            <div class="compare-label">升级至</div>
            <div class="compare-name">{{ targetPlan.name }}</div>
            <div class="compare-days">{{ targetPlan.duration_text || '1 个月' }}</div>
            <div class="compare-price highlight">¥{{ targetPlan.price }}/{{ targetPlan.duration_text || '月' }}</div>
          </div>
        </div>

        <div class="upgrade-diff">
          <div class="diff-title">差价明细</div>
          <div class="diff-row">
            <span>新套餐价格</span>
            <span>¥{{ targetPlan.price }}</span>
          </div>
          <div class="diff-row">
            <span>当前套餐剩余价值</span>
            <span class="deduct">-¥{{ getRemainingValue() }}</span>
          </div>
          <div class="diff-row total">
            <span>应付差价</span>
            <span class="final-price">¥{{ getUpgradeDiff() }}</span>
          </div>
        </div>

        <div class="upgrade-note">
          ⚠️ 升级后将立即激活新套餐，当前套餐将自动取消
        </div>

        <div class="upgrade-footer">
          <a-button @click="upgradeModalVisible = false">取消</a-button>
          <a-button type="primary" :loading="upgrading" @click="confirmUpgrade">确认升级</a-button>
        </div>
      </div>
    </a-modal>

    <!-- 购买确认弹窗 -->
    <a-modal
      v-model:visible="purchaseModalVisible"
      :footer="false"
      :mask-closable="true"
      title="确认订阅"
      width="420px"
    >
      <div v-if="selectedPlan" class="purchase-modal-body">
        <div class="purchase-summary">
          <div class="purchase-row">
            <span class="purchase-label">套餐名称</span>
            <span class="purchase-value">{{ selectedPlan.name }}</span>
          </div>
          <div class="purchase-row">
            <span class="purchase-label">套餐价格</span>
            <span class="purchase-price">¥{{ selectedPlan.price }}/{{ selectedPlan.duration_text || '月' }}</span>
          </div>
          <div class="purchase-row">
            <span class="purchase-label">适用场景</span>
            <span class="purchase-value">{{ selectedPlan.description }}</span>
          </div>
        </div>
        <div class="purchase-footer">
          <a-button @click="purchaseModalVisible = false">取消</a-button>
          <a-button type="primary" :loading="purchasing" @click="confirmPurchase">确认购买</a-button>
        </div>
      </div>
    </a-modal>

    <!-- 支付二维码弹窗 -->
    <a-modal
      v-model:visible="paymentModalVisible"
      :footer="false"
      :mask-closable="true"
      title="微信扫码支付"
      width="420px"
      class="payment-modal"
    >
      <div class="payment-modal-body">
        <div class="qrcode-wrap">
          <img v-if="qrcodeDataUrl" :src="qrcodeDataUrl" alt="微信支付二维码" class="qrcode-img" />
          <div v-else class="qrcode-loading">正在生成支付二维码...</div>
        </div>
        <div class="payment-tip">请使用微信扫码支付</div>
        <div class="payment-tip-sub">（目前只支持微信支付）</div>
        <div class="payment-order-info">
          <div class="payment-info-row">
            <span>套餐</span>
            <span class="payment-info-value">{{ paymentPackageName }}</span>
          </div>
          <div class="payment-info-row">
            <span>金额</span>
            <span class="payment-info-price">¥{{ paymentAmount }}</span>
          </div>
        </div>
        <div class="payment-actions">
          <a-button long @click="paymentModalVisible = false">关闭</a-button>
          <a-button type="primary" long @click="$router.push('/orders')">前往订单中心</a-button>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import QRCode from 'qrcode'
import planApi from '@/api/plan'
import orderApi from '@/api/order'
import paymentApi from '@/api/payment'

const router = useRouter()

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

const NO_PAYMENT_MSG = '系统尚未开通任何支付通道，请设置后开启支付'

const VALID_PLAN_NAMES = ['lite', 'air', 'pro', 'max']

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
    const filtered = list.filter(p => {
      const name = (p.name || '').toLowerCase().trim()
      return VALID_PLAN_NAMES.includes(name)
    })
    planList.value = filtered.map((item, index) => ({
      id: item.id,
      name: item.name,
      price: item.price,
      description: item.description,
      recommended: item.recommended,
      features: parseFeatures(item.features),
      multiplier: item.multiplier,
      sort: item.sort ?? index + 1,
      duration_days: item.duration_days ?? 30,
      duration_text: item.duration_text || '1个月',
    }))
  } catch (e) {
    Message.error('加载套餐列表失败')
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
// safe to continue.
async function ensurePaymentEnabled() {
  try {
    const { data } = await paymentApi.status()
    const anyEnabled = data?.data?.any_enabled
    if (!anyEnabled) {
      Message.error(NO_PAYMENT_MSG)
      return false
    }
    return true
  } catch (e) {
    Message.error('无法获取支付通道状态，请稍后重试')
    return false
  }
}

function handleUpgrade(plan) {
  const currentSort = currentPlan.value?.sort ?? 0
  if (currentSort >= plan.sort) {
    Message.warning({ content: '不能降级到低级别套餐', duration: 3000 })
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
    console.error('生成二维码失败', e)
  }
}

async function openPaymentModal(codeUrl, amount, packageName) {
  paymentCodeUrl.value = codeUrl
  paymentAmount.value = amount
  paymentPackageName.value = packageName
  await generateQRCode(codeUrl)
  paymentModalVisible.value = true
}

async function confirmUpgrade() {
  if (!targetPlan.value) return
  upgrading.value = true
  try {
    const res = await orderApi.createPlan({
      plan_id: targetPlan.value.id,
      pay_method: 'wechat',
    })
    upgradeModalVisible.value = false
    const data = res?.data
    if (!data?.success) {
      // Surface the backend error (e.g. "no payment channel enabled").
      Message.error(data?.message || NO_PAYMENT_MSG)
      return
    }
    const pay = data.pay || {}
    const url = pay.pay_url || pay.qr_code
    if (url) {
      openPaymentModal(url, data.amount || targetPlan.value.price, data.plan_name || targetPlan.value.name)
    } else {
      Message.warning(pay.warning || '获取支付二维码失败，请稍后到订单中心查看')
      router.push('/orders')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || '升级失败')
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
      pay_method: 'wechat',
    })
    purchaseModalVisible.value = false
    const data = res?.data
    if (!data?.success) {
      // Surface the backend error (e.g. "no payment channel enabled").
      Message.error(data?.message || NO_PAYMENT_MSG)
      return
    }
    const pay = data.pay || {}
    const url = pay.pay_url || pay.qr_code
    if (url) {
      openPaymentModal(url, data.amount || selectedPlan.value.price, data.plan_name || selectedPlan.value.name)
    } else {
      Message.warning(pay.warning || '获取支付二维码失败，请稍后到订单中心查看')
      router.push('/orders')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || '创建订单失败')
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
