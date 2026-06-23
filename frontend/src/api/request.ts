// request.ts 封装 axios 实例，配置 baseURL、响应拦截器
// Token 通过 HttpOnly Cookie 自动携带，无需手动注入 Authorization 头
import axios from 'axios'

const request = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '/api/ui',
    timeout: 10000,
    withCredentials: true,
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

// 刷新队列：防止多个并发 401 请求同时触发 refresh
let refreshPromise: Promise<void> | null = null
let pendingRequests: (() => void)[] = []

function doRefresh(): Promise<void> {
    const p = axios.post(
        `${import.meta.env.VITE_API_BASE_URL || '/api/ui'}/refresh`,
        {},
        { withCredentials: true },
    ).then(() => {
        // 刷新成功，重试所有排队的请求
        pendingRequests.forEach((cb) => cb())
        pendingRequests = []
    }).catch(() => {
        // 刷新失败，清空队列，不重定向（由路由守卫处理）
        pendingRequests = []
        throw new Error('登录已过期，请重新登录')
    }).finally(() => {
        refreshPromise = null
    })
    refreshPromise = p
    return p
}

// 响应拦截器：统一处理错误，401 时尝试刷新 token（cookie 自动携带）
request.interceptors.response.use(
    (response) => {
        const data = response.data
        // blob 响应（如文件下载）直接透传
        if (data instanceof Blob) {
            return data
        }
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

        // 401 且未重试过 → 尝试刷新 token（cookie 自动携带）
        if (error.response?.status === 401 && !originalRequest._retry) {
            // 排除登录相关请求
            const isAuthRequest = originalRequest.url?.includes('/login') ||
                originalRequest.url?.includes('/register') ||
                originalRequest.url?.includes('/send-code') ||
                originalRequest.url?.includes('/refresh')
            if (isAuthRequest) {
                return Promise.reject(new Error(getErrorMessage(error)))
            }

            originalRequest._retry = true

            // 如果已有刷新请求在进行中，排队等待
            if (refreshPromise) {
                return new Promise<void>((resolve, reject) => {
                    pendingRequests.push(() => {
                        request(originalRequest).then(() => resolve()).catch(reject)
                    })
                })
            }

            try {
                await doRefresh()
                return request(originalRequest)
            } catch {
                return Promise.reject(new Error('登录已过期，请重新登录'))
            }
        }

        return Promise.reject(new Error(getErrorMessage(error)))
    },
)

export default request
