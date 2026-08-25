<template>
  <a-spin :loading="loading" class="setting-container">
    <a-tabs v-model:active-key="tab" size="large">
      <a-tab-pane key="model" title="模型定价">
        <div class="section-header">
          <h3>模型定价</h3>
          <a-button type="primary" @click="openModelPrice()"><template #icon><icon-plus /></template>添加</a-button>
        </div>
        <a-table
          row-key="id"
          :columns="mpColumns"
          :data="modelPrices"
          :pagination="false"
          size="medium"
          :scroll="mpScroll"
        >
          <template #billing="{ record }">
            <a-tag :color="record.billing_type==='token'?'blue':'green'" size="small">{{ record.billing_type==='token'?'Token':'请求' }}</a-tag>
          </template>
          <template #actions="{ record }">
            <a-space>
              <a-button type="text" size="small" @click="openModelPrice(record)">编辑</a-button>
              <a-popconfirm content="确定删除该定价？" @ok="handleDelMP(record.id)">
                <a-button type="text" size="small" status="danger">删除</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </a-table>
      </a-tab-pane>

      <a-tab-pane key="group" title="分组折扣">
        <div class="section-header">
          <h3>分组折扣</h3>
          <a-button type="primary" @click="openGroupPrice()"><template #icon><icon-plus /></template>添加</a-button>
        </div>
        <a-table
          row-key="id"
          :columns="gpColumns"
          :data="groupPrices"
          :pagination="false"
          size="medium"
          :scroll="gpScroll"
        >
          <template #actions="{ record }">
            <a-space>
              <a-button type="text" size="small" @click="openGroupPrice(record)">编辑</a-button>
              <a-popconfirm content="确定删除该分组折扣？" @ok="handleDelGP(record.id)">
                <a-button type="text" size="small" status="danger">删除</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </a-table>
      </a-tab-pane>
    </a-tabs>

    <a-modal v-model:visible="mpVisible" :title="mpEditing?'编辑模型定价':'添加模型定价'" @ok="handleSaveMP" :ok-loading="mpSaving" width="640">
      <a-form :model="mpForm" layout="vertical">
        <a-form-item label="模型名称" required>
          <a-input v-model="mpForm.model_name" placeholder="gpt-4o" />
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="输入价格 (每 1M Tokens)"><a-input-number v-model="mpForm.input_price" :style="{ width: '100%' }" :precision="6" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="输出价格 (每 1M Tokens)"><a-input-number v-model="mpForm.output_price" :style="{ width: '100%' }" :precision="6" /></a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="缓存价格 (每 1M Tokens)"><a-input-number v-model="mpForm.cached_price" :style="{ width: '100%' }" :precision="6" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="每次请求价格"><a-input-number v-model="mpForm.per_request_price" :style="{ width: '100%' }" :precision="6" /></a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="计费方式">
          <a-select v-model="mpForm.billing_type">
            <a-option value="token">按 Token</a-option>
            <a-option value="per_request">按请求</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:visible="gpVisible" :title="gpEditing?'编辑分组折扣':'添加分组折扣'" @ok="handleSaveGP" :ok-loading="gpSaving" width="520">
      <a-form :model="gpForm" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="分组名称" required><a-input v-model="gpForm.group_name" placeholder="vip" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="模型名称" required><a-input v-model="gpForm.model_name" placeholder="gpt-4o" /></a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="折扣倍数">
          <a-input-number v-model="gpForm.discount" :style="{ width: '100%' }" :precision="4" :min="0" placeholder="1.0(不打折)" />
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

const tab = ref('model')
const loading = ref(false)
const modelPrices = ref([])
const groupPrices = ref([])

