package service

import (
	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"context"
)

type UserService interface {
	ChangePassword(ctx context.Context, userID uint, req *request.UserChangePasswordRequest) error
	ListUsers(ctx context.Context, req *request.UserListRequest) (*common.PaginationResponse, error)
	CreateUser(ctx context.Context, currentUserID uint, req *request.UserCreateRequest) (*response.UserResponse, error)
	GetUser(ctx context.Context, userID uint) (*response.UserResponse, error)
	UpdateUser(ctx context.Context, currentUserID, userID uint, req *request.UserUpdateRequest) (*response.UserResponse, error)
	UpdateUserStatus(ctx context.Context, userID uint, req *request.UserUpdateStatusRequest) error
	UpdateUserPassword(ctx context.Context, userID uint, req *request.UserPasswordUpdateRequest) error
	DeleteUser(ctx context.Context, userID uint) error
	ValidateUserAccess(ctx context.Context, userID uint, resource string) error
	CheckUserQuota(ctx context.Context, userID uint, resourceType string) error
}
