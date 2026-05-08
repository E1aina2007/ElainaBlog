// article.ts 文章相关 API 请求封装 (列表、详情、创建、更新、删除)
import request from './request'

export interface Article {
    id: number
    title: string
    summary: string
    content: string
    cover?: string
    category_id?: number
    category_name?: string
    user_id?: number
    author_name?: string
    author_avatar?: string
    view_count?: number
    comment_count?: number
    is_top?: boolean
    is_draft?: boolean
    created_at?: string
    updated_at?: string
}

interface ArticleListParams {
    page?: number
    pageSize?: number
    categoryId?: number
}

interface ArticleListResult {
    list: Article[]
    total: number
}

// 获取文章列表
export function getArticleList(params?: ArticleListParams): Promise<ArticleListResult> {
    return request.get('/article/list', { params })
}

// 获取文章详情
export function getArticleDetail(id: number): Promise<Article> {
    return request.get(`/article/${id}`)
}

// 创建文章（管理员）
export function createArticle(data: Omit<Article, 'id' | 'created_at' | 'updated_at'>): Promise<{ id: number }> {
    return request.post('/article/create', data)
}

// 更新文章（管理员）
export function updateArticle(data: Partial<Article> & { id: number }): Promise<void> {
    return request.post('/article/update', data)
}

// 删除文章（管理员）
export function deleteArticle(id: number): Promise<void> {
    return request.post('/article/delete', { id })
}