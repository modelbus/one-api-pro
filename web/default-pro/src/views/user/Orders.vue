<template>
  <div class="orders-page">
    <h2 class="page-title">订单</h2>

    <a-tabs v-model:active-key="activeTab" @change="loadData">
      <a-tab-pane key="all" title="全部"></a-tab-pane>
      <a-tab-pane key="1" title="套餐订单"></a-tab-pane>
      <a-tab-pane key="2" title="充值订单"></a-tab-pane>
    </a-tabs>

    <a-spin :loading="loading" style="width: 100%">
      <a-table
        row-key="id"
        :columns="columns"
        :data="orders"
        :pagination="false"
        stripe
      >
        <template #type="{ record }">
          <a-tag :color="record.type === 1 ? 'arcoblue' : 'purple'" size="small">
            {{ record.type === 1 ? '套餐' : '充值' }}
          </a-tag>
        </template>
        <template #source="{ record }">
          <a-tag :color="record.source === 1 ? 'gray' : 'green'" size="small">
            {{ record.source === 1 ? '自助' : '管理员' }}
          </a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="statusColor(record.status)" size="small">
            {{ statusText(record.status) }}
          </a-tag>
        </template>
        <template #pay_method="{ record }">
          <span v-if="record.pay_method">{{ payMethodText(record.pay_method) }}</span>
          <span v-else class="muted">-</span>
        </template>
        <template #amount="{ record }">
          <span class="amount">¥{{ record.amount }}</span>
        </template>
        <template #actions="{ record }">
          <a-button v-if="record.status === 0" type="text" size="small" status="danger" @click="cancelOrder(record)">取消</a-button>
          <a-button v-else type="text" size="small" @click="viewOrder(record)">查看</a-button>
        </template>
      </a-table>
    </a-spin>

    <a-modal
      v-model:visible="detailVisible"
      title="订单详情"
      :footer="false"
      width="520px"
    >
      <a-descriptions v-if="detail" :data="detail" :column="1" bordered>
        <a-descriptions-item label="订单号">{{ detail.order_no }}</a-descriptions-item>
        <a-descriptions-item label="类型">{{ detail.type === 1 ? '套餐订单' : '充值订单' }}</a-descriptions-item>
        <a-descriptions-item label="来源">{{ detail.source === 1 ? '用户自助' : '管理员' }}</a-descriptions-item>
        <a-descriptions-item label="金额">¥{{ detail.amount }}</a-descriptions-item>
        <a-descriptions-item label="支付方式">{{ payMethodText(detail.pay_method) }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="statusColor(detail.status)" size="small">{{ statusText(detail.status) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="下单时间">{{ formatTime(detail.create_time) }}</a-descriptions-item>
        <a-descriptions-item label="支付时间">{{ detail.pay_time ? formatTime(detail.pay_time) : '-' }}</a-descriptions-item>
        <a-descriptions-item label="流水号">{{ detail.pay_trade_no || '-' }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import orderApi from '@/api/order'

const activeTab = ref('all')
const loading = ref(false)
const orders = ref([])
const detailVisible = ref(false)
const detail = ref(null)

const columns = [
  { title: '订单号', dataIndex: 'order_no', width: 220 },
  { title: '类型', slotName: 'type', width: 80 },
  { title: '来源', slotName: 'source', width: 80 },
  { title: '金额', slotName: 'amount', width: 100 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '支付方式', slotName: 'pay_method', width: 100 },
  { title: '下单时间', dataIndex: 'create_time', width: 180, render: ({ record }) => formatTime(record.create_time) },
  { title: '操作', slotName: 'actions', width: 100, fixed: 'right' },
]

function statusText(s) {
  switch (s) {
    case 0: return '待支付'
    case 1: return '已支付'
    case 2: return '已取消'
    case 3: return '已退款'
    default: return '未知'
  }
}

function statusColor(s) {
  switch (s) {
    case 0: return 'orange'
    case 1: return 'green'
    case 2: return 'gray'
    case 3: return 'red'
    default: return 'gray'
  }
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
  const d = new Date(t * 1000)
  return d.toLocaleString('zh-CN')
}

async function loadData() {
  loading.value = true
  try {
    const params = activeTab.value === 'all' ? {} : { type: Number(activeTab.value) }
    const { data } = await orderApi.myOrders(params.type)
    if (data.success) {
      orders.value = data.data || []
    } else {
      Message.error(data.message || '加载失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || '网络错误')
  } finally {
    loading.value = false
  }
}

async function viewOrder(record) {
  try {
    const { data } = await orderApi.myOrder(record.id)
    if (data.success) {
      detail.value = data.data
      detailVisible.value = true
    } else {
      Message.error(data.message || '加载失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || '网络错误')
  }
}

async function cancelOrder(record) {
  try {
    const { data } = await orderApi.cancelMyOrder(record.id)
    if (data.success) {
      Message.success('已取消')
      loadData()
    } else {
      Message.error(data.message || '取消失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || '网络错误')
  }
}

onMounted(loadData)
</script>

<style scoped>
.orders-page { padding: 16px; }
.page-title { font-size: 20px; font-weight: 600; margin-bottom: 16px; color: var(--color-text-1); }
.amount { font-weight: 500; color: #165dff; }
.muted { color: var(--color-text-4); }
</style>
