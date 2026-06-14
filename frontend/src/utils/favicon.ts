/**
 * 从 URL 中提取域名
 */
export function extractDomain(url: string): string {
  try {
    const normalized = url.startsWith('http') ? url : `https://${url}`
    return new URL(normalized).hostname
  } catch {
    return ''
  }
}

/**
 * 获取网站的 favicon URL
 * 优先使用 Google favicon 服务，回退到直接构造的 favicon 路径
 */
export function getFaviconUrl(url: string): string {
  const domain = extractDomain(url)
  if (!domain) return ''

  // 使用 Google favicon 服务（最可靠）
  return `https://www.google.com/s2/favicons?domain=${domain}&sz=64`
}

/**
 * 获取备用 favicon URL（直接访问网站根目录的 favicon.ico）
 */
export function getFallbackFaviconUrl(url: string): string {
  const domain = extractDomain(url)
  if (!domain) return ''

  const normalized = url.startsWith('http') ? url : `https://${url}`
  return `${new URL(normalized).origin}/favicon.ico`
}
