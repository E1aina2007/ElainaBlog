<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getAllSiteConfig, updateSiteConfig } from '@/api/siteConfig'
import { useSiteStore } from '@/stores/site'
import toast from '@/utils/toast'

const siteStore = useSiteStore()
const configs = ref<Record<string, string>>({})
const loading = ref(false)
const saving = ref(false)

// 可编辑的配置项（排除 quotes，用专门的编辑区域）
const editableKeys = ['site_name', 'site_title', 'greeting', 'hero_title', 'icp_beian', 'gov_police_record']
const quotesText = ref('')
const configLabels: Record<string, string> = {
  site_name: '站点名称',
  site_title: '站点标题',
  greeting: '首页问候语',
  hero_title: '首页标题',
  icp_beian: 'ICP 备案号',
  gov_police_record: '公安备案号',
  quotes: '随机语句（JSON 数组）',
}

const fetchConfigs = async () => {
  loading.value = true
  try {
    configs.value = await getAllSiteConfig()
    quotesText.value = configs.value.quotes || '[]'
  } catch {
    toast.error('获取配置失败')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    // 验证 quotes 是合法 JSON
    try {
      JSON.parse(quotesText.value)
    } catch {
      toast.error('随机语句格式不正确，请输入合法的 JSON 数组')
      saving.value = false
      return
    }
    const toSave: Record<string, string> = {}
    for (const key of editableKeys) {
      toSave[key] = configs.value[key] || ''
    }
    toSave.quotes = quotesText.value
    await updateSiteConfig(toSave)
    // 刷新全局站点配置缓存
    await siteStore.fetchConfig()
    toast.success('保存成功')
  } catch {
    toast.error('保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(fetchConfigs)
</script>

<template>
  <div class="site-config-page">
    <div class="page-header">
      <h2>站点配置</h2>
      <button class="btn-primary" :disabled="saving" @click="handleSave">
        {{ saving ? '保存中...' : '保存配置' }}
      </button>
    </div>

    <div v-if="loading" class="loading-state">加载中...</div>

    <div v-else class="config-form">
      <div v-for="key in editableKeys" :key="key" class="config-item">
        <label class="config-label">{{ configLabels[key] || key }}</label>
        <input
          v-model="configs[key]"
          class="config-input"
          :placeholder="'请输入' + (configLabels[key] || key)"
        />
      </div>

      <div class="config-item">
        <label class="config-label">{{ configLabels.quotes }}</label>
        <textarea
          v-model="quotesText"
          class="config-textarea"
          rows="6"
          placeholder='["句子1","句子2","句子3"]'
          v-tab-indent
        ></textarea>
        <p class="config-hint">JSON 数组格式，每个元素为一条随机语句</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.site-config-page {
  max-width: 800px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.page-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.btn-primary {
  padding: 8px 20px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover:not(:disabled) {
  background: var(--primary-dark);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-state {
  text-align: center;
  padding: 40px;
  color: var(--text-muted);
}

.config-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.config-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.config-label {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.config-input {
  padding: 10px 14px;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  font-size: 14px;
  color: var(--text-primary);
  background: var(--bg-card);
  outline: none;
  transition: border-color 0.2s;
}

.config-input:focus {
  border-color: var(--primary);
}

.config-textarea {
  padding: 10px 14px;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  font-size: 14px;
  color: var(--text-primary);
  background: var(--bg-card);
  outline: none;
  resize: vertical;
  font-family: 'Courier New', monospace;
  transition: border-color 0.2s;
}

.config-textarea:focus {
  border-color: var(--primary);
}

.config-hint {
  font-size: 12px;
  color: var(--text-muted);
  margin: 0;
}
</style>
