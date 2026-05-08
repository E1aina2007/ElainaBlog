<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getAuthorStats, getMessageList, createMessage, deleteMessage } from '@/api/message'
import { getAuthorProfile } from '@/api/authorProfile'
import toast from '@/utils/toast'

const router = useRouter()
const userStore = useUserStore()
const canComment = computed(() => !!userStore.accessToken)

const authorInfo = ref({
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
  techStack: {
    frontend: [],
    backend: [],
    engineering: [],
  },
  social: {
    github: '',
    bilibili: '',
  },
  stats: {
    articleCount: 0,
    commentCount: 0,
    daysSinceCreated: 0,
  },
})

const DEFAULT_AUTHOR = {
  nickname: '作者名',
  avatar: '/author/avatar.jpg',
  background: '/author/background.jpg',
  signature: '”签名”',
  location: '城市',
  occupation: '职业',
  school: '院校',
  major: '专业',
  email: '邮箱',
  wechat: '微信',
  bio: '个人简介',
}

const fetchAuthorProfile = async () => {
  try {
    const profile = await getAuthorProfile()
    const safeParse = (str, fallback) => {
      try { return JSON.parse(str) } catch { return fallback }
    }
    authorInfo.value = {
      ...authorInfo.value,
      nickname: profile.nickname || DEFAULT_AUTHOR.nickname,
      avatar: profile.avatar || DEFAULT_AUTHOR.avatar,
      background: profile.background || DEFAULT_AUTHOR.background,
      signature: profile.signature || DEFAULT_AUTHOR.signature,
      location: profile.location || DEFAULT_AUTHOR.location,
      occupation: profile.occupation || DEFAULT_AUTHOR.occupation,
      school: profile.school || DEFAULT_AUTHOR.school,
      major: profile.major || DEFAULT_AUTHOR.major,
      email: profile.email || DEFAULT_AUTHOR.email,
      wechat: profile.wechat || DEFAULT_AUTHOR.wechat,
      bio: profile.bio || DEFAULT_AUTHOR.bio,
      techStack: {
        frontend: safeParse(profile.tech_stack_frontend, []),
        backend: safeParse(profile.tech_stack_backend, []),
        engineering: safeParse(profile.tech_stack_engineering, []),
      },
      social: {
        github: profile.social_github || '',
        bilibili: profile.social_bilibili || '',
      },
    }
  } catch (e) {
    console.error('获取作者信息失败:', e)
    // 接口失败时使用默认值
    authorInfo.value = {
      ...authorInfo.value,
      ...DEFAULT_AUTHOR,
    }
  }
}

// 横幅滚动效果
const scrollProgress = ref(0)
const SCROLL_DISTANCE = 250
let rafId = null
const handleBannerScroll = () => {
  if (rafId) return
  rafId = requestAnimationFrame(() => {
    const scrollTop = window.scrollY || document.documentElement.scrollTop
    scrollProgress.value = Math.min(scrollTop / SCROLL_DISTANCE, 1)
    rafId = null
  })
}

// 留言板
const messages = ref([])
const newMessage = ref('')
const sending = ref(false)

const fetchStats = async () => {
  try {
    const stats = await getAuthorStats()
    authorInfo.value.stats.articleCount = stats.article_count
    authorInfo.value.stats.commentCount = stats.comment_count
    authorInfo.value.stats.daysSinceCreated = stats.days_since
  } catch (e) {
    console.error('获取统计失败:', e)
    toast.error(e?.message || '获取统计数据失败')
  }
}

const fetchMessages = async () => {
  try {
    messages.value = (await getMessageList()) ?? []
  } catch (e) {
    console.error('获取留言失败:', e)
    toast.error(e?.message || '获取留言列表失败')
  }
}

const handleSendMessage = async () => {
  if (!canComment.value) {
    router.push({ name: 'login', query: { redirect: '/author' } })
    return
  }
  const content = newMessage.value.trim()
  if (!content) return

  sending.value = true
  try {
    await createMessage(content)
    newMessage.value = ''
    await fetchMessages()
    toast.success('留言成功')
  } catch (e) {
    toast.error(e?.message || '留言失败')
  } finally {
    sending.value = false
  }
}

