package response

import (
	"time"
)

// UserProfileResponse 定义用户个人资料的响应结构
type UserProfileResponse struct {
	ID          uint       `json:"id" example:"1"`
	Username    string     `json:"username" example:"john_doe"`
	Email       string     `json:"email" example:"john@example.com"`
	Nickname    string     `json:"nickname" example:"John Doe"`
	Avatar      string     `json:"avatar" example:"https://example.com/avatar.jpg"`
	Phone       string     `json:"phone" example:"+1234567890"`
	Gender      int        `json:"gender" example:"1"`
	Signature   string     `json:"signature" example:"This is my signature"`
	Status      int        `json:"status" example:"1"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty" example:"2023-01-01T12:00:00Z"`
	LastLoginIP string     `json:"last_login_ip" example:"192.168.1.1"`
	Timezone    string     `json:"timezone" example:"Asia/Shanghai"`
	Language    string     `json:"language" example:"zh-CN"`
	CreatedAt   time.Time  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	// 关联数据
	Roles []RoleResponse `json:"roles,omitempty"`
}
