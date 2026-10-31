<template>
  <a-spin :loading="loading" class="setting-container">
    <div class="section"><h3>{{ $t('settingPage.personal.affLink') }}</h3>
      <a-form layout="vertical" class="setting-form">
        <a-form-item><a-button @click="getAff" :loading="affLoading" size="large">{{ $t('settingPage.personal.getAffLink') }}</a-button></a-form-item>
        <a-form-item v-if="affLink" :label="$t('settingPage.personal.affLink')">
          <a-input v-model="affLink" readonly size="large" style="max-width:480px" />
          <a-button type="text" size="small" @click="copyIt(affLink)" style="margin-top:8px">{{ $t('settingPage.personal.copy') }}</a-button>
        </a-form-item>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <div class="section"><h3>{{ $t('settingPage.personal.thirdPartyBinding') }}</h3>
      <a-space size="large">
        <a-button @click="bindGH" v-if="statusStore.status?.github_client_id" size="large">{{ $t('settingPage.personal.bindGithub') }}</a-button>
        <a-button @click="bindLark" v-if="statusStore.status?.lark_client_id" size="large">{{ $t('settingPage.personal.bindLark') }}</a-button>
        <a-button @click="showEmail=true" size="large">{{ $t('settingPage.personal.bindEmail') }}</a-button>
      </a-space>
    </div>
    <a-divider :margin="24" />

    <div class="section danger"><h3>{{ $t('settingPage.personal.dangerZone') }}</h3>
      <a-popconfirm :content="$t('settingPage.personal.confirmDeleteAccount')" @ok="delAccount">
        <a-button status="danger" size="large">{{ $t('settingPage.personal.deleteAccount') }}</a-button>
      </a-popconfirm>
    </div>

    <a-modal v-model:visible="showEmail" :title="$t('settingPage.personal.bindEmail')" @ok="submitEmail" :ok-loading="emailBinding">
      <a-form layout="vertical">
        <a-form-item :label="$t('settingPage.personal.email')"><a-input v-model="emailForm.email" /></a-form-item>
        <a-form-item :label="$t('settingPage.personal.verificationCode')">
          <a-space><a-input v-model="emailForm.code" :placeholder="$t('settingPage.personal.verificationCodePlaceholder')" style="width:200px" />
          <a-button @click="sendCode" :loading="emailSending" :disabled="countdown>0">{{ countdown>0?`${countdown}s`:$t('settingPage.personal.sendCode') }}</a-button></a-space>
        </a-form-item>
      </a-form>
    </a-modal>
  </a-spin>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { Message } from '@arco-design/web-vue'
import { useAuthStore } from '@/stores/auth'
import { useStatusStore } from '@/stores/status'
import api from '@/api'

const { t } = useI18n()
const authStore = useAuthStore()
const statusStore = useStatusStore()
const loading = ref(false)
const affLink = ref(''), affLoading = ref(false)
const showEmail = ref(false), emailForm = reactive({ email: '', code: '' })
const emailBinding = ref(false), emailSending = ref(false), countdown = ref(0)

async function getAff() { affLoading.value=true; try { const {data}=await api.get('/api/user/aff'); if(data.success) affLink.value = `${window.location.origin}/register?aff=${data.data}` } catch(e){ Message.error(t('settingPage.personal.fetchFailed')) } finally { affLoading.value=false } }
function copyIt(val) { navigator.clipboard?.writeText(val).then(()=>Message.success(t('settingPage.personal.copied'))) }
function bindGH() { const id=statusStore.status?.github_client_id; if(id) window.location.href=`https://github.com/login/oauth/authorize?client_id=${id}&scope=user:email` }
function bindLark() { const id=statusStore.status?.lark_client_id; if(id) window.location.href=`https://open.feishu.cn/open-apis/authen/v1/authorize?app_id=${id}&redirect_uri=${encodeURIComponent(window.location.origin+'/oauth/lark')}` }
async function sendCode() { if(!emailForm.email)return; emailSending.value=true; try{await api.get('/api/verification',{params:{email:emailForm.email}});countdown.value=60;const timer=setInterval(()=>{countdown.value--;if(countdown.value<=0)clearInterval(timer)},1000)}catch(e){Message.error(t('settingPage.personal.sendFailed'))}finally{emailSending.value=false} }
async function submitEmail() { if(!emailForm.email||!emailForm.code)return; emailBinding.value=true; try{const {data}=await api.get('/api/oauth/email/bind',{params:{email:emailForm.email,code:emailForm.code}});if(data.success){Message.success(t('settingPage.personal.bindSuccess'));showEmail.value=false}else Message.error(data.message)}catch(e){Message.error(t('settingPage.personal.bindFailed'))}finally{emailBinding.value=false} }
async function delAccount() { try{await api.delete('/api/user/self');await authStore.logout();window.location.href='/'}catch(e){Message.error(t('settingPage.personal.deleteFailed'))} }
</script>

<style scoped>
.setting-container { padding: 4px 0; }
.section h3 { font-size: 16px; font-weight: 600; color: var(--color-text-1); margin-bottom: 20px; padding: 0; }
.setting-form { width: 100%; }
.danger h3 { color: rgb(var(--danger-6)); }
</style>
