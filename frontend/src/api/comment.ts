// comment.ts 评论相关 API 请求封装
import request from './request'

export interface Comment {
    id: number
    article_id: number
    user_id: number
    reply_to_user_id: number | null
    reply_to_username: string | null
    reply_to_comment_id: number | null
    reply_to_content: string | null
    username: string
    avatar: string
    is_admin: boolean
    content: string
    created_at: string
}

// 获取文章评论列表
export function getComments(articleId: number): Promise<Comment[]> {
    return request.get(`/comment/${articleId}`)
}

// 创建评论
export function createComment(data: { article_id: number; reply_to_user_id?: number | null; reply_to_comment_id?: number | null; content: string }): Promise<{ id: number }> {
    return request.post('/comment/create', data)
}

// 删除评论
export function deleteComment(id: number): Promise<void> {
    return request.post('/comment/delete', { id })
}