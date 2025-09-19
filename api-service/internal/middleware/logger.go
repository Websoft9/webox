package middleware

import (
	"api-service/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware is a Gin middleware for logging HTTP requests and responses
func LoggerMiddleware(log logger.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		startTime := time.Now()

		// Process the request
		c.Next()

		// Calculate the duration
		duration := time.Since(startTime)

		// Log the request details
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

// Logger returns a Gin logger middleware instance
func Logger() gin.HandlerFunc {
	return gin.Logger()
}
