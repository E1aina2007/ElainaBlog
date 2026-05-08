package user

import (
	"errors"
	"regexp"
)

// ========================
// 正则表达式定义
// ========================

// 邮箱：标准格式 local@domain.tld，长度不超过 100（数据库 VARCHAR(100)）
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// 用户名：2-20 个字符，允许中文、英文字母、数字、下划线
var usernameRegex = regexp.MustCompile(`^[\p{Han}a-zA-Z0-9_]{2,20}$`)

// 密码 — 合法字符：英文字母、数字、常见特殊符号
var passwordCharsRegex = regexp.MustCompile("^[a-zA-Z0-9!@#$%^&*()\\-_=+\\[\\]{};:'\",.<>?/\\\\|`~]+$")

// 密码 — 至少包含一个英文字母
var passwordLetterRegex = regexp.MustCompile(`[a-zA-Z]`)

// 密码 — 至少包含一个数字
var passwordDigitRegex = regexp.MustCompile(`[0-9]`)

// ========================
// 校验错误定义
// ========================

var (
	ErrEmailTooLong       = errors.New("邮箱长度不能超过 100 个字符")
	ErrEmailFormat        = errors.New("邮箱格式不正确")
	ErrUsernameFormat     = errors.New("用户名只能包含中文、英文字母、数字和下划线，长度 2-20")
	ErrPasswordLength     = errors.New("密码长度必须为 8-72 个字符")
	ErrPasswordChars      = errors.New("密码包含不允许的字符")
	ErrPasswordNeedLetter = errors.New("密码必须包含至少一个英文字母")
	ErrPasswordNeedDigit  = errors.New("密码必须包含至少一个数字")
)

// ========================
// 校验函数
// ========================

// ValidateEmail 校验邮箱格式：标准 local@domain.tld，长度 ≤ 100
func ValidateEmail(email string) error {
	if len(email) > 100 {
		return ErrEmailTooLong
	}
	if !emailRegex.MatchString(email) {
		return ErrEmailFormat
	}
	return nil
}

// ValidateUsername 校验用户名格式：中文/英文/数字/下划线，2-20 字符
func ValidateUsername(username string) error {
	if !usernameRegex.MatchString(username) {
		return ErrUsernameFormat
	}
	return nil
}

// ValidatePassword 校验密码格式：8-72 字符，至少含一个字母和一个数字，仅允许合法字符
func ValidatePassword(password string) error {
	if len(password) < 8 || len(password) > 72 {
		return ErrPasswordLength
	}
	if !passwordCharsRegex.MatchString(password) {
		return ErrPasswordChars
	}
	if !passwordLetterRegex.MatchString(password) {
		return ErrPasswordNeedLetter
	}
	if !passwordDigitRegex.MatchString(password) {
		return ErrPasswordNeedDigit
	}
	return nil
}
