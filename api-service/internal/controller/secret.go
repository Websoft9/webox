package controller

import (
	response "api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// SecretKeyController handles secret key related HTTP requests
type SecretKeyController struct {
	secretKeyService service.SecretKeyService
	logger           logger.Logger
	i18n             *i18n.I18n
	validator        *validator.Validate
}

// NewSecretKeyController creates a new secret key controller
func NewSecretKeyController(
	secretKeyService service.SecretKeyService,
	logger logger.Logger,
	i18n *i18n.I18n,
	validator *validator.Validate,
) *SecretKeyController {
	return &SecretKeyController{
		secretKeyService: secretKeyService,
		logger:           logger,
		i18n:             i18n,
		validator:        validator,
	}
}

// CreateSecretKey creates a new secret key
// @Summary Create secret key
// @Description Create a new secret key with encryption
// @Tags Secret Keys
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.SecretKeyCreateRequest true "Create secret key request"
// @Success 201 {object} common.APIResponse{data=response.SecretKeyResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 422 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secret-keys [post]
func (c *SecretKeyController) CreateSecretKey(ctx *gin.Context) {
	// Parse request parameters
	var req request.SecretKeyCreateRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling create secret key request", logger.Uint("userID", currentUserID))

	// Call service layer to create secret key
	secretKey, err := c.secretKeyService.CreateSecretKey(ctx.Request.Context(), &req, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to create secret key", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Secret key created successfully", logger.Uint("userID", currentUserID))
	response.SuccessWithData(ctx, secretKey)
}

// GetSecretKey retrieves a secret key by ID
// @Summary Get secret key
// @Description Get secret key information by ID
// @Tags Secret Keys
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Secret key ID"
// @Success 200 {object} common.APIResponse{data=response.SecretKeyResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secret-keys/{id} [get]
func (c *SecretKeyController) GetSecretKey(ctx *gin.Context) {
	// Parse path parameter
	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling get secret key request", logger.Uint("userID", currentUserID), logger.Uint("secretKeyID", id))

	// Call service layer to get secret key
	secretKey, err := c.secretKeyService.GetSecretKey(ctx.Request.Context(), id, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get secret key", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Secret key retrieved successfully", logger.Uint("userID", currentUserID))
	response.SuccessWithData(ctx, secretKey)
}

// GetSecretKeyValue retrieves the decrypted value of a secret key
// @Summary Get secret key value
// @Description Get the decrypted value of a secret key
// @Tags Secret Keys
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Secret key ID"
// @Success 200 {object} common.APIResponse{data=response.SecretKeyValueResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 422 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secret-keys/{id}/value [get]
func (c *SecretKeyController) GetSecretKeyValue(ctx *gin.Context) {
	// Parse path parameter
	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling get secret key value request", logger.Uint("userID", currentUserID), logger.Uint("secretKeyID", id))

	// Call service layer to get secret key value
	secretKeyValue, err := c.secretKeyService.GetSecretKeyValue(ctx.Request.Context(), id, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get secret key value", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Secret key value retrieved successfully", logger.Uint("userID", currentUserID))
	response.SuccessWithData(ctx, secretKeyValue)
}

// UpdateSecretKey updates an existing secret key
// @Summary Update secret key
// @Description Update an existing secret key
// @Tags Secret Keys
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Secret key ID"
// @Param request body request.SecretKeyUpdateRequest true "Update secret key request"
// @Success 200 {object} common.APIResponse{data=response.SecretKeyResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 422 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secret-keys/{id} [put]
func (c *SecretKeyController) UpdateSecretKey(ctx *gin.Context) {
	// Parse path parameter
	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}

	// Parse request parameters
	var req request.SecretKeyUpdateRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}
	c.logger.InfoContext(ctx, "Handling update secret key request", logger.Uint("userID", currentUserID), logger.Uint("secretKeyID", id))

	// Call service layer to update secret key
	secretKey, err := c.secretKeyService.UpdateSecretKey(ctx.Request.Context(), id, currentUserID, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to update secret key", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Secret key updated successfully", logger.Uint("userID", currentUserID))
	response.SuccessWithData(ctx, secretKey)
}

// DeleteSecretKey deletes a secret key
// @Summary Delete secret key
// @Description Delete a secret key by ID
// @Tags Secret Keys
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Secret key ID"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secret-keys/{id} [delete]
func (c *SecretKeyController) DeleteSecretKey(ctx *gin.Context) {
	// Parse path parameter
	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling delete secret key request", logger.Uint("userID", currentUserID), logger.Uint("secretKeyID", id))

	// Call service layer to delete secret key
	err := c.secretKeyService.DeleteSecretKey(ctx.Request.Context(), id, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to delete secret key", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Secret key deleted successfully", logger.Uint("userID", currentUserID))
	response.Success(ctx)
}

// ListSecretKeys retrieves secret keys with pagination and filtering
// @Summary List secret keys
// @Description Get secret keys with pagination and filtering
// @Tags Secret Keys
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param key_type query string false "Key type filter" Enums(API_KEY,DATABASE,SSH,CERTIFICATE,CUSTOM)
// @Success 200 {object} common.APIResponse{data=response.SecretKeyListResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secret-keys [get]
func (c *SecretKeyController) ListSecretKeys(ctx *gin.Context) {
	// Bind request parameters
	var req request.SecretKeyQueryRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling list secret keys request", logger.Uint("userID", currentUserID))

	// Call service layer to get secret keys
	secretKeys, err := c.secretKeyService.ListSecretKeys(ctx.Request.Context(), &req, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to list secret keys", logger.Uint("user_id", currentUserID), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Secret keys retrieved successfully", logger.Uint("user_id", currentUserID))
	response.SuccessWithData(ctx, secretKeys)
}

// ExportSecretKeys exports secret keys in specified format
// @Summary Export secret keys
// @Description Export secret keys in CSV, JSON, or Excel format
// @Tags Secret Keys
// @Security BearerAuth
// @Accept json
// @Produce application/octet-stream
// @Param format query string false "Export format" Enums(csv,json,excel) default(csv)
// @Success 200 {file} file "Exported file"
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secret-keys/export [get]
func (c *SecretKeyController) ExportSecretKeys(ctx *gin.Context) {
	// Bind request parameters
	var req request.SecretKeyExportRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling export secret keys request", logger.Uint("userID", currentUserID))

	// Call service layer to export secret keys
	data, filename, err := c.secretKeyService.ExportSecretKeys(ctx.Request.Context(), &req, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to export secret keys", logger.Uint("user_id", currentUserID), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	// Set response headers for file download
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Length", strconv.Itoa(len(data)))

	c.logger.InfoContext(ctx, "Secret keys exported successfully", logger.Uint("user_id", currentUserID))
	ctx.Writer.WriteHeader(http.StatusOK)
	_, _ = ctx.Writer.Write(data)
}
