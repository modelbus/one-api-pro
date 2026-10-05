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
      <template #billingType="{ record }">
        <a-tag :color="record.billing_type==='request'?'green':'blue'" size="small">
          {{ record.billing_type==='request' ? $t('settingPage.plan.billingRequest') : $t('settingPage.plan.billingToken') }}
        </a-tag>
      </template>
      <template #virtualAmount="{ record }">
        <span v-if="Number(record.virtual_amount) > 0">¥{{ formatAmount(record.virtual_amount) }}</span>
        <span v-else class="text-muted">{{ $t('settingPage.plan.virtualAmountUnlimited') }}</span>
      </template>
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

    <a-modal v-model:visible="modalVisible" :title="editing ? $t('settingPage.plan.editPlan') : $t('settingPage.plan.addPlan')" @ok="handleSave" :ok-loading="saving" :width="760" modal-class="arco-modal-plan-wide">
      <a-form
        :model="form"
        layout="horizontal"
        :label-col="{ flex: '100px' }"
        :wrapper-col="{ flex: '1' }"
      >
        <div class="form-group-title">{{ $t('settingPage.plan.groupBasic') }}</div>
        <a-form-item :label="$t('settingPage.plan.name')" required>
          <div class="form-field">
            <a-input v-model="form.name" :placeholder="$t('settingPage.plan.namePlaceholder')" />
          </div>
        </a-form-item>
        <a-form-item :label="$t('settingPage.plan.description')">
          <div class="form-field">
            <a-textarea v-model="form.description" :auto-size="{ minRows: 2, maxRows: 4 }" :placeholder="$t('settingPage.plan.descriptionPlaceholder')" />
          </div>
        </a-form-item>
        <a-form-item :label="$t('settingPage.plan.price')">
          <div class="form-field">
            <a-input-number v-model="form.price" :precision="2" />
          </div>
        </a-form-item>
        <a-form-item :label="$t('settingPage.plan.durationDays')">
          <div class="form-field">
            <a-input-group class="duration-group">
              <a-input-number v-model="form.duration_days" :min="1" />
              <a-input v-model="form.duration_text" :placeholder="$t('settingPage.plan.durationTextPlaceholder')" />
            </a-input-group>
          </div>
        </a-form-item>
        <a-form-item :label="$t('settingPage.plan.sort')">
          <div class="form-field">
            <a-input-number v-model="form.sort" />
          </div>
        </a-form-item>
        <a-form-item :label="$t('settingPage.plan.status')">
          <div class="form-field">
            <a-switch
              type="round"
              :model-value="form.status === 1"
              :checked-text="$t('settingPage.plan.enabled')"
              :unchecked-text="$t('settingPage.plan.disabled')"
              @update:model-value="(v) => (form.status = v ? 1 : 0)"
            />
          </div>
        </a-form-item>
        <a-form-item :label="$t('settingPage.plan.recommended')">
          <div class="form-field">
            <a-switch
              type="round"
              v-model="form.recommended"
              :checked-text="$t('settingPage.plan.recommendedOn')"
              :unchecked-text="$t('settingPage.plan.recommendedOff')"
            />
          </div>
        </a-form-item>
        <a-form-item :label="$t('settingPage.plan.features')">
          <div class="form-field">
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
                >
                  <template #icon><icon-delete /></template>
                </a-button>
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
          </div>
        </a-form-item>

        <div class="form-group-title">{{ $t('settingPage.plan.groupBilling') }}</div>
        <a-form-item :label="$t('settingPage.plan.billingType')">
          <div class="form-field">
            <a-radio-group v-model="form.billing_type">
              <a-radio :value="BILLING_TOKEN">{{ $t('settingPage.plan.billingToken') }}</a-radio>
              <a-radio :value="BILLING_REQUEST">{{ $t('settingPage.plan.billingRequest') }}</a-radio>
            </a-radio-group>
          </div>
        </a-form-item>
        <a-form-item :label="$t('settingPage.plan.virtualAmount')" :extra="$t('settingPage.plan.virtualAmountHint')">
          <div class="form-field">
            <a-input-number v-model="form.virtual_amount" :precision="2" :min="0" />
          </div>
        </a-form-item>
        <a-form-item
          :label="$t('settingPage.plan.modelLimits')"
          :extra="limitsUnitHint"
          required
        >
          <div class="limits-editor">
            <div class="limits-header">
              <span>{{ $t('settingPage.plan.modelLimitsColModel') }}</span>
              <span>{{ $t('settingPage.plan.modelLimitsColPeriodH') }}</span>
              <span>{{ $t('settingPage.plan.modelLimitsColLimitPeriod') }}</span>
              <span>{{ $t('settingPage.plan.modelLimitsColLimitWeek') }}</span>
              <span>{{ $t('settingPage.plan.modelLimitsColLimitMonth') }}</span>
              <span />
            </div>
            <div v-for="(row, idx) in limitRows" :key="idx" class="limits-row">
              <a-select
                v-model="row.model"
                :placeholder="$t('settingPage.plan.modelLimitsModelPlaceholder')"
                allow-search
                allow-clear
              >
                <a-option v-for="m in modelOptions" :key="m" :value="m" :label="m" />
              </a-select>
              <a-select v-model="row.periodH">
                <a-option
                  v-for="h in periodHOptions"
                  :key="h"
                  :value="h"
                  :label="$t('settingPage.plan.modelLimitsPeriodHOption', { n: h })"
                />
              </a-select>
              <a-input-number v-model="row.limitPeriod" :min="0" :placeholder="$t('settingPage.plan.modelLimitsLimitPlaceholder')" />
              <a-input-number v-model="row.limitWeek" :min="0" :placeholder="$t('settingPage.plan.modelLimitsLimitPlaceholder')" />
              <a-input-number v-model="row.limitMonth" :min="0" :placeholder="$t('settingPage.plan.modelLimitsLimitPlaceholder')" />
              <a-button
                type="text"
                size="small"
                status="danger"
                class="limits-row-remove"
                @click="removeLimitRow(idx)"
              >
                <template #icon><icon-delete /></template>
              </a-button>
            </div>
            <div v-if="!limitRows.length" class="limits-empty">{{ $t('settingPage.plan.modelLimitsEmpty') }}</div>
            <a-button type="dashed" size="small" long class="limits-add" @click="addLimitRow()">
              <template #icon><icon-plus :size="14" /></template>
              {{ $t('settingPage.plan.modelLimitsAdd') }}
            </a-button>
          </div>
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
import { IconPlus, IconDelete } from '@arco-design/web-vue/es/icon'
import api from '@/api'
import {
  sanitizeFeaturesList,
  featuresFromRecord,
  buildEmptyFeaturesForm,
} from '@/utils/plan'
import {
  BILLING_TOKEN,
  BILLING_REQUEST,
  MAX_PERIOD_H,
  MIN_PERIOD_H,
  createLimitRow,
  modelLimitsToRows,
  rowsToModelLimits,
  validateLimitRows,
} from '@/utils/plan_limits'

