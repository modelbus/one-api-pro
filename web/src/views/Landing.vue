<template>
  <div class="landing">
    <header class="nav-bar">
      <div class="nav-inner">
        <div class="nav-brand">
          <img v-if="statusStore.status?.logo" :src="statusStore.status.logo" class="nav-logo" />
          <span class="nav-name">{{ systemName }}</span>
        </div>
        <div class="nav-menu">
          <a href="#features">优势</a>
          <a href="#scenarios">场景</a>
          <a href="#compare">对比</a>
          <a href="#models">模型</a>
          <a href="#download">下载</a>
          <a href="#faq">问题</a>
          <a href="http://one-api.pro" target="_blank">文档</a>
        </div>
        <div class="nav-actions">
          <a-button type="text" size="small" href="https://github.com/Leon-PanPan/one-api-pro" target="_blank" class="nav-github-btn">
            <template #icon><icon-github /></template>
            GitHub
          </a-button>
          <template v-if="authStore.isLoggedIn">
            <a-button type="primary" @click="$router.push('/dashboard')">控制台</a-button>
          </template>
          <template v-else>
            <a-button type="outline" @click="$router.push('/login')">登录</a-button>
            <a-button type="primary" @click="$router.push('/register')">免费注册</a-button>
          </template>
        </div>
      </div>
    </header>

    <section class="hero">
      <div class="hero-bg">
        <div class="bg-circle c1"></div>
        <div class="bg-circle c2"></div>
        <div class="bg-circle c3"></div>
        <div class="bg-grid"></div>
        <div class="float-particle" v-for="n in 12" :key="n" :style="particleStyle(n)"></div>
      </div>
      <div class="hero-inner">
        <div class="hero-text">
          <div class="hero-badge">
            <a href="https://github.com/Leon-PanPan/one-api-pro" target="_blank" style="color:#165dff;text-decoration:none;">
              <icon-github style="vertical-align:-2px" /> GitHub
            </a>
          </div>
          <h1>企业级 <span class="text-gradient">AI API 网关</span></h1>
          <p>统一管理所有大模型 API，一键接入 30+ 平台<br/>智能路由、集群部署、开箱即用</p>
          <div class="hero-buttons">
            <a-button v-if="authStore.isLoggedIn" type="primary" size="large" shape="round" @click="$router.push('/dashboard')">
              进入控制台 <icon-right />
            </a-button>
            <template v-else>
              <a-button type="primary" size="large" shape="round" @click="$router.push('/register')">
                免费开始 <icon-right />
              </a-button>
              <a-button size="large" shape="round" @click="scrollTo('features')">查看优势</a-button>
            </template>
            <a-button size="large" shape="round" href="https://github.com/Leon-PanPan/one-api-pro" target="_blank" class="btn-gh-dark">
              <template #icon><icon-github /></template>
              GitHub
            </a-button>
            <a-button size="large" shape="round" href="https://github.com/Leon-PanPan/one-api-pro/releases" target="_blank">
              立即下载 <icon-down />
            </a-button>
          </div>
          <div class="hero-counts">
            <div class="count-item"><strong>30+</strong><span>模型平台</span></div>
            <div class="count-divider"></div>
            <div class="count-item"><strong>100%</strong><span>兼容 OpenAI</span></div>
            <div class="count-divider"></div>
            <div class="count-item"><strong>MIT</strong><span>开源协议</span></div>
          </div>
        </div>
        <div class="hero-visual">
          <div class="visual-card">
            <div class="visual-header">
              <span class="visual-dot red"></span><span class="visual-dot yellow"></span><span class="visual-dot green"></span>
            </div>
            <div class="visual-body">
              <div class="visual-line short"></div>
              <div class="visual-line medium"></div>
              <div class="visual-line long"></div>
              <div class="visual-tag">OpenAI</div>
              <div class="visual-tag v2">Anthropic</div>
              <div class="visual-tag v3">Gemini</div>
              <div class="visual-tag v4">DeepSeek</div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section id="features" class="section">
      <div class="section-header">
        <span class="section-label">平台优势</span>
        <h2>为什么选择 One Api Pro</h2>
        <p>不只是模型接入 — 更是企业 AI 基础设施的安全底座</p>
      </div>
      <div class="feature-grid">
        <div class="feature-item" v-for="f in features" :key="f.title">
          <div class="fi-icon" :style="{ color: f.color, background: f.color + '15' }">
            <component :is="f.icon" />
          </div>
          <h3>{{ f.title }}</h3>
          <p>{{ f.desc }}</p>
        </div>
      </div>
    </section>

    <section id="scenarios" class="section section-alt">
      <div class="section-header">
        <span class="section-label">应用场景</span>
        <h2>覆盖企业 AI 全流程</h2>
        <p>从研发测试到生产环境，一套系统满足所有场景</p>
      </div>
      <div class="scenario-grid">
        <div class="scenario-card">
          <icon-code class="sc-icon" />
          <h3>研发团队</h3>
          <p>统一接入多模型 API，开发者无需关心各平台差异，单 API Key 切换所有模型，大幅提升开发效率</p>
        </div>
        <div class="scenario-card">
          <icon-apps class="sc-icon" />
          <h3>SaaS 平台</h3>
          <p>为下游客户提供完整的 AI 能力分发，按用户/按模型计费，支持套餐订阅，快速实现商业化闭环</p>
        </div>
        <div class="scenario-card">
          <icon-cloud class="sc-icon" />
          <h3>企业内部</h3>
          <p>统一管控所有 AI 调用出口，配额分配、成本归集、安全审计，杜绝 API Key 泄露和费用超支</p>
        </div>
        <div class="scenario-card">
          <icon-thunderbolt class="sc-icon" />
          <h3>高并发业务</h3>
          <p>多渠道负载均衡 + 自动故障切换 + 去中心化集群，支撑百万级日调用，服务永不宕机</p>
        </div>
      </div>
    </section>

    <section id="compare" class="section">
      <div class="section-header">
        <span class="section-label">全面升级</span>
        <h2><span class="cmp-old">one-api</span> <span class="cmp-vs">vs</span> <span class="cmp-new">one-api-pro</span></h2>
        <p>在保留原版全部功能的基础上，进行了架构级重构</p>
      </div>
      <div class="compare-container">
        <div class="compare-header-row">
          <div class="compare-col-label"></div>
          <div class="compare-col old-col">
            <span class="col-badge old">原版 one-api</span>
          </div>
          <div class="compare-col new-col">
            <span class="col-badge new">one-api-pro</span>
          </div>
        </div>
        <div class="compare-row" v-for="item in compareItems" :key="item.label">
          <div class="compare-col-label">{{ item.label }}</div>
          <div class="compare-col old-col">
            <p>{{ item.old }}</p>
          </div>
          <div class="compare-col new-col">
            <p><icon-check-circle-fill style="color:#00b42a;vertical-align:-3px;margin-right:6px" />{{ item.new }}</p>
          </div>
        </div>
      </div>
    </section>

    <section id="models" class="section section-alt">
      <div class="section-header">
        <span class="section-label">生态兼容</span>
        <h2>支持 30+ 模型平台</h2>
        <p>覆盖全球主流 AI 服务商，持续增加中</p>
      </div>
      <div class="model-grid">
        <div class="model-card" v-for="m in models" :key="m.name">
          <div class="mc-icon">{{ m.icon }}</div>
          <span>{{ m.name }}</span>
        </div>
      </div>
    </section>

    <section id="download" class="section">
      <div class="section-header">
        <span class="section-label">快速开始</span>
        <h2>立即下载</h2>
        <p>开源免费，支持 Linux / macOS / Windows 多平台</p>
      </div>
      <div class="download-grid">
        <a href="https://github.com/Leon-PanPan/one-api-pro" target="_blank" class="download-card">
          <icon-github style="font-size:36px;color:#1d2129" />
          <h3>GitHub Releases</h3>
          <p>获取最新版本</p>
          <span class="dl-link">前往下载 →</span>
        </a>
        <a href="https://github.com/Leon-PanPan/one-api-pro" target="_blank" class="download-card">
          <icon-code style="font-size:36px;color:#1d2129" />
          <h3>源码编译</h3>
          <p>go build -o one-api-pro</p>
          <span class="dl-link">查看文档 →</span>
        </a>
        <a href="https://github.com/Leon-PanPan/one-api-pro#部署" target="_blank" class="download-card">
          <icon-cloud style="font-size:36px;color:#1d2129" />
          <h3>Docker 部署</h3>
          <p>docker-compose up -d</p>
          <span class="dl-link">部署指南 →</span>
        </a>
      </div>
    </section>

    <section id="faq" class="section section-alt">
      <div class="section-header">
        <span class="section-label">常见问题</span>
        <h2>你可能想了解</h2>
      </div>
      <div class="faq-list">
        <div class="faq-item" v-for="faq in faqs" :key="faq.q">
          <h3>{{ faq.q }}</h3>
          <p>{{ faq.a }}</p>
        </div>
      </div>
    </section>

    <footer class="footer">
      <div class="footer-inner">
        <div class="footer-col">
          <h4>{{ systemName }}</h4>
          <p>企业级 AI API 网关</p>
        </div>
        <div class="footer-col">
          <h4>资源</h4>
          <a href="https://github.com/Leon-PanPan/one-api-pro" target="_blank">GitHub</a>
          <a href="https://github.com/Leon-PanPan/one-api-pro/releases" target="_blank">Releases</a>
        </div>
        <div class="footer-col">
          <h4>文档</h4>
          <a href="#features">平台优势</a>
          <a href="#compare">版本对比</a>
          <a href="http://one-api.pro" target="_blank">在线文档</a>
        </div>
      </div>
      <div class="footer-bottom">
        &copy; {{ new Date().getFullYear() }} One Api Pro &nbsp;|&nbsp; <a href="https://github.com/Leon-PanPan/one-api-pro" target="_blank">GitHub</a> &nbsp;|&nbsp; MIT License
      </div>
    </footer>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useStatusStore } from '@/stores/status'
