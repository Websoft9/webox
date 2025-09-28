package service

import (
	"api-service/internal/dto/request"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// PermissionServiceTestSuite
type PermissionServiceTestSuite struct {
	suite.Suite
	service            *permissionService
	mockPermissionRepo *MockPermissionRepository
	mockRoleRepo       *MockRoleRepository
	logger             logger.Logger
}

func (suite *PermissionServiceTestSuite) SetupTest() {
	suite.mockPermissionRepo = &MockPermissionRepository{}
	suite.mockRoleRepo = &MockRoleRepository{}
	suite.logger = logger.NewZapLogger(logger.InfoLevel, nil)

	suite.service = &permissionService{
		permissionRepo: suite.mockPermissionRepo,
		db:             &gorm.DB{},
		logger:         suite.logger,
	}
}

func (suite *PermissionServiceTestSuite) TestCreatePermission_Success() {
	ctx := context.Background()
	req := &request.CreatePermissionRequest{
		Scope:       "platform",
		Name:        "Test Permission",
		Code:        "test:permission",
		Module:      "test",
		Action:      "manage",
		Resource:    "test",
		Description: "Test permission description",
		IsMenu:      false,
		SortOrder:   10,
	}
	createdBy := uint(1)

	// Mock expectations
	suite.mockPermissionRepo.On("GetByCode", ctx, req.Code).Return(nil, errors.ErrRecordNotFound)
	suite.mockPermissionRepo.On("Create", ctx, mock.AnythingOfType("*model.Permission")).
		Run(func(args mock.Arguments) {
			perm := args.Get(1).(*model.Permission)
			perm.ID = 1 // Simulate database ID assignment
		}).Return(nil)

	// Execute
	result, err := suite.service.CreatePermission(ctx, req, createdBy)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(req.Name, result.Name)
	suite.Equal(req.Code, result.Code)
	suite.Equal(req.Module, result.Module)

	suite.mockPermissionRepo.AssertExpectations(suite.T())
}

func (suite *PermissionServiceTestSuite) TestCreatePermission_CodeAlreadyExists() {
	ctx := context.Background()
	req := &request.CreatePermissionRequest{
		Scope:  "platform",
		Name:   "Test Permission",
		Code:   "existing:permission",
		Module: "test",
		Action: "manage",
	}
	createdBy := uint(1)

	// Mock existing permission
	existingPerm := &model.Permission{
		BaseModel: model.BaseModel{ID: 1},
		Code:      req.Code,
	}
	suite.mockPermissionRepo.On("GetByCode", ctx, req.Code).Return(existingPerm, nil)

	// Execute
	result, err := suite.service.CreatePermission(ctx, req, createdBy)

	// Assert
	suite.Error(err)
	suite.Nil(result)
	suite.Contains(err.Error(), "permission code already exists")

	suite.mockPermissionRepo.AssertExpectations(suite.T())
}

func (suite *PermissionServiceTestSuite) TestGetPermissionTree_Success() {
	ctx := context.Background()
	req := &request.PermissionTreeRequest{
		Scope:  "platform",
		Status: nil,
	}

	// Mock permissions with parent-child relationship
	childPermission := &model.Permission{
		BaseModel:  model.BaseModel{ID: 2},
		ParentCode: "user:manage",
		Name:       "Create User",
		Code:       "user:create",
	}

	permissions := []*model.Permission{
		{
			BaseModel:  model.BaseModel{ID: 1},
			ParentCode: "",
			Name:       "User Management",
			Code:       "user:manage",
			Children:   []*model.Permission{childPermission},
		},
	}

	suite.mockPermissionRepo.On("GetTree", ctx, req).Return(permissions, nil)

	// Execute
	result, err := suite.service.GetPermissionTree(ctx, req)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(1, len(result))
	suite.Equal("User Management", result[0].Name)
	suite.Equal(1, len(result[0].Children))

	suite.mockPermissionRepo.AssertExpectations(suite.T())
}

func (suite *PermissionServiceTestSuite) TestListPermissions_Success() {
	ctx := context.Background()
	req := &request.ListPermissionsRequest{
		PaginationRequest: request.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
		Search: "user",
	}

	permissions := []*model.Permission{
		{
			BaseModel: model.BaseModel{ID: 1},
			Name:      "User Management",
			Code:      "user:manage",
		},
		{
			BaseModel: model.BaseModel{ID: 2},
			Name:      "User Create",
			Code:      "user:create",
		},
	}
	total := int64(2)

	suite.mockPermissionRepo.On("List", ctx, req, mock.AnythingOfType("string")).Return(permissions, total, nil)

	// Execute
	result, err := suite.service.ListPermissions(ctx, req)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	items, ok := result.Items.([]interface{})
	suite.True(ok, "Items should be a slice")
	suite.Equal(2, len(items))
	suite.Equal(total, result.Total)
	suite.Equal(req.GetPage(), result.Page)
	suite.Equal(req.GetPageSize(), result.PageSize)

	suite.mockPermissionRepo.AssertExpectations(suite.T())
}

func (suite *PermissionServiceTestSuite) TestUpdatePermission_Success() {
	ctx := context.Background()
	permissionID := uint(1)
	req := &request.UpdatePermissionRequest{
		Name:        "Updated Permission",
		Description: "Updated description",
		SortOrder:   20,
		Status:      1,
	}
	updatedBy := uint(1)

	// Mock existing permission
	existingPerm := &model.Permission{
		BaseModel: model.BaseModel{ID: permissionID},
		Name:      "Original Permission",
		Code:      "test:permission",
		IsSystem:  false,
	}
	suite.mockPermissionRepo.On("GetByID", ctx, permissionID).Return(existingPerm, nil)
	suite.mockPermissionRepo.On("Update", ctx, mock.AnythingOfType("*model.Permission")).Return(nil)

	// Mock updated permission for return
	updatedPerm := &model.Permission{
		BaseModel:   model.BaseModel{ID: permissionID},
		Name:        req.Name,
		Code:        existingPerm.Code,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
	}
	suite.mockPermissionRepo.On("GetByID", ctx, permissionID).Return(updatedPerm, nil)

	// Execute
	result, err := suite.service.UpdatePermission(ctx, permissionID, req, updatedBy)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(req.Name, result.Name)
	suite.Equal(req.Description, result.Description)

	suite.mockPermissionRepo.AssertExpectations(suite.T())
}

func (suite *PermissionServiceTestSuite) TestUpdatePermission_SystemPermission() {
	ctx := context.Background()
	permissionID := uint(1)
	req := &request.UpdatePermissionRequest{
		Name: "Updated Permission",
	}
	updatedBy := uint(1)

	// Mock system permission
	systemPerm := &model.Permission{
		BaseModel: model.BaseModel{ID: permissionID},
		Name:      "System Permission",
		Code:      "system:permission",
		IsSystem:  true,
	}
	suite.mockPermissionRepo.On("GetByID", ctx, permissionID).Return(systemPerm, nil)

	// Execute
	result, err := suite.service.UpdatePermission(ctx, permissionID, req, updatedBy)

	// Assert
	suite.Error(err)
	suite.Nil(result)
	suite.Contains(err.Error(), "cannot update system permission")

	suite.mockPermissionRepo.AssertExpectations(suite.T())
}

func (suite *PermissionServiceTestSuite) TestDeletePermission_Success() {
	ctx := context.Background()
	permissionID := uint(1)

	// Mock permission that is not system and has no associated roles
	permission := &model.Permission{
		BaseModel: model.BaseModel{ID: permissionID},
		Name:      "Test Permission",
		Code:      "test:permission",
		IsSystem:  false,
	}

	suite.mockPermissionRepo.On("GetByID", ctx, permissionID).Return(permission, nil)
	suite.mockPermissionRepo.On("CountRoles", ctx, permissionID).Return(int64(0), nil)
	suite.mockPermissionRepo.On("Delete", ctx, permissionID).Return(nil)

	// Execute
	err := suite.service.DeletePermission(ctx, permissionID)

	// Assert
	suite.NoError(err)

	suite.mockPermissionRepo.AssertExpectations(suite.T())
}

func (suite *PermissionServiceTestSuite) TestCheckUserPermission_Success() {
	ctx := context.Background()
	userID := uint(1)
	resource := "user"
	action := "create"

	suite.mockPermissionRepo.On("CheckUserPermission", ctx, userID, resource, action).Return(true, nil)

	// Execute
	hasPermission, err := suite.service.CheckUserPermission(ctx, userID, resource, action)

	// Assert
	suite.NoError(err)
	suite.True(hasPermission)

	suite.mockPermissionRepo.AssertExpectations(suite.T())
}

func (suite *PermissionServiceTestSuite) TestGetUserPermissions_Success() {
	ctx := context.Background()
	userID := uint(1)

	permissions := []*model.Permission{
		{
			BaseModel: model.BaseModel{ID: 1},
			Name:      "User Management",
			Code:      "user:manage",
		},
		{
			BaseModel: model.BaseModel{ID: 2},
			Name:      "User Create",
			Code:      "user:create",
		},
	}

	suite.mockPermissionRepo.On("GetUserPermissions", ctx, userID).Return(permissions, nil)

	// Execute
	result, err := suite.service.GetUserPermissions(ctx, userID)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(2, len(result))
	suite.Equal("User Management", result[0].Name)

	suite.mockPermissionRepo.AssertExpectations(suite.T())
}

func (suite *PermissionServiceTestSuite) TestBatchUpdatePermissionStatus_Success() {
	ctx := context.Background()
	ids := []uint{1, 2, 3}
	status := 0

	suite.mockPermissionRepo.On("BatchUpdateStatus", ctx, ids, status).Return(nil)

	// Execute
	err := suite.service.BatchUpdatePermissionStatus(ctx, ids, status)

	// Assert
	suite.NoError(err)

	suite.mockPermissionRepo.AssertExpectations(suite.T())
}

// Run the permission service test suite
func TestPermissionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionServiceTestSuite))
}

