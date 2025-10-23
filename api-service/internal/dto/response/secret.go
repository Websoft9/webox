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
	SecretFields model.SecretFields `json:"secret_fields"`
	// @Schema(example="[1,2,3]")
	AuthorizedUsers []uint     `json:"authorized_users"`
	ExpiresAt       *time.Time `json:"expires_at" example:"2025-12-31T23:59:59Z"`
	OwnerID         uint       `json:"owner_id" example:"1"`
	CreatedAt       time.Time  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt       time.Time  `json:"updated_at" example:"2024-06-15T10:30:00Z"`
}

// SecretKeyValueResponse represents the response containing secret key value
type SecretKeyValueResponse struct {
	KeyType      model.SecretKeyType    `json:"key_type"`
	SecretFields map[string]interface{} `json:"secret_fields"`
	ExpiresAt    *time.Time             `json:"expires_at,omitempty"`
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
		SecretFields: secretKey.SecretFields,
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

// ToSecretKeyValueResponse converts decrypted custom fields to response
func ToSecretKeyValueResponse(keyType model.SecretKeyType, customFields map[string]interface{}, expiresAt *time.Time) *SecretKeyValueResponse {
	return &SecretKeyValueResponse{
		KeyType:      keyType,
		SecretFields: customFields,
		ExpiresAt:    expiresAt,
	}
}

// SecretFileUploadResponse represents the response for secret file upload
type SecretFileUploadResponse struct {
	Filename     string `json:"filename" example:"2345678ioasjhhdvgajdjknasd.key"`
	OriginalName string `json:"original_name" example:"my-certificate.key"`
	FilePath     string `json:"file_path" example:"/home/appuser/data/2345678ioasjhhdvgajdjknasd.key"`
}

// SecretResourceInfo represents basic secret information in resource query response
type SecretResourceInfo struct {
	ID          uint                `json:"id" example:"1"`
	Name        string              `json:"name" example:"数据库连接密钥"`
	KeyType     model.SecretKeyType `json:"key_type" example:"ACCOUNT"`
	Description *string             `json:"description" example:"MySQL数据库连接密钥"`
	IsEncrypted bool                `json:"is_encrypted" example:"true"`
	CreatedAt   time.Time           `json:"created_at" example:"2024-10-01T00:00:00Z"`
}

// AssociatedResource represents a resource associated with a secret
type AssociatedResource struct {
	ID           uint      `json:"id" example:"1"`
	ResourceCode string    `json:"resource_code" example:"mysql-prod-001"`
	CreatedAt    time.Time `json:"created_at" example:"2025-01-22T10:30:00Z"`
}

// SecretResourceResponse represents the response for querying secret's associated resources
type SecretResourceResponse struct {
	SecretInfo          SecretResourceInfo   `json:"secret_info"`
	AssociatedResources []AssociatedResource `json:"associated_resources"`
}

// ToSecretResourceResponse converts secret and references to SecretResourceResponse
func ToSecretResourceResponse(secret *model.SecretKey, references []*model.SecretReference) *SecretResourceResponse {
	if secret == nil {
		return nil
	}

	secretInfo := SecretResourceInfo{
		ID:          secret.ID,
		Name:        secret.Name,
		KeyType:     secret.KeyType,
		Description: secret.Description,
		IsEncrypted: true, // Always true since we encrypt all values
		CreatedAt:   secret.CreatedAt,
	}

	associatedResources := make([]AssociatedResource, 0, len(references))
	for _, ref := range references {
		if ref != nil {
			associatedResources = append(associatedResources, AssociatedResource{
				ID:           ref.ID,
				ResourceCode: ref.ResourceCode,
				CreatedAt:    ref.CreatedAt,
			})
		}
	}

	return &SecretResourceResponse{
		SecretInfo:          secretInfo,
		AssociatedResources: associatedResources,
	}
}
