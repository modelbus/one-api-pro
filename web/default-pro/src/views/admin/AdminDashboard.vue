<!--
  AdminDashboard 管理员运营仪表盘
  Admin operations dashboard: 全站 KPI · 7 日趋势 · 活跃用户排行

  排版参照 web/default-pro/src/views/dashboard/Dashboard.vue：
  顶部 welcome-bar + 左 16 / 右 8 分栏 + 多组 stat-grid + trend-cell + dash-table。
  Layout references Dashboard.vue: welcome-bar header + left 16 / right 8 split
  + multi stat-grid + trend-cell + dash-table.

  版本: v0.0.17
  日期: 2026-09-11
  作者: opencode
-->
<template>
  <div class="dashboard admin-dashboard">
    <!-- 顶部欢迎条 -->
    <!-- Top welcome bar -->
    <div class="welcome-bar">
      <div class="welcome-text">
        <h1 class="welcome-title">{{ t('admin.pageTitle') }}</h1>
        <p class="welcome-desc">{{ t('admin.pageSubtitle') }} — {{ todayStr }}</p>
      </div>
      <div class="welcome-meta">
        <a-radio-group
          v-model="range"
          type="button"
          size="small"
          @change="loadAll"
        >
          <a-radio value="today">{{ t('admin.rangeToday') }}</a-radio>
          <a-radio value="7d">{{ t('admin.range7d') }}</a-radio>
          <a-radio value="30d">{{ t('admin.range30d') }}</a-radio>
          <a-radio value="all">{{ t('admin.rangeAll') }}</a-radio>
        </a-radio-group>
        <span class="meta-chip meta-chip-mono">{{ t('admin.lastRefresh') }} · {{ lastRefreshLabel }}</span>
        <a-button size="mini" @click="loadAll">
          <template #icon><icon-refresh :size="14" /></template>
          {{ t('admin.refresh') }}
        </a-button>
      </div>
    </div>

    <a-row :gutter="16">
      <!-- 左侧主内容 -->
      <!-- Left column: 16/24 -->
      <a-col :xs="24" :lg="16">
        <a-spin :loading="overviewLoading" style="width: 100%">
          <!-- ① 用户层 KPI -->
          <div class="panel stat-grid">
            <div
              v-for="item in userStatItems"
              :key="item.label"
              class="stat-item"
            >
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

          <!-- ② 资源 / 用量 -->
          <div class="panel">
            <div class="panel-head">
              <h2 class="panel-title">{{ t('admin.sectionResources') }}</h2>
              <span class="panel-extra">{{ t('admin.clickToManage') }}</span>
            </div>
            <div class="stat-grid">
              <div
                v-for="item in resourceStatItems"
                :key="item.label"
                class="stat-item stat-item-clickable"
                @click="item.to && $router.push(item.to)"
              >
                <div class="stat-head">
                  <span class="stat-label">{{ item.label }}</span>
                  <span class="stat-icon" :style="{ background: item.bg, color: item.color }">
                    <component :is="item.icon" :size="14" />
                  </span>
                </div>
                <div class="stat-value">{{ item.value }}</div>
                <div class="stat-foot">{{ item.foot }}</div>
              </div>
            </div>
          </div>

          <!-- ③ Token / 请求消耗 + 趋势 -->
          <div class="panel">
            <div class="panel-head">
              <h2 class="panel-title">{{ t('admin.sectionQuota') }}</h2>
              <span class="panel-extra">{{ t('admin.quotaUnit') }} · {{ rangeLabel }}</span>
            </div>
            <div class="stat-grid">
              <div
                v-for="item in quotaStatItems"
                :key="item.label"
                class="stat-item"
              >
                <div class="stat-head">
                  <span class="stat-label">{{ item.label }}</span>
                  <span class="stat-icon" :style="{ background: item.bg, color: item.color }">
                    <component :is="item.icon" :size="14" />
                  </span>
                </div>
                <div class="stat-value">{{ item.value }}</div>
                <div class="stat-foot">{{ item.foot }}</div>
              </div>
            </div>

            <div class="trend-divider" />

            <a-row :gutter="12">
              <a-col :xs="24" :sm="12" v-for="t in trendItems" :key="t.field">
                <div class="trend-cell">
                  <div class="trend-head">
                    <span class="trend-dot" :style="{ background: t.color }"></span>
                    <span class="trend-label">{{ t.label }}</span>
                    <span class="trend-total">{{ t.total }}</span>
                  </div>
                  <v-chart
                    :option="lineOption(t.field, t.color)"
                    :style="{ height: '180px' }"
                    autoresize
                  />
                </div>
              </a-col>
            </a-row>
          </div>

          <!-- ④ 收入 -->
          <div class="panel">
            <div class="panel-head">
              <h2 class="panel-title">{{ t('admin.sectionRevenue') }}</h2>
              <span class="panel-extra">¥ · {{ rangeLabel }}</span>
            </div>
            <div class="stat-grid">
              <div
                v-for="item in revenueStatItems"
                :key="item.label"
                class="stat-item"
              >
                <div class="stat-head">
                  <span class="stat-label">{{ item.label }}</span>
                  <span class="stat-icon" :style="{ background: item.bg, color: item.color }">
                    <component :is="item.icon" :size="14" />
                  </span>
                </div>
                <div class="stat-value">¥ {{ item.value }}</div>
                <div class="stat-foot">{{ item.foot }}</div>
              </div>
            </div>
          </div>
        </a-spin>
      </a-col>

      <!-- 右侧栏 -->
      <!-- Right column: 8/24 -->
      <a-col :xs="24" :lg="8">
        <!-- ⑥ 广告位占位（顶部） -->
        <!-- Ad placeholder (top) -->
        <div class="panel ad-panel">
          <div class="ad-badge">{{ t('admin.adBadge') }}</div>
          <div class="ad-title">{{ t('admin.adTitle') }}</div>
          <div class="ad-desc">{{ t('admin.adDesc') }}</div>
        </div>

        <!-- ⑦ 模型用量分布（堆叠柱图） -->
        <div class="panel">
          <div class="panel-head">
            <h2 class="panel-title">{{ t('admin.sectionModelDist') }}</h2>
            <span class="panel-extra">{{ t('admin.topN', { n: distItems.length || 0 }) }}</span>
          </div>
          <v-chart
            v-if="distItems.length > 0"
            :option="modelBarOption"
            :style="{ height: '220px' }"
            autoresize
          />
          <div v-else class="chart-empty">{{ t('admin.chartEmpty') }}</div>
        </div>

        <!-- ⑧ 活跃用户 TOP 20 -->
        <div class="panel no-pad">
          <div class="panel-head pad-head">
            <h2 class="panel-title">{{ t('admin.sectionLeaderboard') }}</h2>
            <div class="leaderboard-extra">
              <a-radio-group
                v-model="topRange"
                type="button"
                size="small"
                @change="loadTopUsers"
              >
                <a-radio value="today">{{ t('admin.rangeToday') }}</a-radio>
                <a-radio value="7d">{{ t('admin.range7d') }}</a-radio>
                <a-radio value="30d">{{ t('admin.range30d') }}</a-radio>
              </a-radio-group>
            </div>
          </div>
          <a-table
            :columns="topColumns"
            :data="topUsers"
            :pagination="false"
            :loading="topLoading"
            row-key="id"
            :bordered="false"
            :stripe="true"
            size="small"
            class="dash-table"
          >
            <template #rank="{ record }">
              <span class="rank-badge" :class="rankClass(record.rank)">{{ record.rank }}</span>
            </template>
            <template #username="{ record }">
              <span class="user-cell">
                <span class="user-name">{{ record.username }}</span>
                <span v-if="record.email" class="user-email">{{ record.email }}</span>
              </span>
            </template>
            <template #request_count="{ record }">
              <span class="cell-num">{{ numFmt(record.request_count) }}</span>
            </template>
            <template #quota="{ record }">
              <span class="cell-num">{{ formatQuota(record.quota) }}</span>
            </template>
            <template #balance="{ record }">
              <span class="cell-num cell-num-muted">{{ formatQuota(record.balance) }}</span>
            </template>
            <template #current_plan_name="{ record }">
              <a-tag v-if="record.current_plan_name" color="arcoblue" size="small">
                {{ record.current_plan_name }}
              </a-tag>
              <span v-else class="cell-muted">—</span>
            </template>
            <template #empty>
              <div class="empty-state">{{ t('admin.noLeaderboard') }}</div>
            </template>
          </a-table>
        </div>

        <!-- ⑨ 系统公告 -->
        <div class="panel clickable" @click="$router.push('/log')">
          <div class="panel-head">
            <h2 class="panel-title">{{ t('admin.sectionAnnouncement') }}</h2>
            <span class="panel-link">{{ t('admin.viewLog') }} →</span>
          </div>
          <div class="notice-list">
            <a v-for="(n, idx) in announcements" :key="idx" class="notice-item">
              <span class="notice-title">{{ n.title }}</span>
              <span class="notice-time">{{ n.time }}</span>
            </a>
          </div>
        </div>

        <!-- ⑩ 更新日志 -->
        <div class="panel">
          <div class="panel-head">
            <h2 class="panel-title">{{ t('admin.sectionChangelog') }}</h2>
          </div>
          <div class="changelog">
            <div v-for="(c, idx) in changelog" :key="idx" class="cl-item">
              <div class="cl-head">
                <span class="cl-version">{{ c.version }}</span>
                <span class="cl-date">{{ c.date }}</span>
              </div>
              <p class="cl-desc">{{ c.desc }}</p>
            </div>
          </div>
        </div>

        <!-- ⑪ 资源 -->
        <div class="panel">
          <div class="panel-head">
            <h2 class="panel-title">{{ t('admin.sectionResources') }}</h2>
          </div>
          <div class="contact-list">
            <a
              v-for="(r, idx) in contactLinks"
              :key="idx"
              :href="r.href"
              target="_blank"
              class="contact-item"
            >
              <span class="contact-label">{{ r.label }}</span>
              <span class="contact-value is-active">{{ r.value }}</span>
              <icon-launch class="contact-arrow" :size="14" />
            </a>
          </div>
        </div>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import {
  IconRefresh,
  IconUserAdd,
  IconUserGroup,
  IconHeart,
  IconCloseCircle,
  IconDelete,
  IconCode,
  IconApps,
  IconStorage,
  IconGift,
  IconCalendar,
  IconBarChart,
  IconDashboard,
  IconSubscribe,
  IconSubscribeAdd,
  IconArrowFall,
  IconLaunch,
} from '@arco-design/web-vue/es/icon'
import VChart from 'vue-echarts'
import 'echarts'
import adminApi from '@/api/admin'

