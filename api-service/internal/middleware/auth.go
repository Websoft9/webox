package middleware

import (
	"api-service/internal/config"
	"api-service/pkg/auth"
	"api-service/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JWTAuth creates JWT authentication middleware
func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	jwtAuth := auth.NewJWTAuth(cfg.JWT.Secret, cfg.JWT.ExpireTime)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header required", "")
			c.Abort()
			return
		}

		// Extract token from header
		tokenString, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid authorization format", err.Error())
			c.Abort()
			return
		}

		// Validate token
		claims, err := jwtAuth.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid token", err.Error())
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("claims", claims)
		c.Next()
	}
}

// OptionalJWTAuth creates optional JWT authentication middleware
// This middleware extracts user information if a valid token is provided,
// but doesn't require authentication
func OptionalJWTAuth(cfg *config.Config) gin.HandlerFunc {
	jwtAuth := auth.NewJWTAuth(cfg.JWT.Secret, cfg.JWT.ExpireTime)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// Extract token from header
			tokenString, err := auth.ExtractTokenFromHeader(authHeader)
			if err == nil {
				// Validate token
				claims, err := jwtAuth.ValidateToken(tokenString)
				if err == nil {
					// Set user information in context if token is valid
					c.Set("user_id", claims.UserID)
					c.Set("username", claims.Username)
					c.Set("role", claims.Role)
					c.Set("claims", claims)
				}
			}
		}
		c.Next()
	}
}

// RoleBasedAuth creates role-based authorization middleware
func RoleBasedAuth(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context
		role, exists := c.Get("role")
		if !exists {
			response.Error(c, http.StatusForbidden, "User role not found", "")
			c.Abort()
			return
		}

		userRole, ok := role.(string)
		if !ok {
			response.Error(c, http.StatusForbidden, "Invalid user role", "")
			c.Abort()
			return
		}

		// Check if user has required role
		hasRequiredRole := false
		for _, reqRole := range requiredRoles {
			if userRole == reqRole {
				hasRequiredRole = true
				break
			}
		}

		if !hasRequiredRole {
			response.Error(c, http.StatusForbidden, "Insufficient permissions", "")
			c.Abort()
			return
		}

		c.Next()
	}
}

// APIKeyAuth creates API key authentication middleware for API tokens
func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for API key in header
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != "" {
			// TODO: Implement API key validation
			// This would involve checking the API key against the database
			// and setting appropriate context values

			// For now, this is a placeholder
			response.Error(c, http.StatusUnauthorized, "API key authentication not implemented", "")
			c.Abort()
			return
		}

		// If no API key, continue to next middleware (might be JWT auth)
		c.Next()
	}
}

// CombinedAuth creates middleware that supports both JWT and API key authentication
func CombinedAuth(cfg *config.Config) gin.HandlerFunc {
	jwtAuth := auth.NewJWTAuth(cfg.JWT.Secret, cfg.JWT.ExpireTime)

	return func(c *gin.Context) {
		// Try API key authentication first
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != "" {
			// TODO: Implement API key validation
			// For now, reject API key authentication
			response.Error(c, http.StatusUnauthorized, "API key authentication not implemented", "")
			c.Abort()
			return
		}

		// Fall back to JWT authentication
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header or API key required", "")
			c.Abort()
			return
		}

		// Extract and validate JWT token
		tokenString, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid authorization format", err.Error())
			c.Abort()
			return
		}

		claims, err := jwtAuth.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid token", err.Error())
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("claims", claims)
		c.Next()
	}
}

// AdminAuth creates middleware that requires admin role
func AdminAuth() gin.HandlerFunc {
	return RoleBasedAuth("admin")
}

// UserAuth creates middleware that requires user or admin role
func UserAuth() gin.HandlerFunc {
	return RoleBasedAuth("user", "admin")
}
