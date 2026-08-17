<template>
  <div class="dashboard">
    <!-- 顶部欢迎条 -->
    <div class="welcome-bar">
      <div class="welcome-text">
        <h1 class="welcome-title">控制台</h1>
        <p class="welcome-desc">{{ greeting }}，{{ authStore.user?.display_name || authStore.user?.username || '用户' }} · 今天 {{ todayStr }}</p>
      </div>
      <div class="welcome-meta">
        <span class="meta-chip" v-if="version">{{ version }}</span>
        <span class="meta-chip">{{ roleText }}</span>
      </div>
    </div>

    <a-row :gutter="16">
      <!-- 左侧主内容 -->
      <a-col :xs="24" :lg="16">
        <!-- 核心指标 -->
        <div class="panel stat-grid">
          <div class="stat-item" v-for="(item, idx) in statItems" :key="idx">
            <div class="stat-head">
              <span class="stat-label">{{ item.label }}</span>
              <span class="stat-icon" :style="{ background: item.bg, color: item.color }">
                <component :is="item.icon" :size="14" />
              </span>
            </div>
            <div class="stat-value">{{ item.value }}</div>
            <div class="stat-foot" v-if="item.foot">{{ item.foot }}</div>
          </div>
        </div>

        <!-- 用量进度 -->
        <div class="panel usage-grid" v-if="todayPercent > 0 || monthPercent > 0">
          <div class="usage-item" v-for="u in usageItems" :key="u.key">
            <div class="usage-row">
              <span class="usage-label">{{ u.label }}</span>
              <span class="usage-pct" :style="{ color: progressColor(u.percent) }">{{ u.percent }}%</span>
            </div>
            <a-progress
              :percent="u.percent / 100"
              :show-text="false"
              :color="progressColor(u.percent)"
              :stroke-width="6"
              :track-color="'#f2f3f5'"
            />
            <div class="usage-meta">
              <span>{{ fmtTokens(u.tokens) }} tokens</span>
              <span v-if="u.reset">{{ u.reset }}</span>
            </div>
          </div>
        </div>

        <!-- 趋势图 -->
        <div class="panel">
          <div class="panel-head">
            <h2 class="panel-title">使用趋势（7天）</h2>
          </div>
          <a-row :gutter="12">
            <a-col :xs="24" :sm="8" v-for="t in trendItems" :key="t.field">
              <div class="trend-cell">
                <div class="trend-head">
                  <span class="trend-dot" :style="{ background: t.color }"></span>
                  <span class="trend-label">{{ t.label }}</span>
                  <span class="trend-total">{{ fmtCompact(t.total) }}</span>
                </div>
                <v-chart :option="lineOption(t.field, t.color)" :style="{ height: '160px' }" autoresize />
              </div>
            </a-col>
          </a-row>
        </div>

        <!-- 模型分布（7天） -->
        <div class="panel">
          <div class="panel-head">
            <h2 class="panel-title">模型用量分布（7天）</h2>
            <span class="panel-extra">Top {{ barModels.length || 0 }} · 近 7 天</span>
          </div>
          <v-chart v-if="barModels.length > 0" :option="barOption" :style="{ height: '320px' }" autoresize />
          <div v-else class="chart-empty">暂无数据</div>
        </div>

        <!-- 使用明细 -->
        <div class="panel no-pad">
          <div class="panel-head pad-head">
            <h2 class="panel-title">使用明细</h2>
            <span class="panel-extra">共 {{ details.length }} 条</span>
          </div>
          <a-table
            :columns="columns"
            :data="details"
            :pagination="{ pageSize: 8, showTotal: false, showJumper: false }"
            :bordered="false"
            :stripe="true"
            size="medium"
            class="dash-table"
          >
            <template #model="{ record }">
              <div class="model-cell">
                <ModelIcon :slug="record.provider_slug" :name="record.model_name" :size="20" />
                <span class="model-name">{{ record.model_name }}</span>
              </div>
            </template>
            <template #date="{ record }">
              <span class="cell-mono">{{ record.day }}</span>
            </template>
            <template #requests="{ record }">
              <span class="cell-num">{{ fmtNum(record.request_count) }}</span>
            </template>
            <template #quota="{ record }">
              <span class="cell-num">{{ fmtQuota(record.quota) }}</span>
            </template>
            <template #tokens="{ record }">
              <span class="cell-num">{{ fmtTokens(record.prompt_tokens + record.completion_tokens) }}</span>
            </template>
            <template #empty>
              <div class="empty-state">暂无使用记录</div>
            </template>
          </a-table>
        </div>

        <!-- 快捷操作 -->
        <div class="panel">
          <div class="panel-head">
            <h2 class="panel-title">快捷操作</h2>
          </div>
          <div class="quick-grid">
            <button
              v-for="q in quickActions"
              :key="q.path"
              class="quick-btn"
              @click="$router.push(q.path)"
            >
              <span class="quick-icon" :style="{ background: q.bg, color: q.color }">
                <component :is="q.icon" :size="16" />
              </span>
              <span class="quick-label">{{ q.label }}</span>
            </button>
          </div>
        </div>
      </a-col>

      <!-- 右侧栏 -->
      <a-col :xs="24" :lg="8">
        <!-- 余额 -->
        <div class="panel balance-card">
          <div class="balance-head">
            <div class="balance-icon" :style="{ background: 'linear-gradient(135deg, #165dff, #4080ff)' }">
              <icon-archive :size="18" />
            </div>
            <div class="balance-title">
              <span class="balance-label">我的余额</span>
              <span class="balance-unit">额度</span>
            </div>
          </div>
          <div class="balance-value">
            <span class="balance-num">{{ fmtQuota(balance) }}</span>
          </div>
          <div class="balance-actions">
            <a class="balance-btn" @click="goTo('/topup')">充值</a>
            <a class="balance-btn" @click="goTo('/subscription')">订阅</a>
            <a class="balance-btn" @click="goTo('/redemption')">兑换</a>
          </div>
        </div>

        <!-- API Key -->
        <div class="panel" :class="{ 'is-active': isActiveRoute('/token') }">
          <div class="panel-head">
            <h2 class="panel-title">API Key</h2>
            <a class="panel-link" @click="$router.push('/token')">管理 ({{ tokenTotal }}) →</a>
          </div>
          <template v-if="latestApiKey">
            <div class="apikey-box">
              <code class="apikey-masked">{{ maskedApiKey }}</code>
              <a-button type="text" size="mini" @click="copyLatestKey">
                <template #icon><icon-copy :size="14" /></template>
                复制
              </a-button>
            </div>
            <div class="apikey-meta">
              <span class="meta-name">{{ latestApiKey.name || '默认密钥' }}</span>
              <span v-if="latestApiKey.created_time">{{ formatDate(latestApiKey.created_time) }}</span>
            </div>
          </template>
          <div v-else class="apikey-empty">
            <p class="empty-title">尚未创建 API Key</p>
            <p class="empty-desc">创建后即可开始调用 API</p>
            <a-button type="primary" size="small" @click="$router.push('/token')">立即创建</a-button>
          </div>
        </div>

        <!-- 公告 -->
        <div class="panel clickable" @click="$router.push('/log')">
          <div class="panel-head">
            <h2 class="panel-title">系统公告</h2>
            <span class="panel-link">查看日志 →</span>
          </div>
          <div v-if="notice" class="notice-body" v-html="notice"></div>
          <div v-else class="notice-list">
            <a class="notice-item" href="#">
              <span class="notice-title">One Api Pro 企业版正式发布</span>
              <span class="notice-time">2026-08-01</span>
            </a>
            <a class="notice-item" href="#">
              <span class="notice-title">支持 40+ 模型平台统一接入</span>
              <span class="notice-time">2026-07-25</span>
            </a>
            <a class="notice-item" href="#">
              <span class="notice-title">新增订阅套餐与用量管控</span>
              <span class="notice-time">2026-07-15</span>
            </a>
          </div>
        </div>

        <!-- 更新日志 -->
        <div class="panel">
          <div class="panel-head">
            <h2 class="panel-title">更新日志</h2>
          </div>
          <div class="changelog">
            <div class="cl-item">
              <div class="cl-head">
                <span class="cl-version">{{ version || 'v1.0.0' }}</span>
                <span class="cl-date">{{ todayStr }}</span>
              </div>
              <p class="cl-desc">全新仪表盘：核心指标 + 趋势图表 + 模型分布</p>
            </div>
            <div class="cl-item">
              <div class="cl-head">
                <span class="cl-version">v1.0.0</span>
                <span class="cl-date">2026-06-01</span>
              </div>
              <p class="cl-desc">One Api Pro 首个稳定版发布</p>
            </div>
          </div>
        </div>

        <!-- 联系 -->
        <div class="panel">
          <div class="panel-head">
            <h2 class="panel-title">资源</h2>
          </div>
          <div class="contact-list">
            <a class="contact-item" href="http://one-api.pro" target="_blank">
              <span class="contact-label">官方文档</span>
              <span class="contact-value is-active">one-api.pro</span>
              <icon-launch class="contact-arrow" :size="14" />
            </a>
            <a class="contact-item" href="https://github.com/Leon-PanPan/one-api-pro" target="_blank">
              <span class="contact-label">GitHub</span>
              <span class="contact-value is-active">Leon-PanPan/one-api-pro</span>
              <icon-launch class="contact-arrow" :size="14" />
            </a>
          </div>
        </div>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import {
  IconCodeSquare, IconSend, IconArchive, IconCalendar, IconCode, IconFile,
  IconCopy, IconLaunch,
} from '@arco-design/web-vue/es/icon'
import VChart from 'vue-echarts'
import 'echarts'
import api from '@/api'
import { useAuthStore } from '@/stores/auth'
import { useStatusStore } from '@/stores/status'
import ModelIcon from '@/components/ModelIcon.vue'
import { findProviderByName } from '@/constants/providers'

