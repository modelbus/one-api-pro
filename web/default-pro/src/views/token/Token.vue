<template>
  <div class="token-page">
    <!-- 顶部欢迎条 -->
    <div class="welcome-bar">
      <div class="welcome-text">
        <h1 class="welcome-title">{{ $t('token.title') }}</h1>
        <p class="welcome-desc">{{ $t('token.subtitle') }}</p>
      </div>
      <div class="welcome-meta">
        <span class="meta-chip">{{ $t('token.totalCount', { n: total }) }}</span>
      </div>
    </div>

    <!-- 提示卡 -->
    <div class="tip-bar">
      <span class="tip-dot"></span>
      <span class="tip-text">{{ $t('token.tip') }}</span>
    </div>

    <!-- Base URL 卡片（参考 tbus-web） -->
    <div class="url-card">
      <div class="url-head">
        <span class="url-title">Base URL</span>
        <a-button size="small" type="text" @click="showGuide = true">
          <template #icon><icon-book :size="14" /></template>
          {{ $t('token.guide') }}
        </a-button>
      </div>
      <div class="url-row">
        <code class="url-code">{{ baseUrl }}/v1</code>
        <a-button size="small" type="primary" @click="copyText(`${baseUrl}/v1`, $t('token.baseUrlCopied'))">
          <template #icon><icon-copy :size="14" /></template>
          {{ $t('token.copy') }}
        </a-button>
      </div>
      <p class="url-hint">{{ $t('token.urlHint') }}</p>
    </div>

    <!-- 独立搜索栏 -->
    <div class="search-card">
      <div class="search-left">
        <a-input-search
          v-model="keyword"
          :placeholder="$t('token.searchPlaceholder')"
          allow-clear
          @search="handleSearch"
          @clear="handleSearch"
          :style="{ width: '320px' }"
        />
      </div>
      <div class="search-right">
        <a-button type="primary" size="large" @click="openCreateModal">
          <template #icon><icon-plus :size="14" /></template>
          {{ $t('token.addToken') }}
        </a-button>
      </div>
    </div>

    <!-- 列表 -->
    <div class="list-wrap">
      <div v-if="pageItems.length === 0 && !loading && !loadingMore" class="empty-state">
        <div class="empty-icon">
          <icon-lock :size="32" />
        </div>
        <p class="empty-title">{{ $t('token.createFirst') }}</p>
        <p class="empty-desc">{{ $t('token.createFirstDesc') }}</p>
        <a-button type="primary" @click="openCreateModal">
          <template #icon><icon-plus :size="14" /></template>
          {{ $t('token.createNow') }}
        </a-button>
      </div>

      <div v-else class="list-body">
        <div class="list-head">
          <div class="col col-name">{{ $t('token.colToken') }}</div>
          <div class="col col-key">{{ $t('token.colKey') }}</div>
          <div class="col col-quota">{{ $t('token.colQuota') }}</div>
          <div class="col col-models">{{ $t('token.colModels') }}</div>
          <div class="col col-status">{{ $t('token.colStatus') }}</div>
          <div class="col col-expire">{{ $t('token.colExpire') }}</div>
          <div class="col col-action">{{ $t('token.colAction') }}</div>
        </div>

        <a-spin :loading="loading && !loadingMore" style="width: 100%">
          <div
            v-for="t in pageItems"
            :key="t.id"
            class="list-row"
          >
            <div class="col col-name">
              <div class="name-cell">
                <span class="name-id">#{{ t.id }}</span>
                <span class="name-text">{{ t.name }}</span>
              </div>
            </div>

            <div class="col col-key">
              <div class="key-cell">
                <code class="key-text">{{ maskKey(t.key) }}</code>
                <a-tooltip :content="$t('token.viewFullKey')">
                  <a-button type="text" size="mini" @click="openKeyModal(t)">
                    <template #icon><icon-eye :size="14" /></template>
                  </a-button>
                </a-tooltip>
                <a-tooltip :content="$t('token.copyFullKey')">
                  <a-button type="text" size="mini" @click="copyText(withSkPrefix(t.key), $t('token.fullKeyCopied'))">
                    <template #icon><icon-copy :size="14" /></template>
                  </a-button>
                </a-tooltip>
              </div>
            </div>

            <div class="col col-quota">
              <div class="quota-cell">
                <span :class="t.unlimited_quota ? 'quota-unlimited' : 'quota-value'">
                  {{ t.unlimited_quota ? $t('token.unlimited') : formatQuota(t.remain_quota) }}
                </span>
                <span class="quota-sep">/</span>
                <span class="quota-used">{{ formatQuota(t.used_quota) }}</span>
              </div>
            </div>

            <div class="col col-models">
              <button
                type="button"
                class="models-link"
                @click="openModelsModal(t)"
              >
                <span>{{ formatModelCount(t.models) }}</span>
                <icon-eye :size="12" class="models-link-icon" />
              </button>
            </div>

            <div class="col col-status">
              <span class="status-chip" :class="t.status === 1 ? 'status-on' : 'status-off'">
                <span class="status-dot"></span>
                {{ t.status === 1 ? $t('token.enabled') : $t('token.disabled') }}
              </span>
            </div>

            <div class="col col-expire">
              <span class="cell-mono">{{ formatExpiredTime(t.expired_time, $t('token.neverExpire')) }}</span>
            </div>

            <div class="col col-action">
              <a-button type="text" size="small" @click="openEditModal(t)">{{ $t('common.edit') }}</a-button>
              <a-popconfirm
                :content="t.status === 1 ? $t('token.confirmDisable') : $t('token.confirmEnable')"
                @ok="toggleStatus(t)"
              >
                <a-button type="text" size="small">
                  {{ t.status === 1 ? $t('common.disable') : $t('common.enable') }}
                </a-button>
              </a-popconfirm>
              <a-popconfirm :content="$t('token.confirmDelete')" @ok="handleDelete(t.id)">
                <a-button type="text" size="small" class="danger-btn">{{ $t('common.delete') }}</a-button>
              </a-popconfirm>
            </div>
          </div>
        </a-spin>

        <div v-if="loadingMore" class="load-more-row">
          <a-spin :loading="true" :size="14" />
          <span class="load-more-text">{{ $t('token.loadingMore') }}</span>
        </div>
        <div
          v-else-if="isReachedEnd && tokens.length > pageSize && !loading"
          class="load-end-row"
        >
          {{ $t('token.allLoaded', { n: tokens.length }) }}
        </div>
      </div>

      <div v-if="tokens.length > 0" class="list-footer">
        <a-pagination
          :current="activePage"
          :total="totalCountForPager"
          :page-size="pageSize"
          show-total
          show-page-size
          :page-size-options="[10, 20, 50]"
          size="small"
          @change="onPaginationChange"
          @page-size-change="onPageSizeChange"
        />
      </div>
    </div>

    <!-- 编辑 / 新建弹窗 -->
    <a-modal
      v-model:visible="modalVisible"
      :title="modalTitle"
      :width="520"
      @ok="handleSubmit"
      @cancel="closeModal"
      :ok-loading="submitting"
      :ok-text="$t('token.save')"
      :cancel-text="$t('token.cancel')"
    >
      <a-form ref="formRef" :model="form" :rules="rules" layout="vertical" class="token-form">
        <a-form-item field="name" :label="$t('token.nameLabel')" required>
          <a-input v-model="form.name" :placeholder="$t('token.namePlaceholder')" :max-length="50" allow-clear />
        </a-form-item>

        <a-form-item field="models" :label="$t('token.modelsLabel')" :extra="$t('token.modelsExtra')">
          <a-select
            v-model="form.models"
            :placeholder="$t('token.modelsPlaceholder')"
            multiple
            allow-clear
            allow-search
          >
            <a-option v-for="m in availableModels" :key="m" :value="m" :label="m" />
          </a-select>
        </a-form-item>

        <a-form-item field="subnet" :label="$t('token.subnetLabel')" :extra="$t('token.subnetExtra')">
          <a-input v-model="form.subnet" placeholder="192.168.1.0/24, 10.0.0.0/8" allow-clear />
        </a-form-item>

        <a-form-item field="expired_time" :label="$t('token.expireLabel')">
          <div class="input-with-checkbox">
            <a-date-picker
              v-model="form.expired_time"
              show-time
              format="YYYY-MM-DD HH:mm:ss"
              :placeholder="$t('token.selectExpire')"
              style="flex: 1"
              :disabled="form.never_expire"
              value-format="timestamp"
            >
              <template #suffix-icon></template>
            </a-date-picker>
            <a-checkbox v-model="form.never_expire">{{ $t('token.neverExpire') }}</a-checkbox>
          </div>
          <template #extra>
            <span v-if="form.never_expire">{{ $t('token.neverExpireOn') }}</span>
            <span v-else>{{ $t('token.neverExpireOff') }}</span>
          </template>
        </a-form-item>

        <a-form-item field="remain_quota" :label="$t('token.quotaLabel')">
          <div class="input-with-checkbox">
            <a-input-number
              v-model="form.remain_quota"
              :min="-1"
              :precision="0"
              placeholder="500000"
              :disabled="form.unlimited_quota"
              style="flex: 1"
            />
            <a-checkbox v-model="form.unlimited_quota">{{ $t('token.noQuotaLimit') }}</a-checkbox>
          </div>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 密钥查看弹窗 -->
    <a-modal
      :visible="keyModalVisible"
      @update:visible="(v) => (keyModalVisible = v)"
      @cancel="keyModalVisible = false"
      :width="520"
      :footer="false"
      unmount-on-close
      :title="$t('token.viewKeyTitle')"
    >
      <div v-if="keyModalToken" class="key-modal-body">
        <div class="key-modal-meta">
          <span class="key-modal-label">{{ $t('token.nameLabel') }}</span>
          <span class="key-modal-name">{{ keyModalToken.name }} <span class="key-modal-id">#{{ keyModalToken.id }}</span></span>
        </div>
        <div class="key-modal-value-row">
          <code class="key-modal-value">{{ withSkPrefix(keyModalToken.key) }}</code>
        </div>
        <p class="key-modal-warning">
          <icon-lock :size="14" />
          {{ $t('token.keyWarning') }}
        </p>
        <div class="key-modal-actions">
          <a-button type="primary" @click="handleKeyModalCopy">
            <template #icon><icon-copy :size="14" /></template>
            {{ $t('token.copy') }}
          </a-button>
          <a-button @click="keyModalVisible = false">{{ $t('token.close') }}</a-button>
        </div>
      </div>
    </a-modal>

    <!-- 可用模型弹窗 -->
    <a-modal
      :visible="modelsModalVisible"
      @update:visible="(v) => (modelsModalVisible = v)"
      @cancel="modelsModalVisible = false"
      :width="520"
      :footer="false"
      unmount-on-close
      :title="modelsModalToken ? `${modelsModalToken.name} · ${$t('token.modelsModalTitle')}` : $t('token.modelsModalTitle')"
    >
      <div v-if="modelsModalToken" class="models-modal-body">
        <div class="models-modal-meta">
          <span v-if="currentModalModels.length === 0" class="models-modal-empty">
            {{ $t('token.modelsNoRestrict') }}
          </span>
          <span v-else class="models-modal-count">
            {{ $t('token.modelsCount', { n: currentModalModels.length }) }}
          </span>
        </div>
        <div v-if="currentModalModels.length > 0" class="models-modal-list">
          <div
            v-for="m in currentModalModels"
            :key="m"
            class="models-modal-item"
          >
            <span class="models-modal-dot"></span>
            <code class="models-modal-name">{{ m }}</code>
          </div>
        </div>
        <div class="key-modal-actions">
          <a-button @click="modelsModalVisible = false">{{ $t('token.close') }}</a-button>
        </div>
      </div>
    </a-modal>

    <!-- 使用指南弹窗 -->
    <a-modal
      v-model:visible="showGuide"
      :title="$t('token.guide')"
      :width="640"
      :footer="false"
      unmount-on-close
    >
      <a-tabs v-model:active-key="guideTab" size="medium" class="guide-tabs">
        <a-tab-pane key="quickstart" :title="$t('token.guideQuickstart')">
          <div class="guide-section">
            <div class="guide-row">
              <span class="guide-label">Base URL</span>
              <div class="guide-value">
                <code>{{ baseUrl }}/v1</code>
                <a-button size="mini" @click="copyText(`${baseUrl}/v1`, $t('token.copied'))">{{ $t('token.copy') }}</a-button>
              </div>
            </div>
            <div class="guide-row">
              <span class="guide-label">{{ $t('token.guideApiKey') }}</span>
              <div class="guide-value">
                <a-select v-model="guideKeyId" :placeholder="$t('token.selectToken')" allow-clear style="min-width: 220px">
                  <a-option v-for="k in tokens" :key="k.id" :value="k.id" :label="k.name" />
                </a-select>
              </div>
            </div>
            <div class="guide-row">
              <span class="guide-label">{{ $t('token.guideExample') }}</span>
              <pre class="guide-code"><code>{{ curlExample }}</code></pre>
            </div>
          </div>
        </a-tab-pane>

        <a-tab-pane key="python" title="Python">
          <pre class="guide-code"><code>{{ pythonExample }}</code></pre>
        </a-tab-pane>

        <a-tab-pane key="node" title="Node.js">
          <pre class="guide-code"><code>{{ nodeExample }}</code></pre>
        </a-tab-pane>

        <a-tab-pane key="config" title="Config">
          <pre class="guide-code"><code>{{ configExample }}</code></pre>
        </a-tab-pane>
      </a-tabs>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Message } from '@arco-design/web-vue'
