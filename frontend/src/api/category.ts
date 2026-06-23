// category.ts 分类相关 API 请求封装
import request from './request'

export interface Category {
    id: number
    name: string
    is_top: boolean
    article_count: number
    created_at?: string
    updated_at?: string
}

// 获取分类列表
export function getCategoryList(): Promise<Category[]> {
    return request.get('/category/list')
}

// 创建分类（管理员）
export function createCategory(name: string): Promise<Category> {
    return request.post('/category/create', { name })
}

// 更新分类（管理员）
export function updateCategory(id: number, name: string): Promise<Category> {
    return request.post('/category/update', { id, name })
}

// 删除分类（管理员）
export function deleteCategory(id: number): Promise<void> {
    return request.post('/category/delete', { id })
}

// 切换分类置顶状态（管理员）
export function toggleCategoryTop(id: number, isTop: boolean): Promise<void> {
    return request.post('/category/toggle-top', { id, is_top: isTop })
}
