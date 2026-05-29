<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useSiteStore } from '@/stores/site'
import { useTheme } from '@/composables/useTheme'

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
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
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
</style>