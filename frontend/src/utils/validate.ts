// ========================
// 前端输入校验 —— 正则表达式与规则说明
// 与后端 backend/internal/user/validate.go 保持一致
// ========================

// ---------- 邮箱 ----------
// 规则：标准格式 local@domain.tld，长度不超过 100 个字符
export const EMAIL_MAX_LENGTH = 100
export const EMAIL_REGEX = /^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/

// ---------- 用户名 ----------
// 规则：2-20 个字符，仅允许中文、英文字母、数字、下划线
export const USERNAME_REGEX = /^[\u4e00-\u9fa5a-zA-Z0-9_]{2,20}$/

// ---------- 密码 ----------
// 规则：
//   1. 长度 8-72 个字符（bcrypt 上限 72 字节）
//   2. 仅允许英文字母、数字和常见特殊符号 !@#$%^&*()-_=+[]{};:'",.<>?/\|`~
//   3. 必须包含至少一个英文字母
//   4. 必须包含至少一个数字
export const PASSWORD_MIN_LENGTH = 8
export const PASSWORD_MAX_LENGTH = 72
export const PASSWORD_CHARS_REGEX = /^[a-zA-Z0-9!@#$%^&*()\-_=+\[\]{};:'",.<>?/\\|`~]+$/
export const PASSWORD_LETTER_REGEX = /[a-zA-Z]/
export const PASSWORD_DIGIT_REGEX = /[0-9]/

// ========================
// 校验函数
// ========================

export function validateEmail(email: string): string | null {
  if (email.length > EMAIL_MAX_LENGTH) {
    return '邮箱长度不能超过 100 个字符'
  }
  if (!EMAIL_REGEX.test(email)) {
    return '邮箱格式不正确'
  }
  return null
}

export function validateUsername(username: string): string | null {
  if (!USERNAME_REGEX.test(username)) {
    return '用户名只能包含中文、英文字母、数字和下划线，长度 2-20'
  }
  return null
}

export function validatePassword(password: string): string | null {
  if (password.length < PASSWORD_MIN_LENGTH || password.length > PASSWORD_MAX_LENGTH) {
    return '密码长度必须为 8-72 个字符'
  }
  if (!PASSWORD_CHARS_REGEX.test(password)) {
    return '密码包含不允许的字符'
  }
  if (!PASSWORD_LETTER_REGEX.test(password)) {
    return '密码必须包含至少一个英文字母'
  }
  if (!PASSWORD_DIGIT_REGEX.test(password)) {
    return '密码必须包含至少一个数字'
  }
  return null
}
