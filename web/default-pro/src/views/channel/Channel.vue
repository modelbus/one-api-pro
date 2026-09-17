<template>
  <div class="channel-page">
    <!-- 顶部欢迎条 -->
    <div class="welcome-bar">
      <div class="welcome-text">
        <h1 class="welcome-title">{{ $t('channelPage.title') }}</h1>
        <p class="welcome-desc">{{ $t('channelPage.subtitle') }}</p>
      </div>
      <div class="welcome-meta">
        <span class="meta-chip">{{ $t('channelPage.total', { n: channels.length }) }}</span>
      </div>
    </div>

    <!-- 独立搜索栏 -->
    <div class="search-card">
      <div class="search-left">
        <a-input-search
          v-model="searchKeyword"
          :placeholder="$t('channelPage.searchPlaceholder')"
          allow-clear
          @search="handleSearch"
          @clear="handleSearchClear"
          :style="{ width: '320px' }"
        />
      </div>
      <div class="search-right">
        <a-button @click="handleTestAll" :loading="testingAll">{{ $t('channelPage.testAll') }}</a-button>
        <a-button @click="handleUpdateBalance" :loading="updatingBalance">{{ $t('channelPage.updateBalance') }}</a-button>
        <a-popconfirm :content="$t('channelPage.confirmDeleteDisabled')" @ok="handleDeleteDisabled">
          <a-button status="danger" :loading="deletingDisabled">{{ $t('channelPage.deleteDisabled') }}</a-button>
        </a-popconfirm>
        <a-button type="primary" size="large" @click="openCreateModal">
          <template #icon><icon-plus :size="14" /></template>
          {{ $t('channelPage.addChannel') }}
        </a-button>
      </div>
    </div>

    <!-- 列表 -->
    <div class="list-wrap">
      <div v-if="pageItems.length === 0 && !loading && !loadingMore" class="empty-state">
        <div class="empty-icon">
          <icon-apps :size="32" />
        </div>
        <p class="empty-title">{{ $t('channelPage.emptyTitle') }}</p>
        <p class="empty-desc">{{ $t('channelPage.emptyDesc') }}</p>
        <a-button type="primary" @click="openCreateModal">
          <template #icon><icon-plus :size="14" /></template>
          {{ $t('channelPage.addNow') }}
        </a-button>
      </div>

      <div v-else class="list-body">
        <div class="list-head">
          <div class="col">ID</div>
          <div class="col">{{ $t('channelPage.colType') }}</div>
          <div class="col">{{ $t('channelPage.colName') }}</div>
          <div class="col">Base URL</div>
          <div class="col">{{ $t('channelPage.colModels') }}</div>
          <div class="col">{{ $t('channelPage.colGroups') }}</div>
          <div class="col">{{ $t('channelPage.colStatus') }}</div>
          <div class="col">{{ $t('channelPage.colResponseTime') }}</div>
          <div class="col">{{ $t('channelPage.colFallback') }}</div>
          <div class="col col-action">{{ $t('channelPage.colAction') }}</div>
        </div>

        <a-spin :loading="loading && !loadingMore" style="width: 100%">
          <div v-for="c in pageItems" :key="c.id" class="list-row">
            <div class="col"><span class="cell-mono">#{{ c.id }}</span></div>

            <div class="col">
              <a-tag :color="typeColorMap[c.type] || 'gray'" size="small">
                {{ typeNameMap[c.type] || `Type ${c.type}` }}
              </a-tag>
            </div>

            <div class="col">
              <span class="cell-strong ellipsis" :title="c.name">{{ c.name }}</span>
            </div>

            <div class="col">
              <code class="cell-mono ellipsis" :title="c.base_url">{{ c.base_url || '-' }}</code>
            </div>

            <div class="col">
              <span class="cell-muted ellipsis" :title="c.models">{{ c.models || '-' }}</span>
            </div>

            <div class="col">
              <span class="cell-muted ellipsis" :title="c.groups">{{ c.groups || '-' }}</span>
            </div>

            <div class="col">
              <span class="status-chip" :class="statusClass(c.status)">
                <span class="status-dot"></span>
                {{ statusText(c.status) }}
              </span>
            </div>

            <div class="col">
              <span class="cell-mono">{{ c.response_time ? `${c.response_time}ms` : '-' }}</span>
            </div>

            <div class="col">
              <div v-if="c.is_fallback" class="fallback-cell">
                <a-tooltip :content="$t('channelPage.fallbackHint')">
                  <a-tag color="orangered" size="small" class="fallback-tag">
                    <template #icon><icon-swap :size="12" /></template>
                    {{ $t('channelPage.fallbackTag') }}
                  </a-tag>
                </a-tooltip>
                <span v-if="c.fallback_priority" class="fallback-priority" :title="$t('channelPage.fallbackPriorityHint')">
                  P{{ c.fallback_priority }}
                </span>
              </div>
              <span v-else class="cell-muted">-</span>
            </div>

            <div class="col col-action">
              <a-button type="text" size="small" @click="openEditModal(c)">{{ $t('channelPage.edit') }}</a-button>
              <a-popconfirm
                :content="c.status === 1 ? $t('channelPage.confirmDisable') : $t('channelPage.confirmEnable')"
                @ok="handleToggleStatus(c)"
              >
                <a-button type="text" size="small">
                  {{ c.status === 1 ? $t('channelPage.disable') : $t('channelPage.enable') }}
                </a-button>
              </a-popconfirm>
              <a-button type="text" size="small" :loading="testingIds.includes(c.id)" @click="handleTest(c)">{{ $t('channelPage.test') }}</a-button>
              <a-popconfirm :content="$t('channelPage.confirmDelete')" @ok="handleDelete(c.id)">
                <a-button type="text" size="small" class="danger-btn">{{ $t('channelPage.delete') }}</a-button>
              </a-popconfirm>
            </div>
          </div>
        </a-spin>

        <div v-if="loadingMore" class="load-more-row">
          <a-spin :loading="true" :size="14" />
          <span class="load-more-text">{{ $t('channelPage.loadingMore') }}</span>
        </div>
        <div
          v-else-if="isReachedEnd && channels.length > pageSize && !loading"
          class="load-end-row"
        >
          {{ $t('channelPage.allLoaded', { n: channels.length }) }}
        </div>
      </div>

      <div v-if="channels.length > 0" class="list-footer">
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
      :width="560"
      @ok="handleSubmit"
      @cancel="closeModal"
      :ok-loading="submitting"
      :ok-text="$t('channelPage.save')"
      :cancel-text="$t('channelPage.cancel')"
    >
      <a-form ref="formRef" :model="form" layout="vertical" class="channel-form">
        <a-form-item field="type" :label="$t('channelPage.labelType')" required>
          <a-select v-model="form.type" :placeholder="$t('channelPage.selectType')">
            <a-option v-for="(name, key) in typeNameMap" :key="key" :value="Number(key)" :label="name" />
          </a-select>
        </a-form-item>
        <a-form-item field="name" :label="$t('channelPage.labelName')" required>
          <a-input v-model="form.name" :placeholder="$t('channelPage.channelName')" allow-clear />
        </a-form-item>
        <a-form-item field="groups" :label="$t('channelPage.labelGroups')">
          <a-select
            v-model="form.groups"
            :placeholder="$t('channelPage.selectGroups')"
            multiple
            allow-clear
            allow-search
          >
            <a-option v-for="g in availableGroups" :key="g" :value="g" :label="g" />
          </a-select>
        </a-form-item>
        <a-form-item field="base_url" label="Base URL">
          <a-input v-model="form.base_url" placeholder="https://api.openai.com" allow-clear />
        </a-form-item>
        <a-form-item field="models" :label="$t('channelPage.labelModels')">
          <a-select
            v-model="form.models"
            :placeholder="$t('channelPage.selectModels')"
            multiple
            allow-clear
            allow-search
          >
            <a-option v-for="m in availableModelOptions" :key="m" :value="m" :label="m" />
          </a-select>
        </a-form-item>
        <a-form-item field="model_mapping" :label="$t('channelPage.labelModelMapping')">
          <div class="form-stack">
            <a-textarea
              v-model="form.model_mapping"
              :placeholder="$t('channelPage.modelMappingPlaceholder')"
              :auto-size="{ minRows: 4, maxRows: 8 }"
              allow-clear
            />
            <span class="form-hint form-hint-block">{{ $t('channelPage.modelMappingExample') }}&#123;"gpt-3.5-turbo-0301":"gpt-3.5-turbo","gpt-4-0314":"gpt-4"&#125;</span>
          </div>
        </a-form-item>
        <a-form-item field="key" :label="$t('channelPage.labelKey')" :required="!isEdit">
          <a-input-password v-model="form.key" :placeholder="isEdit ? $t('channelPage.keyKeepHint') : 'API Key'" />
          <span v-if="isEdit" class="form-hint">{{ $t('channelPage.keyKeepHint') }}</span>
        </a-form-item>
        <a-divider :margin="6" />
        <a-form-item field="is_fallback" :label="$t('channelPage.fallback')">
          <a-switch v-model="form.is_fallback" />
          <span class="form-hint">{{ $t('channelPage.fallbackHint') }}</span>
        </a-form-item>
        <a-form-item v-if="form.is_fallback" field="fallback_priority" :label="$t('channelPage.fallbackPriority')">
          <a-input-number v-model="form.fallback_priority" :min="0" :step="1" />
          <span class="form-hint">{{ $t('channelPage.fallbackPriorityHint') }}</span>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Message } from '@arco-design/web-vue'