const handleDeleteMessage = async (id) => {
  if (!confirm('确定删除这条留言？')) return
  try {
    await deleteMessage(id)
    await fetchMessages()
    toast.success('留言已删除')
  } catch (e) {
    toast.error('删除失败')
  }
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(() => {
  fetchAuthorProfile()
  fetchStats()
  fetchMessages()
  window.addEventListener('scroll', handleBannerScroll, { passive: true })
  handleBannerScroll()
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleBannerScroll)
  if (rafId) {
    cancelAnimationFrame(rafId)
    rafId = null
  }
})
</script>

<template>
  <main class="author-page">
    <!-- 顶部背景横幅 -->
    <section class="banner-section">
      <img
        v-if="authorInfo.background"
        :src="authorInfo.background"
        class="banner-bg-img"
        :style="{
          opacity: 1 - scrollProgress * 0.4,
          transform: `scale(${1 + (1 - scrollProgress) * 0.05})`,
        }"
        alt=""
        @error="authorInfo.background = ''"
      />
      <div
        class="banner-bg"
        :style="{
          opacity: 1 - scrollProgress * 0.4,
          transform: `scale(${1 + (1 - scrollProgress) * 0.05})`,
        }"
      ></div>
      <div class="banner-overlay" :style="{ opacity: scrollProgress }"></div>
    </section>

    <div class="page-container">
      <!-- 上方：头像卡片 + 统计卡片 横向排列 -->
      <div class="top-row">
        <!-- 左：头像与昵称小卡片 -->
        <section class="profile-card">
          <div class="avatar">
            <img v-if="authorInfo.avatar" :src="authorInfo.avatar" class="avatar-img" alt="头像" @error="authorInfo.avatar = ''" />
            <span v-else class="avatar-placeholder">{{ (authorInfo.nickname || '?').charAt(0) }}</span>
          </div>
          <h1 class="nickname">{{ authorInfo.nickname }}</h1>
        </section>

        <!-- 右：统计卡片 -->
        <section class="stats-card">
          <div class="stat-item">
            <span class="stat-number">{{ authorInfo.stats.articleCount }}</span>
            <span class="stat-label">文章数</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-number">{{ authorInfo.stats.commentCount }}</span>
            <span class="stat-label">评论数</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-number">{{ authorInfo.stats.daysSinceCreated }}天</span>
            <span class="stat-label">建站天数</span>
          </div>
        </section>
      </div>

      <!-- 签名卡片 -->
      <section class="signature-card">
        <p class="signature"> {{ authorInfo.signature }} </p>
      </section>

      <!-- 个人信息卡片 -->
      <section class="info-card">
        <div class="info-item">
          <span class="info-label">城市</span>
          <span class="info-value">{{ authorInfo.location }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">职业</span>
          <span class="info-value">{{ authorInfo.occupation }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">院校</span>
          <span class="info-value">{{ authorInfo.school }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">专业</span>
          <span class="info-value">{{ authorInfo.major }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">邮箱</span>
          <span class="info-value">{{ authorInfo.email }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">微信</span>
          <span class="info-value">{{ authorInfo.wechat }}</span>
        </div>
      </section>

      <!-- 个人简介 -->
      <section class="bio-card">
        <h2 class="section-title">个人简介</h2>
        <div class="bio-divider"></div>
        <p class="bio-text">{{ authorInfo.bio }}</p>
      </section>

      <!-- 技术栈 -->
      <section class="tech-card">
        <h2 class="section-title">技术栈</h2>
        <div class="bio-divider"></div>
        <div class="tech-groups">
          <div class="tech-group">
            <h3 class="tech-group-title">前端</h3>
            <div class="tech-tags">
              <span v-for="tech in authorInfo.techStack.frontend" :key="tech" class="tech-tag">{{ tech }}</span>
            </div>
          </div>
          <div class="tech-group">
            <h3 class="tech-group-title">后端</h3>
            <div class="tech-tags">
              <span v-for="tech in authorInfo.techStack.backend" :key="tech" class="tech-tag">{{ tech }}</span>
            </div>
          </div>
          <div class="tech-group">
            <h3 class="tech-group-title">工程化</h3>
            <div class="tech-tags">
              <span v-for="tech in authorInfo.techStack.engineering" :key="tech" class="tech-tag">{{ tech }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- 社交链接 -->
      <section class="social-card">
        <h2 class="section-title">社交链接</h2>
        <div class="social-links">
          <a
            v-if="authorInfo.social.github"
            :href="authorInfo.social.github"
            target="_blank"
            rel="noopener noreferrer"
            class="social-btn"
            title="GitHub"
          >
            <svg class="social-svg" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/>
            </svg>
            <span class="social-name">GitHub</span>
          </a>
          <a
            v-if="authorInfo.social.bilibili"
            :href="authorInfo.social.bilibili"
            target="_blank"
            rel="noopener noreferrer"
            class="social-btn"
            title="Bilibili"
          >
            <svg class="social-svg" viewBox="0 0 24 24" fill="currentColor">
              <path d="M17.813 4.653h.854c1.51.044 2.769.578 3.773 1.608 1.004 1.029 1.524 2.29 1.56 3.777v7.36c-.044 1.487-.565 2.748-1.56 3.777-1.004 1.03-2.262 1.564-3.773 1.608H5.333c-1.51-.044-2.769-.578-3.773-1.608C.556 19.748.036 18.487 0 17V9.653c.036-1.487.556-2.748 1.56-3.777 1.004-1.03 2.262-1.564 3.773-1.608h.774l-1.174-1.12a1.234 1.234 0 01-.373-.906c0-.356.124-.658.373-.907l.027-.027c.267-.249.573-.373.92-.373.347 0 .653.124.92.373l3.307 3.16h2.507l3.307-3.16c.267-.249.573-.373.92-.373.347 0 .653.124.92.373l.027.027c.249.249.373.551.373.907 0 .355-.124.657-.373.906L17.813 4.653zM5.333 6.12c-.8.027-1.466.311-2 .987-.533.676-.8 1.458-.8 2.347v7.36c0 .889.267 1.671.8 2.347.534.676 1.2 1.014 2 1.014h13.334c.8 0 1.466-.338 2-1.014.533-.676.8-1.458.8-2.347v-7.36c0-.889-.267-1.671-.8-2.347-.534-.676-1.2-1.014-2-1.014H5.333zM8 11.107c.733 0 1.333.6 1.333 1.334v1.333c0 .733-.6 1.334-1.333 1.334s-1.334-.6-1.334-1.334v-1.333c0-.734.6-1.334 1.334-1.334zm8 0c.733 0 1.333.6 1.333 1.334v1.333c0 .733-.6 1.334-1.333 1.334s-1.334-.6-1.334-1.334v-1.333c0-.734.6-1.334 1.334-1.334z"/>
            </svg>
            <span class="social-name">Bilibili</span>
          </a>
        </div>
      </section>

      <!-- 留言板 -->
      <section class="message-card">
        <h2 class="section-title">留言板</h2>
        <div class="bio-divider"></div>

        <!-- 留言输入 -->
        <div class="message-input-area">
          <textarea
            v-model="newMessage"
            class="message-textarea"
            rows="3"
            :placeholder="canComment ? '写下你的留言...' : '登录后即可留言'"
            :disabled="!canComment"
          ></textarea>
          <div class="message-actions">
            <span v-if="!canComment" class="login-hint">
              <router-link to="/login" class="login-link">登录</router-link>后即可留言
            </span>
            <button
              class="btn-send"
              :disabled="sending || !newMessage.trim() || !canComment"
              @click="handleSendMessage"
            >
              {{ sending ? '发送中...' : '发表留言' }}
            </button>
          </div>
        </div>

        <!-- 留言列表 -->
        <div class="message-list">
          <div v-if="messages.length === 0" class="message-empty">
            还没有留言，快来写一条吧~
          </div>
          <div v-for="msg in messages" :key="msg.id" class="message-item">
            <div class="message-header">
              <div class="message-user">
                <img v-if="msg.avatar" :src="msg.avatar" class="message-avatar" alt="" />
                <span v-else class="message-avatar-placeholder">{{ msg.username.charAt(0) }}</span>
                <span class="message-username">{{ msg.username }}</span>
              </div>
              <div class="message-meta">
                <span class="message-time">{{ formatDate(msg.created_at) }}</span>
                <button
                  v-if="userStore.isAdmin || userStore.userInfo?.id === msg.user_id"
                  class="btn-delete-msg"
                  @click="handleDeleteMessage(msg.id)"
                  title="删除"
                >
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3 6 5 6 21 6"></polyline>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                  </svg>
                </button>
              </div>
            </div>
            <p class="message-content">{{ msg.content }}</p>
          </div>
        </div>
      </section>
    </div>
  </main>
</template>

<style scoped>
.author-page {
  min-height: 60vh;
  background: linear-gradient(135deg, var(--bg-primary) 0%, var(--bg-secondary) 100%);
}

/* ===== 顶部背景横幅 ===== */
.banner-section {
  width: 100%;
  height: 500px;
  position: relative;
  overflow: hidden;
}

.banner-bg {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, var(--primary-lighter) 0%, var(--primary-light) 50%, var(--accent) 100%);
  transition: opacity 0.3s ease-out, transform 0.3s ease-out;
}

.banner-bg-img {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: opacity 0.3s ease-out, transform 0.3s ease-out;
  z-index: 1;
}

.banner-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, var(--primary-lighter) 0%, var(--primary-light) 50%, var(--accent) 100%);
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.3s ease-out;
}

.page-container {
  max-width: 800px;
  margin: -80px auto 0;
  padding: 0 24px 60px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  position: relative;
  z-index: 1;
}

/* ===== 上方横排：头像 + 统计 ===== */
.top-row {
  display: flex;
  gap: 20px;
  align-items: stretch;
}

/* 左：头像与昵称小卡片 */
.profile-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 28px 24px;
  box-shadow: var(--shadow-soft);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
  width: 200px;
}

.avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  font-size: 40px;
  font-weight: 600;
  color: white;
  line-height: 1;
}

.nickname {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

/* 右：统计卡片 */
.stats-card {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0;
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(10px);
  border-radius: var(--radius-lg);
  padding: 24px 32px;
  box-shadow: var(--shadow-soft);
}

.stat-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.stat-number {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--primary);
}

.stat-label {
  font-size: 0.8125rem;
  color: var(--text-muted);
}

.stat-divider {
  width: 1px;
  height: 36px;
  background: var(--border);
}

/* ===== 签名卡片 ===== */
.signature-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 20px 28px;
  box-shadow: var(--shadow-soft);
  text-align: center;
}

