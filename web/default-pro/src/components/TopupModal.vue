<template>
  <!--
    TopupModal 在线充值弹窗（可复用组件）
    版本: v0.0.10
    日期: 2026-09-06
    作者: opencode

    Props:
      modelValue: Boolean  - v-model 控制显示
      settings: Object      - { enabled, allow_custom, exchange_rate, presets: [{amount, bonus_quota}] }
      title: String         - 弹窗标题（默认"在线充值"）
    Emits:
      update:modelValue     - v-model 关闭
      success(order)       - 下单成功（订单对象）
      error(msg)            - 下单失败
      pay({url|note, payMethod, amount}) - 需要展示支付二维码/转账信息时
  -->
  <a-modal
    :visible="modelValue"
    :title="title"
    :footer="false"
    :mask-closable="true"
    :width="480"
    class="topup-modal"
    @cancel="close"
    @update:visible="(v) => emit('update:modelValue', v)"
  >
    <div v-if="hasUsableConfig" class="topup-body">
      <!-- 快捷金额 chips + 可选的自定义金额 chip（仅开关开启时显示） -->
      <div class="topup-section">
        <div class="topup-label">选择金额</div>
        <div class="topup-presets">
          <button
            v-for="(p, idx) in settings.presets"
            :key="idx"
            type="button"
            class="preset-chip"
            :class="{ active: selectedIdx === idx }"
            @click="selectPreset(idx)"
          >
            <span class="preset-amount">¥{{ p.amount }}</span>
            <span class="preset-bonus">送 {{ formatNumber(p.bonus_quota) }}</span>
          </button>
          <!-- 最后一个固定为「自定义金额」chip（前提：allow_custom=true） -->
          <button
            v-if="settings.allow_custom"
            type="button"
            class="preset-chip preset-custom"
            :class="{ active: selectedIdx === CUSTOM_IDX }"
            @click="selectCustom"
          >
            <span class="preset-amount">自定义</span>
            <span class="preset-bonus">输入金额</span>
          </button>
        </div>
      </div>

      <!-- 自定义金额输入框：仅当选中「自定义」chip 时出现 -->
      <div v-if="isCustomSelected" class="topup-section">
        <a-input-number
          v-model="customAmount"
          :min="0.01"
          :step="1"
          :precision="2"
          placeholder="请输入充值金额"
          size="large"
          style="width: 100%"
          @change="onCustomChange"
        />
        <div class="topup-hint">将获得 {{ formatNumber(calcCustomBonus) }} 额度</div>
      </div>

      <!-- 支付方式 -->
      <div class="topup-section">
        <div class="topup-label">支付方式</div>
        <div class="pay-picker-list">
          <button
            v-for="m in availableMethods"
            :key="m.name"
            type="button"
            class="pay-picker-item"
            :class="[m.name, { active: selectedPayMethod === m.name }]"
            @click="selectedPayMethod = m.name"
          >
            <component :is="iconOf(m.name)" :size="28" class="pay-picker-icon" />
            <span class="pay-picker-name">{{ m.label }}</span>
            <span v-if="selectedPayMethod === m.name" class="pay-picker-check">✓</span>
          </button>
          <div v-if="availableMethods.length === 0" class="pay-picker-empty">
            暂无可用的支付方式
          </div>
        </div>
      </div>

      <div class="topup-summary">
        <span>支付金额</span>
        <span class="summary-amount">¥{{ formatAmount(finalPayAmount) }}</span>
      </div>

      <div class="topup-footer">
        <a-button @click="close">取消</a-button>
        <a-button
          type="primary"
          :loading="submitting"
          :disabled="!canSubmit"
          @click="onConfirm"
        >
          确认充值
        </a-button>
      </div>
    </div>

    <!-- 加载中 / 无配置时 -->
    <div v-else class="topup-loading">正在加载充值配置…</div>
  </a-modal>
</template>

