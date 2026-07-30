<template>
  <div class="setting-layout">
    <a-card :bordered="false" class="setting-card">
      <div class="setting-sidebar">
        <div class="sidebar-title">设置</div>
        <a-menu :selected-keys="[activeKey]" @menu-item-click="onMenuClick" class="setting-menu">
          <a-menu-item key="system"><template #icon><icon-settings /></template>系统</a-menu-item>
          <a-menu-item key="cluster" v-if="authStore.isRoot"><template #icon><icon-storage /></template>集群</a-menu-item>
          <a-menu-item key="operation" v-if="authStore.isRoot"><template #icon><icon-tool /></template>运营</a-menu-item>
          <a-menu-item key="pricing" v-if="authStore.isRoot"><template #icon><icon-tags /></template>定价</a-menu-item>
          <a-menu-item key="plan" v-if="authStore.isRoot"><template #icon><icon-calendar /></template>套餐</a-menu-item>
          <a-menu-item key="personal"><template #icon><icon-user /></template>个人</a-menu-item>
        </a-menu>
      </div>
      <div class="setting-content">
        <router-view />
      </div>
    </a-card>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { IconSettings, IconStorage, IconTool, IconTags, IconCalendar, IconUser } from '@arco-design/web-vue/es/icon'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const activeKey = computed(() => route.path.split('/').pop() || 'system')

function onMenuClick(key) {
  router.push(`/setting/${key}`)
}
</script>

<style scoped>
.setting-layout { height: 100%; }
.setting-card { height: 100%; display: flex; padding: 0; }
.setting-card :deep(.arco-card-body) { display: flex; padding: 0; width: 100%; }
.setting-sidebar { width: 160px; flex-shrink: 0; border-right: 1px solid var(--color-border-2); padding: 20px 0; }
.sidebar-title { padding: 0 20px 16px; font-size: 16px; font-weight: 700; color: var(--color-text-1); }
.setting-menu { border: none; }
.setting-content { flex: 1; padding: 24px 28px; overflow-y: auto; overflow-x: hidden; min-width: 0; }
</style>
