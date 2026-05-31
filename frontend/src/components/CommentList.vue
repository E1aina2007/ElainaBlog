<script setup lang="ts">
import { useUserStore } from '@/stores/user'
import type { Comment } from '@/api/comment'

interface Props {
  comments: Comment[]
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
})

const emit = defineEmits<{
  delete: [id: number]
}>()

const userStore = useUserStore()

function formatDate(date: string): string {
  return new Date(date).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function canDelete(comment: Comment): boolean {
  // 本人或管理员可删除
  return userStore.userInfo?.id === comment.user_id || userStore.isAdmin
}

function handleDelete(id: number) {
  if (confirm('确定要删除这条评论吗？')) {
    emit('delete', id)
  }
}
</script>

<template>
  <div class="comment-list">
    <h3 class="list-title">
      评论区
      <span v-if="comments.length > 0" class="comment-count">({{ comments.length }})</span>
    </h3>

    <div v-if="loading" class="loading-state">
      <span>加载中...</span>
    </div>

    <div v-else-if="comments.length === 0" class="empty-state">
      <p>暂无评论，来说两句吧~</p>
    </div>

    <div v-else class="comments">
      <div
        v-for="comment in comments"
        :key="comment.id"
        class="comment-item"
      >
        <div class="comment-header">
          <div class="user-meta">
            <img v-if="comment.avatar" :src="comment.avatar" class="user-avatar-img" alt="头像" />
            <span v-else class="user-avatar">{{ comment.username.charAt(0).toUpperCase() }}</span>
            <span class="username">{{ comment.username }}</span>
            <span v-if="comment.is_admin" class="admin-badge">管理员</span>
          </div>
          <time class="comment-time">{{ formatDate(comment.created_at) }}</time>
        </div>
        <p class="comment-content">{{ comment.content }}</p>
        <div v-if="canDelete(comment)" class="comment-actions">
          <button class="delete-btn" @click="handleDelete(comment.id)">
            删除
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.comment-list {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 24px;
  box-shadow: var(--shadow-soft);
}

.list-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 20px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.comment-count {
  font-size: 0.875rem;
  color: var(--text-muted);
  font-weight: 400;
}

.loading-state,
.empty-state {
  padding: 40px 24px;
  text-align: center;
  color: var(--text-muted);
}

.empty-state p {
  margin: 0;
  font-size: 0.9375rem;
}

.comments {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.comment-item {
  padding: 16px;
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
}

.comment-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.user-meta {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
}

.user-avatar-img {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
}

.username {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-primary);
}

.admin-badge {
  display: inline-flex;
  align-items: center;
  padding: 1px 8px;
  background: linear-gradient(135deg, var(--color-warning) 0%, var(--color-warning-dark) 100%);
  color: white;
  font-size: 11px;
  font-weight: 600;
  border-radius: 10px;
  line-height: 1.6;
}

.comment-time {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.comment-content {
  font-size: 0.9375rem;
  color: var(--text-secondary);
  line-height: 1.7;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}

.comment-actions {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
}

.delete-btn {
  font-size: 0.75rem;
  color: var(--color-danger);
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  transition: background 0.15s ease;
}

.delete-btn:hover {
  background: rgba(231, 76, 60, 0.1);
}
</style>