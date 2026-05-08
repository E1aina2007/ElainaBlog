// siteConfig.ts 站点配置相关 API
import request from './request'

// 获取公开站点配置（key-value 映射）
export function getSiteConfig(): Promise<Record<string, string>> {
  return request.get('/site-config')
}

// 管理员：获取所有配置
export function getAllSiteConfig(): Promise<Record<string, string>> {
  return request.get('/site-config/all')
}

// 管理员：批量更新配置
export function updateSiteConfig(configs: Record<string, string>): Promise<void> {
  return request.post('/site-config/update', { configs })
}

// 管理员：删除配置
export function deleteSiteConfig(key: string): Promise<void> {
  return request.post('/site-config/delete', { key })
}
