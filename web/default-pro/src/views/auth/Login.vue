<template>
  <AuthLayout>
    <div class="auth-heading">
      <h1 id="auth-title">{{ $t('auth.loginWelcome') }}</h1>
      <p>{{ $t('auth.loginSubtitle') }}</p>
    </div>
    <a-form class="auth-form" :model="form" layout="vertical" size="large" @submit="handleLogin" aria-labelledby="auth-title">
      <a-form-item field="username" hide-label>
        <a-input v-model="form.username" :placeholder="$t('auth.username')" allow-clear :aria-label="$t('auth.username')">
          <template #prefix><icon-user /></template>
        </a-input>
      </a-form-item>
      <a-form-item field="password" hide-label>
        <a-input-password v-model="form.password" :placeholder="$t('auth.password')" allow-clear :aria-label="$t('auth.password')">
          <template #prefix><icon-lock /></template>
        </a-input-password>
      </a-form-item>
      <a-form-item>
        <a-button type="primary" html-type="submit" long :loading="loading" size="large">
          {{ $t('auth.login') }}
        </a-button>
      </a-form-item>
      <div class="form-extra">
        <a-link @click="$router.push('/register')">{{ $t('auth.noAccountLink') }}</a-link>
        <a-link class="forgot-link" @click="$router.push('/reset')">{{ $t('auth.forgotPasswordLink') }}</a-link>
      </div>
      <div v-if="errorMsg" class="form-alert">
        <a-alert type="error" :show-icon="false" closable @close="errorMsg = ''">{{ errorMsg }}</a-alert>
      </div>
    </a-form>
  </AuthLayout>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AuthLayout from '@/layouts/AuthLayout.vue'
import { useAuthStore } from '@/stores/auth'
import { IconUser, IconLock } from '@arco-design/web-vue/es/icon'

const router = useRouter()
const authStore = useAuthStore()
const { t } = useI18n()

const form = ref({ username: '', password: '' })
const loading = ref(false)
const errorMsg = ref('')

async function handleLogin() {
  if (!form.value.username || !form.value.password) {
    errorMsg.value = t('auth.enterUsernamePassword')
    return
  }
  loading.value = true
  errorMsg.value = ''
  try {
    await authStore.login(form.value.username, form.value.password)
    router.push('/dashboard')
  } catch (e) {
    errorMsg.value = e.response?.data?.message || e.message || t('auth.loginFailed')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.form-extra {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 4px;
}

.form-extra :deep(.arco-link) {
  font-size: 13px;
  font-weight: 500;
}

.forgot-link {
  color: #7b8699;
}

.form-alert {
  margin-top: 20px;
}
</style>
