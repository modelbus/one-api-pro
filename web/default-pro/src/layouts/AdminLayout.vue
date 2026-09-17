<template>
  <a-layout class="admin-layout">
    <a-layout-sider
      v-model:collapsed="collapsed"
      :width="180"
      collapsible
      breakpoint="lg"
      :default-collapsed="false"
    >
      <div class="sidebar-logo">
        <img :src="logoSrc" class="logo-img" @error="onLogoError" />
      </div>
      <a-menu
        v-model:selected-keys="selectedKeys"
        @menu-item-click="handleMenuClick"
        class="sidebar-menu"
      >
        <a-menu-item key="/dashboard">
          <template #icon><icon-dashboard /></template>
          {{ $t('menu.dashboard') }}
        </a-menu-item>
        <a-menu-item key="/chat" v-if="chatLink">
          <template #icon><icon-message /></template>
          {{ $t('menu.chat') }}
        </a-menu-item>
        <a-menu-item key="/token">
          <template #icon><icon-code /></template>
          {{ $t('menu.token') }}
        </a-menu-item>
        <a-menu-item key="/redeem">
          <template #icon><icon-gift /></template>
          {{ $t('menu.redeem') }}
        </a-menu-item>
        <a-menu-item key="/plans">
          <template #icon><icon-gift /></template>
          {{ $t('menu.plans') }}
        </a-menu-item>
        <a-menu-item key="/orders">
          <template #icon><icon-storage /></template>
          {{ $t('menu.orders') }}
        </a-menu-item>
        <a-menu-item key="/log">
          <template #icon><icon-file /></template>
          {{ $t('menu.log') }}
        </a-menu-item>

        <template v-if="authStore.isAdmin">
          <div class="menu-divider" />
          <a-menu-item key="/admin/dashboard">
            <template #icon><icon-bar-chart /></template>
            {{ $t('menu.adminDashboard') }}
          </a-menu-item>
          <a-menu-item key="/channel">
            <template #icon><icon-apps /></template>
            {{ $t('menu.channel') }}
          </a-menu-item>
          <a-menu-item key="/admin/orders">
            <template #icon><icon-storage /></template>
            {{ $t('menu.orders') }}
          </a-menu-item>
          <a-menu-item key="/redemption">
            <template #icon><icon-gift /></template>
            {{ $t('menu.redemption') }}
          </a-menu-item>
          <a-menu-item key="/user">
            <template #icon><icon-user-group /></template>
            {{ $t('menu.user') }}
          </a-menu-item>
          <a-menu-item key="/subscription">
            <template #icon><icon-calendar /></template>
            {{ $t('menu.subscription') }}
          </a-menu-item>
          <a-menu-item key="/setting">
            <template #icon><icon-settings /></template>
            {{ $t('menu.setting') }}
          </a-menu-item>
        </template>
      </a-menu>
    </a-layout-sider>

    <a-layout class="main-area">
      <a-layout-header class="top-navbar">
        <div class="navbar-left">
          <a-breadcrumb>
            <a-breadcrumb-item>{{ $t('layout.breadcrumbHome') }}</a-breadcrumb-item>
            <a-breadcrumb-item v-if="currentTitle">{{ currentTitle }}</a-breadcrumb-item>
          </a-breadcrumb>
        </div>
        <div class="navbar-right">
          <LangSwitcher size="small" width="110px" />
          <a-dropdown trigger="hover" position="br">
            <a-space class="user-trigger">
              <span class="username">{{ authStore.user?.username }}</span>
              <a-tag :color="roleColor" size="small">{{ roleText }}</a-tag>
              <icon-down style="font-size:12px;color:var(--color-text-3)" />
            </a-space>
            <template #content>
              <a-doption @click="showProfile = true">
                <template #icon><icon-user /></template>
                {{ $t('layout.profile') }}
              </a-doption>
              <a-doption @click="handleGenToken">
                <template #icon><icon-code /></template>
                {{ $t('layout.accessToken') }}
              </a-doption>
              <a-doption @click="handleLogout">
                <template #icon><icon-export /></template>
                {{ $t('layout.logout') }}
              </a-doption>
            </template>
          </a-dropdown>
        </div>
      </a-layout-header>

      <a-layout-content class="content-area">
        <router-view />
      </a-layout-content>

      <a-layout-footer class="admin-footer">
        {{ $t('layout.footer') }}
      </a-layout-footer>
    </a-layout>

    <!-- Profile Modal -->
    <a-modal v-model:visible="showProfile" :title="$t('layout.profile')" @ok="saveProfile" :ok-loading="profileSaving" width="480">
      <a-form :model="profileForm" layout="vertical">
        <a-form-item :label="$t('layout.displayName')">
          <a-input v-model="profileForm.display_name" :placeholder="$t('layout.displayNamePlaceholder')" />
        </a-form-item>
        <a-form-item :label="$t('layout.newPassword')">
          <a-input-password v-model="profileForm.password" :placeholder="$t('layout.passwordKeepPlaceholder')" />
        </a-form-item>
        <a-form-item :label="$t('layout.confirmPassword')">
          <a-input-password v-model="profileForm.password_confirm" :placeholder="$t('layout.confirmPasswordPlaceholder')" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- Token Modal -->
    <a-modal v-model:visible="showToken" :title="$t('layout.accessToken')" :footer="false" width="520">
      <a-input v-if="accessToken" :model-value="accessToken" readonly size="large" />
      <a-button v-if="accessToken" type="text" size="small" @click="copyToken" style="margin-top:8px">{{ $t('common.copy') }}</a-button>
      <div v-if="!accessToken && !tokenLoading" style="text-align:center;padding:20px 0">
        <p style="color:var(--color-text-3);margin-bottom:16px">{{ $t('layout.generateConfirm') }}</p>
        <a-button type="primary" @click="genToken">{{ $t('layout.generate') }}</a-button>
      </div>
      <a-spin v-if="tokenLoading" style="display:flex;justify-content:center;padding:40px 0" />
    </a-modal>
  </a-layout>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Message } from '@arco-design/web-vue'