<script setup>
// 版本: v0.0.10
// 日期: 2026-09-06
// 作者: opencode
import { ref, computed, watch } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconWechatpay, IconAlipayCircle, IconSafe } from '@arco-design/web-vue/es/icon'
import topupApi from '@/api/topup'
import paymentApi from '@/api/payment'
import { formatNumber, formatAmount, calcCustomBonus as calcCustomBonusFn } from '@/utils/topup'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  settings: { type: Object, default: () => null },
  title: { type: String, default: '在线充值' },
})
const emit = defineEmits(['update:modelValue', 'success', 'error', 'pay'])

// 选中状态：
//   selectedIdx === CUSTOM_IDX(-1)  表示选中「自定义金额」chip
//   selectedIdx >= 0                 表示选中第 N 个 preset chip
// 当 allow_custom=false 时 CUSTOM_IDX 不可达，selectedIdx 必定 >= 0。
const CUSTOM_IDX = -1
const selectedIdx = ref(0)
const customAmount = ref(null)
const submitting = ref(false)
const selectedPayMethod = ref('wechat')
const availableMethods = ref([])

// 计算属性：是否至少有 preset 或自定义金额可用
const hasUsableConfig = computed(() => {
  if (!props.settings) return false
  const hasPresets = props.settings.presets && props.settings.presets.length > 0
  return hasPresets || props.settings.allow_custom
})

// 是否当前选中了「自定义金额」chip
const isCustomSelected = computed(() => selectedIdx.value === CUSTOM_IDX)

// 计算属性
const canSubmit = computed(() => {
  if (!selectedPayMethod.value) return false
  if (isCustomSelected.value) {
    return customAmount.value && Number(customAmount.value) > 0
  }
  return selectedIdx.value >= 0
})

const finalPayAmount = computed(() => {
  if (isCustomSelected.value) {
    return customAmount.value && Number(customAmount.value) > 0
      ? Number(customAmount.value)
      : 0
  }
  if (selectedIdx.value >= 0 && props.settings?.presets?.[selectedIdx.value]) {
    return props.settings.presets[selectedIdx.value].amount
  }
  return 0
})

const calcCustomBonus = computed(() => calcCustomBonusFn(customAmount.value, props.settings?.exchange_rate))

// 方法
function close() {
  emit('update:modelValue', false)
}

function selectPreset(idx) {
  selectedIdx.value = idx
  customAmount.value = null
}

function selectCustom() {
  selectedIdx.value = CUSTOM_IDX
}

function onCustomChange() {
  // 用户在输入框里改了金额，自动保持 CUSTOM_IDX 选中态
  if (customAmount.value && Number(customAmount.value) > 0) {
    selectedIdx.value = CUSTOM_IDX
  }
}

function iconOf(name) {
  if (name === 'wechat') return IconWechatpay
  if (name === 'alipay') return IconAlipayCircle
  return IconSafe
}

async function loadPaymentStatus() {
  try {
    const { data } = await paymentApi.status()
    const d = data?.data || {}
    availableMethods.value = (d.methods || []).filter(m => m.enabled && ['wechat', 'alipay'].includes(m.name))
    if (availableMethods.value.length > 0 && !availableMethods.value.find(m => m.name === selectedPayMethod.value)) {
      selectedPayMethod.value = availableMethods.value[0].name
    }
    return !!d.any_enabled
  } catch (e) {
    availableMethods.value = []
    return false
  }
}

async function onConfirm() {
  if (!canSubmit.value) return
  submitting.value = true
  try {
    const payload = { pay_method: selectedPayMethod.value }
    if (isCustomSelected.value) {
      payload.amount = Number(customAmount.value)
    } else if (selectedIdx.value >= 0 && props.settings?.presets?.[selectedIdx.value]) {
      payload.preset_amount = props.settings.presets[selectedIdx.value].amount
    }
    const res = await topupApi.createOrder(payload)
    const data = res?.data
    if (!data?.success) {
      Message.error(data?.message || '下单失败')
      emit('error', data?.message || '下单失败')
      return
    }
    emit('success', data)
    close()
    // 触发支付二维码弹窗（在父组件中）
    const pay = data?.pay || {}
    const url = pay.pay_url || pay.qr_code
    if (pay.status === 'success' && url) {
      emit('pay', { url, payMethod: selectedPayMethod.value, amount: data.amount })
    } else if (pay.status === 'success' && pay.note) {
      emit('pay', { note: pay.note, payMethod: selectedPayMethod.value, amount: data.amount })
    } else if (pay.status === 'warning') {
      Message.error(pay.warning || '发起支付失败，请稍后重试')
    }
  } catch (e) {
    const msg = e.response?.data?.message || e.message || '下单失败'
    Message.error(msg)
    emit('error', msg)
  } finally {
    submitting.value = false
  }
}

