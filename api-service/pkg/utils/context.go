package utils

import (
	"api-service/internal/constants"
	"api-service/pkg/redis"
	"context"

	"github.com/gin-gonic/gin"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
)

// GetUserIDFromContext attempts to get user ID from context
// It tries different approaches:
// 1. First, try to get from Gin context if available
// 2. Then, try to get from standard context if it was set there
// Returns the user ID and a boolean indicating if it was found
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	// Try to get from Gin context if it's a Gin context
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if userID, exists := ginCtx.Get("user_id"); exists {
			if id, ok := userID.(uint); ok {
				return id, true
			}
		}
	}

	// Try to get from standard context
	if userID := ctx.Value(UserIDKey); userID != nil {
		if id, ok := userID.(uint); ok {
			return id, true
		}
	}

	// Also check the original key format in case it was set differently
	if userID := ctx.Value("user_id"); userID != nil {
		if id, ok := userID.(uint); ok {
			return id, true
		}
	}

	return 0, false
}

// SetUserIDInContext sets user ID in context
func SetUserIDInContext(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// ContextWithUserID creates a context with user ID from Gin context
// This is a helper function for controllers to pass user information to services
func ContextWithUserID(ginCtx *gin.Context) context.Context {
	ctx := ginCtx.Request.Context()
	if userID, exists := ginCtx.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			ctx = SetUserIDInContext(ctx, id)
		}
	}
	return ctx
}

// GetUserLangFromRedis gets user language from Redis
func GetUserLangFromRedis(ctx *gin.Context) string {
	userID, exists := ctx.Get("user_id")
	if exists {
		redisKey := redis.FormatRedisKeyWithID(redis.RK_USER_PREFERENCES, userID.(uint))
		language, err := redis.HGet(ctx, redisKey, constants.UserLanguage)

		if err == nil && language != "" {
			return language
		}
	}
	return constants.DefaultLanguage
}