import { IconPlus, IconCopy, IconBook, IconLock, IconEye } from '@arco-design/web-vue/es/icon'
import api from '@/api'
import { useStatusStore } from '@/stores/status'
import { buildTokenExpiredTime, formatExpiredTime } from '@/utils/token'

const statusStore = useStatusStore()
const { t } = useI18n()

const loading = ref(false)
const loadingMore = ref(false)
const isReachedEnd = ref(false)
const submitting = ref(false)
const tokens = ref([])
const keyword = ref('')
const activePage = ref(1)
const pageSize = ref(10)
const modalVisible = ref(false)
const isEdit = ref(false)
const editingId = ref(null)
const availableModels = ref([])

const pageItems = computed(() => {
  const start = (activePage.value - 1) * pageSize.value
  return tokens.value.slice(start, start + pageSize.value)
})

const totalCountForPager = computed(() => {
  if (isReachedEnd.value) return tokens.value.length
  return tokens.value.length + pageSize.value
})

const currentModalModels = computed(() => {
  if (!modelsModalToken.value) return []
  const arr = parseModelArray(modelsModalToken.value.models)
  if (arr.length > 0) return arr
  return availableModels.value || []
})

const showGuide = ref(false)
const guideTab = ref('quickstart')
const guideKeyId = ref(null)

