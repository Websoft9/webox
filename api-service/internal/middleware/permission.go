package middleware

import (
	"api-service/internal/config"
	"api-service/internal/interface/service"
	"api-service/pkg/auth"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"api-service/pkg/utils"
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

	// ID validation constants
	maxIDLength     = 10
	minLongIDLength = 8
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

// whiteListRoutes defines routes that require authentication but bypass permission checks
// These routes are accessible to any authenticated user regardless of their specific permissions
var whiteListRoutes = []string{
	"/api/v1/api-tokens/refresh", // API token refresh
	"/api/v1/api-tokens/revoke",  // API token revocation

	"/api/v1/notifications/records",
	"/api/v1/notifications/records/{id}",

	// Notification channel management - allow all authenticated users
	"/api/v1/notifications/channels",              // GET/POST channels list
	"/api/v1/notifications/channels/",             // All channel operations with parameters
	"/api/v1/notifications/channels/email/test",   // Test email channel
	"/api/v1/notifications/channels/webhook/test", // Test webhook channel

	// Notification template management - allow all authenticated users
	"/api/v1/notifications/templates",  // GET/POST templates list
	"/api/v1/notifications/templates/", // All template operations with parameters

	"/api/v1/notifications/templates/{id}/test", // Test template
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

// isWhiteListRoutes checks if a given path is a whitelisted route that requires authentication
// but bypasses permission validation. These routes are accessible to any authenticated user.
func isWhiteListRoutes(path string) bool {
	for _, route := range whiteListRoutes {
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

		if !isWhiteListRoutes(c.Request.URL.Path) {
			// Check if user has permission for the requested resource and action
			if !checkUserPermission(c, permissionService, claims.UserID, log) {
				return
			}
		}
		// Continue to the next handler
		c.Next()
	}
}

// authenticateRequest handles JWT token extraction and validation
func authenticateRequest(c *gin.Context) (*auth.Claims, string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, "", errors.NewAppError(errors.CodeRequiredParameterMissing)
	}

	token, err := auth.ExtractTokenFromHeader(authHeader)
	if err != nil {
		return nil, "", err
	}

	claims, err := auth.GetGlobalJWT().ValidateToken(token)
	if err != nil {
		return nil, "", err
	}

	return claims, token, nil
}

// handleAuthError handles authentication errors
func handleAuthError(c *gin.Context, err error, log logger.Logger) {
	log.WarnContext(c, err.Error())
	c.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"code":    http.StatusUnauthorized,
		"message": i18n.T("auth.user_not_authenticated", utils.GetUserLangFromRedis(c)),
	})
	c.Abort()
}

// handleTokenValidationError handles token validation errors
func handleTokenValidationError(c *gin.Context, log logger.Logger) {
	log.WarnContext(c, "JWT token not found in storage")
	c.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"code":    errors.CodeInvalidToken,
		"message": i18n.T("auth.token_invalid", utils.GetUserLangFromRedis(c)),
	})
	c.Abort()
}

// checkUserPermission checks if user has required permission
func checkUserPermission(c *gin.Context, permissionService service.PermissionService, userID uint, log logger.Logger) bool {
	path := c.Request.URL.Path
	method := c.Request.Method
	resource, action := buildResourceAction(c, path, method)

	hasPermission, err := permissionService.CheckUserPermission(c.Request.Context(), userID, resource, action)
	if err != nil {
		log.Error("Failed to check user permission", logger.ErrorField(err))
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"code":    errors.CodeRecordQueryFailed,
			"message": i18n.T("resource.record_query_failed", utils.GetUserLangFromRedis(c)),
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
			"code":    errors.CodeInsufficientPermissions,
			"message": i18n.T("auth.permission_denied", utils.GetUserLangFromRedis(c)),
		})
		c.Abort()
		return false
	}

	return true
}

// buildResourceAction constructs resource and action identifiers from HTTP path and method
// This function uses Gin's route matching to accurately map dynamic parameters
//
// Examples:
//   - /api/v1/users -> resource: "/users", action: "query"
//   - /api/v1/users/123 -> resource: "/users/*", action: "query"
//   - /api/v1/audit-logs/export -> resource: "/audit-logs/export", action: "query"
//   - /api/v1/roles/5/permissions -> resource: "/roles/*/permissions", action: "query"
func buildResourceAction(c *gin.Context, path, method string) (resource, action string) {
	// Get the matched route pattern from Gin context
	routePattern := c.FullPath()
	if routePattern != "" {
		// Remove API version prefix to get the core resource path
		resource = strings.TrimPrefix(routePattern, "/api/v1")
		// Normalize Gin's :id parameter format to * for consistency with permissions table
		resource = strings.ReplaceAll(resource, ":id", "*")
		if resource == "" {
			resource = "/"
		}
	} else {
		// Fallback to original logic if route pattern is not available
		resource = buildResourceFromPath(path)
	}

	if resource == "" || resource == "/" {
		return "", ""
	}

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

	return resource, action
}

// buildResourceFromPath is a fallback function that constructs resource from path
// when Gin route pattern is not available
func buildResourceFromPath(path string) string {
	// Remove API version prefix to get the core resource path
	path = strings.TrimPrefix(path, "/api/v1")

	if path == "" || path == "/" {
		return ""
	}

	// Ensure path starts with "/"
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Parse path components to build the resource path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return path
	}

	// Build resource path with dynamic parameter normalization
	resourceParts := make([]string, len(parts))
	for i, part := range parts {
		// Check if this part looks like a dynamic parameter
		if isPathParameter(part) {
			resourceParts[i] = "*"
		} else {
			resourceParts[i] = part
		}
	}

	// Construct the full resource path
	return "/" + strings.Join(resourceParts, "/")
}

// isPathParameter determines if a path segment is likely a dynamic parameter
// It checks for numeric IDs, UUIDs, and other common parameter patterns
func isPathParameter(segment string) bool {
	if segment == "" {
		return false
	}

	// Check for numeric ID (positive integers)
	if isNumericID(segment) {
		return true
	}

	return false
}

// isNumericID checks if the segment is a positive integer (common for database IDs)
func isNumericID(segment string) bool {
	if segment == "" || len(segment) > maxIDLength { // Reasonable ID length limit
		return false
	}

	// For longer numeric strings (8+ digits), be more conservative
	// They might be legitimate path segments rather than IDs
	if len(segment) >= minLongIDLength {
		return false
	}

	for _, char := range segment {
		if char < '0' || char > '9' {
			return false
		}
	}

	// Avoid treating "0" as an ID, and ensure it's not just leading zeros
	return segment != "0" && segment[0] != '0'
}
