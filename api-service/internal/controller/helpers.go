package controller

import (
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"api-service/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindAndValidateRequest binds JSON request and validates it
func BindAndValidateRequest(ctx *gin.Context, req interface{}, validator *validator.Validate, log logger.Logger) bool {
	// Bind request parameters
	if err := ctx.ShouldBindJSON(req); err != nil {
		log.ErrorContext(ctx.Request.Context(), "Invalid request format", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": i18n.T("validation.invalid_request_format", utils.GetUserLangFromRedis(ctx)),
			"error":   err.Error(),
		})
		return false
	}

	// Validate request parameters
	if err := validator.Struct(req); err != nil {
		log.ErrorContext(ctx.Request.Context(), "Request validation failed", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": i18n.T("validation.request_validation_failed", utils.GetUserLangFromRedis(ctx)),
			"error":   err.Error(),
		})
		return false
	}

	return true
}

// BindAndValidateQuery binds query parameters and validates them
func BindAndValidateQuery(ctx *gin.Context, req interface{}, validator *validator.Validate, log logger.Logger) bool {
	// Bind query parameters
	if err := ctx.ShouldBindQuery(req); err != nil {
		log.ErrorContext(ctx.Request.Context(), "Invalid query parameters", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": i18n.T("validation.invalid_query_parameters", utils.GetUserLangFromRedis(ctx)),
			"error":   err.Error(),
		})
		return false
	}

	// Validate request parameters
	if err := validator.Struct(req); err != nil {
		log.ErrorContext(ctx.Request.Context(), "Query validation failed", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": i18n.T("validation.invalid_query_parameters", utils.GetUserLangFromRedis(ctx)),
			"error":   err.Error(),
		})
		return false
	}

	return true
}

// GetUserID extracts and validates user ID from context
func GetUserID(ctx *gin.Context) (uint, bool) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": i18n.T("auth.user_not_authenticated", utils.GetUserLangFromRedis(ctx)),
		})
		return 0, false
	}
	return userID.(uint), true
}

// ParseIDParam parses ID parameter from URL
func ParseIDParam(ctx *gin.Context, paramName string) (uint, bool) {
	idStr := ctx.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": i18n.T("validation.invalid_query_parameters", utils.GetUserLangFromRedis(ctx)),
		})
		return 0, false
	}
	return uint(id), true
}

// GetPaginationParams extracts pagination parameters from query string
func GetPaginationParams(ctx *gin.Context) (page, pageSize int) {
	page = 1
	pageSize = 20

	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, parseErr := strconv.Atoi(pageStr); parseErr == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := ctx.Query("page_size"); pageSizeStr != "" {
		if ps, parseErr := strconv.Atoi(pageSizeStr); parseErr == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	return page, pageSize
}
