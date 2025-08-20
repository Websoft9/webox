package middleware

import (
	"api-service/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 结构化日志中间件
func LoggerMiddleware(log logger.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		startTime := time.Now()
		
		// 处理请求
		c.Next()
		
		// 计算处理时间
		duration := time.Since(startTime)
		
		// 记录请求日志
		log.InfoContext(c, "HTTP请求",
			logger.String("method", c.Request.Method),
			logger.String("path", c.Request.URL.Path),
			logger.String("client_ip", c.ClientIP()),
			logger.Int("status_code", c.Writer.Status()),
			logger.String("user_agent", c.Request.UserAgent()),
			logger.String("duration", duration.String()),
		)
	})
}

// Logger 保持向后兼容的简单日志中间件
func Logger() gin.HandlerFunc {
	return gin.Logger()
}
