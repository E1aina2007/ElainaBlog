// site.ts Pinia 站点配置状态管理
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getSiteConfig } from '@/api/siteConfig'

// 默认兜底配置
const DEFAULT_CONFIG: Record<string, string> = {
  site_name: '网站名',
  site_title: '网站标题',
  greeting: '问候语1',
  hero_title: '问候语2',
  icp_beian: '备案号',
  quotes: '["首页语句/签名1", "首页语句/签名2", "首页语句/签名3"]',
}

export const useSiteStore = defineStore('site', () => {
  const config = ref<Record<string, string>>({ ...DEFAULT_CONFIG })
  const loaded = ref(false)

  async function fetchConfig() {
    try {
      const data = await getSiteConfig()
      // 合并：后端返回的配置覆盖默认值
      config.value = { ...DEFAULT_CONFIG, ...data }
      loaded.value = true
    } catch {
      // 接口失败时使用默认配置，不做任何报错
      config.value = { ...DEFAULT_CONFIG }
      loaded.value = true
    }
    // 动态更新页面标题
    const title = config.value.site_name || DEFAULT_CONFIG.site_name
    if (title) {
      document.title = title
    }
  }

  function get(key: string): string {
    return config.value[key] ?? DEFAULT_CONFIG[key] ?? ''
  }

  function getQuotes(): string[] {
    try {
      const raw: string = config.value.quotes || DEFAULT_CONFIG.quotes || '[]'
      return JSON.parse(raw)
    } catch {
      return JSON.parse(DEFAULT_CONFIG.quotes || '[]')
    }
  }

  return { config, loaded, fetchConfig, get, getQuotes }
})
