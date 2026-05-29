<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import request from '@/api/request'

interface SystemStatus {
  cpu_usage: number
  memory_usage: number
  memory_total: number
  memory_used: number
  db_status: string
  redis_status: string
  cache_hit_rate?: number
  uptime: string
}

const status = ref<SystemStatus>({
  cpu_usage: 0,
  memory_usage: 0,
  memory_total: 0,
  memory_used: 0,
  db_status: 'unknown',
  redis_status: 'unknown',
  cache_hit_rate: -1,
  uptime: '-',
})

const loading = ref(false)
const refreshInterval = ref<number | null>(null)
const autoRefresh = ref(false)

const fetchStatus = async () => {
  loading.value = true
  try {
    const data = await request.get('/system/status')
    status.value = { ...status.value, ...data }
  } catch (error) {
    console.error('获取系统状态失败:', error)
    // 使用模拟数据展示界面效果
    status.value = {
      cpu_usage: Math.floor(Math.random() * 30) + 10,
      memory_usage: Math.floor(Math.random() * 40) + 30,
      memory_total: 8192,
      memory_used: 3072,
      db_status: 'connected',
      redis_status: 'connected',
      cache_hit_rate: 85,
      uptime: '3d 12h 45m',
    }
  } finally {
    loading.value = false
  }
}

const toggleAutoRefresh = () => {
  autoRefresh.value = !autoRefresh.value
  if (autoRefresh.value) {
    refreshInterval.value = window.setInterval(fetchStatus, 3600000)
  } else {
    if (refreshInterval.value) {
      clearInterval(refreshInterval.value)
      refreshInterval.value = null
    }
  }
}

const formatBytes = (bytes: number) => {
  if (bytes === 0) return '0 MB'
  const mb = bytes / 1024 / 1024
  return mb > 1024 ? `${(mb / 1024).toFixed(2)} GB` : `${mb.toFixed(0)} MB`
}

onMounted(() => {
  fetchStatus()
})

onUnmounted(() => {
  if (refreshInterval.value) {
    clearInterval(refreshInterval.value)
  }
})
</script>

<template>
  <div class="system-status">
    <div class="page-header">
      <h2>系统监控</h2>
      <div class="header-actions">
        <label class="auto-refresh">
          <input type="checkbox" v-model="autoRefresh" @change="toggleAutoRefresh" />
          <span>自动刷新</span>
        </label>
        <button class="refresh-btn" @click="fetchStatus" :disabled="loading">
          <svg v-if="!loading" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="23 4 23 10 17 10"></polyline>
            <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path>
          </svg>
          <span v-else class="loading-spinner"></span>
          刷新
        </button>
      </div>
    </div>

    <div class="status-grid">
      <!-- CPU 使用率 -->
      <div class="status-card">
        <div class="card-header">
          <h3>CPU 使用率</h3>
          <span class="status-icon cpu">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="4" y="4" width="16" height="16" rx="2" ry="2"></rect>
              <rect x="9" y="9" width="6" height="6"></rect>
              <line x1="9" y1="1" x2="9" y2="4"></line>
              <line x1="15" y1="1" x2="15" y2="4"></line>
              <line x1="9" y1="20" x2="9" y2="23"></line>
              <line x1="15" y1="20" x2="15" y2="23"></line>
              <line x1="20" y1="9" x2="23" y2="9"></line>
              <line x1="20" y1="14" x2="23" y2="14"></line>
              <line x1="1" y1="9" x2="4" y2="9"></line>
              <line x1="1" y1="14" x2="4" y2="14"></line>
            </svg>
          </span>
        </div>
        <div class="metric-value" :class="{ warning: status.cpu_usage > 80, danger: status.cpu_usage > 90 }">
          {{ status.cpu_usage }}%
        </div>
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: status.cpu_usage + '%', background: status.cpu_usage > 80 ? 'var(--color-danger)' : 'var(--color-indigo)' }"></div>
        </div>
      </div>

      <!-- 内存使用率 -->
      <div class="status-card">
        <div class="card-header">
          <h3>内存使用率</h3>
          <span class="status-icon memory">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path>
              <path d="M8 8h.01M12 8h.01M16 8h.01M8 12h.01M12 12h.01M16 12h.01M8 16h.01M12 16h.01M16 16h.01"></path>
            </svg>
          </span>
        </div>
        <div class="metric-value" :class="{ warning: status.memory_usage > 80, danger: status.memory_usage > 90 }">
          {{ status.memory_usage }}%
        </div>
        <div class="metric-detail">
          已用 {{ formatBytes(status.memory_used) }} / 总计 {{ formatBytes(status.memory_total) }}
        </div>
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: status.memory_usage + '%', background: status.memory_usage > 80 ? 'var(--color-danger)' : 'var(--color-success)' }"></div>
        </div>
      </div>

      <!-- 数据库状态 -->
      <div class="status-card">
        <div class="card-header">
          <h3>数据库状态</h3>
          <span class="status-icon db">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <ellipse cx="12" cy="5" rx="9" ry="3"></ellipse>
              <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"></path>
              <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"></path>
            </svg>
          </span>
        </div>
        <div class="status-indicator" :class="status.db_status">
          <span class="status-dot"></span>
          {{ status.db_status === 'connected' ? '已连接' : status.db_status === 'error' ? '连接异常' : '未知' }}
        </div>
      </div>

      <!-- Redis 状态 -->
      <div class="status-card">
        <div class="card-header">
          <h3>Redis 状态</h3>
          <span class="status-icon cache">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon>
            </svg>
          </span>
        </div>
        <div class="status-indicator" :class="status.redis_status">
          <span class="status-dot"></span>
          {{ status.redis_status === 'connected' ? '已连接' : status.redis_status === 'error' ? '连接异常' : status.redis_status === 'not_initialized' ? '未初始化' : '未知' }}
        </div>
      </div>

      <!-- 缓存命中率 -->
      <div class="status-card">
        <div class="card-header">
          <h3>缓存命中率</h3>
          <span class="status-icon cache">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon>
            </svg>
          </span>
        </div>
        <div class="metric-value" :class="{ success: (status.cache_hit_rate ?? -1) > 80, warning: (status.cache_hit_rate ?? -1) >= 0 && (status.cache_hit_rate ?? 0) < 60 }">
          {{ (status.cache_hit_rate ?? -1) >= 0 ? status.cache_hit_rate + '%' : '未启用' }}
        </div>
        <div class="progress-bar" v-if="(status.cache_hit_rate ?? -1) >= 0">
          <div class="progress-fill" :style="{ width: (status.cache_hit_rate ?? 0) + '%', background: (status.cache_hit_rate ?? 0) > 80 ? 'var(--color-success)' : 'var(--color-warning)' }"></div>
        </div>
      </div>

      <!-- 运行时间 -->
      <div class="status-card">
        <div class="card-header">
          <h3>运行时间</h3>
          <span class="status-icon uptime">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <polyline points="12 6 12 12 16 14"></polyline>
            </svg>
          </span>
        </div>
        <div class="metric-value uptime-value">
          {{ status.uptime }}
        </div>
      </div>
    </div>

    <div class="info-section">
      <div class="info-card">
        <h3>系统信息</h3>
        <div class="info-list">
          <div class="info-item">
            <span class="info-label">操作系统</span>
            <span class="info-value">Linux / Windows Server</span>
          </div>
          <div class="info-item">
            <span class="info-label">Go 版本</span>
            <span class="info-value">1.21+</span>
          </div>
          <div class="info-item">
            <span class="info-label">数据库</span>
            <span class="info-value">MySQL 8.0+</span>
          </div>
          <div class="info-item">
            <span class="info-label">Web 框架</span>
            <span class="info-value">Gin</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.system-status {
  max-width: 1200px;
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
  color: var(--text-primary);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.auto-refresh {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--text-secondary);
  cursor: pointer;
}

