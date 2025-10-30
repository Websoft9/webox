package request

import (
	"api-service/internal/dto/common"
	"mime/multipart"
)

// CreateTextSecretRequest represents the request to create a text type secret
type CreateTextSecretRequest struct {
	ResourceGroupID uint    `json:"resource_group_id" validate:"required"`                              // Resource group ID
	ResourceCode    *string `json:"resource_code" validate:"omitempty,max=128"`                         // Resource code (optional)
	Name            string  `json:"name" validate:"required,min=1,max=64"`                              // Secret name
	Description     *string `json:"description" validate:"omitempty,max=500"`                           // Secret description
	SecretText      string  `json:"secret_text" validate:"required"`                                    // Secret text value
	AuthorizedUsers []uint  `json:"authorized_users" validate:"omitempty,dive,required"`                // Authorized user IDs
	ExpiresAt       *string `json:"expires_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"` // Expiration time
}

// CreateAccountSecretRequest represents the request to create an account type secret
type CreateAccountSecretRequest struct {
	ResourceGroupID uint    `json:"resource_group_id" validate:"required"`                              // Resource group ID
	ResourceCode    *string `json:"resource_code" validate:"omitempty,max=128"`                         // Resource code (optional)
	Name            string  `json:"name" validate:"required,min=1,max=64"`                              // Secret name
	Description     *string `json:"description" validate:"omitempty,max=500"`                           // Secret description
	SecretUsername  string  `json:"secret_username" validate:"required"`                                // Username
	SecretPassword  string  `json:"secret_password" validate:"required"`                                // Password
	AuthorizedUsers []uint  `json:"authorized_users" validate:"omitempty,dive,required"`                // Authorized user IDs
	ExpiresAt       *string `json:"expires_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"` // Expiration time
}

// CreateFileSecretRequest represents the request to create a file type secret
type CreateFileSecretRequest struct {
	ResourceGroupID uint                  `form:"resource_group_id" binding:"required"` // Resource group ID
	ResourceCode    *string               `form:"resource_code" binding:"omitempty"`    // Resource code (optional)
	Name            string                `form:"name" binding:"required,min=1,max=64"` // Secret name
	Description     *string               `form:"description" binding:"omitempty"`      // Secret description
	SecretFile      *multipart.FileHeader `form:"secret_file" binding:"required"`       // File data
	SecretPassword  *string               `form:"secret_password" binding:"omitempty"`  // Private key password (optional)
	AuthorizedUsers []uint                `form:"authorized_users" binding:"omitempty"` // Authorized user IDs
	ExpiresAt       *string               `form:"expires_at" binding:"omitempty"`       // Expiration time
}

// UpdateSecretRequest represents the request to update a secret
type UpdateSecretRequest struct {
	ResourceGroupID *uint   `json:"resource_group_id" validate:"omitempty"`                             // Resource group ID
	Name            *string `json:"name" validate:"omitempty,min=1,max=64"`                             // Secret name
	Description     *string `json:"description" validate:"omitempty,max=500"`                           // Secret description
	AuthorizedUsers *[]uint `json:"authorized_users" validate:"omitempty,dive,required"`                // Authorized user IDs (覆盖模式), nil表示不更新
	ExpiresAt       *string `json:"expires_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"` // Expiration time
}

// ListSecretsRequest represents the request to list secrets
type ListSecretsRequest struct {
	common.BaseListRequest
	ResourceCode *string `form:"resource_code" binding:"omitempty"` // Filter by resource code
	common.TimeRangeRequest
}

// CreateReferenceRequest represents the request to create a secret reference
type CreateReferenceRequest struct {
	SecretID     uint   `json:"secret_id" validate:"required"`     // Secret ID
	ResourceCode string `json:"resource_code" validate:"required"` // Resource code
}
