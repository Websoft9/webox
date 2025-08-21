package response

import "time"

// UserResponse 用户响应结构
type UserResponse struct {
	ID          uint       `json:"id" example:"1"`
	GroupID     uint       `json:"group_id" example:"1"`
	Username    string     `json:"username" example:"john_doe"`
	Email       string     `json:"email" example:"john@example.com"`
	Nickname    string     `json:"nickname" example:"John"`
	Avatar      string     `json:"avatar" example:"https://example.com/avatar.jpg"`
	Phone       string     `json:"phone" example:"13800138000"`
	Gender      int        `json:"gender" example:"1"`
	Signature   string     `json:"signature" example:"个性签名"`
	Status      int        `json:"status" example:"1"`
	Timezone    string     `json:"timezone" example:"Asia/Shanghai"`
	Language    string     `json:"language" example:"zh-CN"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty" example:"2025-07-15T10:30:00Z"`
	LastLoginIP string     `json:"last_login_ip" example:"192.168.1.100"`
	CreatedAt   time.Time  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// UserLoginResponse 用户登录响应
type UserLoginResponse struct {
	Token     string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt time.Time    `json:"expires_at" example:"2023-01-02T00:00:00Z"`
	User      UserResponse `json:"user"`
}

// UserProfileResponse 用户资料响应（包含更多详细信息）
type UserProfileResponse struct {
	UserResponse
	LoginCount       int `json:"login_count" example:"10"`
	ApplicationCount int `json:"application_count" example:"5"`
	WorkflowCount    int `json:"workflow_count" example:"3"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users []UserResponse `json:"users"`
}
