package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse 修复参数名冲突问题，避免与内置 error 类型冲突
func ErrorResponse(c *gin.Context, code int, message, errMsg string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Error:   errMsg,
	})
}

// Error 保持向后兼容性的别名
func Error(c *gin.Context, code int, message, errMsg string) {
	ErrorResponse(c, code, message, errMsg)
}