const authStore = useAuthStore()
const statusStore = useStatusStore()
const route = useRoute()

function isActiveRoute(prefix) {
  return route.path === prefix || route.path.startsWith(prefix + '/')
}

const loading = ref(false)
const notice = ref('')
const version = ref('')
const latestApiKey = ref(null)
const tokenTotal = ref(0)
const balance = ref(0)

const statData = reactive({ total_tokens: 0, total_requests: 0, total_quota: 0 })
const dailyUsage = ref([])
const logStat = reactive({
  quota: 0, token: 0, normal_quota: 0, subscription_quota: 0,
  prompt_tokens: 0, completion_tokens: 0,
  total_users: 0, active_channels: 0, request_count: 0,
})

const currentPlan = ref({ id: 0, name: '-', expireDate: '-' })
const planExpired = ref(false)
const todayPercent = ref(0)
const monthPercent = ref(0)
const todayTokens = ref(0)
const sevendayTokens = ref(0)
const todayResetTime = ref('')
const sevendayResetTime = ref('')

const chartData = ref({ requests: [], quota: [], tokens: [] })
const details = ref([])
const barModels = ref([])
const modelColors = ['#165dff', '#00b42a', '#ff7d00', '#f53f3f', '#722ed1', '#0fc6c2']

const maskedApiKey = computed(() => {
  if (!latestApiKey.value) return ''
  const raw = latestApiKey.value.key || latestApiKey.value.show_key || ''
  if (!raw) return ''
  const key = raw.startsWith('sk-') ? raw : 'sk-' + raw
  if (key.length <= 12) return key.slice(0, 6) + '****'
  return key.slice(0, 5) + '****' + key.slice(-4)
})

