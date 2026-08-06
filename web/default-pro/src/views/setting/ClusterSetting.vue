<template>
  <a-spin :loading="loading" class="setting-container">
    <div class="section-header">
      <h3>集群节点</h3>
      <a-button type="primary" size="small" @click="openModal()"><template #icon><icon-plus /></template>添加节点</a-button>
    </div>
    <a-table :columns="columns" :data="nodes" :pagination="false" row-key="node_id" size="medium">
      <template #status="{ record }">
        <a-tag :color="record.status===1?'green':'red'" size="small">{{ record.status===1?'存活':'失败' }}</a-tag>
      </template>
      <template #last_heartbeat="{ record }">{{ formatTime(record.last_heartbeat) }}</template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="handlePing(record)" :loading="pingSet.has(record.node_id)">Ping</a-button>
          <a-button type="text" size="small" @click="openModal(record)">编辑</a-button>
          <a-popconfirm :content="record.status===1?'确定禁用？':'确定启用？'" @ok="toggleNode(record)">
            <a-button type="text" size="small" :status="record.status===1?'warning':'success'">{{ record.status===1?'禁用':'启用' }}</a-button>
          </a-popconfirm>
          <a-popconfirm content="确定删除？" @ok="delNode(record.node_id)">
            <a-button type="text" size="small" status="danger">删除</a-button>
          </a-popconfirm>
        </a-space>
      </template>
    </a-table>

    <a-modal v-model:visible="modalVisible" :title="editing?'编辑节点':'添加节点'" @ok="saveNode" :ok-loading="saving" width="500">
      <a-form :model="nodeForm" layout="vertical">
        <a-form-item label="节点编号" required>
          <a-input-number v-model="nodeForm.node_id" :disabled="editing" :min="1" :max="49" :style="{width:'100%'}" placeholder="1-49" />
        </a-form-item>
        <a-form-item label="节点名称" required>
          <a-input v-model="nodeForm.node_name" placeholder="node-cn" />
        </a-form-item>
        <a-form-item label="地址" required>
          <a-input v-model="nodeForm.address" placeholder="https://cn.example.com" />
        </a-form-item>
        <a-form-item label="密钥">
          <a-input-password v-model="nodeForm.secret" placeholder="集群通信密钥" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-spin>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'
import api from '@/api'

const loading = ref(false), nodes = ref([])
const modalVisible = ref(false), editing = ref(false), saving = ref(false)
const nodeForm = reactive({ node_id: 1, node_name: '', address: '', secret: '' })
const pingSet = ref(new Set())

const columns = [
  { title: 'ID', dataIndex: 'node_id', width: 60 },
  { title: '名称', dataIndex: 'node_name' },
  { title: '地址', dataIndex: 'address', ellipsis: true },
  { title: '状态', slotName: 'status', width: 80 },
  { title: '最后心跳', slotName: 'last_heartbeat', width: 170 },
  { title: '操作', slotName: 'actions', width: 260 },
]

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
  } catch (e) { Message.error('保存失败') } finally { saving.value = false }
}

async function handlePing(record) {
  pingSet.value.add(record.node_id)
  try { const { data } = await api.get(`/api/cluster_node/ping/${record.node_id}`); if (data.success) { Message.success(`Ping ${record.node_name} 成功`); loadData() } else Message.warning(data.message) }
  catch (e) { Message.error('Ping失败') } finally { pingSet.value.delete(record.node_id) }
}

async function toggleNode(record) {
  try {
    const url = record.status === 1 ? `/api/cluster_node/${record.node_id}/` : `/api/cluster_node/${record.node_id}/enable`
    const method = record.status === 1 ? 'delete' : 'post'
    const { data } = await api[method](url)
    if (data.success) { Message.success(record.status === 1 ? '已禁用' : '已启用'); loadData() } else Message.error(data.message)
  } catch (e) { Message.error('操作失败') }
}

async function delNode(id) { try { await api.delete(`/api/cluster_node/${id}/`); loadData() } catch (e) { Message.error('删除失败') } }

function formatTime(ts) { if (!ts) return '-'; return new Date(ts * 1000).toLocaleString() }

onMounted(() => { loadData() })
</script>

<style scoped>
.setting-container { padding: 4px 0; }
.section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px; }
.section-header h3 { font-size: 16px; font-weight: 600; color: var(--color-text-1); margin: 0; padding: 0; }
</style>
