<template>
  <a-card :bordered="false" class="table-card">
    <div class="action-bar">
      <a-input-search
        v-model="keyword"
        placeholder="搜索名称或密钥..."
        allow-clear
        @search="handleSearch"
        @clear="handleSearch"
        :style="{ width: '220px' }"
      />
      <div class="action-bar-right">
        <a-button type="primary" @click="openCreateModal">
          <template #icon><icon-plus /></template>
          新建令牌
        </a-button>
      </div>
    </div>

    <a-table
      :columns="columns"
      :data="tokens"
      :loading="loading"
      :pagination="false"
      row-key="id"
      size="medium"
      :bordered="{ wrapper: true, cell: false }"
      :scroll="{ x: 900 }"
    >
      <template #key="{ record }">
        <a-space>
          <span class="key-text">{{ maskKey(record.key) }}</span>
          <a-tooltip content="复制完整密钥">
            <a-button
              type="text"
              size="mini"
              @click="copyText(record.key, '完整密钥已复制')"
            >
              <template #icon><icon-copy /></template>
            </a-button>
          </a-tooltip>
        </a-space>
      </template>

      <template #status="{ record }">
        <a-tag :color="record.status === 1 ? 'green' : 'red'">
          {{ record.status === 1 ? '已启用' : '已禁用' }}
        </a-tag>
      </template>

      <template #remain_quota="{ record }">
        <span v-if="record.unlimited_quota" class="quota-unlimited">无限制</span>
        <span v-else>{{ formatQuota(record.remain_quota) }}</span>
      </template>

      <template #used_quota="{ record }">
        {{ formatQuota(record.used_quota) }}
      </template>

      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="openEditModal(record)">编辑</a-button>
          <a-popconfirm
            :content="record.status === 1 ? '确定要禁用该令牌？' : '确定要启用该令牌？'"
            @ok="toggleStatus(record)"
          >
            <a-button type="text" size="small" :status="record.status === 1 ? 'warning' : 'success'">
              {{ record.status === 1 ? '禁用' : '启用' }}
            </a-button>
          </a-popconfirm>
          <a-popconfirm content="确定要删除该令牌？此操作不可撤销" @ok="handleDelete(record.id)">
            <a-button type="text" size="small" status="danger">删除</a-button>
          </a-popconfirm>
        </a-space>
      </template>
    </a-table>

    <div class="table-page-footer">
      <a-pagination
        :current="page + 1"
        :total="total"
        :page-size="pageSize"
        show-total
        show-page-size
        @change="handlePageChange"
        @page-size-change="handlePageSizeChange"
      />
    </div>
  </a-card>

  <a-modal
    v-model:visible="modalVisible"
    :title="modalTitle"
    :width="600"
    @ok="handleSubmit"
    @cancel="closeModal"
    :ok-loading="submitting"
  >
    <a-form ref="formRef" :model="form" :rules="rules" layout="vertical" auto-label-width>
      <a-form-item field="name" label="名称" required>
        <a-input v-model="form.name" placeholder="令牌名称" :max-length="50" />
      </a-form-item>

      <a-form-item field="models" label="可用模型">
        <a-select
          v-model="form.models"
          placeholder="选择可用模型（留空表示全部）"
          multiple
          allow-clear
          allow-search
        >
          <a-option
            v-for="m in availableModels"
            :key="m"
            :value="m"
            :label="m"
          />
        </a-select>
      </a-form-item>

      <a-form-item field="subnet" label="IP 白名单">
        <a-input v-model="form.subnet" placeholder="192.168.1.0/24 多个用逗号分隔" />
      </a-form-item>

      <a-form-item field="expired_time" label="过期时间">
        <a-date-picker
          v-model="form.expired_time"
          show-time
          format="YYYY-MM-DD HH:mm:ss"
          placeholder="选择过期时间（留空永不过期）"
          style="width: 100%"
          value-format="timestamp"
        />
      </a-form-item>

      <a-form-item field="remain_quota" label="剩余额度">
        <a-input-number
          v-model="form.remain_quota"
          :min="-1"
          :precision="0"
          placeholder="500000"
          :disabled="form.unlimited_quota"
          style="width: 100%"
        />
      </a-form-item>

      <a-form-item field="unlimited_quota" label="无限额度">
        <a-checkbox v-model="form.unlimited_quota">
          无限制
        </a-checkbox>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus, IconCopy } from '@arco-design/web-vue/es/icon'
import api from '@/api'

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '名称', dataIndex: 'name', width: 160, ellipsis: true, tooltip: true },
  { title: '密钥', slotName: 'key', width: 200 },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '剩余额度', slotName: 'remain_quota', width: 130 },
  { title: '已使用', slotName: 'used_quota', width: 110 },
  { title: '过期时间', dataIndex: 'expired_time', width: 170, render: ({ record }) => formatExpiredTime(record.expired_time) },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' },
]

const loading = ref(false)
const submitting = ref(false)
const tokens = ref([])
const keyword = ref('')
const page = ref(0)
const pageSize = ref(10)
const total = ref(0)
const modalVisible = ref(false)
const isEdit = ref(false)
const editingId = ref(null)
const availableModels = ref([])

const formRef = ref(null)
const form = reactive({
  name: '',
  models: [],
  subnet: '',
  expired_time: null,
  remain_quota: 500000,
  unlimited_quota: false,
})

const rules = {
  name: [{ required: true, message: '请输入令牌名称' }],
}

