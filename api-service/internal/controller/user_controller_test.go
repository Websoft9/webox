package controller

import (
	"api-service/internal/model"
	"api-service/internal/service"
	"api-service/pkg/response"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService 模拟用户服务
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(username, email, password string) (*model.User, error) {
	args := m.Called(username, email, password)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) GetProfile(userID uint) (*model.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) UpdateProfile(userID uint, updates map[string]interface{}) error {
	args := m.Called(userID, updates)
	return args.Error(0)
}

func (m *MockUserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	args := m.Called(userID, oldPassword, newPassword)
	return args.Error(0)
}

func (m *MockUserService) ListUsers(page, pageSize int, keyword, status, groupID string) ([]*model.User, int64, error) {
	args := m.Called(page, pageSize, keyword, status, groupID)
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) CreateUser(req *service.CreateUserRequest) (*model.User, error) {
	args := m.Called(req)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(userID uint, req *service.UpdateUserRequest) (*model.User, error) {
	args := m.Called(userID, req)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestUserController_ListUsers(t *testing.T) {
	// 准备
	mockService := new(MockUserService)
	controller := NewUserController(mockService)
	router := setupRouter()

	// 模拟数据
	users := []*model.User{
		{
			ID:       1,
			Username: "admin",
			Email:    "admin@example.com",
			Nickname: "管理员",
			Status:   1,
		},
		{
			ID:       2,
			Username: "user1",
			Email:    "user1@example.com",
			Nickname: "普通用户",
			Status:   1,
		},
	}

	// 设置预期
	mockService.On("ListUsers", 1, 20, "", "", "").Return(users, int64(2), nil)

	// 设置路由
	router.GET("/api/v1/users", controller.ListUsers)

	// 执行请求
	req, _ := http.NewRequest("GET", "/api/v1/users?page=1&page_size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resp.Message)

	mockService.AssertExpectations(t)
}

func TestUserController_CreateUser(t *testing.T) {
	// 准备
	mockService := new(MockUserService)
	controller := NewUserController(mockService)
	router := setupRouter()

	// 模拟数据
	newUser := &model.User{
		ID:       3,
		Username: "newuser",
		Email:    "newuser@example.com",
		Nickname: "新用户",
		Status:   1,
		GroupID:  1,
	}

	// 设置预期
	mockService.On("CreateUser", mock.AnythingOfType("*service.CreateUserRequest")).Return(newUser, nil)

	// 设置路由
	router.POST("/api/v1/users", controller.CreateUser)

	// 准备请求数据
	reqData := CreateUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
		Nickname: "新用户",
		GroupID:  1,
	}

	jsonData, _ := json.Marshal(reqData)

	// 执行请求
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resp.Message)

	mockService.AssertExpectations(t)
}

func TestUserController_GetUserByID(t *testing.T) {
	// 准备
	mockService := new(MockUserService)
	controller := NewUserController(mockService)
	router := setupRouter()

	// 模拟数据
	user := &model.User{
		ID:       1,
		Username: "admin",
		Email:    "admin@example.com",
		Nickname: "管理员",
		Status:   1,
	}

	// 设置预期
	mockService.On("GetProfile", uint(1)).Return(user, nil)

	// 设置路由
	router.GET("/api/v1/users/:id", controller.GetUserByID)

	// 执行请求
	req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resp.Message)

	mockService.AssertExpectations(t)
}

func TestUserController_UpdateUser(t *testing.T) {
	// 准备
	mockService := new(MockUserService)
	controller := NewUserController(mockService)
	router := setupRouter()

	// 模拟数据
	updatedUser := &model.User{
		ID:       1,
		Username: "admin",
		Email:    "admin@example.com",
		Nickname: "超级管理员",
		Status:   1,
	}

	// 设置预期
	mockService.On("UpdateUser", uint(1), mock.AnythingOfType("*service.UpdateUserRequest")).Return(updatedUser, nil)

	// 设置路由
	router.PUT("/api/v1/users/:id", controller.UpdateUser)

	// 准备请求数据
	reqData := UpdateUserRequest{
		Nickname: "超级管理员",
	}

	jsonData, _ := json.Marshal(reqData)

	// 执行请求
	req, _ := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resp.Message)

	mockService.AssertExpectations(t)
}

func TestUserController_DeleteUser(t *testing.T) {
	// 准备
	mockService := new(MockUserService)
	controller := NewUserController(mockService)
	router := setupRouter()

	// 设置预期
	mockService.On("DeleteUser", uint(1)).Return(nil)

	// 设置路由
	router.DELETE("/api/v1/users/:id", controller.DeleteUser)

	// 执行请求
	req, _ := http.NewRequest("DELETE", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resp.Message)

	mockService.AssertExpectations(t)
}

func TestUserController_ChangePassword(t *testing.T) {
	// 准备
	mockService := new(MockUserService)
	controller := NewUserController(mockService)
	router := setupRouter()

	// 设置预期
	mockService.On("ChangePassword", uint(1), "oldpass123", "newpass123").Return(nil)

	// 设置路由
	router.PUT("/api/v1/users/:id/password", controller.ChangePassword)

	// 准备请求数据
	reqData := ChangePasswordRequest{
		OldPassword: "oldpass123",
		NewPassword: "newpass123",
	}

	jsonData, _ := json.Marshal(reqData)

	// 执行请求
	req, _ := http.NewRequest("PUT", "/api/v1/users/1/password", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resp.Message)

	mockService.AssertExpectations(t)
}
