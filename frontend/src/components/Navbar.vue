<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useSiteStore } from '@/stores/site'
import { useTheme } from '@/composables/useTheme'
import {
  getNotificationList,
  getUnreadCount,
  markAsRead,
  markAllAsRead,
  type Notification,
} from '@/api/notification'

const router = useRouter()
const userStore = useUserStore()
const siteStore = useSiteStore()
const { isDark, toggleDark } = useTheme()

function handleToggleDark() {
  document.documentElement.classList.add('transitioning')
  toggleDark()
  setTimeout(() => {
    document.documentElement.classList.remove('transitioning')
  }, 350)
}
const mobileMenuOpen = ref(false)
const userDropdownOpen = ref(false)

const siteName = computed(() => siteStore.get('site_name'))

const navLinks = [
  { name: '首页', path: '/' },
  { name: '工具', path: '/tools' },
  { name: '关于作者', path: '/author' },
]

const isActive = (path: string) => router.currentRoute.value.path === path

const isLoggedIn = computed(() => userStore.isLoggedIn)
const isAdmin = computed(() => userStore.isAdmin)
const username = computed(() => userStore.userInfo?.username || '')
const avatarUrl = computed(() => userStore.userInfo?.avatar || '')

async function handleLogout() {
  await userStore.logout()
  userDropdownOpen.value = false
  router.push('/')
}

// 点击外部关闭下拉菜单
function handleClickOutside(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (!target.closest('.user-dropdown-wrapper')) {
    userDropdownOpen.value = false
  }
  if (!target.closest('.notif-wrapper')) {
    notificationOpen.value = false
  }
}

// 通知相关
const notificationOpen = ref(false)
const unreadCount = ref(0)
const notifications = ref<Notification[]>([])
let pollTimer: ReturnType<typeof setInterval> | null = null

const fetchUnreadCount = async () => {
  if (!userStore.isLoggedIn) return
  try {
    const res = await getUnreadCount()
    unreadCount.value = res.count ?? 0
  } catch {
    // 静默失败
  }
}

const fetchNotifications = async () => {
  try {
    notifications.value = (await getNotificationList()) ?? []
  } catch {
    // 静默失败
  }
}

const handleToggleNotifications = async () => {
  notificationOpen.value = !notificationOpen.value
  if (notificationOpen.value) {
    await fetchNotifications()
  }
}

const handleMarkAsRead = async (id: number) => {
  try {
    await markAsRead(id)
    const n = notifications.value.find((n) => n.id === id)
    if (n) n.is_read = true
    unreadCount.value = Math.max(0, unreadCount.value - 1)
  } catch {
    // 静默失败
  }
}

const handleMarkAllAsRead = async () => {
  try {
    await markAllAsRead()
    notifications.value.forEach((n) => (n.is_read = true))
    unreadCount.value = 0
  } catch {
    // 静默失败
  }
}

const handleNotificationClick = async (n: Notification) => {
  if (!n.is_read) {
    await handleMarkAsRead(n.id)
  }
  if (n.target_id > 0) {
    notificationOpen.value = false
    router.push(`/article/${n.target_id}`)
  }
}

const formatNotifTime = (date: string) => {
  const d = new Date(date)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  const minutes = Math.floor(diff / 60000)
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  if (days < 30) return `${days}天前`
  return d.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
}

const notifTypeIcon = (type: string) => {
  if (type === 'comment') return '💬'
  if (type === 'message') return '📝'
  return '🔔'
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  // 轮询未读通知
  if (userStore.isLoggedIn) {
    fetchUnreadCount()
    pollTimer = setInterval(fetchUnreadCount, 60000)
  }
})

onUnmounted(() => {
  if (pollTimer) {
    clearInterval(pollTimer)
  }
})
</script>

