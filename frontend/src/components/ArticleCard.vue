<script setup lang="ts">
interface Article {
  id: number
  title: string
  summary: string
  cover?: string
  category: string
  author: string
  createdAt: string
  viewCount: number
  isPinned?: boolean
}

interface Props {
  article: Article
  index?: number
}

withDefaults(defineProps<Props>(), {
  index: 0,
})

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric',
  })
}
</script>

<template>
  <article
    class="article-card"
    :class="{ pinned: article.isPinned, 'stagger-enter': true }"
    :style="`animation-delay: ${index * 100}ms`"
  >
    <!-- 置顶标记 -->
    <div v-if="article.isPinned" class="pin-badge">
      <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
        <path d="M16 12V4H17V2H7V4H8V12L6 14V16H11.2V22H12.8V16H18V14L16 12Z"/>
      </svg>
      <span>置顶</span>
    </div>

    <router-link :to="`/article/${article.id}`" class="card-link">
      <!-- 内容区 -->
      <div class="card-content">
        <!-- 顶部信息：作者 + 日期 -->
        <div class="card-header">
          <div class="author-info">
            <span class="author-avatar">👤</span>
            <span class="author-name">{{ article.author }}</span>
          </div>
          <time class="publish-time">{{ formatDate(article.createdAt) }}</time>
        </div>

        <!-- 标题 -->
        <h3 class="card-title">{{ article.title }}</h3>

        <!-- 分类标签 -->
        <div class="card-category">
          <span class="category-tag">{{ article.category }}</span>
        </div>

        <!-- 摘要 -->
        <p class="card-summary">{{ article.summary }}</p>
      </div>
    </router-link>
  </article>
</template>

<style scoped>
.article-card {
  position: relative;
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-soft);
  transition: all var(--transition-base);
  opacity: 0;
  animation: fadeUp 0.5s cubic-bezier(0.4, 0, 0.2, 1) forwards;
}

.article-card:hover {
  transform: translateY(-6px);
  box-shadow: var(--shadow-hover);
}

.article-card.pinned {
  border: 2px solid var(--primary-light);
}

/* 置顶标记 */
.pin-badge {
  position: absolute;
  top: 12px;
  left: 12px;
  z-index: 10;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  background: linear-gradient(135deg, var(--accent) 0%, var(--highlight) 100%);
  color: white;
  font-size: 12px;
  font-weight: 600;
  border-radius: var(--radius-sm);
  box-shadow: 0 2px 8px rgba(255, 154, 162, 0.3);
}

/* 卡片链接 */
.card-link {
  display: block;
  text-decoration: none;
  color: inherit;
  height: 100%;
}

/* 内容区域 */
.card-content {
  padding: 20px;
  display: flex;
  flex-direction: column;
}

/* 顶部信息：作者 + 日期 */
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.author-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.author-avatar {
  font-size: 16px;
}

.author-name {
  font-size: 13px;
  color: var(--text-secondary);
  font-weight: 500;
}

.publish-time {
  font-size: 12px;
  color: var(--text-muted);
}

/* 标题 */
.card-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  line-height: 1.5;
  margin-bottom: 10px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color var(--transition-fast);
}

.article-card:hover .card-title {
  color: var(--primary);
}

/* 分类标签 */
.card-category {
  margin-bottom: 12px;
}

.category-tag {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  background: var(--primary-lighter);
  color: var(--primary-dark);
  font-size: 12px;
  font-weight: 500;
  border-radius: var(--radius-sm);
}

/* 摘要 */
.card-summary {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.8;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 动画 */
@keyframes fadeUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 响应式 */
@media (max-width: 768px) {
  .card-content {
    padding: 16px;
  }

  .card-title {
    font-size: 16px;
  }

  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 6px;
  }
}
</style>