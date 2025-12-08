package service

import (
	"context"

	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
)

// CredentialService defines the interface for credential business logic
type CredentialService interface {
	// CreateCredential creates a new credential
	CreateCredential(ctx context.Context, req *request.CreateCredentialRequest, ownerID uint) (*response.CredentialResponse, error)

	// UpdateCredential updates an existing credential
	UpdateCredential(ctx context.Context, id uint, req *request.UpdateCredentialRequest, userID uint) (*response.CredentialResponse, error)

	// DeleteCredential deletes a credential
	DeleteCredential(ctx context.Context, id uint, userID uint) error

	// GetCredential retrieves a credential by ID
	GetCredential(ctx context.Context, id uint) (*response.CredentialDetailResponse, error)

	// ListCredentials retrieves a paginated list of credentials
	ListCredentials(ctx context.Context, req *request.ListCredentialsRequest) (*common.PaginationResponse, error)

	// ListCredentialTemplates retrieves credential templates
	ListCredentialTemplates(ctx context.Context, req *request.ListCredentialTemplatesRequest) ([]response.CredentialTemplateResponse, error)

	// ListCredentialCategories retrieves credential categories
	ListCredentialCategories(ctx context.Context) ([]response.CredentialCategoryResponse, error)

	// ResolveCredentialReference resolves credential reference expression
	// Input: {{ credentials.my_db.username }}
	// Output: decrypted actual value
	ResolveCredentialReference(ctx context.Context, expression string) (string, error)
}
