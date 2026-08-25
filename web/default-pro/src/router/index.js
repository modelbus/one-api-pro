import { createRouter, createWebHistory } from 'vue-router'
import AdminLayout from '@/layouts/AdminLayout.vue'

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
      { path: 'dashboard', name: 'Dashboard', component: () => import('@/views/dashboard/Dashboard.vue'), meta: { title: '仪表盘' } },
      { path: 'channel', name: 'Channel', component: () => import('@/views/channel/Channel.vue'), meta: { title: '渠道', admin: true } },
      { path: 'token', name: 'Token', component: () => import('@/views/token/Token.vue'), meta: { title: '令牌' } },
      { path: 'user', name: 'User', component: () => import('@/views/user/User.vue'), meta: { title: '用户', admin: true } },
      { path: 'redemption', name: 'Redemption', component: () => import('@/views/redemption/Redemption.vue'), meta: { title: '兑换码', admin: true } },
      { path: 'log', name: 'Log', component: () => import('@/views/log/Log.vue'), meta: { title: '日志' } },
      { path: 'subscription', name: 'Subscription', component: () => import('@/views/subscription/Subscription.vue'), meta: { title: '订阅' } },
      { path: 'plans', name: 'Plans', component: () => import('@/views/user/Plans.vue'), meta: { title: '套餐' } },
      { path: 'orders', name: 'Orders', component: () => import('@/views/user/Orders.vue'), meta: { title: '订单' } },
      { path: 'setting', name: 'Setting', component: () => import('@/views/setting/Setting.vue'), meta: { title: '设置' }, redirect: '/setting/system', children: [
        { path: 'system', name: 'SystemSetting', component: () => import('@/views/setting/SystemSetting.vue'), meta: { title: '系统设置' } },
        { path: 'cluster', name: 'ClusterSetting', component: () => import('@/views/setting/ClusterSetting.vue'), meta: { title: '集群设置' } },
        { path: 'operation', name: 'OperationSetting', component: () => import('@/views/setting/OperationSetting.vue'), meta: { title: '运营设置' } },
        { path: 'payment', name: 'PaymentSetting', component: () => import('@/views/setting/PaymentSetting.vue'), meta: { title: '支付设置' } },
        { path: 'pricing', name: 'PricingSetting', component: () => import('@/views/setting/PricingSetting.vue'), meta: { title: '定价管理' } },
        { path: 'plan', name: 'PlanSetting', component: () => import('@/views/setting/PlanSetting.vue'), meta: { title: '套餐管理' } },
        { path: 'personal', name: 'PersonalSetting', component: () => import('@/views/setting/PersonalSetting.vue'), meta: { title: '个人设置' } },
      ] },
      { path: 'redeem', name: 'Redeem', component: () => import('@/views/redeem/Redeem.vue'), meta: { title: '兑换' } },
      { path: 'chat', name: 'Chat', component: () => import('@/views/chat/Chat.vue'), meta: { title: '对话' } },
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
  const isAdminRoute = to.meta.admin || to.name === 'Channel' || to.name === 'User' || to.name === 'Redemption'

  // 根路径始终展示落地页
  if (to.path === '/') {
    return next()
  }

  // session 过期被踢回登录页：清除本地 user，停留登录页，避免与 /dashboard 死循环
  if (to.path === '/login' && to.query.expired) {
    localStorage.removeItem('user')
    return next()
  }

  if (!isLoggedIn && !['/login', '/register', '/reset'].includes(to.path) && !to.path.startsWith('/reset/') && !to.path.startsWith('/oauth/')) {
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

export default router
