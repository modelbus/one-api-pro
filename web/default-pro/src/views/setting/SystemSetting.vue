<template>
  <a-spin :loading="loading" class="setting-container">
    <!-- Server Address -->
    <div class="section">
      <h3>服务器地址</h3>
      <a-form :model="form" layout="vertical" class="setting-form">
        <a-row :gutter="16">
          <a-col :span="18"><a-form-item hide-label><a-input v-model="form.ServerAddress" placeholder="https://api.example.com" size="large" /></a-form-item></a-col>
          <a-col><a-form-item hide-label><a-button type="primary" size="large" @click="saveKey('ServerAddress')">保存</a-button></a-form-item></a-col>
        </a-row>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <!-- Login & Registration Switches -->
    <div class="section">
      <h3>登录与注册</h3>
      <a-form :model="form" layout="vertical" class="setting-form">
        <a-row :gutter="[32, 16]">
          <a-col :span="6"><a-form-item label="密码登录"><a-switch v-model="form.PasswordLoginEnabled" @change="saveSwitch('PasswordLoginEnabled')" /></a-form-item></a-col>
          <a-col :span="6"><a-form-item label="密码注册"><a-switch v-model="form.PasswordRegisterEnabled" @change="saveSwitch('PasswordRegisterEnabled')" /></a-form-item></a-col>
          <a-col :span="6"><a-form-item label="允许注册"><a-switch v-model="form.RegisterEnabled" @change="saveSwitch('RegisterEnabled')" /></a-form-item></a-col>
          <a-col :span="6"><a-form-item label="邮箱验证"><a-switch v-model="form.EmailVerificationEnabled" @change="saveSwitch('EmailVerificationEnabled')" /></a-form-item></a-col>
          <a-col :span="6"><a-form-item label="GitHub OAuth"><a-switch v-model="form.GitHubOAuthEnabled" @change="saveSwitch('GitHubOAuthEnabled')" /></a-form-item></a-col>
          <a-col :span="6"><a-form-item label="飞书 OAuth"><a-switch v-model="form.LarkOAuthEnabled" @change="saveSwitch('LarkOAuthEnabled')" /></a-form-item></a-col>
          <a-col :span="6"><a-form-item label="微信登录"><a-switch v-model="form.WeChatAuthEnabled" @change="saveSwitch('WeChatAuthEnabled')" /></a-form-item></a-col>
          <a-col :span="6"><a-form-item label="Turnstile"><a-switch v-model="form.TurnstileCheckEnabled" @change="saveSwitch('TurnstileCheckEnabled')" /></a-form-item></a-col>
        </a-row>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <!-- GitHub OAuth -->
    <div class="section">
      <h3>GitHub OAuth 配置</h3>
      <a-form :model="form" layout="vertical" class="setting-form">
        <a-row :gutter="16">
          <a-col :span="11"><a-form-item label="Client ID"><a-input v-model="form.GitHubClientId" placeholder="GitHub OAuth Client ID" size="large" /></a-form-item></a-col>
          <a-col :span="11"><a-form-item label="Client Secret"><a-input-password v-model="form.GitHubClientSecret" placeholder="GitHub OAuth Client Secret" size="large" /></a-form-item></a-col>
        </a-row>
        <a-form-item><a-button type="primary" @click="saveKeys(['GitHubClientId','GitHubClientSecret'])">保存 GitHub OAuth</a-button></a-form-item>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <!-- Lark OAuth -->
    <div class="section">
      <h3>飞书 OAuth 配置</h3>
      <a-form :model="form" layout="vertical" class="setting-form">
        <a-row :gutter="16">
          <a-col :span="11"><a-form-item label="Client ID"><a-input v-model="form.LarkClientId" placeholder="飞书 Client ID" size="large" /></a-form-item></a-col>
          <a-col :span="11"><a-form-item label="Client Secret"><a-input-password v-model="form.LarkClientSecret" placeholder="飞书 Client Secret" size="large" /></a-form-item></a-col>
        </a-row>
        <a-form-item><a-button type="primary" @click="saveKeys(['LarkClientId','LarkClientSecret'])">保存飞书 OAuth</a-button></a-form-item>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <!-- WeChat -->
    <div class="section">
      <h3>微信配置</h3>
      <a-form :model="form" layout="vertical" class="setting-form">
        <a-row :gutter="[16, 0]">
          <a-col :span="7"><a-form-item label="服务器地址"><a-input v-model="form.WeChatServerAddress" placeholder="微信服务器地址" size="large" /></a-form-item></a-col>
          <a-col :span="7"><a-form-item label="服务器Token"><a-input-password v-model="form.WeChatServerToken" placeholder="微信服务器 Token" size="large" /></a-form-item></a-col>
          <a-col :span="8"><a-form-item label="二维码URL"><a-input v-model="form.WeChatAccountQRCodeImageURL" placeholder="二维码图片 URL" size="large" /></a-form-item></a-col>
        </a-row>
        <a-form-item><a-button type="primary" @click="saveKeys(['WeChatServerAddress','WeChatServerToken','WeChatAccountQRCodeImageURL'])">保存微信配置</a-button></a-form-item>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <!-- Turnstile -->
    <div class="section">
      <h3>Turnstile 配置</h3>
      <a-form :model="form" layout="vertical" class="setting-form">
        <a-row :gutter="16">
          <a-col :span="11"><a-form-item label="Site Key"><a-input v-model="form.TurnstileSiteKey" placeholder="Turnstile Site Key" size="large" /></a-form-item></a-col>
          <a-col :span="11"><a-form-item label="Secret Key"><a-input-password v-model="form.TurnstileSecretKey" placeholder="Turnstile Secret Key" size="large" /></a-form-item></a-col>
        </a-row>
        <a-form-item><a-button type="primary" @click="saveKeys(['TurnstileSiteKey','TurnstileSecretKey'])">保存 Turnstile</a-button></a-form-item>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <!-- SMTP -->
    <div class="section">
      <h3>SMTP 配置</h3>
      <a-form :model="form" layout="vertical" class="setting-form">
        <a-row :gutter="[16, 0]">
          <a-col :span="6"><a-form-item label="服务器"><a-input v-model="form.SMTPServer" placeholder="SMTP服务器地址" size="large" /></a-form-item></a-col>
          <a-col :span="4"><a-form-item label="端口"><a-input v-model="form.SMTPPort" placeholder="587" size="large" /></a-form-item></a-col>
          <a-col :span="6"><a-form-item label="账号"><a-input v-model="form.SMTPAccount" placeholder="SMTP账号" size="large" /></a-form-item></a-col>
          <a-col :span="6"><a-form-item label="发件人"><a-input v-model="form.SMTPFrom" placeholder="noreply@example.com" size="large" /></a-form-item></a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="11"><a-form-item label="Token"><a-input-password v-model="form.SMTPToken" placeholder="SMTP密码/Token" size="large" /></a-form-item></a-col>
        </a-row>
        <a-form-item><a-button type="primary" @click="saveKeys(['SMTPServer','SMTPPort','SMTPAccount','SMTPFrom','SMTPToken'])">保存 SMTP</a-button></a-form-item>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <!-- Appearance -->
    <div class="section">
      <h3>系统外观</h3>
      <a-form :model="form" layout="vertical" class="setting-form">
        <a-row :gutter="[16, 0]">
          <a-col :span="7"><a-form-item label="系统名称"><a-input v-model="form.SystemName" placeholder="One Api Pro" size="large" /></a-form-item></a-col>
          <a-col :span="7"><a-form-item label="Logo URL"><a-input v-model="form.Logo" placeholder="Logo图片地址" size="large" /></a-form-item></a-col>
          <a-col :span="4"><a-form-item label="主题"><a-input v-model="form.Theme" placeholder="default" size="large" /></a-form-item></a-col>
        </a-row>
        <a-form-item><a-button type="primary" @click="saveKeys(['SystemName','Logo','Theme'])">保存外观设置</a-button></a-form-item>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <!-- Content -->
    <div class="section">
      <h3>内容</h3>
      <a-form :model="form" layout="vertical" class="setting-form">
        <a-row :gutter="16">
          <a-col :span="11"><a-form-item label="系统公告"><a-textarea v-model="form.Notice" :auto-size="{minRows:3,maxRows:6}" placeholder="系统公告内容" /></a-form-item></a-col>
          <a-col :span="11"><a-form-item label="首页内容"><a-textarea v-model="form.HomePageContent" :auto-size="{minRows:3,maxRows:6}" placeholder="首页HTML/Markdown" /></a-form-item></a-col>
        </a-row>
        <a-form-item><a-button type="primary" @click="saveKeys(['Notice','HomePageContent'])">保存内容</a-button></a-form-item>
      </a-form>
    </div>
  </a-spin>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import api from '@/api'

