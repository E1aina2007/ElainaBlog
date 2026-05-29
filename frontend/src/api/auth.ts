// auth.ts 登录和注册相关 API 请求封装
import request from './request'

export interface LoginParams {
    email: string
    password: string
}

export interface LoginResult {
    user_id: number
    email: string
}

interface RegisterParams {
    username: string
    email: string
    password: string
    code: string
}

export function login(params: LoginParams): Promise<LoginResult> {
    return request.post('/login', params)
}

export function register(params: RegisterParams): Promise<{ user_id: number }> {
    return request.post('/register', params)
}

export function sendCode(email: string): Promise<void> {
    return request.post('/send-code', { email })
}

export function logout(): Promise<void> {
    return request.post('/logout')
}

export function resetPassword(email: string, code: string, newPassword: string): Promise<void> {
    return request.post('/reset-password', { email, code, new_password: newPassword })
}