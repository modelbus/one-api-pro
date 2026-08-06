<template>
  <div class="sub-page">
    <!-- 顶部欢迎条 -->
    <div class="welcome-bar">
      <div class="welcome-text">
        <h1 class="welcome-title">订阅</h1>
        <p class="welcome-desc">管理用户的订阅套餐、用量配额与计费周期</p>
      </div>
      <div class="welcome-meta">
        <span class="meta-chip">共 {{ subscriptions.length }} 条</span>
      </div>
    </div>

    <!-- 独立搜索栏 -->
    <div class="search-card">
      <div class="search-left">
        <a-input-search
          v-model="searchKeyword"
          placeholder="搜索用户、套餐或备注..."
          allow-clear
          @search="searchSubscriptions"
          @clear="handleClearSearch"
          :style="{ width: '320px' }"
        />
      </div>
      <div class="search-right">
        <a-button @click="refresh" :loading="loading">
          <template #icon><icon-refresh :size="14" /></template>
          刷新
        </a-button>
        <a-button v-if="authStore.isAdmin" type="primary" size="large" @click="openAddModal">
          <template #icon><icon-plus :size="14" /></template>
          新建订阅
        </a-button>
      </div>
    </div>

    <!-- 列表 -->
    <div class="list-wrap">
      <div v-if="pageItems.length === 0 && !loading" class="empty-state">
        <div class="empty-icon">
          <icon-calendar :size="32" />
        </div>
        <p class="empty-title">还没有任何订阅</p>
        <p v-if="authStore.isAdmin" class="empty-desc">为用户创建第一个订阅即可开始</p>
        <p v-else class="empty-desc">尚无订阅记录</p>
        <a-button v-if="authStore.isAdmin" type="primary" @click="openAddModal">
          <template #icon><icon-plus :size="14" /></template>
          立即创建
        </a-button>
      </div>

      <div v-else class="list-body">
        <div class="list-head" :style="{ gridTemplateColumns: headGrid }">
          <div class="col">ID</div>
          <div v-if="authStore.isAdmin" class="col">用户</div>
          <div class="col">套餐</div>
          <div class="col">计费</div>
          <div class="col">开始时间</div>
          <div class="col">结束时间</div>
          <div class="col">状态</div>
          <div class="col col-action">操作</div>
        </div>

        <a-spin :loading="loading" style="width: 100%">
          <div
            v-for="s in pageItems"
            :key="s.id"
            class="list-row"
            :style="{ gridTemplateColumns: headGrid }"
          >
            <div class="col"><span class="cell-mono">#{{ s.id }}</span></div>

            <div v-if="authStore.isAdmin" class="col">
              <span class="cell-strong">{{ getUserName(s.user_id) }}</span>
            </div>

            <div class="col">
              <span class="cell-strong">{{ s.plan?.name || s.plan_name || s.plan_id }}</span>
            </div>

            <div class="col">
              <span class="billing-chip">{{ renderBillingType(s.billing_type) }}</span>
            </div>

            <div class="col">
              <span class="cell-mono">{{ formatTime(s.start_time) }}</span>
            </div>

            <div class="col">
              <span class="cell-mono">{{ formatTime(s.end_time) }}</span>
            </div>

            <div class="col">
              <span class="status-chip" :class="statusClass(s.status)">
                <span class="status-dot"></span>
                {{ renderStatus(s.status) }}
              </span>
            </div>

            <div class="col col-action">
              <a-button type="text" size="small" @click="viewUsage(s)">用量</a-button>
              <template v-if="authStore.isAdmin">
                <a-button
                  v-if="isActive(s.status)"
                  type="text"
                  size="small"
                  @click="openRenewModal(s)"
                >
                  续费
                </a-button>
                <a-popconfirm
                  v-if="isActive(s.status)"
                  content="确定将该订阅置为过期？"
                  @ok="expireSubscription(s.id)"
                >
                  <a-button type="text" size="small">过期</a-button>
                </a-popconfirm>
                <a-popconfirm
                  content="确定要删除该订阅？删除后无法恢复"
                  @ok="deleteSubscription(s.id)"
                >
                  <a-button type="text" size="small" class="danger-btn">删除</a-button>
                </a-popconfirm>
              </template>
            </div>
          </div>
        </a-spin>
      </div>

      <div v-if="subscriptions.length > 0" class="list-footer">
        <a-pagination
          :current="activePage"
          :total="subscriptions.length"
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

    <!-- 新建订阅 -->
    <a-modal
      v-model:visible="showAddModal"
      title="新建订阅"
      :width="520"
      :ok-loading="addSubmitting"
      @ok="addSubscription"
      @cancel="showAddModal = false"
      ok-text="确定创建"
      cancel-text="取消"
    >
      <a-form :model="addForm" layout="vertical" class="sub-form">
        <a-form-item label="用户">
          <a-select
            v-model="addForm.user_id"
            placeholder="选择用户"
            allow-search
            :loading="userListLoading"
          >
            <a-option
              v-for="u in users"
              :key="u.id"
              :value="u.id"
              :label="u.username"
            >
              #{{ u.id }} {{ u.username }}
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="套餐">
          <a-select
            v-model="addForm.plan_id"
            placeholder="选择套餐"
            allow-search
            :loading="planListLoading"
          >
            <a-option
              v-for="p in plans"
              :key="p.id"
              :value="p.id"
              :label="p.name"
            >
              {{ p.name }}<span v-if="p.recommended" class="opt-recommend"> *</span>
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="计费类型">
          <a-select v-model="addForm.billing_type" placeholder="选择计费类型">
            <a-option value="token">按 Token</a-option>
            <a-option value="request">按次数</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="时长（天）">
          <a-input-number
            v-model="addForm.duration_days"
            :min="1"
            :max="3650"
            placeholder="例如 30"
            style="width: 100%"
          />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea
            v-model="addForm.notes"
            placeholder="可选备注..."
            :auto-size="{ minRows: 2, maxRows: 4 }"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 续费 -->
    <a-modal
      v-model:visible="showRenewModal"
      title="续费订阅"
      :width="460"
      :ok-loading="renewSubmitting"
      @ok="renewSubscription"
      @cancel="showRenewModal = false"
      ok-text="确定续费"
      cancel-text="取消"
    >
      <a-form :model="renewForm" layout="vertical" class="sub-form">
        <a-form-item label="新的结束时间">
          <a-date-picker
            v-model="renewForm.end_time"
            show-time
            format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%"
          />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea
            v-model="renewForm.notes"
            placeholder="续费备注..."
            :auto-size="{ minRows: 2, maxRows: 4 }"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 用量弹窗 -->
    <a-modal
      v-model:visible="showUsageModal"
      title="订阅用量"
      @cancel="showUsageModal = false"
      :footer="false"
      :width="780"
      unmount-on-close
    >
      <a-spin :loading="usageLoading">
        <template v-if="usageData">
          <!-- 基础信息 -->
          <div class="usage-info">
            <div class="ud-item">
              <span class="ud-label">套餐</span>
              <span class="ud-value">{{ usageData.subscription?.plan?.name || usageData.plan_name || '-' }}</span>
            </div>
            <div class="ud-item">
              <span class="ud-label">计费类型</span>
              <span class="ud-value">{{ renderBillingType(usageData.billing_type || usageData.subscription?.billing_type) }}</span>
            </div>
            <div class="ud-item">
              <span class="ud-label">开始时间</span>
              <span class="ud-value ud-mono">{{ formatTime(usageData.subscription?.start_time) }}</span>
            </div>
            <div class="ud-item">
              <span class="ud-label">结束时间</span>
              <span class="ud-value ud-mono">{{ formatTime(usageData.subscription?.end_time) }}</span>
            </div>
          </div>

          <!-- 窗口用量进度 -->
          <div class="usage-section" v-if="usageData.weighted">
            <h4 class="section-title">用量窗口</h4>

            <div
              v-for="wt in windowTypes"
              :key="wt"
              class="window-block"
            >
              <div class="window-head" @click="toggleWindow(wt)">
                <span class="window-toggle">
                  <icon-caret-right v-if="!expandedWindows[wt]" :size="12" />
                  <icon-caret-down v-else :size="12" />
                </span>
                <span class="window-name">{{ renderWindowType(wt) }}</span>
                <span class="window-pct" :class="windowPctClass(usageData.weighted[wt])">
                  {{ formatPercent(usageData.weighted[wt]) }}
                </span>
                <span
                  v-if="usageData.next_reset && usageData.next_reset[wt]"
                  class="window-reset"
                >
                  下次重置：{{ formatTime(usageData.next_reset[wt]) }}
                </span>
              </div>

              <!-- 分段进度条（按模型分段） -->
              <div
                v-if="getSegments(wt).length > 0"
                class="window-segments"
              >
                <div class="segment-bar">
                  <div
                    v-for="(seg, i) in getSegments(wt)"
                    :key="i"
                    class="segment"
                    :style="{
                      left: seg.offset + '%',
                      width: seg.width + '%',
                      background: seg.color,
                    }"
                    :title="`${seg.label}: ${formatNumber(seg.value)}`"
                  ></div>
                </div>
                <div class="segment-legend">
                  <span
                    v-for="(seg, i) in getSegments(wt)"
                    :key="i"
                    class="legend-tag"
                  >
                    <span class="legend-dot" :style="{ background: seg.color }"></span>
                    <span class="legend-name">{{ seg.label }}</span>
                    <span class="legend-value">{{ formatNumber(seg.value) }}</span>
                  </span>
                </div>
              </div>
              <div v-else class="window-empty">暂无模型用量</div>

              <!-- 展开模型明细 -->
              <div v-if="expandedWindows[wt]" class="window-detail">
                <table class="detail-table" v-if="getModelDetails(wt).length > 0">
                  <thead>
                    <tr>
                      <th>模型</th>
                      <th>请求数</th>
                      <th>占比</th>
                      <th>Prompt Tokens</th>
                      <th>Completion Tokens</th>
                      <th>Cached Tokens</th>
                      <th>Token 占比</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(md, idx) in getModelDetails(wt)" :key="idx">
                      <td>
                        <span class="detail-model">
                          <span class="legend-dot" :style="{ background: MODEL_COLORS[idx % MODEL_COLORS.length] }"></span>
                          {{ md.model }}
                        </span>
                      </td>
                      <td><span class="cell-num">{{ formatNumber(md.requests) }}</span></td>
                      <td><span class="cell-muted">{{ md.request_percent > 0 ? formatPercent(md.request_percent) : '-' }}</span></td>
                      <td><span class="cell-num">{{ formatNumber(md.prompt_tokens) }}</span></td>
                      <td><span class="cell-num">{{ formatNumber(md.completion_tokens) }}</span></td>
                      <td><span class="cell-num">{{ formatNumber(md.cached_tokens) }}</span></td>
                      <td><span class="cell-muted">{{ md.token_percent > 0 ? formatPercent(md.token_percent) : '-' }}</span></td>
                    </tr>
                  </tbody>
                </table>
                <div v-else class="window-empty">暂无模型明细</div>
              </div>
            </div>
          </div>

          <div v-else class="usage-empty">暂无用量数据</div>
        </template>
      </a-spin>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  IconPlus, IconCalendar, IconRefresh,
  IconCaretRight, IconCaretDown,
} from '@arco-design/web-vue/es/icon'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'