<template>
  <nav class="navbar">
    <div class="navbar-container">
      <!-- Logo -->
      <router-link to="/" class="logo">
        <img class="logo-icon" src="/favicon.ico" alt="logo" />
        <span class="logo-text">{{ siteName }}</span>
      </router-link>

      <!-- Desktop Navigation -->
      <div class="nav-links hide-mobile">
        <router-link
          v-for="link in navLinks"
          :key="link.path"
          :to="link.path"
          class="nav-link"
          :class="{ active: isActive(link.path) }"
        >
          {{ link.name }}
        </router-link>
        <!-- 管理员面板入口 -->
        <router-link
          v-if="isLoggedIn && isAdmin"
          to="/admin"
          class="nav-link admin-link"
          :class="{ active: isActive('/admin') }"
        >
          管理面板
        </router-link>
        <!-- 我的文章入口 -->
        <router-link
          v-if="isLoggedIn"
          to="/my-articles"
          class="nav-link"
          :class="{ active: isActive('/my-articles') }"
        >
          我的文章
        </router-link>
      </div>

      <!-- Right Actions -->
      <div class="nav-actions">
        <!-- 主题切换 -->
        <button class="theme-toggle" :title="isDark ? '切换到浅色模式' : '切换到深色模式'" @click="handleToggleDark()">
          <svg v-if="isDark" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="5" />
            <line x1="12" y1="1" x2="12" y2="3" />
            <line x1="12" y1="21" x2="12" y2="23" />
            <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
            <line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
            <line x1="1" y1="12" x2="3" y2="12" />
            <line x1="21" y1="12" x2="23" y2="12" />
            <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
            <line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
          </svg>
          <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
          </svg>
        </button>

        <!-- 通知铃铛（已登录） -->
        <div v-if="isLoggedIn" class="notif-wrapper">
          <button class="notif-btn" title="通知" @click.stop="handleToggleNotifications">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
              <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
            </svg>
            <span v-if="unreadCount > 0" class="notif-badge">
              {{ unreadCount > 99 ? '99+' : unreadCount }}
            </span>
          </button>

          <transition name="dropdown">
            <div v-if="notificationOpen" class="notif-dropdown" @click.stop>
              <div class="notif-header">
                <span class="notif-title">通知</span>
                <button
                  v-if="unreadCount > 0"
                  class="notif-read-all"
                  @click="handleMarkAllAsRead"
                >
                  全部已读
                </button>
              </div>
              <div class="notif-divider"></div>
              <div class="notif-list">
                <div v-if="notifications.length === 0" class="notif-empty">
                  暂无新通知
                </div>
                <div
                  v-for="n in notifications"
                  :key="n.id"
                  class="notif-item"
                  :class="{ unread: !n.is_read }"
                  @click="handleNotificationClick(n)"
                >
                  <span class="notif-icon">{{ notifTypeIcon(n.type) }}</span>
                  <div class="notif-body">
                    <p class="notif-text">{{ n.title }}</p>
                    <p v-if="n.content" class="notif-content">{{ n.content }}</p>
                    <span class="notif-time">{{ formatNotifTime(n.created_at) }}</span>
                  </div>
                </div>
              </div>
            </div>
          </transition>
        </div>

        <!-- 未登录：显示登录按钮 -->
        <router-link v-if="!isLoggedIn" to="/login" class="login-btn">
          <span>登录</span>
        </router-link>

        <!-- 已登录：用户头像+用户名下拉 -->
        <div v-else class="user-dropdown-wrapper">
          <button class="user-info-btn" @click="userDropdownOpen = !userDropdownOpen">
            <img v-if="avatarUrl" :src="avatarUrl" class="user-avatar-img" alt="头像" />
            <span v-else class="user-avatar">{{ username.charAt(0).toUpperCase() }}</span>
            <span class="user-name">{{ username }}</span>
          </button>

          <transition name="dropdown">
            <div v-if="userDropdownOpen" class="user-dropdown">
              <div class="dropdown-header">
                <span class="dropdown-username">{{ username }}</span>
              </div>
              <div class="dropdown-divider"></div>
              <router-link
                v-if="isAdmin"
                to="/admin"
                class="dropdown-item"
                @click="userDropdownOpen = false"
              >
                管理面板
              </router-link>
              <router-link to="/my-articles" class="dropdown-item" @click="userDropdownOpen = false">
                我的文章
              </router-link>
              <router-link to="/profile" class="dropdown-item" @click="userDropdownOpen = false">
                个人中心
              </router-link>
              <button class="dropdown-item logout" @click="handleLogout">
                退出登录
              </button>
            </div>
          </transition>
        </div>

        <!-- Mobile Menu Button -->
        <button class="menu-toggle hide-desktop" @click="mobileMenuOpen = !mobileMenuOpen">
          <svg v-if="!mobileMenuOpen" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="3" y1="12" x2="21" y2="12"></line>
            <line x1="3" y1="6" x2="21" y2="6"></line>
            <line x1="3" y1="18" x2="21" y2="18"></line>
          </svg>
          <svg v-else width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>
    </div>

    <!-- Mobile Menu -->
    <transition name="slide">
      <div v-if="mobileMenuOpen" class="mobile-menu hide-desktop">
        <router-link
          v-for="link in navLinks"
          :key="link.path"
          :to="link.path"
          class="mobile-nav-link"
          :class="{ active: isActive(link.path) }"
          @click="mobileMenuOpen = false"
        >
          {{ link.name }}
        </router-link>
        <router-link
          v-if="isLoggedIn && isAdmin"
          to="/admin"
          class="mobile-nav-link admin-link"
          :class="{ active: isActive('/admin') }"
          @click="mobileMenuOpen = false"
        >
          管理面板
        </router-link>
        <router-link
          v-if="isLoggedIn"
          to="/my-articles"
          class="mobile-nav-link"
          :class="{ active: isActive('/my-articles') }"
          @click="mobileMenuOpen = false"
        >
          我的文章
        </router-link>
        <div class="mobile-nav-link theme-row">
          <span>{{ isDark ? '深色模式' : '浅色模式' }}</span>
          <button class="theme-toggle-mobile" @click="handleToggleDark()">
            <svg v-if="isDark" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="5" />
              <line x1="12" y1="1" x2="12" y2="3" />
              <line x1="12" y1="21" x2="12" y2="23" />
              <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
              <line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
              <line x1="1" y1="12" x2="3" y2="12" />
              <line x1="21" y1="12" x2="23" y2="12" />
              <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
              <line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
            </svg>
            <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
            </svg>
          </button>
        </div>
        <div v-if="isLoggedIn" class="mobile-menu-divider"></div>
        <button v-if="isLoggedIn" class="mobile-nav-link logout" @click="handleLogout(); mobileMenuOpen = false">
          退出登录 ({{ username }})
        </button>
        <router-link v-if="!isLoggedIn" to="/login" class="mobile-nav-link" @click="mobileMenuOpen = false">
          登录
        </router-link>
      </div>
    </transition>
  </nav>
