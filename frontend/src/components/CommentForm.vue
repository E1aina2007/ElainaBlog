<script setup lang="ts">
import { ref } from 'vue'
import { useUserStore } from '@/stores/user'

interface Props {
  articleId: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  submit: [content: string]
}>()

const userStore = useUserStore()
const content = ref('')
const isSubmitting = ref(false)

const isLoggedIn = () => userStore.isLoggedIn

async function handleSubmit() {
  if (!content.value.trim()) return

  isSubmitting.value = true
  try {
    emit('submit', content.value.trim())
    content.value = ''
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="comment-form">
    <h3 class="form-title">发表评论</h3>

    <!-- 未登录提示 -->
    <div v-if="!isLoggedIn()" class="login-tip">
      <p>请先<router-link to="/login" class="login-link">登录</router-link>后发表评论</p>
    </div>

    <!-- 评论表单 -->
    <template v-else>
      <div class="user-info">
        <span class="username">{{ userStore.userInfo?.username }}</span>
      </div>
      <textarea
        v-model="content"
        class="comment-input"
        placeholder="写下你的评论..."
        rows="4"
        :disabled="isSubmitting"
      />
      <div class="form-actions">
        <button
          class="submit-btn"
          :disabled="!content.trim() || isSubmitting"
          @click="handleSubmit"
        >
          <span v-if="isSubmitting">提交中...</span>
          <span v-else>发表评论</span>
        </button>
      </div>
    </template>
  </div>
</template>

<style scoped>
.comment-form {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: var(--shadow-soft);
}

.form-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 16px;
}

.login-tip {
  padding: 24px;
  text-align: center;
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
}

.login-link {
  color: var(--primary);
  text-decoration: none;
  font-weight: 500;
  margin: 0 4px;
}

.login-link:hover {
  text-decoration: underline;
}

.user-info {
  margin-bottom: 12px;
}

.username {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
}

.comment-input {
  width: 100%;
  padding: 12px 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  font-size: 0.9375rem;
  color: var(--text-primary);
  line-height: 1.7;
  resize: vertical;
  outline: none;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.comment-input::placeholder {
  color: var(--text-muted);
}

.comment-input:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(126, 215, 193, 0.15);
}

.comment-input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.form-actions {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
}

.submit-btn {
  padding: 10px 24px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease-out;
}

.submit-btn:hover:not(:disabled) {
  background: var(--primary-dark);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(126, 215, 193, 0.3);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>