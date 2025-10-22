package request

import (
	"mime/multipart"
	"time"

	"api-service/internal/dto/common"
	"api-service/internal/model"
	"api-service/pkg/errors"
)

// SecretKeyCreateTextRequest represents the request to create a secret key
type SecretKeyCreateTextRequest struct {
	Name            string                 `json:"name" binding:"required,min=1,max=64"`
	Description     *string                `json:"description" binding:"omitempty,max=500"`
	SecretFields    map[string]interface{} `json:"secret_fields" binding:"required"`
	ResourceGroupID *uint                  `json:"resource_group_id" binding:"omitempty"`
	ExpiresAt       *time.Time             `json:"expires_at" binding:"omitempty"`
	AuthorizedUsers []uint                 `json:"authorized_users" binding:"omitempty"`
	ResourceCode    *string                `json:"resource_code" binding:"omitempty,max=64"`
}

// Validate validates the secret_fields for text-based secret keys
func (r *SecretKeyCreateTextRequest) Validate() error {
	if r.SecretFields == nil {
		return errors.NewAppError(errors.CodeValidationFailed)
	}

	// Check if it's a TEXT type (secret_key field exists)
	if secretKey, ok := r.SecretFields["secret_key"]; ok {
		secretKeyStr, ok := secretKey.(string)
		if !ok || secretKeyStr == "" {
			return errors.NewAppError(errors.CodeValidationFailed)
		}
		return nil
	}

	// Check if it's an ACCOUNT type (username and password fields exist)
	username, hasUsername := r.SecretFields["username"]
	password, hasPassword := r.SecretFields["password"]

	if hasUsername && hasPassword {
		usernameStr, ok1 := username.(string)
		passwordStr, ok2 := password.(string)
		if !ok1 || usernameStr == "" || !ok2 || passwordStr == "" {
			return errors.NewAppError(errors.CodeValidationFailed)
		}
		return nil
	}

	return errors.NewAppError(errors.CodeValidationFailed)
}

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

// GetKeyType determines the key type based on secret_fields content
func (r *SecretKeyCreateTextRequest) GetKeyType() model.SecretKeyType {
	if r.SecretFields == nil {
		return model.SecretKeyTypeText
	}

	// Check for secret_key field (TEXT type)
	if _, ok := r.SecretFields["secret_key"]; ok {
		return model.SecretKeyTypeText
	}

	// Check for username/password fields (ACCOUNT type)
	if _, hasUsername := r.SecretFields["username"]; hasUsername {
		if _, hasPassword := r.SecretFields["password"]; hasPassword {
			return model.SecretKeyTypeAccount
		}
	}

	return model.SecretKeyTypeText
}

// SecretKeyCreateFileRequest represents the request to create a file-based secret key
type SecretKeyCreateFileRequest struct {
	Name            string                 `form:"name" binding:"required,min=1,max=64"`
	Description     *string                `form:"description" binding:"omitempty,max=500"`
	SecretFields    map[string]interface{} `form:"-"` // Will be populated from form data
	File            *multipart.FileHeader  `form:"file" binding:"required"`
	ResourceGroupID *uint                  `form:"resource_group_id" binding:"omitempty"`
	ExpiresAt       *time.Time             `form:"expires_at" binding:"omitempty"`
	AuthorizedUsers []uint                 `form:"authorized_users" binding:"omitempty"`
	ResourceCode    *string                `form:"resource_code" binding:"omitempty,max=64"`
}

// Validate validates the file-based secret key request
func (r *SecretKeyCreateFileRequest) Validate() error {
	if r.File == nil {
		return errors.NewAppError(errors.CodeValidationFailed)
	}
	return nil
}

// GetKeyType returns FILE type for file-based secret keys
func (r *SecretKeyCreateFileRequest) GetKeyType() model.SecretKeyType {
	return model.SecretKeyTypeFile
}

// BuildSecretFields builds the secret_fields map from form data
func (r *SecretKeyCreateFileRequest) BuildSecretFields() {
	if r.SecretFields == nil {
		r.SecretFields = make(map[string]interface{})
	}
	r.SecretFields["filename"] = r.File.Filename
}

// SecretAssignRequest represents the request to assign a secret to a resource
type SecretAssignRequest struct {
	SecretID     uint   `json:"secret_id" binding:"required,gt=0"`
	ResourceCode string `json:"resource_code" binding:"required,min=1,max=64"`
}

// SecretUnassignRequest represents the request to unassign a secret from a resource
type SecretUnassignRequest struct {
	SecretID     uint   `json:"secret_id" binding:"required,gt=0"`
	ResourceCode string `json:"resource_code" binding:"required,min=1,max=64"`
}

// SecretResourceQueryRequest represents the request to query resources associated with a secret
type SecretResourceQueryRequest struct {
	SecretID uint `form:"secret_id" binding:"required,gt=0"`
}