const greeting = computed(() => {
  const h = new Date().getHours()
  if (h < 6) return '夜深了'
  if (h < 12) return '早上好'
  if (h < 14) return '中午好'
  if (h < 18) return '下午好'
  return '晚上好'
})

const todayStr = computed(() => {
  const d = new Date()
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
})

const roleText = computed(() => {
  if (authStore.isRoot) return '超级管理员'
  if (authStore.isAdmin) return '管理员'
  return '普通用户'
})

const statItems = computed(() => {
  const planFoot = (() => {
    if (currentPlan.value.id > 0 && !planExpired.value) {
      return currentPlan.value.expireDate && currentPlan.value.expireDate !== '-'
        ? `${currentPlan.value.expireDate} 到期`
        : '当前生效中'
    }
    if (planExpired.value) return '已过期 · 请续费'
    return null
  })()

  return [
    {
      label: '总 Token',
      value: fmtNum(statData.total_tokens),
      foot: '累计消耗',
      icon: IconCodeSquare,
      bg: 'rgba(22,93,255,0.08)',
      color: '#165dff',
    },
    {
      label: '总请求',
      value: fmtNum(statData.total_requests),
      foot: '累计调用',
      icon: IconSend,
      bg: 'rgba(0,180,42,0.08)',
      color: '#00b42a',
    },
    {
      label: '总配额',
      value: fmtQuota(statData.total_quota),
      foot: '累计消耗',
      icon: IconArchive,
      bg: 'rgba(255,125,0,0.08)',
      color: '#ff7d00',
    },
    {
      label: '当前套餐',
      value: currentPlan.value.id > 0 ? currentPlan.value.name : (planExpired.value ? '已过期' : '未订阅'),
      foot: planFoot,
      icon: IconCalendar,
      bg: 'rgba(114,46,209,0.08)',
      color: '#722ed1',
    },
  ]
})

const usageItems = computed(() => [
  { key: 'today', label: '今日用量', percent: todayPercent.value, tokens: todayTokens.value, reset: todayResetTime.value ? `重置 ${todayResetTime.value}` : '' },
  { key: 'week', label: '7 日用量', percent: monthPercent.value, tokens: sevendayTokens.value, reset: sevendayResetTime.value ? `重置 ${sevendayResetTime.value}` : '' },
])

const trendItems = computed(() => [
  { field: 'requests', label: '请求数', color: '#165dff', total: sumField('requests') },
  { field: 'quota', label: '额度', color: '#00b42a', total: sumField('quota') },
  { field: 'tokens', label: 'Token', color: '#ff7d00', total: sumField('tokens') },
])

function sumField(field) {
  return (chartData.value[field] || []).reduce((s, d) => s + (d.value || 0), 0)
}

const quickActions = computed(() => [
  { label: '管理令牌', path: '/token', icon: IconCode, bg: 'rgba(22,93,255,0.08)', color: '#165dff' },
  { label: '充值额度', path: '/topup', icon: IconArchive, bg: 'rgba(255,125,0,0.08)', color: '#ff7d00' },
  { label: '使用日志', path: '/log', icon: IconFile, bg: 'rgba(15,198,194,0.08)', color: '#0fc6c2' },
])

