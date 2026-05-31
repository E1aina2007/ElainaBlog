// notification.ts 通知相关 API 请求封装
import request from './request'

export interface Notification {
  id: number
  type: string
  title: string
  content: string
  target_id: number
  is_read: boolean
  created_at: string
}

// 获取通知列表
export function getNotificationList(unread?: boolean): Promise<Notification[]> {
  const params = unread ? '?unread=1' : ''
  return request.get(`/notification/list${params}`)
}

// 获取未读通知数量
export function getUnreadCount(): Promise<{ count: number }> {
  return request.get('/notification/unread-count')
}

// 标记单条通知为已读
export function markAsRead(id: number): Promise<void> {
  return request.post('/notification/read', { id })
}

// 标记所有通知为已读
export function markAllAsRead(): Promise<void> {
  return request.post('/notification/read-all')
}

// 删除通知
export function deleteNotification(id: number): Promise<void> {
  return request.post('/notification/delete', { id })
}
