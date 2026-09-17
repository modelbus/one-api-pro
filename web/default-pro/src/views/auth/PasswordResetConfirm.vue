<template>
  <AuthLayout>
    <div class="auth-heading">
      <h1 id="auth-title">{{ $t('auth.resetPassword') }}</h1>
      <p>{{ $t('auth.resetConfirmSubtitle') }}</p>
    </div>
    <a-form class="auth-form" :model="form" layout="vertical" size="large" @submit="handleSubmit" aria-labelledby="auth-title">
      <a-form-item field="password" hide-label>
        <a-input-password v-model="form.password" :placeholder="$t('auth.newPasswordPlaceholder')" allow-clear :aria-label="$t('auth.newPasswordPlaceholder')">
          <template #prefix><icon-lock /></template>
        </a-input-password>
      </a-form-item>
      <a-form-item field="password2" hide-label>
        <a-input-password v-model="form.password2" :placeholder="$t('auth.confirmNewPassword')" allow-clear :aria-label="$t('auth.confirmNewPassword')">
          <template #prefix><icon-lock /></template>
        </a-input-password>
      </a-form-item>
      <a-form-item>
        <a-button type="primary" html-type="submit" long :loading="loading" size="large">
          {{ $t('auth.resetPassword') }}
        </a-button>
      </a-form-item>
      <div class="form-extra">
        <a-link @click="$router.push('/login')">{{ $t('auth.backToLogin') }}</a-link>
      </div>
      <div v-if="message" class="form-alert">
        <a-alert :type="success ? 'success' : 'error'" :show-icon="false">{{ message }}</a-alert>
      </div>
    </a-form>
  </AuthLayout>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AuthLayout from '@/layouts/AuthLayout.vue'
import { IconLock } from '@arco-design/web-vue/es/icon'
import api from '@/api'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const form = ref({ password: '', password2: '' })
const loading = ref(false)
const message = ref('')
const success = ref(false)

async function handleSubmit() {
  if (form.value.password !== form.value.password2) { message.value = t('auth.passwordMismatch'); return }
  loading.value = true
  message.value = ''
  try {
    await api.post('/api/user/reset', {
      token: route.params.token,
      password: form.value.password,
      password2: form.value.password2,
    })
    success.value = true
    message.value = t('auth.resetAndLogin')
    setTimeout(() => router.push('/login'), 1500)
  } catch (e) {
    success.value = false
    message.value = e.response?.data?.message || t('auth.resetFailed')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.form-extra {
  display: flex;
  justify-content: center;
  margin-top: 4px;
}

.form-extra :deep(.arco-link) {
  font-size: 13px;
  font-weight: 500;
}

.form-alert {
  margin-top: 20px;
}
</style>
