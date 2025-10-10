package repository

import (
	"context"
	"fmt"
	"time"

	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"

	"gorm.io/gorm"
)

// serverAgentRepository implements repository.ServerAgentRepository
type serverAgentRepository struct {
	db *gorm.DB
}

// NewServerAgentRepository creates a new server agent repository instance
func NewServerAgentRepository(db *gorm.DB) repository.ServerAgentRepository {
	return &serverAgentRepository{
		db: db,
	}
}

// CreateAgent creates a new server agent record
func (r *serverAgentRepository) CreateAgent(ctx context.Context, agent *model.ServerAgent) error {
	if err := r.db.WithContext(ctx).Create(agent).Error; err != nil {
		return fmt.Errorf("failed to create server agent: %w", err)
	}
	return nil
}

// GetAgentByID retrieves a server agent by its ID
func (r *serverAgentRepository) GetAgentByID(ctx context.Context, id uint) (*model.ServerAgent, error) {
	var agent model.ServerAgent
	err := r.db.WithContext(ctx).
		Preload("Server").
		Where("id = ?", id).
		First(&agent).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrServerAgentNotFound
		}
		return nil, fmt.Errorf("failed to get server agent by ID: %w", err)
	}
	return &agent, nil
}

// GetAgentByServerID retrieves the agent for a specific server
func (r *serverAgentRepository) GetAgentByServerID(ctx context.Context, serverID uint) (*model.ServerAgent, error) {
	var agent model.ServerAgent
	err := r.db.WithContext(ctx).
		Preload("Server").
		Where("server_id = ?", serverID).
		First(&agent).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrServerAgentNotFound
		}
		return nil, fmt.Errorf("failed to get server agent by server ID: %w", err)
	}
	return &agent, nil
}

// UpdateAgent updates an existing server agent record
func (r *serverAgentRepository) UpdateAgent(ctx context.Context, agent *model.ServerAgent) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", agent.ID).
		Updates(agent)

	if result.Error != nil {
		return fmt.Errorf("failed to update server agent: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.ErrServerAgentNotFound
	}

	return nil
}

// DeleteAgent physically deletes a server agent record
func (r *serverAgentRepository) DeleteAgent(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Unscoped().
		Where("id = ?", id).
		Delete(&model.ServerAgent{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete server agent: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.ErrServerAgentNotFound
	}

	return nil
}

// UpdateAgentLastSeen updates a server agent's last seen timestamp
func (r *serverAgentRepository) UpdateAgentLastSeen(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Model(&model.ServerAgent{}).
		Where("id = ?", id).
		Update("last_seen_at", time.Now())

	if result.Error != nil {
		return fmt.Errorf("failed to update server agent last seen: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.ErrServerAgentNotFound
	}

	return nil
}
