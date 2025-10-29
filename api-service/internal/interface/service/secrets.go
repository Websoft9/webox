package service

import (
	"context"

	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
)

// SecretService defines the interface for secret business logic
type SecretService interface {
	// CreateTextSecret creates a text type secret
	CreateTextSecret(ctx context.Context, req *request.CreateTextSecretRequest, ownerID uint) (*response.SecretResponse, error)

	// CreateAccountSecret creates an account type secret
	CreateAccountSecret(ctx context.Context, req *request.CreateAccountSecretRequest, ownerID uint) (*response.SecretResponse, error)

	// CreateFileSecret creates a file type secret
	CreateFileSecret(ctx context.Context, req *request.CreateFileSecretRequest, ownerID uint) (*response.SecretResponse, error)

	// ListSecrets retrieves a paginated list of secrets
	ListSecrets(ctx context.Context, req *request.ListSecretsRequest, userID uint) (*common.PaginationResponse, error)

	// GetSecret retrieves a secret by ID with full details
	GetSecret(ctx context.Context, id uint, userID uint) (*response.SecretDetailResponse, error)

	// UpdateSecret updates an existing secret
	UpdateSecret(ctx context.Context, id uint, req *request.UpdateSecretRequest, userID uint) (*response.SecretResponse, error)

	// DeleteSecret deletes a secret
	DeleteSecret(ctx context.Context, id uint, userID uint) error

	// CreateReference creates a secret reference
	CreateReference(ctx context.Context, req *request.CreateReferenceRequest) (*response.ReferenceResponse, error)
}
