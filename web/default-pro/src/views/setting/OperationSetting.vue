<template>
  <div class="operation-page">
    <a-spin :loading="loading">
          <div class="section">
            <h3>额度设置</h3>
            <a-form :model="form" layout="vertical" class="setting-form">
              <a-row :gutter="[24, 8]">
                <a-col :span="6"><a-form-item label="新用户初始额度"><a-input-number v-model="form.QuotaForNewUser" :style="{width:'100%'}" size="large" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item label="预扣额度"><a-input-number v-model="form.PreConsumedQuota" :style="{width:'100%'}" size="large" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item label="邀请人奖励额度"><a-input-number v-model="form.QuotaForInviter" :style="{width:'100%'}" size="large" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item label="被邀请人奖励额度"><a-input-number v-model="form.QuotaForInvitee" :style="{width:'100%'}" size="large" /></a-form-item></a-col>
              </a-row>
              <a-form-item><a-button type="primary" @click="saveSection(['QuotaForNewUser','PreConsumedQuota','QuotaForInviter','QuotaForInvitee'])">保存额度设置</a-button></a-form-item>
            </a-form>
          </div>
          <a-divider :margin="24" />

          <div class="section">
            <h3>监控设置</h3>
            <a-form :model="form" layout="vertical" class="setting-form">
              <a-row :gutter="[24, 8]">
                <a-col :span="6"><a-form-item label="渠道禁用响应时间阈值(ms)"><a-input-number v-model="form.ChannelDisableThreshold" :style="{width:'100%'}" size="large" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item label="额度提醒阈值"><a-input-number v-model="form.QuotaRemindThreshold" :style="{width:'100%'}" size="large" /></a-form-item></a-col>
              </a-row>
              <a-row :gutter="[32, 8]">
                <a-col :span="6"><a-form-item label="自动禁用低成功率渠道"><a-switch v-model="form.AutomaticDisableChannelEnabled" @change="saveSwitch('AutomaticDisableChannelEnabled')" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item label="自动启用恢复渠道"><a-switch v-model="form.AutomaticEnableChannelEnabled" @change="saveSwitch('AutomaticEnableChannelEnabled')" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item label="启用消费日志"><a-switch v-model="form.LogConsumeEnabled" @change="saveSwitch('LogConsumeEnabled')" /></a-form-item></a-col>
              </a-row>
              <a-form-item><a-button type="primary" @click="saveSection(['ChannelDisableThreshold','QuotaRemindThreshold'])">保存监控设置</a-button></a-form-item>
            </a-form>
          </div>
          <a-divider :margin="24" />

          <div class="section">
            <h3>日志清理</h3>
            <a-form layout="vertical" class="setting-form">
              <a-row :gutter="16" align="center">
                <a-col :span="14"><a-form-item label="清理指定日期之前的日志"><a-date-picker v-model="logCleanDate" style="width:100%" size="large" /></a-form-item></a-col>
                <a-col style="margin-top:28px"><a-button size="large" @click="cleanLogs" :loading="logCleaning">清理日志</a-button></a-col>
              </a-row>
            </a-form>
          </div>
          <a-divider :margin="24" />

          <div class="section">
            <h3>通用运营</h3>
            <a-form :model="form" layout="vertical" class="setting-form">
              <a-row :gutter="[24, 8]">
                <a-col :span="8"><a-form-item label="充值链接"><a-input v-model="form.TopUpLink" placeholder="充值页面URL" size="large" /></a-form-item></a-col>
                <a-col :span="8"><a-form-item label="Chat链接"><a-input v-model="form.ChatLink" placeholder="Chat页面URL" size="large" /></a-form-item></a-col>
                <a-col :span="4"><a-form-item label="每单位额度价格"><a-input-number v-model="form.QuotaPerUnit" :style="{width:'100%'}" :precision="2" size="large" /></a-form-item></a-col>
                <a-col :span="4"><a-form-item label="重试次数"><a-input-number v-model="form.RetryTimes" :style="{width:'100%'}" size="large" /></a-form-item></a-col>
              </a-row>
              <a-row :gutter="[32, 8]">
                <a-col :span="6"><a-form-item label="按货币显示额度"><a-switch v-model="form.DisplayInCurrencyEnabled" @change="saveSwitch('DisplayInCurrencyEnabled')" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item label="显示Token统计"><a-switch v-model="form.DisplayTokenStatEnabled" @change="saveSwitch('DisplayTokenStatEnabled')" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item label="近似Token计数"><a-switch v-model="form.ApproximateTokenEnabled" @change="saveSwitch('ApproximateTokenEnabled')" /></a-form-item></a-col>
              </a-row>
            </a-form>
          </div>
          <a-divider :margin="24" />

          <div class="section">
            <h3>通道路由</h3>
            <a-form :model="form" layout="vertical" class="setting-form">
              <a-row :gutter="[24, 8]">
                <a-col :span="6"><a-form-item label="默认冷却时间(秒)"><a-input-number v-model="form.ChannelDefaultCooldownSeconds" :style="{width:'100%'}" size="large" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item label="最大冷却时间(秒)"><a-input-number v-model="form.ChannelMaxCooldownSeconds" :style="{width:'100%'}" size="large" /></a-form-item></a-col>
              </a-row>
              <a-row :gutter="[32, 8]">
                <a-col :span="6"><a-form-item label="启用渠道并发限制"><a-switch v-model="form.ChannelConcurrencyEnabled" @change="saveSwitch('ChannelConcurrencyEnabled')" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item label="启用粘性会话"><a-switch v-model="form.ChannelStickySessionEnabled" @change="saveSwitch('ChannelStickySessionEnabled')" /></a-form-item></a-col>
              </a-row>
            </a-form>
          </div>
          <a-divider :margin="24" />

          <div class="section">
            <h3>错误响应策略</h3>
            <a-form layout="vertical" class="setting-form">
              <a-row :gutter="[32, 8]">
                <a-col :span="6"><a-form-item label="透传 (400/404/422)"><a-switch v-model="errorNext.passthrough" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item label="重试 (500/502/503)"><a-switch v-model="errorNext.retry" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item label="禁用渠道 (401/402/403)"><a-switch v-model="errorNext.disable" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item label="冷却+重试 (429/529)"><a-switch v-model="errorNext.cooldown" /></a-form-item></a-col>
              </a-row>
              <a-form-item><a-button type="primary" @click="saveAll">保存全部运营设置</a-button></a-form-item>
            </a-form>
          </div>
        </a-spin>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'
