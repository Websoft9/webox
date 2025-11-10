package service

import (
	"api-service/internal/model"
	"context"
)

// ConfigService defines the interface for configuration center business logic
type ConfigService interface {
	// System Configuration Operations
	GetSystemConfig(ctx context.Context, configKey string) (*model.SystemConfig, error)
	CreateSystemConfig(ctx context.Context, config *model.SystemConfig) error
	UpdateSystemConfig(ctx context.Context, configKey string, config *model.SystemConfig) error
	DeleteSystemConfig(ctx context.Context, configKey string) error

	// Service Configuration Operations
	GetServiceConfig(ctx context.Context, code string) (*model.ServiceConfig, error)
	CreateServiceConfig(ctx context.Context, config *model.ServiceConfig) error
	UpdateServiceConfig(ctx context.Context, code string, config *model.ServiceConfig) error
	DeleteServiceConfig(ctx context.Context, code string) error

	// User Profile Configuration Operations
	GetUserProfile(ctx context.Context, userID uint) (map[string]interface{}, error)
	UpdateUserProfile(ctx context.Context, userID uint, configKey string, configValue interface{}) error
	SyncUserProfileOnLogin(ctx context.Context, userID uint) error

	// User Permission Configuration Operations
	GetUserPermissions(ctx context.Context, userID uint) (map[string]string, error)
	SyncUserPermissionsOnLogin(ctx context.Context, userID uint) error
	SyncUserPermissionsOnChange(ctx context.Context, userID uint) error

	// Configuration Preload Operations
	PreloadConfigs(ctx context.Context) error
}
