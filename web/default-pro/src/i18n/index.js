import { watch } from 'vue'
import { createI18n } from 'vue-i18n'
import messages from './locales.js'

// SUPPORTED_LOCALES 支持的语言列表（供语言切换组件使用）
// Supported locale list (consumed by the language switcher).
// 版本: v0.0.21
// 日期: 2026-09-17
// 作者: opencode
export const SUPPORTED_LOCALES = [
  { value: 'zh', label: '简体中文' },
  { value: 'en', label: 'English' },
]

const SUPPORTED_VALUES = SUPPORTED_LOCALES.map((l) => l.value)

function normalizeLocale(locale) {
  return SUPPORTED_VALUES.includes(locale) ? locale : 'zh'
}

const savedLang = normalizeLocale(localStorage.getItem('lang'))

// syncHtmlLang 同步 <html lang>，让浏览器字体/断行等行为跟随语言
// Keeps <html lang> in sync so browser behaviors follow the active locale.
// 版本: v0.0.21
// 日期: 2026-09-17
function syncHtmlLang(locale) {
  if (typeof document !== 'undefined') {
    document.documentElement.setAttribute('lang', locale === 'zh' ? 'zh-CN' : locale)
  }
}

const i18n = createI18n({
  legacy: false,
  locale: savedLang,
  fallbackLocale: 'zh',
  messages,
})

syncHtmlLang(savedLang)

watch(
  () => i18n.global.locale.value,
  (val) => syncHtmlLang(val),
)

// setLocale 切换应用语言并持久化到 localStorage
// Switches the app locale and persists it to localStorage.
// 版本: v0.0.21
// 日期: 2026-09-17
export function setLocale(locale) {
  const next = normalizeLocale(locale)
  i18n.global.locale.value = next
  localStorage.setItem('lang', next)
  syncHtmlLang(next)
}

export default i18n