import { IconPlus, IconApps, IconSwap } from '@arco-design/web-vue/es/icon'
import api from '@/api'

const { t } = useI18n()

// 渠道类型映射（对应后端 relay/registry 的 LegacyType）
const channelTypeIds = [
  1, 50, 14, 33, 3, 11, 24, 51, 28, 41, 40, 15,
  47, 17, 49, 18, 48, 16, 19, 25, 23, 26, 27, 29,
  30, 31, 32, 34, 35, 36, 37, 38, 39, 42, 43, 44,
  45, 46, 8, 22, 21, 20, 2, 5, 7, 10, 4, 6, 9, 12, 13,
]
const typeNameMap = computed(() => {
  const map = {}
  channelTypeIds.forEach((id) => {
    map[id] = t(`channelPage.types.${id}`)
  })
  return map
})
const typeColorMap = {
  1: 'arcoblue', 50: 'green', 14: 'gray', 33: 'gray',
  3: 'cyan', 11: 'orange', 24: 'orange', 51: 'orange',
  28: 'orange', 41: 'purple', 40: 'arcoblue', 15: 'arcoblue',
  47: 'arcoblue', 17: 'orange', 49: 'orange',
  18: 'arcoblue', 48: 'arcoblue', 16: 'purple',
  19: 'arcoblue', 25: 'gray', 23: 'green', 26: 'orange',
  27: 'red', 29: 'orange', 30: 'gray', 31: 'green',
  32: 'arcoblue', 34: 'arcoblue', 35: 'arcoblue', 36: 'gray',
  37: 'orange', 38: 'gray', 39: 'arcoblue', 42: 'arcoblue',
  43: 'arcoblue', 44: 'arcoblue', 45: 'arcoblue', 46: 'arcoblue',
  8: 'pinkpurple', 22: 'arcoblue', 21: 'purple',
  20: 'gray', 2: 'arcoblue', 5: 'gold', 7: 'purple',
  10: 'purple', 4: 'cyan', 6: 'purple', 9: 'gold', 12: 'arcoblue',
  13: 'purple',
}

