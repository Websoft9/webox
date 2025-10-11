package request

import (
	"api-service/internal/dto/common"
	"api-service/internal/model"
	"time"
)

// SecretKeyCreateRequest represents the request to create a secret key
type SecretKeyCreateRequest struct {
	Name            string              `json:"name" binding:"required,min=1,max=64" example:"Database Connection Key"`
	KeyType         model.SecretKeyType `json:"key_type" binding:"required,oneof=API_KEY DATABASE SSH CERTIFICATE CUSTOM" example:"DATABASE"`
	EncryptedValue  string              `json:"encrypted_value" binding:"required" example:"username:password or key content"`
	ResourceGroupID *uint               `json:"resource_group_id" example:"1"`
	Description     *string             `json:"description" example:"MySQL database connection credentials"`
	// @Schema(example="{\"rotation_interval\":90}")
	CustomFields model.CustomFields `json:"custom_fields"`
	// @Schema(example="[1,2,3]")
	AuthorizedUsers []uint     `json:"authorized_users"`
	ExpiresAt       *time.Time `json:"expires_at" example:"2025-12-31T23:59:59Z"`
}

// SecretKeyUpdateRequest represents the request to update a secret key
// According to API design, only two parameters are needed for update
type SecretKeyUpdateRequest struct {
	EncryptedValue string              `json:"encrypted_value" binding:"required" example:"new_username:new_password"`
	KeyType        model.SecretKeyType `json:"key_type" binding:"required,oneof=API_KEY DATABASE SSH CERTIFICATE CUSTOM" example:"DATABASE"`
}

// SecretKeyQueryRequest represents the request to query secret keys
// According to API design, only 3 parameters: page, page_size, key_type
type SecretKeyQueryRequest struct {
	common.PaginationRequest
	KeyType *model.SecretKeyType `form:"key_type" binding:"omitempty,oneof=API_KEY DATABASE SSH CERTIFICATE CUSTOM"`
}

// SecretKeyExportRequest represents the request to export secret keys
// According to API design, only 1 parameter: format
type SecretKeyExportRequest struct {
	Format string `form:"format" binding:"omitempty,oneof=csv excel json"`
}
