<script setup lang="ts">
import { ref, nextTick, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, onBeforeRouteLeave } from 'vue-router'
import { getArticleDetail, getAdminArticleDetail, getMyArticleDetail } from '@/api/article'
import { getCategoryList, type Category } from '@/api/category'
import { uploadImage } from '@/api/upload'
import ByteMDEditor from '@/components/ByteMDEditor.vue'
import toast from '@/utils/toast'

export interface ArticleSubmitData {
  title: string
  summary: string
  content: string
  tags: string
  category_id: number | null
  is_top: boolean
  is_draft: boolean
}

const props = withDefaults(defineProps<{
  showTopOption?: boolean
  adminMode?: boolean
  userMode?: boolean
}>(), {
  showTopOption: false,
  adminMode: false,
  userMode: false,
})

const emit = defineEmits<{
  submit: [data: ArticleSubmitData]
  cancel: []
}>()

const route = useRoute()

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

// 关键词标签（内部用数组，提交时转 | 分隔字符串）
const tags = ref<string[]>([])
const tagInput = ref('')
const showTagInput = ref(false)
const tagInputRef = ref<HTMLInputElement | null>(null)

// 离开确认相关状态
const isDirty = ref(false)
const isSubmitting = ref(false)
let initialFormSnapshot = ''

// 保存当前表单快照
const saveSnapshot = () => {
  initialFormSnapshot = JSON.stringify({ ...form.value, tags: tags.value })
}

// 检查表单是否有变化
const checkDirty = () => {
  const currentSnapshot = JSON.stringify({ ...form.value, tags: tags.value })
  isDirty.value = currentSnapshot !== initialFormSnapshot
}

// 监听表单和标签变化
watch(form, checkDirty, { deep: true })
watch(tags, checkDirty, { deep: true })

// 路由离开拦截
onBeforeRouteLeave((_to, _from, next) => {
  if (isSubmitting.value || !isDirty.value) {
    next()
    return
  }
  const answer = window.confirm('文章内容尚未保存，确定要离开吗？')
  if (answer) {
    next()
  } else {
    next(false)
  }
})

// 浏览器关闭/刷新拦截
const handleBeforeUnload = (e: BeforeUnloadEvent) => {
  if (isDirty.value) {
    e.preventDefault()
    e.returnValue = ''
  }
}

onMounted(() => {
  window.addEventListener('beforeunload', handleBeforeUnload)
})

onUnmounted(() => {
  window.removeEventListener('beforeunload', handleBeforeUnload)
})

// 标记为已保存（父组件调用）
const markSaved = () => {
  isDirty.value = false
  isSubmitting.value = true
}

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
    let article
    if (props.adminMode) {
      article = await getAdminArticleDetail(id)
    } else if (props.userMode) {
      article = await getMyArticleDetail(id)
    } else {
      article = await getArticleDetail(id)
    }
    form.value = {
      title: article.title,
      summary: article.summary || '',
      content: article.content,
      category_id: article.category_id || null,
      is_top: article.is_top || false,
    }
    // 解析关键词标签
    if (article.tags) {
      tags.value = article.tags.split('|').map(t => t.trim()).filter(Boolean)
    }
    saveSnapshot()
  } catch {
    toast.error('获取文章失败')
  }
}

// 点击 + 按钮展开输入框
const openTagInput = () => {
  showTagInput.value = true
  nextTick(() => {
    tagInputRef.value?.focus()
  })
}

// 确认添加关键词标签
const confirmTag = () => {
  const val = tagInput.value.trim()
  if (val && !tags.value.includes(val)) {
    tags.value.push(val)
  }
  tagInput.value = ''
  showTagInput.value = false
}

// 取消添加
const cancelTag = () => {
  tagInput.value = ''
  showTagInput.value = false
}

// 删除关键词标签
const removeTag = (index: number) => {
  tags.value.splice(index, 1)
}

// 标签输入框回车确认
const handleTagKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Enter') {
    e.preventDefault()
    confirmTag()
  } else if (e.key === 'Escape') {
    cancelTag()
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
    tags: tags.value.join('|'),
    is_draft: !publish,
    category_id: form.value.category_id || null,
  })
}

// 取消操作
const handleCancel = () => {
  isDirty.value = false
  emit('cancel')
}