import { useAuthStore } from '@/stores/auth'
import { useStatusStore } from '@/stores/status'
import { IconDashboard, IconMessage, IconApps, IconCode, IconGift, IconArchive, IconUserGroup, IconCalendar, IconFile, IconSettings, IconExport, IconUser, IconDown, IconStorage, IconBarChart } from '@arco-design/web-vue/es/icon'
import LangSwitcher from '@/components/LangSwitcher.vue'
import api from '@/api'
import logoPng from '@/assets/logo.png'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const statusStore = useStatusStore()

const collapsed = ref(false)
const selectedKeys = ref([route.path])

const systemName = computed(() => statusStore.status?.system_name || 'One Api Pro')
const chatLink = computed(() => statusStore.status?.chat_link || '')
const currentTitle = computed(() => (route.meta?.title ? t(route.meta.title) : ''))

const logoFallback = computed(() => statusStore.status?.logo || '')
const logoSrc = ref(logoPng)
function onLogoError() {
  if (logoSrc.value !== logoFallback.value && logoFallback.value) {
    logoSrc.value = logoFallback.value
  }
}

const roleText = computed(() => {
  if (authStore.isRoot) return t('layout.roleRoot')
  if (authStore.isAdmin) return t('layout.roleAdmin')
  return t('layout.roleUser')
})
const roleColor = computed(() => authStore.isRoot ? 'orangered' : authStore.isAdmin ? 'arcoblue' : 'gray')

const handleMenuClick = (key) => { router.push(key) }
const handleLogout = async () => { await authStore.logout(); router.push('/login') }

// Profile Modal
const showProfile = ref(false)
const profileSaving = ref(false)
const profileForm = reactive({ display_name: '', password: '', password_confirm: '' })

watch(showProfile, (val) => {
  if (val) {
    profileForm.display_name = authStore.user?.display_name || ''
    profileForm.password = ''
    profileForm.password_confirm = ''
  }
})

async function saveProfile() {
  if (profileForm.password && profileForm.password !== profileForm.password_confirm) { Message.warning(t('layout.passwordMismatch')); return }
  profileSaving.value = true
  try {
    const b = { display_name: profileForm.display_name }
    if (profileForm.password) b.password = profileForm.password
    const { data } = await api.put('/api/user/self', b)
    if (data.success) { showProfile.value = false; Message.success(t('layout.saved')); if (data.data) authStore.user = data.data }
    else Message.error(data.message)
  } catch (e) { Message.error(t('layout.saveFailed')) } finally { profileSaving.value = false }
}

// Token
const showToken = ref(false)
const accessToken = ref('')
const tokenLoading = ref(false)

async function handleGenToken() {
  showToken.value = true
  accessToken.value = ''
  tokenLoading.value = false
}

async function genToken() {
  tokenLoading.value = true
  try {
    const { data } = await api.get('/api/user/token')
    if (data.success) accessToken.value = data.data || data.message
    else Message.error(data.message)
  } catch (e) { Message.error(t('layout.tokenLoadFailed')) } finally { tokenLoading.value = false }
}

function copyToken() {
  navigator.clipboard?.writeText(accessToken.value).then(() => Message.success(t('common.copied')))
}

watch(() => route.path, (path) => { selectedKeys.value = [path] })
</script>

<style scoped>
.admin-layout { height: 100vh; }

.sidebar-logo {
  height: 52px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  border-bottom: 1px solid var(--color-border-2);
  padding: 0 16px;
}
.logo-img { height: 36px; border-radius: 4px; }
.logo-text { font-size: 16px; font-weight: 700; white-space: nowrap; overflow: hidden; }
.sidebar-menu { border-right: none; }
.menu-divider { height: 1px; background: var(--color-border-2); margin: 16px 16px; }

.main-area { background: var(--color-fill-2); min-width: 0; overflow-x: hidden; }

.top-navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 52px;
  padding: 0 20px;
  background: var(--color-bg-2);
  border-bottom: 1px solid var(--color-border-2);
}
.navbar-left { display: flex; align-items: center; }
.navbar-right { display: flex; align-items: center; gap: 12px; }
.username { color: var(--color-text-2); font-size: 14px; }
.user-trigger { cursor: pointer; padding: 4px 8px; border-radius: 4px; }
.user-trigger:hover { background: var(--color-fill-2); }

.content-area { padding: 20px; overflow-y: auto; min-width: 0; overflow-x: hidden; }

.admin-footer {
  text-align: center;
  font-size: 12px;
  color: var(--color-text-4);
  padding: 12px;
  background: var(--color-bg-2);
  border-top: 1px solid var(--color-border-2);
}
</style>
