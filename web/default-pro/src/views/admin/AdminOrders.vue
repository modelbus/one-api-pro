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
        <h1 class="welcome-title">订单</h1>
        <p class="welcome-desc">管理全部用户订单 · 含套餐订阅与充值</p>
      </div>
      <div class="welcome-meta">
        <span class="meta-chip">共 {{ totalCountForPager }} 条</span>
      </div>
    </div>

    <!-- 独立搜索栏 -->
    <div class="search-card">
      <div class="search-left">
        <a-input-search
          v-model="searchKeyword"
          placeholder="订单号 / 支付流水号"
          allow-clear
          style="width: 240px"
          @search="handleSearch"
          @clear="handleSearchClear"
        />
        <a-select v-model="filterType" placeholder="类型" allow-clear style="width: 110px" @change="handleFilterChange">
          <a-option :value="1">套餐</a-option>
          <a-option :value="2">充值</a-option>
        </a-select>
        <a-select v-model="filterStatus" placeholder="状态" allow-clear style="width: 120px" @change="handleFilterChange">
          <a-option :value="0">待支付</a-option>
          <a-option :value="1">已支付</a-option>
          <a-option :value="2">已取消</a-option>
          <a-option :value="3">已退款</a-option>
        </a-select>
        <a-select v-model="filterSource" placeholder="来源" allow-clear style="width: 120px" @change="handleFilterChange">
          <a-option :value="1">用户自助</a-option>
          <a-option :value="2">管理员</a-option>
        </a-select>
      </div>
      <div class="search-right">
        <a-button @click="handleResetFilters" :disabled="!hasActiveFilter">重置筛选</a-button>
        <a-button @click="handleRefresh">
          <template #icon><icon-refresh :size="14" /></template>
          刷新
        </a-button>
      </div>
    </div>

    <!-- 列表 -->
    <div class="list-wrap">
      <div v-if="orders.length === 0 && !loading" class="empty-state">
        <div class="empty-icon">
          <icon-storage :size="32" />
        </div>
        <p class="empty-title">暂无订单</p>
        <p class="empty-desc">没有匹配的订单记录</p>
        <a-button v-if="hasActiveFilter" type="primary" @click="handleResetFilters">重置筛选</a-button>
      </div>

      <div v-else class="list-body">
        <div class="list-head">
          <div class="col">ID</div>
          <div class="col">订单号</div>
          <div class="col">类型</div>
          <div class="col">用户</div>
          <div class="col">套餐/额度</div>
          <div class="col">金额</div>
          <div class="col">支付方式</div>
          <div class="col">状态</div>
          <div class="col">来源</div>
          <div class="col">创建时间</div>
          <div class="col col-action">操作</div>
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
              <a-button type="text" size="small" @click="openDetail(o)">查看</a-button>
              <a-button
                v-if="o.status === 0"
                type="text"
                size="small"
                @click="openMarkPaid(o)"
              >标记已付</a-button>
              <a-popconfirm
                v-if="o.status === 1"
                content="确认将该订单标记为退款？"
                @ok="handleMarkRefunded(o)"
              >
                <a-button type="text" size="small">退款</a-button>
              </a-popconfirm>
              <a-popconfirm
                v-if="authStore.isRoot && o.status !== 1"
                :content="`确认删除订单 ${o.order_no}？`"
                @ok="handleDelete(o)"
              >
                <a-button type="text" size="small" class="danger-btn">删除</a-button>
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
      title="订单详情"
      width="560"
    >
      <a-spin :loading="detailLoading" style="width: 100%">
        <div v-if="currentOrder" class="detail-grid">
          <div class="detail-row">
            <span class="detail-label">订单号</span>
            <span class="detail-value cell-mono">{{ currentOrder.order_no }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">类型</span>
            <span class="detail-value">{{ typeNameMap[currentOrder.type] || '-' }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">来源</span>
            <span class="detail-value">{{ sourceNameMap[currentOrder.source] || '-' }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">用户</span>
            <span class="detail-value">
              {{ userDisplay(currentOrder) }} <span class="cell-muted">#{{ currentOrder.user_id }}</span>
            </span>
          </div>
          <div class="detail-row">
            <span class="detail-label">套餐/额度</span>
            <span class="detail-value">{{ planOrAmount(currentOrder) }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">金额</span>
            <span class="detail-value cell-mono">¥{{ formatAmount(currentOrder.amount) }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">支付方式</span>
            <span class="detail-value">{{ payNameMap[currentOrder.pay_method] || currentOrder.pay_method || '-' }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">支付流水号</span>
            <span class="detail-value cell-mono">{{ currentOrder.pay_trade_no || '-' }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">状态</span>
            <span class="detail-value">{{ statusNameMap[currentOrder.status] || '-' }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">创建时间</span>
            <span class="detail-value">{{ formatTime(currentOrder.create_time) }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">支付时间</span>
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
      title="标记订单为已支付"
      width="480"
      ok-text="确认标记"
      cancel-text="取消"
    >
      <a-form :model="markPaidForm" layout="vertical">
        <a-form-item label="支付方式" required>
          <a-select v-model="markPaidForm.pay_method" placeholder="选择支付方式">
            <a-option value="wechat">微信</a-option>
            <a-option value="alipay">支付宝</a-option>
            <a-option value="bank">银行转账</a-option>
            <a-option value="offline">线下</a-option>
            <a-option value="free">免费</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="支付流水号（选填）">
          <a-input v-model="markPaidForm.pay_trade_no" placeholder="留空则使用订单号" allow-clear />
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
import { Message } from '@arco-design/web-vue'
import { IconRefresh, IconStorage } from '@arco-design/web-vue/es/icon'
import orderApi from '@/api/order'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// 常量映射
const typeNameMap = { 1: '套餐', 2: '充值' }
const typeColorMap = { 1: 'arcoblue', 2: 'orange' }
const statusNameMap = { 0: '待支付', 1: '已支付', 2: '已取消', 3: '已退款' }
const sourceNameMap = { 1: '用户自助', 2: '管理员' }
const payNameMap = {
  '': '-',
  wechat: '微信',
  alipay: '支付宝',
  bank: '银行',
  offline: '线下',
  free: '免费',
}

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
      Message.error(data.message || '加载失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '加载失败')
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
    Message.warning('请选择支付方式')
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
      Message.success('订单已支付')
      closeMarkPaid()
      fetchOrders()
    } else {
      Message.error(data.message || '操作失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '操作失败')
  } finally {
    markPaidSubmitting.value = false
  }
}

async function handleMarkRefunded(o) {
  try {
    const { data } = await orderApi.markPaid(o.id, { status: 3 })
    if (data.success) {
      Message.success('订单已标记为退款')
      fetchOrders()
    } else {
      Message.error(data.message || '操作失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '操作失败')
  }
}

async function handleDelete(o) {
  try {
    const { data } = await orderApi.delete(o.id)
    if (data.success) {
      Message.success('订单已删除')
      fetchOrders()
    } else {
      Message.error(data.message || '删除失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '删除失败')
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