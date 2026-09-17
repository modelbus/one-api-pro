<template>
  <div class="operation-page">
    <a-spin :loading="loading">
          <div class="section">
            <h3>{{ $t('settingPage.operation.quotaSection') }}</h3>
            <a-form :model="form" layout="vertical" class="setting-form">
              <a-row :gutter="[24, 8]">
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.quotaForNewUser')"><a-input-number v-model="form.QuotaForNewUser" :style="{width:'100%'}" size="large" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.preConsumedQuota')"><a-input-number v-model="form.PreConsumedQuota" :style="{width:'100%'}" size="large" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.quotaForInviter')"><a-input-number v-model="form.QuotaForInviter" :style="{width:'100%'}" size="large" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.quotaForInvitee')"><a-input-number v-model="form.QuotaForInvitee" :style="{width:'100%'}" size="large" /></a-form-item></a-col>
              </a-row>
              <a-form-item><a-button type="primary" @click="saveSection(['QuotaForNewUser','PreConsumedQuota','QuotaForInviter','QuotaForInvitee'])">{{ $t('settingPage.operation.saveQuota') }}</a-button></a-form-item>
            </a-form>
          </div>
          <a-divider :margin="24" />

          <div class="section">
            <h3>{{ $t('settingPage.operation.monitorSection') }}</h3>
            <a-form :model="form" layout="vertical" class="setting-form">
              <a-row :gutter="[24, 8]">
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.channelDisableThreshold')"><a-input-number v-model="form.ChannelDisableThreshold" :style="{width:'100%'}" size="large" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.quotaRemindThreshold')"><a-input-number v-model="form.QuotaRemindThreshold" :style="{width:'100%'}" size="large" /></a-form-item></a-col>
              </a-row>
              <a-row :gutter="[32, 8]">
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.autoDisableChannel')"><a-switch v-model="form.AutomaticDisableChannelEnabled" @change="saveSwitch('AutomaticDisableChannelEnabled')" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.autoEnableChannel')"><a-switch v-model="form.AutomaticEnableChannelEnabled" @change="saveSwitch('AutomaticEnableChannelEnabled')" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.logConsumeEnabled')"><a-switch v-model="form.LogConsumeEnabled" @change="saveSwitch('LogConsumeEnabled')" /></a-form-item></a-col>
              </a-row>
              <a-form-item><a-button type="primary" @click="saveSection(['ChannelDisableThreshold','QuotaRemindThreshold'])">{{ $t('settingPage.operation.saveMonitor') }}</a-button></a-form-item>
            </a-form>
          </div>
          <a-divider :margin="24" />

          <div class="section">
            <h3>{{ $t('settingPage.operation.logCleanSection') }}</h3>
            <a-form layout="vertical" class="setting-form">
              <a-row :gutter="16" align="center">
                <a-col :span="14"><a-form-item :label="$t('settingPage.operation.cleanLogsBefore')"><a-date-picker v-model="logCleanDate" style="width:100%" size="large" /></a-form-item></a-col>
                <a-col style="margin-top:28px"><a-button size="large" @click="cleanLogs" :loading="logCleaning">{{ $t('settingPage.operation.cleanLogs') }}</a-button></a-col>
              </a-row>
            </a-form>
          </div>
          <a-divider :margin="24" />

          <div class="section">
            <h3>{{ $t('settingPage.operation.generalSection') }}</h3>
            <a-form :model="form" layout="vertical" class="setting-form">
              <a-row :gutter="[24, 8]">
                <a-col :span="8"><a-form-item :label="$t('settingPage.operation.topUpLink')"><a-input v-model="form.TopUpLink" :placeholder="$t('settingPage.operation.topUpLinkPlaceholder')" size="large" /></a-form-item></a-col>
                <a-col :span="8"><a-form-item :label="$t('settingPage.operation.chatLink')"><a-input v-model="form.ChatLink" :placeholder="$t('settingPage.operation.chatLinkPlaceholder')" size="large" /></a-form-item></a-col>
                <a-col :span="4"><a-form-item :label="$t('settingPage.operation.quotaPerUnit')"><a-input-number v-model="form.QuotaPerUnit" :style="{width:'100%'}" :precision="2" size="large" /></a-form-item></a-col>
                <a-col :span="4"><a-form-item :label="$t('settingPage.operation.retryTimes')"><a-input-number v-model="form.RetryTimes" :style="{width:'100%'}" size="large" /></a-form-item></a-col>
              </a-row>
              <a-row :gutter="[32, 8]">
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.displayInCurrency')"><a-switch v-model="form.DisplayInCurrencyEnabled" @change="saveSwitch('DisplayInCurrencyEnabled')" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.displayTokenStat')"><a-switch v-model="form.DisplayTokenStatEnabled" @change="saveSwitch('DisplayTokenStatEnabled')" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.approximateToken')"><a-switch v-model="form.ApproximateTokenEnabled" @change="saveSwitch('ApproximateTokenEnabled')" /></a-form-item></a-col>
              </a-row>
            </a-form>
          </div>
          <a-divider :margin="24" />

          <div class="section">
            <h3>{{ $t('settingPage.operation.channelRouteSection') }}</h3>
            <a-form :model="form" layout="vertical" class="setting-form">
              <a-row :gutter="[24, 8]">
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.defaultCooldown')"><a-input-number v-model="form.ChannelDefaultCooldownSeconds" :style="{width:'100%'}" size="large" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.maxCooldown')"><a-input-number v-model="form.ChannelMaxCooldownSeconds" :style="{width:'100%'}" size="large" /></a-form-item></a-col>
              </a-row>
              <a-row :gutter="[32, 8]">
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.channelConcurrency')"><a-switch v-model="form.ChannelConcurrencyEnabled" @change="saveSwitch('ChannelConcurrencyEnabled')" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.stickySession')"><a-switch v-model="form.ChannelStickySessionEnabled" @change="saveSwitch('ChannelStickySessionEnabled')" /></a-form-item></a-col>
              </a-row>
            </a-form>
          </div>
          <a-divider :margin="24" />

          <div class="section">
            <h3>{{ $t('settingPage.operation.errorStrategySection') }}</h3>
            <a-form layout="vertical" class="setting-form">
              <a-row :gutter="[32, 8]">
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.passthrough')"><a-switch v-model="errorNext.passthrough" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.retry')"><a-switch v-model="errorNext.retry" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.disableChannel')"><a-switch v-model="errorNext.disable" /></a-form-item></a-col>
                <a-col :span="6"><a-form-item :label="$t('settingPage.operation.cooldownRetry')"><a-switch v-model="errorNext.cooldown" /></a-form-item></a-col>
              </a-row>
              <a-form-item><a-button type="primary" @click="saveAll">{{ $t('settingPage.operation.saveAll') }}</a-button></a-form-item>
            </a-form>
          </div>
          <a-divider :margin="24" />

          <div class="section">
            <h3>{{ $t('settingPage.operation.planSection') }}</h3>
            <a-form layout="vertical" class="setting-form">
              <a-row :gutter="[24, 8]">
                <a-col :span="8">
                  <a-form-item :label="$t('settingPage.operation.upgradeMode')">
                    <a-radio-group v-model="planSettings.upgrade_mode" type="button">
                      <a-radio value="price_diff">{{ $t('settingPage.operation.upgradePriceDiff') }}</a-radio>
                      <a-radio value="stack">{{ $t('settingPage.operation.upgradeStack') }}</a-radio>
                    </a-radio-group>
                  </a-form-item>
                </a-col>
              </a-row>
              <a-form-item><a-button type="primary" @click="savePlanSettings">{{ $t('settingPage.operation.savePlan') }}</a-button></a-form-item>
            </a-form>
          </div>
        </a-spin>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'
