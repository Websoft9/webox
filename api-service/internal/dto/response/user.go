package response

import "time"

// UserResponse 用户响应结构
type UserResponse struct {
	ID        uint      `json:"id" example:"1"`
	Username  string    `json:"username" example:"john_doe"`
	Email     string    `json:"email" example:"john@example.com"`
	FirstName string    `json:"first_name" example:"John"`
	LastName  string    `json:"last_name" example:"Doe"`
	Avatar    string    `json:"avatar" example:"https://example.com/avatar.jpg"`
	Status    string    `json:"status" example:"active"`
	Role      string    `json:"role" example:"user"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
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
	LastLoginAt    *time.Time `json:"last_login_at,omitempty" example:"2023-01-01T12:00:00Z"`
	LoginCount     int        `json:"login_count" example:"10"`
	ApplicationCount int      `json:"application_count" example:"5"`
	WorkflowCount    int      `json:"workflow_count" example:"3"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users []UserResponse `json:"users"`
}
