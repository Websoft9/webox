package response

import "time"

// UserResponse 用户响应结构
type UserResponse struct {
	ID          uint       `json:"id" example:"1"`
	GroupID     uint       `json:"group_id" example:"1"`
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
	Group *UserGroupResponse `json:"group,omitempty"`
	Roles []RoleResponse     `json:"roles,omitempty"`
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
	Total int64          `json:"total" example:"100"`
}

// UserGroupResponse 用户组响应结构
type UserGroupResponse struct {
	ID          uint      `json:"id" example:"1"`
	Name        string    `json:"name" example:"管理员组"`
	Code        string    `json:"code" example:"admin"`
	Description string    `json:"description" example:"系统管理员用户组"`
	SortOrder   int       `json:"sort_order" example:"0"`
	Status      int       `json:"status" example:"1"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}
