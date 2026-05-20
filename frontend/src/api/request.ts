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

// 后端错误码 → 中文提示映射
const errorCodeMap: Record<number, string> = {
    // 通用参数错误
    400001: '请求参数错误',
    // 认证相关
    401001: '未授权，请重新登录',
    401002: '邮箱或密码错误',
    401003: '登录已过期，请重新登录',
    401004: '登录状态无效，请重新登录',
    401005: '刷新令牌无效，请重新登录',
    // 权限
    403001: '没有权限执行此操作',
    // 资源不存在
    404001: '资源不存在',
    404101: '分类不存在',
    404102: '文章不存在',
    404103: '评论不存在',
    404104: '留言不存在',
    // 冲突
    409001: '资源冲突',
    409101: '用户名已存在',
    409102: '邮箱已被注册',
    409103: '分类已存在',
    // 用户相关
    400101: '用户不存在',
    400102: '邮箱格式不正确',
    400103: '用户名格式不正确',
    400104: '密码长度不符合要求',
    400105: '密码必须包含字母和数字',
    // 验证码
    400201: '验证码已过期，请重新获取',
    400202: '验证码错误',
    429201: '验证码发送过于频繁，请稍后再试',
    // 上传
    400301: '文件大小超出限制',
    400302: '不支持的文件类型',
    // 限流
    429001: '请求过于频繁，请稍后再试',
    // 服务器错误
    500001: '服务器内部错误，请稍后再试',
}

// 从后端错误响应中提取中文提示
function getErrorMessage(error: any): string {
    const status = error.response?.status
    const data = error.response?.data

    // 网络错误 / 超时
    if (!error.response) {
        if (error.code === 'ECONNABORTED') return '请求超时，请重试'
        return '网络错误，请检查网络连接'
    }

    // 根据后端错误码映射
    const code = data?.code
    if (code && errorCodeMap[code]) {
        return errorCodeMap[code]
    }

    // HTTP 状态码兜底
    const statusCodeMap: Record<number, string> = {
        400: '请求参数错误',
        401: '登录已过期，请重新登录',
        403: '没有权限执行此操作',
        404: '请求的资源不存在',
        429: '请求过于频繁，请稍后再试',
        500: '服务器内部错误，请稍后再试',
        502: '服务暂时不可用，请稍后再试',
        503: '服务暂时不可用，请稍后再试',
    }
    if (status && statusCodeMap[status]) {
        return statusCodeMap[status]
    }

    return '请求失败'
}

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
        // 业务错误（HTTP 200 但 code !== 0）
        const message = errorCodeMap[data.code] || data.message || '请求失败'
        return Promise.reject(new Error(message))
    },
    async (error) => {
        const originalRequest = error.config

        // 401 且未重试过 → 尝试刷新 token
        if (error.response?.status === 401 && !originalRequest._retry) {
            // 排除登录相关请求，直接返回错误码映射的中文提示
            const isAuthRequest = originalRequest.url?.includes('/login') ||
                originalRequest.url?.includes('/register') ||
                originalRequest.url?.includes('/send-code')
            if (isAuthRequest) {
                return Promise.reject(new Error(getErrorMessage(error)))
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

        return Promise.reject(new Error(getErrorMessage(error)))
    },
)

export default request