const { t } = useI18n()

// formatAmount 金额显示：保留两位小数。
// 套餐额度由后端以「元」返回，前端不做任何 1e6 换算。
//
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
function formatAmount(v) {
  const n = Number(v)
  if (!Number.isFinite(n)) return '0.00'
  return n.toFixed(2)
}

const loading = ref(false)
const plans = ref([])
const modalVisible = ref(false)
const editing = ref(false)
const saving = ref(false)
const form = reactive({
  name: '', description: '', price: 0, duration_days: 30, duration_text: '',
  status: 1, recommended: false, sort: 0, features: [],
  billing_type: BILLING_TOKEN, virtual_amount: 0, model_limits: '',
})

// 弹窗中"特性说明"动态行（与 form.features 解耦，避免双向 v-model 在行删除时的索引问题）。
//
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode
const formFeatures = ref(buildEmptyFeaturesForm())

// 模型限制在表单内部是「限制行」列表，每行 = 一个模型 + 窗口小时数 + 三个窗口限额
// （窗口期 / 周 / 月），提交前才由 rowsToModelLimits 合成后端 model_limits JSON。
// 限额字段名随计费维度切换：token_period/week/month 或 request_period/week/month。
//
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
const limitRows = ref([])
const modelOptions = ref([])

// periodHOptions 窗口小时数下拉选项（1-24 小时）。
const periodHOptions = computed(() => Array.from({ length: MAX_PERIOD_H }, (_, i) => i + MIN_PERIOD_H))

// limitsUnitHint 提示三个限额列的数值含义（随计费维度变化）。
const limitsUnitHint = computed(() =>
  form.billing_type === BILLING_REQUEST
    ? t('settingPage.plan.modelLimitsUnitHintRequest')
    : t('settingPage.plan.modelLimitsUnitHintToken'),
)

// limitErrorMessage 把 validateLimitRows 返回的错误码翻译成本地化提示。
//
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
function limitErrorMessage(err) {
  return t(`settingPage.plan.${err?.code || 'errLimitsEmpty'}`, { ...(err?.params || {}) })
}

// addLimitRow 追加一行空的模型限制。
//
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
function addLimitRow() {
  limitRows.value.push(createLimitRow())
}