const loading = ref(false)
const form = reactive({
  ServerAddress: '', PasswordLoginEnabled: false, PasswordRegisterEnabled: false,
  RegisterEnabled: false, EmailVerificationEnabled: false, GitHubOAuthEnabled: false,
  LarkOAuthEnabled: false, WeChatAuthEnabled: false, TurnstileCheckEnabled: false,
  GitHubClientId: '', GitHubClientSecret: '', LarkClientId: '', LarkClientSecret: '',
  WeChatServerAddress: '', WeChatServerToken: '', WeChatAccountQRCodeImageURL: '',
  TurnstileSiteKey: '', TurnstileSecretKey: '',
  SMTPServer: '', SMTPPort: '', SMTPAccount: '', SMTPFrom: '', SMTPToken: '',
  SystemName: '', Logo: '', Theme: '', Notice: '', HomePageContent: ''
})

const allKeys = ['ServerAddress','PasswordLoginEnabled','PasswordRegisterEnabled','RegisterEnabled','EmailVerificationEnabled','GitHubOAuthEnabled','LarkOAuthEnabled','WeChatAuthEnabled','TurnstileCheckEnabled','GitHubClientId','GitHubClientSecret','LarkClientId','LarkClientSecret','WeChatServerAddress','WeChatServerToken','WeChatAccountQRCodeImageURL','TurnstileSiteKey','TurnstileSecretKey','SMTPServer','SMTPPort','SMTPAccount','SMTPFrom','SMTPToken','SystemName','Logo','Theme','Notice','HomePageContent']