.signature {
  font-size: 0.9375rem;
  color: var(--text-secondary);
  margin: 0;
  line-height: 1.6;
}

/* ===== 个人信息 ===== */
.info-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 24px 28px;
  box-shadow: var(--shadow-soft);
  display: flex;
  flex-direction: column;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 0;
  border-bottom: 1px solid var(--border);
}

.info-item:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.info-item:first-child {
  padding-top: 0;
}

.info-icon {
  font-size: 18px;
  width: 24px;
  text-align: center;
  flex-shrink: 0;
}

.info-label {
  font-size: 0.875rem;
  color: var(--text-muted);
  width: 48px;
  flex-shrink: 0;
}

.info-value {
  font-size: 0.9375rem;
  color: var(--text-primary);
  font-weight: 500;
}

/* ===== 个人简介 ===== */
.bio-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 28px;
  box-shadow: var(--shadow-soft);
}

/* ===== 技术栈 ===== */
.tech-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 28px;
  box-shadow: var(--shadow-soft);
}

.tech-groups {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.tech-group {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.tech-group-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--primary);
  margin: 0;
  padding-left: 8px;
  border-left: 3px solid var(--primary);
}

.tech-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tech-tag {
  display: inline-block;
  padding: 6px 14px;
  background: var(--primary);
  border: 1px solid var(--primary);
  border-radius: var(--radius-md);
  font-size: 0.8125rem;
  font-weight: 500;
  color: white;
  transition: all 0.2s ease-out;
}

