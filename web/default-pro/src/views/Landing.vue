<template>
  <div class="landing">
    <header class="nav-bar">
      <div class="nav-inner">
        <div class="nav-brand">
          <img :src="logoSrc" class="nav-logo" @error="onLogoError" />
        </div>
        <div class="nav-menu">
          <a href="#features">{{ $t('landing.navFeatures') }}</a>
          <a href="#scenarios">{{ $t('landing.navScenarios') }}</a>
          <a href="#compare">{{ $t('landing.navCompare') }}</a>
          <a href="#models">{{ $t('landing.navModels') }}</a>
          <a href="#download">{{ $t('landing.navDownload') }}</a>
          <a href="#faq">{{ $t('landing.navFaq') }}</a>
          <a href="http://one-api.pro" target="_blank">{{ $t('landing.navDocs') }}</a>
        </div>
        <div class="nav-actions">
          <LangSwitcher size="small" width="110px" />
          <a-button type="text" size="small" href="https://github.com/modelbus/one-api-pro" target="_blank" class="nav-github-btn">
            <template #icon><icon-github /></template>
            GitHub
          </a-button>
          <template v-if="authStore.isLoggedIn">
            <a-button type="primary" @click="$router.push('/dashboard')">{{ $t('landing.console') }}</a-button>
          </template>
          <template v-else>
            <a-button type="outline" @click="$router.push('/login')">{{ $t('landing.login') }}</a-button>
            <a-button type="primary" @click="$router.push('/register')">{{ $t('landing.register') }}</a-button>
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
            <a href="https://github.com/modelbus/one-api-pro" target="_blank" style="color:#165dff;text-decoration:none;">
              <icon-github style="vertical-align:-2px" /> GitHub
            </a>
          </div>
          <h1>{{ $t('landing.heroTitlePrefix') }} <span class="text-gradient">{{ $t('landing.heroTitleGradient') }}</span></h1>
          <p>{{ $t('landing.heroLine1') }}<br/>{{ $t('landing.heroLine2') }}</p>
          <div class="hero-buttons">
            <a-button v-if="authStore.isLoggedIn" type="primary" size="large" shape="round" @click="$router.push('/dashboard')">
              {{ $t('landing.enterConsole') }} <icon-right />
            </a-button>
            <template v-else>
              <a-button type="primary" size="large" shape="round" @click="$router.push('/register')">
                {{ $t('landing.getStartedFree') }} <icon-right />
              </a-button>
              <a-button size="large" shape="round" @click="scrollTo('features')">{{ $t('landing.viewFeatures') }}</a-button>
            </template>
            <a-button size="large" shape="round" href="https://github.com/modelbus/one-api-pro" target="_blank" class="btn-gh-dark">
              <template #icon><icon-github /></template>
              GitHub
            </a-button>
            <a-button size="large" shape="round" href="https://github.com/modelbus/one-api-pro/releases" target="_blank">
              {{ $t('landing.downloadNow') }} <icon-down />
            </a-button>
          </div>
          <div class="hero-counts">
            <div class="count-item"><strong>30+</strong><span>{{ $t('landing.countPlatforms') }}</span></div>
            <div class="count-divider"></div>
            <div class="count-item"><strong>100%</strong><span>{{ $t('landing.countOpenAI') }}</span></div>
            <div class="count-divider"></div>
            <div class="count-item"><strong>MIT</strong><span>{{ $t('landing.countLicense') }}</span></div>
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
        <span class="section-label">{{ $t('landing.featuresLabel') }}</span>
        <h2>{{ $t('landing.featuresTitle') }}</h2>
        <p>{{ $t('landing.featuresSubtitle') }}</p>
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
        <span class="section-label">{{ $t('landing.scenariosLabel') }}</span>
        <h2>{{ $t('landing.scenariosTitle') }}</h2>
        <p>{{ $t('landing.scenariosSubtitle') }}</p>
      </div>
      <div class="scenario-grid">
        <div class="scenario-card">
          <icon-code class="sc-icon" />
          <h3>{{ $t('landing.scenarioRnd') }}</h3>
          <p>{{ $t('landing.scenarioRndDesc') }}</p>
        </div>
        <div class="scenario-card">
          <icon-apps class="sc-icon" />
          <h3>{{ $t('landing.scenarioSaas') }}</h3>
          <p>{{ $t('landing.scenarioSaasDesc') }}</p>
        </div>
        <div class="scenario-card">
          <icon-cloud class="sc-icon" />
          <h3>{{ $t('landing.scenarioEnterprise') }}</h3>
          <p>{{ $t('landing.scenarioEnterpriseDesc') }}</p>
        </div>
        <div class="scenario-card">
          <icon-thunderbolt class="sc-icon" />
          <h3>{{ $t('landing.scenarioHighConcurrency') }}</h3>
          <p>{{ $t('landing.scenarioHighConcurrencyDesc') }}</p>
        </div>
      </div>
    </section>

    <section id="compare" class="section">
      <div class="section-header">
        <span class="section-label">{{ $t('landing.compareLabel') }}</span>
        <h2><span class="cmp-old">one-api</span> <span class="cmp-vs">vs</span> <span class="cmp-new">one-api-pro</span></h2>
        <p>{{ $t('landing.compareSubtitle') }}</p>
      </div>
      <div class="compare-container">
        <div class="compare-header-row">
          <div class="compare-col-label"></div>
          <div class="compare-col old-col">
            <span class="col-badge old">{{ $t('landing.compareOldBadge') }}</span>
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
        <span class="section-label">{{ $t('landing.modelsLabel') }}</span>
        <h2>{{ $t('landing.modelsTitle') }}</h2>
        <p>{{ $t('landing.modelsSubtitle') }}</p>
      </div>
      <div class="model-grid">
        <div class="model-card" v-for="m in models" :key="m.slug">
          <img v-if="m.slug" :src="modelIconSrc(m.slug)" class="mc-icon" :alt="m.name" loading="lazy" @error="onIconError" />
          <span v-else class="mc-icon mc-icon-fallback" :style="{ background: m.color + '15', color: m.color }">{{ m.name.charAt(0) }}</span>
          <span class="mc-name">{{ m.name }}</span>
          <span v-if="m.tag === '国产'" class="mc-tag mc-tag-cn">{{ $t('landing.modelDomestic') }}</span>
          <span v-else class="mc-tag mc-tag-intl">{{ $t('landing.modelOverseas') }}</span>
        </div>
      </div>
    </section>

    <section id="download" class="section">
      <div class="section-header">
        <span class="section-label">{{ $t('landing.downloadLabel') }}</span>
        <h2>{{ $t('landing.downloadTitle') }}</h2>
        <p>{{ $t('landing.downloadSubtitle') }}</p>
      </div>
      <div class="download-grid">
        <a href="https://github.com/modelbus/one-api-pro" target="_blank" class="download-card">
          <icon-github style="font-size:36px;color:#1d2129" />
          <h3>GitHub Releases</h3>
          <p>{{ $t('landing.downloadGithubDesc') }}</p>
          <span class="dl-link">{{ $t('landing.downloadGo') }}</span>
        </a>
        <a href="https://github.com/modelbus/one-api-pro" target="_blank" class="download-card">
          <icon-code style="font-size:36px;color:#1d2129" />
          <h3>{{ $t('landing.downloadSourceTitle') }}</h3>
          <p>{{ $t('landing.downloadSourceDesc') }}</p>
          <span class="dl-link">{{ $t('landing.downloadDocs') }}</span>
        </a>
        <a href="https://github.com/modelbus/one-api-pro#部署" target="_blank" class="download-card">
          <icon-cloud style="font-size:36px;color:#1d2129" />
          <h3>{{ $t('landing.downloadDockerTitle') }}</h3>
          <p>{{ $t('landing.downloadDockerDesc') }}</p>
          <span class="dl-link">{{ $t('landing.downloadDeploy') }}</span>
        </a>
      </div>
    </section>

    <section id="faq" class="section section-alt">
      <div class="section-header">
        <span class="section-label">{{ $t('landing.faqLabel') }}</span>
        <h2>{{ $t('landing.faqTitle') }}</h2>
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
          <p>{{ $t('landing.footerDesc') }}</p>
        </div>
        <div class="footer-col">
          <h4>{{ $t('landing.footerResources') }}</h4>
          <a href="https://github.com/modelbus/one-api-pro" target="_blank">GitHub</a>
          <a href="https://github.com/modelbus/one-api-pro/releases" target="_blank">Releases</a>
        </div>
        <div class="footer-col">
          <h4>{{ $t('landing.footerDocs') }}</h4>
          <a href="#features">{{ $t('landing.footerAdvantages') }}</a>
          <a href="#compare">{{ $t('landing.footerCompare') }}</a>
          <a href="http://one-api.pro" target="_blank">{{ $t('landing.footerOnlineDocs') }}</a>
        </div>
      </div>
      <div class="footer-bottom">
        &copy; {{ new Date().getFullYear() }} One Api Pro &nbsp;|&nbsp; <a href="https://github.com/modelbus/one-api-pro" target="_blank">GitHub</a> &nbsp;|&nbsp; MIT License
      </div>
    </footer>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useStatusStore } from '@/stores/status'
