package request

import "api-service/internal/dto"

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin@websoft9.com"`
	Password string `json:"password" binding:"required" example:"Websoft9"`
}

// UserChangePasswordRequest 用户修改密码请求
type UserChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" example:"oldpass123"`
	NewPassword string `json:"new_password" binding:"required" example:"newpass123"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	dto.BaseListRequest
	Status   *int    `form:"status" json:"status" binding:"omitempty,min=0,max=1" example:"1"`
	Keyword  *string `form:"keyword" json:"keyword" binding:"omitempty" example:"john"`
	Gender   *int    `form:"gender" json:"gender" binding:"omitempty,min=0,max=2" example:"1"`
	Language *string `form:"language" json:"language" binding:"omitempty" example:"zh-CN"`
	RoleID   *uint   `form:"role_id" json:"role_id" binding:"omitempty,gt=0" example:"2"` // Role ID filter (optional)
}

// UserCreateRequest defines the request structure for creating a user.
// Includes user basic information and supports assigning multiple roles via role_ids.
type UserCreateRequest struct {
	Username  string  `json:"username" binding:"required,max=64" example:"johndoe"`                              // Username, 3-32 characters, letters, numbers, underscore
	Email     string  `json:"email" binding:"required,email" example:"john@example.com"`                         // Email address, must be unique
	Password  string  `json:"password" binding:"required" example:"password123"`                                 // Password, 8-32 characters, must contain uppercase
	Nickname  *string `json:"nickname,omitempty" binding:"omitempty,max=64" example:"John Doe"`                  // User nickname (optional)
	Phone     *string `json:"phone,omitempty" binding:"omitempty,max=20" example:"+1234567890"`                  // Phone number (optional)
	Avatar    *string `json:"avatar,omitempty" binding:"omitempty,url" example:"https://example.com/avatar.jpg"` // Avatar URL (optional)
	Gender    *int    `json:"gender,omitempty" binding:"omitempty,min=0,max=2" example:"1"`                      // Gender: 1-male, 2-female, 0-unknown (optional)
	Signature *string `json:"signature,omitempty" binding:"omitempty,max=255" example:"This is my signature"`    // Personal signature (optional)
	Status    *int    `json:"status,omitempty" binding:"omitempty,min=0,max=1" example:"1"`                      // Status: 0-disabled, 1-enabled (optional)
	Timezone  *string `json:"timezone,omitempty" binding:"omitempty,max=64" example:"Asia/Shanghai"`             // Timezone, default UTC (optional)
	Language  *string `json:"language,omitempty" binding:"omitempty,max=10" example:"zh-CN"`                     // Language, default zh-CN (optional)
	RoleIDs   []uint  `json:"role_ids,omitempty" binding:"omitempty,dive,gt=0" example:"[1,2]"`                  // Array of role IDs to assign (optional)
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	Username  *string `json:"username,omitempty" binding:"omitempty,max=64" example:"johndoe"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email" example:"newemail@example.com"`
	Nickname  *string `json:"nickname,omitempty" binding:"omitempty,max=64" example:"John Doe"`
	Phone     *string `json:"phone,omitempty" binding:"omitempty,max=20" example:"+1234567890"`
	Avatar    *string `json:"avatar,omitempty" binding:"omitempty,url" example:"https://example.com/avatar.jpg"`
	Gender    *int    `json:"gender,omitempty" binding:"omitempty,min=0,max=2" example:"1"`
	Signature *string `json:"signature,omitempty" binding:"omitempty,max=255" example:"This is my signature"`
	Timezone  *string `json:"timezone,omitempty" binding:"omitempty,max=64" example:"Asia/Shanghai"`
	Language  *string `json:"language,omitempty" binding:"omitempty,max=10" example:"zh-CN"`
	RoleIDs   []uint  `json:"role_ids,omitempty" binding:"omitempty,dive,gt=0" example:"[1,2]"` // Array of role IDs to assign (optional)
}

// UserUpdateStatusRequest 用户状态更新请求
type UserUpdateStatusRequest struct {
	Status int `json:"status" binding:"required,min=0,max=1" example:"1"`
}

// UserPasswordUpdateRequest 管理员修改用户密码请求
type UserPasswordUpdateRequest struct {
	NewPassword string `json:"new_password" binding:"required" example:"newpassword123"`
}

// ForgotPasswordRequest 忘记密码请求
type ForgotPasswordRequest struct {
	Username string `json:"username" binding:"required,email" example:"john@example.com"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required" example:"abc123def456"`
	NewPassword string `json:"new_password" binding:"required" example:"newpassword123"`
}

// VerifyEmailRequest 验证邮箱请求
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required" example:"abc123def456"`
}

// ResendVerificationRequest 重新发送验证邮件请求
type ResendVerificationRequest struct {
	Username string `json:"username" binding:"required,email" example:"john@example.com"`
}

// OAuth2LoginRequest OAuth2登录请求
type OAuth2LoginRequest struct {
	Provider string `json:"provider" binding:"required" example:"github"`
	Code     string `json:"code" binding:"required" example:"authorization_code"`
	State    string `json:"state" binding:"required" example:"random_state"`
}
