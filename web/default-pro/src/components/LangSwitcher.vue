<!--
  LangSwitcher 语言切换下拉组件（AdminLayout / Landing 共用）
  Reusable language switcher dropdown shared by AdminLayout and Landing.

  版本: v0.0.21
  日期: 2026-09-17
  作者: opencode
-->
<template>
  <a-select
    :model-value="locale"
    :size="size"
    :style="{ width }"
    :trigger-props="{ autoFitPopupMinWidth: true }"
    @change="handleChange"
  >
    <a-option v-for="l in SUPPORTED_LOCALES" :key="l.value" :value="l.value">
      {{ l.label }}
    </a-option>
  </a-select>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import { SUPPORTED_LOCALES, setLocale } from '@/i18n'

defineProps({
  size: { type: String, default: 'small' },
  width: { type: String, default: '110px' },
})

const { locale } = useI18n()

function handleChange(val) {
  setLocale(val)
}
</script>
