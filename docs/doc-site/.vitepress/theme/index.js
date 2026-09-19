// doc-site/.vitepress/theme/index.js
// Custom VitePress theme: extends the default theme and layers a modern
// visual style (brand gradient, rounded cards, refined typography) plus a
// custom Layout that renders the "Docs" badge next to the navbar logo.

import DefaultTheme from 'vitepress/theme'
import Layout from './Layout.vue'
import './custom.css'

export default {
  extends: DefaultTheme,
  Layout,
}