import { IconDashboard, IconApps, IconLock, IconCloud, IconThunderbolt, IconTool, IconCode, IconGithub, IconRight, IconDown, IconCheckCircleFill } from '@arco-design/web-vue/es/icon'
import { PROVIDERS } from '@/constants/providers'
import LangSwitcher from '@/components/LangSwitcher.vue'
import logoPng from '@/assets/logo.png'

const router = useRouter()
const authStore = useAuthStore()
const statusStore = useStatusStore()
const { t } = useI18n()
const systemName = computed(() => statusStore.status?.system_name || 'One Api Pro')

const logoFallback = computed(() => statusStore.status?.logo || '')
const logoSrc = ref(logoPng)
function onLogoError(e) {
  if (logoFallback.value && logoSrc.value !== logoFallback.value) {
    logoSrc.value = logoFallback.value
  } else {
    e.target.style.visibility = 'hidden'
  }
}

const features = computed(() => [
  { icon: IconLock, title: t('landing.featureGovernance'), desc: t('landing.featureGovernanceDesc'), color: '#165DFF' },
  { icon: IconTool, title: t('landing.featureCost'), desc: t('landing.featureCostDesc'), color: '#00B42A' },
  { icon: IconApps, title: t('landing.featureSecurity'), desc: t('landing.featureSecurityDesc'), color: '#FF7D00' },
  { icon: IconThunderbolt, title: t('landing.featureHA'), desc: t('landing.featureHADesc'), color: '#722ED1' },
  { icon: IconCloud, title: t('landing.featureApi'), desc: t('landing.featureApiDesc'), color: '#0FC6C2' },
  { icon: IconDashboard, title: t('landing.featureOps'), desc: t('landing.featureOpsDesc'), color: '#F53F3F' },
])

