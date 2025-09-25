package service

import (
	"api-service/internal/dto/request"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"context"
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

// Mock repositories
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) Create(ctx context.Context, role *model.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, role *model.Role) error {
	args := m.Called(ctx, tx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) GetByID(ctx context.Context, id uint) (*model.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Role), args.Error(1)
}

func (m *MockRoleRepository) GetByCode(ctx context.Context, code string) (*model.Role, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Role), args.Error(1)
}

func (m *MockRoleRepository) Update(ctx context.Context, role *model.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, role *model.Role) error {
	args := m.Called(ctx, tx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoleRepository) List(ctx context.Context, req *request.ListRolesRequest) ([]*model.Role, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*model.Role), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoleRepository) GetWithPermissions(ctx context.Context, id uint) (*model.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Role), args.Error(1)
}

func (m *MockRoleRepository) GetUsers(ctx context.Context, id uint, page, pageSize int) ([]model.User, int64, error) {
	args := m.Called(ctx, id, page, pageSize)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoleRepository) GetWithUsers(ctx context.Context, id uint) (*model.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Role), args.Error(1)
}

func (m *MockRoleRepository) GetPermissions(ctx context.Context, roleID uint) ([]model.Permission, error) {
	args := m.Called(ctx, roleID)
	return args.Get(0).([]model.Permission), args.Error(1)
}

func (m *MockRoleRepository) AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint, grantedBy uint) error {
	args := m.Called(ctx, roleID, permissionIDs, grantedBy)
	return args.Error(0)
}

func (m *MockRoleRepository) AssignPermissionsWithTx(ctx context.Context, tx *gorm.DB, roleID uint, permissionIDs []uint, grantedBy uint) error {
	args := m.Called(ctx, tx, roleID, permissionIDs, grantedBy)
	return args.Error(0)
}

func (m *MockRoleRepository) RemovePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	args := m.Called(ctx, roleID, permissionIDs)
	return args.Error(0)
}

func (m *MockRoleRepository) BatchUpdateStatus(ctx context.Context, ids []uint, status int) error {
	args := m.Called(ctx, ids, status)
	return args.Error(0)
}

func (m *MockRoleRepository) CountPermissions(ctx context.Context, roleID uint) (int64, error) {
	args := m.Called(ctx, roleID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRoleRepository) CountUsers(ctx context.Context, roleID uint) (int64, error) {
	args := m.Called(ctx, roleID)
	return args.Get(0).(int64), args.Error(1)
}

type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) Create(ctx context.Context, permission *model.Permission) error {
	args := m.Called(ctx, permission)
	return args.Error(0)
}

func (m *MockPermissionRepository) GetByID(ctx context.Context, id uint) (*model.Permission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetByCode(ctx context.Context, code string) (*model.Permission, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) Update(ctx context.Context, permission *model.Permission) error {
	args := m.Called(ctx, permission)
	return args.Error(0)
}

func (m *MockPermissionRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPermissionRepository) List(ctx context.Context, req *request.ListPermissionsRequest) ([]*model.Permission, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*model.Permission), args.Get(1).(int64), args.Error(2)
}

