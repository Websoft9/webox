package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"

	"api-service/internal/interface/service"
	"api-service/pkg/logger"
)

const (
	// Maximum size for request/response body to be fully captured (10MB)
	maxAuditBodySize = 10 * 1024 * 1024
	// Marker for truncated body
	bodyTruncatedMarker = "[BODY_TOO_LARGE_TRUNCATED]"
)

// limitedResponseWriter custom response writer with size limit
type limitedResponseWriter struct {
	gin.ResponseWriter
	body      *bytes.Buffer
	maxSize   int64
	truncated bool
}

func (rw *limitedResponseWriter) Write(b []byte) (int, error) {
	// Check if adding this data would exceed the limit
	if int64(rw.body.Len()+len(b)) > rw.maxSize {
		rw.truncated = true
		// Don't write to buffer if already truncated
		if rw.body.Len() == 0 {
			rw.body.WriteString(bodyTruncatedMarker)
		}
	} else if !rw.truncated {
		rw.body.Write(b)
	}
	// Always write to actual response
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

		// Check request body size before reading
		var requestBody []byte
		shouldCaptureRequestBody := true
		if ctx.Request.ContentLength > maxAuditBodySize {
			shouldCaptureRequestBody = false
			log.Warn("Request body too large for audit, skipping body capture",
				logger.Field{Key: "content_length", Value: ctx.Request.ContentLength},
				logger.Field{Key: "url", Value: ctx.Request.RequestURI})
		}

		// Wrap response writer to capture response content (with size limit)
		body := &bytes.Buffer{}
		writer := &limitedResponseWriter{
			ResponseWriter: ctx.Writer,
			body:           body,
			maxSize:        maxAuditBodySize,
			truncated:      false,
		}
		ctx.Writer = writer

		// Read request body if size is acceptable
		if shouldCaptureRequestBody && ctx.Request.Body != nil {
			var err error
			requestBody, err = io.ReadAll(io.LimitReader(ctx.Request.Body, maxAuditBodySize+1))
			if err != nil {
				log.Warn("Failed to read request body for audit",
					logger.Field{Key: "error", Value: err},
					logger.Field{Key: "url", Value: ctx.Request.RequestURI})
			} else if len(requestBody) > maxAuditBodySize {
				// Body exceeded limit
				requestBody = []byte(bodyTruncatedMarker)
			}
			// Restore request body for downstream handlers
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

			// Store request body in context for audit logging (if not truncated)
			if string(requestBody) != bodyTruncatedMarker {
				ctx.Set("audit_request_body", requestBody)
			}
		}

		// Continue processing request
		ctx.Next()

		// Calculate request processing time
		duration := time.Since(startTime)
		responseTime := int(duration.Milliseconds())

		// Get response body (may be truncated)
		responseBody := body.Bytes()
		if writer.truncated {
			responseBody = []byte(bodyTruncatedMarker)
		}

		// Submit audit job to worker pool (non-blocking)
		auditLogService.SubmitAuditJob(ctx, responseBody, responseTime)
	}
}
