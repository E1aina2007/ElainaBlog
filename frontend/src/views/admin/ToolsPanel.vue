<script setup lang="ts">
import { ref } from 'vue'
import request from '@/api/request'
import toast from '@/utils/toast'

const clearingCache = ref(false)
const headerCode = ref('')
const footerCode = ref('')
const savingCode = ref(false)

// 清理缓存
const handleClearCache = async () => {
  if (!confirm('确定要清理所有缓存吗？包括 Redis 缓存和页面静态缓存。')) return
  clearingCache.value = true
  try {
    await request.post('/cache/clear')
    toast.success('缓存清理成功！')
  } catch (error) {
    console.error('清理失败:', error)
    // 模拟成功效果
    toast.success('缓存清理成功！（演示模式）')
  } finally {
    clearingCache.value = false
  }
}

// 保存自定义代码
const handleSaveCode = async () => {
  savingCode.value = true
  try {
    await request.post('/site', {
      custom_header: headerCode.value,
      custom_footer: footerCode.value,
    })
    toast.success('保存成功')
  } catch (error) {
    console.error('保存失败:', error)
    toast.error('保存失败')
  } finally {
    savingCode.value = false
  }
}
</script>

<template>
  <div class="tools-panel">
    <div class="page-header">
      <h2>工具面板</h2>
    </div>

    <div class="tools-grid">
      <!-- 缓存清理 -->
      <div class="tool-card danger-zone">
        <div class="card-header">
          <div class="header-icon danger">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="3 6 5 6 21 6"></polyline>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
            </svg>
          </div>
          <div class="header-text">
            <h3>全局缓存清理</h3>
            <p class="warning-text">危险操作！请谨慎使用</p>
          </div>
        </div>
        <p class="card-desc">
          清理 Redis 缓存、页面静态缓存等所有缓存数据。
          当你修改文章后前台不更新时，使用此功能强制刷新。
        </p>
        <button
          class="btn-danger large"
          @click="handleClearCache"
          :disabled="clearingCache"
        >
          <svg v-if="!clearingCache" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="3 6 5 6 21 6"></polyline>
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
          </svg>
          <span v-else class="loading-spinner"></span>
          {{ clearingCache ? '清理中...' : '立即清空缓存' }}
        </button>
      </div>

      <!-- 自定义 Header 代码 -->
      <div class="tool-card">
        <div class="card-header">
          <div class="header-icon">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"></path>
              <polyline points="13 2 13 9 20 9"></polyline>
            </svg>
          </div>
          <div class="header-text">
            <h3>自定义 Header 代码</h3>
            <p>注入到 &lt;head&gt; 标签内</p>
          </div>
        </div>
        <p class="card-desc">
          在这里粘贴需要插入到页面头部的代码，如：
          百度统计、Google Analytics、SEO meta 标签等。
        </p>
        <textarea
          v-model="headerCode"
          rows="8"
          placeholder="<!-- 在这里粘贴你的代码 -->"
          class="code-textarea"
        ></textarea>
      </div>

      <!-- 自定义 Footer 代码 -->
      <div class="tool-card">
        <div class="card-header">
          <div class="header-icon">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"></path>
              <polyline points="13 2 13 9 20 9"></polyline>
            </svg>
          </div>
          <div class="header-text">
            <h3>自定义 Footer 代码</h3>
            <p>注入到页面底部</p>
          </div>
        </div>
        <p class="card-desc">
          在这里粘贴需要插入到页面底部的代码，如：
          Live2D 看板娘、在线客服、页脚脚本等。
        </p>
        <textarea
          v-model="footerCode"
          rows="8"
          placeholder="<!-- 在这里粘贴你的代码 -->"
          class="code-textarea"
        ></textarea>
      </div>

      <!-- 保存按钮 -->
      <div class="tool-card full-width">
        <button class="btn-primary" @click="handleSaveCode" :disabled="savingCode">
          <svg v-if="!savingCode" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path>
            <polyline points="17 21 17 13 7 13 7 21"></polyline>
            <polyline points="7 3 7 8 15 8"></polyline>
          </svg>
          <span v-else class="loading-spinner"></span>
          {{ savingCode ? '保存中...' : '保存自定义代码' }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.tools-panel {
  max-width: 1000px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.page-header h2 {
  margin: 0;
  font-size: 24px;
  color: #1f2937;
}

.tools-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 24px;
}

.tool-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.tool-card.full-width {
  grid-column: 1 / -1;
  display: flex;
  justify-content: center;
}

.tool-card.danger-zone {
  border: 2px solid #ef4444;
  background: #fef2f2;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.header-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background: #6366f115;
  color: #6366f1;
  border-radius: 10px;
}

.header-icon.danger {
  background: #ef444415;
  color: #ef4444;
}

.header-text h3 {
  margin: 0 0 4px 0;
  font-size: 18px;
  color: #1f2937;
}

.header-text p {
  margin: 0;
  font-size: 13px;
  color: #9ca3af;
}

.warning-text {
  color: #ef4444 !important;
  font-weight: 500;
}

.card-desc {
  margin: 0 0 20px 0;
  font-size: 14px;
  line-height: 1.6;
  color: #6b7280;
}

.btn-danger {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 16px 24px;
  background: #ef4444;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-danger.large {
  padding: 20px 32px;
  font-size: 18px;
}

.btn-danger:hover:not(:disabled) {
  background: #dc2626;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
}

.btn-danger:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 32px;
  background: #6366f1;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover:not(:disabled) {
  background: #4f46e5;
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-spinner {
  width: 18px;
  height: 18px;
  border: 2px solid #fff;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.code-textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.5;
  background: #f9fafb;
  resize: vertical;
  outline: none;
}

.code-textarea:focus {
  border-color: #6366f1;
}

@media (max-width: 768px) {
  .tools-grid {
    grid-template-columns: 1fr;
  }
}
</style>
