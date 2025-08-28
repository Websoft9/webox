package response

import "time"

// ProfileResponse 个人资料响应
type ProfileResponse struct {
	ID                   uint                  `json:"id" example:"1"`
	Username             string                `json:"username" example:"admin"`
	Email                string                `json:"email" example:"admin@example.com"`
	Nickname             string                `json:"nickname" example:"系统管理员"`
	Avatar               string                `json:"avatar" example:"https://example.com/avatars/admin.jpg"`
	Phone                string                `json:"phone" example:"13800138000"`
	Gender               int                   `json:"gender" example:"1"`
	Signature            string                `json:"signature" example:"系统管理员账户"`
	Timezone             string                `json:"timezone" example:"Asia/Shanghai"`
	Language             string                `json:"language" example:"zh-CN"`
	TwoFactorEnabled     bool                  `json:"two_factor_enabled" example:"true"`
	NotificationSettings *NotificationSettings `json:"notification_settings,omitempty"`
	SecuritySettings     *SecuritySettings     `json:"security_settings,omitempty"`
	LastLoginAt          *time.Time            `json:"last_login_at,omitempty" example:"2025-07-15T10:30:00Z"`
	LastLoginIP          string                `json:"last_login_ip" example:"192.168.1.100"`
	CreatedAt            time.Time             `json:"created_at" example:"2025-07-15T10:30:00Z"`
	UpdatedAt            time.Time             `json:"updated_at" example:"2025-07-15T10:30:00Z"`
}

// NotificationSettings 通知设置
type NotificationSettings struct {
	EmailNotifications bool `json:"email_notifications" example:"true"`
	SMSNotifications   bool `json:"sms_notifications" example:"false"`
	PushNotifications  bool `json:"push_notifications" example:"true"`
	MarketingEmails    bool `json:"marketing_emails" example:"false"`
}

// SecuritySettings 安全设置
type SecuritySettings struct {
	LoginAlerts         bool       `json:"login_alerts" example:"true"`
	SessionTimeout      int        `json:"session_timeout" example:"3600"`
	PasswordLastChanged *time.Time `json:"password_last_changed,omitempty" example:"2025-07-15T10:30:00Z"`
}

// TwoFactorEnableResponse 启用双因子认证响应
type TwoFactorEnableResponse struct {
	QRCode      string   `json:"qr_code" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
	Secret      string   `json:"secret" example:"JBSWY3DPEHPK3PXP"`
	BackupCodes []string `json:"backup_codes" example:"[\"12345678\",\"87654321\"]"`
}

// LoginHistoryResponse 登录历史响应
type LoginHistoryResponse struct {
	ID         uint       `json:"id" example:"1"`
	IPAddress  string     `json:"ip_address" example:"192.168.1.100"`
	UserAgent  string     `json:"user_agent" example:"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"`
	Location   string     `json:"location" example:"北京市"`
	Device     string     `json:"device" example:"Windows PC"`
	Browser    string     `json:"browser" example:"Chrome 120.0"`
	OS         string     `json:"os" example:"Windows 10"`
	LoginTime  time.Time  `json:"login_time" example:"2025-07-15T10:30:00Z"`
	LogoutTime *time.Time `json:"logout_time,omitempty" example:"2025-07-15T18:30:00Z"`
	Status     string     `json:"status" example:"ACTIVE"`
}

// LoginHistoryListResponse 登录历史列表响应
type LoginHistoryListResponse struct {
	Items []LoginHistoryResponse `json:"items"`
	Total int64                  `json:"total" example:"100"`
}

// PreferencesResponse 个人偏好设置响应
type PreferencesResponse struct {
	Theme    string `json:"theme" example:"light"`
	Language string `json:"language" example:"zh-CN"`
	Timezone string `json:"timezone" example:"Asia/Shanghai"`
}

// AvatarUploadResponse 头像上传响应
type AvatarUploadResponse struct {
	AvatarURL string `json:"avatar_url" example:"https://example.com/avatars/admin_new.jpg"`
}

// NotificationSettingsResponse 通知设置响应
type NotificationSettingsResponse struct {
	EmailEnabled bool `json:"email_enabled" example:"true"`
	SMSEnabled   bool `json:"sms_enabled" example:"false"`
	PushEnabled  bool `json:"push_enabled" example:"true"`
}

// SecuritySettingsResponse 安全设置响应
type SecuritySettingsResponse struct {
	TwoFactorEnabled  bool `json:"two_factor_enabled" example:"true"`
	LoginNotification bool `json:"login_notification" example:"true"`
	SessionTimeout    int  `json:"session_timeout" example:"30"`
}