async function loadData() {
  loading.value = true
  try {
    const { data } = await api.get('/api/option/')
    if (data.success && data.data) {
      const items = Array.isArray(data.data) ? data.data : Object.entries(data.data).map(([k,v]) => ({ key: k, value: String(v) }))
      items.forEach(item => { if (allKeys.includes(item.key)) { form[item.key] = item.value === 'true' ? true : item.value === 'false' ? false : item.value } })
    }
  } catch (e) { /* ignore */ } finally { loading.value = false }
}

async function saveSwitch(key) {
  try { await api.put('/api/option/', { key, value: form[key] ? 'true' : 'false' }); Message.success('已保存') }
  catch (e) { Message.error('保存失败') }
}

async function saveKey(key) {
  try { await api.put('/api/option/', { key, value: String(form[key] ?? '') }); Message.success('已保存') }
  catch (e) { Message.error('保存失败') }
}

async function saveKeys(keys) {
  for (const k of keys) {
    try { await api.put('/api/option/', { key: k, value: String(form[k] ?? '') }) } catch (e) { /* continue */ }
  }
  Message.success('已保存')
}

onMounted(() => { loadData() })
</script>

<style scoped>
.setting-container { padding: 4px 0; }
.section h3 { font-size: 16px; font-weight: 600; color: var(--color-text-1); margin-bottom: 20px; padding: 0; }
.setting-form { width: 100%; }
</style>