func (m *MockPermissionRepository) GetTree(ctx context.Context, req *request.PermissionTreeRequest) ([]*model.Permission, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetWithRoles(ctx context.Context, id uint) (*model.Permission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetChildren(ctx context.Context, parentID uint) ([]*model.Permission, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).([]*model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetRoles(ctx context.Context, permissionID uint, page, pageSize int) ([]model.Role, int64, error) {
	args := m.Called(ctx, permissionID, page, pageSize)
	return args.Get(0).([]model.Role), args.Get(1).(int64), args.Error(2)
}

func (m *MockPermissionRepository) CountRoles(ctx context.Context, permissionID uint) (int64, error) {
	args := m.Called(ctx, permissionID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPermissionRepository) BatchUpdateStatus(ctx context.Context, ids []uint, status int) error {
	args := m.Called(ctx, ids, status)
	return args.Error(0)
}

func (m *MockPermissionRepository) GetByIDs(ctx context.Context, ids []uint) ([]*model.Permission, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetUserPermissions(ctx context.Context, userID uint) ([]*model.Permission, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.Permission), args.Error(1)
}

func (m *MockPermissionRepository) CheckUserPermission(ctx context.Context, userID uint, resource, action string) (bool, error) {
	args := m.Called(ctx, userID, resource, action)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockPermissionRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, permission *model.Permission) error {
	args := m.Called(ctx, tx, permission)
	return args.Error(0)
}

func (m *MockPermissionRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, permission *model.Permission) error {
	args := m.Called(ctx, tx, permission)
	return args.Error(0)
}

// Mock database
type MockDB struct {
	mock.Mock
}

func (m *MockDB) WithContext(ctx context.Context) *gorm.DB {
	return &gorm.DB{}
}

func (m *MockDB) Transaction(fn func(tx *gorm.DB) error) error {
	args := m.Called(fn)
	return args.Error(0)
}

// Test setup
func setupRoleServiceTest() (*roleService, *MockRoleRepository, *MockPermissionRepository, *MockDB) {
	mockRoleRepo := &MockRoleRepository{}
	mockPermissionRepo := &MockPermissionRepository{}
	mockDB := &MockDB{}
	mockLogger := logger.NewZapLogger(logger.InfoLevel, nil)

	// Create a test database connection
	testDB := setupTestDB()

	service := &roleService{
		roleRepo:       mockRoleRepo,
		permissionRepo: mockPermissionRepo,
		db:             testDB,
		logger:         mockLogger,
	}

	return service, mockRoleRepo, mockPermissionRepo, mockDB
}

// setupTestDB creates a test database connection
func setupTestDB() *gorm.DB {
	// Use SQLite in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: gorm_logger.Default.LogMode(gorm_logger.Silent),
	})
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}
	return db
}

func TestRoleService_CreateRole_Success(t *testing.T) {
	service, mockRoleRepo, mockPermissionRepo, _ := setupRoleServiceTest()
	ctx := context.Background()

	req := &request.CreateRoleRequest{
		Name:          "Test Role",
		Code:          "test_role",
		Description:   "Test role description",
		PermissionIDs: []uint{1, 2, 3},
		SortOrder:     1,
	}
	createdBy := uint(1)

	// Mock expectations
	mockRoleRepo.On("GetByCode", ctx, req.Code).Return(nil, errors.ErrRecordNotFound)
	mockPermissionRepo.On("GetByIDs", ctx, req.PermissionIDs).Return([]*model.Permission{
		{BaseModel: model.BaseModel{ID: 1}},
		{BaseModel: model.BaseModel{ID: 2}},
		{BaseModel: model.BaseModel{ID: 3}},
	}, nil)

	// Mock transaction
	mockRoleRepo.On("CreateWithTx", ctx, mock.AnythingOfType("*gorm.DB"), mock.AnythingOfType("*model.Role")).
		Run(func(args mock.Arguments) {
			role := args.Get(2).(*model.Role)
			role.ID = 1 // Simulate database ID assignment
		}).Return(nil)
	mockRoleRepo.On("AssignPermissionsWithTx", ctx, mock.AnythingOfType("*gorm.DB"), uint(1), req.PermissionIDs, createdBy).Return(nil)

	// Mock GetRole call for return value
	expectedRole := &model.Role{
		BaseModel:       model.BaseModel{ID: 1},
		Name:            req.Name,
		Code:            req.Code,
		Description:     req.Description,
		IsSystem:        false,
		SortOrder:       req.SortOrder,
		PermissionCount: 3,
		UserCount:       0,
	}
	mockRoleRepo.On("GetByID", ctx, uint(1)).Return(expectedRole, nil)
	mockRoleRepo.On("CountPermissions", ctx, uint(1)).Return(int64(3), nil)
	mockRoleRepo.On("CountUsers", ctx, uint(1)).Return(int64(0), nil)

	// Execute
	result, err := service.CreateRole(ctx, req, createdBy)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Code, result.Code)
	assert.Equal(t, req.Description, result.Description)

	mockRoleRepo.AssertExpectations(t)
	mockPermissionRepo.AssertExpectations(t)
}

