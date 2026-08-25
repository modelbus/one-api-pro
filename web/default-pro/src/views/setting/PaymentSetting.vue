<template>
  <div class="payment-setting-page">
    <a-spin :loading="loading">
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
import settingApi from '@/api/setting'

const loading = ref(false)

const paymentSettings = reactive({
  wechat: { enabled: false, config: { app_id: '', mch_id: '', api_key: '', notify_url: '', cert_file: '', key_file: '' } },
  alipay: { enabled: false, config: { app_id: '', private_key: '', public_key: '', private_key_file: '', public_key_file: '', notify_url: '', gateway: '' } },
  bank:   { enabled: false, config: { account_name: '', account_no: '', bank_name: '', branch: '', notes: '' } },
})

async function loadPaymentSettings() {
  loading.value = true
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
  } catch (e) {} finally { loading.value = false }
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

onMounted(() => { loadPaymentSettings() })
</script>

<style scoped>
.payment-setting-page { padding: 0 4px; }
.section h3 { font-size: 16px; font-weight: 600; color: var(--color-text-1); margin-bottom: 20px; padding: 0; }
.section-hint { color: var(--color-text-3); font-size: 13px; margin: -10px 0 16px; }
.setting-form { width: 100%; }
</style>
