<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getUserList, deleteUser, type UserProfile } from '@/api/user'
import toast from '@/utils/toast'

const users = ref<UserProfile[]>([])
const loading = ref(false)
const searchQuery = ref('')

const fetchUsers = async () => {
  loading.value = true
  try {
    users.value = await getUserList()
  } catch {
    toast.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

const filteredUsers = computed(() => {
  if (!searchQuery.value) return users.value
  const q = searchQuery.value.toLowerCase()
  return users.value.filter(u =>
    u.username.toLowerCase().includes(q) || u.email.toLowerCase().includes(q)
  )
})

const handleDelete = async (user: UserProfile) => {
  if (!confirm(`确定要删除用户 "${user.username}"（${user.email}）？此操作不可撤销。`)) return
  try {
    await deleteUser(user.id)
    users.value = users.value.filter(u => u.id !== user.id)
    toast.success('用户已删除')
  } catch {
    toast.error('删除失败')
  }
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(fetchUsers)
</script>

<template>
  <div class="user-list-page">
    <div class="page-header">
      <h2>用户管理</h2>
      <span class="user-count">共 {{ filteredUsers.length }} 位用户</span>
    </div>

    <!-- 搜索 -->
    <div class="search-bar">
      <input
        v-model="searchQuery"
        class="search-input"
        placeholder="搜索用户名或邮箱..."
      />
    </div>

    <!-- 加载 -->
    <div v-if="loading" class="loading-state">加载中...</div>

    <!-- 用户列表 -->
    <div v-else class="user-table-wrapper">
      <table class="user-table">
        <thead>
          <tr>
            <th>用户</th>
            <th>邮箱</th>
            <th>角色</th>
            <th>注册时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in filteredUsers" :key="user.id">
            <td>
              <div class="user-cell">
                <img v-if="user.avatar" :src="user.avatar" class="user-avatar" alt="" />
                <span v-else class="user-avatar-placeholder">{{ user.username.charAt(0).toUpperCase() }}</span>
                <span class="user-name">{{ user.username }}</span>
              </div>
            </td>
            <td class="email-cell">{{ user.email }}</td>
            <td>
              <span class="role-badge" :class="{ admin: user.is_admin }">
                {{ user.is_admin ? '管理员' : '普通用户' }}
              </span>
            </td>
            <td class="date-cell">{{ formatDate(user.created_at) }}</td>
            <td>
              <button
                class="btn-delete"
                :disabled="user.is_admin"
                :title="user.is_admin ? '不能删除管理员' : '删除用户'"
                @click="handleDelete(user)"
              >
                删除
              </button>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-if="!loading && filteredUsers.length === 0" class="empty-state">
        {{ searchQuery ? '没有匹配的用户' : '暂无用户' }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.user-list-page {
  max-width: 1000px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.page-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.user-count {
  font-size: 14px;
  color: var(--text-muted);
}

.search-bar {
  margin-bottom: 20px;
}

.search-input {
  width: 100%;
  max-width: 400px;
  padding: 8px 14px;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  font-size: 14px;
  color: var(--text-primary);
  background: var(--bg-card);
  outline: none;
  transition: border-color 0.2s;
}

.search-input:focus {
  border-color: var(--primary);
}

.loading-state {
  text-align: center;
  padding: 40px;
  color: var(--text-muted);
}

.user-table-wrapper {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-soft);
  overflow: hidden;
}

.user-table {
  width: 100%;
  border-collapse: collapse;
}

.user-table th {
  text-align: left;
  padding: 12px 16px;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border);
}

.user-table td {
  padding: 12px 16px;
  font-size: 14px;
  color: var(--text-primary);
  border-bottom: 1px solid var(--border);
}

.user-table tr:last-child td {
  border-bottom: none;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
}

.user-avatar-placeholder {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
  flex-shrink: 0;
}

.user-name {
  font-weight: 500;
}

.email-cell {
  color: var(--text-secondary);
}

.date-cell {
  font-size: 13px;
  color: var(--text-muted);
  white-space: nowrap;
}

.role-badge {
  display: inline-block;
  padding: 2px 10px;
  border-radius: var(--radius-full);
  font-size: 12px;
  font-weight: 500;
  background: var(--bg-secondary);
  color: var(--text-secondary);
}

.role-badge.admin {
  background: var(--primary);
  color: white;
}

.btn-delete {
  padding: 4px 12px;
  border: 1px solid #e05555;
  border-radius: var(--radius-md);
  background: transparent;
  color: #e05555;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-delete:hover:not(:disabled) {
  background: #e05555;
  color: white;
}

.btn-delete:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.empty-state {
  text-align: center;
  padding: 40px;
  color: var(--text-muted);
  font-size: 14px;
}

@media (max-width: 768px) {
  .user-table th:nth-child(4),
  .user-table td:nth-child(4) {
    display: none;
  }
}
</style>
