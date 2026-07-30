<template>
  <AuthLayout>
    <a-form :model="form" layout="vertical" size="large" @submit="handleLogin">
      <a-form-item field="username" hide-label>
        <a-input v-model="form.username" placeholder="用户名" allow-clear>
          <template #prefix><icon-user /></template>
        </a-input>
      </a-form-item>
      <a-form-item field="password" hide-label>
        <a-input-password v-model="form.password" placeholder="密码" allow-clear>
          <template #prefix><icon-lock /></template>
        </a-input-password>
      </a-form-item>
      <a-form-item>
        <a-button type="primary" html-type="submit" long :loading="loading" size="large">
          登 录
        </a-button>
      </a-form-item>
      <div class="form-extra">
        <a-link @click="$router.push('/register')">注册账号</a-link>
        <a-link @click="$router.push('/reset')">忘记密码？</a-link>
      </div>
      <div v-if="errorMsg" style="margin-top: 16px">
        <a-alert type="error" :show-icon="false" closable @close="errorMsg = ''">{{ errorMsg }}</a-alert>
      </div>
    </a-form>
  </AuthLayout>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import AuthLayout from '@/layouts/AuthLayout.vue'
import { useAuthStore } from '@/stores/auth'
import { IconUser, IconLock } from '@arco-design/web-vue/es/icon'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({ username: '', password: '' })
const loading = ref(false)
const errorMsg = ref('')

async function handleLogin() {
  if (!form.value.username || !form.value.password) {
    errorMsg.value = '请输入用户名和密码'
    return
  }
  loading.value = true
  errorMsg.value = ''
  try {
    await authStore.login(form.value.username, form.value.password)
    if (authStore.user.username === 'root' && form.value.password === '123456') {
      router.push('/setting')
    } else {
      router.push('/')
    }
  } catch (e) {
    errorMsg.value = e.response?.data?.message || e.message || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.form-extra {
  display: flex;
  justify-content: space-between;
}
</style>