const columns = [
  { title: '模型', dataIndex: 'model_name', slotName: 'model', width: 220, ellipsis: true, tooltip: true },
  { title: '日期', dataIndex: 'day', slotName: 'date', width: 120 },
  { title: '请求数', dataIndex: 'request_count', slotName: 'requests', width: 110, align: 'right' },
  { title: '额度消耗', dataIndex: 'quota', slotName: 'quota', width: 120, align: 'right' },
  { title: 'Token 用量', slotName: 'tokens', width: 130, align: 'right' },
]

function fmtNum(n) { return Number(n || 0).toLocaleString() }
function fmtQuota(n) {
  const v = Number(n) || 0
  if (v >= 10000) return (v / 10000).toFixed(2) + 'w'
  return v.toLocaleString(undefined, { minimumFractionDigits: 1, maximumFractionDigits: 1 })
}
function fmtTokens(n) {
  const v = Number(n) || 0
  if (v >= 1000000) return (v / 1000000).toFixed(1) + 'M'
  if (v >= 1000) return (v / 1000).toFixed(1) + 'K'
  return v.toLocaleString()
}
function fmtCompact(n) {
  const v = Number(n) || 0
  if (v >= 1000000) return (v / 1000000).toFixed(1) + 'M'
  if (v >= 1000) return (v / 1000).toFixed(1) + 'K'
  return String(v)
}
function progressColor(p) {
  if (p >= 80) return '#f53f3f'
  if (p >= 50) return '#ff7d00'
  return '#165dff'
}
function formatDate(ts) {
  if (!ts) return ''
  const d = typeof ts === 'number' ? new Date(ts * 1000) : new Date(ts)
  if (isNaN(d.getTime())) return ''
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
}

function lineOption(field, color) {
  const data = chartData.value[field] || []
  return {
    grid: { left: 32, right: 8, top: 12, bottom: 22 },
    tooltip: { trigger: 'axis', confine: true, backgroundColor: '#1d2129', borderWidth: 0, textStyle: { color: '#fff', fontSize: 12 } },
    xAxis: {
      type: 'category', boundaryGap: false,
      data: data.map((d) => (d.date || '').slice(5)),
      axisLabel: { fontSize: 11, color: '#86909c' },
      axisLine: { show: false },
      axisTick: { show: false },
    },
    yAxis: {
      type: 'value',
      axisLabel: { fontSize: 11, color: '#86909c', formatter: (v) => (v >= 1000 ? `${(v / 1000).toFixed(0)}k` : v) },
      splitLine: { lineStyle: { color: '#f2f3f5', type: 'dashed' } },
      axisLine: { show: false },
      axisTick: { show: false },
    },
    series: [{
      type: 'line', smooth: true, symbol: 'circle', symbolSize: 5,
      data: data.map((d) => d.value),
      lineStyle: { color, width: 2 },
      itemStyle: { color },
      areaStyle: {
        color: {
          type: 'linear', x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [{ offset: 0, color: color + '35' }, { offset: 1, color: color + '02' }],
        },
      },
    }],
  }
}

const barOption = computed(() => {
  // 参考 web-back StatisticalBarChart：x 轴 = 7 天，每模型一个 series，按天聚合 quota，堆叠展示
  const fmtDay = (d) => {
    const y = d.getFullYear()
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    return `${y}-${m}-${day}`
  }
  const days = []
  const today = new Date()
  for (let i = 6; i >= 0; i--) {
    const d = new Date(today)
    d.setDate(d.getDate() - i)
    days.push(fmtDay(d))
  }

  const grouped = {}
  for (const item of details.value) {
    const day = item.day
    if (!day) continue
    if (!grouped[day]) grouped[day] = {}
    const model = item.model_name
    if (!grouped[day][model]) grouped[day][model] = 0
    grouped[day][model] += item.quota || 0
  }

  const models = barModels.value.slice(0, 8)
  const series = models.map((m, idx) => ({
    name: m,
    type: 'bar',
    stack: 'total',
    barMaxWidth: 36,
    itemStyle: { color: modelColors[idx % modelColors.length], borderRadius: [3, 3, 0, 0] },
    emphasis: { focus: 'series' },
    data: days.map((d) => (grouped[d] && grouped[d][m]) || 0),
  }))

  return {
    color: modelColors,
    grid: { left: 48, right: 16, top: 16, bottom: 36 },
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' }, backgroundColor: '#1d2129', borderWidth: 0, textStyle: { color: '#fff', fontSize: 12 } },
    legend: {
      show: models.length > 0,
      bottom: 0,
      left: 'center',
      type: 'scroll',
      icon: 'roundRect',
      itemWidth: 12,
      itemHeight: 12,
      textStyle: { fontSize: 11, color: '#86909c' },
    },
    xAxis: {
      type: 'category', data: days.map((d) => d.slice(5)),
      axisLabel: { fontSize: 11, color: '#86909c' },
      axisLine: { show: false },
      axisTick: { show: false },
    },
    yAxis: {
      type: 'value',
      axisLabel: { fontSize: 11, color: '#86909c', formatter: (v) => (v >= 1000 ? `${(v / 1000).toFixed(0)}k` : v) },
      splitLine: { lineStyle: { color: '#f2f3f5', type: 'dashed' } },
      axisLine: { show: false },
      axisTick: { show: false },
    },
    series,
  }
})

