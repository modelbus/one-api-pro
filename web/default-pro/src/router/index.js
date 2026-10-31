import { watch } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import AdminLayout from '@/layouts/AdminLayout.vue'
import i18n from '@/i18n'
import { buildDocumentTitle, resolvePageTitle } from '@/utils/title'

const routes = [
  {
    path: '/',
    name: 'Landing',
    component: () => import('@/views/Landing.vue'),
    meta: { title: 'route.home' },
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { title: 'route.login' },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/Register.vue'),
    meta: { title: 'route.register' },
  },
  {
    path: '/reset',
    name: 'PasswordReset',
    component: () => import('@/views/auth/PasswordReset.vue'),
    meta: { title: 'route.passwordReset' },
  },
  {
    path: '/reset/:token',
    name: 'PasswordResetConfirm',
    component: () => import('@/views/auth/PasswordResetConfirm.vue'),
    meta: { title: 'route.passwordReset' },
  },
  {
    path: '/terms',
    name: 'TermsOfService',
    component: () => import('@/views/TermsOfService.vue'),
    meta: { title: 'route.terms' },
  },
  {
    path: '/privacy',
    name: 'PrivacyPolicy',
    component: () => import('@/views/PrivacyPolicy.vue'),
    meta: { title: 'route.privacy' },
  },
  {
    path: '/oauth/github',
    name: 'GitHubOAuth',
    component: () => import('@/views/auth/GitHubOAuth.vue'),
    meta: { title: 'route.login' },
  },
  {
    path: '/oauth/lark',
    name: 'LarkOAuth',
    component: () => import('@/views/auth/LarkOAuth.vue'),
    meta: { title: 'route.login' },
  },
  {
    path: '/',
    component: AdminLayout,
    children: [
      { path: 'dashboard', name: 'Dashboard', component: () => import('@/views/dashboard/Dashboard.vue'), meta: { title: 'route.dashboard' } },
      { path: 'admin/dashboard', name: 'AdminDashboard', component: () => import('@/views/admin/AdminDashboard.vue'), meta: { title: 'route.adminDashboard', admin: true } },
      { path: 'channel', name: 'Channel', component: () => import('@/views/channel/Channel.vue'), meta: { title: 'route.channel', admin: true } },
      { path: 'admin/orders', name: 'AdminOrders', component: () => import('@/views/admin/AdminOrders.vue'), meta: { title: 'route.adminOrders', admin: true } },
      { path: 'token', name: 'Token', component: () => import('@/views/token/Token.vue'), meta: { title: 'route.token' } },
      { path: 'user', name: 'User', component: () => import('@/views/user/User.vue'), meta: { title: 'route.user', admin: true } },
      { path: 'redemption', name: 'Redemption', component: () => import('@/views/redemption/Redemption.vue'), meta: { title: 'route.redemption', admin: true } },
      { path: 'log', name: 'Log', component: () => import('@/views/log/Log.vue'), meta: { title: 'route.log' } },
      { path: 'subscription', name: 'Subscription', component: () => import('@/views/subscription/Subscription.vue'), meta: { title: 'route.subscription' } },
      { path: 'plans', name: 'Plans', component: () => import('@/views/user/Plans.vue'), meta: { title: 'route.plans' } },
      { path: 'orders', name: 'Orders', component: () => import('@/views/user/Orders.vue'), meta: { title: 'route.orders' } },
      { path: 'setting', name: 'Setting', component: () => import('@/views/setting/Setting.vue'), meta: { title: 'route.setting' }, redirect: '/setting/system', children: [
        { path: 'system', name: 'SystemSetting', component: () => import('@/views/setting/SystemSetting.vue'), meta: { title: 'route.systemSetting' } },
        { path: 'cluster', name: 'ClusterSetting', component: () => import('@/views/setting/ClusterSetting.vue'), meta: { title: 'route.clusterSetting' } },
        { path: 'operation', name: 'OperationSetting', component: () => import('@/views/setting/OperationSetting.vue'), meta: { title: 'route.operationSetting' } },
        { path: 'payment', name: 'PaymentSetting', component: () => import('@/views/setting/PaymentSetting.vue'), meta: { title: 'route.paymentSetting' } },
        { path: 'pricing', name: 'PricingSetting', component: () => import('@/views/setting/PricingSetting.vue'), meta: { title: 'route.pricingSetting' } },
        { path: 'plan', name: 'PlanSetting', component: () => import('@/views/setting/PlanSetting.vue'), meta: { title: 'route.planSetting' } },
        { path: 'topup', name: 'TopupSetting', component: () => import('@/views/setting/TopupSetting.vue'), meta: { title: 'route.topupSetting' } },
        { path: 'personal', name: 'PersonalSetting', component: () => import('@/views/setting/PersonalSetting.vue'), meta: { title: 'route.personalSetting' } },
      ] },
      { path: 'redeem', name: 'Redeem', component: () => import('@/views/redeem/Redeem.vue'), meta: { title: 'route.redeem' } },
      { path: 'chat', name: 'Chat', component: () => import('@/views/chat/Chat.vue'), meta: { title: 'route.chat' } },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { title: 'route.notFound' },
  },
]

const router = createRouter({ history: createWebHistory(), routes })

router.beforeEach((to, from, next) => {
  const user = JSON.parse(localStorage.getItem('user'))
  const isLoggedIn = !!user
  const isAdminRoute = to.meta.admin || to.name === 'Channel' || to.name === 'User' || to.name === 'Redemption' || to.name === 'AdminOrders' || to.name === 'AdminDashboard'

  // 根路径始终展示落地页
  if (to.path === '/') {
    return next()
  }

  // session 过期被踢回登录页：清除本地 user，停留登录页，避免与 /dashboard 死循环
  if (to.path === '/login' && to.query.expired) {
    localStorage.removeItem('user')
    return next()
  }

  if (!isLoggedIn && !['/login', '/register', '/reset', '/terms', '/privacy'].includes(to.path) && !to.path.startsWith('/reset/') && !to.path.startsWith('/oauth/')) {
    return next('/login')
  }

  if (isLoggedIn && (to.path === '/login' || to.path === '/register')) {
    return next('/dashboard')
  }

  if (isAdminRoute && user && user.role < 10) {
    return next('/dashboard')
  }

  next()
})

// applyDocumentTitle 依据「当前路由 meta.title → i18n → 站点名」渲染 Tab 标题
// 标题只由路由与页面元信息驱动，任何全局初始化（如 /api/status）都不得改写它。
// Renders the tab title from the active route meta; global bootstrap must never override it.
function applyDocumentTitle(to) {
  if (typeof document === 'undefined') return
  const pageTitle = resolvePageTitle(to?.meta?.title, (key) => i18n.global.t(key))
  document.title = buildDocumentTitle(pageTitle, i18n.global.t('route.baseTitle'))
}

router.afterEach((to) => {
  applyDocumentTitle(to)
})

// 切换语言后同一路由的标题文案需要跟随刷新（路由不会重新导航）
// Keep the title in sync when the locale changes without a navigation.
watch(
  () => i18n.global.locale.value,
  () => applyDocumentTitle(router.currentRoute.value),
)

export default router
