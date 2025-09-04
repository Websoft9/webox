package middleware

import (
	"api-service/internal/config"
	"api-service/internal/interface/service"
	"api-service/pkg/auth"
	"api-service/pkg/logger"
	"errors"
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

// publicRoutes defines the list of routes that don't require authentication
// These routes are accessible without valid JWT tokens
var publicRoutes = []string{
	"/api/v1/auth/register",            // User registration endpoint
	"/api/v1/auth/login",               // User login endpoint
	"/api/v1/auth/logout",              // User logout endpoint
	"/api/v1/auth/forgot-password",     // Password reset request
	"/api/v1/auth/reset-password",      // Password reset confirmation
	"/api/v1/auth/verify-email",        // Email verification
	"/api/v1/auth/resend-verification", // Resend verification email
	"/api/v1/auth/oauth2/login",        // OAuth2 authentication
	"/api/v1/i18n/",                    // Internationalization resources
	"/health",                          // Health check endpoint
	"/ping",                            // Ping endpoint
	"/readiness",                       // Readiness probe
	"/liveness",                        // Liveness probe
	"/swagger/",                        // API documentation
}

// isPublicRoute checks if a given path is a public route that doesn't require authentication
// It uses prefix matching to allow for dynamic path parameters
func isPublicRoute(path string) bool {
	for _, route := range publicRoutes {
		if strings.HasPrefix(path, route) {
			return true
		}
	}
	return false
}

// PermissionMiddleware creates a middleware for JWT authentication and permission validation
// It performs the following operations:
// 1. Checks if the route is public (no authentication required)
// 2. Validates JWT token signature and format
// 3. Verifies token exists in storage (Redis/Database)
// 4. Checks user permissions for the requested resource and action
func PermissionMiddleware(permissionService service.PermissionService, apiTokenService service.APITokenService, cfg *config.Config, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for public routes
		if isPublicRoute(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Authenticate the request and extract JWT claims
		claims, token, authErr := authenticateRequest(c)
		if authErr != nil {
			handleAuthError(c, authErr, log)
			return
		}

		// Validate token exists in storage
		exists := apiTokenService.CheckTokenIsExists(c.Request.Context(), token)
		if !exists {
			handleTokenValidationError(c, log)
			return
		}

		// Set user information in context for downstream handlers
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("claims", claims)

		// Check if user has permission for the requested resource and action
		if !checkUserPermission(c, permissionService, claims.UserID, log) {
			return
		}

		// Continue to the next handler
		c.Next()
	}
}

// authenticateRequest handles JWT token extraction and validation
func authenticateRequest(c *gin.Context) (*auth.Claims, string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, "", errors.New("authorization header missing")
	}

	token, err := extractBearerToken(authHeader)
	if err != nil {
		return nil, "", err
	}

	claims, err := auth.ValidateToken(token)
	if err != nil {
		return nil, "", errors.New("invalid token signature")
	}

	return claims, token, nil
}

// extractBearerToken extracts token from Bearer authorization header
func extractBearerToken(authHeader string) (string, error) {
	if len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
		return "", errors.New("invalid authorization header format")
	}
	return authHeader[7:], nil
}

// handleAuthError handles authentication errors
func handleAuthError(c *gin.Context, err error, log logger.Logger) {
	log.WarnContext(c, err.Error())
	c.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"code":    http.StatusUnauthorized,
		"message": "User not authenticated",
	})
	c.Abort()
}

// handleTokenValidationError handles token validation errors
func handleTokenValidationError(c *gin.Context, log logger.Logger) {
	log.WarnContext(c, "JWT token not found in storage")
	c.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"code":    http.StatusUnauthorized,
		"message": "Invalid token",
	})
	c.Abort()
}

// checkUserPermission checks if user has required permission
func checkUserPermission(c *gin.Context, permissionService service.PermissionService, userID uint, log logger.Logger) bool {
	path := c.Request.URL.Path
	method := c.Request.Method
	resource, action := buildResourceAction(path, method)

	hasPermission, err := permissionService.CheckUserPermission(c.Request.Context(), userID, resource, action)
	if err != nil {
		log.Error("Failed to check user permission", logger.ErrorField(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": "Permission check failed",
		})
		c.Abort()
		return false
	}

	if !hasPermission {
		log.Warn("User permission denied",
			logger.Uint("user_id", userID),
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
		return false
	}

	return true
}

// buildResourceAction constructs resource and action identifiers from HTTP path and method
// This function maps REST API endpoints to permission system resources and actions
// It supports both standard CRUD operations and special administrative actions
func buildResourceAction(path, method string) (resource, action string) {
	// Remove API version prefix to get the core resource path
	path = strings.TrimPrefix(path, "/api/v1")

	// Parse path components to identify the resource
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return "", ""
	}

	// The first path component is typically the resource name
	resource = parts[0]

	// Map HTTP methods to permission actions
	switch method {
	case methodGET:
		action = actionRead // Read/query operations
	case methodPOST:
		action = actionCreate // Create operations
	case methodPUT, methodPATCH:
		action = actionUpdate // Update/modify operations
	case methodDELETE:
		action = actionDelete // Delete operations
	default:
		action = actionRead // Default to read for unknown methods
	}

	// Handle special administrative endpoints with custom actions
	if len(parts) >= minResourceParts {
		switch parts[2] {
		case "permissions":
			// Permission management endpoints
			switch method {
			case methodPOST:
				action = "assign_permission"
			case methodDELETE:
				action = "remove_permission"
			}
		case "users":
			// User management endpoints
			action = "manage_users"
		case "roles":
			// Role management endpoints
			action = "manage_roles"
		}
	}

	return resource, action
}
