<template>
  <a-spin :loading="loading" class="setting-container">
    <div class="section-header">
      <h3>{{ $t('settingPage.cluster.title') }}</h3>
      <a-button type="primary" size="small" @click="openModal()"><template #icon><icon-plus /></template>{{ $t('settingPage.cluster.addNode') }}</a-button>
    </div>
    <a-table :columns="columns" :data="nodes" :pagination="false" row-key="node_id" size="medium">
      <template #status="{ record }">
        <a-tag :color="record.status===1?'green':'red'" size="small">{{ record.status===1?$t('settingPage.cluster.alive'):$t('settingPage.cluster.failed') }}</a-tag>
      </template>
      <template #last_heartbeat="{ record }">{{ formatTime(record.last_heartbeat) }}</template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="handlePing(record)" :loading="pingSet.has(record.node_id)">Ping</a-button>
          <a-button type="text" size="small" @click="openModal(record)">{{ $t('settingPage.cluster.edit') }}</a-button>
          <a-popconfirm :content="record.status===1?$t('settingPage.cluster.confirmDisable'):$t('settingPage.cluster.confirmEnable')" @ok="toggleNode(record)">
            <a-button type="text" size="small" :status="record.status===1?'warning':'success'">{{ record.status===1?$t('settingPage.cluster.disable'):$t('settingPage.cluster.enable') }}</a-button>
          </a-popconfirm>
          <a-popconfirm :content="$t('settingPage.cluster.confirmDelete')" @ok="delNode(record.node_id)">
            <a-button type="text" size="small" status="danger">{{ $t('settingPage.cluster.delete') }}</a-button>
          </a-popconfirm>
        </a-space>
      </template>
    </a-table>

    <a-modal v-model:visible="modalVisible" :title="editing?$t('settingPage.cluster.editNode'):$t('settingPage.cluster.addNode')" @ok="saveNode" :ok-loading="saving" width="500">
      <a-form :model="nodeForm" layout="vertical">
        <a-form-item :label="$t('settingPage.cluster.nodeId')" required>
          <a-input-number v-model="nodeForm.node_id" :disabled="editing" :min="1" :max="49" :style="{width:'100%'}" placeholder="1-49" />
        </a-form-item>
        <a-form-item :label="$t('settingPage.cluster.nodeName')" required>
          <a-input v-model="nodeForm.node_name" placeholder="node-cn" />
        </a-form-item>
        <a-form-item :label="$t('settingPage.cluster.address')" required>
          <a-input v-model="nodeForm.address" placeholder="https://cn.example.com" />
        </a-form-item>
        <a-form-item :label="$t('settingPage.cluster.secret')">
          <a-input-password v-model="nodeForm.secret" :placeholder="$t('settingPage.cluster.secretPlaceholder')" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-spin>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'
import api from '@/api'

const { t } = useI18n()

const loading = ref(false), nodes = ref([])
const modalVisible = ref(false), editing = ref(false), saving = ref(false)
const nodeForm = reactive({ node_id: 1, node_name: '', address: '', secret: '' })
const pingSet = ref(new Set())

const columns = computed(() => [
  { title: t('settingPage.cluster.colId'), dataIndex: 'node_id', width: 60 },
  { title: t('settingPage.cluster.colName'), dataIndex: 'node_name' },
  { title: t('settingPage.cluster.colAddress'), dataIndex: 'address', ellipsis: true },
  { title: t('settingPage.cluster.colStatus'), slotName: 'status', width: 80 },
  { title: t('settingPage.cluster.colLastHeartbeat'), slotName: 'last_heartbeat', width: 170 },
  { title: t('settingPage.cluster.colAction'), slotName: 'actions', width: 260 },
])

async function loadData() {
  loading.value = true
  try { const { data } = await api.get('/api/cluster_node/'); if (data.success) nodes.value = Array.isArray(data.data) ? data.data : [] } catch (e) { /* ignore */ }
  finally { loading.value = false }
}

function openModal(record) {
  editing.value = !!record
  Object.assign(nodeForm, record ? { node_id: record.node_id, node_name: record.node_name || '', address: record.address || '', secret: record.secret_key || '' } : { node_id: 1, node_name: '', address: '', secret: '' })
  modalVisible.value = true
}

async function saveNode() {
  saving.value = true
  try {
    const body = { node_id: nodeForm.node_id, node_name: nodeForm.node_name, address: nodeForm.address, secret: nodeForm.secret }
    const { data } = editing.value ? await api.put('/api/cluster_node/', body) : await api.post('/api/cluster_node/', body)
    if (data.success) { modalVisible.value = false; loadData() } else Message.error(data.message)
  } catch (e) { Message.error(t('settingPage.cluster.saveFailed')) } finally { saving.value = false }
}

async function handlePing(record) {
  pingSet.value.add(record.node_id)
  try { const { data } = await api.get(`/api/cluster_node/ping/${record.node_id}`); if (data.success) { Message.success(t('settingPage.cluster.pingSuccess', { name: record.node_name })); loadData() } else Message.warning(data.message) }
  catch (e) { Message.error(t('settingPage.cluster.pingFailed')) } finally { pingSet.value.delete(record.node_id) }
}

async function toggleNode(record) {
  try {
    const url = record.status === 1 ? `/api/cluster_node/${record.node_id}/` : `/api/cluster_node/${record.node_id}/enable`
    const method = record.status === 1 ? 'delete' : 'post'
    const { data } = await api[method](url)
    if (data.success) { Message.success(record.status === 1 ? t('settingPage.cluster.disabled') : t('settingPage.cluster.enabled')); loadData() } else Message.error(data.message)
  } catch (e) { Message.error(t('settingPage.cluster.opFailed')) }
}

async function delNode(id) { try { await api.delete(`/api/cluster_node/${id}/`); loadData() } catch (e) { Message.error(t('settingPage.cluster.deleteFailed')) } }

function formatTime(ts) { if (!ts) return '-'; return new Date(ts * 1000).toLocaleString() }

onMounted(() => { loadData() })
</script>

<style scoped>
.setting-container { padding: 4px 0; }
.section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px; }
.section-header h3 { font-size: 16px; font-weight: 600; color: var(--color-text-1); margin: 0; padding: 0; }
</style>