import api from '@/api'

const loading = ref(false), logCleanDate = ref(''), logCleaning = ref(false)
const form = reactive({
  QuotaForNewUser: '', PreConsumedQuota: '', QuotaForInviter: '', QuotaForInvitee: '',
  ChannelDisableThreshold: '', QuotaRemindThreshold: '',
  AutomaticDisableChannelEnabled: false, AutomaticEnableChannelEnabled: false, LogConsumeEnabled: false,
  TopUpLink: '', ChatLink: '', QuotaPerUnit: '', RetryTimes: '',
  DisplayInCurrencyEnabled: false, DisplayTokenStatEnabled: false, ApproximateTokenEnabled: false,
  ChannelDefaultCooldownSeconds: '', ChannelMaxCooldownSeconds: '',
  ChannelConcurrencyEnabled: false, ChannelStickySessionEnabled: false
})
const errorNext = reactive({ passthrough: true, retry: true, disable: true, cooldown: true })

const opKeys = ['QuotaForNewUser','PreConsumedQuota','QuotaForInviter','QuotaForInvitee','ChannelDisableThreshold','QuotaRemindThreshold','AutomaticDisableChannelEnabled','AutomaticEnableChannelEnabled','LogConsumeEnabled','TopUpLink','ChatLink','QuotaPerUnit','RetryTimes','DisplayInCurrencyEnabled','DisplayTokenStatEnabled','ApproximateTokenEnabled','ChannelDefaultCooldownSeconds','ChannelMaxCooldownSeconds','ChannelConcurrencyEnabled','ChannelStickySessionEnabled','ErrorNext']

const numberKeys = ['QuotaForNewUser','PreConsumedQuota','QuotaForInviter','QuotaForInvitee','ChannelDisableThreshold','QuotaRemindThreshold','QuotaPerUnit','RetryTimes','ChannelDefaultCooldownSeconds','ChannelMaxCooldownSeconds']

async function loadOps() {
  loading.value = true
  try {
    const { data } = await api.get('/api/option/')
    if (data.success && data.data) {
      const items = Array.isArray(data.data) ? data.data : Object.entries(data.data).map(([k,v])=>({key:k,value:String(v)}))
      items.forEach(i => {
        if (!opKeys.includes(i.key)) return
        if (i.value === 'true') form[i.key] = true
        else if (i.value === 'false') form[i.key] = false
        else if (numberKeys.includes(i.key)) form[i.key] = Number(i.value) || 0
        else form[i.key] = i.value
      })
    }
  } catch(e){} finally { loading.value = false }
}

async function saveSwitch(key) { try { await api.put('/api/option/', { key, value: form[key] ? 'true' : 'false' }); Message.success('已保存') } catch(e){ Message.error('保存失败') } }

async function saveSection(keys) {
  for (const k of keys) { try { await api.put('/api/option/', { key: k, value: String(form[k]??'') }) } catch(e){ /* continue */ } }
  Message.success('已保存')
}

async function cleanLogs() {
  if (!logCleanDate.value) return; logCleaning.value = true
  try { const ts = Math.floor(new Date(logCleanDate.value).getTime()/1000); await api.delete(`/api/log/?target_timestamp=${ts}`); Message.success('已清理') }
  catch(e){ Message.error('清理失败') } finally { logCleaning.value = false }
}

async function saveAll() {
  await saveSection(['TopUpLink','ChatLink','QuotaPerUnit','RetryTimes','ChannelDefaultCooldownSeconds','ChannelMaxCooldownSeconds'])
  try { await api.put('/api/option/', { key: 'ErrorNext', value: JSON.stringify({...errorNext}) }); Message.success('全部已保存') } catch(e){ Message.error('保存失败') }
}

onMounted(() => { loadOps() })
</script>

<style scoped>
.section h3 { font-size: 16px; font-weight: 600; color: var(--color-text-1); margin-bottom: 20px; padding: 0; }
.section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; }
.section-header h3 { margin: 0; padding: 0; }
.setting-form { width: 100%; }
</style>
