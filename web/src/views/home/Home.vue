<template>
  <div class="dashboard-page">
    <!-- Welcome Section -->
    <div class="welcome-section">
      <div class="welcome-left">
        <div class="welcome-avatar">{{ firstChar }}</div>
        <div class="welcome-text">
          <h2>{{ greeting }}，{{ authStore.user?.username || '用户' }}</h2>
          <p>{{ dateStr }} · 欢迎回到控制台</p>
        </div>
      </div>
      <div class="welcome-right">
        <a-space>
          <a-tag color="arcoblue" size="medium">角色：{{ roleLabel }}</a-tag>
          <a-tag v-if="version" color="green" size="medium">v{{ version }}</a-tag>
        </a-space>
      </div>
    </div>

    <!-- Stat Cards -->
    <a-row :gutter="[16, 16]" class="stat-grid">
      <a-col :span="authStore.isAdmin ? 6 : 8">
        <div class="stat-card">
          <div class="stat-icon" style="background:linear-gradient(135deg,#165dff,#4080ff)">
            <icon-code-square :size="22" />
          </div>
          <div class="stat-body">
            <div class="stat-label">总 Token</div>
            <div class="stat-value">{{ fmtNum(statData.total_tokens) }}</div>
          </div>
        </div>
      </a-col>
      <a-col :span="authStore.isAdmin ? 6 : 8">
        <div class="stat-card">
          <div class="stat-icon" style="background:linear-gradient(135deg,#0fc6c2,#36cfc9)">
            <icon-send :size="22" />
          </div>
          <div class="stat-body">
            <div class="stat-label">总请求</div>
            <div class="stat-value">{{ fmtNum(statData.total_requests) }}</div>
          </div>
        </div>
      </a-col>
      <a-col :span="authStore.isAdmin ? 6 : 8">
        <div class="stat-card">
          <div class="stat-icon" style="background:linear-gradient(135deg,#f7ba1e,#ffcf3f)">
            <icon-archive :size="22" />
          </div>
          <div class="stat-body">
            <div class="stat-label">总配额</div>
            <div class="stat-value">{{ fmtQuota(statData.total_quota) }}</div>
          </div>
        </div>
      </a-col>
      <template v-if="authStore.isAdmin">
        <a-col :span="6">
          <div class="stat-card">
            <div class="stat-icon" style="background:linear-gradient(135deg,#00b42a,#36cf6a)">
              <icon-user-group :size="22" />
            </div>
            <div class="stat-body">
              <div class="stat-label">用户数</div>
              <div class="stat-value">{{ fmtNum(logStat.total_users || 0) }}</div>
            </div>
          </div>
        </a-col>
        <a-col :span="6">
          <div class="stat-card">
            <div class="stat-icon" style="background:linear-gradient(135deg,#722ed1,#9254de)">
              <icon-apps :size="22" />
            </div>
            <div class="stat-body">
              <div class="stat-label">活跃渠道</div>
              <div class="stat-value">{{ fmtNum(logStat.active_channels || 0) }}</div>
            </div>
          </div>
        </a-col>
        <a-col :span="6">
          <div class="stat-card">
            <div class="stat-icon" style="background:linear-gradient(135deg,#ff5722,#ff7a45)">
              <icon-storage :size="22" />
            </div>
            <div class="stat-body">
              <div class="stat-label">今日配额</div>
              <div class="stat-value">{{ fmtQuota(logStat.quota || 0) }}</div>
            </div>
          </div>
        </a-col>
      </template>
    </a-row>

    <!-- Charts Row -->
    <a-row :gutter="16" class="chart-row">
      <a-col :span="14">
        <a-card :bordered="false" class="chart-card">
          <template #title><span class="card-hd">📈 Token 用量趋势</span></template>
          <v-chart :option="lineOption" :style="{ height: '340px' }" />
        </a-card>
      </a-col>
      <a-col :span="10">
        <a-card :bordered="false" class="chart-card">
          <template #title><span class="card-hd">📊 用量分布</span></template>
          <v-chart :option="pieOption" :style="{ height: '340px' }" autoresize />
        </a-card>
      </a-col>
    </a-row>

    <!-- Bottom Row: Log Stats + Quick Actions -->
    <a-row :gutter="16" class="bottom-row">
      <a-col :span="14">
        <a-card :bordered="false" class="chart-card">
          <template #title><span class="card-hd">📋 统计明细</span></template>
          <a-row :gutter="16">
            <a-col :span="8">
              <a-statistic title="Prompt Tokens" :value="fmtNum(logStat.prompt_tokens || 0)" :value-style="{fontSize:'24px',fontWeight:600}" />
            </a-col>
            <a-col :span="8">
              <a-statistic title="Completion Tokens" :value="fmtNum(logStat.completion_tokens || 0)" :value-style="{fontSize:'24px',fontWeight:600}" />
            </a-col>
            <a-col :span="8">
              <a-statistic title="普通配额" :value="fmtQuota(logStat.normal_quota || 0)" :value-style="{fontSize:'24px',fontWeight:600}" />
            </a-col>
          </a-row>
          <a-divider :margin="16" />
          <a-row :gutter="16">
            <a-col :span="8">
              <a-statistic title="订阅配额" :value="fmtQuota(logStat.subscription_quota || 0)" :value-style="{fontSize:'24px',fontWeight:600}" />
            </a-col>
            <a-col :span="8">
              <a-statistic title="Token" :value="fmtNum(logStat.token || 0)" :value-style="{fontSize:'24px',fontWeight:600}" />
            </a-col>
            <a-col :span="8">
              <a-statistic title="请求数" :value="fmtNum(logStat.request_count || 0)" :value-style="{fontSize:'24px',fontWeight:600}" />
            </a-col>
          </a-row>
        </a-card>
      </a-col>
      <a-col :span="10">
        <a-card :bordered="false" class="chart-card">
          <template #title><span class="card-hd">⚡ 快捷操作</span></template>
          <a-space direction="vertical" fill>
            <a-button type="outline" long @click="$router.push('/token')">
              <template #icon><icon-code /></template> 管理令牌
            </a-button>
            <a-button type="outline" long @click="$router.push('/topup')">
              <template #icon><icon-archive /></template> 充值额度
            </a-button>
            <a-button type="outline" long @click="$router.push('/log')">
              <template #icon><icon-file /></template> 查看日志
            </a-button>
            <a-button v-if="authStore.isAdmin" type="outline" long @click="$router.push('/channel')">
              <template #icon><icon-apps /></template> 管理渠道
            </a-button>
            <a-button v-if="authStore.isAdmin" type="outline" long @click="$router.push('/user')">
              <template #icon><icon-user-group /></template> 管理用户
            </a-button>
          </a-space>
        </a-card>
      </a-col>
    </a-row>

    <!-- Notice -->
    <a-card v-if="notice" :bordered="false" class="notice-card">
      <template #title><span class="card-hd">📢 系统公告</span></template>
      <div v-html="notice"></div>
    </a-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useStatusStore } from '@/stores/status'