// 兜底：dashboard 接口无数据时，拉取最近消费日志聚合图表数据
async function loadLogsFallback() {
  try {
    const basePath = authStore.isAdmin ? '/api/log' : '/api/log/self'
    const startTs = sevenDayRange.value.start
    const endTs = sevenDayRange.value.end
    const collected = []
    // 拉前 5 页（每页 10 条，共 50 条），只取近 7 天的消费日志
    for (let p = 0; p < 5; p++) {
      const { data } = await api.get(`${basePath}/`, {
        params: { p, type: 2, start_timestamp: startTs, end_timestamp: endTs },
      }).catch(() => ({ data: { success: false } }))
      if (!data.success || !Array.isArray(data.data) || data.data.length === 0) break
      collected.push(...data.data)
      if (data.data.length < 10) break
    }
    // 转成 buildCharts 需要的 { day, model_name, quota, prompt_tokens, completion_tokens, request_count }
    const stats = collected
      .filter((l) => l.created_at && l.created_at >= startTs && l.created_at <= endTs)
      .map((l) => {
        const d = new Date(l.created_at * 1000)
        const pad = (n) => String(n).padStart(2, '0')
        return {
          day: `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`,
          model_name: l.model_name || '其他',
          quota: l.quota || 0,
          prompt_tokens: l.prompt_tokens || 0,
          completion_tokens: l.completion_tokens || 0,
          request_count: 1,
        }
      })
    buildCharts(stats)
  } catch (e) {
    buildCharts([])
  }
}

async function loadApiKey() {
  // 通过令牌列表接口获取所有令牌（id 倒序，第一条为最新），展示第一条
  try {
    const { data } = await api.get('/api/token/', { params: { p: 0 } })
    if (data.success && Array.isArray(data.data)) {
      const tokens = data.data
      latestApiKey.value = tokens.length > 0 ? tokens[0] : null
      // 再用 search 接口拿全部数量（列表接口有 size 限制）
      try {
        const { data: searchRes } = await api.get('/api/token/search', { params: { keyword: '' } })
        if (searchRes.success && Array.isArray(searchRes.data)) {
          tokenTotal.value = searchRes.data.length
        }
      } catch { /* ignore */ }
    }
  } catch (e) { /* ignore */ }
}

function copyLatestKey() {
  const raw = latestApiKey.value?.key || latestApiKey.value?.show_key
  if (raw) {
    const key = raw.startsWith('sk-') ? raw : 'sk-' + raw
    navigator.clipboard?.writeText(key)
      .then(() => Message.success('API Key 已复制到剪贴板'))
      .catch(() => Message.error('复制失败，请手动复制'))
  }
}

// 充值/订阅入口（保留入口，暂不跳转）
function goTo() {
  /* 占位，暂无跳转 */
}

// 获取当前用户余额（/api/user/self 返回 quota）
async function loadSelf() {
  try {
    const { data } = await api.get('/api/user/self')
    if (data.success && data.data) {
      balance.value = data.data.quota || 0
    }
  } catch (e) { /* ignore */ }
}

async function loadSubscription() {
  try {
    const { data } = await api.get('/api/subscription/self')
    if (data.success && data.data) {
      const p = data.data
      if (Array.isArray(p) && p.length > 0) {
        const first = p[0]
        currentPlan.value = {
          id: first.id || 1,
          name: first.plan_name || first.name || '套餐',
          expireDate: first.end_time ? formatDate(first.end_time) : '-',
        }
        planExpired.value = !!(first.end_time && new Date(first.end_time) < new Date())
        todayPercent.value = Number(first.used_percent || 0)
        monthPercent.value = Number(first.week_percent || 0)
        todayTokens.value = Number(first.today_tokens || 0)
        sevendayTokens.value = Number(first.week_tokens || 0)
        todayResetTime.value = first.today_reset || ''
        sevendayResetTime.value = first.week_reset || ''
      }
    }
  } catch (e) { /* ignore */ }
}

// 计算「近 7 天」时间范围（对齐后端 SearchLogsByDayAndModel 语义）
const sevenDayRange = computed(() => {
  const now = new Date()
  const end = Math.floor(now.getTime() / 1000)
  const startOfDay = Math.floor(new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime() / 1000)
  const start = startOfDay - 6 * 86400 // 包含今天共 7 天
  return { start, end }
})

