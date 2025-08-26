package controller

import (
	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// TwoFactorController two-factor authentication controller
type TwoFactorController struct {
	twoFactorService service.TwoFactorService
	validator        *validator.Validate
	logger           logger.Logger
	i18n             *i18n.I18n
}

// NewTwoFactorController creates a new two-factor authentication controller instance
func NewTwoFactorController(
	twoFactorService service.TwoFactorService,
	validator *validator.Validate,
	logger logger.Logger,
	i18n *i18n.I18n,
) *TwoFactorController {
	return &TwoFactorController{
		twoFactorService: twoFactorService,
		validator:        validator,
		logger:           logger,
		i18n:             i18n,
	}
}

// EnableTOTP enables TOTP two-factor authentication
// @Summary Enable TOTP 2FA
// @Description Enable TOTP two-factor authentication for user
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=response.TOTPSetupResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/2fa/totp/enable [post]
func (c *TwoFactorController) EnableTOTP(ctx *gin.Context) {
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

	// Enable TOTP
	setup, err := c.twoFactorService.EnableTOTP(ctx.Request.Context(), userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to enable TOTP", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "two_factor.enable_totp_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "two_factor.enable_totp_success"),
		"data":    setup,
	})
}

// ConfirmTOTP confirms TOTP setup
// @Summary Confirm TOTP setup
// @Description Confirm TOTP setup with verification code
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Param request body request.ConfirmTOTPRequest true "Confirm TOTP request"
// @Success 200 {object} response.APIResponse{data=response.TOTPConfirmResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/2fa/totp/confirm [post]
func (c *TwoFactorController) ConfirmTOTP(ctx *gin.Context) {
	var req request.ConfirmTOTPRequest

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

	// Confirm TOTP
	result, err := c.twoFactorService.ConfirmTOTP(ctx.Request.Context(), userID.(uint), req.Code)
	if err != nil {
		if err.Error() == "invalid code" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"code":    http.StatusBadRequest,
				"message": c.i18n.T(ctx, "two_factor.invalid_code"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to confirm TOTP", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "two_factor.confirm_totp_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "two_factor.confirm_totp_success"),
		"data":    result,
	})
}

// DisableTOTP disables TOTP two-factor authentication
// @Summary Disable TOTP 2FA
// @Description Disable TOTP two-factor authentication for user
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Param request body request.DisableTOTPRequest true "Disable TOTP request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/2fa/totp/disable [post]
func (c *TwoFactorController) DisableTOTP(ctx *gin.Context) {
	var req request.DisableTOTPRequest

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

	// Disable TOTP
	err := c.twoFactorService.DisableTOTP(ctx.Request.Context(), userID.(uint), req.Code)
	if err != nil {
		if err.Error() == "invalid code" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"code":    http.StatusBadRequest,
				"message": c.i18n.T(ctx, "two_factor.invalid_code"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to disable TOTP", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "two_factor.disable_totp_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "two_factor.disable_totp_success"),
	})
}

// EnableEmailTwoFactor enables email two-factor authentication
// @Summary Enable Email 2FA
// @Description Enable email two-factor authentication for user
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Param request body request.EnableEmailTwoFactorRequest true "Enable email 2FA request"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/2fa/email/enable [post]
func (c *TwoFactorController) EnableEmailTwoFactor(ctx *gin.Context) {
	var req request.EnableEmailTwoFactorRequest

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

	// Enable email two-factor
	err := c.twoFactorService.EnableEmailTwoFactor(ctx.Request.Context(), userID.(uint), req.Email)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to enable email 2FA", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "two_factor.enable_email_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "two_factor.enable_email_success"),
	})
}

// DisableEmailTwoFactor disables email two-factor authentication
// @Summary Disable Email 2FA
// @Description Disable email two-factor authentication for user
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/2fa/email/disable [post]
func (c *TwoFactorController) DisableEmailTwoFactor(ctx *gin.Context) {
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

	// Disable email two-factor
	err := c.twoFactorService.DisableEmailTwoFactor(ctx.Request.Context(), userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to disable email 2FA", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "two_factor.disable_email_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "two_factor.disable_email_success"),
	})
}

// SendEmailCode sends email verification code
// @Summary Send email verification code
// @Description Send verification code to user's email for 2FA
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/2fa/email/send-code [post]
func (c *TwoFactorController) SendEmailCode(ctx *gin.Context) {
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

	// Send email code
	err := c.twoFactorService.SendEmailCode(ctx.Request.Context(), userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to send email code", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "two_factor.send_email_code_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "two_factor.send_email_code_success"),
	})
}

// VerifyTwoFactor verifies two-factor authentication code
// @Summary Verify 2FA code
// @Description Verify two-factor authentication code
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Param request body request.VerifyTwoFactorRequest true "Verify 2FA request"
// @Success 200 {object} response.APIResponse{data=response.TwoFactorVerificationResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /api/v1/2fa/verify [post]
func (c *TwoFactorController) VerifyTwoFactor(ctx *gin.Context) {
	var req request.VerifyTwoFactorRequest

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

	// Verify two-factor code
	result, err := c.twoFactorService.VerifyTwoFactor(ctx.Request.Context(), req.UserID, req.Code, req.Method)
	if err != nil {
		if err.Error() == "invalid code" || err.Error() == "code expired" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": c.i18n.T(ctx, "two_factor.invalid_code"),
			})
			return
		}

		c.logger.ErrorContext(ctx.Request.Context(), "Failed to verify 2FA code", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "two_factor.verify_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "two_factor.verify_success"),
		"data":    result,
	})
}

// GetTwoFactorStatus gets user's two-factor authentication status
// @Summary Get 2FA status
// @Description Get user's two-factor authentication status
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=response.TwoFactorStatusResponse}
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/2fa/status [get]
func (c *TwoFactorController) GetTwoFactorStatus(ctx *gin.Context) {
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

	// Get two-factor status
	status, err := c.twoFactorService.GetTwoFactorStatus(ctx.Request.Context(), userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to get 2FA status", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "two_factor.status_failed"),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "common.success"),
		"data":    status,
	})
}

// GenerateBackupCodes generates backup codes for two-factor authentication
// @Summary Generate backup codes
// @Description Generate backup codes for two-factor authentication recovery
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=response.BackupCodesResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/2fa/backup-codes [post]
func (c *TwoFactorController) GenerateBackupCodes(ctx *gin.Context) {
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

	// Generate backup codes
	codes, err := c.twoFactorService.GenerateBackupCodes(ctx.Request.Context(), userID.(uint))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to generate backup codes", logger.ErrorField(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"message": c.i18n.T(ctx, "two_factor.backup_codes_failed"),
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"message": c.i18n.T(ctx, "two_factor.backup_codes_success"),
		"data":    codes,
	})
}
