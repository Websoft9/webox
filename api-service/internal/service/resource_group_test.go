package service

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/model"
	pkgErrors "api-service/pkg/errors"
	"api-service/pkg/logger"
)

// MockResourceGroupRepository mocks ResourceGroupRepository interface
type MockResourceGroupRepository struct {
	mock.Mock
}

func (m *MockResourceGroupRepository) Create(ctx context.Context, rg *model.ResourceGroup) error {
	args := m.Called(ctx, rg)
	if args.Error(0) == nil && rg.ID == 0 {
		// Simulate ID generation
		rg.ID = 1
		rg.Code = "rg_test123456"
	}
	return args.Error(0)
}

func (m *MockResourceGroupRepository) GetByID(ctx context.Context, id uint) (*model.ResourceGroup, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ResourceGroup), args.Error(1)
}

func (m *MockResourceGroupRepository) GetByCode(ctx context.Context, code string) (*model.ResourceGroup, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ResourceGroup), args.Error(1)
}

func (m *MockResourceGroupRepository) GetList(ctx context.Context, req *request.GetResourceGroupListRequest, ownerID uint) ([]*model.ResourceGroup, int64, error) {
	args := m.Called(ctx, req, ownerID)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*model.ResourceGroup), args.Get(1).(int64), args.Error(2)
}

func (m *MockResourceGroupRepository) Update(ctx context.Context, rg *model.ResourceGroup) error {
	args := m.Called(ctx, rg)
	return args.Error(0)
}

func (m *MockResourceGroupRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockResourceGroupRepository) CheckNameExists(ctx context.Context, projectID uint, name string, excludeID *uint) (bool, error) {
	args := m.Called(ctx, projectID, name, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockResourceGroupRepository) GetDefaultResourceGroup(ctx context.Context, projectID uint) (*model.ResourceGroup, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ResourceGroup), args.Error(1)
}

func (m *MockResourceGroupRepository) GetProjectIDFromResourceCode(ctx context.Context, resourceCode string) (uint, error) {
	args := m.Called(ctx, resourceCode)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockResourceGroupRepository) GetResourcesByGroupID(ctx context.Context, groupID uint, req *request.GetResourceGroupResourcesRequest) ([]response.ResourceItemResponse, int64, error) {
	args := m.Called(ctx, groupID, req)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]response.ResourceItemResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockResourceGroupRepository) MoveResourcesToGroup(ctx context.Context, resourceCodes []string, targetGroupID *uint, projectID uint) error {
	args := m.Called(ctx, resourceCodes, targetGroupID, projectID)
	return args.Error(0)
}

func (m *MockResourceGroupRepository) GetResourceStatistics(ctx context.Context, projectID, resourceGroupID *uint) (map[string]int, error) {
	args := m.Called(ctx, projectID, resourceGroupID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int), args.Error(1)
}

// setupResourceGroupTest sets up test dependencies
func setupResourceGroupTest(t *testing.T) (*resourceGroupService, *MockResourceGroupRepository) {
	mockRepo := new(MockResourceGroupRepository)
	mockLogger := logger.NewZapLogger(logger.InfoLevel, io.Discard)

	service := &resourceGroupService{
		repo:   mockRepo,
		logger: mockLogger,
	}

	return service, mockRepo
}

