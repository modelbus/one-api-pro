import { createRouter, createWebHistory } from 'vue-router'
import AdminLayout from '@/layouts/AdminLayout.vue'
import i18n from '@/i18n'

const routes = [
  {
    path: '/',
    name: 'Landing',
    component: () => import('@/views/Landing.vue'),
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/Register.vue'),
  },
  {
    path: '/reset',
    name: 'PasswordReset',
    component: () => import('@/views/auth/PasswordReset.vue'),
  },
  {
    path: '/reset/:token',
    name: 'PasswordResetConfirm',
    component: () => import('@/views/auth/PasswordResetConfirm.vue'),
  },
  {
    path: '/terms',
    name: 'TermsOfService',
    component: () => import('@/views/TermsOfService.vue'),
  },
  {
    path: '/privacy',
    name: 'PrivacyPolicy',
    component: () => import('@/views/PrivacyPolicy.vue'),
  },
  {
    path: '/oauth/github',
    name: 'GitHubOAuth',
    component: () => import('@/views/auth/GitHubOAuth.vue'),
  },
  {
    path: '/oauth/lark',
    name: 'LarkOAuth',
    component: () => import('@/views/auth/LarkOAuth.vue'),
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

const BASE_TITLE = 'ONE-API-PRO'

router.afterEach((to) => {
  if (to.meta && to.meta.title) {
    document.title = `${i18n.global.t(to.meta.title)}—${BASE_TITLE}`
  } else {
    document.title = `${BASE_TITLE}—${i18n.global.t('route.baseTitle')}`
  }
})

export default router
