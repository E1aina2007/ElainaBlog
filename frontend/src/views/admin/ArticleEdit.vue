<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { createArticle, updateArticle } from '@/api/article'
import ArticleEditor from '@/components/ArticleEditor.vue'
import type { ArticleSubmitData } from '@/components/ArticleEditor.vue'
import toast from '@/utils/toast'

const route = useRoute()
const router = useRouter()
const editorRef = ref<InstanceType<typeof ArticleEditor>>()

const handleSubmit = async (data: ArticleSubmitData) => {
  if (editorRef.value) editorRef.value.setSaving(true)
  try {
    const payload = {
      ...data,
      category_id: data.category_id || undefined,
    }

    const id = route.params.id
    if (id) {
      await updateArticle({ ...payload, id: parseInt(id as string, 10) })
      toast.success(data.is_draft ? '草稿已保存' : '文章已更新并发布')
    } else {
      await createArticle(payload as any)
      toast.success(data.is_draft ? '草稿已保存' : '文章已创建并发布')
    }
    editorRef.value?.markSaved()
    setTimeout(() => router.push('/admin/articles'), 1000)
  } catch {
    toast.error('保存失败')
  } finally {
    if (editorRef.value) editorRef.value.setSaving(false)
  }
}
</script>

<template>
  <div class="article-edit">
    <ArticleEditor
      ref="editorRef"
      :show-top-option="true"
      :admin-mode="true"
      @submit="handleSubmit"
      @cancel="router.push('/admin/articles')"
    />
  </div>
</template>

<style scoped>
.article-edit {
  max-width: 1000px;
}
</style>
