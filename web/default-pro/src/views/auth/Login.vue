<template>
  <AuthLayout>
    <div class="auth-heading">
      <h1 id="auth-title">欢迎回来</h1>
      <p>登录您的账户</p>
    </div>
    <a-form class="auth-form" :model="form" layout="vertical" size="large" @submit="handleLogin" aria-labelledby="auth-title">
      <a-form-item field="username" hide-label>
        <a-input v-model="form.username" placeholder="用户名" allow-clear aria-label="用户名">
          <template #prefix><icon-user /></template>
        </a-input>
      </a-form-item>
      <a-form-item field="password" hide-label>
        <a-input-password v-model="form.password" placeholder="密码" allow-clear aria-label="密码">
          <template #prefix><icon-lock /></template>
        </a-input-password>
      </a-form-item>
      <a-form-item>
        <a-button type="primary" html-type="submit" long :loading="loading" size="large">
          登录
        </a-button>
      </a-form-item>
      <div class="form-extra">
        <a-link @click="$router.push('/register')">注册账号</a-link>
        <a-link class="forgot-link" @click="$router.push('/reset')">忘记密码</a-link>
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
    router.push('/dashboard')
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
