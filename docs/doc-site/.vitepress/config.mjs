import { defineConfig } from 'vitepress'

// ---------------------------------------------------------------------------
// Sidebar definitions
// ---------------------------------------------------------------------------
// Category meta: slug -> [zh label, en label] (12 categories per spec)
const CATEGORY_LABELS = {
  start: ['快速开始', 'Getting Started'],
  install: ['安装与更新', 'Install & Upgrade'],
  user: ['用户', 'User'],
  channel: ['渠道', 'Channel'],
  subscription: ['订阅（Token Plan）', 'Subscription'],
  redemption: ['兑换码', 'Redemption'],
  pricing: ['定价', 'Pricing'],
  decentralization: ['去中心化', 'Decentralization'],
  schema: ['数据结构', 'Schema'],
  api: ['API', 'API'],
  changelog: ['更新日志', 'Changelog'],
  contribute: ['贡献', 'Contribute'],
  misc: ['其他', 'Misc'],
}

// Per-category page metadata: slug -> [zh title, en title]
const PAGES = {
  start: {
    overview: ['项目简介', 'Overview'],
    showcase: ['展示案例', 'Showcase'],
    'quick-tour': ['5 分钟快速体验', 'Quick Tour'],
    architecture: ['架构总览', 'Architecture'],
  },
  install: {
    requirements: ['系统要求', 'Requirements'],
    'docker-deploy': ['Docker 单实例部署', 'Docker Deploy'],
    'docker-compose': ['docker-compose 部署', 'Docker Compose'],
    'source-build': ['源码编译', 'Source Build'],
    config: ['配置项与环境变量', 'Configuration'],
    upgrade: ['版本升级', 'Upgrade'],
    'backup-restore': ['备份与恢复', 'Backup & Restore'],
    'reverse-proxy': ['反向代理', 'Reverse Proxy'],
  },
  user: {
    registration: ['注册与登录', 'Registration & Login'],
    dashboard: ['个人仪表盘', 'Dashboard'],
    profile: ['个人资料', 'Profile'],
    'access-token': ['Access Token', 'Access Token'],
    orders: ['我的订单', 'My Orders'],
    chat: ['Chat Playground', 'Chat Playground'],
    'user-management': ['用户管理（管理员）', 'User Management (Admin)'],
  },
  channel: {
    overview: ['渠道概览', 'Channel Overview'],
    'add-channel': ['新建渠道', 'Add a Channel'],
    'channel-routing': ['渠道路由策略', 'Channel Routing'],
    'channel-test': ['渠道连通性测试', 'Channel Test'],
    'balance-update': ['余额自动更新', 'Balance Update'],
    'provider-list': ['Provider 全清单', 'Provider List'],
  },
  subscription: {
    overview: ['订阅概览', 'Subscription Overview'],
    'user-guide': ['用户使用指南', 'User Guide'],
    'admin-guide': ['管理员指南', 'Admin Guide'],
    'billing-rules': ['计费规则', 'Billing Rules'],
    'upgrade-downgrade': ['套餐升降级', 'Upgrade & Downgrade'],
    'plan-management': ['套餐管理', 'Plan Management'],
    'subscription-management': ['订阅管理', 'Subscription Management'],
    'plan-settings': ['套餐业务配置', 'Plan Settings'],
  },
  redemption: {
    overview: ['兑换码概览', 'Redemption Overview'],
    'admin-guide': ['管理员指南', 'Admin Guide'],
    'user-guide': ['用户兑换流程', 'User Guide'],
    quotas: ['兑换额度计算', 'Quota Calculation'],
    'redemption-management': ['兑换码管理', 'Redemption Management'],
  },
  pricing: {
    'model-price': ['模型定价', 'Model Price'],
    'group-price': ['分组折扣', 'Group Price'],
    plan: ['套餐定价', 'Plan Pricing'],
    topup: ['充值', 'Topup'],
    payment: ['支付通道', 'Payment Channels'],
    'model-price-management': ['模型定价管理', 'Model Price Management'],
    'group-price-management': ['分组折扣管理', 'Group Price Management'],
    'topup-management': ['充值管理', 'Topup Management'],
    'order-management': ['订单管理', 'Order Management'],
    'payment-settings': ['支付配置', 'Payment Settings'],
    'topup-settings': ['充值业务配置', 'Topup Settings'],
  },
  decentralization: {
    overview: ['Cluster 概览', 'Cluster Overview'],
    'node-management': ['节点管理', 'Node Management'],
    'config-sync': ['配置同步', 'Config Sync'],
    'node-health': ['节点健康', 'Node Health'],
    deployment: ['多节点部署', 'Multi-node Deployment'],
    'cluster-settings': ['集群设置', 'Cluster Settings'],
  },
  schema: {
    overview: ['数据结构总览', 'Schema Overview'],
    user: ['用户', 'User'],
    token: ['访问令牌', 'Access Token'],
    channel: ['渠道', 'Channel'],
    'model-price': ['模型定价', 'Model Price'],
    'group-price': ['分组折扣', 'Group Price'],
    plan: ['套餐', 'Plan'],
    subscription: ['订阅', 'Subscription'],
    order: ['订单', 'Order'],
    topup: ['充值', 'Topup'],
    redemption: ['兑换码', 'Redemption'],
    'cluster-node': ['集群节点', 'Cluster Node'],
    log: ['调用日志', 'Log'],
    option: ['系统设置', 'Option'],
  },
  api: {
    README: ['总览与鉴权', 'Overview & Auth'],
    channel: ['渠道 API', 'Channel API'],
    token: ['令牌 API', 'Token API'],
    user: ['用户 API', 'User API'],
    redemption: ['兑换码 API', 'Redemption API'],
    subscription: ['订阅 API', 'Subscription API'],
    plan: ['套餐 API', 'Plan API'],
    order: ['订单 API', 'Order API'],
    topup: ['充值 API', 'Topup API'],
    payment: ['支付 API', 'Payment API'],
    'model-price': ['模型定价 API', 'Model Price API'],
    'group-price': ['分组折扣 API', 'Group Price API'],
    log: ['日志 API', 'Log API'],
    cluster: ['集群 API', 'Cluster API'],
    'admin-dashboard': ['运营仪表盘 API', 'Admin Dashboard API'],
    diag: ['诊断 API', 'Diagnostic API'],
    misc: ['其他公共 API', 'Misc Public API'],
    'openai-compatible': ['OpenAI 兼容接口', 'OpenAI Compatible'],
    'anthropic-compatible': ['Anthropic 兼容接口', 'Anthropic Compatible'],
    'error-code': ['错误码表', 'Error Codes'],
    'token-management': ['API Token 管理', 'API Token Management'],
  },
  changelog: {
    index: ['版本索引', 'All Releases'],
  },
  contribute: {
    'dev-setup': ['开发环境搭建', 'Dev Setup'],
    'code-style': ['编码规范', 'Code Style'],
    'commit-convention': ['提交规范', 'Commit Convention'],
    'add-provider': ['新增 Provider', 'Add a Provider'],
    'add-payment': ['新增支付通道', 'Add a Payment Channel'],
    'release-process': ['发版流程', 'Release Process'],
  },
  misc: {
    'admin-dashboard': ['运营仪表盘', 'Admin Dashboard'],
    logs: ['日志查询', 'Logs'],
    'operation-settings': ['运营设置', 'Operation Settings'],
    'system-settings': ['系统设置', 'System Settings'],
    faq: ['常见问题', 'FAQ'],
    troubleshooting: ['故障排查', 'Troubleshooting'],
    glossary: ['术语表', 'Glossary'],
  },
}

