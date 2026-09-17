<template>
  <div class="not-found-page">
    <a-result status="404" title="404" :subtitle="$t('notFoundPage.subtitle')">
      <template #extra>
        <a-button type="primary" @click="$router.push('/')">{{ $t('notFoundPage.backHome') }}</a-button>
      </template>
    </a-result>
  </div>
</template>

<script setup>
</script>

<style scoped>
.not-found-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
