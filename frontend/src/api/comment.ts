// comment.ts 评论相关 API 请求封装
import request from './request'

export interface Comment {
    id: number
    article_id: number
    user_id: number
    username: string
    avatar: string
    content: string
    created_at: string
}

// 获取文章评论列表
export function getComments(articleId: number): Promise<Comment[]> {
    return request.get(`/comment/${articleId}`)
}

// 创建评论
export function createComment(data: { article_id: number; content: string }): Promise<{ id: number }> {
    return request.post('/comment/create', data)
}

// 删除评论
export function deleteComment(id: number): Promise<void> {
    return request.post('/comment/delete', { id })
}