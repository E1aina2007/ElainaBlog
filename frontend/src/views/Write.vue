<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { createArticle, updateArticle, getArticleDetail } from '@/api/article'
import { getCategoryList, type Category } from '@/api/category'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import toast from '@/utils/toast'

const md = new MarkdownIt({
  html: false,
  linkify: true,
  typographer: true,
  highlight(str: string, lang: string) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        const result = hljs.highlight(str, { language: lang, ignoreIllegals: true })
        return `<pre class="hljs-code-block"><div class="code-lang">${lang}</div><code class="hljs language-${lang}">${result.value}</code></pre>`
      } catch (_) {}
    }
    const escaped = md.utils.escapeHtml(str)
    return `<pre class="hljs-code-block"><code class="hljs">${escaped}</code></pre>`
  },
})

const renderedPreview = computed(() => md.render(form.value.content || ''))

const route = useRoute()
const router = useRouter()

const isEdit = ref(false)
const articleId = ref<number | null>(null)
const categories = ref<Category[]>([])
const saving = ref(false)
const preview = ref(false)

const form = ref({
  title: '',
  summary: '',
  content: '',
  category_id: null as number | null,
})

const fetchCategories = async () => {
  try {
    categories.value = await getCategoryList()
  } catch (e) {
    console.error('获取分类失败:', e)
  }
}

const fetchArticle = async (id: number) => {
  try {
    const article = await getArticleDetail(id)
    form.value = {
      title: article.title,
      summary: article.summary || '',
      content: article.content,
      category_id: article.category_id || null,
    }
  } catch (e) {
    toast.error('获取文章失败')
    router.push('/')
  }
}

const insertMarkdown = (syntax: string) => {
  const textarea = document.querySelector('.content-editor') as HTMLTextAreaElement
  if (!textarea) return
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const text = form.value.content
  form.value.content = text.substring(0, start) + syntax + text.substring(end)
  textarea.focus()
}

const handleSubmit = async (publish: boolean) => {
  if (!form.value.title.trim() || !form.value.content.trim()) {
    toast.warning('标题和内容不能为空')
    return
  }
  if (publish && !form.value.category_id) {
    toast.warning('发布文章请选择分类')
    return
  }

  saving.value = true
  try {
    const data = {
      ...form.value,
      is_draft: !publish,
      category_id: form.value.category_id || undefined,
    }

    if (isEdit.value && articleId.value) {
      await updateArticle({ ...data, id: articleId.value })
      toast.success(publish ? '文章已更新' : '草稿已保存')
    } else {
      await createArticle(data as any)
      toast.success(publish ? '文章已发布' : '草稿已保存')
    }
    setTimeout(() => {
      router.push('/')
    }, 1000)
  } catch (e) {
    toast.error('保存失败')
  } finally {
    saving.value = false
  }
}

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
  <main class="write-page">
    <div class="write-container">
      <!-- 顶部操作栏 -->
      <div class="write-header">
        <h1 class="write-title">{{ isEdit ? '编辑文章' : '写文章' }}</h1>
        <div class="header-actions">
          <button class="btn-outline" @click="router.push('/')">取消</button>
          <button class="btn-outline" @click="preview = !preview">
            {{ preview ? '编辑' : '预览' }}
          </button>
          <button class="btn-secondary" :disabled="saving" @click="handleSubmit(false)">
            保存草稿
          </button>
          <button class="btn-primary" :disabled="saving" @click="handleSubmit(true)">
            {{ saving ? '发布中...' : '发布文章' }}
          </button>
        </div>
      </div>

      <!-- 编辑表单 -->
      <div class="write-body">
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
        ></textarea>

        <!-- Markdown 工具栏 -->
        <div class="editor-toolbar">
          <button @click="insertMarkdown('**粗体**')" title="粗体"><strong>B</strong></button>
          <button @click="insertMarkdown('*斜体*')" title="斜体"><em>I</em></button>
          <button @click="insertMarkdown('# 标题')" title="标题">H</button>
          <button @click="insertMarkdown('[链接文字](url)')" title="链接">🔗</button>
          <button @click="insertMarkdown('![图片描述](url)')" title="图片">🖼</button>
          <button @click="insertMarkdown('```\n代码\n```')" title="代码块">&lt;/&gt;</button>
          <button @click="insertMarkdown('- 列表项')" title="列表">•</button>
          <button @click="insertMarkdown('> 引用内容')" title="引用">"</button>
          <button @click="insertMarkdown('---')" title="分割线">—</button>
        </div>

        <!-- 内容编辑 / 预览 -->
        <div v-if="!preview" class="editor-wrapper">
          <textarea
            v-model="form.content"
            class="content-editor"
            placeholder="请输入文章内容（支持 Markdown 语法）..."
          ></textarea>
        </div>
        <div v-else class="preview-wrapper">
          <div v-if="form.content" class="preview-content markdown-body" v-html="renderedPreview"></div>
          <div v-else class="preview-empty">暂无内容</div>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped>
.write-page {
  min-height: 100vh;
  background: var(--bg-primary);
  padding: 80px 24px 40px;
}

.write-container {
  max-width: 900px;
  margin: 0 auto;
}

/* 顶部操作栏 */
.write-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.write-title {
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
.write-body {
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

/* 分类与封面 */
.meta-row {
  display: flex;
  gap: 12px;
}

.meta-select,
.meta-input {
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

.meta-select:focus,
.meta-input:focus {
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

/* Markdown 工具栏 */
.editor-toolbar {
  display: flex;
  gap: 4px;
  padding: 8px 10px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-bottom: none;
  border-radius: var(--radius-md) var(--radius-md) 0 0;
}

.editor-toolbar button {
  padding: 6px 12px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 4px;
  font-size: 0.8125rem;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.15s;
}

.editor-toolbar button:hover {
  background: var(--primary-lighter);
  color: var(--primary-dark);
  border-color: var(--primary-light);
}

/* 编辑器 */
.editor-wrapper {
  flex: 1;
}

.content-editor {
  width: 100%;
  min-height: 500px;
  padding: 16px;
  font-size: 0.9375rem;
  line-height: 1.8;
  color: var(--text-primary);
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-top: none;
  border-radius: 0 0 var(--radius-md) var(--radius-md);
  outline: none;
  resize: vertical;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  tab-size: 2;
  transition: border-color 0.2s;
}

.content-editor:focus {
  border-color: var(--primary);
}

.content-editor::placeholder {
  color: var(--text-muted);
}

/* 预览 */
.preview-wrapper {
  min-height: 500px;
  padding: 24px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
}

.preview-content {
  font-size: 0.9375rem;
  line-height: 1.8;
  color: var(--text-primary);
  white-space: pre-wrap;
  word-break: break-word;
}

.preview-empty {
  text-align: center;
  padding: 80px 0;
  color: var(--text-muted);
  font-size: 0.875rem;
}

/* 响应式 */
@media (max-width: 768px) {
  .write-page {
    padding: 60px 16px 24px;
  }

  .write-header {
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

  .write-body {
    padding: 20px;
  }

  .title-input {
    font-size: 1.25rem;
  }

  .meta-row {
    flex-direction: column;
  }

  .content-editor {
    min-height: 350px;
  }
}
</style>
