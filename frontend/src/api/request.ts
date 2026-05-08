// request.ts 封装 axios 实例，配置 baseURL、请求/响应拦截器、自动注入 token
import axios from 'axios'
import { getAccessToken, getRefreshToken, setAccessToken, clearTokens } from '@/utils/auth'

const request = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '/api/ui',
    timeout: 10000,
})

// 请求拦截器：自动注入 Authorization
request.interceptors.request.use((config) => {
    const token = getAccessToken()
    if (token) {
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
})

// 响应拦截器：统一处理错误，401 时尝试刷新 token
let isRefreshing = false
let pendingRequests: ((token: string) => void)[] = []

request.interceptors.response.use(
    (response) => {
        const data = response.data
        // 后端返回 { code: 0, data: ..., message: ... }
        if (data.code === 0) {
            return data.data
        }
        return Promise.reject(new Error(data.message || '请求失败'))
    },
    async (error) => {
        const originalRequest = error.config

        // 401 且未重试过 → 尝试刷新 token
        if (error.response?.status === 401 && !originalRequest._retry) {
            // 排除登录相关请求，让调用方处理错误
            const isAuthRequest = originalRequest.url?.includes('/login') ||
                originalRequest.url?.includes('/register') ||
                originalRequest.url?.includes('/send-code')
            if (isAuthRequest) {
                // 从后端响应中提取错误消息并转为中文
                const backendMessage = error.response?.data?.message || ''
                const errorMap: Record<string, string> = {
                    'unauthorized': '邮箱或密码错误',
                    'invalid credentials': '邮箱或密码错误',
                    'user not found': '用户不存在',
                    'password mismatch': '邮箱或密码错误',
                }
                const message = errorMap[String(backendMessage).toLowerCase()] ||
                    (backendMessage || '邮箱或密码错误')
                return Promise.reject(new Error(message))
            }

            const refreshToken = getRefreshToken()
            if (!refreshToken) {
                clearTokens()
                window.location.href = '/login'
                return Promise.reject(new Error('登录已过期，请重新登录'))
            }

            if (isRefreshing) {
                // 正在刷新，排队等待
                return new Promise((resolve) => {
                    pendingRequests.push((token: string) => {
                        originalRequest.headers.Authorization = `Bearer ${token}`
                        resolve(request(originalRequest))
                    })
                })
            }

            originalRequest._retry = true
            isRefreshing = true

            try {
                const res = await axios.post(
                    `${import.meta.env.VITE_API_BASE_URL || '/api/ui'}/refresh`,
                    { refresh_token: refreshToken },
                )
                const newAccessToken = res.data?.data?.access_token
                if (newAccessToken) {
                    setAccessToken(newAccessToken)
                    originalRequest.headers.Authorization = `Bearer ${newAccessToken}`
                    pendingRequests.forEach((cb) => cb(newAccessToken))
                    pendingRequests = []
                    return request(originalRequest)
                }
                throw new Error('刷新 token 失败')
            } catch {
                clearTokens()
                window.location.href = '/login'
                return Promise.reject(new Error('登录已过期，请重新登录'))
            } finally {
                isRefreshing = false
            }
        }

        // 从 error 中提取原始错误消息
        // 优先读取后端响应体中的 message（401 错误时）
        let rawMessage = '网络错误'
        if (error.response?.data?.message) {
            rawMessage = error.response.data.message
        } else if (error.response?.data) {
            // 有些后端直接返回 data 字符串
            rawMessage = typeof error.response.data === 'string'
                ? error.response.data
                : JSON.stringify(error.response.data)
        } else if (error.message) {
            rawMessage = error.message
        }

        // 常见错误消息映射为中文
        const errorMap: Record<string, string> = {
            'invalid credentials': '邮箱或密码错误',
            'user not found': '用户不存在',
            'password incorrect': '密码错误',
            'account disabled': '账号已被禁用',
            'email not verified': '邮箱未验证',
            'too many attempts': '登录次数过多，请稍后再试',
            'network error': '网络错误，请检查网络连接',
            'timeout': '请求超时，请重试',
            'unauthorized': '未授权，请重新登录',
            'request failed with status code 401': '邮箱或密码错误',
        }
        const lowerMessage = String(rawMessage).toLowerCase()
        const message = errorMap[lowerMessage] || rawMessage
        return Promise.reject(new Error(message))
    },
)

export default request