</template>

<style scoped>
.navbar {
  position: sticky;
  top: 0;
  z-index: 100;
  background: var(--bg-glass);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--border);
}

.navbar-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

/* Logo */
.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  text-decoration: none;
  transition: transform var(--transition-fast);
}

.logo:hover {
  transform: scale(1.02);
}

.logo-icon {
  width: 28px;
  height: 28px;
  object-fit: contain;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -0.5px;
}

.logo-highlight {
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* Desktop Nav Links */
.nav-links {
  display: flex;
  align-items: center;
  gap: 8px;
}

.nav-link {
  position: relative;
  padding: 8px 16px;
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 15px;
  font-weight: 500;
  border-radius: var(--radius-sm);
  transition: all var(--transition-fast);
}

.nav-link:hover {
  color: var(--primary);
  background: var(--bg-secondary);
}

.nav-link.active {
  color: var(--primary);
  background: linear-gradient(135deg, var(--primary-lighter) 0%, rgba(126, 215, 193, 0.1) 100%);
}

.nav-link.active::after {
  content: '';
  position: absolute;
  bottom: 4px;
  left: 50%;
  transform: translateX(-50%);
  width: 16px;
  height: 3px;
  background: var(--primary);
  border-radius: 2px;
}

/* Actions */
.nav-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* Theme Toggle */
.theme-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  color: var(--text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.theme-toggle:hover {
  color: var(--primary);
  background: var(--bg-secondary);
}

/* User Dropdown */
.user-dropdown-wrapper {
  position: relative;
}

.user-info-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 14px 6px 6px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 24px;
  cursor: pointer;
  transition: all 0.2s ease-out;
}

.user-info-btn:hover {
  border-color: var(--primary-light);
  box-shadow: 0 2px 12px rgba(126, 215, 193, 0.2);
}

.user-avatar {
  width: 28px;
  height: 28px;
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

.user-avatar-img {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
}

.user-name {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
  max-width: 80px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-dropdown {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  min-width: 160px;
  background: var(--bg-card);
  border-radius: var(--radius-md);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  border: 1px solid var(--border);
  padding: 8px 0;
  z-index: 200;
}

.dropdown-header {
  padding: 10px 16px;
}

.dropdown-username {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-primary);
}

.dropdown-divider {
  height: 1px;
  background: var(--border);
  margin: 4px 0;
}

.dropdown-item {
  display: block;
  width: 100%;
  padding: 10px 16px;
  font-size: 0.875rem;
  color: var(--text-secondary);
  text-decoration: none;
  background: none;
  border: none;
  cursor: pointer;
  text-align: left;
  transition: background 0.15s ease, color 0.15s ease;
}

.dropdown-item:hover {
  background: var(--bg-secondary);
  color: var(--primary);
}

.dropdown-item.logout {
  color: var(--color-danger);
}

.dropdown-item.logout:hover {
  background: rgba(239, 68, 68, 0.08);
  color: var(--color-danger-hover);
}

/* Admin link */
.admin-link {
  color: var(--primary) !important;
  font-weight: 600;
}

/* Dropdown transition */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

.mobile-menu-divider {
  height: 1px;
  background: var(--border);
  margin: 8px 16px;
}

.mobile-nav-link.logout {
  color: var(--color-danger);
  background: none;
  border: none;
  cursor: pointer;
  font-size: 15px;
  font-weight: 500;
  text-align: left;
  width: 100%;
}

.mobile-nav-link.theme-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: default;
}

.theme-toggle-mobile {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  color: var(--text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.theme-toggle-mobile:hover {
  color: var(--primary);
  background: var(--bg-secondary);
}

.search-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  color: var(--text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.search-btn:hover {
  color: var(--primary);
  background: var(--bg-secondary);
}

.login-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 20px;
  background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
  color: white;
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  border-radius: var(--radius-sm);
  transition: all var(--transition-fast);
  box-shadow: var(--shadow-soft);
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-hover);
}

.menu-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  color: var(--text-primary);
  background: transparent;
  border: none;
  cursor: pointer;
}

/* Mobile Menu */
.mobile-menu {
  position: absolute;
  top: 64px;
  left: 0;
  right: 0;
  background: var(--bg-glass);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--border);
  padding: 16px 24px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.mobile-nav-link {
  padding: 12px 16px;
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 15px;
  font-weight: 500;
  border-radius: var(--radius-sm);
  transition: all var(--transition-fast);
}

.mobile-nav-link:hover,
.mobile-nav-link.active {
  color: var(--primary);
  background: var(--bg-secondary);
}

/* Transitions */
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}

