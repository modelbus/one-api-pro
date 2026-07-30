<template>
  <a-card :bordered="false" class="table-card">
    <div class="action-bar">
      <a-input-search
        v-model="keyword"
        placeholder="Search subscriptions..."
        :style="{ width: '220px' }"
        @search="handleSearch"
        @clear="handleClearSearch"
        allow-clear
      />
      <div class="action-bar-right">
        <a-space>
          <a-button v-if="authStore.isAdmin" type="primary" @click="openAddModal">
            <template #icon><icon-plus /></template>
            New Subscription
          </a-button>
        </a-space>
      </div>
    </div>

    <a-table
      :columns="computedColumns"
      :data="subscriptions"
      :loading="loading"
      :pagination="false"
      row-key="id"
      size="medium"
      :bordered="{ wrapper: true, cell: false }"
      :scroll="{ x: 1000 }"
    >
      <template #user_name="{ record }">
        <span>{{ record.user_name ?? record.username ?? '-' }}</span>
      </template>
      <template #plan_name="{ record }">
        <span>{{ record.plan_name ?? record.plan ?? '-' }}</span>
      </template>
      <template #billing_type="{ record }">
        <a-tag :color="record.billing_type === 'token' ? 'blue' : 'purple'" size="small">
          {{ record.billing_type ?? '-' }}
        </a-tag>
      </template>
      <template #status="{ record }">
        <a-tag :color="getStatusColor(record.status)">
          {{ getStatusText(record.status) }}
        </a-tag>
      </template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="openUsageModal(record)">
            View Usage
          </a-button>
          <template v-if="authStore.isAdmin">
            <a-button
              v-if="isActive(record.status)"
              type="text"
              size="small"
              status="warning"
              @click="openRenewModal(record)"
            >
              Renew
            </a-button>
            <a-popconfirm
              v-if="isActive(record.status)"
              content="Are you sure you want to expire this subscription?"
              @ok="handleExpire(record)"
            >
              <a-button type="text" size="small" status="danger">
                Expire
              </a-button>
            </a-popconfirm>
            <a-popconfirm
              content="Are you sure you want to delete this subscription? This cannot be undone."
              @ok="handleDelete(record)"
            >
              <a-button type="text" size="small" status="danger">
                {{ $t('common.delete') }}
              </a-button>
            </a-popconfirm>
          </template>
        </a-space>
      </template>
      <template #empty>
        <a-empty description="No subscriptions found" />
      </template>
    </a-table>

    <div class="table-page-footer">
      <a-pagination
        v-model:current="pagination.current"
        :total="pagination.total"
        :page-size="pagination.pageSize"
        show-total
        @change="handlePageChange"
      />
    </div>
  </a-card>

  <a-modal
    v-model:visible="usageVisible"
    title="Subscription Usage"
    @cancel="usageVisible = false"
    :footer="false"
    :width="720"
  >
    <a-spin :loading="usageLoading">
      <template v-if="usageData">
        <a-descriptions :column="2" bordered size="small" class="usage-desc">
          <a-descriptions-item label="Plan">
            {{ usageData.plan_name ?? '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="Billing Type">
            {{ usageData.billing_type ?? '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="Start Time">
            {{ formatTime(usageData.start_time) }}
          </a-descriptions-item>
          <a-descriptions-item label="End Time">
            {{ formatTime(usageData.end_time) }}
          </a-descriptions-item>
        </a-descriptions>

        <a-divider style="margin: 16px 0" />

        <div v-if="usageData.windows" class="usage-windows">
          <h4 class="section-title">Usage Windows</h4>
          <a-row :gutter="16">
            <a-col :span="8" v-for="win in usageWindows" :key="win.key">
              <a-card :bordered="true" size="small" class="window-card">
                <template #title>
                  <span style="text-transform: capitalize">{{ win.key }}</span>
                </template>
                <a-statistic
                  :title="'Used'"
                  :value="formatQuota(win.used ?? 0)"
                  :precision="2"
                />
                <a-progress
                  v-if="win.limit"
                  :percent="getUsagePercent(win.used, win.limit)"
                  :status="getUsagePercent(win.used, win.limit) >= 90 ? 'danger' : getUsagePercent(win.used, win.limit) >= 70 ? 'warning' : 'success'"
                  :style="{ marginTop: '8px' }"
                />
                <div v-if="win.limit" style="font-size: 12px; color: var(--color-text-3); margin-top: 4px">
                  {{ formatQuota(win.used) }} / {{ formatQuota(win.limit) }}
                </div>
              </a-card>
            </a-col>
          </a-row>
        </div>

        <a-divider v-if="usageData.models" style="margin: 16px 0" />

        <div v-if="usageData.models && usageData.models.length" class="usage-models">
          <h4 class="section-title">Model Breakdown</h4>
          <a-table
            :columns="modelColumns"
            :data="usageData.models"
            :pagination="false"
            size="small"
            row-key="model"
          >
            <template #quota_used="{ record }">
              {{ formatQuota(record.quota_used ?? 0) }}
            </template>
            <template #count="{ record }">
              {{ formatNumber(record.count ?? 0) }}
            </template>
          </a-table>
        </div>

        <a-empty v-if="!usageData.windows && (!usageData.models || !usageData.models.length)" description="No usage data available" />
      </template>
      <a-empty v-else-if="!usageLoading" description="No usage data available" />
    </a-spin>
  </a-modal>

  <a-modal
    v-model:visible="renewVisible"
    title="Renew Subscription"
    @ok="handleRenew"
    @cancel="renewVisible = false"
    :ok-loading="renewSubmitting"
  >
    <a-form ref="renewFormRef" :model="renewForm" layout="vertical">
      <a-form-item field="end_time" label="New End Time">
        <a-date-picker
          v-model="renewForm.end_time"
          show-time
          format="YYYY-MM-DD HH:mm:ss"
          style="width: 100%"
        />
      </a-form-item>
      <a-form-item field="notes" label="Notes">
        <a-textarea
          v-model="renewForm.notes"
          placeholder="Renewal notes..."
          :auto-size="{ minRows: 2, maxRows: 4 }"
        />
      </a-form-item>
    </a-form>
  </a-modal>

  <a-modal
    v-model:visible="addVisible"
    title="Add Subscription"
    @ok="handleAdd"
    @cancel="addVisible = false"
    :ok-loading="addSubmitting"
  >
    <a-form ref="addFormRef" :model="addForm" layout="vertical">
      <a-form-item field="user_id" label="User">
        <a-select
          v-model="addForm.user_id"
          placeholder="Select user"
          allow-search
          :loading="userListLoading"
        >
          <a-option
            v-for="u in userList"
            :key="u.id"
            :value="u.id"
            :label="u.display_name || u.username"
          >
            #{{ u.id }} {{ u.display_name || u.username }}
          </a-option>
        </a-select>
      </a-form-item>
      <a-form-item field="plan_id" label="Plan">
        <a-select
          v-model="addForm.plan_id"
          placeholder="Select plan"
          allow-search
          :loading="planListLoading"
        >
          <a-option
            v-for="p in planList"
            :key="p.id"
            :value="p.id"
            :label="p.name"
          >
            {{ p.name }}
          </a-option>
        </a-select>
      </a-form-item>
      <a-form-item field="billing_type" label="Billing Type">
        <a-select v-model="addForm.billing_type" placeholder="Select billing type">
          <a-option value="token">Token</a-option>
          <a-option value="request">Request</a-option>
        </a-select>
      </a-form-item>
      <a-form-item field="duration_days" label="Duration (Days)">
        <a-input-number
          v-model="addForm.duration_days"
          :min="1"
          :max="3650"
          placeholder="e.g. 30"
          style="width: 100%"
        />
      </a-form-item>
      <a-form-item field="notes" label="Notes">
        <a-textarea
          v-model="addForm.notes"
          placeholder="Optional notes..."
          :auto-size="{ minRows: 2, maxRows: 4 }"
        />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'

const authStore = useAuthStore()

const loading = ref(false)
const subscriptions = ref([])
const keyword = ref('')
const isSearchMode = ref(false)

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showTotal: true,
})

const baseColumns = [
  { title: 'ID', dataIndex: 'id', width: 70 },
  { title: 'User', slotName: 'user_name', width: 140, ellipsis: true },
  { title: 'Plan', slotName: 'plan_name', width: 160, ellipsis: true },
  { title: 'Billing Type', slotName: 'billing_type', width: 110 },
  { title: 'Start Time', dataIndex: 'start_time', width: 170, render: ({ record }) => formatTime(record.start_time) },
  { title: 'End Time', dataIndex: 'end_time', width: 170, render: ({ record }) => formatTime(record.end_time) },
  { title: 'Status', slotName: 'status', width: 100 },
  { title: 'Actions', slotName: 'actions', width: authStore.isAdmin ? 280 : 80, fixed: 'right' },
]

const columnsNoUser = baseColumns.filter(c => c.dataIndex !== 'user_name' && c.slotName !== 'user_name')

const computedColumns = computed(() => (authStore.isAdmin ? baseColumns : columnsNoUser))

const usageVisible = ref(false)
const usageLoading = ref(false)
const usageData = ref(null)
const currentSubscription = ref(null)

const renewVisible = ref(false)
const renewSubmitting = ref(false)
const renewFormRef = ref(null)
const renewForm = reactive({
  end_time: '',
  notes: '',
})

const addVisible = ref(false)
const addSubmitting = ref(false)
const addFormRef = ref(null)
const addForm = reactive({
  user_id: null,
  plan_id: null,
  billing_type: 'token',
  duration_days: 30,
  notes: '',
})

const userList = ref([])
const userListLoading = ref(false)
const planList = ref([])
const planListLoading = ref(false)

const modelColumns = [
  { title: 'Model', dataIndex: 'model', ellipsis: true },
  { title: 'Quota Used', slotName: 'quota_used', width: 140 },
  { title: 'Requests', slotName: 'count', width: 100 },
]

const usageWindows = computed(() => {
  if (!usageData.value?.windows) return []
  return Object.entries(usageData.value.windows).map(([key, val]) => ({
    key,
    used: val?.used ?? 0,
    limit: val?.limit ?? 0,
  }))
})

function isActive(status) {
  if (typeof status === 'number') return status === 1
  return status === 'active'
}

function getStatusColor(status) {
  if (typeof status === 'number') {
    return status === 1 ? 'green' : 'red'
  }
  switch (status) {
    case 'active': return 'green'
    case 'expired': return 'red'
    case 'pending': return 'orange'
    case 'cancelled': return 'gray'
    default: return 'gray'
  }
}

function getStatusText(status) {
  if (typeof status === 'number') {
    return status === 1 ? 'Active' : 'Expired'
  }
  return status ?? '-'
}

function formatTime(val) {
  if (!val) return '-'
  const t = Number(val)
  if (!isNaN(t) && t > 0) {
    return new Date(t * 1000).toLocaleString()
  }
  return val
}

function formatQuota(val) {
  if (val == null) return '0'
  const n = Number(val)
  if (isNaN(n)) return String(val)
  if (n >= 1000000) return (n / 1000000).toFixed(2) + 'M'
  if (n >= 1000) return (n / 1000).toFixed(1) + 'k'
  return n.toFixed(2)
}

function formatNumber(val) {
  if (val == null) return '0'
  return Number(val).toLocaleString()
}

function getUsagePercent(used, limit) {
  if (!limit || limit <= 0) return 0
  return Math.min(100, Math.round((used / limit) * 100))
}

async function fetchData() {
  loading.value = true
  try {
    let url
    if (isSearchMode.value && keyword.value) {
      url = `/api/subscription/search?keyword=${encodeURIComponent(keyword.value)}`
    } else if (authStore.isAdmin) {
      url = `/api/subscription/?p=${pagination.current - 1}`
    } else {
      url = '/api/subscription/self'
    }
    const { data } = await api.get(url)
    if (data.success) {
      const items = Array.isArray(data.data) ? data.data : (data.data?.items ?? data.data ?? [])
      subscriptions.value = items
      pagination.total = data.total ?? data.data?.total ?? items.length
    } else {
      Message.error(data.message || 'Failed to fetch subscriptions')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || 'Failed to fetch subscriptions')
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.current = 1
  if (keyword.value.trim()) {
    isSearchMode.value = true
  }
  fetchData()
}

function handleClearSearch() {
  keyword.value = ''
  isSearchMode.value = false
  pagination.current = 1
  fetchData()
}

function handlePageChange(page) {
  pagination.current = page
  fetchData()
}

async function openUsageModal(record) {
  currentSubscription.value = record
  usageVisible.value = true
  usageData.value = null
  usageLoading.value = true
  try {
    const { data } = await api.get(`/api/subscription/${record.id}/usage`)
    if (data.success) {
      usageData.value = data.data
    } else {
      Message.error(data.message || 'Failed to fetch usage')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || 'Failed to fetch usage')
  } finally {
    usageLoading.value = false
  }
}

function openRenewModal(record) {
  currentSubscription.value = record
  renewForm.end_time = ''
  renewForm.notes = ''
  renewFormRef.value?.clearValidate()
  renewVisible.value = true
}

async function handleRenew() {
  renewSubmitting.value = true
  try {
    const payload = {
      end_time: renewForm.end_time ? Math.floor(new Date(renewForm.end_time).getTime() / 1000) : undefined,
      notes: renewForm.notes,
    }
    const { data } = await api.post(`/api/subscription/${currentSubscription.value.id}/renew`, payload)
    if (data.success) {
      Message.success('Subscription renewed successfully')
      renewVisible.value = false
      fetchData()
    } else {
      Message.error(data.message || 'Failed to renew')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || 'Failed to renew')
  } finally {
    renewSubmitting.value = false
  }
}

async function handleExpire(record) {
  try {
    const { data } = await api.post(`/api/subscription/${record.id}/expire`)
    if (data.success) {
      Message.success('Subscription expired')
      fetchData()
    } else {
      Message.error(data.message || 'Failed to expire')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || 'Failed to expire')
  }
}

async function handleDelete(record) {
  try {
    const { data } = await api.delete(`/api/subscription/${record.id}/`)
    if (data.success) {
      Message.success('Subscription deleted')
      fetchData()
    } else {
      Message.error(data.message || 'Failed to delete')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || 'Failed to delete')
  }
}

async function openAddModal() {
  addForm.user_id = null
  addForm.plan_id = null
  addForm.billing_type = 'token'
  addForm.duration_days = 30
  addForm.notes = ''
  addFormRef.value?.clearValidate()
  addVisible.value = true
  await Promise.all([fetchUserList(), fetchPlanList()])
}

async function fetchUserList() {
  userListLoading.value = true
  try {
    const { data } = await api.get('/api/user/?p=0')
    if (data.success) {
      userList.value = Array.isArray(data.data) ? data.data : (data.data?.items ?? [])
    }
  } catch {
    // ignore
  } finally {
    userListLoading.value = false
  }
}

async function fetchPlanList() {
  planListLoading.value = true
  try {
    const { data } = await api.get('/api/plan/?p=0')
    if (data.success) {
      planList.value = Array.isArray(data.data) ? data.data : (data.data?.items ?? [])
    }
  } catch {
    // ignore
  } finally {
    planListLoading.value = false
  }
}

async function handleAdd() {
  addSubmitting.value = true
  try {
    const { data } = await api.post('/api/subscription/', {
      user_id: addForm.user_id,
      plan_id: addForm.plan_id,
      billing_type: addForm.billing_type,
      duration_days: addForm.duration_days,
      notes: addForm.notes,
    })
    if (data.success) {
      Message.success('Subscription created')
      addVisible.value = false
      fetchData()
    } else {
      Message.error(data.message || 'Failed to create subscription')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || 'Failed to create subscription')
  } finally {
    addSubmitting.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.table-card { border-radius: 6px; }
.action-bar { display: flex; align-items: center; gap: 12px; padding: 12px 16px; background: var(--color-fill-2); border-radius: 6px; margin-bottom: 15px; }
.action-bar-right { display: flex; align-items: center; gap: 8px; margin-left: auto; }
.table-page-footer { display: flex; justify-content: flex-end; margin-top: 20px; padding-top: 16px; border-top: 1px solid var(--color-border-2); }

.section-title {
  margin: 0 0 12px 0;
  font-size: 15px;
  font-weight: 600;
}

.usage-desc {
  margin-bottom: 8px;
}

.usage-windows {
  margin-bottom: 4px;
}

.window-card {
  text-align: center;
}

.usage-models {
  margin-top: 4px;
}
</style>
