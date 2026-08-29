<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/api/request'
import toast from '@/utils/toast'

const backupLoading = ref(false)
const bannedIPs = ref<BannedIPInfo[]>([])
const banInput = ref('')
const banning = ref(false)

// 后端 /security/banned-ips 返回的封禁条目；ttl_seconds 为 -1 表示永久封禁
interface BannedIPInfo {
  ip: string
  banned_at: number
  ttl_seconds: number
}

const formatTime = (unix: number): string =>
  unix > 0 ? new Date(unix * 1000).toLocaleString('zh-CN', { hour12: false }) : '--'

const formatTTL = (ttl: number): string => {
  if (ttl < 0) return '永久'
  const minutes = Math.ceil(ttl / 60)
  if (minutes < 60) return `${minutes} 分钟`
  const hours = Math.floor(minutes / 60)
  const rest = minutes % 60
  return rest > 0 ? `${hours} 小时 ${rest} 分钟` : `${hours} 小时`
}

// 手动封禁 IP（永久封禁，可解封）
const handleBan = async () => {
  const ip = banInput.value.trim()
  if (!ip) return
  banning.value = true
  try {
    await request.post('/security/ban', { ip })
    toast.success('IP 已封禁')
    banInput.value = ''
    await fetchBannedIPs()
  } catch {
    toast.error('封禁失败，请检查 IP 格式')
  } finally {
    banning.value = false
  }
}

// 一键备份
const handleBackup = async () => {
  if (!confirm('确定要备份全站数据吗？备份文件将包含所有数据库内容。')) return
  backupLoading.value = true
  try {
    const response = await request.get('/backup/export', { responseType: 'blob' }) as Blob
    // 创建下载链接
    const blob = new Blob([response], { type: 'application/sql' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `backup_${new Date().toISOString().split('T')[0]}.sql`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    toast.success('备份成功！文件已下载')
  } catch (error) {
    console.error('备份失败:', error)
    toast.error('备份失败，请稍后重试')
  } finally {
    backupLoading.value = false
  }
}

// 获取被封禁的IP列表
const fetchBannedIPs = async () => {
  try {
    const data = await request.get('/security/banned-ips') as BannedIPInfo[]
    bannedIPs.value = data || []
  } catch {
    console.log('封禁功能需要后端支持')
    bannedIPs.value = []
  }
}

// 解封IP
const handleUnban = async (ip: string) => {
  if (!confirm(`确定要解封 IP ${ip} 吗？`)) return
  try {
    await request.post('/security/unban', { ip })
    bannedIPs.value = bannedIPs.value.filter(item => item.ip !== ip)
    toast.success('解封成功')
  } catch {
    toast.error('解封失败')
  }
}

onMounted(() => {
  fetchBannedIPs()
})
</script>

<template>
  <div class="security-panel">
    <div class="page-header">
      <h2>安全与备份</h2>
    </div>

    <div class="security-grid">
      <!-- 一键备份 -->
      <div class="security-card backup-card">
        <div class="card-icon">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
            <polyline points="8 12 12 16 16 12"></polyline>
            <line x1="12" y1="8" x2="12" y2="16"></line>
          </svg>
        </div>
        <h3>一键备份</h3>
        <p class="card-desc">
          导出全站数据库备份，包含所有文章、评论、用户和配置数据。
          <br><strong style="color: #ef4444;">这是底线功能！定期备份以防数据丢失。</strong>
        </p>
        <button class="btn-primary large" @click="handleBackup" :disabled="backupLoading">
          <svg v-if="!backupLoading" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
            <polyline points="7 10 12 15 17 10"></polyline>
            <line x1="12" y1="15" x2="12" y2="3"></line>
          </svg>
          <span v-else class="loading-spinner"></span>
          {{ backupLoading ? '备份中...' : '立即备份数据库' }}
        </button>
      </div>

      <!-- 登录防护 -->
      <div class="security-card">
        <div class="card-icon warning">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
            <line x1="12" y1="8" x2="12" y2="12"></line>
            <line x1="12" y1="16" x2="12.01" y2="16"></line>
          </svg>
        </div>
        <h3>登录防爆破</h3>
        <p class="card-desc">
          系统会自动封禁多次登录失败的IP地址。
          <br>默认规则：<strong>15 分钟内失败 10 次后封禁 1 小时</strong>，管理员可手动永久封禁
        </p>
        <div class="stats-row">
          <div class="stat">
            <span class="stat-value">{{ bannedIPs.length }}</span>
            <span class="stat-label">当前封禁IP</span>
          </div>
        </div>
      </div>

      <!-- 封禁IP列表 -->
      <div class="security-card full-width">
        <div class="list-header">
          <h3>封禁IP列表</h3>
          <div class="ban-form">
            <input
              v-model="banInput"
              class="ban-input"
              placeholder="输入要封禁的 IP 地址"
              @keyup.enter="handleBan"
            />
            <button class="btn-primary" :disabled="banning || !banInput.trim()" @click="handleBan">
              {{ banning ? '封禁中...' : '封禁' }}
            </button>
          </div>
        </div>
        <div v-if="bannedIPs.length === 0" class="empty-state">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
            <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
          </svg>
          <p>当前没有被封禁的IP</p>
          <span class="hint">系统会自动封禁多次登录失败的IP</span>
        </div>
        <table v-else class="ip-table">
          <thead>
            <tr>
              <th>IP地址</th>
              <th>封禁时间</th>
              <th>剩余时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in bannedIPs" :key="item.ip">
              <td class="ip-address">{{ item.ip }}</td>
              <td>{{ formatTime(item.banned_at) }}</td>
              <td>{{ formatTTL(item.ttl_seconds) }}</td>
              <td>
                <button class="btn-text" @click="handleUnban(item.ip)">解封</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 安全建议 -->
      <div class="security-card full-width tips-card">
        <h3>安全建议</h3>
        <ul class="tips-list">
          <li>
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            <span>定期更换管理员密码，使用强密码（包含大小写字母、数字和特殊字符）</span>
          </li>
          <li>
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            <span>开启登录防爆破功能，防止暴力破解攻击</span>
          </li>
          <li>
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            <span>每周至少进行一次数据库备份</span>
          </li>
          <li>
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="8" x2="12" y2="12"></line>
              <line x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>
            <span>不要在公共网络环境下登录管理后台</span>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<style scoped>
.security-panel {
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

.security-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 24px;
}

.security-card {
  background: var(--bg-card);
  border-radius: 12px;
  padding: 24px;
  box-shadow: var(--shadow-card);
}

.security-card.full-width {
  grid-column: 1 / -1;
}

.security-card.backup-card {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: #fff;
}

.card-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 60px;
  height: 60px;
  background: color-mix(in srgb, var(--color-indigo) 10%, transparent);
  color: var(--color-indigo);
  border-radius: 12px;
  margin-bottom: 16px;
}

.card-icon.warning {
  background: color-mix(in srgb, var(--color-warning) 10%, transparent);
  color: var(--color-warning);
}

.security-card.backup-card .card-icon {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
}

.security-card h3 {
  margin: 0 0 12px 0;
  font-size: 18px;
  font-weight: 600;
}

.card-desc {
  margin: 0 0 20px 0;
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-secondary);
}

