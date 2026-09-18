<template>
  <div class="user-page">
    <!-- 顶部欢迎条 -->
    <div class="welcome-bar">
      <div class="welcome-text">
        <h1 class="welcome-title">{{ $t('userPage.title') }}</h1>
        <p class="welcome-desc">{{ $t('userPage.subtitle') }}</p>
      </div>
      <div class="welcome-meta">
        <span class="meta-chip">{{ $t('userPage.totalPeople', { n: pageTotal || users.length }) }}</span>
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
      <div v-if="pageItems.length === 0 && !loading && !loadingMore" class="empty-state">
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
                    {{ u.unlimited_quota ? $t('userPage.unlimited') : formatNumber(u.quota) }}
                  </span>
                </a-tooltip>
                <div class="quota-meta">
                  <a-tooltip :content="$t('userPage.usedQuota')">
                    <span class="quota-used">{{ $t('userPage.used', { value: formatNumber(u.used_quota) }) }}</span>
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

        <div v-if="loadingMore" class="load-more-row">
          <a-spin :loading="true" :size="14" />
          <span class="load-more-text">{{ $t('userPage.loadingMore') }}</span>
        </div>
        <div
          v-else-if="isReachedEnd && users.length > pageSize && !loading"
          class="load-end-row"
        >
          {{ $t('userPage.allLoaded', { n: users.length }) }}
        </div>
      </div>

      <div v-if="users.length > 0 && !isSearchMode" class="list-footer">
        <div class="batch-bar">
          <span class="batch-count">{{ $t('userPage.selectedCount', { n: selectedUsernames.length }) }}</span>
          <a-dropdown trigger="click" position="top">
            <a-button
              :disabled="selectedUsernames.length === 0 || batchOperating"
              :loading="batchOperating"
            >
              {{ $t('userPage.batchOperate') }}
              <template #icon><icon-down :size="12" /></template>
            </a-button>
            <template #content>
              <a-doption
                :disabled="selectedUsernames.length === 0 || batchOperating"
                @click="runBatch('batch-delete')"
              >
                <template #icon><icon-delete /></template>
                {{ $t('userPage.batchDelete') }}
              </a-doption>
              <a-doption
                :disabled="selectedUsernames.length === 0 || batchOperating"
                @click="runBatch('batch-disable')"
              >
                <template #icon><icon-stop /></template>
                {{ $t('userPage.batchDisable') }}
              </a-doption>
            </template>
          </a-dropdown>
        </div>
        <a-pagination
          :current="activePage"
          :total="totalCountForPager"
          :page-size="pageSize"
          show-total
          show-page-size
          :page-size-options="[10, 20, 50]"
          size="small"
          @change="onPaginationChange"
          @page-size-change="handlePageSizeChange"
        />
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
          <a-input v-model="addForm.username" :placeholder="$t('userPage.usernamePlaceholder')" :max-length="32" allow-clear />
        </a-form-item>
        <a-form-item field="display_name" :label="$t('userPage.displayName')">
          <a-input v-model="addForm.display_name" :placeholder="$t('userPage.displayNamePlaceholder')" :max-length="64" allow-clear />
        </a-form-item>
        <a-form-item field="password" :label="$t('userPage.password')">
          <a-input-password v-model="addForm.password" :placeholder="$t('userPage.passwordPlaceholder')" :max-length="64" />
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
          <a-input v-model="editForm.display_name" :placeholder="$t('userPage.displayNamePlaceholder')" :max-length="64" allow-clear />
        </a-form-item>
        <a-form-item field="password" :label="$t('userPage.password')">
          <a-input-password v-model="editForm.password" :placeholder="$t('userPage.passwordKeepPlaceholder')" :max-length="64" />
        </a-form-item>
        <a-form-item field="group" :label="$t('userPage.group')">
          <a-select v-model="editForm.group" :placeholder="$t('userPage.groupPlaceholder')" allow-clear>
            <a-option v-for="g in groups" :key="g" :value="g">{{ g }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="quota" :label="$t('userPage.quota')">
          <a-input-number
            v-model="editForm.quota"
            :min="0"
            :max="99999999999"
            :placeholder="$t('userPage.quotaPlaceholder')"
            style="width: 100%"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { IconPlus, IconUserGroup, IconDown, IconDelete, IconStop } from '@arco-design/web-vue/es/icon'
import { Message } from '@arco-design/web-vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'

const authStore = useAuthStore()
const { t } = useI18n()

const ITEMS_PER_PAGE = 10

const loading = ref(true)
const loadingMore = ref(false)
const isReachedEnd = ref(false)
const submitting = ref(false)
const users = ref([])
const groups = ref([])
const subscriptions = ref([])
const keyword = ref('')
const isSearchMode = ref(false)
const orderBy = ref('')
const activePage = ref(1)
const pageSize = ref(20)

const pageItems = computed(() => {
  const start = (activePage.value - 1) * pageSize.value
  return users.value.slice(start, start + pageSize.value)
})

// 批量操作选中状态：跨页持久化
// Batch selection: persisted across pagination.
const selectedUsernames = ref([])
const batchOperating = ref(false)

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

// 操作成功后保留仍存在的 username，清理已不存在的
// After a successful operation, prune usernames no longer in the list.
function pruneSelection() {
  const live = new Set(users.value.map((u) => u.username))
  selectedUsernames.value = selectedUsernames.value.filter((n) => live.has(n))
}

const totalCountForPager = computed(() => {
  if (isReachedEnd.value) return users.value.length
  return users.value.length + pageSize.value
})

const addVisible = ref(false)
const addFormRef = ref(null)
const addForm = reactive({
  username: '',
  display_name: '',
  password: '',
})

const addRules = computed(() => ({
  username: [
    { required: true, message: t('userPage.usernameRequired') },
    { minLength: 3, message: t('userPage.usernameMin') },
  ],
  password: [
    { required: true, message: t('userPage.passwordRequired') },
    { minLength: 6, message: t('userPage.passwordMin') },
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

onMounted(async () => {
  await Promise.all([fetchUsers(), fetchGroups(), fetchSubscriptions()])
  loading.value = false
})

async function fetchUsers({ append = false, pageIdx = 0 } = {}) {
  if (append) loadingMore.value = true
  else loading.value = true
  try {
    const params = { p: pageIdx }
    if (orderBy.value) params.order = orderBy.value
    const { data } = await api.get('/api/user/', { params })
    if (data.success) {
      const list = Array.isArray(data.data) ? data.data : (data.data?.items || [])
      if (append) {
        if (list.length === 0) {
          isReachedEnd.value = true
        } else {
          users.value = [...users.value, ...list]
          if (list.length < pageSize.value) isReachedEnd.value = true
        }
      } else {
        users.value = list
        activePage.value = 1
        isReachedEnd.value = list.length < pageSize.value
      }
    } else {
      Message.error(data.message || t('userPage.fetchUsersFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('userPage.fetchUsersFailed'))
  } finally {
    loading.value = false
    loadingMore.value = false
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

async function handleSearch(val) {
  const term = (val || '').trim()
  if (!term) {
    handleClearSearch()
    return
  }
  loading.value = true
  isSearchMode.value = true
  keyword.value = term
  activePage.value = 1
  isReachedEnd.value = true // 搜索结果不支持追加
  try {
    const { data } = await api.get('/api/user/search', { params: { keyword: term } })
    if (data.success) {
      const items = Array.isArray(data.data) ? data.data : (data.data?.items || [])
      users.value = items
    } else {
      Message.error(data.message || t('userPage.searchFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('userPage.searchFailed'))
  } finally {
    loading.value = false
  }
}

async function handleClearSearch() {
  keyword.value = ''
  isSearchMode.value = false
  activePage.value = 1
  isReachedEnd.value = false
  await fetchUsers({ append: false })
}

function onPaginationChange(page) {
  activePage.value = page
  const totalPages = Math.ceil(users.value.length / pageSize.value)
  if (page > totalPages && !isReachedEnd.value && !loadingMore.value && !keyword.value) {
    const nextPageIdx = totalPages
    fetchUsers({ append: true, pageIdx: nextPageIdx })
  }
}

function handlePageSizeChange(s) {
  pageSize.value = s
  activePage.value = 1
}

function handleOrderByChange() {
  activePage.value = 1
  isReachedEnd.value = false
  fetchUsers({ append: false })
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
      activePage.value = 1
      await Promise.all([fetchUsers(), fetchSubscriptions()])
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
  submitting.value = true
  try {
    const payload = {
      username: editForm.username,
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
      await fetchUsers()
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
      if (pageItems.value.length === 1 && activePage.value > 1) {
        activePage.value -= 1
      }
      await Promise.all([fetchUsers(), fetchSubscriptions()])
    } else {
      Message.error(data.message || t('userPage.deleteFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('userPage.deleteFailed'))
  }
}

// 批量操作：复用 /api/user/manage，action = batch-delete / batch-disable
// Batch: reuses /api/user/manage with action = batch-delete / batch-disable.
async function runBatch(action) {
  const usernames = selectedUsernames.value.slice()
  if (usernames.length === 0) {
    Message.warning(t('userPage.batchEmptySelection'))
    return
  }
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
      await Promise.all([fetchUsers(), fetchSubscriptions()])
      // 保留仍存在的 username，清理已不存在的
      // Keep usernames still in the list, drop missing ones.
      pruneSelection()
      if (pageItems.value.length === 0 && activePage.value > 1) {
        activePage.value -= 1
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
  }
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
  color: var(--color-text-2);
}
.danger-btn:hover {
  color: #f53f3f !important;
  background: rgba(245, 63, 63, 0.06) !important;
}

.list-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 20px;
  border-top: 1px solid var(--color-fill-3);
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

.user-form :deep(.arco-form-item) {
  margin-bottom: 16px;
}
.user-form :deep(.arco-form-item-label) {
  font-weight: 500;
  font-size: 13px;
  color: var(--color-text-2);
}
</style>
