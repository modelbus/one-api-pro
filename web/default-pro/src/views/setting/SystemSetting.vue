<template>
  <a-spin :loading="loading" class="setting-container">
    <!-- Server Address -->
    <div class="section">
      <h3>{{ $t('settingPage.system.serverAddress') }}</h3>
      <a-form :model="form" layout="vertical" class="setting-form">
        <a-row :gutter="16">
          <a-col :span="18"><a-form-item hide-label><a-input v-model="form.ServerAddress" placeholder="https://api.example.com" size="large" /></a-form-item></a-col>
          <a-col><a-form-item hide-label><a-button type="primary" size="large" @click="saveKey('ServerAddress')">{{ $t('settingPage.system.save') }}</a-button></a-form-item></a-col>
        </a-row>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <!-- Login & Registration Switches -->
    <div class="section">
      <h3>{{ $t('settingPage.system.loginRegister') }}</h3>
      <a-form :model="form" layout="vertical" class="setting-form">
        <a-row :gutter="[32, 16]">
          <a-col :span="6"><a-form-item :label="$t('settingPage.system.passwordLogin')"><a-switch v-model="form.PasswordLoginEnabled" @change="saveSwitch('PasswordLoginEnabled')" /></a-form-item></a-col>
          <a-col :span="6"><a-form-item :label="$t('settingPage.system.passwordRegister')"><a-switch v-model="form.PasswordRegisterEnabled" @change="saveSwitch('PasswordRegisterEnabled')" /></a-form-item></a-col>
          <a-col :span="6"><a-form-item :label="$t('settingPage.system.allowRegister')"><a-switch v-model="form.RegisterEnabled" @change="saveSwitch('RegisterEnabled')" /></a-form-item></a-col>
          <a-col :span="6"><a-form-item :label="$t('settingPage.system.emailVerification')"><a-switch v-model="form.EmailVerificationEnabled" @change="saveSwitch('EmailVerificationEnabled')" /></a-form-item></a-col>
          <a-col :span="6"><a-form-item :label="$t('settingPage.system.githubOAuth')"><a-switch v-model="form.GitHubOAuthEnabled" @change="saveSwitch('GitHubOAuthEnabled')" /></a-form-item></a-col>
          <a-col :span="6"><a-form-item :label="$t('settingPage.system.larkOAuth')"><a-switch v-model="form.LarkOAuthEnabled" @change="saveSwitch('LarkOAuthEnabled')" /></a-form-item></a-col>
          <a-col :span="6"><a-form-item :label="$t('settingPage.system.wechatLogin')"><a-switch v-model="form.WeChatAuthEnabled" @change="saveSwitch('WeChatAuthEnabled')" /></a-form-item></a-col>
          <a-col :span="6"><a-form-item :label="$t('settingPage.system.turnstile')"><a-switch v-model="form.TurnstileCheckEnabled" @change="saveSwitch('TurnstileCheckEnabled')" /></a-form-item></a-col>
        </a-row>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <!-- GitHub OAuth -->
    <div class="section">
      <h3>{{ $t('settingPage.system.githubOauthConfig') }}</h3>
      <a-form :model="form" layout="vertical" class="setting-form">
        <a-row :gutter="16">
          <a-col :span="11"><a-form-item label="Client ID"><a-input v-model="form.GitHubClientId" placeholder="GitHub OAuth Client ID" size="large" /></a-form-item></a-col>
          <a-col :span="11"><a-form-item label="Client Secret"><a-input-password v-model="form.GitHubClientSecret" placeholder="GitHub OAuth Client Secret" size="large" /></a-form-item></a-col>
        </a-row>
        <a-form-item><a-button type="primary" @click="saveKeys(['GitHubClientId','GitHubClientSecret'])">{{ $t('settingPage.system.saveGithubOauth') }}</a-button></a-form-item>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <!-- Lark OAuth -->
    <div class="section">
      <h3>{{ $t('settingPage.system.larkOauthConfig') }}</h3>
      <a-form :model="form" layout="vertical" class="setting-form">
        <a-row :gutter="16">
          <a-col :span="11"><a-form-item label="Client ID"><a-input v-model="form.LarkClientId" :placeholder="$t('settingPage.system.larkClientIdPlaceholder')" size="large" /></a-form-item></a-col>
          <a-col :span="11"><a-form-item label="Client Secret"><a-input-password v-model="form.LarkClientSecret" :placeholder="$t('settingPage.system.larkClientSecretPlaceholder')" size="large" /></a-form-item></a-col>
        </a-row>
        <a-form-item><a-button type="primary" @click="saveKeys(['LarkClientId','LarkClientSecret'])">{{ $t('settingPage.system.saveLarkOauth') }}</a-button></a-form-item>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <!-- WeChat -->
    <div class="section">
      <h3>{{ $t('settingPage.system.wechatConfig') }}</h3>
      <a-form :model="form" layout="vertical" class="setting-form">
        <a-row :gutter="[16, 0]">
          <a-col :span="7"><a-form-item :label="$t('settingPage.system.serverAddress')"><a-input v-model="form.WeChatServerAddress" :placeholder="$t('settingPage.system.wechatServerAddressPlaceholder')" size="large" /></a-form-item></a-col>
          <a-col :span="7"><a-form-item :label="$t('settingPage.system.wechatServerToken')"><a-input-password v-model="form.WeChatServerToken" :placeholder="$t('settingPage.system.wechatServerTokenPlaceholder')" size="large" /></a-form-item></a-col>
          <a-col :span="8"><a-form-item :label="$t('settingPage.system.qrCodeUrl')"><a-input v-model="form.WeChatAccountQRCodeImageURL" :placeholder="$t('settingPage.system.qrCodeUrlPlaceholder')" size="large" /></a-form-item></a-col>
        </a-row>
        <a-form-item><a-button type="primary" @click="saveKeys(['WeChatServerAddress','WeChatServerToken','WeChatAccountQRCodeImageURL'])">{{ $t('settingPage.system.saveWechat') }}</a-button></a-form-item>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <!-- Turnstile -->
    <div class="section">
      <h3>{{ $t('settingPage.system.turnstileConfig') }}</h3>
      <a-form :model="form" layout="vertical" class="setting-form">
        <a-row :gutter="16">
          <a-col :span="11"><a-form-item label="Site Key"><a-input v-model="form.TurnstileSiteKey" placeholder="Turnstile Site Key" size="large" /></a-form-item></a-col>
          <a-col :span="11"><a-form-item label="Secret Key"><a-input-password v-model="form.TurnstileSecretKey" placeholder="Turnstile Secret Key" size="large" /></a-form-item></a-col>
        </a-row>
        <a-form-item><a-button type="primary" @click="saveKeys(['TurnstileSiteKey','TurnstileSecretKey'])">{{ $t('settingPage.system.saveTurnstile') }}</a-button></a-form-item>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <!-- SMTP -->
    <div class="section">
      <h3>{{ $t('settingPage.system.smtpConfig') }}</h3>
      <a-form :model="form" layout="vertical" class="setting-form">
        <a-row :gutter="[16, 0]">
          <a-col :span="6"><a-form-item :label="$t('settingPage.system.smtpServer')"><a-input v-model="form.SMTPServer" :placeholder="$t('settingPage.system.smtpServerPlaceholder')" size="large" /></a-form-item></a-col>
          <a-col :span="4"><a-form-item :label="$t('settingPage.system.smtpPort')"><a-input v-model="form.SMTPPort" placeholder="587" size="large" /></a-form-item></a-col>
          <a-col :span="6"><a-form-item :label="$t('settingPage.system.smtpAccount')"><a-input v-model="form.SMTPAccount" :placeholder="$t('settingPage.system.smtpAccountPlaceholder')" size="large" /></a-form-item></a-col>
          <a-col :span="6"><a-form-item :label="$t('settingPage.system.smtpFrom')"><a-input v-model="form.SMTPFrom" placeholder="noreply@example.com" size="large" /></a-form-item></a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="11"><a-form-item label="Token"><a-input-password v-model="form.SMTPToken" :placeholder="$t('settingPage.system.smtpTokenPlaceholder')" size="large" /></a-form-item></a-col>
        </a-row>
        <a-form-item><a-button type="primary" @click="saveKeys(['SMTPServer','SMTPPort','SMTPAccount','SMTPFrom','SMTPToken'])">{{ $t('settingPage.system.saveSmtp') }}</a-button></a-form-item>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <!-- Appearance -->
    <div class="section">
      <h3>{{ $t('settingPage.system.appearance') }}</h3>
      <a-form :model="form" layout="vertical" class="setting-form">
        <a-row :gutter="[16, 0]">
          <a-col :span="7"><a-form-item :label="$t('settingPage.system.systemName')"><a-input v-model="form.SystemName" placeholder="One Api Pro" size="large" /></a-form-item></a-col>
          <a-col :span="7"><a-form-item label="Logo URL"><a-input v-model="form.Logo" :placeholder="$t('settingPage.system.logoPlaceholder')" size="large" /></a-form-item></a-col>
          <a-col :span="4"><a-form-item :label="$t('settingPage.system.theme')"><a-input v-model="form.Theme" placeholder="default" size="large" /></a-form-item></a-col>
        </a-row>
        <a-form-item><a-button type="primary" @click="saveKeys(['SystemName','Logo','Theme'])">{{ $t('settingPage.system.saveAppearance') }}</a-button></a-form-item>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <!-- Content -->
    <div class="section">
      <h3>{{ $t('settingPage.system.content') }}</h3>
      <a-form :model="form" layout="vertical" class="setting-form">
        <a-row :gutter="16">
          <a-col :span="11"><a-form-item :label="$t('settingPage.system.notice')"><a-textarea v-model="form.Notice" :auto-size="{minRows:3,maxRows:6}" :placeholder="$t('settingPage.system.noticePlaceholder')" /></a-form-item></a-col>
          <a-col :span="11"><a-form-item :label="$t('settingPage.system.homePageContent')"><a-textarea v-model="form.HomePageContent" :auto-size="{minRows:3,maxRows:6}" :placeholder="$t('settingPage.system.homePageContentPlaceholder')" /></a-form-item></a-col>
        </a-row>
        <a-form-item><a-button type="primary" @click="saveKeys(['Notice','HomePageContent'])">{{ $t('settingPage.system.saveContent') }}</a-button></a-form-item>
      </a-form>
    </div>
  </a-spin>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Message } from '@arco-design/web-vue'
import api from '@/api'

const { t } = useI18n()

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
  try { await api.put('/api/option/', { key, value: form[key] ? 'true' : 'false' }); Message.success(t('settingPage.system.saved')) }
  catch (e) { Message.error(t('settingPage.system.saveFailed')) }
}

async function saveKey(key) {
  try { await api.put('/api/option/', { key, value: String(form[key] ?? '') }); Message.success(t('settingPage.system.saved')) }
  catch (e) { Message.error(t('settingPage.system.saveFailed')) }
}

async function saveKeys(keys) {
  for (const k of keys) {
    try { await api.put('/api/option/', { key: k, value: String(form[k] ?? '') }) } catch (e) { /* continue */ }
  }
  Message.success(t('settingPage.system.saved'))
}

onMounted(() => { loadData() })
</script>

<style scoped>
.setting-container { padding: 4px 0; }
.section h3 { font-size: 16px; font-weight: 600; color: var(--color-text-1); margin-bottom: 20px; padding: 0; }
.setting-form { width: 100%; }
</style>
