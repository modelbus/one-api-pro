<template>
  <a-spin :loading="loading" class="setting-container">
    <div class="section-header">
      <h3>{{ $t('settingPage.plan.listTitle') }}</h3>
      <a-space>
        <a-button type="primary" @click="openModal()"><template #icon><icon-plus /></template>{{ $t('settingPage.plan.addPlan') }}</a-button>
      </a-space>
    </div>
    <a-table :columns="columns" :data="plans" :pagination="false" row-key="id" size="medium" :scroll="{ x: 960 }">
      <template #recommended="{ record }">
        <span v-if="record.recommended" style="color:#f5a623;font-size:16px">★</span>
        <span v-else style="color:#c9cdd4">-</span>
      </template>
      <template #status="{ record }">
        <a-tag :color="record.status===1?'green':'gray'" size="small">{{ record.status===1 ? $t('settingPage.plan.enabled') : $t('settingPage.plan.disabled') }}</a-tag>
      </template>
      <template #default_model="{ record }">{{ record.default_model || '-' }}</template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="openModal(record)">{{ $t('settingPage.plan.edit') }}</a-button>
          <a-button type="text" size="small" @click="toggleStatus(record)" :status="record.status===1?'warning':'success'">
            {{ record.status===1 ? $t('settingPage.plan.disabled') : $t('settingPage.plan.enabled') }}
          </a-button>
          <a-popconfirm :content="$t('settingPage.plan.confirmDelete')" @ok="handleDelete(record.id)">
            <a-button type="text" size="small" status="danger">{{ $t('settingPage.plan.delete') }}</a-button>
          </a-popconfirm>
        </a-space>
      </template>
    </a-table>

    <a-modal v-model:visible="modalVisible" :title="editing ? $t('settingPage.plan.editPlan') : $t('settingPage.plan.addPlan')" @ok="handleSave" :ok-loading="saving" :width="900" modal-class="arco-modal-plan-wide">
      <a-form :model="form" layout="vertical">
        <a-form-item :label="$t('settingPage.plan.name')" required>
          <a-input v-model="form.name" :placeholder="$t('settingPage.plan.namePlaceholder')" />
        </a-form-item>
        <a-form-item :label="$t('settingPage.plan.description')">
          <a-textarea v-model="form.description" :auto-size="{ minRows: 2, maxRows: 4 }" :placeholder="$t('settingPage.plan.descriptionPlaceholder')" />
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item :label="$t('settingPage.plan.price')"><a-input-number v-model="form.price" :style="{ width: '100%' }" :precision="2" /></a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item :label="$t('settingPage.plan.durationDays')"><a-input-number v-model="form.duration_days" :style="{ width: '100%' }" /></a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item :label="$t('settingPage.plan.durationText')"><a-input v-model="form.duration_text" :placeholder="$t('settingPage.plan.durationTextPlaceholder')" /></a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item :label="$t('settingPage.plan.sort')"><a-input-number v-model="form.sort" :style="{ width: '100%' }" /></a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item :label="$t('settingPage.plan.status')">
              <a-select v-model="form.status">
                <a-option :value="1" :label="$t('settingPage.plan.enabled')" />
                <a-option :value="0" :label="$t('settingPage.plan.disabled')" />
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item :label="$t('settingPage.plan.recommended')">
              <a-switch v-model="form.recommended" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item :label="$t('settingPage.plan.features')">
          <div class="features-editor">
            <div
              v-for="(item, idx) in formFeatures"
              :key="idx"
              class="features-editor-row"
            >
              <a-input
                :model-value="item"
                :placeholder="$t('settingPage.plan.featurePlaceholder')"
                allow-clear
                @update:model-value="(v) => updateFeature(idx, v)"
                @keyup.enter="addFeatureAfter(idx)"
              />
              <a-button
                type="text"
                size="small"
                status="danger"
                class="features-editor-remove"
                :disabled="formFeatures.length === 1 && !item"
                @click="removeFeature(idx)"
              >{{ $t('settingPage.plan.delete') }}</a-button>
            </div>
            <a-button
              type="dashed"
              size="small"
              long
              class="features-editor-add"
              @click="addFeature()"
            >
              <template #icon><icon-plus :size="14" /></template>
              {{ $t('settingPage.plan.addFeature') }}
            </a-button>
          </div>
        </a-form-item>
        <a-form-item :label="$t('settingPage.plan.defaultModel')">
          <a-input v-model="form.default_model" :placeholder="$t('settingPage.plan.defaultModelPlaceholder')" />
        </a-form-item>
        <a-form-item :label="$t('settingPage.plan.modelLimits')">
          <a-textarea v-model="form.model_limits" :auto-size="{ minRows: 2, maxRows: 6 }" placeholder='{"gpt-4": 1000, "gpt-3.5": 5000}' />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-spin>
</template>

<script setup>
// 套餐设置：增删改、启用/禁用、特性说明（动态行）
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'
import api from '@/api'
import {
  sanitizeFeaturesList,
  featuresFromRecord,
  buildEmptyFeaturesForm,
} from '@/utils/plan'

