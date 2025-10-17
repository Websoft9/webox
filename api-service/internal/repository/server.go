package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"

	"gorm.io/gorm"
)

// serverRepository implements repository.ServerRepository
type serverRepository struct {
	db *gorm.DB
}

// NewServerRepository creates a new server repository instance
func NewServerRepository(db *gorm.DB) repository.ServerRepository {
	return &serverRepository{
		db: db,
	}
}

// CreateServer creates a new server record
func (r *serverRepository) CreateServer(ctx context.Context, server *model.Server) error {
	if err := r.db.WithContext(ctx).Create(server).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return errors.ErrServerNameExists
		}
		return fmt.Errorf("failed to create server: %w", err)
	}
	return nil
}

// GetServerByID retrieves a server by its ID
func (r *serverRepository) GetServerByID(ctx context.Context, id uint) (*model.Server, error) {
	var server model.Server
	err := r.db.WithContext(ctx).
		Preload("Agents").
		Where("id = ?", id).
		First(&server).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrServerNotFound
		}
		return nil, fmt.Errorf("failed to get server by ID: %w", err)
	}
	return &server, nil
}

// GetServerByName retrieves a server by its name
func (r *serverRepository) GetServerByName(ctx context.Context, name string) (*model.Server, error) {
	var server model.Server
	err := r.db.WithContext(ctx).
		Preload("Agents").
		Where("name = ?", name).
		First(&server).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrServerNotFound
		}
		return nil, fmt.Errorf("failed to get server by name: %w", err)
	}
	return &server, nil
}

// UpdateServer updates an existing server record
func (r *serverRepository) UpdateServer(ctx context.Context, server *model.Server) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", server.ID).
		Updates(server)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate") || strings.Contains(result.Error.Error(), "UNIQUE constraint") {
			return errors.ErrServerNameExists
		}
		return fmt.Errorf("failed to update server: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.ErrServerNotFound
	}

	return nil
}

// DeleteServer soft deletes a server record
func (r *serverRepository) DeleteServer(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Server{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete server: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.ErrServerNotFound
	}

	return nil
}

// ListServers retrieves servers with pagination and filtering
func (r *serverRepository) ListServers(ctx context.Context, req *request.ListServersRequest) ([]*model.Server, int64, error) {
	var servers []*model.Server
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Server{})

	// Apply filters
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.OS != "" {
		query = query.Where("os LIKE ?", "%"+req.OS+"%")
	}

	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		query = query.Where("name LIKE ? OR host LIKE ? OR description LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count servers: %w", err)
	}

	// Apply pagination and sorting
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	if req.SortBy != "" {
		order := "ASC"
		if req.SortOrder == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", req.SortBy, order))
	} else {
		query = query.Order("created_at DESC")
	}

	// Execute query with preloading
	err := query.Preload("Agents").Find(&servers).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list servers: %w", err)
	}

	return servers, total, nil
}

// UpdateServerStatus updates a server's status
func (r *serverRepository) UpdateServerStatus(ctx context.Context, id uint, status string) error {
	result := r.db.WithContext(ctx).
		Model(&model.Server{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update server status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.ErrServerNotFound
	}

	return nil
}

// UpdateServerLastSeen updates a server's last seen timestamp
func (r *serverRepository) UpdateServerLastSeen(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Model(&model.Server{}).
		Where("id = ?", id).
		Update("last_seen_at", time.Now())

	if result.Error != nil {
		return fmt.Errorf("failed to update server last seen: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.ErrServerNotFound
	}

	return nil
}

// GetServersByStatus retrieves servers by status
func (r *serverRepository) GetServersByStatus(ctx context.Context, status string) ([]*model.Server, error) {
	var servers []*model.Server
	err := r.db.WithContext(ctx).
		Preload("Agents").
		Where("status = ?", status).
		Find(&servers).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get servers by status: %w", err)
	}

	return servers, nil
}

// GetServersByIDs retrieves servers by a list of IDs
func (r *serverRepository) GetServersByIDs(ctx context.Context, ids []uint) ([]*model.Server, error) {
	var servers []*model.Server
	err := r.db.WithContext(ctx).
		Preload("Agents").
		Where("id IN ?", ids).
		Find(&servers).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get servers by IDs: %w", err)
	}

	return servers, nil
}

// BatchUpdateServerStatus updates status for multiple servers
func (r *serverRepository) BatchUpdateServerStatus(ctx context.Context, ids []uint, status string) error {
	result := r.db.WithContext(ctx).
		Model(&model.Server{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to batch update server status: %w", result.Error)
	}

	return nil
}

// CountServersByStatus counts servers by status
func (r *serverRepository) CountServersByStatus(ctx context.Context) (map[string]int64, error) {
	type StatusCount struct {
		Status string
		Count  int64
	}

	var results []StatusCount
	err := r.db.WithContext(ctx).
		Model(&model.Server{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to count servers by status: %w", err)
	}

	counts := make(map[string]int64)
	for _, result := range results {
		counts[result.Status] = result.Count
	}

	return counts, nil
}

// ExistsServerByName checks if a server with the given name exists
func (r *serverRepository) ExistsServerByName(ctx context.Context, name string, excludeID ...uint) (bool, error) {
	query := r.db.WithContext(ctx).
		Model(&model.Server{}).
		Where("name = ?", name)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check server name existence: %w", err)
	}

	return count > 0, nil
}

// GetServerAgentsByServerID retrieves all agents for a specific server
func (r *serverRepository) GetServerAgentsByServerID(ctx context.Context, serverID uint) ([]*model.ServerAgent, error) {
	var agents []*model.ServerAgent
	err := r.db.WithContext(ctx).
		Where("server_id = ?", serverID).
		Find(&agents).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get server agents: %w", err)
	}

	return agents, nil
}
