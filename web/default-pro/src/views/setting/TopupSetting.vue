<template>
  <!--
    TopupSetting 充值设置页面
    版本: v0.0.10
    日期: 2026-09-06
    作者: opencode

    功能：
      - 总开关
      - 是否允许用户自定义金额
      - 自定义金额换算比例（1 元 = X quota）
      - 快捷金额列表（可增删改）
  -->
  <div class="topup-setting-page">
    <a-spin :loading="loading" style="width: 100%">
      <div class="section">
        <h3>在线充值</h3>
        <p class="section-hint">开启后，用户可在控制台通过微信/支付宝在线支付充值余额。</p>

        <!-- 基础开关（竖排） -->
        <a-form layout="vertical" class="setting-form">
          <a-form-item label="启用在线充值">
            <a-switch v-model="form.enabled" />
          </a-form-item>
          <a-form-item label="允许用户自定义金额">
            <a-switch v-model="form.allow_custom" />
          </a-form-item>
          <a-form-item label="自定义金额 1 元 = (quota)">
            <a-input-number
              v-model="form.exchange_rate"
              :min="1"
              :step="1"
              :precision="0"
              placeholder="默认 1（即 1:1 换算）"
              style="width: 320px"
            />
          </a-form-item>
        </a-form>

        <a-divider :margin="20" />

        <!-- 快捷金额列表 -->
        <div class="section-head">
          <h4 class="section-sub-title">快捷金额</h4>
          <a-button type="primary" size="small" @click="addPreset">
            <template #icon><icon-plus :size="14" /></template>
            添加
          </a-button>
        </div>

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
          <template #bonus_quota="{ record, rowIndex }">
            <a-input-number
              :model-value="record.bonus_quota"
              :min="0"
              :precision="0"
              :step="1000"
              size="small"
              style="width: 100%"
              @change="(v) => updatePreset(rowIndex, 'bonus_quota', v)"
            />
          </template>
          <template #action="{ rowIndex }">
            <a-button type="text" size="mini" status="danger" @click="removePreset(rowIndex)">
              删除
            </a-button>
          </template>
          <template #empty>
            <div class="table-empty">暂无快捷金额，点击右上角"添加"新增</div>
          </template>
        </a-table>

        <div class="form-actions">
          <a-button type="primary" :loading="saving" @click="save">保存设置</a-button>
        </div>
      </div>
    </a-spin>
  </div>
</template>

<script setup>
// 版本: v0.0.10
// 日期: 2026-09-06
// 作者: opencode
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'
import settingApi from '@/api/setting'
import { validateTopupPresets } from '@/utils/topup'

const loading = ref(false)
const saving = ref(false)

const form = reactive({
  enabled: false,
  allow_custom: true,
  exchange_rate: 1,
  presets: [],
})

const columns = [
  { title: '支付金额 (元)', slotName: 'amount', width: 200 },
  { title: '获得额度 (quota)', slotName: 'bonus_quota', width: 220 },
  { title: '操作', slotName: 'action', width: 100, align: 'center' },
]

function addPreset() {
  form.presets.push({ amount: 10, bonus_quota: 10 })
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
      form.exchange_rate = Number(d.exchange_rate || 1)
      form.presets = Array.isArray(d.presets) ? d.presets.map(p => ({
        amount: Number(p.amount),
        bonus_quota: Number(p.bonus_quota),
      })) : []
    }
  } catch (e) {
    Message.error('加载设置失败')
  } finally {
    loading.value = false
  }
}

async function save() {
  // 前端预校验：快捷金额金额不允许重复（v0.0.10 2026-09-06 新增）
  const err = validateTopupPresets(form.presets)
  if (err) {
    Message.error(err)
    return
  }
  saving.value = true
  try {
    const payload = {
      enabled: form.enabled,
      allow_custom: form.allow_custom,
      exchange_rate: Number(form.exchange_rate || 1),
      presets: form.presets.map(p => ({
        amount: Number(p.amount),
        bonus_quota: Number(p.bonus_quota),
      })),
    }
    const { data } = await settingApi.putTopup(payload)
    if (data.success) {
      Message.success('已保存')
    } else {
      Message.error(data.message || '保存失败')
    }
  } catch (e) {
    Message.error(e.response?.data?.message || '保存失败')
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