watch(() => props.modelValue, async (v) => {
  if (v) {
    // 默认选中第一项 preset；如果没有 preset（仅允许自定义）则进入自定义模式
    selectedIdx.value = props.settings?.presets?.length > 0 ? 0 : -1
    customAmount.value = null
    await loadPaymentStatus()
  }
})
</script>

<style scoped>
.topup-body { padding: 4px 0; }
.topup-section { margin-bottom: 16px; }
.topup-label {
  font-size: 13px;
  font-weight: 600;
  color: #1D1D1F;
  margin-bottom: 10px;
}
.topup-presets {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}
.preset-chip {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  padding: 14px 8px;
  border: 1.5px solid #E5E5E7;
  border-radius: 10px;
  background: #fff;
  cursor: pointer;
  transition: all 0.15s;
}
.preset-chip:hover { border-color: #C7C7CC; }
.preset-chip.active {
  border-color: #007AFF;
  background: rgba(0, 122, 255, 0.04);
  box-shadow: 0 0 0 2px rgba(0, 122, 255, 0.12);
}
.preset-custom .preset-amount { color: #007AFF; }
.preset-amount { font-size: 18px; font-weight: 700; color: #1D1D1F; }
.preset-bonus { font-size: 11px; color: #86868B; }
.topup-hint { font-size: 12px; color: #86868B; margin-top: 6px; }

.pay-picker-list {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}
.pay-picker-item {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px 10px;
  border: 1.5px solid #E5E5E7;
  border-radius: 10px;
  background: #fff;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  color: #1D1D1F;
  transition: all 0.15s;
}
.pay-picker-item:hover:not(:disabled) { border-color: #C7C7CC; }
.pay-picker-item.active { border-width: 1.5px; }
.pay-picker-icon {
  flex-shrink: 0;
  font-size: 28px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}
.pay-picker-item.wechat .pay-picker-icon { color: #07C160; }
.pay-picker-item.wechat.active {
  border-color: #07C160;
  background: rgba(7, 193, 96, 0.04);
  box-shadow: 0 0 0 2px rgba(7, 193, 96, 0.12);
}
.pay-picker-item.wechat.active .pay-picker-check { color: #07C160; }
.pay-picker-item.alipay .pay-picker-icon { color: #1677FF; }
.pay-picker-item.alipay.active {
  border-color: #1677FF;
  background: rgba(22, 119, 255, 0.04);
  box-shadow: 0 0 0 2px rgba(22, 119, 255, 0.12);
}
.pay-picker-item.alipay.active .pay-picker-check { color: #1677FF; }
.pay-picker-name { flex: 1; text-align: left; }
.pay-picker-check {
  flex-shrink: 0;
  font-size: 14px;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
.pay-picker-empty {
  grid-column: 1 / -1;
  text-align: center;
  color: #86868B;
  font-size: 13px;
  padding: 20px 0;
}

.topup-summary {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 14px;
  background: #F9FAFB;
  border-radius: 8px;
  margin-bottom: 16px;
}
.topup-summary > span:first-child { font-size: 13px; color: #86868B; }
.summary-amount { font-size: 20px; font-weight: 700; color: #007AFF; }

.topup-footer { display: flex; justify-content: flex-end; gap: 12px; }
.topup-loading { text-align: center; color: #86868B; padding: 40px 0; }
</style>
