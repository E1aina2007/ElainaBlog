<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getAdminArticleList, deleteArticle, toggleArticleTop, type Article } from '@/api/article'
import { getCategoryList, type Category } from '@/api/category'
import toast from '@/utils/toast'

const router = useRouter()

const articles = ref<Article[]>([])
const categories = ref<Category[]>([])
const loading = ref(false)
const searchQuery = ref('')
const selectedCategory = ref<number | ''>('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const fetchArticles = async () => {
  loading.value = true
  try {
    const result = await getAdminArticleList({
      page: currentPage.value,
      pageSize: pageSize.value,
      categoryId: selectedCategory.value || undefined,
    })
    articles.value = result.list
    total.value = result.total
  } catch (error) {
    console.error('获取文章列表失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchCategories = async () => {
  try {
    categories.value = await getCategoryList()
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

const filteredArticles = computed(() => {
  if (!searchQuery.value) return articles.value
  const query = searchQuery.value.toLowerCase()
  return articles.value.filter(article =>
    article.title.toLowerCase().includes(query) ||
    article.summary?.toLowerCase().includes(query)
  )
})

const handleCreate = () => {
  router.push('/admin/article/create')
}

const handleEdit = (id: number) => {
  router.push(`/admin/article/edit/${id}`)
}

const handleDelete = async (article: Article) => {
  if (!confirm(`确定要删除文章 "${article.title}" 吗？此操作不可恢复。`)) {
    return
  }
  try {
    await deleteArticle(article.id)
    fetchArticles()
    toast.success('删除成功')
  } catch (error) {
    console.error('删除失败:', error)
    toast.error('删除失败')
  }
}

const handleToggleTop = async (article: Article) => {
  try {
    await toggleArticleTop(article.id, !article.is_top)
    article.is_top = !article.is_top
    toast.success(article.is_top ? '已置顶' : '已取消置顶')
  } catch (error) {
    console.error('操作失败:', error)
    toast.error('操作失败')
  }
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('zh-CN')
}

watch([selectedCategory, currentPage], fetchArticles)

onMounted(() => {
  fetchArticles()
  fetchCategories()
})
</script>

<template>
  <div class="article-list">
    <div class="page-header">
      <h2>文章管理</h2>
      <button class="btn-primary" @click="handleCreate">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="12" y1="5" x2="12" y2="19"></line>
          <line x1="5" y1="12" x2="19" y2="12"></line>
        </svg>
        写文章
      </button>
    </div>

    <div class="filter-bar">
      <div class="search-box">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
        </svg>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="搜索文章标题或摘要..."
        />
      </div>
      <select v-model="selectedCategory" class="category-select">
        <option value="">全部分类</option>
        <option v-for="cat in categories" :key="cat.id" :value="cat.id">
          {{ cat.name }}
        </option>
      </select>
    </div>

    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>标题</th>
            <th>分类</th>
            <th>作者</th>
            <th>阅读量</th>
            <th>UV</th>
            <th>状态</th>
            <th>发布时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="article in filteredArticles" :key="article.id">
            <td class="title-cell">
              <div class="article-title">
                <span v-if="article.is_top" class="top-badge">置顶</span>
                <span v-if="article.is_draft" class="draft-badge">草稿</span>
                {{ article.title }}
              </div>
            </td>
            <td>{{ article.category_name || '未分类' }}</td>
            <td>{{ article.author_name }}</td>
            <td>{{ article.view_count || 0 }}</td>
            <td>{{ article.uv_count || 0 }}</td>
            <td>
              <span :class="['status-badge', article.is_draft ? 'draft' : 'published']">
                {{ article.is_draft ? '草稿' : '已发布' }}
              </span>
            </td>
            <td>{{ formatDate(article.created_at!) }}</td>
            <td class="actions">
              <button class="action-btn" @click="handleToggleTop(article)" title="置顶/取消置顶">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
                </svg>
              </button>
              <button class="action-btn edit" @click="handleEdit(article.id)" title="编辑">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                </svg>
              </button>
              <button class="action-btn delete" @click="handleDelete(article)" title="删除">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3 6 5 6 21 6"></polyline>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                </svg>
              </button>
            </td>
          </tr>
          <tr v-if="filteredArticles.length === 0">
            <td colspan="8" class="empty-cell">
              <div class="empty-state">
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                  <polyline points="14 2 14 8 20 8"></polyline>
                </svg>
                <p>暂无文章</p>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="pagination" v-if="total > pageSize">
      <button
        :disabled="currentPage === 1"
        @click="currentPage--"
        class="page-btn"
      >
        上一页
      </button>
      <span class="page-info">第 {{ currentPage }} 页，共 {{ Math.ceil(total / pageSize) }} 页</span>
      <button
        :disabled="currentPage >= Math.ceil(total / pageSize)"
        @click="currentPage++"
        class="page-btn"
      >
        下一页
      </button>
    </div>
  </div>
</template>

<style scoped>
.article-list {
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
  color: var(--text-primary);
}

.btn-primary {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: var(--color-indigo);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover {
  background: var(--color-indigo-hover);
  transform: translateY(-1px);
}

.filter-bar {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.search-box {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  background: var(--bg-card);
  border: 1px solid var(--input-border);
  border-radius: 8px;
}

.search-box svg {
  color: var(--text-muted);
}

.search-box input {
  flex: 1;
  border: none;
  outline: none;
  font-size: 14px;
  color: var(--text-primary);
}

.category-select {
  padding: 10px 16px;
  background: var(--bg-card);
  border: 1px solid var(--input-border);
  border-radius: 8px;
  font-size: 14px;
  color: var(--text-primary);
  cursor: pointer;
  min-width: 140px;
}

.table-container {
  background: var(--bg-card);
  border-radius: 12px;
  box-shadow: var(--shadow-card);
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
  background: var(--toolbar-bg);
  font-weight: 600;
  color: var(--text-primary);
  border-bottom: 1px solid var(--input-border);
}

.data-table td {
  color: var(--text-secondary);
  border-bottom: 1px solid var(--border);
}

.data-table tr:hover td {
  background: var(--toolbar-bg);
}

.title-cell {
  max-width: 300px;
}

.article-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  color: var(--text-primary);
}

.top-badge {
  padding: 2px 8px;
  background: var(--color-warning);
  color: #fff;
  font-size: 11px;
  border-radius: 4px;
  font-weight: 500;
}

.draft-badge {
  padding: 2px 8px;
  background: var(--text-secondary);
  color: #fff;
  font-size: 11px;
  border-radius: 4px;
  font-weight: 500;
}

.status-badge {
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.published {
  background: color-mix(in srgb, var(--color-success) 10%, transparent);
  color: var(--color-success);
}

.status-badge.draft {
  background: color-mix(in srgb, var(--text-secondary) 10%, transparent);
  color: var(--text-secondary);
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
  background: var(--border);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:hover {
  background: var(--input-border);
  color: var(--text-primary);
}

.action-btn.edit:hover {
  background: color-mix(in srgb, var(--color-indigo) 10%, transparent);
  color: var(--color-indigo);
}

.action-btn.delete:hover {
  background: color-mix(in srgb, var(--color-danger) 10%, transparent);
  color: var(--color-danger);
}

.empty-cell {
  padding: 60px 16px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: var(--text-muted);
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-top: 24px;
  padding: 16px;
}

.page-btn {
  padding: 8px 16px;
  background: var(--bg-card);
  border: 1px solid var(--input-border);
  border-radius: 6px;
  font-size: 14px;
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.2s;
}

.page-btn:hover:not(:disabled) {
  background: var(--toolbar-bg);
  border-color: #d1d5db;
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  font-size: 14px;
  color: var(--text-secondary);
}
</style>