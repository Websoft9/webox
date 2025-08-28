package request

import "api-service/internal/dto"

// ProfileUpdateRequest 更新个人资料请求
type ProfileUpdateRequest struct {
	Nickname  *string `json:"nickname,omitempty" binding:"omitempty,max=64" example:"John Doe"`
	Avatar    *string `json:"avatar,omitempty" binding:"omitempty,url" example:"https://example.com/avatar.jpg"`
	Phone     *string `json:"phone,omitempty" binding:"omitempty,max=20" example:"+1234567890"`
	Gender    *int    `json:"gender,omitempty" binding:"omitempty,min=0,max=2" example:"1"`
	Signature *string `json:"signature,omitempty" binding:"omitempty,max=255" example:"This is my signature"`
	Timezone  *string `json:"timezone,omitempty" binding:"omitempty,max=64" example:"Asia/Shanghai"`
	Language  *string `json:"language,omitempty" binding:"omitempty,max=10" example:"zh-CN"`
}

// PasswordChangeRequest 修改密码请求
type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" example:"oldpass123"`
	NewPassword     string `json:"new_password" binding:"required,min=6" example:"newpass123"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6" example:"newpass123"`
}

// NotificationSettingsRequest 通知设置请求
type NotificationSettingsRequest struct {
	EmailNotifications *bool `json:"email_notifications,omitempty" example:"true"`
	SMSNotifications   *bool `json:"sms_notifications,omitempty" example:"false"`
	PushNotifications  *bool `json:"push_notifications,omitempty" example:"true"`
	LoginAlerts        *bool `json:"login_alerts,omitempty" example:"true"`
	SecurityAlerts     *bool `json:"security_alerts,omitempty" example:"true"`
}

// SecuritySettingsRequest 安全设置请求
type SecuritySettingsRequest struct {
	SessionTimeout        *int  `json:"session_timeout,omitempty" binding:"omitempty,min=5,max=1440" example:"30"`
	LoginNotifications    *bool `json:"login_notifications,omitempty" example:"true"`
	PasswordStrengthCheck *bool `json:"password_strength_check,omitempty" example:"true"`
}

// TwoFactorEnableRequest 启用两因子认证请求
type TwoFactorEnableRequest struct {
	Type string `json:"type" binding:"required,oneof=totp" example:"totp"`
}

// TwoFactorVerifyRequest 验证两因子认证请求
type TwoFactorVerifyRequest struct {
	Type string `json:"type" binding:"required,oneof=totp" example:"totp"`
	Code string `json:"code" binding:"required,len=6" example:"123456"`
}

// TwoFactorDisableRequest 禁用两因子认证请求
type TwoFactorDisableRequest struct {
	Type     string `json:"type" binding:"required,oneof=totp" example:"totp"`
	Password string `json:"password,omitempty" example:"password123"`
	Code     string `json:"code,omitempty" binding:"omitempty,len=6" example:"123456"`
}

// LoginHistoryRequest 登录历史请求
type LoginHistoryRequest struct {
	dto.PaginationRequest
	StartDate *string `json:"start_date,omitempty" binding:"omitempty,datetime=2006-01-02" example:"2023-01-01"`
	EndDate   *string `json:"end_date,omitempty" binding:"omitempty,datetime=2006-01-02" example:"2023-12-31"`
	Status    *string `json:"status,omitempty" example:"success"`
}

// LoginHistoryListRequest 登录历史列表请求 (别名，保持兼容性)
type LoginHistoryListRequest = LoginHistoryRequest
