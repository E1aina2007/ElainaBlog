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
 * 获取所有候选 favicon URL（按优先级排列）
 * 顺序：站点直接图标 → DuckDuckGo → Google
 */
export function getFaviconFallbacks(url: string): string[] {
  const domain = extractDomain(url)
  if (!domain) return []

  return [
    `https://${domain}/favicon.ico`,
    `https://${domain}/apple-touch-icon.png`,
    `https://icons.duckduckgo.com/ip3/${domain}.ico`,
    `https://www.google.com/s2/favicons?domain=${domain}&sz=64`,
  ]
}

/**
 * 获取主 favicon URL（第一个候选）
 */
export function getFaviconUrl(url: string): string {
  const fallbacks = getFaviconFallbacks(url)
  return fallbacks[0] ?? ''
}

/**
 * 创建 favicon 加载失败时的回退处理器
 * 用法：@error="createFaviconErrorHandler(url)"
 *
 * @param url - 目标网站 URL
 * @param onAllFailed - 所有候选都失败时的回调（可选，默认隐藏图片）
 * @returns onerror 事件处理函数
 */
export function createFaviconErrorHandler(
  url: string,
  onAllFailed?: (e: Event) => void
) {
  return (e: Event) => {
    const img = e.target as HTMLImageElement
    const fallbacks = getFaviconFallbacks(url)
    const currentSrc = img.src
    const currentIndex = fallbacks.indexOf(currentSrc)
    const nextIndex = currentIndex + 1

    const nextSrc = fallbacks[nextIndex]
    if (nextSrc) {
      // 还有候选，尝试下一个
      img.src = nextSrc
    } else {
      // 全部失败
      if (onAllFailed) {
        onAllFailed(e)
      } else {
        img.style.display = 'none'
      }
    }
  }
}
