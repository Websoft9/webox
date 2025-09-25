package utils

import (
	"context"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetUserIDFromContext(t *testing.T) {
	tests := []struct {
		name     string
		setupCtx func() context.Context
		expected uint
		found    bool
	}{
		{
			name: "Get user ID from standard context",
			setupCtx: func() context.Context {
				return SetUserIDInContext(context.Background(), 123)
			},
			expected: 123,
			found:    true,
		},
		{
			name: "Get user ID from Gin context",
			setupCtx: func() context.Context {
				gin.SetMode(gin.TestMode)
				ctx, _ := gin.CreateTestContext(nil)
				ctx.Set("user_id", uint(456))
				return ctx
			},
			expected: 456,
			found:    true,
		},
		{
			name: "No user ID in context",
			setupCtx: func() context.Context {
				return context.Background()
			},
			expected: 0,
			found:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx()
			userID, found := GetUserIDFromContext(ctx)
			assert.Equal(t, tt.expected, userID)
			assert.Equal(t, tt.found, found)
		})
	}
}

func TestContextWithUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ginCtx, _ := gin.CreateTestContext(nil)
	ginCtx.Request = &http.Request{}
	ginCtx.Request = ginCtx.Request.WithContext(context.Background())
	ginCtx.Set("user_id", uint(789))

	ctx := ContextWithUserID(ginCtx)
	userID, found := GetUserIDFromContext(ctx)

	assert.True(t, found)
	assert.Equal(t, uint(789), userID)
}

func TestSetUserIDInContext(t *testing.T) {
	ctx := context.Background()
	ctx = SetUserIDInContext(ctx, 999)

	userID, found := GetUserIDFromContext(ctx)
	assert.True(t, found)
	assert.Equal(t, uint(999), userID)
}
