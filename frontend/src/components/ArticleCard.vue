<script setup lang="ts">
interface Article {
  id: number
  title: string
  summary: string
  cover?: string
  category_name?: string
  author_name?: string
  author_avatar?: string
  author_is_admin?: boolean
  created_at?: string
  view_count?: number
  comment_count?: number
  is_top?: boolean
}

interface Props {
  article: Article
  index?: number
}

withDefaults(defineProps<Props>(), {
  index: 0,
})

const formatDate = (date?: string) => {
  if (!date) return ''
  return new Date(date).toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric',
  })
}
</script>

<template>
  <article
    class="article-card"
    :class="{ pinned: article.is_top, 'stagger-enter': true }"
    :style="`animation-delay: ${index * 100}ms`"
  >
    <router-link :to="`/article/${article.id}`" class="card-link">
      <!-- 内容区 -->
      <div class="card-content">
        <!-- 顶部信息：作者 + 日期 -->
        <div class="card-header">
          <div class="author-info">
            <img v-if="article.author_avatar" :src="article.author_avatar" class="author-avatar-img" alt="" />
            <span v-else class="author-avatar">👤</span>
            <span class="author-name">{{ article.author_name }}</span>
            <span v-if="article.author_is_admin" class="admin-badge">管理员</span>
          </div>
          <time class="publish-time">{{ formatDate(article.created_at) }}</time>
        </div>

        <!-- 标题 + 置顶标记 -->
        <div class="title-row">
          <h3 class="card-title">{{ article.title }}</h3>
          <span v-if="article.is_top" class="pin-badge">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor">
              <path d="M16 12V4H17V2H7V4H8V12L6 14V16H11.2V22H12.8V16H18V14L16 12Z"/>
            </svg>
            <span>置顶</span>
          </span>
        </div>

        <!-- 分类标签 + 统计 -->
        <div class="card-category">
          <span class="category-tag">{{ article.category_name }}</span>
          <div class="card-stats">
            <span class="stat-item">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                <circle cx="12" cy="12" r="3"/>
              </svg>
              {{ article.view_count || 0 }}
            </span>
            <span class="stat-item">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
              </svg>
              {{ article.comment_count || 0 }}
            </span>
          </div>
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

/* 标题行（标题 + 置顶标记） */
.title-row {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 10px;
}

/* 置顶标记 */
.pin-badge {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: linear-gradient(135deg, var(--accent) 0%, var(--highlight) 100%);
  color: white;
  font-size: 11px;
  font-weight: 600;
  border-radius: var(--radius-sm);
  margin-top: 2px;
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

.author-avatar-img {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  object-fit: cover;
}

.author-name {
  font-size: 13px;
  color: var(--text-secondary);
  font-weight: 500;
}

.admin-badge {
  display: inline-flex;
  align-items: center;
  padding: 1px 8px;
  background: linear-gradient(135deg, var(--color-warning) 0%, var(--color-warning-dark) 100%);
  color: white;
  font-size: 10px;
  font-weight: 600;
  border-radius: 10px;
  line-height: 1.6;
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
  flex: 1;
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

/* 分类标签 + 统计 */
.card-category {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.card-stats {
  display: flex;
  align-items: center;
  gap: 12px;
}

.card-stats .stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--text-muted);
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