const keyModalVisible = ref(false)
const keyModalToken = ref(null)
const modelsModalVisible = ref(false)
const modelsModalToken = ref(null)

const formRef = ref(null)
const form = reactive({
  name: '',
  models: [],
  subnet: '',
  expired_time: null,
  never_expire: false,
  remain_quota: 500000,
  unlimited_quota: false,
})

const rules = {
  name: [{ required: true, message: t('token.nameRequired') }],
}

const baseUrl = computed(() => {
  const u = statusStore.status?.server_address || ''
  if (!u) return ''
  return u.replace(/\/+$/, '')
})

const modalTitle = computed(() => (isEdit.value ? t('token.editToken') : t('token.addToken')))

const selectedKey = computed(() => tokens.value.find((k) => k.id === guideKeyId.value) || tokens.value[0])

const curlExample = computed(() => {
  const key = selectedKey.value?.key || 'sk-your-key'
  return `curl ${baseUrl.value}/v1/chat/completions \\
  -H "Authorization: Bearer ${key}" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "gpt-4o-mini",
    "messages": [{ "role": "user", "content": "Hello" }]
  }'`
})

const pythonExample = computed(() => {
  const key = selectedKey.value?.key || 'sk-your-key'
  return `from openai import OpenAI

client = OpenAI(
    api_key="${key}",
    base_url="${baseUrl.value}/v1",
)

resp = client.chat.completions.create(
    model="gpt-4o-mini",
    messages=[{"role": "user", "content": "Hello"}],
)
print(resp.choices[0].message.content)`
})

