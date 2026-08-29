<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { getProfile, updateProfile, updatePassword } from '@/api/user'
import AvatarUpload from '@/components/AvatarUpload.vue'

const userStore = useUserStore()

const activeTab = ref<'info' | 'password'>('info')
const isLoading = ref(false)
const message = ref('')
const error = ref('')

// 用户信息表单
const profileForm = ref({
  username: '',
  email: '',
  avatar: '',
})

// 密码表单
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

// 加载用户信息
async function loadProfile() {
  try {
    const profile = await getProfile()
    profileForm.value = {
      username: profile.username,
      email: profile.email,
      avatar: profile.avatar || '',
    }
  } catch (err: any) {
    error.value = err?.message || '加载用户信息失败'
  }
}

// 更新个人信息
async function handleUpdateProfile() {
  error.value = ''
  message.value = ''

  if (!profileForm.value.username.trim()) {
    error.value = '用户名不能为空'
    return
  }

  isLoading.value = true
  try {
    await updateProfile({
      username: profileForm.value.username,
      email: profileForm.value.email,
      avatar: profileForm.value.avatar,
    })
    message.value = '个人信息更新成功'
    // 刷新用户信息
    await userStore.fetchProfile()
  } catch (err: any) {
    error.value = err?.message || '更新失败'
  } finally {
    isLoading.value = false
  }
}

// 修改密码
async function handleUpdatePassword() {
  error.value = ''
  message.value = ''

  if (!passwordForm.value.oldPassword) {
    error.value = '请输入原密码'
    return
  }
  if (!passwordForm.value.newPassword) {
    error.value = '请输入新密码'
    return
  }
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    error.value = '两次输入的新密码不一致'
    return
  }

  isLoading.value = true
  try {
    await updatePassword({
      old_password: passwordForm.value.oldPassword,
      new_password: passwordForm.value.newPassword,
    })
    message.value = '密码修改成功'
    passwordForm.value = { oldPassword: '', newPassword: '', confirmPassword: '' }
  } catch (err: any) {
    error.value = err?.message || '密码修改失败'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  loadProfile()
})
</script>

