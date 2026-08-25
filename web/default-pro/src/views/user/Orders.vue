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
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useRouter } from 'vue-router'
import orderApi from '@/api/order'

const router = useRouter()

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

function payOrder(record) {
  Message.loading({ content: '跳转支付中…', duration: 1500 })
  router.push(`/plans?order=${record.orderNo}`)
}

onMounted(() => { loadOrders() })
</script>

<style scoped>
.page-container {
  padding: 24px;
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
</style>
