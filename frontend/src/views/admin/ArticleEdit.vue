<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { createArticle, updateArticle, getArticleDetail, type Article } from '@/api/article'
import { getCategoryList, type Category } from '@/api/category'
import toast from '@/utils/toast'

const route = useRoute()
const router = useRouter()

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
  is_draft: false,
})

const fetchArticle = async (id: number) => {
  try {
    const article = await getArticleDetail(id)
    form.value = {
      title: article.title,
      summary: article.summary || '',
      content: article.content,
      category_id: article.category_id || null,
      is_top: article.is_top || false,
      is_draft: article.is_draft || false,
    }
  } catch (error) {
    console.error('获取文章失败:', error)
    toast.error('获取文章失败')
  }
}

const fetchCategories = async () => {
  try {
    categories.value = await getCategoryList()
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

const handleSubmit = async (publish: boolean = true) => {
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
    const data: any = {
      ...form.value,
      is_draft: !publish,
      category_id: form.value.category_id || undefined,
    }

    if (isEdit.value && articleId.value) {
      await updateArticle({ ...data, id: articleId.value })
      toast.success(publish ? '文章已更新并发布' : '草稿已保存')
    } else {
      await createArticle(data as any)
      toast.success(publish ? '文章已创建并发布' : '草稿已保存')
    }
    setTimeout(() => {
      router.push('/admin/articles')
    }, 1000)
  } catch (error) {
    console.error('保存失败:', error)
    toast.error('保存失败')
  } finally {
    saving.value = false
  }
}

const handleCancel = () => {
  router.push('/admin/articles')
}

const insertMarkdown = (syntax: string) => {
  const textarea = document.querySelector('textarea[name="content"]') as HTMLTextAreaElement
  if (!textarea) return

  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const text = form.value.content

  form.value.content = text.substring(0, start) + syntax + text.substring(end)
  textarea.focus()
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
  <div class="article-edit">
    <div class="page-header">
      <h2>{{ isEdit ? '编辑文章' : '写文章' }}</h2>
      <div class="header-actions">
        <button class="btn-secondary" @click="handleCancel">
          取消
        </button>
        <button class="btn-secondary" @click="handleSubmit(false)" :disabled="saving">
          保存草稿
        </button>
        <button class="btn-primary" @click="handleSubmit(true)" :disabled="saving">
          {{ saving ? '保存中...' : (isEdit ? '更新发布' : '立即发布') }}
        </button>
      </div>
    </div>

    <div class="edit-form">
      <div class="form-group">
        <label>文章标题</label>
        <input
          v-model="form.title"
          type="text"
          placeholder="请输入文章标题..."
          class="title-input"
        />
      </div>

      <div class="form-group">
        <label>所属分类</label>
        <select v-model="form.category_id" class="form-select">
          <option :value="null">未分类</option>
          <option v-for="cat in categories" :key="cat.id" :value="cat.id">
            {{ cat.name }}
          </option>
        </select>
      </div>

      <div class="form-group">
        <label>文章摘要</label>
        <textarea
          v-model="form.summary"
          rows="3"
          placeholder="请输入文章摘要（可选）..."
          class="summary-input"
        ></textarea>
      </div>

      <div class="form-group">
        <label>文章内容</label>
        <div class="editor-toolbar">
          <button @click="insertMarkdown('**粗体**')" title="粗体">B</button>
          <button @click="insertMarkdown('*斜体*')" title="斜体">I</button>
          <button @click="insertMarkdown('# 标题')" title="标题">H</button>
          <button @click="insertMarkdown('[链接](url)')" title="链接">🔗</button>
          <button @click="insertMarkdown('```\n代码\n```')" title="代码块">&lt;/&gt;</button>
          <button @click="insertMarkdown('- 列表项')" title="列表">•</button>
          <button @click="insertMarkdown('> 引用')" title="引用">"</button>
        </div>
        <textarea
          v-model="form.content"
          name="content"
          rows="20"
          placeholder="请输入文章内容（支持 Markdown 语法）..."
          class="content-input"
        ></textarea>
      </div>

      <div class="form-group">
        <label class="checkbox-label">
          <input type="checkbox" v-model="form.is_top" />
          <span>置顶文章</span>
        </label>
      </div>
    </div>
  </div>
</template>

<style scoped>
.article-edit {
  max-width: 1000px;
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

.header-actions {
  display: flex;
  gap: 12px;
}

.btn-primary {
  padding: 10px 24px;
  background: #6366f1;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover:not(:disabled) {
  background: #4f46e5;
}

.btn-secondary {
  padding: 10px 20px;
  background: #fff;
  color: #374151;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary:hover:not(:disabled) {
  background: #f9fafb;
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.edit-form {
  background: #fff;
  border-radius: 12px;
  padding: 32px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.form-group {
  margin-bottom: 24px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.title-input {
  width: 100%;
  padding: 12px 16px;
  font-size: 18px;
  font-weight: 600;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  outline: none;
  transition: border-color 0.2s;
}

.title-input:focus {
  border-color: #6366f1;
}

.form-select,
.form-input {
  width: 100%;
  padding: 10px 14px;
  font-size: 14px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  outline: none;
  background: #fff;
  transition: border-color 0.2s;
}

.form-select:focus,
.form-input:focus {
  border-color: #6366f1;
}

.summary-input,
.content-input {
  width: 100%;
  padding: 12px 16px;
  font-size: 14px;
  line-height: 1.6;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  outline: none;
  resize: vertical;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  transition: border-color 0.2s;
}

.summary-input:focus,
.content-input:focus {
  border-color: #6366f1;
}

.editor-toolbar {
  display: flex;
  gap: 4px;
  padding: 8px;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-bottom: none;
  border-radius: 8px 8px 0 0;
}

.editor-toolbar button {
  padding: 6px 12px;
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 4px;
  font-size: 13px;
  font-weight: 600;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s;
}

.editor-toolbar button:hover {
  background: #f3f4f6;
}

.content-input {
  border-radius: 0 0 8px 8px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  accent-color: #6366f1;
}

.checkbox-label span {
  font-size: 14px;
  color: #374151;
}

@media (max-width: 768px) {
  .form-row {
    grid-template-columns: 1fr;
  }

  .header-actions {
    flex-wrap: wrap;
  }
}
</style>