<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import ArticleCard from '../components/ArticleCard.vue'

// 更多测试文章数据
const allArticles = [
  {
    id: 1,
    title: '春日里的小确幸：如何发现生活中的美好瞬间',
    summary: '生活中不缺少美，只是缺少发现美的眼睛。在这篇文章中，我将分享如何在日常琐碎中发现那些治愈心灵的小确幸，让每一天都充满温暖与期待。',
    cover: '',
    category: '随笔',
    author: 'Elaina',
    createdAt: '2024-05-15',
    viewCount: 328,
    isPinned: true,
  },
  {
    id: 2,
    title: 'Vue 3 组合式 API 最佳实践指南',
    summary: '深入探讨 Vue 3 Composition API 的核心概念与实际应用，分享我在项目开发中总结的宝贵经验与技巧，帮助你写出更优雅的 Vue 代码。',
    cover: '',
    category: '技术',
    author: 'Elaina',
    createdAt: '2024-05-10',
    viewCount: 256,
    isPinned: false,
  },
  {
    id: 3,
    title: '治愈系前端设计：打造让人舒适的界面',
    summary: '颜色、间距、动效...每一个细节都影响着用户的感受。本文将介绍清新治愈系设计的核心要素，让你的网站也能传递温暖与美好。',
    cover: '',
    category: '技术',
    author: 'Elaina',
    createdAt: '2024-05-05',
    viewCount: 189,
    isPinned: false,
  },
  {
    id: 4,
    title: '我的摄影日记：记录四季的色彩',
    summary: '用镜头捕捉时光，用照片定格美好。分享我一年四季拍摄的照片，以及每张照片背后的故事与感悟。',
    cover: '',
    category: '随笔',
    author: 'Elaina',
    createdAt: '2024-04-28',
    viewCount: 145,
    isPinned: false,
  },
  {
    id: 5,
    title: '读《小王子》：关于爱与责任的思考',
    summary: '重读经典，总有新的感悟。《小王子》不仅是一本童话，更是一本关于人生哲学的书。分享我的阅读笔记与思考。',
    cover: '',
    category: '随笔',
    author: 'Elaina',
    createdAt: '2024-04-20',
    viewCount: 234,
    isPinned: false,
  },
  {
    id: 6,
    title: 'TypeScript 进阶技巧：从入门到精通',
    summary: 'TypeScript 让 JavaScript 开发更加安全和高效。本文将介绍泛型、类型推断、条件类型等高级特性，帮助你掌握 TS 的核心技能。',
    cover: '',
    category: '技术',
    author: 'Elaina',
    createdAt: '2024-04-15',
    viewCount: 412,
    isPinned: false,
  },
  {
    id: 7,
    title: '周末的咖啡馆时光',
    summary: '一个人，一杯咖啡，一本书。周末的午后总是过得特别惬意。记录下那些在咖啡馆度过的美好时光，以及独处时的思考与感悟。',
    cover: '',
    category: '随笔',
    author: 'Elaina',
    createdAt: '2024-04-12',
    viewCount: 178,
    isPinned: false,
  },
  {
    id: 8,
    title: 'CSS Grid 布局完全指南',
    summary: 'Grid 布局是 CSS 最强大的布局系统之一。从基础概念到实战案例，手把手教你用 Grid 创建复杂的响应式页面布局。',
    cover: '',
    category: '技术',
    author: 'Elaina',
    createdAt: '2024-04-08',
    viewCount: 356,
    isPinned: false,
  },
  {
    id: 9,
    title: '雨天的治愈时刻',
    summary: '窗外的雨滴轻轻落下，室内的暖光温馨舒适。雨天总能让人放慢脚步，静静感受生活的美好。分享一些雨天里的小确幸。',
    cover: '',
    category: '随笔',
    author: 'Elaina',
    createdAt: '2024-04-05',
    viewCount: 267,
    isPinned: false,
  },
  {
    id: 10,
    title: '前端性能优化实战',
    summary: '页面加载慢？卡顿？用户体验差？本文从代码分割、懒加载、缓存策略等多个维度，分享前端性能优化的实战经验。',
    cover: '',
    category: '技术',
    author: 'Elaina',
    createdAt: '2024-04-01',
    viewCount: 523,
    isPinned: false,
  },
  {
    id: 11,
    title: '秋日散步：发现城市的另一面',
    summary: '秋天的阳光温暖而不炽热，正是散步的好时节。走过熟悉的街道，却发现了一些从未注意过的美好细节。',
    cover: '',
    category: '随笔',
    author: 'Elaina',
    createdAt: '2024-03-28',
    viewCount: 198,
    isPinned: false,
  },
  {
    id: 12,
    title: 'Node.js 后端开发入门',
    summary: '想要成为全栈开发者？Node.js 是前端开发者进入后端领域的最佳选择。从 Express 到数据库操作，带你快速上手后端开发。',
    cover: '',
    category: '技术',
    author: 'Elaina',
    createdAt: '2024-03-25',
    viewCount: 445,
    isPinned: false,
  },
  {
    id: 13,
    title: '手作时光：做一件属于自己的东西',
    summary: 'DIY 的魅力在于创造的过程。无论是一件手工饰品还是一顿精心准备的晚餐，亲手制作的东西总是格外珍贵。',
    cover: '',
    category: '随笔',
    author: 'Elaina',
    createdAt: '2024-03-20',
    viewCount: 156,
    isPinned: false,
  },
  {
    id: 14,
    title: 'Git 工作流与团队协作',
    summary: '代码版本管理是每个开发者的必修课。本文介绍 Git 常用命令、分支策略以及团队协作的最佳实践，让你的开发流程更加顺畅。',
    cover: '',
    category: '技术',
    author: 'Elaina',
    createdAt: '2024-03-18',
    viewCount: 389,
    isPinned: false,
  },
  {
    id: 15,
    title: '深夜食堂：一人食的小确幸',
    summary: '一个人也要好好吃饭。深夜的厨房，简单的食材，却能烹饪出满满的幸福感。分享几道简单美味的深夜食谱。',
    cover: '',
    category: '随笔',
    author: 'Elaina',
    createdAt: '2024-03-15',
    viewCount: 312,
    isPinned: false,
  },
  {
    id: 16,
    title: 'React Hooks 深度解析',
    summary: 'useState、useEffect、useContext...React Hooks 改变了我们编写组件的方式。深入理解 Hooks 的工作原理，写出更优雅的函数组件。',
    cover: '',
    category: '技术',
    author: 'Elaina',
    createdAt: '2024-03-12',
    viewCount: 478,
    isPinned: false,
  },
  {
    id: 17,
    title: '晨间日记：开启美好一天的仪式',
    summary: '早晨的第一缕阳光，一杯温水，几行文字。晨间日记帮助我整理思绪，设定目标，以更积极的心态迎接新的一天。',
    cover: '',
    category: '随笔',
    author: 'Elaina',
    createdAt: '2024-03-08',
    viewCount: 223,
    isPinned: false,
  },
  {
    id: 18,
    title: 'Docker 容器化部署实战',
    summary: '从开发环境到生产环境，Docker 让部署变得简单可靠。学习 Dockerfile 编写、镜像构建和容器编排，实现一键部署。',
    cover: '',
    category: '技术',
    author: 'Elaina',
    createdAt: '2024-03-05',
    viewCount: 567,
    isPinned: false,
  },
]

