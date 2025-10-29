package repository

import (
	"context"

	"api-service/internal/dto/request"
	"api-service/internal/model"
)

// SecretRepository defines the interface for secret data access
type SecretRepository interface {
	// Create creates a new secret
	Create(ctx context.Context, secret *model.Secret) error

	// GetByID retrieves a secret by ID
	GetByID(ctx context.Context, id uint) (*model.Secret, error)

	// GetByCode retrieves a secret by code
	GetByCode(ctx context.Context, code string) (*model.Secret, error)

	// List retrieves a paginated list of secrets with filters
	List(ctx context.Context, req *request.ListSecretsRequest, userID uint) ([]*model.Secret, int64, error)

	// Update updates an existing secret
	Update(ctx context.Context, secret *model.Secret) error

	// Delete deletes a secret by ID
	Delete(ctx context.Context, id uint) error

	// ExistsByName checks if a secret name exists in a resource group
	ExistsByName(ctx context.Context, resourceGroupID uint, name string, excludeID *uint) (bool, error)

	// ExistsByCode checks if a secret code exists
	ExistsByCode(ctx context.Context, code string) (bool, error)

	// GetReferenceCount gets the count of references for a secret
	GetReferenceCount(ctx context.Context, secretID uint) (int, error)
}

// SecretReferenceRepository defines the interface for secret reference data access
type SecretReferenceRepository interface {
	// Create creates a new secret reference
	Create(ctx context.Context, reference *model.SecretReference) error

	// ListBySecretID retrieves all references for a secret
	ListBySecretID(ctx context.Context, secretID uint) ([]*model.SecretReference, error)

	// ListByResourceCode retrieves all references for a resource
	ListByResourceCode(ctx context.Context, resourceCode string) ([]*model.SecretReference, error)

	// Exists checks if a reference exists
	Exists(ctx context.Context, secretID uint, resourceCode string) (bool, error)

	// Delete deletes a specific reference
	Delete(ctx context.Context, id uint) error

	// DeleteBySecretID deletes all references for a secret
	DeleteBySecretID(ctx context.Context, secretID uint) error

	// HasActiveReferences checks if a secret has any active references
	HasActiveReferences(ctx context.Context, secretID uint) (bool, error)
}

// SecretAuthorizeRepository defines the interface for secret authorization data access
type SecretAuthorizeRepository interface {
	// Create creates a new secret authorization
	Create(ctx context.Context, authorize *model.SecretAuthorize) error

	// BatchCreate creates multiple secret authorizations
	BatchCreate(ctx context.Context, authorizes []*model.SecretAuthorize) error

	// ListBySecretID retrieves all authorizations for a secret
	ListBySecretID(ctx context.Context, secretID uint) ([]*model.SecretAuthorize, error)

	// Delete deletes a specific authorization
	Delete(ctx context.Context, id uint) error

	// DeleteBySecretID deletes all authorizations for a secret
	DeleteBySecretID(ctx context.Context, secretID uint) error

	// IsAuthorized checks if a user is authorized to access a secret
	IsAuthorized(ctx context.Context, secretID uint, userID uint) (bool, error)
}
