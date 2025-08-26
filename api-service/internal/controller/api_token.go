package controller

import (
	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// APITokenController API token controller
type APITokenController struct {
	apiTokenService service.APITokenService
	validator       *validator.Validate
	logger          logger.Logger
	i18n            *i18n.I18n
}

// NewAPITokenController creates a new API token controller instance
func NewAPITokenController(
	apiTokenService service.APITokenService,
	validator *validator.Validate,
	logger logger.Logger,
	i18n *i18n.I18n,
) *APITokenController {
	return &APITokenController{
		apiTokenService: apiTokenService,
		validator:       validator,
		logger:          logger,
		i18n:            i18n,
	}
}

// CreateAPIToken creates a new API token
// @Summary Create API token
// @Description Create a new API token
// @Tags API Token Management
// @Accept json
// @Produce json
// @Param request body request.CreateAPITokenRequest true "Create API token request"
// @Success 201 {object} response.APIResponse{data=response.APITokenResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/api-tokens [post]
func (c *APITokenController) CreateAPIToken(ctx *gin.Context) {
	var req request.CreateAPITokenRequest

	// Bind request parameters
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid request format", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_request_format"),
			"error":   err.Error(),
		})
		return
	}

	// Validate request parameters
	if err := c.validator.Struct(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Request validation failed", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.request_validation_failed"),
			"error":   err.Error(),
		})
		return
	}

	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Create API token
	token, err := c.apiTokenService.CreateAPIToken(ctx.Request.Context(), &req, userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to create API token", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "api_token.create_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"code":    http.StatusCreated,
		"message": c.i18n.T(ctx, "api_token.create_success"),
		"data":    token,
	})
}

// GetAPIToken gets API token details
// @Summary Get API token details
// @Description Get API token details by ID
// @Tags API Token Management
// @Accept json
// @Produce json
// @Param id path int true "API Token ID"
// @Success 200 {object} response.APIResponse{data=response.APITokenResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/api-tokens/{id} [get]
func (c *APITokenController) GetAPIToken(ctx *gin.Context) {
	// Get token ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_token_id"),
		})
		return
	}

	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Get API token
	token, err := c.apiTokenService.GetAPIToken(ctx.Request.Context(), uint(id), userID.(uint))
	if err != nil {
		if err.Error() == "token not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    http.StatusNotFound,
				"message": c.i18n.T(ctx, "api_token.not_found"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to get API token", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "api_token.get_failed"),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "common.success"),
		"data":    token,
	})
}

// UpdateAPIToken updates API token
// @Summary Update API token
// @Description Update API token information
// @Tags API Token Management
// @Accept json
// @Produce json
// @Param id path int true "API Token ID"
// @Param request body request.UpdateAPITokenRequest true "Update API token request"
// @Success 200 {object} response.APIResponse{data=response.APITokenResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/api-tokens/{id} [put]
func (c *APITokenController) UpdateAPIToken(ctx *gin.Context) {
	// Get token ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_token_id"),
		})
		return
	}

	var req request.UpdateAPITokenRequest

	// Bind request parameters
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid request format", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_request_format"),
			"error":   err.Error(),
		})
		return
	}

	// Validate request parameters
	if err := c.validator.Struct(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Request validation failed", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.request_validation_failed"),
			"error":   err.Error(),
		})
		return
	}

	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Update API token
	token, err := c.apiTokenService.UpdateAPIToken(ctx.Request.Context(), uint(id), &req, userID.(uint))
	if err != nil {
		if err.Error() == "token not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    http.StatusNotFound,
				"message": c.i18n.T(ctx, "api_token.not_found"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to update API token", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "api_token.update_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "api_token.update_success"),
		"data":    token,
	})
}

