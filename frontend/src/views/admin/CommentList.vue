<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import request from '@/api/request'
import toast from '@/utils/toast'

interface Comment {
  id: number
  article_id: number
  article_title?: string
  user_id: number
  username: string
  content: string
  created_at: string
  is_pending?: boolean
}

const comments = ref<Comment[]>([])
const loading = ref(false)
const searchQuery = ref('')

const fetchComments = async () => {
  loading.value = true
  try {
    // 获取所有评论列表（需要后端支持）
    const data = await request.get('/comment/list') as Comment[]
    comments.value = data || []
  } catch (error) {
    console.error('获取评论列表失败:', error)
    // 如果接口不存在，使用空数组
    comments.value = []
  } finally {
    loading.value = false
  }
}

const filteredComments = computed(() => {
  if (!searchQuery.value) return comments.value
  const query = searchQuery.value.toLowerCase()
  return comments.value.filter(comment =>
    comment.content.toLowerCase().includes(query) ||
    comment.username.toLowerCase().includes(query) ||
    comment.article_title?.toLowerCase().includes(query)
  )
})

const handleDelete = async (comment: Comment) => {
  if (!confirm('确定要删除这条评论吗？')) return
  try {
    await request.post('/comment/delete', { id: comment.id })
    comments.value = comments.value.filter(c => c.id !== comment.id)
    toast.success('删除成功')
  } catch (error) {
    console.error('删除失败:', error)
    toast.error('删除失败')
  }
}

const handleViewArticle = (articleId: number) => {
  window.open(`/article/${articleId}`, '_blank')
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleString('zh-CN')
}

const formatContent = (content: string) => {
  return content.length > 100 ? content.substring(0, 100) + '...' : content
}

onMounted(fetchComments)
</script>

<template>
  <div class="comment-list">
    <div class="page-header">
      <h2>评论管理</h2>
      <div class="search-box">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
        </svg>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="搜索评论内容、用户名..."
        />
      </div>
    </div>

    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>评论内容</th>
            <th>评论者</th>
            <th>所属文章</th>
            <th>评论时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="comment in filteredComments" :key="comment.id">
            <td class="content-cell">
              <div class="comment-content">{{ formatContent(comment.content) }}</div>
            </td>
            <td>{{ comment.username }}</td>
            <td>
              <a @click="handleViewArticle(comment.article_id)" class="article-link">
                {{ comment.article_title || '查看文章' }}
              </a>
            </td>
            <td>{{ formatDate(comment.created_at) }}</td>
            <td class="actions">
              <button class="action-btn delete" @click="handleDelete(comment)" title="删除">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3 6 5 6 21 6"></polyline>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                </svg>
              </button>
            </td>
          </tr>
          <tr v-if="filteredComments.length === 0">
            <td colspan="5" class="empty-cell">
              <div class="empty-state">
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"></path>
                </svg>
                <p>暂无评论</p>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.comment-list {
  max-width: 1400px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.page-header h2 {
  margin: 0;
  font-size: 24px;
  color: #1f2937;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  min-width: 300px;
}

.search-box svg {
  color: #9ca3af;
}

.search-box input {
  flex: 1;
  border: none;
  outline: none;
  font-size: 14px;
  color: #374151;
}

.table-container {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 16px;
  text-align: left;
  font-size: 14px;
}

.data-table th {
  background: #f9fafb;
  font-weight: 600;
  color: #374151;
  border-bottom: 1px solid #e5e7eb;
}

.data-table td {
  color: #6b7280;
  border-bottom: 1px solid #f3f4f6;
}

.data-table tr:hover td {
  background: #f9fafb;
}

.content-cell {
  max-width: 400px;
}

.comment-content {
  line-height: 1.5;
  color: #374151;
}

.article-link {
  color: #6366f1;
  cursor: pointer;
  text-decoration: none;
}

.article-link:hover {
  text-decoration: underline;
}

.actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 6px;
  background: #f3f4f6;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:hover {
  background: #e5e7eb;
  color: #374151;
}

.action-btn.delete:hover {
  background: #ef444415;
  color: #ef4444;
}

.empty-cell {
  padding: 60px 16px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: #9ca3af;
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}
</style>
