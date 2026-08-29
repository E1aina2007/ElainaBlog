// tools.ts 运维工具相关 API
import request from './request'

// 管理员：清空 Redis 缓存
export function clearCache(): Promise<void> {
  return request.post('/cache/clear')
}

// 管理员：保存自定义页头/页脚代码
export function saveCustomCode(data: { custom_header: string; custom_footer: string }): Promise<void> {
  return request.post('/site', data)
}