const { t } = useI18n()

const loading = ref(false)
const plans = ref([])
const modalVisible = ref(false)
const editing = ref(false)
const saving = ref(false)
const form = reactive({
  name: '', description: '', price: 0, duration_days: 30, duration_text: '',
  status: 1, recommended: false, sort: 0, features: [], model_limits: '', default_model: '',
})

// 弹窗中"特性说明"动态行（与 form.features 解耦，避免双向 v-model 在行删除时的索引问题）。
//
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode
const formFeatures = ref(buildEmptyFeaturesForm())

const columns = computed(() => [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: t('settingPage.plan.name'), dataIndex: 'name', width: 140 },
  { title: t('settingPage.plan.price'), dataIndex: 'price', width: 80 },
  { title: t('settingPage.plan.durationDays'), dataIndex: 'duration_days', width: 90 },
  { title: t('settingPage.plan.defaultModel'), slotName: 'default_model', width: 120 },
  { title: t('settingPage.plan.recommended'), slotName: 'recommended', width: 60, align: 'center' },
  { title: t('settingPage.plan.status'), slotName: 'status', width: 70 },
  { title: t('settingPage.plan.sort'), dataIndex: 'sort', width: 60 },
  { title: t('settingPage.plan.action'), slotName: 'actions', width: 220, fixed: 'right' },
])

async function loadData() {
  loading.value = true
  try {
    const { data } = await api.get('/api/plan/', { params: { p: 0 } })
    if (data.success) plans.value = Array.isArray(data.data) ? data.data : data.data?.items || []
  } catch (e) { /* ignore */ } finally { loading.value = false }
}

// formFeatures 操作：动态行的增删改
//
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode
function updateFeature(idx, value) {
  if (idx < 0 || idx >= formFeatures.value.length) return
  formFeatures.value[idx] = value ?? ''
}

function addFeature() {
  formFeatures.value.push('')
}

function addFeatureAfter(idx) {
  formFeatures.value.splice(idx + 1, 0, '')
}

function removeFeature(idx) {
  if (idx < 0 || idx >= formFeatures.value.length) return
  formFeatures.value.splice(idx, 1)
  // 至少保留一行空输入框，方便继续新增
  if (formFeatures.value.length === 0) {
    formFeatures.value.push('')
  }
}

function resetForm() {
  form.name = ''
  form.description = ''
  form.price = 0
  form.duration_days = 30
  form.duration_text = ''
  form.status = 1
  form.recommended = false
  form.sort = 0
  form.features = []
  form.model_limits = ''
  form.default_model = ''
}

function openModal(record) {
  editing.value = !!record
  resetForm()
  if (record) {
    Object.assign(form, {
      ...record,
      recommended: record.recommended || false,
      model_limits: typeof record.model_limits === 'string'
        ? record.model_limits
        : JSON.stringify(record.model_limits || {}, null, 2),
      default_model: record.default_model || '',
      features: Array.isArray(record.features) ? record.features : [],
    })
  }
  formFeatures.value = featuresFromRecord(form.features)
  modalVisible.value = true
}

async function handleSave() {
  saving.value = true
  try {
    const cleanedFeatures = sanitizeFeaturesList(formFeatures.value)
    const body = {
      ...form,
      features: cleanedFeatures,
    }
    if (editing.value) body.id = form.id
    const { data } = editing.value ? await api.put('/api/plan/', body) : await api.post('/api/plan/', body)
    if (data.success) { modalVisible.value = false; Message.success(editing.value ? t('settingPage.plan.updated') : t('settingPage.plan.added')); loadData() }
    else Message.error(data.message || t('settingPage.plan.opFailed'))
  } catch (e) { Message.error(t('settingPage.plan.opFailed')) } finally { saving.value = false }
}

async function handleDelete(id) { try { await api.delete(`/api/plan/${id}/`); Message.success(t('settingPage.plan.deleted')); loadData() } catch (e) { Message.error(t('settingPage.plan.deleteFailed')) } }

async function toggleStatus(plan) {
  try {
    const { data } = await api.put('/api/plan/', { ...plan, status: plan.status === 1 ? 0 : 1 })
    if (data.success) { Message.success(t('settingPage.plan.statusToggled')); loadData() } else Message.error(data.message)
  } catch (e) { Message.error(t('settingPage.plan.opFailed')) }
}

onMounted(() => { loadData() })
</script>

<style scoped>
.setting-container { padding: 4px 0; }
.section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px; }
.section-header h3 { font-size: 16px; font-weight: 600; color: var(--color-text-1); margin: 0; padding: 0; }
.features-editor { display: flex; flex-direction: column; gap: 8px; }
.features-editor-row { display: flex; align-items: center; gap: 8px; }
.features-editor-row > :first-child { flex: 1; min-width: 0; }
.features-editor-remove { flex-shrink: 0; }
.features-editor-add { margin-top: 4px; }
</style>

<style>
.arco-modal-plan-wide { width: 900px !important; }
</style>
