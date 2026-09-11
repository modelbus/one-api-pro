<!--
  AdminDashboard 管理员运营仪表盘
  Admin operations dashboard: KPIs + 7-day trends + leaderboard.

  版本: v0.0.16
  日期: 2026-09-11
  作者: opencode
-->
<template>
  <div class="admin-dashboard">
    <!-- 顶栏：时间筛选 + 刷新 -->
    <!-- Top bar: range filter + refresh -->
    <div class="dashboard-toolbar">
      <a-radio-group
        v-model="range"
        type="button"
        @change="loadOverview"
      >
        <a-radio value="today">{{ t('admin.rangeToday') }}</a-radio>
        <a-radio value="7d">{{ t('admin.range7d') }}</a-radio>
        <a-radio value="30d">{{ t('admin.range30d') }}</a-radio>
        <a-radio value="all">{{ t('admin.rangeAll') }}</a-radio>
      </a-radio-group>
      <div class="toolbar-right">
        <span class="last-refresh">{{ t('admin.lastRefresh') }}: {{ lastRefreshLabel }}</span>
        <a-button size="small" @click="loadAll">
          <template #icon><icon-refresh /></template>
          {{ t('admin.refresh') }}
        </a-button>
      </div>
    </div>

    <a-spin :loading="overviewLoading" style="width: 100%">
      <!-- ① 用户 -->
      <a-card class="panel" :title="t('admin.sectionUsers')">
        <a-grid :cols="7" :col-gap="12" :row-gap="12">
          <a-grid-item v-for="item in userCards" :key="item.label">
            <div class="kpi-cell">
              <div class="kpi-label">{{ item.label }}</div>
              <div class="kpi-value" :style="{ color: item.color }">{{ item.value }}</div>
              <div v-if="item.foot" class="kpi-foot">{{ item.foot }}</div>
            </div>
          </a-grid-item>
        </a-grid>
      </a-card>

      <!-- ② 资源 / 用量 -->
      <a-card class="panel" :title="t('admin.sectionResources')">
        <a-grid :cols="5" :col-gap="12" :row-gap="12">
          <a-grid-item v-for="item in resourceCards" :key="item.label">
            <div
              class="kpi-cell"
              :class="{ 'kpi-cell-clickable': item.to }"
              @click="item.to && $router.push(item.to)"
            >
              <div class="kpi-label">{{ item.label }}</div>
              <div class="kpi-value">{{ item.value }}</div>
              <div v-if="item.foot" class="kpi-foot">{{ item.foot }}</div>
            </div>
          </a-grid-item>
        </a-grid>
      </a-card>

      <!-- ③ Token / 请求消耗 -->
      <a-card class="panel" :title="t('admin.sectionQuota')">
        <a-grid :cols="4" :col-gap="12" :row-gap="12">
          <a-grid-item v-for="item in quotaCards" :key="item.label">
            <div class="kpi-cell">
              <div class="kpi-label">{{ item.label }}</div>
              <div class="kpi-value">{{ item.value }}</div>
              <div v-if="item.foot" class="kpi-foot">{{ item.foot }}</div>
            </div>
          </a-grid-item>
        </a-grid>

        <a-grid :cols="2" :col-gap="16" :row-gap="16" class="chart-row">
          <a-grid-item>
            <div class="chart-title">{{ t('admin.chartRequests') }}</div>
            <v-chart :option="requestTrendOption" :style="{ height: '240px' }" autoresize />
          </a-grid-item>
          <a-grid-item>
            <div class="chart-title">{{ t('admin.chartQuota') }}</div>
            <v-chart :option="quotaTrendOption" :style="{ height: '240px' }" autoresize />
          </a-grid-item>
        </a-grid>
      </a-card>

      <!-- ④ 收入 -->
      <a-card class="panel" :title="t('admin.sectionRevenue')">
        <a-grid :cols="4" :col-gap="12" :row-gap="12">
          <a-grid-item v-for="item in revenueCards" :key="item.label">
            <div class="kpi-cell">
              <div class="kpi-label">{{ item.label }}</div>
              <div class="kpi-value">¥ {{ item.value }}</div>
              <div v-if="item.foot" class="kpi-foot">{{ item.foot }}</div>
            </div>
          </a-grid-item>
        </a-grid>
      </a-card>
    </a-spin>

    <!-- ⑤ 活跃用户 TOP 20 -->
    <a-card class="panel" :title="t('admin.sectionLeaderboard')">
      <template #extra>
        <a-radio-group v-model="topRange" type="button" @change="loadTopUsers">
          <a-radio value="today">{{ t('admin.rangeToday') }}</a-radio>
          <a-radio value="7d">{{ t('admin.range7d') }}</a-radio>
          <a-radio value="30d">{{ t('admin.range30d') }}</a-radio>
        </a-radio-group>
      </template>
      <a-table
        :columns="topColumns"
        :data="topUsers"
        :pagination="false"
        :loading="topLoading"
        row-key="id"
        stripe
      >
        <template #request_count="{ record }">
          {{ record.request_count.toLocaleString() }}
        </template>
        <template #quota="{ record }">
          {{ formatQuota(record.quota) }}
        </template>
        <template #balance="{ record }">
          {{ formatQuota(record.balance) }}
        </template>
        <template #current_plan_name="{ record }">
          <a-tag v-if="record.current_plan_name" color="arcoblue">{{ record.current_plan_name }}</a-tag>
          <span v-else class="muted">-</span>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { IconRefresh } from '@arco-design/web-vue/es/icon'