func TestRoleService_CreateRole_CodeAlreadyExists(t *testing.T) {
	service, mockRoleRepo, _, _ := setupRoleServiceTest()
	ctx := context.Background()

	req := &request.CreateRoleRequest{
		Name: "Test Role",
		Code: "existing_role",
	}
	createdBy := uint(1)

	// Mock existing role
	existingRole := &model.Role{
		BaseModel: model.BaseModel{ID: 1},
		Code:      req.Code,
	}
	mockRoleRepo.On("GetByCode", ctx, req.Code).Return(existingRole, nil)

	// Execute
	result, err := service.CreateRole(ctx, req, createdBy)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "role code already exists")

	mockRoleRepo.AssertExpectations(t)
}

func TestRoleService_CreateRole_InvalidPermissionIDs(t *testing.T) {
	service, mockRoleRepo, mockPermissionRepo, _ := setupRoleServiceTest()
	ctx := context.Background()

	req := &request.CreateRoleRequest{
		Name:          "Test Role",
		Code:          "test_role",
		PermissionIDs: []uint{1, 2, 999}, // 999 doesn't exist
	}
	createdBy := uint(1)

	// Mock expectations
	mockRoleRepo.On("GetByCode", ctx, req.Code).Return(nil, errors.ErrRecordNotFound)
	mockPermissionRepo.On("GetByIDs", ctx, req.PermissionIDs).Return([]*model.Permission{
		{BaseModel: model.BaseModel{ID: 1}},
		{BaseModel: model.BaseModel{ID: 2}}, // Only 2 permissions found, not 3
	}, nil)

	// Execute
	result, err := service.CreateRole(ctx, req, createdBy)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "some permission IDs are invalid")

	mockRoleRepo.AssertExpectations(t)
	mockPermissionRepo.AssertExpectations(t)
}

func TestRoleService_UpdateRole_Success(t *testing.T) {
	service, mockRoleRepo, mockPermissionRepo, _ := setupRoleServiceTest()
	ctx := context.Background()

	roleID := uint(1)
	req := &request.UpdateRoleRequest{
		Name:          "Updated Role",
		Description:   "Updated description",
		PermissionIDs: []uint{1, 2},
		SortOrder:     2,
		Status:        1,
	}
	updatedBy := uint(1)

	// Mock existing role
	existingRole := &model.Role{
		BaseModel: model.BaseModel{ID: roleID},
		Name:      "Original Role",
		Code:      "test_role",
		IsSystem:  false,
	}
	mockRoleRepo.On("GetByID", ctx, roleID).Return(existingRole, nil)
	mockPermissionRepo.On("GetByIDs", ctx, req.PermissionIDs).Return([]*model.Permission{
		{BaseModel: model.BaseModel{ID: 1}},
		{BaseModel: model.BaseModel{ID: 2}},
	}, nil)

	// Mock update operations
	mockRoleRepo.On("UpdateWithTx", ctx, mock.AnythingOfType("*gorm.DB"), mock.AnythingOfType("*model.Role")).Return(nil)
	mockRoleRepo.On("AssignPermissionsWithTx", ctx, mock.AnythingOfType("*gorm.DB"), roleID, req.PermissionIDs, updatedBy).Return(nil)

	// Mock GetRole call for return value
	updatedRole := &model.Role{
		BaseModel:       model.BaseModel{ID: roleID},
		Name:            req.Name,
		Code:            existingRole.Code,
		Description:     req.Description,
		IsSystem:        false,
		SortOrder:       req.SortOrder,
		Status:          req.Status,
		PermissionCount: 2,
		UserCount:       0,
	}
	mockRoleRepo.On("GetByID", ctx, roleID).Return(updatedRole, nil)
	mockRoleRepo.On("CountPermissions", ctx, roleID).Return(int64(2), nil)
	mockRoleRepo.On("CountUsers", ctx, roleID).Return(int64(0), nil)

	// Execute
	result, err := service.UpdateRole(ctx, roleID, req, updatedBy)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Description, result.Description)

	mockRoleRepo.AssertExpectations(t)
	mockPermissionRepo.AssertExpectations(t)
}

