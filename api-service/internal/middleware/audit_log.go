package middleware

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/gin-gonic/gin"

	"api-service/internal/interface/service"
	"api-service/pkg/logger"
)

// responseWriter custom response writer for capturing response content
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (rw responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// AuditLogMiddleware creates audit log middleware following the functional pattern
func AuditLogMiddleware(auditLogService service.AuditLogService, log logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Delegate skip logic to service
		if auditLogService.ShouldSkipAudit(ctx.Request.Method, ctx.Request.RequestURI) {
			ctx.Next()
			return
		}

		startTime := time.Now()

		// Wrap response writer to capture response content
		body := &bytes.Buffer{}
		writer := responseWriter{
			ResponseWriter: ctx.Writer,
			body:           body,
		}
		ctx.Writer = writer

		// Read request body
		var requestBody []byte
		if ctx.Request.Body != nil {
			requestBody, _ = io.ReadAll(ctx.Request.Body)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

			// Store request body in context for audit logging
			ctx.Set("audit_request_body", requestBody)
		}

		// Continue processing request
		ctx.Next()

		// Calculate request processing time
		duration := time.Since(startTime)
		responseTime := int(duration.Milliseconds())

		// Delegate audit logging to service with background context
		go func() {
			// Use background context to avoid cancellation when HTTP request ends
			backgroundCtx := context.Background()
			if err := auditLogService.RecordAuditFromRequest(
				backgroundCtx,
				ctx,
				body.Bytes(),
				responseTime,
			); err != nil {
				log.Error("Failed to record audit log", logger.Field{Key: "error", Value: err})
			}
		}()
	}
}