import api from '@/api'
import settingApi from '@/api/setting'

const { t } = useI18n()

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

// 套餐运营设置（v0.0.10 移除 allow_topup，已迁移到 设置-充值）
const planSettings = reactive({ upgrade_mode: 'price_diff' })
async function loadPlanSettings() {
  try {
    const res = await settingApi.getPlan()
    const data = res?.data?.data
    if (res?.data?.success && data) {
      planSettings.upgrade_mode = data.upgrade_mode || 'price_diff'
    }
  } catch (e) {}
}
async function savePlanSettings() {
  try {
    const res = await settingApi.putPlan({ upgrade_mode: planSettings.upgrade_mode })
    if (res?.data?.success) Message.success(t('settingPage.operation.saved'))
    else Message.error(res?.data?.message || t('settingPage.operation.saveFailed'))
  } catch (e) { Message.error(t('settingPage.operation.saveFailed')) }
}

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

async function saveSwitch(key) { try { await api.put('/api/option/', { key, value: form[key] ? 'true' : 'false' }); Message.success(t('settingPage.operation.saved')) } catch(e){ Message.error(t('settingPage.operation.saveFailed')) } }

async function saveSection(keys) {
  for (const k of keys) { try { await api.put('/api/option/', { key: k, value: String(form[k]??'') }) } catch(e){ /* continue */ } }
  Message.success(t('settingPage.operation.saved'))
}

async function cleanLogs() {
  if (!logCleanDate.value) return; logCleaning.value = true
  try { const ts = Math.floor(new Date(logCleanDate.value).getTime()/1000); await api.delete(`/api/log/?target_timestamp=${ts}`); Message.success(t('settingPage.operation.cleaned')) }
  catch(e){ Message.error(t('settingPage.operation.cleanFailed')) } finally { logCleaning.value = false }
}

async function saveAll() {
  await saveSection(['TopUpLink','ChatLink','QuotaPerUnit','RetryTimes','ChannelDefaultCooldownSeconds','ChannelMaxCooldownSeconds'])
  try { await api.put('/api/option/', { key: 'ErrorNext', value: JSON.stringify({...errorNext}) }); Message.success(t('settingPage.operation.allSaved')) } catch(e){ Message.error(t('settingPage.operation.saveFailed')) }
}

onMounted(() => { loadOps(); loadPlanSettings() })
</script>

<style scoped>
.section h3 { font-size: 16px; font-weight: 600; color: var(--color-text-1); margin-bottom: 20px; padding: 0; }
.section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; }
.section-header h3 { margin: 0; padding: 0; }
.setting-form { width: 100%; }
</style>