const nodeExample = computed(() => {
  const key = selectedKey.value?.key || 'sk-your-key'
  return `import OpenAI from "openai";

const client = new OpenAI({
  apiKey: "${key}",
  baseURL: "${baseUrl.value}/v1",
});

const resp = await client.chat.completions.create({
  model: "gpt-4o-mini",
  messages: [{ role: "user", content: "Hello" }],
});
console.log(resp.choices[0].message.content);`
})

const configExample = computed(() => {
  const key = selectedKey.value?.key || 'sk-your-key'
  return JSON.stringify(
    {
      baseUrl: `${baseUrl.value}/v1`,
      apiKey: key,
    },
    null,
    2,
  )
})

function withSkPrefix(key) {
  if (!key) return key
  return key.startsWith('sk-') ? key : 'sk-' + key
}

function openKeyModal(token) {
  keyModalToken.value = token
  keyModalVisible.value = true
}

function handleKeyModalCopy() {
  if (!keyModalToken.value) {
    keyModalVisible.value = false
    return
  }
  copyText(withSkPrefix(keyModalToken.value.key), t('token.fullKeyCopied'))
}

function formatModelCount(val) {
  const arr = parseModelArray(val)
  if (arr.length === 0) return t('token.notLimited')
  return t('token.modelCountValue', { n: arr.length })
}