// 编辑器图片上传（ByteMD: (files: File[]) => Promise<{url: string}[]>）
const handleImageUpload = async (files: File[]): Promise<{ url: string }[]> => {
  const results = await Promise.allSettled(
    files.map(async (file) => {
      const res = await uploadImage(file)
      return { url: res.url }
    })
  )
  results.forEach((r, i) => {
    if (r.status === 'rejected') {
      toast.error(`图片 ${files[i]?.name ?? '未知'} 上传失败`)
    }
  })
  return results
    .filter((r): r is PromiseFulfilledResult<{ url: string }> => r.status === 'fulfilled')
    .map((r) => r.value)
}

// 暴露给父组件
const setSaving = (v: boolean) => { saving.value = v }
defineExpose({ setSaving, markSaved })

onMounted(() => {
  fetchCategories()
  const id = route.params.id
  if (id) {
    isEdit.value = true
    articleId.value = parseInt(id as string, 10)
    fetchArticle(articleId.value)
  } else {
    // 新建文章，保存初始快照
    saveSnapshot()
  }
})
</script>

<template>
  <div class="article-editor">
    <!-- 顶部操作栏 -->
    <div class="editor-header">
      <h2 class="editor-title">{{ isEdit ? '编辑文章' : '写文章' }}</h2>
      <div class="header-actions">
        <button class="btn-outline" @click="handleCancel">取消</button>
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

      <!-- 关键词标签 -->
      <div class="tags-input-row">
        <div class="tags-list">
          <div v-for="(tag, index) in tags" :key="index" class="tag-card">
            <span class="tag-card-text">{{ tag }}</span>
            <button type="button" class="tag-card-remove" @click="removeTag(index)">×</button>
          </div>
          <button v-if="!showTagInput" type="button" class="tag-add-btn" @click="openTagInput" title="添加关键词">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
              <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
            <span>添加关键词</span>
          </button>
        </div>
        <div v-if="showTagInput" class="tag-input-wrapper">
          <input
            ref="tagInputRef"
            v-model="tagInput"
            type="text"
            class="tag-input"
            placeholder="输入关键词"
            @keydown="handleTagKeydown"
          />
          <button type="button" class="tag-confirm-btn" @click="confirmTag">确定</button>
          <button type="button" class="tag-cancel-btn" @click="cancelTag">取消</button>
        </div>
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

      <!-- Markdown 编辑器（ByteMD 分屏预览） -->
      <ByteMDEditor
        v-model="form.content"
        :upload-images="handleImageUpload"
        editor-height="500px"
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

/* 关键词标签输入 */
.tags-input-row {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-card {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: var(--bg-card);
  color: var(--primary-dark);
  font-size: 13px;
  font-weight: 500;
  border: 1px solid var(--primary-light);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-soft);
  transition: all 0.2s;
}

.tag-card:hover {
  border-color: var(--primary);
  box-shadow: 0 2px 8px rgba(126, 215, 193, 0.25);
}

.tag-card-text {
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tag-card-remove {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  border: none;
  background: var(--bg-secondary);
  color: var(--text-muted);
  font-size: 14px;
  cursor: pointer;
  border-radius: 50%;
  transition: all 0.15s;
  flex-shrink: 0;
}

.tag-card-remove:hover {
  background: color-mix(in srgb, var(--color-danger) 15%, transparent);
  color: var(--color-danger);
}

.tag-input-wrapper {
  display: flex;
  gap: 8px;
}

.tag-input {
  flex: 1;
  padding: 8px 12px;
  font-size: 0.875rem;
  color: var(--text-primary);
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  outline: none;
  transition: border-color 0.2s;
}

.tag-input:focus {
  border-color: var(--primary);
}

.tag-input::placeholder {
  color: var(--text-muted);
}

.tag-add-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 14px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  flex-shrink: 0;
}

.tag-add-btn:hover {
  background: var(--primary-dark);
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(126, 215, 193, 0.4);
}

.tag-confirm-btn {
  padding: 8px 14px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
  white-space: nowrap;
}

.tag-confirm-btn:hover {
  background: var(--primary-dark);
}

.tag-cancel-btn {
  padding: 8px 14px;
  background: transparent;
  color: var(--text-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.tag-cancel-btn:hover {
  color: var(--text-primary);
  border-color: var(--primary-light);
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
.editor-body :deep(.bytemd-editor-wrapper) {
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
