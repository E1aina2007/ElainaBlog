package response

import (
	"fmt"
	"net/http"
)

// AppError 全局业务错误定义，通过 Code 区分错误类型并映射到对应的 HTTP 状态码。
//
// 错误码段约定：
//   - 400xxx: 客户端参数错误 → HTTP 400
//   - 401xxx: 未认证         → HTTP 401
//   - 403xxx: 无权限         → HTTP 403
//   - 404xxx: 资源不存在     → HTTP 404
//   - 409xxx: 资源冲突       → HTTP 409
//   - 500xxx: 服务器内部错误 → HTTP 500
type AppError struct {
	Code    int    `json:"code"`              // 业务错误码
	Type    string `json:"type,omitempty"`    // 错误类型标识（如 UNAUTHORIZED）
	Message string `json:"message"`           // 面向用户的错误描述
	I18nKey string `json:"i18n_key,omitempty"` // 国际化键名，用于前端多语言
	Detail  any    `json:"detail,omitempty"`   // 附加详情，便于调试
}

// Error 实现 error 接口。
func (e *AppError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("code=%d type=%s message=%s", e.Code, e.Type, e.Message)
}

// WithDetail 返回一个携带附加详情的 AppError 副本，不会修改原始实例。
func (e *AppError) WithDetail(detail any) *AppError {
	if e == nil {
		return nil
	}
	cp := *e
	cp.Detail = detail
	return &cp
}