function openModelsModal(token) {
  modelsModalToken.value = token
  modelsModalVisible.value = true
}

function maskKey(key) {
  if (!key) return '-'
  const k = withSkPrefix(key)
  if (k.length <= 8) return k
  return `${k.substring(0, 4)}****${k.substring(k.length - 4)}`
}

function formatQuota(val) {
  if (val == null || val === '') return '-'
  const n = Number(val)
  if (isNaN(n)) return val
  if (n < 0) return t('token.unlimited')
  if (n >= 1000000) return `${(n / 1000000).toFixed(2)}M`
  if (n >= 1000) return `${(n / 1000).toFixed(1)}k`
  return String(n)
}

async function copyText(text, successMsg) {
  try {
    await navigator.clipboard.writeText(text)
    Message.success(successMsg)
  } catch {
    Message.warning(t('token.copyFailed'))
  }
}

async function fetchTokens({ append = false, pageIdx = 0 } = {}) {
  if (append) loadingMore.value = true
  else loading.value = true
  try {
    const params = { p: pageIdx, size: pageSize.value }
    let url = '/api/token/'
    if (keyword.value) {
      url = '/api/token/search'
      params.keyword = keyword.value
    }
    const { data } = await api.get(url, { params })
    if (data.success) {
      const list = data.data || []
      if (append) {
        if (list.length === 0) {
          isReachedEnd.value = true
        } else {
          tokens.value = [...tokens.value, ...list]
          if (list.length < pageSize.value) isReachedEnd.value = true
        }
      } else {
        tokens.value = list
        activePage.value = 1
        isReachedEnd.value = list.length < pageSize.value
      }
    } else {
      Message.error(data.message || t('token.loadFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('token.loadFailed'))
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

async function fetchAvailableModels() {
  try {
    const { data } = await api.get('/api/user/available_models')
    if (data.success) {
      availableModels.value = data.data || []
    }
  } catch {
    /* ignore */
  }
}

function handleSearch() {
  isReachedEnd.value = false
  fetchTokens({ append: false })
}

function onPaginationChange(page) {
  activePage.value = page
  const totalPages = Math.ceil(tokens.value.length / pageSize.value)
  if (page > totalPages && !isReachedEnd.value && !loadingMore.value && !keyword.value) {
    const nextPageIdx = totalPages
    fetchTokens({ append: true, pageIdx: nextPageIdx })
  }
}

function onPageSizeChange(s) {
  pageSize.value = s
  activePage.value = 1
}

function openCreateModal() {
  isEdit.value = false
  editingId.value = null
  form.name = ''
  form.models = []
  form.subnet = ''
  form.expired_time = null
  form.never_expire = false
  form.remain_quota = 500000
  form.unlimited_quota = false
  modalVisible.value = true
}

function openEditModal(record) {
  isEdit.value = true
  editingId.value = record.id
  form.name = record.name || ''
  form.models = parseModelArray(record.models)
  form.subnet = record.subnet || ''
  const exp = record.expired_time ? record.expired_time * 1000 : null
  form.expired_time = exp
  form.never_expire = !exp
  form.remain_quota = record.remain_quota ?? 500000
  form.unlimited_quota = !!record.unlimited_quota
  modalVisible.value = true
}

function parseModelArray(val) {
  if (!val) return []
  if (Array.isArray(val)) return val
  if (typeof val === 'string') return val.split(',').map((s) => s.trim()).filter(Boolean)
  return []
}

function closeModal() {
  modalVisible.value = false
  formRef.value?.resetFields?.()
}

async function handleSubmit() {
  const valid = await formRef.value?.validate()
  if (valid !== undefined) return

  submitting.value = true
  try {
    const payload = {
      name: form.name,
      models: form.models.length ? form.models.join(',') : '',
      subnet: form.subnet,
      expired_time: buildTokenExpiredTime(form),
      remain_quota: form.unlimited_quota ? -1 : form.remain_quota,
      unlimited_quota: form.unlimited_quota,
    }

    if (isEdit.value) {
      payload.id = editingId.value
      const { data } = await api.put('/api/token/', payload)
      if (data.success) {
        Message.success(t('token.tokenUpdated'))
        closeModal()
        fetchTokens()
      } else {
        Message.error(data.message || t('token.updateFailed'))
      }
    } else {
      const { data } = await api.post('/api/token/', payload)
      if (data.success) {
        Message.success(t('token.tokenCreated'))
        closeModal()
        activePage.value = 1
        fetchTokens()
      } else {
        Message.error(data.message || t('token.createFailed'))
      }
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('token.opFailed'))
  } finally {
    submitting.value = false
  }
}

async function toggleStatus(record) {
  const newStatus = record.status === 1 ? 2 : 1
  try {
    const { data } = await api.put(
      '/api/token/',
      { id: record.id, status: newStatus },
      { params: { status_only: true } },
    )
    if (data.success) {
      Message.success(newStatus === 1 ? t('token.enabled') : t('token.disabled'))
      fetchTokens()
    } else {
      Message.error(data.message || t('token.opFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('token.opFailed'))
  }
}

async function handleDelete(id) {
  try {
    const { data } = await api.delete(`/api/token/${id}/`)
    if (data.success) {
      Message.success(t('token.deleteOk'))
      fetchTokens()
    } else {
      Message.error(data.message || t('token.deleteFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('token.deleteFailed'))
  }
}

onMounted(async () => {
  if (!statusStore.loaded) await statusStore.fetchStatus()
  fetchTokens()
  fetchAvailableModels()
})
</script>

<style scoped>
.token-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* ============ 顶部欢迎条 ============ */
.welcome-bar {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  padding: 4px 4px 0;
}
.welcome-title {
  font-size: 24px;
  font-weight: 600;
  color: var(--color-text-1);
  margin: 0 0 4px;
  letter-spacing: -0.2px;
}
.welcome-desc {
  font-size: 13px;
  color: var(--color-text-3);
  margin: 0;
}
.welcome-meta {
  display: flex;
  gap: 6px;
}
.meta-chip {
  font-size: 12px;
  color: var(--color-text-3);
  background: var(--color-fill-2);
  padding: 3px 10px;
  border-radius: 4px;
}

/* ============ 提示卡 ============ */
.tip-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  background: rgba(22, 93, 255, 0.06);
  border: 1px solid rgba(22, 93, 255, 0.12);
  border-radius: 6px;
  font-size: 13px;
  color: var(--color-text-2);
}
.tip-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: rgb(var(--primary-6));
  flex-shrink: 0;
}
.tip-text {
  line-height: 1.6;
}

/* ============ Base URL 卡片（参考 tbus-web） ============ */
.url-card {
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-2);
  border-radius: 8px;
  padding: 16px 20px;
}
.url-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}
.url-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-1);
}
.url-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}
.url-code {
  flex: 1;
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 13px;
  color: var(--color-text-1);
  background: var(--color-fill-1);
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid var(--color-fill-3);
  font-variant-numeric: tabular-nums;
  word-break: break-all;
}
.url-hint {
  font-size: 12px;
  color: var(--color-text-3);
  margin: 0;
  line-height: 1.6;
}

/* ============ 独立搜索栏（白底独立卡） ============ */
.search-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 14px 20px;
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-2);
  border-radius: 8px;
}
.search-left,
.search-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* ============ 列表 ============ */
.list-wrap {
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-2);
  border-radius: 8px;
  overflow: hidden;
}
.list-body {
  padding: 0;
}
.list-head,
.list-row {
  display: grid;
  grid-template-columns: 1.2fr 2fr 1.2fr 1fr 1fr 1.2fr 1.2fr;
  align-items: center;
  padding: 0 20px;
}
.list-head {
  height: 40px;
  background: var(--color-fill-1);
  border-bottom: 1px solid var(--color-fill-3);
  font-size: 12px;
  font-weight: 500;
  color: var(--color-text-3);
}
.list-row {
  min-height: 56px;
  border-bottom: 1px solid var(--color-fill-3);
  transition: background 0.15s;
}
.list-row:last-child {
  border-bottom: none;
}
.list-row:hover {
  background: var(--color-fill-1);
}

/* ============ 单元格 ============ */
.col {
  font-size: 13px;
  color: var(--color-text-2);
  min-width: 0;
  padding-right: 16px;
}
.col:last-child {
  padding-right: 0;
  display: flex;
  justify-content: flex-end;
}

.cell-mono {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 12px;
  color: var(--color-text-3);
  font-variant-numeric: tabular-nums;
}

.name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}
.name-id {
  font-size: 11px;
  color: var(--color-text-4);
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-variant-numeric: tabular-nums;
}
.name-text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--color-text-1);
  font-weight: 500;
}

