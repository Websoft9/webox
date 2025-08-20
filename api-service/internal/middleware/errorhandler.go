package middleware

import (
	"api-service/pkg/errors"
	"api-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 错误处理中间件
func ErrorHandler(log logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.ErrorContext(c, "系统发生panic", logger.String("error", err))
			errors.HandleError(c, errors.NewAppError(errors.CodeInternalError, err))
		} else if err, ok := recovered.(error); ok {
			log.ErrorContext(c, "系统发生panic", logger.ErrorField(err))
			errors.HandleError(c, errors.WrapError(err, errors.CodeInternalError, "服务器内部错误"))
		} else {
			log.ErrorContext(c, "系统发生panic", logger.Any("error", recovered))
			errors.HandleError(c, errors.ErrInternalError)
		}
		c.Abort()
	})
}

// RequestValidator 请求验证中间件
func RequestValidator(log logger.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 可以在这里添加通用的请求验证逻辑
		// 比如检查Content-Type、请求大小限制等
		
		c.Next()
		
		// 请求处理完成后的清理工作
	})
}
