<template>
  <div class="redemption-page">
    <!-- 顶部欢迎条 -->
    <div class="welcome-bar">
      <div class="welcome-text">
        <h1 class="welcome-title">兑换码</h1>
        <p class="welcome-desc">生成和管理可用于兑换额度的兑换码</p>
      </div>
      <div class="welcome-meta">
        <span class="meta-chip">共 {{ redemptions.length }} 个</span>
      </div>
    </div>

    <!-- 独立搜索栏 -->
    <div class="search-card">
      <div class="search-left">
        <a-input-search
          v-model="keyword"
          placeholder="搜索名称或密钥..."
          allow-clear
          @search="handleSearch"
          @clear="handleSearchClear"
          :style="{ width: '320px' }"
        />
      </div>
      <div class="search-right">
        <a-button @click="refresh" :loading="loading && !loadingMore">
          <template #icon><icon-refresh :size="14" /></template>
          刷新
        </a-button>
        <a-button type="primary" size="large" @click="openCreateModal">
          <template #icon><icon-plus :size="14" /></template>
          生成兑换码
        </a-button>
      </div>
    </div>

    <!-- 列表 -->
    <div class="list-wrap">
      <div v-if="pageItems.length === 0 && !loading && !loadingMore" class="empty-state">
        <div class="empty-icon">
          <icon-gift :size="32" />
        </div>
        <p class="empty-title">还没有任何兑换码</p>
        <p class="empty-desc">生成第一个兑换码即可分发</p>
        <a-button type="primary" @click="openCreateModal">
          <template #icon><icon-plus :size="14" /></template>
          立即生成
        </a-button>
      </div>

      <div v-else class="list-body">
        <div class="list-head">
          <div class="col col-sort" :class="sortClass('id')" @click="sortBy('id')">ID</div>
          <div class="col col-sort" :class="sortClass('name')" @click="sortBy('name')">名称</div>
          <div class="col col-sort" :class="sortClass('status')" @click="sortBy('status')">状态</div>
          <div class="col col-num col-sort" :class="sortClass('quota')" @click="sortBy('quota')">额度</div>
          <div class="col col-sort" :class="sortClass('created_time')" @click="sortBy('created_time')">创建时间</div>
          <div class="col col-sort" :class="sortClass('redeemed_time')" @click="sortBy('redeemed_time')">兑换时间</div>
          <div class="col col-action">操作</div>
        </div>

        <a-spin :loading="loading && !loadingMore" style="width: 100%">
          <div v-for="r in pageItems" :key="r.id" class="list-row">
            <div class="col"><span class="cell-mono">#{{ r.id }}</span></div>

            <div class="col">
              <span class="cell-strong ellipsis" :title="r.name">{{ r.name || '未命名' }}</span>
            </div>

            <div class="col">
              <span class="status-chip" :class="statusClass(r.status)">
                <span class="status-dot"></span>
                {{ statusText(r.status) }}
              </span>
            </div>

            <div class="col col-num">
              <span class="cell-num">{{ formatQuota(r.quota) }}</span>
            </div>

            <div class="col">
              <span class="cell-mono">{{ formatTime(r.created_time) }}</span>
            </div>

            <div class="col">
              <span v-if="r.redeemed_time" class="cell-mono">{{ formatTime(r.redeemed_time) }}</span>
              <span v-else class="cell-muted">未兑换</span>
            </div>

            <div class="col col-action">
              <a-button type="text" size="small" @click="copyKey(r)">复制</a-button>
              <a-popconfirm
                v-if="r.status !== 3"
                :content="r.status === 1 ? '确定要禁用该兑换码？' : '确定要启用该兑换码？'"
                @ok="toggleStatus(r)"
              >
                <a-button type="text" size="small">
                  {{ r.status === 1 ? '禁用' : '启用' }}
                </a-button>
              </a-popconfirm>
              <a-button type="text" size="small" @click="openEditModal(r)">编辑</a-button>
              <a-popconfirm content="确定要删除该兑换码吗？" @ok="handleDelete(r)">
                <a-button type="text" size="small" class="danger-btn">删除</a-button>
              </a-popconfirm>
            </div>
          </div>
        </a-spin>

        <!-- 追加加载中提示 -->
        <div v-if="loadingMore" class="load-more-row">
          <a-spin :loading="true" :size="14" />
          <span class="load-more-text">正在加载更多…</span>
        </div>

        <!-- 末尾提示：已加载全部 -->
        <div
          v-else-if="isReachedEnd && redemptions.length > pageSize && !loading"
          class="load-end-row"
        >
          已显示全部 {{ redemptions.length }} 条数据
        </div>
      </div>

      <div v-if="redemptions.length > 0" class="list-footer">
        <a-pagination
          :current="activePage"
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

    <!-- 编辑 / 新建弹窗 -->
    <a-modal
      v-model:visible="modalVisible"
      :title="editingRecord ? '编辑兑换码' : '生成兑换码'"
      :width="460"
      @ok="handleSubmit"
      @cancel="modalVisible = false"
      :ok-loading="submitting"
      ok-text="保存"
      cancel-text="取消"
    >
      <a-form ref="formRef" :model="form" layout="vertical" class="redemption-form">
        <a-form-item field="name" label="名称" :rules="[{ required: true, message: '请输入名称' }]">
          <a-input v-model="form.name" placeholder="兑换码名称" allow-clear />
        </a-form-item>
        <a-form-item v-if="!editingRecord" field="count" label="生成数量" :rules="[{ required: true, message: '请输入数量' }]">
          <a-input-number
            v-model="form.count"
            :min="1"
            placeholder="生成数量"
            style="width: 100%"
          />
        </a-form-item>
        <a-form-item field="quota" label="额度" :rules="[{ required: true, message: '请输入额度' }]">
          <a-input-number
            v-model="form.quota"
            :min="0"
            :precision="2"
            placeholder="额度"
            style="width: 100%"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus, IconCopy, IconGift, IconRefresh } from '@arco-design/web-vue/es/icon'
