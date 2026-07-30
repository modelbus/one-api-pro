import { createI18n } from 'vue-i18n'
import messages from './locales.js'

const savedLang = localStorage.getItem('lang') || 'zh'

const i18n = createI18n({
  legacy: false,
  locale: savedLang,
  fallbackLocale: 'zh',
  messages,
})

export default i18n