const { t } = useI18n()
const router = useRouter()

// ---------- 状态 ----------
const range = ref('7d')
const topRange = ref('7d')
const overviewLoading = ref(false)
const topLoading = ref(false)
const distLoading = ref(false)
const overview = ref(null)
const topUsers = ref([])
const distDays = ref([])
const distItems = ref([])
const lastRefreshAt = ref(0)

// ---------- 公告 / 更新日志 / 资源（管理员仪表盘固定列表） ----------
const announcements = [
  { title: t('admin.announcement1Title'), time: t('admin.announcement1Time') },
  { title: t('admin.announcement2Title'), time: t('admin.announcement2Time') },
  { title: t('admin.announcement3Title'), time: t('admin.announcement3Time') },
]
const changelog = [
  { version: 'v0.0.16', date: '2026-09-11', desc: t('admin.changelogCurrentDesc') },
  { version: 'v0.0.10', date: '2026-09-06', desc: t('admin.changelogTopupDesc') },
  { version: 'v0.0.0', date: '2026-06-01', desc: t('admin.changelogInitialDesc') },
]
const contactLinks = [
  { label: t('admin.contactDocsLabel'), value: 'one-api.pro', href: 'http://one-api.pro' },
  { label: t('admin.contactGithubLabel'), value: 'modelbus/one-api-pro', href: 'https://github.com/modelbus/one-api-pro' },
]

