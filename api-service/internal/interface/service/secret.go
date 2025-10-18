package service

import (
	"context"

	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
)

// SecretKeyService defines the interface for secret key business logic
type SecretKeyService interface {
	// CreateSecretKeyText creates a new text-based secret key with encryption
	CreateSecretKeyText(ctx context.Context, req *request.SecretKeyCreateTextRequest, userID uint) (*response.SecretKeyResponse, error)

	// CreateSecretKeyFile creates a new file-based secret key with encryption
	CreateSecretKeyFile(ctx context.Context, req *request.SecretKeyCreateFileRequest, userID uint) (*response.SecretKeyResponse, error)

	// GetSecretKey retrieves a secret key by ID (without sensitive data)
	GetSecretKey(ctx context.Context, id, userID uint) (*response.SecretKeyResponse, error)

	// GetSecretKeyValue retrieves the decrypted value of a secret key
	GetSecretKeyValue(ctx context.Context, id, userID uint) (*response.SecretKeyValueResponse, error)

	// UpdateSecretKey updates an existing secret key
	UpdateSecretKey(ctx context.Context, id, userID uint, req *request.SecretKeyUpdateRequest) (*response.SecretKeyResponse, error)

	// DeleteSecretKey deletes a secret key
	DeleteSecretKey(ctx context.Context, id, userID uint) error

	// ListSecretKeys retrieves secret keys with pagination and filtering
	ListSecretKeys(ctx context.Context, req *request.SecretKeyQueryRequest, userID uint) (*common.PaginationResponse, error)

	// ExportSecretKeys exports secret keys in specified format
	ExportSecretKeys(ctx context.Context, req *request.SecretKeyExportRequest, userID uint) ([]byte, string, error)

	// ValidateSecretKeyOwnership checks if user owns the secret key
	ValidateSecretKeyOwnership(ctx context.Context, secretKeyID, userID uint) error
}
