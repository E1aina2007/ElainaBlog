// message.ts 留言板相关 API 请求封装
import request from './request'

export interface Message {
    id: number
    user_id: number
    username: string
    avatar: string
    is_admin: boolean
    content: string
    created_at: string
}

export interface AuthorStats {
    article_count: number
    comment_count: number
    total_views: number
    days_since: number
}

// 获取留言列表
export function getMessageList(): Promise<Message[]> {
    return request.get('/message/list')
}

// 发表留言
export function createMessage(content: string): Promise<{ id: number }> {
    return request.post('/message/create', { content })
}

// 删除留言
export function deleteMessage(id: number): Promise<void> {
    return request.post('/message/delete', { id })
}

// 获取作者页统计
export function getAuthorStats(): Promise<AuthorStats> {
    return request.get('/author/stats')
}