.tech-tag:hover {
  background: var(--primary-dark);
  border-color: var(--primary-dark);
  transform: translateY(-4px);
  box-shadow: 0 4px 16px rgba(126, 215, 193, 0.35);
}

.section-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 12px;
}

.bio-divider {
  width: 40px;
  height: 3px;
  background: var(--primary);
  border-radius: 2px;
  margin-bottom: 16px;
}

.bio-text {
  font-size: 0.9375rem;
  color: var(--text-secondary);
  line-height: 1.8;
  margin: 0;
  white-space: pre-line;
}

/* ===== 社交链接 ===== */
.social-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 28px;
  box-shadow: var(--shadow-soft);
}

.social-links {
  display: flex;
  gap: 16px;
  margin-top: 4px;
}

.social-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  text-decoration: none;
  color: var(--text-secondary);
  font-size: 0.875rem;
  font-weight: 500;
  transition: transform 0.2s ease-out, box-shadow 0.2s ease-out, background 0.2s ease-out, color 0.2s ease-out;
}

.social-btn:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-hover);
  background: var(--primary-lighter);
  color: var(--primary-dark);
  border-color: var(--primary-light);
}

.social-svg {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.social-name {
  font-weight: 500;
}

/* ===== 留言板 ===== */
.message-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 28px;
  box-shadow: var(--shadow-soft);
}