.key-cell {
  display: flex;
  align-items: center;
  gap: 4px;
  min-width: 0;
}
.key-text {
  flex: 1;
  min-width: 0;
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 12px;
  color: var(--color-text-2);
  background: var(--color-fill-1);
  padding: 4px 8px;
  border-radius: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-variant-numeric: tabular-nums;
  border: 1px solid var(--color-fill-3);
}

.quota-cell {
  display: flex;
  flex-direction: row;
  align-items: baseline;
  gap: 4px;
  white-space: nowrap;
}
.quota-value {
  font-variant-numeric: tabular-nums;
  font-weight: 500;
  color: var(--color-text-1);
  font-size: 13px;
}
.quota-unlimited {
  color: rgb(var(--primary-6));
  font-weight: 500;
  font-size: 13px;
}
.quota-sep {
  color: var(--color-text-4);
  font-size: 12px;
  font-variant-numeric: tabular-nums;
}
.quota-used {
  font-size: 12px;
  color: var(--color-text-3);
  font-variant-numeric: tabular-nums;
}

/* ============ 可用模型链接 ============ */
.models-link {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  background: var(--color-fill-1);
  border: 1px solid var(--color-fill-3);
  border-radius: 4px;
  font-size: 12px;
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-variant-numeric: tabular-nums;
  color: var(--color-text-2);
  cursor: pointer;
  transition: color 0.15s, border-color 0.15s, background 0.15s;
}
.models-link:hover {
  color: rgb(var(--primary-6));
  border-color: rgb(var(--primary-6));
  background: rgba(var(--primary-6), 0.06);
}
.models-link-icon {
  opacity: 0.5;
  transition: opacity 0.15s;
}
.models-link:hover .models-link-icon {
  opacity: 1;
}

