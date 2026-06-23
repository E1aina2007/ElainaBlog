<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import MarkdownRenderer from '@/components/MarkdownRenderer.vue'
import type { TocItem } from '@/components/MarkdownRenderer.vue'
import TableOfContents from '@/components/TableOfContents.vue'
import CommentForm from '@/components/CommentForm.vue'
import CommentList from '@/components/CommentList.vue'
import { getArticleDetail, getMyArticleDetail, deleteArticle, toggleArticleTop, type Article } from '@/api/article'
import { getComments, createComment, deleteComment, type Comment } from '@/api/comment'
import toast from '@/utils/toast'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const article = ref<Article | null>(null)
const comments = ref<Comment[]>([])
const isLoading = ref(false)
const isCommentLoading = ref(false)
const replyTo = ref<Comment | null>(null)
const error = ref('')
const mdRenderer = ref<InstanceType<typeof MarkdownRenderer> | null>(null)

const tocItems = computed<TocItem[]>(() => mdRenderer.value?.toc ?? [])

const tagList = computed(() => {
  if (!article.value?.tags) return []
  return article.value.tags.split('|').map(t => t.trim()).filter(Boolean)
})

const articleId = computed(() => {
  const id = Number(route.params.id)
  return isNaN(id) ? 0 : id
})

