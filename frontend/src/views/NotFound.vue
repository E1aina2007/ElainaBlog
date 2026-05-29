<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

function goHome() {
  router.push('/')
}

function goBack() {
  router.back()
}
</script>

<template>
  <main class="not-found-page">
    <div class="content">
      <div class="error-code">404</div>
      <h1 class="title">页面走丢了</h1>
      <p class="desc">你访问的页面不存在或已被移除</p>
      <div class="actions">
        <button class="btn-back" @click="goBack">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M19 12H5M12 19l-7-7 7-7" />
          </svg>
          返回上页
        </button>
        <button class="btn-home" @click="goHome">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
            <polyline points="9 22 9 12 15 12 15 22" />
          </svg>
          回到首页
        </button>
      </div>
    </div>
  </main>
</template>

<style scoped>
.not-found-page {
  min-height: 80vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 24px;
}

.content {
  text-align: center;
}

.error-code {
  font-size: 8rem;
  font-weight: 800;
  line-height: 1;
  background: linear-gradient(135deg, var(--primary) 0%, var(--accent) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 16px;
}

.title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px;
}

.desc {
  font-size: 1rem;
  color: var(--text-muted);
  margin: 0 0 32px;
}

.actions {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.btn-back,
.btn-home {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all var(--transition-fast);
  border: none;
}

.btn-back {
  background: var(--bg-card);
  color: var(--text-primary);
  border: 1px solid var(--border);
}

.btn-back:hover {
  background: var(--bg-secondary);
}

.btn-home {
  background: var(--primary);
  color: white;
}

.btn-home:hover {
  background: var(--primary-dark);
  transform: translateY(-1px);
  box-shadow: var(--shadow-soft);
}
</style>
