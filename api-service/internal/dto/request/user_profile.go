package request

// UserProfileUpdateRequest 更新用户个人资料的请求
type UserProfileUpdateRequest struct {
	Nickname  *string `json:"nickname" example:"John Doe"`
	Avatar    *string `json:"avatar" example:"https://example.com/avatar.jpg"`
	Phone     *string `json:"phone" example:"+1234567890"`
	Gender    *int    `json:"gender" example:"1" binding:"omitempty,oneof=0 1 2"` // 0-未知，1-男，2-女
	Signature *string `json:"signature" example:"This is my signature"`
	Timezone  *string `json:"timezone" example:"Asia/Shanghai"`
	Language  *string `json:"language" example:"zh-CN" binding:"omitempty,len=5"`
}

// ProfileChangePasswordRequest 用户个人中心修改密码的请求
type ProfileChangePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required" example:"oldpass123"`
	NewPassword     string `json:"new_password" binding:"required,min=6" example:"newpass123"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword" example:"newpass123"`
}

type LoginHistoryRequest struct {
	Page     int `form:"page" binding:"min=0"`      // 页码
	PageSize int `form:"page_size" binding:"min=0"` // 每页记录数
}

// NotificationSettingsRequest 通知设置请求
type NotificationSettingsRequest struct {
	EmailNotifications bool `json:"email_notifications" binding:"omitempty"`
	SmsNotifications   bool `json:"sms_notifications" binding:"omitempty"`
	PushNotifications  bool `json:"push_notifications" binding:"omitempty"`
	MarketingEmails    bool `json:"marketing_emails" binding:"omitempty"`
}

// SecuritySettingsRequest 安全设置请求
type SecuritySettingsRequest struct {
	LoginAlerts    bool `json:"login_alerts" binding:"omitempty"`
	SessionTimeout int  `json:"session_timeout" binding:"omitempty,min=0"`
}
