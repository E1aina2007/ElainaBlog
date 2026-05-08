<!-- Profile.vue 个人信息编辑页面 -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { updateProfile, updatePassword } from '@/api/user'
import AvatarUpload from '@/components/AvatarUpload.vue'
import toast from '@/utils/toast'

const userStore = useUserStore()

const profileForm = ref({
  username: '',
  email: '',
  avatar: '',
})

const passwordForm = ref({
  old_password: '',
  new_password: '',
  confirm_password: '',
})

const saving = ref(false)
const changingPassword = ref(false)

onMounted(() => {
  if (userStore.userInfo) {
    profileForm.value.username = userStore.userInfo.username
    profileForm.value.email = userStore.userInfo.email
    profileForm.value.avatar = userStore.userInfo.avatar || ''
  }
})

const handleSaveProfile = async () => {
  saving.value = true
  try {
    await updateProfile({
      username: profileForm.value.username,
      email: profileForm.value.email,
      avatar: profileForm.value.avatar,
    })
    await userStore.fetchProfile()
    toast.success('保存成功')
  } catch (error) {
    toast.error('保存失败')
  } finally {
    saving.value = false
  }
}

const handleChangePassword = async () => {
  if (passwordForm.value.new_password !== passwordForm.value.confirm_password) {
    toast.warning('两次输入的新密码不一致')
    return
  }
  if (passwordForm.value.new_password.length < 6) {
    toast.warning('新密码长度至少6位')
    return
  }

  changingPassword.value = true
  try {
    await updatePassword({
      old_password: passwordForm.value.old_password,
      new_password: passwordForm.value.new_password,
    })
    passwordForm.value = { old_password: '', new_password: '', confirm_password: '' }
    toast.success('密码修改成功')
  } catch (error) {
    toast.error('密码修改失败，请确认原密码正确')
  } finally {
    changingPassword.value = false
  }
}
</script>

<template>
  <div class="profile-page">
    <div class="page-header">
      <h2>个人中心</h2>
    </div>

    <div class="profile-content">
      <!-- 基本信息 -->
      <div class="profile-card">
        <h3>基本信息</h3>
        <div class="avatar-section">
          <AvatarUpload v-model="profileForm.avatar" :size="64" />
          <div class="user-info">
            <div class="username">{{ userStore.userInfo?.username }}</div>
            <div class="user-role">{{ userStore.userInfo?.is_admin ? '管理员' : '普通用户' }}</div>
          </div>
        </div>

        <div class="form-section">
          <div class="form-group">
            <label>用户名</label>
            <input v-model="profileForm.username" type="text" />
          </div>
          <div class="form-group">
            <label>邮箱</label>
            <input v-model="profileForm.email" type="email" />
          </div>
          <button class="btn-primary" @click="handleSaveProfile" :disabled="saving">
            {{ saving ? '保存中...' : '保存修改' }}
          </button>
        </div>
      </div>

      <!-- 修改密码 -->
      <div class="profile-card">
        <h3>修改密码</h3>
        <div class="form-section">
          <div class="form-group">
            <label>原密码</label>
            <input v-model="passwordForm.old_password" type="password" placeholder="请输入原密码" />
          </div>
          <div class="form-group">
            <label>新密码</label>
            <input v-model="passwordForm.new_password" type="password" placeholder="至少6位字符" />
          </div>
          <div class="form-group">
            <label>确认新密码</label>
            <input v-model="passwordForm.confirm_password" type="password" placeholder="再次输入新密码" />
          </div>
          <button class="btn-primary" @click="handleChangePassword" :disabled="changingPassword">
            {{ changingPassword ? '修改中...' : '修改密码' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.profile-page {
  max-width: 800px;
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

.profile-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.profile-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.profile-card h3 {
  margin: 0 0 20px 0;
  font-size: 18px;
  color: #1f2937;
  font-weight: 600;
}

.avatar-section {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
  padding-bottom: 24px;
  border-bottom: 1px solid #f3f4f6;
}

.avatar-section :deep(.avatar-upload) {
  flex-shrink: 0;
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.username {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.user-role {
  font-size: 14px;
  color: #6366f1;
  background: #6366f115;
  padding: 2px 10px;
  border-radius: 20px;
  display: inline-block;
  width: fit-content;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.form-group input {
  padding: 10px 14px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
}

.form-group input:focus {
  border-color: #6366f1;
}

.btn-primary {
  padding: 12px 24px;
  background: #6366f1;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  width: fit-content;
}

.btn-primary:hover:not(:disabled) {
  background: #4f46e5;
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>