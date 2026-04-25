package model

import (
	"encoding/json"
	"time"
)

type ApiResponse struct {
	Success   bool   `json:"success"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	Error     any    `json:"error,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

func ApiSuccessResponse(data any) *ApiResponse {
	return &ApiResponse{
		Success:   true,
		Code:      0,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	}
}

func ApiErrorResponse(code int, message string, err ...any) *ApiResponse {
	res := &ApiResponse{
		Success:   false,
		Code:      code,
		Message:   message,
		Timestamp: time.Now().UnixMilli(),
	}
	if len(err) > 0 {
		res.Error = err[0]
	}
	return res
}

func (r *ApiResponse) ToJson() string {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func (r *ApiResponse) ToJsonBytes() []byte {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	return jsonBytes
}
