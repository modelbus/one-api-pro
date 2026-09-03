<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">我的订单</h1>
      <p class="page-subtitle">查看和管理您的所有订单</p>
    </div>

    <!-- Tabs (按状态过滤) -->
    <div class="tabs">
      <button
        v-for="tab in tabs"
        :key="tab.value"
        class="tab"
        :class="{ active: activeTab === tab.value }"
        @click="activeTab = tab.value; loadOrders()"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- Table -->
    <a-spin :loading="loading" dot tip="加载中..." class="orders-spin">
    <div class="table-wrap">
      <table class="data-table">
        <thead>
          <tr>
            <th>订单号</th>
            <th>套餐</th>
            <th>金额</th>
            <th>创建时间</th>
            <th>状态</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="order in filteredOrders" :key="order.id">
            <td><span class="order-no">{{ order.orderNo }}</span></td>
            <td>{{ order.planName }}</td>
            <td><span class="amount">¥{{ order.amount }}</span></td>
            <td>{{ order.createdAt }}</td>
            <td>
              <span class="badge" :class="statusClass(order.status)">
                {{ order.statusText }}
              </span>
            </td>
            <td>
              <a-space>
                <a-button
                  v-if="order.status === 0"
                  type="primary"
                  size="small"
                  @click="payOrder(order)"
                >支付</a-button>
                <a-button size="small" @click="viewOrder(order)">查看</a-button>
              </a-space>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-if="filteredOrders.length === 0" class="empty-state">
        <div class="empty-illustration">
          <svg width="80" height="80" viewBox="0 0 80 80" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect x="10" y="20" width="60" height="45" rx="4" stroke="#D1D1D6" stroke-width="2" fill="none"/>
            <rect x="16" y="28" width="30" height="4" rx="2" fill="#D1D1D6"/>
            <rect x="16" y="38" width="48" height="4" rx="2" fill="#E5E5E7"/>
            <rect x="16" y="48" width="40" height="4" rx="2" fill="#E5E5E7"/>
            <circle cx="58" cy="52" r="10" fill="#F7F8FA" stroke="#D1D1D6" stroke-width="2"/>
            <line x1="52" y1="52" x2="64" y2="52" stroke="#D1D1D6" stroke-width="2"/>
            <line x1="58" y1="46" x2="58" y2="58" stroke="#D1D1D6" stroke-width="2"/>
          </svg>
        </div>
        <p class="empty-text">暂无订单记录</p>
        <a-button type="primary" @click="$router.push('/plans')">立即订阅</a-button>
      </div>
    </div>
    </a-spin>

    <!-- 详情弹窗 -->
    <div v-if="detailVisible" class="modal-overlay" @click.self="detailVisible = false">
      <div class="modal">
        <div class="modal-header">
          <span class="modal-title">订单详情</span>
          <button class="modal-close" @click="detailVisible = false">×</button>
        </div>
        <div class="modal-body" v-if="currentOrder">
          <div class="detail-row">
            <span class="detail-label">订单号</span>
            <span class="detail-value">{{ currentOrder.orderNo }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">订单类型</span>
            <span class="detail-value">{{ currentOrder.type === 1 ? '套餐订单' : currentOrder.type === 2 ? '充值订单' : '其他' }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">订单来源</span>
            <span class="detail-value">{{ currentOrder.source === 1 ? '用户自助' : currentOrder.source === 2 ? '管理员' : '其他' }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">套餐</span>
            <span class="detail-value">{{ currentOrder.planName }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">金额</span>
            <span class="amount">¥{{ currentOrder.amount }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">支付方式</span>
            <span class="detail-value">{{ currentOrder.payMethodLabel }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">创建时间</span>
            <span class="detail-value">{{ currentOrder.createdAt }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">支付时间</span>
            <span class="detail-value">{{ currentOrder.paidAt }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">状态</span>
            <span class="badge" :class="statusClass(currentOrder.status)">
              {{ currentOrder.statusText }}
            </span>
          </div>
          <div v-if="currentOrder.pay_trade_no" class="detail-row">
            <span class="detail-label">流水号</span>
            <span class="detail-value">{{ currentOrder.pay_trade_no }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 支付方式选择弹窗 -->
    <a-modal
      v-model:visible="payModalVisible"
      :footer="false"
      :mask-closable="true"
      title="选择支付方式"
      width="420px"
    >
      <div class="pay-method-body">
        <div class="pay-method-row">
          <span class="pay-method-label">套餐</span>
          <span class="pay-method-value">{{ payTarget?.planName }}</span>
        </div>
        <div class="pay-method-row">
          <span class="pay-method-label">金额</span>
          <span class="pay-method-amount">¥{{ payTarget?.amount }}</span>
        </div>
        <div class="pay-method-list">
          <button
            v-for="m in payMethods"
            :key="m.name"
            class="pay-method-item"
            :class="[m.name, { active: selectedPayMethod === m.name }]"
            @click="selectedPayMethod = m.name"
          >
            <component :is="methodIconComponent(m.name)" :size="28" class="pay-method-icon" />
            <span class="pay-method-name">{{ m.label }}</span>
            <span v-if="selectedPayMethod === m.name" class="pay-method-check">✓</span>
          </button>
          <div v-if="payMethods.length === 0" class="pay-method-empty">
            暂无可用的支付方式
          </div>
        </div>
        <div class="pay-method-footer">
          <a-button @click="payModalVisible = false">取消</a-button>
          <a-button type="primary" :loading="paySubmitting" :disabled="!selectedPayMethod" @click="confirmPayMethod">
            确认支付
          </a-button>
        </div>
      </div>
    </a-modal>

    <!-- 支付二维码弹窗 -->
    <a-modal
      v-model:visible="qrModalVisible"
      :footer="false"
      :mask-closable="true"
      :title="qrModalTitle"
      width="420px"
      class="payment-modal"
    >
      <div class="payment-modal-body">
        <div class="qrcode-wrap">
          <img v-if="qrcodeDataUrl" :src="qrcodeDataUrl" alt="支付二维码" class="qrcode-img" />
          <div v-else class="qrcode-loading">正在生成支付二维码...</div>
        </div>
        <div class="payment-tip">{{ qrModalTip }}</div>
        <div v-if="qrModalNote" class="payment-tip-sub">{{ qrModalNote }}</div>
        <div class="payment-order-info">
          <div class="payment-info-row">
            <span>套餐</span>
            <span class="payment-info-value">{{ qrPackageName }}</span>
          </div>
          <div class="payment-info-row">
            <span>金额</span>
            <span class="payment-info-price">¥{{ qrAmount }}</span>
          </div>
        </div>
        <div class="payment-actions">
          <a-button long @click="qrModalVisible = false">关闭</a-button>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconWechatpay, IconAlipayCircle, IconSafe } from '@arco-design/web-vue/es/icon'
import QRCode from 'qrcode'
import orderApi from '@/api/order'
import paymentApi from '@/api/payment'

const activeTab = ref('all')
const detailVisible = ref(false)
const currentOrder = ref(null)
const loading = ref(false)
const orders = ref([])

const tabs = [
  { label: '全部', value: 'all' },
  { label: '待支付', value: 'pending' },
  { label: '已支付', value: 'paid' },
  { label: '已取消', value: 'cancelled' },
]

const STATUS_STR = { 0: 'pending', 1: 'paid', 2: 'cancelled', 3: 'refunded' }

const NO_PAYMENT_MSG = '系统尚未开通任何支付通道，请设置后开启支付'

const filteredOrders = computed(() => {
  if (activeTab.value === 'all') return orders.value
  return orders.value.filter(o => STATUS_STR[o.status] === activeTab.value)
})

function statusText(status) {
  const map = { 0: '待支付', 1: '已支付', 2: '已取消', 3: '已退款' }
  return map[status] || String(status)
}

function statusClass(status) {
  const map = { 0: 'badge-warning', 1: 'badge-success', 2: 'badge-muted', 3: 'badge-muted' }
  return map[status] || 'badge-muted'
}

function payMethodText(m) {
  switch (m) {
    case 'wechat': return '微信'
    case 'alipay': return '支付宝'
    case 'bank': return '银行'
    case 'offline': return '线下'
    case 'free': return '免费'
    default: return m || '-'
  }
}

// Returns an Arco icon component matching the pay method so the brand
// color (CSS-driven via currentColor) shows up in the picker.
function methodIconComponent(name) {
  switch (name) {
    case 'wechat': return IconWechatpay
    case 'alipay': return IconAlipayCircle
    default: return IconSafe
  }
}

function formatTime(t) {
  if (!t) return '-'
  const d = new Date(typeof t === 'number' ? t * 1000 : new Date(String(t).replace(' ', 'T')).getTime())
  if (isNaN(d.getTime())) return String(t)
  const pad = n => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

async function loadOrders() {
  loading.value = true
  try {
    // Backend returns { success, message, data: [...] }.
    const res = await orderApi.myOrders()
    const body = res?.data || {}
    if (body.success) {
      const list = Array.isArray(body.data) ? body.data : []
      orders.value = list.map(item => ({
        id: item.id,
        orderNo: item.order_no || '',
        planName: item.plan_info ? safeParseName(item.plan_info) : '',
        amount: item.amount,
        status: item.status,
        statusText: statusText(item.status),
        type: item.type,
        source: item.source,
        payMethod: item.pay_method,
        payMethodLabel: payMethodText(item.pay_method),
        payTradeNo: item.pay_trade_no,
        createdAt: formatTime(item.create_time),
        paidAt: formatTime(item.pay_time),
      }))
    } else {
      Message.error(body.message || '加载失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || '网络错误')
  } finally {
    loading.value = false
  }
}

function safeParseName(planInfo) {
  if (!planInfo) return ''
  try {
    const obj = typeof planInfo === 'string' ? JSON.parse(planInfo) : planInfo
    return obj.name || ''
  } catch { return '' }
}

function viewOrder(record) {
  currentOrder.value = record
  detailVisible.value = true
}

// Payment flow state
const payModalVisible = ref(false)
const payTarget = ref(null)       // { id, planName, amount }
const payMethods = ref([])        // [{ name, label, enabled }]
const selectedPayMethod = ref('')
const paySubmitting = ref(false)

const qrModalVisible = ref(false)
const qrModalTitle = ref('扫码支付')
const qrModalTip = ref('请扫码完成支付')
const qrModalNote = ref('')
const qrcodeDataUrl = ref('')
const qrAmount = ref(0)
const qrPackageName = ref('')

async function loadPaymentStatus() {
  try {
    const { data } = await paymentApi.status()
    const d = data?.data || {}
    payMethods.value = (d.methods || []).filter(m => m.enabled)
    return !!d.any_enabled
  } catch (e) {
    payMethods.value = []
    return false
  }
}

async function payOrder(record) {
  // Pre-flight: refuse the action entirely when no payment channel is
  // enabled. This avoids redirecting the user to the (misleading)
  // plans page.
  const anyEnabled = await loadPaymentStatus()
  if (!anyEnabled) {
    Message.error(NO_PAYMENT_MSG)
    return
  }
  payTarget.value = {
    id: record.id,
    planName: record.planName,
    amount: record.amount,
  }
  // Default to the order's existing method if it's still enabled,
  // otherwise the first enabled method.
  const exists = payMethods.value.find(m => m.name === record.payMethod)
  selectedPayMethod.value = exists ? exists.name : (payMethods.value[0]?.name || '')
  payModalVisible.value = true
}

// Build the QR / redirect modal for an online payment channel.
async function openQrModal(url, payMethod) {
  qrModalTitle.value = payMethod === 'alipay' ? '支付宝扫码支付' : '微信扫码支付'
  qrModalTip.value = payMethod === 'alipay' ? '请使用支付宝扫码支付' : '请使用微信扫码支付'
  qrModalNote.value = ''
  qrcodeDataUrl.value = ''
  try {
    qrcodeDataUrl.value = await QRCode.toDataURL(url, {
      width: 220,
      margin: 2,
      color: { dark: '#000000', light: '#ffffff' },
    })
  } catch (e) {
    qrcodeDataUrl.value = ''
  }
  qrModalVisible.value = true
}

// Build the "transfer info" modal for bank / offline payments.
function openBankNoteModal(note) {
  qrModalTitle.value = '转账信息'
  qrModalTip.value = note
  qrModalNote.value = ''
  qrcodeDataUrl.value = ''
  qrModalVisible.value = true
}

// Inspect the response of /api/order/self/:id/pay and route the UI.
// Mirrors the logic in Plans.vue — driven by pay.status so we never
// silently drop the user into the orders page when the payment
// channel has a problem.
function handlePayMyOrderResult(data) {
  const pay = data?.pay || {}
  const status = pay.status
  const url = pay.pay_url || pay.qr_code
  qrAmount.value = data?.amount || payTarget.value?.amount
  qrPackageName.value = data?.plan_name || payTarget.value?.planName

  if (status === 'success' && url) {
    openQrModal(url, selectedPayMethod.value)
    return
  }
  if (status === 'success' && pay.note) {
    openBankNoteModal(pay.note)
    return
  }
  if (status === 'warning') {
    Message.error(pay.warning || '发起支付失败，请稍后重试')
    return
  }
  // Backwards-compat fallback when pay.status is absent (older backends).
  if (url) {
    openQrModal(url, selectedPayMethod.value)
  } else if (pay.note) {
    openBankNoteModal(pay.note)
  } else {
    Message.error(pay.warning || '发起支付失败，请稍后重试')
  }
}

async function confirmPayMethod() {
  if (!payTarget.value || !selectedPayMethod.value) return
  paySubmitting.value = true
  try {
    const res = await orderApi.payMyOrder(payTarget.value.id, {
      pay_method: selectedPayMethod.value,
    })
    const data = res?.data
    if (!data?.success) {
      Message.error(data?.message || '发起支付失败')
      return
    }
    payModalVisible.value = false
    handlePayMyOrderResult(data)
    // Refresh the order list so the row updates (e.g. once the
    // payment channel marks the order paid asynchronously).
    if (data?.pay?.status === 'success') {
      loadOrders()
    }
  } catch (e) {
    Message.error(e.response?.data?.message || '发起支付失败')
  } finally {
    paySubmitting.value = false
  }
}

onMounted(() => { loadOrders() })
</script>

<style scoped>
.page-container {
  max-width: 1200px;
  margin: 0 auto;
}
.page-header { margin-bottom: 24px; }
.page-title {
  font-size: 22px;
  font-weight: 700;
  color: #1D1D1F;
  margin: 0 0 4px;
}
.page-subtitle { font-size: 14px; color: #86868B; margin: 0; }

.tabs {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 4px;
  border-bottom: 1px solid #E5E5E7;
  margin-bottom: 20px;
}
.tab {
  padding: 10px 18px;
  font-size: 14px;
  font-weight: 500;
  color: #86868B;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;
  transition: all 0.15s;
  background: none;
  border-top: none;
  border-left: none;
  border-right: none;
}
.tab:hover { color: #1D1D1F; }
.tab.active { color: #007AFF; border-bottom-color: #007AFF; }

.orders-spin { display: block; width: 100%; min-height: 300px; }
.table-wrap {
  background: #fff;
  border-radius: 14px;
  border: 1px solid #E5E5E7;
  overflow: hidden;
}
.data-table {
  width: 100%;
  border-collapse: collapse;
}
.data-table th {
  background: #F7F8FA;
  padding: 12px 16px;
  text-align: left;
  font-size: 12px;
  font-weight: 600;
  color: #86868B;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border-bottom: 1px solid #E5E5E7;
}
.data-table td {
  padding: 14px 16px;
  font-size: 14px;
  color: #1D1D1F;
  border-bottom: 1px solid #F2F2F7;
  vertical-align: middle;
}
.data-table tr:last-child td { border-bottom: none; }
.data-table tr:hover td { background: #F7F8FA; }

.order-no {
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 13px;
  color: #86868B;
}
.amount { font-weight: 700; color: #007AFF; }

.badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 100px;
  font-size: 12px;
  font-weight: 500;
}
.badge-success { background: #D1F7D6; color: #1A7A38; }
.badge-warning { background: #FFE8B8; color: #8A5A00; }
.badge-muted { background: #F7F8FA; color: #86868B; }

.empty-state {
  text-align: center;
  padding: 60px 40px;
  color: #86868B;
}
.empty-illustration { margin-bottom: 16px; }
.empty-text {
  font-size: 15px;
  color: #86868B;
  margin-bottom: 16px;
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}
.modal {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 440px;
}
.modal-header {
  padding: 20px 24px 16px;
  border-bottom: 1px solid #F2F2F7;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.modal-title {
  font-size: 17px;
  font-weight: 600;
  color: #1D1D1F;
}
.modal-close {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: #AEAEB2;
  padding: 4px;
  line-height: 1;
}
.modal-close:hover { color: #1D1D1F; }
.modal-body { padding: 24px; }
.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #F2F2F7;
}
.detail-row:last-child { border-bottom: none; }
.detail-label { font-size: 13px; color: #86868B; }
.detail-value { font-size: 13px; color: #1D1D1F; font-weight: 500; }

/* 支付方式选择 */
.pay-method-body { padding: 4px 0; }
.pay-method-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
}
.pay-method-label { font-size: 13px; color: #86868B; }
.pay-method-value { font-size: 14px; color: #1D1D1F; font-weight: 500; }
.pay-method-amount { font-size: 18px; font-weight: 700; color: #007AFF; }

.pay-method-list {
  margin: 12px 0 18px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}
.pay-method-item {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 14px 10px;
  border: 1.5px solid #E5E5E7;
  border-radius: 10px;
  background: #fff;
  cursor: pointer;
  transition: all 0.15s;
  font-size: 14px;
  color: #1D1D1F;
  font-weight: 500;
}
.pay-method-item:hover {
  border-color: #C7C7CC;
}
/* WeChat 官方色：#07C160 */
.pay-method-item.wechat .pay-method-icon { color: #07C160; }
.pay-method-item.wechat.active {
  border-color: #07C160;
  background: rgba(7, 193, 96, 0.04);
  box-shadow: 0 0 0 2px rgba(7, 193, 96, 0.12);
}
.pay-method-item.wechat.active .pay-method-check { color: #07C160; }
/* Alipay 官方色：#1677FF */
.pay-method-item.alipay .pay-method-icon { color: #1677FF; }
.pay-method-item.alipay.active {
  border-color: #1677FF;
  background: rgba(22, 119, 255, 0.04);
  box-shadow: 0 0 0 2px rgba(22, 119, 255, 0.12);
}
.pay-method-item.alipay.active .pay-method-check { color: #1677FF; }
/* Other / bank */
.pay-method-item.bank .pay-method-icon { color: #5856D6; }
.pay-method-item.bank.active {
  border-color: #5856D6;
  background: rgba(88, 86, 214, 0.04);
  box-shadow: 0 0 0 2px rgba(88, 86, 214, 0.12);
}
.pay-method-item.bank.active .pay-method-check { color: #5856D6; }

.pay-method-icon {
  flex-shrink: 0;
  font-size: 28px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}
.pay-method-name {
  flex: 1;
  text-align: left;
  display: inline-flex;
  align-items: center;
}
.pay-method-check {
  flex-shrink: 0;
  font-size: 14px;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}
.pay-method-empty {
  grid-column: 1 / -1;
  text-align: center;
  color: #86868B;
  font-size: 13px;
  padding: 20px 0;
}
.pay-method-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 4px;
}

/* 二维码弹窗 */
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
.payment-info-row:last-child { margin-bottom: 0; }
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