// 渠道状态（对齐后端 model/channel.go）：1=启用 2=手动禁用 3=自动禁用
function statusText(status) {
  if (status === 1) return t('channelPage.statusEnabled')
  if (status === 2) return t('channelPage.statusDisabled')
  if (status === 3) return t('channelPage.statusAutoDisabled')
  return t('channelPage.statusUnknown')
}
function statusClass(status) {
  if (status === 1) return 'status-on'
  if (status === 3) return 'status-warn'
  return 'status-off'
}

const loading = ref(false)
const loadingMore = ref(false)
const isReachedEnd = ref(false)
const submitting = ref(false)
const channels = ref([])
const searchKeyword = ref('')
const activePage = ref(1)
const pageSize = ref(10)
const modalVisible = ref(false)
const isEdit = ref(false)
const editingId = ref(null)
const testingAll = ref(false)
const updatingBalance = ref(false)
const deletingDisabled = ref(false)
const testingIds = ref([])

const formRef = ref(null)
const form = reactive({
  type: 1,
  name: '',
  base_url: '',
  models: [],
  model_mapping: '',
  groups: [],
  key: '',
  is_fallback: false,
  fallback_priority: 0,
})

const modalTitle = computed(() =>
  isEdit.value ? t('channelPage.editChannel') : t('channelPage.addChannel'),
)

