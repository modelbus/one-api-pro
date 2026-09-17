<!--
  AdminOrders 管理员订单列表页
  Admin-side order list page (channel-style list + filters + modals)

  版本: v0.0.13
  日期: 2026-09-08
  作者: opencode
-->
<template>
  <div class="admin-orders-page">
    <!-- 顶部欢迎条 -->
    <div class="welcome-bar">
      <div class="welcome-text">
        <h1 class="welcome-title">{{ $t('adminOrdersPage.title') }}</h1>
        <p class="welcome-desc">{{ $t('adminOrdersPage.subtitle') }}</p>
      </div>
      <div class="welcome-meta">
        <span class="meta-chip">{{ $t('adminOrdersPage.total', { n: totalCountForPager }) }}</span>
      </div>
    </div>

    <!-- 独立搜索栏 -->
    <div class="search-card">
      <div class="search-left">
        <a-input-search
          v-model="searchKeyword"
          :placeholder="$t('adminOrdersPage.searchPlaceholder')"
          allow-clear
          style="width: 240px"
          @search="handleSearch"
          @clear="handleSearchClear"
        />
        <a-select v-model="filterType" :placeholder="$t('adminOrdersPage.typePlaceholder')" allow-clear style="width: 110px" @change="handleFilterChange">
          <a-option :value="1">{{ $t('adminOrdersPage.typePlan') }}</a-option>
          <a-option :value="2">{{ $t('adminOrdersPage.typeTopup') }}</a-option>
        </a-select>
        <a-select v-model="filterStatus" :placeholder="$t('adminOrdersPage.statusPlaceholder')" allow-clear style="width: 120px" @change="handleFilterChange">
          <a-option :value="0">{{ $t('adminOrdersPage.statusPending') }}</a-option>
          <a-option :value="1">{{ $t('adminOrdersPage.statusPaid') }}</a-option>
          <a-option :value="2">{{ $t('adminOrdersPage.statusCanceled') }}</a-option>
          <a-option :value="3">{{ $t('adminOrdersPage.statusRefunded') }}</a-option>
        </a-select>
        <a-select v-model="filterSource" :placeholder="$t('adminOrdersPage.sourcePlaceholder')" allow-clear style="width: 120px" @change="handleFilterChange">
          <a-option :value="1">{{ $t('adminOrdersPage.sourceSelf') }}</a-option>
          <a-option :value="2">{{ $t('adminOrdersPage.sourceAdmin') }}</a-option>
        </a-select>
      </div>
      <div class="search-right">
        <a-button @click="handleResetFilters" :disabled="!hasActiveFilter">{{ $t('adminOrdersPage.resetFilters') }}</a-button>
        <a-button @click="handleRefresh">
          <template #icon><icon-refresh :size="14" /></template>
          {{ $t('adminOrdersPage.refresh') }}
        </a-button>
      </div>
    </div>

    <!-- 列表 -->
    <div class="list-wrap">
      <div v-if="orders.length === 0 && !loading" class="empty-state">
        <div class="empty-icon">
          <icon-storage :size="32" />
        </div>
        <p class="empty-title">{{ $t('adminOrdersPage.emptyTitle') }}</p>
        <p class="empty-desc">{{ $t('adminOrdersPage.emptyDesc') }}</p>
        <a-button v-if="hasActiveFilter" type="primary" @click="handleResetFilters">{{ $t('adminOrdersPage.resetFilters') }}</a-button>
      </div>

      <div v-else class="list-body">
        <div class="list-head">
          <div class="col">{{ $t('adminOrdersPage.colId') }}</div>
          <div class="col">{{ $t('adminOrdersPage.colOrderNo') }}</div>
          <div class="col">{{ $t('adminOrdersPage.colType') }}</div>
          <div class="col">{{ $t('adminOrdersPage.colUser') }}</div>
          <div class="col">{{ $t('adminOrdersPage.colPlanOrAmount') }}</div>
          <div class="col">{{ $t('adminOrdersPage.colAmount') }}</div>
          <div class="col">{{ $t('adminOrdersPage.colPayMethod') }}</div>
          <div class="col">{{ $t('adminOrdersPage.colStatus') }}</div>
          <div class="col">{{ $t('adminOrdersPage.colSource') }}</div>
          <div class="col">{{ $t('adminOrdersPage.colCreateTime') }}</div>
          <div class="col col-action">{{ $t('adminOrdersPage.colAction') }}</div>
        </div>

        <a-spin :loading="loading" style="width: 100%">
          <div v-for="o in orders" :key="o.id" class="list-row">
            <div class="col"><span class="cell-mono">#{{ o.id }}</span></div>

            <div class="col">
              <code class="cell-mono ellipsis" :title="o.order_no">{{ o.order_no }}</code>
            </div>

            <div class="col">
              <a-tag :color="typeColorMap[o.type] || 'gray'" size="small">
                {{ typeNameMap[o.type] || `Type ${o.type}` }}
              </a-tag>
            </div>

            <div class="col">
              <span class="cell-strong" :title="userDisplay(o)">
                {{ userDisplay(o) }}
              </span>
              <span class="cell-muted cell-mono">#{{ o.user_id }}</span>
            </div>

            <div class="col">
              <span class="cell-muted ellipsis" :title="planOrAmount(o)">{{ planOrAmount(o) }}</span>
            </div>

            <div class="col">
              <span class="cell-mono">¥{{ formatAmount(o.amount) }}</span>
            </div>

            <div class="col">
              <span class="cell-muted">{{ payNameMap[o.pay_method] || o.pay_method || '-' }}</span>
            </div>

            <div class="col">
              <span class="status-chip" :class="statusClass(o.status)">
                <span class="status-dot"></span>
                {{ statusNameMap[o.status] || '-' }}
              </span>
            </div>

            <div class="col">
              <span class="cell-muted">{{ sourceNameMap[o.source] || '-' }}</span>
            </div>

            <div class="col">
              <span class="cell-mono">{{ formatTime(o.create_time) }}</span>
            </div>

            <div class="col col-action">
              <a-button type="text" size="small" @click="openDetail(o)">{{ $t('adminOrdersPage.view') }}</a-button>
              <a-button
                v-if="o.status === 0"
                type="text"
                size="small"
                @click="openMarkPaid(o)"
              >{{ $t('adminOrdersPage.markPaid') }}</a-button>
              <a-popconfirm
                v-if="o.status === 1"
                :content="$t('adminOrdersPage.confirmRefund')"
                @ok="handleMarkRefunded(o)"
              >
                <a-button type="text" size="small">{{ $t('adminOrdersPage.refund') }}</a-button>
              </a-popconfirm>
              <a-popconfirm
                v-if="authStore.isRoot && o.status !== 1"
                :content="$t('adminOrdersPage.confirmDelete', { no: o.order_no })"
                @ok="handleDelete(o)"
              >
                <a-button type="text" size="small" class="danger-btn">{{ $t('adminOrdersPage.delete') }}</a-button>
              </a-popconfirm>
            </div>
          </div>
        </a-spin>
      </div>

      <div v-if="orders.length > 0" class="list-footer">
        <a-pagination
          :current="activePage + 1"
          :total="totalCountForPager"
          :page-size="pageSize"
          show-total
          show-page-size
          :page-size-options="[10, 20, 50]"
          size="small"
          @change="onPaginationChange"
          @page-size-change="onPageSizeChange"
        />
      </div>
    </div>

    <!-- 详情弹窗 -->
    <a-modal
      :visible="detailVisible"
      @update:visible="(v) => (detailVisible = v)"
      @cancel="detailVisible = false"
      :footer="false"
      :title="$t('adminOrdersPage.detailTitle')"
      width="560"
    >
      <a-spin :loading="detailLoading" style="width: 100%">
        <div v-if="currentOrder" class="detail-grid">
          <div class="detail-row">
            <span class="detail-label">{{ $t('adminOrdersPage.detailOrderNo') }}</span>
            <span class="detail-value cell-mono">{{ currentOrder.order_no }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">{{ $t('adminOrdersPage.detailType') }}</span>
            <span class="detail-value">{{ typeNameMap[currentOrder.type] || '-' }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">{{ $t('adminOrdersPage.detailSource') }}</span>
            <span class="detail-value">{{ sourceNameMap[currentOrder.source] || '-' }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">{{ $t('adminOrdersPage.detailUser') }}</span>
            <span class="detail-value">
              {{ userDisplay(currentOrder) }} <span class="cell-muted">#{{ currentOrder.user_id }}</span>
            </span>
          </div>
          <div class="detail-row">
            <span class="detail-label">{{ $t('adminOrdersPage.detailPlanOrAmount') }}</span>
            <span class="detail-value">{{ planOrAmount(currentOrder) }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">{{ $t('adminOrdersPage.detailAmount') }}</span>
            <span class="detail-value cell-mono">¥{{ formatAmount(currentOrder.amount) }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">{{ $t('adminOrdersPage.detailPayMethod') }}</span>
            <span class="detail-value">{{ payNameMap[currentOrder.pay_method] || currentOrder.pay_method || '-' }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">{{ $t('adminOrdersPage.detailPayTradeNo') }}</span>
            <span class="detail-value cell-mono">{{ currentOrder.pay_trade_no || '-' }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">{{ $t('adminOrdersPage.detailStatus') }}</span>
            <span class="detail-value">{{ statusNameMap[currentOrder.status] || '-' }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">{{ $t('adminOrdersPage.detailCreateTime') }}</span>
            <span class="detail-value">{{ formatTime(currentOrder.create_time) }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">{{ $t('adminOrdersPage.detailPayTime') }}</span>
            <span class="detail-value">{{ formatTime(currentOrder.pay_time) }}</span>
          </div>
        </div>
      </a-spin>
    </a-modal>

    <!-- 标记已付弹窗 -->
    <a-modal
      :visible="markPaidVisible"
      @update:visible="(v) => (markPaidVisible = v)"
      @cancel="closeMarkPaid"
      @ok="submitMarkPaid"
      :ok-loading="markPaidSubmitting"
      :title="$t('adminOrdersPage.markPaidTitle')"
      width="480"
      :ok-text="$t('adminOrdersPage.markPaidOk')"
      :cancel-text="$t('adminOrdersPage.cancel')"
    >
      <a-form :model="markPaidForm" layout="vertical">
        <a-form-item :label="$t('adminOrdersPage.payMethodLabel')" required>
          <a-select v-model="markPaidForm.pay_method" :placeholder="$t('adminOrdersPage.payMethodPlaceholder')">
            <a-option value="wechat">{{ $t('adminOrdersPage.payWechat') }}</a-option>
            <a-option value="alipay">{{ $t('adminOrdersPage.payAlipay') }}</a-option>
            <a-option value="bank">{{ $t('adminOrdersPage.payBank') }}</a-option>
            <a-option value="offline">{{ $t('adminOrdersPage.payOffline') }}</a-option>
            <a-option value="free">{{ $t('adminOrdersPage.payFree') }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item :label="$t('adminOrdersPage.payTradeNoLabel')">
          <a-input v-model="markPaidForm.pay_trade_no" :placeholder="$t('adminOrdersPage.payTradeNoPlaceholder')" allow-clear />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
// AdminOrders 管理员订单列表页
// Admin-side order list page (channel-style list + filters + modals)
// 版本: v0.0.13
// 日期: 2026-09-08
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Message } from '@arco-design/web-vue'
import { IconRefresh, IconStorage } from '@arco-design/web-vue/es/icon'
import orderApi from '@/api/order'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()

const authStore = useAuthStore()

// 常量映射
const typeNameMap = computed(() => ({ 1: t('adminOrdersPage.typePlan'), 2: t('adminOrdersPage.typeTopup') }))
const typeColorMap = { 1: 'arcoblue', 2: 'orange' }
const statusNameMap = computed(() => ({
  0: t('adminOrdersPage.statusPending'),
  1: t('adminOrdersPage.statusPaid'),
  2: t('adminOrdersPage.statusCanceled'),
  3: t('adminOrdersPage.statusRefunded'),
}))
const sourceNameMap = computed(() => ({
  1: t('adminOrdersPage.sourceSelf'),
  2: t('adminOrdersPage.sourceAdmin'),
}))
const payNameMap = computed(() => ({
  '': '-',
  wechat: t('adminOrdersPage.payWechat'),
  alipay: t('adminOrdersPage.payAlipay'),
  bank: t('adminOrdersPage.payBankShort'),
  offline: t('adminOrdersPage.payOffline'),
  free: t('adminOrdersPage.payFree'),
}))

function statusClass(status) {
  if (status === 1) return 'status-on'
  if (status === 3) return 'status-warn'
  return 'status-off'
}

function formatTime(ts) {
  if (!ts) return '-'
  const d = new Date(ts * 1000)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function formatAmount(v) {
  const n = Number(v || 0)
  return n.toFixed(2)
}

function userDisplay(o) {
  const u = o.user
  if (!u) return ''
  return u.display_name || u.username || ''
}

function planOrAmount(o) {
  // 套餐订单：尝试解析 plan_info JSON 显示套餐名
  if (o.type === 1 && o.plan_info) {
    try {
      const p = JSON.parse(o.plan_info)
      if (p && p.name) return p.name
    } catch (_) {
      // fallthrough
    }
  }
  // 充值订单或解析失败：显示 plan_id 或 user_id
  if (o.plan_id) return `#${o.plan_id}`
  return '-'
}

// 状态
const loading = ref(false)
const orders = ref([])
const userMap = reactive({})
const activePage = ref(0)
const pageSize = ref(20)

const searchKeyword = ref('')
const filterType = ref('')
const filterStatus = ref('')
const filterSource = ref('')

const hasActiveFilter = computed(
  () => !!searchKeyword.value || !!filterType.value || filterStatus.value !== '' || !!filterSource.value,
)

const totalCountForPager = computed(() => {
  // 项目约定：后端裸数组响应，前端用 items.length + pageSize 兜底"至少还有一页"
  // （与 Channel.vue 一致）。
  if (orders.value.length === 0) return 0
  return orders.value.length + pageSize.value
})

// 详情弹窗
const detailVisible = ref(false)
const detailLoading = ref(false)
const currentOrder = ref(null)

// 标记已付弹窗
const markPaidVisible = ref(false)
const markPaidSubmitting = ref(false)
const markPaidForm = reactive({ pay_method: 'offline', pay_trade_no: '' })
const markPaidTarget = ref(null)

async function fetchOrders() {
  loading.value = true
  try {
    const params = {
      p: activePage.value,
      page_size: pageSize.value,
    }
    if (filterType.value) params.type = filterType.value
    if (filterStatus.value !== '') params.status = filterStatus.value
    if (filterSource.value) params.source = filterSource.value
    if (searchKeyword.value) params.keyword = searchKeyword.value

    // 与 Channel.vue 一致：有 keyword 走 search 端点；无 keyword 走 list 端点
    const { data } = searchKeyword.value
      ? await orderApi.search(searchKeyword.value, filterType.value, filterStatus.value, filterSource.value)
      : await orderApi.list(params)
    if (data.success) {
      orders.value = data.data || []
      // 同步 user_map
      const um = {}
      for (const o of orders.value) {
        if (o.user) um[o.user.id] = o.user
      }
      Object.keys(userMap).forEach((k) => delete userMap[k])
      Object.assign(userMap, um)
    } else {
      Message.error(data.message || t('adminOrdersPage.loadFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('adminOrdersPage.loadFailed'))
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  activePage.value = 0
  fetchOrders()
}

function handleSearchClear() {
  searchKeyword.value = ''
  activePage.value = 0
  fetchOrders()
}

function handleFilterChange() {
  activePage.value = 0
  fetchOrders()
}

function handleResetFilters() {
  searchKeyword.value = ''
  filterType.value = ''
  filterStatus.value = ''
  filterSource.value = ''
  activePage.value = 0
  fetchOrders()
}

function handleRefresh() {
  fetchOrders()
}

function onPaginationChange(page) {
  // arco pagination 用 1-based，内部用 0-based
  activePage.value = page - 1
  fetchOrders()
}

function onPageSizeChange(s) {
  pageSize.value = s
  activePage.value = 0
  fetchOrders()
}

async function openDetail(o) {
  currentOrder.value = o
  detailVisible.value = true
  detailLoading.value = true
  try {
    const { data } = await orderApi.get(o.id)
    if (data.success && data.data) {
      currentOrder.value = data.data
      if (data.data.user) userMap[data.data.user.id] = data.data.user
    }
  } catch (e) {
    // 详情加载失败不影响主列表
  } finally {
    detailLoading.value = false
  }
}

function openMarkPaid(o) {
  markPaidTarget.value = o
  markPaidForm.pay_method = 'offline'
  markPaidForm.pay_trade_no = ''
  markPaidVisible.value = true
}

function closeMarkPaid() {
  markPaidVisible.value = false
  markPaidTarget.value = null
}

async function submitMarkPaid() {
  if (!markPaidTarget.value) return
  if (!markPaidForm.pay_method) {
    Message.warning(t('adminOrdersPage.selectPayMethod'))
    return
  }
  markPaidSubmitting.value = true
  try {
    const { data } = await orderApi.markPaid(markPaidTarget.value.id, {
      status: 1,
      pay_method: markPaidForm.pay_method,
      pay_trade_no: markPaidForm.pay_trade_no,
    })
    if (data.success) {
      Message.success(t('adminOrdersPage.orderPaid'))
      closeMarkPaid()
      fetchOrders()
    } else {
      Message.error(data.message || t('adminOrdersPage.opFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('adminOrdersPage.opFailed'))
  } finally {
    markPaidSubmitting.value = false
  }
}

async function handleMarkRefunded(o) {
  try {
    const { data } = await orderApi.markPaid(o.id, { status: 3 })
    if (data.success) {
      Message.success(t('adminOrdersPage.orderRefunded'))
      fetchOrders()
    } else {
      Message.error(data.message || t('adminOrdersPage.opFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('adminOrdersPage.opFailed'))
  }
}

async function handleDelete(o) {
  try {
    const { data } = await orderApi.delete(o.id)
    if (data.success) {
      Message.success(t('adminOrdersPage.orderDeleted'))
      fetchOrders()
    } else {
      Message.error(data.message || t('adminOrdersPage.deleteFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('adminOrdersPage.deleteFailed'))
  }
}

onMounted(() => {
  fetchOrders()
})
</script>

<style scoped>
.admin-orders-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* ============ 顶部欢迎条 ============ */
.welcome-bar {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  padding: 4px 4px 0;
}
.welcome-title {
  font-size: 24px;
  font-weight: 600;
  color: var(--color-text-1);
  margin: 0 0 4px;
  letter-spacing: -0.2px;
}
.welcome-desc {
  font-size: 13px;
  color: var(--color-text-3);
  margin: 0;
}
.welcome-meta { display: flex; gap: 6px; }
.meta-chip {
  font-size: 12px;
  color: var(--color-text-3);
  background: var(--color-fill-2);
  padding: 3px 10px;
  border-radius: 4px;
}

/* ============ 搜索栏 ============ */
.search-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 14px 20px;
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-2);
  border-radius: 8px;
}
.search-left,
.search-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

/* ============ 列表 ============ */
.list-wrap {
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-2);
  border-radius: 8px;
  overflow: hidden;
}
.list-body {
  padding: 0;
  overflow-x: auto;
}
.list-head,
.list-row {
  display: grid;
  grid-template-columns: 70px 160px 80px 160px 130px 90px 90px 90px 90px 140px 240px;
  align-items: center;
  padding: 0 20px;
  min-width: max-content;
}
.list-head {
  height: 40px;
  background: var(--color-fill-1);
  border-bottom: 1px solid var(--color-fill-3);
  font-size: 12px;
  font-weight: 500;
  color: var(--color-text-3);
}
.list-row {
  min-height: 52px;
  border-bottom: 1px solid var(--color-fill-3);
  transition: background 0.15s;
}
.list-row:last-child { border-bottom: none; }
.list-row:hover { background: var(--color-fill-1); }

/* ============ 单元格 ============ */
.col {
  font-size: 13px;
  color: var(--color-text-2);
  min-width: 0;
  padding-right: 12px;
}
.col:last-child { padding-right: 0; }
.col-action {
  display: flex;
  justify-content: flex-end;
  gap: 0;
}
.col-action :deep(.arco-btn) { padding: 0 6px; }

.cell-mono {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 12px;
  color: var(--color-text-2);
  font-variant-numeric: tabular-nums;
}
.cell-strong { color: var(--color-text-1); font-weight: 500; }
.cell-muted { color: var(--color-text-3); }
.ellipsis {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}

/* ============ 状态 chip ============ */
.status-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 2px 10px;
  border-radius: 10px;
  font-size: 12px;
  font-weight: 500;
  width: max-content;
}
.status-dot { width: 6px; height: 6px; border-radius: 50%; }
.status-on { background: rgba(0, 180, 42, 0.08); color: #00b42a; }
.status-on .status-dot { background: #00b42a; }
.status-off { background: var(--color-fill-2); color: var(--color-text-3); }
.status-off .status-dot { background: var(--color-text-4); }
.status-warn { background: rgba(245, 63, 63, 0.08); color: #f53f3f; }
.status-warn .status-dot { background: #f53f3f; }

/* ============ 操作列 ============ */
.danger-btn { color: var(--color-text-2); }
.danger-btn:hover {
  color: #f53f3f !important;
  background: rgba(245, 63, 63, 0.06) !important;
}

/* ============ 分页 ============ */
.list-footer {
  display: flex;
  justify-content: flex-end;
  padding: 14px 20px;
  border-top: 1px solid var(--color-fill-3);
}

/* ============ 空状态 ============ */
.empty-state {
  padding: 80px 20px;
  text-align: center;
}
.empty-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  border-radius: 14px;
  background: var(--color-fill-2);
  color: var(--color-text-3);
  margin-bottom: 12px;
}
.empty-title { font-size: 14px; font-weight: 500; color: var(--color-text-1); margin: 0 0 4px; }
.empty-desc { font-size: 13px; color: var(--color-text-3); margin: 0 0 16px; }

/* ============ 详情弹窗 ============ */
.detail-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 0;
}
.detail-row {
  display: grid;
  grid-template-columns: 110px 1fr;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px dashed var(--color-fill-3);
  font-size: 13px;
}
.detail-row:last-child { border-bottom: none; }
.detail-label { color: var(--color-text-3); }
.detail-value { color: var(--color-text-2); }
</style>