import { IconUserGroup, IconCodeSquare, IconSend, IconFire, IconBarChart, IconBulb, IconCode, IconArchive, IconFile, IconApps, IconStorage } from '@arco-design/web-vue/es/icon'
import VChart from 'vue-echarts'
import 'echarts'
import api from '@/api'

const authStore = useAuthStore()
const statusStore = useStatusStore()

const loading = ref(true)
const notice = ref('')
const version = ref(statusStore.status?.version || '')

const statData = ref({ total_tokens: 0, total_requests: 0, total_quota: 0 })
const dailyUsage = ref([])
const logStat = ref({ quota: 0, token: 0, normal_quota: 0, subscription_quota: 0, prompt_tokens: 0, completion_tokens: 0, total_users: 0, active_channels: 0, request_count: 0 })

const firstChar = computed(() => (authStore.user?.username || 'U')[0].toUpperCase())
const roleLabel = computed(() => authStore.isRoot ? '超级管理员' : authStore.isAdmin ? '管理员' : '普通用户')

const dateStr = computed(() => {
  const d = new Date()
  return `${d.getFullYear()}/${d.getMonth()+1}/${d.getDate()} ${['日','一','二','三','四','五','六'][d.getDay()]}`
})

const greeting = computed(() => {
  const h = new Date().getHours()
  if (h < 6) return '夜深了'
  if (h < 9) return '早上好'
  if (h < 12) return '上午好'
  if (h < 14) return '中午好'
  if (h < 18) return '下午好'
  return '晚上好'
})

function fmtNum(n) { return (n || 0).toLocaleString() }
function fmtQuota(n) { return (n || 0).toLocaleString(undefined, { minimumFractionDigits: 1, maximumFractionDigits: 1 }) }

const lineOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  grid: { left: 40, right: 20, top: 20, bottom: 30 },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: dailyUsage.value.length ? dailyUsage.value.map(d => d.date || d.day || '') : mockDays,
    axisLabel: { fontSize: 11, color: '#86909c' },
  },
  yAxis: {
    type: 'value',
    name: 'Tokens',
    axisLabel: { formatter: v => v > 1000 ? (v/1000).toFixed(0)+'k' : v, fontSize: 11 },
    splitLine: { lineStyle: { color: '#f2f3f5' } },
  },
  series: [{
    name: 'Token 用量',
    type: 'line',
    smooth: true,
    symbol: 'none',
    data: dailyUsage.value.length ? dailyUsage.value.map(d => d.tokens ?? d.token_count ?? 0) : mockData,
    areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{offset:0,color:'rgba(22,93,255,0.25)'},{offset:1,color:'rgba(22,93,255,0.02)'}] } },
    lineStyle: { color: '#165DFF', width: 2 },
    itemStyle: { color: '#165DFF' },
  }],
}))