/* ============ 状态 chip ============ */
.status-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 2px 10px;
  border-radius: 10px;
  font-size: 12px;
  font-weight: 500;
  width: max-content;
}
.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}
.status-on {
  background: rgba(0, 180, 42, 0.08);
  color: #00b42a;
}
.status-on .status-dot {
  background: #00b42a;
}
.status-off {
  background: var(--color-fill-2);
  color: var(--color-text-3);
}
.status-off .status-dot {
  background: var(--color-text-4);
}

/* ============ 操作列 ============ */
.action-cell,
.col-action :deep(.arco-space) {
  display: flex;
  align-items: center;
  gap: 0;
}
.col-action :deep(.arco-btn) {
  padding: 0 6px;
}
.danger-btn {
  color: var(--color-text-2);
}
.danger-btn:hover {
  color: #f53f3f !important;
  background: rgba(245, 63, 63, 0.06) !important;
}

/* ============ 分页 ============ */
.list-footer {
  display: flex;
  justify-content: flex-end;
  padding: 14px 20px;
  border-top: 1px solid var(--color-fill-3);
}

/* ============ 追加加载 / 末尾提示 ============ */
.load-more-row,
.load-end-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px 20px;
  font-size: 12px;
  color: var(--color-text-3);
  border-top: 1px dashed var(--color-fill-3);
}
.load-more-text {
  color: var(--color-text-3);
}
.load-end-row {
  color: var(--color-text-4);
  background: var(--color-fill-1);
  border-top: 1px solid var(--color-fill-3);
}