func TestCreateResourceGroup(t *testing.T) {
	service, mockRepo := setupResourceGroupTest(t)
	ctx := context.Background()
	ownerID := uint(1)

	testCases := []struct {
		name          string
		req           *request.CreateResourceGroupRequest
		setupMock     func()
		expectedError error
		validate      func(t *testing.T, result *response.ResourceGroupResponse)
	}{
		{
			name: "Success - Create resource group",
			req: &request.CreateResourceGroupRequest{
				ProjectID:   10,
				Name:        "Production Environment",
				Description: stringPtr("Production resource group"),
				SortOrder:   intPtr(1),
			},
			setupMock: func() {
				mockRepo.On("CheckNameExists", ctx, uint(10), "Production Environment", (*uint)(nil)).
					Return(false, nil).Once()
				mockRepo.On("Create", ctx, mock.AnythingOfType("*model.ResourceGroup")).
					Return(nil).Once()
			},
			expectedError: nil,
			validate: func(t *testing.T, result *response.ResourceGroupResponse) {
				assert.Equal(t, "Production Environment", result.Name)
				assert.Equal(t, uint(10), result.ProjectID)
				assert.Equal(t, uint(1), result.OwnerID)
				assert.Equal(t, 1, result.SortOrder)
				assert.NotEmpty(t, result.Code)
			},
		},
		{
			name: "Success - Default sort order",
			req: &request.CreateResourceGroupRequest{
				ProjectID: 10,
				Name:      "Test Group",
			},
			setupMock: func() {
				mockRepo.On("CheckNameExists", ctx, uint(10), "Test Group", (*uint)(nil)).
					Return(false, nil).Once()
				mockRepo.On("Create", ctx, mock.AnythingOfType("*model.ResourceGroup")).
					Return(nil).Once()
			},
			expectedError: nil,
			validate: func(t *testing.T, result *response.ResourceGroupResponse) {
				assert.Equal(t, 0, result.SortOrder) // Default value
			},
		},
		{
			name: "Failure - Name already exists",
			req: &request.CreateResourceGroupRequest{
				ProjectID: 10,
				Name:      "Existing Group",
			},
			setupMock: func() {
				mockRepo.On("CheckNameExists", ctx, uint(10), "Existing Group", (*uint)(nil)).
					Return(true, nil).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeResourceAlreadyExists),
		},
		{
			name: "Failure - Check name exists error",
			req: &request.CreateResourceGroupRequest{
				ProjectID: 10,
				Name:      "Test Group",
			},
			setupMock: func() {
				mockRepo.On("CheckNameExists", ctx, uint(10), "Test Group", (*uint)(nil)).
					Return(false, pkgErrors.NewAppError(pkgErrors.CodeRecordQueryFailed)).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordQueryFailed),
		},
		{
			name: "Failure - Repository create error",
			req: &request.CreateResourceGroupRequest{
				ProjectID: 10,
				Name:      "Test Group",
			},
			setupMock: func() {
				mockRepo.On("CheckNameExists", ctx, uint(10), "Test Group", (*uint)(nil)).
					Return(false, nil).Once()
				mockRepo.On("Create", ctx, mock.AnythingOfType("*model.ResourceGroup")).
					Return(pkgErrors.NewAppError(pkgErrors.CodeRecordCreateFailed)).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordCreateFailed),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			result, err := service.CreateResourceGroup(ctx, tc.req, ownerID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tc.validate != nil {
					tc.validate(t, result)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetResourceGroup(t *testing.T) {
	service, mockRepo := setupResourceGroupTest(t)
	ctx := context.Background()

	testCases := []struct {
		name          string
		id            uint
		setupMock     func()
		expectedError error
	}{
		{
			name: "Success - Get resource group",
			id:   1,
			setupMock: func() {
				mockRepo.On("GetByID", ctx, uint(1)).
					Return(&model.ResourceGroup{
						ID:        1,
						ProjectID: 10,
						Name:      "Test Group",
						Code:      "rg_test123",
						OwnerID:   1,
					}, nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Failure - Not found",
			id:   999,
			setupMock: func() {
				mockRepo.On("GetByID", ctx, uint(999)).
					Return(nil, pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound)).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			result, err := service.GetResourceGroup(ctx, tc.id)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetResourceGroupList(t *testing.T) {
	service, mockRepo := setupResourceGroupTest(t)
	ctx := context.Background()
	ownerID := uint(1)

	testCases := []struct {
		name          string
		req           *request.GetResourceGroupListRequest
		setupMock     func()
		expectedError error
		expectedCount int
	}{
		{
			name: "Success - Get list with results",
			req: &request.GetResourceGroupListRequest{
				BaseListRequest: common.BaseListRequest{
					PaginationRequest: common.PaginationRequest{
						Page:     1,
						PageSize: 20,
					},
				},
			},
			setupMock: func() {
				mockGroups := []*model.ResourceGroup{
					{ID: 1, Name: "Group 1"},
					{ID: 2, Name: "Group 2"},
				}
				mockRepo.On("GetList", ctx, mock.AnythingOfType("*request.GetResourceGroupListRequest"), ownerID).
					Return(mockGroups, int64(2), nil).Once()
			},
			expectedError: nil,
			expectedCount: 2,
		},
		{
			name: "Success - Empty list",
			req: &request.GetResourceGroupListRequest{
				BaseListRequest: common.BaseListRequest{
					PaginationRequest: common.PaginationRequest{
						Page:     1,
						PageSize: 20,
					},
				},
			},
			setupMock: func() {
				mockRepo.On("GetList", ctx, mock.AnythingOfType("*request.GetResourceGroupListRequest"), ownerID).
					Return([]*model.ResourceGroup{}, int64(0), nil).Once()
			},
			expectedError: nil,
			expectedCount: 0,
		},
		{
			name: "Failure - Repository error",
			req: &request.GetResourceGroupListRequest{
				BaseListRequest: common.BaseListRequest{
					PaginationRequest: common.PaginationRequest{
						Page:     1,
						PageSize: 20,
					},
				},
			},
			setupMock: func() {
				mockRepo.On("GetList", ctx, mock.AnythingOfType("*request.GetResourceGroupListRequest"), ownerID).
					Return(nil, int64(0), pkgErrors.NewAppError(pkgErrors.CodeRecordQueryFailed)).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordQueryFailed),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			result, err := service.GetResourceGroupList(ctx, tc.req, ownerID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedCount, len(result.Items.([]interface{})))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateResourceGroup(t *testing.T) {
	service, mockRepo := setupResourceGroupTest(t)
	ctx := context.Background()
	ownerID := uint(1)
	groupID := uint(1)

	testCases := []struct {
		name          string
		req           *request.UpdateResourceGroupRequest
		setupMock     func()
		expectedError error
	}{
		{
			name: "Success - Update name",
			req: &request.UpdateResourceGroupRequest{
				Name: stringPtr("Updated Name"),
			},
			setupMock: func() {
				mockRepo.On("GetByID", ctx, groupID).
					Return(&model.ResourceGroup{
						ID:        groupID,
						ProjectID: 10,
						Name:      "Old Name",
						OwnerID:   ownerID,
					}, nil).Once()
				mockRepo.On("CheckNameExists", ctx, uint(10), "Updated Name", &groupID).
					Return(false, nil).Once()
				mockRepo.On("Update", ctx, mock.AnythingOfType("*model.ResourceGroup")).
					Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Success - Update description and sort order",
			req: &request.UpdateResourceGroupRequest{
				Description: stringPtr("New description"),
				SortOrder:   intPtr(5),
			},
			setupMock: func() {
				mockRepo.On("GetByID", ctx, groupID).
					Return(&model.ResourceGroup{
						ID:        groupID,
						ProjectID: 10,
						Name:      "Test Group",
						OwnerID:   ownerID,
					}, nil).Once()
				mockRepo.On("Update", ctx, mock.AnythingOfType("*model.ResourceGroup")).
					Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Failure - Resource group not found",
			req: &request.UpdateResourceGroupRequest{
				Name: stringPtr("Updated Name"),
			},
			setupMock: func() {
				mockRepo.On("GetByID", ctx, groupID).
					Return(nil, pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound)).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound),
		},
		{
			name: "Failure - Name already exists",
			req: &request.UpdateResourceGroupRequest{
				Name: stringPtr("Existing Name"),
			},
			setupMock: func() {
				mockRepo.On("GetByID", ctx, groupID).
					Return(&model.ResourceGroup{
						ID:        groupID,
						ProjectID: 10,
						Name:      "Old Name",
						OwnerID:   ownerID,
					}, nil).Once()
				mockRepo.On("CheckNameExists", ctx, uint(10), "Existing Name", &groupID).
					Return(true, nil).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeResourceAlreadyExists),
		},
		{
			name: "Failure - Update error",
			req: &request.UpdateResourceGroupRequest{
				Description: stringPtr("New description"),
			},
			setupMock: func() {
				mockRepo.On("GetByID", ctx, groupID).
					Return(&model.ResourceGroup{
						ID:        groupID,
						ProjectID: 10,
						Name:      "Test Group",
						OwnerID:   ownerID,
					}, nil).Once()
				mockRepo.On("Update", ctx, mock.AnythingOfType("*model.ResourceGroup")).
					Return(pkgErrors.NewAppError(pkgErrors.CodeRecordUpdateFailed)).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordUpdateFailed),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			result, err := service.UpdateResourceGroup(ctx, groupID, tc.req, ownerID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteResourceGroup(t *testing.T) {
	service, mockRepo := setupResourceGroupTest(t)
	ctx := context.Background()

	testCases := []struct {
		name          string
		id            uint
		setupMock     func()
		expectedError error
	}{
		{
			name: "Success - Delete resource group",
			id:   1,
			setupMock: func() {
				mockRepo.On("GetByID", ctx, uint(1)).
					Return(&model.ResourceGroup{
						ID:   1,
						Name: "Test Group",
					}, nil).Once()
				mockRepo.On("Delete", ctx, uint(1)).
					Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Failure - Resource group not found",
			id:   999,
			setupMock: func() {
				mockRepo.On("GetByID", ctx, uint(999)).
					Return(nil, pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound)).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound),
		},
		{
			name: "Failure - Delete error",
			id:   1,
			setupMock: func() {
				mockRepo.On("GetByID", ctx, uint(1)).
					Return(&model.ResourceGroup{
						ID:   1,
						Name: "Test Group",
					}, nil).Once()
				mockRepo.On("Delete", ctx, uint(1)).
					Return(pkgErrors.NewAppError(pkgErrors.CodeRecordDeleteFailed)).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordDeleteFailed),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			err := service.DeleteResourceGroup(ctx, tc.id)

			if tc.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetResourcesByGroupID(t *testing.T) {
	service, mockRepo := setupResourceGroupTest(t)
	ctx := context.Background()

	testCases := []struct {
		name          string
		groupID       uint
		req           *request.GetResourceGroupResourcesRequest
		setupMock     func()
		expectedError error
		expectedCount int
	}{
		{
			name:    "Success - Get resources with results",
			groupID: 1,
			req: &request.GetResourceGroupResourcesRequest{
				BaseListRequest: common.BaseListRequest{
					PaginationRequest: common.PaginationRequest{
						Page:     1,
						PageSize: 20,
					},
				},
			},
			setupMock: func() {
				mockResources := []response.ResourceItemResponse{
					{ID: 1, Name: "Server 1", ResourceType: "server"},
					{ID: 2, Name: "DB 1", ResourceType: "database"},
				}
				mockRepo.On("GetResourcesByGroupID", ctx, uint(1), mock.AnythingOfType("*request.GetResourceGroupResourcesRequest")).
					Return(mockResources, int64(2), nil).Once()
			},
			expectedError: nil,
			expectedCount: 2,
		},
		{
			name:    "Success - Filter by resource type",
			groupID: 1,
			req: &request.GetResourceGroupResourcesRequest{
				BaseListRequest: common.BaseListRequest{
					PaginationRequest: common.PaginationRequest{
						Page:     1,
						PageSize: 20,
					},
				},
				ResourceType: stringPtr("server"),
			},
			setupMock: func() {
				mockResources := []response.ResourceItemResponse{
					{ID: 1, Name: "Server 1", ResourceType: "server"},
				}
				mockRepo.On("GetResourcesByGroupID", ctx, uint(1), mock.AnythingOfType("*request.GetResourceGroupResourcesRequest")).
					Return(mockResources, int64(1), nil).Once()
			},
			expectedError: nil,
			expectedCount: 1,
		},
		{
			name:    "Success - Empty result",
			groupID: 1,
			req: &request.GetResourceGroupResourcesRequest{
				BaseListRequest: common.BaseListRequest{
					PaginationRequest: common.PaginationRequest{
						Page:     1,
						PageSize: 20,
					},
				},
			},
			setupMock: func() {
				mockRepo.On("GetResourcesByGroupID", ctx, uint(1), mock.AnythingOfType("*request.GetResourceGroupResourcesRequest")).
					Return([]response.ResourceItemResponse{}, int64(0), nil).Once()
			},
			expectedError: nil,
			expectedCount: 0,
		},
		{
			name:    "Failure - Repository error",
			groupID: 1,
			req: &request.GetResourceGroupResourcesRequest{
				BaseListRequest: common.BaseListRequest{
					PaginationRequest: common.PaginationRequest{
						Page:     1,
						PageSize: 20,
					},
				},
			},
			setupMock: func() {
				mockRepo.On("GetResourcesByGroupID", ctx, uint(1), mock.AnythingOfType("*request.GetResourceGroupResourcesRequest")).
					Return(nil, int64(0), pkgErrors.NewAppError(pkgErrors.CodeRecordQueryFailed)).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordQueryFailed),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			result, err := service.GetResourcesByGroupID(ctx, tc.groupID, tc.req)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedCount, len(result.Resources))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMoveResourcesToGroup(t *testing.T) {
	service, mockRepo := setupResourceGroupTest(t)
	ctx := context.Background()
	ownerID := uint(1)

	testCases := []struct {
		name          string
		req           *request.MoveResourcesToGroupRequest
		setupMock     func()
		expectedError error
	}{
		{
			name: "Success - Move to specified group",
			req: &request.MoveResourcesToGroupRequest{
				ResourceCodes:   []string{"server_123", "database_456"},
				ResourceGroupID: uintPtr(10),
			},
			setupMock: func() {
				mockRepo.On("GetByID", ctx, uint(10)).
					Return(&model.ResourceGroup{
						ID:        10,
						ProjectID: 5,
						Name:      "Target Group",
					}, nil).Once()
				mockRepo.On("MoveResourcesToGroup", ctx, []string{"server_123", "database_456"}, uintPtr(10), uint(5)).
					Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Success - Move to default group (null resource_group_id)",
			req: &request.MoveResourcesToGroupRequest{
				ResourceCodes:   []string{"server_123"},
				ResourceGroupID: nil,
			},
			setupMock: func() {
				mockRepo.On("GetProjectIDFromResourceCode", ctx, "server_123").
					Return(uint(5), nil).Once()
				mockRepo.On("MoveResourcesToGroup", ctx, []string{"server_123"}, (*uint)(nil), uint(5)).
					Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Failure - Target group not found",
			req: &request.MoveResourcesToGroupRequest{
				ResourceCodes:   []string{"server_123"},
				ResourceGroupID: uintPtr(999),
			},
			setupMock: func() {
				mockRepo.On("GetByID", ctx, uint(999)).
					Return(nil, pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound)).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound),
		},
		{
			name: "Failure - Get project ID error",
			req: &request.MoveResourcesToGroupRequest{
				ResourceCodes:   []string{"invalid_code"},
				ResourceGroupID: nil,
			},
			setupMock: func() {
				mockRepo.On("GetProjectIDFromResourceCode", ctx, "invalid_code").
					Return(uint(0), pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound)).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordNotFound),
		},
		{
			name: "Failure - Move operation error",
			req: &request.MoveResourcesToGroupRequest{
				ResourceCodes:   []string{"server_123"},
				ResourceGroupID: uintPtr(10),
			},
			setupMock: func() {
				mockRepo.On("GetByID", ctx, uint(10)).
					Return(&model.ResourceGroup{
						ID:        10,
						ProjectID: 5,
					}, nil).Once()
				mockRepo.On("MoveResourcesToGroup", ctx, []string{"server_123"}, uintPtr(10), uint(5)).
					Return(pkgErrors.NewAppError(pkgErrors.CodeRecordUpdateFailed)).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordUpdateFailed),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			err := service.MoveResourcesToGroup(ctx, tc.req, ownerID)

			if tc.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetResourceStatistics(t *testing.T) {
	service, mockRepo := setupResourceGroupTest(t)
	ctx := context.Background()

	testCases := []struct {
		name          string
		req           *request.GetResourceStatisticsRequest
		setupMock     func()
		expectedError error
		expectedTotal int64
	}{
		{
			name: "Success - Statistics by project",
			req: &request.GetResourceStatisticsRequest{
				ProjectID: uintPtr(10),
			},
			setupMock: func() {
				mockStats := map[string]int{
					"server":   12,
					"database": 8,
					"secret":   5,
				}
				mockRepo.On("GetResourceStatistics", ctx, uintPtr(10), (*uint)(nil)).
					Return(mockStats, nil).Once()
			},
			expectedError: nil,
			expectedTotal: 25,
		},
		{
			name: "Success - Statistics by resource group",
			req: &request.GetResourceStatisticsRequest{
				ResourceGroupID: uintPtr(5),
			},
			setupMock: func() {
				mockStats := map[string]int{
					"server":   3,
					"database": 2,
				}
				mockRepo.On("GetResourceStatistics", ctx, (*uint)(nil), uintPtr(5)).
					Return(mockStats, nil).Once()
			},
			expectedError: nil,
			expectedTotal: 5,
		},
		{
			name: "Success - All resources (no filter)",
			req:  &request.GetResourceStatisticsRequest{},
			setupMock: func() {
				mockStats := map[string]int{
					"server":      45,
					"database":    32,
					"secret":      28,
					"application": 15,
				}
				mockRepo.On("GetResourceStatistics", ctx, (*uint)(nil), (*uint)(nil)).
					Return(mockStats, nil).Once()
			},
			expectedError: nil,
			expectedTotal: 120,
		},
		{
			name: "Failure - Both parameters provided (mutual exclusion)",
			req: &request.GetResourceStatisticsRequest{
				ProjectID:       uintPtr(10),
				ResourceGroupID: uintPtr(5),
			},
			setupMock:     func() {},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeInvalidParameterFormat),
		},
		{
			name: "Failure - Repository error",
			req: &request.GetResourceStatisticsRequest{
				ProjectID: uintPtr(10),
			},
			setupMock: func() {
				mockRepo.On("GetResourceStatistics", ctx, uintPtr(10), (*uint)(nil)).
					Return(nil, pkgErrors.NewAppError(pkgErrors.CodeRecordQueryFailed)).Once()
			},
			expectedError: pkgErrors.NewAppError(pkgErrors.CodeRecordQueryFailed),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			result, err := service.GetResourceStatistics(ctx, tc.req)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedTotal, result.Total)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func uintPtr(u uint) *uint {
	return &u
}