import api from '@/api'

const ITEMS_PER_PAGE = 10

// ============ 状态 ============
const redemptions = ref([])
const loading = ref(false)        // 首次/覆盖加载
const loadingMore = ref(false)    // 追加加载
const isReachedEnd = ref(false)   // 触底标记
const keyword = ref('')
const activePage = ref(1)
const pageSize = ref(ITEMS_PER_PAGE)
const sortKey = ref('')
const sortDesc = ref(false)

const modalVisible = ref(false)
const submitting = ref(false)
const editingRecord = ref(null)
const formRef = ref(null)
const form = reactive({
  name: '',
  count: 1,
  quota: 0,
})

// ============ 状态映射（对齐 web-back：1=未用 2=禁用 3=已用）============
const statusMap = {
  1: { text: '未使用', class: 'status-on' },
  2: { text: '已禁用', class: 'status-off' },
  3: { text: '已使用', class: 'status-used' },
}
function statusText(s) {
  return statusMap[s]?.text || '未知'
}
function statusClass(s) {
  return statusMap[s]?.class || 'status-off'
}

// ============ 排序 + 分页 slice ============
const sortedRedemptions = computed(() => {
  if (!sortKey.value) return redemptions.value
  const list = [...redemptions.value]
  const k = sortKey.value
  list.sort((a, b) => {
    const va = a[k]
    const vb = b[k]
    if (va == null && vb == null) return 0
    if (va == null) return 1
    if (vb == null) return -1
    if (typeof va === 'number' && typeof vb === 'number') {
      return sortDesc.value ? vb - va : va - vb
    }
    return sortDesc.value
      ? String(vb).localeCompare(String(va))
      : String(va).localeCompare(String(vb))
  })
  return list
})

const pageItems = computed(() => {
  const start = (activePage.value - 1) * pageSize.value
  return sortedRedemptions.value.slice(start, start + pageSize.value)
})

