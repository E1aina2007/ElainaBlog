// article.ts 文章相关 API 请求封装 (列表、详情、创建、更新、删除)
import request from './request'

export interface Article {
    id: number
    title: string
    summary: string
    content: string
    category_id?: number
    category_name?: string
    user_id?: number
    author_name?: string
    author_avatar?: string
    author_is_admin?: boolean
    view_count?: number
    uv_count?: number
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
    sortBy?: 'latest' | 'popular'
}

interface ArticleListResult {
    list: Article[]
    total: number
}

// 获取文章列表（公开，不含草稿）
export function getArticleList(params?: ArticleListParams): Promise<ArticleListResult> {
    return request.get('/article/list', { params })
}

// 获取文章列表（管理员，含草稿）
export function getAdminArticleList(params?: ArticleListParams): Promise<ArticleListResult> {
    return request.get('/article/admin-list', { params })
}

// 获取当前用户的文章列表（含草稿）
export function getMyArticleList(params?: ArticleListParams): Promise<ArticleListResult> {
    return request.get('/article/mine', { params })
}

// 获取文章详情（公开，不含草稿）
export function getArticleDetail(id: number): Promise<Article> {
    return request.get(`/article/${id}`)
}

// 获取文章详情（管理员，含草稿）
export function getAdminArticleDetail(id: number): Promise<Article> {
    return request.get(`/article/admin/${id}`)
}

// 获取自己的文章详情（含草稿，仅限自己的文章）
export function getMyArticleDetail(id: number): Promise<Article> {
    return request.get(`/article/mine/${id}`)
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

// 全文搜索文章
export function searchArticles(keyword: string, page = 1, pageSize = 10): Promise<ArticleListResult> {
    return request.get('/article/search', { params: { keyword, page, pageSize } })
}