// HTTPStatus 将业务错误码映射为 HTTP 状态码。
func (e *AppError) HTTPStatus() int {
	if e == nil {
		return http.StatusInternalServerError
	}

	switch {
	case e.Code >= 400000 && e.Code < 401000:
		return http.StatusBadRequest
	case e.Code >= 401000 && e.Code < 402000:
		return http.StatusUnauthorized
	case e.Code >= 403000 && e.Code < 404000:
		return http.StatusForbidden
	case e.Code >= 404000 && e.Code < 405000:
		return http.StatusNotFound
	case e.Code >= 409000 && e.Code < 410000:
		return http.StatusConflict
	case e.Code >= 429000 && e.Code < 430000:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

// NewAppError 创建一个基本的 AppError。
func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// 预定义的全局业务错误，handler/middleware 中直接引用，通过 WithDetail 附加上下文信息。
var (
	ErrInvalidParams = &AppError{Code: 400001, Type: "INVALID_PARAMS", Message: "invalid params", I18nKey: "error.invalid_params"}           // 请求参数无效
	ErrUnauthorized  = &AppError{Code: 401001, Type: "UNAUTHORIZED", Message: "unauthorized", I18nKey: "error.unauthorized"}                 // 未认证（token 缺失或无效）
	ErrForbidden     = &AppError{Code: 403001, Type: "FORBIDDEN", Message: "forbidden", I18nKey: "error.forbidden"}                         // 已认证但无权限
	ErrNotFound      = &AppError{Code: 404001, Type: "NOT_FOUND", Message: "resource not found", I18nKey: "error.not_found"}                // 资源不存在
	ErrConflict      = &AppError{Code: 409001, Type: "CONFLICT", Message: "resource conflict", I18nKey: "error.conflict"}                   // 资源冲突（如重复创建）
	ErrInternal      = &AppError{Code: 500001, Type: "INTERNAL_ERROR", Message: "internal server error", I18nKey: "error.internal"}         // 服务器内部错误
	ErrTooManyRequests = &AppError{Code: 429001, Type: "TOO_MANY_REQUESTS", Message: "too many requests", I18nKey: "error.too_many_requests"} // 请求过于频繁

	// 认证相关细分错误码（401xxx）
	ErrPasswordMismatch   = &AppError{Code: 401002, Type: "PASSWORD_MISMATCH", Message: "password mismatch", I18nKey: "error.password_mismatch"}       // 密码错误
	ErrTokenExpired       = &AppError{Code: 401003, Type: "TOKEN_EXPIRED", Message: "token expired", I18nKey: "error.token_expired"}                   // token 已过期
	ErrTokenInvalid       = &AppError{Code: 401004, Type: "TOKEN_INVALID", Message: "token invalid", I18nKey: "error.token_invalid"}                   // token 无效
	ErrRefreshTokenInvalid = &AppError{Code: 401005, Type: "REFRESH_TOKEN_INVALID", Message: "refresh token invalid", I18nKey: "error.refresh_token_invalid"} // 刷新令牌无效

	// 用户相关细分错误码（4001xx）
	ErrUserNotFound       = &AppError{Code: 400101, Type: "USER_NOT_FOUND", Message: "user not found", I18nKey: "error.user_not_found"}               // 用户不存在
	ErrEmailFormat        = &AppError{Code: 400102, Type: "EMAIL_FORMAT", Message: "invalid email format", I18nKey: "error.email_format"}              // 邮箱格式错误
	ErrUsernameFormat     = &AppError{Code: 400103, Type: "USERNAME_FORMAT", Message: "invalid username format", I18nKey: "error.username_format"}     // 用户名格式错误
	ErrPasswordLength     = &AppError{Code: 400104, Type: "PASSWORD_LENGTH", Message: "password length invalid", I18nKey: "error.password_length"}     // 密码长度不符合要求
	ErrPasswordChars      = &AppError{Code: 400105, Type: "PASSWORD_CHARS", Message: "password chars invalid", I18nKey: "error.password_chars"}        // 密码字符不符合要求
	ErrUsernameExists     = &AppError{Code: 409101, Type: "USERNAME_EXISTS", Message: "username already exists", I18nKey: "error.username_exists"}     // 用户名已存在
	ErrEmailExists        = &AppError{Code: 409102, Type: "EMAIL_EXISTS", Message: "email already exists", I18nKey: "error.email_exists"}              // 邮箱已被注册

	// 验证码相关（4002xx）
	ErrCodeExpired        = &AppError{Code: 400201, Type: "CODE_EXPIRED", Message: "verification code expired", I18nKey: "error.code_expired"}         // 验证码已过期
	ErrCodeMismatch       = &AppError{Code: 400202, Type: "CODE_MISMATCH", Message: "verification code mismatch", I18nKey: "error.code_mismatch"}     // 验证码错误
	ErrResendTooFrequent  = &AppError{Code: 429201, Type: "RESEND_TOO_FREQUENT", Message: "resend too frequent", I18nKey: "error.resend_too_frequent"} // 验证码发送过于频繁

	// 资源相关细分错误码（4041xx）
	ErrCategoryNotFound   = &AppError{Code: 404101, Type: "CATEGORY_NOT_FOUND", Message: "category not found", I18nKey: "error.category_not_found"}   // 分类不存在
	ErrArticleNotFound    = &AppError{Code: 404102, Type: "ARTICLE_NOT_FOUND", Message: "article not found", I18nKey: "error.article_not_found"}      // 文章不存在
	ErrCommentNotFound    = &AppError{Code: 404103, Type: "COMMENT_NOT_FOUND", Message: "comment not found", I18nKey: "error.comment_not_found"}      // 评论不存在
	ErrMessageNotFound    = &AppError{Code: 404104, Type: "MESSAGE_NOT_FOUND", Message: "message not found", I18nKey: "error.message_not_found"}      // 留言不存在
	ErrCategoryExists     = &AppError{Code: 409103, Type: "CATEGORY_EXISTS", Message: "category already exists", I18nKey: "error.category_exists"}     // 分类已存在

	// 上传相关（4003xx）
	ErrFileTooLarge       = &AppError{Code: 400301, Type: "FILE_TOO_LARGE", Message: "file too large", I18nKey: "error.file_too_large"}               // 文件过大
	ErrFileTypeInvalid    = &AppError{Code: 400302, Type: "FILE_TYPE_INVALID", Message: "file type invalid", I18nKey: "error.file_type_invalid"}      // 文件类型不支持
)
