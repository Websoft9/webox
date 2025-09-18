package service

import (
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockTagRepository is a mock implementation of TagRepository for testing
type MockTagRepository struct {
	mock.Mock
}

func (m *MockTagRepository) CreateTag(ctx context.Context, tag *model.Tag) error {
	args := m.Called(ctx, tag)
	return args.Error(0)
}

func (m *MockTagRepository) GetTagByID(ctx context.Context, id uint64) (*model.Tag, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Tag), args.Error(1)
}

func (m *MockTagRepository) GetTagByName(ctx context.Context, name string) (*model.Tag, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Tag), args.Error(1)
}

func (m *MockTagRepository) UpdateTag(ctx context.Context, tag *model.Tag) error {
	args := m.Called(ctx, tag)
	return args.Error(0)
}

func (m *MockTagRepository) DeleteTag(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTagRepository) ListTags(ctx context.Context, search string, excludeIDs []uint64) ([]*model.Tag, error) {
	args := m.Called(ctx, search, excludeIDs)
	return args.Get(0).([]*model.Tag), args.Error(1)
}

func (m *MockTagRepository) ListTagsWithUsageCount(ctx context.Context, search string, excludeIDs []uint64) ([]*model.Tag, error) {
	args := m.Called(ctx, search, excludeIDs)
	return args.Get(0).([]*model.Tag), args.Error(1)
}

func (m *MockTagRepository) ExistsTagByName(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

func (m *MockTagRepository) ExistsTagByNameExcludeID(ctx context.Context, name string, excludeID uint64) (bool, error) {
	args := m.Called(ctx, name, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockTagRepository) CreateTagging(ctx context.Context, tagging *model.Tagging) error {
	args := m.Called(ctx, tagging)
	return args.Error(0)
}

func (m *MockTagRepository) GetTaggingsByResourceID(ctx context.Context, resourceID uint64) ([]*model.Tagging, error) {
	args := m.Called(ctx, resourceID)
	return args.Get(0).([]*model.Tagging), args.Error(1)
}

func (m *MockTagRepository) GetTaggingsByTagID(ctx context.Context, tagID uint64) ([]*model.Tagging, error) {
	args := m.Called(ctx, tagID)
	return args.Get(0).([]*model.Tagging), args.Error(1)
}

func (m *MockTagRepository) DeleteTagging(ctx context.Context, tagID, resourceID uint64) error {
	args := m.Called(ctx, tagID, resourceID)
	return args.Error(0)
}

func (m *MockTagRepository) DeleteTaggingsByResourceID(ctx context.Context, resourceID uint64) error {
	args := m.Called(ctx, resourceID)
	return args.Error(0)
}

func (m *MockTagRepository) DeleteTaggingsByTagIDs(ctx context.Context, resourceID uint64, tagIDs []uint64) error {
	args := m.Called(ctx, resourceID, tagIDs)
	return args.Error(0)
}

func (m *MockTagRepository) CreateTaggingsBatch(ctx context.Context, taggings []*model.Tagging) error {
	args := m.Called(ctx, taggings)
	return args.Error(0)
}

func (m *MockTagRepository) ExistsTagging(ctx context.Context, tagID, resourceID uint64) (bool, error) {
	args := m.Called(ctx, tagID, resourceID)
	return args.Bool(0), args.Error(1)
}

func (m *MockTagRepository) SearchResourcesByTags(ctx context.Context, tagIDs []uint64, operation string, offset, limit int) ([]*model.Tagging, int64, error) {
	args := m.Called(ctx, tagIDs, operation, offset, limit)
	return args.Get(0).([]*model.Tagging), args.Get(1).(int64), args.Error(2)
}

func (m *MockTagRepository) SearchTagsByName(ctx context.Context, query string) ([]*model.Tag, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]*model.Tag), args.Error(1)
}

// Test setup
func setupTagServiceTest() (*tagService, *MockTagRepository) {
	mockRepo := &MockTagRepository{}
	mockLogger := logger.NewZapLogger(logger.InfoLevel, nil)
	mockI18n := &i18n.I18n{}

	service := &tagService{
		tagRepo:      mockRepo,
		db:           &gorm.DB{}, // Use real gorm.DB for transaction tests
		logger:       mockLogger,
		i18nInstance: mockI18n,
	}

	return service, mockRepo
}

