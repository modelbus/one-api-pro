<template>
  <a-spin :loading="loading" class="setting-container">
    <a-tabs v-model:active-key="tab" size="large">
      <a-tab-pane key="model" :title="$t('settingPage.pricing.modelTab')">
        <div class="section-header">
          <h3>{{ $t('settingPage.pricing.modelTab') }}</h3>
          <a-button type="primary" @click="openModelPrice()"><template #icon><icon-plus /></template>{{ $t('settingPage.pricing.add') }}</a-button>
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
            <a-tag :color="record.billing_type==='token'?'blue':'green'" size="small">{{ record.billing_type==='token' ? $t('settingPage.pricing.billingToken') : $t('settingPage.pricing.billingRequest') }}</a-tag>
          </template>
          <template #actions="{ record }">
            <a-space>
              <a-button type="text" size="small" @click="openModelPrice(record)">{{ $t('settingPage.pricing.edit') }}</a-button>
              <a-popconfirm :content="$t('settingPage.pricing.confirmDeleteModel')" @ok="handleDelMP(record.id)">
                <a-button type="text" size="small" status="danger">{{ $t('settingPage.pricing.delete') }}</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </a-table>
      </a-tab-pane>

      <a-tab-pane key="group" :title="$t('settingPage.pricing.groupTab')">
        <div class="section-header">
          <h3>{{ $t('settingPage.pricing.groupTab') }}</h3>
          <a-button type="primary" @click="openGroupPrice()"><template #icon><icon-plus /></template>{{ $t('settingPage.pricing.add') }}</a-button>
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
              <a-button type="text" size="small" @click="openGroupPrice(record)">{{ $t('settingPage.pricing.edit') }}</a-button>
              <a-popconfirm :content="$t('settingPage.pricing.confirmDeleteGroup')" @ok="handleDelGP(record.id)">
                <a-button type="text" size="small" status="danger">{{ $t('settingPage.pricing.delete') }}</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </a-table>
      </a-tab-pane>
    </a-tabs>

    <a-modal v-model:visible="mpVisible" :title="mpEditing ? $t('settingPage.pricing.editModelTitle') : $t('settingPage.pricing.addModelTitle')" @ok="handleSaveMP" :ok-loading="mpSaving" width="640">
      <a-form :model="mpForm" layout="vertical">
        <a-form-item :label="$t('settingPage.pricing.modelName')" required>
          <a-input v-model="mpForm.model_name" placeholder="gpt-4o" />
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item :label="$t('settingPage.pricing.inputPrice')"><a-input-number v-model="mpForm.input_price" :style="{ width: '100%' }" :precision="6" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item :label="$t('settingPage.pricing.outputPrice')"><a-input-number v-model="mpForm.output_price" :style="{ width: '100%' }" :precision="6" /></a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item :label="$t('settingPage.pricing.cachedPrice')"><a-input-number v-model="mpForm.cached_price" :style="{ width: '100%' }" :precision="6" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item :label="$t('settingPage.pricing.perRequestPrice')"><a-input-number v-model="mpForm.per_request_price" :style="{ width: '100%' }" :precision="6" /></a-form-item>
          </a-col>
        </a-row>
        <a-form-item :label="$t('settingPage.pricing.billingType')">
          <a-select v-model="mpForm.billing_type">
            <a-option value="token">{{ $t('settingPage.pricing.billingByToken') }}</a-option>
            <a-option value="per_request">{{ $t('settingPage.pricing.billingByRequest') }}</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:visible="gpVisible" :title="gpEditing ? $t('settingPage.pricing.editGroupTitle') : $t('settingPage.pricing.addGroupTitle')" @ok="handleSaveGP" :ok-loading="gpSaving" width="520">
      <a-form :model="gpForm" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item :label="$t('settingPage.pricing.groupName')" required><a-input v-model="gpForm.group_name" placeholder="vip" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item :label="$t('settingPage.pricing.modelName')" required><a-input v-model="gpForm.model_name" placeholder="gpt-4o" /></a-form-item>
          </a-col>
        </a-row>
        <a-form-item :label="$t('settingPage.pricing.discount')">
          <a-input-number v-model="gpForm.discount" :style="{ width: '100%' }" :precision="4" :min="0" :placeholder="$t('settingPage.pricing.discountPlaceholder')" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-spin>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'
import api from '@/api'

const { t } = useI18n()

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
const mpColumns = computed(() => [
  { title: t('settingPage.pricing.colModel'), dataIndex: 'model_name' },
  { title: t('settingPage.pricing.colInputPrice'), dataIndex: 'input_price', width: 110 },
  { title: t('settingPage.pricing.colOutputPrice'), dataIndex: 'output_price', width: 110 },
  { title: t('settingPage.pricing.colCachedPrice'), dataIndex: 'cached_price', width: 110 },
  { title: t('settingPage.pricing.colRequestPrice'), dataIndex: 'per_request_price', width: 110 },
  { title: t('settingPage.pricing.colBilling'), slotName: 'billing', width: 90 },
  { title: t('settingPage.pricing.colAction'), slotName: 'actions', width: 160, fixed: 'right' },
])
const mpScroll = { x: 870 }

const gpColumns = computed(() => [
  { title: t('settingPage.pricing.colGroup'), dataIndex: 'group_name' },
  { title: t('settingPage.pricing.colModel'), dataIndex: 'model_name', width: 220 },
  { title: t('settingPage.pricing.colDiscount'), dataIndex: 'discount', width: 110 },
  { title: t('settingPage.pricing.colAction'), slotName: 'actions', width: 160, fixed: 'right' },
])
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
  Object.assign(mpForm, r ? { ...r } : { id: 0, model_name: '', input_price: 0, output_price: 0, cached_price: 0, per_request_price: 0, billing_type: 'token' })
  mpVisible.value = true
}
async function handleSaveMP() {
  mpSaving.value = true
  try {
    const b = { ...mpForm }; if (mpEditing.value) b.id = mpForm.id
    const { data } = mpEditing.value ? await api.put('/api/model_price/', b) : await api.post('/api/model_price/', b)
    if (data.success) { mpVisible.value = false; loadData() } else Message.error(data.message)
  } catch (e) { Message.error(t('settingPage.pricing.opFailed')) } finally { mpSaving.value = false }
}
async function handleDelMP(id) { try { await api.delete(`/api/model_price/${id}`); loadData() } catch (e) { Message.error(t('settingPage.pricing.deleteFailed')) } }

function openGroupPrice(r) {
  gpEditing.value = !!r
  Object.assign(gpForm, r ? { ...r } : { id: 0, group_name: '', model_name: '', discount: 1 })
  gpVisible.value = true
}
async function handleSaveGP() {
  gpSaving.value = true
  try {
    const b = { ...gpForm }; if (gpEditing.value) b.id = gpForm.id
    const { data } = gpEditing.value ? await api.put('/api/group_price/', b) : await api.post('/api/group_price/', b)
    if (data.success) { gpVisible.value = false; loadData() } else Message.error(data.message)
  } catch (e) { Message.error(t('settingPage.pricing.opFailed')) } finally { gpSaving.value = false }
}
async function handleDelGP(id) { try { await api.delete(`/api/group_price/${id}`); loadData() } catch (e) { Message.error(t('settingPage.pricing.deleteFailed')) } }

onMounted(() => { loadData() })
</script>

<style scoped>
.setting-container { width: 100%; padding: 4px 0; }
.section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px; }
.section-header h3 { font-size: 16px; font-weight: 600; color: var(--color-text-1); margin: 0; padding: 0; }
</style>
