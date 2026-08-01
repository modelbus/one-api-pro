<template>
  <div class="log-page">
    <!-- 顶部欢迎条 -->
    <div class="welcome-bar">
      <div class="welcome-text">
        <h1 class="welcome-title">日志</h1>
        <p class="welcome-desc">查看系统调用记录与操作流水</p>
      </div>
      <div class="welcome-meta">
        <a-select v-model="logType" @change="refresh" size="small" style="width: 110px">
          <a-option :value="0" label="全部类型" />
          <a-option :value="1" label="充值" />
          <a-option :value="2" label="消费" />
          <a-option :value="3" label="管理" />
          <a-option :value="4" label="系统" />
          <a-option :value="5" label="测试" />
        </a-select>
      </div>
    </div>

    <!-- 统计 -->
    <div v-if="showStat" class="tip-bar">
      <span class="tip-dot"></span>
      <span class="tip-text">
        总配额 <strong>{{ stat.quota || 0 }}</strong>
        <span class="meta-sep">·</span>
        普通配额 {{ stat.normal_quota || 0 }}
        <span class="meta-sep">·</span>
        订阅配额 {{ stat.subscription_quota || 0 }}
      </span>
    </div>

    <!-- 独立搜索栏 -->
    <div class="search-card">
      <div class="search-left">
        <a-input-search
          v-model="keyword"
          placeholder="搜索日志..."
          allow-clear
          @search="handleSearch"
          @clear="refresh"
          :style="{ width: '260px' }"
        />
        <a-input v-model="filters.token_name" placeholder="令牌名称" size="medium" allow-clear :style="{ width: '160px' }" />
        <a-input v-model="filters.model_name" placeholder="模型名称" size="medium" allow-clear :style="{ width: '160px' }" />
        <a-input v-if="isAdmin" v-model="filters.username" placeholder="用户名" size="medium" allow-clear :style="{ width: '160px' }" />
        <a-input v-if="isAdmin" v-model="filters.channel" placeholder="渠道 ID" size="medium" allow-clear :style="{ width: '140px' }" />
        <a-date-picker v-model="filters.start_timestamp" placeholder="开始时间" size="medium" :style="{ width: '180px' }" show-time />
        <a-date-picker v-model="filters.end_timestamp" placeholder="结束时间" size="medium" :style="{ width: '180px' }" show-time />
      </div>
      <div class="search-right">
        <a-button @click="toggleStat">{{ showStat ? '隐藏统计' : '统计' }}</a-button>
        <a-button type="primary" size="large" @click="refresh" :loading="loading">
          <template #icon><icon-search :size="14" /></template>
          查询
        </a-button>
      </div>
    </div>

    <!-- 列表 -->
    <div class="list-wrap">
      <div v-if="pageItems.length === 0 && !loading && !loadingMore" class="empty-state">
        <div class="empty-icon">
          <icon-file :size="32" />
        </div>
        <p class="empty-title">暂无日志记录</p>
        <p class="empty-desc">尝试调整筛选条件或刷新</p>
      </div>

      <div v-else class="list-body">
        <div class="list-head" :style="{ gridTemplateColumns: headGrid }">
          <div class="col">时间</div>
          <div v-if="isAdmin" class="col">渠道</div>
          <div class="col">来源</div>
          <div class="col">套餐</div>
          <div class="col">类型</div>
          <div class="col">模型</div>
          <div v-if="isAdmin" class="col">用户</div>
          <div class="col">令牌</div>
          <template v-if="logType !== 5">
            <div class="col col-num">Prompt</div>
            <div class="col col-num">Completion</div>
            <div class="col col-num">配额</div>
          </template>
          <div class="col">详情</div>
        </div>

        <a-spin :loading="loading && !loadingMore" style="width: 100%">
          <div v-for="(r, idx) in pageItems" :key="r.id || idx" class="list-row" :style="{ gridTemplateColumns: headGrid }">
            <div class="col">
              <code class="cell-mono clickable" @click="copyId(r.request_id)" :title="r.request_id || ''">
                {{ formatTime(r.created_at) }}
              </code>
            </div>
            <div v-if="isAdmin" class="col">
              <a-tag v-if="r.channel" color="arcoblue" size="small">{{ r.channel_name || r.channel }}</a-tag>
              <span v-else>-</span>
            </div>
            <div class="col">
              <a-tag :color="r.billing_source === 1 ? 'arcoblue' : 'gray'" size="small">
                {{ r.billing_source === 1 ? '订阅' : '额度' }}
              </a-tag>
            </div>
            <div class="col">
              <a-tag v-if="r.plan_name" color="teal" size="small">{{ r.plan_name }}</a-tag>
              <span v-else>-</span>
            </div>
            <div class="col">
              <a-tag :color="typeColor(r.type)" size="small">{{ typeLabel(r.type) }}</a-tag>
            </div>
            <div class="col">
              <span class="cell-strong">{{ r.model_name || '-' }}</span>
            </div>
            <div v-if="isAdmin" class="col">
              <span class="cell-muted">{{ r.username || '-' }}</span>
            </div>
            <div class="col">
              <span class="cell-muted">{{ r.token_name || '-' }}</span>
            </div>
            <template v-if="logType !== 5">
              <div class="col col-num"><span class="cell-num">{{ r.prompt_tokens || 0 }}</span></div>
              <div class="col col-num"><span class="cell-num">{{ r.completion_tokens || 0 }}</span></div>
              <div class="col col-num"><span class="cell-num">{{ r.quota || 0 }}</span></div>
            </template>
            <div class="col">
              <div class="detail-text">{{ r.content }}</div>
              <div class="detail-tags">
                <a-tag v-if="r.elapsed_time" :color="elapsedColor(r.elapsed_time)" size="small">{{ r.elapsed_time }}ms</a-tag>
                <a-tag v-if="r.is_stream" color="pink" size="small">Stream</a-tag>
                <a-tag v-if="r.system_prompt_reset" color="red" size="small">提示词重置</a-tag>
              </div>
            </div>
          </div>
        </a-spin>

        <div v-if="loadingMore" class="load-more-row">
          <a-spin :loading="true" :size="14" />
          <span class="load-more-text">正在加载更多…</span>
        </div>
        <div
          v-else-if="isReachedEnd && logs.length > pageSize && !loading"
          class="load-end-row"
        >
          已显示全部 {{ logs.length }} 条数据
        </div>
      </div>

      <div v-if="logs.length > 0" class="list-footer">
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
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconSearch, IconFile } from '@arco-design/web-vue/es/icon'
import { useAuthStore } from '@/stores/auth'
import api from '@/api'

