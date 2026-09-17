<template>
  <AuthLayout>
    <div class="auth-heading">
      <h1 id="auth-title">{{ $t('auth.resetHeading') }}</h1>
      <p>{{ $t('auth.resetSubtitle') }}</p>
    </div>
    <a-form class="auth-form" :model="form" layout="vertical" size="large" @submit="handleSubmit" aria-labelledby="auth-title">
      <a-form-item field="email" hide-label>
        <a-input v-model="form.email" :placeholder="$t('auth.emailPlaceholder')" allow-clear :aria-label="$t('auth.emailAriaLabel')">
          <template #prefix><icon-email /></template>
        </a-input>
      </a-form-item>
      <a-form-item>
        <a-button type="primary" html-type="submit" long :loading="loading" size="large">
          {{ $t('auth.sendResetEmail') }}
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
import { useI18n } from 'vue-i18n'
import AuthLayout from '@/layouts/AuthLayout.vue'
import { IconEmail } from '@arco-design/web-vue/es/icon'
import api from '@/api'

const { t } = useI18n()
const form = ref({ email: '' })
const loading = ref(false)
const message = ref('')
const success = ref(false)

async function handleSubmit() {
  if (!form.value.email) { message.value = t('auth.enterEmail'); return }
  loading.value = true
  message.value = ''
  try {
    await api.post('/api/user/reset', { email: form.value.email })
    success.value = true
    message.value = t('auth.resetEmailSent')
  } catch (e) {
    success.value = false
    message.value = e.response?.data?.message || t('auth.sendFailed')
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