import { IconDashboard, IconApps, IconLock, IconCloud, IconThunderbolt, IconTool, IconCode, IconGithub, IconRight, IconDown, IconCheckCircleFill } from '@arco-design/web-vue/es/icon'

const router = useRouter()
const authStore = useAuthStore()
const statusStore = useStatusStore()
const systemName = computed(() => statusStore.status?.system_name || 'One Api Pro')

const features = [
  { icon: IconLock, title: '精细治理', desc: '多级权限管理、令牌按模型限制、分组倍率控制、IP 白名单，支持按渠道/用户/模型三级配额体系', color: '#165DFF' },
  { icon: IconTool, title: '精确成本核算', desc: '按 token 或按次计费，Prompt/Completion/Cached 独立定价，分组折扣叠加、周期用量追踪，分毫不差', color: '#00B42A' },
  { icon: IconApps, title: '企业级安全', desc: '全链路 HTTPS 传输、Token 鉴权、子网 IP 限制、审计日志实时追踪、防泄漏设计', color: '#FF7D00' },
  { icon: IconThunderbolt, title: '弹性伸缩与高可用', desc: '多渠道负载均衡、失败自动重试、冷却/禁用策略、去中心化多活集群，服务永不宕机', color: '#722ED1' },
  { icon: IconCloud, title: '统一 API 入口', desc: '30+ 模型平台统一 OpenAPI 格式接入，零适配切换，开发者只需一次对接，即可访问所有大模型', color: '#0FC6C2' },
  { icon: IconDashboard, title: '开箱即用运维', desc: '单文件部署、Docker 一键启动、可视化后台管理、系统公告推送，零学习成本快速上手', color: '#F53F3F' },
]

