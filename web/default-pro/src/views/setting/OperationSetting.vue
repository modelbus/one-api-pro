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
          <a-divider :margin="24" />

          <div class="section">
            <h3>套餐运营</h3>
            <a-form layout="vertical" class="setting-form">
              <a-row :gutter="[24, 8]">
                <a-col :span="8">
                  <a-form-item label="升级模式">
                    <a-radio-group v-model="planSettings.upgrade_mode" type="button">
                      <a-radio value="price_diff">差价升级（默认）</a-radio>
                      <a-radio value="stack">叠加</a-radio>
                    </a-radio-group>
                  </a-form-item>
                </a-col>
                <a-col :span="8">
                  <a-form-item label="允许余额充值（仅占位）">
                    <a-switch v-model="planSettings.allow_topup" />
                  </a-form-item>
                </a-col>
              </a-row>
              <a-form-item><a-button type="primary" @click="savePlanSettings">保存套餐运营设置</a-button></a-form-item>
            </a-form>
          </div>
          <a-divider :margin="24" />

          <div class="section">
            <h3>支付</h3>
            <p class="section-hint">开启支付方式后，下方会显示对应表单。关闭后表单自动隐藏。</p>

            <!-- 微信支付 -->
            <a-form layout="vertical" class="setting-form">
              <a-form-item label="微信支付">
                <a-switch v-model="paymentSettings.wechat.enabled" @change="autoSavePayment('wechat', paymentSettings.wechat.enabled)" />
              </a-form-item>
              <template v-if="paymentSettings.wechat.enabled">
                <a-row :gutter="[16, 8]">
                  <a-col :span="12"><a-form-item label="AppID"><a-input v-model="paymentSettings.wechat.config.app_id" placeholder="微信开放平台 AppID" /></a-form-item></a-col>
                  <a-col :span="12"><a-form-item label="商户号 (MchID)"><a-input v-model="paymentSettings.wechat.config.mch_id" /></a-form-item></a-col>
                  <a-col :span="12"><a-form-item label="API Key"><a-input-password v-model="paymentSettings.wechat.config.api_key" placeholder="V2 签名密钥" /></a-form-item></a-col>
                  <a-col :span="12"><a-form-item label="回调 URL (notify_url)"><a-input v-model="paymentSettings.wechat.config.notify_url" placeholder="https://your-host/api/payment/wechat/notify" /></a-form-item></a-col>
                </a-row>
                <a-form-item>
                  <a-upload :custom-request="(opt) => uploadPaymentFile('wechat', opt, 'cert_file')" :show-upload-list="false" accept=".pem">
                    <a-button>上传证书 (apiclient_cert.pem)</a-button>
                  </a-upload>
                  <a-upload :custom-request="(opt) => uploadPaymentFile('wechat', opt, 'key_file')" :show-upload-list="false" accept=".pem" style="margin-left:12px">
                    <a-button>上传私钥 (apiclient_key.pem)</a-button>
                  </a-upload>
                </a-form-item>
                <a-form-item><a-button type="primary" @click="savePaymentMethod('wechat')">保存微信支付</a-button></a-form-item>
              </template>
            </a-form>

            <!-- 支付宝 -->
            <a-form layout="vertical" class="setting-form" style="margin-top:24px">
              <a-form-item label="支付宝">
                <a-switch v-model="paymentSettings.alipay.enabled" @change="autoSavePayment('alipay', paymentSettings.alipay.enabled)" />
              </a-form-item>
              <template v-if="paymentSettings.alipay.enabled">
                <a-row :gutter="[16, 8]">
                  <a-col :span="12"><a-form-item label="AppID"><a-input v-model="paymentSettings.alipay.config.app_id" /></a-form-item></a-col>
                  <a-col :span="12"><a-form-item label="网关 (gateway)"><a-input v-model="paymentSettings.alipay.config.gateway" placeholder="https://openapi.alipay.com/gateway.do" /></a-form-item></a-col>
                  <a-col :span="12"><a-form-item label="回调 URL (notify_url)"><a-input v-model="paymentSettings.alipay.config.notify_url" /></a-form-item></a-col>
                </a-row>
                <a-form-item>
                  <a-upload :custom-request="(opt) => uploadPaymentFile('alipay', opt, 'private_key_file')" :show-upload-list="false" accept=".pem,.txt">
                    <a-button>上传应用私钥</a-button>
                  </a-upload>
                  <a-upload :custom-request="(opt) => uploadPaymentFile('alipay', opt, 'public_key_file')" :show-upload-list="false" accept=".pem,.txt" style="margin-left:12px">
                    <a-button>上传支付宝公钥</a-button>
                  </a-upload>
                </a-form-item>
                <a-form-item><a-button type="primary" @click="savePaymentMethod('alipay')">保存支付宝</a-button></a-form-item>
              </template>
            </a-form>

            <!-- 银行转账 -->
            <a-form layout="vertical" class="setting-form" style="margin-top:24px">
              <a-form-item label="银行转账">
                <a-switch v-model="paymentSettings.bank.enabled" @change="autoSavePayment('bank', paymentSettings.bank.enabled)" />
              </a-form-item>
              <template v-if="paymentSettings.bank.enabled">
                <a-row :gutter="[16, 8]">
                  <a-col :span="12"><a-form-item label="收款人"><a-input v-model="paymentSettings.bank.config.account_name" /></a-form-item></a-col>
                  <a-col :span="12"><a-form-item label="账号"><a-input v-model="paymentSettings.bank.config.account_no" /></a-form-item></a-col>
                  <a-col :span="12"><a-form-item label="开户行"><a-input v-model="paymentSettings.bank.config.bank_name" /></a-form-item></a-col>
                  <a-col :span="12"><a-form-item label="备注"><a-input v-model="paymentSettings.bank.config.notes" /></a-form-item></a-col>
                </a-row>
                <a-form-item><a-button type="primary" @click="savePaymentMethod('bank')">保存银行信息</a-button></a-form-item>
              </template>
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
import settingApi from '@/api/setting'

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