async function loadDashboard() {
  loading.value = true
  try {
    const { start, end } = sevenDayRange.value
    const statEndpoint = authStore.isAdmin ? '/api/log/stat' : '/api/log/self/stat'
    const requests = [
      api.get('/api/user/dashboard', { params: { start, end } }).catch(() => ({ data: { success: true, data: [] } })),
      api.get(statEndpoint, {
        params: {
          type: 0,
          username: '',
          token_name: '',
          model_name: '',
          start_timestamp: 0,
          end_timestamp: end,
          channel: 0,
        },
      }).catch(() => ({ data: { success: true, data: {} } })),
      api.get('/api/notice').catch(() => ({ data: { success: true, data: '' } })),
    ]

    const [dashRes, statRes, noticeRes] = await Promise.all(requests)

    // 1. 公告
    if (noticeRes.data?.success) notice.value = noticeRes.data.data || ''

    // 2. 统计信息（来自 stat 接口）→ 总配额卡始终用 stat 值
    if (statRes?.data?.success) {
      const sd = statRes.data.data || {}
      logStat.quota = sd.quota || 0
      logStat.normal_quota = sd.normal_quota || 0
      logStat.subscription_quota = sd.subscription_quota || 0
      statData.total_quota = sd.quota || 0
    }

    // 3. 近 7 天使用明细（来自 dashboard 接口）
    if (dashRes.data?.success && Array.isArray(dashRes.data.data)) {
      buildCharts(dashRes.data.data)
    } else {
      // 兜底：dashboard 无数据时尝试拉取日志聚合（管理员拉全站，普通用户拉自己）
      await loadLogsFallback()
    }

    loadApiKey()
    loadSubscription()
    loadSelf()
  } catch (e) {
    console.warn('Dashboard load error', e)
  } finally {
    loading.value = false
  }
}

function buildCharts(logs) {
  // logs = LogStatistic[]，后端 GORM 序列化字段名为 PascalCase（Day/ModelName/...）
  // 兼容两种命名：PascalCase（实际返回）和小写（fallback）
  const list = Array.isArray(logs) ? logs : []
  const pick = (obj, lower, upper) => obj[lower] ?? obj[upper] ?? 0

  // 1) 聚合核心指标
  let totalTokens = 0
  let totalRequests = 0
  let totalQuota = 0
  const modelTotals = {}
  for (const item of list) {
    const day = item.Day ?? item.day
    if (!day) continue
    const requestCount = pick(item, 'request_count', 'RequestCount')
    const quota = pick(item, 'quota', 'Quota')
    const prompt = pick(item, 'prompt_tokens', 'PromptTokens')
    const completion = pick(item, 'completion_tokens', 'CompletionTokens')
    const modelName = item.ModelName ?? item.model_name ?? '其他'
    const tokens = prompt + completion

    totalRequests += requestCount
    totalQuota += quota
    totalTokens += tokens
    modelTotals[modelName] = (modelTotals[modelName] || 0) + tokens
  }
  statData.total_tokens = totalTokens
  statData.total_requests = totalRequests
  // 总配额以 stat 接口为准（loadDashboard 已赋值），仅当 stat 失败时用 dashboard 聚合兜底
  if (statData.total_quota === 0) {
    statData.total_quota = totalQuota
  }

  // 2) 生成近 7 天日期序列（按 day 分组，用本地日期避免时区偏移）
  const today = new Date()
  const last7Days = []
  const fmtLocal = (d) => {
    const y = d.getFullYear()
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const dd = String(d.getDate()).padStart(2, '0')
    return `${y}-${m}-${dd}`
  }
  for (let i = 6; i >= 0; i--) {
    const d = new Date(today)
    d.setDate(d.getDate() - i)
    last7Days.push(fmtLocal(d))
  }

  const dayMap = {}
  for (const item of list) {
    const day = item.Day ?? item.day
    if (!day || !last7Days.includes(day)) continue
    if (!dayMap[day]) {
      dayMap[day] = { requests: 0, quota: 0, tokens: 0 }
    }
    const requestCount = pick(item, 'request_count', 'RequestCount')
    const quota = pick(item, 'quota', 'Quota')
    const prompt = pick(item, 'prompt_tokens', 'PromptTokens')
    const completion = pick(item, 'completion_tokens', 'CompletionTokens')
    dayMap[day].requests += requestCount
    dayMap[day].quota += quota
    dayMap[day].tokens += prompt + completion
  }

  chartData.value.requests = last7Days.map((d) => ({ date: d, value: dayMap[d]?.requests || 0 }))
  chartData.value.quota = last7Days.map((d) => ({ date: d, value: dayMap[d]?.quota || 0 }))
  chartData.value.tokens = last7Days.map((d) => ({ date: d, value: dayMap[d]?.tokens || 0 }))

  // 3) 选出 Top 模型（按 token 总量）
  const topModels = Object.entries(modelTotals)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 6)
    .map((x) => x[0])
  barModels.value = topModels

  // 4) 填充明细表
  details.value = list
    .slice()
    .sort((a, b) => {
      const dayA = a.Day ?? a.day
      const dayB = b.Day ?? b.day
      const mA = a.ModelName ?? a.model_name
      const mB = b.ModelName ?? b.model_name
      if (dayA !== dayB) return (dayB || '').localeCompare(dayA || '')
      return (mB || '').localeCompare(mA || '')
    })
    .map((d) => {
      const modelName = d.ModelName ?? d.model_name ?? ''
      const provider = findProviderByName(modelName) || findProviderByName(modelName.split('-')[0])
      return {
        model_name: modelName,
        provider_slug: provider?.slug || '',
        day: d.Day ?? d.day,
        request_count: pick(d, 'request_count', 'RequestCount'),
        quota: pick(d, 'quota', 'Quota'),
        prompt_tokens: pick(d, 'prompt_tokens', 'PromptTokens'),
        completion_tokens: pick(d, 'completion_tokens', 'CompletionTokens'),
      }
    })
}