const availableModelOptions = ref([])
const availableGroups = ref([])

const pageItems = computed(() => {
  const start = (activePage.value - 1) * pageSize.value
  return channels.value.slice(start, start + pageSize.value)
})

const totalCountForPager = computed(() => {
  if (isReachedEnd.value) return channels.value.length
  return channels.value.length + pageSize.value
})

async function fetchChannels({ append = false, pageIdx = 0 } = {}) {
  if (append) loadingMore.value = true
  else loading.value = true
  try {
    const params = { p: pageIdx, size: pageSize.value }
    const { data } = await api.get('/api/channel/', { params })
    if (data.success) {
      const list = data.data || []
      if (append) {
        if (list.length === 0) {
          isReachedEnd.value = true
        } else {
          channels.value = [...channels.value, ...list]
          if (list.length < pageSize.value) isReachedEnd.value = true
        }
      } else {
        channels.value = list
        activePage.value = 1
        isReachedEnd.value = list.length < pageSize.value
      }
    } else {
      Message.error(data.message || t('channelPage.loadFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('channelPage.loadFailed'))
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

function handleSearch() {
  if (!searchKeyword.value) { fetchChannels({ append: false }); return }
  loading.value = true
  api.get('/api/channel/search', { params: { keyword: searchKeyword.value } })
    .then(({ data }) => {
      if (data.success) {
        channels.value = data.data || []
        activePage.value = 1
        isReachedEnd.value = true
      }
    })
    .catch((e) => Message.error(e.response?.data?.message || e.message || t('channelPage.searchFailed')))
    .finally(() => { loading.value = false })
}

function handleSearchClear() {
  searchKeyword.value = ''
  isReachedEnd.value = false
  fetchChannels({ append: false })
}

function onPaginationChange(page) {
  activePage.value = page
  const totalPages = Math.ceil(channels.value.length / pageSize.value)
  if (page > totalPages && !isReachedEnd.value && !loadingMore.value && !searchKeyword.value) {
    const nextPageIdx = totalPages
    fetchChannels({ append: true, pageIdx: nextPageIdx })
  }
}

function onPageSizeChange(s) {
  pageSize.value = s
  activePage.value = 1
}

function openCreateModal() {
  isEdit.value = false
  editingId.value = null
  form.type = 1
  form.name = ''
  form.base_url = ''
  form.models = []
  form.model_mapping = ''
  form.groups = []
  form.key = ''
  form.is_fallback = false
  form.fallback_priority = 0
  modalVisible.value = true
}

function openEditModal(record) {
  isEdit.value = true
  editingId.value = record.id
  form.type = record.type || 1
  form.name = record.name || ''
  form.base_url = record.base_url || ''
  form.models = parseModelsField(record.models)
  form.model_mapping = record.model_mapping || ''
  form.groups = parseModelsField(record.groups)
  form.key = record.key || ''
  form.is_fallback = !!record.is_fallback
  form.fallback_priority = record.fallback_priority || 0
  modalVisible.value = true
}

function parseModelsField(val) {
  if (!val) return []
  if (Array.isArray(val)) return val.filter(Boolean)
  return String(val)
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean)
}

function closeModal() {
  modalVisible.value = false
  formRef.value?.clearValidate()
}

async function handleSubmit() {
  const errors = await formRef.value?.validate()
  if (errors) return
  if (form.model_mapping && form.model_mapping.trim()) {
    try {
      const parsed = JSON.parse(form.model_mapping)
      if (typeof parsed !== 'object' || Array.isArray(parsed) || parsed === null) {
        Message.error(t('channelPage.modelMappingMustBeObject'))
        return
      }
    } catch {
      Message.error(t('channelPage.modelMappingInvalidJson'))
      return
    }
  }
  submitting.value = true
  try {
    const payload = {
      type: form.type,
      name: form.name,
      base_url: form.base_url,
      models: Array.isArray(form.models) ? form.models.join(',') : form.models,
      model_mapping: form.model_mapping && form.model_mapping.trim() ? form.model_mapping : '',
      groups: Array.isArray(form.groups) ? form.groups.join(',') : form.groups,
      key: form.key,
      is_fallback: form.is_fallback,
      fallback_priority: form.fallback_priority,
    }
    let res
    if (isEdit.value) {
      payload.id = editingId.value
      res = await api.put('/api/channel/', payload)
    } else {
      res = await api.post('/api/channel/', payload)
    }
    if (res.data.success) {
      Message.success(isEdit.value ? t('channelPage.channelUpdated') : t('channelPage.channelAdded'))
      closeModal()
      fetchChannels()
    } else {
      Message.error(res.data.message || t('channelPage.opFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('channelPage.opFailed'))
  } finally {
    submitting.value = false
  }
}

async function fetchAvailableModels() {
  try {
    const { data } = await api.get('/api/model_price/options')
    const list = Array.isArray(data?.data) ? data.data : []
    const ids = new Set()
    list.forEach((m) => {
      if (typeof m === 'string' && m) ids.add(m)
    })
    if (Array.isArray(form.models)) {
      form.models.forEach((m) => {
        if (m) ids.add(m)
      })
    }
    availableModelOptions.value = Array.from(ids).sort()
  } catch (e) {
    availableModelOptions.value = []
  }
}

async function fetchGroups() {
  try {
    const { data } = await api.get('/api/group/')
    const ids = new Set()
    const list = data?.data || []
    list.forEach((g) => {
      if (g) ids.add(g)
    })
    if (Array.isArray(form.groups)) {
      form.groups.forEach((g) => {
        if (g) ids.add(g)
      })
    }
    availableGroups.value = Array.from(ids).sort()
  } catch (e) {
    availableGroups.value = []
  }
}

async function handleDelete(id) {
  try {
    const { data } = await api.delete(`/api/channel/${id}/`)
    if (data.success) {
      Message.success(t('channelPage.channelDeleted'))
      fetchChannels()
    } else {
      Message.error(data.message || t('channelPage.deleteFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('channelPage.deleteFailed'))
  }
}

async function handleTest(record) {
  testingIds.value.push(record.id)
  try {
    // 单个测试：GET /api/channel/test/:id?model=，返回 time（秒）
    const model = record.models ? record.models.split(',')[0] : ''
    const { data } = await api.get(`/api/channel/test/${record.id}`, { params: { model } })
    if (data.success) {
      Message.success(t('channelPage.testPassed', { time: data.time }))
      fetchChannels()
    } else {
      Message.error(data.message || t('channelPage.testFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('channelPage.testFailed'))
  } finally {
    testingIds.value = testingIds.value.filter((id) => id !== record.id)
  }
}

async function handleTestAll() {
  testingAll.value = true
  try {
    // 全部测试：GET /api/channel/test?scope=all（异步启动）
    const { data } = await api.get('/api/channel/test', { params: { scope: 'all' } })
    if (data.success) {
      Message.info(t('channelPage.testAllStarted'))
    } else {
      Message.error(data.message || t('channelPage.testAllFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('channelPage.testAllFailed'))
  } finally {
    testingAll.value = false
  }
}

async function handleToggleStatus(record) {
  const newStatus = record.status === 1 ? 2 : 1
  try {
    const { data } = await api.put('/api/channel/', { id: record.id, status: newStatus })
    if (data.success) {
      Message.success(newStatus === 1 ? t('channelPage.channelEnabled') : t('channelPage.channelDisabled'))
      record.status = newStatus
    } else {
      Message.error(data.message || t('channelPage.opFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('channelPage.opFailed'))
  }
}

async function handleUpdateBalance() {
  updatingBalance.value = true
  try {
    // 全部更新余额：GET /api/channel/update_balance
    const { data } = await api.get('/api/channel/update_balance')
    if (data.success) {
      Message.success(t('channelPage.balanceUpdateStarted'))
      fetchChannels()
    } else {
      Message.error(data.message || t('channelPage.balanceUpdateFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('channelPage.balanceUpdateFailed'))
  } finally {
    updatingBalance.value = false
  }
}

async function handleDeleteDisabled() {
  deletingDisabled.value = true
  try {
    const { data } = await api.delete('/api/channel/disabled')
    if (data.success) {
      Message.success(t('channelPage.disabledDeleted'))
      fetchChannels()
    } else {
      Message.error(data.message || t('channelPage.deleteFailed'))
    }
  } catch (e) {
    Message.error(e.response?.data?.message || e.message || t('channelPage.deleteFailed'))
  } finally {
    deletingDisabled.value = false
  }
}

onMounted(() => {
  fetchChannels()
  fetchAvailableModels()
  fetchGroups()
})
</script>

<style scoped>
.channel-page {
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

/* ============ 搜索栏 ============ */
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

/* ============ 列表（与 Token 一致） ============ */
.list-wrap {
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-2);
  border-radius: 8px;
  overflow: hidden;
}
.list-body {
  padding: 0;
  overflow-x: auto;
}
.list-head,
.list-row {
  display: grid;
  grid-template-columns: 80px 140px 160px 220px 180px 110px 110px 100px 130px 260px;
  align-items: center;
  padding: 0 20px;
  min-width: max-content;
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
  min-height: 52px;
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
}
.col-action {
  display: flex;
  justify-content: flex-end;
  gap: 0;
}
.col-action :deep(.arco-btn) {
  padding: 0 6px;
}

.cell-mono {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 12px;
  color: var(--color-text-2);
  font-variant-numeric: tabular-nums;
}
.cell-strong {
  color: var(--color-text-1);
  font-weight: 500;
}
.cell-muted {
  color: var(--color-text-3);
}
.ellipsis {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
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
.status-warn {
  background: rgba(245, 63, 63, 0.08);
  color: #f53f3f;
}
.status-warn .status-dot {
  background: #f53f3f;
}

/* ============ 操作列 ============ */
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
.channel-form :deep(.arco-form-item) {
  margin-bottom: 16px;
}
.channel-form :deep(.arco-form-item-label) {
  font-weight: 500;
  font-size: 13px;
  color: var(--color-text-2);
}
.form-hint {
  margin-left: 12px;
  font-size: 12px;
  color: var(--color-text-3);
}
.form-hint-block {
  display: block;
  margin-left: 0;
  margin-top: 8px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
}
.form-stack {
  display: flex;
  flex-direction: column;
  width: 100%;
}

/* ============ 降级渠道标识 ============ */
.fallback-cell {
  display: flex;
  align-items: center;
  gap: 6px;
}
.fallback-tag {
  font-weight: 500;
}
.fallback-priority {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 28px;
  padding: 1px 6px;
  font-size: 11px;
  font-weight: 600;
  color: #d25c1f;
  background: #fff7e8;
  border: 1px solid #ffcca8;
  border-radius: 10px;
  line-height: 1.4;
  cursor: help;
}
</style>
