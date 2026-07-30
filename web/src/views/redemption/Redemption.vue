<template>
  <a-card :bordered="false" class="table-card">
    <div class="action-bar">
      <a-input-search
        v-model="keyword"
        :placeholder="$t('common.search') + '...'"
        :style="{ width: '220px' }"
        @search="handleSearch"
        @clear="fetchData"
        allow-clear
      />
      <div class="action-bar-right">
        <a-space>
          <a-button type="primary" @click="openCreateModal">
            {{ $t('redemption.addRedemption') }}
          </a-button>
        </a-space>
      </div>
    </div>

    <a-table
      :columns="columns"
      :data="redemptions"
      :loading="loading"
      :bordered="{ wrapper: true, cell: false }"
      size="medium"
      :scroll="{ x: 800 }"
      :pagination="false"
      row-key="id"
    >
      <template #status="{ record }">
        <a-tag :color="statusColor(record.status)">
          {{ statusText(record.status) }}
        </a-tag>
      </template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="openEditModal(record)">
            {{ $t('common.edit') }}
          </a-button>
          <a-popconfirm
            content="确定要删除该兑换码吗？"
            @ok="handleDelete(record.id)"
          >
            <a-button type="text" size="small" status="danger">
              {{ $t('common.delete') }}
            </a-button>
          </a-popconfirm>
        </a-space>
      </template>
      <template #empty>
        <a-empty :description="$t('common.noData')" />
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
    v-model:visible="modalVisible"
    :title="editingRecord ? $t('redemption.editRedemption') : $t('redemption.addRedemption')"
    @ok="handleSubmit"
    @cancel="modalVisible = false"
    :ok-loading="submitting"
  >
    <a-form ref="formRef" :model="form" layout="vertical">
      <a-form-item
        field="name"
        :label="$t('redemption.name')"
        :rules="[{ required: true, message: '请输入名称' }]"
      >
        <a-input v-model="form.name" placeholder="兑换码名称" />
      </a-form-item>
      <a-form-item
        v-if="!editingRecord"
        field="count"
        :label="$t('redemption.count')"
        :rules="[{ required: true, message: '请输入数量' }]"
      >
        <a-input-number
          v-model="form.count"
          :min="1"
          placeholder="生成数量"
          :style="{ width: '100%' }"
        />
      </a-form-item>
      <a-form-item
        field="quota"
        :label="$t('redemption.quota')"
        :rules="[{ required: true, message: '请输入额度' }]"
      >
        <a-input-number
          v-model="form.quota"
          :min="0"
          :precision="2"
          placeholder="额度"
          :style="{ width: '100%' }"
        />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { toastSuccess, toastError, toastWarning } from '@/helpers/toast'
import { ref, reactive, onMounted } from 'vue'
import api from '@/api'

const redemptions = ref([])
const loading = ref(false)
const keyword = ref('')
const modalVisible = ref(false)
const submitting = ref(false)
const editingRecord = ref(null)
const formRef = ref(null)
const form = reactive({
  name: '',
  count: 1,
  quota: 0,
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showTotal: true,
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '名称', dataIndex: 'name', ellipsis: true },
  { title: '密钥', dataIndex: 'key', ellipsis: true },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '额度', dataIndex: 'quota', width: 100 },
  { title: '操作', slotName: 'actions', width: 160, fixed: 'right' },
]

const statusMap = {
  1: { text: '已启用', color: 'green' },
  2: { text: '已禁用', color: 'gray' },
  3: { text: '已过期', color: 'red' },
}

const statusColor = (status) => {
  return statusMap[status]?.color || 'gray'
}

const statusText = (status) => {
  return statusMap[status]?.text || status
}

const fetchData = async () => {
  loading.value = true
  try {
    const page = pagination.current - 1
    const url = keyword.value
      ? `/api/redemption/search?keyword=${encodeURIComponent(keyword.value)}`
      : `/api/redemption/?p=${page}`
    const res = await api.get(url)
    if (res.data) {
      redemptions.value = res.data.data ?? res.data.results ?? res.data
      pagination.total = res.data.total ?? res.data.length ?? 0
    }
  } catch (err) {
    import('@arco-design/web-vue').then(m => m.Message.error)(err?.response?.data?.message || '获取兑换码列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
  fetchData()
}

const handlePageChange = (page) => {
  pagination.current = page
  fetchData()
}

const openCreateModal = () => {
  editingRecord.value = null
  form.name = ''
  form.count = 1
  form.quota = 0
  modalVisible.value = true
}

const openEditModal = (record) => {
  editingRecord.value = record
  form.name = record.name
  form.quota = record.quota
  modalVisible.value = true
}

const handleSubmit = async () => {
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
      import('@arco-design/web-vue').then(m => m.Message.success)('兑换码已更新')
    } else {
      await api.post('/api/redemption/', {
        name: form.name,
        count: form.count,
        quota: form.quota,
      })
      import('@arco-design/web-vue').then(m => m.Message.success)('兑换码已生成')
    }
    modalVisible.value = false
    fetchData()
  } catch (err) {
    import('@arco-design/web-vue').then(m => m.Message.error)(err?.response?.data?.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (id) => {
  try {
    await api.delete(`/api/redemption/${id}/`)
    import('@arco-design/web-vue').then(m => m.Message.success)('兑换码已删除')
    fetchData()
  } catch (err) {
    import('@arco-design/web-vue').then(m => m.Message.error)(err?.response?.data?.message || '删除兑换码失败')
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
</style>
