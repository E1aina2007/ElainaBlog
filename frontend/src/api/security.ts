// security.ts 安全管理相关 API（IP 封禁、备份）
import request from './request'

// 后端 /security/banned-ips 返回的封禁条目；ttl_seconds 为 -1 表示永久封禁
export interface BannedIPInfo {
  ip: string
  banned_at: number
  ttl_seconds: number
}

// 管理员：获取封禁 IP 列表
export function getBannedIPs(): Promise<BannedIPInfo[]> {
  return request.get('/security/banned-ips')
}

// 管理员：手动封禁 IP（永久，可解封）
export function banIP(ip: string): Promise<void> {
  return request.post('/security/ban', { ip })
}

// 管理员：解封 IP
export function unbanIP(ip: string): Promise<void> {
  return request.post('/security/unban', { ip })
}

// 管理员：导出全站数据库备份
export function exportBackup(): Promise<Blob> {
  return request.get('/backup/export', { responseType: 'blob' })
}