/* ============ 空状态 ============ */
.empty-state {
  padding: 80px 20px;
  text-align: center;
}
.empty-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  border-radius: 14px;
  background: var(--color-fill-2);
  color: var(--color-text-3);
  margin-bottom: 12px;
}
.empty-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-1);
  margin: 0 0 4px;
}
.empty-desc {
  font-size: 13px;
  color: var(--color-text-3);
  margin: 0 0 16px;
}

/* ============ 表单 ============ */
.token-form :deep(.arco-form-item) {
  margin-bottom: 16px;
}
.token-form :deep(.arco-form-item-label) {
  font-weight: 500;
  font-size: 13px;
  color: var(--color-text-2);
}

.input-with-checkbox {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}
.input-with-checkbox :deep(.arco-checkbox) {
  white-space: nowrap;
  flex-shrink: 0;
}

/* ============ 密钥查看弹窗 ============ */
.key-modal-body {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.key-modal-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}
.key-modal-label {
  color: var(--color-text-3);
  font-weight: 500;
}
.key-modal-name {
  color: var(--color-text-1);
  font-weight: 500;
}
.key-modal-id {
  font-size: 11px;
  color: var(--color-text-4);
  font-family: 'SF Mono', Menlo, Consolas, monospace;
}
.key-modal-value-row {
  display: flex;
  align-items: center;
  gap: 8px;
}
.key-modal-value {
  flex: 1;
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 13px;
  color: var(--color-text-1);
  background: var(--color-fill-1);
  padding: 10px 14px;
  border-radius: 6px;
  border: 1px solid var(--color-fill-3);
  word-break: break-all;
  user-select: all;
}
.key-modal-warning {
  display: flex;
  align-items: center;
  gap: 6px;
  margin: 0;
  font-size: 12px;
  color: var(--color-text-3);
  background: rgba(22, 93, 255, 0.06);
  padding: 8px 12px;
  border-radius: 4px;
  border: 1px solid rgba(22, 93, 255, 0.12);
}
.key-modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

/* ============ 可用模型弹窗 ============ */
.models-modal-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-height: 60vh;
  overflow: hidden;
}
.models-modal-meta {
  font-size: 12px;
  color: var(--color-text-3);
}
.models-modal-empty {
  color: rgb(var(--primary-6));
}
.models-modal-count {
  font-weight: 500;
}
.models-modal-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 8px;
  overflow-y: auto;
  padding: 4px;
  margin: 0 -4px;
}
.models-modal-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  background: var(--color-fill-1);
  border: 1px solid var(--color-fill-3);
  border-radius: 4px;
  min-width: 0;
}
.models-modal-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: rgb(var(--primary-6));
  flex-shrink: 0;
}
.models-modal-name {
  flex: 1;
  min-width: 0;
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 12px;
  color: var(--color-text-2);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* ============ 使用指南 ============ */
.guide-tabs :deep(.arco-tabs-nav) {
  margin-bottom: 16px;
}
.guide-section {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.guide-row {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.guide-label {
  font-size: 12px;
  font-weight: 500;
  color: var(--color-text-3);
}
.guide-value {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}
.guide-value code {
  flex: 1;
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 12px;
  color: var(--color-text-1);
  background: var(--color-fill-1);
  padding: 6px 10px;
  border-radius: 4px;
  border: 1px solid var(--color-fill-3);
  word-break: break-all;
}
.guide-code {
  margin: 0;
  padding: 14px 16px;
  background: #1d2129;
  color: #c9d1d9;
  border-radius: 6px;
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 12px;
  line-height: 1.7;
  overflow-x: auto;
  white-space: pre;
}
.guide-code code {
  color: inherit;
  background: transparent;
  padding: 0;
  font-family: inherit;
  font-size: inherit;
}

/* ============ 响应式 ============ */
@media (max-width: 1280px) {
  .list-head,
  .list-row {
    grid-template-columns: 1fr 1.4fr 1fr 1fr 1fr 1.2fr;
  }
  .col-expire,
  .col-models {
    display: none;
  }
}
</style>
