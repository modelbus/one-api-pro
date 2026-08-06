<template>
  <AuthLayout>
    <a-form :model="form" layout="vertical" size="large" @submit="handleSubmit">
      <a-form-item field="email" hide-label>
        <a-input v-model="form.email" placeholder="请输入注册邮箱" allow-clear>
          <template #prefix><icon-email /></template>
        </a-input>
      </a-form-item>
      <a-form-item>
        <a-button type="primary" html-type="submit" long :loading="loading" size="large">
          发送重置邮件
        </a-button>
      </a-form-item>
      <div class="form-extra">
        <a-link @click="$router.push('/login')">返回登录</a-link>
      </div>
      <div v-if="message" style="margin-top: 16px">
        <a-alert :type="success ? 'success' : 'error'" :show-icon="false">{{ message }}</a-alert>
      </div>
    </a-form>
  </AuthLayout>
</template>

<script setup>
import { ref } from 'vue'
import AuthLayout from '@/layouts/AuthLayout.vue'
import { IconEmail } from '@arco-design/web-vue/es/icon'
import api from '@/api'

const form = ref({ email: '' })
const loading = ref(false)
const message = ref('')
const success = ref(false)

async function handleSubmit() {
  if (!form.value.email) { message.value = '请输入邮箱'; return }
  loading.value = true
  message.value = ''
  try {
    await api.post('/api/user/reset', { email: form.value.email })
    success.value = true
    message.value = '密码重置邮件已发送，请查收'
  } catch (e) {
    success.value = false
    message.value = e.response?.data?.message || '发送失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.form-extra { display: flex; justify-content: center; }
</style>
