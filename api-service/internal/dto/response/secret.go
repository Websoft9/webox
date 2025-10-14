package response

import (
	"api-service/internal/dto/common"
	"api-service/internal/model"
	"time"
)

// SecretKeyResponse represents the response for secret key data
type SecretKeyResponse struct {
	ID          uint                `json:"id" example:"1"`
	Name        string              `json:"name" example:"Database Connection Key"`
	KeyType     model.SecretKeyType `json:"key_type" example:"DATABASE"`
	Usage       string              `json:"usage" example:"MYSQL_CONNECTION"`
	Description *string             `json:"description" example:"MySQL database connection credentials"`
	IsEncrypted bool                `json:"is_encrypted" example:"true"`
	// @Schema(example="{\"rotation_interval\":90}")
	CustomFields model.CustomFields `json:"custom_fields"`
	// @Schema(example="[1,2,3]")
	AuthorizedUsers []uint     `json:"authorized_users"`
	ExpiresAt       *time.Time `json:"expires_at" example:"2025-12-31T23:59:59Z"`
	OwnerID         uint       `json:"owner_id" example:"1"`
	CreatedAt       time.Time  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt       time.Time  `json:"updated_at" example:"2024-06-15T10:30:00Z"`
}

// SecretKeyValueResponse represents the response for secret key value
type SecretKeyValueResponse struct {
	Value     string     `json:"value" example:"sk-1234567890abcdef"`
	ExpiresAt *time.Time `json:"expires_at" example:"2025-12-31T23:59:59Z"`
}

// SecretKeyListResponse represents the paginated list response
// Following the pattern used in AlertRecordListResponse
type SecretKeyListResponse struct {
	Items    []SecretKeyResponse `json:"items"`
	Total    int64               `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
}

// ToSecretKeyResponse converts a SecretKey model to SecretKeyResponse
func ToSecretKeyResponse(secretKey *model.SecretKey) *SecretKeyResponse {
	if secretKey == nil {
		return nil
	}

	response := &SecretKeyResponse{
		ID:           secretKey.ID,
		Name:         secretKey.Name,
		KeyType:      secretKey.KeyType,
		Usage:        string(secretKey.KeyType) + "_CONNECTION", // Derived from KeyType
		IsEncrypted:  true,                                      // Always true since we encrypt all values
		CustomFields: secretKey.CustomFields,
		ExpiresAt:    secretKey.ExpiresAt,
		OwnerID:      secretKey.OwnerID,
		CreatedAt:    secretKey.CreatedAt,
		UpdatedAt:    secretKey.UpdatedAt,
	}

	if secretKey.Description != nil {
		response.Description = secretKey.Description
	}

	// AuthorizedUsers and AuthorizedGroups are simplified for response DTO.
	// In the actual implementation, these would come from authorization logic
	response.AuthorizedUsers = []uint{secretKey.OwnerID}

	return response
}

// ToSecretKeyListResponse converts a slice of SecretKey models to SecretKeyListResponse
// ToSecretKeyListResponse converts a slice of SecretKey models to PaginationResponse
func ToSecretKeyListResponse(secretKeys []*model.SecretKey, total int64, page, pageSize int) *common.PaginationResponse {
	items := make([]SecretKeyResponse, 0, len(secretKeys))
	for _, secretKey := range secretKeys {
		if response := ToSecretKeyResponse(secretKey); response != nil {
			items = append(items, *response)
		}
	}

	return common.NewPaginationResponse(
		page,
		pageSize,
		total,
		items,
	)
}

// ToSecretKeyValueResponse creates a value response
func ToSecretKeyValueResponse(value string, expiresAt *time.Time) *SecretKeyValueResponse {
	return &SecretKeyValueResponse{
		Value:     value,
		ExpiresAt: expiresAt,
	}
}
