package request

import (
	"time"

	"api-service/internal/dto/common"
	"api-service/internal/model"
	"api-service/pkg/errors"
)

// SecretKeyCreateRequest represents the request to create a secret key
type SecretKeyCreateRequest struct {
	Name            string                 `json:"name" binding:"required,min=1,max=64"`
	KeyType         model.SecretKeyType    `json:"key_type" binding:"required,oneof=SECRET_KEY ACCOUNT FILE"`
	Description     *string                `json:"description" binding:"omitempty,max=500"`
	CustomFields    map[string]interface{} `json:"custom_fields" binding:"omitempty"`
	ResourceGroupID *uint                  `json:"resource_group_id" binding:"omitempty"`
	ExpiresAt       *time.Time             `json:"expires_at" binding:"omitempty"`
	AuthorizedUsers []uint                 `json:"authorized_users" binding:"omitempty"`
}

// Validate validates the custom_fields based on key_type
func (r *SecretKeyCreateRequest) Validate() error {
	if r.CustomFields == nil {
		return nil
	}

	switch r.KeyType {
	case model.SecretKeyTypeSecretKey:
		return r.validateSecretKeyFields()
	case model.SecretKeyTypeAccount:
		return r.validateAccountFields()
	case model.SecretKeyTypeFile:
		return r.validateFileFields()
	default:
		return errors.NewAppError(errors.CodeValidationFailed)
	}
}

// validateSecretKeyFields validates custom_fields for SECRET_KEY type
func (r *SecretKeyCreateRequest) validateSecretKeyFields() error {
	secretKey, ok := r.CustomFields["secret_key"]
	if !ok {
		return errors.NewAppError(errors.CodeValidationFailed)
	}

	secretKeyStr, ok := secretKey.(string)
	if !ok || secretKeyStr == "" {
		return errors.NewAppError(errors.CodeValidationFailed)
	}

	return nil
}

// validateAccountFields validates custom_fields for ACCOUNT type
func (r *SecretKeyCreateRequest) validateAccountFields() error {
	username, hasUsername := r.CustomFields["username"]
	password, hasPassword := r.CustomFields["password"]

	if !hasUsername {
		return errors.NewAppError(errors.CodeValidationFailed)
	}

	if !hasPassword {
		return errors.NewAppError(errors.CodeValidationFailed)
	}

	usernameStr, ok := username.(string)
	if !ok || usernameStr == "" {
		return errors.NewAppError(errors.CodeValidationFailed)
	}

	passwordStr, ok := password.(string)
	if !ok || passwordStr == "" {
		return errors.NewAppError(errors.CodeValidationFailed)
	}

	return nil
}

// validateFileFields validates custom_fields for FILE type
func (r *SecretKeyCreateRequest) validateFileFields() error {
	filename, hasFilename := r.CustomFields["filename"]

	if !hasFilename {
		return errors.NewAppError(errors.CodeValidationFailed)
	}

	filenameStr, ok := filename.(string)
	if !ok || filenameStr == "" {
		return errors.NewAppError(errors.CodeValidationFailed)
	}

	// password is optional for FILE type
	if password, hasPassword := r.CustomFields["password"]; hasPassword {
		passwordStr, ok := password.(string)
		if !ok || passwordStr == "" {
			return errors.NewAppError(errors.CodeValidationFailed)
		}
	}

	return nil
}

// SecretKeyUpdateRequest represents the request to update a secret key
type SecretKeyUpdateRequest struct {
	KeyType      model.SecretKeyType    `json:"key_type" binding:"required,oneof=SECRET_KEY ACCOUNT FILE"`
	CustomFields map[string]interface{} `json:"custom_fields" binding:"omitempty"`
}

// Validate validates the custom_fields based on key_type for update request
func (r *SecretKeyUpdateRequest) Validate() error {
	if r.CustomFields == nil {
		return nil
	}

	// Reuse the same validation logic
	createReq := &SecretKeyCreateRequest{
		KeyType:      r.KeyType,
		CustomFields: r.CustomFields,
	}
	return createReq.Validate()
}

// SecretKeyQueryRequest represents the request to query secret keys
type SecretKeyQueryRequest struct {
	common.PaginationRequest
	KeyType *model.SecretKeyType `form:"key_type" binding:"omitempty,oneof=SECRET_KEY ACCOUNT FILE"`
}

// SecretKeyExportRequest represents the request to export secret keys
type SecretKeyExportRequest struct {
	Format string `form:"format" binding:"required,oneof=csv json excel"`
}

// APITokenCreateRequest represents the request to create an API token
type APITokenCreateRequest struct {
	Name        string     `json:"name" binding:"required,min=1,max=100"`
	Description *string    `json:"description" binding:"omitempty,max=500"`
	ExpiresAt   *time.Time `json:"expires_at" binding:"omitempty"`
	Scopes      []string   `json:"scopes" binding:"omitempty"`
}

// APITokenUpdateRequest represents the request to update an API token
type APITokenUpdateRequest struct {
	Name        string     `json:"name" binding:"omitempty,min=1,max=100"`
	Description *string    `json:"description" binding:"omitempty,max=500"`
	ExpiresAt   *time.Time `json:"expires_at" binding:"omitempty"`
	IsActive    *bool      `json:"is_active" binding:"omitempty"`
}

// APITokenQueryRequest represents the request to query API tokens
type APITokenQueryRequest struct {
	common.PaginationRequest
	IsActive *bool `form:"is_active" binding:"omitempty"`
}

// TwoFactorSetupRequest represents the request to setup 2FA
type TwoFactorSetupRequest struct {
	Method string `json:"method" binding:"required,oneof=totp sms email"`
}

// TwoFactorVerifyRequest represents the request to verify 2FA code
type TwoFactorVerifyRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

// TwoFactorDisableRequest represents the request to disable 2FA
type TwoFactorDisableRequest struct {
	Code     string `json:"code" binding:"required,len=6"`
	Password string `json:"password" binding:"required"`
}

// SecurityEventQueryRequest represents the request to query security events
type SecurityEventQueryRequest struct {
	common.PaginationRequest
	EventType *string    `form:"event_type" binding:"omitempty"`
	StartTime *time.Time `form:"start_time" binding:"omitempty"`
	EndTime   *time.Time `form:"end_time" binding:"omitempty"`
}

// PasswordChangeRequest represents the request to change password
type PasswordChangeRequest struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=128"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// PasswordResetRequest represents the request to reset password
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// PasswordResetConfirmRequest represents the request to confirm password reset
type PasswordResetConfirmRequest struct {
	Token           string `json:"token" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=128"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}
