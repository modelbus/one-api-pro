<template>
  <a-config-provider :locale="arcoLocale">
    <router-view />
  </a-config-provider>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { getArcoLocale } from '@/i18n/arco'

const { locale } = useI18n()
const arcoLocale = computed(() => getArcoLocale(locale.value))
</script>
