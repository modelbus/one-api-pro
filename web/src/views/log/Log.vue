<template>
  <a-card :bordered="false" class="table-card">
    <div class="action-bar">
      <a-input-search v-model="keyword" placeholder="搜索日志..." @search="handleSearch" @clear="refresh" allow-clear :style="{ width: '220px' }" />
      <div class="action-bar-right">
        <a-button size="small" @click="toggleStat">{{ showStat ? '隐藏统计' : '统计信息' }}</a-button>
      </div>
    </div>

    <div v-if="showStat" class="stats-row">
      <span>总配额：<strong>{{ stat.quota || 0 }}</strong></span>
      <span class="stats-sub">普通配额：{{ stat.normal_quota || 0 }} | 订阅配额：{{ stat.subscription_quota || 0 }}</span>
    </div>

    <div class="filters">
      <a-row :gutter="[12, 8]">
        <a-col :span="4">
          <a-input v-model="filters.token_name" placeholder="令牌名称" size="small" allow-clear />
        </a-col>
        <a-col :span="4">
          <a-input v-model="filters.model_name" placeholder="模型名称" size="small" allow-clear />
        </a-col>
        <a-col :span="4">
          <a-date-picker v-model="filters.start_timestamp" placeholder="开始时间" size="small" style="width:100%" show-time />
        </a-col>
        <a-col :span="4">
          <a-date-picker v-model="filters.end_timestamp" placeholder="结束时间" size="small" style="width:100%" show-time />
        </a-col>
        <a-col :span="2">
          <a-button type="primary" size="small" @click="refresh" :loading="loading">查询</a-button>
        </a-col>
        <a-col :span="2" v-if="isAdmin">
          <a-input v-model="filters.channel" placeholder="渠道ID" size="small" allow-clear />
        </a-col>
        <a-col :span="2" v-if="isAdmin">
          <a-input v-model="filters.username" placeholder="用户名" size="small" allow-clear />
        </a-col>
      </a-row>
    </div>

    <a-table :columns="columns" :data="logs" :loading="loading" :pagination="false" row-key="id" size="medium" :bordered="{ wrapper: true, cell: false }" :scroll="{ x: 1400 }">
      <template #created_at="{ record }">
        <code class="clickable" @click="copyId(record.request_id)">{{ formatTime(record.created_at) }}</code>
      </template>
      <template #channel="{ record }">
        <a-tag v-if="record.channel" color="arcoblue" size="small">{{ record.channel_name || record.channel }}</a-tag>
      </template>
      <template #billing_source="{ record }">
        <a-tag :color="record.billing_source===1?'blue':''" size="small">{{ record.billing_source===1?'订阅':'额度' }}</a-tag>
      </template>
      <template #plan_name="{ record }">
        <a-tag v-if="record.plan_name" color="teal" size="small">{{ record.plan_name }}</a-tag>
        <span v-else>-</span>
      </template>
      <template #type="{ record }">
        <a-tag :color="typeColor(record.type)" size="small">{{ typeLabel(record.type) }}</a-tag>
      </template>
      <template #model_name="{ record }">
        <a-tag color="arcoblue" size="small" v-if="record.model_name">{{ record.model_name }}</a-tag>
        <span v-else>-</span>
      </template>
      <template #username="{ record }">
        <a-tag v-if="record.username" color="arcoblue" size="small">{{ record.username }}</a-tag>
      </template>
      <template #token_name="{ record }">
        <a-tag v-if="record.token_name" size="small">{{ record.token_name }}</a-tag>
        <span v-else>-</span>
      </template>
      <template #detail="{ record }">
        <div>{{ record.content }}</div>
        <a-space size="mini" style="margin-top:4px">
          <a-tag v-if="record.elapsed_time" :color="elapsedColor(record.elapsed_time)" size="small">{{ record.elapsed_time }}ms</a-tag>
          <a-tag v-if="record.is_stream" color="pink" size="small">Stream</a-tag>
          <a-tag v-if="record.system_prompt_reset" color="red" size="small">提示词重置</a-tag>
        </a-space>
      </template>
    </a-table>

    <div class="table-page-footer table-page-footer--between">
      <a-space>
        <a-select v-model="logType" @change="refresh" size="small" style="width:110px">
          <a-option :value="0" label="全部" />
          <a-option :value="1" label="充值" />
          <a-option :value="2" label="消费" />
          <a-option :value="3" label="管理" />
          <a-option :value="4" label="系统" />
          <a-option :value="5" label="测试" />
        </a-select>
        <a-button size="small" @click="refresh" :loading="loading">刷新</a-button>
      </a-space>
      <a-pagination
        :current="page"
        :total="totalLogs"
        :page-size="10"
        size="small"
        show-total
        @change="onPageChange"
      />
    </div>
  </a-card>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'

const authStore = useAuthStore()
const isAdmin = computed(() => authStore.isAdmin)

const logs = ref([])
const loading = ref(false)
const keyword = ref('')
const logType = ref(0)
const page = ref(1)
const totalLogs = ref(0)
const showStat = ref(false)

const stat = reactive({ quota: 0, token: 0, normal_quota: 0, subscription_quota: 0 })

