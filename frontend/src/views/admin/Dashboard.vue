<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/api/request'

interface DashboardStats {
  article_count: number
  comment_count: number
  user_count: number
  pending_comments?: number
  today_pv?: number
  today_uv?: number
  yesterday_pv?: number
  yesterday_uv?: number
}

const stats = ref<DashboardStats>({
  article_count: 0,
  comment_count: 0,
  user_count: 0,
  pending_comments: 0,
  today_pv: 0,
  today_uv: 0,
  yesterday_pv: 0,
  yesterday_uv: 0,
})

const loading = ref(false)

interface SystemStatus {
  cpu_usage: number
  memory_usage: number
  memory_total: number
  memory_used: number
  db_status: string
  cache_hit_rate: number
  uptime: string
}

const sysStatus = ref<SystemStatus | null>(null)

const fetchStats = async () => {
  loading.value = true
  try {
    const data = await request.get('/dashboard/stats')
    stats.value = { ...stats.value, ...data }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchSystemStatus = async () => {
  try {
    const data = await request.get('/system/status')
    sysStatus.value = data
  } catch (error) {
    console.error('获取系统状态失败:', error)
  }
}

const cards = [
  { key: 'article_count', title: '文章总数', icon: 'article', color: '#6366f1' },
  { key: 'pending_comments', title: '待审评论', icon: 'comment', color: '#f59e0b' },
  { key: 'today_pv', title: '今日访问量(PV)', icon: 'eye', color: '#10b981' },
  { key: 'today_uv', title: '今日访客(UV)', icon: 'user', color: '#ec4899' },
  { key: 'yesterday_pv', title: '昨日访问量(PV)', icon: 'eye', color: '#8b5cf6' },
  { key: 'yesterday_uv', title: '昨日访客(UV)', icon: 'user', color: '#06b6d4' },
  { key: 'user_count', title: '注册用户', icon: 'users', color: '#14b8a6' },
  { key: 'comment_count', title: '评论总数', icon: 'message', color: '#f97316' },
]

onMounted(() => {
  fetchStats()
  fetchSystemStatus()
})
</script>

<template>
  <div class="dashboard">
    <div class="page-header">
      <h2>数据看板</h2>
      <button class="refresh-btn" @click="fetchStats" :disabled="loading">
        <svg v-if="!loading" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="23 4 23 10 17 10"></polyline>
          <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path>
        </svg>
        <span v-else class="loading-spinner"></span>
        刷新数据
      </button>
    </div>

    <div class="stats-grid">
      <div
        v-for="card in cards"
        :key="card.key"
        class="stat-card"
        :style="{ borderLeftColor: card.color }"
      >
        <div class="card-icon" :style="{ background: card.color + '15', color: card.color }">
          <svg v-if="card.icon === 'article'" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
            <polyline points="14 2 14 8 20 8"></polyline>
            <line x1="16" y1="13" x2="8" y2="13"></line>
            <line x1="16" y1="17" x2="8" y2="17"></line>
          </svg>
          <svg v-else-if="card.icon === 'comment'" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"></path>
          </svg>
          <svg v-else-if="card.icon === 'eye'" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
            <circle cx="12" cy="12" r="3"></circle>
          </svg>
          <svg v-else-if="card.icon === 'user'" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
            <circle cx="12" cy="7" r="4"></circle>
          </svg>
          <svg v-else-if="card.icon === 'users'" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
            <circle cx="9" cy="7" r="4"></circle>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
          </svg>
          <svg v-else-if="card.icon === 'message'" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
          </svg>
        </div>
        <div class="card-content">
          <div class="card-title">{{ card.title }}</div>
          <div class="card-value">{{ stats[card.key as keyof DashboardStats] ?? 0 }}</div>
        </div>
      </div>
    </div>

    <div class="dashboard-sections">
      <div class="section-card">
        <h3>快捷操作</h3>
        <div class="quick-actions">
          <router-link to="/admin/article/create" class="action-btn primary">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"></line>
              <line x1="5" y1="12" x2="19" y2="12"></line>
            </svg>
            写文章
          </router-link>
          <router-link to="/admin/comments" class="action-btn warning">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"></path>
            </svg>
            审核评论
          </router-link>
          <router-link to="/admin/tools" class="action-btn danger">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="3 6 5 6 21 6"></polyline>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
            </svg>
            清理缓存
          </router-link>
          <router-link to="/admin/security" class="action-btn success">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
            </svg>
            数据备份
          </router-link>
        </div>
      </div>

      <div class="section-card">
        <h3>系统状态
          <button class="refresh-status-btn" @click="fetchSystemStatus" title="刷新">↻</button>
        </h3>
        <div class="system-status" v-if="sysStatus">
          <div class="status-item">
            <span class="status-label">CPU 使用率</span>
            <span class="status-value" :class="{ warn: sysStatus.cpu_usage > 80 }">{{ sysStatus.cpu_usage }}%</span>
          </div>
          <div class="status-item">
            <span class="status-label">内存使用率</span>
            <span class="status-value" :class="{ warn: sysStatus.memory_usage > 85 }">
              {{ sysStatus.memory_usage }}%（{{ sysStatus.memory_used }}MB / {{ sysStatus.memory_total }}MB）
            </span>
          </div>
          <div class="status-item">
            <span class="status-label">数据库连接</span>
            <span class="status-badge" :class="sysStatus.db_status === 'connected' ? 'success' : 'error'">
              {{ sysStatus.db_status === 'connected' ? '已连接' : '异常' }}
            </span>
          </div>
          <div class="status-item">
            <span class="status-label">运行时间</span>
            <span class="status-value">{{ sysStatus.uptime }}</span>
          </div>
        </div>
        <div v-else class="system-status">
          <div class="status-item">
            <span class="status-label">加载中...</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard {
  max-width: 1400px;
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

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  color: #6b7280;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.refresh-btn:hover {
  background: #f9fafb;
  color: #374151;
}

.refresh-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid #e5e7eb;
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border-left: 4px solid;
  transition: transform 0.2s, box-shadow 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.card-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  border-radius: 12px;
}

.card-content {
  flex: 1;
}

.card-title {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 4px;
}

.card-value {
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
}

.dashboard-sections {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 24px;
}

.section-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.section-card h3 {
  margin: 0 0 20px 0;
  font-size: 18px;
  color: #1f2937;
}

.quick-actions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 12px;
}

.action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 20px 16px;
  border-radius: 12px;
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
}

.action-btn.primary {
  background: #6366f115;
  color: #6366f1;
}

.action-btn.warning {
  background: #f59e0b15;
  color: #f59e0b;
}

.action-btn.danger {
  background: #ef444415;
  color: #ef4444;
}

.action-btn.success {
  background: #10b98115;
  color: #10b981;
}

.action-btn:hover {
  transform: translateY(-2px);
  filter: brightness(0.95);
}

.system-status {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.status-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid #f3f4f6;
}

.status-item:last-child {
  border-bottom: none;
}

.status-label {
  font-size: 14px;
  color: #6b7280;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.success {
  background: #10b98115;
  color: #10b981;
}

.status-badge.error {
  background: #ef444415;
  color: #ef4444;
}

.status-value.warn {
  color: #ef4444;
  font-weight: 600;
}

.refresh-status-btn {
  background: none;
  border: none;
  font-size: 16px;
  cursor: pointer;
  color: #6b7280;
  padding: 0 4px;
  vertical-align: middle;
}

.refresh-status-btn:hover {
  color: #374151;
}

.status-value {
  font-size: 14px;
  color: #9ca3af;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .dashboard-sections {
    grid-template-columns: 1fr;
  }
}
</style>