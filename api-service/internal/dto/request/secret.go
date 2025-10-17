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
	KeyType         model.SecretKeyType    `json:"key_type" binding:"required,oneof=TEXT ACCOUNT FILE"`
	Description     *string                `json:"description" binding:"omitempty,max=500"`
	SecretFields    map[string]interface{} `json:"secret_fields" binding:"omitempty"`
	ResourceGroupID *uint                  `json:"resource_group_id" binding:"omitempty"`
	ExpiresAt       *time.Time             `json:"expires_at" binding:"omitempty"`
	AuthorizedUsers []uint                 `json:"authorized_users" binding:"omitempty"`
	ResourceCode    *string                `json:"resource_code" binding:"omitempty,max=64"`
}

// Validate validates the secret_fields based on key_type
func (r *SecretKeyCreateRequest) Validate() error {
	if r.SecretFields == nil {
		return nil
	}

	switch r.KeyType {
	case model.SecretKeyTypeText:
		return r.validateSecretKeyFields()
	case model.SecretKeyTypeAccount:
		return r.validateAccountFields()
	case model.SecretKeyTypeFile:
		return r.validateFileFields()
	default:
		return errors.NewAppError(errors.CodeValidationFailed)
	}
}

// validateSecretKeyFields validates secret_fields for SECRET_KEY type
func (r *SecretKeyCreateRequest) validateSecretKeyFields() error {
	secretKey, ok := r.SecretFields["secret_key"]
	if !ok {
		return errors.NewAppError(errors.CodeValidationFailed)
	}

	secretKeyStr, ok := secretKey.(string)
	if !ok || secretKeyStr == "" {
		return errors.NewAppError(errors.CodeValidationFailed)
	}

	return nil
}

// validateAccountFields validates secret_fields for ACCOUNT type
func (r *SecretKeyCreateRequest) validateAccountFields() error {
	username, hasUsername := r.SecretFields["username"]
	password, hasPassword := r.SecretFields["password"]

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

// validateFileFields validates secret_fields for FILE type
func (r *SecretKeyCreateRequest) validateFileFields() error {
	filename, hasFilename := r.SecretFields["filename"]

	if !hasFilename {
		return errors.NewAppError(errors.CodeValidationFailed)
	}

	filenameStr, ok := filename.(string)
	if !ok || filenameStr == "" {
		return errors.NewAppError(errors.CodeValidationFailed)
	}

	// password is optional for FILE type
	if password, hasPassword := r.SecretFields["password"]; hasPassword {
		passwordStr, ok := password.(string)
		if !ok || passwordStr == "" {
			return errors.NewAppError(errors.CodeValidationFailed)
		}
	}

	return nil
}

// SecretKeyUpdateRequest represents the request to update a secret key
type SecretKeyUpdateRequest struct {
	Name        *string    `json:"name" binding:"omitempty,min=1,max=64"`
	Description *string    `json:"description" binding:"omitempty,max=500"`
	ExpiresAt   *time.Time `json:"expires_at" binding:"omitempty"`
}

// SecretKeyQueryRequest represents the request to query secret keys
type SecretKeyQueryRequest struct {
	common.PaginationRequest
	KeyType      *model.SecretKeyType `form:"key_type" binding:"omitempty,oneof=TEXT ACCOUNT FILE"`
	Keyword      *string              `form:"keyword" binding:"omitempty,max=100"`
	ResourceCode *string              `form:"resource_code" binding:"omitempty,max=64"`
}

// SecretKeyExportRequest represents the request to export secret keys
type SecretKeyExportRequest struct {
	Format string `form:"format" binding:"required,oneof=csv json excel"`
}