func TestRoleService_UpdateRole_SystemRole(t *testing.T) {
	service, mockRoleRepo, _, _ := setupRoleServiceTest()
	ctx := context.Background()

	roleID := uint(1)
	req := &request.UpdateRoleRequest{
		Name: "Updated Role",
	}
	updatedBy := uint(1)

	// Mock system role
	systemRole := &model.Role{
		BaseModel: model.BaseModel{ID: roleID},
		Name:      "Admin",
		Code:      "admin",
		IsSystem:  true,
	}
	mockRoleRepo.On("GetByID", ctx, roleID).Return(systemRole, nil)

	// Execute
	result, err := service.UpdateRole(ctx, roleID, req, updatedBy)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot update system role")

	mockRoleRepo.AssertExpectations(t)
}

func TestRoleService_DeleteRole_Success(t *testing.T) {
	service, mockRoleRepo, _, _ := setupRoleServiceTest()
	ctx := context.Background()

	roleID := uint(1)

	mockRoleRepo.On("Delete", ctx, roleID).Return(nil)

	// Execute
	err := service.DeleteRole(ctx, roleID)

	// Assert
	assert.NoError(t, err)

	mockRoleRepo.AssertExpectations(t)
}

func TestRoleService_DeleteRole_Error(t *testing.T) {
	service, mockRoleRepo, _, _ := setupRoleServiceTest()
	ctx := context.Background()

	roleID := uint(1)
	expectedError := stderrors.New("database error")

	mockRoleRepo.On("Delete", ctx, roleID).Return(expectedError)

	// Execute
	err := service.DeleteRole(ctx, roleID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	mockRoleRepo.AssertExpectations(t)
}

func TestRoleService_ListRoles_Success(t *testing.T) {
	service, mockRoleRepo, _, _ := setupRoleServiceTest()
	ctx := context.Background()

	req := &request.ListRolesRequest{
		PaginationRequest: request.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
	}

	// Mock data
	roles := []*model.Role{
		{
			BaseModel: model.BaseModel{ID: 1},
			Name:      "Admin",
			Code:      "admin",
		},
		{
			BaseModel: model.BaseModel{ID: 2},
			Name:      "User",
			Code:      "user",
		},
	}
	total := int64(2)

	mockRoleRepo.On("List", ctx, req).Return(roles, total, nil)

	// Execute
	result, err := service.ListRoles(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(roles), len(result.Items))
	assert.Equal(t, total, result.Total)
	assert.Equal(t, req.GetPage(), result.Page)
	assert.Equal(t, req.GetPageSize(), result.PageSize)
	assert.Equal(t, 1, result.TotalPages)

	mockRoleRepo.AssertExpectations(t)
}

func TestRoleService_AssignPermissions_Success(t *testing.T) {
	service, mockRoleRepo, mockPermissionRepo, _ := setupRoleServiceTest()
	ctx := context.Background()

	roleID := uint(1)
	req := &request.RolePermissionRequest{
		PermissionIDs: []uint{1, 2, 3},
	}
	grantedBy := uint(1)

	// Mock role
	role := &model.Role{
		BaseModel: model.BaseModel{ID: roleID},
		Name:      "Test Role",
		IsSystem:  false,
	}
	mockRoleRepo.On("GetByID", ctx, roleID).Return(role, nil)
	mockPermissionRepo.On("GetByIDs", ctx, req.PermissionIDs).Return([]*model.Permission{
		{BaseModel: model.BaseModel{ID: 1}},
		{BaseModel: model.BaseModel{ID: 2}},
		{BaseModel: model.BaseModel{ID: 3}},
	}, nil)
	mockRoleRepo.On("AssignPermissionsWithTx", ctx, mock.AnythingOfType("*gorm.DB"), roleID, req.PermissionIDs, grantedBy).Return(nil)

	// Execute
	err := service.AssignPermissions(ctx, roleID, req, grantedBy)

	// Assert
	assert.NoError(t, err)

	mockRoleRepo.AssertExpectations(t)
	mockPermissionRepo.AssertExpectations(t)
}

func TestRoleService_AssignPermissions_SystemRole(t *testing.T) {
	service, mockRoleRepo, _, _ := setupRoleServiceTest()
	ctx := context.Background()

	roleID := uint(1)
	req := &request.RolePermissionRequest{
		PermissionIDs: []uint{1, 2, 3},
	}
	grantedBy := uint(1)

	// Mock system role
	systemRole := &model.Role{
		BaseModel: model.BaseModel{ID: roleID},
		Name:      "Admin",
		IsSystem:  true,
	}
	mockRoleRepo.On("GetByID", ctx, roleID).Return(systemRole, nil)

	// Execute
	err := service.AssignPermissions(ctx, roleID, req, grantedBy)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot modify system role permissions")

	mockRoleRepo.AssertExpectations(t)
}

func TestRoleService_RemovePermissions_Success(t *testing.T) {
	service, mockRoleRepo, _, _ := setupRoleServiceTest()
	ctx := context.Background()

	roleID := uint(1)
	req := &request.RolePermissionRequest{
		PermissionIDs: []uint{1, 2},
	}

	// Mock role
	role := &model.Role{
		BaseModel: model.BaseModel{ID: roleID},
		Name:      "Test Role",
		IsSystem:  false,
	}
	mockRoleRepo.On("GetByID", ctx, roleID).Return(role, nil)
	mockRoleRepo.On("RemovePermissions", ctx, roleID, req.PermissionIDs).Return(nil)

	// Execute
	err := service.RemovePermissions(ctx, roleID, req)

	// Assert
	assert.NoError(t, err)

	mockRoleRepo.AssertExpectations(t)
}

func TestRoleService_BatchUpdateRoleStatus_Success(t *testing.T) {
	service, mockRoleRepo, _, _ := setupRoleServiceTest()
	ctx := context.Background()

	ids := []uint{1, 2, 3}
	status := 0

	mockRoleRepo.On("BatchUpdateStatus", ctx, ids, status).Return(nil)

	// Execute
	err := service.BatchUpdateRoleStatus(ctx, ids, status)

	// Assert
	assert.NoError(t, err)

	mockRoleRepo.AssertExpectations(t)
}

// Benchmark tests
func BenchmarkRoleService_CreateRole(b *testing.B) {
	service, mockRoleRepo, mockPermissionRepo, _ := setupRoleServiceTest()
	ctx := context.Background()

	req := &request.CreateRoleRequest{
		Name:          "Benchmark Role",
		Code:          "benchmark_role",
		Description:   "Benchmark role description",
		PermissionIDs: []uint{1, 2, 3},
		SortOrder:     1,
	}
	createdBy := uint(1)

	// Setup mocks for benchmark
	mockRoleRepo.On("GetByCode", ctx, mock.AnythingOfType("string")).Return(nil, errors.ErrRecordNotFound)
	mockPermissionRepo.On("GetByIDs", ctx, mock.AnythingOfType("[]uint")).Return([]*model.Permission{
		{BaseModel: model.BaseModel{ID: 1}},
		{BaseModel: model.BaseModel{ID: 2}},
		{BaseModel: model.BaseModel{ID: 3}},
	}, nil)
	mockRoleRepo.On("CreateWithTx", ctx, mock.AnythingOfType("*gorm.DB"), mock.AnythingOfType("*model.Role")).
		Run(func(args mock.Arguments) {
			role := args.Get(2).(*model.Role)
			role.ID = 1
		}).Return(nil)
	mockRoleRepo.On("AssignPermissionsWithTx", ctx, mock.AnythingOfType("*gorm.DB"), mock.AnythingOfType("uint"), mock.AnythingOfType("[]uint"), mock.AnythingOfType("uint")).Return(nil)
	mockRoleRepo.On("GetByID", ctx, mock.AnythingOfType("uint")).Return(&model.Role{
		BaseModel: model.BaseModel{ID: 1},
		Name:      req.Name,
		Code:      req.Code,
	}, nil)
	mockRoleRepo.On("CountPermissions", ctx, mock.AnythingOfType("uint")).Return(int64(3), nil)
	mockRoleRepo.On("CountUsers", ctx, mock.AnythingOfType("uint")).Return(int64(0), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.Code = fmt.Sprintf("benchmark_role_%d", i)
		_, _ = service.CreateRole(ctx, req, createdBy)
	}
}