const authStore = useAuthStore()
const isAdmin = computed(() => authStore.isAdmin)

const logs = ref([])
const loading = ref(false)
const loadingMore = ref(false)
const isReachedEnd = ref(false)
const keyword = ref('')
const logType = ref(0)
const activePage = ref(1)
const pageSize = ref(10)
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

const basePath = computed(() => (isAdmin.value ? '/api/log' : '/api/log/self'))

const headGrid = computed(() => {
  const cols = ['180px']
  if (isAdmin.value) cols.push('120px')
  cols.push('80px', '100px', '80px', '180px')
  if (isAdmin.value) cols.push('110px')
  cols.push('130px')
  if (logType.value !== 5) cols.push('100px', '110px', '100px')
  cols.push('280px')
  return cols.join(' ')
})

const pageItems = computed(() => {
  const start = (activePage.value - 1) * pageSize.value
  return logs.value.slice(start, start + pageSize.value)
})

const totalCountForPager = computed(() => {
  if (isReachedEnd.value) return logs.value.length
  return logs.value.length + pageSize.value
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
  if (!ms) return ''
  if (ms < 1000) return 'green'
  if (ms < 3000) return 'olive'
  if (ms < 5000) return 'yellow'
  if (ms < 10000) return 'orange'
  return 'red'
}

function formatTime(ts) {
  if (!ts) return '-'
  return new Date(ts * 1000).toLocaleString()
}

async function copyId(id) {
  if (!id) return
  try {
    await navigator.clipboard.writeText(id)
    Message.success(`已复制：${id}`)
  } catch {
    Message.warning('复制失败')
  }
}

function getTs(v) {
  return v ? Math.floor(new Date(v).getTime() / 1000) : v
}

async function loadLogs({ append = false, pageIdx = 0 } = {}) {
  if (append) loadingMore.value = true
  else loading.value = true
  try {
    const params = {
      p: pageIdx,
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
      const list = Array.isArray(data.data) ? data.data : (data.data?.items || [])
      if (append) {
        if (list.length === 0) {
          isReachedEnd.value = true
        } else {
          logs.value = [...logs.value, ...list]
          if (list.length < pageSize.value) isReachedEnd.value = true
        }
      } else {
        logs.value = list
        activePage.value = 1
        isReachedEnd.value = list.length < pageSize.value
      }
    }
  } catch {
    /* ignore */
  } finally {
    loading.value = false
    loadingMore.value = false
  }
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
  } catch {
    /* ignore */
  }
}

async function toggleStat() {
  if (!showStat.value) await loadStat()
  showStat.value = !showStat.value
}

function refresh() {
  isReachedEnd.value = false
  loadLogs({ append: false })
}

function onPaginationChange(page) {
  activePage.value = page
  const totalPages = Math.ceil(logs.value.length / pageSize.value)
  // 翻到"最后一页 + 1"且未触底时触发追加加载
  if (page > totalPages && !isReachedEnd.value && !loadingMore.value && !keyword.value) {
    const nextPageIdx = totalPages
    loadLogs({ append: true, pageIdx: nextPageIdx })
  }
}

function onPageSizeChange(s) {
  pageSize.value = s
  activePage.value = 1
}

async function handleSearch() {
  if (!keyword.value) return refresh()
  loading.value = true
  try {
    const { data } = await api.get(`${basePath.value}/search`, { params: { keyword: keyword.value } })
    if (data.success && data.data) {
      logs.value = Array.isArray(data.data) ? data.data : []
      activePage.value = 1
      isReachedEnd.value = true // 搜索结果不支持追加
    }
  } catch {
    /* ignore */
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  refresh()
})
</script>

<style scoped>
.log-page {
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

/* ============ 提示卡 / 统计 ============ */
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
.tip-text strong {
  font-weight: 600;
  color: var(--color-text-1);
  margin: 0 2px;
}
.meta-sep {
  margin: 0 8px;
  color: var(--color-text-4);
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
  flex-wrap: wrap;
}
.search-left {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  flex: 1;
  min-width: 0;
}
.search-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
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
.col-num {
  text-align: right;
}

.cell-mono {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 12px;
  color: var(--color-text-2);
  font-variant-numeric: tabular-nums;
}
.cell-num {
  font-variant-numeric: tabular-nums;
  font-weight: 500;
  color: var(--color-text-1);
}
.cell-strong {
  color: var(--color-text-1);
  font-weight: 500;
}
.cell-muted {
  color: var(--color-text-3);
}
.clickable {
  cursor: pointer;
}
.clickable:hover {
  color: rgb(var(--primary-6));
}

.detail-text {
  font-size: 13px;
  color: var(--color-text-2);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 2px;
}
.detail-tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
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
  margin: 0;
}
</style>