const compareItems = computed(() => [
  { label: t('landing.compareAdaptorLabel'), old: t('landing.compareAdaptorOld'), new: t('landing.compareAdaptorNew') },
  { label: t('landing.compareSubLabel'), old: t('landing.compareSubOld'), new: t('landing.compareSubNew') },
  { label: t('landing.compareDirLabel'), old: t('landing.compareDirOld'), new: t('landing.compareDirNew') },
  { label: t('landing.compareChannelLabel'), old: t('landing.compareChannelOld'), new: t('landing.compareChannelNew') },
  { label: t('landing.compareClusterLabel'), old: t('landing.compareClusterOld'), new: t('landing.compareClusterNew') },
  { label: t('landing.compareAdminLabel'), old: t('landing.compareAdminOld'), new: t('landing.compareAdminNew') },
  { label: t('landing.compareUpdateLabel'), old: t('landing.compareUpdateOld'), new: t('landing.compareUpdateNew') },
])

const faqs = computed(() => [
  { q: t('landing.faq1Q'), a: t('landing.faq1A') },
  { q: t('landing.faq2Q'), a: t('landing.faq2A') },
  { q: t('landing.faq3Q'), a: t('landing.faq3A') },
  { q: t('landing.faq4Q'), a: t('landing.faq4A') },
  { q: t('landing.faq5Q'), a: t('landing.faq5A') },
])

const models = PROVIDERS

const modelIconMap = import.meta.glob('../assets/lobehub/*.svg', { eager: true, query: '?url', import: 'default' })
function modelIconSrc(slug) {
  const key = `../assets/lobehub/${slug}.svg`
  return modelIconMap[key] || ''
}
function onIconError(e) {
  e.target.style.visibility = 'hidden'
}

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
.nav-logo { height: 36px; border-radius: 6px; }
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
.model-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 12px; }
.model-card {
  padding: 14px 16px;
  border-radius: 12px;
  background: #fff;
  border: 1px solid #e5e6eb;
  display: flex;
  align-items: center;
  gap: 10px;
  transition: all .2s;
  position: relative;
}
.model-card:hover { border-color: #165dff40; box-shadow: 0 2px 12px rgba(22,93,255,0.08); transform: translateY(-1px); }
.mc-icon {
  width: 24px;
  height: 24px;
  border-radius: 5px;
  flex-shrink: 0;
  object-fit: contain;
  color: #1d2129;
}
.mc-icon-fallback {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 13px;
}
.mc-name {
  flex: 1;
  font-size: 13px;
  font-weight: 500;
  color: #1d2129;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.mc-tag {
  font-size: 10px;
  padding: 1px 6px;
  border-radius: 3px;
  font-weight: 500;
  flex-shrink: 0;
}
.mc-tag-cn {
  background: rgba(255, 77, 79, 0.1);
  color: #f53f3f;
}
.mc-tag-intl {
  background: rgba(22, 93, 255, 0.1);
  color: #165dff;
}

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
