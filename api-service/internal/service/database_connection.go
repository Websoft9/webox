package service

import (
	"context"
	"encoding/json"

	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	interfaceRepo "api-service/internal/interface/repository"
	interfaceService "api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
)

// databaseConnectionService implements DatabaseConnectionService
type databaseConnectionService struct {
	repo   interfaceRepo.DatabaseConnectionRepository
	logger logger.Logger
}

// NewDatabaseConnectionService creates a new database connection service
func NewDatabaseConnectionService(
	repo interfaceRepo.DatabaseConnectionRepository,
	logger logger.Logger,
) (interfaceService.DatabaseConnectionService, error) {
	return &databaseConnectionService{
		repo:   repo,
		logger: logger,
	}, nil
}

// CreateConnection creates a new database connection (without credentials)
func (s *databaseConnectionService) CreateConnection(
	ctx context.Context,
	req *request.CreateDatabaseConnectionRequest,
	ownerID uint,
) (*response.DatabaseConnectionResponse, error) {
	// Validate config if provided
	var config *model.JSON
	if len(req.Config) > 0 {
		// Validate JSON format
		if err := validateConfig(req.Config); err != nil {
			s.logger.ErrorContext(ctx, "Invalid config format",
				logger.String("name", req.Name),
				logger.ErrorField(err))
			return nil, errors.NewAppError(errors.CodeInvalidParameterFormat)
		}
		cfg := model.JSON(req.Config)
		config = &cfg
	}

	// Create connection model
	conn := &model.DatabaseConnection{
		Name:            req.Name,
		DBType:          req.DBType,
		Host:            req.Host,
		Port:            req.Port,
		Database:        req.Database,
		Description:     req.Description,
		Config:          config,
		OwnerID:         ownerID,
		ResourceGroupID: req.ResourceGroupID,
	}

	// Save to database (code will be generated in repository)
	if err := s.repo.Create(ctx, conn); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create database connection",
			logger.String("name", req.Name),
			logger.Uint("owner_id", ownerID),
			logger.ErrorField(err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "Database connection created successfully",
		logger.Uint("connection_id", conn.ID),
		logger.String("code", conn.Code),
		logger.String("name", conn.Name),
		logger.Uint("owner_id", ownerID))

	return s.toResponse(conn), nil
}

// GetConnection retrieves a database connection by ID
// Credentials are managed separately via secret management system
func (s *databaseConnectionService) GetConnection(ctx context.Context, id, ownerID uint) (*response.DatabaseConnectionDetailResponse, error) {
	conn, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if conn.OwnerID != ownerID {
		return nil, errors.NewAppError(errors.CodeAccessDenied)
	}

	return s.toDetailResponse(conn), nil
}

// GetConnectionList retrieves a paginated list of database connections
func (s *databaseConnectionService) GetConnectionList(ctx context.Context, req *request.GetDatabaseConnectionListRequest, ownerID uint) (*common.PaginationResponse, error) {
	connections, total, err := s.repo.GetList(ctx, req, ownerID)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, len(connections))
	for i, conn := range connections {
		items[i] = s.toResponse(conn)
	}

	return common.NewPaginationResponse(req.Page, req.PageSize, total, items), nil
}

// UpdateConnection updates an existing database connection (without credentials)
func (s *databaseConnectionService) UpdateConnection(
	ctx context.Context,
	id uint,
	req *request.UpdateDatabaseConnectionRequest,
	ownerID uint,
) (*response.DatabaseConnectionResponse, error) {
	conn, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if conn.OwnerID != ownerID {
		return nil, errors.NewAppError(errors.CodeAccessDenied)
	}

	// Update fields
	if req.Name != nil {
		conn.Name = *req.Name
	}
	if req.Host != nil {
		conn.Host = *req.Host
	}
	if req.Port != nil {
		conn.Port = *req.Port
	}
	if req.Database != nil {
		conn.Database = req.Database
	}
	if req.Description != nil {
		conn.Description = req.Description
	}
	if req.Config != nil {
		// Validate config format
		if validateErr := validateConfig(req.Config); validateErr != nil {
			s.logger.ErrorContext(ctx, "Invalid config format",
				logger.Uint("connection_id", id),
				logger.ErrorField(validateErr))
			return nil, errors.NewAppError(errors.CodeInvalidParameterFormat)
		}
		cfg := model.JSON(req.Config)
		conn.Config = &cfg
	}
	if req.ResourceGroupID != nil {
		conn.ResourceGroupID = req.ResourceGroupID
	}

	// Save changes
	if updateErr := s.repo.Update(ctx, conn); updateErr != nil {
		s.logger.ErrorContext(ctx, "Failed to update database connection",
			logger.Uint("connection_id", id),
			logger.ErrorField(updateErr))
		return nil, updateErr
	}

	// Reload from database to get the correct updated_at timestamp set by database trigger
	updatedConn, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to reload updated connection",
			logger.Uint("connection_id", id),
			logger.ErrorField(err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "Database connection updated successfully",
		logger.Uint("connection_id", id),
		logger.Uint("owner_id", ownerID))

	return s.toResponse(updatedConn), nil
}

// DeleteConnection deletes a database connection
func (s *databaseConnectionService) DeleteConnection(ctx context.Context, id, ownerID uint) error {
	conn, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check ownership
	if conn.OwnerID != ownerID {
		return errors.NewAppError(errors.CodeAccessDenied)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete database connection",
			logger.Uint("connection_id", id),
			logger.ErrorField(err))
		return err
	}

	s.logger.InfoContext(ctx, "Database connection deleted successfully",
		logger.Uint("connection_id", id),
		logger.Uint("owner_id", ownerID))

	return nil
}

// validateConfig validates the config JSON format
func validateConfig(config map[string]interface{}) error {
	// Validate JSON can be marshaled
	_, err := json.Marshal(config)
	return err
}

// toResponse converts model to response
func (s *databaseConnectionService) toResponse(conn *model.DatabaseConnection) *response.DatabaseConnectionResponse {
	var config map[string]interface{}
	if conn.Config != nil {
		config = map[string]interface{}(*conn.Config)
	}

	return &response.DatabaseConnectionResponse{
		ID:              conn.ID,
		Name:            conn.Name,
		Code:            conn.Code,
		DBType:          conn.DBType,
		Host:            conn.Host,
		Port:            conn.Port,
		Database:        conn.Database,
		Description:     conn.Description,
		Config:          config,
		OwnerID:         conn.OwnerID,
		ResourceGroupID: conn.ResourceGroupID,
		CreatedAt:       conn.CreatedAt,
		UpdatedAt:       conn.UpdatedAt,
	}
}

// toDetailResponse converts model to detail response
func (s *databaseConnectionService) toDetailResponse(conn *model.DatabaseConnection) *response.DatabaseConnectionDetailResponse {
	return &response.DatabaseConnectionDetailResponse{
		DatabaseConnectionResponse: *s.toResponse(conn),
	}
}
