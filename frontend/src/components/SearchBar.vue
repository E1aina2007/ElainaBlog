<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { searchArticles, type Article } from '@/api/article'

const router = useRouter()

const keyword = ref('')
const results = ref<Article[]>([])
const showDropdown = ref(false)
const isLoading = ref(false)
const searchInput = ref<HTMLInputElement | null>(null)

let debounceTimer: ReturnType<typeof setTimeout> | null = null

const emit = defineEmits<{
  search: [keyword: string]
}>()

// 防抖搜索
watch(keyword, (val) => {
  if (debounceTimer) clearTimeout(debounceTimer)
  if (!val.trim()) {
    results.value = []
    showDropdown.value = false
    return
  }
  debounceTimer = setTimeout(() => {
    doSearch(val.trim())
  }, 300)
})

const doSearch = async (q: string) => {
  if (!q) return
  isLoading.value = true
  try {
    const res = await searchArticles(q, 1, 6)
    results.value = res.list || []
    showDropdown.value = results.value.length > 0
  } catch {
    results.value = []
  } finally {
    isLoading.value = false
  }
}

// 点击搜索结果
const goToArticle = (id: number) => {
  showDropdown.value = false
  keyword.value = ''
  router.push(`/article/${id}`)
}

// 回车触发搜索
const handleEnter = () => {
  const q = keyword.value.trim()
  if (!q) return
  showDropdown.value = false
  emit('search', q)
}

// 点击外部关闭下拉
const handleClickOutside = (e: MouseEvent) => {
  const target = e.target as HTMLElement
  if (!target.closest('.search-bar')) {
    showDropdown.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  if (debounceTimer) clearTimeout(debounceTimer)
})

// 截断文本
const truncate = (text: string, len: number) => {
  if (!text) return ''
  return text.length > len ? text.slice(0, len) + '...' : text
}
</script>

<template>
  <div class="search-bar">
    <div class="search-input-wrapper">
      <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="11" cy="11" r="8"/>
        <line x1="21" y1="21" x2="16.65" y2="16.65"/>
      </svg>
      <input
        ref="searchInput"
        v-model="keyword"
        type="text"
        class="search-input"
        placeholder="搜索文章..."
        @keydown.enter="handleEnter"
        @focus="keyword.trim() && results.length > 0 && (showDropdown = true)"
      />
      <span v-if="isLoading" class="search-loading"></span>
    </div>

    <!-- 搜索结果下拉 -->
    <div v-if="showDropdown" class="search-dropdown">
      <div
        v-for="article in results"
        :key="article.id"
        class="search-result-item"
        @click="goToArticle(article.id)"
      >
        <div class="result-title">{{ article.title }}</div>
        <div class="result-summary">{{ truncate(article.summary || article.content, 80) }}</div>
      </div>
      <div class="search-dropdown-footer" @click="handleEnter">
        查看全部搜索结果 →
      </div>
    </div>
  </div>
</template>

<style scoped>
.search-bar {
  position: relative;
  width: 100%;
  max-width: 400px;
}

.search-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 12px;
  color: var(--text-muted);
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 10px 36px 10px 36px;
  font-size: 14px;
  color: var(--text-primary);
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-full);
  outline: none;
  transition: all 0.2s;
}

.search-input:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(126, 215, 193, 0.15);
}

.search-input::placeholder {
  color: var(--text-muted);
}

.search-loading {
  position: absolute;
  right: 12px;
  width: 16px;
  height: 16px;
  border: 2px solid var(--border);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 下拉结果 */
.search-dropdown {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  right: 0;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg, 0 8px 24px rgba(0,0,0,0.12));
  z-index: 100;
  overflow: hidden;
}

.search-result-item {
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.15s;
  border-bottom: 1px solid var(--border);
}

.search-result-item:last-of-type {
  border-bottom: none;
}

.search-result-item:hover {
  background: var(--bg-secondary);
}

.result-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.result-summary {
  font-size: 12px;
  color: var(--text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.search-dropdown-footer {
  padding: 10px 16px;
  text-align: center;
  font-size: 13px;
  color: var(--primary);
  cursor: pointer;
  border-top: 1px solid var(--border);
  transition: background 0.15s;
}

.search-dropdown-footer:hover {
  background: var(--bg-secondary);
}

@media (max-width: 768px) {
  .search-bar {
    max-width: 100%;
  }
}
</style>