const articles = ref(allArticles)

// 分类筛选
const categories = ['全部', '技术', '随笔']
const currentCategory = ref('全部')

const filteredArticles = computed(() => {
  if (currentCategory.value === '全部') {
    return articles.value
  }
  return articles.value.filter(article => article.category === currentCategory.value)
})

// 无限滚动加载
const displayedCount = ref(6)
const pageSize = 6
const isLoading = ref(false)
const hasMore = computed(() => displayedCount.value < filteredArticles.value.length)

const displayedArticles = computed(() => {
  return filteredArticles.value.slice(0, displayedCount.value)
})

const loadMore = () => {
  if (isLoading.value || !hasMore.value) return
  isLoading.value = true
  // 模拟加载延迟
  setTimeout(() => {
    displayedCount.value += pageSize
    isLoading.value = false
  }, 300)
}

// 滚动监听
const handleScroll = () => {
  const scrollHeight = document.documentElement.scrollHeight
  const scrollTop = window.scrollY || document.documentElement.scrollTop
  const clientHeight = window.innerHeight || document.documentElement.clientHeight

  // 距离底部 200px 时触发加载
  if (scrollTop + clientHeight >= scrollHeight - 200) {
    loadMore()
  }
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})

// 获取分类图标
const getCategoryIcon = (category: string) => {
  const icons: Record<string, string> = {
    '全部': '📚',
    '技术': '💻',
    '随笔': '📝',
  }
  return icons[category] || '📄'
}

