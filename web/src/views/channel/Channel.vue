<template>
  <a-card :bordered="false" class="table-card">
    <div class="action-bar">
      <a-input-search
        v-model="searchKeyword"
        placeholder="搜索渠道名称..."
        :style="{ width: '220px' }"
        @search="handleSearch"
        @clear="handleSearchClear"
        allow-clear
      />
      <div class="action-bar-right">
        <a-button type="primary" @click="openCreateModal">
          <template #icon><icon-plus /></template>
          添加渠道
        </a-button>
        <a-button @click="handleTestAll" :loading="testingAll">测试全部</a-button>
        <a-button @click="handleUpdateBalance" :loading="updatingBalance">更新余额</a-button>
        <a-popconfirm content="确认删除所有已禁用的渠道？" @ok="handleDeleteDisabled">
          <a-button status="danger" :loading="deletingDisabled">删除已禁用</a-button>
        </a-popconfirm>
      </div>
    </div>

    <a-spin :loading="loading">
      <a-table
        :data="channels"
        :columns="columns"
        :pagination="false"
        row-key="id"
        :bordered="{ wrapper: true, cell: false }"
        size="medium"
        :scroll="{ x: 1100 }"
      >
        <template #type="{ record }">
          <a-tag :color="typeColorMap[record.type] || 'gray'" size="small">
            {{ typeNameMap[record.type] || `Type ${record.type}` }}
          </a-tag>
        </template>
        <template #status="{ record }">
          <a-tag :color="record.status === 1 ? 'green' : 'red'" size="small">
            {{ record.status === 1 ? '已启用' : '已禁用' }}
          </a-tag>
        </template>
        <template #response_time="{ record }">
          {{ record.response_time ? record.response_time + 'ms' : '-' }}
        </template>
        <template #actions="{ record }">
          <a-space>
            <a-button type="text" size="small" @click="openEditModal(record)">编辑</a-button>
            <a-popconfirm content="确定删除此渠道？" @ok="handleDelete(record.id)">
              <a-button type="text" size="small" status="danger">删除</a-button>
            </a-popconfirm>
            <a-button type="text" size="small" @click="handleTest(record)" :loading="testingIds.includes(record.id)">测试</a-button>
          </a-space>
        </template>
      </a-table>
    </a-spin>

    <div class="table-page-footer">
      <a-pagination
        :current="currentPage"
        :total="total"
        :page-size="pageSize"
        show-total
        show-page-size
        size="medium"
        @change="handlePageChange"
        @page-size-change="handlePageSizeChange"
      />
    </div>
  </a-card>

  <a-modal
    v-model:visible="modalVisible"
    :title="modalTitle"
    :width="560"
    @ok="handleSubmit"
    @cancel="closeModal"
    :ok-loading="submitting"
  >
    <a-form ref="formRef" :model="form" layout="vertical">
      <a-form-item field="type" label="类型" required>
        <a-select v-model="form.type" placeholder="选择渠道类型">
          <a-option v-for="(name, key) in typeNameMap" :key="key" :value="Number(key)" :label="name" />
        </a-select>
      </a-form-item>
      <a-form-item field="name" label="名称" required>
        <a-input v-model="form.name" placeholder="渠道名称" />
      </a-form-item>
      <a-form-item field="base_url" label="Base URL">
        <a-input v-model="form.base_url" placeholder="https://api.openai.com" />
      </a-form-item>
      <a-form-item field="models" label="模型">
        <a-input v-model="form.models" placeholder="多个模型用逗号分隔" />
      </a-form-item>
      <a-form-item field="groups" label="分组">
        <a-input v-model="form.groups" placeholder="多个分组用逗号分隔" />
      </a-form-item>
      <a-form-item field="key" label="密钥" required>
        <a-input-password v-model="form.key" placeholder="API Key" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'
import api from '@/api'

const typeNameMap = {
  1: 'OpenAI', 2: 'Claude', 3: 'Azure', 4: 'Gemini',
  5: 'Baidu', 6: 'Aliyun', 7: 'Tencent', 8: 'Xunfei',
  9: 'Zhipu', 10: 'DeepSeek', 11: 'Midjourney',
}
const typeColorMap = {
  1: 'blue', 2: 'orange', 3: 'cyan', 4: 'green',
  5: 'red', 6: 'purple', 7: 'magenta', 8: 'gold',
  9: 'lime', 10: 'arcoblue', 11: 'pinkpurple',
}

