package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"api-service/internal/constants"
	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/model"
	"api-service/pkg/logger"
)

// MockNotificationRecordRepository is a mock for NotificationRecordRepository
type MockNotificationRecordRepository struct {
	mock.Mock
}

func (m *MockNotificationRecordRepository) GetByID(ctx context.Context, id uint) (*model.NotificationRecord, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.NotificationRecord), args.Error(1)
}

func (m *MockNotificationRecordRepository) GetList(ctx context.Context, req *request.GetNotificationRecordListRequest) ([]*model.NotificationRecord, int64, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.NotificationRecord), args.Get(1).(int64), args.Error(2)
}

// Test setup helper
func setupNotificationRecordServiceTest() (*notificationRecordService, *MockNotificationRecordRepository) {
	mockRepo := new(MockNotificationRecordRepository)
	mockLogger := logger.NewZapLogger(logger.InfoLevel, nil)
	service := NewNotificationRecordService(mockRepo, mockLogger).(*notificationRecordService)
	return service, mockRepo
}

func TestNotificationRecordService_GetNotificationRecordList(t *testing.T) {
	service, mockRepo := setupNotificationRecordServiceTest()
	ctx := context.Background()

	tests := []struct {
		name      string
		req       *request.GetNotificationRecordListRequest
		mockData  []*model.NotificationRecord
		mockTotal int64
		mockError error
		expectErr bool
	}{
		{
			name: "successful retrieval",
			req: &request.GetNotificationRecordListRequest{
				BaseListRequest: common.BaseListRequest{
					PaginationRequest: common.PaginationRequest{
						Page:     1,
						PageSize: 20,
					},
				},
			},
			mockData: []*model.NotificationRecord{
				{
					ID:          1,
					ChannelType: constants.NotificationChannelEmail,
					Recipient:   "test@example.com",
					Subject:     stringPtrNotification("Test Subject"),
					Content:     "Test Content",
					Status:      constants.NotificationStatusSent,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					ID:          2,
					ChannelType: constants.NotificationChannelInternal,
					Recipient:   "user2",
					Content:     "Internal notification",
					Status:      constants.NotificationStatusPending,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			mockTotal: 2,
			mockError: nil,
			expectErr: false,
		},
		{
			name: "repository error",
			req: &request.GetNotificationRecordListRequest{
				BaseListRequest: common.BaseListRequest{
					PaginationRequest: common.PaginationRequest{
						Page:     1,
						PageSize: 20,
					},
				},
			},
			mockData:  nil,
			mockTotal: 0,
			mockError: errors.New("database connection failed"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetList", ctx, tt.req).Return(tt.mockData, tt.mockTotal, tt.mockError).Once()

			result, err := service.GetNotificationRecordList(ctx, tt.req)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockTotal, result.Total)
				if items, ok := result.Items.([]response.NotificationRecordResponse); ok {
					assert.Equal(t, len(tt.mockData), len(items))
				}
				assert.Equal(t, tt.req.Page, result.Page)
				assert.Equal(t, tt.req.PageSize, result.PageSize)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestNotificationRecordService_GetNotificationRecordByID(t *testing.T) {
	service, mockRepo := setupNotificationRecordServiceTest()
	ctx := context.Background()

	tests := []struct {
		name      string
		id        uint
		mockData  *model.NotificationRecord
		mockError error
		expectErr bool
	}{
		{
			name: "successful retrieval",
			id:   1,
			mockData: &model.NotificationRecord{
				ID:          1,
				ChannelType: constants.NotificationChannelEmail,
				Recipient:   "test@example.com",
				Subject:     stringPtrNotification("Test Subject"),
				Content:     "Test Content",
				Status:      constants.NotificationStatusSent,
				SentAt:      timePtrNotification(time.Now()),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			mockError: nil,
			expectErr: false,
		},
		{
			name:      "record not found",
			id:        999,
			mockData:  nil,
			mockError: gorm.ErrRecordNotFound,
			expectErr: true,
		},
		{
			name:      "repository error",
			id:        1,
			mockData:  nil,
			mockError: errors.New("database connection failed"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetByID", ctx, tt.id).Return(tt.mockData, tt.mockError).Once()

			result, err := service.GetNotificationRecordByID(ctx, tt.id)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockData.ID, result.ID)
				assert.Equal(t, tt.mockData.Content, result.Content)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// Helper functions for tests
func stringPtrNotification(s string) *string {
	return &s
}

func timePtrNotification(t time.Time) *time.Time {
	return &t
}