<template>
  <main class="profile-page">
    <div class="container">
      <h1 class="page-title">个人中心</h1>

      <!-- 消息提示 -->
      <div v-if="message" class="message-banner success">
        {{ message }}
      </div>
      <div v-if="error" class="message-banner error">
        {{ error }}
      </div>

      <div class="profile-layout">
        <!-- 侧边栏 -->
        <aside class="profile-sidebar">
          <div class="user-card">
            <AvatarUpload v-model="profileForm.avatar" :size="80" />
            <h2 class="user-name">{{ userStore.userInfo?.username }}</h2>
            <p class="user-email">{{ userStore.userInfo?.email }}</p>
            <span v-if="userStore.isAdmin" class="admin-badge">管理员</span>
          </div>

          <nav class="profile-nav">
            <button
              class="nav-item"
              :class="{ active: activeTab === 'info' }"
              @click="activeTab = 'info'"
            >
              个人信息
            </button>
            <button
              class="nav-item"
              :class="{ active: activeTab === 'password' }"
              @click="activeTab = 'password'"
            >
              修改密码
            </button>
          </nav>
        </aside>

        <!-- 主内容区 -->
        <div class="profile-content">
          <!-- 个人信息 -->
          <div v-if="activeTab === 'info'" class="content-card">
            <h3 class="card-title">编辑个人信息</h3>
            <form class="form" @submit.prevent="handleUpdateProfile">
              <div class="form-group">
                <label class="form-label">用户名</label>
                <input
                  v-model="profileForm.username"
                  type="text"
                  class="form-input"
                  placeholder="请输入用户名"
                />
              </div>
              <div class="form-group">
                <label class="form-label">邮箱</label>
                <input
                  v-model="profileForm.email"
                  type="email"
                  class="form-input"
                  placeholder="请输入邮箱"
                />
              </div>
              <div class="form-group">
                <label class="form-label">头像</label>
                <AvatarUpload v-model="profileForm.avatar" :size="96" />
              </div>
              <div class="form-actions">
                <button type="submit" class="btn-primary" :disabled="isLoading">
                  {{ isLoading ? '保存中...' : '保存修改' }}
                </button>
              </div>
            </form>
          </div>

          <!-- 修改密码 -->
          <div v-if="activeTab === 'password'" class="content-card">
            <h3 class="card-title">修改密码</h3>
            <form class="form" @submit.prevent="handleUpdatePassword">
              <div class="form-group">
                <label class="form-label">原密码</label>
                <input
                  v-model="passwordForm.oldPassword"
                  type="password"
                  class="form-input"
                  placeholder="请输入原密码"
                />
              </div>
              <div class="form-group">
                <label class="form-label">新密码</label>
                <input
                  v-model="passwordForm.newPassword"
                  type="password"
                  class="form-input"
                  placeholder="请输入新密码（8-72位）"
                />
              </div>
              <div class="form-group">
                <label class="form-label">确认新密码</label>
                <input
                  v-model="passwordForm.confirmPassword"
                  type="password"
                  class="form-input"
                  placeholder="请再次输入新密码"
                />
              </div>
              <div class="form-actions">
                <button type="submit" class="btn-primary" :disabled="isLoading">
                  {{ isLoading ? '修改中...' : '修改密码' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped>
.profile-page {
  min-height: 100vh;
  padding: 40px 24px;
  background: var(--bg-primary);
}

.container {
  max-width: 900px;
  margin: 0 auto;
}

.page-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 24px;
}

/* 消息提示 */
.message-banner {
  padding: 12px 16px;
  border-radius: var(--radius-md);
  margin-bottom: 20px;
  font-size: 0.875rem;
}

.message-banner.success {
  background: rgba(126, 215, 193, 0.15);
  color: #27ae60;
  border: 1px solid var(--primary);
}

.message-banner.error {
  background: rgba(255, 183, 178, 0.15);
  color: #c0392b;
  border: 1px solid var(--accent);
}

/* 布局 */
.profile-layout {
  display: grid;
  grid-template-columns: 260px 1fr;
  gap: 24px;
}

/* 侧边栏 */
.profile-sidebar {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.user-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 24px;
  text-align: center;
  box-shadow: var(--shadow-soft);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.user-name {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 4px;
}

.user-email {
  font-size: 0.8125rem;
  color: var(--text-muted);
  margin: 0 0 12px;
}

.admin-badge {
  display: inline-block;
  padding: 4px 12px;
  background: linear-gradient(135deg, var(--accent) 0%, var(--highlight) 100%);
  color: white;
  font-size: 0.75rem;
  font-weight: 500;
  border-radius: var(--radius-sm);
}

/* 导航 */
.profile-nav {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 12px;
  box-shadow: var(--shadow-soft);
}

.nav-item {
  display: block;
  width: 100%;
  padding: 12px 16px;
  text-align: left;
  font-size: 0.9375rem;
  color: var(--text-secondary);
  background: none;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.2s ease;
}

.nav-item:hover {
  background: var(--bg-secondary);
  color: var(--primary);
}

.nav-item.active {
  background: linear-gradient(135deg, var(--primary-lighter) 0%, rgba(126, 215, 193, 0.1) 100%);
  color: var(--primary);
  font-weight: 500;
}

/* 主内容 */
.profile-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.content-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 24px;
  box-shadow: var(--shadow-soft);
}

.card-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border);
}

/* 表单 */
.form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-secondary);
}

.form-input {
  padding: 12px 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  font-size: 0.9375rem;
  color: var(--text-primary);
  outline: none;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.form-input::placeholder {
  color: var(--text-muted);
}

.form-input:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(126, 215, 193, 0.15);
}

.form-actions {
  margin-top: 8px;
}

.btn-primary {
  padding: 12px 24px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-size: 0.9375rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease-out;
}

.btn-primary:hover:not(:disabled) {
  background: var(--primary-dark);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(126, 215, 193, 0.3);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 响应式 */
@media (max-width: 768px) {
  .profile-layout {
    grid-template-columns: 1fr;
  }

  .profile-sidebar {
    order: -1;
  }

  .user-card {
    display: flex;
    align-items: center;
    gap: 16px;
    text-align: left;
  }

  .user-avatar-large {
    width: 60px;
    height: 60px;
    font-size: 24px;
    margin: 0;
  }

  .profile-nav {
    display: flex;
    gap: 8px;
  }

  .nav-item {
    flex: 1;
    text-align: center;
  }
}
</style>
