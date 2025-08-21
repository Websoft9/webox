package errors

import (
	"api-service/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware 统一错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		switch err := recovered.(type) {
		case string:
			HandleError(c, NewAppError(CodeInternalError, err))
		case error:
			HandleError(c, WrapError(err, CodeInternalError, "服务器内部错误"))
		default:
			HandleError(c, ErrInternalError)
		}
		c.Abort()
	})
}

// HandleError 统一错误处理函数
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		// 自定义应用错误
		response.Error(c, appErr.HTTPStatus, appErr.Message, appErr.Details)
	} else {
		// 标准错误
		response.Error(c, http.StatusInternalServerError, "服务器内部错误", err.Error())
	}
}

// IsAppError 检查是否为应用错误
func IsAppError(err error) (*AppError, bool) {
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}
