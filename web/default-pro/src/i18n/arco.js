// arco.js Arco Design 组件库语言包映射
// Maps the active app locale to an Arco Design locale bundle.
// 版本: v0.0.21
// 日期: 2026-09-17
// 作者: opencode
import zhCN from '@arco-design/web-vue/es/locale/lang/zh-cn'
import enUS from '@arco-design/web-vue/es/locale/lang/en-us'

// arcoLocales 应用语言 -> Arco 语言包
// App locale -> Arco locale bundle
// 版本: v0.0.21
// 日期: 2026-09-17
export const arcoLocales = {
  zh: zhCN,
  en: enUS,
}

// getArcoLocale 根据应用语言返回对应的 Arco 语言包
// Returns the Arco locale bundle for the given app locale.
// 版本: v0.0.21
// 日期: 2026-09-17
export function getArcoLocale(locale) {
  return arcoLocales[locale] || zhCN
}

export default arcoLocales
