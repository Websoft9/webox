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

// UserCreateRequest 创建用户请求（管理员功能）
type UserCreateRequest struct {
	GroupID   uint   `json:"group_id" binding:"required,min=1" example:"1"`
	Username  string `json:"username" binding:"required,min=3,max=32" example:"john_doe"`
	Email     string `json:"email" binding:"required,email" example:"john@example.com"`
	Password  string `json:"password" binding:"required,min=8,max=32" example:"Password123"`
	Nickname  string `json:"nickname" binding:"omitempty,max=100" example:"John"`
	Phone     string `json:"phone" binding:"omitempty" example:"13800138000"`
	Gender    int    `json:"gender" binding:"omitempty,oneof=0 1 2" example:"1"`
	Signature string `json:"signature" binding:"omitempty,max=500" example:"个性签名"`
	Timezone  string `json:"timezone" binding:"omitempty,max=50" example:"Asia/Shanghai"`
	Language  string `json:"language" binding:"omitempty,max=10" example:"zh-CN"`
	Status    int    `json:"status" binding:"omitempty,oneof=0 1" example:"1"`
}

// UserUpdateRequest 更新用户请求（管理员功能）
type UserUpdateRequest struct {
	GroupID   *uint   `json:"group_id,omitempty" binding:"omitempty,min=1" example:"1"`
	Username  *string `json:"username,omitempty" binding:"omitempty,min=3,max=32" example:"john_doe"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email" example:"john@example.com"`
	Nickname  *string `json:"nickname,omitempty" binding:"omitempty,max=100" example:"John"`
	Phone     *string `json:"phone,omitempty" binding:"omitempty" example:"13800138000"`
	Gender    *int    `json:"gender,omitempty" binding:"omitempty,oneof=0 1 2" example:"1"`
	Signature *string `json:"signature,omitempty" binding:"omitempty,max=500" example:"个性签名"`
	Timezone  *string `json:"timezone,omitempty" binding:"omitempty,max=50" example:"Asia/Shanghai"`
	Language  *string `json:"language,omitempty" binding:"omitempty,max=10" example:"zh-CN"`
	Status    *int    `json:"status,omitempty" binding:"omitempty,oneof=0 1" example:"1"`
}

// UserUpdateProfileRequest 用户更新资料请求（用户自己操作）
type UserUpdateProfileRequest struct {
	Email     *string `json:"email,omitempty" binding:"omitempty,email" example:"newemail@example.com"`
	Nickname  *string `json:"nickname,omitempty" binding:"omitempty,max=100" example:"John"`
	Avatar    *string `json:"avatar,omitempty" binding:"omitempty,url" example:"https://example.com/avatar.jpg"`
	Phone     *string `json:"phone,omitempty" binding:"omitempty" example:"13800138000"`
	Gender    *int    `json:"gender,omitempty" binding:"omitempty,oneof=0 1 2" example:"1"`
	Signature *string `json:"signature,omitempty" binding:"omitempty,max=500" example:"个性签名"`
	Timezone  *string `json:"timezone,omitempty" binding:"omitempty,max=50" example:"Asia/Shanghai"`
	Language  *string `json:"language,omitempty" binding:"omitempty,max=10" example:"zh-CN"`
}

// UserChangePasswordRequest 用户修改密码请求
type UserChangePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required" example:"oldpass123"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=32" example:"newpass123"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,max=32" example:"newpass123"`
}

// AdminChangePasswordRequest 管理员修改用户密码请求
type AdminChangePasswordRequest struct {
	NewPassword     string `json:"new_password" binding:"required,min=8,max=32" example:"newpass123"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,max=32" example:"newpass123"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	dto.BaseListRequest
	Status *int `form:"status" json:"status" binding:"omitempty,oneof=0 1" example:"1"`
}

// UserUpdateStatusRequest 用户状态更新请求
type UserUpdateStatusRequest struct {
	Status int `json:"status" binding:"required,oneof=0 1" example:"1"`
}
