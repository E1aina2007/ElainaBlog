import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useSiteStore } from '@/stores/site'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/Home.vue'),
    },
    {
      path: '/article/:id',
      name: 'article-detail',
      component: () => import('../views/ArticleDetail.vue'),
    },
    {
      path: '/author',
      name: 'author',
      component: () => import('../views/Author.vue'),
    },
    {
      path: '/tools',
      name: 'tools',
      component: () => import('../views/Tools.vue'),
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/Profile.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/write',
      name: 'write',
      component: () => import('../views/Write.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/write/:id',
      name: 'write-edit',
      component: () => import('../views/Write.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/Login.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/Register.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/admin',
      component: () => import('../layouts/AdminLayout.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        {
          path: '',
          name: 'admin-dashboard',
          component: () => import('../views/admin/Dashboard.vue'),
        },
        {
          path: 'articles',
          name: 'admin-articles',
          component: () => import('../views/admin/ArticleList.vue'),
        },
        {
          path: 'article/create',
          name: 'admin-article-create',
          component: () => import('../views/admin/ArticleEdit.vue'),
        },
        {
          path: 'article/edit/:id',
          name: 'admin-article-edit',
          component: () => import('../views/admin/ArticleEdit.vue'),
        },
        {
          path: 'comments',
          name: 'admin-comments',
          component: () => import('../views/admin/CommentList.vue'),
        },
        {
          path: 'categories',
          name: 'admin-categories',
          component: () => import('../views/admin/CategoryList.vue'),
        },
        {
          path: 'system',
          name: 'admin-system',
          component: () => import('../views/admin/SystemStatus.vue'),
        },
        {
          path: 'security',
          name: 'admin-security',
          component: () => import('../views/admin/SecurityPanel.vue'),
        },
        {
          path: 'tools',
          name: 'admin-tools',
          component: () => import('../views/admin/ToolsPanel.vue'),
        },
        {
          path: 'profile',
          name: 'admin-profile',
          component: () => import('../views/admin/Profile.vue'),
        },
        {
          path: 'users',
          name: 'admin-users',
          component: () => import('../views/admin/UserList.vue'),
        },
        {
          path: 'site-config',
          name: 'admin-site-config',
          component: () => import('../views/admin/SiteConfig.vue'),
        },
        {
          path: 'author-profile',
          name: 'admin-author-profile',
          component: () => import('../views/admin/AuthorProfile.vue'),
        },
      ],
    },
  ],
  scrollBehavior() {
    return { top: 0 }
  },
})

// 路由守卫
let initialized = false

router.beforeEach(async (to, _from, next) => {
  const userStore = useUserStore()
  const siteStore = useSiteStore()

  // 仅首次加载时尝试恢复登录态和站点配置，失败不阻塞导航
  if (!initialized) {
    initialized = true
    try {
      await Promise.all([
        userStore.init(),
        siteStore.fetchConfig(),
      ])
    } catch {
      // token 无效或后端不可达时忽略，不影响页面访问
    }
  }

  // 需要登录的页面
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({ name: 'login', query: { redirect: to.fullPath } })
    return
  }

  // 需要管理员权限的页面
  if (to.meta.requiresAdmin && !userStore.isAdmin) {
    next({ name: 'home' })
    return
  }

  // 已登录用户不能访问登录/注册页
  if (to.meta.guestOnly && userStore.isLoggedIn) {
    next({ name: 'home' })
    return
  }

  next()
})

export default router
