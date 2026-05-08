// authorProfile.ts 作者信息相关 API
import request from './request'

export interface AuthorProfile {
  nickname: string
  avatar: string
  background: string
  signature: string
  location: string
  occupation: string
  school: string
  major: string
  email: string
  wechat: string
  bio: string
  tech_stack_frontend: string
  tech_stack_backend: string
  tech_stack_engineering: string
  social_github: string
  social_bilibili: string
}

// 获取作者信息（公开）
export function getAuthorProfile(): Promise<AuthorProfile> {
  return request.get('/author/profile')
}

// 管理员：更新作者信息
export function updateAuthorProfile(profile: AuthorProfile): Promise<void> {
  return request.post('/author/profile/update', profile)
}
