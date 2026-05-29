// user.ts Pinia 用户状态管理，存储登录态
// Token 通过 HttpOnly Cookie 管理，前端不存储 token
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, logout as logoutApi, type LoginParams, type LoginResult } from '@/api/auth'
import { getProfile, type UserProfile } from '@/api/user'

export const useUserStore = defineStore('user', () => {
    const userInfo = ref<UserProfile | null>(null)
    const isLoading = ref(false)

    const isLoggedIn = computed(() => !!userInfo.value)
    const isAdmin = computed(() => userInfo.value?.is_admin === true)

    // 登录
    async function login(params: LoginParams): Promise<LoginResult> {
        isLoading.value = true
        try {
            const result = await loginApi(params)
            // Token 通过 HttpOnly Cookie 设置，无需前端存储
            // 登录成功后获取用户信息
            await fetchProfile()
            return result
        } finally {
            isLoading.value = false
        }
    }

    // 获取用户信息
    async function fetchProfile(): Promise<void> {
        try {
            const profile = await getProfile()
            userInfo.value = profile
        } catch {
            userInfo.value = null
        }
    }

    // 退出登录
    async function logout(): Promise<void> {
        try {
            await logoutApi()
        } catch {
            // 即使后端注销失败，也要清理本地状态
        }
        userInfo.value = null
    }

    // 初始化：尝试恢复用户信息（cookie 自动携带）
    async function init(): Promise<void> {
        await fetchProfile()
    }

    return {
        userInfo,
        isLoading,
        isLoggedIn,
        isAdmin,
        login,
        fetchProfile,
        logout,
        init,
    }
})