.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

@media (max-width: 768px) {
  .navbar-container {
    padding: 0 16px;
  }
}

/* 通知铃铛 */
.notif-wrapper {
  position: relative;
}

.notif-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  color: var(--text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  position: relative;
  transition: all var(--transition-fast);
}

.notif-btn:hover {
  color: var(--primary);
  background: var(--bg-secondary);
}

.notif-badge {
  position: absolute;
  top: 2px;
  right: 2px;
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  background: var(--color-danger);
  color: white;
  font-size: 10px;
  font-weight: 600;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}

.notif-dropdown {
  position: absolute;
  top: calc(100% + 8px);
  right: -60px;
  width: 340px;
  max-height: 420px;
  background: var(--bg-card);
  border-radius: var(--radius-md);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  border: 1px solid var(--border);
  z-index: 200;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.notif-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
}

.notif-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-primary);
}

.notif-read-all {
  font-size: 0.75rem;
  color: var(--primary-dark);
  background: none;
  border: none;
  cursor: pointer;
  padding: 2px 8px;
  border-radius: var(--radius-sm);
  transition: all 0.2s;
}

.notif-read-all:hover {
  background: var(--bg-secondary);
}

.notif-divider {
  height: 1px;
  background: var(--border);
}

.notif-list {
  overflow-y: auto;
  max-height: 360px;
}

.notif-empty {
  text-align: center;
  padding: 32px 16px;
  color: var(--text-muted);
  font-size: 0.8125rem;
}

.notif-item {
  display: flex;
  gap: 10px;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.15s;
  border-bottom: 1px solid var(--border);
}

.notif-item:last-child {
  border-bottom: none;
}

.notif-item:hover {
  background: var(--bg-secondary);
}

.notif-item.unread {
  background: rgba(126, 215, 193, 0.06);
}

.notif-item.unread::before {
  content: '';
  display: block;
  width: 6px;
  height: 6px;
  background: var(--primary);
  border-radius: 50%;
  flex-shrink: 0;
  margin-top: 6px;
}

.notif-icon {
  font-size: 18px;
  flex-shrink: 0;
  margin-top: 2px;
}

.notif-body {
  flex: 1;
  min-width: 0;
}

.notif-text {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--text-primary);
  margin: 0;
  line-height: 1.4;
}

.notif-content {
  font-size: 0.75rem;
  color: var(--text-secondary);
  margin: 4px 0 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.notif-time {
  font-size: 0.6875rem;
  color: var(--text-muted);
}

@media (max-width: 768px) {
  .notif-dropdown {
    right: -20px;
    width: 300px;
  }
}
</style>