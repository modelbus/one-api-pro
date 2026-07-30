<template>
  <a-card :bordered="false" class="table-card">
    <div class="action-bar">
      <a-input-search
        v-model="keyword"
        placeholder="搜索用户名或显示名称..."
        allow-clear
        @search="handleSearch"
        @clear="handleClearSearch"
        :style="{ width: '220px' }"
      />
      <div class="action-bar-right">
        <a-space>
          <a-button type="primary" @click="openAddModal">
            <template #icon><icon-plus /></template>
            添加用户
          </a-button>
          <a-button @click="openSelfEditModal">
            <template #icon><icon-edit /></template>
            编辑个人资料
          </a-button>
        </a-space>
      </div>
    </div>

    <a-table
      :columns="columns"
      :data="users"
      :loading="loading"
      :pagination="false"
      :bordered="{ wrapper: true, cell: false }"
      size="medium"
      :scroll="{ x: 1200 }"
      row-key="id"
    >
      <template #role="{ record }">
        <a-tag :color="getRoleColor(record.role)">{{ getRoleLabel(record.role) }}</a-tag>
      </template>
      <template #quota="{ record }">
        {{ formatNumber(record.quota) }}
      </template>
      <template #used_quota="{ record }">
        {{ formatNumber(record.used_quota) }}
      </template>
      <template #status="{ record }">
        <a-tag :color="record.status === 1 ? 'green' : 'red'">
          {{ record.status === 1 ? '启用' : '禁用' }}
        </a-tag>
      </template>
      <template #actions="{ record }">
        <a-space>
          <a-button
            size="small"
            @click="openEditModal(record)"
            :disabled="record.role >= 100"
          >
            <template #icon><icon-edit /></template>
          </a-button>
          <a-popconfirm
            :content="record.status === 1 ? '确定要禁用该用户吗？' : '确定要启用该用户吗？'"
            @ok="toggleStatus(record)"
          >
            <a-button
              size="small"
              :status="record.status === 1 ? 'warning' : 'success'"
              :disabled="record.role >= 100"
            >
              {{ record.status === 1 ? '禁用' : '启用' }}
            </a-button>
          </a-popconfirm>
          <a-popconfirm
            :content="record.role >= 10 ? '确定要降级该用户为普通用户吗？' : '确定要提升该用户为管理员吗？'"
            @ok="togglePromote(record)"
          >
            <a-button
              size="small"
              :disabled="record.role >= 100"
            >
              {{ record.role >= 10 ? '降级' : '提升' }}
            </a-button>
          </a-popconfirm>
          <a-popconfirm
            content="确定要删除该用户吗？此操作不可撤销。"
            @ok="deleteUser(record)"
          >
            <a-button
              size="small"
              status="danger"
              :disabled="record.role >= 100"
            >
              删除
            </a-button>
          </a-popconfirm>
        </a-space>
      </template>
      <template #empty>
        <a-empty description="暂无用户数据" />
      </template>
    </a-table>

    <div v-if="!isSearchMode" class="table-page-footer">
      <a-pagination
        v-model:current="currentPage"
        :total="pageTotal"
        :page-size="pageSize"
        show-total
        @change="handlePageChange"
      />
    </div>
  </a-card>

  <a-modal
    v-model:visible="addVisible"
    title="添加用户"
    @ok="handleAddUser"
    @cancel="resetAddForm"
    :ok-loading="submitting"
  >
    <a-form ref="addFormRef" :model="addForm" :rules="addRules" layout="vertical">
      <a-form-item field="username" label="用户名">
        <a-input v-model="addForm.username" placeholder="请输入用户名" :max-length="32" />
      </a-form-item>
      <a-form-item field="display_name" label="显示名称">
        <a-input v-model="addForm.display_name" placeholder="请输入显示名称" :max-length="64" />
      </a-form-item>
      <a-form-item field="password" label="密码">
        <a-input-password v-model="addForm.password" placeholder="请输入密码" :max-length="64" />
      </a-form-item>
    </a-form>
  </a-modal>

  <a-modal
    v-model:visible="editVisible"
    title="编辑用户"
    @ok="handleEditUser"
    @cancel="resetEditForm"
    :ok-loading="submitting"
  >
    <a-form ref="editFormRef" :model="editForm" layout="vertical">
      <a-form-item field="username" label="用户名">
        <a-input v-model="editForm.username" disabled />
      </a-form-item>
      <a-form-item field="display_name" label="显示名称">
        <a-input v-model="editForm.display_name" placeholder="请输入显示名称" :max-length="64" />
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

  <a-modal
    v-model:visible="selfEditVisible"
    title="编辑个人资料"
    @ok="handleSelfEdit"
    @cancel="resetSelfEditForm"
    :ok-loading="submitting"
  >
    <a-form ref="selfEditFormRef" :model="selfEditForm" layout="vertical">
      <a-form-item field="username" label="用户名">
        <a-input v-model="selfEditForm.username" disabled />
      </a-form-item>
      <a-form-item field="display_name" label="显示名称">
        <a-input v-model="selfEditForm.display_name" placeholder="请输入显示名称" :max-length="64" />
      </a-form-item>
      <a-form-item field="password" label="密码">
        <a-input-password v-model="selfEditForm.password" placeholder="留空则不修改密码" :max-length="64" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { IconPlus, IconEdit } from '@arco-design/web-vue/es/icon'