// Changelog is generated dynamically from the files present.
const CHANGELOG_VERSIONS = [
  'v0.0.21', 'v0.0.20', 'v0.0.18', 'v0.0.17', 'v0.0.16', 'v0.0.15',
  'v0.0.13', 'v0.0.12', 'v0.0.10', 'v0.0.9', 'v0.0.8', 'v0.0.7',
  'v0.0.6', 'v0.0.5', 'v0.0.4', 'v0.0.3', 'v0.0.2', 'v0.0.1',
]

// Build a sidebar for a locale: zh lives at root, en under /en/.
function buildSidebar(lang) {
  const prefix = lang === 'zh' ? '' : '/en'
  const li = lang === 'zh' ? 0 : 1
  const sidebar = {}

  for (const [cat, labelPair] of Object.entries(CATEGORY_LABELS)) {
    const catLabel = labelPair[li]
    const items = Object.entries(PAGES[cat] || {}).map(([slug, titles]) => ({
      text: titles[li],
      link: `${prefix}/${cat}/${slug}`,
    }))

    if (cat === 'changelog') {
      items.push({
        text: lang === 'zh' ? '历史版本' : 'All versions',
        collapsed: true,
        items: CHANGELOG_VERSIONS.map((v) => ({
          text: v,
          link: `${prefix}/changelog/${v}`,
        })),
      })
    }

    sidebar[`${prefix}/${cat}/`] = [{ text: catLabel, items }]
  }
  return sidebar
}