const pieOption = computed(() => ({
  tooltip: { trigger: 'item' },
  legend: { bottom: 0, textStyle: { fontSize: 12 } },
  series: [{
    type: 'pie',
    radius: ['55%', '78%'],
    center: ['50%', '45%'],
    avoidLabelOverlap: false,
    itemStyle: { borderRadius: 4, borderColor: '#fff', borderWidth: 2 },
    label: { show: false },
    data: dailyUsage.value.length ? buildPieData() : [
      { value: 45, name: 'Chat', itemStyle: { color: '#165DFF' } },
      { value: 25, name: 'Completion', itemStyle: { color: '#0FC6C2' } },
      { value: 18, name: 'Embedding', itemStyle: { color: '#F7BA1E' } },
      { value: 12, name: 'Image', itemStyle: { color: '#722ED1' } },
    ],
  }],
}))

function buildPieData() {
  const usage = dailyUsage.value
  const types = {}
  usage.forEach(d => { const t = d.type || d.model || 'Other'; types[t] = (types[t]||0)+(d.tokens||d.token_count||0) })
  return Object.entries(types).map(([name,value]) => ({ name, value }))
}

const mockDays = ['07/19','07/20','07/21','07/22','07/23','07/24','07/25']
const mockData = [820, 932, 901, 934, 1290, 1330, 1520]

onMounted(async () => {
  await statusStore.fetchStatus()
  version.value = statusStore.status?.version || ''
  await fetchData()
})

async function fetchData() {
  loading.value = true
  try {
    const requests = [
      api.get('/api/user/dashboard').catch(() => ({ data: { success: true, data: {} } })),
      api.get('/api/notice').catch(() => ({ data: { success: true, data: '' } })),
    ]
    if (authStore.isAdmin) {
      requests.push(api.get('/api/log/stat').catch(() => ({ data: { success: true, data: {} } })))
    } else {
      requests.push(api.get('/api/log/self/stat').catch(() => ({ data: { success: true, data: {} } })))
    }

    const [dashboardRes, noticeRes, statRes] = await Promise.all(requests)
    if (dashboardRes.data?.success) {
      const d = dashboardRes.data.data || {}
      statData.value.total_tokens = d.total_tokens ?? d.token_count ?? 0
      statData.value.total_requests = d.total_requests ?? d.request_count ?? 0
      statData.value.total_quota = d.total_quota ?? d.quota ?? 0
      dailyUsage.value = d.daily_usage ?? d.daily_data ?? []
    }
    if (noticeRes.data?.success) notice.value = noticeRes.data.data || ''
    if (statRes?.data?.success) Object.assign(logStat.value, statRes.data.data || {})
  } catch (e) { /* ignore */ } finally { loading.value = false }
}
</script>

<style scoped>
.dashboard-page { max-width: 1280px; margin: 0 auto; }

.welcome-section {
  display: flex; align-items: center; justify-content: space-between;
  padding: 20px 24px; background: var(--color-bg-2); border-radius: 8px; margin-bottom: 20px;
}
.welcome-left { display: flex; align-items: center; gap: 16px; }
.welcome-avatar {
  width: 48px; height: 48px; border-radius: 12px;
  background: linear-gradient(135deg, #165dff, #722ed1);
  color: #fff; font-size: 20px; font-weight: 700;
  display: flex; align-items: center; justify-content: center;
}
.welcome-text h2 { font-size: 18px; font-weight: 600; margin: 0 0 4px; }
.welcome-text p { font-size: 13px; color: var(--color-text-3); margin: 0; }

.stat-grid { margin-bottom: 20px; }
.stat-card {
  display: flex; align-items: center; gap: 16px;
  padding: 20px; background: var(--color-bg-2); border-radius: 8px;
  transition: box-shadow .2s;
}
.stat-card:hover { box-shadow: 0 2px 12px rgba(0,0,0,0.06); }
.stat-icon {
  width: 46px; height: 46px; border-radius: 12px;
  display: flex; align-items: center; justify-content: center; color: #fff; flex-shrink: 0;
}
.stat-body { flex: 1; min-width: 0; }
.stat-label { font-size: 13px; color: var(--color-text-3); margin-bottom: 4px; }
.stat-value { font-size: 24px; font-weight: 700; color: var(--color-text-1); }

.chart-row { margin-bottom: 20px; }
.bottom-row { margin-bottom: 20px; }
.chart-card { border-radius: 8px; height: 100%; }
.card-hd { font-size: 15px; font-weight: 600; }
.notice-card { border-radius: 8px; margin-bottom: 20px; }
</style>
