<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { register, sendCode } from '@/api/auth'
import { validateEmail, validateUsername, validatePassword } from '@/utils/validate'

const router = useRouter()

const form = ref({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  code: '',
})

const errors = ref({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  code: '',
  general: '',
})

const isSubmitting = ref(false)
const isSendingCode = ref(false)
const countdown = ref(0)

let countdownTimer: ReturnType<typeof setInterval> | null = null

const canSendCode = computed(() => {
  return form.value.email && !errors.value.email && countdown.value === 0 && !isSendingCode.value
})

function startCountdown() {
  countdown.value = 60
  countdownTimer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      countdown.value = 0
      if (countdownTimer) clearInterval(countdownTimer)
    }
  }, 1000)
}

async function handleSendCode() {
  const emailErr = validateEmail(form.value.email)
  if (emailErr) {
    errors.value.email = emailErr
    return
  }

  isSendingCode.value = true
  errors.value.general = ''

  try {
    await sendCode(form.value.email)
    startCountdown()
  } catch (err: any) {
    errors.value.general = err?.message || '发送验证码失败，请稍后重试'
  } finally {
    isSendingCode.value = false
  }
}

function validateForm(): boolean {
  errors.value = { username: '', email: '', password: '', confirmPassword: '', code: '', general: '' }
  let valid = true

  const usernameErr = validateUsername(form.value.username)
  if (usernameErr) {
    errors.value.username = usernameErr
    valid = false
  }

  const emailErr = validateEmail(form.value.email)
  if (emailErr) {
    errors.value.email = emailErr
    valid = false
  }

  const passwordErr = validatePassword(form.value.password)
  if (passwordErr) {
    errors.value.password = passwordErr
    valid = false
  }

  if (form.value.password !== form.value.confirmPassword) {
    errors.value.confirmPassword = '两次输入的密码不一致'
    valid = false
  }

  if (!form.value.code) {
    errors.value.code = '请输入验证码'
    valid = false
  }

  return valid
}

async function handleRegister() {
  if (!validateForm()) return

  isSubmitting.value = true
  errors.value.general = ''

  try {
    await register({
      username: form.value.username,
      email: form.value.email,
      password: form.value.password,
      code: form.value.code,
    })
    router.push('/login')
  } catch (err: any) {
    errors.value.general = err?.message || '注册失败，请稍后重试'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <main class="register-page">
    <div class="register-container">

      <div class="deco-top"><img src="/favicon.ico" alt="logo" />
      <h1 class="register-title">创建账号</h1>
      <p class="register-subtitle">注册后即可开始你的博客之旅</p>

      <div v-if="errors.general" class="error-banner">
        {{ errors.general }}
        </div>

        <form class="register-form" @submit.prevent="handleRegister">
          <div class="form-group">
            <label class="form-label">用户名</label>
            <input
              v-model="form.username"
              type="text"
              class="form-input"
              :class="{ error: errors.username }"
              placeholder="2-20位，中文/英文/数字/下划线"
              autocomplete="username"
              @blur="errors.username = validateUsername(form.username) || ''"
            />
            <span v-if="errors.username" class="form-error">{{ errors.username }}</span>
          </div>

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
              placeholder="8-72位，需含字母和数字"
              autocomplete="new-password"
              @blur="errors.password = validatePassword(form.password) || ''"
            />
            <span v-if="errors.password" class="form-error">{{ errors.password }}</span>
          </div>

          <div class="form-group">
            <label class="form-label">确认密码</label>
            <input
              v-model="form.confirmPassword"
              type="password"
              class="form-input"
              :class="{ error: errors.confirmPassword }"
              placeholder="再次输入密码"
              autocomplete="new-password"
              @blur="errors.confirmPassword = form.password === form.confirmPassword ? '' : '两次输入的密码不一致'"
            />
            <span v-if="errors.confirmPassword" class="form-error">{{ errors.confirmPassword }}</span>
          </div>

          <div class="form-group">
            <label class="form-label">验证码</label>
            <div class="code-row">
              <input
                v-model="form.code"
                type="text"
                class="form-input code-input"
                :class="{ error: errors.code }"
                placeholder="邮箱验证码"
                @blur="errors.code = form.code ? '' : '请输入验证码'"
              />
              <button
                type="button"
                class="btn-send-code"
                :disabled="!canSendCode"
                @click="handleSendCode"
              >
                <span v-if="isSendingCode">发送中...</span>
                <span v-else-if="countdown > 0">{{ countdown }}s</span>
                <span v-else>发送验证码</span>
              </button>
            </div>
            <span v-if="errors.code" class="form-error">{{ errors.code }}</span>
          </div>

          <button type="submit" class="btn-submit" :disabled="isSubmitting">
            <span v-if="isSubmitting">注册中...</span>
            <span v-else>注 册</span>
          </button>
        </form>

        <p class="register-footer">
          已有账号？
          <router-link to="/login" class="link">立即登录</router-link>
        </p>
      </div>
    </div>
  </main>
</template>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--bg-primary) 0%, var(--bg-secondary) 100%);
  padding: 24px;
}

.register-container {
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

.register-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 40px 32px;
  box-shadow: var(--shadow-soft);
}

.register-title {
  font-size: 1.75rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px;
  text-align: center;
}

.register-subtitle {
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

.register-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
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

.code-row {
  display: flex;
  gap: 10px;
}

.code-input {
  flex: 1;
}

.btn-send-code {
  flex-shrink: 0;
  padding: 12px 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--primary);
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.2s ease-out;
}

.btn-send-code:hover:not(:disabled) {
  background: var(--primary-lighter);
  border-color: var(--primary-light);
}

.btn-send-code:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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

.register-footer {
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