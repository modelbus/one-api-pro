import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ArcoVue from '@arco-design/web-vue'
import App from './App.vue'
import router from './router'
import i18n from './i18n'
import '@arco-design/web-vue/dist/arco.css'
import './style.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(ArcoVue)
app.use(i18n)
app.mount('#app')
