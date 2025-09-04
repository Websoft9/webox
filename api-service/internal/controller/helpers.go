package controller

import (
	"net/http"
	"strconv"

	"api-service/pkg/i18n"
	"api-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindAndValidateRequest binds JSON request and validates it
func BindAndValidateRequest(ctx *gin.Context, req interface{}, validator *validator.Validate, log logger.Logger, i18n *i18n.I18n) bool {
	// Bind request parameters
	if err := ctx.ShouldBindJSON(req); err != nil {
		log.ErrorContext(ctx.Request.Context(), "Invalid request format", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": i18n.T(ctx, "validation.invalid_request_format"),
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
			"message": i18n.T(ctx, "validation.request_validation_failed"),
			"error":   err.Error(),
		})
		return false
	}

	return true
}

// BindAndValidateQuery binds query parameters and validates them
func BindAndValidateQuery(ctx *gin.Context, req interface{}, validator *validator.Validate, log logger.Logger, i18n *i18n.I18n) bool {
	// Bind query parameters
	if err := ctx.ShouldBindQuery(req); err != nil {
		log.ErrorContext(ctx.Request.Context(), "Invalid query parameters", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": i18n.T(ctx, "validation.invalid_query_parameters"),
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
			"message": i18n.T(ctx, "validation.query_validation_failed"),
			"error":   err.Error(),
		})
		return false
	}

	return true
}

// GetUserID extracts and validates user ID from context
func GetUserID(ctx *gin.Context, i18n *i18n.I18n) (uint, bool) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return 0, false
	}
	return userID.(uint), true
}

// ParseIDParam parses ID parameter from URL
func ParseIDParam(ctx *gin.Context, paramName, errorKey string, i18n *i18n.I18n) (uint, bool) {
	idStr := ctx.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": i18n.T(ctx, errorKey),
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

// ResponseCreated sends a created response
func ResponseCreated(ctx *gin.Context, data interface{}, messageKey string, i18n *i18n.I18n) {
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"code":    http.StatusCreated,
		"message": i18n.T(ctx, messageKey),
		"data":    data,
	})
}

// ResponseBadRequest sends a bad request response
func ResponseBadRequest(ctx *gin.Context, err error, messageKey string, i18n *i18n.I18n) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"code":    http.StatusBadRequest,
		"message": i18n.T(ctx, messageKey),
		"error":   err.Error(),
	})
}

// ResponseOK sends a success response
func ResponseOK(ctx *gin.Context, data interface{}, messageKey string, i18n *i18n.I18n) {
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": i18n.T(ctx, messageKey),
		"data":    data,
	})
}

// ResponseNotFound sends a not found response
func ResponseNotFound(ctx *gin.Context, messageKey string, i18n *i18n.I18n) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"code":    http.StatusNotFound,
		"message": i18n.T(ctx, messageKey),
	})
}

// ResponseUnauthorized sends an unauthorized response
func ResponseUnauthorized(ctx *gin.Context, messageKey string, i18n *i18n.I18n) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"code":    http.StatusUnauthorized,
		"message": i18n.T(ctx, messageKey),
	})
}

// ResponseForbidden sends a forbidden response
func ResponseForbidden(ctx *gin.Context, messageKey string, i18n *i18n.I18n) {
	ctx.JSON(http.StatusForbidden, gin.H{
		"success": false,
		"code":    http.StatusForbidden,
		"message": i18n.T(ctx, messageKey),
	})
}

// ResponseInternalError sends an internal server error response
func ResponseInternalError(ctx *gin.Context, err error, messageKey string, log logger.Logger, i18n *i18n.I18n) {
	log.ErrorContext(ctx.Request.Context(), "Internal server error", logger.ErrorField(err))
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"code":    http.StatusInternalServerError,
		"message": i18n.T(ctx, messageKey),
		"error":   err.Error(),
	})
}