.security-card.backup-card .card-desc {
  color: rgba(255, 255, 255, 0.9);
}

.btn-primary {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  background: var(--color-indigo);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary.large {
  padding: 16px 32px;
  font-size: 16px;
}

.btn-primary:hover:not(:disabled) {
  background: var(--color-indigo-hover);
  transform: translateY(-1px);
}

.security-card.backup-card .btn-primary {
  background: #fff;
  color: var(--color-indigo);
}

.security-card.backup-card .btn-primary:hover:not(:disabled) {
  background: var(--toolbar-bg);
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

.stats-row {
  display: flex;
  gap: 24px;
}

.stat {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
}

.stat-label {
  font-size: 13px;
  color: var(--text-muted);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 40px;
  color: var(--text-muted);
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}

.empty-state .hint {
  font-size: 13px;
  color: #d1d5db;
}

.ip-table {
  width: 100%;
  border-collapse: collapse;
}

.ip-table th,
.ip-table td {
  padding: 12px;
  text-align: left;
  font-size: 14px;
  border-bottom: 1px solid var(--border);
}

.ip-table th {
  font-weight: 600;
  color: var(--text-primary);
  background: var(--toolbar-bg);
}

.ip-table td {
  color: var(--text-secondary);
}

.ip-address {
  font-family: monospace;
  color: var(--text-primary);
}

.btn-text {
  padding: 6px 12px;
  background: none;
  border: none;
  color: var(--color-indigo);
  cursor: pointer;
  font-size: 14px;
}

.btn-text:hover {
  text-decoration: underline;
}

.tips-card {
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
}

.tips-card h3 {
  color: #166534;
}

.tips-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.tips-list li {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 14px;
  color: #166534;
}

.list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
}

.ban-form {
  display: flex;
  gap: 8px;
}

.ban-input {
  width: 220px;
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 14px;
  outline: none;
}

.ban-input:focus {
  border-color: #3b82f6;
}
</style>