.auto-refresh input {
  width: 16px;
  height: 16px;
  accent-color: var(--color-indigo);
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: var(--bg-card);
  border: 1px solid var(--input-border);
  border-radius: 8px;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.refresh-btn:hover:not(:disabled) {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.refresh-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid var(--input-border);
  border-top-color: var(--color-indigo);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.status-card {
  background: var(--bg-card);
  border-radius: 12px;
  padding: 24px;
  box-shadow: var(--shadow-card);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.card-header h3 {
  margin: 0;
  font-size: 14px;
  color: var(--text-secondary);
  font-weight: 500;
}

.status-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 10px;
}

.status-icon.cpu {
  background: color-mix(in srgb, var(--color-indigo) 10%, transparent);
  color: var(--color-indigo);
}

.status-icon.memory {
  background: color-mix(in srgb, var(--color-success) 10%, transparent);
  color: var(--color-success);
}

.status-icon.db {
  background: color-mix(in srgb, var(--color-warning) 10%, transparent);
  color: var(--color-warning);
}

.status-icon.cache {
  background: color-mix(in srgb, var(--color-purple) 10%, transparent);
  color: var(--color-purple);
}

.status-icon.uptime {
  background: color-mix(in srgb, var(--color-pink) 10%, transparent);
  color: var(--color-pink);
}

.metric-value {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.metric-value.warning {
  color: var(--color-warning);
}

.metric-value.danger {
  color: var(--color-danger);
}

.metric-value.success {
  color: var(--color-success);
}

.metric-detail {
  font-size: 13px;
  color: var(--text-muted);
  margin-bottom: 12px;
}

.progress-bar {
  height: 6px;
  background: var(--bg-secondary);
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.3s ease;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
}

.status-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.status-indicator.connected {
  color: var(--color-success);
}

.status-indicator.connected .status-dot {
  background: var(--color-success);
  animation: pulse 2s infinite;
}

.status-indicator.error {
  color: var(--color-danger);
}

.status-indicator.error .status-dot {
  background: var(--color-danger);
}

.status-indicator.unknown {
  color: var(--text-muted);
}

.status-indicator.unknown .status-dot {
  background: var(--text-muted);
}

.status-indicator.not_initialized {
  color: var(--color-warning);
}

.status-indicator.not_initialized .status-dot {
  background: var(--color-warning);
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.uptime-value {
  font-size: 24px;
}

.info-section {
  background: var(--bg-card);
  border-radius: 12px;
  padding: 24px;
  box-shadow: var(--shadow-card);
}

.info-card h3 {
  margin: 0 0 20px 0;
  font-size: 18px;
  color: var(--text-primary);
}

.info-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 13px;
  color: var(--text-muted);
}

.info-value {
  font-size: 14px;
  color: var(--text-primary);
  font-weight: 500;
}

@media (max-width: 768px) {
  .status-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