const compareItems = [
  { label: 'Adaptor 扩展', old: '需修改 4 个框架文件', new: '自注册机制，仅需新增包' },
  { label: '订阅模式', old: '无订阅/套餐体系', new: '完整套餐订阅 + 周期限频' },
  { label: '目录结构', old: 'adaptor/ 平铺 40 个目录', new: 'base/基础协议 + provider/供应商' },
  { label: 'Channel 类型', old: '56 行 iota 整数常量', new: '字符串 ID，语义清晰' },
  { label: '集群支持', old: '无独立集群', new: '去中心化多活集群' },
  { label: '管理后台', old: 'Semantic UI React', new: 'Vue 3 + Arco Design 全新重构' },
  { label: '持续更新', old: '2024 年停止更新', new: '持续维护，企业级优化' },
]

const faqs = [
  { q: 'One Api Pro 与原版 one-api 的关系是什么？', a: 'One Api Pro 基于 one-api（by JustSong）深度重构，在保留全部原有功能的基础上，对 Adaptor 架构、目录结构、管理后台进行了全面升级，新增了订阅模式、去中心化集群等企业级特性。' },
  { q: '如何从原版 one-api 迁移到 One Api Pro？', a: '数据库结构完全兼容，直接替换二进制即可运行。如有自定义 Adaptor，需按新注册机制（register.go）进行适配，改动量极小。' },
  { q: '支持哪些数据库？', a: '默认使用 SQLite（零配置），也支持 MySQL 5.7+ 和 PostgreSQL，生产环境推荐 MySQL。Redis 为可选项，启用后可获得更好的缓存性能。' },
  { q: '是否支持私有化部署？', a: '完全支持。One Api Pro 本身就是为私有化部署设计的，所有数据存储在你的服务器上，不外泄任何信息。' },
  { q: '多个 API Key 如何负载均衡？', a: '支持按权重随机分配、自动故障切换、渠道冷却/禁用策略，还可按用户组和模型维度精细化配置优先级。' },
]

