<template>
  <a-card :bordered="false" class="data-card">
    <PageHeader :title="title" :description="description">
      <template #extra>
        <a-space>
          <a-input-search
            v-if="searchable"
            v-model="keyword"
            :placeholder="searchPlaceholder"
            allow-clear
            @search="handleSearch"
            @clear="handleClear"
            style="width: 240px"
          />
          <slot name="extra" />
        </a-space>
      </template>
    </PageHeader>

    <a-spin :loading="loading">
      <div v-if="!loading" class="table-wrapper">
        <a-table
          :data="showData"
          :columns="columns"
          :pagination="false"
          :bordered="{ wrapper: true, cell: false }"
          :stripe="true"
          size="medium"
          row-key="id"
        >
          <template v-for="(_, slot) in $slots" :key="slot" #[slot]="scope">
            <slot :name="slot" v-bind="scope" />
          </template>
        </a-table>
        <div class="table-footer" v-if="allData.length > pageSize">
          <a-pagination
            :total="total"
            :current="currentPage"
            :page-size="pageSize"
            :page-size-options="[10, 20, 50]"
            show-total
            show-page-size
            @change="handlePageChange"
            @page-size-change="handlePageSizeChange"
          />
        </div>
      </div>
    </a-spin>

    <div v-if="!loading && !showData.length" class="empty-state">
      <a-empty :description="emptyText" />
    </div>
  </a-card>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import PageHeader from './PageHeader.vue'

const props = defineProps({
  title: { type: String, default: '' },
  description: { type: String, default: '' },
  columns: { type: Array, default: () => [] },
  data: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
  searchable: { type: Boolean, default: true },
  searchPlaceholder: { type: String, default: '搜索...' },
  emptyText: { type: String, default: '暂无数据' },
})

const keyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const filtered = ref([])

const allData = computed(() => filtered.value.length ? filtered.value : props.data)
const total = computed(() => allData.value.length)
const showData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return allData.value.slice(start, start + pageSize.value)
})

const handleSearch = () => {
  currentPage.value = 1
  if (keyword.value) {
    filtered.value = props.data.filter(item =>
      Object.values(item).some(v => String(v).toLowerCase().includes(keyword.value.toLowerCase()))
    )
  } else {
    filtered.value = []
  }
}

const handleClear = () => {
  keyword.value = ''
  filtered.value = []
  currentPage.value = 1
}

const handlePageChange = (page) => { currentPage.value = page }
const handlePageSizeChange = (size) => { pageSize.value = size; currentPage.value = 1 }

defineExpose({ keyword, currentPage, pageSize })
</script>

<style scoped>
.data-card {
  border-radius: 4px;
}

.table-wrapper {
  min-height: 200px;
}

.table-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid var(--color-border-2);
}

.empty-state {
  padding: 60px 0;
}
</style>