// ---------- 工具 ----------
// numFmt: 安全 toLocaleString，null/undefined 视为 0。
// numFmt: safe toLocaleString, treat null/undefined as 0.
function numFmt(n) {
  const v = Number(n) || 0
  return v.toLocaleString()
}

// formatQuota: 按 1e3/1e6/1e9 自动换算到 K/M/B。
// formatQuota: auto-suffix K/M/B for big numbers.
function formatQuota(n) {
  const num = Number(n) || 0
  if (num < 10000) return num.toLocaleString()
  if (num < 1e6) return (num / 1e3).toFixed(2) + 'K'
  if (num < 1e9) return (num / 1e6).toFixed(2) + 'M'
  return (num / 1e9).toFixed(2) + 'B'
}

function fmtMoney(n) {
  const num = Number(n) || 0
  return num.toFixed(2)
}

const todayStr = computed(() => {
  const d = new Date()
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
})

const lastRefreshLabel = computed(() => {
  if (!lastRefreshAt.value) return '--:--:--'
  const d = new Date(lastRefreshAt.value * 1000)
  const pad = (n) => String(n).padStart(2, '0')
  return `${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
})

const rangeLabel = computed(() => {
  switch (range.value) {
    case 'today': return t('admin.rangeToday')
    case '30d': return t('admin.range30d')
    case 'all': return t('admin.rangeAll')
    default: return t('admin.range7d')
  }
})

// 卡片配色：与 Dashboard.vue 一致的 4 色。
// KPI palette consistent with Dashboard.vue.
const palette = {
  blue: { bg: 'rgba(22,93,255,0.08)', color: '#165dff' },
  green: { bg: 'rgba(0,180,42,0.08)', color: '#00b42a' },
  purple: { bg: 'rgba(114,46,209,0.08)', color: '#722ed1' },
  orange: { bg: 'rgba(255,125,0,0.08)', color: '#ff7d00' },
  red: { bg: 'rgba(245,63,63,0.08)', color: '#f53f3f' },
  cyan: { bg: 'rgba(15,198,194,0.08)', color: '#0fc6c2' },
  gray: { bg: 'rgba(134,144,156,0.08)', color: '#86909c' },
}

// 模型分布柱图配色（堆叠）
const modelColors = ['#165dff', '#00b42a', '#722ed1', '#ff7d00', '#0fc6c2', '#f53f3f', '#f7ba1e', '#86909c']

// ---------- KPI 卡数据 ----------
// getNum: 防御性读取，undefined → 0。
// getNum: defensive read, undefined → 0.
function getNum(obj, key) {
  return Number(obj?.[key] ?? 0)
}

const userStatItems = computed(() => {
  const u = overview.value?.users
  return [
    { label: t('admin.userTotal'), value: numFmt(getNum(u, 'total')), icon: IconUserGroup, ...palette.blue },
    { label: t('admin.userNewToday'), value: numFmt(getNum(u, 'new_today')), icon: IconUserAdd, ...palette.green },
    { label: t('admin.userNew7d'), value: numFmt(getNum(u, 'new_7d')), icon: IconUserAdd, ...palette.green },
    { label: t('admin.userNew30d'), value: numFmt(getNum(u, 'new_30d')), icon: IconUserAdd, ...palette.green },
    { label: t('admin.userActive7d'), value: numFmt(getNum(u, 'active_7d')), icon: IconHeart, ...palette.purple },
    { label: t('admin.userDisabled'), value: numFmt(getNum(u, 'disabled')), icon: IconCloseCircle, ...palette.gray },
    { label: t('admin.userDeleted'), value: numFmt(getNum(u, 'deleted')), icon: IconDelete, ...palette.gray },
    { label: t('admin.userTodayRatio'), value: ratioText(getNum(u, 'new_today'), getNum(u, 'total')), icon: IconDashboard, ...palette.cyan },
  ]
})

const resourceStatItems = computed(() => {
  const ov = overview.value
  return [
    {
      label: t('admin.tokenTotal'), icon: IconCode, ...palette.blue,
      value: numFmt(getNum(ov?.tokens, 'total')),
      foot: `${t('admin.enabled')}: ${numFmt(getNum(ov?.tokens, 'enabled'))}`,
    },
    {
      label: t('admin.channelTotal'), icon: IconApps, ...palette.cyan,
      value: numFmt(getNum(ov?.channels, 'total')),
      foot: `${t('admin.enabled')}: ${numFmt(getNum(ov?.channels, 'enabled'))}`,
    },
    {
      label: t('admin.planTotal'), icon: IconStorage, ...palette.purple,
      value: numFmt(getNum(ov?.plans, 'total')),
      foot: `${t('admin.enabled')}: ${numFmt(getNum(ov?.plans, 'enabled'))}`,
      to: '/setting/plan',
    },
    {
      label: t('admin.redemptionTotal'), icon: IconGift, ...palette.orange,
      value: numFmt(getNum(ov?.redemptions, 'total')),
      foot: `${t('admin.used')}: ${numFmt(getNum(ov?.redemptions, 'used'))} / ${t('admin.unused')}: ${numFmt(getNum(ov?.redemptions, 'unused'))}`,
      to: '/redemption',
    },
    {
      label: t('admin.subscriptionTotal'), icon: IconCalendar, ...palette.green,
      value: numFmt(getNum(ov?.subscriptions, 'total')),
      foot: `${t('admin.active')}: ${numFmt(getNum(ov?.subscriptions, 'active'))} / ${t('admin.expired')}: ${numFmt(getNum(ov?.subscriptions, 'expired'))}`,
    },
  ]
})

const quotaStatItems = computed(() => {
  const q = overview.value?.quota
  return [
    { label: t('admin.quotaToday'), icon: IconBarChart, ...palette.blue,
      value: formatQuota(getNum(q, 'today')),
      foot: t('admin.quotaUnit') },
    { label: t('admin.quota7d'), icon: IconBarChart, ...palette.purple,
      value: formatQuota(getNum(q, 'week')),
      foot: t('admin.quotaUnit') },
    { label: t('admin.quota30d'), icon: IconBarChart, ...palette.orange,
      value: formatQuota(getNum(q, 'month')),
      foot: t('admin.quotaUnit') },
    { label: t('admin.quotaTotal'), icon: IconBarChart, ...palette.cyan,
      value: formatQuota(getNum(q, 'total')),
      foot: t('admin.quotaUnit') },
  ]
})

const revenueStatItems = computed(() => {
  const r = overview.value?.revenue
  return [
    { label: t('admin.revenueTotal'), icon: IconDashboard, ...palette.blue,
      value: fmtMoney(getNum(r, 'total')),
      foot: t('admin.revenueTotalFoot') },
    { label: t('admin.revenueTopup'), icon: IconSubscribeAdd, ...palette.green,
      value: fmtMoney(getNum(r, 'topup')),
      foot: t('admin.revenueTopupFoot') },
    { label: t('admin.revenueSubscription'), icon: IconSubscribe, ...palette.purple,
      value: fmtMoney(getNum(r, 'subscription')),
      foot: t('admin.revenueSubscriptionFoot') },
    { label: t('admin.revenueRefund'), icon: IconArrowFall, ...palette.red,
      value: fmtMoney(getNum(r, 'refund')),
      foot: t('admin.revenueRefundFoot') },
  ]
})

const trendItems = computed(() => {
  const trends = overview.value?.trends || {}
  const req = trends.requests || []
  const qua = trends.quota || []
  return [
    { field: 'requests', label: t('admin.chartRequests'), color: '#165dff', total: numFmt(sumBy(req, 'count')) },
    { field: 'quota', label: t('admin.chartQuota'), color: '#722ed1', total: formatQuota(sumBy(qua, 'quota')) },
  ]
})

function sumBy(arr, key) {
  return (arr || []).reduce((acc, p) => acc + (Number(p?.[key]) || 0), 0)
}

function ratioText(num, denom) {
  const a = Number(num) || 0
  const b = Number(denom) || 0
  if (b <= 0) return '0%'
  const pct = (a / b) * 100
  return pct >= 10 ? pct.toFixed(0) + '%' : pct.toFixed(1) + '%'
}

function rankClass(rank) {
  if (rank === 1) return 'rank-1'
  if (rank === 2) return 'rank-2'
  if (rank === 3) return 'rank-3'
  return 'rank-default'
}

// ---------- 折线图（与 Dashboard.vue lineOption 同构） ----------
// 修复：当 7 天数据全是 0 时，y 轴 auto-scale 会把线压成一条不可见的线。
// 显式设 min=0 并加大 symbol，保证哪怕是 0 数据也能看到「今天的点」。
function lineOption(field, color) {
  const points = overview.value?.trends?.[field] || []
  const values = points.map((d) => (field === 'quota' ? d.quota : d.count))
  const maxV = values.length ? Math.max(...values) : 0
  return {
    grid: { left: 36, right: 12, top: 18, bottom: 22 },
    tooltip: {
      trigger: 'axis',
      confine: true,
      backgroundColor: '#1d2129',
      borderWidth: 0,
      textStyle: { color: '#fff', fontSize: 12 },
    },
    xAxis: {
      type: 'category', boundaryGap: false,
      data: points.map((d) => (d.day || '').slice(5)),
      axisLabel: { fontSize: 11, color: '#86909c' },
      axisLine: { show: false },
      axisTick: { show: false },
    },
    yAxis: {
      type: 'value',
      min: 0,
      // 数据全 0 时给一个最小刻度上限，避免线被压扁到看不见
      max: maxV > 0 ? undefined : 1,
      splitNumber: 4,
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
      type: 'line', smooth: true, symbol: 'circle', symbolSize: 6,
      connectNulls: true,
      data: values,
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

// ---------- 模型分布堆叠柱图 ----------
const modelBarOption = computed(() => {
  const days = distDays.value
  const items = distItems.value
  const series = items.map((m, idx) => ({
    name: m.model_name,
    type: 'bar',
    stack: 'total',
    barMaxWidth: 22,
    itemStyle: { color: modelColors[idx % modelColors.length], borderRadius: [2, 2, 0, 0] },
    emphasis: { focus: 'series' },
    data: m.day_quota || [],
  }))
  return {
    color: modelColors,
    grid: { left: 32, right: 8, top: 8, bottom: 28 },
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' },
      backgroundColor: '#1d2129',
      borderWidth: 0,
      textStyle: { color: '#fff', fontSize: 12 },
    },
    legend: {
      show: items.length > 0,
      bottom: 0,
      left: 'center',
      type: 'scroll',
      icon: 'roundRect',
      itemWidth: 10,
      itemHeight: 10,
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
      axisLabel: {
        fontSize: 11,
        color: '#86909c',
        formatter: (v) => (v >= 1000 ? `${(v / 1000).toFixed(0)}k` : v),
      },
      splitLine: { lineStyle: { color: '#f2f3f5', type: 'dashed' } },
      axisLine: { show: false },
      axisTick: { show: false },
    },
    series,
  }
})

// ---------- 排行榜 ----------
const topColumns = computed(() => [
  { title: '#', slotName: 'rank', width: 52, align: 'center' },
  { title: t('admin.colUsername'), slotName: 'username' },
  { title: t('admin.colRequestCount'), slotName: 'request_count', width: 88, align: 'right' },
  { title: t('admin.colQuota'), slotName: 'quota', width: 88, align: 'right' },
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

async function loadModelDistribution() {
  distLoading.value = true
  try {
    const { data } = await adminApi.modelDistribution(range.value, 8)
    if (data?.success) {
      distDays.value = data.data?.days || []
      distItems.value = data.data?.items || []
    } else {
      Message.error(data?.message || t('admin.loadFailed'))
    }
  } catch (e) {
    Message.error(t('admin.loadFailed') + ': ' + (e?.message || ''))
  } finally {
    distLoading.value = false
  }
}

function loadAll() {
  loadOverview()
  loadTopUsers()
  loadModelDistribution()
}

onMounted(loadAll)
</script>

<style scoped>
/* ============================================================
   注意：Vue <style scoped> 会给每个组件生成独立的 class hash，
   Dashboard.vue 里的同名样式不会作用到本组件。
   因此把本组件实际用到的样式内联在这里（与 Dashboard.vue 同步）。
   Note: <style scoped> produces per-component class hashes.
   Dashboard.vue's same-name classes won't apply here, so we
   inline the same style rules in this file.
   ============================================================ */

.admin-dashboard {
  padding: 4px 4px 8px;
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
  align-items: center;
  gap: 8px;
}
.meta-chip {
  font-size: 12px;
  color: var(--color-text-3);
  background: var(--color-fill-2);
  padding: 3px 10px;
  border-radius: 4px;
}
.meta-chip-mono {
  font-variant-numeric: tabular-nums;
  font-family: 'SF Mono', Menlo, Consolas, monospace;
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

/* ============ 核心指标 ============ */
.stat-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  padding: 18px 20px;
}
/* 用户层 8 个 stat-item 排成 4 列 × 2 行：覆盖 stat-grid 缺省 4 列 + nth-child 边线 */
.admin-dashboard :deep(.stat-grid) {
  grid-template-columns: repeat(4, 1fr);
  padding: 18px 20px;
}
.admin-dashboard :deep(.stat-grid .stat-item:nth-child(4n)) {
  border-right: none;
  padding-right: 0;
}
.admin-dashboard :deep(.stat-grid .stat-item:nth-child(4n+1)) {
  padding-left: 0;
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

/* stat-item-clickable：点击感的 hover 高亮 */
.stat-item-clickable {
  cursor: pointer;
  border-radius: 6px;
  margin: -6px;
  padding: 6px;
  transition: background 0.15s;
}
.stat-item-clickable:hover {
  background: var(--color-fill-1);
}

/* ============ 趋势图 ============ */
.trend-divider {
  height: 1px;
  background: var(--color-fill-3);
  margin: 18px 0 16px;
}
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
.cell-num {
  font-variant-numeric: tabular-nums;
  font-weight: 500;
  color: var(--color-text-1);
}
.cell-num-muted {
  font-variant-numeric: tabular-nums;
  font-weight: 400;
  color: var(--color-text-3);
}
.cell-mono {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 12px;
  color: var(--color-text-2);
}
.cell-muted {
  color: var(--color-text-4);
  font-size: 13px;
}
.empty-state {
  padding: 40px 0;
  text-align: center;
  color: var(--color-text-4);
  font-size: 13px;
}

/* ============ 排行榜 extra ============ */
.leaderboard-extra {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* ============ 排行榜 rank 徽章 ============ */
.rank-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  font-size: 12px;
  font-weight: 600;
  background: var(--color-fill-2);
  color: var(--color-text-3);
}
.rank-badge.rank-1 {
  background: rgba(255, 125, 0, 0.12);
  color: #ff7d00;
}
.rank-badge.rank-2 {
  background: rgba(22, 93, 255, 0.12);
  color: #165dff;
}
.rank-badge.rank-3 {
  background: rgba(0, 180, 42, 0.12);
  color: #00b42a;
}

/* ============ 排行榜：用户名 + 邮箱堆叠 ============ */
.user-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}
.user-name {
  font-size: 13px;
  color: var(--color-text-1);
  font-weight: 500;
}
.user-email {
  font-size: 12px;
  color: var(--color-text-3);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* ============ 右侧栏：通用 ============ */
.panel.clickable {
  cursor: pointer;
  transition: border-color 0.15s, box-shadow 0.15s, transform 0.15s;
}
.panel.clickable:hover {
  border-color: rgb(var(--primary-5));
}
.panel-link {
  font-size: 12px;
  color: rgb(var(--primary-6));
  cursor: pointer;
  text-decoration: none;
}
.chart-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 180px;
  color: var(--color-text-4);
  font-size: 13px;
  background: var(--color-fill-1);
  border-radius: 6px;
}

/* ============ 广告位占位 ============ */
.ad-panel {
  background: linear-gradient(135deg, #165dff 0%, #4080ff 50%, #722ed1 100%);
  border: none;
  color: #fff;
  position: relative;
  overflow: hidden;
}
.ad-panel::before {
  content: '';
  position: absolute;
  inset: 0;
  background:
    radial-gradient(circle at 90% 10%, rgba(255, 255, 255, 0.18), transparent 40%),
    radial-gradient(circle at 10% 90%, rgba(255, 255, 255, 0.12), transparent 50%);
  pointer-events: none;
}
.ad-badge {
  display: inline-block;
  font-size: 11px;
  background: rgba(255, 255, 255, 0.18);
  padding: 2px 8px;
  border-radius: 3px;
  margin-bottom: 8px;
  position: relative;
  z-index: 1;
}
.ad-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 4px;
  position: relative;
  z-index: 1;
}
.ad-desc {
  font-size: 12px;
  opacity: 0.85;
  position: relative;
  z-index: 1;
}

/* ============ 系统公告 ============ */
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
.notice-title {
  font-size: 13px;
  color: var(--color-text-2);
}
.notice-item:hover .notice-title {
  color: rgb(var(--primary-6));
}
.notice-time {
  font-size: 11px;
  color: var(--color-text-4);
  flex-shrink: 0;
  margin-left: 8px;
  font-variant-numeric: tabular-nums;
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
.contact-value.is-active {
  color: rgb(var(--primary-6));
  font-weight: 500;
}
.contact-arrow {
  color: rgb(var(--primary-6));
  opacity: 0.6;
  flex-shrink: 0;
  transition: opacity 0.15s;
}
.contact-item:hover .contact-value {
  color: rgb(var(--primary-6));
}
.contact-item:hover .contact-arrow {
  opacity: 1;
}

/* ============ 响应式：lg 以下左侧占满，右侧下沉 ============ */
@media (max-width: 1200px) {
  .admin-dashboard :deep(.stat-grid) {
    grid-template-columns: repeat(4, 1fr);
  }
}
@media (max-width: 992px) {
  .admin-dashboard :deep(.stat-grid) {
    grid-template-columns: repeat(2, 1fr);
    gap: 20px 16px;
  }
  .admin-dashboard :deep(.stat-grid .stat-item) {
    border-right: none;
    padding-right: 0;
  }
  .admin-dashboard :deep(.stat-grid .stat-item:nth-child(odd)) {
    padding-right: 16px;
    border-right: 1px solid var(--color-fill-3);
  }
}
</style>
