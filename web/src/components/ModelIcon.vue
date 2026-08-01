<template>
  <span class="model-icon-wrap" :style="{ width: sizeCss, height: sizeCss }">
    <img
      v-if="resolvedSrc"
      :src="resolvedSrc"
      :alt="name"
      class="model-icon-img"
      :style="{ width: sizeCss, height: sizeCss }"
      @error="onError"
    />
    <span v-else class="model-icon-fallback" :style="{ background: color || '#86909c', fontSize: fallbackFont }">
      {{ firstChar }}
    </span>
  </span>
</template>

<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  slug: { type: String, default: '' },
  name: { type: String, default: '' },
  size: { type: [Number, String], default: 18 },
  color: { type: String, default: '' },
})

const errored = ref(false)

const allSvgModules = import.meta.glob('../assets/lobehub/*.svg', { eager: true, query: '?url', import: 'default' })

function resolveAsset() {
  if (!props.slug) return null
  const key = `../assets/lobehub/${props.slug}.svg`
  return allSvgModules[key] || null
}

const resolvedSrc = computed(() => {
  if (errored.value) return null
  return resolveAsset()
})

function onError() {
  errored.value = true
}

const sizeCss = computed(() => (typeof props.size === 'number' ? `${props.size}px` : props.size))
const fallbackFont = computed(() => {
  const n = typeof props.size === 'number' ? props.size : parseInt(props.size, 10) || 18
  return `${Math.round(n * 0.5)}px`
})
const firstChar = computed(() => (props.name || props.slug || '?').charAt(0).toUpperCase())
</script>

<style scoped>
.model-icon-wrap {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  vertical-align: middle;
}
.model-icon-img {
  object-fit: contain;
  border-radius: 4px;
}
.model-icon-fallback {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  border-radius: 4px;
  color: #fff;
  font-weight: 700;
  line-height: 1;
}
</style>
