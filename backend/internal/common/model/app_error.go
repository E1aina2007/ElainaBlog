package model

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
	Code    int    `json:"code"`             // 业务错误码
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
	default:
		return http.StatusInternalServerError
	}
}

// HTTPStatusFromError 从任意 error 推导 HTTP 状态码；若为 *AppError 则走 HTTPStatus()，否则返回 500。
func HTTPStatusFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if ae, ok := err.(*AppError); ok {
		return ae.HTTPStatus()
	}
	return http.StatusInternalServerError
}

// NewAppError 创建一个基本的 AppError。
func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// NewAppTypedError 创建一个带类型标识的 AppError。
func NewAppTypedError(code int, errType, message string) *AppError {
	return &AppError{Code: code, Type: errType, Message: message}
}

// 预定义的全局业务错误，handler/middleware 中直接引用，通过 WithDetail 附加上下文信息。
var (
	ErrInvalidParams = &AppError{Code: 400001, Type: "INVALID_PARAMS", Message: "invalid params", I18nKey: "error.invalid_params"}       // 请求参数无效
	ErrUnauthorized  = &AppError{Code: 401001, Type: "UNAUTHORIZED", Message: "unauthorized", I18nKey: "error.unauthorized"}              // 未认证（token 缺失或无效）
	ErrForbidden     = &AppError{Code: 403001, Type: "FORBIDDEN", Message: "forbidden", I18nKey: "error.forbidden"}                       // 已认证但无权限
	ErrNotFound      = &AppError{Code: 404001, Type: "NOT_FOUND", Message: "resource not found", I18nKey: "error.not_found"}              // 资源不存在
	ErrConflict      = &AppError{Code: 409001, Type: "CONFLICT", Message: "resource conflict", I18nKey: "error.conflict"}                 // 资源冲突（如重复创建）
	ErrInternal      = &AppError{Code: 500001, Type: "INTERNAL_ERROR", Message: "internal server error", I18nKey: "error.internal"}       // 服务器内部错误
)