const authStore = useAuthStore()

const ITEMS_PER_PAGE = 10
const pageSize = ref(ITEMS_PER_PAGE)
const MODEL_COLORS = ['#165dff', '#00b42a', '#ff7d00', '#f53f3f', '#722ed1', '#0fc6c2', '#af52de', '#ff5722', '#5856d6', '#a2845e']
const WINDOW_TYPES = ['period', 'week', 'month']

// ============ 数据状态 ============
const subscriptions = ref([])
const users = ref([])
const plans = ref([])
const loading = ref(false)
const searching = ref(false)
const searchKeyword = ref('')
const activePage = ref(1)

// ============ 弹窗状态 ============
const showAddModal = ref(false)
const addSubmitting = ref(false)
const addForm = reactive({
  user_id: null,
  plan_id: null,
  billing_type: 'token',
  duration_days: 30,
  notes: '',
})

const showRenewModal = ref(false)
const renewSubmitting = ref(false)
const renewForm = reactive({
  id: null,
  end_time: '',
  notes: '',
})

const showUsageModal = ref(false)
const usageLoading = ref(false)
const usageData = ref(null)
const expandedWindows = ref({})

// ============ 计算属性 ============
const pageItems = computed(() => {
  const start = (activePage.value - 1) * ITEMS_PER_PAGE
  return subscriptions.value.slice(start, start + ITEMS_PER_PAGE)
})

