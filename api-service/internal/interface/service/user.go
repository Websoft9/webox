package service

import (
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"context"
)

// UserService 用户业务逻辑接口
type UserService interface {
	// 用户资料管理
	ChangePassword(ctx context.Context, userID uint, req *request.UserChangePasswordRequest) error

	// 用户管理（管理员功能）
	ListUsers(ctx context.Context, req *request.UserListRequest) (*response.UserListResponse, int64, error)
	CreateUser(ctx context.Context, req *request.UserCreateRequest) (*response.UserResponse, error)
	GetUser(ctx context.Context, userID uint) (*response.UserResponse, error)
	UpdateUser(ctx context.Context, userID uint, req *request.UserUpdateRequest) (*response.UserResponse, error)
	UpdateUserStatus(ctx context.Context, userID uint, req *request.UserUpdateStatusRequest) error
	UpdateUserPassword(ctx context.Context, userID uint, req *request.UserPasswordUpdateRequest) error
	DeleteUser(ctx context.Context, userID uint) error

	// 业务验证
	ValidateUserAccess(ctx context.Context, userID uint, resource string) error
	CheckUserQuota(ctx context.Context, userID uint, resourceType string) error
}
