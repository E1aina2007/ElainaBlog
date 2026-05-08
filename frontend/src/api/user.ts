// user.ts 用户信息相关 API 请求封装
import request from './request'

export interface UserProfile {
    id: number
    username: string
    email: string
    avatar: string
    is_admin: boolean
    created_at: string
    updated_at?: string
}

export function getProfile(): Promise<UserProfile> {
    return request.get('/user/profile')
}

export function updateProfile(data: { username?: string; email?: string; avatar?: string }): Promise<void> {
    return request.post('/user/profile', data)
}

export function updatePassword(data: { old_password: string; new_password: string }): Promise<void> {
    return request.post('/user/password', data)
}

// 管理员：获取用户列表
export function getUserList(): Promise<UserProfile[]> {
    return request.get('/user/list')
}

// 管理员：删除用户
export function deleteUser(userId: number): Promise<void> {
    return request.post('/user/delete', { user_id: userId })
}