// removeLimitRow 删除指定下标（从 0 开始）的限制行。
//
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
function removeLimitRow(idx) {
  if (idx < 0 || idx >= limitRows.value.length) return
  limitRows.value.splice(idx, 1)
}

// fetchModelOptions 拉取可选模型列表；与渠道编辑弹窗使用同一个接口。
//
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
async function fetchModelOptions() {
  try {
    const { data } = await api.get('/api/model_price/options')
    const list = Array.isArray(data?.data) ? data.data : []
    modelOptions.value = Array.from(new Set(list.filter((m) => typeof m === 'string' && m))).sort()
  } catch (e) { /* ignore */ }
}

const columns = computed(() => [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: t('settingPage.plan.name'), dataIndex: 'name', width: 140 },
  { title: t('settingPage.plan.price'), dataIndex: 'price', width: 80 },
  { title: t('settingPage.plan.billingType'), slotName: 'billingType', width: 100 },
  { title: t('settingPage.plan.virtualAmount'), slotName: 'virtualAmount', width: 110 },
  { title: t('settingPage.plan.durationDays'), dataIndex: 'duration_days', width: 90 },
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
  form.billing_type = BILLING_TOKEN
  form.virtual_amount = 0
  form.model_limits = ''
  limitRows.value = []
}

function openModal(record) {
  editing.value = !!record
  resetForm()
  let rawLimits = ''
  if (record) {
    Object.assign(form, {
      ...record,
      recommended: record.recommended || false,
      billing_type: record.billing_type === BILLING_REQUEST ? BILLING_REQUEST : BILLING_TOKEN,
      virtual_amount: Number(record.virtual_amount) || 0,
      features: Array.isArray(record.features) ? record.features : [],
    })
    rawLimits = record.model_limits
  }
  // model_limits 在表单里只以限制行形态存在，提交前才合成回 JSON。
  const parsed = modelLimitsToRows(rawLimits, form.billing_type)
  limitRows.value = parsed.ok ? parsed.rows : []
  if (!parsed.ok && rawLimits) Message.warning(t('settingPage.plan.modelLimitsParseFailed'))
  // 历史套餐里的模型可能已从模型定价中下架，仍要并入候选列表，否则编辑会丢模型。
  const used = limitRows.value.map((row) => row.model).filter(Boolean)
  modelOptions.value = Array.from(new Set([...modelOptions.value, ...used])).sort()
  formFeatures.value = featuresFromRecord(form.features)
  modalVisible.value = true
}

async function handleSave() {
  // 保存前先在前端拦掉非法模型限制：历史上格式写错会被后端当成
  // 「不限量」静默放行，变成套餐免费用。
  const rowError = validateLimitRows(limitRows.value)
  if (rowError) {
    Message.error(limitErrorMessage(rowError))
    return
  }
  saving.value = true
  try {
    const cleanedFeatures = sanitizeFeaturesList(formFeatures.value)
    const body = {
      ...form,
      features: cleanedFeatures,
      model_limits: rowsToModelLimits(limitRows.value, form.billing_type).json,
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

onMounted(() => { loadData(); fetchModelOptions() })
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
.text-muted { color: var(--color-text-3); }
.form-group-title {
  margin: 4px 0 16px;
  padding-left: 8px;
  border-left: 3px solid rgb(var(--primary-6));
  font-size: 14px;
  font-weight: 600;
  line-height: 1.4;
  color: var(--color-text-1);
}
.form-group-title:first-child { margin-top: 0; }
.form-group-title:not(:first-child) { margin-top: 24px; }
.limits-editor { display: flex; flex-direction: column; gap: 8px; width: 100%; }
.limits-header,
.limits-row {
  display: grid;
  grid-template-columns: minmax(0, 1.7fr) minmax(0, 0.9fr) minmax(0, 1fr) minmax(0, 1fr) minmax(0, 1fr) 32px;
  gap: 8px;
  align-items: center;
}
.limits-header {
  padding: 0 2px;
  font-size: 12px;
  line-height: 1.3;
  color: var(--color-text-3);
}
.limits-empty { padding: 4px 0; font-size: 12px; color: var(--color-text-3); }
.limits-add { margin-top: 4px; }
.limits-row-remove { padding: 0; }
/* 有效期组合：有效天数固定宽度 + 时长文本自适应，两框边框连体 */
.duration-group { display: flex; width: 100%; }
.duration-group > :first-child { flex: 0 0 140px; }
.duration-group > :last-child { flex: 1; min-width: 0; }
/* 单列排版：一行一个表单项，控件统一撑满内容区，与模型限制表格左右边界对齐 */
.form-field { width: 100%; }
.arco-form-item-extra { font-size: 12px; line-height: 1.4; }
</style>

<style>
.arco-modal-plan-wide { width: 760px !important; }
</style>