// 获取分类文章数量
const getCategoryCount = (category: string) => {
  if (category === '全部') {
    return articles.value.length
  }
  return articles.value.filter(article => article.category === category).length
}

// 监听分类变化，重置显示数量
const resetDisplayCount = () => {
  displayedCount.value = pageSize
}

// 快速滚动到顶部
const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// 快速滚动到底部
const scrollToBottom = () => {
  const scrollHeight = document.documentElement.scrollHeight
  window.scrollTo({ top: scrollHeight, behavior: 'smooth' })
}
</script>

<template>
  <main class="home-page">
    <!-- 英雄区域 -->
    <section class="hero-section">
      <div class="hero-decoration">
        <span class="deco-cloud">☁️</span>
        <span class="deco-leaf">🍃</span>
        <span class="deco-flower">🌸</span>
        <span class="deco-star">✨</span>
      </div>

      <div class="hero-content">
        <h1 class="hero-title">
          <span class="title-greeting">你好呀，</span>
          <span class="title-highlight">欢迎来到我的小世界</span>
          <span class="title-emoji">🌿</span>
        </h1>
        <p class="hero-subtitle">
          这里记录着代码与生活的点滴美好
        </p>
        <div class="hero-stats">
          <div class="stat-item">
            <span class="stat-number">{{ articles.length }}</span>
            <span class="stat-label">篇文章</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-number">{{ categories.length - 1 }}</span>
            <span class="stat-label">个分类</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-number">1.2k</span>
            <span class="stat-label">次阅读</span>
          </div>
        </div>
      </div>

      <div class="hero-scroll">
        <span class="scroll-text">向下探索</span>
        <div class="scroll-indicator">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 5v14M19 12l-7 7-7-7"/>
          </svg>
        </div>
      </div>
    </section>

    <!-- 左右分栏文章区域 -->
    <section class="articles-section">
      <div class="section-container two-column">
        <!-- 左侧：分类 + 导航 -->
        <aside class="left-sidebar">
          <div class="category-sidebar">
            <div class="category-list">
              <button
                v-for="category in categories"
                :key="category"
                class="category-card"
                :class="{ active: currentCategory === category }"
                @click="currentCategory = category; resetDisplayCount()"
              >
                <span class="category-icon">{{ getCategoryIcon(category) }}</span>
                <span class="category-name">{{ category }}</span>
                <span class="category-count">{{ getCategoryCount(category) }}篇</span>
              </button>
            </div>
          </div>

          <div class="quick-nav-card">
            <button class="nav-btn" @click="scrollToTop" title="回到顶部">
              <span class="nav-icon">⬆️</span>
            </button>
            <button class="nav-btn" @click="scrollToBottom" title="前往底部">
              <span class="nav-icon">⬇️</span>
            </button>
          </div>
        </aside>

        <!-- 右侧：文章列表 -->
        <div class="articles-main">
          <div class="section-header">
            <h2 class="section-title">
              <span class="title-icon">{{ getCategoryIcon(currentCategory) }}</span>
              {{ currentCategory === '全部' ? '全部文章' : currentCategory + '文章' }}
            </h2>
            <span class="article-count">共 {{ filteredArticles.length }} 篇</span>
          </div>

          <!-- 文章列表 -->
          <div class="articles-list">
            <ArticleCard
              v-for="(article, index) in displayedArticles"
              :key="article.id"
              :article="article"
              :index="index"
            />
          </div>

          <!-- 空状态 -->
          <div v-if="displayedArticles.length === 0" class="empty-state">
            <span class="empty-icon">🌱</span>
            <p class="empty-text">这个分类下还没有文章哦</p>
            <button class="btn-primary" @click="currentCategory = '全部'">
              查看全部文章
            </button>
          </div>

          <!-- 加载更多 -->
          <div v-if="hasMore" class="loading-more">
            <span class="loading-spinner"></span>
            <span class="loading-text">加载中...</span>
          </div>

          <!-- 已加载全部 -->
          <div v-else-if="displayedArticles.length > 0" class="all-loaded">
            <span>已加载全部文章</span>
          </div>
        </div>
      </div>
    </section>
  </main>
