<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getAuthorProfile, updateAuthorProfile, type AuthorProfile as AuthorProfileType } from '@/api/authorProfile'
import toast from '@/utils/toast'

const profile = ref<AuthorProfileType>({
  nickname: '',
  avatar: '',
  background: '',
  signature: '',
  location: '',
  occupation: '',
  school: '',
  major: '',
  email: '',
  wechat: '',
  bio: '',
  tech_stack_frontend: '[]',
  tech_stack_backend: '[]',
  tech_stack_engineering: '[]',
  social_github: '',
  social_bilibili: '',
})
const loading = ref(false)
const saving = ref(false)

const fetchProfile = async () => {
  loading.value = true
  try {
    profile.value = await getAuthorProfile()
  } catch {
    toast.error('获取作者信息失败')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    // 验证 JSON 字段
    const jsonFields = ['tech_stack_frontend', 'tech_stack_backend', 'tech_stack_engineering']
    for (const field of jsonFields) {
      try {
        JSON.parse(profile.value[field as keyof AuthorProfileType] as string)
      } catch {
        toast.error(`${fieldLabel(field)} 格式不正确，请输入合法的 JSON 数组`)
        saving.value = false
        return
      }
    }
    await updateAuthorProfile(profile.value)
    toast.success('保存成功')
  } catch {
    toast.error('保存失败')
  } finally {
    saving.value = false
  }
}

const fieldLabel = (key: string): string => {
  const labels: Record<string, string> = {
    nickname: '昵称', avatar: '头像URL', background: '背景图URL',
    signature: '个性签名', location: '所在城市', occupation: '职业',
    school: '院校', major: '专业', email: '邮箱', wechat: '微信',
    bio: '个人简介', tech_stack_frontend: '前端技术栈',
    tech_stack_backend: '后端技术栈', tech_stack_engineering: '工程化技术栈',
    social_github: 'GitHub 链接', social_bilibili: 'Bilibili 链接',
  }
  return labels[key] || key
}

onMounted(fetchProfile)
</script>

<template>
  <div class="author-profile-page">
    <div class="page-header">
      <h2>作者信息管理</h2>
      <button class="btn-primary" :disabled="saving || loading" @click="handleSave">
        {{ saving ? '保存中...' : '保存修改' }}
      </button>
    </div>

    <div v-if="loading" class="loading-state">加载中...</div>

    <div v-else class="profile-form">
      <!-- 基本信息 -->
      <section class="form-section">
        <h3 class="section-title">基本信息</h3>
        <div class="form-grid">
          <div class="form-item">
            <label>昵称</label>
            <input v-model="profile.nickname" placeholder="请输入昵称" />
          </div>
          <div class="form-item">
            <label>个性签名</label>
            <input v-model="profile.signature" placeholder="请输入个性签名" />
          </div>
          <div class="form-item">
            <label>头像 URL</label>
            <input v-model="profile.avatar" placeholder="/dist/author/avatar.jpg" />
          </div>
          <div class="form-item">
            <label>背景图 URL</label>
            <input v-model="profile.background" placeholder="/dist/author/background.jpg" />
          </div>
        </div>
      </section>

      <!-- 个人信息 -->
      <section class="form-section">
        <h3 class="section-title">个人信息</h3>
        <div class="form-grid">
          <div class="form-item">
            <label>所在城市</label>
            <input v-model="profile.location" placeholder="例如：中国 · 南宁" />
          </div>
          <div class="form-item">
            <label>职业</label>
            <input v-model="profile.occupation" placeholder="例如：学生/开发者" />
          </div>
          <div class="form-item">
            <label>院校</label>
            <input v-model="profile.school" placeholder="例如：广西大学" />
          </div>
          <div class="form-item">
            <label>专业</label>
            <input v-model="profile.major" placeholder="例如：计算机科学与技术" />
          </div>
          <div class="form-item">
            <label>邮箱</label>
            <input v-model="profile.email" placeholder="例如：example@qq.com" />
          </div>
          <div class="form-item">
            <label>微信</label>
            <input v-model="profile.wechat" placeholder="例如：wechat_id" />
          </div>
        </div>
      </section>

      <!-- 个人简介 -->
      <section class="form-section">
        <h3 class="section-title">个人简介</h3>
        <textarea v-model="profile.bio" class="bio-textarea" rows="5" placeholder="请输入个人简介"></textarea>
      </section>

      <!-- 技术栈 -->
      <section class="form-section">
        <h3 class="section-title">技术栈（JSON 数组格式）</h3>
        <div class="form-grid">
          <div class="form-item full-width">
            <label>前端技术栈</label>
            <textarea v-model="profile.tech_stack_frontend" class="json-textarea" rows="2" placeholder='["HTML5","CSS 3","JavaScript"]'></textarea>
          </div>
          <div class="form-item full-width">
            <label>后端技术栈</label>
            <textarea v-model="profile.tech_stack_backend" class="json-textarea" rows="2" placeholder='["Go","MySQL","Redis"]'></textarea>
          </div>
          <div class="form-item full-width">
            <label>工程化技术栈</label>
            <textarea v-model="profile.tech_stack_engineering" class="json-textarea" rows="2" placeholder='["Docker","Git","Nginx"]'></textarea>
          </div>
        </div>
      </section>

      <!-- 社交链接 -->
      <section class="form-section">
        <h3 class="section-title">社交链接</h3>
        <div class="form-grid">
          <div class="form-item">
            <label>GitHub</label>
            <input v-model="profile.social_github" placeholder="https://github.com/username" />
          </div>
          <div class="form-item">
            <label>Bilibili</label>
            <input v-model="profile.social_bilibili" placeholder="https://space.bilibili.com/xxx" />
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.author-profile-page {
  max-width: 900px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.page-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.btn-primary {
  padding: 8px 20px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover:not(:disabled) {
  background: var(--primary-dark);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-state {
  text-align: center;
  padding: 40px;
  color: var(--text-muted);
}

.profile-form {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.form-section {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 24px;
  box-shadow: var(--shadow-soft);
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border);
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-item.full-width {
  grid-column: 1 / -1;
}

.form-item label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
}

.form-item input {
  padding: 8px 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  font-size: 14px;
  color: var(--text-primary);
  background: var(--bg-primary);
  outline: none;
  transition: border-color 0.2s;
}

.form-item input:focus {
  border-color: var(--primary);
}

.bio-textarea,
.json-textarea {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  font-size: 14px;
  color: var(--text-primary);
  background: var(--bg-primary);
  outline: none;
  resize: vertical;
  transition: border-color 0.2s;
}

.json-textarea {
  font-family: 'Courier New', monospace;
  font-size: 13px;
}

.bio-textarea:focus,
.json-textarea:focus {
  border-color: var(--primary);
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