function buildNav(lang) {
  const prefix = lang === 'zh' ? '' : '/en'
  const items = lang === 'zh'
    ? [
        { text: '快速开始', link: '/start/overview' },
        { text: '安装', link: '/install/docker-deploy' },
        { text: '用户', link: '/user/registration' },
        { text: '渠道', link: '/channel/overview' },
        { text: '订阅', link: '/subscription/overview' },
        { text: '兑换码', link: '/redemption/overview' },
        { text: '定价', link: '/pricing/model-price' },
        { text: '去中心化', link: '/decentralization/overview' },
        { text: '数据结构', link: '/schema/overview' },
        { text: 'API', link: '/api/README' },
        { text: '更新日志', link: '/changelog/index' },
        { text: '贡献', link: '/contribute/dev-setup' },
        { text: '其他', link: '/misc/faq' },
      ]
    : [
        { text: 'Start', link: '/en/start/overview' },
        { text: 'Install', link: '/en/install/docker-deploy' },
        { text: 'User', link: '/en/user/registration' },
        { text: 'Channel', link: '/en/channel/overview' },
        { text: 'Subscription', link: '/en/subscription/overview' },
        { text: 'Redemption', link: '/en/redemption/overview' },
        { text: 'Pricing', link: '/en/pricing/model-price' },
        { text: 'Cluster', link: '/en/decentralization/overview' },
        { text: 'Schema', link: '/en/schema/overview' },
        { text: 'API', link: '/en/api/README' },
        { text: 'Changelog', link: '/en/changelog/index' },
        { text: 'Contribute', link: '/en/contribute/dev-setup' },
        { text: 'Misc', link: '/en/misc/faq' },
      ]
  return items
}

export default defineConfig({
  title: 'One API Pro Docs',
  description: 'LLM 聚合网关与商业化平台 · 文档',
  srcDir: 'docs',
  cleanUrls: true,
  // Synced files are untracked by git, so lastUpdated would be empty.
  lastUpdated: false,
  head: [
    ['link', { rel: 'icon', type: 'image/png', href: '/logo.png' }],
    ['meta', { name: 'theme-color', content: '#6d5efc' }],
    [
      'meta',
      { property: 'og:title', content: 'One API Pro Docs' },
    ],
    [
      'meta',
      {
        property: 'og:description',
        content: '统一接入 40+ 大模型，内置渠道路由、订阅计费、充值支付与多节点集群。',
      },
    ],
  ],
  markdown: {
    // API docs are sanitised by docs/scripts/sync-docs.mjs (angle brackets escaped)
    // so raw HTML can stay enabled for card grids on the home page.
    html: true,
    theme: { light: 'github-light', dark: 'github-dark' },
  },
  ignoreDeadLinks: true,
  locales: {
    root: {
      label: '简体中文',
      lang: 'zh-CN',
      link: '/',
      themeConfig: {
        logo: { src: '/logo.png', alt: 'One API Pro' },
        nav: buildNav('zh'),
        sidebar: buildSidebar('zh'),
        outline: { level: [2, 3], label: '本页目录' },
        docFooter: { prev: '上一篇', next: '下一篇' },
        darkModeSwitchLabel: '外观',
        sidebarMenuLabel: '目录',
        returnToTopLabel: '回到顶部',
        lastUpdatedText: '最后更新',
        langMenuLabel: '切换语言',
        footer: {
          message: '基于 MIT 许可发布 · 文档持续更新中',
          copyright: 'Copyright © 2024-present One API Pro',
        },
      },
    },
    en: {
      label: 'English',
      lang: 'en-US',
      link: '/en/',
      themeConfig: {
        logo: { src: '/logo.png', alt: 'One API Pro' },
        nav: buildNav('en'),
        sidebar: buildSidebar('en'),
        outline: { level: [2, 3], label: 'On this page' },
        docFooter: { prev: 'Previous', next: 'Next' },
        darkModeSwitchLabel: 'Appearance',
        sidebarMenuLabel: 'Menu',
        returnToTopLabel: 'Return to top',
        lastUpdatedText: 'Last updated',
        langMenuLabel: 'Change language',
        footer: {
          message: 'Released under the MIT License · Docs are a work in progress',
          copyright: 'Copyright © 2024-present One API Pro',
        },
      },
    },
  },
})
