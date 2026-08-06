<template>
  <a-spin :loading="loading" class="setting-container">
    <div class="section-header">
      <h3>套餐列表</h3>
      <a-space>
        <a-button type="primary" @click="openModal()"><template #icon><icon-plus /></template>添加套餐</a-button>
      </a-space>
    </div>
    <a-table :columns="columns" :data="plans" :pagination="false" row-key="id" size="medium" :scroll="{ x: 960 }">
      <template #recommended="{ record }">
        <span v-if="record.recommended" style="color:#f5a623;font-size:16px">★</span>
        <span v-else style="color:#c9cdd4">-</span>
      </template>
      <template #status="{ record }">
        <a-tag :color="record.status===1?'green':'gray'" size="small">{{ record.status===1?'启用':'禁用' }}</a-tag>
      </template>
      <template #default_model="{ record }">{{ record.default_model || '-' }}</template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="openModal(record)">编辑</a-button>
          <a-button type="text" size="small" @click="toggleStatus(record)" :status="record.status===1?'warning':'success'">
            {{ record.status===1?'禁用':'启用' }}
          </a-button>
          <a-popconfirm content="确定删除该套餐？" @ok="handleDelete(record.id)">
            <a-button type="text" size="small" status="danger">删除</a-button>
          </a-popconfirm>
        </a-space>
      </template>
    </a-table>

    <a-modal v-model:visible="modalVisible" :title="editing ? '编辑套餐' : '添加套餐'" @ok="handleSave" :ok-loading="saving" :width="900" modal-class="arco-modal-plan-wide">
      <a-form :model="form" layout="vertical">
        <a-form-item label="名称" required>
          <a-input v-model="form.name" placeholder="请输入套餐名称" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model="form.description" :auto-size="{ minRows: 2, maxRows: 4 }" placeholder="套餐描述" />
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item label="价格"><a-input-number v-model="form.price" :style="{ width: '100%' }" :precision="2" /></a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="有效天数"><a-input-number v-model="form.duration_days" :style="{ width: '100%' }" /></a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="时长文本"><a-input v-model="form.duration_text" placeholder="如：月" /></a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item label="排序"><a-input-number v-model="form.sort" :style="{ width: '100%' }" /></a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="状态">
              <a-select v-model="form.status">
                <a-option :value="1" label="启用" />
                <a-option :value="0" label="禁用" />
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="推荐">
              <a-switch v-model="form.recommended" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="特性说明">
          <a-textarea v-model="form.features" :auto-size="{ minRows: 2, maxRows: 4 }" placeholder="每行一个特性" />
        </a-form-item>
        <a-form-item label="默认模型">
          <a-input v-model="form.default_model" placeholder="该套餐的默认模型，如：gpt-4o" />
        </a-form-item>
        <a-form-item label="模型限制 (JSON)">
          <a-textarea v-model="form.model_limits" :auto-size="{ minRows: 2, maxRows: 6 }" placeholder='{"gpt-4": 1000, "gpt-3.5": 5000}' />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-spin>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'
import api from '@/api'

const loading = ref(false)
const plans = ref([])
const modalVisible = ref(false)
const editing = ref(false)
const saving = ref(false)
const form = reactive({
  name: '', description: '', price: 0, duration_days: 30, duration_text: '',
  status: 1, recommended: false, sort: 0, features: '', model_limits: '', default_model: '',
})

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '名称', dataIndex: 'name', width: 140 },
  { title: '价格', dataIndex: 'price', width: 80 },
  { title: '有效天数', dataIndex: 'duration_days', width: 90 },
  { title: '默认模型', slotName: 'default_model', width: 120 },
  { title: '推荐', slotName: 'recommended', width: 60, align: 'center' },
  { title: '状态', slotName: 'status', width: 70 },
  { title: '排序', dataIndex: 'sort', width: 60 },
  { title: '操作', slotName: 'actions', width: 220, fixed: 'right' },
]

async function loadData() {
  loading.value = true
  try {
    const { data } = await api.get('/api/plan/', { params: { p: 0 } })
    if (data.success) plans.value = Array.isArray(data.data) ? data.data : data.data?.items || []
  } catch (e) { /* ignore */ } finally { loading.value = false }
}

function openModal(record) {
  editing.value = !!record
  if (record) {
    Object.assign(form, {
      ...record,
      recommended: record.recommended || false,
      model_limits: typeof record.model_limits === 'string' ? record.model_limits : JSON.stringify(record.model_limits || {}, null, 2),
      default_model: record.default_model || '',
      features: record.features || '',
    })
  } else {
    form.name = ''; form.description = ''; form.price = 0; form.duration_days = 30; form.duration_text = ''
    form.status = 1; form.recommended = false; form.sort = 0; form.features = ''; form.model_limits = ''; form.default_model = ''
  }
  modalVisible.value = true
}

async function handleSave() {
  saving.value = true
  try {
    const body = { ...form }
    if (editing.value) body.id = form.id
    const { data } = editing.value ? await api.put('/api/plan/', body) : await api.post('/api/plan/', body)
    if (data.success) { modalVisible.value = false; Message.success(editing.value ? '套餐已更新' : '套餐已添加'); loadData() }
    else Message.error(data.message || '操作失败')
  } catch (e) { Message.error('操作失败') } finally { saving.value = false }
}

async function handleDelete(id) { try { await api.delete(`/api/plan/${id}/`); Message.success('已删除'); loadData() } catch (e) { Message.error('删除失败') } }

async function toggleStatus(plan) {
  try {
    const { data } = await api.put('/api/plan/', { ...plan, status: plan.status === 1 ? 0 : 1 })
    if (data.success) { Message.success('状态已切换'); loadData() } else Message.error(data.message)
  } catch (e) { Message.error('操作失败') }
}

onMounted(() => { loadData() })
</script>

<style scoped>
.setting-container { padding: 4px 0; }
.section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px; }
.section-header h3 { font-size: 16px; font-weight: 600; color: var(--color-text-1); margin: 0; padding: 0; }
</style>

<style>
.arco-modal-plan-wide { width: 900px !important; }
</style>