onMounted(async () => {
  if (!statusStore.loaded) await statusStore.fetchStatus()
  version.value = statusStore.status?.version || ''
  nextTick(() => loadDashboard())
})
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* ============ 顶部欢迎条 ============ */
.welcome-bar {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  padding: 4px 4px 8px;
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

/* ============ 通用 Panel ============ */
.panel {
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-2);
  border-radius: 8px;
  padding: 18px 20px;
  margin-bottom: 16px;
  transition: border-color 0.15s, box-shadow 0.15s;
}
.panel:last-child {
  margin-bottom: 0;
}
.panel.is-active {
  border-color: rgb(var(--primary-6));
  box-shadow: 0 0 0 1px rgb(var(--primary-6)) inset;
}
.panel.clickable {
  cursor: pointer;
  transition: border-color 0.15s, box-shadow 0.15s, transform 0.15s;
}
.panel.clickable:hover {
  border-color: rgb(var(--primary-5));
}
.panel.no-pad {
  padding: 18px 0 0;
}
.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14px;
}
.panel-head.pad-head {
  padding: 0 20px;
}
.panel-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-1);
  margin: 0;
}
.panel-extra {
  font-size: 12px;
  color: var(--color-text-3);
}
.panel-link {
  font-size: 12px;
  color: rgb(var(--primary-6));
  cursor: pointer;
  text-decoration: none;
}

/* ============ 核心指标 ============ */
.stat-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  padding: 18px 20px;
}
.stat-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding-right: 16px;
  border-right: 1px solid var(--color-fill-3);
}
.stat-item:last-child {
  border-right: none;
  padding-right: 0;
}
.stat-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.stat-label {
  font-size: 13px;
  color: var(--color-text-3);
}
.stat-icon {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: var(--color-text-1);
  line-height: 1.2;
  letter-spacing: -0.3px;
  word-break: break-all;
}
.stat-foot {
  font-size: 12px;
  color: var(--color-text-3);
}

/* ============ 用量进度 ============ */
.usage-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  padding: 18px 20px;
}
.usage-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.usage-row {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
}
.usage-label {
  font-size: 13px;
  color: var(--color-text-2);
}
.usage-pct {
  font-size: 18px;
  font-weight: 600;
  letter-spacing: -0.2px;
}
.usage-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--color-text-3);
}

/* ============ 趋势图 ============ */
.trend-cell {
  background: var(--color-fill-1);
  border-radius: 6px;
  padding: 12px 14px;
}
.trend-head {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 6px;
}
.trend-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}
.trend-label {
  font-size: 12px;
  color: var(--color-text-2);
}
.trend-total {
  margin-left: auto;
  font-size: 12px;
  color: var(--color-text-1);
  font-weight: 600;
  font-variant-numeric: tabular-nums;
}

/* ============ 图表空态 ============ */
.chart-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 320px;
  color: var(--color-text-4);
  font-size: 13px;
  background: var(--color-fill-1);
  border-radius: 6px;
}

/* ============ 表格 ============ */
.dash-table {
  background: transparent;
}
.dash-table :deep(.arco-table) {
  background: transparent;
}
.dash-table :deep(.arco-table-th) {
  background: var(--color-fill-1);
  font-size: 12px;
  color: var(--color-text-3);
  font-weight: 500;
}
.dash-table :deep(.arco-table-td) {
  font-size: 13px;
  color: var(--color-text-2);
}
.dash-table :deep(.arco-table-tr:hover .arco-table-td) {
  background: var(--color-fill-1);
}

.model-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}
.model-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--color-text-1);
}
.cell-num {
  font-variant-numeric: tabular-nums;
  font-weight: 500;
  color: var(--color-text-1);
}
.cell-mono {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 12px;
  color: var(--color-text-2);
}
.empty-state {
  padding: 40px 0;
  text-align: center;
  color: var(--color-text-4);
  font-size: 13px;
}