// 分页器显示的总条数：
// - 若已触底：使用真实 redemptions.length
// - 否则额外加一页（让"下一页"按钮可点，触发追加加载，对齐 web-back）
const totalCountForPager = computed(() => {
  if (isReachedEnd.value) return redemptions.value.length
  return redemptions.value.length + pageSize.value
})

// ============ 数据加载（核心：reset / append 两种模式）============
async function fetchData({ append = false, pageIdx = 0 } = {}) {
  if (append) loadingMore.value = true
  else loading.value = true

  try {
    let url
    if (keyword.value) {
      // 搜索直接覆盖（不支持追加）
      url = `/api/redemption/search?keyword=${encodeURIComponent(keyword.value)}`
    } else {
      url = `/api/redemption/?p=${pageIdx}`
    }

    const { data } = await api.get(url)
    const items = data.data ?? data.results ?? []
    const list = Array.isArray(items) ? items : []

    if (append && !keyword.value) {
      // 追加模式
      if (list.length === 0) {
        isReachedEnd.value = true
      } else {
        redemptions.value = [...redemptions.value, ...list]
        // 如果返回数量小于 pageSize，说明已是末页
        if (list.length < pageSize.value) {
          isReachedEnd.value = true
        }
      }
    } else {
      // 覆盖模式（首次 / 刷新 / 搜索）
      redemptions.value = list
      activePage.value = 1
      isReachedEnd.value = list.length < pageSize.value
    }
  } catch (err) {
    Message.error(err?.response?.data?.message || '获取兑换码列表失败')
    if (!append) {
      redemptions.value = []
      activePage.value = 1
      isReachedEnd.value = true
    }
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

function handleSearch() {
  // 搜索：直接覆盖
  activePage.value = 1
  isReachedEnd.value = false
  fetchData({ append: false })
}

function handleSearchClear() {
  keyword.value = ''
  handleSearch()
}

function refresh() {
  // 刷新：覆盖
  isReachedEnd.value = false
  fetchData({ append: false })
}

function onPaginationChange(page) {
  activePage.value = page
  const totalPages = Math.ceil(redemptions.value.length / pageSize.value)
  // 关键：翻到"最后一页 + 1"且未触底时，触发追加加载（对齐 web-back）
  if (page > totalPages && !isReachedEnd.value && !loadingMore.value && !keyword.value) {
    const nextPageIdx = totalPages // 后端下一页索引 = 当前已加载的页数（从 0 开始）
    fetchData({ append: true, pageIdx: nextPageIdx })
  }
}

function onPageSizeChange(s) {
  pageSize.value = s
  // 切换每页大小时，保留追加数据不重置
  activePage.value = 1
}

// ============ 排序 ============
function sortBy(key) {
  if (sortKey.value === key) {
    sortDesc.value = !sortDesc.value
  } else {
    sortKey.value = key
    sortDesc.value = false
  }
  activePage.value = 1
}
function sortClass(key) {
  if (sortKey.value !== key) return ''
  return sortDesc.value ? 'active-desc' : 'active-asc'
}

// ============ 操作 ============
async function copyKey(record) {
  const key = record.key || record.show_key
  if (!key) {
    Message.warning('该兑换码无完整密钥')
    return
  }
  try {
    await navigator.clipboard.writeText(key)
    Message.success('已复制到剪贴板')
  } catch {
    Message.warning('复制失败，请手动复制')
    keyword.value = key
  }
}

async function openCreateModal() {
  editingRecord.value = null
  form.name = ''
  form.count = 1
  form.quota = 0
  modalVisible.value = true
}

async function openEditModal(record) {
  editingRecord.value = record
  form.name = record.name || ''
  form.quota = record.quota ?? 0
  modalVisible.value = true
}

async function handleSubmit() {
  const errors = await formRef.value?.validate()
  if (errors) return

  submitting.value = true
  try {
    if (editingRecord.value) {
      await api.put('/api/redemption/', {
        id: editingRecord.value.id,
        name: form.name,
        quota: form.quota,
      })
      Message.success('兑换码已更新')
    } else {
      await api.post('/api/redemption/', {
        name: form.name,
        count: form.count,
        quota: form.quota,
      })
      Message.success('兑换码已生成')
    }
    modalVisible.value = false
    refresh()
  } catch (err) {
    Message.error(err?.response?.data?.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

async function toggleStatus(record) {
  const newStatus = record.status === 1 ? 2 : 1
  try {
    await api.put('/api/redemption/', {
      id: record.id,
      status: newStatus,
    }, { params: { status_only: true } })
    Message.success(newStatus === 1 ? '已启用' : '已禁用')
    record.status = newStatus
  } catch (e) {
    Message.error(e?.response?.data?.message || e.message || '操作失败')
  }
}

async function handleDelete(record) {
  try {
    await api.delete(`/api/redemption/${record.id}/`)
    Message.success('兑换码已删除')
    // 当前页删除空时回退到上一页
    if (pageItems.value.length === 1 && activePage.value > 1) {
      activePage.value -= 1
    }
    refresh()
  } catch (err) {
    Message.error(err?.response?.data?.message || '删除兑换码失败')
  }
}

// ============ 格式化 ============
function formatQuota(val) {
  if (val == null) return '-'
  const n = Number(val)
  if (isNaN(n)) return val
  if (n >= 1000000) return `${(n / 1000000).toFixed(2)}M`
  if (n >= 10000) return `${(n / 10000).toFixed(2)}w`
  if (n >= 1000) return `${(n / 1000).toFixed(1)}k`
  return n.toLocaleString()
}

function formatTime(ts) {
  if (!ts) return '-'
  const t = Number(ts)
  if (!isNaN(t) && t > 0) return new Date(t * 1000).toLocaleString()
  return String(ts)
}

onMounted(() => {
  fetchData({ append: false })
})
</script>

<style scoped>
.redemption-page {
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

/* ============ 列表（与 Token 一致） ============ */
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
  grid-template-columns: 80px 160px 110px 130px 170px 170px 240px;
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
.col-num {
  text-align: right;
}
.col-action {
  display: flex;
  justify-content: flex-end;
  gap: 0;
}
.col-action :deep(.arco-btn) {
  padding: 0 6px;
}

.col-sort {
  cursor: pointer;
  user-select: none;
  transition: color 0.15s;
}
.col-sort:hover {
  color: var(--color-text-1);
}
.col-sort::after {
  content: '';
  display: inline-block;
  width: 0;
  height: 0;
  margin-left: 4px;
  vertical-align: middle;
  border-left: 4px solid transparent;
  border-right: 4px solid transparent;
  border-top: 4px solid currentColor;
  opacity: 0.25;
}
.col-sort.active-asc::after {
  border-top: none;
  border-bottom: 4px solid currentColor;
  opacity: 0.7;
}
.col-sort.active-desc::after {
  border-top: 4px solid currentColor;
  border-bottom: none;
  opacity: 0.7;
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
.ellipsis {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* ============ 追加加载 / 末尾提示 ============ */
.load-more-row,
.load-end-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px 20px;
  font-size: 12px;
  color: var(--color-text-3);
  border-top: 1px dashed var(--color-fill-3);
}
.load-more-text {
  color: var(--color-text-3);
}
.load-end-row {
  color: var(--color-text-4);
  background: var(--color-fill-1);
  border-top: 1px solid var(--color-fill-3);
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
.status-used {
  background: rgba(134, 144, 156, 0.10);
  color: #4e5969;
}
.status-used .status-dot {
  background: #86909c;
}

/* ============ 操作列 ============ */
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
.redemption-form :deep(.arco-form-item) {
  margin-bottom: 16px;
}
.redemption-form :deep(.arco-form-item-label) {
  font-weight: 500;
  font-size: 13px;
  color: var(--color-text-2);
}
</style>
