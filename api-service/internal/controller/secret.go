package controller

import (
	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	pkg_response "api-service/pkg/response"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SecretKeyController handles secret key related HTTP requests
type SecretKeyController struct {
	secretKeyService service.SecretKeyService
	logger           logger.Logger
	i18n             *i18n.I18n
}

// NewSecretKeyController creates a new secret key controller
func NewSecretKeyController(
	secretKeyService service.SecretKeyService,
	logger logger.Logger,
	i18n *i18n.I18n,
) *SecretKeyController {
	return &SecretKeyController{
		secretKeyService: secretKeyService,
		logger:           logger,
		i18n:             i18n,
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
// @Success 201 {object} response.APIResponse{data=response.SecretKeyResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 422 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/secret-keys [post]
func (c *SecretKeyController) CreateSecretKey(ctx *gin.Context) {
	// Parse request parameters
	var req request.SecretKeyCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WarnContext(ctx, "Invalid request parameters", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationFailed, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	// Get current logged in user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ResponseUnauthorized(ctx, "auth.user_not_authenticated", c.i18n)
		return
	}

	c.logger.InfoContext(ctx, "Handling create secret key request", logger.Uint("userID", userID.(uint)))

	// Call service layer to create secret key
	secretKey, err := c.secretKeyService.CreateSecretKey(ctx.Request.Context(), &req, userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to create secret key", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Secret key created successfully", logger.Uint("userID", userID.(uint)))
	pkg_response.Success(ctx, c.i18n.T(ctx, "secret_key.create_success"), secretKey)
}

// GetSecretKey retrieves a secret key by ID
// @Summary Get secret key
// @Description Get secret key information by ID
// @Tags Secret Keys
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Secret key ID"
// @Success 200 {object} response.APIResponse{data=response.SecretKeyResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/secret-keys/{id} [get]
func (c *SecretKeyController) GetSecretKey(ctx *gin.Context) {
	// Parse path parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.logger.WarnContext(ctx, "Invalid secret key ID", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationFailed, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	// Get current logged in user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ResponseUnauthorized(ctx, "auth.user_not_authenticated", c.i18n)
		return
	}

	c.logger.InfoContext(ctx, "Handling get secret key request", logger.Uint("userID", userID.(uint)), logger.Uint("secretKeyID", uint(id)))

	// Call service layer to get secret key
	secretKey, err := c.secretKeyService.GetSecretKey(ctx.Request.Context(), uint(id), userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get secret key", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Secret key retrieved successfully", logger.Uint("userID", userID.(uint)))
	pkg_response.Success(ctx, c.i18n.T(ctx, "secret_key.get_success"), secretKey)
}

// GetSecretKeyValue retrieves the decrypted value of a secret key
// @Summary Get secret key value
// @Description Get the decrypted value of a secret key
// @Tags Secret Keys
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Secret key ID"
// @Success 200 {object} response.APIResponse{data=response.SecretKeyValueResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 422 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/secret-keys/{id}/value [get]
func (c *SecretKeyController) GetSecretKeyValue(ctx *gin.Context) {
	// Parse path parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.logger.WarnContext(ctx, "Invalid secret key ID", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationFailed, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	// Get current logged in user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ResponseUnauthorized(ctx, "auth.user_not_authenticated", c.i18n)
		return
	}

	c.logger.InfoContext(ctx, "Handling get secret key value request", logger.Uint("userID", userID.(uint)), logger.Uint("secretKeyID", uint(id)))

	// Call service layer to get secret key value
	secretKeyValue, err := c.secretKeyService.GetSecretKeyValue(ctx.Request.Context(), uint(id), userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get secret key value", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Secret key value retrieved successfully", logger.Uint("userID", userID.(uint)))
	pkg_response.Success(ctx, c.i18n.T(ctx, "secret_key.get_value_success"), secretKeyValue)
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
// @Success 200 {object} response.APIResponse{data=response.SecretKeyResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 422 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/secret-keys/{id} [put]
func (c *SecretKeyController) UpdateSecretKey(ctx *gin.Context) {
	// Parse path parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.logger.WarnContext(ctx, "Invalid secret key ID", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationFailed, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	// Parse request parameters
	var req request.SecretKeyUpdateRequest
	if bindErr := ctx.ShouldBindJSON(&req); bindErr != nil {
		c.logger.WarnContext(ctx, "Invalid request parameters", logger.ErrorField(bindErr))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationFailed, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	// Get current logged in user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ResponseUnauthorized(ctx, "auth.user_not_authenticated", c.i18n)
		return
	}

	c.logger.InfoContext(ctx, "Handling update secret key request", logger.Uint("userID", userID.(uint)), logger.Uint("secretKeyID", uint(id)))

	// Call service layer to update secret key
	secretKey, err := c.secretKeyService.UpdateSecretKey(ctx.Request.Context(), uint(id), userID.(uint), &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to update secret key", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Secret key updated successfully", logger.Uint("userID", userID.(uint)))
	pkg_response.Success(ctx, c.i18n.T(ctx, "secret_key.update_success"), secretKey)
}

// DeleteSecretKey deletes a secret key
// @Summary Delete secret key
// @Description Delete a secret key by ID
// @Tags Secret Keys
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Secret key ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/secret-keys/{id} [delete]
func (c *SecretKeyController) DeleteSecretKey(ctx *gin.Context) {
	// Parse path parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.logger.WarnContext(ctx, "Invalid secret key ID", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationFailed, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	// Get current logged in user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ResponseUnauthorized(ctx, "auth.user_not_authenticated", c.i18n)
		return
	}

	c.logger.InfoContext(ctx, "Handling delete secret key request", logger.Uint("userID", userID.(uint)), logger.Uint("secretKeyID", uint(id)))

	// Call service layer to delete secret key
	err = c.secretKeyService.DeleteSecretKey(ctx.Request.Context(), uint(id), userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to delete secret key", logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Secret key deleted successfully", logger.Uint("userID", userID.(uint)))
	pkg_response.Success(ctx, c.i18n.T(ctx, "secret_key.delete_success"), nil)
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
// @Success 200 {object} response.APIResponse{data=response.SecretKeyListResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/secret-keys [get]
func (c *SecretKeyController) ListSecretKeys(ctx *gin.Context) {
	// Bind request parameters
	var req request.SecretKeyQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.logger.WarnContext(ctx, "Secret key list request parameter binding failed", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationFailed, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	// Get current logged in user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ResponseUnauthorized(ctx, "auth.user_not_authenticated", c.i18n)
		return
	}

	c.logger.InfoContext(ctx, "Handling list secret keys request", logger.Uint("userID", userID.(uint)))

	// Call service layer to get secret keys
	secretKeys, err := c.secretKeyService.ListSecretKeys(ctx.Request.Context(), &req, userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to list secret keys", logger.Uint("user_id", userID.(uint)), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Secret keys retrieved successfully", logger.Uint("user_id", userID.(uint)))
	pkg_response.Success(ctx, c.i18n.T(ctx, "secret_key.list_success"), secretKeys)
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
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/secret-keys/export [get]
func (c *SecretKeyController) ExportSecretKeys(ctx *gin.Context) {
	// Bind request parameters
	var req request.SecretKeyExportRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.logger.WarnContext(ctx, "Secret key export request parameter binding failed", logger.ErrorField(err))
		errors.HandleError(ctx, errors.NewAppError(errors.CodeValidationFailed, c.i18n.T(ctx, "common.validation_failed")))
		return
	}

	// Get current logged in user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ResponseUnauthorized(ctx, "auth.user_not_authenticated", c.i18n)
		return
	}

	c.logger.InfoContext(ctx, "Handling export secret keys request", logger.Uint("userID", userID.(uint)))

	// Call service layer to export secret keys
	data, filename, err := c.secretKeyService.ExportSecretKeys(ctx.Request.Context(), &req, userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to export secret keys", logger.Uint("user_id", userID.(uint)), logger.ErrorField(err))
		errors.HandleError(ctx, err)
		return
	}

	// Set response headers for file download
	contentType := getContentTypeByFormat(req.Format)
	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Length", strconv.Itoa(len(data)))

	c.logger.InfoContext(ctx, "Secret keys exported successfully", logger.Uint("user_id", userID.(uint)))
	ctx.Writer.WriteHeader(http.StatusOK)
	_, _ = ctx.Writer.Write(data)
}

// getContentTypeByFormat returns the correct content type for the given format
func getContentTypeByFormat(format string) string {
	switch format {
	case "json":
		return "application/json"
	case "excel":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	default:
		return "text/csv"
	}
}