const models = [
  { icon: '🤖', name: 'OpenAI' }, { icon: '☁️', name: 'Azure' }, { icon: '🧠', name: 'Anthropic' },
  { icon: '☁️', name: 'AWS Bedrock' }, { icon: '🌐', name: 'Google Gemini' }, { icon: '🔍', name: 'DeepSeek' },
  { icon: '📚', name: '百度文心' }, { icon: '💡', name: '阿里通义' }, { icon: '⭐', name: '讯飞星火' },
  { icon: '🎓', name: '智谱 ChatGLM' }, { icon: '🛡️', name: '360 智脑' }, { icon: '💬', name: '腾讯混元' },
  { icon: '🌙', name: 'Moonshot' }, { icon: '🏔️', name: '百川' }, { icon: '🎯', name: 'MINIMAX' },
  { icon: '⚡', name: 'Groq' }, { icon: '🐫', name: 'Ollama' }, { icon: '🌊', name: '零一万物' },
  { icon: '✨', name: '阶跃星辰' }, { icon: '🤖', name: 'Coze' }, { icon: '🧩', name: 'Cohere' },
  { icon: '☁️', name: 'Cloudflare' }, { icon: '🌍', name: 'DeepL' }, { icon: '🤝', name: 'together.ai' },
  { icon: '💎', name: '硅基流动' }, { icon: '❌', name: 'xAI' }, { icon: '🔗', name: 'OpenRouter' },
  { icon: '🫘', name: '豆包' }, { icon: '🎨', name: 'Novita' },
]

const particleStyle = (n) => ({
  left: `${(n * 37 + 13) % 100}%`,
  top: `${(n * 23 + 7) % 100}%`,
  animationDelay: `${(n * 0.7) % 5}s`,
  animationDuration: `${4 + (n % 4)}s`,
  width: `${4 + (n % 6)}px`,
  height: `${4 + (n % 6)}px`,
})

const scrollTo = (id) => document.getElementById(id)?.scrollIntoView({ behavior: 'smooth' })

onMounted(async () => {
  if (!statusStore.loaded) await statusStore.fetchStatus()
  authStore.loadUser()
})
</script>

