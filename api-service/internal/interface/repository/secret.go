package repository

import (
	"context"

	"api-service/internal/dto/request"
	"api-service/internal/model"
)

// SecretKeyRepository defines the interface for secret key data access
type SecretKeyRepository interface {
	// Create creates a new secret key
	Create(ctx context.Context, secretKey *model.SecretKey) error

	// GetByID retrieves a secret key by ID
	GetByID(ctx context.Context, id uint) (*model.SecretKey, error)

	// Update updates an existing secret key
	Update(ctx context.Context, secretKey *model.SecretKey) error

	// Delete soft deletes a secret key by ID
	Delete(ctx context.Context, id uint) error

	// List retrieves secret keys with pagination and filtering
	List(ctx context.Context, req *request.SecretKeyQueryRequest, userID uint) ([]*model.SecretKey, int64, error)

	// ListAll retrieves all secret keys for a user without pagination (for export)
	ListAll(ctx context.Context, userID uint) ([]*model.SecretKey, error)

	// GetByOwnerID retrieves secret keys by owner ID
	GetByOwnerID(ctx context.Context, ownerID uint) ([]*model.SecretKey, error)

	// ExistsByName checks if a secret key with the given name exists for a user
	ExistsByName(ctx context.Context, name string, ownerID uint, excludeID ...uint) (bool, error)

	// CountByType counts secret keys by type for a user
	CountByType(ctx context.Context, keyType model.SecretKeyType, ownerID uint) (int64, error)

	// CreateUserSecret creates a user secret relationship
	CreateUserSecret(ctx context.Context, userSecret *model.UserSecret) error

	// DeleteUserSecretsBySecretKeyID deletes all user secret relationships for a secret key
	DeleteUserSecretsBySecretKeyID(ctx context.Context, secretKeyID uint) error

	// CheckUserSecretAccess checks if a user has access to a secret key
	CheckUserSecretAccess(ctx context.Context, userID, secretKeyID uint) (bool, error)

	// CreateSecretReference creates a secret reference record
	CreateSecretReference(ctx context.Context, reference *model.SecretReference) error

	// DeleteSecretReferencesBySecretID deletes all references for a secret key
	DeleteSecretReferencesBySecretID(ctx context.Context, secretID uint) error
}
