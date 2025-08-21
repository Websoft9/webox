package request

import "api-service/internal/dto"

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20" example:"john_doe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"123456"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" binding:"required" example:"john_doe"`
	Password string `json:"password" binding:"required" example:"123456"`
}

// UserUpdateProfileRequest 用户更新资料请求
type UserUpdateProfileRequest struct {
	Email     *string `json:"email,omitempty" binding:"omitempty,email" example:"newemail@example.com"`
	FirstName *string `json:"first_name,omitempty" binding:"omitempty,max=50" example:"John"`
	LastName  *string `json:"last_name,omitempty" binding:"omitempty,max=50" example:"Doe"`
	Avatar    *string `json:"avatar,omitempty" binding:"omitempty,url" example:"https://example.com/avatar.jpg"`
}

// UserChangePasswordRequest 用户修改密码请求
type UserChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" example:"oldpass123"`
	NewPassword string `json:"new_password" binding:"required,min=6" example:"newpass123"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	dto.BaseListRequest
	Status string `form:"status" json:"status" binding:"omitempty,oneof=active inactive banned" example:"active"`
	Role   string `form:"role" json:"role" binding:"omitempty,oneof=admin user guest" example:"user"`
}

// UserUpdateStatusRequest 用户状态更新请求
type UserUpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive banned" example:"active"`
}