const columns = [
  { title: '类型', slotName: 'type', width: 100 },
  { title: '名称', dataIndex: 'name', width: 160, ellipsis: true, tooltip: true },
  { title: 'Base URL', dataIndex: 'base_url', width: 200, ellipsis: true, tooltip: true },
  { title: '模型', dataIndex: 'models', width: 180, ellipsis: true, tooltip: true },
  { title: '分组', dataIndex: 'groups', width: 120, ellipsis: true, tooltip: true },
  { title: '状态', slotName: 'status', width: 90 },
  { title: '响应时间', slotName: 'response_time', width: 100 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' },
]

const loading = ref(false)
const submitting = ref(false)
const channels = ref([])
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const modalVisible = ref(false)
const isEdit = ref(false)
const editingId = ref(null)
const testingAll = ref(false)
const updatingBalance = ref(false)
const deletingDisabled = ref(false)
const testingIds = ref([])

const formRef = ref(null)
const form = reactive({
  type: 1,
  name: '',
  base_url: '',
  models: '',
  groups: '',
  key: '',
})

const modalTitle = ref('添加渠道')

async function fetchChannels() {
  loading.value = true
  try {
    const params = { p: currentPage.value - 1, size: pageSize.value }
    const { data } = await api.get('/api/channel/', { params })
    if (data.success) {
      channels.value = data.data || []
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

function handleSearch() {
  currentPage.value = 1
  if (!searchKeyword.value) { fetchChannels(); return }
  loading.value = true
  api.get('/api/channel/search', { params: { keyword: searchKeyword.value } })
    .then(({ data }) => {
      if (data.success) { channels.value = data.data || []; total.value = data.total || 0 }
    })
    .catch(e => Message.error(e.response?.data?.message || e.message || '搜索失败'))
    .finally(() => { loading.value = false })
}

function handleSearchClear() {
  searchKeyword.value = ''
  currentPage.value = 1
  fetchChannels()
}

function handlePageChange(p) {
  currentPage.value = p
  fetchChannels()
}

function handlePageSizeChange(s) {
  pageSize.value = s
  currentPage.value = 1
  fetchChannels()
}

function openCreateModal() {
  isEdit.value = false
  editingId.value = null
  form.type = 1
  form.name = ''
  form.base_url = ''
  form.models = ''
  form.groups = ''
  form.key = ''
  modalTitle.value = '添加渠道'
  modalVisible.value = true
}

function openEditModal(record) {
  isEdit.value = true
  editingId.value = record.id
  form.type = record.type || 1
  form.name = record.name || ''
  form.base_url = record.base_url || ''
  form.models = record.models || ''
  form.groups = record.groups || ''
  form.key = record.key || ''
  modalTitle.value = '编辑渠道'
  modalVisible.value = true
}

function closeModal() {
  modalVisible.value = false
  formRef.value?.clearValidate()
}

async function handleSubmit() {
  const errors = await formRef.value?.validate()
  if (errors) return
  submitting.value = true
  try {
    const payload = {
      type: form.type,
      name: form.name,
      base_url: form.base_url,
      models: form.models,
      groups: form.groups,
      key: form.key,
    }
    let res
    if (isEdit.value) {
      payload.id = editingId.value
      res = await api.put('/api/channel/', payload)
    } else {
      res = await api.post('/api/channel/', payload)
    }
    if (res.data.success) {
      Message.success(isEdit.value ? '渠道已更新' : '渠道已添加')
      closeModal()
      fetchChannels()
    } else {
      Message.error(res.data.message || '操作失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

async function handleDelete(id) {
  try {
    const { data } = await api.delete(`/api/channel/${id}/`)
    if (data.success) {
      Message.success('渠道已删除')
      fetchChannels()
    } else {
      Message.error(data.message || '删除失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '删除失败')
  }
}

async function handleTest(record) {
  testingIds.value.push(record.id)
  try {
    const { data } = await api.post('/api/channel/test', { id: record.id })
    if (data.success) {
      Message.success('测试通过')
      fetchChannels()
    } else {
      Message.error(data.message || '测试失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '测试失败')
  } finally {
    testingIds.value = testingIds.value.filter(id => id !== record.id)
  }
}

async function handleTestAll() {
  testingAll.value = true
  try {
    const { data } = await api.post('/api/channel/test_all')
    if (data.success) {
      Message.success('全部测试完成')
      fetchChannels()
    } else {
      Message.error(data.message || '批量测试失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '批量测试失败')
  } finally {
    testingAll.value = false
  }
}

async function handleUpdateBalance() {
  updatingBalance.value = true
  try {
    const { data } = await api.put('/api/channel/update_balance')
    if (data.success) {
      Message.success('余额已更新')
      fetchChannels()
    } else {
      Message.error(data.message || '更新余额失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '更新余额失败')
  } finally {
    updatingBalance.value = false
  }
}

async function handleDeleteDisabled() {
  deletingDisabled.value = true
  try {
    const { data } = await api.delete('/api/channel/disabled')
    if (data.success) {
      Message.success('已删除所有禁用的渠道')
      fetchChannels()
    } else {
      Message.error(data.message || '删除失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || '删除失败')
  } finally {
    deletingDisabled.value = false
  }
}

onMounted(() => {
  fetchChannels()
})
</script>

<style scoped>
.table-card { border-radius: 6px; }
.action-bar { display: flex; align-items: center; gap: 12px; padding: 12px 16px; background: var(--color-fill-2); border-radius: 6px; margin-bottom: 15px; }
.action-bar-right { display: flex; align-items: center; gap: 8px; margin-left: auto; }
.table-page-footer { display: flex; justify-content: flex-end; margin-top: 20px; padding-top: 16px; border-top: 1px solid var(--color-border-2); }
</style>
