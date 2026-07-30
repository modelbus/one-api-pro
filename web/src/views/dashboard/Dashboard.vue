<template>
  <div class="dashboard-page">
    <a-spin :loading="loading" tip="Loading dashboard...">
      <a-row :gutter="16" class="stats-row">
        <a-col :span="8">
          <a-card :bordered="false" class="stat-card">
            <a-statistic
              title="Total Tokens Used"
              :value="dashboard.total_tokens ?? 0"
              :value-style="{ color: '#165DFF' }"
            />
          </a-card>
        </a-col>
        <a-col :span="8">
          <a-card :bordered="false" class="stat-card">
            <a-statistic
              title="Total Quota Consumed"
              :value="dashboard.total_quota ?? 0"
              :precision="4"
              :value-style="{ color: '#00B42A' }"
            />
          </a-card>
        </a-col>
        <a-col :span="8">
          <a-card :bordered="false" class="stat-card">
            <a-statistic
              title="Request Count"
              :value="dashboard.request_count ?? 0"
              :value-style="{ color: '#F7BA1E' }"
            />
          </a-card>
        </a-col>
      </a-row>

      <a-divider />

      <a-row :gutter="16" class="section-row" v-if="subscription">
        <a-col :span="24">
          <a-card title="Subscription" :bordered="false">
            <template #extra>
              <a-tag :color="subStatusColor(subscription.status)">
                {{ subscription.status }}
              </a-tag>
            </template>

            <a-row :gutter="16">
              <a-col :span="8">
                <a-statistic title="Plan" :value="subscription.plan_name ?? '-'" />
              </a-col>
              <a-col :span="8">
                <a-statistic title="Quota Used" :value="subscription.quota_used ?? 0" :precision="2" />
              </a-col>
              <a-col :span="8">
                <a-statistic title="Quota Limit" :value="subscription.quota_limit ?? 0" :precision="2" />
              </a-col>
            </a-row>

            <a-row :gutter="16" style="margin-top: 16px">
              <a-col :span="12">
                <a-statistic title="Start Time" :value="subscription.start_time ?? '-'" />
              </a-col>
              <a-col :span="12">
                <a-statistic title="End Time" :value="subscription.end_time ?? '-'" />
              </a-col>
            </a-row>
          </a-card>
        </a-col>
      </a-row>

      <a-row v-else-if="!loading" :gutter="16" class="section-row">
        <a-col :span="24">
          <a-card :bordered="false">
            <a-empty description="No active subscription" />
          </a-card>
        </a-col>
      </a-row>
    </a-spin>
  </div>
</template>

<script setup>
import { toastSuccess, toastError, toastWarning } from '@/helpers/toast'
import { ref, reactive, onMounted } from 'vue'
import api from '@/api'

const loading = ref(true)
const dashboard = reactive({
  total_tokens: 0,
  total_quota: 0,
  request_count: 0,
})
const subscription = ref(null)

const subStatusColor = (status) => {
  switch (status) {
    case 'active': return 'green'
    case 'expired': return 'red'
    case 'pending': return 'orange'
    default: return 'gray'
  }
}

onMounted(async () => {
  try {
    const [dashRes, subRes] = await Promise.all([
      api.get('/api/user/dashboard'),
      api.get('/api/subscription/self'),
    ])
    if (dashRes.data) {
      Object.assign(dashboard, dashRes.data)
    }
    if (subRes.data) {
      subscription.value = Array.isArray(subRes.data)
        ? subRes.data[0]
        : subRes.data
    }
  } catch (err) {
    import('@arco-design/web-vue').then(m => m.Message.error)(err?.response?.data?.message || 'Failed to load dashboard')
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.dashboard-page {
  padding: 16px;
}
.stats-row {
  margin-bottom: 8px;
}
.stat-card {
  text-align: center;
}
.section-row {
  margin-top: 8px;
}
</style>
