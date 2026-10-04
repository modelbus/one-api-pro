<template>
  <!--
    TopupSetting 充值设置页面
    版本: v0.0.24
    日期: 2026-10-04
    作者: opencode

    功能：
      - 总开关
      - 是否允许用户自定义金额（恒 1:1，支付多少到账多少）
      - 快捷金额列表（支付金额 / 到账金额，均以元输入，支持「充 10 得 15」赠送）
  -->
  <div class="topup-setting-page">
    <a-spin :loading="loading" style="width: 100%">
      <div class="section">
        <h3>{{ $t('settingPage.topup.section') }}</h3>
        <p class="section-hint">{{ $t('settingPage.topup.hint') }}</p>

        <!-- 基础开关（竖排） -->
        <a-form layout="vertical" class="setting-form">
          <a-form-item :label="$t('settingPage.topup.enabled')">
            <a-switch v-model="form.enabled" />
          </a-form-item>
          <a-form-item :label="$t('settingPage.topup.allowCustom')">
            <a-switch v-model="form.allow_custom" />
          </a-form-item>
        </a-form>

        <a-divider :margin="20" />

        <!-- 快捷金额列表 -->
        <div class="section-head">
          <h4 class="section-sub-title">{{ $t('settingPage.topup.presets') }}</h4>
          <a-button type="primary" size="small" @click="addPreset">
            <template #icon><icon-plus :size="14" /></template>
            {{ $t('settingPage.topup.add') }}
          </a-button>
        </div>
        <p class="section-hint">{{ $t('settingPage.topup.presetsHint') }}</p>

        <a-table
          :columns="columns"
          :data="form.presets"
          :pagination="false"
          :bordered="{ cell: true }"
          size="small"
          row-key="__idx"
          class="preset-table"
        >
          <template #amount="{ record, rowIndex }">
            <a-input-number
              :model-value="record.amount"
              :min="0.01"
              :precision="2"
              :step="1"
              size="small"
              style="width: 100%"
              @change="(v) => updatePreset(rowIndex, 'amount', v)"
            />
          </template>
          <template #credit="{ record, rowIndex }">
            <a-input-number
              :model-value="record.credit"
              :min="0"
              :precision="yuanInputPrecision(record.credit)"
              :step="1"
              size="small"
              style="width: 100%"
              @change="(v) => updatePreset(rowIndex, 'credit', v)"
            />
          </template>
          <template #action="{ rowIndex }">
            <a-button type="text" size="mini" status="danger" @click="removePreset(rowIndex)">
              {{ $t('settingPage.topup.delete') }}
            </a-button>
          </template>
          <template #empty>
            <div class="table-empty">{{ $t('settingPage.topup.empty') }}</div>
          </template>
        </a-table>

        <div class="form-actions">
          <a-button type="primary" :loading="saving" @click="save">{{ $t('settingPage.topup.save') }}</a-button>
        </div>
      </div>
    </a-spin>
  </div>
</template>

<script setup>
// 版本: v0.0.24
// 日期: 2026-10-04
// 作者: opencode
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'
import settingApi from '@/api/setting'
import { validateTopupPresets } from '@/utils/topup'
import { yuanInputPrecision } from '@/utils/quota'

const { t } = useI18n()

// validateTopupPresets 内部使用短 key（amountPositive / duplicate / bonusTooLow），
// 而 i18n 将它们定义在 settingPage.topup 命名空间下；这里绑定前缀后再传入，
// 否则 vue-i18n 找不到 key 会原样返回 key 字符串（历史上界面会直接显示 "bonusTooLow"）。
// The validator uses short keys while i18n defines them under settingPage.topup,
// so bind the namespace here.
const tTopup = (key, params) => t(`settingPage.topup.${key}`, params)

const loading = ref(false)
const saving = ref(false)

// presets 在前端以「元」为单位编辑（amount = 支付金额，credit = 到账金额）。
// v0.0.25 起后端 API 同样使用「元」，前后端不再需要任何 1e6 换算。
const form = reactive({
  enabled: false,
  allow_custom: true,
  presets: [],
})

const columns = computed(() => [
  { title: t('settingPage.topup.colAmount'), slotName: 'amount', width: 200 },
  { title: t('settingPage.topup.colCredit'), slotName: 'credit', width: 220 },
  { title: t('settingPage.topup.colAction'), slotName: 'action', width: 100, align: 'center' },
])

// quotaToYuanText 保留 6 位小数避免小额点零；输入框再按精度自适应展示。
// 入参已是「元」（后端口径），此处仅做浮点误差收敛。
// 版本: v0.0.25
// 日期: 2026-10-04
function quotaToYuanText(y) {
  const n = Number(y)
  return Number.isFinite(n) ? Number(n.toFixed(6)) : 0
}

function addPreset() {
  // 默认 1:1：充 10 元到账 10 元
  form.presets.push({ amount: 10, credit: 10 })
}

function removePreset(idx) {
  form.presets.splice(idx, 1)
}

function updatePreset(idx, key, val) {
  form.presets[idx][key] = val
}

async function loadSettings() {
  loading.value = true
  try {
    const { data } = await settingApi.getTopup()
    if (data.success && data.data) {
      const d = data.data
      form.enabled = !!d.enabled
      form.allow_custom = !!d.allow_custom
      form.presets = Array.isArray(d.presets) ? d.presets.map(p => ({
        amount: Number(Number(p.amount).toFixed(2)),
        credit: quotaToYuanText(p.bonus_quota),
      })) : []
    }
  } catch (e) {
    Message.error(t('settingPage.topup.loadFailed'))
  } finally {
    loading.value = false
  }
}

async function save() {
  // 前端预校验：金额 > 0、金额不重复、到账金额不低于支付金额
  // 两个字段均为「元」口径（后端 API 自 v0.0.25 起统一返回并接收元）。
  const payloadPresets = form.presets.map(p => ({
    amount: Number(p.amount),
    bonus_quota: Number(p.credit),
  }))
  const err = validateTopupPresets(payloadPresets, tTopup)
  if (err) {
    Message.error(err)
    return
  }
  saving.value = true
  try {
    const payload = {
      enabled: form.enabled,
      allow_custom: form.allow_custom,
      presets: payloadPresets,
    }
    const { data } = await settingApi.putTopup(payload)
    if (data.success) {
      Message.success(t('settingPage.topup.saved'))
    } else {
      Message.error(data.message || t('settingPage.topup.saveFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || t('settingPage.topup.saveFailed'))
  } finally {
    saving.value = false
  }
}

onMounted(() => { loadSettings() })
</script>

<style scoped>
.topup-setting-page { padding: 0 4px; }
.section h3 { font-size: 16px; font-weight: 600; color: var(--color-text-1); margin-bottom: 8px; padding: 0; }
.section-hint { color: var(--color-text-3); font-size: 13px; margin: 0 0 20px; }
.section-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; }
.section-sub-title { font-size: 14px; font-weight: 600; color: var(--color-text-1); margin: 0; }
.setting-form { width: 100%; }
.preset-table { margin-bottom: 20px; }
.table-empty { padding: 30px 0; text-align: center; color: var(--color-text-3); font-size: 13px; }
.form-actions { display: flex; justify-content: flex-end; gap: 12px; }
</style>
