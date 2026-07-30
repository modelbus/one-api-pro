<template>
  <div class="chat-page">
    <a-spin :loading="loading">
      <div v-if="loadError" class="chat-error">
        <a-result status="error" title="加载失败" :subtitle="loadError">
          <template #extra>
            <a-button @click="fetchChatLink">重试</a-button>
          </template>
        </a-result>
      </div>

      <div v-else-if="chatLink" class="chat-frame-wrapper">
        <iframe :src="chatLink" class="chat-iframe" frameborder="0" allowfullscreen />
      </div>

      <div v-else-if="loaded" class="chat-empty">
        <a-empty description="暂无聊天链接配置" />
      </div>
    </a-spin>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useStatusStore } from '@/stores/status'
import api from '@/api'

const statusStore = useStatusStore()

const loading = ref(true)
const loaded = ref(false)
const loadError = ref('')
const chatLink = ref('')

onMounted(async () => {
  await fetchChatLink()
})

async function fetchChatLink() {
  loading.value = true
  loadError.value = ''
  try {
    if (!statusStore.loaded) {
      await statusStore.fetchStatus()
    }
    chatLink.value = statusStore.status?.chat_link || ''
    loaded.value = true
  } catch (e) {
    loadError.value = e.response?.data?.message || e.message || '加载失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.chat-page {
  height: calc(100vh - 64px);
  display: flex;
  flex-direction: column;
}

.chat-frame-wrapper {
  flex: 1;
  display: flex;
}

.chat-iframe {
  width: 100%;
  height: 100%;
  border: none;
}

.chat-error,
.chat-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
}
</style>
