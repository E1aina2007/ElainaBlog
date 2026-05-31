<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { sendCode, resetPassword } from '@/api/auth'
import { validateEmail, validatePassword } from '@/utils/validate'

const router = useRouter()

const form = ref({
  email: '',
  code: '',
  newPassword: '',
  confirmPassword: '',
})

const errors = ref({
  email: '',
  code: '',
  newPassword: '',
  confirmPassword: '',
  general: '',
})

const isSubmitting = ref(false)
const isSendingCode = ref(false)
const countdown = ref(0)
const step = ref(1) // 1: 输入邮箱, 2: 输入验证码和新密码

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
    step.value = 2
  } catch (err: any) {
    errors.value.general = err?.message || '发送验证码失败，请稍后重试'
  } finally {
    isSendingCode.value = false
  }
}

function validateForm(): boolean {
  errors.value = { email: '', code: '', newPassword: '', confirmPassword: '', general: '' }
  let valid = true

  const emailErr = validateEmail(form.value.email)
  if (emailErr) {
    errors.value.email = emailErr
    valid = false
  }

  if (!form.value.code) {
    errors.value.code = '请输入验证码'
    valid = false
  }

  const passwordErr = validatePassword(form.value.newPassword)
  if (passwordErr) {
    errors.value.newPassword = passwordErr
    valid = false
  }

  if (form.value.newPassword !== form.value.confirmPassword) {
    errors.value.confirmPassword = '两次输入的密码不一致'
    valid = false
  }

  return valid
}

async function handleResetPassword() {
  if (!validateForm()) return

  isSubmitting.value = true
  errors.value.general = ''

  try {
    await resetPassword(form.value.email, form.value.code, form.value.newPassword)
    router.push('/login')
  } catch (err: any) {
    errors.value.general = err?.message || '重置密码失败，请稍后重试'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <main class="forgot-password-page">
    <div class="forgot-password-container">
      <!-- 装饰 -->
      <div class="deco-top"><img src="/favicon.ico" alt="logo" /></div>

      <div class="forgot-password-card">
        <h1 class="forgot-password-title">重置密码</h1>
        <p class="forgot-password-subtitle">
          {{ step === 1 ? '输入你的邮箱地址，我们将发送验证码' : '输入验证码和新密码' }}
        </p>

        <!-- 通用错误 -->
        <div v-if="errors.general" class="error-banner">
          {{ errors.general }}
        </div>

        <!-- 步骤1：输入邮箱 -->
        <form v-if="step === 1" class="forgot-password-form" @submit.prevent="handleSendCode">
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

          <button type="submit" class="btn-submit" :disabled="!canSendCode">
            <span v-if="isSendingCode">发送中...</span>
            <span v-else>发送验证码</span>
          </button>
        </form>

        <!-- 步骤2：输入验证码和新密码 -->
        <form v-else class="forgot-password-form" @submit.prevent="handleResetPassword">
          <div class="form-group">
            <label class="form-label">邮箱</label>
            <input
              v-model="form.email"
              type="email"
              class="form-input"
              :class="{ error: errors.email }"
              placeholder="请输入邮箱"
              autocomplete="email"
              disabled
            />
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
                <span v-else>重新发送</span>
              </button>
            </div>
            <span v-if="errors.code" class="form-error">{{ errors.code }}</span>
          </div>

          <div class="form-group">
            <label class="form-label">新密码</label>
            <input
              v-model="form.newPassword"
              type="password"
              class="form-input"
              :class="{ error: errors.newPassword }"
              placeholder="8-72位，需含字母和数字"
              autocomplete="new-password"
              @blur="errors.newPassword = validatePassword(form.newPassword) || ''"
            />
            <span v-if="errors.newPassword" class="form-error">{{ errors.newPassword }}</span>
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
              @blur="errors.confirmPassword = form.newPassword === form.confirmPassword ? '' : '两次输入的密码不一致'"
            />
            <span v-if="errors.confirmPassword" class="form-error">{{ errors.confirmPassword }}</span>
          </div>

          <button type="submit" class="btn-submit" :disabled="isSubmitting">
            <span v-if="isSubmitting">重置中...</span>
            <span v-else>重置密码</span>
          </button>
        </form>

        <p class="forgot-password-footer">
          <router-link to="/login" class="link">返回登录</router-link>
        </p>
      </div>
    </div>
  </main>
</template>

<style scoped>
.forgot-password-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--bg-primary) 0%, var(--bg-secondary) 100%);
  padding: 24px;
}

.forgot-password-container {
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

.forgot-password-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 40px 32px;
  box-shadow: var(--shadow-soft);
}

.forgot-password-title {
  font-size: 1.75rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px;
  text-align: center;
}

.forgot-password-subtitle {
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

.forgot-password-form {
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

.form-input:disabled {
  opacity: 0.7;
  cursor: not-allowed;
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

.forgot-password-footer {
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