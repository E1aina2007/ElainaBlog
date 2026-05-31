// friendlink.ts 友情链接相关 API 请求封装
import request from './request'

export interface FriendLink {
  id: number
  name: string
  url: string
  avatar: string
  description: string
  sort_order: number
}

// 获取友情链接列表（公开）
export function getFriendLinkList(): Promise<FriendLink[]> {
  return request.get('/friend-link/list')
}

// 创建友情链接（管理员）
export function createFriendLink(data: {
  name: string
  url: string
  avatar?: string
  description?: string
  sort_order?: number
}): Promise<FriendLink> {
  return request.post('/friend-link/create', data)
}

// 更新友情链接（管理员）
export function updateFriendLink(data: {
  id: number
  name: string
  url: string
  avatar?: string
  description?: string
  sort_order?: number
}): Promise<FriendLink> {
  return request.post('/friend-link/update', data)
}

// 删除友情链接（管理员）
export function deleteFriendLink(id: number): Promise<void> {
  return request.post('/friend-link/delete', { id })
}