const modalTitle = computed(() => isEdit.value ? '编辑令牌' : '新建令牌')

function maskKey(key) {
  if (!key) return '-'
  if (key.length <= 8) return key
  return key.substring(0, 4) + '****' + key.substring(key.length - 4)
}

function formatQuota(val) {
  if (val == null || val === '') return '-'
  const n = Number(val)
  if (isNaN(n)) return val
  if (n < 0) return '无限制'
  if (n >= 1000000) return (n / 1000000).toFixed(2) + 'M'
  if (n >= 1000) return (n / 1000).toFixed(1) + 'k'
  return String(n)
}

function formatExpiredTime(ts) {
  if (!ts) return '永不过期'
  const t = Number(ts)
  if (isNaN(t) || t <= 0) return '永不过期'
  return new Date(t * 1000).toLocaleString()
}

async function copyText(text, successMsg) {
  try {
    await navigator.clipboard.writeText(text)
    Message.success(successMsg)
  } catch {
    Message.warning('复制失败，请手动复制')
  }
}

async function fetchTokens() {
  loading.value = true
  try {
    const params = { p: page.value, size: pageSize.value }
    let url = '/api/token/'
    if (keyword.value) {
      url = '/api/token/search'
      params.keyword = keyword.value
    }
    const { data } = await api.get(url, { params })
    if (data.success) {
      tokens.value = data.data || []
      total.value = data.total || 0
    } else {
      Message.error(data.message || '加载失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function fetchAvailableModels() {
  try {
    const { data } = await api.get('/api/user/available_models')
    if (data.success) {
      availableModels.value = data.data || []
    }
  } catch {
    // ignore
  }
}

function handleSearch() {
  page.value = 0
  fetchTokens()
}

function handlePageChange(p) {
  page.value = p - 1
  fetchTokens()
}

function handlePageSizeChange(s) {
  pageSize.value = s
  page.value = 0
  fetchTokens()
}

function openCreateModal() {
  isEdit.value = false
  editingId.value = null
  form.name = ''
  form.models = []
  form.subnet = ''
  form.expired_time = null
  form.remain_quota = 500000
  form.unlimited_quota = false
  modalVisible.value = true
}

function openEditModal(record) {
  isEdit.value = true
  editingId.value = record.id
  form.name = record.name || ''
  form.models = parseModelArray(record.models)
  form.subnet = record.subnet || ''
  form.expired_time = record.expired_time ? record.expired_time * 1000 : null
  form.remain_quota = record.remain_quota ?? 500000
  form.unlimited_quota = !!record.unlimited_quota
  modalVisible.value = true
}

function parseModelArray(val) {
  if (!val) return []
  if (Array.isArray(val)) return val
  if (typeof val === 'string') {
    return val.split(',').map(s => s.trim()).filter(Boolean)
  }
  return []
}

function closeModal() {
  modalVisible.value = false
  formRef.value?.resetFields?.()
}

async function handleSubmit() {
  const valid = await formRef.value?.validate()
  if (valid !== undefined) return

  submitting.value = true
  try {
    const payload = {
      name: form.name,
      models: form.models.length ? form.models.join(',') : '',
      subnet: form.subnet,
      expired_time: form.expired_time ? Math.floor(form.expired_time / 1000) : 0,
      remain_quota: form.unlimited_quota ? -1 : form.remain_quota,
      unlimited_quota: form.unlimited_quota,
    }

    if (isEdit.value) {
      payload.id = editingId.value
      const { data } = await api.put('/api/token/', payload)
      if (data.success) {
        Message.success('令牌已更新')
        closeModal()
        fetchTokens()
      } else {
        Message.error(data.message || '更新失败')
      }
    } else {
      const { data } = await api.post('/api/token/', payload)
      if (data.success) {
        Message.success('令牌已创建')
        closeModal()
        page.value = 0
        fetchTokens()
      } else {
        Message.error(data.message || '创建失败')
      }
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

async function toggleStatus(record) {
  const newStatus = record.status === 1 ? 2 : 1
  try {
    const { data } = await api.put('/api/token/', {
      id: record.id,
      status: newStatus,
    }, { params: { status_only: true } })
    if (data.success) {
      Message.success(newStatus === 1 ? '已启用' : '已禁用')
      fetchTokens()
    } else {
      Message.error(data.message || '操作失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '操作失败')
  }
}

async function handleDelete(id) {
  try {
    const { data } = await api.delete(`/api/token/${id}/`)
    if (data.success) {
      Message.success('令牌已删除')
      fetchTokens()
    } else {
      Message.error(data.message || '删除失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '删除失败')
  }
}

onMounted(() => {
  fetchTokens()
  fetchAvailableModels()
})
</script>

<style scoped>
.table-card { border-radius: 6px; }
.action-bar { display: flex; align-items: center; gap: 12px; padding: 12px 16px; background: var(--color-fill-2); border-radius: 6px; margin-bottom: 15px; }
.action-bar-right { display: flex; align-items: center; gap: 8px; margin-left: auto; }
.table-page-footer { display: flex; justify-content: flex-end; margin-top: 20px; padding-top: 16px; border-top: 1px solid var(--color-border-2); }

.key-text {
  font-family: 'Courier New', Courier, monospace;
  font-size: 13px;
  color: var(--color-text-2);
}

.quota-unlimited {
  color: var(--color-primary-6);
  font-weight: 500;
}
</style>
