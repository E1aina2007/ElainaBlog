<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getArticleDetail } from '@/api/article'
import { getCategoryList, type Category } from '@/api/category'
import { useTheme } from '@/composables/useTheme'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import toast from '@/utils/toast'

export interface ArticleSubmitData {
  title: string
  summary: string
  content: string
  category_id: number | null
  is_top: boolean
  is_draft: boolean
}

const props = withDefaults(defineProps<{
  showTopOption?: boolean
}>(), {
  showTopOption: false,
})

const emit = defineEmits<{
  submit: [data: ArticleSubmitData]
  cancel: []
}>()

const route = useRoute()
const { isDark } = useTheme()

const isEdit = ref(false)
const articleId = ref<number | null>(null)
const categories = ref<Category[]>([])
const saving = ref(false)

const form = ref({
  title: '',
  summary: '',
  content: '',
  category_id: null as number | null,
  is_top: false,
})

const editorTheme = computed(() => isDark.value ? 'dark' : 'light')

// 加载分类
const fetchCategories = async () => {
  try {
    categories.value = await getCategoryList()
  } catch (e) {
    console.error('获取分类失败:', e)
  }
}

// 加载文章（编辑模式）
const fetchArticle = async (id: number) => {
  try {
    const article = await getArticleDetail(id)
    form.value = {
      title: article.title,
      summary: article.summary || '',
      content: article.content,
      category_id: article.category_id || null,
      is_top: article.is_top || false,
    }
  } catch {
    toast.error('获取文章失败')
  }
}

// 提交
const handleSubmit = (publish: boolean) => {
  if (!form.value.title.trim() || !form.value.content.trim()) {
    toast.warning('标题和内容不能为空')
    return
  }
  if (publish && !form.value.category_id) {
    toast.warning('发布文章请选择分类')
    return
  }
  emit('submit', {
    ...form.value,
    is_draft: !publish,
    category_id: form.value.category_id || null,
  })
}

// 暴露给父组件：设置 saving 状态
const setSaving = (v: boolean) => { saving.value = v }
defineExpose({ setSaving })

onMounted(() => {
  fetchCategories()
  const id = route.params.id
  if (id) {
    isEdit.value = true
    articleId.value = parseInt(id as string, 10)
    fetchArticle(articleId.value)
  }
})
</script>

<template>
  <div class="article-editor">
    <!-- 顶部操作栏 -->
    <div class="editor-header">
      <h2 class="editor-title">{{ isEdit ? '编辑文章' : '写文章' }}</h2>
      <div class="header-actions">
        <button class="btn-outline" @click="emit('cancel')">取消</button>
        <button class="btn-secondary" :disabled="saving" @click="handleSubmit(false)">
          保存草稿
        </button>
        <button class="btn-primary" :disabled="saving" @click="handleSubmit(true)">
          {{ saving ? '保存中...' : (isEdit ? '更新发布' : '发布文章') }}
        </button>
      </div>
    </div>

    <!-- 编辑表单 -->
    <div class="editor-body">
      <!-- 标题 -->
      <input
        v-model="form.title"
        type="text"
        class="title-input"
        placeholder="请输入文章标题..."
      />

      <!-- 分类 -->
      <div class="meta-row">
        <select v-model="form.category_id" class="meta-select">
          <option :value="null">选择分类</option>
          <option v-for="cat in categories" :key="cat.id" :value="cat.id">
            {{ cat.name }}
          </option>
        </select>
      </div>

      <!-- 摘要 -->
      <textarea
        v-model="form.summary"
        class="summary-input"
        rows="2"
        placeholder="文章摘要（可选，会显示在文章列表中）"
        v-tab-indent
      ></textarea>

      <!-- 置顶选项 -->
      <label v-if="showTopOption" class="checkbox-label">
        <input type="checkbox" v-model="form.is_top" />
        <span>置顶文章</span>
      </label>

      <!-- Markdown 编辑器（内置 CodeMirror 语法高亮 + 实时预览） -->
      <MdEditor
        v-model="form.content"
        :theme="editorTheme"
        :show-code-row-number="true"
        style="height: 500px"
        placeholder="请输入文章内容（支持 Markdown 语法）..."
      />
    </div>
  </div>
</template>

<style scoped>
.article-editor {
  width: 100%;
}

/* 顶部操作栏 */
.editor-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.editor-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.btn-primary {
  padding: 10px 24px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover:not(:disabled) {
  background: var(--primary-dark);
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(126, 215, 193, 0.4);
}

.btn-secondary {
  padding: 10px 20px;
  background: var(--bg-card);
  color: var(--text-primary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary:hover:not(:disabled) {
  background: var(--bg-secondary);
  border-color: var(--primary-light);
}

.btn-outline {
  padding: 10px 20px;
  background: transparent;
  color: var(--text-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-outline:hover {
  color: var(--text-primary);
  border-color: var(--primary-light);
  background: var(--bg-secondary);
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 编辑表单 */
.editor-body {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 32px;
  box-shadow: var(--shadow-soft);
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* 标题 */
.title-input {
  width: 100%;
  padding: 14px 0;
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--text-primary);
  background: transparent;
  border: none;
  border-bottom: 2px solid var(--border);
  outline: none;
  transition: border-color 0.2s;
}

.title-input:focus {
  border-color: var(--primary);
}

.title-input::placeholder {
  color: var(--text-muted);
  font-weight: 400;
}

/* 分类 */
.meta-row {
  display: flex;
  gap: 12px;
}

.meta-select {
  flex: 1;
  padding: 10px 14px;
  font-size: 0.875rem;
  color: var(--text-primary);
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  outline: none;
  transition: border-color 0.2s;
}

.meta-select:focus {
  border-color: var(--primary);
}

/* 摘要 */
.summary-input {
  width: 100%;
  padding: 12px 14px;
  font-size: 0.875rem;
  color: var(--text-primary);
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  outline: none;
  resize: vertical;
  line-height: 1.6;
  transition: border-color 0.2s;
}

.summary-input:focus {
  border-color: var(--primary);
}

.summary-input::placeholder {
  color: var(--text-muted);
}

/* 置顶选项 */
.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  accent-color: var(--primary);
}

.checkbox-label span {
  font-size: 0.875rem;
  color: var(--text-primary);
}

/* 编辑器圆角适配 */
.editor-body :deep(.md-editor) {
  border-radius: var(--radius-md);
  overflow: hidden;
}

/* 响应式 */
@media (max-width: 768px) {
  .editor-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .header-actions {
    flex-wrap: wrap;
    width: 100%;
  }

  .header-actions button {
    flex: 1;
    min-width: 0;
  }

  .editor-body {
    padding: 20px;
  }

  .title-input {
    font-size: 1.25rem;
  }

  .meta-row {
    flex-direction: column;
  }
}
</style>
