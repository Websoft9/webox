package service

import (
	"context"
	"mime/multipart"

	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
)

// SecretKeyService defines the interface for secret key business logic
type SecretKeyService interface {
	// CreateSecretKey creates a new secret key with encryption
	CreateSecretKey(ctx context.Context, req *request.SecretKeyCreateRequest, userID uint) (*response.SecretKeyResponse, error)

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

	// UploadSecretFile uploads a secret key file
	UploadSecretFile(ctx context.Context, file *multipart.FileHeader, fileType string, userID uint) (*response.SecretFileUploadResponse, error)

	// DownloadSecretFile downloads a secret key file
	DownloadSecretFile(ctx context.Context, filename string, userID uint) (filePath string, originalName string, err error)

	// DeleteSecretFile deletes a secret key file
	DeleteSecretFile(ctx context.Context, filename string, userID uint) error
}
