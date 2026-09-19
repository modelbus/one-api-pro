<!--
  Layout.vue 自定义布局组件
  Custom layout that adds a "Docs" badge next to the navbar logo.

  版本: v0.0.21
  日期: 2026-09-19
  作者: opencode
-->
<script setup>
import DefaultTheme from 'vitepress/theme'

const { Layout } = DefaultTheme
</script>

<template>
  <Layout>
    <template #nav-bar-title-after>
      <span class="op-docs-badge">Docs</span>
    </template>
  </Layout>
</template>
