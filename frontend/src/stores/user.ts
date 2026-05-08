// user.ts Pinia 用户状态管理，存储登录态和 token
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, type LoginParams, type LoginResult } from '@/api/auth'
import { getProfile, type UserProfile } from '@/api/user'
import { setAccessToken, setRefreshToken, getAccessToken, clearTokens } from '@/utils/auth'

export const useUserStore = defineStore('user', () => {
    const userInfo = ref<UserProfile | null>(null)
    const isLoading = ref(false)
    // 使用响应式 token 变量，确保登录状态正确更新
    const accessToken = ref<string | null>(getAccessToken())

    const isLoggedIn = computed(() => !!accessToken.value && !!userInfo.value)
    const isAdmin = computed(() => userInfo.value?.is_admin === true)

    // 登录
    async function login(params: LoginParams): Promise<LoginResult> {
        isLoading.value = true
        try {
            const result = await loginApi(params)
            setAccessToken(result.access_token)
            setRefreshToken(result.refresh_token)
            // 更新响应式 token，触发登录状态更新
            accessToken.value = result.access_token
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
    function logout(): void {
        clearTokens()
        accessToken.value = null
        userInfo.value = null
    }

    // 初始化：如果本地有 token，尝试恢复用户信息
    async function init(): Promise<void> {
        const token = getAccessToken()
        if (token) {
            accessToken.value = token
            await fetchProfile()
        }
    }

    return {
        userInfo,
        isLoading,
        isLoggedIn,
        isAdmin,
        accessToken,
        login,
        fetchProfile,
        logout,
        init,
    }
})