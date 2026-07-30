import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  server: {
    port: 3001,
    proxy: {
      '/api': 'http://localhost:3000',
      '/v1': 'http://localhost:3000',
    },
  },
  build: {
    outDir: 'build',
    emptyOutDir: true,
  },
})