import VChart from 'vue-echarts'
import 'echarts'
import adminApi from '@/api/admin'

const { t } = useI18n()
const router = useRouter()

const range = ref('7d')
const topRange = ref('7d')
const overviewLoading = ref(false)
const topLoading = ref(false)
const overview = ref(null)
const topUsers = ref([])
const lastRefreshAt = ref(0)

// ---------- 顶部刷新时间 ----------
const lastRefreshLabel = computed(() => {
  if (!lastRefreshAt.value) return '-'
  const d = new Date(lastRefreshAt.value * 1000)
  const pad = (n) => String(n).padStart(2, '0')
  return `${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
})

// ---------- 配额/数值单位格式化 ----------
// formatQuota: 按 1e3/1e6/1e9 自动换算到 K/M/B，2 位小数。
// formatQuota: auto-suffix K/M/B for big numbers with 2 decimals.
function formatQuota(n) {
  const num = Number(n) || 0
  if (num < 10000) return String(num)
  if (num < 1e6) return (num / 1e3).toFixed(2) + 'K'
  if (num < 1e9) return (num / 1e6).toFixed(2) + 'M'
  return (num / 1e9).toFixed(2) + 'B'
}

function fmtMoney(n) {
  const num = Number(n) || 0
  return num.toFixed(2)
}

// ---------- KPI 卡数据 ----------
const userCards = computed(() => {
  const u = overview.value?.users
  if (!u) return []
  return [
    { label: t('admin.userTotal'), value: u.total.toLocaleString(), color: '#165dff' },
    { label: t('admin.userNewToday'), value: u.new_today.toLocaleString(), color: '#00b42a' },
    { label: t('admin.userNew7d'), value: u.new_7d.toLocaleString(), color: '#00b42a' },
    { label: t('admin.userNew30d'), value: u.new_30d.toLocaleString(), color: '#00b42a' },
    { label: t('admin.userActive7d'), value: u.active_7d.toLocaleString(), color: '#722ed1' },
    { label: t('admin.userDisabled'), value: u.disabled.toLocaleString(), color: '#86909c' },
    { label: t('admin.userDeleted'), value: u.deleted.toLocaleString(), color: '#86909c' },
  ]
})

const resourceCards = computed(() => {
  const ov = overview.value
  if (!ov) return []
  return [
    {
      label: t('admin.tokenTotal'),
      value: ov.tokens.total.toLocaleString(),
      foot: `${t('admin.enabled')}: ${ov.tokens.enabled.toLocaleString()}`,
    },
    {
      label: t('admin.channelTotal'),
      value: ov.channels.total.toLocaleString(),
      foot: `${t('admin.enabled')}: ${ov.channels.enabled.toLocaleString()}`,
    },
    {
      label: t('admin.planTotal'),
      value: ov.plans.total.toLocaleString(),
      foot: `${t('admin.enabled')}: ${ov.plans.enabled.toLocaleString()}`,
      to: '/setting/plan',
    },
    {
      label: t('admin.redemptionTotal'),
      value: ov.redemptions.total.toLocaleString(),
      foot: `${t('admin.used')}: ${ov.redemptions.used.toLocaleString()} / ${t('admin.unused')}: ${ov.redemptions.unused.toLocaleString()}`,
      to: '/redemption',
    },
    {
      label: t('admin.subscriptionTotal'),
      value: ov.subscriptions.total.toLocaleString(),
      foot: `${t('admin.active')}: ${ov.subscriptions.active.toLocaleString()} / ${t('admin.expired')}: ${ov.subscriptions.expired.toLocaleString()}`,
    },
  ]
})

const quotaCards = computed(() => {
  const q = overview.value?.quota
  if (!q) return []
  return [
    { label: t('admin.quotaToday'), value: formatQuota(q.today) },
    { label: t('admin.quota7d'), value: formatQuota(q.week), foot: t('admin.quotaUnit') },
    { label: t('admin.quota30d'), value: formatQuota(q.month) },
    { label: t('admin.quotaTotal'), value: formatQuota(q.total) },
  ]
})

const revenueCards = computed(() => {
  const r = overview.value?.revenue
  if (!r) return []
  return [
    { label: t('admin.revenueTotal'), value: fmtMoney(r.total) },
    { label: t('admin.revenueTopup'), value: fmtMoney(r.topup) },
    { label: t('admin.revenueSubscription'), value: fmtMoney(r.subscription) },
    { label: t('admin.revenueRefund'), value: fmtMoney(r.refund) },
  ]
})

// ---------- 趋势图 ----------
function buildLineOption(field, color) {
  const points = overview.value?.trends?.[field] || []
  return {
    grid: { left: 40, right: 16, top: 16, bottom: 28 },
    tooltip: { trigger: 'axis', confine: true, backgroundColor: '#1d2129', borderWidth: 0, textStyle: { color: '#fff', fontSize: 12 } },
    xAxis: {
      type: 'category', boundaryGap: false,
      data: points.map((d) => (d.day || '').slice(5)),
      axisLabel: { fontSize: 11, color: '#86909c' },
      axisLine: { show: false },
      axisTick: { show: false },
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        fontSize: 11,
        color: '#86909c',
        formatter: (v) => (v >= 1000 ? `${(v / 1000).toFixed(0)}k` : v),
      },
      splitLine: { lineStyle: { color: '#f2f3f5', type: 'dashed' } },
      axisLine: { show: false },
      axisTick: { show: false },
    },
    series: [{
      type: 'line', smooth: true, symbol: 'circle', symbolSize: 5,
      data: points.map((d) => (field === 'quota' ? d.quota : d.count)),
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

const requestTrendOption = computed(() => buildLineOption('requests', '#165dff'))
const quotaTrendOption = computed(() => buildLineOption('quota', '#722ed1'))

// ---------- 排行榜 ----------
const topColumns = computed(() => [
  { title: '#', dataIndex: 'rank', width: 56 },
  { title: t('admin.colUsername'), dataIndex: 'username', width: 160 },
  { title: t('admin.colEmail'), dataIndex: 'email', ellipsis: true },
  { title: t('admin.colRequestCount'), slotName: 'request_count', width: 140, align: 'right' },
  { title: t('admin.colQuota'), slotName: 'quota', width: 140, align: 'right' },
  { title: t('admin.colBalance'), slotName: 'balance', width: 140, align: 'right' },
  { title: t('admin.colPlan'), slotName: 'current_plan_name', width: 120, align: 'center' },
])

// ---------- 数据加载 ----------
async function loadOverview() {
  overviewLoading.value = true
  try {
    const { data } = await adminApi.overview(range.value)
    if (data?.success) {
      overview.value = data.data
      lastRefreshAt.value = Math.floor(Date.now() / 1000)
    } else {
      Message.error(data?.message || t('admin.loadFailed'))
    }
  } catch (e) {
    Message.error(t('admin.loadFailed') + ': ' + (e?.message || ''))
  } finally {
    overviewLoading.value = false
  }
}

async function loadTopUsers() {
  topLoading.value = true
  try {
    const { data } = await adminApi.topUsers(topRange.value, 20)
    if (data?.success) {
      const rows = data.data?.items || []
      topUsers.value = rows.map((r, idx) => ({ ...r, rank: idx + 1 }))
    } else {
      Message.error(data?.message || t('admin.loadFailed'))
    }
  } catch (e) {
    Message.error(t('admin.loadFailed') + ': ' + (e?.message || ''))
  } finally {
    topLoading.value = false
  }
}

function loadAll() {
  loadOverview()
  loadTopUsers()
}

onMounted(loadAll)
</script>

<style scoped>
.admin-dashboard {
  padding: 16px 0;
}

.dashboard-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.last-refresh {
  color: var(--color-text-3);
  font-size: 12px;
}

.panel {
  margin-bottom: 16px;
}

.kpi-cell {
  border: 1px solid var(--color-border-2);
  border-radius: 6px;
  padding: 12px;
  background: var(--color-fill-1);
  transition: border-color 0.2s;
}

.kpi-cell-clickable {
  cursor: pointer;
}

.kpi-cell-clickable:hover {
  border-color: rgb(var(--arcoblue-6));
}

.kpi-label {
  font-size: 12px;
  color: var(--color-text-3);
  margin-bottom: 6px;
}

.kpi-value {
  font-size: 22px;
  font-weight: 600;
  color: var(--color-text-1);
  line-height: 1.2;
}

.kpi-foot {
  margin-top: 4px;
  font-size: 11px;
  color: var(--color-text-3);
}

.chart-row {
  margin-top: 16px;
}

.chart-title {
  font-size: 13px;
  color: var(--color-text-2);
  margin-bottom: 8px;
}

.muted {
  color: var(--color-text-4);
}
</style>
