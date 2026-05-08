import request from './request'

interface UploadResult {
  url: string
}

export function uploadAvatar(file: File): Promise<UploadResult> {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/upload/avatar', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}
