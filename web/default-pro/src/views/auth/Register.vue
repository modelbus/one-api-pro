<template>
  <AuthLayout>
    <div class="auth-heading">
      <h1 id="auth-title">创建账号</h1>
      <p>填写以下信息完成注册</p>
    </div>
    <a-form class="auth-form" :model="form" layout="vertical" size="large" @submit="handleRegister" aria-labelledby="auth-title">
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
      <a-form-item field="password2" hide-label>
        <a-input-password v-model="form.password2" placeholder="确认密码" allow-clear aria-label="确认密码">
          <template #prefix><icon-lock /></template>
        </a-input-password>
      </a-form-item>
      <a-form-item v-if="statusStore.status?.email_verification" field="email" hide-label>
        <a-input v-model="form.email" placeholder="邮箱" allow-clear aria-label="邮箱">
          <template #prefix><icon-email /></template>
        </a-input>
      </a-form-item>
      <a-form-item v-if="statusStore.status?.email_verification" field="verification_code" hide-label>
        <div class="verify-row">
          <a-input v-model="form.verification_code" placeholder="验证码" :style="{ flex: 1 }" aria-label="验证码" />
          <a-button type="outline" size="large" @click="sendVerifyCode" :loading="sending" :disabled="countdown > 0">
            {{ countdown > 0 ? `${countdown}s` : '发送验证码' }}
          </a-button>
        </div>
      </a-form-item>
      <a-form-item field="aff_code" hide-label>
        <a-input v-model="form.aff_code" placeholder="邀请码（选填）" allow-clear aria-label="邀请码（选填）">
          <template #prefix><icon-gift /></template>
        </a-input>
      </a-form-item>
      <div class="agreement-row">
        <a-checkbox v-model="agreedToTerms" aria-label="同意服务条款和隐私政策" />
        <span class="agreement-copy">
          我同意
          <router-link to="/terms" class="agreement-link">《服务条款》</router-link>
          和
          <router-link to="/privacy" class="agreement-link">《隐私政策》</router-link>
        </span>
      </div>
      <a-form-item>
        <a-button type="primary" html-type="submit" long :loading="loading" size="large">
          创建账号
        </a-button>
      </a-form-item>
      <div class="form-extra">
        <a-link @click="$router.push('/login')">已有账号？去登录</a-link>
      </div>
      <div v-if="errorMsg" class="form-alert">
        <a-alert type="error" :show-icon="false" closable @close="errorMsg = ''">{{ errorMsg }}</a-alert>
      </div>
    </a-form>
  </AuthLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import AuthLayout from '@/layouts/AuthLayout.vue'
import { useAuthStore } from '@/stores/auth'
import { useStatusStore } from '@/stores/status'
import { IconUser, IconLock, IconEmail, IconGift } from '@arco-design/web-vue/es/icon'
import api from '@/api'

const router = useRouter()
const authStore = useAuthStore()
const statusStore = useStatusStore()

const form = ref({ username: '', password: '', password2: '', email: '', verification_code: '', aff_code: '' })
const agreedToTerms = ref(false)
const loading = ref(false)
const sending = ref(false)
const countdown = ref(0)
const errorMsg = ref('')

onMounted(async () => {
  if (!statusStore.loaded) await statusStore.fetchStatus()
})

async function sendVerifyCode() {
  if (!form.value.email) return
  sending.value = true
  try {
    await api.get('/api/verification', { params: { email: form.value.email } })
    countdown.value = 60
    const timer = setInterval(() => { countdown.value--; if (countdown.value <= 0) clearInterval(timer) }, 1000)
  } catch (e) {
    errorMsg.value = '发送验证码失败'
  } finally {
    sending.value = false
  }
}

async function handleRegister() {
  if (!agreedToTerms.value) {
    errorMsg.value = '请先阅读并同意服务条款和隐私政策'
    return
  }
  if (form.value.password !== form.value.password2) {
    errorMsg.value = '两次密码不一致'
    return
  }
  loading.value = true
  errorMsg.value = ''
  try {
    await authStore.register({ ...form.value })
    router.push('/login')
  } catch (e) {
    errorMsg.value = e.response?.data?.message || e.message || '注册失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.verify-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.verify-row .arco-input {
  min-width: 0;
}

.verify-row :deep(.arco-btn) {
  flex-shrink: 0;
  min-width: 122px;
}

.agreement-row {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  min-height: 24px;
  margin: 4px 0 20px;
  color: #7b8699;
  font-size: 13px;
  line-height: 22px;
}

.agreement-row :deep(.arco-checkbox) {
  margin-top: 2px;
}

.agreement-copy {
  flex: 1;
}

.agreement-link {
  color: rgb(var(--primary-6));
  transition: color 0.2s ease;
}

.agreement-link:hover {
  color: #0e4fd6;
}

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
