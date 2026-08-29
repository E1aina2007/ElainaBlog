// dashboard.ts 仪表盘与系统状态相关 API
import request from './request'

// 管理员：仪表盘统计数据（视图按需合并字段）
export function getDashboardStats() {
  return request.get('/dashboard/stats')
}

// 管理员：系统运行状态（CPU / 内存等）
export function getSystemStatus() {
  return request.get('/system/status')
}
