<template>
  <div class="user-page">
    <!-- 顶部欢迎条 -->
    <div class="welcome-bar">
      <div class="welcome-text">
        <h1 class="welcome-title">{{ $t('userPage.title') }}</h1>
        <p class="welcome-desc">{{ $t('userPage.subtitle') }}</p>
      </div>
      <div class="welcome-meta">
        <span class="meta-chip">{{ $t('userPage.totalPeople', { n: pagerTotal }) }}</span>
      </div>
    </div>

    <!-- 独立搜索栏 -->
    <div class="search-card">
      <div class="search-left">
        <a-input-search
          v-model="keyword"
          :placeholder="$t('userPage.searchPlaceholder')"
          allow-clear
          @search="handleSearch"
          @clear="handleClearSearch"
          :style="{ width: '320px' }"
        />
        <a-select
          v-model="orderBy"
          size="medium"
          :style="{ width: '160px' }"
          @change="handleOrderByChange"
        >
          <a-option value="">{{ $t('userPage.sortDefault') }}</a-option>
          <a-option value="quota">{{ $t('userPage.sortQuota') }}</a-option>
          <a-option value="used_quota">{{ $t('userPage.sortUsedQuota') }}</a-option>
        </a-select>
      </div>
      <div class="search-right">
        <a-button type="primary" size="large" @click="openAddModal">
          <template #icon><icon-plus :size="14" /></template>
          {{ $t('userPage.addUser') }}
        </a-button>
      </div>
    </div>

    <!-- 列表 -->
    <div class="list-wrap">
      <div v-if="pageItems.length === 0 && !loading" class="empty-state">
        <div class="empty-icon">
          <icon-user-group :size="32" />
        </div>
        <p class="empty-title">{{ $t('userPage.emptyTitle') }}</p>
        <p class="empty-desc">{{ $t('userPage.emptyDesc') }}</p>
        <a-button type="primary" @click="openAddModal">
          <template #icon><icon-plus :size="14" /></template>
          {{ $t('userPage.addNow') }}
        </a-button>
      </div>

      <div v-else class="list-body">
        <div class="list-head">
          <div class="col col-check">
            <a-checkbox
              :model-value="isAllSelected"
              :indeterminate="isPartiallySelected"
              :disabled="pageItems.length === 0"
              @change="toggleSelectAll"
            />
          </div>
          <div class="col">{{ $t('userPage.colId') }}</div>
          <div class="col">{{ $t('userPage.colUsername') }}</div>
          <div class="col">{{ $t('userPage.colDisplayName') }}</div>
          <div class="col">{{ $t('userPage.colGroup') }}</div>
          <div class="col">{{ $t('userPage.colPlan') }}</div>
          <div class="col">{{ $t('userPage.colQuota') }}</div>
          <div class="col">{{ $t('userPage.colRole') }}</div>
          <div class="col">{{ $t('userPage.colStatus') }}</div>
          <div class="col">{{ $t('userPage.colRegisteredAt') }}</div>
          <div class="col col-action">{{ $t('userPage.colAction') }}</div>
        </div>

        <a-spin :loading="loading" style="width: 100%">
          <div v-for="u in pageItems" :key="u.id" class="list-row">
            <div class="col col-check">
              <a-tooltip
                v-if="u.role >= 100"
                :content="$t('userPage.batchRootSkipped')"
                position="top"
              >
                <a-checkbox :model-value="false" disabled />
              </a-tooltip>
              <a-checkbox
                v-else
                :model-value="isUserSelected(u.username)"
                @change="(v) => toggleSelectOne(u.username, v)"
              />
            </div>
            <div class="col">
              <span class="cell-mono">#{{ u.id }}</span>
            </div>

            <div class="col">
              <a-tooltip :content="u.email || $t('userPage.noEmail')">
                <span class="cell-strong username-link">{{ u.username }}</span>
              </a-tooltip>
            </div>

            <div class="col">
              <span class="cell-muted ellipsis" :title="u.display_name">{{ u.display_name || '-' }}</span>
            </div>

            <div class="col">
              <span class="cell-muted ellipsis" :title="u.group">{{ u.group || '-' }}</span>
            </div>

            <!-- 套餐列 -->
            <div class="col">
              <div class="plan-cell">
                <span
                  v-for="(p, i) in getUserPlans(u.id)"
                  :key="i"
                  class="plan-chip"
                >
                  <span class="plan-name">{{ p.plan_name || '-' }}</span>
                  <span class="plan-sep">·</span>
                  <span class="plan-billing">{{ renderBilling(p.billing_type) }}</span>
                  <span class="plan-sep">·</span>
                  <span class="plan-expire">{{ p.end_time ? formatTime(p.end_time) : $t('userPage.noExpiry') }}</span>
                </span>
                <span v-if="getUserPlans(u.id).length === 0" class="cell-muted">-</span>
              </div>
            </div>

            <!-- 额度列（优化样式：分两行展示剩余/已用/请求） -->
            <div class="col">
              <div class="quota-cell">
                <a-tooltip :content="$t('userPage.remainingQuota')">
                  <span :class="u.unlimited_quota ? 'quota-unlimited' : 'quota-value'">
                    {{ u.unlimited_quota ? $t('userPage.unlimited') : fmtQuota(u.quota) }}
                  </span>
                </a-tooltip>
                <div class="quota-meta">
                  <a-tooltip :content="$t('userPage.usedQuota')">
                    <span class="quota-used">{{ $t('userPage.used', { value: fmtQuota(u.used_quota) }) }}</span>
                  </a-tooltip>
                  <a-tooltip v-if="u.request_count !== undefined" :content="$t('userPage.requestCount')">
                    <span class="quota-used">{{ $t('userPage.requestTimes', { n: formatNumber(u.request_count) }) }}</span>
                  </a-tooltip>
                </div>
              </div>
            </div>

            <div class="col">
              <span class="role-chip" :class="`role-${roleClass(u.role)}`">
                <span class="status-dot"></span>
                {{ getRoleLabel(u.role) }}
              </span>
            </div>

            <div class="col">
              <span class="status-chip" :class="u.status === 1 ? 'status-on' : 'status-off'">
                <span class="status-dot"></span>
                {{ u.status === 1 ? $t('userPage.enabled') : $t('userPage.disabled') }}
              </span>
            </div>

            <div class="col">
              <a-tooltip :content="formatDateTime(u.created_at)">
                <span class="cell-mono">{{ formatDate(u.created_at) }}</span>
              </a-tooltip>
            </div>

            <div class="col col-action">
              <a-button type="text" size="small" :disabled="u.role >= 100" @click="openEditModal(u)">{{ $t('userPage.edit') }}</a-button>
              <a-popconfirm
                :content="u.status === 1 ? $t('userPage.confirmDisable') : $t('userPage.confirmEnable')"
                @ok="toggleStatus(u)"
              >
                <a-button type="text" size="small" :disabled="u.role >= 100">
                  {{ u.status === 1 ? $t('userPage.disable') : $t('userPage.enable') }}
                </a-button>
              </a-popconfirm>
              <a-popconfirm
                :content="u.role >= 10 ? $t('userPage.confirmDemote') : $t('userPage.confirmPromote')"
                @ok="togglePromote(u)"
              >
                <a-button type="text" size="small" :disabled="u.role >= 100">
                  {{ u.role >= 10 ? $t('userPage.demote') : $t('userPage.promote') }}
                </a-button>
              </a-popconfirm>
              <a-popconfirm :content="$t('userPage.confirmDelete')" @ok="deleteUser(u)">
                <a-button type="text" size="small" class="danger-btn" :disabled="u.role >= 100">{{ $t('userPage.delete') }}</a-button>
              </a-popconfirm>
            </div>
          </div>
        </a-spin>

      </div>

      <div
        v-if="pagerTotal > 0 && !loading"
        class="list-footer"
        :class="{ 'list-footer--search': isSearchMode }"
      >
        <div v-if="!isSearchMode" class="batch-bar">
          <span class="batch-count">{{ $t('userPage.selectedCount', { n: selectedUsernames.length }) }}</span>
          <a-dropdown trigger="click" position="top">
            <a-button
              :disabled="selectedUsernames.length === 0 || batchOperating"
              :loading="batchOperating"
            >
              {{ $t('userPage.batchOperate') }}
              <icon-down :size="12" class="batch-trigger-arrow" />
            </a-button>
            <template #content>
              <a-doption
                class="batch-opt-danger"
                :disabled="selectedUsernames.length === 0 || batchOperating"
                @click="openBatchConfirm('batch-delete')"
              >
                <template #icon><icon-delete /></template>
                {{ $t('userPage.batchDelete') }}
              </a-doption>
              <a-doption
                :disabled="selectedUsernames.length === 0 || batchOperating"
                @click="openBatchConfirm('batch-disable')"
              >
                <template #icon><icon-stop /></template>
                {{ $t('userPage.batchDisable') }}
              </a-doption>
            </template>
          </a-dropdown>
        </div>
        <!-- 后端 /api/user/?page=&pageSize= 返回真实 total/page/page_size，
             分页器直接用服务端总数，不需要再推算页数；搜索模式为本地分页。 -->
        <a-pagination
          :current="activePage"
          :total="pagerTotal"
          :page-size="pageSize"
          :page-size-options="PAGE_SIZE_OPTIONS"
          :auto-adjust="false"
          show-total
          show-page-size
          size="small"
          @change="onPaginationChange"
          @page-size-change="onPageSizeChange"
        >
          <template #total>
            <span class="pager-summary">
              {{ $t('userPage.pagerSummary', { total: pagerTotal, page: activePage, pages: pagerPages }) }}
            </span>
          </template>
        </a-pagination>
      </div>
    </div>

    <!-- 添加用户 -->
    <a-modal
      v-model:visible="addVisible"
      :title="$t('userPage.addUser')"
      :width="460"
      @ok="handleAddUser"
      @cancel="resetAddForm"
      :ok-loading="submitting"
      :ok-text="$t('userPage.save')"
      :cancel-text="$t('userPage.cancel')"
    >
      <a-form ref="addFormRef" :model="addForm" :rules="addRules" layout="vertical" class="user-form">
        <a-form-item field="username" :label="$t('userPage.username')">
          <a-input v-model="addForm.username" :placeholder="$t('userPage.usernamePlaceholder')" :max-length="12" allow-clear />
        </a-form-item>
        <a-form-item field="display_name" :label="$t('userPage.displayName')">
          <a-input v-model="addForm.display_name" :placeholder="$t('userPage.displayNamePlaceholder')" :max-length="20" allow-clear />
        </a-form-item>
        <a-form-item field="password" :label="$t('userPage.password')">
          <a-input-password v-model="addForm.password" :placeholder="$t('userPage.passwordPlaceholder')" :max-length="20" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 编辑用户 -->
    <a-modal
      v-model:visible="editVisible"
      :title="$t('userPage.editUser')"
      :width="460"
      @ok="handleEditUser"
      @cancel="resetEditForm"
      :ok-loading="submitting"
      :ok-text="$t('userPage.save')"
      :cancel-text="$t('userPage.cancel')"
    >
      <a-form ref="editFormRef" :model="editForm" layout="vertical" class="user-form">
        <a-form-item field="username" :label="$t('userPage.username')">
          <a-input v-model="editForm.username" disabled />
        </a-form-item>
        <a-form-item field="display_name" :label="$t('userPage.displayName')">
          <a-input v-model="editForm.display_name" :placeholder="$t('userPage.displayNamePlaceholder')" :max-length="20" allow-clear />
        </a-form-item>
        <a-form-item field="password" :label="$t('userPage.password')">
          <a-input-password v-model="editForm.password" :placeholder="$t('userPage.passwordKeepPlaceholder')" :max-length="20" />
        </a-form-item>
        <a-form-item field="group" :label="$t('userPage.group')">
          <!-- 不允许清空：分组为空会让该用户匹配不到任何渠道，后端也会跳过空值写入 -->
          <a-select v-model="editForm.group" :placeholder="$t('userPage.groupPlaceholder')">
            <a-option v-for="g in groups" :key="g" :value="g">{{ g }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="quota" :label="$t('userPage.quota')">
          <a-input-number
            :model-value="editQuotaYuan"
            :min="0"
            :precision="yuanInputPrecision(editQuotaYuan)"
            :step="1"
            :placeholder="$t('userPage.quotaPlaceholder')"
            style="width: 100%"
            @change="onEditQuotaYuanChange"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 批量操作二次确认 Modal（arco Modal，per arco.design/vue/component/modal） -->
    <!-- Batch-action confirm Modal (arco a-modal, per arco docs). -->
    <a-modal
      :visible="batchConfirmVisible"
      :title="batchConfirmTitle"
      :ok-text="batchConfirmOkText"
      :cancel-text="$t('userPage.cancel')"
      :ok-button-props="batchConfirmOkProps"
      :ok-loading="batchOperating"
      @ok="confirmBatchOk"
      @cancel="cancelBatchConfirm"
      @update:visible="(v) => { if (!v) cancelBatchConfirm() }"
      :width="460"
    >
      <div class="batch-confirm-body">
        <icon-exclamation-circle :size="22" class="batch-confirm-icon" />
        <span>{{ batchConfirmContent }}</span>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { IconPlus, IconUserGroup, IconDown, IconDelete, IconStop, IconExclamationCircle } from '@arco-design/web-vue/es/icon'
import { Message } from '@arco-design/web-vue'
import { useAuthStore } from '@/stores/auth'
import { formatQuota as formatQuotaYuan, yuanInputPrecision } from '@/utils/quota'
import api from '@/api'

const authStore = useAuthStore()
const { t } = useI18n()

const ITEMS_PER_PAGE = 10
// 与后端 parseUserListPagination 的 maxUserPageSize(100) 对齐：允许的每页条数。
// Keep every option within the backend cap (maxUserPageSize = 100).
const PAGE_SIZE_OPTIONS = [10, 20, 50]

const loading = ref(true)
// 服务端返回的真实总数（搜索模式为本地结果条数）。
// Real total returned by the server (local result count in search mode).
const totalCount = ref(0)
const submitting = ref(false)
const users = ref([])
const groups = ref([])
const subscriptions = ref([])
const keyword = ref('')
const isSearchMode = ref(false)
const orderBy = ref('')
const activePage = ref(1)
// 后端支持 page(1 基) + pageSize（上限 maxUserPageSize=100），默认沿用 config.ItemsPerPage。
// Backend accepts 1-based `page` + `pageSize`; the default matches config.ItemsPerPage.
const pageSize = ref(ITEMS_PER_PAGE)

// 服务端分页下 users 即当前页；搜索结果一次返回全量，需要本地切片。
// Server-side paging: users already holds the current page. Search results are
// returned in full, so they are sliced locally.
const pageItems = computed(() => {
  if (!isSearchMode.value) return users.value
  const start = (activePage.value - 1) * pageSize.value
  return users.value.slice(start, start + pageSize.value)
})

// 批量操作选中状态：跨页持久化
// Batch selection: persisted across pagination.
const selectedUsernames = ref([])
const batchOperating = ref(false)
// 批量操作二次确认 Modal：单一 Modal 复用，按 action 切换文案/样式
// Batch-confirm modal: single modal reused, swap copy + style by action.
const batchConfirmVisible = ref(false)
const batchConfirmAction = ref('') // 'batch-delete' | 'batch-disable'

// 复选框工具函数：基于 username 切分，避免 id 在批量接口里被误用
// Checkbox helpers: use username as key (matches the batch API contract).
function isUserSelected(username) {
  return selectedUsernames.value.includes(username)
}
function toggleSelectOne(username, checked) {
  if (checked) {
    if (!selectedUsernames.value.includes(username)) {
      selectedUsernames.value.push(username)
    }
  } else {
    selectedUsernames.value = selectedUsernames.value.filter((n) => n !== username)
  }
}
function toggleSelectAll(checked) {
  if (checked) {
    // 仅勾选当前页里可操作的行（root 用户跳过）
    // Only tick rows that are operable on this page (skip root).
    const names = pageItems.value
      .filter((u) => u.role < 100)
      .map((u) => u.username)
    const merged = new Set(selectedUsernames.value)
    names.forEach((n) => merged.add(n))
    selectedUsernames.value = Array.from(merged)
  } else {
    const namesOnPage = new Set(pageItems.value.map((u) => u.username))
    selectedUsernames.value = selectedUsernames.value.filter((n) => !namesOnPage.has(n))
  }
}
const isAllSelected = computed(() => {
  const operable = pageItems.value.filter((u) => u.role < 100)
  if (operable.length === 0) return false
  return operable.every((u) => selectedUsernames.value.includes(u.username))
})
const isPartiallySelected = computed(() => {
  if (isAllSelected.value) return false
  return pageItems.value.some((u) => u.role < 100 && selectedUsernames.value.includes(u.username))
})

// 搜索模式下清空已选，避免跨上下文误操作
// Reset selection when entering search mode.
watch(isSearchMode, (val) => {
  if (val) selectedUsernames.value = []
})

// 批量操作成功后移除已成功处理的 username（跨页选择时不能用当前页裁剪）。
// Drop usernames reported as succeeded by the batch API.
function pruneSelection(succeeded) {
  if (!Array.isArray(succeeded) || succeeded.length === 0) return
  const done = new Set(succeeded)
  selectedUsernames.value = selectedUsernames.value.filter((n) => !done.has(n))
}

// 分页器总数：服务端直接返回 total，搜索模式（/api/user/search 返回全量）用结果条数。
// Pager total: the server returns `total`; in search mode (full result set) use its length.
const pagerTotal = computed(() => {
  if (isSearchMode.value) return users.value.length
  return totalCount.value
})
const pagerPages = computed(() =>
  Math.max(1, Math.ceil(pagerTotal.value / Math.max(1, pageSize.value))),
)

const addVisible = ref(false)
const addFormRef = ref(null)
const addForm = reactive({
  username: '',
  display_name: '',
  password: '',
})

const addRules = computed(() => ({
  // 与后端 model.User 的 validate 对齐：username max=12，password min=8/max=20。
  username: [
    { required: true, message: t('userPage.usernameRequired') },
    { minLength: 3, message: t('userPage.usernameMin') },
    { maxLength: 12, message: t('userPage.usernameMaxLength') },
  ],
  password: [
    { required: true, message: t('userPage.passwordRequired') },
    { minLength: 8, message: t('userPage.passwordMin') },
    { maxLength: 20, message: t('userPage.passwordMaxLength') },
  ],
}))

const editVisible = ref(false)
const editFormRef = ref(null)
const editingUser = ref(null)
const editForm = reactive({
  username: '',
  display_name: '',
  password: '',
  group: '',
  quota: 0,
})

// 管理端额度统一以「元」输入与提交（后端 API 自 v0.0.25 起同样使用元）。
// 版本: v0.0.25
// 日期: 2026-10-04
const editQuotaYuan = computed(() => Number(editForm.quota) || 0)

// fmtQuota 额度列直接展示（后端口径已是「元」）。
function fmtQuota(val) {
  if (val == null || val === '') return '-'
  return formatQuotaYuan(Number(val))
}

// onEditQuotaYuanChange 直接写入表单中的额度（元）。
function onEditQuotaYuanChange(v) {
  editForm.quota = Number(v) || 0
}

onMounted(async () => {
  await Promise.all([fetchUsers(), fetchGroups(), fetchSubscriptions()])
  loading.value = false
})

// 请求序号：快速切页/切搜索时只保留最后一次响应，避免旧响应把页码改回去。
// Request sequence: only the latest response wins when pages are switched quickly.
let fetchSeq = 0

// 后端契约：GET /api/user/?page=<1 基页号>&pageSize=<每页条数>&order=<quota|used_quota>
// 响应：data 为当前页数组，另有 total / page / page_size 三个字段，据此驱动分页器。
// Backend contract: 1-based `page` + `pageSize`; the response carries
// total / page / page_size (data stays an array).
async function fetchUsers({ page = activePage.value, size = pageSize.value } = {}) {
  const seq = ++fetchSeq
  loading.value = true
  try {
    const targetPage = Math.max(1, Math.floor(Number(page)) || 1)
    const targetSize = Math.max(1, Math.floor(Number(size)) || ITEMS_PER_PAGE)
    const params = { page: targetPage, pageSize: targetSize }
    if (orderBy.value) params.order = orderBy.value
    const { data } = await api.get('/api/user/', { params })
    if (seq !== fetchSeq) return
    if (data.success) {
      const list = Array.isArray(data.data) ? data.data : (data.data?.items || [])
      // 以后端回显的分页信息为准，避免前后端页号/页大小不一致。
      // Trust the echoed page/page_size so the pager never drifts from the server.
      const resolvedPage = Number(data.page) > 0 ? Number(data.page) : targetPage
      const resolvedSize = Number(data.page_size) > 0 ? Number(data.page_size) : targetSize
      const resolvedTotal = Number.isFinite(Number(data.total))
        ? Number(data.total)
        : (resolvedPage - 1) * resolvedSize + list.length
      // 当前页已越界（如批量删除后末页变空）时回退到最后一页，避免停在空页。
      // Fall back to the last page when the requested page went out of range.
      const lastPage = Math.max(1, Math.ceil(resolvedTotal / resolvedSize))
      if (list.length === 0 && resolvedPage > lastPage) {
        fetchUsers({ page: lastPage, size: resolvedSize })
        return
      }
      users.value = list
      activePage.value = resolvedPage
      pageSize.value = resolvedSize
      totalCount.value = resolvedTotal
    } else {
      Message.error(data.message || t('userPage.fetchUsersFailed'))
    }
  } catch (e) {
    if (seq !== fetchSeq) return
    Message.error(e.response?.data?.message || e.message || t('userPage.fetchUsersFailed'))
  } finally {
    if (seq === fetchSeq) loading.value = false
  }
}

async function fetchSubscriptions() {
  try {
    const { data } = await api.get('/api/subscription/?p=0')
    if (data.success) {
      subscriptions.value = Array.isArray(data.data) ? data.data : []
    }
  } catch (e) {
    subscriptions.value = []
  }
}

async function fetchGroups() {
  try {
    const { data } = await api.get('/api/group/')
    if (data.success) {
      const list = data.data || []
      groups.value = Array.isArray(list) ? list : list.map((g) => g.name || g.key || g)
    }
  } catch (e) {
    groups.value = []
  }
}

function getUserPlans(userId) {
  if (!userId) return []
  return subscriptions.value
    .filter((s) => s.user_id === userId || s.user?.id === userId)
    .map((s) => ({
      plan_name: s.plan?.name || s.plan_name || '-',
      billing_type: s.billing_type,
      end_time: s.end_time,
    }))
}

// 搜索结果接口一次返回全量，前端做本地分页。
// /api/user/search returns the full match set, so it is paginated client-side.
async function loadSearchResults(term) {
  const seq = ++fetchSeq
  loading.value = true
  isSearchMode.value = true
  keyword.value = term
  activePage.value = 1
  try {
    const { data } = await api.get('/api/user/search', { params: { keyword: term } })
    if (seq !== fetchSeq) return
    if (data.success) {
      users.value = Array.isArray(data.data) ? data.data : (data.data?.items || [])
    } else {
      Message.error(data.message || t('userPage.searchFailed'))
    }
  } catch (e) {
    if (seq !== fetchSeq) return
    Message.error(e.response?.data?.message || e.message || t('userPage.searchFailed'))
  } finally {
    if (seq === fetchSeq) loading.value = false
  }
}

async function handleSearch(val) {
  const term = (val || '').trim()
  if (!term) {
    await handleClearSearch()
    return
  }
  await loadSearchResults(term)
}

async function handleClearSearch() {
  keyword.value = ''
  isSearchMode.value = false
  activePage.value = 1
  await fetchUsers({ page: 1 })
}

// 列表刷新：搜索模式重跑搜索（保持本地分页），否则重新拉取指定页。
// Refresh helper: re-run search in search mode, otherwise re-fetch the given page.
async function refreshUsers({ page = activePage.value } = {}) {
  if (isSearchMode.value) {
    await loadSearchResults(keyword.value)
    return
  }
  await fetchUsers({ page })
}

function onPaginationChange(page) {
  const target = Math.floor(Number(page))
  if (!Number.isFinite(target) || target < 1 || target === activePage.value) return
  if (isSearchMode.value) {
    activePage.value = target
    return
  }
  fetchUsers({ page: target })
}

// 切换每页条数：回到第 1 页重新拉取（搜索模式只需本地重算）。
// Page-size switch: back to page 1 and re-fetch (search mode only re-slices locally).
function onPageSizeChange(size) {
  const target = Math.floor(Number(size))
  if (!Number.isFinite(target) || target < 1 || target === pageSize.value) return
  pageSize.value = target
  activePage.value = 1
  if (isSearchMode.value) return
  fetchUsers({ page: 1, size: target })
}

function handleOrderByChange() {
  // 排序只在服务端列表生效，若处于搜索状态先退出，避免两种数据源混用。
  if (isSearchMode.value) {
    keyword.value = ''
    isSearchMode.value = false
  }
  activePage.value = 1
  fetchUsers({ page: 1 })
}

function openAddModal() {
  addForm.username = ''
  addForm.display_name = ''
  addForm.password = ''
  addFormRef.value?.clearValidate()
  addVisible.value = true
}

async function handleAddUser() {
  const errors = await addFormRef.value?.validate()
  if (errors) return
  submitting.value = true
  try {
    const { data } = await api.post('/api/user/', {
      username: addForm.username,
      display_name: addForm.display_name,
      password: addForm.password,
    })
    if (data.success) {
      Message.success(t('userPage.addSuccess'))
      addVisible.value = false
      resetAddForm()
      // 新用户按 id desc 排序后位于第 1 页；搜索状态下先退出搜索再回到首页。
      keyword.value = ''
      isSearchMode.value = false
      activePage.value = 1
      await Promise.all([fetchUsers({ page: 1 }), fetchSubscriptions()])
    } else {
      Message.error(data.message || t('userPage.addFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('userPage.addFailed'))
  } finally {
    submitting.value = false
  }
}

function resetAddForm() {
  addForm.username = ''
  addForm.display_name = ''
  addForm.password = ''
  addFormRef.value?.clearValidate()
}

function openEditModal(record) {
  editingUser.value = record
  editForm.username = record.username
  editForm.display_name = record.display_name || ''
  editForm.password = ''
  editForm.group = record.group || ''
  editForm.quota = record.quota || 0
  editFormRef.value?.clearValidate()
  editVisible.value = true
}

async function handleEditUser() {
  // 后端 UpdateUser 以 body.id 定位目标用户（req.Id == 0 直接返回 invalid_parameter），
  // 因此必须回传被编辑用户的 id，否则永远提示「无效的参数」。
  if (!editingUser.value?.id) {
    Message.error(t('userPage.editFailed'))
    return
  }
  // 与后端 validate 对齐的前置校验：密码 8-20 字符（留空表示不修改）。
  if (editForm.password && editForm.password.length < 8) {
    Message.error(t('userPage.passwordMin'))
    return
  }
  submitting.value = true
  try {
    // username 不在提交范围：用户名不允许修改，且后端对 username 有 max=12 校验，
    // 历史超长用户名会被误判为「输入不合法」；不传即保持原值（后端跳过零值字段）。
    const payload = {
      id: editingUser.value.id,
      display_name: editForm.display_name,
      group: editForm.group,
      quota: editForm.quota,
    }
    if (editForm.password) {
      payload.password = editForm.password
    }
    const { data } = await api.put('/api/user/', payload)
    if (data.success) {
      Message.success(t('userPage.editSuccess'))
      editVisible.value = false
      resetEditForm()
      await refreshUsers()
    } else {
      Message.error(data.message || t('userPage.editFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('userPage.editFailed'))
  } finally {
    submitting.value = false
  }
}

function resetEditForm() {
  editingUser.value = null
  editForm.username = ''
  editForm.display_name = ''
  editForm.password = ''
  editForm.group = ''
  editForm.quota = 0
  editFormRef.value?.clearValidate()
}

async function toggleStatus(record) {
  try {
    const action = record.status === 1 ? 'disable' : 'enable'
    const { data } = await api.post('/api/user/manage', {
      username: record.username,
      action,
    })
    if (data.success) {
      Message.success(action === 'disable' ? t('userPage.userDisabled') : t('userPage.userEnabled'))
      record.status = record.status === 1 ? 0 : 1
    } else {
      Message.error(data.message || t('userPage.opFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('userPage.opFailed'))
  }
}

async function togglePromote(record) {
  try {
    const action = record.role >= 10 ? 'demote' : 'promote'
    const { data } = await api.post('/api/user/manage', {
      username: record.username,
      action,
    })
    if (data.success) {
      Message.success(action === 'promote' ? t('userPage.promoted') : t('userPage.demoted'))
      record.role = action === 'promote' ? 10 : 0
    } else {
      Message.error(data.message || t('userPage.opFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('userPage.opFailed'))
  }
}

async function deleteUser(record) {
  try {
    const { data } = await api.post('/api/user/manage', {
      username: record.username,
      action: 'delete',
    })
    if (data.success) {
      Message.success(t('userPage.deleteSuccess'))
      // 当前页删空后回到上一页，否则停留在当前页刷新。
      const targetPage =
        !isSearchMode.value && pageItems.value.length === 1 && activePage.value > 1
          ? activePage.value - 1
          : activePage.value
      await Promise.all([refreshUsers({ page: targetPage }), fetchSubscriptions()])
    } else {
      Message.error(data.message || t('userPage.deleteFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('userPage.deleteFailed'))
  }
}

// 批量操作：复用 /api/user/manage，action = batch-delete / batch-disable
// Batch: reuses /api/user/manage with action = batch-delete / batch-disable.
const batchConfirmTitle = computed(() => {
  return batchConfirmAction.value === 'batch-delete'
    ? t('userPage.batchConfirmDeleteTitle')
    : t('userPage.batchConfirmDisableTitle')
})
const batchConfirmContent = computed(() => {
  const n = selectedUsernames.value.length
  return batchConfirmAction.value === 'batch-delete'
    ? t('userPage.confirmBatchDelete', { n })
    : t('userPage.confirmBatchDisable', { n })
})
const batchConfirmOkText = computed(() => {
  return batchConfirmAction.value === 'batch-delete'
    ? t('userPage.batchDelete')
    : t('userPage.batchDisable')
})
const batchConfirmOkProps = computed(() => {
  if (batchConfirmAction.value === 'batch-delete') {
    return { status: 'danger' }
  }
  return {}
})

function openBatchConfirm(action) {
  if (selectedUsernames.value.length === 0) {
    Message.warning(t('userPage.batchEmptySelection'))
    return
  }
  batchConfirmAction.value = action
  batchConfirmVisible.value = true
}
function cancelBatchConfirm() {
  if (batchOperating.value) return
  batchConfirmVisible.value = false
  batchConfirmAction.value = ''
}

async function runBatch(action) {
  const usernames = selectedUsernames.value.slice()
  batchOperating.value = true
  try {
    const { data } = await api.post('/api/user/manage', { action, usernames })
    if (data.success && data.data) {
      const ok = data.data.succeeded_count || 0
      const fail = data.data.failed_count || 0
      if (fail > 0) {
        // 部分失败：提示并展示失败明细
        // Partial failure: show summary + per-user failure reasons.
        const failedMap = data.data.failed || {}
        const detail = Object.entries(failedMap)
          .map(([name, reason]) => `${name}: ${reason}`)
          .join('\n')
        Message.warning(
          `${t('userPage.batchPartialResult', { ok, fail })}\n${detail}`,
          6000,
        )
      } else {
        Message.success(
          action === 'batch-delete'
            ? t('userPage.deleteSuccess', { n: ok })
            : t('userPage.userDisabled', { n: ok }),
        )
      }
      // 刷新当前页；若本页已被清空则回退一页重新拉取。
      // Refresh current page, fall back one page if it became empty.
      const targetPage = activePage.value
      await Promise.all([refreshUsers({ page: targetPage }), fetchSubscriptions()])
      // 仅移除批量接口确认成功的 username（跨页选择不能被当前页裁剪）。
      // Only drop the usernames the batch API reported as succeeded.
      pruneSelection(data.data.succeeded)
      if (!isSearchMode.value && users.value.length === 0 && targetPage > 1) {
        await fetchUsers({ page: targetPage - 1 })
      }
    } else {
      // 全部失败：data.data 仍包含明细
      // All failed: data.data still carries the breakdown.
      const failedMap = (data.data && data.data.failed) || {}
      const detail = Object.entries(failedMap)
        .map(([name, reason]) => `${name}: ${reason}`)
        .join('\n')
      Message.error(
        `${data.message || t('userPage.userPageBatchFailed')}${detail ? '\n' + detail : ''}`,
        6000,
      )
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('userPage.userPageBatchFailed'))
  } finally {
    batchOperating.value = false
    batchConfirmVisible.value = false
    batchConfirmAction.value = ''
  }
}

async function confirmBatchOk() {
  if (batchOperating.value) return
  await runBatch(batchConfirmAction.value)
}

function getRoleLabel(role) {
  if (role >= 100) return t('userPage.roleRoot')
  if (role >= 10) return t('userPage.roleAdmin')
  return t('userPage.roleUser')
}

function roleClass(role) {
  if (role >= 100) return 'root'
  if (role >= 10) return 'admin'
  return 'user'
}

function formatNumber(num) {
  if (num == null || num === undefined) return '-'
  return Number(num).toLocaleString()
}

function formatTime(ts) {
  if (!ts) return ''
  const t = Number(ts)
  if (!isNaN(t) && t > 0) return new Date(t * 1000).toLocaleDateString()
  return ''
}

function formatDate(ts) {
  if (!ts) return '-'
  const t = Number(ts)
  if (isNaN(t) || t <= 0) return '-'
  const d = new Date(t * 1000)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

function formatDateTime(ts) {
  if (!ts) return ''
  const t = Number(ts)
  if (isNaN(t) || t <= 0) return ''
  return new Date(t * 1000).toLocaleString()
}

function renderBilling(type) {
  if (type === 'token') return t('userPage.billingToken')
  if (type === 'request') return t('userPage.billingRequest')
  return type || '-'
}
</script>

<style scoped>
.user-page {
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
.search-left {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
}
.search-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
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
  grid-template-columns: 40px 80px 130px 130px 110px 220px 170px 110px 90px 120px 240px;
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
.ellipsis {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.username-link {
  cursor: pointer;
}
.username-link:hover {
  color: rgb(var(--primary-6));
}

/* ============ 套餐列（参考 web-back 风格） ============ */
.plan-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
  align-items: flex-start;
  max-height: 56px;
  overflow: hidden;
  padding: 6px 0;
}
.plan-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
  background: rgba(0, 180, 167, 0.08);
  color: #00b4a7;
  border: 1px solid rgba(0, 180, 167, 0.18);
  cursor: default;
  max-width: 100%;
  line-height: 1.5;
  white-space: nowrap;
}
.plan-name {
  font-weight: 600;
  color: #00b4a7;
}
.plan-sep {
  color: rgba(0, 180, 167, 0.4);
  margin: 0 1px;
}
.plan-billing {
  color: var(--color-text-2);
  font-weight: 400;
}
.plan-expire {
  color: var(--color-text-3);
  font-weight: 400;
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 10.5px;
  font-variant-numeric: tabular-nums;
}

/* ============ 额度列（优化） ============ */
.quota-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}
.quota-value {
  font-variant-numeric: tabular-nums;
  font-weight: 600;
  font-size: 14px;
  color: var(--color-text-1);
  letter-spacing: -0.2px;
}
.quota-unlimited {
  color: rgb(var(--primary-6));
  font-weight: 600;
  font-size: 13px;
}
.quota-meta {
  display: flex;
  align-items: center;
  gap: 2px;
  font-size: 11px;
  color: var(--color-text-4);
  font-variant-numeric: tabular-nums;
}
.quota-used {
  cursor: default;
}
.quota-used:hover {
  color: var(--color-text-3);
}

/* ============ 状态 / 角色 chip ============ */
.status-chip,
.role-chip {
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

.role-root {
  background: rgba(245, 63, 63, 0.08);
  color: #f53f3f;
}
.role-root .status-dot {
  background: #f53f3f;
}
.role-admin {
  background: rgba(22, 93, 255, 0.08);
  color: #165dff;
}
.role-admin .status-dot {
  background: #165dff;
}
.role-user {
  background: var(--color-fill-2);
  color: var(--color-text-3);
}
.role-user .status-dot {
  background: var(--color-text-4);
}

.danger-btn {
  color: #f53f3f;
}
.danger-btn:hover {
  color: #f53f3f !important;
  background: rgba(245, 63, 63, 0.06) !important;
}
.danger-btn:disabled,
.danger-btn.arco-btn-disabled {
  color: var(--color-text-4) !important;
  background: transparent !important;
}

/* 下拉箭头放在按钮文字右侧（按 arco 推荐用法） */
/* Dropdown trigger arrow lives to the right of the button text. */
.batch-trigger-arrow {
  margin-left: 6px;
  color: var(--color-text-3);
}

/* 下拉项「批量删除」红色高亮（与单条删除按钮颜色一致） */
/* Batch-delete doption highlighted in red, matching single-row delete. */
.batch-opt-danger {
  color: #f53f3f !important;
}
.batch-opt-danger :deep(.arco-icon) {
  color: #f53f3f !important;
}
.batch-opt-danger:hover {
  background: rgba(245, 63, 63, 0.06) !important;
}

/* 批量确认 Modal 内容布局 */
/* Batch-confirm modal body. */
.batch-confirm-body {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  font-size: 14px;
  color: var(--color-text-2);
  line-height: 1.6;
}
.batch-confirm-icon {
  color: #f53f3f;
  flex-shrink: 0;
  margin-top: 1px;
}

.list-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 20px;
  border-top: 1px solid var(--color-fill-3);
}

/* 搜索模式隐藏批量栏，分页器靠右对齐。 */
.list-footer--search {
  justify-content: flex-end;
}

.pager-summary {
  font-variant-numeric: tabular-nums;
}

.batch-bar {
  display: flex;
  align-items: center;
  gap: 12px;
}
.batch-count {
  font-size: 12px;
  color: var(--color-text-3);
}

.col-check {
  display: flex;
  justify-content: center;
  align-items: center;
  padding-right: 0;
}

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


.user-form :deep(.arco-form-item) {
  margin-bottom: 16px;
}
.user-form :deep(.arco-form-item-label) {
  font-weight: 500;
  font-size: 13px;
  color: var(--color-text-2);
}
</style>
