<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { createArticle, updateArticle } from '@/api/article'
import ArticleEditor from '@/components/ArticleEditor.vue'
import type { ArticleSubmitData } from '@/components/ArticleEditor.vue'
import { useUserStore } from '@/stores/user'
import toast from '@/utils/toast'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const editorRef = ref<InstanceType<typeof ArticleEditor>>()
const showTopOption = computed(() => userStore.isAdmin)

const handleSubmit = async (data: ArticleSubmitData) => {
  if (editorRef.value) editorRef.value.setSaving(true)
  try {
    const payload = {
      ...data,
      // 非管理员不允许置顶
      is_top: userStore.isAdmin ? data.is_top : false,
      category_id: data.category_id || undefined,
    }

    const id = route.params.id
    if (id) {
      await updateArticle({ ...payload, id: parseInt(id as string, 10) })
      toast.success(data.is_draft ? '草稿已保存' : '文章已更新')
    } else {
      await createArticle(payload as any)
      toast.success(data.is_draft ? '草稿已保存' : '文章已发布')
    }
    editorRef.value?.markSaved()
    setTimeout(() => router.push('/'), 1000)
  } catch {
    toast.error('保存失败')
  } finally {
    if (editorRef.value) editorRef.value.setSaving(false)
  }
}
</script>

<template>
  <main class="write-page">
    <div class="write-container">
      <ArticleEditor
        ref="editorRef"
        :show-top-option="showTopOption"
        :user-mode="!!route.params.id"
        @submit="handleSubmit"
        @cancel="router.push('/')"
      />
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

@media (max-width: 768px) {
  .write-page {
    padding: 60px 16px 24px;
  }
}
</style>