// Columns + scroll per arco-design table#scroll convention:
//  - the first column (模型 / 分组) has NO width so it absorbs leftover
//    horizontal space and the table element stretches to fill the container
//  - the other columns keep fixed widths for a predictable layout
//  - scroll.x is a numeric minimum table width (no scroll.minWidth, which
//    would override the component's CSS min-width:100% and pin the table
//    to a fixed pixel width, producing a blank trailing column)
//  - fixed: 'right' keeps the 操作 column pinned while scrolling
const mpColumns = [
  { title: '模型', dataIndex: 'model_name' },
  { title: '输入价格', dataIndex: 'input_price', width: 110 },
  { title: '输出价格', dataIndex: 'output_price', width: 110 },
  { title: '缓存价格', dataIndex: 'cached_price', width: 110 },
  { title: '请求价格', dataIndex: 'per_request_price', width: 110 },
  { title: '方式', slotName: 'billing', width: 90 },
  { title: '操作', slotName: 'actions', width: 160, fixed: 'right' },
]
const mpScroll = { x: 870 }

const gpColumns = [
  { title: '分组', dataIndex: 'group_name' },
  { title: '模型', dataIndex: 'model_name', width: 220 },
  { title: '折扣', dataIndex: 'discount', width: 110 },
  { title: '操作', slotName: 'actions', width: 160, fixed: 'right' },
]
const gpScroll = { x: 650 }

const mpVisible = ref(false), mpEditing = ref(false), mpSaving = ref(false)
const mpForm = reactive({ model_name: '', input_price: 0, output_price: 0, cached_price: 0, per_request_price: 0, billing_type: 'token' })

const gpVisible = ref(false), gpEditing = ref(false), gpSaving = ref(false)
const gpForm = reactive({ group_name: '', model_name: '', discount: 1 })

async function loadData() {
  loading.value = true
  try {
    const [mp, gp] = await Promise.all([api.get('/api/model_price/'), api.get('/api/group_price/')])
    modelPrices.value = Array.isArray(mp.data.data) ? mp.data.data : mp.data.data?.items || []
    groupPrices.value = Array.isArray(gp.data.data) ? gp.data.data : gp.data.data?.items || []
  } catch (e) { /* ignore */ } finally { loading.value = false }
}

function openModelPrice(r) {
  mpEditing.value = !!r
  Object.assign(mpForm, r ? { ...r } : { model_name: '', input_price: 0, output_price: 0, cached_price: 0, per_request_price: 0, billing_type: 'token' })
  mpVisible.value = true
}
async function handleSaveMP() {
  mpSaving.value = true
  try {
    const b = { ...mpForm }; if (mpEditing.value) b.id = mpForm.id
    const { data } = mpEditing.value ? await api.put('/api/model_price/', b) : await api.post('/api/model_price/', b)
    if (data.success) { mpVisible.value = false; loadData() } else Message.error(data.message)
  } catch (e) { Message.error('操作失败') } finally { mpSaving.value = false }
}
async function handleDelMP(id) { try { await api.delete(`/api/model_price/${id}`); loadData() } catch (e) { Message.error('删除失败') } }

function openGroupPrice(r) {
  gpEditing.value = !!r
  Object.assign(gpForm, r ? { ...r } : { group_name: '', model_name: '', discount: 1 })
  gpVisible.value = true
}
async function handleSaveGP() {
  gpSaving.value = true
  try {
    const b = { ...gpForm }; if (gpEditing.value) b.id = gpForm.id
    const { data } = gpEditing.value ? await api.put('/api/group_price/', b) : await api.post('/api/group_price/', b)
    if (data.success) { gpVisible.value = false; loadData() } else Message.error(data.message)
  } catch (e) { Message.error('操作失败') } finally { gpSaving.value = false }
}
async function handleDelGP(id) { try { await api.delete(`/api/group_price/${id}`); loadData() } catch (e) { Message.error('删除失败') } }

onMounted(() => { loadData() })
</script>

<style scoped>
.setting-container { width: 100%; padding: 4px 0; }
.section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px; }
.section-header h3 { font-size: 16px; font-weight: 600; color: var(--color-text-1); margin: 0; padding: 0; }
</style>