// Benchmark tests
func BenchmarkPermissionService_CreatePermission(b *testing.B) {
	mockPermissionRepo := &MockPermissionRepository{}
	logger := logger.NewZapLogger(logger.InfoLevel, nil)

	service := &permissionService{
		permissionRepo: mockPermissionRepo,
		db:             &gorm.DB{},
		logger:         logger,
	}

	ctx := context.Background()
	req := &request.CreatePermissionRequest{
		Scope:  "platform",
		Name:   "Benchmark Permission",
		Code:   "benchmark:permission",
		Module: "benchmark",
		Action: "manage",
	}
	createdBy := uint(1)

	// Setup mocks for benchmark
	mockPermissionRepo.On("GetByCode", ctx, mock.AnythingOfType("string")).Return(nil, errors.ErrRecordNotFound)
	mockPermissionRepo.On("Create", ctx, mock.AnythingOfType("*model.Permission")).
		Run(func(args mock.Arguments) {
			perm := args.Get(1).(*model.Permission)
			perm.ID = 1
		}).Return(nil)
	mockPermissionRepo.On("GetByID", ctx, mock.AnythingOfType("uint")).Return(&model.Permission{
		BaseModel: model.BaseModel{ID: 1},
		Name:      req.Name,
		Code:      req.Code,
	}, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.Code = fmt.Sprintf("benchmark:permission_%d", i)
		_, _ = service.CreatePermission(ctx, req, createdBy)
	}
}
