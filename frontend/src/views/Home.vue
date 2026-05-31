<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useSiteStore } from '@/stores/site'
import { getArticleList, type Article } from '@/api/article'
import { getCategoryList, type Category } from '@/api/category'
import ArticleCard from '../components/ArticleCard.vue'

const router = useRouter()
const userStore = useUserStore()
const siteStore = useSiteStore()

const articles = ref<Article[]>([])
const categories = ref<Category[]>([])
const totalArticles = ref(0)
const totalViews = ref(0)
const currentCategoryId = ref<number | null>(null)
const allArticlesTotal = ref(0)
const currentPage = ref(1)
const pageSize = 12
const isLoading = ref(false)
const isInitialLoading = ref(true)
const hasMore = ref(true)

const greeting = computed(() => siteStore.get('greeting'))
const heroTitle = computed(() => siteStore.get('hero_title'))

// 随机选择一句话
const randomQuote = ref('"风静黄昏浪息间"')
onMounted(() => {
  const quotes = siteStore.getQuotes()
  if (quotes.length > 0) {
    randomQuote.value = quotes[Math.floor(Math.random() * quotes.length)] ?? '"风静黄昏浪息间"'
  }
})

// 拉取分类列表
const fetchCategories = async () => {
  try {
    categories.value = await getCategoryList()
  } catch (e) {
    console.error('获取分类失败:', e)
  }
}

// 拉取文章列表（追加模式）
const fetchArticles = async (page: number, append = false) => {
  isLoading.value = true
  try {
    const params: { page: number; pageSize: number; categoryId?: number } = { page, pageSize }
    if (currentCategoryId.value !== null) {
      params.categoryId = currentCategoryId.value
    }
    const res = await getArticleList(params)
    const list = res.list || []
    if (append) {
      articles.value = [...articles.value, ...list]
    } else {
      articles.value = list
    }
    totalArticles.value = res.total ?? articles.value.length
    if (currentCategoryId.value === null) {
      allArticlesTotal.value = totalArticles.value
    }
    totalViews.value = articles.value.reduce((sum, a) => sum + (a.view_count || 0), 0)
    hasMore.value = list.length >= pageSize
  } catch (e) {
    console.error('获取文章列表失败:', e)
    if (!append) {
      articles.value = []
      totalArticles.value = 0
      totalViews.value = 0
    }
    hasMore.value = false
  } finally {
    isLoading.value = false
    isInitialLoading.value = false
  }
}

// 切换分类时重新加载
const switchCategory = (categoryId: number | null) => {
  currentCategoryId.value = categoryId
  currentPage.value = 1
  hasMore.value = true
  fetchArticles(1)
}

// 无限滚动加载更多
const loadMore = () => {
  if (isLoading.value || !hasMore.value) return
  currentPage.value++
  fetchArticles(currentPage.value, true)
}

// 写文章（需登录）
const goToWrite = () => {
  if (!userStore.isLoggedIn) {
    router.push({ name: 'login', query: { redirect: '/write' } })
    return
  }
  router.push('/write')
}

// 滚动监听
const handleScroll = () => {
  const scrollTop = window.scrollY || document.documentElement.scrollTop
  const clientHeight = window.innerHeight || document.documentElement.clientHeight
  const scrollHeight = document.documentElement.scrollHeight

  if (scrollTop + clientHeight >= scrollHeight - 200) {
    loadMore()
  }
}