.message-input-area {
  margin-bottom: 24px;
}

.message-textarea {
  width: 100%;
  padding: 14px 16px;
  font-size: 0.9375rem;
  line-height: 1.6;
  color: var(--text-primary);
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  outline: none;
  resize: vertical;
  transition: border-color 0.2s;
}

.message-textarea:focus {
  border-color: var(--primary);
}

.message-textarea::placeholder {
  color: var(--text-muted);
}

.message-textarea:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.message-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 10px;
}

.login-hint {
  font-size: 0.8125rem;
  color: var(--text-muted);
}

.login-link {
  color: var(--primary-dark);
  text-decoration: none;
  font-weight: 500;
}

.login-link:hover {
  text-decoration: underline;
}

.btn-send {
  padding: 8px 24px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-send:hover:not(:disabled) {
  background: var(--primary-dark);
  transform: translateY(-1px);
}

.btn-send:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.message-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.message-empty {
  text-align: center;
  padding: 32px 0;
  color: var(--text-muted);
  font-size: 0.875rem;
}

.message-item {
  padding: 16px;
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  border: 1px solid var(--border);
}

.message-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.message-user {
  display: flex;
  align-items: center;
  gap: 10px;
}

.message-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
}

.message-avatar-placeholder {
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

.message-username {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-primary);
}

.message-meta {
  display: flex;
  align-items: center;
  gap: 10px;
}

.message-time {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.btn-delete-msg {
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: all 0.2s;
  display: flex;
  align-items: center;
}

.btn-delete-msg:hover {
  color: #e05555;
  background: rgba(224, 85, 85, 0.1);
}

.message-content {
  font-size: 0.9375rem;
  color: var(--text-secondary);
  line-height: 1.7;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}

/* ===== 响应式 ===== */
@media (max-width: 768px) {
  .banner-section {
    height: 440px;
  }

  .page-container {
    margin-top: -60px;
    padding: 0 16px 40px;
    gap: 16px;
  }

  .top-row {
    flex-direction: column;
    align-items: stretch;
  }

  .profile-card {
    width: 100%;
    flex-direction: row;
    padding: 20px 24px;
    gap: 16px;
  }

  .avatar {
    width: 72px;
    height: 72px;
  }

  .avatar-placeholder {
    font-size: 28px;
  }

  .nickname {
    font-size: 1.125rem;
  }

  .stats-card {
    padding: 20px 16px;
  }

  .stat-number {
    font-size: 1.25rem;
  }

  .social-links {
    flex-wrap: wrap;
  }

  .social-btn {
    padding: 10px 16px;
  }
}
</style>
