<template>
  <a-spin :loading="loading" class="setting-container">
    <div class="section"><h3>个人信息</h3>
      <a-form :model="form" layout="vertical" class="setting-form">
        <a-row :gutter="[24,8]">
          <a-col :span="8"><a-form-item label="显示名称"><a-input v-model="form.display_name" placeholder="请输入显示名称" size="large" /></a-form-item></a-col>
          <a-col :span="8"><a-form-item label="新密码"><a-input-password v-model="form.password" placeholder="留空则不修改" size="large" /></a-form-item></a-col>
          <a-col :span="8"><a-form-item label="确认密码"><a-input-password v-model="form.password_confirm" placeholder="请再次输入新密码" size="large" /></a-form-item></a-col>
        </a-row>
        <a-form-item><a-button type="primary" :loading="saving" @click="saveProfile" size="large">保存修改</a-button></a-form-item>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <div class="section"><h3>访问令牌</h3>
      <a-form layout="vertical" class="setting-form">
        <a-form-item><a-button @click="genToken" :loading="tokenLoading" size="large">生成系统访问令牌</a-button></a-form-item>
        <a-form-item v-if="accessToken" label="令牌（仅显示一次，请复制保存）">
          <a-input v-model="accessToken" readonly size="large" style="max-width:480px" />
          <a-button type="text" size="small" @click="copyIt(accessToken)" style="margin-top:8px">复制</a-button>
        </a-form-item>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <div class="section"><h3>邀请链接</h3>
      <a-form layout="vertical" class="setting-form">
        <a-form-item><a-button @click="getAff" :loading="affLoading" size="large">获取邀请链接</a-button></a-form-item>
        <a-form-item v-if="affLink" label="邀请链接">
          <a-input v-model="affLink" readonly size="large" style="max-width:480px" />
          <a-button type="text" size="small" @click="copyIt(affLink)" style="margin-top:8px">复制</a-button>
        </a-form-item>
      </a-form>
    </div>
    <a-divider :margin="24" />

    <div class="section"><h3>第三方绑定</h3>
      <a-space size="large">
        <a-button @click="bindGH" v-if="statusStore.status?.github_client_id" size="large">绑定 GitHub</a-button>
        <a-button @click="bindLark" v-if="statusStore.status?.lark_client_id" size="large">绑定飞书</a-button>
        <a-button @click="showEmail=true" size="large">绑定邮箱</a-button>
      </a-space>
    </div>
    <a-divider :margin="24" />

    <div class="section danger"><h3>危险操作</h3>
      <a-popconfirm content="确定删除账号？不可撤销！" @ok="delAccount">
        <a-button status="danger" size="large">删除我的账号</a-button>
      </a-popconfirm>
    </div>

    <a-modal v-model:visible="showEmail" title="绑定邮箱" @ok="submitEmail" :ok-loading="emailBinding">
      <a-form layout="vertical">
        <a-form-item label="邮箱"><a-input v-model="emailForm.email" /></a-form-item>
        <a-form-item label="验证码">
          <a-space><a-input v-model="emailForm.code" placeholder="验证码" style="width:200px" />
          <a-button @click="sendCode" :loading="emailSending" :disabled="countdown>0">{{ countdown>0?`${countdown}s`:'发送验证码' }}</a-button></a-space>
        </a-form-item>
      </a-form>
    </a-modal>
  </a-spin>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useAuthStore } from '@/stores/auth'
import { useStatusStore } from '@/stores/status'
import api from '@/api'

const authStore = useAuthStore()
const statusStore = useStatusStore()
const loading = ref(false), saving = ref(false)
const form = reactive({ display_name: '', password: '', password_confirm: '' })
const accessToken = ref(''), tokenLoading = ref(false)
const affLink = ref(''), affLoading = ref(false)
const showEmail = ref(false), emailForm = reactive({ email: '', code: '' })
const emailBinding = ref(false), emailSending = ref(false), countdown = ref(0)

function loadData() { form.display_name = authStore.user?.display_name || ''; form.password = ''; form.password_confirm = '' }
async function saveProfile() {
  if (form.password && form.password !== form.password_confirm) { Message.warning('两次密码不一致'); return }
  saving.value = true
  try { const b = { display_name: form.display_name }; if (form.password) b.password = form.password; const { data } = await api.put('/api/user/self', b); if (data.success) { Message.success('已保存'); form.password=''; form.password_confirm=''; if(data.data) authStore.user = data.data } else Message.error(data.message) } catch (e) { Message.error('保存失败') } finally { saving.value = false }
}
async function genToken() { tokenLoading.value=true; try { const {data}=await api.get('/api/user/token'); if(data.success) accessToken.value=data.data||data.message } catch(e){ Message.error('获取失败') } finally { tokenLoading.value=false } }
async function getAff() { affLoading.value=true; try { const {data}=await api.get('/api/user/aff'); if(data.success) affLink.value = `${window.location.origin}/register?aff=${data.data}` } catch(e){ Message.error('获取失败') } finally { affLoading.value=false } }
function copyIt(t) { navigator.clipboard?.writeText(t).then(()=>Message.success('已复制')) }
function bindGH() { const id=statusStore.status?.github_client_id; if(id) window.location.href=`https://github.com/login/oauth/authorize?client_id=${id}&scope=user:email` }
function bindLark() { const id=statusStore.status?.lark_client_id; if(id) window.location.href=`https://open.feishu.cn/open-apis/authen/v1/authorize?app_id=${id}&redirect_uri=${encodeURIComponent(window.location.origin+'/oauth/lark')}` }
async function sendCode() { if(!emailForm.email)return; emailSending.value=true; try{await api.get('/api/verification',{params:{email:emailForm.email}});countdown.value=60;const t=setInterval(()=>{countdown.value--;if(countdown.value<=0)clearInterval(t)},1000)}catch(e){Message.error('发送失败')}finally{emailSending.value=false} }
async function submitEmail() { if(!emailForm.email||!emailForm.code)return; emailBinding.value=true; try{const {data}=await api.get('/api/oauth/email/bind',{params:{email:emailForm.email,code:emailForm.code}});if(data.success){Message.success('绑定成功');showEmail.value=false}else Message.error(data.message)}catch(e){Message.error('绑定失败')}finally{emailBinding.value=false} }
async function delAccount() { try{await api.delete('/api/user/self');await authStore.logout();window.location.href='/'}catch(e){Message.error('删除失败')} }
onMounted(() => { loadData() })
</script>

<style scoped>
.setting-container { padding: 4px 0; }
.section h3 { font-size: 16px; font-weight: 600; color: var(--color-text-1); margin-bottom: 20px; padding: 0; }
.setting-form { width: 100%; }
.danger h3 { color: rgb(var(--danger-6)); }
</style>
