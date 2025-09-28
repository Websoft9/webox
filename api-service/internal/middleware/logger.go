package middleware

import (
	"api-service/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware structured logging middleware
func LoggerMiddleware(log logger.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate processing time
		duration := time.Since(startTime)

		// Log request information
		log.InfoContext(c, "HTTP Request",
			logger.String("method", c.Request.Method),
			logger.String("path", c.Request.URL.Path),
			logger.String("client_ip", c.ClientIP()),
			logger.Int("status_code", c.Writer.Status()),
			logger.String("user_agent", c.Request.UserAgent()),
			logger.String("duration", duration.String()),
		)
	})
}

// Logger simple logging middleware for backward compatibility
func Logger() gin.HandlerFunc {
	return gin.Logger()
}