const filters = reactive({
  token_name: '',
  model_name: '',
  start_timestamp: '',
  end_timestamp: '',
  channel: '',
  username: '',
})

const basePath = computed(() => isAdmin.value ? '/api/log' : '/api/log/self')

const columns = computed(() => {
  const cols = [
    { title: '时间', slotName: 'created_at', width: 160 },
    { title: '来源', slotName: 'billing_source', width: 70 },
    { title: '套餐', slotName: 'plan_name', width: 80 },
    { title: '类型', slotName: 'type', width: 70 },
    { title: '模型', slotName: 'model_name', width: 140 },
  ]
  if (isAdmin.value) {
    cols.splice(1, 0, { title: '渠道', slotName: 'channel', width: 100 })
  }
  if (logType.value !== 5) {
    if (isAdmin.value) cols.push({ title: '用户', slotName: 'username', width: 100 })
    cols.push(
      { title: 'Prompt', dataIndex: 'prompt_tokens', width: 80 },
      { title: 'Completion', dataIndex: 'completion_tokens', width: 100 },
      { title: '配额', dataIndex: 'quota', width: 90 },
    )
  }
  cols.push({ title: '详情', slotName: 'detail', width: 200 })
  return cols
})

function typeLabel(t) {
  const m = { 1: '充值', 2: '消费', 3: '管理', 4: '系统', 5: '测试' }
  return m[t] || '未知'
}
function typeColor(t) {
  const m = { 1: 'green', 2: 'olive', 3: 'orange', 4: 'purple', 5: 'violet' }
  return m[t] || ''
}
function elapsedColor(ms) {
  if (!ms) return ''; if (ms < 1000) return 'green'; if (ms < 3000) return 'olive'; if (ms < 5000) return 'yellow'; if (ms < 10000) return 'orange'; return 'red'
}

function formatTime(ts) { if (!ts) return '-'; return new Date(ts * 1000).toLocaleString() }

async function copyId(id) {
  if (!id) return
  try { await navigator.clipboard.writeText(id); Message.success(`已复制：${id}`) }
  catch (e) { Message.warning('复制失败') }
}

function getTs(v) { return v ? Math.floor(new Date(v).getTime() / 1000) : v }

async function loadLogs() {
  loading.value = true
  try {
    const params = {
      p: page.value - 1,
      type: logType.value,
      token_name: filters.token_name,
      model_name: filters.model_name,
      start_timestamp: getTs(filters.start_timestamp) || '',
      end_timestamp: getTs(filters.end_timestamp) || '',
    }
    if (isAdmin.value) {
      params.username = filters.username
      params.channel = filters.channel
    }
    const { data } = await api.get(`${basePath.value}/`, { params })
    if (data.success && data.data) {
      logs.value = Array.isArray(data.data) ? data.data : data.data?.items || []
      totalLogs.value = data.total || logs.value.length
    }
  } catch (e) { /* ignore */ } finally { loading.value = false }
}

async function loadStat() {
  try {
    const params = {
      type: logType.value,
      token_name: filters.token_name,
      model_name: filters.model_name,
      start_timestamp: getTs(filters.start_timestamp) || '',
      end_timestamp: getTs(filters.end_timestamp) || '',
    }
    if (isAdmin.value) {
      params.username = filters.username
      params.channel = filters.channel
    }
    const { data } = await api.get(`${basePath.value}/stat`, { params })
    if (data.success) Object.assign(stat, data.data)
  } catch (e) { /* ignore */ }
}

async function toggleStat() {
  if (!showStat.value) await loadStat()
  showStat.value = !showStat.value
}

function refresh() { page.value = 1; loadLogs() }
function onPageChange(p) { page.value = p; loadLogs() }

async function handleSearch() {
  if (!keyword.value) return refresh()
  loading.value = true
  try {
    const { data } = await api.get(`${basePath.value}/search`, { params: { keyword: keyword.value } })
    if (data.success && data.data) {
      logs.value = Array.isArray(data.data) ? data.data : []
      totalLogs.value = logs.value.length
    }
  } catch (e) { /* ignore */ } finally { loading.value = false }
}

onMounted(() => { refresh() })
</script>

<style scoped>
.table-card { border-radius: 6px; }
.action-bar { display: flex; align-items: center; gap: 12px; padding: 12px 16px; background: var(--color-fill-2); border-radius: 6px; margin-bottom: 15px; }
.action-bar-right { display: flex; align-items: center; gap: 8px; margin-left: auto; }
.table-page-footer { display: flex; justify-content: flex-end; margin-top: 20px; padding-top: 16px; border-top: 1px solid var(--color-border-2); }
.table-page-footer--between { justify-content: space-between; }

.stats-row { font-size: 14px; color: var(--color-text-2); margin-bottom: 12px; }
.stats-row strong { font-weight: 600; color: var(--color-text-1); }
.stats-sub { margin-left: 12px; font-size: 13px; color: var(--color-text-3); }
.filters { margin-bottom: 12px; }
.clickable { cursor: pointer; font-size: 13px; }
.clickable:hover { color: rgb(var(--primary-6)); }
</style>
