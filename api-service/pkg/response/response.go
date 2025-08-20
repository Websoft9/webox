package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo 错误信息
type ErrorInfo struct {
	Type    string                  `json:"type"`
	Code    string                  `json:"code"`
	Details []ValidationErrorDetail `json:"details,omitempty"`
}

// ValidationErrorDetail 验证错误详情
type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// Pagination 分页信息
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// Success 成功响应
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// SuccessWithPagination 带分页的成功响应
func SuccessWithPagination(c *gin.Context, items interface{}, pagination Pagination) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"items":      items,
			"pagination": pagination,
		},
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string, err interface{}) {
	var errorInfo *ErrorInfo

	switch e := err.(type) {
	case string:
		errorInfo = &ErrorInfo{
			Type: "error",
			Code: "GENERAL_ERROR",
		}
	case []ValidationErrorDetail:
		errorInfo = &ErrorInfo{
			Type:    "validation_error",
			Code:    "VALIDATION_FAILED",
			Details: e,
		}
	default:
		errorInfo = &ErrorInfo{
			Type: "error",
			Code: "GENERAL_ERROR",
		}
	}

	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Error:   errorInfo,
	})
}

// ValidationError 验证错误响应
func ValidationError(c *gin.Context, errors []ValidationErrorDetail) {
	Error(c, http.StatusBadRequest, "Validation failed", errors)
}

// ErrorResponse 修复参数名冲突问题，避免与内置 error 类型冲突
func ErrorResponse(c *gin.Context, code int, message, errMsg string) {
	Error(c, code, message, errMsg)
}
