package middleware

import (
	"api-service/internal/interface/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	// HTTP methods
	methodPOST   = "POST"
	methodDELETE = "DELETE"
	methodGET    = "GET"
	methodPUT    = "PUT"

	// Actions
	actionRead   = "read"
	actionCreate = "create"
	actionUpdate = "update"
	actionDelete = "delete"

	// Constants
	minResourceParts = 3
)

// PermissionMiddleware 权限验证中间件
func PermissionMiddleware(permissionService service.PermissionService, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			logger.Error("User ID not found in context")
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
			logger.WithError(err).Error("Failed to check user permission")
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    http.StatusInternalServerError,
				"message": "Permission check failed",
			})
			c.Abort()
			return
		}

		if !hasPermission {
			logger.WithFields(logrus.Fields{
				"user_id":  userID,
				"resource": resource,
				"action":   action,
				"path":     path,
				"method":   method,
			}).Warn("User permission denied")

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
func RequirePermission(permissionService service.PermissionService, resource, action string, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			logger.Error("User ID not found in context")
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
			logger.WithError(err).Error("Failed to check user permission")
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    http.StatusInternalServerError,
				"message": "Permission check failed",
			})
			c.Abort()
			return
		}

		if !hasPermission {
			logger.WithFields(logrus.Fields{
				"user_id":  userID,
				"resource": resource,
				"action":   action,
			}).Warn("User permission denied")

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
	case "PUT", "PATCH":
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

// PermissionConfig 权限配置
type PermissionConfig struct {
	Resource string
	Action   string
}

// RoutePermissions 路由权限映射
var RoutePermissions = map[string]PermissionConfig{
	// 用户管理
	"GET:/users":        {Resource: "user", Action: "read"},
	"POST:/users":       {Resource: "user", Action: "create"},
	"GET:/users/:id":    {Resource: "user", Action: "read"},
	"PUT:/users/:id":    {Resource: "user", Action: "update"},
	"DELETE:/users/:id": {Resource: "user", Action: "delete"},

	// 角色管理
	"GET:/roles":        {Resource: "role", Action: "read"},
	"POST:/roles":       {Resource: "role", Action: "create"},
	"GET:/roles/:id":    {Resource: "role", Action: "read"},
	"PUT:/roles/:id":    {Resource: "role", Action: "update"},
	"DELETE:/roles/:id": {Resource: "role", Action: "delete"},

	// 权限管理
	"GET:/permissions":        {Resource: "permission", Action: "read"},
	"POST:/permissions":       {Resource: "permission", Action: "create"},
	"GET:/permissions/:id":    {Resource: "permission", Action: "read"},
	"PUT:/permissions/:id":    {Resource: "permission", Action: "update"},
	"DELETE:/permissions/:id": {Resource: "permission", Action: "delete"},
}

// GetRoutePermission 获取路由权限配置
func GetRoutePermission(method, path string) (PermissionConfig, bool) {
	// 标准化路径，将具体ID替换为参数占位符
	normalizedPath := normalizePath(path)
	key := method + ":" + normalizedPath

	config, exists := RoutePermissions[key]
	return config, exists
}

// normalizePath 标准化路径
func normalizePath(path string) string {
	// 移除 API 版本前缀
	path = strings.TrimPrefix(path, "/api/v1")

	parts := strings.Split(strings.Trim(path, "/"), "/")
	normalizedParts := make([]string, 0, len(parts))

	for i, part := range parts {
		// 如果是数字ID，替换为参数占位符
		if i > 0 && isNumeric(part) {
			normalizedParts = append(normalizedParts, ":id")
		} else {
			normalizedParts = append(normalizedParts, part)
		}
	}

	return "/" + strings.Join(normalizedParts, "/")
}

// isNumeric 检查字符串是否为数字
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
