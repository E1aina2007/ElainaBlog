<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { validateEmail } from '@/utils/validate'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const form = ref({
  email: '',
  password: '',
})

const errors = ref({
  email: '',
  password: '',
  general: '',
})

const isSubmitting = ref(false)

function validateForm(): boolean {
  errors.value = { email: '', password: '', general: '' }
  let valid = true

  const emailErr = validateEmail(form.value.email)
  if (emailErr) {
    errors.value.email = emailErr
    valid = false
  }

  if (!form.value.password) {
    errors.value.password = '请输入密码'
    valid = false
  }

  return valid
}

async function handleLogin() {
  if (!validateForm()) return

  isSubmitting.value = true
  errors.value.general = ''

  try {
    await userStore.login({
      email: form.value.email,
      password: form.value.password,
    })
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (err: any) {
    errors.value.general = err?.message || '登录失败，请检查邮箱和密码'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <main class="login-page">
    <div class="login-container">
      <!-- 装饰 -->
      <div class="deco-top"><img src="/favicon.ico" alt="logo" /></div>

      <div class="login-card">
        <h1 class="login-title">欢迎回来</h1>
        <p class="login-subtitle">登录你的账号，继续探索</p>

        <!-- 通用错误 -->
        <div v-if="errors.general" class="error-banner">
          {{ errors.general }}
        </div>

        <form class="login-form" @submit.prevent="handleLogin">
          <div class="form-group">
            <label class="form-label">邮箱</label>
            <input
              v-model="form.email"
              type="email"
              class="form-input"
              :class="{ error: errors.email }"
              placeholder="请输入邮箱"
              autocomplete="email"
              @blur="errors.email = validateEmail(form.email) || ''"
            />
            <span v-if="errors.email" class="form-error">{{ errors.email }}</span>
          </div>

          <div class="form-group">
            <label class="form-label">密码</label>
            <input
              v-model="form.password"
              type="password"
              class="form-input"
              :class="{ error: errors.password }"
              placeholder="请输入密码"
              autocomplete="current-password"
              @blur="errors.password = form.password ? '' : '请输入密码'"
            />
            <span v-if="errors.password" class="form-error">{{ errors.password }}</span>
          </div>

          <button type="submit" class="btn-submit" :disabled="isSubmitting">
            <span v-if="isSubmitting">登录中...</span>
            <span v-else>登 录</span>
          </button>
        </form>

        <p class="login-footer">
          还没有账号？
          <router-link to="/register" class="link">立即注册</router-link>
        </p>
      </div>
    </div>
  </main>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--bg-primary) 0%, var(--bg-secondary) 100%);
  padding: 24px;
}

.login-container {
  width: 100%;
  max-width: 420px;
  position: relative;
}

.deco-top {
  text-align: center;
  margin-bottom: 16px;
  opacity: 0.6;
}

.deco-top img {
  width: 48px;
  height: 48px;
  object-fit: contain;
}

.login-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 40px 32px;
  box-shadow: var(--shadow-soft);
}

.login-title {
  font-size: 1.75rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px;
  text-align: center;
}

.login-subtitle {
  font-size: 0.875rem;
  color: var(--text-muted);
  margin: 0 0 28px;
  text-align: center;
}

.error-banner {
  background: rgba(255, 183, 178, 0.15);
  border: 1px solid var(--accent);
  color: #c0392b;
  padding: 10px 16px;
  border-radius: var(--radius-md);
  font-size: 0.8125rem;
  margin-bottom: 20px;
  text-align: center;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-secondary);
}

.form-input {
  padding: 12px 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  font-size: 0.9375rem;
  color: var(--text-primary);
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
  outline: none;
}

.form-input::placeholder {
  color: var(--text-muted);
}

.form-input:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(126, 215, 193, 0.15);
}

.form-input.error {
  border-color: #e74c3c;
}

.form-error {
  font-size: 0.75rem;
  color: #e74c3c;
}

.btn-submit {
  width: 100%;
  padding: 12px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-size: 0.9375rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease-out;
  margin-top: 8px;
}

.btn-submit:hover:not(:disabled) {
  background: var(--primary-dark);
  transform: translateY(-1px);
  box-shadow: 0 4px 16px rgba(126, 215, 193, 0.35);
}

.btn-submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.login-footer {
  text-align: center;
  font-size: 0.875rem;
  color: var(--text-muted);
  margin: 24px 0 0;
}

.link {
  color: var(--primary);
  text-decoration: none;
  font-weight: 500;
}

.link:hover {
  text-decoration: underline;
}
</style>