<style scoped>
.landing { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; color: #1d2129; overflow-x: hidden; }

/* Nav */
.nav-bar { position: fixed; top: 0; left: 0; right: 0; z-index: 1000; background: rgba(255,255,255,0.92); backdrop-filter: blur(12px); border-bottom: 1px solid rgba(0,0,0,0.06); }
.nav-inner { max-width: 1280px; margin: 0 auto; padding: 0 32px; height: 64px; display: flex; align-items: center; gap: 40px; }
.nav-brand { display: flex; align-items: center; gap: 10px; flex-shrink: 0; }
.nav-logo { width: 32px; height: 32px; border-radius: 6px; }
.nav-name { font-size: 18px; font-weight: 700; color: #1d2129; }
.nav-menu { display: flex; gap: 32px; flex: 1; }
.nav-menu a { color: #4e5969; text-decoration: none; font-size: 14px; font-weight: 500; transition: color .2s; }
.nav-menu a:hover { color: #165dff; }
.nav-actions { display: flex; gap: 12px; flex-shrink: 0; align-items: center; }
.nav-github-btn { color: #4e5969 !important; }
.nav-github-btn:hover { color: #1d2129 !important; background: #f2f3f5 !important; }

/* Hero */
.hero { position: relative; padding: 160px 32px 120px; overflow: hidden; min-height: 100vh; display: flex; align-items: center; background: linear-gradient(180deg, #f0f5ff 0%, #fff 100%); }
.hero-bg { position: absolute; inset: 0; overflow: hidden; }
.bg-circle { position: absolute; border-radius: 50%; opacity: 0.08; }
.c1 { width: 800px; height: 800px; background: #165dff; top: -300px; right: -200px; animation: float1 20s ease-in-out infinite; }
.c2 { width: 500px; height: 500px; background: #722ed1; bottom: -100px; left: -100px; animation: float2 25s ease-in-out infinite; }
.c3 { width: 300px; height: 300px; background: #00b42a; top: 200px; left: 40%; animation: float3 18s ease-in-out infinite; }
.bg-grid { position: absolute; inset: 0; background-image: radial-gradient(#e5e6eb 1px, transparent 1px); background-size: 32px 32px; opacity: 0.4; }
.float-particle { position: absolute; border-radius: 50%; background: linear-gradient(135deg, #165dff, #722ed1); opacity: 0.15; animation: floatUp linear infinite; }
@keyframes float1 { 0%,100% { transform: translate(0,0) scale(1); } 50% { transform: translate(40px,-30px) scale(1.05); } }
@keyframes float2 { 0%,100% { transform: translate(0,0); } 50% { transform: translate(-30px,20px) scale(1.08); } }
@keyframes float3 { 0%,100% { transform: translate(0,0) scale(1); } 50% { transform: translate(-20px,-25px) scale(1.04); } }
@keyframes floatUp { 0% { transform: translateY(0) scale(0); opacity: 0; } 10% { opacity: 0.2; } 90% { opacity: 0.1; } 100% { transform: translateY(-600px) scale(1.5); opacity: 0; } }

.hero-inner { position: relative; z-index: 1; max-width: 1280px; margin: 0 auto; display: flex; align-items: center; gap: 80px; }
.hero-text { flex: 1; }
.hero-badge { display: inline-block; padding: 6px 16px; background: rgba(22,93,255,0.08); color: #165dff; border-radius: 20px; font-size: 13px; font-weight: 500; margin-bottom: 24px; }
.hero-text h1 { font-size: 52px; font-weight: 800; line-height: 1.2; margin-bottom: 20px; color: #1d2129; }
.text-gradient { background: linear-gradient(135deg, #165dff, #722ed1); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; }
.hero-text > p { font-size: 17px; color: #4e5969; line-height: 1.8; margin-bottom: 36px; }
.hero-buttons { display: flex; gap: 16px; flex-wrap: wrap; }
.btn-gh-dark { background: #1b1f23 !important; color: #fff !important; border-color: #1b1f23 !important; }
.btn-gh-dark:hover { background: #2d333b !important; border-color: #2d333b !important; }
.hero-counts { display: flex; align-items: center; gap: 24px; margin-top: 60px; }
.count-item { display: flex; flex-direction: column; gap: 2px; }
.count-item strong { font-size: 24px; font-weight: 700; color: #1d2129; }
.count-item span { font-size: 13px; color: #86909c; }
.count-divider { width: 1px; height: 32px; background: #e5e6eb; }

/* Hero Visual */
.hero-visual { flex: 0 0 420px; }
.visual-card { background: #fff; border-radius: 16px; box-shadow: 0 20px 60px rgba(0,0,0,0.08); border: 1px solid #e5e6eb; overflow: hidden; animation: visualFloat 6s ease-in-out infinite; }
@keyframes visualFloat { 0%,100% { transform: translateY(0); } 50% { transform: translateY(-10px); } }
.visual-header { padding: 14px 20px; background: #f7f8fa; display: flex; gap: 8px; border-bottom: 1px solid #e5e6eb; }
.visual-dot { width: 10px; height: 10px; border-radius: 50%; }
.visual-dot.red { background: #f53f3f; }
.visual-dot.yellow { background: #ff7d00; }
.visual-dot.green { background: #00b42a; }
.visual-body { padding: 32px 24px; position: relative; display: flex; flex-direction: column; gap: 12px; min-height: 200px; }
.visual-line { height: 8px; border-radius: 4px; background: #f2f3f5; }
.visual-line.short { width: 40%; }
.visual-line.medium { width: 65%; }
.visual-line.long { width: 80%; }
.visual-tag { position: absolute; padding: 8px 16px; background: rgba(22,93,255,0.08); color: #165dff; border-radius: 8px; font-size: 13px; font-weight: 500; }
.visual-tag { right: 24px; top: 40px; }
.visual-tag.v2 { top: 80px; background: rgba(114,46,209,0.08); color: #722ed1; right: 80px; }
.visual-tag.v3 { top: 120px; background: rgba(0,180,42,0.08); color: #00b42a; right: 120px; }
.visual-tag.v4 { top: 160px; background: rgba(255,125,0,0.08); color: #ff7d00; right: 40px; }

/* Sections */
.section { padding: 100px 32px; max-width: 1280px; margin: 0 auto; }
.section-alt { background: #f7f8fa; max-width: 100%; }
.section-alt .section-header, .section-alt .feature-grid, .section-alt .scenario-grid, .section-alt .model-grid, .section-alt .download-grid, .section-alt .faq-list { max-width: 1280px; margin-left: auto; margin-right: auto; }
.section-header { text-align: center; margin-bottom: 60px; }
.section-label { display: inline-block; padding: 4px 14px; background: rgba(22,93,255,0.06); color: #165dff; border-radius: 20px; font-size: 13px; font-weight: 600; margin-bottom: 16px; }
.section-header h2 { font-size: 36px; font-weight: 700; margin-bottom: 12px; color: #1d2129; }
.section-header p { font-size: 16px; color: #86909c; }

/* Features */
.feature-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 20px; }
.feature-item { padding: 36px 28px; border-radius: 16px; background: #fff; border: 1px solid #f2f3f5; transition: all .3s; }
.feature-item:hover { border-color: #165dff20; box-shadow: 0 8px 30px rgba(0,0,0,0.06); transform: translateY(-2px); }
.fi-icon { width: 44px; height: 44px; border-radius: 12px; display: flex; align-items: center; justify-content: center; font-size: 22px; margin-bottom: 20px; }
.feature-item h3 { font-size: 17px; font-weight: 600; margin-bottom: 8px; }
.feature-item p { font-size: 14px; color: #86909c; line-height: 1.7; }

/* Scenarios */
.scenario-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; }
.scenario-card { padding: 36px 24px; border-radius: 16px; background: #fff; border: 1px solid #e5e6eb; text-align: center; transition: all .3s; }
.scenario-card:hover { border-color: #165dff40; box-shadow: 0 4px 20px rgba(0,0,0,0.06); transform: translateY(-2px); }
.sc-icon { font-size: 32px; color: #165dff; margin-bottom: 16px; }
.scenario-card h3 { font-size: 17px; font-weight: 600; margin-bottom: 8px; }
.scenario-card p { font-size: 13px; color: #86909c; line-height: 1.7; }

/* Models */
.model-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: 12px; }
.model-card { padding: 16px; border-radius: 12px; background: #fff; border: 1px solid #e5e6eb; display: flex; align-items: center; gap: 10px; transition: all .2s; }
.model-card:hover { border-color: #165dff40; box-shadow: 0 2px 12px rgba(22,93,255,0.08); }
.mc-icon { font-size: 20px; }

/* Compare */
.cmp-old { color: #86909c; font-weight: 600; }
.cmp-vs { color: #c9cdd4; font-size: 0.7em; font-weight: 400; margin: 0 4px; vertical-align: middle; }
.cmp-new { background: linear-gradient(135deg, #165dff, #722ed1); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; font-weight: 700; }
.compare-container { max-width: 900px; margin: 0 auto; background: #fff; border-radius: 20px; overflow: hidden; box-shadow: 0 2px 20px rgba(0,0,0,0.04); border: 1px solid #e5e6eb; }
.compare-header-row { display: flex; background: #f7f8fa; border-bottom: 1px solid #e5e6eb; }
.compare-col-label { flex: 0 0 160px; padding: 18px 24px; font-weight: 600; font-size: 14px; color: #1d2129; }
.compare-col { flex: 1; padding: 18px 24px; text-align: center; }
.old-col { background: #fafafa; }
.new-col { background: linear-gradient(135deg, #f0f5ff, #f9f0ff); }
.col-badge { display: inline-block; padding: 4px 16px; border-radius: 12px; font-size: 13px; font-weight: 600; }
.col-badge.old { background: #f2f3f5; color: #86909c; }
.col-badge.new { background: linear-gradient(135deg, #165dff15, #722ed115); color: #165dff; }
.compare-row { display: flex; border-bottom: 1px solid #f2f3f5; transition: background .2s; }
.compare-row:last-child { border-bottom: none; }
.compare-row:hover { background: #f9fafb; }
.compare-row .compare-col-label { font-weight: 500; color: #4e5969; }
.compare-row .compare-col { padding: 22px 24px; display: flex; align-items: center; justify-content: center; }
.compare-row p { margin: 0; font-size: 14px; color: #4e5969; }
.compare-row .new-col p { color: #1d2129; font-weight: 500; }

/* Download */
.download-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 20px; }
.download-card { padding: 40px 28px; border-radius: 16px; background: #fff; border: 1px solid #e5e6eb; text-align: center; text-decoration: none; color: #1d2129; transition: all .3s; display: block; }
.download-card:hover { border-color: #165dff40; box-shadow: 0 8px 30px rgba(0,0,0,0.08); transform: translateY(-2px); }
.download-card h3 { font-size: 17px; font-weight: 600; margin: 16px 0 8px; }
.download-card p { font-size: 14px; color: #86909c; font-family: monospace; }
.dl-link { display: inline-block; margin-top: 12px; font-size: 14px; color: #165dff; font-weight: 500; }

/* FAQ */
.faq-list { max-width: 800px; margin: 0 auto; }
.faq-item { padding: 28px 0; border-bottom: 1px solid #e5e6eb; }
.faq-item:first-child { padding-top: 0; }
.faq-item h3 { font-size: 17px; font-weight: 600; margin-bottom: 8px; color: #1d2129; }
.faq-item p { font-size: 14px; color: #86909c; line-height: 1.7; }

/* Footer */
.footer { background: #1d2129; color: rgba(255,255,255,0.6); padding: 60px 32px 28px; }
.footer-inner { max-width: 1280px; margin: 0 auto; display: flex; gap: 80px; margin-bottom: 40px; }
.footer-col h4 { color: #fff; font-size: 15px; font-weight: 600; margin-bottom: 16px; }
.footer-col p { font-size: 13px; line-height: 1.6; }
.footer-col a { display: block; color: rgba(255,255,255,0.5); text-decoration: none; font-size: 13px; margin-bottom: 10px; transition: color .2s; }
.footer-col a:hover { color: #fff; }
.footer-bottom { max-width: 1280px; margin: 0 auto; text-align: center; font-size: 12px; border-top: 1px solid rgba(255,255,255,0.08); padding-top: 24px; }
.footer-bottom a { color: rgba(255,255,255,0.4); }

@media (max-width: 768px) {
  .nav-menu, .nav-actions .a-btn:first-child { display: none; }
  .hero-inner { flex-direction: column; gap: 40px; }
  .hero-text h1 { font-size: 34px; }
  .hero-visual { flex: none; width: 100%; max-width: 400px; }
  .feature-grid, .scenario-grid, .download-grid { grid-template-columns: 1fr; }
  .compare-col-label { flex: 0 0 100px; font-size: 12px; padding: 14px 12px; }
  .compare-col { padding: 14px 12px; }
  .footer-inner { flex-direction: column; gap: 32px; }
  .hero-counts { flex-wrap: wrap; gap: 16px; }
}
</style>