/* ============ 快捷操作 ============ */
.quick-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 8px;
}
.quick-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: var(--color-fill-1);
  border: 1px solid transparent;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s;
  text-align: left;
}
.quick-btn:hover {
  background: var(--color-fill-2);
  border-color: var(--color-border-2);
}
.quick-icon {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.quick-label {
  font-size: 13px;
  color: var(--color-text-1);
}

/* ============ 余额卡 ============ */
.balance-card {
  padding: 18px 20px;
}
.balance-head {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 14px;
}
.balance-icon {
  width: 34px;
  height: 34px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}
.balance-title {
  display: flex;
  flex-direction: column;
  gap: 1px;
}
.balance-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-1);
}
.balance-unit {
  font-size: 11px;
  color: var(--color-text-4);
}
.balance-value {
  margin-bottom: 16px;
}
.balance-num {
  font-size: 34px;
  font-weight: 700;
  color: var(--color-text-1);
  letter-spacing: -0.5px;
  line-height: 1.2;
  font-variant-numeric: tabular-nums;
}
.balance-actions {
  display: flex;
  gap: 8px;
  padding-top: 14px;
  border-top: 1px solid var(--color-fill-3);
}
.balance-btn {
  flex: 1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 7px 0;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  color: rgb(var(--primary-6));
  background: rgba(22, 93, 255, 0.06);
  border: 1px solid rgba(22, 93, 255, 0.12);
  cursor: pointer;
  text-decoration: none;
  transition: background 0.15s, border-color 0.15s;
}
.balance-btn:hover {
  background: rgba(22, 93, 255, 0.12);
  border-color: rgba(22, 93, 255, 0.22);
}
.balance-btn:active {
  transform: scale(0.98);
}

/* ============ API Key ============ */
.apikey-box {
  display: flex;
  align-items: center;
  gap: 4px;
  background: var(--color-fill-1);
  border-radius: 6px;
  padding: 6px 6px 6px 12px;
}
.apikey-masked {
  flex: 1;
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 12px;
  color: var(--color-text-1);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.apikey-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--color-text-3);
  margin-top: 8px;
}
.meta-name {
  color: var(--color-text-2);
}
.apikey-empty {
  padding: 8px 0 4px;
  text-align: center;
}
.empty-title {
  font-size: 13px;
  color: var(--color-text-1);
  font-weight: 500;
  margin: 0 0 4px;
}
.empty-desc {
  font-size: 12px;
  color: var(--color-text-3);
  margin: 0 0 12px;
}

/* ============ 公告 ============ */
.notice-body {
  font-size: 13px;
  color: var(--color-text-2);
  line-height: 1.7;
}
.notice-list {
  display: flex;
  flex-direction: column;
}
.notice-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid var(--color-fill-3);
  text-decoration: none;
  transition: color 0.15s;
}
.notice-item:last-child {
  border-bottom: none;
}
.notice-item:hover .notice-title {
  color: rgb(var(--primary-6));
}
.notice-title {
  font-size: 13px;
  color: var(--color-text-2);
}
.notice-time {
  font-size: 11px;
  color: var(--color-text-4);
  font-variant-numeric: tabular-nums;
  flex-shrink: 0;
  margin-left: 12px;
}

/* ============ 更新日志 ============ */
.changelog {
  display: flex;
  flex-direction: column;
}
.cl-item {
  padding: 10px 0;
  border-bottom: 1px solid var(--color-fill-3);
}
.cl-item:last-child {
  border-bottom: none;
}
.cl-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 4px;
}
.cl-version {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-1);
  font-family: 'SF Mono', Menlo, Consolas, monospace;
}
.cl-date {
  font-size: 11px;
  color: var(--color-text-4);
  font-variant-numeric: tabular-nums;
}
.cl-desc {
  font-size: 12px;
  color: var(--color-text-3);
  line-height: 1.5;
  margin: 0;
}

/* ============ 资源 ============ */
.contact-list {
  display: flex;
  flex-direction: column;
}
.contact-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid var(--color-fill-3);
  text-decoration: none;
  transition: color 0.15s;
}
.contact-item:last-child {
  border-bottom: none;
}
.contact-label {
  font-size: 13px;
  color: var(--color-text-2);
}
.contact-value {
  font-size: 12px;
  color: var(--color-text-3);
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-left: 12px;
}
.contact-item:hover .contact-value {
  color: rgb(var(--primary-6));
}
.contact-value.is-active {
  color: rgb(var(--primary-6));
}
.contact-value.is-active {
  font-weight: 500;
}
.contact-arrow {
  color: rgb(var(--primary-6));
  opacity: 0.6;
  flex-shrink: 0;
  transition: opacity 0.15s;
}
.contact-item:hover .contact-arrow {
  opacity: 1;
}

/* ============ 响应式 ============ */
@media (max-width: 1200px) {
  .stat-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 20px 16px;
  }
  .stat-item {
    border-right: none;
    padding-right: 0;
  }
  .stat-item:nth-child(odd) {
    padding-right: 16px;
    border-right: 1px solid var(--color-fill-3);
  }
}
@media (max-width: 768px) {
  .usage-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  .welcome-bar {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>