const headGrid = computed(() => {
  if (authStore.isAdmin) {
    return '80px 130px 160px 100px 170px 170px 110px 240px'
  }
  return '80px 160px 100px 170px 170px 110px 130px'
})

const windowTypes = WINDOW_TYPES

// ============ API 调用 ============
async function loadSubscriptions(reset = true) {
  loading.value = true
  try {
    let url
    if (authStore.isAdmin) {
      url = `/api/subscription/?p=0`
    } else {
      url = '/api/subscription/self'
    }
    const { data } = await api.get(url)
    if (data.success) {
      const items = Array.isArray(data.data) ? data.data : []
      subscriptions.value = items
      activePage.value = 1
    } else {
      Message.error(data.message || '加载订阅失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '加载订阅失败')
  } finally {
    loading.value = false
  }
}

async function loadPlans() {
  try {
    const { data } = await api.get('/api/plan/?p=0')
    if (data.success) {
      plans.value = Array.isArray(data.data) ? data.data : []
    }
  } catch {
    plans.value = []
  }
}

async function loadUsers() {
  try {
    const { data } = await api.get('/api/user/?p=0')
    if (data.success) {
      users.value = Array.isArray(data.data) ? data.data : []
    }
  } catch {
    users.value = []
  }
}

async function searchSubscriptions() {
  if (!authStore.isAdmin) return
  if (!searchKeyword.value) {
    await loadSubscriptions()
    return
  }
  searching.value = true
  try {
    const { data } = await api.get(`/api/subscription/search?keyword=${encodeURIComponent(searchKeyword.value)}`)
    if (data.success) {
      subscriptions.value = Array.isArray(data.data) ? data.data : []
      activePage.value = 1
    } else {
      Message.error(data.message || '搜索失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '搜索失败')
  } finally {
    searching.value = false
  }
}

function handleClearSearch() {
  searchKeyword.value = ''
  loadSubscriptions()
}

function onPaginationChange(page) {
  activePage.value = page
}

function onPageSizeChange(s) {
  pageSize.value = s
  activePage.value = 1
}

function refresh() {
  loadSubscriptions()
}

// ============ 用户名查找（参考 web-back） ============
function getUserName(userId) {
  if (!userId) return '-'
  const user = users.value.find((u) => u.id === userId)
  return user ? user.username : `#${userId}`
}

// ============ 渲染函数（参考 web-back） ============
function renderStatus(status) {
  if (typeof status === 'number') return status === 1 ? '生效中' : '已过期'
  switch (status) {
    case 'active': return '生效中'
    case 'expired': return '已过期'
    case 'pending': return '待生效'
    case 'cancelled': return '已取消'
    default: return status || '未知'
  }
}

function statusClass(status) {
  if (typeof status === 'number') return status === 1 ? 'status-on' : 'status-off'
  switch (status) {
    case 'active': return 'status-on'
    case 'expired': return 'status-warn'
    case 'pending': return 'status-pending'
    default: return 'status-off'
  }
}

function renderBillingType(type) {
  if (type === 'token') return '按 Token'
  if (type === 'request') return '按次数'
  return type || '-'
}

function renderWindowType(wt) {
  const map = { period: '周期', week: '周', month: '月' }
  return map[wt] || wt
}

function isActive(status) {
  if (typeof status === 'number') return status === 1
  return status === 'active'
}

// ============ 格式化 ============
function formatTime(val) {
  if (!val) return '-'
  const t = Number(val)
  if (!isNaN(t) && t > 0) return new Date(t * 1000).toLocaleString()
  return val
}

function formatNumber(num) {
  if (num === null || num === undefined) return '-'
  const n = Number(num)
  if (isNaN(n)) return String(num)
  if (n >= 1000000) return `${(n / 1000000).toFixed(1)}M`
  if (n >= 1000) return `${(n / 1000).toFixed(1)}K`
  return String(n)
}

function formatPercent(val) {
  if (val === null || val === undefined) return '-'
  const n = Number(val)
  if (isNaN(n)) return '-'
  return `${n.toFixed(2)}%`
}

function windowPctClass(val) {
  const n = Number(val) || 0
  if (n >= 100) return 'pct-danger'
  if (n >= 80) return 'pct-warning'
  return 'pct-normal'
}

// ============ 进度条分段（参考 web-back MultiSegmentProgress） ============
function getLimitForWindow(wt, limits, billingType) {
  if (!limits) return 0
  let rule = limits.other
  if (!rule) {
    const keys = Object.keys(limits)
    if (keys.length === 0) return 0
    rule = limits[keys[0]]
  }
  if (billingType === 'token') {
    if (wt === 'period') return rule.token_period
    if (wt === 'week') return rule.token_week
    if (wt === 'month') return rule.token_month
  } else {
    if (wt === 'period') return rule.request_period
    if (wt === 'week') return rule.request_week
    if (wt === 'month') return rule.request_month
  }
  return 0
}

function getSegments(wt) {
  if (!usageData.value?.model_usage) return []
  const details = usageData.value.model_usage[wt] || []
  if (details.length === 0) return []

  const billingType = usageData.value.billing_type || usageData.value.subscription?.billing_type
  const limit = getLimitForWindow(wt, usageData.value.limits, billingType)
  const total = limit > 0 ? limit : 1

  const isToken = billingType === 'token'
  const segments = []
  let offset = 0

  details.forEach((md, i) => {
    const value = isToken
      ? (md.prompt_tokens || 0) + (md.completion_tokens || 0)
      : (md.requests || 0)
    const width = (value / total) * 100
    segments.push({
      label: md.model,
      value,
      color: MODEL_COLORS[i % MODEL_COLORS.length],
      offset: Math.min(offset, 100),
      width: Math.max(0, Math.min(width, 100 - offset)),
    })
    offset += width
  })

  return segments
}

function getModelDetails(wt) {
  if (!usageData.value?.model_usage) return []
  return usageData.value.model_usage[wt] || []
}

function toggleWindow(wt) {
  expandedWindows.value = { ...expandedWindows.value, [wt]: !expandedWindows.value[wt] }
}

// ============ 操作 ============
async function openAddModal() {
  addForm.user_id = null
  addForm.plan_id = null
  addForm.billing_type = 'token'
  addForm.duration_days = 30
  addForm.notes = ''
  showAddModal.value = true
  await Promise.all([loadUsers(), loadPlans()])
}

async function addSubscription() {
  if (!addForm.user_id || !addForm.plan_id) {
    Message.warning('请选择用户和套餐')
    return
  }
  addSubmitting.value = true
  try {
    const { data } = await api.post('/api/subscription/', {
      user_id: Number(addForm.user_id),
      plan_id: Number(addForm.plan_id),
      billing_type: addForm.billing_type,
      duration_days: addForm.duration_days ? Number(addForm.duration_days) : 0,
      notes: addForm.notes,
    })
    if (data.success) {
      Message.success('订阅已创建')
      showAddModal.value = false
      loadSubscriptions()
    } else {
      Message.error(data.message || '创建失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '创建失败')
  } finally {
    addSubmitting.value = false
  }
}

function openRenewModal(s) {
  // 参考 web-back：默认 +30 天
  const newEnd = (Number(s.end_time) || Math.floor(Date.now() / 1000)) + 30 * 86400
  renewForm.id = s.id
  renewForm.end_time = newEnd * 1000 // arco date-picker 用毫秒
  renewForm.notes = ''
  showRenewModal.value = true
}

async function renewSubscription() {
  if (!renewForm.end_time) {
    Message.warning('请选择新的结束时间')
    return
  }
  renewSubmitting.value = true
  try {
    const endTs = typeof renewForm.end_time === 'number'
      ? Math.floor(renewForm.end_time / 1000)
      : Math.floor(new Date(renewForm.end_time).getTime() / 1000)
    const { data } = await api.post(`/api/subscription/${renewForm.id}/renew`, {
      end_time: endTs,
      notes: renewForm.notes,
    })
    if (data.success) {
      Message.success('续费成功')
      showRenewModal.value = false
      loadSubscriptions()
    } else {
      Message.error(data.message || '续费失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '续费失败')
  } finally {
    renewSubmitting.value = false
  }
}

async function expireSubscription(id) {
  try {
    const { data } = await api.post(`/api/subscription/${id}/expire`)
    if (data.success) {
      Message.success('已置为过期')
      loadSubscriptions()
    } else {
      Message.error(data.message || '操作失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '操作失败')
  }
}

async function deleteSubscription(id) {
  try {
    const { data } = await api.delete(`/api/subscription/${id}/`)
    if (data.success) {
      Message.success('订阅已删除')
      // 当前页删除空时回退到上一页
      if (pageItems.value.length === 1 && activePage.value > 1) {
        activePage.value -= 1
      }
      loadSubscriptions()
    } else {
      Message.error(data.message || '删除失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '删除失败')
  }
}

async function viewUsage(record) {
  showUsageModal.value = true
  usageData.value = null
  expandedWindows.value = {}
  usageLoading.value = true
  try {
    const { data } = await api.get(`/api/subscription/${record.id}/usage`)
    if (data.success) {
      usageData.value = data.data
    } else {
      Message.error(data.message || '加载用量失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '加载用量失败')
  } finally {
    usageLoading.value = false
  }
}

onMounted(() => {
  loadSubscriptions()
  if (authStore.isAdmin) {
    loadUsers()
    loadPlans()
  }
})
</script>

<style scoped>
.sub-page {
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
.welcome-meta {
  display: flex;
  gap: 6px;
}
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
.list-row:last-child {
  border-bottom: none;
}
.list-row:hover {
  background: var(--color-fill-1);
}

/* ============ 单元格 ============ */
.col {
  font-size: 13px;
  color: var(--color-text-2);
  min-width: 0;
  padding-right: 16px;
}
.col:last-child {
  padding-right: 0;
}
.col-action {
  display: flex;
  justify-content: flex-end;
  gap: 0;
}
.col-action :deep(.arco-btn) {
  padding: 0 6px;
}

.cell-mono {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 12px;
  color: var(--color-text-2);
  font-variant-numeric: tabular-nums;
}
.cell-strong {
  color: var(--color-text-1);
  font-weight: 500;
}
.cell-muted {
  color: var(--color-text-3);
}
.cell-num {
  font-variant-numeric: tabular-nums;
  font-weight: 500;
  color: var(--color-text-1);
}

.billing-chip {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 10px;
  font-size: 12px;
  font-weight: 500;
  background: var(--color-fill-2);
  color: var(--color-text-2);
}

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
.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}
.status-on {
  background: rgba(0, 180, 42, 0.08);
  color: #00b42a;
}
.status-on .status-dot {
  background: #00b42a;
}
.status-off {
  background: var(--color-fill-2);
  color: var(--color-text-3);
}
.status-off .status-dot {
  background: var(--color-text-4);
}
.status-warn {
  background: rgba(245, 63, 63, 0.08);
  color: #f53f3f;
}
.status-warn .status-dot {
  background: #f53f3f;
}
.status-pending {
  background: rgba(255, 125, 0, 0.08);
  color: #ff7d00;
}
.status-pending .status-dot {
  background: #ff7d00;
}

.danger-btn {
  color: var(--color-text-2);
}
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
.empty-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-1);
  margin: 0 0 4px;
}
.empty-desc {
  font-size: 13px;
  color: var(--color-text-3);
  margin: 0 0 16px;
}

/* ============ 表单 ============ */
.sub-form :deep(.arco-form-item) {
  margin-bottom: 16px;
}
.sub-form :deep(.arco-form-item-label) {
  font-weight: 500;
  font-size: 13px;
  color: var(--color-text-2);
}
.opt-recommend {
  color: #ff7d00;
  font-weight: 600;
  margin-left: 4px;
}

/* ============ 用量弹窗 ============ */
.usage-info {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px 24px;
  padding: 12px 16px;
  background: var(--color-fill-1);
  border-radius: 6px;
  margin-bottom: 16px;
}
.ud-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.ud-label {
  font-size: 12px;
  color: var(--color-text-3);
}
.ud-value {
  font-size: 13px;
  color: var(--color-text-1);
  font-weight: 500;
}
.ud-mono {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 12px;
  font-variant-numeric: tabular-nums;
  font-weight: 400;
}

.usage-section {
  margin-top: 4px;
}
.section-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-1);
  margin: 0 0 12px;
}

.window-block {
  margin-bottom: 12px;
  border: 1px solid var(--color-fill-3);
  border-radius: 6px;
  overflow: hidden;
}
.window-head {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  background: var(--color-fill-1);
  cursor: pointer;
  user-select: none;
}
.window-head:hover {
  background: var(--color-fill-2);
}
.window-toggle {
  display: inline-flex;
  align-items: center;
  color: var(--color-text-3);
}
.window-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-1);
}
.window-pct {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 13px;
  font-weight: 600;
  margin-left: 4px;
}
.window-pct.pct-normal {
  color: #00b42a;
}
.window-pct.pct-warning {
  color: #ff7d00;
}
.window-pct.pct-danger {
  color: #f53f3f;
}
.window-reset {
  margin-left: auto;
  font-size: 12px;
  color: var(--color-text-3);
}

.window-segments {
  padding: 12px 14px;
  border-top: 1px solid var(--color-fill-3);
}
.segment-bar {
  position: relative;
  height: 20px;
  background: var(--color-fill-2);
  border-radius: 4px;
  overflow: hidden;
}
.segment {
  position: absolute;
  top: 0;
  bottom: 0;
  transition: width 0.3s ease;
  opacity: 0.85;
}
.segment:hover {
  opacity: 1;
}

.segment-legend {
  display: flex;
  flex-wrap: wrap;
  gap: 6px 14px;
  margin-top: 10px;
  font-size: 12px;
}
.legend-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--color-text-2);
}
.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 2px;
  display: inline-block;
  flex-shrink: 0;
}
.legend-name {
  font-weight: 500;
}
.legend-value {
  font-variant-numeric: tabular-nums;
  color: var(--color-text-3);
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 11px;
}

.window-detail {
  padding: 0 14px 12px;
  border-top: 1px solid var(--color-fill-3);
}
.detail-table {
  width: 100%;
  margin-top: 10px;
  border-collapse: separate;
  border-spacing: 0;
  font-size: 12px;
}
.detail-table th {
  background: var(--color-fill-1);
  color: var(--color-text-3);
  font-weight: 500;
  padding: 8px 10px;
  text-align: left;
  border-bottom: 1px solid var(--color-fill-3);
  white-space: nowrap;
}
.detail-table td {
  padding: 8px 10px;
  color: var(--color-text-2);
  border-bottom: 1px solid var(--color-fill-3);
  font-variant-numeric: tabular-nums;
}
.detail-table tbody tr:last-child td {
  border-bottom: none;
}
.detail-table tbody tr:hover {
  background: var(--color-fill-1);
}
.detail-model {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--color-text-1);
  font-weight: 500;
}

.window-empty {
  padding: 14px;
  font-size: 12px;
  color: var(--color-text-4);
  text-align: center;
  border-top: 1px solid var(--color-fill-3);
}

.usage-empty {
  text-align: center;
  padding: 40px 0;
  color: var(--color-text-3);
  font-size: 13px;
}
</style>