// RevokeAPIToken revokes API token
// @Summary Revoke API token
// @Description Revoke specified API token
// @Tags API Token Management
// @Accept json
// @Produce json
// @Param id path int true "API Token ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/api-tokens/{id}/revoke [post]
func (c *APITokenController) RevokeAPIToken(ctx *gin.Context) {
	// Get token ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_token_id"),
		})
		return
	}

	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Revoke API token
	err = c.apiTokenService.RevokeAPIToken(ctx.Request.Context(), uint(id), userID.(uint))
	if err != nil {
		if err.Error() == "token not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    http.StatusNotFound,
				"message": c.i18n.T(ctx, "api_token.not_found"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to revoke API token", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "api_token.revoke_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "api_token.revoke_success"),
	})
}

// RefreshAPIToken refreshes API token
// @Summary Refresh API token
// @Description Refresh API token to extend expiration
// @Tags API Token Management
// @Accept json
// @Produce json
// @Param id path int true "API Token ID"
// @Success 200 {object} response.APIResponse{data=response.APITokenResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/api-tokens/{id}/refresh [post]
func (c *APITokenController) RefreshAPIToken(ctx *gin.Context) {
	// Get token ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_token_id"),
		})
		return
	}

	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Refresh API token
	token, err := c.apiTokenService.RefreshAPIToken(ctx.Request.Context(), uint(id), userID.(uint))
	if err != nil {
		if err.Error() == "token not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    http.StatusNotFound,
				"message": c.i18n.T(ctx, "api_token.not_found"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to refresh API token", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "api_token.refresh_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "api_token.refresh_success"),
		"data":    token,
	})
}

// ListAPITokens gets API token list
// @Summary Get API token list
// @Description Get paginated API token list
// @Tags API Token Management
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search keyword"
// @Param status query string false "Token status" Enums(active, revoked, expired)
// @Param start_time query string false "Start time" format(datetime)
// @Param end_time query string false "End time" format(datetime)
// @Success 200 {object} response.APIResponse{data=response.APITokenListResponse}
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/api-tokens [get]
func (c *APITokenController) ListAPITokens(ctx *gin.Context) {
	var req request.ListAPITokensRequest

	// Bind query parameters
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid query parameters", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_query_parameters"),
			"error":   err.Error(),
		})
		return
	}

	// Validate request parameters
	if err := c.validator.Struct(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Query validation failed", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.query_validation_failed"),
			"error":   err.Error(),
		})
		return
	}

	// Get current user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": c.i18n.T(ctx, "auth.user_not_authenticated"),
		})
		return
	}

	// Get API token list
	tokens, err := c.apiTokenService.ListAPITokens(ctx.Request.Context(), &req, userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to list API tokens", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "api_token.list_failed"),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "common.success"),
		"data":    tokens,
	})
}

// ValidateAPIToken validates API token
// @Summary Validate API token
// @Description Validate API token and return token information
// @Tags API Token Management
// @Accept json
// @Produce json
// @Param request body request.ValidateAPITokenRequest true "Validate API token request"
// @Success 200 {object} response.APIResponse{data=response.APITokenValidationResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /api/v1/api-tokens/validate [post]
func (c *APITokenController) ValidateAPIToken(ctx *gin.Context) {
	var req request.ValidateAPITokenRequest

	// Bind request parameters
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Invalid request format", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.invalid_request_format"),
			"error":   err.Error(),
		})
		return
	}

	// Validate request parameters
	if err := c.validator.Struct(&req); err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Request validation failed", logger.ErrorField(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"message": c.i18n.T(ctx, "validation.request_validation_failed"),
			"error":   err.Error(),
		})
		return
	}

	// Validate API token
	validation, err := c.apiTokenService.ValidateAPIToken(ctx.Request.Context(), req.Token)
	if err != nil {
		if err.Error() == "invalid token" || err.Error() == "token expired" || err.Error() == "token revoked" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": c.i18n.T(ctx, "api_token.invalid_token"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to validate API token", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "api_token.validate_failed"),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "api_token.validate_success"),
		"data":    validation,
	})
}
