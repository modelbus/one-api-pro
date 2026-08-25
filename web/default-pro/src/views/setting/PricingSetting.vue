<template>
  <a-spin :loading="loading" class="setting-container">
    <a-tabs v-model:active-key="tab" size="large">
      <a-tab-pane key="model" title="模型定价">
        <div class="section-header">
          <h3>模型定价</h3>
          <a-button type="primary" @click="openModelPrice()"><template #icon><icon-plus /></template>添加</a-button>
        </div>
        <div class="table-wrap">
          <div class="table-body">
            <div class="table-row table-head mp-cols">
              <div class="col">模型</div>
              <div class="col col-num">输入价格</div>
              <div class="col col-num">输出价格</div>
              <div class="col col-num">缓存价格</div>
              <div class="col col-num">请求价格</div>
              <div class="col">方式</div>
              <div class="col col-action">操作</div>
            </div>
            <div v-for="r in modelPrices" :key="r.id" class="table-row mp-cols">
              <div class="col cell-strong ellipsis" :title="r.model_name">{{ r.model_name }}</div>
              <div class="col col-num cell-mono">{{ r.input_price }}</div>
              <div class="col col-num cell-mono">{{ r.output_price }}</div>
              <div class="col col-num cell-mono">{{ r.cached_price }}</div>
              <div class="col col-num cell-mono">{{ r.per_request_price }}</div>
              <div class="col">
                <a-tag :color="r.billing_type==='token'?'blue':'green'" size="small">{{ r.billing_type==='token'?'Token':'请求' }}</a-tag>
              </div>
              <div class="col col-action">
                <a-button type="text" size="small" @click="openModelPrice(r)">编辑</a-button>
                <a-popconfirm content="确定删除该定价？" @ok="handleDelMP(r.id)">
                  <a-button type="text" size="small" status="danger">删除</a-button>
                </a-popconfirm>
              </div>
            </div>
            <div v-if="modelPrices.length === 0" class="table-empty">暂无模型定价</div>
          </div>
        </div>
      </a-tab-pane>

      <a-tab-pane key="group" title="分组折扣">
        <div class="section-header">
          <h3>分组折扣</h3>
          <a-button type="primary" @click="openGroupPrice()"><template #icon><icon-plus /></template>添加</a-button>
        </div>
        <div class="table-wrap">
          <div class="table-body">
            <div class="table-row table-head gp-cols">
              <div class="col">分组</div>
              <div class="col">模型</div>
              <div class="col col-num">折扣</div>
              <div class="col col-action">操作</div>
            </div>
            <div v-for="r in groupPrices" :key="r.id" class="table-row gp-cols">
              <div class="col cell-strong ellipsis" :title="r.group_name">{{ r.group_name }}</div>
              <div class="col cell-muted ellipsis" :title="r.model_name">{{ r.model_name }}</div>
              <div class="col col-num cell-mono">{{ r.discount }}</div>
              <div class="col col-action">
                <a-button type="text" size="small" @click="openGroupPrice(r)">编辑</a-button>
                <a-popconfirm content="确定删除该分组折扣？" @ok="handleDelGP(r.id)">
                  <a-button type="text" size="small" status="danger">删除</a-button>
                </a-popconfirm>
              </div>
            </div>
            <div v-if="groupPrices.length === 0" class="table-empty">暂无分组折扣</div>
          </div>
        </div>
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
.setting-container { padding: 4px 0; }
.section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px; }
.section-header h3 { font-size: 16px; font-weight: 600; color: var(--color-text-1); margin: 0; padding: 0; }

/* 自定义表格（与渠道列表一致）：外层裁剪圆角，内层横向滚动 */
.table-wrap {
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-2);
  border-radius: 8px;
  overflow: hidden;
}
.table-body {
  overflow-x: auto;
}
.table-row {
  display: grid;
  align-items: center;
  padding: 0 16px;
  border-bottom: 1px solid var(--color-fill-3);
}
.table-row:last-child {
  border-bottom: none;
}
.mp-cols {
  grid-template-columns: minmax(200px, 1fr) 120px 120px 120px 120px 90px 180px;
  min-width: 900px;
}
.gp-cols {
  grid-template-columns: minmax(180px, 1fr) 260px 130px 180px;
  min-width: 720px;
}
.table-head {
  height: 40px;
  background: var(--color-fill-1);
  font-size: 12px;
  font-weight: 500;
  color: var(--color-text-3);
}
.table-row:not(.table-head) {
  min-height: 52px;
  transition: background 0.15s;
}
.table-row:not(.table-head):hover {
  background: var(--color-fill-1);
}

/* 单元格 */
.col {
  font-size: 13px;
  color: var(--color-text-2);
  min-width: 0;
  padding-right: 16px;
}
.col:last-child {
  padding-right: 0;
}
.col-num {
  text-align: right;
  font-variant-numeric: tabular-nums;
}
/* 操作列固定在右侧 */
.col-action {
  display: flex;
  justify-content: flex-end;
  gap: 0;
  position: sticky;
  right: 0;
  background: var(--color-bg-2);
  box-shadow: -8px 0 8px -8px rgba(0, 0, 0, 0.12);
  padding-left: 16px;
}
.col-action :deep(.arco-btn) {
  padding: 0 6px;
}
.table-head .col-action {
  background: var(--color-fill-1);
}
.table-row:not(.table-head):hover .col-action {
  background: var(--color-fill-1);
}

.cell-strong {
  color: var(--color-text-1);
  font-weight: 500;
}
.cell-muted {
  color: var(--color-text-3);
}
.cell-mono {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 12px;
  color: var(--color-text-2);
}
.ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: block;
}
.table-empty {
  padding: 40px 16px;
  text-align: center;
  color: var(--color-text-4);
  font-size: 13px;
}
</style>