func TestTagService_CreateTag(t *testing.T) {
	service, mockRepo := setupTagServiceTest()
	ctx := context.Background()

	tests := []struct {
		name        string
		req         *request.TagCreateRequest
		userID      uint64
		setupMocks  func()
		expectedErr error
	}{
		{
			name: "successful tag creation",
			req: &request.TagCreateRequest{
				Name:        "test-tag",
				Color:       "#ff0000",
				Description: "Test tag description",
			},
			userID: 1,
			setupMocks: func() {
				mockRepo.On("ExistsTagByName", ctx, "test-tag").Return(false, nil)
				mockRepo.On("CreateTag", ctx, mock.AnythingOfType("*model.Tag")).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "tag name already exists",
			req: &request.TagCreateRequest{
				Name:        "existing-tag",
				Color:       "#ff0000",
				Description: "Test tag description",
			},
			userID: 1,
			setupMocks: func() {
				mockRepo.On("ExistsTagByName", ctx, "existing-tag").Return(true, nil)
			},
			expectedErr: errors.NewAppError(errors.CodeResourceAlreadyExists, "tag name already exists"),
		},
		{
			name: "database error on check existence",
			req: &request.TagCreateRequest{
				Name:        "test-tag",
				Color:       "#ff0000",
				Description: "Test tag description",
			},
			userID: 1,
			setupMocks: func() {
				mockRepo.On("ExistsTagByName", ctx, "test-tag").Return(false, gorm.ErrInvalidDB)
			},
			expectedErr: errors.NewAppError(errors.CodeInternalError, "failed to check tag existence"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			// Setup mocks
			tt.setupMocks()

			// Execute test
			result, err := service.CreateTag(ctx, tt.req, tt.userID)

			// Assertions
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.req.Name, result.Name)
				assert.Equal(t, tt.req.Color, result.Color)
				assert.Equal(t, tt.req.Description, result.Description)
				assert.Equal(t, tt.userID, result.CreatedBy)
			}

			// Verify all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTagService_GetTag(t *testing.T) {
	service, mockRepo := setupTagServiceTest()
	ctx := context.Background()

	now := time.Now()
	expectedTag := &model.Tag{
		ID:          1,
		Name:        "test-tag",
		Color:       "#ff0000",
		Description: "Test tag",
		CreatedBy:   1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name        string
		tagID       uint64
		setupMocks  func()
		expectedErr error
		expectedTag *response.TagResponse
	}{
		{
			name:  "successful tag retrieval",
			tagID: 1,
			setupMocks: func() {
				mockRepo.On("GetTagByID", ctx, uint64(1)).Return(expectedTag, nil)
			},
			expectedErr: nil,
			expectedTag: &response.TagResponse{
				ID:          1,
				Name:        "test-tag",
				Color:       "#ff0000",
				Description: "Test tag",
				CreatedBy:   1,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
		{
			name:  "tag not found",
			tagID: 999,
			setupMocks: func() {
				mockRepo.On("GetTagByID", ctx, uint64(999)).Return((*model.Tag)(nil), gorm.ErrRecordNotFound)
			},
			expectedErr: errors.NewAppError(errors.CodeResourceNotFound, "tag not found"),
			expectedTag: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			// Setup mocks
			tt.setupMocks()

			// Execute test
			result, err := service.GetTag(ctx, tt.tagID)

			// Assertions
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedTag.ID, result.ID)
				assert.Equal(t, tt.expectedTag.Name, result.Name)
				assert.Equal(t, tt.expectedTag.Color, result.Color)
				assert.Equal(t, tt.expectedTag.Description, result.Description)
			}

			// Verify all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTagService_UpdateTag(t *testing.T) {
	service, mockRepo := setupTagServiceTest()
	ctx := context.Background()

	now := time.Now()
	existingTag := &model.Tag{
		ID:          1,
		Name:        "old-name",
		Color:       "#ff0000",
		Description: "Old description",
		CreatedBy:   1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name        string
		tagID       uint64
		req         *request.TagUpdateRequest
		userID      uint64
		setupMocks  func()
		expectedErr error
	}{
		{
			name:  "successful tag update",
			tagID: 1,
			req: &request.TagUpdateRequest{
				Name:        "new-name",
				Color:       "#00ff00",
				Description: "New description",
			},
			userID: 1,
			setupMocks: func() {
				mockRepo.On("GetTagByID", ctx, uint64(1)).Return(existingTag, nil)
				mockRepo.On("ExistsTagByNameExcludeID", ctx, "new-name", uint64(1)).Return(false, nil)
				mockRepo.On("UpdateTag", ctx, mock.AnythingOfType("*model.Tag")).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:  "tag not found",
			tagID: 999,
			req: &request.TagUpdateRequest{
				Name:        "new-name",
				Color:       "#00ff00",
				Description: "New description",
			},
			userID: 1,
			setupMocks: func() {
				mockRepo.On("GetTagByID", ctx, uint64(999)).Return((*model.Tag)(nil), gorm.ErrRecordNotFound)
			},
			expectedErr: errors.NewAppError(errors.CodeResourceNotFound, "tag not found"),
		},
		{
			name:  "name conflict",
			tagID: 1,
			req: &request.TagUpdateRequest{
				Name:        "conflicting-name",
				Color:       "#00ff00",
				Description: "New description",
			},
			userID: 1,
			setupMocks: func() {
				mockRepo.On("GetTagByID", ctx, uint64(1)).Return(existingTag, nil)
				mockRepo.On("ExistsTagByNameExcludeID", ctx, "conflicting-name", uint64(1)).Return(true, nil)
			},
			expectedErr: errors.NewAppError(errors.CodeResourceAlreadyExists, "tag name already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			// Setup mocks
			tt.setupMocks()

			// Execute test
			result, err := service.UpdateTag(ctx, tt.tagID, tt.req, tt.userID)

			// Assertions
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.req.Name, result.Name)
				assert.Equal(t, tt.req.Color, result.Color)
				assert.Equal(t, tt.req.Description, result.Description)
			}

			// Verify all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTagService_DeleteTag(t *testing.T) {
	service, mockRepo := setupTagServiceTest()
	ctx := context.Background()

	now := time.Now()
	existingTag := &model.Tag{
		ID:          1,
		Name:        "test-tag",
		Color:       "#ff0000",
		Description: "Test tag",
		CreatedBy:   1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name        string
		tagID       uint64
		userID      uint64
		setupMocks  func()
		expectedErr error
	}{
		{
			name:   "successful tag deletion",
			tagID:  1,
			userID: 1,
			setupMocks: func() {
				mockRepo.On("GetTagByID", ctx, uint64(1)).Return(existingTag, nil)
				mockRepo.On("DeleteTag", ctx, uint64(1)).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "tag not found",
			tagID:  999,
			userID: 1,
			setupMocks: func() {
				mockRepo.On("GetTagByID", ctx, uint64(999)).Return((*model.Tag)(nil), gorm.ErrRecordNotFound)
			},
			expectedErr: errors.NewAppError(errors.CodeResourceNotFound, "tag not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			// Setup mocks
			tt.setupMocks()

			// Execute test
			err := service.DeleteTag(ctx, tt.tagID, tt.userID)

			// Assertions
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTagService_ListTags(t *testing.T) {
	service, mockRepo := setupTagServiceTest()
	ctx := context.Background()

	now := time.Now()
	expectedTags := []*model.Tag{
		{
			ID:          1,
			Name:        "tag1",
			Color:       "#ff0000",
			Description: "First tag",
			CreatedBy:   1,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          2,
			Name:        "tag2",
			Color:       "#00ff00",
			Description: "Second tag",
			CreatedBy:   1,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	tests := []struct {
		name        string
		req         *request.TagListRequest
		setupMocks  func()
		expectedErr error
		expectedLen int
	}{
		{
			name: "successful tag listing",
			req: &request.TagListRequest{
				Search:     "",
				ExcludeIDs: "",
			},
			setupMocks: func() {
				mockRepo.On("ListTags", ctx, "", []uint64{}).Return(expectedTags, nil)
			},
			expectedErr: nil,
			expectedLen: 2,
		},
		{
			name: "tag listing with search",
			req: &request.TagListRequest{
				Search:     "tag1",
				ExcludeIDs: "",
			},
			setupMocks: func() {
				mockRepo.On("ListTags", ctx, "tag1", []uint64{}).Return([]*model.Tag{expectedTags[0]}, nil)
			},
			expectedErr: nil,
			expectedLen: 1,
		},
		{
			name: "tag listing with exclusions",
			req: &request.TagListRequest{
				Search:     "",
				ExcludeIDs: "1,2",
			},
			setupMocks: func() {
				mockRepo.On("ListTags", ctx, "", []uint64{1, 2}).Return([]*model.Tag{}, nil)
			},
			expectedErr: nil,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			// Setup mocks
			tt.setupMocks()

			// Execute test
			result, err := service.ListTags(ctx, tt.req)

			// Assertions
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, tt.expectedLen)
			}

			// Verify all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTagService_GetResourceTags(t *testing.T) {
	service, mockRepo := setupTagServiceTest()
	ctx := context.Background()

	now := time.Now()
	expectedTaggings := []*model.Tagging{
		{
			ID:         1,
			TagID:      1,
			ResourceID: 123,
			CreatedBy:  1,
			CreatedAt:  now,
			Tag: &model.Tag{
				ID:          1,
				Name:        "tag1",
				Color:       "#ff0000",
				Description: "First tag",
			},
		},
		{
			ID:         2,
			TagID:      2,
			ResourceID: 123,
			CreatedBy:  1,
			CreatedAt:  now,
			Tag: &model.Tag{
				ID:          2,
				Name:        "tag2",
				Color:       "#00ff00",
				Description: "Second tag",
			},
		},
	}

	tests := []struct {
		name        string
		req         *request.TaggingListRequest
		setupMocks  func()
		expectedErr error
		expectedLen int
	}{
		{
			name: "successful resource tags retrieval",
			req: &request.TaggingListRequest{
				ResourceID: 123,
			},
			setupMocks: func() {
				mockRepo.On("GetTaggingsByResourceID", ctx, uint64(123)).Return(expectedTaggings, nil)
			},
			expectedErr: nil,
			expectedLen: 2,
		},
		{
			name: "no tags for resource",
			req: &request.TaggingListRequest{
				ResourceID: 456,
			},
			setupMocks: func() {
				mockRepo.On("GetTaggingsByResourceID", ctx, uint64(456)).Return([]*model.Tagging{}, nil)
			},
			expectedErr: nil,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			// Setup mocks
			tt.setupMocks()

			// Execute test
			result, err := service.GetResourceTags(ctx, tt.req)

			// Assertions
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, tt.expectedLen)
			}

			// Verify all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}