</template>

<style scoped>
.home-page {
  min-height: 100vh;
}

/* ===== 英雄区域 ===== */
.hero-section {
  position: relative;
  min-height: 70vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 24px 40px;
  overflow: hidden;
}

/* 装饰元素 */
.hero-decoration {
  position: absolute;
  inset: 0;
  pointer-events: none;
  overflow: hidden;
}

.hero-decoration span {
  position: absolute;
  font-size: 32px;
  opacity: 0.15;
  animation: float 6s ease-in-out infinite;
}

.deco-cloud { top: 15%; left: 10%; animation-delay: 0s; }
.deco-leaf { top: 25%; right: 15%; animation-delay: 1.5s; }
.deco-flower { bottom: 30%; left: 8%; animation-delay: 3s; }
.deco-star { top: 40%; right: 10%; animation-delay: 4.5s; }

@keyframes float {
  0%, 100% {
    transform: translateY(0) rotate(0deg);
  }
  50% {
    transform: translateY(-20px) rotate(5deg);
  }
}

/* 英雄内容 */
.hero-content {
  text-align: center;
  max-width: 600px;
  z-index: 1;
}

.hero-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 16px;
  line-height: 1.3;
}

.title-greeting {
  display: block;
  font-size: 1.25rem;
  font-weight: 400;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.title-highlight {
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.title-emoji {
  display: inline-block;
  margin-left: 8px;
  animation: gentleBounce 2s ease-in-out infinite;
}

@keyframes gentleBounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
}

.hero-subtitle {
  font-size: 1.125rem;
  color: var(--text-secondary);
  margin-bottom: 32px;
  line-height: 1.8;
}

/* 统计数据 */
.hero-stats {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 24px;
  padding: 20px 32px;
  background: var(--bg-glass);
  backdrop-filter: blur(10px);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border);
  box-shadow: var(--shadow-soft);
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.stat-number {
  font-size: 1.75rem;
  font-weight: 700;
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-muted);
}

.stat-divider {
  width: 1px;
  height: 40px;
  background: var(--divider);
}

/* 滚动指示器 */
.hero-scroll {
  position: absolute;
  bottom: 24px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  color: var(--text-muted);
  animation: fadeUp 1s ease 0.5s both;
}

.scroll-text {
  font-size: 0.75rem;
  letter-spacing: 1px;
}

.scroll-indicator {
  animation: bounce 2s infinite;
}

@keyframes bounce {
  0%, 20%, 50%, 80%, 100% { transform: translateY(0); }
  40% { transform: translateY(-8px); }
  60% { transform: translateY(-4px); }
}