// 套餐运营设置
const planSettings = reactive({ upgrade_mode: 'price_diff', allow_topup: false })
async function loadPlanSettings() {
  try {
    const res = await settingApi.getPlan()
    const data = res?.data?.data
    if (res?.data?.success && data) {
      planSettings.upgrade_mode = data.upgrade_mode || 'price_diff'
      planSettings.allow_topup = !!data.allow_topup
    }
  } catch (e) {}
}
async function savePlanSettings() {
  try {
    const res = await settingApi.putPlan({ upgrade_mode: planSettings.upgrade_mode, allow_topup: planSettings.allow_topup })
    if (res?.data?.success) Message.success('已保存')
    else Message.error(res?.data?.message || '保存失败')
  } catch (e) { Message.error('保存失败') }
}

// 支付设置
const paymentSettings = reactive({
  wechat: { enabled: false, config: { app_id: '', mch_id: '', api_key: '', notify_url: '', cert_file: '', key_file: '' } },
  alipay: { enabled: false, config: { app_id: '', private_key: '', public_key: '', private_key_file: '', public_key_file: '', notify_url: '', gateway: '' } },
  bank:   { enabled: false, config: { account_name: '', account_no: '', bank_name: '', branch: '', notes: '' } },
})
async function loadPaymentSettings() {
  try {
    const res = await settingApi.getPayment()
    if (res?.data?.success && res.data.data) {
      for (const method of ['wechat', 'alipay', 'bank']) {
        const m = res.data.data[method + '_enabled']
        const c = res.data.data[method + '_config']
        if (m) paymentSettings[method].enabled = !!m.enabled
        if (c) paymentSettings[method].config = { ...paymentSettings[method].config, ...(c.config || {}) }
      }
    }
  } catch (e) {}
}
function buildPaymentFormData(method) {
  const fd = new FormData()
  fd.append('config', JSON.stringify({ enabled: paymentSettings[method].enabled, config: paymentSettings[method].config }))
  return fd
}
async function savePaymentMethod(method) {
  try {
    const fd = buildPaymentFormData(method)
    const res = await settingApi.putPaymentMethod(method, fd)
    if (res?.data?.success) Message.success('已保存')
    else Message.error(res?.data?.message || '保存失败')
  } catch (e) { Message.error(e.response?.data?.message || '保存失败') }
}
async function autoSavePayment(method, enabled) {
  // Just persist the enabled flag immediately so the toggle is durable
  // even if the user closes the page before clicking Save.
  try {
    const fd = new FormData()
    paymentSettings[method].config = { ...paymentSettings[method].config }
    fd.append('config', JSON.stringify({ enabled, config: paymentSettings[method].config }))
    await settingApi.putPaymentMethod(method, fd)
  } catch (e) { /* non-fatal */ }
}
async function uploadPaymentFile(method, opt, field) {
  const file = opt.fileItem?.file || opt.file
  if (!file) { Message.error('未选择文件'); return }
  const fd = new FormData()
  fd.append('config', JSON.stringify({ enabled: paymentSettings[method].enabled, config: paymentSettings[method].config }))
  fd.append(field || (method === 'wechat' ? 'cert_file' : 'private_key_file'), file)
  try {
    const { data } = await settingApi.putPaymentMethod(method, fd)
    if (data.success) {
      Message.success('上传成功')
      loadPaymentSettings()
    } else {
      Message.error(data.message || '上传失败')
    }
  } catch (e) { Message.error('上传失败') }
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

onMounted(() => { loadOps(); loadPlanSettings(); loadPaymentSettings() })
</script>

<style scoped>
.section h3 { font-size: 16px; font-weight: 600; color: var(--color-text-1); margin-bottom: 20px; padding: 0; }
.section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; }
.section-header h3 { margin: 0; padding: 0; }
.section-hint { color: var(--color-text-3); font-size: 13px; margin: -10px 0 16px; }
.setting-form { width: 100%; }
</style>
