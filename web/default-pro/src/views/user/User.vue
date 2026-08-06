<template>
  <div class="user-page">
    <!-- 顶部欢迎条 -->
    <div class="welcome-bar">
      <div class="welcome-text">
        <h1 class="welcome-title">用户</h1>
        <p class="welcome-desc">管理系统用户、角色、额度与订阅套餐</p>
      </div>
      <div class="welcome-meta">
        <span class="meta-chip">共 {{ pageTotal || users.length }} 人</span>
      </div>
    </div>

    <!-- 独立搜索栏 -->
    <div class="search-card">
      <div class="search-left">
        <a-input-search
          v-model="keyword"
          placeholder="搜索用户名或显示名称..."
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
          <a-option value="">默认排序</a-option>
          <a-option value="quota">按剩余额度</a-option>
          <a-option value="used_quota">按已使用额度</a-option>
        </a-select>
      </div>
      <div class="search-right">
        <a-button @click="openSelfEditModal">
          <template #icon><icon-edit :size="14" /></template>
          编辑个人资料
        </a-button>
        <a-button type="primary" size="large" @click="openAddModal">
          <template #icon><icon-plus :size="14" /></template>
          添加用户
        </a-button>
      </div>
    </div>

    <!-- 列表 -->
    <div class="list-wrap">
      <div v-if="pageItems.length === 0 && !loading && !loadingMore" class="empty-state">
        <div class="empty-icon">
          <icon-user-group :size="32" />
        </div>
        <p class="empty-title">还没有任何用户</p>
        <p class="empty-desc">添加第一个用户即可开始</p>
        <a-button type="primary" @click="openAddModal">
          <template #icon><icon-plus :size="14" /></template>
          立即添加
        </a-button>
      </div>

      <div v-else class="list-body">
        <div class="list-head">
          <div class="col">ID</div>
          <div class="col">用户名</div>
          <div class="col">显示名称</div>
          <div class="col">分组</div>
          <div class="col">套餐</div>
          <div class="col">额度</div>
          <div class="col">角色</div>
          <div class="col">状态</div>
          <div class="col col-action">操作</div>
        </div>

        <a-spin :loading="loading" style="width: 100%">
          <div v-for="u in pageItems" :key="u.id" class="list-row">
            <div class="col">
              <span class="cell-mono">#{{ u.id }}</span>
            </div>

            <div class="col">
              <a-tooltip :content="u.email || '未绑定邮箱'">
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
                  <span class="plan-expire">{{ p.end_time ? formatTime(p.end_time) : '无期限' }}</span>
                </span>
                <span v-if="getUserPlans(u.id).length === 0" class="cell-muted">-</span>
              </div>
            </div>

            <!-- 额度列（优化样式：分两行展示剩余/已用/请求） -->
            <div class="col">
              <div class="quota-cell">
                <a-tooltip content="剩余额度">
                  <span :class="u.unlimited_quota ? 'quota-unlimited' : 'quota-value'">
                    {{ u.unlimited_quota ? '无限制' : formatNumber(u.quota) }}
                  </span>
                </a-tooltip>
                <div class="quota-meta">
                  <a-tooltip content="已使用">
                    <span class="quota-used">已用 {{ formatNumber(u.used_quota) }}</span>
                  </a-tooltip>
                  <a-tooltip v-if="u.request_count !== undefined" content="请求次数">
                    <span class="quota-used">· {{ formatNumber(u.request_count) }} 次</span>
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
                {{ u.status === 1 ? '启用' : '禁用' }}
              </span>
            </div>

            <div class="col col-action">
              <a-button type="text" size="small" :disabled="u.role >= 100" @click="openEditModal(u)">编辑</a-button>
              <a-popconfirm
                :content="u.status === 1 ? '确定要禁用该用户吗？' : '确定要启用该用户吗？'"
                @ok="toggleStatus(u)"
              >
                <a-button type="text" size="small" :disabled="u.role >= 100">
                  {{ u.status === 1 ? '禁用' : '启用' }}
                </a-button>
              </a-popconfirm>
              <a-popconfirm
                :content="u.role >= 10 ? '确定要降级该用户为普通用户吗？' : '确定要提升该用户为管理员吗？'"
                @ok="togglePromote(u)"
              >
                <a-button type="text" size="small" :disabled="u.role >= 100">
                  {{ u.role >= 10 ? '降级' : '提升' }}
                </a-button>
              </a-popconfirm>
              <a-popconfirm content="确定要删除该用户吗？此操作不可撤销。" @ok="deleteUser(u)">
                <a-button type="text" size="small" class="danger-btn" :disabled="u.role >= 100">删除</a-button>
              </a-popconfirm>
            </div>
          </div>
        </a-spin>

        <div v-if="loadingMore" class="load-more-row">
          <a-spin :loading="true" :size="14" />
          <span class="load-more-text">正在加载更多…</span>
        </div>
        <div
          v-else-if="isReachedEnd && users.length > pageSize && !loading"
          class="load-end-row"
        >
          已显示全部 {{ users.length }} 条数据
        </div>
      </div>

      <div v-if="users.length > 0 && !isSearchMode" class="list-footer">
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
      title="添加用户"
      :width="460"
      @ok="handleAddUser"
      @cancel="resetAddForm"
      :ok-loading="submitting"
      ok-text="保存"
      cancel-text="取消"
    >
      <a-form ref="addFormRef" :model="addForm" :rules="addRules" layout="vertical" class="user-form">
        <a-form-item field="username" label="用户名">
          <a-input v-model="addForm.username" placeholder="请输入用户名" :max-length="32" allow-clear />
        </a-form-item>
        <a-form-item field="display_name" label="显示名称">
          <a-input v-model="addForm.display_name" placeholder="请输入显示名称" :max-length="64" allow-clear />
        </a-form-item>
        <a-form-item field="password" label="密码">
          <a-input-password v-model="addForm.password" placeholder="请输入密码" :max-length="64" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 编辑用户 -->
    <a-modal
      v-model:visible="editVisible"
      title="编辑用户"
      :width="460"
      @ok="handleEditUser"
      @cancel="resetEditForm"
      :ok-loading="submitting"
      ok-text="保存"
      cancel-text="取消"
    >
      <a-form ref="editFormRef" :model="editForm" layout="vertical" class="user-form">
        <a-form-item field="username" label="用户名">
          <a-input v-model="editForm.username" disabled />
        </a-form-item>
        <a-form-item field="display_name" label="显示名称">
          <a-input v-model="editForm.display_name" placeholder="请输入显示名称" :max-length="64" allow-clear />
        </a-form-item>
        <a-form-item field="password" label="密码">
          <a-input-password v-model="editForm.password" placeholder="留空则不修改密码" :max-length="64" />
        </a-form-item>
        <a-form-item field="group" label="分组">
          <a-select v-model="editForm.group" placeholder="请选择分组" allow-clear>
            <a-option v-for="g in groups" :key="g" :value="g">{{ g }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="quota" label="额度">
          <a-input-number
            v-model="editForm.quota"
            :min="0"
            :max="99999999999"
            placeholder="请输入额度"
            style="width: 100%"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 编辑个人资料 -->
    <a-modal
      v-model:visible="selfEditVisible"
      title="编辑个人资料"
      :width="460"
      @ok="handleSelfEdit"
      @cancel="resetSelfEditForm"
      :ok-loading="submitting"
      ok-text="保存"
      cancel-text="取消"
    >
      <a-form ref="selfEditFormRef" :model="selfEditForm" layout="vertical" class="user-form">
        <a-form-item field="username" label="用户名">
          <a-input v-model="selfEditForm.username" disabled />
        </a-form-item>
        <a-form-item field="display_name" label="显示名称">
          <a-input v-model="selfEditForm.display_name" placeholder="请输入显示名称" :max-length="64" allow-clear />
        </a-form-item>
        <a-form-item field="password" label="密码">
          <a-input-password v-model="selfEditForm.password" placeholder="留空则不修改密码" :max-length="64" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { IconPlus, IconEdit, IconUserGroup } from '@arco-design/web-vue/es/icon'
import { Message } from '@arco-design/web-vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'

const authStore = useAuthStore()

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

const addRules = {
  username: [
    { required: true, message: '请输入用户名' },
    { minLength: 3, message: '用户名至少3个字符' },
  ],
  password: [
    { required: true, message: '请输入密码' },
    { minLength: 6, message: '密码至少6个字符' },
  ],
}

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

const selfEditVisible = ref(false)
const selfEditFormRef = ref(null)
const selfEditForm = reactive({
  username: '',
  display_name: '',
  password: '',
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
      Message.error(data.message || '获取用户列表失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '获取用户列表失败')
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
      Message.error(data.message || '搜索失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '搜索失败')
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
      Message.success('用户添加成功')
      addVisible.value = false
      resetAddForm()
      activePage.value = 1
      await Promise.all([fetchUsers(), fetchSubscriptions()])
    } else {
      Message.error(data.message || '添加用户失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '添加用户失败')
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
      Message.success('用户编辑成功')
      editVisible.value = false
      resetEditForm()
      await fetchUsers()
    } else {
      Message.error(data.message || '编辑用户失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '编辑用户失败')
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

function openSelfEditModal() {
  const user = authStore.user
  if (!user) {
    Message.warning('用户信息未加载')
    return
  }
  selfEditForm.username = user.username || ''
  selfEditForm.display_name = user.display_name || ''
  selfEditForm.password = ''
  selfEditFormRef.value?.clearValidate()
  selfEditVisible.value = true
}

async function handleSelfEdit() {
  submitting.value = true
  try {
    const payload = {
      display_name: selfEditForm.display_name,
    }
    if (selfEditForm.password) {
      payload.password = selfEditForm.password
    }
    const { data } = await api.put('/api/user/self', payload)
    if (data.success) {
      Message.success('个人资料更新成功')
      selfEditVisible.value = false
      resetSelfEditForm()
      if (data.data) {
        authStore.user = data.data
        localStorage.setItem('user', JSON.stringify(data.data))
      }
      await fetchUsers()
    } else {
      Message.error(data.message || '更新失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '更新失败')
  } finally {
    submitting.value = false
  }
}

function resetSelfEditForm() {
  selfEditForm.username = ''
  selfEditForm.display_name = ''
  selfEditForm.password = ''
  selfEditFormRef.value?.clearValidate()
}

async function toggleStatus(record) {
  try {
    const action = record.status === 1 ? 'disable' : 'enable'
    const { data } = await api.post('/api/user/manage', {
      username: record.username,
      action,
    })
    if (data.success) {
      Message.success(action === 'disable' ? '用户已禁用' : '用户已启用')
      record.status = record.status === 1 ? 0 : 1
    } else {
      Message.error(data.message || '操作失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '操作失败')
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
      Message.success(action === 'promote' ? '用户已提升为管理员' : '用户已降级为普通用户')
      record.role = action === 'promote' ? 10 : 0
    } else {
      Message.error(data.message || '操作失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '操作失败')
  }
}

async function deleteUser(record) {
  try {
    const { data } = await api.post('/api/user/manage', {
      username: record.username,
      action: 'delete',
    })
    if (data.success) {
      Message.success('用户已删除')
      if (pageItems.value.length === 1 && activePage.value > 1) {
        activePage.value -= 1
      }
      await Promise.all([fetchUsers(), fetchSubscriptions()])
    } else {
      Message.error(data.message || '删除失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '删除失败')
  }
}

function getRoleLabel(role) {
  if (role >= 100) return '超级管理员'
  if (role >= 10) return '管理员'
  return '普通用户'
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

function renderBilling(type) {
  if (type === 'token') return '按 Token'
  if (type === 'request') return '按请求'
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
  grid-template-columns: 80px 130px 130px 110px 220px 170px 110px 90px 240px;
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
  justify-content: flex-end;
  padding: 14px 20px;
  border-top: 1px solid var(--color-fill-3);
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