/* ===== 文章列表区域 - 左右分栏 ===== */
.articles-section {
  padding: 40px 24px 60px;
}

.section-container {
  max-width: 1200px;
  margin: 0 auto;
}

.section-container.two-column {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 32px;
  align-items: start;
}

/* 左侧整体侧边栏（sticky） */
.left-sidebar {
  position: sticky;
  top: 80px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* 分类卡片 */
.category-sidebar {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 24px;
  box-shadow: var(--shadow-soft);
}

.category-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.category-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  background: transparent;
  border: 1px solid transparent;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--transition-fast);
  text-align: left;
}

.category-card:hover {
  background: var(--bg-secondary);
  border-color: var(--border);
}

.category-card.active {
  background: linear-gradient(135deg, var(--primary-lighter) 0%, rgba(126, 215, 193, 0.1) 100%);
  border-color: var(--primary-light);
}

.category-icon {
  font-size: 20px;
}

.category-name {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.category-count {
  font-size: 12px;
  color: var(--text-muted);
  background: var(--bg-secondary);
  padding: 2px 8px;
  border-radius: var(--radius-full);
}

.category-card.active .category-count {
  background: var(--primary);
  color: white;
}

/* 快速定位卡片（宽度减小75%） */
.quick-nav-card {
  width: 70px;
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 12px 8px;
  box-shadow: var(--shadow-soft);
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-left: auto;
  margin-right: auto;
}

.nav-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--transition-fast);
  font-size: 14px;
  color: var(--text-secondary);
}

.nav-btn:hover {
  background: var(--primary-lighter);
  border-color: var(--primary-light);
  color: var(--primary-dark);
}

.nav-icon {
  font-size: 16px;
}

.nav-label {
  font-weight: 500;
}

/* 右侧文章区域 */
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--border);
}

.section-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
}

.title-icon {
  font-size: 24px;
}

.article-count {
  font-size: 14px;
  color: var(--text-muted);
}

/* 文章列表 */
.articles-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
  margin-bottom: 40px;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 60px 24px;
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-soft);
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
  display: block;
}

.empty-text {
  font-size: 1rem;
  color: var(--text-secondary);
  margin-bottom: 24px;
}

/* 加载更多 */
.loading-more {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 24px;
  color: var(--text-muted);
  font-size: 14px;
}

.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid var(--border);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* 已加载全部 */
.all-loaded {
  text-align: center;
  padding: 24px;
  color: var(--text-muted);
  font-size: 14px;
}

/* ===== 响应式 ===== */
@media (max-width: 968px) {
  .section-container.two-column {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .left-sidebar {
    position: static;
    flex-direction: row;
    gap: 12px;
  }

  .category-sidebar {
    flex: 1;
    padding: 20px;
  }

  .category-list {
    flex-direction: row;
    flex-wrap: wrap;
    gap: 8px;
  }

  .category-card {
    padding: 10px 16px;
  }

  .category-name {
    white-space: nowrap;
  }

  .quick-nav-card {
    width: auto;
    flex-direction: row;
    justify-content: center;
    padding: 12px;
    gap: 12px;
    margin: 0;
  }

  .nav-btn {
    padding: 10px 20px;
  }
}

@media (max-width: 768px) {
  .hero-section {
    min-height: auto;
    padding: 60px 16px 30px;
  }

  .hero-title {
    font-size: 1.75rem;
  }

  .title-greeting {
    font-size: 1rem;
  }

  .hero-subtitle {
    font-size: 1rem;
  }

  .hero-stats {
    gap: 16px;
    padding: 16px 24px;
  }

  .stat-number {
    font-size: 1.25rem;
  }

  .stat-divider {
    height: 30px;
  }

  .hero-scroll {
    display: none;
  }

  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .articles-section {
    padding: 24px 16px 40px;
  }
}
</style>