onMounted(() => {
  fetchCategories()
  fetchArticles(1)
  window.addEventListener('scroll', handleScroll)
  // 安全兜底：确保加载状态不会永远卡住
  setTimeout(() => {
    isInitialLoading.value = false
  }, 8000)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<template>
  <main class="home-page">
    <!-- 英雄区域 -->
    <section class="hero-section">
      <div class="hero-decoration">
        <!-- 飘落的花叶装饰 -->
        <span class="deco-leaf leaf-1">🍃</span>
        <span class="deco-leaf leaf-2">🍃</span>
        <span class="deco-leaf leaf-3">🍃</span>
        <span class="deco-flower flower-1">🌸</span>
        <span class="deco-flower flower-2">🌸</span>
        <span class="deco-flower flower-3">🌸</span>
        <span class="deco-petal petal-1">🌿</span>
        <span class="deco-petal petal-2">🌿</span>
        <span class="deco-petal petal-3">🌿</span>
        <span class="deco-flower flower-4">🌼</span>
        <span class="deco-flower flower-5">🌺</span>
      </div>

      <div class="hero-content">
        <h1 class="hero-title">
          <span class="title-greeting">{{ greeting }}</span>
          <span class="title-highlight">{{ heroTitle }}</span>
          <span class="title-emoji"><img src="/favicon.ico" alt="logo" /></span>
        </h1>
        <p class=”hero-subtitle”>
          {{ randomQuote }}
        </p>
        <div class="hero-stats">
          <div class="stat-item">
            <span class="stat-number">{{ allArticlesTotal }}</span>
            <span class="stat-label">篇文章</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-number">{{ categories.length }}</span>
            <span class="stat-label">个分类</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-number">{{ totalViews }}</span>
            <span class="stat-label">次阅读</span>
          </div>
        </div>
      </div>

      <div class="hero-scroll">
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
        <!-- 左侧：分类 -->
        <aside class="left-sidebar">
          <div class="category-sidebar">
            <div class="category-list">
              <button
                class="category-card"
                :class="{ active: currentCategoryId === null }"
                @click="switchCategory(null)"
              >
                <span class="category-name">全部</span>
                <span class="category-count">{{ allArticlesTotal }}篇</span>
              </button>
              <button
                v-for="cat in categories"
                :key="cat.id"
                class="category-card"
                :class="{ active: currentCategoryId === cat.id }"
                @click="switchCategory(cat.id)"
              >
                <span class="category-name">{{ cat.name }}</span>
                <span class="category-count">{{ cat.article_count }}篇</span>
              </button>
            </div>
          </div>
        </aside>

        <!-- 右侧：文章列表 -->
        <div class="articles-main">
          <div class="section-header">
            <h2 class="section-title">{{ currentCategoryId === null ? '全部文章' : (categories.find(c => c.id === currentCategoryId)?.name ?? '') + '文章' }}</h2>
            <div class="header-actions">
              <span class="article-count">共 {{ totalArticles }} 篇</span>
              <button class="btn-write" title="写新文章" @click="goToWrite">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="12" y1="5" x2="12" y2="19"/>
                  <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                <span>写文章</span>
              </button>
            </div>
          </div>

          <!-- 初始加载 -->
          <div v-if="isInitialLoading" class="loading-more">
            <span class="loading-spinner"></span>
            <span class="loading-text">加载中...</span>
          </div>

          <!-- 文章列表 -->
          <div v-else class="articles-list">
            <ArticleCard
              v-for="(article, index) in articles"
              :key="article.id"
              :article="article"
              :index="index"
            />
          </div>

          <!-- 空状态 -->
          <div v-if="!isInitialLoading && articles.length === 0" class="empty-state">
            <p class="empty-text">这个分类下还没有文章哦</p>
            <button class="btn-primary" @click="switchCategory(null)">
              查看全部文章
            </button>
          </div>

          <!-- 加载更多 -->
          <div v-if="!isInitialLoading && isLoading" class="loading-more">
            <span class="loading-spinner"></span>
            <span class="loading-text">加载中...</span>
          </div>

          <!-- 已加载全部 -->
          <div v-else-if="!hasMore && articles.length > 0" class="all-loaded">
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
  font-size: 28px;
  opacity: 0.2;
  animation: fall 8s ease-in-out infinite;
}

/* 叶子位置与动画延迟 */
.deco-leaf.leaf-1 { top: -10%; left: 10%; animation-delay: 0s; font-size: 24px; }
.deco-leaf.leaf-2 { top: -15%; left: 35%; animation-delay: 2s; font-size: 20px; }
.deco-leaf.leaf-3 { top: -5%; right: 20%; animation-delay: 4s; font-size: 26px; }

/* 樱花位置与动画延迟 */
.deco-flower.flower-1 { top: -12%; left: 55%; animation-delay: 1s; font-size: 22px; }
.deco-flower.flower-2 { top: -8%; right: 15%; animation-delay: 3s; font-size: 24px; }
.deco-flower.flower-3 { top: -18%; left: 25%; animation-delay: 5s; font-size: 20px; }
.deco-flower.flower-4 { top: -6%; right: 40%; animation-delay: 6s; font-size: 22px; }
.deco-flower.flower-5 { top: -14%; left: 70%; animation-delay: 7s; font-size: 26px; }

/* 飘落叶片 */
.deco-petal.petal-1 { top: -20%; left: 45%; animation-delay: 1.5s; font-size: 18px; }
.deco-petal.petal-2 { top: -10%; right: 30%; animation-delay: 3.5s; font-size: 20px; }
.deco-petal.petal-3 { top: -16%; left: 5%; animation-delay: 5.5s; font-size: 18px; }

/* 飘落动画 - 从上往下飘落并摇摆 */
@keyframes fall {
  0% {
    transform: translateY(0) translateX(0) rotate(0deg);
    opacity: 0;
  }
  10% {
    opacity: 0.25;
  }
  25% {
    transform: translateY(25vh) translateX(15px) rotate(45deg);
  }
  50% {
    transform: translateY(50vh) translateX(-10px) rotate(90deg);
  }
  75% {
    transform: translateY(75vh) translateX(20px) rotate(135deg);
  }
  90% {
    opacity: 0.25;
  }
  100% {
    transform: translateY(110vh) translateX(-5px) rotate(180deg);
    opacity: 0;
  }
}

/* 保留轻微的浮动效果 */
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

.title-emoji img {
  width: 32px;
  height: 32px;
  object-fit: contain;
}

@keyframes gentleBounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
}

.hero-subtitle {
  font-size: 1.125rem;
  font-family: 'Microsoft YaHei', sans-serif;
  font-weight: 500;
  color: var(--text-muted);
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
  z-index: 10;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  color: var(--text-muted);
  animation: fadeUp 1s ease 0.5s both;
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
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.article-count {
  font-size: 14px;
  color: var(--text-muted);
}

/* 写文章按钮 */
.btn-write {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  transition: all var(--transition-fast);
}

.btn-write:hover {
  background: var(--primary-dark);
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(126, 215, 193, 0.4);
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