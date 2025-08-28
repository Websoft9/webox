package middleware

import (
	"api-service/internal/interface/service"
	"api-service/pkg/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// HTTP methods
	methodPOST   = "POST"
	methodDELETE = "DELETE"
	methodGET    = "GET"
	methodPUT    = "PUT"
	methodPATCH  = "PATCH"

	// Actions
	actionRead   = "query"
	actionCreate = "create"
	actionUpdate = "update"
	actionDelete = "delete"

	// Constants
	minResourceParts = 3
)

// PermissionMiddleware 权限验证中间件
func PermissionMiddleware(permissionService service.PermissionService, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			log.Error("User ID not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		// 获取请求路径和方法
		path := c.Request.URL.Path
		method := c.Request.Method

		// 构建资源和操作
		resource, action := buildResourceAction(path, method)

		// 检查用户权限
		hasPermission, err := permissionService.CheckUserPermission(
			c.Request.Context(),
			userID.(uint),
			resource,
			action,
		)

		if err != nil {
			log.Error("Failed to check user permission", logger.ErrorField(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    http.StatusInternalServerError,
				"message": "Permission check failed",
			})
			c.Abort()
			return
		}

		if !hasPermission {
			log.Warn("User permission denied",
				logger.Uint("user_id", userID.(uint)),
				logger.String("resource", resource),
				logger.String("action", action),
				logger.String("path", path),
				logger.String("method", method),
			)

			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    http.StatusForbidden,
				"message": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission 要求特定权限的中间件
func RequirePermission(permissionService service.PermissionService, resource, action string, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			log.Error("User ID not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		// 检查用户权限
		hasPermission, err := permissionService.CheckUserPermission(
			c.Request.Context(),
			userID.(uint),
			resource,
			action,
		)

		if err != nil {
			log.Error("Failed to check user permission", logger.ErrorField(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    http.StatusInternalServerError,
				"message": "Permission check failed",
			})
			c.Abort()
			return
		}

		if !hasPermission {
			log.Warn("User permission denied",
				logger.Uint("user_id", userID.(uint)),
				logger.String("resource", resource),
				logger.String("action", action),
			)

			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    http.StatusForbidden,
				"message": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// buildResourceAction 根据路径和方法构建资源和操作
func buildResourceAction(path, method string) (resource, action string) {
	// 移除 API 版本前缀
	path = strings.TrimPrefix(path, "/api/v1")

	// 解析路径获取资源
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return "", ""
	}

	resource = parts[0]

	// 根据 HTTP 方法确定操作
	switch method {
	case methodGET:
		action = actionRead
	case methodPOST:
		action = actionCreate
	case methodPUT, methodPATCH:
		action = actionUpdate
	case methodDELETE:
		action = actionDelete
	default:
		action = actionRead
	}

	// 特殊路径处理
	if len(parts) >= minResourceParts {
		switch parts[2] {
		case "permissions":
			switch method {
			case methodPOST:
				action = "assign_permission"
			case methodDELETE:
				action = "remove_permission"
			}
		case "users":
			action = "manage_users"
		case "roles":
			action = "manage_roles"
		}
	}

	return resource, action
}