import { Message } from '@arco-design/web-vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'

const authStore = useAuthStore()

const loading = ref(true)
const submitting = ref(false)
const users = ref([])
const groups = ref([])
const keyword = ref('')
const isSearchMode = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const pageTotal = ref(0)
const reachedEnd = ref(false)

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

const columns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: '用户名', dataIndex: 'username', width: 140, ellipsis: true, tooltip: true },
  { title: '显示名称', dataIndex: 'display_name', width: 140, ellipsis: true, tooltip: true },
  { title: '分组', dataIndex: 'group', width: 120, ellipsis: true, tooltip: true },
  { title: '角色', slotName: 'role', width: 110 },
  { title: '额度', slotName: 'quota', width: 140 },
  { title: '已使用额度', slotName: 'used_quota', width: 140 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '操作', slotName: 'actions', width: 320, fixed: 'right' },
]

onMounted(async () => {
  await Promise.all([fetchUsers(), fetchGroups()])
  loading.value = false
})

async function fetchUsers() {
  try {
    const params = { p: currentPage.value - 1 }
    const { data } = await api.get('/api/user/', { params })
    if (data.success) {
      const items = Array.isArray(data.data) ? data.data : (data.data?.items || [])
      users.value = items
      if (items.length < pageSize.value) {
        reachedEnd.value = true
        pageTotal.value = (currentPage.value - 1) * pageSize.value + items.length
      } else {
        reachedEnd.value = false
        pageTotal.value = currentPage.value * pageSize.value + 1
      }
    } else {
      Message.error(data.message || '获取用户列表失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '获取用户列表失败')
  }
}

async function fetchGroups() {
  try {
    const { data } = await api.get('/api/group/')
    if (data.success) {
      const list = data.data || []
      groups.value = Array.isArray(list) ? list : list.map(g => g.name || g.key || g)
    }
  } catch (e) {
    groups.value = []
  }
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
  currentPage.value = 1
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
  currentPage.value = 1
  pageTotal.value = 0
  reachedEnd.value = false
  loading.value = true
  await fetchUsers()
  loading.value = false
}

async function handlePageChange(page) {
  currentPage.value = page
  loading.value = true
  await fetchUsers()
  loading.value = false
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
      currentPage.value = 1
      await fetchUsers()
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
      if (users.value.length === 1 && currentPage.value > 1) {
        currentPage.value -= 1
      }
      await fetchUsers()
    } else {
      Message.error(data.message || '删除失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '删除失败')
  }
}

function isSelf(record) {
  return authStore.user && authStore.user.username === record.username
}

function getRoleLabel(role) {
  if (role >= 100) return '超级管理员'
  if (role >= 10) return '管理员'
  return '普通用户'
}

function getRoleColor(role) {
  if (role >= 100) return 'orange'
  if (role >= 10) return 'blue'
  return 'gray'
}

function formatNumber(num) {
  if (num == null || num === undefined) return '-'
  return Number(num).toLocaleString()
}
</script>

<style scoped>
.table-card { border-radius: 6px; }
.action-bar { display: flex; align-items: center; gap: 12px; padding: 12px 16px; background: var(--color-fill-2); border-radius: 6px; margin-bottom: 15px; }
.action-bar-right { display: flex; align-items: center; gap: 8px; margin-left: auto; }
.table-page-footer { display: flex; justify-content: flex-end; margin-top: 20px; padding-top: 16px; border-top: 1px solid var(--color-border-2); }
</style>