// 格式化日期
function formatDate(date?: string): string {
  if (!date) return ''
  return new Date(date).toLocaleString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

// 加载文章详情
async function loadArticle() {
  if (!articleId.value) {
    error.value = '无效的文章 ID'
    return
  }

  isLoading.value = true
  error.value = ''

  try {
    // 已登录用户使用 getMyArticleDetail，确保返回 user_id 字段
    if (userStore.isLoggedIn) {
      try {
        article.value = await getMyArticleDetail(articleId.value)
      } catch {
        // 如果是 404（不是自己的文章），降级到公开接口
        article.value = await getArticleDetail(articleId.value)
      }
    } else {
      article.value = await getArticleDetail(articleId.value)
    }
  } catch (err: any) {
    error.value = err?.message || '加载文章失败'
  } finally {
    isLoading.value = false
  }
}

// 加载评论
async function loadComments() {
  if (!articleId.value) return

  isCommentLoading.value = true
  try {
    const res = await getComments(articleId.value)
    comments.value = Array.isArray(res) ? res : []
  } catch (err: any) {
    console.error('加载评论失败:', err)
    comments.value = []
  } finally {
    isCommentLoading.value = false
  }
}

// 提交评论
async function handleSubmitComment(data: { content: string; replyToUserId?: number; replyToCommentId?: number }) {
  try {
    await createComment({
      article_id: articleId.value,
      reply_to_user_id: data.replyToUserId,
      reply_to_comment_id: data.replyToCommentId,
      content: data.content,
    })
    replyTo.value = null
    await loadComments()
    toast.success('评论发表成功')
  } catch (err: any) {
    toast.error(err?.message || '发表评论失败')
  }
}

// 回复评论
function handleReply(comment: Comment) {
  replyTo.value = comment
  document.querySelector('.comment-form')?.scrollIntoView({ behavior: 'smooth' })
}

// 取消回复
function cancelReply() {
  replyTo.value = null
}

// 删除评论
async function handleDeleteComment(id: number) {
  try {
    await deleteComment(id)
    // 重新加载评论
    await loadComments()
    toast.success('评论已删除')
  } catch (err: any) {
    toast.error(err?.message || '删除评论失败')
  }
}

// 是否有权编辑/删除文章（作者本人或管理员）
const canEditArticle = computed(() => {
  if (!article.value || !userStore.isLoggedIn) return false
  // 管理员可以编辑所有文章
  if (userStore.isAdmin) return true
  // 作者本人可以编辑自己的文章（user_id 可能是 undefined 或 0，需要安全比较）
  return article.value.user_id != null && article.value.user_id > 0 && article.value.user_id === userStore.userInfo?.id
})

// 删除文章
async function handleDeleteArticle() {
  if (!article.value) return
  if (!confirm(`确定要删除文章「${article.value.title}」吗？文章下的评论将一并删除，此操作不可恢复。`)) return
  try {
    await deleteArticle(article.value.id)
    toast.success('文章已删除')
    setTimeout(() => {
      window.location.href = '/'
    }, 1000)
  } catch (err) {
    toast.error('删除失败')
  }
}

// 编辑文章
function handleEditArticle() {
  if (!article.value) return
  router.push(`/write/${article.value.id}`)
}

// 切换置顶
async function handleTogglePin() {
  if (!article.value) return
  try {
    await toggleArticleTop(article.value.id, !article.value.is_top)
    article.value.is_top = !article.value.is_top
    toast.success(article.value.is_top ? '已置顶' : '已取消置顶')
  } catch (err) {
    toast.error('操作失败')
  }
}

// 返回首页
function goBack() {
  window.location.href = '/'
}

// 分享文章
function handleShare() {
  navigator.clipboard.writeText(window.location.href).then(() => {
    toast.success('链接已复制到剪贴板')
  }).catch(() => {
    toast.error('复制失败，请手动复制')
  })
}

onMounted(() => {
  loadArticle()
  loadComments()
})
</script>

<template>
  <main class="article-detail-page">
    <div class="container">
      <!-- 加载状态 -->
      <div v-if="isLoading" class="loading-state">
        <div class="spinner" />
        <p>加载中...</p>
      </div>

      <!-- 错误状态 -->
      <div v-else-if="error" class="error-state">
        <p>{{ error }}</p>
        <button class="back-btn" @click="goBack">
          返回首页
        </button>
      </div>

      <!-- 文章内容 -->
      <article v-else-if="article" class="article">
        <!-- 文章头部 -->
        <header class="article-header">
          <h1 class="article-title">{{ article.title }}</h1>
          <div class="article-meta">
            <div class="meta-left">
              <span class="meta-item author">
                <img v-if="article.author_avatar" :src="article.author_avatar" class="avatar-img" alt="" />
                <span v-else class="avatar">{{ (article.author_name || '作者').charAt(0) }}</span>
                <span>{{ article.author_name || '作者' }}</span>
                <span v-if="article.author_is_admin" class="admin-badge">管理员</span>
              </span>
              <span class="meta-separator">·</span>
              <time class="meta-item">{{ formatDate(article.created_at) }}</time>
              <span class="meta-separator">·</span>
              <span class="meta-item">{{ article.view_count || 0 }} 阅读</span>
              <template v-if="article.category_name">
                <span class="meta-separator">·</span>
                <span class="meta-item category">{{ article.category_name }}</span>
              </template>
            </div>
            <div v-if="tagList.length > 0" class="article-tags">
              <span v-for="(tag, i) in tagList" :key="i" class="article-tag">{{ tag }}</span>
            </div>
            <div class="article-actions">
              <button class="share-article-btn" @click="handleShare">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="18" cy="5" r="3" />
                  <circle cx="6" cy="12" r="3" />
                  <circle cx="18" cy="19" r="3" />
                  <line x1="8.59" y1="13.51" x2="15.42" y2="17.49" />
                  <line x1="15.41" y1="6.51" x2="8.59" y2="10.49" />
                </svg>
                分享
              </button>
            </div>
            <div class="article-actions">
              <!-- 置顶按钮：仅管理员可见 -->
              <button
                v-if="userStore.isAdmin"
                class="pin-article-btn"
                :class="{ pinned: article.is_top }"
                @click="handleTogglePin"
              >
                <svg width="14" height="14" viewBox="0 0 24 24" :fill="article.is_top ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2">
                  <path d="M16 12V4H17V2H7V4H8V12L6 14V16H11.2V22H12.8V16H18V14L16 12Z"/>
                </svg>
                {{ article.is_top ? '取消置顶' : '置顶' }}
              </button>
              <!-- 编辑和删除按钮：作者本人或管理员可见 -->
              <template v-if="canEditArticle">
                <button class="edit-article-btn" @click="handleEditArticle">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                  </svg>
                  修改文章
                </button>
                <button class="delete-article-btn" @click="handleDeleteArticle">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3 6 5 6 21 6"></polyline>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                  </svg>
                  删除文章
                </button>
              </template>
            </div>
          </div>
        </header>

        <!-- 文章内容 -->
        <div class="article-content">
          <MarkdownRenderer ref="mdRenderer" :content="article.content" />
        </div>

        <!-- 文章底部 -->
        <footer class="article-footer">
          <button class="back-link" @click="goBack">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M19 12H5M12 19l-7-7 7-7" />
            </svg>
            返回首页
          </button>
        </footer>
      </article>

      <!-- 文章目录 -->
      <TableOfContents :items="tocItems" />

      <!-- 评论区 -->
      <section v-if="article" class="comments-section">
        <CommentForm
          :article-id="articleId"
          :reply-to="replyTo"
          @submit="handleSubmitComment"
          @cancel-reply="cancelReply"
        />
        <CommentList
          :comments="comments"
          :loading="isCommentLoading"
          @reply="handleReply"
          @delete="handleDeleteComment"
        />
      </section>
    </div>
  </main>
</template>

<style scoped>
.article-detail-page {
  min-height: 100vh;
  padding: 40px 24px;
  background: var(--bg-primary);
}

.container {
  max-width: 800px;
  margin: 0 auto;
}

/* 加载和错误状态 */
.loading-state,
.error-state {
  text-align: center;
  padding: 80px 24px;
  color: var(--text-secondary);
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--border);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.back-btn {
  margin-top: 16px;
  padding: 10px 24px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.back-btn:hover {
  background: var(--primary-dark);
}

/* 文章卡片 */
.article {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 40px;
  box-shadow: var(--shadow-soft);
  margin-bottom: 24px;
}

/* 文章头部 */
.article-header {
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--border);
}

.article-title {
  font-size: 1.875rem;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.4;
  margin: 0 0 16px;
}

.article-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 8px;
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.meta-left {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.article-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.edit-article-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 14px;
  background: transparent;
  color: var(--primary);
  border: 1px solid var(--primary);
  border-radius: var(--radius-md);
  font-size: 0.8125rem;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.edit-article-btn:hover {
  background: var(--primary);
  color: white;
}

.share-article-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 14px;
  background: transparent;
  color: var(--primary);
  border: 1px solid var(--primary);
  border-radius: var(--radius-md);
  font-size: 0.8125rem;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.share-article-btn:hover {
  background: var(--primary);
  color: white;
}

.pin-article-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 14px;
  background: transparent;
  color: var(--color-warning);
  border: 1px solid var(--color-warning);
  border-radius: var(--radius-md);
  font-size: 0.8125rem;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.pin-article-btn:hover {
  background: var(--color-warning);
  color: white;
}

.pin-article-btn.pinned {
  background: var(--color-warning);
  color: white;
}

.delete-article-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 14px;
  background: transparent;
  color: var(--color-danger);
  border: 1px solid var(--color-danger);
  border-radius: var(--radius-md);
  font-size: 0.8125rem;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.delete-article-btn:hover {
  background: var(--color-danger);
  color: white;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.meta-item.author {
  font-weight: 500;
}

.admin-badge {
  display: inline-flex;
  align-items: center;
  padding: 1px 8px;
  background: linear-gradient(135deg, var(--color-warning) 0%, var(--color-warning-dark) 100%);
  color: white;
  font-size: 11px;
  font-weight: 600;
  border-radius: 10px;
  line-height: 1.6;
}

.avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
}

.avatar-img {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  object-fit: cover;
}

.meta-separator {
  color: var(--text-muted);
}

.meta-item.category {
  padding: 2px 10px;
  background: var(--primary-lighter);
  color: var(--primary-dark);
  border-radius: var(--radius-sm);
  font-size: 0.75rem;
  font-weight: 500;
}

.article-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.article-tag {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  background: var(--bg-secondary);
  color: var(--text-muted);
  font-size: 12px;
  font-weight: 500;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
}

/* 文章内容 */
.article-content {
  line-height: 1.8;
  color: var(--text-secondary);
}

/* 文章底部 */
.article-footer {
  margin-top: 40px;
  padding-top: 24px;
  border-top: 1px solid var(--border);
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--text-secondary);
  font-size: 0.875rem;
  text-decoration: none;
  cursor: pointer;
  background: none;
  border: none;
  padding: 8px 12px;
  border-radius: var(--radius-md);
  transition: all 0.2s ease;
}

.back-link:hover {
  color: var(--primary);
  background: var(--bg-secondary);
}

/* 评论区 */
.comments-section {
  margin-top: 32px;
}

/* 响应式 */
@media (max-width: 768px) {
  .article-detail-page {
    padding: 24px 16px;
  }

  .article {
    padding: 24px;
  }

  .article-title {
    font-size: 1.5rem;
  }

  .article-meta {
    font-size: 0.8125rem;
  }
}
</style>